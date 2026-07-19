package cmd

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/config"
)

var migrateCmd = &cobra.Command{Use: "migrate", Short: "Apply or roll back database migrations", RunE: runMigration}

func runMigration(command *cobra.Command, _ []string) error {
	serviceConfig, err := config.Load()
	if err != nil {
		return err
	}
	directory, _ := command.Flags().GetString("dir")
	action, _ := command.Flags().GetString("action")
	runner, err := migrate.New("file://"+directory, serviceConfig.Database.MigrationURL)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer runner.Close()
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
