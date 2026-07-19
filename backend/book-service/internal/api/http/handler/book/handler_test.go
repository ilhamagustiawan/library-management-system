package book

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	bookusecase "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/book"
)

type fakeUsecase struct {
	createdInput bookusecase.CreateInput
	book         *entity.Book
}

func (u *fakeUsecase) Create(_ context.Context, input bookusecase.CreateInput) (*entity.Book, error) {
	u.createdInput = input
	return u.book, nil
}
func (u *fakeUsecase) Get(context.Context, string) (*entity.Book, error) { return u.book, nil }
func (u *fakeUsecase) List(context.Context, bookusecase.ListInput) (bookusecase.Page, error) {
	return bookusecase.Page{Items: []*entity.Book{}, Page: 1, PageSize: 20}, nil
}
func (u *fakeUsecase) Update(context.Context, string, bookusecase.UpdateInput) (*entity.Book, error) {
	return u.book, nil
}
func (u *fakeUsecase) Archive(context.Context, string) error                { return nil }
func (u *fakeUsecase) Stock(context.Context, string) (*entity.Stock, error) { return nil, nil }
func (u *fakeUsecase) Reserve(context.Context, string, string) (*entity.Reservation, bool, error) {
	return nil, false, nil
}
func (u *fakeUsecase) Release(context.Context, string, string) error { return nil }

func TestCreateRejectsUnknownJSONField(t *testing.T) {
	app := fiber.New()
	handler := NewHandler(&fakeUsecase{}, validator.New())
	app.Post("/api/v1/books", handler.Create)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(
		`{"isbn":"9780132350884","title":"Clean Code","author":"Robert C. Martin","totalCopies":3,"role":"admin"}`,
	))
	request.Header.Set("Content-Type", "application/json")
	result, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", result.StatusCode)
	}
}

func TestCreateReturnsStandardBookEnvelope(t *testing.T) {
	bookID := "f81d4fae-7dec-41d0-a765-00a0c91e6bf6"
	coverURL := "https://covers.openlibrary.org/b/id/10521270-M.jpg"
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	usecase := &fakeUsecase{book: &entity.Book{
		ID: bookID, ISBN: "9780132350884", Title: "Clean Code", Author: "Robert C. Martin",
		CoverURL: &coverURL, TotalCopies: 3, AvailableCopies: 3, CreatedAt: now, UpdatedAt: now,
	}}
	app := fiber.New()
	handler := NewHandler(usecase, validator.New())
	app.Post("/api/v1/books", handler.Create)

	request := httptest.NewRequest(http.MethodPost, "/api/v1/books", strings.NewReader(
		`{"isbn":"9780132350884","title":"Clean Code","author":"Robert C. Martin","coverUrl":"https://covers.openlibrary.org/b/id/10521270-M.jpg","totalCopies":3}`,
	))
	request.Header.Set("Content-Type", "application/json")
	result, err := app.Test(request)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer result.Body.Close()
	if result.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want 201", result.StatusCode)
	}
	var payload response.SuccessResponse
	if err := json.NewDecoder(result.Body).Decode(&payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	data, ok := payload.Data.(map[string]any)
	if payload.Code != "LMS-200000" || usecase.createdInput.TotalCopies != 3 ||
		usecase.createdInput.CoverURL == nil || *usecase.createdInput.CoverURL != coverURL ||
		!ok || data["coverUrl"] != coverURL {
		t.Fatalf("payload = %#v, input = %#v", payload, usecase.createdInput)
	}
}
