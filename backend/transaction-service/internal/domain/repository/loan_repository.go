package repository

import (
	"context"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
)

type ReturnCommand struct {
	LoanID         string
	MemberID       string
	AllowAnyMember bool
	EventID        string
	TransactionID  string
	FineID         string
	ReturnedAt     time.Time
	DailyFineMinor int64
}

type PageFilter struct {
	MemberID string
	Page     int
	PageSize int
}

type TransactionPage struct {
	Items      []*entity.Transaction
	TotalItems int
}

type OutboxMessage struct {
	ID         string
	EventType  string
	RoutingKey string
	Payload    []byte
	Attempts   int
}

type LoanRepository interface {
	CreatePending(context.Context, *entity.Loan) error
	Activate(context.Context, string, string, time.Time) (*entity.Loan, error)
	CancelPending(context.Context, string, time.Time) error
	Return(context.Context, ReturnCommand) (*entity.Loan, bool, error)
	ListTransactions(context.Context, PageFilter) (TransactionPage, error)
	StockSyncStatus(context.Context, string) (entity.StockSyncStatus, error)
	ConfirmStock(context.Context, string, string, string, string, time.Time) error
	PendingOutbox(context.Context, int) ([]OutboxMessage, error)
	MarkOutboxPublished(context.Context, string, time.Time) error
	MarkOutboxFailed(context.Context, string, string, time.Time) error
}

type BookStock interface {
	Reserve(context.Context, string, string) error
	Release(context.Context, string, string) error
}
