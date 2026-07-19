package entity

import "testing"

func TestParseRoleAcceptsCatalogRoles(t *testing.T) {
	for _, value := range []string{"member", "admin"} {
		role, err := ParseRole(value)
		if err != nil || role.String() != value {
			t.Fatalf("ParseRole(%q) = (%q, %v)", value, role, err)
		}
	}
}

func TestParseRoleRejectsUnknownRole(t *testing.T) {
	if _, err := ParseRole("librarian"); err == nil {
		t.Fatal("ParseRole() error = nil, want unknown-role error")
	}
}
