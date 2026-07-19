package server

import (
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
)

func TestFiberConfigTrustsConfiguredProxyForClientIP(t *testing.T) {
	cfg := fiberConfig(config.Config{Service: config.ServiceConfig{
		Name: "library-auth-service", TrustedProxies: []string{"172.30.0.2"},
	}})

	if cfg.ProxyHeader != "X-Real-IP" {
		t.Fatalf("ProxyHeader = %q, want X-Real-IP", cfg.ProxyHeader)
	}
	if !cfg.EnableTrustedProxyCheck || !cfg.EnableIPValidation {
		t.Fatalf("trusted proxy controls disabled: %#v", cfg)
	}
	if len(cfg.TrustedProxies) != 1 || cfg.TrustedProxies[0] != "172.30.0.2" {
		t.Fatalf("TrustedProxies = %#v, want configured gateway", cfg.TrustedProxies)
	}
}

func TestFiberConfigShowsStartupMessage(t *testing.T) {
	cfg := fiberConfig(config.Config{})

	if cfg.DisableStartupMessage {
		t.Fatal("DisableStartupMessage = true, want Fiber startup banner enabled")
	}
}
