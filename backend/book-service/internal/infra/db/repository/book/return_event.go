package book

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/repository"
)

const stockUpdatedRoutingKey = "books.stock.updated.v1"

func (r *Repository) ApplyLoanReturned(ctx context.Context, command repository.StockReturnCommand) error {
	transaction, err := r.database.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()
	result, err := transaction.ExecContext(ctx, `INSERT IGNORE INTO message_inbox (event_id, event_type, processed_at) VALUES (?, ?, ?)`, command.Incoming.EventID, command.Incoming.Type, command.Ack.OccurredAt)
	if err != nil {
		return err
	}
	if affected, _ := result.RowsAffected(); affected == 0 {
		return transaction.Commit()
	}
	reservation, err := getReservation(ctx, transaction, command.Incoming.Data.LoanID, true)
	if err != nil {
		return err
	}
	if reservation.BookID != command.Incoming.Data.BookID {
		return fmt.Errorf("loan-returned book does not match reservation")
	}
	if reservation.Status == entity.ReservationActive {
		var counts struct {
			Total     int `db:"total_copies"`
			Available int `db:"available_copies"`
		}
		if err := transaction.GetContext(ctx, &counts, `SELECT total_copies, available_copies FROM books WHERE id = ? FOR UPDATE`, reservation.BookID); err != nil {
			return err
		}
		if counts.Available >= counts.Total {
			return fmt.Errorf("release returned stock: stored copy counts are inconsistent")
		}
		if _, err := transaction.ExecContext(ctx, `UPDATE books SET available_copies = available_copies + 1, updated_at = ? WHERE id = ?`, command.Ack.OccurredAt, reservation.BookID); err != nil {
			return err
		}
		if _, err := transaction.ExecContext(ctx, `UPDATE book_reservations SET status = 'released', released_at = ? WHERE transaction_id = ?`, command.Ack.OccurredAt, reservation.TransactionID); err != nil {
			return err
		}
	}
	payload, err := json.Marshal(command.Ack)
	if err != nil {
		return fmt.Errorf("marshal BookStockUpdated event: %w", err)
	}
	if _, err := transaction.ExecContext(ctx, `INSERT INTO message_outbox (id, event_type, routing_key, payload, next_attempt_at, created_at) VALUES (?, ?, ?, ?, ?, ?)`, command.Ack.EventID, command.Ack.Type, stockUpdatedRoutingKey, payload, command.Ack.OccurredAt, command.Ack.OccurredAt); err != nil {
		return err
	}
	return transaction.Commit()
}

func (r *Repository) PendingMessages(ctx context.Context, limit int) ([]repository.OutboxMessage, error) {
	rows, err := r.database.QueryxContext(ctx, `SELECT id, event_type, routing_key, payload, attempts FROM message_outbox WHERE status = 'pending' AND next_attempt_at <= NOW(6) ORDER BY created_at LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]repository.OutboxMessage, 0)
	for rows.Next() {
		var message repository.OutboxMessage
		if err := rows.Scan(&message.ID, &message.EventType, &message.RoutingKey, &message.Payload, &message.Attempts); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, rows.Err()
}

func (r *Repository) MarkPublished(ctx context.Context, id string, now time.Time) error {
	_, err := r.database.ExecContext(ctx, `UPDATE message_outbox SET status = 'published', published_at = ?, last_error = NULL WHERE id = ?`, now, id)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id, message string, now time.Time) error {
	message = strings.TrimSpace(message)
	if len(message) > 1000 {
		message = message[:1000]
	}
	_, err := r.database.ExecContext(ctx, `UPDATE message_outbox SET attempts = attempts + 1, last_error = ?, next_attempt_at = DATE_ADD(?, INTERVAL LEAST(POW(2, attempts + 1), 60) SECOND) WHERE id = ?`, message, now, id)
	return err
}
