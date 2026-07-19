package entity

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrInvalidScope = errors.New("invalid scope")

type Scope struct {
	Code     string `db:"code"`
	Audience string `db:"audience"`
}

type ScopeGrant struct {
	Codes    []string
	Audience string
}

func (g ScopeGrant) String() string {
	return strings.Join(g.Codes, " ")
}

func ResolveUserScopeGrant(requested string, role Role, roleScopes, clientScopes []Scope) (ScopeGrant, error) {
	if _, err := ParseRole(role.String()); err != nil {
		return ScopeGrant{}, fmt.Errorf("%w: %v", ErrInvalidScope, err)
	}
	return resolveScopeGrant(requested, clientScopes, roleScopes, true)
}

func ResolveServiceScopeGrant(requested string, clientScopes []Scope) (ScopeGrant, error) {
	return resolveScopeGrant(requested, clientScopes, nil, false)
}

func resolveScopeGrant(requested string, clientScopes, requiredScopes []Scope, requireRole bool) (ScopeGrant, error) {
	requestedCodes := strings.Fields(requested)
	if len(requestedCodes) == 0 {
		return ScopeGrant{}, fmt.Errorf("%w: at least one scope is required", ErrInvalidScope)
	}

	clients := scopeIndex(clientScopes)
	required := scopeIndex(requiredScopes)
	unique := make(map[string]struct{}, len(requestedCodes))
	audience := ""
	for _, code := range requestedCodes {
		clientScope, ok := clients[code]
		if !ok || clientScope.Audience == "" {
			return ScopeGrant{}, fmt.Errorf("%w: scope %q is not granted to client", ErrInvalidScope, code)
		}
		if requireRole {
			roleScope, roleAllows := required[code]
			if !roleAllows || roleScope.Audience != clientScope.Audience {
				return ScopeGrant{}, fmt.Errorf("%w: scope %q is not granted to role", ErrInvalidScope, code)
			}
		}
		if audience != "" && audience != clientScope.Audience {
			return ScopeGrant{}, fmt.Errorf("%w: requested scopes target multiple audiences", ErrInvalidScope)
		}
		audience = clientScope.Audience
		unique[code] = struct{}{}
	}

	codes := make([]string, 0, len(unique))
	for code := range unique {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	return ScopeGrant{Codes: codes, Audience: audience}, nil
}

func scopeIndex(scopes []Scope) map[string]Scope {
	result := make(map[string]Scope, len(scopes))
	for _, scope := range scopes {
		result[scope.Code] = scope
	}
	return result
}
