package helper

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) && fiberErr.Code == fiber.StatusNotFound {
		return response.Write(c, fiber.StatusNotFound, errs.CodeNotFound, "endpoint not found")
	}
	return response.Error(c, err)
}

func RequestLogger(c *fiber.Ctx) error {
	started := time.Now()
	err := c.Next()
	slog.Info("http request", "method", c.Method(), "path", c.Path(), "status", c.Response().StatusCode(), "duration", time.Since(started))
	return err
}

func ClientIP(c *fiber.Ctx) string { return c.IP() }

func RegistrationRateLimitReached(c *fiber.Ctx) error {
	return response.Write(c, fiber.StatusTooManyRequests, errs.CodeRateLimited, "too many registration attempts; retry after the rate-limit window")
}
