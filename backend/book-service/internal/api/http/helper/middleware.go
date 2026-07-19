package helper

import (
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var fiberError *fiber.Error
	if errors.As(err, &fiberError) && fiberError.Code == fiber.StatusNotFound {
		return response.Write(ctx, fiber.StatusNotFound, errs.CodeNotFound, "endpoint not found")
	}
	return response.Error(ctx, err)
}

func RequestLogger(ctx *fiber.Ctx) error {
	started := time.Now()
	err := ctx.Next()
	slog.Info("http request", "method", ctx.Method(), "path", ctx.Path(), "status", ctx.Response().StatusCode(), "duration", time.Since(started))
	return err
}
