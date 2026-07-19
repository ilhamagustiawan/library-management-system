package book

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type OAuthConfig struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
}

type OAuthTokenSource struct {
	config    OAuthConfig
	http      *http.Client
	mutex     sync.Mutex
	token     string
	expiresAt time.Time
	now       func() time.Time
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func NewOAuthTokenSource(config OAuthConfig) (*OAuthTokenSource, error) {
	parsed, err := url.Parse(strings.TrimSpace(config.TokenURL))
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return nil, fmt.Errorf("AUTH_TOKEN_URL must be an absolute HTTP(S) URL")
	}
	if strings.TrimSpace(config.ClientID) == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("Book Service OAuth client credentials are required")
	}
	if config.Timeout <= 0 {
		config.Timeout = 2 * time.Second
	}
	return &OAuthTokenSource{config: config, http: &http.Client{Timeout: config.Timeout}, now: func() time.Time { return time.Now().UTC() }}, nil
}

func (s *OAuthTokenSource) Token(ctx context.Context) (string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.token != "" && s.expiresAt.After(s.now().Add(30*time.Second)) {
		return s.token, nil
	}
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {"book-stock:reserve book-stock:release"}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.config.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.SetBasicAuth(s.config.ClientID, s.config.ClientSecret)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := s.http.Do(request)
	if err != nil {
		return "", fmt.Errorf("request OAuth token: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OAuth token endpoint returned HTTP %d", response.StatusCode)
	}
	var token tokenResponse
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes))
	if err := decoder.Decode(&token); err != nil {
		return "", fmt.Errorf("decode OAuth token: %w", err)
	}
	if token.AccessToken == "" || !strings.EqualFold(token.TokenType, "bearer") || token.ExpiresIn <= 30 {
		return "", fmt.Errorf("OAuth token response is incomplete")
	}
	s.token, s.expiresAt = token.AccessToken, s.now().Add(time.Duration(token.ExpiresIn)*time.Second)
	return s.token, nil
}
