package db

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type fakeSeedExecutor struct {
	queries []string
	err     error
}

func (e *fakeSeedExecutor) ExecContext(_ context.Context, query string, _ ...any) (sql.Result, error) {
	e.queries = append(e.queries, query)
	return driver.RowsAffected(2), e.err
}

func TestApplySeedFileExecutesSQL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oauth_clients.sql")
	if err := os.WriteFile(path, []byte("INSERT INTO oauth_clients (id) VALUES ('member-nextjs-web')"), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := &fakeSeedExecutor{}

	if err := ApplySeedFile(context.Background(), executor, path); err != nil {
		t.Fatalf("ApplySeedFile() error = %v", err)
	}
	if len(executor.queries) != 1 || !strings.Contains(executor.queries[0], "INSERT INTO oauth_clients") {
		t.Fatalf("executed queries = %q, want OAuth client insert", executor.queries)
	}
}

func TestApplySeedFileExecutesStatementsInOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oauth_clients.sql")
	if err := os.WriteFile(path, []byte("INSERT INTO oauth_clients (id) VALUES ('client-1');\nDELETE FROM oauth_client_scopes;"), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := &fakeSeedExecutor{}

	if err := ApplySeedFile(context.Background(), executor, path); err != nil {
		t.Fatalf("ApplySeedFile() error = %v", err)
	}
	if len(executor.queries) != 2 || !strings.HasPrefix(executor.queries[0], "INSERT") || !strings.HasPrefix(executor.queries[1], "DELETE") {
		t.Fatalf("executed queries = %q, want insert then delete", executor.queries)
	}
}

func TestApplySeedFilePreservesSemicolonsInsideSQLText(t *testing.T) {
	path := filepath.Join(t.TempDir(), "seed.sql")
	contents := "-- comment; stays attached\nINSERT INTO example (value) VALUES ('one;two'); /* block; comment */\nDELETE FROM example;"
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := &fakeSeedExecutor{}

	if err := ApplySeedFile(context.Background(), executor, path); err != nil {
		t.Fatalf("ApplySeedFile() error = %v", err)
	}
	if len(executor.queries) != 2 || !strings.Contains(executor.queries[0], "'one;two'") || !strings.Contains(executor.queries[1], "block; comment") {
		t.Fatalf("executed queries = %q, want two intact statements", executor.queries)
	}
}

func TestApplySeedFilesExecutesEachFileInOrder(t *testing.T) {
	dir := t.TempDir()
	usersPath := filepath.Join(dir, "users.sql")
	clientsPath := filepath.Join(dir, "oauth_clients.sql")
	if err := os.WriteFile(usersPath, []byte("INSERT INTO users (id) VALUES ('user-1')"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(clientsPath, []byte("INSERT INTO oauth_clients (id) VALUES ('client-1')"), 0o600); err != nil {
		t.Fatal(err)
	}
	executor := &fakeSeedExecutor{}

	if err := ApplySeedFiles(context.Background(), executor, usersPath, clientsPath); err != nil {
		t.Fatalf("ApplySeedFiles() error = %v", err)
	}
	if len(executor.queries) != 2 || !strings.Contains(executor.queries[0], "INSERT INTO users") || !strings.Contains(executor.queries[1], "INSERT INTO oauth_clients") {
		t.Fatalf("executed queries = %q, want users then OAuth clients", executor.queries)
	}
}

func TestOAuthClientSeedContainsRequiredPrincipalsAndScopes(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "seeds", "oauth_clients.sql"))
	if err != nil {
		t.Fatalf("read OAuth client seed: %v", err)
	}
	seed := string(contents)
	for _, required := range []string{
		"'member-nextjs-web'", "'kong-gateway'", "'book-service'", "'user-service'", "'transaction-service'",
		"'identities:create'", "'book-stock:read'", "'book-stock:reserve'", "'book-stock:release'",
		"'transactions:read:any'", "'loans:return:any'", "'fines:manage'", "'books:manage'",
	} {
		if !strings.Contains(seed, required) {
			t.Errorf("OAuth client seed missing %s", required)
		}
	}
	if strings.Contains(seed, "library:read") || strings.Contains(seed, "library:write") {
		t.Fatal("OAuth client seed retains coarse library scopes")
	}
}

func TestUserSeedAssignsExplicitRoles(t *testing.T) {
	contents, err := os.ReadFile(filepath.Join("..", "..", "..", "db", "seeds", "users.sql"))
	if err != nil {
		t.Fatalf("read user seed: %v", err)
	}
	seed := string(contents)
	for _, required := range []string{"role_code", "'admin'", "'member'"} {
		if !strings.Contains(seed, required) {
			t.Errorf("user seed missing %s", required)
		}
	}
}

func TestApplySeedFileRejectsEmptySQL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oauth_clients.sql")
	if err := os.WriteFile(path, []byte("  \n"), 0o600); err != nil {
		t.Fatal(err)
	}

	err := ApplySeedFile(context.Background(), &fakeSeedExecutor{}, path)

	if err == nil || !strings.Contains(err.Error(), "empty") {
		t.Fatalf("ApplySeedFile() error = %v, want empty-file error", err)
	}
}

func TestApplySeedFileReportsMissingFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.sql")

	err := ApplySeedFile(context.Background(), &fakeSeedExecutor{}, path)

	if err == nil || !strings.Contains(err.Error(), path) {
		t.Fatalf("ApplySeedFile() error = %v, want path context", err)
	}
}

func TestApplySeedFileReportsDatabaseFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "oauth_clients.sql")
	if err := os.WriteFile(path, []byte("SELECT 1"), 0o600); err != nil {
		t.Fatal(err)
	}
	databaseErr := errors.New("database unavailable")

	err := ApplySeedFile(context.Background(), &fakeSeedExecutor{err: databaseErr}, path)

	if !errors.Is(err, databaseErr) || !strings.Contains(err.Error(), path) {
		t.Fatalf("ApplySeedFile() error = %v, want database error and path", err)
	}
}
