package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

type Config struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Scope        string
	Attempts     int
}

type Client struct {
	config Config
	http   *http.Client

	mu        sync.Mutex
	token     string
	expiresAt time.Time
	now       func() time.Time
	wait      func(context.Context, time.Duration) error
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type identityResponse struct {
	Code string       `json:"code"`
	Data identityData `json:"data"`
}

type identityData struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Email string      `json:"email"`
	Role  entity.Role `json:"role"`
}

type errorResponse struct {
	Code string `json:"code"`
}

func NewClient(config Config, httpClient *http.Client) *Client {
	if config.Attempts < 1 {
		config.Attempts = 3
	}
	return &Client{
		config: config,
		http:   httpClient,
		now:    func() time.Time { return time.Now().UTC() },
		wait:   wait,
	}
}

func (c *Client) Create(ctx context.Context, key string, input registrationusecase.IdentityInput) (*registrationusecase.Identity, error) {
	var lastErr error
	for attempt := 0; attempt < c.config.Attempts; attempt++ {
		identity, retry, err := c.createOnce(ctx, key, input)
		if err == nil {
			return identity, nil
		}
		if errors.Is(err, errs.ErrIdentityConflict) || errors.Is(err, errs.ErrIdentityIdempotency) || errors.Is(err, errs.ErrInvalidIdentity) {
			return nil, err
		}
		lastErr = err
		if !retry || attempt == c.config.Attempts-1 {
			break
		}
		if err := c.wait(ctx, time.Duration(1<<attempt)*100*time.Millisecond); err != nil {
			return nil, errs.ErrIdentityUnavailable
		}
	}
	return nil, fmt.Errorf("%w: %v", errs.ErrIdentityUnavailable, lastErr)
}

func (c *Client) createOnce(ctx context.Context, key string, input registrationusecase.IdentityInput) (*registrationusecase.Identity, bool, error) {
	token, err := c.accessToken(ctx)
	if err != nil {
		return nil, true, err
	}
	body, err := json.Marshal(input)
	if err != nil {
		return nil, false, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.config.BaseURL, "/")+"/internal/identities", bytes.NewReader(body))
	if err != nil {
		return nil, false, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Idempotency-Key", key)

	response, err := c.http.Do(request)
	if err != nil {
		return nil, true, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusUnauthorized {
		c.clearToken()
		return nil, true, errors.New("identity token rejected")
	}
	if response.StatusCode == http.StatusConflict {
		var payload errorResponse
		if err := decodeBounded(response.Body, &payload); err == nil && payload.Code == errs.CodeIdempotencyConflict {
			return nil, false, errs.ErrIdentityIdempotency
		}
		return nil, false, errs.ErrIdentityConflict
	}
	if response.StatusCode == http.StatusUnprocessableEntity || response.StatusCode == http.StatusForbidden {
		return nil, false, errs.ErrInvalidIdentity
	}
	if response.StatusCode >= http.StatusInternalServerError {
		return nil, true, fmt.Errorf("identity service status %d", response.StatusCode)
	}
	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected identity service status %d", response.StatusCode)
	}
	var payload identityResponse
	if err := decodeBounded(response.Body, &payload); err != nil || payload.Code != errs.CodeSuccess {
		return nil, false, errs.ErrInvalidIdentity
	}
	if payload.Data.ID == "" || payload.Data.Name == "" || payload.Data.Email == "" || payload.Data.Role != entity.RoleMember {
		return nil, false, errs.ErrInvalidIdentity
	}
	return &registrationusecase.Identity{
		ID: payload.Data.ID, Name: payload.Data.Name, Email: payload.Data.Email, Role: payload.Data.Role,
	}, false, nil
}

func (c *Client) accessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.token != "" && c.expiresAt.After(c.now().Add(30*time.Second)) {
		return c.token, nil
	}

	form := url.Values{"grant_type": {"client_credentials"}, "scope": {c.config.Scope}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.config.BaseURL, "/")+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)
	response, err := c.http.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("token service status %d", response.StatusCode)
	}
	var payload tokenResponse
	if err := decodeBounded(response.Body, &payload); err != nil {
		return "", err
	}
	if payload.AccessToken == "" || !strings.EqualFold(payload.TokenType, "Bearer") || payload.ExpiresIn <= 30 {
		return "", errors.New("invalid service token response")
	}
	c.token = payload.AccessToken
	c.expiresAt = c.now().Add(time.Duration(payload.ExpiresIn) * time.Second)
	return c.token, nil
}

func (c *Client) clearToken() {
	c.mu.Lock()
	c.token = ""
	c.expiresAt = time.Time{}
	c.mu.Unlock()
}

func decodeBounded(reader io.Reader, output any) error {
	decoder := json.NewDecoder(io.LimitReader(reader, 1<<20))
	return decoder.Decode(output)
}

func wait(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
