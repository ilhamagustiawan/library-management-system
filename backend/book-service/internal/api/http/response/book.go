package response

import (
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	bookusecase "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/usecase/book"
)

type Book struct {
	ID              string    `json:"id" example:"f81d4fae-7dec-41d0-a765-00a0c91e6bf6"`
	ISBN            string    `json:"isbn" example:"9780132350884"`
	Title           string    `json:"title" example:"Clean Code"`
	Author          string    `json:"author" example:"Robert C. Martin"`
	Description     *string   `json:"description" example:"A handbook of agile software craftsmanship."`
	CoverURL        *string   `json:"coverUrl" example:"https://covers.openlibrary.org/b/id/10521270-M.jpg"`
	PublicationYear *int      `json:"publicationYear" example:"2008"`
	TotalCopies     int       `json:"totalCopies" example:"3"`
	AvailableCopies int       `json:"availableCopies" example:"2"`
	CreatedAt       time.Time `json:"createdAt" example:"2026-07-19T08:00:00Z"`
	UpdatedAt       time.Time `json:"updatedAt" example:"2026-07-19T08:00:00Z"`
}

type Pagination struct {
	Page       int `json:"page" example:"1"`
	PageSize   int `json:"pageSize" example:"20"`
	TotalItems int `json:"totalItems" example:"42"`
	TotalPages int `json:"totalPages" example:"3"`
}

type BookPage struct {
	Items      []Book     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type BookSuccess struct {
	Code string `json:"code" example:"LMS-200000"`
	Data Book   `json:"data"`
}

type BookPageSuccess struct {
	Code string   `json:"code" example:"LMS-200000"`
	Data BookPage `json:"data"`
}

type StockSuccess struct {
	Code string       `json:"code" example:"LMS-200000"`
	Data entity.Stock `json:"data"`
}

type ReservationSuccess struct {
	Code string             `json:"code" example:"LMS-200000"`
	Data entity.Reservation `json:"data"`
}

func NewBook(value *entity.Book) Book {
	return Book{
		ID: value.ID, ISBN: value.ISBN, Title: value.Title, Author: value.Author,
		Description: value.Description, CoverURL: value.CoverURL, PublicationYear: value.PublicationYear,
		TotalCopies: value.TotalCopies, AvailableCopies: value.AvailableCopies,
		CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt,
	}
}

func NewBookPage(value bookusecase.Page) BookPage {
	items := make([]Book, 0, len(value.Items))
	for _, book := range value.Items {
		items = append(items, NewBook(book))
	}
	return BookPage{Items: items, Pagination: Pagination{
		Page: value.Page, PageSize: value.PageSize, TotalItems: value.TotalItems, TotalPages: value.TotalPages,
	}}
}
