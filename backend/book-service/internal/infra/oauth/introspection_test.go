package oauth

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }

func TestClientIntrospectsBearerToken(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(request *http.Request) (*http.Response, error) {
		clientID, secret, ok := request.BasicAuth()
		if !ok || clientID != "book-service" || secret != "secret" {
			t.Fatal("missing introspection client credentials")
		}
		if request.FormValue("token") != "access-token" {
			t.Fatalf("token = %q", request.FormValue("token"))
		}
		return &http.Response{
			StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{"active":true,"client_id":"transaction-service","sub":"transaction-service","scope":"book-stock:read","token_type":"Bearer","iss":"http://auth-service:8081","iat":1,"exp":9999999999,"aud":["book-service"]}`)),
		}, nil
	})}

	client := newClient(Config{
		URL: "http://auth-service:8081/oauth/introspect", ClientID: "book-service", ClientSecret: "secret", Timeout: time.Second,
	}, httpClient)
	principal, err := client.Introspect(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("Introspect() error = %v", err)
	}
	if !principal.Active || principal.Subject != "transaction-service" || principal.Audience[0] != "book-service" {
		t.Fatalf("principal = %#v", principal)
	}
}

func TestClientRejectsOversizedResponse(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}},
			Body: io.NopCloser(strings.NewReader(strings.Repeat("x", 70_000))),
		}, nil
	})}
	client := newClient(Config{
		URL: "http://auth-service:8081/oauth/introspect", ClientID: "book-service", ClientSecret: "secret", Timeout: time.Second,
	}, httpClient)
	if _, err := client.Introspect(context.Background(), "access-token"); err == nil {
		t.Fatal("Introspect() error = nil, want oversized response error")
	}
}

func TestClientRejectsTrailingJSON(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusOK, Header: http.Header{"Content-Type": {"application/json"}},
			Body: io.NopCloser(strings.NewReader(`{"active":false}{"active":true}`)),
		}, nil
	})}
	client := newClient(Config{
		URL: "http://auth-service:8081/oauth/introspect", ClientID: "book-service", ClientSecret: "secret", Timeout: time.Second,
	}, httpClient)
	if _, err := client.Introspect(context.Background(), "access-token"); err == nil {
		t.Fatal("Introspect() error = nil, want trailing JSON error")
	}
}
