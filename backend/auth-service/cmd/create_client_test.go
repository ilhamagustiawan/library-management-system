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
