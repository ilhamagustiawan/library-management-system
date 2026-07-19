package entity

import (
	"math"
	"time"
)

type LoanStatus string

const (
	LoanPendingReservation LoanStatus = "pending_reservation"
	LoanActive             LoanStatus = "active"
	LoanReturned           LoanStatus = "returned"
	LoanCancelled          LoanStatus = "cancelled"
)

type StockSyncStatus string

const (
	StockSyncNotApplicable StockSyncStatus = "not_applicable"
	StockSyncPending       StockSyncStatus = "pending"
	StockSyncConfirmed     StockSyncStatus = "confirmed"
)

type Loan struct {
	ID              string          `json:"id" db:"id"`
	MemberID        string          `json:"memberId" db:"member_id"`
	BookID          string          `json:"bookId" db:"book_id"`
	Status          LoanStatus      `json:"status" db:"status"`
	StockSyncStatus StockSyncStatus `json:"stockSyncStatus" db:"stock_sync_status"`
	BorrowedAt      time.Time       `json:"borrowedAt" db:"borrowed_at"`
	DueAt           time.Time       `json:"dueAt" db:"due_at"`
	ReturnedAt      *time.Time      `json:"returnedAt,omitempty" db:"returned_at"`
	CreatedAt       time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" db:"updated_at"`
	Fine            *Fine           `json:"fine,omitempty"`
}

type FineStatus string

const FineUnpaid FineStatus = "unpaid"

type Fine struct {
	ID               string     `json:"id" db:"fine_id"`
	LoanID           string     `json:"loanId" db:"loan_id"`
	MemberID         string     `json:"memberId" db:"fine_member_id"`
	OverdueDays      int        `json:"overdueDays" db:"overdue_days"`
	DailyRateMinor   int64      `json:"dailyRateMinor" db:"daily_rate_minor"`
	TotalAmountMinor int64      `json:"totalAmountMinor" db:"total_amount_minor"`
	Currency         string     `json:"currency" db:"currency"`
	Status           FineStatus `json:"status" db:"fine_status"`
	AssessedAt       time.Time  `json:"assessedAt" db:"assessed_at"`
}

type FineQuote struct {
	OverdueDays      int    `json:"overdueDays"`
	DailyRateMinor   int64  `json:"dailyRateMinor"`
	TotalAmountMinor int64  `json:"totalAmountMinor"`
	Currency         string `json:"currency"`
}

type ReturnQuote struct {
	LoanID   string     `json:"loanId"`
	BookID   string     `json:"bookId"`
	DueAt    time.Time  `json:"dueAt"`
	QuotedAt time.Time  `json:"quotedAt"`
	Fine     *FineQuote `json:"fine"`
}

func QuoteFine(dueAt, quotedAt time.Time, dailyRateMinor int64) *FineQuote {
	if !quotedAt.After(dueAt) {
		return nil
	}
	overdueDays := int(math.Ceil(quotedAt.Sub(dueAt).Hours() / 24))
	return &FineQuote{
		OverdueDays: overdueDays, DailyRateMinor: dailyRateMinor,
		TotalAmountMinor: int64(overdueDays) * dailyRateMinor, Currency: "IDR",
	}
}

type TransactionType string

const (
	TransactionBorrow TransactionType = "borrow"
	TransactionReturn TransactionType = "return"
)

type Transaction struct {
	ID         string          `json:"id" db:"id"`
	LoanID     string          `json:"loanId" db:"loan_id"`
	MemberID   string          `json:"memberId" db:"member_id"`
	BookID     string          `json:"bookId" db:"book_id"`
	Type       TransactionType `json:"type" db:"type"`
	OccurredAt time.Time       `json:"occurredAt" db:"occurred_at"`
	Fine       *Fine           `json:"fine,omitempty"`
}
