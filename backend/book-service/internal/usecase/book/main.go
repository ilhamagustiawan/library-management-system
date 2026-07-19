package book

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type CreateInput struct {
	ISBN            string
	Title           string
	Author          string
	Description     *string
	CoverURL        *string
	PublicationYear *int
	TotalCopies     int
}

type UpdateInput = entity.BookPatch
type ListInput = repository.ListBooksFilter

type Page struct {
	Items      []*entity.Book
	Page       int
	PageSize   int
	TotalItems int
	TotalPages int
}

type Usecase interface {
	Create(context.Context, CreateInput) (*entity.Book, error)
	Get(context.Context, string) (*entity.Book, error)
	List(context.Context, ListInput) (Page, error)
	Update(context.Context, string, UpdateInput) (*entity.Book, error)
	Archive(context.Context, string) error
	Stock(context.Context, string) (*entity.Stock, error)
	Reserve(context.Context, string, string) (*entity.Reservation, bool, error)
	Release(context.Context, string, string) error
}
