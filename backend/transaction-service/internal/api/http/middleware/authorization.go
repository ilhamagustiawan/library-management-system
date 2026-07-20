package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/domain/errs"
)

const credentialLocal = "credential"

type Credential struct {
	Subject string
	Scopes  map[string]struct{}
}

// RequireScope trusts headers created by Kong on the private backend network.
// The service must never publish this listener directly to an untrusted network.
func RequireScope(scope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		subject := strings.TrimSpace(c.Get("X-Credential-Sub"))
		if subject == "" {
			return response.Write(c, http.StatusUnauthorized, errs.CodeInvalidToken, "authentication required")
		}
		scopes := make(map[string]struct{})
		for _, value := range strings.Fields(c.Get("X-Credential-Scope")) {
			scopes[value] = struct{}{}
		}
		if _, allowed := scopes[scope]; !allowed {
			return response.Write(c, http.StatusForbidden, errs.CodeForbidden, "required OAuth scope is missing")
		}
		c.Locals(credentialLocal, Credential{Subject: subject, Scopes: scopes})
		return c.Next()
	}
}

func Current(c *fiber.Ctx) (Credential, bool) {
	credential, ok := c.Locals(credentialLocal).(Credential)
	return credential, ok
}
