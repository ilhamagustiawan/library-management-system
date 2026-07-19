package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type Repository struct {
	database *sqlx.DB
}

func NewRepository(database *sqlx.DB) *Repository { return &Repository{database: database} }

func (r *Repository) Create(ctx context.Context, book *entity.Book) error {
	const query = `
		INSERT INTO books (
			id, isbn, title, author, description, cover_url, publication_year,
			total_copies, available_copies, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.database.ExecContext(ctx, query,
		book.ID, book.ISBN, book.Title, book.Author, book.Description, book.CoverURL, book.PublicationYear,
		book.TotalCopies, book.AvailableCopies, book.CreatedAt, book.UpdatedAt,
	)
	if duplicateKey(err) {
		return errs.ErrISBNExists
	}
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*entity.Book, error) {
	const query = `
		SELECT id, isbn, title, author, description, cover_url, publication_year,
		       total_copies, available_copies, created_at, updated_at, archived_at
		FROM books
		WHERE id = ? AND archived_at IS NULL`
	return getBook(ctx, r.database, query, id)
}

func (r *Repository) List(ctx context.Context, filter repository.ListBooksFilter) (repository.BookPage, error) {
	where := "archived_at IS NULL"
	arguments := make([]any, 0, 8)
	if filter.Query != "" {
		where += " AND (title LIKE CONCAT('%', ?, '%') OR author LIKE CONCAT('%', ?, '%') OR isbn LIKE CONCAT('%', ?, '%'))"
		arguments = append(arguments, filter.Query, filter.Query, filter.Query)
	}
	if filter.AvailableOnly {
		where += " AND available_copies > 0"
	}

	var total int
	if err := r.database.GetContext(ctx, &total, "SELECT COUNT(*) FROM books WHERE "+where, arguments...); err != nil {
		return repository.BookPage{}, err
	}
	listArguments := append([]any(nil), arguments...)
	listArguments = append(listArguments, filter.PageSize, (filter.Page-1)*filter.PageSize)
	query := `
		SELECT id, isbn, title, author, description, cover_url, publication_year,
		       total_copies, available_copies, created_at, updated_at, archived_at
		FROM books WHERE ` + where + ` ORDER BY ` + orderClause(filter.SortBy, filter.SortOrder) + ` LIMIT ? OFFSET ?`
	var books []*entity.Book
	if err := r.database.SelectContext(ctx, &books, query, listArguments...); err != nil {
		return repository.BookPage{}, err
	}
	if books == nil {
		books = make([]*entity.Book, 0)
	}
	return repository.BookPage{Items: books, TotalItems: total}, nil
}

func (r *Repository) Update(
	ctx context.Context,
	id string,
	patch entity.BookPatch,
	now time.Time,
) (*entity.Book, error) {
	transaction, err := r.database.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	book, err := getBook(ctx, transaction, `
		SELECT id, isbn, title, author, description, cover_url, publication_year,
		       total_copies, available_copies, created_at, updated_at, archived_at
		FROM books WHERE id = ? AND archived_at IS NULL FOR UPDATE`, id)
	if err != nil {
		return nil, err
	}
	applyPatch(book, patch)
	checkedOut := book.TotalCopies - book.AvailableCopies
	if patch.TotalCopies.Set {
		if *patch.TotalCopies.Value < checkedOut {
			return nil, errs.ErrConflict
		}
		book.TotalCopies = *patch.TotalCopies.Value
		book.AvailableCopies = book.TotalCopies - checkedOut
	}
	book.UpdatedAt = now
	const query = `
		UPDATE books
		SET isbn = ?, title = ?, author = ?, description = ?, cover_url = ?, publication_year = ?,
		    total_copies = ?, available_copies = ?, updated_at = ?
		WHERE id = ?`
	_, err = transaction.ExecContext(ctx, query,
		book.ISBN, book.Title, book.Author, book.Description, book.CoverURL, book.PublicationYear,
		book.TotalCopies, book.AvailableCopies, book.UpdatedAt, book.ID,
	)
	if duplicateKey(err) {
		return nil, errs.ErrISBNExists
	}
	if err != nil {
		return nil, err
	}
	if err := transaction.Commit(); err != nil {
		return nil, err
	}
	return book, nil
}

func (r *Repository) Archive(ctx context.Context, id string, now time.Time) error {
	// TODO: Add restore/purge operations after catalog retention policy is defined.
	transaction, err := r.database.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	if _, err := getBook(ctx, transaction, `
		SELECT id, isbn, title, author, description, cover_url, publication_year,
		       total_copies, available_copies, created_at, updated_at, archived_at
		FROM books WHERE id = ? AND archived_at IS NULL FOR UPDATE`, id); err != nil {
		return err
	}
	var active int
	if err := transaction.GetContext(ctx, &active, `
		SELECT COUNT(*) FROM book_reservations WHERE book_id = ? AND status = 'active'`, id); err != nil {
		return err
	}
	if active > 0 {
		return errs.ErrConflict
	}
	if _, err := transaction.ExecContext(ctx, `UPDATE books SET archived_at = ?, updated_at = ? WHERE id = ?`, now, now, id); err != nil {
		return err
	}
	return transaction.Commit()
}

func (r *Repository) Stock(ctx context.Context, id string) (*entity.Stock, error) {
	var stock entity.Stock
	err := r.database.GetContext(ctx, &stock, `
		SELECT id AS book_id, total_copies, available_copies
		FROM books WHERE id = ? AND archived_at IS NULL`, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &stock, nil
}

func (r *Repository) Reserve(
	ctx context.Context,
	bookID string,
	transactionID string,
	now time.Time,
) (*entity.Reservation, bool, error) {
	for attempt := 0; attempt < 3; attempt++ {
		reservation, created, retry, err := r.reserveOnce(ctx, bookID, transactionID, now)
		if !retry && !retryableTransaction(err) {
			return reservation, created, err
		}
	}
	return nil, false, fmt.Errorf("reserve transaction %s remained concurrent after three attempts", transactionID)
}

func (r *Repository) reserveOnce(
	ctx context.Context,
	bookID string,
	transactionID string,
	now time.Time,
) (*entity.Reservation, bool, bool, error) {
	transaction, err := r.database.BeginTxx(ctx, nil)
	if err != nil {
		return nil, false, false, err
	}
	defer transaction.Rollback()

	var available int
	err = transaction.GetContext(ctx, &available, `
		SELECT available_copies FROM books
		WHERE id = ? AND archived_at IS NULL FOR UPDATE`, bookID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, false, false, errs.ErrNotFound
	}
	if err != nil {
		return nil, false, false, err
	}

	reservation, err := getReservation(ctx, transaction, transactionID, true)
	if err == nil {
		if reservation.BookID != bookID || reservation.Status != entity.ReservationActive {
			return nil, false, false, errs.ErrConflict
		}
		if err := transaction.Commit(); err != nil {
			return nil, false, false, err
		}
		return reservation, false, false, nil
	}
	if !errors.Is(err, errs.ErrNotFound) {
		return nil, false, false, err
	}

	if available < 1 {
		return nil, false, false, errs.ErrStockUnavailable
	}
	reservation = &entity.Reservation{
		TransactionID: transactionID, BookID: bookID, Status: entity.ReservationActive, CreatedAt: now,
	}
	_, err = transaction.ExecContext(ctx, `
		INSERT INTO book_reservations (transaction_id, book_id, status, created_at)
		VALUES (?, ?, 'active', ?)`, transactionID, bookID, now)
	if duplicateKey(err) {
		return nil, false, true, nil
	}
	if err != nil {
		return nil, false, false, err
	}
	result, err := transaction.ExecContext(ctx, `
		UPDATE books SET available_copies = available_copies - 1, updated_at = ?
		WHERE id = ? AND available_copies > 0`, now, bookID)
	if err != nil {
		return nil, false, false, err
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return nil, false, false, fmt.Errorf("decrement book stock: rows affected %d: %w", affected, err)
	}
	if err := transaction.Commit(); err != nil {
		return nil, false, false, err
	}
	return reservation, true, false, nil
}

func (r *Repository) Release(ctx context.Context, bookID, transactionID string, now time.Time) error {
	// TODO: Prune released reservations only after the transaction retry window is defined.
	transaction, err := r.database.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	var counts struct {
		Total     int `db:"total_copies"`
		Available int `db:"available_copies"`
	}
	if err := transaction.GetContext(ctx, &counts, `
		SELECT total_copies, available_copies FROM books WHERE id = ? FOR UPDATE`, bookID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errs.ErrNotFound
		}
		return err
	}
	reservation, err := getReservation(ctx, transaction, transactionID, true)
	if err != nil {
		return err
	}
	if reservation.BookID != bookID {
		return errs.ErrConflict
	}
	if reservation.Status == entity.ReservationReleased {
		return transaction.Commit()
	}
	if counts.Available >= counts.Total {
		return fmt.Errorf("release book stock: stored copy counts are inconsistent")
	}
	if _, err := transaction.ExecContext(ctx, `
		UPDATE books SET available_copies = available_copies + 1, updated_at = ? WHERE id = ?`, now, bookID); err != nil {
		return err
	}
	if _, err := transaction.ExecContext(ctx, `
		UPDATE book_reservations SET status = 'released', released_at = ? WHERE transaction_id = ?`, now, transactionID); err != nil {
		return err
	}
	return transaction.Commit()
}

type sqlxGetter interface {
	GetContext(context.Context, any, string, ...any) error
}

func getBook(ctx context.Context, database sqlxGetter, query string, arguments ...any) (*entity.Book, error) {
	var book entity.Book
	if err := database.GetContext(ctx, &book, query, arguments...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &book, nil
}

func getReservation(
	ctx context.Context,
	database sqlxGetter,
	transactionID string,
	lock bool,
) (*entity.Reservation, error) {
	query := `
		SELECT transaction_id, book_id, status, created_at, released_at
		FROM book_reservations WHERE transaction_id = ?`
	if lock {
		query += " FOR UPDATE"
	}
	var reservation entity.Reservation
	if err := database.GetContext(ctx, &reservation, query, transactionID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, err
	}
	return &reservation, nil
}

func applyPatch(book *entity.Book, patch entity.BookPatch) {
	if patch.ISBN.Set {
		book.ISBN = *patch.ISBN.Value
	}
	if patch.Title.Set {
		book.Title = *patch.Title.Value
	}
	if patch.Author.Set {
		book.Author = *patch.Author.Value
	}
	if patch.Description.Set {
		book.Description = patch.Description.Value
	}
	if patch.CoverURL.Set {
		book.CoverURL = patch.CoverURL.Value
	}
	if patch.PublicationYear.Set {
		book.PublicationYear = patch.PublicationYear.Value
	}
}

func orderClause(field, order string) string {
	columns := map[string]string{"title": "title", "author": "author", "createdAt": "created_at"}
	column, ok := columns[field]
	if !ok || order != "asc" && order != "desc" {
		return "title ASC"
	}
	return column + " " + strings.ToUpper(order)
}

func duplicateKey(err error) bool {
	var mysqlError *mysql.MySQLError
	return errors.As(err, &mysqlError) && mysqlError.Number == 1062
}

func retryableTransaction(err error) bool {
	var mysqlError *mysql.MySQLError
	return errors.As(err, &mysqlError) && (mysqlError.Number == 1205 || mysqlError.Number == 1213)
}
