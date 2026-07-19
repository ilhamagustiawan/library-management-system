package oauth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	oauth2lib "github.com/go-oauth2/oauth2/v4"
	"github.com/golang-jwt/jwt/v5"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

const minimumHMACKeyBytes = 32
const tokenAudienceExtension = "audience"
const tokenRoleExtension = "role"

type AccessTokenConfig struct {
	Issuer     string
	SigningKey []byte
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	ClientID string      `json:"client_id"`
	Scope    string      `json:"scope"`
	Role     entity.Role `json:"role,omitempty"`
}

type AccessTokenGenerator struct {
	issuer     string
	signingKey []byte
}

func NewAccessTokenGenerator(config AccessTokenConfig) (*AccessTokenGenerator, error) {
	issuer := strings.TrimRight(strings.TrimSpace(config.Issuer), "/")
	if issuer == "" {
		return nil, fmt.Errorf("JWT issuer is required")
	}
	if len(config.SigningKey) < minimumHMACKeyBytes {
		return nil, fmt.Errorf("JWT signing key must contain at least %d bytes", minimumHMACKeyBytes)
	}
	return &AccessTokenGenerator{
		issuer: issuer, signingKey: append([]byte(nil), config.SigningKey...),
	}, nil
}

func (g *AccessTokenGenerator) Token(
	_ context.Context,
	data *oauth2lib.GenerateBasic,
	generateRefresh bool,
) (string, string, error) {
	if data == nil || data.TokenInfo == nil || data.Client == nil {
		return "", "", fmt.Errorf("generate JWT access token: OAuth token data is incomplete")
	}

	issuedAt := data.TokenInfo.GetAccessCreateAt()
	expiresIn := data.TokenInfo.GetAccessExpiresIn()
	if issuedAt.IsZero() || expiresIn <= 0 {
		return "", "", fmt.Errorf("generate JWT access token: positive access-token lifetime is required")
	}
	extendable, ok := data.TokenInfo.(oauth2lib.ExtendableTokenInfo)
	if !ok {
		return "", "", fmt.Errorf("generate JWT access token: OAuth token context is unavailable")
	}
	audience := strings.TrimSpace(extendable.GetExtension().Get(tokenAudienceExtension))
	if audience == "" {
		return "", "", fmt.Errorf("generate JWT access token: audience is required")
	}
	subject := data.UserID
	role := entity.Role(extendable.GetExtension().Get(tokenRoleExtension))
	if subject == "" {
		subject = data.Client.GetID()
		if role != "" {
			return "", "", fmt.Errorf("generate JWT access token: service token cannot contain role")
		}
	} else if _, err := entity.ParseRole(role.String()); err != nil {
		return "", "", fmt.Errorf("generate JWT access token: user role is required: %w", err)
	}
	jti, err := randomToken(16)
	if err != nil {
		return "", "", fmt.Errorf("generate JWT ID: %w", err)
	}
	claims := AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: g.issuer, Subject: subject, Audience: jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(expiresIn)), IssuedAt: jwt.NewNumericDate(issuedAt), ID: jti,
		},
		ClientID: data.Client.GetID(), Scope: data.TokenInfo.GetScope(), Role: role,
	}
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(g.signingKey)
	if err != nil {
		return "", "", fmt.Errorf("sign JWT access token: %w", err)
	}
	if !generateRefresh {
		return access, "", nil
	}
	refresh, err := randomToken(32)
	if err != nil {
		return "", "", fmt.Errorf("generate refresh token: %w", err)
	}
	return access, refresh, nil
}

func randomToken(size int) (string, error) {
	buffer := make([]byte, size)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
