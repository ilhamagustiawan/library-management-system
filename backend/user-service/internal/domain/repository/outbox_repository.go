package repository

import (
	"context"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type OutboxRepository interface {
	Claim(ctx context.Context, workerID string, now, leaseUntil time.Time, limit int) ([]entity.OutboxMessage, error)
	MarkPublished(ctx context.Context, eventID, workerID string, publishedAt time.Time) error
	MarkFailed(ctx context.Context, eventID, workerID, message string, availableAt time.Time) error
}
