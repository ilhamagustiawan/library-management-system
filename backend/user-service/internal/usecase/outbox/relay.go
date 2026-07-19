package outbox

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/repository"
)

type Publisher interface {
	Publish(ctx context.Context, message entity.OutboxMessage) error
	Close() error
}

type Config struct {
	WorkerID     string
	BatchSize    int
	Lease        time.Duration
	PollInterval time.Duration
	BaseRetry    time.Duration
	MaxRetry     time.Duration
}

type Relay struct {
	repository repository.OutboxRepository
	publisher  Publisher
	config     Config
	now        func() time.Time
}

func NewRelay(repository repository.OutboxRepository, publisher Publisher, config Config) *Relay {
	if config.BatchSize <= 0 {
		config.BatchSize = 50
	}
	if config.Lease <= 0 {
		config.Lease = 30 * time.Second
	}
	if config.PollInterval <= 0 {
		config.PollInterval = 500 * time.Millisecond
	}
	if config.BaseRetry <= 0 {
		config.BaseRetry = time.Second
	}
	if config.MaxRetry <= 0 {
		config.MaxRetry = time.Minute
	}
	return &Relay{repository: repository, publisher: publisher, config: config, now: func() time.Time { return time.Now().UTC() }}
}

func (r *Relay) Run(ctx context.Context) {
	ticker := time.NewTicker(r.config.PollInterval)
	defer ticker.Stop()
	for {
		if err := r.Dispatch(ctx); err != nil && ctx.Err() == nil {
			slog.Error("outbox dispatch failed", "operation", "claim or update outbox", "error", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (r *Relay) Dispatch(ctx context.Context) error {
	now := r.now()
	messages, err := r.repository.Claim(ctx, r.config.WorkerID, now, now.Add(r.config.Lease), r.config.BatchSize)
	if err != nil {
		return fmt.Errorf("claim outbox events: %w", err)
	}
	for _, message := range messages {
		if err := r.publisher.Publish(ctx, message); err != nil {
			failure := truncate(err.Error(), 1000)
			if markErr := r.repository.MarkFailed(ctx, message.ID, r.config.WorkerID, failure, now.Add(r.retryDelay(message.Attempts))); markErr != nil {
				return fmt.Errorf("record event %s publish failure: %w", message.ID, markErr)
			}
			continue
		}
		if err := r.repository.MarkPublished(ctx, message.ID, r.config.WorkerID, now); err != nil {
			return fmt.Errorf("mark event %s published: %w", message.ID, err)
		}
	}
	return nil
}

func (r *Relay) retryDelay(attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	delay := r.config.BaseRetry
	for step := 1; step < attempts && delay < r.config.MaxRetry; step++ {
		delay *= 2
		if delay > r.config.MaxRetry {
			return r.config.MaxRetry
		}
	}
	return delay
}

func truncate(value string, limit int) string {
	value = strings.TrimSpace(value)
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
