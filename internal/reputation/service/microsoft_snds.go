package service

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
	"go.uber.org/zap"
)

// Microsoft SNDS API Response Structure

type sndsDataResponse struct {
	XMLName xml.Name  `xml:"SNDSDataResponse"`
	Records []sndsRecord `xml:"Record"`
}

type sndsRecord struct {
	Date            string  `xml:"Date"`
	IPAddress       string  `xml:"IPAddress"`
	MessageCount    int     `xml:"MessageCount"`
	FilterLevel     string  `xml:"FilterLevel"`     // Green, Yellow, Red
	ComplaintRate   float64 `xml:"ComplaintRate"`
	TrapHits        int     `xml:"TrapHits"`
}

// MicrosoftSNDSService integrates with Microsoft SNDS (Smart Network Data Services)

type MicrosoftSNDSService struct {
	client      *http.Client
	apiKey      string
	metricsRepo repository.SNDSMetricsRepository
	alertsRepo  repository.AlertsRepository
	logger      *zap.Logger
}

func NewMicrosoftSNDSService(
	apiKey string,
	metricsRepo repository.SNDSMetricsRepository,
	alertsRepo repository.AlertsRepository,
	logger *zap.Logger,
) *MicrosoftSNDSService {
	return &MicrosoftSNDSService{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:      apiKey,
		metricsRepo: metricsRepo,
		alertsRepo:  alertsRepo,
		logger:      logger,
	}
}

