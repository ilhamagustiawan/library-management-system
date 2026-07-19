package server

import (
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
)

func TestFiberConfigDoesNotTrustForwardedClientIP(t *testing.T) {
	cfg := fiberConfig(config.Config{Service: config.ServiceConfig{Name: "library-auth-service"}})

	if cfg.ProxyHeader != "" || cfg.EnableTrustedProxyCheck || len(cfg.TrustedProxies) != 0 {
		t.Fatalf("proxy trust enabled: %#v", cfg)
	}
}

func TestFiberConfigShowsStartupMessage(t *testing.T) {
	cfg := fiberConfig(config.Config{})

	if cfg.DisableStartupMessage {
		t.Fatal("DisableStartupMessage = true, want Fiber startup banner enabled")
	}
}
