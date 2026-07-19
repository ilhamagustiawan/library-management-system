package server

import (
	"context"
	"fmt"
	"log/slog"
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

	bookhandler "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/handler/book"
	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/handler/healthcheck"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/route"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/config"
	dbinfra "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/db"
	bookrepository "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/db/repository/book"
	healthrepository "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/db/repository/healthcheck"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/messaging"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/oauth"
	bookusecase "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/book"
	healthusecase "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/healthcheck"
)

type Server struct {
	config config.Config
	db     *sqlx.DB
	app    *fiber.App
	rabbit *messaging.Rabbit
	worker *messaging.OutboxWorker
}

func New(ctx context.Context) (*Server, error) {
	configValue, err := config.LoadServer()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	database, err := dbinfra.Connect(ctx, dbinfra.Config{
		DSN: configValue.Database.DSN, MaxOpenConns: configValue.Database.MaxOpenConns,
		MaxIdleConns: configValue.Database.MaxIdleConns, ConnMaxLifetime: configValue.Database.ConnMaxLifetime,
		ConnMaxIdleTime: configValue.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return nil, err
	}
	bookRepository := bookrepository.NewRepository(database)
	rabbit, err := messaging.NewRabbit(messaging.RabbitConfig{
		URL: configValue.Rabbit.URL, Exchange: configValue.Rabbit.Exchange, DeadExchange: configValue.Rabbit.DeadExchange,
		ReturnQueue: configValue.Rabbit.ReturnQueue, AckQueue: configValue.Rabbit.AckQueue,
		DeadQueue: configValue.Rabbit.DeadQueue, ConfirmTimeout: configValue.Rabbit.ConfirmTimeout,
	}, messaging.NewReturnProcessor(bookRepository))
	if err != nil {
		_ = database.Close()
		return nil, err
	}
	app := buildApp(configValue, database)
	return &Server{config: configValue, db: database, app: app, rabbit: rabbit, worker: messaging.NewOutboxWorker(bookRepository, rabbit)}, nil
}

func buildApp(configValue config.Config, database *sqlx.DB) *fiber.App {
	app := fiber.New(fiberConfig(configValue))
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{ReferrerPolicy: "no-referrer"}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: configValue.Service.AllowedOrigin, AllowCredentials: false,
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(helper.RequestLogger)

	books := bookhandler.NewHandler(bookusecase.NewUsecase(bookrepository.NewRepository(database)), validator.New())
	health := healthhandler.NewHandler(healthusecase.NewUsecase(healthrepository.NewRepository(database)))
	introspector := oauth.NewClient(oauth.Config{
		URL: configValue.OAuth.IntrospectionURL, ClientID: configValue.OAuth.ClientID,
		ClientSecret: configValue.OAuth.ClientSecret, Timeout: configValue.OAuth.Timeout,
	})
	routes := route.New(route.Config{
		Issuer: configValue.OAuth.Issuer, Audience: configValue.OAuth.Audience,
		TransactionServiceID: configValue.OAuth.ServiceClientID,
	}, route.Dependency{Books: books, Health: health, Introspector: introspector})
	routes.Register(app)
	return app
}

func fiberConfig(configValue config.Config) fiber.Config {
	return fiber.Config{
		AppName: configValue.Service.Name, BodyLimit: 1 << 20,
		ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second,
		DisableStartupMessage: false, ErrorHandler: helper.ErrorHandler,
	}
}

func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	var background sync.WaitGroup
	background.Add(2)
	go func() { defer background.Done(); s.worker.Run(ctx) }()
	go func() { defer background.Done(); s.rabbit.RunConsumer(ctx) }()
	errorChannel := make(chan error, 1)
	go func() {
		slog.Info("book service started", "port", s.config.Service.Port, "environment", s.config.Service.Environment)
		errorChannel <- s.app.Listen(s.config.Service.Port)
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	select {
	case err := <-errorChannel:
		if err != nil {
			cancel()
			s.rabbit.Close()
			background.Wait()
			return s.close(fmt.Errorf("serve HTTP: %w", err))
		}
	case <-signals:
	}
	cancel()
	s.rabbit.Close()
	background.Wait()
	if err := s.app.ShutdownWithTimeout(15 * time.Second); err != nil {
		return s.close(fmt.Errorf("shutdown HTTP server: %w", err))
	}
	return s.close(nil)
}

func (s *Server) close(result error) error {
	s.rabbit.Close()
	if err := s.db.Close(); err != nil && result == nil {
		return fmt.Errorf("close database: %w", err)
	}
	return result
}
