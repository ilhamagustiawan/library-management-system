package cmd

import "testing"

func TestValidateRedirectURIAllowsHTTPSAndLoopbackHTTP(t *testing.T) {
	for _, redirectURI := range []string{
		"https://client.example/callback",
		"http://localhost:3000/callback",
		"http://127.0.0.1:3000/callback",
	} {
		if err := validateRedirectURI(redirectURI); err != nil {
			t.Errorf("validateRedirectURI(%q) error = %v", redirectURI, err)
		}
	}
}

func TestValidateRedirectURIRejectsUnsafeSchemes(t *testing.T) {
	for _, redirectURI := range []string{
		"http://client.example/callback",
		"ftp://localhost/callback",
		"library://localhost/callback",
	} {
		if err := validateRedirectURI(redirectURI); err == nil {
			t.Errorf("validateRedirectURI(%q) error = nil, want rejection", redirectURI)
		}
	}
}

func TestValidateClientProvisioningEnforcesKindRedirectAndScopes(t *testing.T) {
	for _, test := range []struct {
		kind        string
		redirectURI string
		scopes      string
		wantError   bool
	}{
		{kind: "authorization_code", redirectURI: "https://client.example/callback", scopes: "books:read"},
		{kind: "client_credentials", scopes: "book-stock:read"},
		{kind: "resource_server"},
		{kind: "authorization_code", scopes: "books:read", wantError: true},
		{kind: "client_credentials", redirectURI: "https://client.example/callback", scopes: "book-stock:read", wantError: true},
		{kind: "client_credentials", wantError: true},
		{kind: "unknown", scopes: "books:read", wantError: true},
	} {
		err := validateClientProvisioning(test.kind, test.redirectURI, test.scopes)
		if (err != nil) != test.wantError {
			t.Errorf("validateClientProvisioning(%q, %q, %q) error = %v", test.kind, test.redirectURI, test.scopes, err)
		}
	}
}
