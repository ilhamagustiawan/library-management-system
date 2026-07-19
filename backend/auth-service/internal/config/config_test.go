package config

import (
	"strings"
	"testing"
)

func TestLoadAllowsAdministrativeCommandsWithoutClientSecret(t *testing.T) {
	t.Setenv("OAUTH_CLIENT_SECRET", "")
	t.Setenv("OAUTH_INTROSPECTION_CLIENT_SECRET", "")

	if _, err := Load(); err != nil {
		t.Fatalf("Load() error = %v, want config without seeded-client credentials", err)
	}
}

func TestLoadConfiguresSupportedOAuthScopes(t *testing.T) {
	t.Setenv("OAUTH_CLIENT_SCOPES", "library:read  library:write")

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
