package entity

import "strings"

type OAuthMetadata struct {
	Issuer                      string
	AuthorizationEndpoint       string
	TokenEndpoint               string
	IntrospectionEndpoint       string
	ResponseTypesSupported      []string
	GrantTypesSupported         []string
	TokenAuthMethodsSupported   []string
	CodeChallengeMethodsSupport []string
	ScopesSupported             []string
}

func NewOAuthMetadata(issuerRaw string, supportedScopes []string) OAuthMetadata {
	issuer := strings.TrimRight(issuerRaw, "/")
	return OAuthMetadata{
		Issuer:                      issuer,
		AuthorizationEndpoint:       issuer + "/oauth/authorize",
		TokenEndpoint:               issuer + "/oauth/token",
		IntrospectionEndpoint:       issuer + "/oauth/introspect",
		ResponseTypesSupported:      []string{"code"},
		GrantTypesSupported:         []string{"authorization_code", "client_credentials", "refresh_token"},
		TokenAuthMethodsSupported:   []string{"client_secret_basic"},
		CodeChallengeMethodsSupport: []string{"S256"},
		ScopesSupported:             supportedScopes,
	}
}
