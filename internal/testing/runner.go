package testing

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// TestScenario defines the interface for test scenarios
type TestScenario interface {
	// Name returns the scenario name
	Name() string

	// Description returns a description of the scenario
	Description() string

	// Setup prepares the scenario for execution
	Setup(ctx context.Context) error

	// Execute runs the main scenario logic
	Execute(ctx context.Context) error

	// Verify checks that the scenario executed correctly
	Verify(ctx context.Context) error

	// Cleanup cleans up after the scenario
	Cleanup(ctx context.Context) error
}

// TestRunner orchestrates test execution
type TestRunner struct {
	config   TestConfig
	tracer   *TraceCollector
	smtpConn *smtp.Client
	imapConn *client.Client
}

// NewTestRunner creates a new test runner
func NewTestRunner(config TestConfig) *TestRunner {
	return &TestRunner{
		config: config,
		tracer: NewTraceCollector(),
	}
}

// Run executes a test scenario
func (r *TestRunner) Run(ctx context.Context, scenario TestScenario) (*TestResult, error) {
	result := NewTestResult(scenario.Name(), scenario.Description())

	// Setup phase
	setupTrace := r.tracer.StartWithDetails("setup", "test_runner", "initialize")
	if err := scenario.Setup(ctx); err != nil {
		setupTrace.EndWithError(err)
		result.Complete(false, "Setup failed: "+err.Error())
		result.AddError(err)
		return result, err
	}
	setupTrace.End()

	// Connect to services
	if err := r.connectServices(ctx); err != nil {
		result.Complete(false, "Service connection failed: "+err.Error())
		result.AddError(err)
		return result, err
	}

	// Execute phase
	executeTrace := r.tracer.StartWithDetails("execute", "test_runner", "run_scenario")
	err := scenario.Execute(ctx)
	if err != nil {
		executeTrace.EndWithError(err)
		result.Complete(false, "Execution failed: "+err.Error())
		result.AddError(err)
		// Still run cleanup even if execution failed
		r.cleanupScenario(ctx, scenario, result)
		return result, err
	}
	executeTrace.End()

	// Verify phase
	verifyTrace := r.tracer.StartWithDetails("verify", "test_runner", "verify_results")
	verifyErr := scenario.Verify(ctx)
	if verifyErr != nil {
		verifyTrace.EndWithError(verifyErr)
		result.Complete(false, "Verification failed: "+verifyErr.Error())
		result.AddError(verifyErr)
		r.cleanupScenario(ctx, scenario, result)
		return result, verifyErr
	}
	verifyTrace.End()

	// Cleanup phase
	r.cleanupScenario(ctx, scenario, result)

	// Success
	result.Complete(true, "Test passed successfully")
	return result, nil
}

// connectServices establishes connections to required services
func (r *TestRunner) connectServices(ctx context.Context) error {
	trace := r.tracer.StartWithDetails("connect", "test_runner", "connect_services")

	// Connect to SMTP server
	if err := r.connectSMTP(ctx); err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("SMTP connection failed: %w", err)
	}

	// Connect to IMAP server
	if err := r.connectIMAP(ctx); err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("IMAP connection failed: %w", err)
	}

	trace.End()
	return nil
}

