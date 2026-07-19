package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "transaction-service", Short: "Library loan and transaction service", SilenceUsage: true}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() { rootCmd.AddCommand(serveCmd, migrateCmd, swaggerCmd) }
