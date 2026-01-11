package main

import (
	"os"

	"github.com/btafoya/gomailserver/cmd/gomailtest/commands"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gomailtest",
		Short: "GoMailServer verification and testing tool",
		Long:  "Verify gomailserver configuration, connectivity, and mail flow",
	}

	rootCmd.AddCommand(commands.VerifyCmd())
	rootCmd.AddCommand(commands.TestCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
