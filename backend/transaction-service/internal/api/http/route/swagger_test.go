package route

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegisterSwaggerServesUIAndTransactionContract(t *testing.T) {
	app := fiber.New()
	registerSwagger(app)

	uiResponse, err := app.Test(swaggerRequest(t, "/api/v1/docs/transactions/swagger"))
	if err != nil {
		t.Fatalf("request Swagger UI: %v", err)
	}
	defer uiResponse.Body.Close()
	uiBody, err := io.ReadAll(uiResponse.Body)
	if err != nil {
		t.Fatalf("read Swagger UI: %v", err)
	}
	if uiResponse.StatusCode != http.StatusOK || !strings.Contains(string(uiBody), "swagger-ui-dist@5.32.9") {
		t.Fatalf("Swagger UI status = %d or assets are not version-pinned", uiResponse.StatusCode)
	}
	if !strings.Contains(string(uiBody), `url: '\/api\/v1\/docs\/transactions\/swagger.json'`) {
		t.Fatal("Swagger UI spec URL is not proxy-safe")
	}

	specResponse, err := app.Test(swaggerRequest(t, "/api/v1/docs/transactions/swagger.json"))
	if err != nil {
		t.Fatalf("request Swagger spec: %v", err)
	}
	defer specResponse.Body.Close()
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
		"/api/v1/transactions/admin",
		"/api/v1/transactions/admin/loans/{loanId}/return",
		"/api/v1/transactions/loans",
		"/api/v1/transactions/loans/{loanId}/return",
		"/api/v1/transactions/me",
	} {
		if _, ok := spec.Paths[path]; !ok {
			t.Errorf("Swagger spec missing %s", path)
		}
	}
}

func swaggerRequest(t *testing.T, target string) *http.Request {
	t.Helper()
	request, err := http.NewRequest(http.MethodGet, target, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	return request
}
