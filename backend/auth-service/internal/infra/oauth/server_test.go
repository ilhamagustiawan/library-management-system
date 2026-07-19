package oauth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	oauth2lib "github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type fakeSessions struct{}

func (fakeSessions) AuthenticateSession(_ context.Context, token string) (*entity.User, error) {
	if token != "valid-session" {
		return nil, context.Canceled
	}
	return &entity.User{ID: "user-123"}, nil
}

type allowScopes struct {
	grants        map[string]map[oauth2lib.GrantType]bool
	introspectors map[string]bool
}

func (allowScopes) AllowsScopes(context.Context, string, string) (bool, error) { return true, nil }
func (p allowScopes) AllowsGrant(_ context.Context, clientID string, grant oauth2lib.GrantType) (bool, error) {
	return p.grants[clientID][grant], nil
}
func (p allowScopes) CanIntrospect(_ context.Context, clientID string) (bool, error) {
	return p.introspectors[clientID], nil
}

func TestAuthorizationServerRequiresS256PKCEAndExactRedirect(t *testing.T) {
	authServer := newTestServer(t)
	verifier := strings.Repeat("a", 64)
	challenge := s256(verifier)

	tests := []struct {
		name  string
		query url.Values
	}{
		{
			name:  "missing PKCE",
			query: authorizeQuery("http://client.example/callback", ""),
		},
		{
			name: "plain PKCE",
			query: func() url.Values {
				q := authorizeQuery("http://client.example/callback", verifier)
				q.Set("code_challenge_method", "plain")
				return q
			}(),
		},
		{
			name:  "non-exact redirect",
			query: authorizeQuery("http://client.example/callback/extra", challenge),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/oauth/authorize?"+test.query.Encode(), nil)
			request.AddCookie(&http.Cookie{Name: "lms_session", Value: "valid-session"})

			authServer.AuthorizeHandler().ServeHTTP(recorder, request)

			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400; body=%s", recorder.Code, recorder.Body.String())
			}
			if recorder.Header().Get("Location") != "" {
				t.Fatalf("unsafe redirect Location = %q", recorder.Header().Get("Location"))
			}
		})
	}
}

func TestAuthorizationServerCompletesCodeExchangeWithPKCE(t *testing.T) {
	authServer := newTestServer(t)
	verifier := strings.Repeat("a", 64)
	query := authorizeQuery("http://client.example/callback", s256(verifier))

	authorizeRecorder := httptest.NewRecorder()
	authorizeRequest := httptest.NewRequest(http.MethodGet, "/oauth/authorize?"+query.Encode(), nil)
	authorizeRequest.AddCookie(&http.Cookie{Name: "lms_session", Value: "valid-session"})
	authServer.AuthorizeHandler().ServeHTTP(authorizeRecorder, authorizeRequest)

	if authorizeRecorder.Code != http.StatusFound {
		t.Fatalf("authorize status = %d, want 302; body=%s", authorizeRecorder.Code, authorizeRecorder.Body.String())
	}
	location, err := url.Parse(authorizeRecorder.Header().Get("Location"))
	if err != nil {
		t.Fatalf("parse authorize redirect: %v", err)
	}
	code := location.Query().Get("code")
	if code == "" || location.Query().Get("state") != "state-123" {
		t.Fatalf("authorize redirect = %q, want code and original state", location.String())
	}

	form := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://client.example/callback"},
		"code_verifier": {verifier},
	}
	tokenRecorder := httptest.NewRecorder()
	tokenRequest := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	tokenRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	tokenRequest.SetBasicAuth("nextjs", "client-secret")
	authServer.TokenHandler().ServeHTTP(tokenRecorder, tokenRequest)

	if tokenRecorder.Code != http.StatusOK {
		t.Fatalf("token status = %d, want 200; body=%s", tokenRecorder.Code, tokenRecorder.Body.String())
	}
	var tokenResponse map[string]any
	if err := json.Unmarshal(tokenRecorder.Body.Bytes(), &tokenResponse); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	refreshToken, ok := tokenResponse["refresh_token"].(string)
	if tokenResponse["access_token"] == "" || !ok || refreshToken == "" {
		t.Fatalf("token response = %#v, want access and refresh tokens", tokenResponse)
	}

	refreshForm := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {refreshToken}}
	refreshRecorder := httptest.NewRecorder()
	refreshRequest := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(refreshForm.Encode()))
	refreshRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	refreshRequest.SetBasicAuth("other-client", "other-secret")
	authServer.TokenHandler().ServeHTTP(refreshRecorder, refreshRequest)
	if refreshRecorder.Code != http.StatusBadRequest {
		t.Fatalf("cross-client refresh status = %d, want 400; body=%s", refreshRecorder.Code, refreshRecorder.Body.String())
	}
}

