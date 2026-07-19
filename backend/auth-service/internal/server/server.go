package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"

	authhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/auth"
	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/healthcheck"
	identityhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/identity"
	oauthhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/oauth"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/helper"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/route"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
	authinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/auth"
	infraDB "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db"
	healthRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/healthcheck"
	identityRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/identity"
	oauthClientRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/oauthclient"
	oauthTokenRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/oauthtoken"
	sessionRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/session"
	userRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/user"
	oauthinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/oauth"
	authusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/auth"
	healthusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/healthcheck"
	identityusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/identity"
	oauthusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/oauth"
)

type Server struct {
	config config.Config
	db     *sqlx.DB
	app    *fiber.App
}

func New(ctx context.Context) (*Server, error) {
	cfg, err := config.LoadServer()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	db, err := infraDB.Connect(ctx, infraDB.Config{
		DSN: cfg.Database.DSN, MaxOpenConns: cfg.Database.MaxOpenConns, MaxIdleConns: cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime, ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return nil, err
	}

	app, err := buildApp(ctx, cfg, db)
	if err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Server{config: cfg, db: db, app: app}, nil
}

func buildApp(ctx context.Context, cfg config.Config, db *sqlx.DB) (*fiber.App, error) {
	hasher := authinfra.NewPasswordHasher(cfg.Auth.BcryptCost)
	dummyHash, err := hasher.Hash("not-the-user-password")
	if err != nil {
		return nil, fmt.Errorf("create dummy password hash: %w", err)
	}
	clients := oauthClientRepo.NewStore(db)
	users := userRepo.NewRepository(db)
	sessions := sessionRepo.NewRepository(db)
	authUC := authusecase.NewAuthUsecase(
		users, sessions, hasher, authinfra.NewTokenGenerator(32),
		authusecase.Config{SessionTTL: cfg.Auth.SessionTTL, DummyPasswordHash: dummyHash},
	)

	tokens := oauthTokenRepo.NewStore(db)
	manager := manage.NewDefaultManager()
	manager.MapClientStorage(clients)
	manager.MapTokenStorage(tokens)
	oauthServer, err := oauthinfra.NewAuthorizationServer(manager, authUC, clients, oauthinfra.Config{
		Issuer: cfg.OAuth.Issuer, LoginURL: cfg.OAuth.LoginURL, SessionCookieName: cfg.Auth.SessionCookieName,
		CodeTTL: cfg.OAuth.CodeTTL, AccessTokenTTL: cfg.OAuth.AccessTokenTTL, RefreshTokenTTL: cfg.OAuth.RefreshTokenTTL,
		SupportedScopes: cfg.OAuth.SupportedScopes, JWTSigningKey: cfg.OAuth.JWTSigningKey,
	})
	if err != nil {
		return nil, err
	}

	app := fiber.New(fiberConfig(cfg))
	app.Use(requestid.New())
	app.Use(recover.New())
	app.Use(helmet.New(helmet.Config{ReferrerPolicy: "no-referrer"}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.Service.AllowedOrigin, AllowCredentials: true,
		AllowMethods: "GET,POST,OPTIONS", AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))
	app.Use(helper.RequestLogger)

	authHTTP := authhandler.NewHandler(authUC, validator.New(), authhandler.Config{
		Issuer: cfg.OAuth.Issuer, SessionCookieName: cfg.Auth.SessionCookieName,
		SessionCookieDomain: cfg.Auth.SessionCookieDomain, SessionCookieSecure: cfg.Auth.SessionCookieSecure,
	})
	oauthUC := oauthusecase.NewUsecase(oauthServer, users)
	healthUC := healthusecase.NewUsecase(healthRepo.NewRepository(db))
	oauthHTTP := oauthhandler.NewHandler(oauthUC)
	identityUC := identityusecase.New(identityRepo.NewRepository(db), hasher)
	identityHTTP := identityhandler.NewHandler(identityUC, oauthServer, validator.New())
	routes := route.New(route.Config{
		AllowedOrigin: cfg.Service.AllowedOrigin, AuthRateMax: cfg.Auth.RateLimitMax, AuthRateWindow: cfg.Auth.RateLimitWindow,
	}, route.Dependency{Auth: authHTTP, Identity: identityHTTP, Health: healthhandler.NewHandler(healthUC), OAuth: oauthServer, OAuthAPI: oauthHTTP})
	routes.Register(app)
	return app, nil
}

func fiberConfig(cfg config.Config) fiber.Config {
	return fiber.Config{
		AppName: cfg.Service.Name, BodyLimit: 1 << 20, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second,
		IdleTimeout: 60 * time.Second, DisableStartupMessage: false, ErrorHandler: helper.ErrorHandler,
		ProxyHeader: "X-Real-IP", EnableIPValidation: true, EnableTrustedProxyCheck: true,
		TrustedProxies: cfg.Service.TrustedProxies,
	}
}

func (s *Server) Start() error {
	errCh := make(chan error, 1)
	go s.listen(errCh)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)

	select {
	case err := <-errCh:
		if err != nil {
			_ = s.db.Close()
			return fmt.Errorf("serve HTTP: %w", err)
		}
	case <-signals:
	}

	if err := s.app.ShutdownWithTimeout(15 * time.Second); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close database: %w", err)
	}
	return nil
}

func (s *Server) listen(errCh chan<- error) {
	slog.Info("auth service started", "port", s.config.Service.Port, "environment", s.config.Service.Environment)
	errCh <- s.app.Listen(s.config.Service.Port)
}
