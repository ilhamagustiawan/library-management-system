package transaction

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

const (
	defaultLoanTerm     = 7 * 24 * time.Hour
	defaultAckTimeout   = 5 * time.Second
	defaultPollInterval = 100 * time.Millisecond
	defaultDailyFineIDR = int64(5000)
)

type Config struct {
	LoanTerm       time.Duration
	AckTimeout     time.Duration
	PollInterval   time.Duration
	DailyFineMinor int64
	Now            func() time.Time
	NewID          func() string
}

type usecase struct {
	repository     repository.LoanRepository
	stock          repository.BookStock
	loanTerm       time.Duration
	ackTimeout     time.Duration
	pollInterval   time.Duration
	dailyFineMinor int64
	now            func() time.Time
	newID          func() string
}

func NewUsecase(repository repository.LoanRepository, stock repository.BookStock, config Config) *usecase {
	if config.LoanTerm <= 0 {
		config.LoanTerm = defaultLoanTerm
	}
	if config.AckTimeout <= 0 {
		config.AckTimeout = defaultAckTimeout
	}
	if config.PollInterval <= 0 {
		config.PollInterval = defaultPollInterval
	}
	if config.DailyFineMinor <= 0 {
		config.DailyFineMinor = defaultDailyFineIDR
	}
	if config.Now == nil {
		config.Now = func() time.Time { return time.Now().UTC() }
	}
	if config.NewID == nil {
		config.NewID = uuid.NewString
	}
	return &usecase{
		repository: repository, stock: stock, loanTerm: config.LoanTerm,
		ackTimeout: config.AckTimeout, pollInterval: config.PollInterval,
		dailyFineMinor: config.DailyFineMinor, now: config.Now, newID: config.NewID,
	}
}

func (u *usecase) Borrow(ctx context.Context, input BorrowInput) (*entity.Loan, error) {
	memberID, bookID := strings.TrimSpace(input.MemberID), strings.TrimSpace(input.BookID)
	if memberID == "" || bookID == "" {
		return nil, validation("member and book IDs are required")
	}
	now := u.now()
	loan := &entity.Loan{
		ID: u.newID(), MemberID: memberID, BookID: bookID,
		Status: entity.LoanPendingReservation, StockSyncStatus: entity.StockSyncNotApplicable,
		BorrowedAt: now, DueAt: now.Add(u.loanTerm), CreatedAt: now, UpdatedAt: now,
	}
	if err := u.repository.CreatePending(ctx, loan); err != nil {
		return nil, mapLoanError("create pending loan", err)
	}
	if err := u.stock.Reserve(ctx, bookID, loan.ID); err != nil {
		if cancelErr := u.repository.CancelPending(ctx, loan.ID, u.now()); cancelErr != nil {
			return nil, fmt.Errorf("reserve book: %w; cancel pending loan: %v", err, cancelErr)
		}
		return nil, mapBookError(err)
	}
	activated, err := u.repository.Activate(ctx, loan.ID, u.newID(), u.now())
	if err == nil {
		return activated, nil
	}
	releaseErr := u.stock.Release(ctx, bookID, loan.ID)
	cancelErr := u.repository.CancelPending(ctx, loan.ID, u.now())
	if releaseErr != nil || cancelErr != nil {
		return nil, fmt.Errorf("activate reserved loan: %w; release compensation: %v; cancel pending loan: %v", err, releaseErr, cancelErr)
	}
	return nil, fmt.Errorf("activate reserved loan: %w; stock reservation released and pending loan cancelled", err)
}

func (u *usecase) Return(ctx context.Context, input ReturnInput) (*entity.Loan, bool, error) {
	if strings.TrimSpace(input.LoanID) == "" || strings.TrimSpace(input.MemberID) == "" {
		return nil, false, validation("loan and member IDs are required")
	}
	now := u.now()
	loan, _, err := u.repository.Return(ctx, repository.ReturnCommand{
		LoanID: input.LoanID, MemberID: input.MemberID, AllowAnyMember: input.AllowAnyMember,
		EventID: u.newID(), TransactionID: u.newID(), FineID: u.newID(), ReturnedAt: now,
		DailyFineMinor: u.dailyFineMinor,
	})
	if err != nil {
		return nil, false, mapLoanError("return loan", err)
	}
	confirmed, err := u.waitForStockAck(ctx, loan.ID)
	if err != nil {
		return nil, false, err
	}
	if confirmed {
		loan.StockSyncStatus = entity.StockSyncConfirmed
	}
	return loan, confirmed, nil
}

func (u *usecase) waitForStockAck(ctx context.Context, loanID string) (bool, error) {
	waitCtx, cancel := context.WithTimeout(ctx, u.ackTimeout)
	defer cancel()
	ticker := time.NewTicker(u.pollInterval)
	defer ticker.Stop()
	for {
		status, err := u.repository.StockSyncStatus(waitCtx, loanID)
		if err != nil {
			return false, fmt.Errorf("read stock update status: %w", err)
		}
		if status == entity.StockSyncConfirmed {
			return true, nil
		}
		select {
		case <-waitCtx.Done():
			if errors.Is(waitCtx.Err(), context.DeadlineExceeded) {
				return false, nil
			}
			return false, waitCtx.Err()
		case <-ticker.C:
		}
	}
}

func (u *usecase) List(ctx context.Context, input ListInput) (Page, error) {
	if strings.TrimSpace(input.MemberID) == "" && !input.AllMembers {
		return Page{}, validation("member ID is required")
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
	result, err := u.repository.ListTransactions(ctx, repository.PageFilter{MemberID: input.MemberID, Page: input.Page, PageSize: input.PageSize})
	if err != nil {
		return Page{}, fmt.Errorf("list transactions: %w", err)
	}
	return Page{
		Items: result.Items, Page: input.Page, PageSize: input.PageSize, TotalItems: result.TotalItems,
		TotalPages: int(math.Ceil(float64(result.TotalItems) / float64(input.PageSize))),
	}, nil
}

func mapLoanError(operation string, err error) error {
	switch {
	case errors.Is(err, errs.ErrLoanLimit):
		return errs.New(http.StatusConflict, errs.CodeLoanLimit, "member already has the maximum of three active loans", nil, err)
	case errors.Is(err, errs.ErrActiveLoan):
		return errs.New(http.StatusConflict, errs.CodeActiveLoan, "member already has an active loan for this book", nil, err)
	case errors.Is(err, errs.ErrNotFound):
		return errs.New(http.StatusNotFound, errs.CodeLoanNotFound, "loan was not found", nil, err)
	case errors.Is(err, errs.ErrForbidden):
		return errs.New(http.StatusForbidden, errs.CodeForbidden, "loan belongs to another member", nil, err)
	default:
		return fmt.Errorf("%s: %w", operation, err)
	}
}

func mapBookError(err error) error {
	switch {
	case errors.Is(err, errs.ErrBookNotFound):
		return errs.New(http.StatusNotFound, errs.CodeBookNotFound, "book was not found", nil, err)
	case errors.Is(err, errs.ErrStockUnavailable):
		return errs.New(http.StatusConflict, errs.CodeStockUnavailable, "book has no available copies", nil, err)
	case errors.Is(err, errs.ErrDependency):
		return errs.New(http.StatusServiceUnavailable, errs.CodeDependency, "Book Service is unavailable; loan was not created", nil, err)
	default:
		return fmt.Errorf("reserve book stock: %w", err)
	}
}

func validation(message string) error {
	return errs.New(http.StatusUnprocessableEntity, errs.CodeValidation, message, nil, nil)
}