// connectSMTP connects to the SMTP server
func (r *TestRunner) connectSMTP(ctx context.Context) error {
	trace := r.tracer.StartWithDetails("smtp_connect", "smtp", "dial")

	client, err := smtp.Dial(r.config.SMTPAddr)
	if err != nil {
		trace.EndWithError(err)
		return err
	}

	// Authenticate if credentials provided
	if r.config.Username != "" && r.config.Password != "" {
		auth := smtp.PlainAuth("", r.config.Username, r.config.Password, strings.Split(r.config.SMTPAddr, ":")[0])
		if err := client.Auth(auth); err != nil {
			client.Close()
			trace.EndWithError(err)
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	r.smtpConn = client
	trace.WithDetail("server", r.config.SMTPAddr).End()
	return nil
}

// connectIMAP connects to the IMAP server
func (r *TestRunner) connectIMAP(ctx context.Context) error {
	trace := r.tracer.StartWithDetails("imap_connect", "imap", "dial")

	c, err := client.Dial(r.config.IMAPAddr)
	if err != nil {
		trace.EndWithError(err)
		return err
	}

	// Login if credentials provided
	if r.config.Username != "" && r.config.Password != "" {
		if err := c.Login(r.config.Username, r.config.Password); err != nil {
			c.Close()
			trace.EndWithError(err)
			return fmt.Errorf("IMAP login failed: %w", err)
		}
	}

	r.imapConn = c
	trace.WithDetail("server", r.config.IMAPAddr).End()
	return nil
}

// cleanupScenario runs the scenario cleanup and disconnects services
func (r *TestRunner) cleanupScenario(ctx context.Context, scenario TestScenario, result *TestResult) {
	cleanupTrace := r.tracer.StartWithDetails("cleanup", "test_runner", "cleanup_scenario")

	// Run scenario cleanup
	if cleanupErr := scenario.Cleanup(ctx); cleanupErr != nil {
		cleanupTrace.WithDetail("scenario_cleanup_error", cleanupErr.Error())
		result.AddError(cleanupErr)
	}

	// Disconnect services
	r.disconnectServices()

	cleanupTrace.End()
}

// disconnectServices closes all service connections
func (r *TestRunner) disconnectServices() {
	if r.smtpConn != nil {
		r.smtpConn.Close()
		r.smtpConn = nil
	}

	if r.imapConn != nil {
		r.imapConn.Logout()
		r.imapConn.Close()
		r.imapConn = nil
	}
}

// SMTPClient returns the SMTP client for test scenarios
func (r *TestRunner) SMTPClient() *smtp.Client {
	return r.smtpConn
}

// IMAPClient returns the IMAP client for test scenarios
func (r *TestRunner) IMAPClient() *client.Client {
	return r.imapConn
}

// Config returns the test configuration
func (r *TestRunner) Config() TestConfig {
	return r.config
}

// Tracer returns the trace collector
func (r *TestRunner) Tracer() *TraceCollector {
	return r.tracer
}

// SendEmail is a convenience method for sending email via SMTP
func (r *TestRunner) SendEmail(from, to, subject, body string) error {
	if r.smtpConn == nil {
		return fmt.Errorf("SMTP client not connected")
	}

	trace := r.tracer.StartWithDetails("send_email", "smtp", "send")

	// Send MAIL FROM command
	if err := r.smtpConn.Mail(from); err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Send RCPT TO command
	if err := r.smtpConn.Rcpt(to); err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send DATA command and message content
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	wc, err := r.smtpConn.Data()
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to start data: %w", err)
	}

	_, err = fmt.Fprint(wc, msg)
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to write message data: %w", err)
	}

	err = wc.Close()
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	trace.WithDetails(map[string]interface{}{
		"from":      from,
		"to":        to,
		"subject":   subject,
		"body_size": len(body),
	}).End()

	return nil
}

// FetchLatestMessage is a convenience method for fetching the latest IMAP message
func (r *TestRunner) FetchLatestMessage() (*imap.Message, error) {
	if r.imapConn == nil {
		return nil, fmt.Errorf("IMAP client not connected")
	}

	trace := r.tracer.StartWithDetails("fetch_latest", "imap", "fetch")

	// Select INBOX
	_, err := r.imapConn.Select("INBOX", false)
	if err != nil {
		trace.EndWithError(err)
		return nil, fmt.Errorf("failed to select INBOX: %w", err)
	}

	// Get mailbox status
	mbox := r.imapConn.Mailbox()
	if mbox == nil || mbox.Messages == 0 {
		trace.EndWithError(fmt.Errorf("no messages in mailbox"))
		return nil, fmt.Errorf("no messages in mailbox")
	}

	// Fetch the latest message
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(mbox.Messages)

	messages := make(chan *imap.Message, 1)
	done := make(chan error, 1)

	go func() {
		done <- r.imapConn.Fetch(seqSet, []imap.FetchItem{imap.FetchEnvelope, imap.FetchRFC822Text}, messages)
	}()

	msg := <-messages
	err = <-done

	if err != nil {
		trace.EndWithError(err)
		return nil, fmt.Errorf("failed to fetch message: %w", err)
	}

	trace.WithDetail("message_uid", msg.Uid).End()
	return msg, nil
}

// WaitForMessage waits for a message to arrive in the mailbox
func (r *TestRunner) WaitForMessage(timeout time.Duration) (*imap.Message, error) {
	if r.imapConn == nil {
		return nil, fmt.Errorf("IMAP client not connected")
	}

	trace := r.tracer.StartWithDetails("wait_message", "imap", "wait")

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		// Check for new messages
		msg, err := r.FetchLatestMessage()
		if err == nil {
			trace.WithDetail("message_uid", msg.Uid).End()
			return msg, nil
		}

		// Wait a bit before checking again
		time.Sleep(100 * time.Millisecond)
	}

	trace.EndWithError(fmt.Errorf("timeout waiting for message"))
	return nil, fmt.Errorf("timeout waiting for message after %v", timeout)
}

// Close closes the test runner and cleans up resources
func (r *TestRunner) Close() error {
	r.disconnectServices()
	return nil
}
