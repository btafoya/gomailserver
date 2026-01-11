package scenarios

import (
	"context"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/testing"
)

// BasicDeliveryTest implements a basic email delivery test scenario
type BasicDeliveryTest struct {
	*BaseScenario

	from      string
	to        string
	subject   string
	body      string
	messageID string
	runner    *testing.TestRunner
}

// NewBasicDeliveryTest creates a new basic delivery test
func NewBasicDeliveryTest(from, to, subject, body string) *BasicDeliveryTest {
	return &BasicDeliveryTest{
		BaseScenario: NewBaseScenario(
			"basic_delivery",
			"Tests basic email delivery from SMTP to IMAP",
			nil, // Will be set by runner
		),
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

// Setup prepares the test scenario
func (t *BasicDeliveryTest) Setup(ctx context.Context) error {
	// Generate a unique message ID for this test
	t.messageID = fmt.Sprintf("test-%d@example.com", time.Now().UnixNano())

	// Get the test runner (injected during execution)
	// This would be set by the runner when executing

	return nil
}

// Execute runs the main test logic
func (t *BasicDeliveryTest) Execute(ctx context.Context) error {
	if t.runner == nil {
		return fmt.Errorf("test runner not set")
	}

	// Phase 1: Send email via SMTP
	trace := t.Trace().StartWithComponent("smtp_send", "smtp")
	err := t.runner.SendEmail(t.from, t.to, t.subject, t.body)
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	trace.End()

	// Phase 2: Wait for message to be processed (brief pause for async processing)
	trace = t.Trace().StartWithComponent("wait_processing", "queue")
	time.Sleep(500 * time.Millisecond) // Allow time for queue processing
	trace.End()

	// Phase 3: Fetch message via IMAP
	trace = t.Trace().StartWithComponent("imap_fetch", "imap")
	message, err := t.runner.WaitForMessage(5 * time.Second)
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to fetch message via IMAP: %w", err)
	}

	// Store the fetched message for verification
	t.messageID = fmt.Sprintf("imap-%d", message.Uid)
	trace.WithDetail("message_uid", message.Uid).End()

	return nil
}

// Verify checks that the test executed correctly
func (t *BasicDeliveryTest) Verify(ctx context.Context) error {
	if t.runner == nil {
		return fmt.Errorf("test runner not set")
	}

	// Verify message was received via IMAP
	trace := t.Trace().StartWithComponent("verify_delivery", "imap")

	message, err := t.runner.FetchLatestMessage()
	if err != nil {
		trace.EndWithError(err)
		return fmt.Errorf("failed to verify message delivery: %w", err)
	}

	// Check message content
	if message.Envelope.Subject != t.subject {
		trace.EndWithError(fmt.Errorf("subject mismatch"))
		return fmt.Errorf("subject mismatch: expected %q, got %q", t.subject, message.Envelope.Subject)
	}

	// Check sender
	if len(message.Envelope.From) == 0 || message.Envelope.From[0].Address() != t.from {
		trace.EndWithError(fmt.Errorf("sender mismatch"))
		return fmt.Errorf("sender mismatch: expected %q", t.from)
	}

	// Check recipient
	if len(message.Envelope.To) == 0 || message.Envelope.To[0].Address() != t.to {
		trace.EndWithError(fmt.Errorf("recipient mismatch"))
		return fmt.Errorf("recipient mismatch: expected %q", t.to)
	}

	trace.WithDetails(map[string]interface{}{
		"subject_match":   true,
		"sender_match":    true,
		"recipient_match": true,
		"message_uid":     message.Uid,
	}).End()

	return nil
}

// Cleanup performs cleanup after the test
func (t *BasicDeliveryTest) Cleanup(ctx context.Context) error {
	// Cleanup would be handled by the test runner
	// For now, just log completion
	if t.Trace() != nil {
		trace := t.Trace().StartWithComponent("cleanup", "test")
		trace.End()
	}
	return nil
}

// SetRunner sets the test runner for this scenario
func (t *BasicDeliveryTest) SetRunner(runner *testing.TestRunner) {
	t.runner = runner
	t.BaseScenario = NewBaseScenario(
		"basic_delivery",
		"Tests basic email delivery from SMTP to IMAP",
		runner.Tracer(),
	)
}
