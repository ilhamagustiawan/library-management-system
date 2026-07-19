package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
	infraDB "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db"
)

const developmentOAuthSeedPath = "./db/seeds/oauth_clients.sql"

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Apply or roll back database migrations",
	RunE:  runMigration,
}

func runMigration(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	dir, _ := cmd.Flags().GetString("dir")
	action, _ := cmd.Flags().GetString("action")
	migration, err := migrate.New("file://"+dir, cfg.Database.MigrationURL)
	if err != nil {
		return fmt.Errorf("create migration runner: %w", err)
	}
	defer closeMigration(migration)

	return applyMigration(cmd.Context(), migration, action, cfg.Service.Environment, func(ctx context.Context) error {
		db, err := infraDB.Connect(ctx, infraDB.Config{
			DSN: cfg.Database.DSN, MaxOpenConns: cfg.Database.MaxOpenConns, MaxIdleConns: cfg.Database.MaxIdleConns,
			ConnMaxLifetime: cfg.Database.ConnMaxLifetime, ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
		})
		if err != nil {
			return err
		}
		defer db.Close()
		return infraDB.ApplySeedFile(ctx, db, developmentOAuthSeedPath)
	})
}

type migrationAction interface {
	Up() error
	Steps(int) error
}

func applyMigration(
	ctx context.Context,
	migration migrationAction,
	action string,
	environment string,
	seed func(context.Context) error,
) error {
	switch action {
	case "up":
		err := migration.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("migrate up: %w", err)
		}
		if environment == "development" {
			if err := seed(ctx); err != nil {
				return fmt.Errorf("apply development seed: %w; schema migrations remain applied; fix the seed and rerun migrate up", err)
			}
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

func closeMigration(migration *migrate.Migrate) {
	_, _ = migration.Close()
}

func init() {
	migrateCmd.Flags().String("dir", "./db/migrations", "migration directory")
	migrateCmd.Flags().String("action", "up", "migration action: up or down")
}
