package repository

import (
	"context"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
)

type ListBooksFilter struct {
	Query         string
	AvailableOnly bool
	Page          int
	PageSize      int
	SortBy        string
	SortOrder     string
}

type BookPage struct {
	Items      []*entity.Book
	TotalItems int
}

type BookRepository interface {
	Create(context.Context, *entity.Book) error
	FindByID(context.Context, string) (*entity.Book, error)
	List(context.Context, ListBooksFilter) (BookPage, error)
	Update(context.Context, string, entity.BookPatch, time.Time) (*entity.Book, error)
	Archive(context.Context, string, time.Time) error
	Stock(context.Context, string) (*entity.Stock, error)
	Reserve(context.Context, string, string, time.Time) (*entity.Reservation, bool, error)
	Release(context.Context, string, string, time.Time) error
}
