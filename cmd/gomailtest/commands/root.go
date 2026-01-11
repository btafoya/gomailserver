package commands

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/btafoya/gomailserver/internal/testing/checks"
	"github.com/btafoya/gomailserver/internal/testing/types"
	"github.com/btafoya/gomailserver/internal/testing/verifier"
	"github.com/spf13/cobra"
)

func VerifyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Run all verification checks",
		Long:  "Run all verification checks against a gomailserver instance",
		Run: func(cmd *cobra.Command, args []string) {
			executeVerify(cmd)
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to gomailserver.conf (local mode)")
	cmd.Flags().StringP("profile", "p", "", "Profile name to use (remote mode)")
	cmd.Flags().BoolP("dry-run", "d", false, "Non-invasive checks only (no test messages)")
	cmd.Flags().BoolP("quiet", "q", false, "Quiet mode (exit code only, no output)")
	cmd.Flags().CountP("verbose", "v", "Verbose output with detailed logging")
	cmd.Flags().StringP("report-html", "", "", "Save HTML report to file")
	cmd.Flags().StringP("report-json", "", "", "Save JSON report to file")
	cmd.Flags().BoolP("warnings-ok", "w", false, "Don't fail on warnings")
	cmd.Flags().IntP("rate-limit", "r", 10, "Rate limit (operations per second)")
	cmd.Flags().BoolP("no-cleanup", "", false, "Don't auto-delete test messages")

	return cmd
}

func TestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [category]",
		Short: "Run specific check categories",
		Long:  "Run specific check categories: config, smtp, imap, mailflow, security",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			executeTest(cmd, args)
		},
	}

	cmd.Flags().StringP("config", "c", "", "Path to gomailserver.conf (local mode)")
	cmd.Flags().StringP("profile", "p", "", "Profile name to use (remote mode)")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output with detailed logging")

	return cmd
}

func executeVerify(cmd *cobra.Command) {
	config, serverConfig, err := loadConfig(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	v := verifier.NewVerifier(config, serverConfig)
	consoleReporter := verifier.NewConsoleReporter(config.OutputMode)

	v.InitializeChecks()

	if _, err := v.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error running verification: %v\n", err)
		os.Exit(1)
	}

	consoleReporter.PrintFinal(v.GetReport())

	if config.ReportHTML != "" {
		htmlReporter := verifier.NewHTMLReporter(config.ReportHTML)
		htmlReporter.PrintFinal(v.GetReport())
	}

	if config.ReportJSON != "" {
		jsonReporter := verifier.NewJSONReporter(config.ReportJSON)
		jsonReporter.PrintFinal(v.GetReport())
	}

	if v.GetReport().ExitCode != 0 {
		os.Exit(1)
	}
}

