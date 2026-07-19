package oauth

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	oauth2lib "github.com/go-oauth2/oauth2/v4"
	oauth2errors "github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2server "github.com/go-oauth2/oauth2/v4/server"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"
)

type Config struct {
	Issuer            string
	LoginURL          string
	SessionCookieName string
	CodeTTL           time.Duration
	AccessTokenTTL    time.Duration
	RefreshTokenTTL   time.Duration
	SupportedScopes   []string
	JWTSigningKey     []byte
}

type SessionAuthenticator interface {
	AuthenticateSession(ctx context.Context, token string) (*entity.User, error)
	FindUser(ctx context.Context, id string) (*entity.User, error)
}

type ClientPolicyStore interface {
	GetScopes(ctx context.Context, clientID string) ([]entity.Scope, error)
	GetRoleScopes(ctx context.Context, role entity.Role) ([]entity.Scope, error)
	AllowsGrant(ctx context.Context, clientID string, grant oauth2lib.GrantType) (bool, error)
	CanIntrospect(ctx context.Context, clientID string) (bool, error)
}

type AuthorizationServer struct {
	server   *oauth2server.Server
	manager  *manage.Manager
	sessions SessionAuthenticator
	policies ClientPolicyStore
	config   Config
}

var ErrInvalidServiceToken = errors.New("invalid service token")
var ErrInsufficientServiceGrant = errors.New("insufficient service grant")

func NewAuthorizationServer(
	manager *manage.Manager,
	sessions SessionAuthenticator,
	policies ClientPolicyStore,
	config Config,
) (*AuthorizationServer, error) {
	accessTokens, err := NewAccessTokenGenerator(AccessTokenConfig{Issuer: config.Issuer, SigningKey: config.JWTSigningKey})
	if err != nil {
		return nil, fmt.Errorf("configure OAuth access tokens: %w", err)
	}
	manager.MapAccessGenerate(accessTokens)
	manager.SetExtractExtensionHandler(extractTokenGrant)
	manager.SetAuthorizeCodeExp(config.CodeTTL)
	manager.SetAuthorizeCodeTokenCfg(&manage.Config{
		AccessTokenExp: config.AccessTokenTTL, RefreshTokenExp: config.RefreshTokenTTL, IsGenerateRefresh: true,
	})
	manager.SetClientTokenCfg(&manage.Config{AccessTokenExp: config.AccessTokenTTL, IsGenerateRefresh: false})
	manager.SetRefreshTokenCfg(&manage.RefreshingConfig{
		AccessTokenExp: config.AccessTokenTTL, RefreshTokenExp: config.RefreshTokenTTL,
		IsGenerateRefresh: true, IsResetRefreshTime: true, IsRemoveAccess: true, IsRemoveRefreshing: true,
	})
	manager.SetValidateURIHandler(validateRedirectURI)

	serverConfig := oauth2server.NewConfig()
	serverConfig.AllowGetAccessRequest = false
	serverConfig.AllowedResponseTypes = []oauth2lib.ResponseType{oauth2lib.Code}
	serverConfig.AllowedGrantTypes = []oauth2lib.GrantType{
		oauth2lib.AuthorizationCode,
		oauth2lib.ClientCredentials,
		oauth2lib.Refreshing,
	}
	serverConfig.AllowedCodeChallengeMethods = []oauth2lib.CodeChallengeMethod{oauth2lib.CodeChallengeS256}
	serverConfig.ForcePKCE = true

	srv := oauth2server.NewServer(serverConfig, manager)
	srv.SetClientInfoHandler(oauth2server.ClientBasicHandler)
	srv.SetAccessTokenResolveHandler(bearerToken)

	result := &AuthorizationServer{server: srv, manager: manager, sessions: sessions, policies: policies, config: config}
	srv.SetClientAuthorizedHandler(result.clientAuthorized)
	srv.SetClientScopeHandler(result.clientScope)
	srv.SetAuthorizeScopeHandler(result.authorizeScope)
	srv.SetRefreshingScopeHandler(refreshingScope)
	srv.SetRefreshingValidationHandler(result.validateRefresh)
	srv.SetUserAuthorizationHandler(result.authorizeUser)
	return result, nil
}

func validateRedirectURI(registered, requested string) error {
	if registered != requested {
		return oauth2errors.ErrInvalidRedirectURI
	}
	return nil
}

