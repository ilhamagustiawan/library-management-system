package config

import (
	"strings"
	"testing"
)

func TestLoadRequiresServiceCredentials(t *testing.T) {
	t.Setenv("AUTH_CLIENT_SECRET", "short")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "AUTH_CLIENT_SECRET") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestProductionRequiresSecureDependencyURLs(t *testing.T) {
	t.Setenv("SERVICE_ENV", "production")
	t.Setenv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("AUTH_SERVICE_URL", "http://auth.internal")
	t.Setenv("RABBITMQ_URL", "amqp://rabbit.internal")
	_, err := Load()
	if err == nil || !strings.Contains(err.Error(), "HTTPS") {
		t.Fatalf("Load() error = %v", err)
	}
}

func TestDevelopmentLoadsExplicitLocalDefaults(t *testing.T) {
	t.Setenv("SERVICE_ENV", "development")
	t.Setenv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("DATABASE_MIGRATION_URL", "")
	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Service.Port != ":8082" || !strings.Contains(config.Database.DSN, ":3307") || config.Auth.Scope != "identities:create" || config.Outbox.BatchSize != 50 {
		t.Fatalf("config = %#v", config)
	}
}
