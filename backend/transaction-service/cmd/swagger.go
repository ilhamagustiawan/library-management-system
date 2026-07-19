package cmd

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/spf13/cobra"
)

type swaggerGenerator func(context.Context, io.Writer, io.Writer) error

var swaggerCmd = newSwaggerCommand(generateSwagger)

func newSwaggerCommand(generate swaggerGenerator) *cobra.Command {
	return &cobra.Command{
		Use:   "swagger",
		Short: "Generate Swagger documentation",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return generate(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}
}

func generateSwagger(ctx context.Context, stdout, stderr io.Writer) error {
	command := exec.CommandContext(
		ctx,
		"go",
		"run",
		"github.com/swaggo/swag/cmd/swag",
		"init",
		"--parseInternal",
		"--output",
		"docs",
	)
	command.Stdout = stdout
	command.Stderr = stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("generate Swagger documentation: %w; inspect generator output, resolve the cause, then rerun", err)
	}
	return nil
}
