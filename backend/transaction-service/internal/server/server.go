package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"

	transactionhandler "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/handler/transaction"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/route"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/config"
	bookinfra "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/infra/book"
	dbinfra "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/infra/db"
	loanrepo "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/infra/db/repository/loan"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/infra/messaging"
	transactionusecase "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/usecase/transaction"
)

type Server struct {
	config    config.Config
	db        *sqlx.DB
	app       *fiber.App
	publisher *messaging.RabbitPublisher
	worker    *messaging.OutboxWorker
	consumer  *messaging.RabbitAckConsumer
}

func New(ctx context.Context) (*Server, error) {
	cfg, err := config.LoadServer()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	database, err := dbinfra.Connect(ctx, dbinfra.Config{DSN: cfg.Database.DSN, MaxOpenConns: cfg.Database.MaxOpenConns, MaxIdleConns: cfg.Database.MaxIdleConns, ConnMaxLifetime: cfg.Database.ConnMaxLifetime, ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime})
	if err != nil {
		return nil, err
	}
	closeOnError := func(err error) (*Server, error) { _ = database.Close(); return nil, err }
	tokens, err := bookinfra.NewOAuthTokenSource(bookinfra.OAuthConfig{TokenURL: cfg.OAuth.TokenURL, ClientID: cfg.OAuth.ClientID, ClientSecret: cfg.OAuth.ClientSecret, Timeout: cfg.Book.Timeout})
	if err != nil {
		return closeOnError(err)
	}
	stock, err := bookinfra.NewClient(bookinfra.Config{BaseURL: cfg.Book.URL, Timeout: cfg.Book.Timeout}, tokens)
	if err != nil {
		return closeOnError(err)
	}
	repository := loanrepo.NewRepository(database)
	rabbitConfig := messaging.RabbitConfig{URL: cfg.Rabbit.URL, Exchange: cfg.Rabbit.Exchange, DeadExchange: cfg.Rabbit.DeadExchange, BookReturnQueue: cfg.Rabbit.BookReturnQueue, StockAckQueue: cfg.Rabbit.StockAckQueue, DeadLetterQueue: cfg.Rabbit.DeadLetterQueue, ConfirmTimeout: cfg.Rabbit.ConfirmTimeout}
	publisher, err := messaging.NewRabbitPublisher(rabbitConfig)
	if err != nil {
		return closeOnError(err)
	}
	consumer, err := messaging.NewRabbitAckConsumer(rabbitConfig, messaging.NewAckProcessor(repository))
	if err != nil {
		return closeOnError(err)
	}
	usecase := transactionusecase.NewUsecase(repository, stock, transactionusecase.Config{LoanTerm: cfg.Loan.Term, AckTimeout: cfg.Loan.AckTimeout, PollInterval: cfg.Loan.PollInterval, DailyFineMinor: cfg.Loan.DailyFineMinor})
	app := buildApp(cfg, database, usecase)
	return &Server{config: cfg, db: database, app: app, publisher: publisher, worker: messaging.NewOutboxWorker(repository, publisher), consumer: consumer}, nil
}

func buildApp(cfg config.Config, database *sqlx.DB, usecase transactionusecase.Usecase) *fiber.App {
	app := fiber.New(fiber.Config{AppName: cfg.Service.Name, BodyLimit: 1 << 20, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, ErrorHandler: helper.ErrorHandler})
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{ReferrerPolicy: "no-referrer"}))
	app.Use(cors.New(cors.Config{AllowOrigins: cfg.Service.AllowedOrigin, AllowMethods: "GET,POST,OPTIONS", AllowHeaders: "Origin,Content-Type,Accept,Authorization", MaxAge: 600}))
	app.Use(helper.RequestLogger)
	ready := func(c *fiber.Ctx) error {
		if err := database.PingContext(c.UserContext()); err != nil {
			return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
		}
		return c.JSON(fiber.Map{"status": "ready"})
	}
	route.New(transactionhandler.NewHandler(usecase, validator.New()), ready).Register(app)
	return app
}

func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	var background sync.WaitGroup
	background.Add(2)
	go func() { defer background.Done(); s.worker.Run(ctx, 500*time.Millisecond) }()
	go func() { defer background.Done(); s.consumer.Run(ctx) }()
	errCh := make(chan error, 1)
	go func() {
		slog.Info("transaction service started", "port", s.config.Service.Port, "environment", s.config.Service.Environment)
		errCh <- s.app.Listen(s.config.Service.Port)
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			_ = s.publisher.Close()
			background.Wait()
			return s.close(fmt.Errorf("serve HTTP: %w", err))
		}
	case <-signals:
	}
	cancel()
	_ = s.publisher.Close()
	background.Wait()
	if err := s.app.ShutdownWithTimeout(15 * time.Second); err != nil {
		return s.close(fmt.Errorf("shutdown HTTP server: %w", err))
	}
	return s.close(nil)
}

func (s *Server) close(result error) error {
	_ = s.publisher.Close()
	if err := s.db.Close(); err != nil && result == nil {
		return fmt.Errorf("close transaction database: %w", err)
	}
	return result
}
