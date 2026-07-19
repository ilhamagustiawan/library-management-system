package oauthtoken

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, info oauth2.TokenInfo) error {
	payload, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("marshal oauth token: %w", err)
	}

	var code, access, refresh any
	var expiresAt time.Time
	if info.GetCode() != "" {
		code = info.GetCode()
		expiresAt = expiry(info.GetCodeCreateAt(), info.GetCodeExpiresIn())
	} else {
		access = info.GetAccess()
		if info.GetRefresh() != "" {
			refresh = info.GetRefresh()
			expiresAt = expiry(info.GetRefreshCreateAt(), info.GetRefreshExpiresIn())
		} else {
			expiresAt = expiry(info.GetAccessCreateAt(), info.GetAccessExpiresIn())
		}
	}

	const query = `
		INSERT INTO oauth_tokens (id, code, access_token, refresh_token, payload, expires_at)
		VALUES (?, ?, ?, ?, ?, ?)`
	_, err = s.db.ExecContext(ctx, query, uuid.NewString(), code, access, refresh, payload, expiresAt)
	return err
}

func (s *Store) RemoveByCode(ctx context.Context, code string) error {
	return s.remove(ctx, `DELETE FROM oauth_tokens WHERE code = ?`, code)
}

func (s *Store) RemoveByAccess(ctx context.Context, access string) error {
	return s.remove(ctx, `DELETE FROM oauth_tokens WHERE access_token = ?`, access)
}

func (s *Store) RemoveByRefresh(ctx context.Context, refresh string) error {
	return s.remove(ctx, `DELETE FROM oauth_tokens WHERE refresh_token = ?`, refresh)
}

func (s *Store) GetByCode(ctx context.Context, code string) (oauth2.TokenInfo, error) {
	return s.get(ctx, `SELECT payload FROM oauth_tokens WHERE code = ? AND expires_at > NOW(6)`, code)
}

func (s *Store) GetByAccess(ctx context.Context, access string) (oauth2.TokenInfo, error) {
	return s.get(ctx, `SELECT payload FROM oauth_tokens WHERE access_token = ? AND expires_at > NOW(6)`, access)
}

func (s *Store) GetByRefresh(ctx context.Context, refresh string) (oauth2.TokenInfo, error) {
	return s.get(ctx, `SELECT payload FROM oauth_tokens WHERE refresh_token = ? AND expires_at > NOW(6)`, refresh)
}

func (s *Store) DeleteExpired(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM oauth_tokens WHERE expires_at <= NOW(6)`)
	return err
}

func (s *Store) remove(ctx context.Context, query, value string) error {
	_, err := s.db.ExecContext(ctx, query, value)
	return err
}

func (s *Store) get(ctx context.Context, query, value string) (oauth2.TokenInfo, error) {
	var payload []byte
	if err := s.db.GetContext(ctx, &payload, query, value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	var token models.Token
	if err := json.Unmarshal(payload, &token); err != nil {
		return nil, fmt.Errorf("unmarshal oauth token: %w", err)
	}
	return &token, nil
}

func expiry(createdAt time.Time, ttl time.Duration) time.Time {
	if ttl <= 0 {
		return createdAt.AddDate(100, 0, 0)
	}
	return createdAt.Add(ttl)
}
