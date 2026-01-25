package smtp

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	netmail "net/smtp"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/emersion/go-message"
	"github.com/miekg/dns"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	repService "github.com/btafoya/gomailserver/internal/reputation/service"
	"github.com/btafoya/gomailserver/internal/service"
)

// DeliveryWorker handles outbound SMTP message delivery
type DeliveryWorker struct {
	queueService     *service.QueueService
	domainRepo       repository.DomainRepository
	telemetryService *repService.TelemetryService
	logger           *zap.Logger
	config           *Config
	retryDelay       time.Duration
	maxRetries       int
}

// Config holds delivery worker configuration
type Config struct {
	Hostname    string
	DNSServer   string
	TLSMode     string // "none", "starttls", "tls"
	Username    string
	Password    string
	InsecureTLS bool
}

// NewDeliveryWorker creates a new delivery worker
func NewDeliveryWorker(
	queueService *service.QueueService,
	domainRepo repository.DomainRepository,
	telemetryService *repService.TelemetryService,
	logger *zap.Logger,
	config *Config,
) *DeliveryWorker {
	return &DeliveryWorker{
		queueService:     queueService,
		domainRepo:       domainRepo,
		telemetryService: telemetryService,
		logger:           logger,
		config:           config,
		retryDelay:       5 * time.Minute,
		maxRetries:       3,
	}
}

// ProcessQueue processes pending queue items
func (w *DeliveryWorker) ProcessQueue(ctx context.Context) error {
	// Get all pending queue items
	items, err := w.queueService.GetPendingItems(ctx)
	if err != nil {
		w.logger.Error("failed to get pending queue items", zap.Error(err))
		return err
	}

	w.logger.Info("processing delivery queue",
		zap.Int("pending_count", len(items)),
	)

	// Process each item
	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := w.processItem(ctx, item); err != nil {
				w.logger.Error("failed to process queue item",
					zap.Int64("id", item.ID),
					zap.String("sender", item.Sender),
					zap.Error(err),
				)

				// Handle retry vs permanent failure
				if item.RetryCount >= item.MaxRetries {
					// Mark as permanently failed
					if failErr := w.handlePermanentFailure(ctx, item, err); failErr != nil {
						w.logger.Error("failed to mark item as failed",
							zap.Int64("id", item.ID),
							zap.Error(failErr),
						)
					}
				} else {
					// Schedule retry
					if retryErr := w.scheduleRetry(ctx, item); retryErr != nil {
						w.logger.Error("failed to schedule retry",
							zap.Int64("id", item.ID),
							zap.Error(retryErr),
						)
					}
				}
			} else {
				// Mark as successfully delivered
				if delErr := w.handleSuccess(ctx, item); delErr != nil {
					w.logger.Error("failed to mark item as delivered",
						zap.Int64("id", item.ID),
						zap.Error(delErr),
					)
				}
			}
		}
	}

	return nil
}

// processItem handles delivery of a single queue item
func (w *DeliveryWorker) processItem(ctx context.Context, item *domain.QueueItem) error {
	w.logger.Debug("processing queue item",
		zap.Int64("id", item.ID),
		zap.String("sender", item.Sender),
	)

	// Parse message from file
	msg, err := w.readMessage(item.MessagePath)
	if err != nil {
		return fmt.Errorf("failed to read message: %w", err)
	}

	// Get recipients from stored JSON
	recipients := w.parseRecipients(item.Recipients)
	if len(recipients) == 0 {
		return fmt.Errorf("no recipients in queue item")
	}

	// For each recipient domain, attempt delivery
	for _, recipient := range recipients {
		recipientDomain := extractDomain(recipient)
		if recipientDomain == "" {
			w.logger.Warn("invalid recipient format",
				zap.String("recipient", recipient),
			)
			continue
		}

		// Check if this is local domain
		if w.isLocalDomain(recipientDomain) {
			// Handle local delivery
			if err := w.deliverLocal(ctx, recipient, msg); err != nil {
				w.logger.Warn("local delivery failed",
					zap.String("recipient", recipient),
					zap.Error(err),
				)
			}
			continue
		}

		// Attempt remote delivery
		if err := w.deliverRemote(ctx, recipientDomain, recipient, msg); err != nil {
			w.logger.Warn("remote delivery failed",
				zap.String("recipient", recipient),
				zap.String("domain", recipientDomain),
				zap.Error(err),
			)
			return err // Return first error for retry logic
		}
	}

	return nil
}

