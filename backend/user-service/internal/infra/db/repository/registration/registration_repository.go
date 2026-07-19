package registration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
)

type Repository struct{ db *sqlx.DB }

func NewRepository(db *sqlx.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Prepare(ctx context.Context, candidate *entity.Registration) (*entity.Registration, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO registration_operations (id, name, email, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		candidate.ID, candidate.Name, candidate.Email, candidate.Status, candidate.CreatedAt, candidate.UpdatedAt,
	)
	if err != nil && !isDuplicate(err) {
		return nil, err
	}
	var registration entity.Registration
	if err := r.db.GetContext(ctx, &registration, `
		SELECT id, name, email, status, created_at, updated_at
		FROM registration_operations WHERE email = ?`, candidate.Email); err != nil {
		return nil, err
	}
	return &registration, nil
}

func (r *Repository) Complete(ctx context.Context, registrationID string, user *entity.User, event *entity.OutboxEvent) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status entity.RegistrationStatus
	if err := tx.GetContext(ctx, &status, `
		SELECT status FROM registration_operations WHERE id = ? FOR UPDATE`, registrationID); err != nil {
		return err
	}
	if status != entity.RegistrationPending {
		return errs.ErrConflict
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO users (id, name, email, role_code, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.Role, user.CreatedAt, user.UpdatedAt,
	); err != nil {
		if isDuplicate(err) {
			return errs.ErrConflict
		}
		return err
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO outbox_events
			(id, event_type, aggregate_id, payload, occurred_at, available_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		event.ID, event.Type, event.AggregateID, event.Payload, event.OccurredAt, event.OccurredAt, event.OccurredAt,
	); err != nil {
		return err
	}
	result, err := tx.ExecContext(ctx, `
		UPDATE registration_operations
		SET status = 'completed', identity_id = ?, updated_at = ?
		WHERE id = ? AND status = 'pending'`, user.ID, user.UpdatedAt, registrationID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return fmt.Errorf("complete registration %s: expected one pending operation, updated %d", registrationID, changed)
	}
	return tx.Commit()
}

func (r *Repository) MarkConflict(ctx context.Context, registrationID string) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE registration_operations SET status = 'conflict', updated_at = UTC_TIMESTAMP(6)
		WHERE id = ? AND status = 'pending'`, registrationID)
	if err != nil {
		return err
	}
	changed, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if changed != 1 {
		return sql.ErrNoRows
	}
	return nil
}

func isDuplicate(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}