func executeTest(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: category required")
		fmt.Println("Available categories: config, smtp, imap, mailflow, security")
		os.Exit(1)
	}

	category := args[0]

	serverConfig, err := loadTestConfig(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var testCategory types.Category
	switch category {
	case "config":
		testCategory = types.CategoryConfig
	case "smtp":
		testCategory = types.CategoryMailFlow
	case "imap":
		testCategory = types.CategoryMailFlow
	case "mailflow":
		testCategory = types.CategoryMailFlow
	case "security":
		testCategory = types.CategorySecurity
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown category '%s'\n", category)
		fmt.Println("Available categories: config, smtp, imap, mailflow, security")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	cfg := verifier.DefaultConfig()
	cfg.OutputMode = types.OutputSummary

	if cmd.Flags().Changed("verbose") {
		cfg.OutputMode = types.OutputVerbose
	}

	registry := checks.NewRegistry()
	registry.RegisterAll()

	testChecks := registry.GetByCategory(testCategory)

	fmt.Printf("Running %s checks...\n", category)
	for _, check := range testChecks {
		fmt.Printf("  %s: ", check.Name())
		result, err := check.Run(ctx, serverConfig)
		if err != nil {
			fmt.Printf("FAILED (%v)\n", err)
		} else {
			fmt.Printf("%s\n", result.Status)
		}
	}
}

func loadConfig(cmd *cobra.Command) (*verifier.Config, *types.ServerConfig, error) {
	config := verifier.DefaultConfig()
	serverConfig := &types.ServerConfig{}

	configFile, _ := cmd.Flags().GetString("config")
	profileName, _ := cmd.Flags().GetString("profile")

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	quiet, _ := cmd.Flags().GetBool("quiet")
	verbose, _ := cmd.Flags().GetCount("verbose")
	reportHTML, _ := cmd.Flags().GetString("report-html")
	reportJSON, _ := cmd.Flags().GetString("report-json")
	warningsOk, _ := cmd.Flags().GetBool("warnings-ok")
	rateLimit, _ := cmd.Flags().GetInt("rate-limit")
	noCleanup, _ := cmd.Flags().GetBool("no-cleanup")

	config.ConfigFile = configFile
	config.ProfileName = profileName
	config.DryRun = dryRun

	if quiet {
		config.OutputMode = types.OutputQuiet
	} else if verbose > 0 {
		config.OutputMode = types.OutputVerbose
	} else {
		config.OutputMode = types.OutputSummary
	}

	config.ReportHTML = reportHTML
	config.ReportJSON = reportJSON
	config.WarningsOk = warningsOk
	config.RateLimit = rateLimit
	config.NoCleanup = noCleanup

	serverConfig.DryRun = dryRun
	serverConfig.AutoCleanup = !noCleanup

	if configFile != "" {
		config.Mode = verifier.ModeLocal
		serverConfig.SMTPHost = "localhost"
		serverConfig.SMTPPort = 587
		serverConfig.IMAPHost = "localhost"
		serverConfig.IMAPPort = 143
		serverConfig.Domains = []string{"localhost"}
		serverConfig.TLS = false
		serverConfig.Timeout = 30 * time.Second
	} else if profileName != "" {
		config.Mode = verifier.ModeRemote
		profile, err := verifier.LoadProfile(profileName)
		if err != nil {
			return nil, nil, fmt.Errorf("loading profile '%s': %w", profileName, err)
		}

		serverConfig.SMTPHost = profile.SMTPHost
		serverConfig.SMTPPort = profile.SMTPPort
		serverConfig.IMAPHost = profile.IMAPHost
		serverConfig.IMAPPort = profile.IMAPPort
		serverConfig.Domains = profile.Domains
		serverConfig.TestUser = profile.TestUser
		serverConfig.PasswordEnv = profile.PasswordEnv
		serverConfig.TLS = profile.Options.TLS
		serverConfig.StartTLS = profile.Options.StartTLS
		serverConfig.Timeout = profile.Options.Timeout
		serverConfig.ConfigPath = fmt.Sprintf("~/.gomailserver/profiles/%s.yaml", profileName)
	} else {
		config.Mode = verifier.ModeLocal
		serverConfig.SMTPHost = "localhost"
		serverConfig.SMTPPort = 587
		serverConfig.IMAPHost = "localhost"
		serverConfig.IMAPPort = 143
		serverConfig.Domains = []string{"localhost"}
		serverConfig.TLS = false
		serverConfig.Timeout = 30 * time.Second
	}

	return config, serverConfig, nil
}

func loadTestConfig(cmd *cobra.Command) (*types.ServerConfig, error) {
	profileName, _ := cmd.Flags().GetString("profile")

	if profileName != "" {
		profile, err := verifier.LoadProfile(profileName)
		if err != nil {
			return nil, fmt.Errorf("loading profile '%s': %w", profileName, err)
		}

		serverConfig := &types.ServerConfig{
			SMTPHost: profile.SMTPHost,
			SMTPPort: profile.SMTPPort,
			IMAPHost: profile.IMAPHost,
			IMAPPort: profile.IMAPPort,
			Domains:  profile.Domains,
			TestUser: profile.TestUser,
			TLS:      profile.Options.TLS,
			StartTLS: profile.Options.StartTLS,
			Timeout:  profile.Options.Timeout,
		}

		return serverConfig, nil
	}

	serverConfig := &types.ServerConfig{
		SMTPHost: "localhost",
		SMTPPort: 587,
		IMAPHost: "localhost",
		IMAPPort: 143,
		Domains:  []string{"localhost"},
		TLS:      false,
		Timeout:  30 * time.Second,
	}

	return serverConfig, nil
}
