package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"go.uber.org/zap"
)

// MTASTSService handles MTA-STS policy fetching, caching, and enforcement
type MTASTSService struct {
	db         *database.DB
	logger     *zap.Logger
	httpClient *http.Client
}

// NewMTASTSService creates a new MTA-STS service
func NewMTASTSService(db *database.DB, logger *zap.Logger) *MTASTSService {
	return &MTASTSService{
		db:     db,
		logger: logger,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchPolicy fetches and caches an MTA-STS policy for a domain
func (s *MTASTSService) FetchPolicy(ctx context.Context, domainName string) (*domain.MTASTSPolicy, error) {
	// Check cache first
	cached, err := s.getCachedPolicy(ctx, domainName)
	if err == nil && cached != nil {
		// Verify cache hasn't expired
		if time.Now().Before(cached.ExpiresAt) {
			s.logger.Debug("MTA-STS policy cache hit",
				zap.String("domain", domainName),
			)
			return cached, nil
		}
	}

	// Fetch policy from well-known URL
	policy, err := s.fetchPolicyFromWeb(ctx, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch MTA-STS policy: %w", err)
	}

	// Cache the policy
	if err := s.cachePolicy(ctx, policy); err != nil {
		s.logger.Error("failed to cache MTA-STS policy",
			zap.Error(err),
			zap.String("domain", domainName),
		)
	}

	return policy, nil
}

// fetchPolicyFromWeb fetches the MTA-STS policy from the well-known URL
func (s *MTASTSService) fetchPolicyFromWeb(ctx context.Context, domainName string) (*domain.MTASTSPolicy, error) {
	// MTA-STS policy URL: https://mta-sts.{domain}/.well-known/mta-sts.txt
	url := fmt.Sprintf("https://mta-sts.%s/.well-known/mta-sts.txt", domainName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	policyText := string(body)
	policy, err := s.parsePolicy(domainName, policyText)
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy: %w", err)
	}

	s.logger.Info("MTA-STS policy fetched",
		zap.String("domain", domainName),
		zap.String("mode", policy.Mode),
		zap.Int("max_age", policy.MaxAge),
	)

	return policy, nil
}

// parsePolicy parses an MTA-STS policy text file
func (s *MTASTSService) parsePolicy(domainName, policyText string) (*domain.MTASTSPolicy, error) {
	policy := &domain.MTASTSPolicy{
		Domain:     domainName,
		FetchedAt:  time.Now(),
		PolicyText: policyText,
	}

	var mxPatterns []string

	lines := strings.Split(policyText, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "version":
			policy.Version = value
		case "mode":
			policy.Mode = value
		case "max_age":
			var maxAge int
			fmt.Sscanf(value, "%d", &maxAge)
			policy.MaxAge = maxAge
		case "mx":
			mxPatterns = append(mxPatterns, value)
		}
	}

	// Validate required fields
	if policy.Version == "" || policy.Mode == "" || policy.MaxAge == 0 {
		return nil, fmt.Errorf("invalid policy: missing required fields")
	}

	// Store MX patterns as JSON
	mxJSON, err := json.Marshal(mxPatterns)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal MX patterns: %w", err)
	}
	policy.MXPatterns = string(mxJSON)

	// Calculate expiration
	policy.ExpiresAt = policy.FetchedAt.Add(time.Duration(policy.MaxAge) * time.Second)

	return policy, nil
}