func TestAuthorizationServerMetadataUsesConfiguredScopes(t *testing.T) {
	authServer := newTestServerWithScopes(t, []string{"catalog:read"})

	metadata := authServer.Metadata()
	if len(metadata.ScopesSupported) != 1 || metadata.ScopesSupported[0] != "catalog:read" {
		t.Fatalf("scopes_supported = %#v, want configured scopes", metadata.ScopesSupported)
	}
	if !slices.Contains(metadata.GrantTypesSupported, "client_credentials") {
		t.Fatalf("grant_types_supported = %#v, want client_credentials", metadata.GrantTypesSupported)
	}
	if metadata.IntrospectionEndpoint != "http://auth.example/oauth/introspect" {
		t.Fatalf("introspection_endpoint = %#v, want configured issuer endpoint", metadata.IntrospectionEndpoint)
	}
}

func TestAuthorizationServerIssuesClientCredentialsWithoutRefreshToken(t *testing.T) {
	authServer := newTestServer(t)
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {"library:read"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("machine-client", "machine-secret")

	authServer.TokenHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("token status = %d, want 200; body=%s", recorder.Code, recorder.Body.String())
	}
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode token response: %v", err)
	}
	if response["access_token"] == "" {
		t.Fatalf("token response = %#v, want access token", response)
	}
	if _, exists := response["refresh_token"]; exists {
		t.Fatalf("token response = %#v, client credentials must not issue refresh token", response)
	}
}

func TestAuthorizationServerRejectsClientCredentialsForAuthorizationCodeClient(t *testing.T) {
	authServer := newTestServer(t)
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {"library:read"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("nextjs", "client-secret")

	authServer.TokenHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("token status = %d, want 401; body=%s", recorder.Code, recorder.Body.String())
	}
}

func TestIntrospectionReturnsActiveClientTokenWithoutUserSubject(t *testing.T) {
	authServer := newTestServer(t)
	accessToken := issueClientCredentialsToken(t, authServer)
	form := url.Values{"token": {accessToken}, "token_type_hint": {"access_token"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/introspect", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("kong-gateway", "gateway-secret")

	authServer.IntrospectionHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("introspection status = %d, want 200; body=%s", recorder.Code, recorder.Body.String())
	}
	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode introspection response: %v", err)
	}
	if response["active"] != true || response["client_id"] != "machine-client" || response["scope"] != "library:read" {
		t.Fatalf("introspection response = %#v, want active machine token", response)
	}
	if _, exists := response["sub"]; exists {
		t.Fatalf("introspection response = %#v, machine token must not have user subject", response)
	}
	for _, field := range []string{"token_type", "iat", "exp", "iss"} {
		if _, exists := response[field]; !exists {
			t.Fatalf("introspection response = %#v, missing %s", response, field)
		}
	}
}

