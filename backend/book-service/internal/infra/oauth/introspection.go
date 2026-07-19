package oauth

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
	"time"
)

const maxIntrospectionResponseBytes = 64 << 10

var ErrUnavailable = errors.New("authentication service unavailable")

type Config struct {
	URL          string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
}

type Principal struct {
	Active    bool     `json:"active"`
	ClientID  string   `json:"client_id"`
	Subject   string   `json:"sub"`
	Scope     string   `json:"scope"`
	TokenType string   `json:"token_type"`
	Issuer    string   `json:"iss"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
	Audience  []string `json:"aud"`
}

type Client struct {
	config Config
	http   *http.Client
}

func NewClient(config Config) *Client {
	return newClient(config, &http.Client{Timeout: config.Timeout})
}

func newClient(config Config, httpClient *http.Client) *Client {
	return &Client{config: config, http: httpClient}
}

func (c *Client) Introspect(ctx context.Context, token string) (Principal, error) {
	form := url.Values{"token": {token}, "token_type_hint": {"access_token"}}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return Principal{}, fmt.Errorf("create introspection request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(c.config.ClientID, c.config.ClientSecret)

	response, err := c.http.Do(request)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: request introspection: %v", ErrUnavailable, err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, maxIntrospectionResponseBytes))
		return Principal{}, fmt.Errorf("%w: introspection returned HTTP %d", ErrUnavailable, response.StatusCode)
	}

	limited := io.LimitReader(response.Body, maxIntrospectionResponseBytes+1)
	payload, err := io.ReadAll(limited)
	if err != nil {
		return Principal{}, fmt.Errorf("%w: read introspection response: %v", ErrUnavailable, err)
	}
	if len(payload) > maxIntrospectionResponseBytes {
		return Principal{}, fmt.Errorf("%w: introspection response exceeds %d bytes", ErrUnavailable, maxIntrospectionResponseBytes)
	}
	var principal Principal
	decoder := json.NewDecoder(bytes.NewReader(payload))
	if err := decoder.Decode(&principal); err != nil {
		return Principal{}, fmt.Errorf("%w: decode introspection response: %v", ErrUnavailable, err)
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Principal{}, fmt.Errorf("%w: introspection response must contain one JSON object", ErrUnavailable)
	}
	return principal, nil
}
