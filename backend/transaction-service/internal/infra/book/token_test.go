package book

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOAuthTokenSourceCachesUsableToken(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		requests++
		clientID, secret, ok := request.BasicAuth()
		if !ok || clientID != "transaction-service" || secret != "secret" {
			t.Fatalf("unexpected client authentication")
		}
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"access_token":"service-token","token_type":"Bearer","expires_in":300,"scope":"book-stock:reserve book-stock:release"}`))
	}))
	defer server.Close()
	source, err := NewOAuthTokenSource(OAuthConfig{TokenURL: server.URL, ClientID: "transaction-service", ClientSecret: "secret"})
	if err != nil {
		t.Fatalf("NewOAuthTokenSource() error = %v", err)
	}
	for range 2 {
		token, err := source.Token(context.Background())
		if err != nil || token != "service-token" {
			t.Fatalf("Token() = %q, %v", token, err)
		}
	}
	if requests != 1 {
		t.Fatalf("token requests = %d, want 1", requests)
	}
}
