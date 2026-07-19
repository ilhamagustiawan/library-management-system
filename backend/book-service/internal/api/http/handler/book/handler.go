package book

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/request"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	bookusecase "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/book"
)

type Handler interface {
	List(*fiber.Ctx) error
	Get(*fiber.Ctx) error
	Create(*fiber.Ctx) error
	Update(*fiber.Ctx) error
	Archive(*fiber.Ctx) error
	Stock(*fiber.Ctx) error
	Reserve(*fiber.Ctx) error
	Release(*fiber.Ctx) error
}

type handler struct {
	usecase  bookusecase.Usecase
	validate *validator.Validate
}

func NewHandler(usecase bookusecase.Usecase, validate *validator.Validate) Handler {
	return &handler{usecase: usecase, validate: validate}
}

// List returns the active book catalog.
// @Summary List books
// @Tags Books
// @Produce json
// @Param q query string false "Title, author, or ISBN search"
// @Param availableOnly query bool false "Only books with available stock"
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Page size" default(20) maximum(100)
// @Param sortBy query string false "Sort field" Enums(title,author,createdAt)
// @Param sortOrder query string false "Sort order" Enums(asc,desc)
// @Success 200 {object} response.BookPageSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /api/v1/books [get]
func (h *handler) List(ctx *fiber.Ctx) error {
	input, err := listInput(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}
	page, err := h.usecase.List(ctx.UserContext(), input)
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, http.StatusOK, response.NewBookPage(page))
}

// Get returns one active book.
// @Summary Get book
// @Tags Books
// @Produce json
// @Param id path string true "Book UUID"
// @Success 200 {object} response.BookSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /api/v1/books/{id} [get]
func (h *handler) Get(ctx *fiber.Ctx) error {
	book, err := h.usecase.Get(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, http.StatusOK, response.NewBook(book))
}

// Create adds a catalog book and its initial stock.
// @Summary Create book
// @Tags Books
// @Accept json
// @Produce json
// @Param request body request.CreateBook true "Book data"
// @Success 201 {object} response.BookSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/books [post]
func (h *handler) Create(ctx *fiber.Ctx) error {
	var input request.CreateBook
	if err := request.DecodeStrictJSON(ctx, &input); err != nil || h.validate.Struct(input) != nil {
		return response.Error(ctx, validation("invalid book data"))
	}
	book, err := h.usecase.Create(ctx.UserContext(), bookusecase.CreateInput{
		ISBN: input.ISBN, Title: input.Title, Author: input.Author, Description: input.Description,
		CoverURL: input.CoverURL, PublicationYear: input.PublicationYear, TotalCopies: input.TotalCopies,
	})
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, http.StatusCreated, response.NewBook(book))
}

// Update changes provided book fields only.
// @Summary Update book
// @Tags Books
// @Accept json
// @Produce json
// @Param id path string true "Book UUID"
// @Param request body request.UpdateBook true "Partial book data"
// @Success 200 {object} response.BookSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Failure 422 {object} response.ErrorResponse
// @Router /api/v1/books/{id} [patch]
func (h *handler) Update(ctx *fiber.Ctx) error {
	var input request.UpdateBook
	if err := request.DecodeStrictJSON(ctx, &input); err != nil {
		return response.Error(ctx, validation("invalid book data"))
	}
	book, err := h.usecase.Update(ctx.UserContext(), ctx.Params("id"), patch(input))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, http.StatusOK, response.NewBook(book))
}

// Archive removes a book from active catalog results.
// @Summary Archive book
// @Tags Books
// @Param id path string true "Book UUID"
// @Success 204
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /api/v1/books/{id} [delete]
func (h *handler) Archive(ctx *fiber.Ctx) error {
	if err := h.usecase.Archive(ctx.UserContext(), ctx.Params("id")); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

// Stock returns current stock for Transaction Service.
// @Summary Read internal book stock
// @Tags Internal stock
// @Produce json
// @Param id path string true "Book UUID"
// @Success 200 {object} response.StockSuccess
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /internal/v1/books/{id}/stock [get]
func (h *handler) Stock(ctx *fiber.Ctx) error {
	stock, err := h.usecase.Stock(ctx.UserContext(), ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, err)
	}
	return response.Success(ctx, http.StatusOK, stock)
}

// Reserve atomically verifies and consumes one available copy.
// @Summary Reserve book stock
// @Tags Internal stock
// @Produce json
// @Param id path string true "Book UUID"
// @Param transactionId path string true "Transaction UUID"
// @Success 200 {object} response.ReservationSuccess "Existing active reservation"
// @Success 201 {object} response.ReservationSuccess "Reservation created"
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /internal/v1/books/{id}/reservations/{transactionId} [put]
func (h *handler) Reserve(ctx *fiber.Ctx) error {
	reservation, created, err := h.usecase.Reserve(ctx.UserContext(), ctx.Params("id"), ctx.Params("transactionId"))
	if err != nil {
		return response.Error(ctx, err)
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	return response.Success(ctx, status, reservation)
}

// Release restores one reserved copy. Repeated calls remain successful.
// @Summary Release book stock
// @Tags Internal stock
// @Param id path string true "Book UUID"
// @Param transactionId path string true "Transaction UUID"
// @Success 204
// @Failure 401 {object} response.ErrorResponse
// @Failure 403 {object} response.ErrorResponse
// @Failure 404 {object} response.ErrorResponse
// @Failure 409 {object} response.ErrorResponse
// @Router /internal/v1/books/{id}/reservations/{transactionId} [delete]
func (h *handler) Release(ctx *fiber.Ctx) error {
	if err := h.usecase.Release(ctx.UserContext(), ctx.Params("id"), ctx.Params("transactionId")); err != nil {
		return response.Error(ctx, err)
	}
	return ctx.SendStatus(http.StatusNoContent)
}

func listInput(ctx *fiber.Ctx) (bookusecase.ListInput, error) {
	page, err := optionalInt(ctx.Query("page"))
	if err != nil {
		return bookusecase.ListInput{}, validation("invalid pagination")
	}
	pageSize, err := optionalInt(ctx.Query("pageSize"))
	if err != nil {
		return bookusecase.ListInput{}, validation("invalid pagination")
	}
	available := false
	if raw := ctx.Query("availableOnly"); raw != "" {
		available, err = strconv.ParseBool(raw)
		if err != nil {
			return bookusecase.ListInput{}, validation("invalid availability filter")
		}
	}
	return bookusecase.ListInput{
		Query: ctx.Query("q"), AvailableOnly: available, Page: page, PageSize: pageSize,
		SortBy: strings.TrimSpace(ctx.Query("sortBy")), SortOrder: strings.ToLower(strings.TrimSpace(ctx.Query("sortOrder"))),
	}, nil
}

func optionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}

func patch(input request.UpdateBook) entity.BookPatch {
	return entity.BookPatch{
		ISBN:            entity.Change[string]{Set: input.ISBN.Present, Value: input.ISBN.Value},
		Title:           entity.Change[string]{Set: input.Title.Present, Value: input.Title.Value},
		Author:          entity.Change[string]{Set: input.Author.Present, Value: input.Author.Value},
		Description:     entity.Change[string]{Set: input.Description.Present, Value: input.Description.Value},
		CoverURL:        entity.Change[string]{Set: input.CoverURL.Present, Value: input.CoverURL.Value},
		PublicationYear: entity.Change[int]{Set: input.PublicationYear.Present, Value: input.PublicationYear.Value},
		TotalCopies:     entity.Change[int]{Set: input.TotalCopies.Present, Value: input.TotalCopies.Value},
	}
}

func validation(message string) error {
	return errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, message, nil, nil)
}
