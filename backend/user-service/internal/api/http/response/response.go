package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
)

type SuccessResponse struct {
	Code string `json:"code"`
	Data any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(SuccessResponse{Code: errs.CodeSuccess, Data: data})
}

func Error(c *fiber.Ctx, err error) error {
	var domainErr *errs.Error
	if errors.As(err, &domainErr) {
		return c.Status(domainErr.HTTPStatus).JSON(ErrorResponse{
			Code: domainErr.ErrorCode, Message: domainErr.Message, Data: domainErr.Data,
		})
	}
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{
		Code: errs.CodeInternal, Message: "an unexpected error occurred; pending registration state remains preserved",
	})
}

func ValidationError(c *fiber.Ctx) error {
	return Error(c, errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, "invalid registration data", nil, nil))
}

func Write(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorResponse{Code: code, Message: message})
}
