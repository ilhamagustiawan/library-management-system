package loan

import (
	"context"
	"errors"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

func TestCreatePendingAtomicallyLimitsMemberToThreeLoans(t *testing.T) {
	database := integrationDatabase(t)
	store := NewRepository(database)
	memberID := uuid.NewString()
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	loans := make([]*entity.Loan, 4)
	for index := range loans {
		loans[index] = pendingLoan(memberID, uuid.NewString(), now)
	}
	cleanupMember(t, database, memberID)

	results := make(chan error, len(loans))
	var group sync.WaitGroup
	for _, loan := range loans {
		group.Add(1)
		go func() {
			defer group.Done()
			results <- store.CreatePending(context.Background(), loan)
		}()
	}
	group.Wait()
	close(results)

	successes, limitFailures := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, errs.ErrLoanLimit):
			limitFailures++
		default:
			t.Fatalf("CreatePending() error = %v", err)
		}
	}
	if successes != 3 || limitFailures != 1 {
		t.Fatalf("successes = %d, limit failures = %d; want 3 and 1", successes, limitFailures)
	}
}

func TestReturnAssessesFineAndConfirmsMatchingBookAck(t *testing.T) {
	database := integrationDatabase(t)
	store := NewRepository(database)
	memberID, bookID := uuid.NewString(), uuid.NewString()
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	loan := pendingLoan(memberID, bookID, now)
	cleanupMember(t, database, memberID)

	if err := store.CreatePending(context.Background(), loan); err != nil {
		t.Fatalf("CreatePending() error = %v", err)
	}
	if _, err := store.Activate(context.Background(), loan.ID, uuid.NewString(), now); err != nil {
		t.Fatalf("Activate() error = %v", err)
	}
	eventID := uuid.NewString()
	returnedAt := loan.DueAt.Add(time.Second)
	returned, _, err := store.Return(context.Background(), repository.ReturnCommand{
		LoanID: loan.ID, MemberID: memberID, EventID: eventID, TransactionID: uuid.NewString(),
		FineID: uuid.NewString(), ReturnedAt: returnedAt, DailyFineMinor: 5000,
	})
	if err != nil {
		t.Fatalf("Return() error = %v", err)
	}
	if returned.Fine == nil || returned.Fine.OverdueDays != 1 || returned.Fine.TotalAmountMinor != 5000 {
		t.Fatalf("fine = %#v, want one-day IDR 5000 fine", returned.Fine)
	}
	mismatchedAckID, matchingAckID := uuid.NewString(), uuid.NewString()
	t.Cleanup(func() {
		_, _ = database.Exec("DELETE FROM inbox_events WHERE event_id IN (?, ?)", mismatchedAckID, matchingAckID)
	})
	if err := store.ConfirmStock(context.Background(), mismatchedAckID, eventID, loan.ID, uuid.NewString(), returnedAt); !errors.Is(err, errs.ErrNotFound) {
		t.Fatalf("mismatched ConfirmStock() error = %v, want not found", err)
	}
	if err := store.ConfirmStock(context.Background(), matchingAckID, eventID, loan.ID, bookID, returnedAt); err != nil {
		t.Fatalf("ConfirmStock() error = %v", err)
	}
	status, err := store.StockSyncStatus(context.Background(), loan.ID)
	if err != nil || status != entity.StockSyncConfirmed {
		t.Fatalf("stock status = %q, error = %v", status, err)
	}
}

func integrationDatabase(t *testing.T) *sqlx.DB {
	t.Helper()
	dsn := os.Getenv("TRANSACTION_TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("TRANSACTION_TEST_DATABASE_DSN not set")
	}
	database, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatalf("connect integration database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })
	return database
}

func pendingLoan(memberID, bookID string, now time.Time) *entity.Loan {
	return &entity.Loan{
		ID: uuid.NewString(), MemberID: memberID, BookID: bookID,
		Status: entity.LoanPendingReservation, StockSyncStatus: entity.StockSyncNotApplicable,
		BorrowedAt: now, DueAt: now.Add(7 * 24 * time.Hour), CreatedAt: now, UpdatedAt: now,
	}
}

func cleanupMember(t *testing.T, database *sqlx.DB, memberID string) {
	t.Helper()
	t.Cleanup(func() {
		_, _ = database.Exec(`DELETE FROM outbox_events WHERE payload ->> '$.data.memberId' = ?`, memberID)
		_, _ = database.Exec("DELETE FROM fines WHERE member_id = ?", memberID)
		_, _ = database.Exec("DELETE FROM loan_transactions WHERE member_id = ?", memberID)
		_, _ = database.Exec("DELETE FROM loans WHERE member_id = ?", memberID)
		_, _ = database.Exec("DELETE FROM member_loan_counters WHERE member_id = ?", memberID)
	})
}
