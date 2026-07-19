package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	domainerrs "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/errs"
)

func RequireTrustedOrigin(allowedOrigin string) fiber.Handler {
	return trustedOrigin{allowedOrigin: allowedOrigin}.Handle
}

type trustedOrigin struct {
	allowedOrigin string
}

func (m trustedOrigin) Handle(c *fiber.Ctx) error {
	origin := c.Get(fiber.HeaderOrigin)
	if origin != "" && origin != m.allowedOrigin {
		return response.Write(c, http.StatusForbidden, domainerrs.CodeInvalidToken, "untrusted request origin")
	}
	return c.Next()
}
