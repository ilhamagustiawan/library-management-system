package response

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	transactionusecase "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/usecase/transaction"
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

type LoanSuccess struct {
	Code string       `json:"code" example:"LMS-200000"`
	Data *entity.Loan `json:"data"`
}

type PageSuccess struct {
	Code string                  `json:"code" example:"LMS-200000"`
	Data transactionusecase.Page `json:"data"`
}

func Success(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(SuccessResponse{Code: errs.CodeSuccess, Data: data})
}
func Error(c *fiber.Ctx, err error) error {
	var domainErr *errs.Error
	if errors.As(err, &domainErr) {
		return c.Status(domainErr.HTTPStatus).JSON(ErrorResponse{Code: domainErr.ErrorCode, Message: domainErr.Message, Data: domainErr.Data})
	}
	return c.Status(http.StatusInternalServerError).JSON(ErrorResponse{Code: errs.CodeInternal, Message: "an unexpected error occurred"})
}
func Write(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorResponse{Code: code, Message: message})
}
func ValidationError(c *fiber.Ctx, message string) error {
	return Write(c, http.StatusUnprocessableEntity, errs.CodeValidation, message)
}
