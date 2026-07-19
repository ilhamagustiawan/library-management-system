package book

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type fakeRepository struct {
	created *entity.Book
	page    repository.BookPage
	err     error
}

func (r *fakeRepository) Create(_ context.Context, value *entity.Book) error {
	r.created = value
	return r.err
}
func (r *fakeRepository) FindByID(context.Context, string) (*entity.Book, error) { return nil, r.err }
func (r *fakeRepository) List(context.Context, repository.ListBooksFilter) (repository.BookPage, error) {
	return r.page, r.err
}
func (r *fakeRepository) Update(context.Context, string, entity.BookPatch, time.Time) (*entity.Book, error) {
	return nil, r.err
}
func (r *fakeRepository) Archive(context.Context, string, time.Time) error     { return r.err }
func (r *fakeRepository) Stock(context.Context, string) (*entity.Stock, error) { return nil, r.err }
func (r *fakeRepository) Reserve(context.Context, string, string, time.Time) (*entity.Reservation, bool, error) {
	return nil, false, r.err
}
func (r *fakeRepository) Release(context.Context, string, string, time.Time) error { return r.err }

func TestNormalizeISBN(t *testing.T) {
	for _, test := range []struct{ raw, want string }{
		{raw: "978-0-13-235088-4", want: "9780132350884"},
		{raw: "0-13-235088-2", want: "0132350882"},
	} {
		got, err := normalizeISBN(test.raw)
		if err != nil || got != test.want {
			t.Fatalf("normalizeISBN(%q) = %q, %v; want %q", test.raw, got, err, test.want)
		}
	}
}

func TestNormalizeISBNRejectsInvalidChecksum(t *testing.T) {
	if _, err := normalizeISBN("9780132350885"); err == nil {
		t.Fatal("normalizeISBN() error = nil, want invalid checksum")
	}
}

func TestCreateNormalizesBookAndInitializesAvailability(t *testing.T) {
	repo := &fakeRepository{}
	usecase := NewUsecase(repo)
	usecase.now = func() time.Time { return time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC) }
	usecase.newID = func() string { return "book-id" }

	created, err := usecase.Create(context.Background(), CreateInput{
		ISBN: "978-0-13-235088-4", Title: "  Clean Code ", Author: " Robert C. Martin ",
		CoverURL: stringPointer("  https://covers.openlibrary.org/b/id/10521270-M.jpg  "), TotalCopies: 3,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ISBN != "9780132350884" || created.Title != "Clean Code" || created.Author != "Robert C. Martin" {
		t.Fatalf("created book = %#v", created)
	}
	if created.CoverURL == nil || *created.CoverURL != "https://covers.openlibrary.org/b/id/10521270-M.jpg" {
		t.Fatalf("created cover URL = %v", created.CoverURL)
	}
	if created.TotalCopies != 3 || created.AvailableCopies != 3 || repo.created != created {
		t.Fatalf("created stock = %#v", created)
	}
}

func TestCreateRejectsInsecureCoverURL(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{})
	_, err := usecase.Create(context.Background(), CreateInput{
		ISBN: "9780132350884", Title: "Clean Code", Author: "Robert C. Martin",
		CoverURL: stringPointer("http://covers.openlibrary.org/b/id/10521270-M.jpg"), TotalCopies: 1,
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeValidation {
		t.Fatalf("Create() error = %v, want cover URL validation error", err)
	}
}

func stringPointer(value string) *string { return &value }

func TestCreateMapsDuplicateISBN(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{err: errs.ErrConflict})
	_, err := usecase.Create(context.Background(), CreateInput{
		ISBN: "9780132350884", Title: "Clean Code", Author: "Robert C. Martin", TotalCopies: 1,
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeISBNExists {
		t.Fatalf("Create() error = %v, want ISBN conflict", err)
	}
}

func TestListAppliesSafeDefaults(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{page: repository.BookPage{TotalItems: 41}})
	page, err := usecase.List(context.Background(), ListInput{})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if page.Page != 1 || page.PageSize != 20 || page.TotalPages != 3 {
		t.Fatalf("page = %#v", page)
	}
}
