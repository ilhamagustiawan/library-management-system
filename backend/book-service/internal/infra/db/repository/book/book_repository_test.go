package book

import "testing"

func TestOrderClauseAllowlist(t *testing.T) {
	for _, test := range []struct {
		field string
		order string
		want  string
	}{
		{field: "title", order: "asc", want: "title ASC"},
		{field: "author", order: "desc", want: "author DESC"},
		{field: "createdAt", order: "desc", want: "created_at DESC"},
		{field: "title; DROP TABLE books", order: "desc", want: "title ASC"},
	} {
		if got := orderClause(test.field, test.order); got != test.want {
			t.Fatalf("orderClause(%q, %q) = %q, want %q", test.field, test.order, got, test.want)
		}
	}
}