func TestIntrospectionReturnsOnlyInactiveForUnknownToken(t *testing.T) {
	authServer := newTestServer(t)
	form := url.Values{"token": {"unknown-token"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/introspect", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("kong-gateway", "gateway-secret")

	authServer.IntrospectionHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK || strings.TrimSpace(recorder.Body.String()) != `{"active":false}` {
		t.Fatalf("introspection = %d %s, want 200 active=false only", recorder.Code, recorder.Body.String())
	}
}

func TestIntrospectionReturnsUserSubjectForAuthorizationCodeToken(t *testing.T) {
	authServer := newTestServer(t)
	accessToken := issueAuthorizationCodeToken(t, authServer)
	form := url.Values{"token": {accessToken}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/introspect", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("kong-gateway", "gateway-secret")

	authServer.IntrospectionHandler().ServeHTTP(recorder, request)

	var response map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode introspection response: %v", err)
	}
	if recorder.Code != http.StatusOK || response["active"] != true || response["sub"] != "user-123" {
		t.Fatalf("introspection = %d %#v, want active user subject", recorder.Code, response)
	}
}

func TestIntrospectionRejectsUnauthorizedCaller(t *testing.T) {
	for _, caller := range []struct{ id, secret string }{
		{id: "kong-gateway", secret: "wrong-secret"},
		{id: "nextjs", secret: "client-secret"},
	} {
		t.Run(caller.id+caller.secret, func(t *testing.T) {
			authServer := newTestServer(t)
			form := url.Values{"token": {"unknown-token"}}
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/oauth/introspect", strings.NewReader(form.Encode()))
			request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			request.SetBasicAuth(caller.id, caller.secret)

			authServer.IntrospectionHandler().ServeHTTP(recorder, request)

			if recorder.Code != http.StatusUnauthorized {
				t.Fatalf("introspection status = %d, want 401; body=%s", recorder.Code, recorder.Body.String())
			}
		})
	}
}

func TestIntrospectionRequiresToken(t *testing.T) {
	authServer := newTestServer(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/introspect", strings.NewReader("token="))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("kong-gateway", "gateway-secret")

	authServer.IntrospectionHandler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("introspection status = %d, want 400; body=%s", recorder.Code, recorder.Body.String())
	}
}

func issueClientCredentialsToken(t *testing.T, authServer *AuthorizationServer) string {
	t.Helper()
	form := url.Values{"grant_type": {"client_credentials"}, "scope": {"library:read"}}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("machine-client", "machine-secret")
	authServer.TokenHandler().ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("issue client token status = %d; body=%s", recorder.Code, recorder.Body.String())
	}
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || response.AccessToken == "" {
		t.Fatalf("decode issued token: %v; body=%s", err, recorder.Body.String())
	}
	return response.AccessToken
}

func issueAuthorizationCodeToken(t *testing.T, authServer *AuthorizationServer) string {
	t.Helper()
	verifier := strings.Repeat("a", 64)
	authorizeRecorder := httptest.NewRecorder()
	authorizeRequest := httptest.NewRequest(
		http.MethodGet,
		"/oauth/authorize?"+authorizeQuery("http://client.example/callback", s256(verifier)).Encode(),
		nil,
	)
	authorizeRequest.AddCookie(&http.Cookie{Name: "lms_session", Value: "valid-session"})
	authServer.AuthorizeHandler().ServeHTTP(authorizeRecorder, authorizeRequest)
	location, err := url.Parse(authorizeRecorder.Header().Get("Location"))
	if err != nil || location.Query().Get("code") == "" {
		t.Fatalf("authorize user token: status=%d location=%q error=%v", authorizeRecorder.Code, location, err)
	}

	form := url.Values{
		"grant_type": {"authorization_code"}, "code": {location.Query().Get("code")},
		"redirect_uri": {"http://client.example/callback"}, "code_verifier": {verifier},
	}
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/oauth/token", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth("nextjs", "client-secret")
	authServer.TokenHandler().ServeHTTP(recorder, request)
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil || response.AccessToken == "" {
		t.Fatalf("issue user token: status=%d error=%v body=%s", recorder.Code, err, recorder.Body.String())
	}
	return response.AccessToken
}

func newTestServer(t *testing.T) *AuthorizationServer {
	t.Helper()
	return newTestServerWithScopes(t, []string{"library:read", "library:write"})
}

func newTestServerWithScopes(t *testing.T, scopes []string) *AuthorizationServer {
	t.Helper()
	manager := manage.NewDefaultManager()
	tokenStore, err := store.NewMemoryTokenStore()
	if err != nil {
		t.Fatalf("create memory token store: %v", err)
	}
	manager.MapTokenStorage(tokenStore)
	clientStore := store.NewClientStore()
	_ = clientStore.Set("nextjs", &models.Client{
		ID: "nextjs", Secret: "client-secret", Domain: "http://client.example/callback",
	})
	_ = clientStore.Set("other-client", &models.Client{
		ID: "other-client", Secret: "other-secret", Domain: "http://other.example/callback",
	})
	_ = clientStore.Set("machine-client", &models.Client{
		ID: "machine-client", Secret: "machine-secret",
	})
	_ = clientStore.Set("kong-gateway", &models.Client{
		ID: "kong-gateway", Secret: "gateway-secret",
	})
	manager.MapClientStorage(clientStore)
	policies := allowScopes{grants: map[string]map[oauth2lib.GrantType]bool{
		"nextjs": {
			oauth2lib.AuthorizationCode: true,
			oauth2lib.Refreshing:        true,
		},
		"other-client": {
			oauth2lib.AuthorizationCode: true,
			oauth2lib.Refreshing:        true,
		},
		"machine-client": {
			oauth2lib.ClientCredentials: true,
		},
	}, introspectors: map[string]bool{"kong-gateway": true}}

	return NewAuthorizationServer(manager, fakeSessions{}, policies, Config{
		Issuer: "http://auth.example", LoginURL: "http://client.example/login", SessionCookieName: "lms_session",
		CodeTTL: 5 * time.Minute, AccessTokenTTL: 15 * time.Minute, RefreshTokenTTL: 24 * time.Hour, SupportedScopes: scopes,
	})
}

func authorizeQuery(redirectURI, challenge string) url.Values {
	query := url.Values{
		"response_type": {"code"}, "client_id": {"nextjs"}, "redirect_uri": {redirectURI},
		"scope": {"library:read"}, "state": {"state-123"},
	}
	if challenge != "" {
		query.Set("code_challenge", challenge)
		query.Set("code_challenge_method", oauth2lib.CodeChallengeS256.String())
	}
	return query
}

func s256(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
