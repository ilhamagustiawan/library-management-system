package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/handler/healthcheck"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/route"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/config"
	infraDB "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/db"
	healthrepository "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/db/repository/healthcheck"
	outboxrepository "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/db/repository/outbox"
	registrationrepository "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/db/repository/registration"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/events/rabbitmq"
	identityinfra "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/infra/identity"
	healthusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/healthcheck"
	outboxusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/outbox"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

type Server struct {
	config    config.Config
	db        *sqlx.DB
	app       *fiber.App
	relay     *outboxusecase.Relay
	publisher *rabbitmq.Publisher
}

func New(ctx context.Context) (*Server, error) {
	serviceConfig, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	database, err := infraDB.Connect(ctx, infraDB.Config{
		DSN: serviceConfig.Database.DSN, MaxOpenConns: serviceConfig.Database.MaxOpenConns,
		MaxIdleConns: serviceConfig.Database.MaxIdleConns, ConnMaxLifetime: serviceConfig.Database.ConnMaxLifetime,
		ConnMaxIdleTime: serviceConfig.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout:       serviceConfig.Auth.Timeout,
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse },
	}
	identities := identityinfra.NewClient(identityinfra.Config{
		BaseURL: serviceConfig.Auth.BaseURL, ClientID: serviceConfig.Auth.ClientID,
		ClientSecret: serviceConfig.Auth.ClientSecret, Scope: serviceConfig.Auth.Scope, Attempts: serviceConfig.Auth.Attempts,
	}, httpClient)
	registration := registrationusecase.New(registrationrepository.NewRepository(database), identities)
	health := healthhandler.NewHandler(healthusecase.New(healthrepository.NewRepository(database)))

	app := fiber.New(fiberConfig(serviceConfig))
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{ReferrerPolicy: "no-referrer"}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: serviceConfig.Service.AllowedOrigin, AllowCredentials: true,
		AllowMethods: "GET,POST,OPTIONS", AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(helper.RequestLogger)
	route.Register(app, route.Config{
		AllowedOrigin: serviceConfig.Service.AllowedOrigin, RateMax: serviceConfig.Rate.Max, RateWindow: serviceConfig.Rate.Window,
	}, health, registration)

	publisher := rabbitmq.NewPublisher(rabbitmq.Config{
		URL: serviceConfig.RabbitMQ.URL, Exchange: serviceConfig.RabbitMQ.Exchange,
		RoutingKey: serviceConfig.RabbitMQ.RoutingKey, Queue: serviceConfig.RabbitMQ.Queue,
		ConfirmTimeout: serviceConfig.RabbitMQ.ConfirmTimeout,
	})
	relay := outboxusecase.NewRelay(outboxrepository.NewRepository(database), publisher, outboxusecase.Config{
		WorkerID: uuid.NewString(), BatchSize: serviceConfig.Outbox.BatchSize, Lease: serviceConfig.Outbox.Lease,
		PollInterval: serviceConfig.Outbox.PollInterval, BaseRetry: serviceConfig.Outbox.BaseRetry, MaxRetry: serviceConfig.Outbox.MaxRetry,
	})
	return &Server{config: serviceConfig, db: database, app: app, relay: relay, publisher: publisher}, nil
}

func fiberConfig(config config.Config) fiber.Config {
	return fiber.Config{
		AppName: config.Service.Name, BodyLimit: 1 << 20, ReadTimeout: 10 * time.Second,
		WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, ErrorHandler: helper.ErrorHandler,
		ProxyHeader: "X-Real-IP", EnableIPValidation: true, EnableTrustedProxyCheck: true,
		TrustedProxies: config.Service.TrustedProxies,
	}
}

func (s *Server) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	relayDone := make(chan struct{})
	go func() {
		defer close(relayDone)
		s.relay.Run(ctx)
	}()

	errCh := make(chan error, 1)
	go func() {
		slog.Info("user service started", "port", s.config.Service.Port, "environment", s.config.Service.Environment)
		errCh <- s.app.Listen(s.config.Service.Port)
	}()
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case err := <-errCh:
		if err != nil {
			cancel()
			<-relayDone
			return s.close(fmt.Errorf("serve HTTP: %w", err))
		}
	case <-signals:
	}
	cancel()
	if err := s.app.ShutdownWithTimeout(15 * time.Second); err != nil {
		<-relayDone
		return s.close(fmt.Errorf("shutdown HTTP server: %w", err))
	}
	<-relayDone
	return s.close(nil)
}

func (s *Server) close(result error) error {
	if err := s.publisher.Close(); err != nil && result == nil {
		result = fmt.Errorf("close RabbitMQ publisher: %w", err)
	}
	if err := s.db.Close(); err != nil && result == nil {
		result = fmt.Errorf("close database: %w", err)
	}
	return result
}
