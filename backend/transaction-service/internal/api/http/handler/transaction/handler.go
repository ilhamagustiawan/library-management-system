package transaction

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/middleware"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/request"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	transactionusecase "github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/usecase/transaction"
)

type Handler struct {
	usecase  transactionusecase.Usecase
	validate *validator.Validate
}

func NewHandler(usecase transactionusecase.Usecase, validate *validator.Validate) *Handler {
	return &Handler{usecase: usecase, validate: validate}
}

// Borrow creates a seven-day loan after Book Service atomically reserves stock.
// @Summary Borrow book
// @Tags Transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body request.Borrow true "Book to borrow"
// @Success 201 {object} response.LoanSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 503 {object} response.ErrorResponse
// @Router /api/v1/transactions/loans [post]
func (h *Handler) Borrow(c *fiber.Ctx) error {
	credential, ok := middleware.Current(c)
	if !ok {
		return response.Write(c, http.StatusUnauthorized, errs.CodeInvalidToken, "authentication required")
	}
	var input request.Borrow
	if err := request.DecodeStrictJSON(c, &input); err != nil || h.validate.Struct(input) != nil {
		return response.ValidationError(c, "invalid borrow data")
	}
	loan, err := h.usecase.Borrow(c.UserContext(), transactionusecase.BorrowInput{MemberID: credential.Subject, BookID: input.BookID})
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusCreated, loan)
}

// ReturnSelf returns the authenticated member's loan and waits up to five seconds for stock acknowledgement.
// @Summary Return own loan
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param loanId path string true "Loan UUID" Format(uuid)
// @Success 200 {object} response.LoanSuccess "Book stock update confirmed"
// @Success 202 {object} response.LoanSuccess "Return committed; stock acknowledgement pending"
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/transactions/loans/{loanId}/return [post]
func (h *Handler) ReturnSelf(c *fiber.Ctx) error { return h.returnLoan(c, false) }

// ReturnAny lets an authorized librarian return any active loan.
// @Summary Return any member loan
// @Tags Admin transactions
// @Produce json
// @Security BearerAuth
// @Param loanId path string true "Loan UUID" Format(uuid)
// @Success 200 {object} response.LoanSuccess "Book stock update confirmed"
// @Success 202 {object} response.LoanSuccess "Return committed; stock acknowledgement pending"
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/transactions/admin/loans/{loanId}/return [post]
func (h *Handler) ReturnAny(c *fiber.Ctx) error { return h.returnLoan(c, true) }

func (h *Handler) returnLoan(c *fiber.Ctx, allowAny bool) error {
	credential, ok := middleware.Current(c)
	if !ok {
		return response.Write(c, http.StatusUnauthorized, errs.CodeInvalidToken, "authentication required")
	}
	loanID := c.Params("loanId")
	if _, err := uuid.Parse(loanID); err != nil {
		return response.ValidationError(c, "invalid loan ID")
	}
	loan, confirmed, err := h.usecase.Return(c.UserContext(), transactionusecase.ReturnInput{LoanID: loanID, MemberID: credential.Subject, AllowAnyMember: allowAny})
	if err != nil {
		return response.Error(c, err)
	}
	if !confirmed {
		c.Set(fiber.HeaderRetryAfter, "2")
		return response.Success(c, http.StatusAccepted, loan)
	}
	return response.Success(c, http.StatusOK, loan)
}

// ListSelf returns only the authenticated member's transaction history.
// @Summary List own transaction history
// @Tags Transactions
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1) minimum(1)
// @Param pageSize query int false "Page size" default(20) minimum(1) maximum(100)
// @Success 200 {object} response.PageSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/transactions/me [get]
func (h *Handler) ListSelf(c *fiber.Ctx) error {
	credential, ok := middleware.Current(c)
	if !ok {
		return response.Write(c, http.StatusUnauthorized, errs.CodeInvalidToken, "authentication required")
	}
	return h.list(c, transactionusecase.ListInput{MemberID: credential.Subject})
}

// ListAny returns paginated transaction history across members.
// @Summary List all member transactions
// @Tags Admin transactions
// @Produce json
// @Security BearerAuth
// @Param memberId query string false "Member UUID" Format(uuid)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param pageSize query int false "Page size" default(20) minimum(1) maximum(100)
// @Success 200 {object} response.PageSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/transactions/admin [get]
func (h *Handler) ListAny(c *fiber.Ctx) error {
	memberID := c.Query("memberId")
	if memberID != "" {
		if _, err := uuid.Parse(memberID); err != nil {
			return response.ValidationError(c, "invalid member ID")
		}
	}
	return h.list(c, transactionusecase.ListInput{MemberID: memberID, AllMembers: true})
}

func (h *Handler) list(c *fiber.Ctx, input transactionusecase.ListInput) error {
	page, pageSize, err := pagination(c)
	if err != nil {
		return response.ValidationError(c, "invalid pagination")
	}
	input.Page, input.PageSize = page, pageSize
	result, err := h.usecase.List(c.UserContext(), input)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Success(c, http.StatusOK, result)
}

func pagination(c *fiber.Ctx) (int, int, error) {
	page, pageSize := 1, 20
	var err error
	if raw := c.Query("page"); raw != "" {
		page, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, err
		}
	}
	if raw := c.Query("pageSize"); raw != "" {
		pageSize, err = strconv.Atoi(raw)
		if err != nil {
			return 0, 0, err
		}
	}
	return page, pageSize, nil
}
