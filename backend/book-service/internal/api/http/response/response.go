package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
)

type SuccessResponse struct {
	Code string `json:"code" example:"LMS-200000"`
	Data any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code" example:"LMS-422001"`
	Message string `json:"message" example:"invalid request data"`
	Data    any    `json:"data,omitempty"`
}

func Success(ctx *fiber.Ctx, status int, data any) error {
	return ctx.Status(status).JSON(SuccessResponse{Code: errs.CodeSuccess, Data: data})
}

func Error(ctx *fiber.Ctx, err error) error {
	var domainError *errs.Error
	if errors.As(err, &domainError) {
		return ctx.Status(domainError.HTTPStatus).JSON(ErrorResponse{
			Code: domainError.ErrorCode, Message: domainError.Message, Data: domainError.Data,
		})
	}
	return ctx.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Code: errs.CodeInternal, Message: "an unexpected error occurred",
	})
}

func Write(ctx *fiber.Ctx, status int, code, message string) error {
	return ctx.Status(status).JSON(ErrorResponse{Code: code, Message: message})
}
