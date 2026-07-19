package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/ilhamagustiawan/library-management-system/backend/book-service/internal/server"
)

var serveCmd = &cobra.Command{Use: "serve", Short: "Start the book service", RunE: runServe}

func runServe(_ *cobra.Command, _ []string) error {
	app, err := server.New(context.Background())
	if err != nil {
		return err
	}
	return app.Start()
}
