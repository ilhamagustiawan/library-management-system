package entity

import (
	"errors"
	"fmt"
)

var ErrUnknownRole = errors.New("unknown role")

type Role string

const (
	RoleMember Role = "member"
	RoleAdmin  Role = "admin"
)

func ParseRole(value string) (Role, error) {
	role := Role(value)
	switch role {
	case RoleMember, RoleAdmin:
		return role, nil
	default:
		return "", fmt.Errorf("%w: %q", ErrUnknownRole, value)
	}
}

func (r Role) String() string {
	return string(r)
}
