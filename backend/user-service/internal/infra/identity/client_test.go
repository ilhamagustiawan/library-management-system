package identity

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/entity"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
	registrationusecase "github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/usecase/registration"
)

func TestCreateGetsScopedTokenAndForwardsStableKey(t *testing.T) {
	var tokenCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			tokenCalls++
			clientID, secret, ok := r.BasicAuth()
			if !ok || clientID != "user-service" || secret != "client-secret" {
				t.Fatalf("Basic auth = (%q, %q, %v)", clientID, secret, ok)
			}
			_ = r.ParseForm()
			if r.Form.Get("grant_type") != "client_credentials" || r.Form.Get("scope") != "identities:create" {
				t.Fatalf("token form = %v", r.Form)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "service-token", "token_type": "Bearer", "expires_in": 300,
			})
		case "/internal/identities":
			if r.Header.Get("Authorization") != "Bearer service-token" || r.Header.Get("Idempotency-Key") != "registration-123" {
				t.Fatalf("identity headers = %v", r.Header)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"code": "LMS-200000",
				"data": map[string]any{"id": "user-123", "name": "Ada", "email": "ada@example.com", "role": "member"},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL, ClientID: "user-service", ClientSecret: "client-secret",
		Scope: "identities:create", Attempts: 3,
	}, server.Client())
	client.now = func() time.Time { return time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC) }
	client.wait = func(context.Context, time.Duration) error { return nil }

	for range 2 {
		identity, err := client.Create(context.Background(), "registration-123", registrationusecase.IdentityInput{
			Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
		})
		if err != nil || identity.ID != "user-123" || identity.Role != entity.RoleMember {
			t.Fatalf("Create() = (%#v, %v)", identity, err)
		}
	}
	if tokenCalls != 1 {
		t.Fatalf("token calls = %d, want cached token", tokenCalls)
	}
}

func TestCreateRetriesTransientFailureWithSameKey(t *testing.T) {
	var identityKeys []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/oauth/token" {
			_ = json.NewEncoder(w).Encode(map[string]any{"access_token": "token", "token_type": "Bearer", "expires_in": 300})
			return
		}
		identityKeys = append(identityKeys, r.Header.Get("Idempotency-Key"))
		if len(identityKeys) == 1 {
			http.Error(w, "temporary", http.StatusServiceUnavailable)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"code": "LMS-200000", "data": map[string]any{"id": "user-123", "name": "Ada", "email": "ada@example.com", "role": "member"},
		})
	}))
	defer server.Close()
	client := NewClient(Config{BaseURL: server.URL, ClientID: "user-service", ClientSecret: "secret", Scope: "identities:create", Attempts: 3}, server.Client())
	client.wait = func(context.Context, time.Duration) error { return nil }

	_, err := client.Create(context.Background(), "registration-123", registrationusecase.IdentityInput{
		Name: "Ada", Email: "ada@example.com", Password: "correct horse battery staple",
	})
	if err != nil || strings.Join(identityKeys, ",") != "registration-123,registration-123" {
		t.Fatalf("Create() error = %v, keys = %v", err, identityKeys)
	}
}

func TestTokenRequestUsesFormEncoding(t *testing.T) {
	values := url.Values{"grant_type": {"client_credentials"}, "scope": {"identities:create"}}
	if values.Encode() != "grant_type=client_credentials&scope=identities%3Acreate" {
		t.Fatalf("encoded form = %q", values.Encode())
	}
}

func TestCreateDistinguishesIdempotencyConflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if request.URL.Path == "/oauth/token" {
			response.Header().Set("Content-Type", "application/json")
			_, _ = response.Write([]byte(`{"access_token":"service-token","token_type":"Bearer","expires_in":300}`))
			return
		}
		response.Header().Set("Content-Type", "application/json")
		response.WriteHeader(http.StatusConflict)
		_, _ = response.Write([]byte(`{"code":"LMS-409002","message":"changed replay"}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		BaseURL: server.URL, ClientID: "user-service", ClientSecret: "secret", Scope: "identities:create", Attempts: 1,
	}, server.Client())
	_, err := client.Create(context.Background(), "registration-123", registrationusecase.IdentityInput{
		Name: "Ada", Email: "ada@example.com", Password: "different secure password",
	})
	if !errors.Is(err, errs.ErrIdentityIdempotency) {
		t.Fatalf("Create() error = %v", err)
	}
}
