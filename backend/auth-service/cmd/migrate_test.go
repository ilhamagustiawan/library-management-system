package cmd

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/golang-migrate/migrate/v4"
)

type fakeMigration struct {
	upErr     error
	upCalls   int
	downCalls int
	downStep  int
}

func (m *fakeMigration) Up() error {
	m.upCalls++
	return m.upErr
}

func (m *fakeMigration) Steps(step int) error {
	m.downCalls++
	m.downStep = step
	return nil
}

func TestApplyMigrationSeedsDevelopmentAfterNoChange(t *testing.T) {
	migration := &fakeMigration{upErr: migrate.ErrNoChange}
	seedCalls := 0

	err := applyMigration(context.Background(), migration, "up", "development", func(context.Context) error {
		seedCalls++
		return nil
	})

	if err != nil {
		t.Fatalf("applyMigration() error = %v", err)
	}
	if migration.upCalls != 1 || seedCalls != 1 {
		t.Fatalf("up calls = %d, seed calls = %d; want 1, 1", migration.upCalls, seedCalls)
	}
}

func TestApplyMigrationSkipsDevelopmentSeedInProduction(t *testing.T) {
	migration := &fakeMigration{}
	seedCalls := 0

	err := applyMigration(context.Background(), migration, "up", "production", func(context.Context) error {
		seedCalls++
		return nil
	})

	if err != nil {
		t.Fatalf("applyMigration() error = %v", err)
	}
	if seedCalls != 0 {
		t.Fatalf("seed calls = %d, want 0", seedCalls)
	}
}

func TestApplyMigrationDoesNotSeedAfterDown(t *testing.T) {
	migration := &fakeMigration{}
	seedCalls := 0

	err := applyMigration(context.Background(), migration, "down", "development", func(context.Context) error {
		seedCalls++
		return nil
	})

	if err != nil {
		t.Fatalf("applyMigration() error = %v", err)
	}
	if migration.downCalls != 1 || migration.downStep != -1 || seedCalls != 0 {
		t.Fatalf("down calls = %d, down step = %d, seed calls = %d; want 1, -1, 0", migration.downCalls, migration.downStep, seedCalls)
	}
}

func TestApplyMigrationDoesNotSeedAfterMigrationFailure(t *testing.T) {
	migrationErr := errors.New("migration failed")
	seedCalls := 0

	err := applyMigration(context.Background(), &fakeMigration{upErr: migrationErr}, "up", "development", func(context.Context) error {
		seedCalls++
		return nil
	})

	if !errors.Is(err, migrationErr) || seedCalls != 0 {
		t.Fatalf("applyMigration() error = %v, seed calls = %d; want migration error and 0", err, seedCalls)
	}
}

func TestApplyMigrationReportsSeedFailureAfterMigration(t *testing.T) {
	migration := &fakeMigration{}

	err := applyMigration(context.Background(), migration, "up", "development", func(context.Context) error {
		return errors.New("invalid seed")
	})

	if err == nil || !strings.Contains(err.Error(), "schema migrations remain applied") {
		t.Fatalf("applyMigration() error = %v, want applied-schema recovery context", err)
	}
}
