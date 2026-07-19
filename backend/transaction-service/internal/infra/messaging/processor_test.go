package messaging

import (
	"context"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/repository"
)

type fakeStore struct {
	confirmed bool
	published bool
}

func (s *fakeStore) CreatePending(context.Context, *entity.Loan) error { return nil }
func (s *fakeStore) Activate(context.Context, string, string, time.Time) (*entity.Loan, error) {
	return nil, nil
}
func (s *fakeStore) CancelPending(context.Context, string, time.Time) error { return nil }
func (s *fakeStore) Get(context.Context, string) (*entity.Loan, error)      { return nil, nil }
func (s *fakeStore) Return(context.Context, repository.ReturnCommand) (*entity.Loan, bool, error) {
	return nil, false, nil
}
func (s *fakeStore) ListTransactions(context.Context, repository.PageFilter) (repository.TransactionPage, error) {
	return repository.TransactionPage{}, nil
}
func (s *fakeStore) StockSyncStatus(context.Context, string) (entity.StockSyncStatus, error) {
	return entity.StockSyncPending, nil
}
func (s *fakeStore) ConfirmStock(context.Context, string, string, string, string, time.Time) error {
	s.confirmed = true
	return nil
}
func (s *fakeStore) PendingOutbox(context.Context, int) ([]repository.OutboxMessage, error) {
	return []repository.OutboxMessage{{ID: "event-1", RoutingKey: "transactions.loan.returned.v1", Payload: []byte(`{}`)}}, nil
}
func (s *fakeStore) MarkOutboxPublished(context.Context, string, time.Time) error {
	s.published = true
	return nil
}
func (s *fakeStore) MarkOutboxFailed(context.Context, string, string, time.Time) error { return nil }

type fakePublisher struct{ called bool }

func (p *fakePublisher) Publish(context.Context, repository.OutboxMessage) error {
	p.called = true
	return nil
}

func TestAckProcessorConfirmsMatchingLoan(t *testing.T) {
	store := &fakeStore{}
	processor := NewAckProcessor(store)
	payload := []byte(`{"eventId":"0b99be15-1904-4df1-8c06-f5f95c248d6d","type":"BookStockUpdated.v1","occurredAt":"2026-07-19T10:00:00Z","causationId":"d840ca41-ed41-4da9-b5ea-95449a998dde","data":{"loanId":"52a88672-a4c2-4876-be5a-65863aeb35e4","bookId":"7b36fe43-f31d-4861-884f-42ed7386b1e9","updatedAt":"2026-07-19T10:00:00Z"}}`)
	if err := processor.Process(context.Background(), payload); err != nil {
		t.Fatalf("Process() error = %v", err)
	}
	if !store.confirmed {
		t.Fatal("stock update was not confirmed")
	}
}

func TestAckProcessorRejectsUnknownFields(t *testing.T) {
	processor := NewAckProcessor(&fakeStore{})
	payload := []byte(`{"eventId":"0b99be15-1904-4df1-8c06-f5f95c248d6d","type":"BookStockUpdated.v1","occurredAt":"2026-07-19T10:00:00Z","causationId":"d840ca41-ed41-4da9-b5ea-95449a998dde","data":{"loanId":"52a88672-a4c2-4876-be5a-65863aeb35e4","bookId":"7b36fe43-f31d-4861-884f-42ed7386b1e9","updatedAt":"2026-07-19T10:00:00Z"},"unexpected":true}`)
	if err := processor.Process(context.Background(), payload); err == nil {
		t.Fatal("Process() error = nil, want strict decoding error")
	}
}

func TestOutboxWorkerMarksConfirmedPublish(t *testing.T) {
	store, publisher := &fakeStore{}, &fakePublisher{}
	worker := NewOutboxWorker(store, publisher)
	if err := worker.PublishBatch(context.Background()); err != nil {
		t.Fatalf("PublishBatch() error = %v", err)
	}
	if !publisher.called || !store.published {
		t.Fatalf("publish = %v, marked = %v", publisher.called, store.published)
	}
}
