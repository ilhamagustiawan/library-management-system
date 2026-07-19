package book

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type usecase struct {
	repository repository.BookRepository
	now        func() time.Time
	newID      func() string
}

func NewUsecase(repository repository.BookRepository) *usecase {
	return &usecase{repository: repository, now: utcNow, newID: uuid.NewString}
}

func utcNow() time.Time { return time.Now().UTC() }

func (u *usecase) Create(ctx context.Context, input CreateInput) (*entity.Book, error) {
	isbn, err := normalizeISBN(input.ISBN)
	if err != nil {
		return nil, validation("invalid ISBN")
	}
	title, author := strings.TrimSpace(input.Title), strings.TrimSpace(input.Author)
	if title == "" || len(title) > 255 || author == "" || len(author) > 255 || input.TotalCopies < 1 {
		return nil, validation("invalid book data")
	}
	description, err := normalizeDescription(input.Description)
	if err != nil {
		return nil, validation("invalid book data")
	}
	coverURL, err := normalizeCoverURL(input.CoverURL)
	if err != nil || !validPublicationYear(input.PublicationYear, u.now()) {
		return nil, validation("invalid book data")
	}
	now := u.now()
	book := &entity.Book{
		ID: u.newID(), ISBN: isbn, Title: title, Author: author,
		Description: description, CoverURL: coverURL, PublicationYear: input.PublicationYear,
		TotalCopies: input.TotalCopies, AvailableCopies: input.TotalCopies,
		CreatedAt: now, UpdatedAt: now,
	}
	if err := u.repository.Create(ctx, book); err != nil {
		if errors.Is(err, errs.ErrISBNExists) || errors.Is(err, errs.ErrConflict) {
			return nil, errs.New(http.StatusConflict, errs.CodeISBNExists, "ISBN is already registered", nil, err)
		}
		return nil, fmt.Errorf("create book: %w", err)
	}
	return book, nil
}

func (u *usecase) Get(ctx context.Context, id string) (*entity.Book, error) {
	if !validID(id) {
		return nil, validation("invalid book ID")
	}
	book, err := u.repository.FindByID(ctx, id)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, bookNotFound(err)
	}
	if err != nil {
		return nil, fmt.Errorf("find book: %w", err)
	}
	return book, nil
}

func (u *usecase) List(ctx context.Context, input ListInput) (Page, error) {
	input.Query = strings.TrimSpace(input.Query)
	if len(input.Query) > 200 {
		return Page{}, validation("book search query is too long")
	}
	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}
	if input.Page < 1 || input.PageSize < 1 || input.PageSize > 100 {
		return Page{}, validation("invalid pagination")
	}
	if input.SortBy == "" {
		input.SortBy = "title"
	}
	if input.SortOrder == "" {
		input.SortOrder = "asc"
	}
	if !allowedSort(input.SortBy, input.SortOrder) {
		return Page{}, validation("invalid book sort")
	}
	result, err := u.repository.List(ctx, input)
	if err != nil {
		return Page{}, fmt.Errorf("list books: %w", err)
	}
	return Page{
		Items: result.Items, Page: input.Page, PageSize: input.PageSize, TotalItems: result.TotalItems,
		TotalPages: int(math.Ceil(float64(result.TotalItems) / float64(input.PageSize))),
	}, nil
}

func (u *usecase) Update(ctx context.Context, id string, input UpdateInput) (*entity.Book, error) {
	if !validID(id) {
		return nil, validation("invalid book ID")
	}
	if !hasChanges(input) {
		return nil, validation("book update must contain at least one field")
	}
	if err := u.normalizePatch(&input); err != nil {
		return nil, err
	}
	book, err := u.repository.Update(ctx, id, input, u.now())
	if errors.Is(err, errs.ErrNotFound) {
		return nil, bookNotFound(err)
	}
	if errors.Is(err, errs.ErrConflict) {
		return nil, errs.New(http.StatusConflict, errs.CodeInventoryConflict, "total copies cannot be lower than reserved copies", nil, err)
	}
	if errors.Is(err, errs.ErrISBNExists) {
		return nil, errs.New(http.StatusConflict, errs.CodeISBNExists, "ISBN is already registered", nil, err)
	}
	if err != nil {
		return nil, fmt.Errorf("update book: %w", err)
	}
	return book, nil
}

func (u *usecase) Archive(ctx context.Context, id string) error {
	if !validID(id) {
		return validation("invalid book ID")
	}
	err := u.repository.Archive(ctx, id, u.now())
	if errors.Is(err, errs.ErrNotFound) {
		return bookNotFound(err)
	}
	if errors.Is(err, errs.ErrConflict) {
		return errs.New(http.StatusConflict, errs.CodeInventoryConflict, "book has active reservations and remains available", nil, err)
	}
	if err != nil {
		return fmt.Errorf("archive book: %w", err)
	}
	return nil
}

func (u *usecase) Stock(ctx context.Context, id string) (*entity.Stock, error) {
	if !validID(id) {
		return nil, validation("invalid book ID")
	}
	stock, err := u.repository.Stock(ctx, id)
	if errors.Is(err, errs.ErrNotFound) {
		return nil, bookNotFound(err)
	}
	if err != nil {
		return nil, fmt.Errorf("read book stock: %w", err)
	}
	return stock, nil
}

