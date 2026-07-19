package middleware

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRequireScopeAcceptsCredentialHeaders(t *testing.T) {
	app := fiber.New()
	app.Get("/", RequireScope("books:read"), func(c *fiber.Ctx) error { return c.SendStatus(http.StatusNoContent) })
	request, _ := http.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("X-Credential-Sub", "member-1")
	request.Header.Set("X-Credential-Scope", "books:read")

	response, err := app.Test(request)

	if err != nil || response.StatusCode != http.StatusNoContent {
		t.Fatalf("response = %v, error = %v", response, err)
	}
}
