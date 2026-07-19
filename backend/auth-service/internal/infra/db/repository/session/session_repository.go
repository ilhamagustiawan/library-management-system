package session

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
	domainerrs "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, session *entity.Session) error {
	const query = `
		INSERT INTO sessions (id, user_id, token_hash, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, session.ID, session.UserID, session.TokenHash, session.ExpiresAt, session.CreatedAt)
	return err
}

func (r *Repository) FindByTokenHash(ctx context.Context, tokenHash string) (*entity.Session, error) {
	const query = `
		SELECT id, user_id, token_hash, expires_at, created_at
		FROM sessions WHERE token_hash = ? AND expires_at > NOW(6)`
	var session entity.Session
	if err := r.db.GetContext(ctx, &session, query, tokenHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domainerrs.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *Repository) DeleteByTokenHash(ctx context.Context, tokenHash string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	return err
}
