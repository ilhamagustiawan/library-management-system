package entity

import (
	"errors"
	"reflect"
	"testing"
)

var memberScopes = []Scope{
	{Code: "books:read", Audience: "library-api"},
	{Code: "loans:borrow:self", Audience: "library-api"},
}

func TestResolveUserScopeGrantRequiresRoleAndClientPermission(t *testing.T) {
	grant, err := ResolveUserScopeGrant(
		"loans:borrow:self books:read", RoleMember, memberScopes, memberScopes,
	)
	if err != nil {
		t.Fatalf("ResolveUserScopeGrant() error = %v", err)
	}
	if grant.Audience != "library-api" || !reflect.DeepEqual(grant.Codes, []string{"books:read", "loans:borrow:self"}) {
		t.Fatalf("grant = %#v", grant)
	}
}

func TestResolveUserScopeGrantRejectsEscalation(t *testing.T) {
	clientScopes := append(append([]Scope(nil), memberScopes...), Scope{Code: "transactions:read:any", Audience: "library-api"})
	_, err := ResolveUserScopeGrant("transactions:read:any", RoleMember, memberScopes, clientScopes)
	if !errors.Is(err, ErrInvalidScope) {
		t.Fatalf("error = %v, want invalid scope", err)
	}
}

func TestResolveUserScopeGrantRejectsRoleWithoutScopes(t *testing.T) {
	_, err := ResolveUserScopeGrant("books:read", RoleMember, nil, memberScopes)
	if !errors.Is(err, ErrInvalidScope) {
		t.Fatalf("error = %v, want invalid scope", err)
	}
}

func TestResolveServiceScopeGrantRejectsMixedAudiences(t *testing.T) {
	_, err := ResolveServiceScopeGrant("identities:create book-stock:read", []Scope{
		{Code: "identities:create", Audience: "auth-service"},
		{Code: "book-stock:read", Audience: "book-service"},
	})
	if !errors.Is(err, ErrInvalidScope) {
		t.Fatalf("error = %v, want invalid scope", err)
	}
}

func TestResolveServiceScopeGrantRejectsEmptyAndUnknownScopes(t *testing.T) {
	for _, requested := range []string{"", "unknown:scope"} {
		if _, err := ResolveServiceScopeGrant(requested, memberScopes); !errors.Is(err, ErrInvalidScope) {
			t.Fatalf("ResolveServiceScopeGrant(%q) error = %v", requested, err)
		}
	}
}
