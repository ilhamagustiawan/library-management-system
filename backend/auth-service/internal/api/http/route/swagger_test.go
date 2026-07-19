package route

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterSwaggerServesUIAndDocumentsRoutes(t *testing.T) {
	app := fiber.New()
	registerSwagger(app)

	uiResponse, err := app.Test(newRequest(t, http.MethodGet, "/api/v1/docs/auth/swagger"))
	if err != nil {
		t.Fatalf("request Swagger UI: %v", err)
	}
	defer uiResponse.Body.Close()
	if uiResponse.StatusCode != http.StatusOK {
		t.Fatalf("Swagger UI status = %d, want 200", uiResponse.StatusCode)
	}
	uiBody, err := io.ReadAll(uiResponse.Body)
	if err != nil {
		t.Fatalf("read Swagger UI: %v", err)
	}
	if !strings.Contains(string(uiBody), "swagger-ui-dist@5.32.9") {
		t.Fatal("Swagger UI assets are not version-pinned")
	}
	if !strings.Contains(string(uiBody), `url: '\/api\/v1\/docs\/auth\/swagger.json'`) {
		t.Fatal("Swagger UI spec URL is not proxy-safe")
	}

	specResponse, err := app.Test(newRequest(t, http.MethodGet, "/api/v1/docs/auth/swagger.json"))
	if err != nil {
		t.Fatalf("request Swagger spec: %v", err)
	}
	defer specResponse.Body.Close()
	if specResponse.StatusCode != http.StatusOK {
		t.Fatalf("Swagger spec status = %d, want 200", specResponse.StatusCode)
	}

	var spec struct {
		Swagger string                     `json:"swagger"`
		Paths   map[string]json.RawMessage `json:"paths"`
	}
	if err := json.NewDecoder(specResponse.Body).Decode(&spec); err != nil {
		t.Fatalf("decode Swagger spec: %v", err)
	}
	if spec.Swagger != "2.0" {
		t.Fatalf("Swagger version = %q, want 2.0", spec.Swagger)
	}
	for _, path := range []string{
		"/.well-known/oauth-authorization-server",
		"/api/v1/auth/login",
		"/api/v1/auth/logout",
		"/api/v1/auth/me",
		"/api/v1/oauth/userinfo",
		"/health/liveness",
		"/health/readiness",
		"/oauth/authorize",
		"/oauth/introspect",
		"/oauth/token",
		"/internal/identities",
	} {
		if _, ok := spec.Paths[path]; !ok {
			t.Errorf("Swagger spec missing %s", path)
		}
	}
}

func newRequest(t *testing.T, method, target string) *http.Request {
	t.Helper()
	request, err := http.NewRequest(method, target, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	return request
}