// getCachedPolicy retrieves a cached MTA-STS policy
func (s *MTASTSService) getCachedPolicy(ctx context.Context, domainName string) (*domain.MTASTSPolicy, error) {
	policy := &domain.MTASTSPolicy{}

	query := `
		SELECT id, domain, version, mode, max_age, mx_patterns,
		       fetched_at, expires_at, policy_text
		FROM mtasts_policy_cache
		WHERE domain = ?
	`

	err := s.db.QueryRowContext(ctx, query, domainName).Scan(
		&policy.ID,
		&policy.Domain,
		&policy.Version,
		&policy.Mode,
		&policy.MaxAge,
		&policy.MXPatterns,
		&policy.FetchedAt,
		&policy.ExpiresAt,
		&policy.PolicyText,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return policy, nil
}

// cachePolicy stores an MTA-STS policy in the cache
func (s *MTASTSService) cachePolicy(ctx context.Context, policy *domain.MTASTSPolicy) error {
	query := `
		INSERT OR REPLACE INTO mtasts_policy_cache (
			domainName, version, mode, max_age, mx_patterns,
			fetched_at, expires_at, policy_text
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := s.db.ExecContext(ctx, query,
		policy.Domain,
		policy.Version,
		policy.Mode,
		policy.MaxAge,
		policy.MXPatterns,
		policy.FetchedAt,
		policy.ExpiresAt,
		policy.PolicyText,
	)

	return err
}

// EnforcePolicy checks if an MX hostname is allowed by the MTA-STS policy
func (s *MTASTSService) EnforcePolicy(ctx context.Context, domainName, mxHostname string) (bool, error) {
	policy, err := s.FetchPolicy(ctx, domainName)
	if err != nil {
		return false, fmt.Errorf("failed to fetch policy: %w", err)
	}

	// If mode is "none", don't enforce
	if policy.Mode == domain.MTASTSModeNone {
		return true, nil
	}

	// Parse MX patterns
	var mxPatterns []string
	if err := json.Unmarshal([]byte(policy.MXPatterns), &mxPatterns); err != nil {
		return false, fmt.Errorf("failed to parse MX patterns: %w", err)
	}

	// Check if MX hostname matches any pattern
	for _, pattern := range mxPatterns {
		if s.matchMXPattern(mxHostname, pattern) {
			return true, nil
		}
	}

	// In testing mode, log but allow
	if policy.Mode == domain.MTASTSModeTesting {
		s.logger.Warn("MTA-STS policy violation (testing mode)",
			zap.String("domain", domainName),
			zap.String("mx_hostname", mxHostname),
		)
		return true, nil
	}

	// In enforce mode, reject
	s.logger.Error("MTA-STS policy violation",
		zap.String("domain", domainName),
		zap.String("mx_hostname", mxHostname),
		zap.String("mode", policy.Mode),
	)

	return false, nil
}

// matchMXPattern checks if an MX hostname matches a pattern
// Patterns can use wildcards: *.example.com matches mail.example.com
func (s *MTASTSService) matchMXPattern(hostname, pattern string) bool {
	// Convert wildcard pattern to regex
	pattern = strings.ReplaceAll(pattern, ".", "\\.")
	pattern = strings.ReplaceAll(pattern, "*", ".*")
	pattern = "^" + pattern + "$"

	matched, err := regexp.MatchString(pattern, hostname)
	if err != nil {
		s.logger.Error("invalid MX pattern",
			zap.Error(err),
			zap.String("pattern", pattern),
		)
		return false
	}

	return matched
}

// ClearCache removes expired MTA-STS policies from the cache
func (s *MTASTSService) ClearCache(ctx context.Context) (int64, error) {
	query := "DELETE FROM mtasts_policy_cache WHERE expires_at < ?"

	result, err := s.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return 0, fmt.Errorf("failed to clear MTA-STS cache: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	if deleted > 0 {
		s.logger.Info("cleared expired MTA-STS policies",
			zap.Int64("deleted", deleted),
		)
	}

	return deleted, nil
}

// CreateTLSReport creates a TLS report entry for TLSRPT (RFC 8460)
func (s *MTASTSService) CreateTLSReport(ctx context.Context, report *domain.TLSReport) error {
	if report.CreatedAt.IsZero() {
		report.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO tls_reports (
			report_id, domainName, date_range_start, date_range_end,
			contact_info, report_json, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		report.ReportID,
		report.Domain,
		report.DateRangeStart,
		report.DateRangeEnd,
		report.ContactInfo,
		report.ReportJSON,
		report.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create TLS report: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		report.ID = id
	}

	s.logger.Info("TLS report created",
		zap.String("domain", report.Domain),
		zap.String("report_id", report.ReportID),
	)

	return nil
}

// GetPendingReports retrieves unsent TLS reports
func (s *MTASTSService) GetPendingReports(ctx context.Context) ([]*domain.TLSReport, error) {
	query := `
		SELECT id, report_id, domainName, date_range_start, date_range_end,
		       contact_info, report_json, created_at
		FROM tls_reports
		WHERE sent_at IS NULL
		ORDER BY created_at ASC
	`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query TLS reports: %w", err)
	}
	defer rows.Close()

	var reports []*domain.TLSReport
	for rows.Next() {
		report := &domain.TLSReport{}
		var contactInfo sql.NullString

		err := rows.Scan(
			&report.ID,
			&report.ReportID,
			&report.Domain,
			&report.DateRangeStart,
			&report.DateRangeEnd,
			&contactInfo,
			&report.ReportJSON,
			&report.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan TLS report: %w", err)
		}

		report.ContactInfo = contactInfo.String
		reports = append(reports, report)
	}

	return reports, rows.Err()
}

// MarkReportSent marks a TLS report as sent
func (s *MTASTSService) MarkReportSent(ctx context.Context, reportID int64) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx,
		"UPDATE tls_reports SET sent_at = ? WHERE id = ?",
		now, reportID,
	)
	if err != nil {
		return fmt.Errorf("failed to mark report as sent: %w", err)
	}

	s.logger.Info("TLS report marked as sent", zap.Int64("report_id", reportID))

	return nil
}