func (s *AuthorizationServer) clientAuthorized(clientID string, grant oauth2lib.GrantType) (bool, error) {
	return s.policies.AllowsGrant(context.Background(), clientID, grant)
}

func (s *AuthorizationServer) clientScope(request *oauth2lib.TokenGenerateRequest) (bool, error) {
	ctx := context.Background()
	if request.Request != nil {
		ctx = request.Request.Context()
	}
	if request.UserID != "" {
		_, ok := tokenGrantFromContext(ctx)
		return ok, nil
	}
	clientScopes, err := s.policies.GetScopes(ctx, request.ClientID)
	if err != nil {
		return false, err
	}
	grant, err := entity.ResolveServiceScopeGrant(request.Scope, clientScopes)
	if errors.Is(err, entity.ErrInvalidScope) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	request.Scope = grant.String()
	if request.Request != nil {
		setTokenGrantContext(request.Request, tokenGrant{ScopeGrant: grant})
	}
	return true, nil
}

func (s *AuthorizationServer) authorizeScope(_ http.ResponseWriter, r *http.Request) (string, error) {
	user, ok := authorizedUserFromContext(r.Context())
	if !ok {
		return "", oauth2errors.ErrAccessDenied
	}
	roleScopes, err := s.policies.GetRoleScopes(r.Context(), user.Role)
	if err != nil {
		return "", err
	}
	clientScopes, err := s.policies.GetScopes(r.Context(), r.FormValue("client_id"))
	if err != nil {
		return "", err
	}
	grant, err := entity.ResolveUserScopeGrant(r.FormValue("scope"), user.Role, roleScopes, clientScopes)
	if errors.Is(err, entity.ErrInvalidScope) {
		return "", oauth2errors.ErrInvalidScope
	}
	if err != nil {
		return "", err
	}
	setTokenGrantContext(r, tokenGrant{ScopeGrant: grant, Role: user.Role})
	return grant.String(), nil
}

func refreshingScope(request *oauth2lib.TokenGenerateRequest, oldScope string) (bool, error) {
	return scopeSubset(request.Scope, oldScope), nil
}

func (s *AuthorizationServer) validateRefresh(token oauth2lib.TokenInfo) (bool, error) {
	clientScopes, err := s.policies.GetScopes(context.Background(), token.GetClientID())
	if err != nil {
		return false, err
	}
	extendable, ok := token.(oauth2lib.ExtendableTokenInfo)
	if !ok {
		return false, nil
	}
	extension := extendable.GetExtension()
	storedAudience := extension.Get(tokenAudienceExtension)
	if token.GetUserID() == "" {
		grant, err := entity.ResolveServiceScopeGrant(token.GetScope(), clientScopes)
		return err == nil && grant.String() == token.GetScope() && grant.Audience == storedAudience, nil
	}

	user, err := s.sessions.FindUser(context.Background(), token.GetUserID())
	if err != nil {
		return false, err
	}
	roleScopes, err := s.policies.GetRoleScopes(context.Background(), user.Role)
	if err != nil {
		return false, err
	}
	grant, err := entity.ResolveUserScopeGrant(token.GetScope(), user.Role, roleScopes, clientScopes)
	if err != nil || grant.String() != token.GetScope() || grant.Audience != storedAudience {
		return false, nil
	}
	extension.Set(tokenRoleExtension, user.Role.String())
	extendable.SetExtension(extension)
	return true, nil
}

// AuthorizeHandler starts an OAuth Authorization Code flow.
//
// @Summary Authorize OAuth client
// @Description Validates the client, redirect URI, state, and S256 PKCE challenge. Redirects unauthenticated users to login.
// @Tags OAuth
// @Param response_type query string true "OAuth response type" Enums(code)
// @Param client_id query string true "OAuth client ID"
// @Param redirect_uri query string true "Exact registered callback URL"
// @Param scope query string true "Space-delimited scopes"
// @Param state query string true "CSRF correlation value"
// @Param code_challenge query string true "Base64url-encoded SHA-256 PKCE challenge"
// @Param code_challenge_method query string true "PKCE method" Enums(S256)
// @Success 302 {string} string "Redirect to login or callback with authorization code"
// @Failure 400 {object} OAuthErrorResponse "Invalid authorization request"
// @Router /oauth/authorize [get]
func (s *AuthorizationServer) AuthorizeHandler() http.Handler {
	return http.HandlerFunc(s.handleAuthorize)
}

