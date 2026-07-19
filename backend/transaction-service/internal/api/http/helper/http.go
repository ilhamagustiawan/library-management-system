package helper

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if fiberError, ok := err.(*fiber.Error); ok {
		if fiberError.Code == http.StatusNotFound {
			return response.Write(c, fiberError.Code, errs.CodeLoanNotFound, "route was not found")
		}
		return response.Write(c, fiberError.Code, errs.CodeValidation, fiberError.Message)
	}
	return response.Error(c, err)
}

func RequestLogger(c *fiber.Ctx) error {
	started := time.Now()
	err := c.Next()
	slog.Info("HTTP request", "method", c.Method(), "path", c.Path(), "status", c.Response().StatusCode(), "duration_ms", time.Since(started).Milliseconds(), "request_id", c.GetRespHeader(fiber.HeaderXRequestID))
	return err
}
