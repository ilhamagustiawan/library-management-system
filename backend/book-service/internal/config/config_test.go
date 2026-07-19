package config

import (
	"strings"
	"testing"
	"time"
)

func TestValidateServerRejectsMissingIntrospectionSecret(t *testing.T) {
	config := validConfig()
	config.OAuth.ClientSecret = ""
	if err := config.validateServer(); err == nil || !strings.Contains(err.Error(), "INTROSPECTION_CLIENT_SECRET") {
		t.Fatalf("validateServer() error = %v", err)
	}
}

func TestLoadDefaultsToDedicatedLocalPorts(t *testing.T) {
	t.Setenv("SERVICE_PORT", "")
	t.Setenv("DATABASE_DSN", "")
	t.Setenv("DATABASE_MIGRATION_URL", "")

	config, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if config.Service.Port != ":8083" {
		t.Fatalf("Load() port = %q, want :8083", config.Service.Port)
	}
	if !strings.Contains(config.Database.DSN, ":3308") || !strings.Contains(config.Database.MigrationURL, ":3308") {
		t.Fatalf("Load() database config = %#v, want port 3308", config.Database)
	}
}

func validConfig() Config {
	return Config{
		Service:  ServiceConfig{Name: "book", Port: ":8080", AllowedOrigin: "http://localhost:3000"},
		Database: DatabaseConfig{DSN: "book:password@tcp(localhost:3308)/book", MigrationURL: "mysql://book:password@tcp(localhost:3308)/book"},
		OAuth: OAuthConfig{
			IntrospectionURL: "http://auth-service:8081/oauth/introspect", ClientID: "book-service",
			ClientSecret: "secret", Issuer: "http://localhost:8081", Audience: "book-service",
			ServiceClientID: "transaction-service", Timeout: 2 * time.Second,
		},
		Rabbit: RabbitConfig{URL: "amqp://localhost", ConfirmTimeout: 5 * time.Second},
	}
}
