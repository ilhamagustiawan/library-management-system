package transaction

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

type fakeRepository struct {
	createErr error
	loan      *entity.Loan
	confirmed bool
}

func (r *fakeRepository) CreatePending(_ context.Context, loan *entity.Loan) error {
	r.loan = loan
	return r.createErr
}
func (r *fakeRepository) Activate(_ context.Context, id, _ string, _ time.Time) (*entity.Loan, error) {
	if r.loan != nil {
		r.loan.Status = entity.LoanActive
		return r.loan, nil
	}
	return &entity.Loan{ID: id, Status: entity.LoanActive}, nil
}
func (r *fakeRepository) CancelPending(context.Context, string, time.Time) error { return nil }
func (r *fakeRepository) Return(_ context.Context, command repository.ReturnCommand) (*entity.Loan, bool, error) {
	if r.loan != nil {
		return r.loan, false, nil
	}
	return &entity.Loan{ID: command.LoanID, Status: entity.LoanReturned, StockSyncStatus: entity.StockSyncPending}, false, nil
}
func (r *fakeRepository) ListTransactions(context.Context, repository.PageFilter) (repository.TransactionPage, error) {
	return repository.TransactionPage{}, nil
}
func (r *fakeRepository) StockSyncStatus(context.Context, string) (entity.StockSyncStatus, error) {
	if r.confirmed {
		return entity.StockSyncConfirmed, nil
	}
	return entity.StockSyncPending, nil
}
func (r *fakeRepository) ConfirmStock(context.Context, string, string, string, string, time.Time) error {
	return nil
}
func (r *fakeRepository) PendingOutbox(context.Context, int) ([]repository.OutboxMessage, error) {
	return nil, nil
}
func (r *fakeRepository) MarkOutboxPublished(context.Context, string, time.Time) error { return nil }
func (r *fakeRepository) MarkOutboxFailed(context.Context, string, string, time.Time) error {
	return nil
}

type fakeStock struct{ reserveErr error }

func (s fakeStock) Reserve(context.Context, string, string) error { return s.reserveErr }
func (s fakeStock) Release(context.Context, string, string) error { return nil }

func TestBorrowCreatesSevenDayLoanAfterAtomicReserve(t *testing.T) {
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	repository := &fakeRepository{}
	usecase := NewUsecase(repository, fakeStock{}, Config{Now: func() time.Time { return now }, NewID: func() string { return "loan-1" }})

	loan, err := usecase.Borrow(context.Background(), BorrowInput{MemberID: "member-1", BookID: "book-1"})
	if err != nil {
		t.Fatalf("Borrow() error = %v", err)
	}
	if loan.ID != "loan-1" || !loan.DueAt.Equal(now.Add(7*24*time.Hour)) {
		t.Fatalf("loan = %#v, want seven-day term", loan)
	}
}

func TestBorrowPreservesThreeLoanLimitConflict(t *testing.T) {
	usecase := NewUsecase(&fakeRepository{createErr: errs.ErrLoanLimit}, fakeStock{}, Config{})
	_, err := usecase.Borrow(context.Background(), BorrowInput{MemberID: "member-1", BookID: "book-1"})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.ErrorCode != errs.CodeLoanLimit {
		t.Fatalf("Borrow() error = %v, want loan-limit domain error", err)
	}
}

func TestReturnReportsPendingWhenBookAckTimesOut(t *testing.T) {
	repository := &fakeRepository{}
	usecase := NewUsecase(repository, fakeStock{}, Config{AckTimeout: time.Millisecond, PollInterval: time.Millisecond})
	loan, confirmed, err := usecase.Return(context.Background(), ReturnInput{LoanID: "loan-1", MemberID: "member-1"})
	if err != nil {
		t.Fatalf("Return() error = %v", err)
	}
	if confirmed || loan.StockSyncStatus != entity.StockSyncPending {
		t.Fatalf("Return() = confirmed %v, loan %#v", confirmed, loan)
	}
}

func TestReturnReportsConfirmedBookAck(t *testing.T) {
	repository := &fakeRepository{confirmed: true}
	usecase := NewUsecase(repository, fakeStock{}, Config{AckTimeout: time.Second, PollInterval: time.Millisecond})
	_, confirmed, err := usecase.Return(context.Background(), ReturnInput{LoanID: "loan-1", MemberID: "member-1"})
	if err != nil || !confirmed {
		t.Fatalf("Return() = confirmed %v, error %v", confirmed, err)
	}
}
