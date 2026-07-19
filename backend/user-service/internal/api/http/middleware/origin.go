package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/domain/errs"
)

func RequireTrustedOrigin(allowedOrigin string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		origin := c.Get(fiber.HeaderOrigin)
		if origin != "" && origin != allowedOrigin {
			return response.Write(c, http.StatusForbidden, errs.CodeValidation, "untrusted request origin; no registration data changed")
		}
		return c.Next()
	}
}
