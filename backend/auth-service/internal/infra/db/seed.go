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

func ApplySeedFiles(ctx context.Context, executor SQLExecutor, paths ...string) error {
	for _, path := range paths {
		if err := ApplySeedFile(ctx, executor, path); err != nil {
			return err
		}
	}
	return nil
}

func ApplySeedFile(ctx context.Context, executor SQLExecutor, path string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read seed file %q: %w", path, err)
	}
	statements, err := splitSeedStatements(string(contents))
	if err != nil {
		return fmt.Errorf("parse seed file %q: %w", path, err)
	}
	if len(statements) == 0 {
		return fmt.Errorf("seed file %q is empty", path)
	}
	for index, statement := range statements {
		if _, err := executor.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("execute seed file %q statement %d: %w", path, index+1, err)
		}
	}
	return nil
}

func splitSeedStatements(contents string) ([]string, error) {
	var statements []string
	var statement strings.Builder
	var quote byte
	lineComment := false
	blockComment := false

	for index := 0; index < len(contents); index++ {
		current := contents[index]
		next := byte(0)
		if index+1 < len(contents) {
			next = contents[index+1]
		}

		if lineComment {
			statement.WriteByte(current)
			if current == '\n' {
				lineComment = false
			}
			continue
		}
		if blockComment {
			statement.WriteByte(current)
			if current == '*' && next == '/' {
				statement.WriteByte(next)
				index++
				blockComment = false
			}
			continue
		}
		if quote != 0 {
			statement.WriteByte(current)
			if current == '\\' && quote != '`' && next != 0 {
				statement.WriteByte(next)
				index++
				continue
			}
			if current == quote {
				if next == quote {
					statement.WriteByte(next)
					index++
				} else {
					quote = 0
				}
			}
			continue
		}

		switch {
		case current == '-' && next == '-':
			statement.WriteByte(current)
			statement.WriteByte(next)
			index++
			lineComment = true
		case current == '#':
			statement.WriteByte(current)
			lineComment = true
		case current == '/' && next == '*':
			statement.WriteByte(current)
			statement.WriteByte(next)
			index++
			blockComment = true
		case current == '\'' || current == '"' || current == '`':
			quote = current
			statement.WriteByte(current)
		case current == ';':
			if query := strings.TrimSpace(statement.String()); query != "" {
				statements = append(statements, query)
			}
			statement.Reset()
		default:
			statement.WriteByte(current)
		}
	}
	if quote != 0 || blockComment {
		return nil, fmt.Errorf("unterminated SQL quote or comment")
	}
	if query := strings.TrimSpace(statement.String()); query != "" {
		statements = append(statements, query)
	}
	return statements, nil
}
