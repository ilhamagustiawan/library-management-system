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
	OAuth    OAuthConfig
	Rabbit   RabbitConfig
}

type RabbitConfig struct {
	URL            string
	Exchange       string
	DeadExchange   string
	ReturnQueue    string
	AckQueue       string
	DeadQueue      string
	ConfirmTimeout time.Duration
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

type OAuthConfig struct {
	IntrospectionURL string
	ClientID         string
	ClientSecret     string
	Issuer           string
	Audience         string
	ServiceClientID  string
	Timeout          time.Duration
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}
	config := Config{
		Service: ServiceConfig{
			Name: env("SERVICE_NAME", "library-book-service"), Environment: env("SERVICE_ENV", "development"),
			Port: env("SERVICE_PORT", ":8083"), AllowedOrigin: env("ALLOWED_ORIGIN", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			DSN:          env("DATABASE_DSN", "book:book_password@tcp(localhost:3308)/book?parseTime=true&loc=UTC&charset=utf8mb4"),
			MigrationURL: env("DATABASE_MIGRATION_URL", "mysql://book:book_password@tcp(localhost:3308)/book"),
			MaxOpenConns: envInt("DATABASE_MAX_OPEN_CONNS", 20), MaxIdleConns: envInt("DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: envDuration("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: envDuration("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		OAuth: OAuthConfig{
			IntrospectionURL: env("AUTH_INTROSPECTION_URL", "http://localhost:8081/oauth/introspect"),
			ClientID:         env("INTROSPECTION_CLIENT_ID", "book-service"), ClientSecret: os.Getenv("INTROSPECTION_CLIENT_SECRET"),
			Issuer: env("OAUTH_ISSUER", "http://localhost:8000"), Audience: env("OAUTH_AUDIENCE", "book-service"),
			ServiceClientID: env("TRANSACTION_SERVICE_CLIENT_ID", "transaction-service"),
			Timeout:         envDuration("AUTH_INTROSPECTION_TIMEOUT", 2*time.Second),
		},
		Rabbit: RabbitConfig{
			URL:      env("RABBITMQ_URL", "amqp://library:library_password@localhost:5672/"),
			Exchange: env("RABBITMQ_EXCHANGE", "library.events"), DeadExchange: env("RABBITMQ_DEAD_EXCHANGE", "library.events.dlx"),
			ReturnQueue: env("RABBITMQ_BOOK_RETURN_QUEUE", "book-service.loan-returned.v1"),
			AckQueue:    env("RABBITMQ_STOCK_ACK_QUEUE", "transaction-service.book-stock-updated.v1"),
			DeadQueue:   env("RABBITMQ_BOOK_DEAD_QUEUE", "book-service.dead-letter"), ConfirmTimeout: envDuration("RABBITMQ_CONFIRM_TIMEOUT", 5*time.Second),
		},
	}
	if err := config.validate(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func LoadServer() (Config, error) {
	config, err := Load()
	if err != nil {
		return Config{}, err
	}
	if err := config.validateServer(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Database.DSN) == "" || !strings.HasPrefix(c.Database.MigrationURL, "mysql://") {
		return fmt.Errorf("DATABASE_DSN and a mysql:// DATABASE_MIGRATION_URL are required")
	}
	if c.Database.MaxOpenConns < 1 || c.Database.MaxIdleConns < 0 || c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		return fmt.Errorf("database connection limits are invalid")
	}
	if c.Database.ConnMaxLifetime <= 0 || c.Database.ConnMaxIdleTime <= 0 {
		return fmt.Errorf("database connection lifetimes must be positive")
	}
	origin, err := url.Parse(c.Service.AllowedOrigin)
	if err != nil || origin.Scheme == "" || origin.Host == "" || origin.Path != "" || origin.RawQuery != "" || origin.Fragment != "" {
		return fmt.Errorf("ALLOWED_ORIGIN must be an absolute origin without a path")
	}
	return nil
}

func (c Config) validateServer() error {
	if strings.TrimSpace(c.OAuth.ClientSecret) == "" {
		return fmt.Errorf("INTROSPECTION_CLIENT_SECRET is required")
	}
	if strings.TrimSpace(c.OAuth.ClientID) == "" || strings.TrimSpace(c.OAuth.Audience) == "" || strings.TrimSpace(c.OAuth.ServiceClientID) == "" {
		return fmt.Errorf("introspection client, audience, and Transaction Service client ID are required")
	}
	endpoint, err := url.Parse(c.OAuth.IntrospectionURL)
	if err != nil || endpoint.Scheme == "" || endpoint.Host == "" || endpoint.User != nil || endpoint.RawQuery != "" || endpoint.Fragment != "" {
		return fmt.Errorf("AUTH_INTROSPECTION_URL must be an absolute URL without credentials, query, or fragment")
	}
	issuer, err := url.Parse(c.OAuth.Issuer)
	if err != nil || issuer.Scheme == "" || issuer.Host == "" || issuer.Path != "" || issuer.RawQuery != "" || issuer.Fragment != "" {
		return fmt.Errorf("OAUTH_ISSUER must be an absolute origin without a path")
	}
	if c.OAuth.Timeout <= 0 || c.OAuth.Timeout > 10*time.Second {
		return fmt.Errorf("AUTH_INTROSPECTION_TIMEOUT must be between zero and 10s")
	}
	if strings.TrimSpace(c.Rabbit.URL) == "" || c.Rabbit.ConfirmTimeout <= 0 {
		return fmt.Errorf("RabbitMQ URL and positive confirm timeout are required")
	}
	return nil
}

func loadDotEnv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	return nil
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
