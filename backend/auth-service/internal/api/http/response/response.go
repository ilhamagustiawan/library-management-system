package response

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	domainerrs "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
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

// UserSuccess documents endpoints returning the standard success envelope.
type UserSuccess struct {
	Code string `json:"code" example:"LMS-200000"`
	Data User   `json:"data"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(SuccessResponse{Code: domainerrs.CodeSuccess, Data: data})
}

func Error(c *fiber.Ctx, err error) error {
	var domainErr *domainerrs.Error
	if errors.As(err, &domainErr) {
		return c.Status(domainErr.HTTPStatus).JSON(ErrorResponse{
			Code: domainErr.ErrorCode, Message: domainErr.Message, Data: domainErr.Data,
		})
	}
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Code: domainerrs.CodeInternal, Message: "an unexpected error occurred",
	})
}

func Write(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorResponse{Code: code, Message: message})
}

func ValidationError(c *fiber.Ctx, message string) error {
	return Error(c, domainerrs.New(
		http.StatusUnprocessableEntity,
		domainerrs.CodeValidation,
		strings.TrimSpace(message),
		nil,
		nil,
	))
}
