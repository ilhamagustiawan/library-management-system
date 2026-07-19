package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

type fakeRepository struct{ applied repository.StockReturnCommand }

func (r *fakeRepository) ApplyLoanReturned(_ context.Context, command repository.StockReturnCommand) error {
	r.applied = command
	return nil
}
func (r *fakeRepository) PendingMessages(context.Context, int) ([]repository.OutboxMessage, error) {
	return nil, nil
}
func (r *fakeRepository) MarkPublished(context.Context, string, time.Time) error      { return nil }
func (r *fakeRepository) MarkFailed(context.Context, string, string, time.Time) error { return nil }

func TestReturnProcessorCreatesCorrelatedStockAck(t *testing.T) {
	repository := &fakeRepository{}
	processor := NewReturnProcessor(repository)
	payload := []byte(`{"eventId":"dd8cd46e-41e1-4583-aab2-c0342884201e","type":"LoanReturned.v1","occurredAt":"2026-07-19T10:00:00Z","data":{"loanId":"52a88672-a4c2-4876-be5a-65863aeb35e4","bookId":"7b36fe43-f31d-4861-884f-42ed7386b1e9","memberId":"31c73b2e-0640-49bd-8f06-3bb7272921fe","returnedAt":"2026-07-19T10:00:00Z"}}`)
	if err := processor.Process(context.Background(), payload); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if repository.applied.Ack.CausationID != repository.applied.Incoming.EventID || repository.applied.Ack.Type != entity.BookStockUpdatedEventType {
		t.Fatalf("ack = %#v", repository.applied.Ack)
	}
}