func (s *AuthorizationServer) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "unable to parse authorization request")
		return
	}
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "authorization request must use GET or POST")
		return
	}
	if r.FormValue("response_type") != oauth2lib.Code.String() || r.FormValue("state") == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "response_type=code and state are required")
		return
	}
	if r.FormValue("code_challenge") == "" || r.FormValue("code_challenge_method") != oauth2lib.CodeChallengeS256.String() {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "PKCE with code_challenge_method=S256 is required")
		return
	}

	client, err := s.manager.GetClient(r.Context(), r.FormValue("client_id"))
	if err != nil || client.GetDomain() != r.FormValue("redirect_uri") {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "invalid client or redirect_uri")
		return
	}

	writer := &trackingResponseWriter{ResponseWriter: w}
	if err := s.server.HandleAuthorizeRequest(writer, r); err != nil && !writer.wroteHeader {
		if errors.Is(err, oauth2errors.ErrInvalidScope) {
			writeOAuthError(w, http.StatusBadRequest, "invalid_scope", "requested scope is not allowed")
			return
		}
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "authorization request rejected")
	}
}

// TokenHandler exchanges grants for JWT access tokens.
//
// @Summary Issue OAuth token
// @Description Supports authorization_code, refresh_token, and provisioned client_credentials clients. Client authentication uses HTTP Basic.
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Security BasicAuth
// @Param grant_type formData string true "OAuth grant type" Enums(authorization_code,refresh_token,client_credentials)
// @Param code formData string false "Authorization code; required for authorization_code"
// @Param redirect_uri formData string false "Exact callback URL; required for authorization_code"
// @Param code_verifier formData string false "PKCE verifier; required for authorization_code"
// @Param refresh_token formData string false "Refresh token; required for refresh_token"
// @Param scope formData string false "Space-delimited requested scopes"
// @Success 200 {object} OAuthTokenResponse "Token issued"
// @Failure 400 {object} OAuthErrorResponse "Invalid request or grant"
// @Failure 401 {object} OAuthErrorResponse "Client authentication failed"
// @Failure 429 {object} APIErrorResponse "Rate limit exceeded"
// @Router /oauth/token [post]
func (s *AuthorizationServer) TokenHandler() http.Handler {
	return http.HandlerFunc(s.handleToken)
}

func (s *AuthorizationServer) handleToken(w http.ResponseWriter, r *http.Request) {
	if err := s.validateTokenClient(r); err != nil {
		if errors.Is(err, oauth2errors.ErrInvalidClient) {
			w.Header().Set("WWW-Authenticate", `Basic realm="oauth/token"`)
			writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "client authentication failed")
			return
		}
		writeOAuthError(w, http.StatusBadRequest, "invalid_grant", "the supplied grant is invalid")
		return
	}
	writer := &trackingResponseWriter{ResponseWriter: w}
	if err := s.server.HandleTokenRequest(writer, r); err != nil && !writer.wroteHeader {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "token request rejected")
	}
}

type introspectionResponse struct {
	Active    bool        `json:"active"`
	ClientID  string      `json:"client_id,omitempty"`
	Scope     string      `json:"scope,omitempty"`
	TokenType string      `json:"token_type,omitempty"`
	Subject   string      `json:"sub,omitempty"`
	Issuer    string      `json:"iss,omitempty"`
	IssuedAt  int64       `json:"iat,omitempty"`
	ExpiresAt int64       `json:"exp,omitempty"`
	Audience  []string    `json:"aud,omitempty"`
	Role      entity.Role `json:"role,omitempty"`
}

// IntrospectionHandler reports JWT access-token state.
//
// @Summary Introspect access token
// @Description RFC 7662 endpoint restricted to provisioned resource-server clients using HTTP Basic.
// @Tags OAuth
// @Accept x-www-form-urlencoded
// @Produce json
// @Security BasicAuth
// @Param token formData string true "JWT access token"
// @Param token_type_hint formData string false "Token type hint" Enums(access_token)
// @Success 200 {object} OAuthIntrospectionResponse "Token state; inactive tokens include only active=false"
// @Failure 400 {object} OAuthErrorResponse "Token missing or request invalid"
// @Failure 401 {object} OAuthErrorResponse "Client authentication or authorization failed"
// @Failure 500 {object} OAuthErrorResponse "Token state could not be determined"
// @Router /oauth/introspect [post]
func (s *AuthorizationServer) IntrospectionHandler() http.Handler {
	return http.HandlerFunc(s.handleIntrospection)
}

