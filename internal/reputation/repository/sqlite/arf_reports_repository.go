package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type arfReportsRepository struct {
	db *sql.DB
}

// NewARFReportsRepository creates a new SQLite ARF reports repository
func NewARFReportsRepository(db *sql.DB) repository.ARFReportsRepository {
	return &arfReportsRepository{db: db}
}

// Create stores a new ARF complaint report
func (r *arfReportsRepository) Create(ctx context.Context, report *domain.ARFReport) error {
	query := `
		INSERT INTO arf_reports (
			received_at, feedback_type, user_agent, version,
			original_rcpt_to, arrival_date, reporting_mta, source_ip,
			authentication_results, message_id, subject, raw_report,
			processed, suppressed_recipient
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		report.ReceivedAt,
		report.FeedbackType,
		report.UserAgent,
		report.Version,
		report.OriginalRcptTo,
		report.ArrivalDate,
		report.ReportingMTA,
		report.SourceIP,
		report.AuthenticationResults,
		report.MessageID,
		report.Subject,
		report.RawReport,
		report.Processed,
		report.SuppressedRecipient,
	)
	if err != nil {
		return fmt.Errorf("failed to create ARF report: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	report.ID = id

	return nil
}

// GetByID retrieves a report by ID
func (r *arfReportsRepository) GetByID(ctx context.Context, id int64) (*domain.ARFReport, error) {
	query := `
		SELECT
			id, received_at, feedback_type, user_agent, version,
			original_rcpt_to, arrival_date, reporting_mta, source_ip,
			authentication_results, message_id, subject, raw_report,
			processed, suppressed_recipient
		FROM arf_reports
		WHERE id = ?
	`

	report := &domain.ARFReport{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.ReceivedAt,
		&report.FeedbackType,
		&report.UserAgent,
		&report.Version,
		&report.OriginalRcptTo,
		&report.ArrivalDate,
		&report.ReportingMTA,
		&report.SourceIP,
		&report.AuthenticationResults,
		&report.MessageID,
		&report.Subject,
		&report.RawReport,
		&report.Processed,
		&report.SuppressedRecipient,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("ARF report not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get ARF report: %w", err)
	}

	return report, nil
}

// ListUnprocessed returns unprocessed complaints
func (r *arfReportsRepository) ListUnprocessed(ctx context.Context, limit int) ([]*domain.ARFReport, error) {
	query := `
		SELECT
			id, received_at, feedback_type, user_agent, version,
			original_rcpt_to, arrival_date, reporting_mta, source_ip,
			authentication_results, message_id, subject, raw_report,
			processed, suppressed_recipient
		FROM arf_reports
		WHERE processed = 0
		ORDER BY received_at ASC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list unprocessed ARF reports: %w", err)
	}
	defer rows.Close()

	var reports []*domain.ARFReport
	for rows.Next() {
		report := &domain.ARFReport{}
		err := rows.Scan(
			&report.ID,
			&report.ReceivedAt,
			&report.FeedbackType,
			&report.UserAgent,
			&report.Version,
			&report.OriginalRcptTo,
			&report.ArrivalDate,
			&report.ReportingMTA,
			&report.SourceIP,
			&report.AuthenticationResults,
			&report.MessageID,
			&report.Subject,
			&report.RawReport,
			&report.Processed,
			&report.SuppressedRecipient,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ARF report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ARF reports: %w", err)
	}

	return reports, nil
}

// MarkProcessed marks a report as processed
func (r *arfReportsRepository) MarkProcessed(ctx context.Context, id int64, suppressedRecipient string) error {
	query := `
		UPDATE arf_reports
		SET processed = 1, suppressed_recipient = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, suppressedRecipient, id)
	if err != nil {
		return fmt.Errorf("failed to mark ARF report as processed: %w", err)
	}

	return nil
}

// ListByTimeRange returns complaints within a time range
func (r *arfReportsRepository) ListByTimeRange(ctx context.Context, startTime, endTime int64, limit, offset int) ([]*domain.ARFReport, error) {
	query := `
		SELECT
			id, received_at, feedback_type, user_agent, version,
			original_rcpt_to, arrival_date, reporting_mta, source_ip,
			authentication_results, message_id, subject, raw_report,
			processed, suppressed_recipient
		FROM arf_reports
		WHERE received_at >= ? AND received_at <= ?
		ORDER BY received_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list ARF reports by time: %w", err)
	}
	defer rows.Close()

	var reports []*domain.ARFReport
	for rows.Next() {
		report := &domain.ARFReport{}
		err := rows.Scan(
			&report.ID,
			&report.ReceivedAt,
			&report.FeedbackType,
			&report.UserAgent,
			&report.Version,
			&report.OriginalRcptTo,
			&report.ArrivalDate,
			&report.ReportingMTA,
			&report.SourceIP,
			&report.AuthenticationResults,
			&report.MessageID,
			&report.Subject,
			&report.RawReport,
			&report.Processed,
			&report.SuppressedRecipient,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ARF report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating ARF reports: %w", err)
	}

	return reports, nil
}

// GetComplaintRate calculates complaint rate for a domain/IP
func (r *arfReportsRepository) GetComplaintRate(ctx context.Context, domainName string, hours int) (float64, error) {
	// Get complaint count
	complaintQuery := `
		SELECT COUNT(*)
		FROM arf_reports
		WHERE received_at >= ?
		AND (
			original_rcpt_to LIKE '%@' || ? || '%'
			OR reporting_mta LIKE '%' || ? || '%'
		)
	`

	startTime := time.Now().Add(-time.Duration(hours) * time.Hour).Unix()

	var complaintCount int
	err := r.db.QueryRowContext(ctx, complaintQuery, startTime, domainName, domainName).Scan(&complaintCount)
	if err != nil {
		return 0, fmt.Errorf("failed to get complaint count: %w", err)
	}

	// For now, return a simplified rate
	// In a complete implementation, this would need to query sent message counts
	// from the events repository or another source
	return float64(complaintCount), nil
}
