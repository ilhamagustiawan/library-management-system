package oauth

import (
	"context"
	"encoding/base64"
	"net/url"
	"strings"
	"testing"
	"time"

	oauth2lib "github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/golang-jwt/jwt/v5"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

func TestAccessTokenGeneratorIssuesSignedResourceJWTAndOpaqueRefreshToken(t *testing.T) {
	key := []byte("test-signing-key-with-at-least-32-bytes")
	generator, err := NewAccessTokenGenerator(AccessTokenConfig{
		Issuer: "https://auth.example", SigningKey: key,
	})
	if err != nil {
		t.Fatalf("NewAccessTokenGenerator() error = %v", err)
	}

	createdAt := time.Now().UTC().Truncate(time.Second)
	tokenInfo := models.NewToken()
	tokenInfo.SetClientID("member-nextjs-web")
	tokenInfo.SetUserID("user-123")
	tokenInfo.SetScope("books:read loans:borrow:self")
	tokenInfo.SetAccessCreateAt(createdAt)
	tokenInfo.SetAccessExpiresIn(15 * time.Minute)
	tokenInfo.SetExtension(url.Values{"audience": {"library-api"}, "role": {"member"}})

	access, refresh, err := generator.Token(context.Background(), &oauth2lib.GenerateBasic{
		Client: &models.Client{ID: "member-nextjs-web"}, UserID: "user-123", CreateAt: createdAt, TokenInfo: tokenInfo,
	}, true)
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}

	claims := &AccessTokenClaims{}
	parsed, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			t.Fatalf("signing method = %s, want HS256", token.Method.Alg())
		}
		return key, nil
	}, jwt.WithAudience("library-api"), jwt.WithIssuer("https://auth.example"), jwt.WithExpirationRequired())
	if err != nil || !parsed.Valid {
		t.Fatalf("parse access token = (%v, %v), want valid JWT", parsed.Valid, err)
	}
	if claims.Subject != "user-123" || claims.ClientID != "member-nextjs-web" || claims.Scope != "books:read loans:borrow:self" {
		t.Fatalf("claims = %#v, want resource owner, client, and scope", claims)
	}
	if claims.Role != entity.RoleMember {
		t.Fatalf("role = %q, want member", claims.Role)
	}
	if claims.IssuedAt == nil || !claims.IssuedAt.Time.Equal(createdAt) {
		t.Fatalf("iat = %v, want %v", claims.IssuedAt, createdAt)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.Time.Equal(createdAt.Add(15*time.Minute)) {
		t.Fatalf("exp = %v, want %v", claims.ExpiresAt, createdAt.Add(15*time.Minute))
	}
	if claims.ID == "" {
		t.Fatal("jti is empty")
	}
	if strings.Count(refresh, ".") == 2 {
		t.Fatalf("refresh token = %q, want opaque token", refresh)
	}
	decodedRefresh, err := base64.RawURLEncoding.DecodeString(refresh)
	if err != nil || len(decodedRefresh) != 32 {
		t.Fatalf("refresh token decoded length = %d, err = %v; want 32", len(decodedRefresh), err)
	}
}

func TestAccessTokenGeneratorUsesClientSubjectAndOmitsRoleAndRefreshForClientCredentials(t *testing.T) {
	generator, err := NewAccessTokenGenerator(AccessTokenConfig{
		Issuer: "https://auth.example", SigningKey: []byte("test-signing-key-with-at-least-32-bytes"),
	})
	if err != nil {
		t.Fatalf("NewAccessTokenGenerator() error = %v", err)
	}
	tokenInfo := models.NewToken()
	tokenInfo.SetClientID("inventory-worker")
	tokenInfo.SetScope("books:read")
	tokenInfo.SetAccessCreateAt(time.Now())
	tokenInfo.SetAccessExpiresIn(15 * time.Minute)
	tokenInfo.SetExtension(url.Values{"audience": {"book-service"}})

	access, refresh, err := generator.Token(context.Background(), &oauth2lib.GenerateBasic{
		Client: &models.Client{ID: "inventory-worker"}, TokenInfo: tokenInfo,
	}, false)
	if err != nil {
		t.Fatalf("Token() error = %v", err)
	}
	claims := &AccessTokenClaims{}
	_, _, err = jwt.NewParser().ParseUnverified(access, claims)
	if err != nil {
		t.Fatalf("parse access token: %v", err)
	}
	if claims.Subject != "inventory-worker" || claims.Role != "" || refresh != "" {
		t.Fatalf("subject, role, refresh = %q, %q, %q", claims.Subject, claims.Role, refresh)
	}
}

func TestNewAccessTokenGeneratorRejectsUnsafeConfig(t *testing.T) {
	tests := []struct {
		name   string
		config AccessTokenConfig
	}{
		{name: "missing issuer", config: AccessTokenConfig{SigningKey: []byte(strings.Repeat("k", 32))}},
		{name: "short signing key", config: AccessTokenConfig{Issuer: "https://auth.example", SigningKey: []byte("short")}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := NewAccessTokenGenerator(test.config); err == nil {
				t.Fatal("NewAccessTokenGenerator() error = nil, want configuration error")
			}
		})
	}
}
