package oauthclient

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

var ErrNotFound = errors.New("oauth client not found")

type Kind string

const (
	KindAuthorizationCode Kind = "authorization_code"
	KindClientCredentials Kind = "client_credentials"
	KindResourceServer    Kind = "resource_server"
)

type Client struct {
	ID          string         `db:"id"`
	Name        string         `db:"name"`
	Kind        Kind           `db:"kind"`
	SecretHash  string         `db:"secret_hash"`
	RedirectURI string         `db:"redirect_uri"`
	Scopes      []string       `db:"-"`
	Public      bool           `db:"is_public"`
	UserID      sql.NullString `db:"user_id"`
	CreatedAt   time.Time      `db:"created_at"`
}

func (c *Client) GetID() string     { return c.ID }
func (c *Client) GetSecret() string { return "" }
func (c *Client) GetDomain() string { return c.RedirectURI }
func (c *Client) IsPublic() bool    { return c.Public }
func (c *Client) GetUserID() string { return c.UserID.String }
func (c *Client) VerifyPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.SecretHash), []byte(raw)) == nil
}
func (c *Client) ScopeList() []string { return append([]string(nil), c.Scopes...) }
func (c *Client) AllowsGrant(grant oauth2.GrantType) bool {
	switch c.Kind {
	case KindAuthorizationCode:
		return grant == oauth2.AuthorizationCode || grant == oauth2.Refreshing
	case KindClientCredentials:
		return grant == oauth2.ClientCredentials
	case KindResourceServer:
		return false
	default:
		return false
	}
}
func (c *Client) CanIntrospect() bool { return c.Kind == KindResourceServer }

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetByID(ctx context.Context, id string) (oauth2.ClientInfo, error) {
	return s.GetClient(ctx, id)
}

func (s *Store) GetClient(ctx context.Context, id string) (*Client, error) {
	const query = `
		SELECT id, name, kind, secret_hash, redirect_uri, is_public, user_id, created_at
		FROM oauth_clients WHERE id = ?`
	var client Client
	if err := s.db.GetContext(ctx, &client, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	scopes, err := s.GetScopes(ctx, id)
	if err != nil {
		return nil, err
	}
	client.Scopes = make([]string, 0, len(scopes))
	for _, scope := range scopes {
		client.Scopes = append(client.Scopes, scope.Code)
	}
	return &client, nil
}

func (s *Store) GetScopes(ctx context.Context, clientID string) ([]entity.Scope, error) {
	var scopes []entity.Scope
	err := s.db.SelectContext(ctx, &scopes, `
		SELECT scopes.code, scopes.audience
		FROM oauth_client_scopes
		JOIN scopes ON scopes.code = oauth_client_scopes.scope_code
		WHERE oauth_client_scopes.client_id = ?
		ORDER BY scopes.code`, clientID)
	return scopes, err
}

func (s *Store) GetRoleScopes(ctx context.Context, role entity.Role) ([]entity.Scope, error) {
	var scopes []entity.Scope
	err := s.db.SelectContext(ctx, &scopes, `
		SELECT scopes.code, scopes.audience
		FROM role_scopes
		JOIN scopes ON scopes.code = role_scopes.scope_code
		WHERE role_scopes.role_code = ?
		ORDER BY scopes.code`, role)
	return scopes, err
}

func (s *Store) Create(ctx context.Context, client *Client) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	const query = `
		INSERT INTO oauth_clients (id, name, kind, secret_hash, redirect_uri, is_public, user_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	if _, err = tx.ExecContext(ctx, query,
		client.ID, client.Name, client.Kind, client.SecretHash, client.RedirectURI, client.Public, client.UserID, client.CreatedAt,
	); err != nil {
		return err
	}
	for _, scope := range client.Scopes {
		if _, err = tx.ExecContext(ctx,
			`INSERT INTO oauth_client_scopes (client_id, scope_code) VALUES (?, ?)`, client.ID, scope,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) AllowsGrant(ctx context.Context, clientID string, grant oauth2.GrantType) (bool, error) {
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return false, err
	}
	return client.AllowsGrant(grant), nil
}

func (s *Store) CanIntrospect(ctx context.Context, clientID string) (bool, error) {
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return false, err
	}
	return client.CanIntrospect(), nil
}

func (s *Store) AllowsScopes(ctx context.Context, clientID, requested string) (bool, error) {
	client, err := s.GetClient(ctx, clientID)
	if err != nil {
		return false, err
	}
	allowed := make(map[string]struct{}, len(client.ScopeList()))
	for _, scope := range client.ScopeList() {
		allowed[scope] = struct{}{}
	}
	for _, scope := range strings.Fields(requested) {
		if _, ok := allowed[scope]; !ok {
			return false, nil
		}
	}
	return true, nil
}
