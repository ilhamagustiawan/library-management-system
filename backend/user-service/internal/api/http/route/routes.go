package route

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"

	apidocs "github.com/ilhamagustiawan/library-management-system/backend/user-service/docs"
	healthhandler "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/handler/healthcheck"
	userhandler "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/handler/user"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/helper"
	appmiddleware "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/middleware"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

type Config struct {
	AllowedOrigin string
	RateMax       int
	RateWindow    time.Duration
}

func Register(app *fiber.App, config Config, health *healthhandler.Handler, registration registrationusecase.Usecase) {
	app.Use(swagger.New(swagger.Config{
		BasePath: "/api/v1/docs/users", FilePath: "swagger.json", FileContent: []byte(apidocs.SwaggerInfo.ReadDoc()),
		Path: "swagger", Title: "Library Management User API",
		SwaggerURL:       "https://unpkg.com/swagger-ui-dist@5.32.9/swagger-ui-bundle.js",
		SwaggerPresetURL: "https://unpkg.com/swagger-ui-dist@5.32.9/swagger-ui-standalone-preset.js",
		SwaggerStylesURL: "https://unpkg.com/swagger-ui-dist@5.32.9/swagger-ui.css",
	}))
	limit := limiter.New(limiter.Config{
		Max: config.RateMax, Expiration: config.RateWindow, KeyGenerator: helper.ClientIP,
		LimitReached: helper.RegistrationRateLimitReached,
	})
	app.Get("/health/liveness", health.Liveness)
	app.Get("/health/readiness", health.Readiness)
	app.Post("/api/v1/users", appmiddleware.RequireTrustedOrigin(config.AllowedOrigin), limit, userhandler.NewHandler(registration, validator.New()).Register)
}