func (s *AuthorizationServer) handleIntrospection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "introspection request must use POST")
		return
	}
	if err := r.ParseForm(); err != nil || r.FormValue("token") == "" {
		writeOAuthError(w, http.StatusBadRequest, "invalid_request", "token is required")
		return
	}
	clientID, err := s.authenticateClient(r)
	if err != nil {
		writeIntrospectionAuthError(w)
		return
	}
	allowed, err := s.policies.CanIntrospect(r.Context(), clientID)
	if err != nil || !allowed {
		writeIntrospectionAuthError(w)
		return
	}

	token, err := s.manager.LoadAccessToken(r.Context(), r.FormValue("token"))
	if errors.Is(err, oauth2errors.ErrInvalidAccessToken) || errors.Is(err, oauth2errors.ErrExpiredAccessToken) {
		writeIntrospectionResponse(w, introspectionResponse{Active: false})
		return
	}
	if err != nil || token == nil {
		writeOAuthError(w, http.StatusInternalServerError, "server_error", "token state could not be determined")
		return
	}

	response := introspectionResponse{
		Active: true, ClientID: token.GetClientID(), Scope: token.GetScope(), TokenType: "Bearer",
		Subject: token.GetUserID(), Issuer: strings.TrimRight(s.config.Issuer, "/"),
		IssuedAt: token.GetAccessCreateAt().Unix(),
	}
	if extendable, ok := token.(oauth2lib.ExtendableTokenInfo); ok {
		extension := extendable.GetExtension()
		if audience := extension.Get(tokenAudienceExtension); audience != "" {
			response.Audience = []string{audience}
		}
		response.Role = entity.Role(extension.Get(tokenRoleExtension))
	}
	if response.Subject == "" && response.Role == "" {
		response.Subject = response.ClientID
	}
	if ttl := token.GetAccessExpiresIn(); ttl > 0 {
		response.ExpiresAt = token.GetAccessCreateAt().Add(ttl).Unix()
	}
	writeIntrospectionResponse(w, response)
}

func writeIntrospectionAuthError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", `Basic realm="oauth/introspect"`)
	writeOAuthError(w, http.StatusUnauthorized, "invalid_client", "client authentication failed")
}

func writeIntrospectionResponse(w http.ResponseWriter, response introspectionResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	_ = json.NewEncoder(w).Encode(response)
}

func (s *AuthorizationServer) validateTokenClient(r *http.Request) error {
	if r.Method != http.MethodPost {
		return oauth2errors.ErrInvalidRequest
	}
	if err := r.ParseForm(); err != nil {
		return oauth2errors.ErrInvalidRequest
	}
	clientID, err := s.authenticateClient(r)
	if err != nil {
		return err
	}

	if r.FormValue("grant_type") == oauth2lib.Refreshing.String() {
		token, err := s.manager.LoadRefreshToken(r.Context(), r.FormValue("refresh_token"))
		if err != nil || token.GetClientID() != clientID {
			return oauth2errors.ErrInvalidGrant
		}
	}
	return nil
}

func (s *AuthorizationServer) authenticateClient(r *http.Request) (string, error) {
	clientID, clientSecret, ok := r.BasicAuth()
	if !ok || clientID == "" || clientSecret == "" {
		return "", oauth2errors.ErrInvalidClient
	}
	client, err := s.manager.GetClient(r.Context(), clientID)
	if err != nil || client == nil || client.IsPublic() {
		return "", oauth2errors.ErrInvalidClient
	}
	if verifier, ok := client.(oauth2lib.ClientPasswordVerifier); ok {
		if !verifier.VerifyPassword(clientSecret) {
			return "", oauth2errors.ErrInvalidClient
		}
	} else if subtle.ConstantTimeCompare([]byte(client.GetSecret()), []byte(clientSecret)) != 1 {
		return "", oauth2errors.ErrInvalidClient
	}
	return clientID, nil
}
func (s *AuthorizationServer) LoadAccessToken(ctx context.Context, rawToken string) (oauth2lib.TokenInfo, error) {
	return s.manager.LoadAccessToken(ctx, rawToken)
}

func (s *AuthorizationServer) FindUserIDByAccessToken(ctx context.Context, rawToken string) (string, error) {
	token, err := s.manager.LoadAccessToken(ctx, rawToken)
	if err != nil {
		return "", err
	}
	return token.GetUserID(), nil
}

