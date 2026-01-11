package main

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/btafoya/gomailserver/internal/testing"
	"github.com/btafoya/gomailserver/internal/testing/reporters"
	"github.com/btafoya/gomailserver/internal/testing/scenarios"
)

// mailtestCmd represents the mailtest command
type mailtestCmd struct {
	// Server configuration
	smtpAddr string
	imapAddr string
	httpAddr string

	// Authentication
	username string
	password string

	// Test configuration
	outputDir  string
	htmlReport bool
	jsonReport bool
	verbose    bool
	debug      bool

	// Test parameters
	testName string
	fromAddr string
	toAddr   string
	subject  string
	body     string
}

// newMailtestCmd creates a new mailtest command
func newMailtestCmd() *mailtestCmd {
	return &mailtestCmd{
		smtpAddr:   "localhost:587",
		imapAddr:   "localhost:143",
		httpAddr:   "http://localhost:8980",
		username:   "test@example.com",
		password:   "password",
		outputDir:  "./test-reports",
		htmlReport: true,
		jsonReport: false,
		verbose:    false,
		debug:      false,
		fromAddr:   "test@example.com",
		toAddr:     "test@example.com",
		subject:    "Test Email",
		body:       "This is a test email sent by gomailserver test suite.",
	}
}

// run executes the mailtest command
func (cmd *mailtestCmd) run(args []string) error {
	// Parse arguments
	if len(args) > 0 {
		cmd.testName = args[0]
	}

	// Create test configuration
	config := testing.TestConfig{
		SMTPAddr:     cmd.smtpAddr,
		IMAPAddr:     cmd.imapAddr,
		HTTPAddr:     cmd.httpAddr,
		DatabasePath: "./data/mailserver.db",
		Username:     cmd.username,
		Password:     cmd.password,
		OutputDir:    cmd.outputDir,
		HTMLReport:   cmd.htmlReport,
		JSONReport:   cmd.jsonReport,
		Verbose:      cmd.verbose,
		Debug:        cmd.debug,
		Timeout:      30 * time.Second,
		MaxRetries:   3,
		RetryDelay:   1 * time.Second,
	}

	// Create test runner
	runner := testing.NewTestRunner(config)

	// Create test scenario
	var scenario testing.TestScenario
	switch cmd.testName {
	case "basic", "delivery", "":
		basicTest := scenarios.NewBasicDeliveryTest(
			cmd.fromAddr,
			cmd.toAddr,
			cmd.subject,
			cmd.body,
		)
		basicTest.SetRunner(runner)
		scenario = basicTest
	default:
		return fmt.Errorf("unknown test: %s. Available tests: basic, delivery", cmd.testName)
	}

	// Run the test
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	if cmd.verbose {
		fmt.Printf("Starting test: %s\n", scenario.Name())
		fmt.Printf("SMTP: %s, IMAP: %s\n", config.SMTPAddr, config.IMAPAddr)
		fmt.Printf("From: %s, To: %s\n", cmd.fromAddr, cmd.toAddr)
	}

	result, err := runner.Run(ctx, scenario)
	if err != nil {
		return fmt.Errorf("test execution failed: %w", err)
	}

	// Generate reports
	if cmd.htmlReport {
		htmlReporter := reporters.NewHTMLReporter()

		// Generate detailed report
		reportPath := filepath.Join(config.OutputDir,
			fmt.Sprintf("%s-%s.html", result.Name, time.Now().Format("20060102-150405")))

		if err := htmlReporter.GenerateReport(result, reportPath); err != nil {
			log.Printf("Failed to generate HTML report: %v", err)
		} else if cmd.verbose {
			fmt.Printf("HTML report generated: %s\n", reportPath)
		}

		// Also generate quick report to stdout if verbose
		if cmd.verbose {
			quickReport, err := htmlReporter.GenerateQuickReport(result)
			if err == nil {
				fmt.Println("\nQuick Report:")
				fmt.Println(quickReport)
			}
		}
	}

	// Print result summary
	fmt.Printf("\nTest Result: %s\n", result.Name)
	fmt.Printf("Status: %s\n", map[bool]string{true: "PASSED", false: "FAILED"}[result.Passed])
	fmt.Printf("Duration: %v\n", result.Duration)
	fmt.Printf("Trace Events: %d\n", len(result.Trace))

	if !result.Passed {
		fmt.Printf("Errors: %d\n", len(result.Errors))
		for i, err := range result.Errors {
			fmt.Printf("  %d. %v\n", i+1, err)
		}
		return fmt.Errorf("test failed")
	}

	if cmd.verbose {
		fmt.Printf("Summary: %s\n", result.Summary)
	}

	return nil
}

// printUsage prints the command usage
func (cmd *mailtestCmd) printUsage() {
	fmt.Println("Usage: gomailserver test [options] [test-name]")
	fmt.Println()
	fmt.Println("Run mail delivery tests against a running gomailserver instance.")
	fmt.Println()
	fmt.Println("Available tests:")
	fmt.Println("  basic, delivery  - Basic email delivery test (SMTP → IMAP)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --smtp-addr string     SMTP server address (default \"localhost:587\")")
	fmt.Println("  --imap-addr string     IMAP server address (default \"localhost:143\")")
	fmt.Println("  --http-addr string     HTTP API address (default \"http://localhost:8980\")")
	fmt.Println("  --username string      Authentication username (default \"test@example.com\")")
	fmt.Println("  --password string      Authentication password (default \"password\")")
	fmt.Println("  --from string          From email address (default \"test@example.com\")")
	fmt.Println("  --to string            To email address (default \"test@example.com\")")
	fmt.Println("  --subject string       Email subject (default \"Test Email\")")
	fmt.Println("  --body string          Email body (default \"This is a test email...\")")
	fmt.Println("  --output-dir string    Output directory for reports (default \"./test-reports\")")
	fmt.Println("  --html-report          Generate HTML report (default true)")
	fmt.Println("  --no-html-report       Disable HTML report generation")
	fmt.Println("  --json-report          Generate JSON report")
	fmt.Println("  --verbose, -v          Verbose output")
	fmt.Println("  --debug                Debug output")
	fmt.Println("  --help, -h             Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gomailserver test basic")
	fmt.Println("  gomailserver test --smtp-addr mail.example.com:587 --verbose")
	fmt.Println("  gomailserver test basic --from sender@example.com --to recipient@example.com")
}
