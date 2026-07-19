package oauthclient

import (
	"testing"

	"github.com/go-oauth2/oauth2/v4"
)

func TestClientGrantPolicy(t *testing.T) {
	tests := []struct {
		name    string
		kind    Kind
		grant   oauth2.GrantType
		allowed bool
	}{
		{name: "authorization code", kind: KindAuthorizationCode, grant: oauth2.AuthorizationCode, allowed: true},
		{name: "authorization refresh", kind: KindAuthorizationCode, grant: oauth2.Refreshing, allowed: true},
		{name: "authorization rejects client credentials", kind: KindAuthorizationCode, grant: oauth2.ClientCredentials},
		{name: "client credentials", kind: KindClientCredentials, grant: oauth2.ClientCredentials, allowed: true},
		{name: "client credentials rejects refresh", kind: KindClientCredentials, grant: oauth2.Refreshing},
		{name: "resource server rejects token grants", kind: KindResourceServer, grant: oauth2.ClientCredentials},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := Client{Kind: test.kind}
			if got := client.AllowsGrant(test.grant); got != test.allowed {
				t.Fatalf("AllowsGrant(%q) = %t, want %t", test.grant, got, test.allowed)
			}
		})
	}
}

func TestOnlyResourceServerCanIntrospect(t *testing.T) {
	for _, kind := range []Kind{KindAuthorizationCode, KindClientCredentials, KindResourceServer} {
		client := Client{Kind: kind}
		want := kind == KindResourceServer
		if got := client.CanIntrospect(); got != want {
			t.Errorf("kind %q CanIntrospect() = %t, want %t", kind, got, want)
		}
	}
}