func (s *AuthorizationServer) AuthenticateServiceToken(
	ctx context.Context,
	rawToken string,
	expectedClientID string,
	expectedAudience string,
	requiredScope string,
) error {
	token, err := s.manager.LoadAccessToken(ctx, rawToken)
	if err != nil || token == nil || token.GetUserID() != "" || token.GetClientID() != expectedClientID {
		return ErrInvalidServiceToken
	}
	extendable, ok := token.(oauth2lib.ExtendableTokenInfo)
	if !ok {
		return ErrInvalidServiceToken
	}
	extension := extendable.GetExtension()
	if extension.Get(tokenAudienceExtension) != expectedAudience || extension.Get(tokenRoleExtension) != "" {
		return ErrInvalidServiceToken
	}
	if !scopeSubset(requiredScope, token.GetScope()) {
		return ErrInsufficientServiceGrant
	}
	return nil
}

func (s *AuthorizationServer) Metadata() entity.OAuthMetadata {
	return entity.NewOAuthMetadata(s.config.Issuer, s.config.SupportedScopes)
}

func (s *AuthorizationServer) authorizeUser(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(s.config.SessionCookieName)
	if err == nil {
		user, authErr := s.sessions.AuthenticateSession(r.Context(), cookie.Value)
		if authErr == nil {
			setAuthorizedUserContext(r, user)
			return user.ID, nil
		}
	}

	loginURL, err := url.Parse(s.config.LoginURL)
	if err != nil {
		return "", err
	}
	issuer, err := url.Parse(s.config.Issuer)
	if err != nil {
		return "", err
	}
	issuer.Path = r.URL.Path
	issuer.RawQuery = r.URL.RawQuery
	query := loginURL.Query()
	query.Set("return_to", issuer.String())
	loginURL.RawQuery = query.Encode()
	http.Redirect(w, r, loginURL.String(), http.StatusFound)
	return "", nil
}

type requestContextKey string

const authorizedUserContextKey requestContextKey = "authorized-user"
const tokenGrantContextKey requestContextKey = "token-grant"

type tokenGrant struct {
	entity.ScopeGrant
	Role entity.Role
}

func setAuthorizedUserContext(r *http.Request, user *entity.User) {
	*r = *r.WithContext(context.WithValue(r.Context(), authorizedUserContextKey, user))
}

func authorizedUserFromContext(ctx context.Context) (*entity.User, bool) {
	user, ok := ctx.Value(authorizedUserContextKey).(*entity.User)
	return user, ok && user != nil
}

func setTokenGrantContext(r *http.Request, grant tokenGrant) {
	*r = *r.WithContext(context.WithValue(r.Context(), tokenGrantContextKey, grant))
}

func tokenGrantFromContext(ctx context.Context) (tokenGrant, bool) {
	grant, ok := ctx.Value(tokenGrantContextKey).(tokenGrant)
	return grant, ok
}

func extractTokenGrant(request *oauth2lib.TokenGenerateRequest, token oauth2lib.ExtendableTokenInfo) {
	if request.Request == nil || token.GetExtension().Get(tokenAudienceExtension) != "" {
		return
	}
	grant, ok := tokenGrantFromContext(request.Request.Context())
	if !ok {
		return
	}
	extension := token.GetExtension()
	if extension == nil {
		extension = make(url.Values)
	}
	extension.Set(tokenAudienceExtension, grant.Audience)
	if grant.Role != "" {
		extension.Set(tokenRoleExtension, grant.Role.String())
	}
	token.SetExtension(extension)
}

func bearerToken(r *http.Request) (string, bool) {
	parts := strings.Fields(r.Header.Get("Authorization"))
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", false
	}
	return parts[1], true
}

func scopeSubset(requested, allowed string) bool {
	allowedSet := make(map[string]struct{})
	for _, value := range strings.Fields(allowed) {
		allowedSet[value] = struct{}{}
	}
	for _, value := range strings.Fields(requested) {
		if _, ok := allowedSet[value]; !ok {
			return false
		}
	}
	return true
}

func writeOAuthError(w http.ResponseWriter, status int, code, description string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": code, "error_description": description})
}

type trackingResponseWriter struct {
	http.ResponseWriter
	wroteHeader bool
}

func (w *trackingResponseWriter) WriteHeader(status int) {
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(status)
}

func (w *trackingResponseWriter) Write(body []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(body)
}
