package transaction

import (
	"context"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
)

type BorrowInput struct {
	MemberID string
	BookID   string
}

type ReturnInput struct {
	LoanID                  string
	MemberID                string
	AllowAnyMember          bool
	AcceptedFineAmountMinor *int64
}

type ReturnQuoteInput struct {
	LoanID   string
	MemberID string
}

type ListInput struct {
	MemberID   string
	AllMembers bool
	Page       int
	PageSize   int
}

type Page struct {
	Items      []*entity.Transaction `json:"items"`
	Page       int                   `json:"page"`
	PageSize   int                   `json:"pageSize"`
	TotalItems int                   `json:"totalItems"`
	TotalPages int                   `json:"totalPages"`
}

type Usecase interface {
	Borrow(context.Context, BorrowInput) (*entity.Loan, error)
	QuoteReturn(context.Context, ReturnQuoteInput) (*entity.ReturnQuote, error)
	Return(context.Context, ReturnInput) (*entity.Loan, bool, error)
	List(context.Context, ListInput) (Page, error)
}
