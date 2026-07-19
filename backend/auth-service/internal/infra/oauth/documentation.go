package oauth

// OAuthTokenResponse documents the OAuth token endpoint response. RefreshToken
// is omitted for client credentials grants.
type OAuthTokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsibGlicmFyeS1hcGkiXX0.signature"`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" example:"900"`
	RefreshToken string `json:"refresh_token,omitempty" example:"opaque-refresh-token"`
	Scope        string `json:"scope" example:"books:read loans:borrow:self"`
}

// OAuthErrorResponse follows the OAuth 2.0 error response shape.
type OAuthErrorResponse struct {
	Error            string `json:"error" example:"invalid_request"`
	ErrorDescription string `json:"error_description" example:"the request is invalid"`
}

// APIErrorResponse documents service-level errors produced before OAuth handling.
type APIErrorResponse struct {
	Code    string `json:"code" example:"LMS-429001"`
	Message string `json:"message" example:"too many authentication attempts"`
}

// OAuthIntrospectionResponse documents RFC 7662 token state. Inactive tokens
// return only Active; other fields are omitted.
type OAuthIntrospectionResponse struct {
	Active    bool     `json:"active" example:"true"`
	ClientID  string   `json:"client_id,omitempty" example:"member-nextjs-web"`
	Scope     string   `json:"scope,omitempty" example:"books:read"`
	TokenType string   `json:"token_type,omitempty" example:"Bearer"`
	Subject   string   `json:"sub,omitempty" example:"f81d4fae-7dec-11d0-a765-00a0c91e6bf6"`
	Issuer    string   `json:"iss,omitempty" example:"http://localhost:8081"`
	IssuedAt  int64    `json:"iat,omitempty" example:"1784448000"`
	ExpiresAt int64    `json:"exp,omitempty" example:"1784448900"`
	Audience  []string `json:"aud,omitempty" example:"library-api"`
	Role      string   `json:"role,omitempty" example:"member"`
}
