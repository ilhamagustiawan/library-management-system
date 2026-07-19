package outbox

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type fakeRepository struct {
	messages  []entity.OutboxMessage
	published []string
	failed    []string
	retryAt   time.Time
}

func (r *fakeRepository) Claim(context.Context, string, time.Time, time.Time, int) ([]entity.OutboxMessage, error) {
	return append([]entity.OutboxMessage(nil), r.messages...), nil
}

func (r *fakeRepository) MarkPublished(_ context.Context, eventID, _ string, _ time.Time) error {
	r.published = append(r.published, eventID)
	return nil
}

func (r *fakeRepository) MarkFailed(_ context.Context, eventID, _ string, _ string, availableAt time.Time) error {
	r.failed = append(r.failed, eventID)
	r.retryAt = availableAt
	return nil
}

type fakePublisher struct{ err error }

func (p fakePublisher) Publish(context.Context, entity.OutboxMessage) error { return p.err }
func (fakePublisher) Close() error                                          { return nil }

func TestDispatchMarksConfirmedEventPublished(t *testing.T) {
	repository := &fakeRepository{messages: []entity.OutboxMessage{{
		OutboxEvent: entity.OutboxEvent{ID: "event-1", Type: entity.UserRegisteredV1}, Attempts: 1,
	}}}
	relay := NewRelay(repository, fakePublisher{}, Config{WorkerID: "worker-1", BatchSize: 50})
	relay.now = func() time.Time { return time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC) }

	if err := relay.Dispatch(context.Background()); err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if len(repository.published) != 1 || repository.published[0] != "event-1" || len(repository.failed) != 0 {
		t.Fatalf("published = %v, failed = %v", repository.published, repository.failed)
	}
}

func TestDispatchSchedulesBoundedRetryAfterPublishFailure(t *testing.T) {
	repository := &fakeRepository{messages: []entity.OutboxMessage{{
		OutboxEvent: entity.OutboxEvent{ID: "event-1", Type: entity.UserRegisteredV1}, Attempts: 3,
	}}}
	relay := NewRelay(repository, fakePublisher{err: errors.New("broker down")}, Config{
		WorkerID: "worker-1", BatchSize: 50, BaseRetry: time.Second, MaxRetry: time.Minute,
	})
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	relay.now = func() time.Time { return now }

	if err := relay.Dispatch(context.Background()); err != nil {
		t.Fatalf("Dispatch() error = %v", err)
	}
	if len(repository.failed) != 1 || !repository.retryAt.Equal(now.Add(4*time.Second)) {
		t.Fatalf("failed = %v, retryAt = %v", repository.failed, repository.retryAt)
	}
}
