package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/api/http/response"
	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/config"
	authinfra "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/auth"
	infraDB "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db"
	userRepo "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/infra/db/repository/user"
	adminusecase "github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/usecase/admin"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create an administrator identity offline",
	RunE:  runCreateAdmin,
}

func runCreateAdmin(cmd *cobra.Command, _ []string) error {
	name, _ := cmd.Flags().GetString("name")
	email, _ := cmd.Flags().GetString("email")
	if _, err := fmt.Fprint(cmd.ErrOrStderr(), "Password: "); err != nil {
		return err
	}
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	_, _ = fmt.Fprintln(cmd.ErrOrStderr())
	if err != nil {
		return fmt.Errorf("read password from terminal: %w", err)
	}
	defer clear(password)

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	db, err := infraDB.Connect(context.Background(), infraDB.Config{
		DSN: cfg.Database.DSN, MaxOpenConns: cfg.Database.MaxOpenConns, MaxIdleConns: cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime, ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return err
	}
	defer db.Close()

	usecase := adminusecase.New(userRepo.NewRepository(db), authinfra.NewPasswordHasher(cfg.Auth.BcryptCost))
	user, err := usecase.Create(cmd.Context(), adminusecase.Input{Name: name, Email: email, Password: string(password)})
	if errors.Is(err, adminusecase.ErrExistingMember) {
		return fmt.Errorf("create admin: %w; existing member remains unchanged", err)
	}
	if err != nil {
		return err
	}
	return json.NewEncoder(cmd.OutOrStdout()).Encode(response.NewUser(user))
}

func init() {
	createAdminCmd.Flags().String("name", "", "administrator name")
	createAdminCmd.Flags().String("email", "", "administrator email")
}
