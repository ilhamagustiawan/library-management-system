package book

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	domainrepo "github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

func TestReserveIsAtomicAndRetrySafe(t *testing.T) {
	database := integrationDatabase(t)
	repository := NewRepository(database)
	bookID := uuid.NewString()
	transactionIDs := []string{uuid.NewString(), uuid.NewString()}
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	book := &entity.Book{
		ID: bookID, ISBN: testISBN(bookID), Title: "Concurrent Systems", Author: "Ada Lovelace",
		TotalCopies: 1, AvailableCopies: 1, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.Create(context.Background(), book); err != nil {
		t.Fatalf("create book: %v", err)
	}
	t.Cleanup(func() {
		_, _ = database.Exec("DELETE FROM book_reservations WHERE book_id = ?", bookID)
		_, _ = database.Exec("DELETE FROM books WHERE id = ?", bookID)
	})

	type result struct {
		transactionID string
		err           error
	}
	results := make(chan result, len(transactionIDs))
	var group sync.WaitGroup
	for _, transactionID := range transactionIDs {
		group.Add(1)
		go func() {
			defer group.Done()
			_, _, err := repository.Reserve(context.Background(), bookID, transactionID, now)
			results <- result{transactionID: transactionID, err: err}
		}()
	}
	group.Wait()
	close(results)

	winner := ""
	stockFailures := 0
	for result := range results {
		if result.err == nil {
			winner = result.transactionID
		} else if errors.Is(result.err, errs.ErrStockUnavailable) {
			stockFailures++
		} else {
			t.Fatalf("Reserve() unexpected error = %v", result.err)
		}
	}
	if winner == "" || stockFailures != 1 {
		t.Fatalf("winner = %q, stock failures = %d; want one each", winner, stockFailures)
	}
	if _, created, err := repository.Reserve(context.Background(), bookID, winner, now); err != nil || created {
		t.Fatalf("retry Reserve() created = %t, error = %v", created, err)
	}
	stock, err := repository.Stock(context.Background(), bookID)
	if err != nil || stock.AvailableCopies != 0 {
		t.Fatalf("reserved stock = %#v, error = %v", stock, err)
	}
	concurrentErrors := make(chan error, 2)
	go func() {
		_, _, err := repository.Reserve(context.Background(), bookID, winner, now.Add(time.Minute))
		concurrentErrors <- err
	}()
	go func() {
		concurrentErrors <- repository.Release(context.Background(), bookID, winner, now.Add(time.Minute))
	}()
	firstError, secondError := <-concurrentErrors, <-concurrentErrors
	for _, err := range []error{firstError, secondError} {
		if err != nil && !errors.Is(err, errs.ErrConflict) {
			t.Fatalf("concurrent reserve/release error = %v", err)
		}
	}
	if err := repository.Release(context.Background(), bookID, winner, now.Add(2*time.Minute)); err != nil {
		t.Fatalf("retry Release() error = %v", err)
	}
	stock, err = repository.Stock(context.Background(), bookID)
	if err != nil || stock.AvailableCopies != 1 {
		t.Fatalf("released stock = %#v, error = %v", stock, err)
	}
}

func TestApplyLoanReturnedRestoresStockAndQueuesAckOnce(t *testing.T) {
	database := integrationDatabase(t)
	store := NewRepository(database)
	bookID, loanID, memberID := uuid.NewString(), uuid.NewString(), uuid.NewString()
	incomingID, ackID := uuid.NewString(), uuid.NewString()
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	book := &entity.Book{
		ID: bookID, ISBN: testISBN(bookID), Title: "Reliable Messaging", Author: "Grace Hopper",
		TotalCopies: 1, AvailableCopies: 1, CreatedAt: now, UpdatedAt: now,
	}
	if err := store.Create(context.Background(), book); err != nil {
		t.Fatalf("create book: %v", err)
	}
	if _, _, err := store.Reserve(context.Background(), bookID, loanID, now); err != nil {
		t.Fatalf("reserve book: %v", err)
	}
	t.Cleanup(func() {
		_, _ = database.Exec("DELETE FROM message_outbox WHERE id = ?", ackID)
		_, _ = database.Exec("DELETE FROM message_inbox WHERE event_id = ?", incomingID)
		_, _ = database.Exec("DELETE FROM book_reservations WHERE transaction_id = ?", loanID)
		_, _ = database.Exec("DELETE FROM books WHERE id = ?", bookID)
	})

	incoming := entity.Event[entity.LoanReturnedData]{
		EventID: incomingID, Type: entity.LoanReturnedEventType, OccurredAt: now,
		Data: entity.LoanReturnedData{LoanID: loanID, BookID: bookID, MemberID: memberID, ReturnedAt: now},
	}
	ack := entity.Event[entity.BookStockUpdatedData]{
		EventID: ackID, Type: entity.BookStockUpdatedEventType, OccurredAt: now, CausationID: incomingID,
		Data: entity.BookStockUpdatedData{LoanID: loanID, BookID: bookID, UpdatedAt: now},
	}
	command := domainrepo.StockReturnCommand{Incoming: incoming, Ack: ack}
	if err := store.ApplyLoanReturned(context.Background(), command); err != nil {
		t.Fatalf("ApplyLoanReturned() error = %v", err)
	}
	if err := store.ApplyLoanReturned(context.Background(), command); err != nil {
		t.Fatalf("duplicate ApplyLoanReturned() error = %v", err)
	}
	stock, err := store.Stock(context.Background(), bookID)
	if err != nil || stock.AvailableCopies != 1 {
		t.Fatalf("stock = %#v, error = %v; want one available copy", stock, err)
	}
	var ackCount int
	if err := database.Get(&ackCount, "SELECT COUNT(*) FROM message_outbox WHERE id = ?", ackID); err != nil || ackCount != 1 {
		t.Fatalf("ack count = %d, error = %v; want 1", ackCount, err)
	}
}

func integrationDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("BOOK_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("BOOK_TEST_DATABASE_DSN not set")
	}
	database, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("connect integration database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return database
}

func testISBN(id string) string {
	// Integration uniqueness matters; checksum validation belongs to use case tests.
	return "9" + id[:12]
}
