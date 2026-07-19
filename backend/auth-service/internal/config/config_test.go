package config

import (
	"strings"
	"testing"
)

func TestLoadAllowsAdministrativeCommandsWithoutClientSecret(t *testing.T) {
	t.Setenv("OAUTH_CLIENT_SECRET", "")
	t.Setenv("OAUTH_INTROSPECTION_CLIENT_SECRET", "")
	t.Setenv("OAUTH_JWT_SIGNING_KEY", "")

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error = %v, want config without runtime credentials", err)
	}
}

func TestLoadDefaultsOAuthIssuerToGateway(t *testing.T) {
	t.Setenv("OAUTH_ISSUER", "")

	cfg, err := Load()

	if err != nil || cfg.OAuth.Issuer != "http://localhost:8000" {
		t.Fatalf("Load() OAuth issuer = %q, error = %v; want gateway origin", cfg.OAuth.Issuer, err)
	}
}

func TestLoadServerRequiresStrongJWTSigningKey(t *testing.T) {
	t.Setenv("OAUTH_JWT_SIGNING_KEY", "short")

	_, err := LoadServer()

	if err == nil || !strings.Contains(err.Error(), "OAUTH_JWT_SIGNING_KEY") {
		t.Fatalf("LoadServer() error = %v, want JWT signing-key error", err)
	}
}

func TestLoadServerConfiguresJWTSigningKey(t *testing.T) {
	key := strings.Repeat("k", 32)
	t.Setenv("OAUTH_JWT_SIGNING_KEY", key)

	cfg, err := LoadServer()

	if err != nil || string(cfg.OAuth.JWTSigningKey) != key {
		t.Fatalf("LoadServer() JWT key = %q, error = %v", cfg.OAuth.JWTSigningKey, err)
	}
}

func TestLoadConfiguresSupportedOAuthScopes(t *testing.T) {
	t.Setenv("OAUTH_CLIENT_SCOPES", "books:read  loans:borrow:self")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.OAuth.SupportedScopes) != 2 {
		t.Fatalf("OAuth.SupportedScopes = %#v, want two normalized scopes", cfg.OAuth.SupportedScopes)
	}
}

func TestLoadRejectsEmptySupportedOAuthScopes(t *testing.T) {
	t.Setenv("OAUTH_CLIENT_SCOPES", " ")

	_, err := Load()

	if err == nil || !strings.Contains(err.Error(), "OAUTH_CLIENT_SCOPES") {
		t.Fatalf("Load() error = %v, want supported-scopes error", err)
	}
}

func TestLoadConfiguresTrustedProxies(t *testing.T) {
	t.Setenv("TRUSTED_PROXIES", "172.30.0.2 10.0.0.0/24")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Service.TrustedProxies) != 2 || cfg.Service.TrustedProxies[0] != "172.30.0.2" || cfg.Service.TrustedProxies[1] != "10.0.0.0/24" {
		t.Fatalf("Service.TrustedProxies = %#v, want normalized proxy list", cfg.Service.TrustedProxies)
	}
}
