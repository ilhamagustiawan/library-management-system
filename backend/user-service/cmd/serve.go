package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/user-service/internal/server"
)

var serveCmd = &cobra.Command{
	Use: "serve", Short: "Start the user service", Aliases: []string{"server"},
	RunE: func(_ *cobra.Command, _ []string) error {
		application, err := server.New(context.Background())
		if err != nil {
			return err
		}
		return application.Start()
	},
}
