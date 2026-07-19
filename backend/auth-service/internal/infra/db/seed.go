package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
)

type SQLExecutor interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func ApplySeedFile(ctx context.Context, executor SQLExecutor, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read seed file %q: %w", path, err)
	}
	query := strings.TrimSpace(string(contents))
	if query == "" {
		return fmt.Errorf("seed file %q is empty", path)
	}
	if _, err := executor.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("execute seed file %q: %w", path, err)
	}
	return nil
}
