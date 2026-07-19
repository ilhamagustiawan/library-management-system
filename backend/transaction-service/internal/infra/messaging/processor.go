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

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

var ErrInvalidMessage = errors.New("invalid message")

type AckProcessor struct{ repository repository.LoanRepository }

func NewAckProcessor(repository repository.LoanRepository) *AckProcessor {
	return &AckProcessor{repository: repository}
}

func (p *AckProcessor) Process(ctx context.Context, payload []byte) error {
	var event entity.Event[entity.BookStockUpdatedData]
	decoder := json.NewDecoder(io.LimitReader(bytes.NewReader(payload), 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&event); err != nil {
		return fmt.Errorf("%w: decode stock-updated event: %v", ErrInvalidMessage, err)
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("%w: stock-updated event contains multiple JSON values", ErrInvalidMessage)
	}
	if event.Type != entity.BookStockUpdatedEventType || event.OccurredAt.IsZero() || event.CausationID == "" || event.Data.UpdatedAt.IsZero() {
		return fmt.Errorf("%w: stock-updated event is incomplete", ErrInvalidMessage)
	}
	for _, id := range []string{event.EventID, event.CausationID, event.Data.LoanID, event.Data.BookID} {
		if _, err := uuid.Parse(id); err != nil {
			return fmt.Errorf("%w: stock-updated event contains invalid ID", ErrInvalidMessage)
		}
	}
	return p.repository.ConfirmStock(ctx, event.EventID, event.CausationID, event.Data.LoanID, event.Data.BookID, event.Data.UpdatedAt)
}

type Publisher interface {
	Publish(context.Context, repository.OutboxMessage) error
}

type OutboxWorker struct {
	repository repository.LoanRepository
	publisher  Publisher
	now        func() time.Time
}

func NewOutboxWorker(repository repository.LoanRepository, publisher Publisher) *OutboxWorker {
	return &OutboxWorker{repository: repository, publisher: publisher, now: func() time.Time { return time.Now().UTC() }}
}

func (w *OutboxWorker) PublishBatch(ctx context.Context) error {
	messages, err := w.repository.PendingOutbox(ctx, 50)
	if err != nil {
		return fmt.Errorf("load pending outbox: %w", err)
	}
	for _, message := range messages {
		if err := w.publisher.Publish(ctx, message); err != nil {
			if markErr := w.repository.MarkOutboxFailed(ctx, message.ID, err.Error(), w.now()); markErr != nil {
				return fmt.Errorf("publish outbox %s: %v; mark failed: %w", message.ID, err, markErr)
			}
			continue
		}
		if err := w.repository.MarkOutboxPublished(ctx, message.ID, w.now()); err != nil {
			return fmt.Errorf("mark outbox %s published: %w", message.ID, err)
		}
	}
	return nil
}

func (w *OutboxWorker) Run(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 500 * time.Millisecond
	}
	ticker := time.NewTicker(interval)
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