// FetchIPData retrieves reputation metrics for an IP address
func (s *MicrosoftSNDSService) FetchIPData(ctx context.Context, ipAddress string) (*domain.SNDSMetrics, error) {
	// Microsoft SNDS API endpoint
	url := fmt.Sprintf("https://postmaster.live.com/snds/data.aspx?key=%s&ip=%s", s.apiKey, ipAddress)

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Error("Failed to fetch SNDS data",
			zap.String("ip", ipAddress),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		s.logger.Error("SNDS API returned error",
			zap.String("ip", ipAddress),
			zap.Int("status", resp.StatusCode),
			zap.String("body", string(body)),
		)
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse XML response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var sndsResp sndsDataResponse
	if err := xml.Unmarshal(body, &sndsResp); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	if len(sndsResp.Records) == 0 {
		s.logger.Warn("No SNDS data available for IP", zap.String("ip", ipAddress))
		return nil, nil
	}

	// Get most recent record
	record := sndsResp.Records[0]

	// Parse metric date
	metricDate, err := time.Parse("2006-01-02", record.Date)
	if err != nil {
		s.logger.Warn("Failed to parse metric date", zap.Error(err))
		metricDate = time.Now()
	}

	// Convert to domain model
	metrics := &domain.SNDSMetrics{
		IPAddress:     ipAddress,
		FetchedAt:     time.Now().Unix(),
		MetricDate:    metricDate.Unix(),
		SpamTrapHits:  record.TrapHits,
		ComplaintRate: record.ComplaintRate,
		FilterLevel:   strings.ToUpper(record.FilterLevel), // GREEN, YELLOW, RED
		MessageCount:  record.MessageCount,
		RawResponse:   string(body),
	}

	s.logger.Info("Fetched Microsoft SNDS metrics",
		zap.String("ip", ipAddress),
		zap.String("filter_level", metrics.FilterLevel),
		zap.Float64("complaint_rate", metrics.ComplaintRate),
		zap.Int("trap_hits", metrics.SpamTrapHits),
	)

	return metrics, nil
}

// SyncIP syncs metrics for a single IP address
func (s *MicrosoftSNDSService) SyncIP(ctx context.Context, ipAddress string) error {
	metrics, err := s.FetchIPData(ctx, ipAddress)
	if err != nil {
		return err
	}

	if metrics == nil {
		return nil // No data available
	}

	// Store metrics
	if err := s.metricsRepo.Create(ctx, metrics); err != nil {
		s.logger.Error("Failed to store SNDS metrics", zap.Error(err))
		return fmt.Errorf("failed to store metrics: %w", err)
	}

	// Check for issues and create alerts
	if err := s.checkSNDSAlerts(ctx, metrics); err != nil {
		s.logger.Error("Failed to check SNDS alerts", zap.Error(err))
		// Don't fail sync on alert error
	}

	return nil
}

// SyncAll syncs metrics for all configured IPs
func (s *MicrosoftSNDSService) SyncAll(ctx context.Context, ipAddresses []string) error {
	s.logger.Info("Starting Microsoft SNDS sync for all IPs", zap.Int("count", len(ipAddresses)))

	errors := make([]error, 0)
	for _, ipAddress := range ipAddresses {
		if err := s.SyncIP(ctx, ipAddress); err != nil {
			s.logger.Error("Failed to sync IP",
				zap.String("ip", ipAddress),
				zap.Error(err),
			)
			errors = append(errors, err)
		}

		// Add small delay to avoid rate limiting
		time.Sleep(1 * time.Second)
	}

	if len(errors) > 0 {
		return fmt.Errorf("sync completed with %d errors", len(errors))
	}

	s.logger.Info("Microsoft SNDS sync completed successfully")
	return nil
}

// checkSNDSAlerts checks for filtering issues and creates alerts
func (s *MicrosoftSNDSService) checkSNDSAlerts(ctx context.Context, metrics *domain.SNDSMetrics) error {
	// Check filter level
	if metrics.FilterLevel == "RED" || metrics.FilterLevel == "YELLOW" {
		severity := domain.SeverityHigh
		if metrics.FilterLevel == "RED" {
			severity = domain.SeverityCritical
		}

		alert := &domain.ReputationAlert{
			Domain:    metrics.IPAddress, // Use IP as domain for SNDS
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  severity,
			Title:     fmt.Sprintf("Microsoft SNDS Filter Level: %s", metrics.FilterLevel),
			Message:   fmt.Sprintf("IP %s is being filtered at %s level by Microsoft", metrics.IPAddress, metrics.FilterLevel),
			Details: map[string]interface{}{
				"ip_address":     metrics.IPAddress,
				"filter_level":   metrics.FilterLevel,
				"complaint_rate": metrics.ComplaintRate,
				"trap_hits":      metrics.SpamTrapHits,
				"message_count":  metrics.MessageCount,
				"metric_date":    metrics.MetricDate,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	// Check spam trap hits
	if metrics.SpamTrapHits > 0 {
		severity := domain.SeverityMedium
		if metrics.SpamTrapHits > 10 {
			severity = domain.SeverityHigh
		}
		if metrics.SpamTrapHits > 50 {
			severity = domain.SeverityCritical
		}

		alert := &domain.ReputationAlert{
			Domain:    metrics.IPAddress,
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  severity,
			Title:     fmt.Sprintf("Spam Trap Hits Detected: %d", metrics.SpamTrapHits),
			Message:   fmt.Sprintf("IP %s hit %d spam traps at Microsoft", metrics.IPAddress, metrics.SpamTrapHits),
			Details: map[string]interface{}{
				"ip_address":  metrics.IPAddress,
				"trap_hits":   metrics.SpamTrapHits,
				"metric_date": metrics.MetricDate,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	// Check complaint rate
	if metrics.ComplaintRate > 0.001 { // 0.1%
		severity := domain.SeverityMedium
		if metrics.ComplaintRate > 0.003 {
			severity = domain.SeverityHigh
		}
		if metrics.ComplaintRate > 0.01 {
			severity = domain.SeverityCritical
		}

		alert := &domain.ReputationAlert{
			Domain:    metrics.IPAddress,
			AlertType: domain.AlertTypeExternalFeedback,
			Severity:  severity,
			Title:     fmt.Sprintf("High Complaint Rate: %.3f%%", metrics.ComplaintRate*100),
			Message:   fmt.Sprintf("IP %s has elevated complaint rate at Microsoft", metrics.IPAddress),
			Details: map[string]interface{}{
				"ip_address":     metrics.IPAddress,
				"complaint_rate": metrics.ComplaintRate,
				"metric_date":    metrics.MetricDate,
				"threshold":      0.001,
			},
			CreatedAt: time.Now().Unix(),
		}

		if err := s.alertsRepo.Create(ctx, alert); err != nil {
			return err
		}
	}

	return nil
}

// GetLatestMetrics returns the latest metrics for an IP
func (s *MicrosoftSNDSService) GetLatestMetrics(ctx context.Context, ipAddress string) (*domain.SNDSMetrics, error) {
	return s.metricsRepo.GetLatest(ctx, ipAddress)
}

// GetMetricsHistory returns metrics history for an IP
func (s *MicrosoftSNDSService) GetMetricsHistory(ctx context.Context, ipAddress string, days int) ([]*domain.SNDSMetrics, error) {
	return s.metricsRepo.ListByIP(ctx, ipAddress, days)
}

// GetFilterLevelTrend returns filter level trend over time
func (s *MicrosoftSNDSService) GetFilterLevelTrend(ctx context.Context, ipAddress string, days int) ([]string, error) {
	return s.metricsRepo.GetFilterLevelTrend(ctx, ipAddress, days)
}
