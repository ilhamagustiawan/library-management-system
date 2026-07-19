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
	OAuth    OAuthConfig
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
	BcryptCost          int
	SessionTTL          time.Duration
	SessionCookieName   string
	SessionCookieDomain string
	SessionCookieSecure bool
	RateLimitMax        int
	RateLimitWindow     time.Duration
}

type OAuthConfig struct {
	Issuer          string
	LoginURL        string
	CodeTTL         time.Duration
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	SupportedScopes []string
	JWTSigningKey   []byte
}

func loadDotEnv() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load .env: %w", err)
	}
	return nil
}

func Load() (Config, error) {
	if err := loadDotEnv(); err != nil {
		return Config{}, err
	}

	cfg := Config{
		Service: ServiceConfig{
			Name: env("SERVICE_NAME", "library-auth-service"), Environment: env("SERVICE_ENV", "development"),
			Port: env("SERVICE_PORT", ":8081"), AllowedOrigin: env("ALLOWED_ORIGIN", "http://localhost:3000"),
		},
		Database: DatabaseConfig{
			DSN:          env("DATABASE_DSN", "auth:auth_password@tcp(localhost:3306)/auth?parseTime=true&loc=UTC&charset=utf8mb4"),
			MigrationURL: env("DATABASE_MIGRATION_URL", "mysql://auth:auth_password@tcp(localhost:3306)/auth"),
			MaxOpenConns: envInt("DATABASE_MAX_OPEN_CONNS", 20), MaxIdleConns: envInt("DATABASE_MAX_IDLE_CONNS", 10),
			ConnMaxLifetime: envDuration("DATABASE_CONN_MAX_LIFETIME", 30*time.Minute),
			ConnMaxIdleTime: envDuration("DATABASE_CONN_MAX_IDLE_TIME", 5*time.Minute),
		},
		Auth: AuthConfig{
			BcryptCost: envInt("BCRYPT_COST", 12), SessionTTL: envDuration("SESSION_TTL", 24*time.Hour),
			SessionCookieName:   env("SESSION_COOKIE_NAME", "lms_session"),
			SessionCookieDomain: os.Getenv("SESSION_COOKIE_DOMAIN"),
			SessionCookieSecure: envBool("SESSION_COOKIE_SECURE", false),
			RateLimitMax:        envInt("AUTH_RATE_LIMIT_MAX", 10), RateLimitWindow: envDuration("AUTH_RATE_LIMIT_WINDOW", 15*time.Minute),
		},
		OAuth: OAuthConfig{
			Issuer: env("OAUTH_ISSUER", "http://localhost:8000"), LoginURL: env("LOGIN_URL", "http://localhost:3000/login"),
			CodeTTL: envDuration("OAUTH_CODE_TTL", 5*time.Minute), AccessTokenTTL: envDuration("OAUTH_ACCESS_TOKEN_TTL", 15*time.Minute),
			RefreshTokenTTL: envDuration("OAUTH_REFRESH_TOKEN_TTL", 7*24*time.Hour),
			SupportedScopes: strings.Fields(env("OAUTH_CLIENT_SCOPES", "loans:borrow:self loans:return:self transactions:read:self books:read transactions:read:any loans:return:any fines:manage books:manage identities:create book-stock:read book-stock:reserve book-stock:release")),
			JWTSigningKey:   []byte(os.Getenv("OAUTH_JWT_SIGNING_KEY")),
		},
	}
	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func LoadServer() (Config, error) {
	cfg, err := Load()
	if err != nil {
		return Config{}, err
	}
	if len(cfg.OAuth.JWTSigningKey) < 32 {
		return Config{}, fmt.Errorf("OAUTH_JWT_SIGNING_KEY must contain at least 32 bytes")
	}
	return cfg, nil
}

func (c Config) validate() error {
	if strings.TrimSpace(c.Database.DSN) == "" || !strings.HasPrefix(c.Database.MigrationURL, "mysql://") {
		return fmt.Errorf("DATABASE_DSN and a mysql:// DATABASE_MIGRATION_URL are required")
	}
	if c.Auth.BcryptCost < 10 || c.Auth.BcryptCost > 14 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 14")
	}
	if c.Auth.SessionTTL <= 0 || c.OAuth.CodeTTL <= 0 || c.OAuth.AccessTokenTTL <= 0 || c.OAuth.RefreshTokenTTL <= 0 {
		return fmt.Errorf("authentication and OAuth TTL values must be positive")
	}
	if len(c.OAuth.SupportedScopes) == 0 {
		return fmt.Errorf("OAUTH_CLIENT_SCOPES must contain at least one scope")
	}
	issuer, err := url.Parse(c.OAuth.Issuer)
	if err != nil || issuer.Scheme == "" || issuer.Host == "" || issuer.Path != "" {
		return fmt.Errorf("OAUTH_ISSUER must be an absolute origin without a path")
	}
	login, err := url.Parse(c.OAuth.LoginURL)
	if err != nil || login.Scheme == "" || login.Host == "" {
		return fmt.Errorf("LOGIN_URL must be an absolute URL")
	}
	origin, err := url.Parse(c.Service.AllowedOrigin)
	if err != nil || origin.Scheme == "" || origin.Host == "" || origin.Path != "" {
		return fmt.Errorf("ALLOWED_ORIGIN must be an absolute origin without a path")
	}
	if c.Service.Environment == "production" && (issuer.Scheme != "https" || !c.Auth.SessionCookieSecure) {
		return fmt.Errorf("production requires an HTTPS issuer and secure session cookies")
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

func envBool(key string, fallback bool) bool {
	value, err := strconv.ParseBool(os.Getenv(key))
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
