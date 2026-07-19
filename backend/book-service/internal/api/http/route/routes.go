package route

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"

	bookhandler "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/handler/book"
	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/handler/healthcheck"
	appmiddleware "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/middleware"
)

const swaggerUIBaseURL = "https://unpkg.com/swagger-ui-dist@5.32.9"

type Config struct {
	Issuer               string
	Audience             string
	TransactionServiceID string
}

type Dependency struct {
	Books        bookhandler.Handler
	Health       *healthhandler.Handler
	Introspector appmiddleware.Introspector
}

type Router struct {
	config Config
	deps   Dependency
}

func New(config Config, dependencies Dependency) *Router {
	return &Router{config: config, deps: dependencies}
}

func (r *Router) Register(app *fiber.App) {
	registerSwagger(app)
	app.Get("/health/liveness", r.deps.Health.Liveness)
	app.Get("/health/readiness", r.deps.Health.Readiness)

	read := appmiddleware.RequireGatewayScopes("books:read", "books:manage")
	manage := appmiddleware.RequireGatewayScopes("books:manage")
	books := app.Group("/api/v1/books")
	books.Get("/", read, r.deps.Books.List)
	books.Get("/:id", read, r.deps.Books.Get)
	books.Post("/", manage, r.deps.Books.Create)
	books.Patch("/:id", manage, r.deps.Books.Update)
	books.Delete("/:id", manage, r.deps.Books.Archive)

	internal := app.Group("/internal/v1/books")
	internal.Get("/:id/stock", r.internalAuth("book-stock:read"), r.deps.Books.Stock)
	internal.Put("/:id/reservations/:transactionId", r.internalAuth("book-stock:reserve"), r.deps.Books.Reserve)
	internal.Delete("/:id/reservations/:transactionId", r.internalAuth("book-stock:release"), r.deps.Books.Release)
}

func (r *Router) internalAuth(scope string) fiber.Handler {
	return appmiddleware.RequireInternal(r.deps.Introspector, appmiddleware.InternalPolicy{
		Issuer: r.config.Issuer, Audience: r.config.Audience,
		ClientID: r.config.TransactionServiceID, Scope: scope,
	})
}

func registerSwagger(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/docs/books", FilePath: "docs/swagger.json", Path: "swagger",
		Title: "Library Management Book API", SwaggerURL: swaggerUIBaseURL + "/swagger-ui-bundle.js",
		SwaggerPresetURL: swaggerUIBaseURL + "/swagger-ui-standalone-preset.js",
		SwaggerStylesURL: swaggerUIBaseURL + "/swagger-ui.css",
		Favicon32:        swaggerUIBaseURL + "/favicon-32x32.png", Favicon16: swaggerUIBaseURL + "/favicon-16x16.png",
	}))
}