func (u *usecase) Reserve(ctx context.Context, bookID, transactionID string) (*entity.Reservation, bool, error) {
	if !validID(bookID) || !validID(transactionID) {
		return nil, false, validation("invalid book or transaction ID")
	}
	reservation, created, err := u.repository.Reserve(ctx, bookID, transactionID, u.now())
	if errors.Is(err, errs.ErrNotFound) {
		return nil, false, bookNotFound(err)
	}
	if errors.Is(err, errs.ErrStockUnavailable) {
		return nil, false, errs.New(http.StatusConflict, errs.CodeStockUnavailable, "book has no available copies", nil, err)
	}
	if errors.Is(err, errs.ErrConflict) {
		return nil, false, errs.New(http.StatusConflict, errs.CodeReservationConflict, "transaction ID cannot reserve this book", nil, err)
	}
	if err != nil {
		return nil, false, fmt.Errorf("reserve book stock: %w", err)
	}
	return reservation, created, nil
}

func (u *usecase) Release(ctx context.Context, bookID, transactionID string) error {
	if !validID(bookID) || !validID(transactionID) {
		return validation("invalid book or transaction ID")
	}
	err := u.repository.Release(ctx, bookID, transactionID, u.now())
	if errors.Is(err, errs.ErrNotFound) {
		return errs.New(http.StatusNotFound, errs.CodeReservationMissing, "stock reservation was not found", nil, err)
	}
	if errors.Is(err, errs.ErrConflict) {
		return errs.New(http.StatusConflict, errs.CodeReservationConflict, "transaction ID does not belong to this book", nil, err)
	}
	if err != nil {
		return fmt.Errorf("release book stock: %w", err)
	}
	return nil
}

func (u *usecase) normalizePatch(input *UpdateInput) error {
	for _, field := range []*entity.Change[string]{&input.ISBN, &input.Title, &input.Author} {
		if field.Set && field.Value == nil {
			return validation("required book fields cannot be null")
		}
	}
	if input.ISBN.Set {
		isbn, err := normalizeISBN(*input.ISBN.Value)
		if err != nil {
			return validation("invalid ISBN")
		}
		input.ISBN.Value = &isbn
	}
	for _, field := range []*entity.Change[string]{&input.Title, &input.Author} {
		if field.Set {
			value := strings.TrimSpace(*field.Value)
			if value == "" || len(value) > 255 {
				return validation("invalid book data")
			}
			field.Value = &value
		}
	}
	if input.Description.Set {
		description, err := normalizeDescription(input.Description.Value)
		if err != nil {
			return validation("invalid book data")
		}
		input.Description.Value = description
	}
	if input.CoverURL.Set {
		coverURL, err := normalizeCoverURL(input.CoverURL.Value)
		if err != nil {
			return validation("invalid cover URL")
		}
		input.CoverURL.Value = coverURL
	}
	if input.PublicationYear.Set && !validPublicationYear(input.PublicationYear.Value, u.now()) {
		return validation("invalid publication year")
	}
	if input.TotalCopies.Set && (input.TotalCopies.Value == nil || *input.TotalCopies.Value < 0) {
		return validation("total copies must be zero or greater")
	}
	return nil
}

func normalizeISBN(raw string) (string, error) {
	value := strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) || r == '-' {
			return -1
		}
		return unicode.ToUpper(r)
	}, strings.TrimSpace(raw))
	if len(value) == 10 && validISBN10(value) || len(value) == 13 && validISBN13(value) {
		return value, nil
	}
	return "", errors.New("invalid ISBN")
}

func validISBN10(value string) bool {
	sum := 0
	for index, r := range value {
		digit := int(r - '0')
		if index == 9 && r == 'X' {
			digit = 10
		} else if r < '0' || r > '9' {
			return false
		}
		sum += (10 - index) * digit
	}
	return sum%11 == 0
}

func validISBN13(value string) bool {
	sum := 0
	for index, r := range value {
		if r < '0' || r > '9' {
			return false
		}
		weight := 1
		if index%2 == 1 {
			weight = 3
		}
		sum += int(r-'0') * weight
	}
	return sum%10 == 0
}

func normalizeDescription(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	normalized := strings.TrimSpace(*value)
	if len(normalized) > 5000 {
		return nil, errors.New("description too long")
	}
	if normalized == "" {
		return nil, nil
	}
	return &normalized, nil
}

func normalizeCoverURL(value *string) (*string, error) {
	if value == nil {
		return nil, nil
	}
	normalized := strings.TrimSpace(*value)
	if normalized == "" {
		return nil, nil
	}
	parsed, err := url.ParseRequestURI(normalized)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || len(normalized) > 512 {
		return nil, errors.New("invalid cover URL")
	}
	return &normalized, nil
}

func validPublicationYear(value *int, now time.Time) bool {
	return value == nil || *value >= 1000 && *value <= now.Year()+1
}

func validID(value string) bool {
	_, err := uuid.Parse(value)
	return err == nil
}

func allowedSort(sortBy, sortOrder string) bool {
	validField := sortBy == "title" || sortBy == "author" || sortBy == "createdAt"
	return validField && (sortOrder == "asc" || sortOrder == "desc")
}

func hasChanges(input UpdateInput) bool {
	return input.ISBN.Set || input.Title.Set || input.Author.Set || input.Description.Set || input.CoverURL.Set ||
		input.PublicationYear.Set || input.TotalCopies.Set
}

func validation(message string) error {
	return errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, message, nil, nil)
}

func bookNotFound(cause error) error {
	return errs.New(http.StatusNotFound, errs.CodeBookNotFound, "book was not found", nil, cause)
}
