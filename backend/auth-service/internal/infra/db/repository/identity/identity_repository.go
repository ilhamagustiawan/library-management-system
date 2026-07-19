package identity

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, key string, user *entity.User) (*entity.User, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	existing, err := findByKey(ctx, tx, key)
	if err == nil {
		return matchReplay(existing, user.Email)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO users (id, name, email, password_hash, role_code, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Name, user.Email, user.PasswordHash, user.Role, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		_ = tx.Rollback()
		return r.resolveConflict(ctx, key, user.Email, err)
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO identity_requests (idempotency_key, user_id, request_email, created_at)
		VALUES (?, ?, ?, ?)`, key, user.ID, user.Email, user.CreatedAt)
	if err != nil {
		_ = tx.Rollback()
		return r.resolveConflict(ctx, key, user.Email, err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	copy := *user
	return &copy, nil
}

type queryer interface {
	GetContext(ctx context.Context, destination any, query string, args ...any) error
}

func findByKey(ctx context.Context, db queryer, key string) (*entity.User, error) {
	var user entity.User
	err := db.GetContext(ctx, &user, `
		SELECT users.id, users.name, users.email, users.password_hash, users.role_code, users.created_at, users.updated_at
		FROM identity_requests
		JOIN users ON users.id = identity_requests.user_id
		WHERE identity_requests.idempotency_key = ?`, key)
	return &user, err
}

func (r *Repository) resolveConflict(ctx context.Context, key, email string, cause error) (*entity.User, error) {
	existing, err := findByKey(ctx, r.db, key)
	if err == nil {
		return matchReplay(existing, email)
	}
	var mysqlErr *mysql.MySQLError
	if errors.As(cause, &mysqlErr) && mysqlErr.Number == 1062 {
		return nil, errs.New(http.StatusConflict, errs.CodeEmailExists, "email is already registered; no identity was changed", nil, cause)
	}
	return nil, cause
}

func matchReplay(existing *entity.User, email string) (*entity.User, error) {
	if existing.Email != email {
		return nil, errs.New(http.StatusConflict, errs.CodeIdempotencyConflict, "idempotency key was already used for another email; existing identity remains unchanged", nil, nil)
	}
	return existing, nil
}