// deliverLocal handles delivery to local recipients
func (w *DeliveryWorker) deliverLocal(ctx context.Context, recipient string, msg *message.Entity) error {
	w.logger.Debug("delivering locally",
		zap.String("recipient", recipient),
	)

	// For local delivery, we can use the MessageService directly
	// This would store the message directly in the recipient's mailbox
	// TODO: Implement local delivery using MessageService
	return fmt.Errorf("local delivery not yet implemented")
}

// deliverRemote handles delivery to external domains
func (w *DeliveryWorker) deliverRemote(ctx context.Context, domain, recipient string, msg *message.Entity) error {
	w.logger.Debug("delivering remotely",
		zap.String("domain", domain),
		zap.String("recipient", recipient),
	)

	// Lookup MX records for recipient domain
	mxServers, err := w.lookupMX(domain)
	if err != nil {
		return fmt.Errorf("MX lookup failed for %s: %w", domain, err)
	}

	if len(mxServers) == 0 {
		return fmt.Errorf("no MX records found for domain %s", domain)
	}

	// Try each MX server in order
	for _, mx := range mxServers {
		err := w.deliverViaMX(ctx, mx, recipient, msg)
		if err == nil {
			return nil // Success
		}
		w.logger.Warn("delivery via MX failed, trying next",
			zap.String("mx", mx),
			zap.String("recipient", recipient),
			zap.Error(err),
		)
	}

	return fmt.Errorf("all MX servers failed for domain %s", domain)
}

