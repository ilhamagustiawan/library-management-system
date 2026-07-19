package server

import (
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/config"
)

func TestFiberConfigDoesNotTrustForwardedClientIP(t *testing.T) {
	result := fiberConfig(config.Config{Service: config.ServiceConfig{Name: "library-user-service"}})

	if result.ProxyHeader != "" || result.EnableTrustedProxyCheck || len(result.TrustedProxies) != 0 {
		t.Fatalf("proxy trust enabled: %#v", result)
	}
}
