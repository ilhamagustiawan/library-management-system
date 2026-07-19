package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/auth-service/internal/server"
)

var serveCmd = &cobra.Command{
	Use:     "serve",
	Short:   "Start the OAuth service",
	Aliases: []string{"server"},
	RunE:    runServe,
}

func runServe(_ *cobra.Command, _ []string) error {
	app, err := server.New(context.Background())
	if err != nil {
		return err
	}
	return app.Start()
}
