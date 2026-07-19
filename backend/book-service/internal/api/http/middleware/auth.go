package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/domain/errs"
	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/infra/oauth"
)

const (
	credentialScopeHeader = "X-Credential-Scope"
	credentialSubHeader   = "X-Credential-Sub"
)

type Introspector interface {
	Introspect(context.Context, string) (oauth.Principal, error)
}

type InternalPolicy struct {
	Issuer   string
	Audience string
	ClientID string
	Scope    string
	Now      func() time.Time
}

// TODO: Authenticate gateway-injected headers before exposing this service outside local development.
func RequireGatewayScopes(scopes ...string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err := authorizeGateway(
			ctx.Get(credentialScopeHeader), ctx.Get(credentialSubHeader), scopes,
		)
		if err != nil {
			return response.Error(ctx, err)
		}
		return ctx.Next()
	}
}

func RequireInternal(introspector Introspector, policy InternalPolicy) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token := bearerToken(ctx.Get(fiber.HeaderAuthorization))
		if token == "" {
			return response.Error(ctx, unauthorized("Bearer token is required"))
		}
		principal, err := introspector.Introspect(ctx.UserContext(), token)
		if err != nil {
			return response.Error(ctx, errs.New(
				http.StatusServiceUnavailable, errs.CodeAuthUnavailable,
				"authentication service is unavailable; stock remains unchanged; retry later", nil, err,
			))
		}
		if err := authorizeInternal(principal, policy); err != nil {
			return response.Error(ctx, err)
		}
		return ctx.Next()
	}
}

func authorizeGateway(scopes, subject string, required []string) error {
	if strings.TrimSpace(subject) == "" {
		return unauthorized("authentication is required")
	}
	if !hasAnyScope(scopes, required) {
		return forbidden("token lacks required book scope")
	}
	return nil
}

func authorizeInternal(principal oauth.Principal, policy InternalPolicy) error {
	now := time.Now
	if policy.Now != nil {
		now = policy.Now
	}
	if !principal.Active || principal.ExpiresAt <= now().Unix() {
		return unauthorized("active Bearer token is required")
	}
	if !strings.EqualFold(principal.TokenType, "Bearer") || principal.IssuedAt <= 0 ||
		principal.IssuedAt > principal.ExpiresAt || principal.IssuedAt > now().Add(30*time.Second).Unix() {
		return forbidden("token metadata is invalid")
	}
	if principal.Issuer != policy.Issuer || principal.ClientID != policy.ClientID || principal.Subject != policy.ClientID ||
		!contains(principal.Audience, policy.Audience) || !hasAnyScope(principal.Scope, []string{policy.Scope}) {
		return forbidden("token is not authorized for this stock operation")
	}
	return nil
}

func hasAnyScope(granted string, required []string) bool {
	set := make(map[string]struct{})
	for _, scope := range strings.Fields(granted) {
		set[scope] = struct{}{}
	}
	for _, scope := range required {
		if _, ok := set[scope]; ok {
			return true
		}
	}
	return false
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func bearerToken(header string) string {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return parts[1]
}

func unauthorized(message string) error {
	return errs.New(http.StatusUnauthorized, errs.CodeUnauthorized, message, nil, nil)
}

func forbidden(message string) error {
	return errs.New(http.StatusForbidden, errs.CodeForbidden, message, nil, nil)
}
