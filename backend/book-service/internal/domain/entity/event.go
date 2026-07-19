package entity

import "time"

const (
	LoanReturnedEventType     = "LoanReturned.v1"
	BookStockUpdatedEventType = "BookStockUpdated.v1"
)

type Event[T any] struct {
	EventID     string    `json:"eventId"`
	Type        string    `json:"type"`
	OccurredAt  time.Time `json:"occurredAt"`
	CausationID string    `json:"causationId,omitempty"`
	Data        T         `json:"data"`
}

type LoanReturnedData struct {
	LoanID     string    `json:"loanId"`
	BookID     string    `json:"bookId"`
	MemberID   string    `json:"memberId"`
	ReturnedAt time.Time `json:"returnedAt"`
}

type BookStockUpdatedData struct {
	LoanID    string    `json:"loanId"`
	BookID    string    `json:"bookId"`
	UpdatedAt time.Time `json:"updatedAt"`
}