// lookupMX performs DNS MX record lookup
func (w *DeliveryWorker) lookupMX(domain string) ([]string, error) {
	client := &dns.Client{
		Net:          "tcp",
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	m.SetEdns0(4096, false)

	// Use configured DNS server or default
	dnsServer := w.config.DNSServer
	if dnsServer == "" {
		dnsServer = "8.8.8.8:53"
	}

	reply, _, err := client.Exchange(m, dnsServer)
	if err != nil {
		return nil, fmt.Errorf("DNS query failed: %w", err)
	}

	var mxServers []string
	for _, ans := range reply.Answer {
		if mx, ok := ans.(*dns.MX); ok {
			mxServers = append(mxServers, mx.Mx)
		}
	}

	if len(mxServers) == 0 {
		return nil, fmt.Errorf("no MX records found")
	}

	return mxServers, nil
}

// deliverViaMX attempts delivery via a specific MX server
func (w *DeliveryWorker) deliverViaMX(ctx context.Context, mx, recipient string, msg *message.Entity) error {
	// Convert message to bytes for sending
	var buf bytes.Buffer
	if err := msg.WriteTo(&buf); err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	// Create SMTP client
	client, err := w.createSMTPClient(mx)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	// Send message using standard library
	return w.sendWithClient(client, msg, recipient)
}

// createSMTPClient creates an SMTP client with appropriate TLS configuration
func (w *DeliveryWorker) createSMTPClient(mx string) (*netmail.Client, error) {
	switch w.config.TLSMode {
	case "none":
		return netmail.Dial(mx + ":25")
	case "tls":
		tlsConfig := &tls.Config{
			InsecureSkipVerify: w.config.InsecureTLS,
			ServerName:         mx,
		}
		// Dial with TLS for implicit TLS (port 465)
		conn, err := tls.Dial("tcp", mx+":465", tlsConfig)
		if err != nil {
			return nil, fmt.Errorf("TLS dial failed: %w", err)
		}
		return netmail.NewClient(conn, mx)
	case "starttls":
		client, err := netmail.Dial(mx + ":587")
		if err != nil {
			return nil, err
		}

		tlsConfig := &tls.Config{
			InsecureSkipVerify: w.config.InsecureTLS,
			ServerName:         mx,
		}

		if ok, _ := client.Extension("STARTTLS"); ok {
			if err := client.StartTLS(tlsConfig); err != nil {
				client.Close()
				return nil, fmt.Errorf("STARTTLS failed: %w", err)
			}
		}
		return client, nil
	default:
		return nil, fmt.Errorf("unknown TLS mode: %s", w.config.TLSMode)
	}
}

// sendWithClient sends message using standard Go SMTP client
func (w *DeliveryWorker) sendWithClient(client *netmail.Client, msg *message.Entity, recipient string) error {
	// Extract sender from message
	sender := msg.Header.Get("From")
	if sender == "" {
		return fmt.Errorf("no From header in message")
	}

	// Authenticate if configured
	if w.config.Username != "" && w.config.Password != "" {
		auth := smtp.PlainAuth("", w.config.Username, w.config.Password, w.config.Hostname)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP authentication failed: %w", err)
		}
	}

	// Set sender
	if err := client.Mail(sender); err != nil {
		return fmt.Errorf("MAIL FROM failed: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(recipient); err != nil {
		return fmt.Errorf("RCPT TO failed: %w", err)
	}

	// Send message data
	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed: %w", err)
	}

	// Convert message to RFC822 format
	var buf bytes.Buffer
	if err := msg.WriteTo(&buf); err != nil {
		wc.Close()
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	if _, err := wc.Write(buf.Bytes()); err != nil {
		wc.Close()
		return fmt.Errorf("failed to write message data: %w", err)
	}

	wc.Close()
	return nil
}

// readMessage reads message from file
func (w *DeliveryWorker) readMessage(path string) (*message.Entity, error) {
	if path == "" {
		return nil, fmt.Errorf("empty message path")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read message file: %w", err)
	}

	reader := strings.NewReader(string(data))
	msg, err := message.Read(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message: %w", err)
	}
	return msg, nil
}

// parseRecipients parses recipient list from JSON string
func (w *DeliveryWorker) parseRecipients(recipientsJSON string) []string {
	if recipientsJSON == "" {
		return []string{}
	}

	var recipients []string
	if err := json.Unmarshal([]byte(recipientsJSON), &recipients); err != nil {
		w.logger.Warn("failed to parse recipients JSON",
			zap.Error(err),
			zap.String("json", recipientsJSON),
		)
		return []string{}
	}
	return recipients
}

// isLocalDomain checks if domain is configured as local
func (w *DeliveryWorker) isLocalDomain(domain string) bool {
	// TODO: Check against configured local domains
	// For now, return false
	return false
}

// handleSuccess marks a queue item as successfully delivered
func (w *DeliveryWorker) handleSuccess(ctx context.Context, item *domain.QueueItem) error {
	// Mark as delivered in queue
	if err := w.queueService.MarkDelivered(item.ID); err != nil {
		return err
	}

	// Record delivery telemetry
	if w.telemetryService != nil {
		for _, recipient := range w.parseRecipients(item.Recipients) {
			recipientDomain := extractDomain(recipient)
			senderDomain := extractDomain(item.Sender)

			if recipientDomain != "" && senderDomain != "" {
				if err := w.queueService.RecordDeliveryTelemetry(ctx, senderDomain, recipientDomain, ""); err != nil {
					w.logger.Warn("failed to record delivery telemetry",
						zap.Error(err),
					)
				}
			}
		}
	}

	w.logger.Info("message delivered successfully",
		zap.Int64("id", item.ID),
		zap.String("sender", item.Sender),
	)

	return nil
}

// handlePermanentFailure marks item as failed and generates bounce
func (w *DeliveryWorker) handlePermanentFailure(ctx context.Context, item *domain.QueueItem, deliveryErr error) error {
	// Mark as failed in queue
	errorMsg := deliveryErr.Error()
	if err := w.queueService.MarkFailed(item.ID, errorMsg); err != nil {
		return err
	}

	// Generate bounce message
	if err := w.generateBounce(ctx, item, deliveryErr); err != nil {
		w.logger.Error("failed to generate bounce message",
			zap.Error(err),
		)
	}

	// Record bounce telemetry
	if w.telemetryService != nil {
		recipients := w.parseRecipients(item.Recipients)
		for _, recipient := range recipients {
			recipientDomain := extractDomain(recipient)
			senderDomain := extractDomain(item.Sender)

			if recipientDomain != "" && senderDomain != "" {
				if err := w.queueService.RecordBounceTelemetry(ctx, senderDomain, recipientDomain, "", "permanent", "550", errorMsg); err != nil {
					w.logger.Warn("failed to record bounce telemetry",
						zap.Error(err),
					)
				}
			}
		}
	}

	w.logger.Info("message marked as permanently failed",
		zap.Int64("id", item.ID),
		zap.String("error", errorMsg),
	)

	return nil
}

// scheduleRetry schedules next retry attempt
func (w *DeliveryWorker) scheduleRetry(ctx context.Context, item *domain.QueueItem) error {
	failedAt := time.Now() // Use current time since FailedAt field doesn't exist
	nextRetry := w.queueService.CalculateNextRetry(item.RetryCount, failedAt)
	if nextRetry.IsZero() {
		// No more retries
		return w.handlePermanentFailure(ctx, item, fmt.Errorf("maximum retries exceeded"))
	}

	if err := w.queueService.IncrementRetry(item.ID, item.RetryCount, failedAt); err != nil {
		return err
	}

	w.logger.Info("item scheduled for retry",
		zap.Int64("id", item.ID),
		zap.Time("next_retry", nextRetry),
		zap.Int("retry_count", item.RetryCount+1),
	)

	return nil
}

// generateBounce creates a bounce message (DSN)
func (w *DeliveryWorker) generateBounce(ctx context.Context, item *domain.QueueItem, deliveryErr error) error {
	w.logger.Debug("generating bounce message",
		zap.Int64("id", item.ID),
	)

	// Parse original message
	originalMsg, err := w.readMessage(item.MessagePath)
	if err != nil {
		return fmt.Errorf("failed to read original message for bounce: %w", err)
	}

	// Build a simple DSN bounce message
	originalMessageID := originalMsg.Header.Get("Message-ID")
	originalSubject := originalMsg.Header.Get("Subject")

	bounceBody := fmt.Sprintf(`This is an automatically generated Delivery Status Notification.

Delivery to the following recipients failed:

    %s

Reason: %s

--- Original Message Headers ---
Message-ID: %s
Subject: %s
`, item.Recipients, deliveryErr.Error(), originalMessageID, originalSubject)

	// Build RFC5322 message
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: mailer-daemon@%s\r\n", w.config.Hostname))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", item.Sender))
	buf.WriteString(fmt.Sprintf("Subject: Delivery Status Notification (Failure)\r\n"))
	buf.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	buf.WriteString(fmt.Sprintf("Message-ID: <%d.%d.bounce@%s>\r\n", item.ID, time.Now().Unix(), w.config.Hostname))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(bounceBody)

	// Queue the bounce message
	recipients := []string{item.Sender}

	messageID, err := w.queueService.Enqueue(
		"mailer-daemon@"+w.config.Hostname,
		recipients,
		buf.Bytes(),
	)
	if err != nil {
		return fmt.Errorf("failed to queue bounce message: %w", err)
	}

	w.logger.Info("bounce message generated",
		zap.String("bounce_message_id", messageID),
		zap.Int64("original_id", item.ID),
	)

	return nil
}
