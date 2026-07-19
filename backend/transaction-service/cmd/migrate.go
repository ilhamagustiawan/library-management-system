package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/transaction-service/internal/config"
)

var migrateCmd = &cobra.Command{Use: "migrate", Short: "Apply or roll back transaction migrations", RunE: runMigration}

func runMigration(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	dir, _ := cmd.Flags().GetString("dir")
	action, _ := cmd.Flags().GetString("action")
	runner, err := migrate.New("file://"+dir, cfg.Database.MigrationURL)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer func() { _, _ = runner.Close() }()
	switch action {
	case "up":
		err = runner.Up()
	case "down":
		err = runner.Steps(-1)
	default:
		return fmt.Errorf("invalid migration action %q", action)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate %s: %w", action, err)
	}
	return nil
}

func init() {
	migrateCmd.Flags().String("dir", "./db/migrations", "migration directory")
	migrateCmd.Flags().String("action", "up", "migration action: up or down")
}
