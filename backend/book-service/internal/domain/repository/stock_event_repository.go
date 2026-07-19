package repository

import (
	"context"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
)

type StockReturnCommand struct {
	Incoming entity.Event[entity.LoanReturnedData]
	Ack      entity.Event[entity.BookStockUpdatedData]
}

type OutboxMessage struct {
	ID, EventType, RoutingKey string
	Payload                   []byte
	Attempts                  int
}

type StockEventRepository interface {
	ApplyLoanReturned(context.Context, StockReturnCommand) error
	PendingMessages(context.Context, int) ([]OutboxMessage, error)
	MarkPublished(context.Context, string, time.Time) error
	MarkFailed(context.Context, string, string, time.Time) error
}
