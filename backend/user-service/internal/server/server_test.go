package server

import (
	"testing"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/config"
)

func TestFiberConfigTrustsOnlyConfiguredProxy(t *testing.T) {
	result := fiberConfig(config.Config{Service: config.ServiceConfig{
		Name: "library-user-service", TrustedProxies: []string{"172.30.0.2"},
	}})
	if result.ProxyHeader != "X-Real-IP" || !result.EnableTrustedProxyCheck || !result.EnableIPValidation {
		t.Fatalf("fiber config = %#v", result)
	}
	if len(result.TrustedProxies) != 1 || result.TrustedProxies[0] != "172.30.0.2" {
		t.Fatalf("trusted proxies = %v", result.TrustedProxies)
	}
}
