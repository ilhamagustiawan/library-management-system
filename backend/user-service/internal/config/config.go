package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Service  ServiceConfig
	Database DatabaseConfig
	Auth     AuthConfig
	RabbitMQ RabbitMQConfig
	Outbox   OutboxConfig
	Rate     RateConfig
}

type ServiceConfig struct {
	Name          string
	Environment   string
	Port          string
	AllowedOrigin string
}

type DatabaseConfig struct {
	DSN             string
	MigrationURL    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type AuthConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	Scope        string
	Timeout      time.Duration
	Attempts     int
}

type RabbitMQConfig struct {
	URL            string
	Exchange       string
	RoutingKey     string
	Queue          string
	ConfirmTimeout time.Duration
}

type OutboxConfig struct {
	BatchSize    int
	Lease        time.Duration
	PollInterval time.Duration
	BaseRetry    time.Duration
	MaxRetry     time.Duration
}

type RateConfig struct {
	Max    int
	Window time.Duration
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	config := Config{
		Service: ServiceConfig{
			Name: env("SERVICE_NAME", "library-user-service"), Environment: env("SERVICE_ENV", "development"),
			Port: env("SERVICE_PORT", ":8082"), AllowedOrigin: env("ALLOWED_ORIGIN", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			DSN:          env("DATABASE_DSN", "users:users_password@tcp(127.0.0.1:3307)/users?parseTime=true&loc=UTC&charset=utf8mb4"),
			MigrationURL: env("DATABASE_MIGRATION_URL", "mysql://users:users_password@tcp(127.0.0.1:3307)/users"),
			MaxOpenConns: envInt("DATABASE_MAX_OPEN_CONNS", 20), MaxIdleConns: envInt("DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: envDuration("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: envDuration("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			BaseURL: env("AUTH_SERVICE_URL", "http://127.0.0.1:8081"), ClientID: env("AUTH_CLIENT_ID", "user-service"),
			ClientSecret: os.Getenv("AUTH_CLIENT_SECRET"), Scope: env("AUTH_SCOPE", "identities:create"),
			Timeout: envDuration("AUTH_REQUEST_TIMEOUT", 5*time.Second), Attempts: envInt("AUTH_REQUEST_ATTEMPTS", 3),
		},
		RabbitMQ: RabbitMQConfig{
			URL:      env("RABBITMQ_URL", "amqp://library:library_password@127.0.0.1:5672/"),
			Exchange: env("RABBITMQ_EXCHANGE", "library.events"), RoutingKey: env("RABBITMQ_ROUTING_KEY", "user.registered.v1"),
			Queue: env("RABBITMQ_QUEUE", "library.user-registered.v1"), ConfirmTimeout: envDuration("RABBITMQ_CONFIRM_TIMEOUT", 5*time.Second),
		},
		Outbox: OutboxConfig{
			BatchSize: envInt("OUTBOX_BATCH_SIZE", 50), Lease: envDuration("OUTBOX_LEASE", 30*time.Second),
			PollInterval: envDuration("OUTBOX_POLL_INTERVAL", 500*time.Millisecond),
			BaseRetry:    envDuration("OUTBOX_BASE_RETRY", time.Second), MaxRetry: envDuration("OUTBOX_MAX_RETRY", time.Minute),
		},
		Rate: RateConfig{Max: envInt("REGISTRATION_RATE_LIMIT_MAX", 10), Window: envDuration("REGISTRATION_RATE_LIMIT_WINDOW", 15*time.Minute)},
	}
	if err := config.validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Database.DSN) == "" || !strings.HasPrefix(c.Database.MigrationURL, "mysql://") {
		return fmt.Errorf("DATABASE_DSN and a mysql:// DATABASE_MIGRATION_URL are required")
	}
	if len(c.Auth.ClientSecret) < 32 {
		return fmt.Errorf("AUTH_CLIENT_SECRET must contain at least 32 bytes")
	}
	if c.Auth.ClientID != "user-service" || c.Auth.Scope != "identities:create" {
		return fmt.Errorf("AUTH_CLIENT_ID must be user-service and AUTH_SCOPE must be identities:create")
	}
	authURL, err := absoluteOrigin(c.Auth.BaseURL, "AUTH_SERVICE_URL")
	if err != nil {
		return err
	}
	origin, err := absoluteOrigin(c.Service.AllowedOrigin, "ALLOWED_ORIGIN")
	if err != nil {
		return err
	}
	rabbitURL, err := url.Parse(c.RabbitMQ.URL)
	if err != nil || rabbitURL.Host == "" || (rabbitURL.Scheme != "amqp" && rabbitURL.Scheme != "amqps") {
		return fmt.Errorf("RABBITMQ_URL must be an absolute amqp or amqps URL")
	}
	if c.Service.Environment == "production" && (authURL.Scheme != "https" || origin.Scheme != "https" || rabbitURL.Scheme != "amqps") {
		return fmt.Errorf("production requires HTTPS Auth/origin URLs and AMQPS RabbitMQ URL")
	}
	if c.Auth.Timeout <= 0 || c.Auth.Attempts < 1 || c.Auth.Attempts > 5 || c.Rate.Max <= 0 || c.Rate.Window <= 0 {
		return fmt.Errorf("Auth timeout/attempts and registration rate limits must be positive")
	}
	if c.Outbox.BatchSize <= 0 || c.Outbox.BatchSize > 500 || c.Outbox.Lease <= 0 || c.Outbox.PollInterval <= 0 || c.Outbox.BaseRetry <= 0 || c.Outbox.MaxRetry < c.Outbox.BaseRetry {
		return fmt.Errorf("outbox settings are invalid")
	}
	return nil
}

func absoluteOrigin(raw, field string) (*url.URL, error) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" || (parsed.Path != "" && parsed.Path != "/") || parsed.RawQuery != "" || parsed.Fragment != "" || parsed.User != nil {
		return nil, fmt.Errorf("%s must be an absolute origin", field)
	}
	return parsed, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value, err := time.ParseDuration(os.Getenv(key))
	if err != nil {
		return fallback
	}
	return value
}
