package service

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/mail"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// ARFParserService parses Abuse Reporting Format (ARF) complaint reports

type ARFParserService struct {
	arfRepo     repository.ARFReportsRepository
	eventsRepo  repository.EventsRepository
	alertsRepo  repository.AlertsRepository
	logger      *zap.Logger
}

func NewARFParserService(
	arfRepo repository.ARFReportsRepository,
	eventsRepo repository.EventsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *ARFParserService {
	return &ARFParserService{
		arfRepo:    arfRepo,
		eventsRepo: eventsRepo,
		alertsRepo: alertsRepo,
		logger:     logger,
	}
}

// ParseARFReport parses an ARF complaint report
func (s *ARFParserService) ParseARFReport(ctx context.Context, rawMessage []byte) (*domain.ARFReport, error) {
	msg, err := mail.ReadMessage(strings.NewReader(string(rawMessage)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email: %w", err)
	}

	report := &domain.ARFReport{
		ReceivedAt: time.Now().Unix(),
		RawReport:  string(rawMessage),
		Processed:  false,
	}

	// Parse ARF headers
	report.FeedbackType = msg.Header.Get("Feedback-Type")
	report.UserAgent = msg.Header.Get("User-Agent")
	report.Version = msg.Header.Get("Version")
	report.OriginalRcptTo = msg.Header.Get("Original-Rcpt-To")

	// Parse arrival date
	if arrivalDate := msg.Header.Get("Arrival-Date"); arrivalDate != "" {
		if t, err := mail.ParseDate(arrivalDate); err == nil {
			report.ArrivalDate = t.Unix()
		}
	}

	report.ReportingMTA = msg.Header.Get("Reporting-MTA")
	report.SourceIP = msg.Header.Get("Source-IP")
	report.AuthenticationResults = msg.Header.Get("Authentication-Results")

	// Try to extract original message details from body
	body, _ := io.ReadAll(msg.Body)
	report.MessageID, report.Subject = s.extractOriginalMessageDetails(string(body))

	s.logger.Info("Parsed ARF report",
		zap.String("feedback_type", report.FeedbackType),
		zap.String("recipient", report.OriginalRcptTo),
		zap.String("source_ip", report.SourceIP),
	)

	return report, nil
}

// extractOriginalMessageDetails extracts details from the original message
func (s *ARFParserService) extractOriginalMessageDetails(body string) (messageID, subject string) {
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Message-ID:") {
			messageID = strings.TrimSpace(strings.TrimPrefix(line, "Message-ID:"))
		} else if strings.HasPrefix(line, "Subject:") {
			subject = strings.TrimSpace(strings.TrimPrefix(line, "Subject:"))
		}
	}
	return
}

// StoreReport saves an ARF report to the database
func (s *ARFParserService) StoreReport(ctx context.Context, report *domain.ARFReport) error {
	return s.arfRepo.Create(ctx, report)
}

// ProcessComplaint processes an ARF complaint and suppresses the recipient
func (s *ARFParserService) ProcessComplaint(ctx context.Context, report *domain.ARFReport) error {
	if report.OriginalRcptTo == "" {
		s.logger.Warn("ARF report missing recipient, cannot suppress")
		return nil
	}

	// Record complaint event in telemetry
	senderDomain := s.extractDomain(report.OriginalRcptTo)
	if senderDomain != "" {
		event := &domain.SendingEvent{
			Timestamp:       report.ReceivedAt,
			Domain:          senderDomain,
			RecipientDomain: s.extractDomain(report.OriginalRcptTo),
			EventType:       domain.EventComplaint,
			IPAddress:       report.SourceIP,
			Metadata: map[string]interface{}{
				"feedback_type": report.FeedbackType,
				"message_id":    report.MessageID,
			},
		}

		if err := s.eventsRepo.RecordEvent(ctx, event); err != nil {
			s.logger.Error("Failed to record complaint event", zap.Error(err))
		}
	}

	// Mark as processed with suppressed recipient
	if err := s.arfRepo.MarkProcessed(ctx, report.ID, report.OriginalRcptTo); err != nil {
		return fmt.Errorf("failed to mark ARF report as processed: %w", err)
	}

	s.logger.Info("Processed ARF complaint and suppressed recipient",
		zap.String("recipient", report.OriginalRcptTo),
		zap.String("feedback_type", report.FeedbackType),
	)

	return nil
}

// ProcessUnprocessedReports processes all unprocessed ARF reports
func (s *ARFParserService) ProcessUnprocessedReports(ctx context.Context) error {
	reports, err := s.arfRepo.ListUnprocessed(ctx, 100)
	if err != nil {
		return fmt.Errorf("failed to list unprocessed reports: %w", err)
	}

	for _, report := range reports {
		if err := s.ProcessComplaint(ctx, report); err != nil {
			s.logger.Error("Failed to process ARF report",
				zap.Int64("id", report.ID),
				zap.Error(err),
			)
			// Continue processing other reports
		}
	}

	s.logger.Info("Processed unprocessed ARF reports", zap.Int("count", len(reports)))
	return nil
}

// extractDomain extracts domain from an email address
func (s *ARFParserService) extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// GetComplaintRate returns complaint rate for a domain
func (s *ARFParserService) GetComplaintRate(ctx context.Context, domain string, hours int) (float64, error) {
	return s.arfRepo.GetComplaintRate(ctx, domain, hours)
}
