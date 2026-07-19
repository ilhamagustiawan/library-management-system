package messaging

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReturnProcessor struct {
	repository repository.StockEventRepository
	now        func() time.Time
	newID      func() string
}

func NewReturnProcessor(repository repository.StockEventRepository) *ReturnProcessor {
	return &ReturnProcessor{repository: repository, now: func() time.Time { return time.Now().UTC() }, newID: uuid.NewString}
}

func (p *ReturnProcessor) Process(ctx context.Context, payload []byte) error {
	var incoming entity.Event[entity.LoanReturnedData]
	decoder := json.NewDecoder(io.LimitReader(bytes.NewReader(payload), 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&incoming); err != nil {
		return fmt.Errorf("%w: decode LoanReturned event: %v", ErrInvalidMessage, err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("%w: LoanReturned event contains multiple JSON values", ErrInvalidMessage)
	}
	if incoming.Type != entity.LoanReturnedEventType || incoming.OccurredAt.IsZero() || incoming.Data.ReturnedAt.IsZero() {
		return fmt.Errorf("%w: LoanReturned event is incomplete", ErrInvalidMessage)
	}
	for _, id := range []string{incoming.EventID, incoming.Data.LoanID, incoming.Data.BookID, incoming.Data.MemberID} {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("%w: LoanReturned event contains invalid ID", ErrInvalidMessage)
		}
	}
	now := p.now()
	ack := entity.Event[entity.BookStockUpdatedData]{EventID: p.newID(), Type: entity.BookStockUpdatedEventType, OccurredAt: now, CausationID: incoming.EventID, Data: entity.BookStockUpdatedData{LoanID: incoming.Data.LoanID, BookID: incoming.Data.BookID, UpdatedAt: now}}
	return p.repository.ApplyLoanReturned(ctx, repository.StockReturnCommand{Incoming: incoming, Ack: ack})
}

type Publisher interface {
	Publish(context.Context, repository.OutboxMessage) error
}

type OutboxWorker struct {
	repository repository.StockEventRepository
	publisher  Publisher
	now        func() time.Time
}

func NewOutboxWorker(repository repository.StockEventRepository, publisher Publisher) *OutboxWorker {
	return &OutboxWorker{repository: repository, publisher: publisher, now: func() time.Time { return time.Now().UTC() }}
}

func (w *OutboxWorker) PublishBatch(ctx context.Context) error {
	messages, err := w.repository.PendingMessages(ctx, 50)
	if err != nil {
		return err
	}
	for _, message := range messages {
		if err := w.publisher.Publish(ctx, message); err != nil {
			if markErr := w.repository.MarkFailed(ctx, message.ID, err.Error(), w.now()); markErr != nil {
				return markErr
			}
			continue
		}
		if err := w.repository.MarkPublished(ctx, message.ID, w.now()); err != nil {
			return err
		}
	}
	return nil
}

func (w *OutboxWorker) Run(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	for {
		_ = w.PublishBatch(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
