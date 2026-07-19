package middleware

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/oauth"
)

func TestAuthorizeGatewayAcceptsEitherCatalogScope(t *testing.T) {
	err := authorizeGateway("books:manage", "admin-id", []string{"books:read", "books:manage"})
	if err != nil {
		t.Fatalf("authorizeGateway() error = %v", err)
	}
}

func TestAuthorizeGatewayDoesNotCheckPeerAddress(t *testing.T) {
	err := authorizeGateway("books:manage", "admin-id", []string{"books:manage"})
	if err != nil {
		t.Fatalf("authorizeGateway() error = %v", err)
	}
}

func TestAuthorizeInternalRejectsWrongAudience(t *testing.T) {
	principal := oauth.Principal{
		Active: true, ClientID: "transaction-service", Subject: "transaction-service",
		Audience: []string{"library-api"}, Scope: "book-stock:reserve", TokenType: "Bearer",
		Issuer: "http://auth-service:8081", IssuedAt: time.Now().Add(-time.Minute).Unix(),
		ExpiresAt: time.Now().Add(time.Minute).Unix(),
	}
	err := authorizeInternal(principal, InternalPolicy{
		Issuer: "http://auth-service:8081", Audience: "book-service",
		ClientID: "transaction-service", Scope: "book-stock:reserve", Now: time.Now,
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.HTTPStatus != http.StatusForbidden {
		t.Fatalf("authorizeInternal() error = %v, want forbidden", err)
	}
}

func TestAuthorizeInternalRejectsFutureIssuedToken(t *testing.T) {
	now := time.Date(2026, 7, 19, 10, 0, 0, 0, time.UTC)
	principal := oauth.Principal{
		Active: true, ClientID: "transaction-service", Subject: "transaction-service",
		Audience: []string{"book-service"}, Scope: "book-stock:read", TokenType: "Bearer",
		Issuer: "http://auth-service:8081", IssuedAt: now.Add(time.Minute).Unix(), ExpiresAt: now.Add(2 * time.Minute).Unix(),
	}
	err := authorizeInternal(principal, InternalPolicy{
		Issuer: "http://auth-service:8081", Audience: "book-service",
		ClientID: "transaction-service", Scope: "book-stock:read", Now: func() time.Time { return now },
	})
	var domainErr *errs.Error
	if !errors.As(err, &domainErr) || domainErr.HTTPStatus != http.StatusForbidden {
		t.Fatalf("authorizeInternal() error = %v, want forbidden", err)
	}
}
