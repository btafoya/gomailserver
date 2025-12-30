package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	cfgFile string
	version string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "gomailserver",
	Short: "A modern, all-in-one mail server written in Go",
	Long: `gomailserver is a composable, all-in-one mail server implementing
SMTP, IMAP, CalDAV, CardDAV with built-in security features including
DKIM, SPF, DMARC, antivirus, and anti-spam capabilities.`,
}

// Execute runs the root command
func Execute(v string) error {
	version = v
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is /etc/gomailserver/gomailserver.yaml)")

	rootCmd.Version = version
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gomailserver version %s\n", version)
	},
}
