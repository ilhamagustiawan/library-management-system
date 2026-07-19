package route

import (
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"

	apidocs "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/docs"
	transactionhandler "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/handler/transaction"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/middleware"
)

const swaggerUIBaseURL = "https://unpkg.com/swagger-ui-dist@5.32.9"

type Router struct {
	transaction *transactionhandler.Handler
	ready       fiber.Handler
}

func New(transaction *transactionhandler.Handler, ready fiber.Handler) *Router {
	return &Router{transaction: transaction, ready: ready}
}

func (r *Router) Register(app *fiber.App) {
	registerSwagger(app)
	app.Get("/health/liveness", func(c *fiber.Ctx) error { return c.JSON(fiber.Map{"status": "ok"}) })
	app.Get("/health/readiness", r.ready)
	transactions := app.Group("/api/v1/transactions")
	transactions.Post("/loans", middleware.RequireScope("loans:borrow:self"), r.transaction.Borrow)
	transactions.Post("/loans/:loanId/return", middleware.RequireScope("loans:return:self"), r.transaction.ReturnSelf)
	transactions.Get("/me", middleware.RequireScope("transactions:read:self"), r.transaction.ListSelf)
	transactions.Get("/admin", middleware.RequireScope("transactions:read:any"), r.transaction.ListAny)
	transactions.Post("/admin/loans/:loanId/return", middleware.RequireScope("loans:return:any"), r.transaction.ReturnAny)
}

func registerSwagger(app *fiber.App) {
	app.Use(swagger.New(swagger.Config{
		BasePath:         "/api/v1/docs/transactions",
		FilePath:         "swagger.json",
		FileContent:      []byte(apidocs.SwaggerInfo.ReadDoc()),
		Path:             "swagger",
		Title:            "Library Management Transaction API",
		SwaggerURL:       swaggerUIBaseURL + "/swagger-ui-bundle.js",
		SwaggerPresetURL: swaggerUIBaseURL + "/swagger-ui-standalone-preset.js",
		SwaggerStylesURL: swaggerUIBaseURL + "/swagger-ui.css",
		Favicon32:        swaggerUIBaseURL + "/favicon-32x32.png",
		Favicon16:        swaggerUIBaseURL + "/favicon-16x16.png",
	}))
}
