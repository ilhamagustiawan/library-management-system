package route

import (
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	apidocs "github.com/ilhamagustiawan/library-management-system/backend/auth-service/docs"
	authhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/auth"
	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/healthcheck"
	oauthhandler "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/handler/oauth"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/helper"
	appmiddleware "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/middleware"
	oauthinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/oauth"
)

type Config struct {
	AllowedOrigin  string
	AuthRateMax    int
	AuthRateWindow time.Duration
}

type Dependency struct {
	Auth     authhandler.Handler
	Health   *healthhandler.Handler
	OAuth    *oauthinfra.AuthorizationServer
	OAuthAPI *oauthhandler.Handler
}

type Router struct {
	config Config
	deps   Dependency
}

const swaggerUIBaseURL = "https://unpkg.com/swagger-ui-dist@5.32.9"

func New(config Config, deps Dependency) *Router {
	return &Router{config: config, deps: deps}
}

func (r *Router) Register(app *fiber.App) {
	registerSwagger(app)

	trustedOrigin := appmiddleware.RequireTrustedOrigin(r.config.AllowedOrigin)
	authLimit := limiter.New(limiter.Config{
		Max: r.config.AuthRateMax, Expiration: r.config.AuthRateWindow,
		KeyGenerator: helper.ClientIP,
		LimitReached: helper.AuthRateLimitReached,
	})

	app.Get("/health/liveness", r.deps.Health.Liveness)
	app.Get("/health/readiness", r.deps.Health.Readiness)
	app.Get("/.well-known/oauth-authorization-server", r.deps.OAuthAPI.Metadata)

	auth := app.Group("/api/v1/auth")
	auth.Post("/register", trustedOrigin, authLimit, r.deps.Auth.Register)
	auth.Post("/login", trustedOrigin, authLimit, r.deps.Auth.Login)
	auth.Post("/logout", trustedOrigin, r.deps.Auth.Logout)
	auth.Get("/me", r.deps.Auth.Me)

	app.Get("/api/v1/oauth/userinfo", r.deps.OAuthAPI.UserInfo)
	app.Get("/oauth/authorize", adaptor.HTTPHandler(r.deps.OAuth.AuthorizeHandler()))
	app.Post("/oauth/token", authLimit, adaptor.HTTPHandler(r.deps.OAuth.TokenHandler()))
	app.Post("/oauth/introspect", adaptor.HTTPHandler(r.deps.OAuth.IntrospectionHandler()))
}

func registerSwagger(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath:         "/api/v1/docs/auth",
		FilePath:         "swagger.json",
		FileContent:      []byte(apidocs.SwaggerInfo.ReadDoc()),
		Path:             "swagger",
		Title:            "Library Management Auth API",
		SwaggerURL:       swaggerUIBaseURL + "/swagger-ui-bundle.js",
		SwaggerPresetURL: swaggerUIBaseURL + "/swagger-ui-standalone-preset.js",
		SwaggerStylesURL: swaggerUIBaseURL + "/swagger-ui.css",
		Favicon32:        swaggerUIBaseURL + "/favicon-32x32.png",
		Favicon16:        swaggerUIBaseURL + "/favicon-16x16.png",
	}))
}
