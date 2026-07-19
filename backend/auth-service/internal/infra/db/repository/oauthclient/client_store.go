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
)

var ErrNotFound = errors.New("oauth client not found")

type Kind string

const (
	KindAuthorizationCode Kind = "authorization_code"
	KindClientCredentials Kind = "client_credentials"
	KindResourceServer    Kind = "resource_server"
)

type Client struct {
	ID            string         `db:"id"`
	Name          string         `db:"name"`
	Kind          Kind           `db:"kind"`
	SecretHash    string         `db:"secret_hash"`
	RedirectURI   string         `db:"redirect_uri"`
	AllowedScopes string         `db:"allowed_scopes"`
	Public        bool           `db:"is_public"`
	UserID        sql.NullString `db:"user_id"`
	CreatedAt     time.Time      `db:"created_at"`
}

func (c *Client) GetID() string     { return c.ID }
func (c *Client) GetSecret() string { return "" }
func (c *Client) GetDomain() string { return c.RedirectURI }
func (c *Client) IsPublic() bool    { return c.Public }
func (c *Client) GetUserID() string { return c.UserID.String }
func (c *Client) VerifyPassword(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(c.SecretHash), []byte(raw)) == nil
}
func (c *Client) ScopeList() []string { return strings.Fields(c.AllowedScopes) }
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
		SELECT id, name, kind, secret_hash, redirect_uri, allowed_scopes, is_public, user_id, created_at
		FROM oauth_clients WHERE id = ?`
	var client Client
	if err := s.db.GetContext(ctx, &client, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &client, nil
}

func (s *Store) Create(ctx context.Context, client *Client) error {
	const query = `
		INSERT INTO oauth_clients (id, name, kind, secret_hash, redirect_uri, allowed_scopes, is_public, user_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := s.db.ExecContext(ctx, query,
		client.ID, client.Name, client.Kind, client.SecretHash, client.RedirectURI, client.AllowedScopes,
		client.Public, client.UserID, client.CreatedAt,
	)
	return err
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
