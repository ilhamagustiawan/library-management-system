package response

import "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/domain/entity"

type OAuthMetadata struct {
	Issuer                            string   `json:"issuer" example:"http://localhost:8081"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint" example:"http://localhost:8081/oauth/authorize"`
	TokenEndpoint                     string   `json:"token_endpoint" example:"http://localhost:8081/oauth/token"`
	IntrospectionEndpoint             string   `json:"introspection_endpoint" example:"http://localhost:8081/oauth/introspect"`
	ResponseTypesSupported            []string `json:"response_types_supported" example:"code"`
	GrantTypesSupported               []string `json:"grant_types_supported" example:"authorization_code,refresh_token,client_credentials"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported" example:"client_secret_basic"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported" example:"S256"`
	ScopesSupported                   []string `json:"scopes_supported" example:"books:read,loans:borrow:self"`
}

func NewOAuthMetadata(metadata entity.OAuthMetadata) OAuthMetadata {
	return OAuthMetadata{
		Issuer:                            metadata.Issuer,
		AuthorizationEndpoint:             metadata.AuthorizationEndpoint,
		TokenEndpoint:                     metadata.TokenEndpoint,
		IntrospectionEndpoint:             metadata.IntrospectionEndpoint,
		ResponseTypesSupported:            metadata.ResponseTypesSupported,
		GrantTypesSupported:               metadata.GrantTypesSupported,
		TokenEndpointAuthMethodsSupported: metadata.TokenAuthMethodsSupported,
		CodeChallengeMethodsSupported:     metadata.CodeChallengeMethodsSupport,
		ScopesSupported:                   metadata.ScopesSupported,
	}
}
