package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var swaggerCmd = &cobra.Command{
	Use: "swagger", Short: "Generate Swagger documentation", Args: cobra.NoArgs,
	RunE: func(command *cobra.Command, _ []string) error {
		generator := exec.CommandContext(command.Context(), "go", "run", "github.com/swaggo/swag/cmd/swag", "init", "--parseInternal", "--output", "docs")
		generator.Stdout = command.OutOrStdout()
		generator.Stderr = command.ErrOrStderr()
		if err := generator.Run(); err != nil {
			return fmt.Errorf("generate Swagger documentation: %w", err)
		}
		return nil
	},
}
