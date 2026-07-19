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
	Book     BookConfig
	Rabbit   RabbitConfig
	Loan     LoanConfig
}

type ServiceConfig struct {
	Name, Environment, Port, AllowedOrigin string
}
type DatabaseConfig struct {
	DSN, MigrationURL                string
	MaxOpenConns, MaxIdleConns       int
	ConnMaxLifetime, ConnMaxIdleTime time.Duration
}
type OAuthConfig struct{ TokenURL, ClientID, ClientSecret string }
type BookConfig struct {
	URL     string
	Timeout time.Duration
}
type RabbitConfig struct {
	URL, Exchange, DeadExchange, BookReturnQueue, StockAckQueue, DeadLetterQueue string
	ConfirmTimeout                                                               time.Duration
}
type LoanConfig struct {
	Term, AckTimeout, PollInterval time.Duration
	DailyFineMinor                 int64
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("load .env: %w", err)
	}
	config := Config{
		Service:  ServiceConfig{Name: env("SERVICE_NAME", "library-transaction-service"), Environment: env("SERVICE_ENV", "development"), Port: env("SERVICE_PORT", ":8084"), AllowedOrigin: env("ALLOWED_ORIGIN", "http://localhost:3000")},
		Database: DatabaseConfig{DSN: env("DATABASE_DSN", "transactions:transactions_password@tcp(localhost:3309)/transactions?parseTime=true&loc=UTC&charset=utf8mb4"), MigrationURL: env("DATABASE_MIGRATION_URL", "mysql://transactions:transactions_password@tcp(localhost:3309)/transactions"), MaxOpenConns: envInt("DATABASE_MAX_OPEN_CONNS", 20), MaxIdleConns: envInt("DATABASE_MAX_IDLE_CONNS", 10), ConnMaxLifetime: envDuration("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute), ConnMaxIdleTime: envDuration("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute)},
		OAuth:    OAuthConfig{TokenURL: env("AUTH_TOKEN_URL", "http://localhost:8081/oauth/token"), ClientID: env("BOOK_SERVICE_CLIENT_ID", "transaction-service"), ClientSecret: os.Getenv("BOOK_SERVICE_CLIENT_SECRET")},
		Book:     BookConfig{URL: env("BOOK_SERVICE_URL", "http://localhost:8083"), Timeout: envDuration("BOOK_SERVICE_TIMEOUT", 2*time.Second)},
		Rabbit:   RabbitConfig{URL: env("RABBITMQ_URL", "amqp://library:library_password@localhost:5672/"), Exchange: env("RABBITMQ_EXCHANGE", "library.events"), DeadExchange: env("RABBITMQ_DEAD_EXCHANGE", "library.events.dlx"), BookReturnQueue: env("RABBITMQ_BOOK_RETURN_QUEUE", "book-service.loan-returned.v1"), StockAckQueue: env("RABBITMQ_STOCK_ACK_QUEUE", "transaction-service.book-stock-updated.v1"), DeadLetterQueue: env("RABBITMQ_DEAD_LETTER_QUEUE", "transaction-service.dead-letter"), ConfirmTimeout: envDuration("RABBITMQ_CONFIRM_TIMEOUT", 5*time.Second)},
		Loan:     LoanConfig{Term: envDuration("LOAN_TERM", 7*24*time.Hour), AckTimeout: envDuration("STOCK_ACK_TIMEOUT", 5*time.Second), PollInterval: envDuration("STOCK_ACK_POLL_INTERVAL", 100*time.Millisecond), DailyFineMinor: envInt64("FINE_DAILY_RATE_MINOR", 5000)},
	}
	if err := config.validateBase(); err != nil {
		return Config{}, err
	}
	return config, nil
}

func LoadServer() (Config, error) {
	config, err := Load()
	if err != nil {
		return Config{}, err
	}
	for name, value := range map[string]string{"BOOK_SERVICE_CLIENT_SECRET": config.OAuth.ClientSecret, "RABBITMQ_URL": config.Rabbit.URL} {
		if strings.TrimSpace(value) == "" {
			return Config{}, fmt.Errorf("%s is required", name)
		}
	}
	for name, raw := range map[string]string{"AUTH_TOKEN_URL": config.OAuth.TokenURL, "BOOK_SERVICE_URL": config.Book.URL} {
		parsed, err := url.Parse(raw)
		if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
			return Config{}, fmt.Errorf("%s must be an absolute HTTP(S) URL", name)
		}
	}
	return config, nil
}

func (c Config) validateBase() error {
	if strings.TrimSpace(c.Database.DSN) == "" || !strings.HasPrefix(c.Database.MigrationURL, "mysql://") {
		return fmt.Errorf("DATABASE_DSN and a mysql:// DATABASE_MIGRATION_URL are required")
	}
	if c.Loan.Term <= 0 || c.Loan.AckTimeout <= 0 || c.Loan.PollInterval <= 0 || c.Loan.DailyFineMinor <= 0 {
		return fmt.Errorf("loan timing and fine values must be positive")
	}
	if c.Book.Timeout <= 0 || c.Rabbit.ConfirmTimeout <= 0 {
		return fmt.Errorf("dependency timeouts must be positive")
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
func envInt64(key string, fallback int64) int64 {
	value, err := strconv.ParseInt(os.Getenv(key), 10, 64)
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
