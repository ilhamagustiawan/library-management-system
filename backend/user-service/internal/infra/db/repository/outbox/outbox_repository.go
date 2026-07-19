package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
)

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Claim(ctx context.Context, workerID string, now, leaseUntil time.Time, limit int) ([]entity.OutboxMessage, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var messages []entity.OutboxMessage
	if err := tx.SelectContext(ctx, &messages, `
		SELECT id, event_type, aggregate_id, payload, occurred_at, attempts + 1 AS attempts
		FROM outbox_events
		WHERE published_at IS NULL
		  AND available_at <= ?
		  AND (claimed_until IS NULL OR claimed_until < ?)
		ORDER BY occurred_at, id
		LIMIT ?
		FOR UPDATE SKIP LOCKED`, now, now, limit); err != nil {
		return nil, err
	}
	for _, message := range messages {
		result, err := tx.ExecContext(ctx, `
			UPDATE outbox_events
			SET claimed_by = ?, claimed_until = ?, attempts = attempts + 1
			WHERE id = ? AND published_at IS NULL`, workerID, leaseUntil, message.ID)
		if err != nil {
			return nil, err
		}
		changed, err := result.RowsAffected()
		if err != nil || changed != 1 {
			return nil, fmt.Errorf("claim outbox event %s: updated %d rows: %w", message.ID, changed, err)
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *Repository) MarkPublished(ctx context.Context, eventID, workerID string, publishedAt time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE outbox_events
		SET published_at = ?, claimed_by = NULL, claimed_until = NULL, last_error = NULL
		WHERE id = ? AND claimed_by = ? AND published_at IS NULL`, publishedAt, eventID, workerID)
	return exactlyOne(result, err, "mark published", eventID)
}

func (r *Repository) MarkFailed(ctx context.Context, eventID, workerID, message string, availableAt time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE outbox_events
		SET available_at = ?, claimed_by = NULL, claimed_until = NULL, last_error = ?
		WHERE id = ? AND claimed_by = ? AND published_at IS NULL`, availableAt, message, eventID, workerID)
	return exactlyOne(result, err, "mark failed", eventID)
}

func exactlyOne(result sql.Result, err error, operation, eventID string) error {
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return fmt.Errorf("%s outbox event %s: updated %d rows", operation, eventID, changed)
	}
	return nil
}
