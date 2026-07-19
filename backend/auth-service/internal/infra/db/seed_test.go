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
	query string
	err   error
}

func (e *fakeSeedExecutor) ExecContext(_ context.Context, query string, _ ...any) (sql.Result, error) {
	e.query = query
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
	if !strings.Contains(executor.query, "INSERT INTO oauth_clients") {
		t.Fatalf("executed query = %q, want OAuth client insert", executor.query)
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
