package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/config"
)

var migrateCmd = &cobra.Command{Use: "migrate", Short: "Apply or roll back database migrations", RunE: runMigration}

type migrationAction interface {
	Up() error
	Steps(int) error
}

func runMigration(command *cobra.Command, _ []string) error {
	configValue, err := config.Load()
	if err != nil {
		return err
	}
	directory, _ := command.Flags().GetString("dir")
	action, _ := command.Flags().GetString("action")
	migration, err := migrate.New("file://"+directory, configValue.Database.MigrationURL)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer func() { _, _ = migration.Close() }()
	return applyMigration(migration, action)
}

func applyMigration(migration migrationAction, action string) error {
	switch action {
	case "up":
		if err := migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
		return nil
	case "down":
		if err := migration.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate down: %w", err)
		}
		return nil
	default:
		return fmt.Errorf("invalid migration action %q", action)
	}
}

func init() {
	migrateCmd.Flags().String("dir", "./db/migrations", "migration directory")
	migrateCmd.Flags().String("action", "up", "migration action: up or down")
}
