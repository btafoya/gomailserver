package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type dmarcReportsRepository struct {
	db *sql.DB
}

// NewDMARCReportsRepository creates a new SQLite DMARC reports repository
func NewDMARCReportsRepository(db *sql.DB) repository.DMARCReportsRepository {
	return &dmarcReportsRepository{db: db}
}

// Create stores a new DMARC report
func (r *dmarcReportsRepository) Create(ctx context.Context, report *domain.DMARCReport) error {
	query := `
		INSERT INTO dmarc_reports (
			domain, report_id, begin_time, end_time, organization,
			total_messages, spf_pass, dkim_pass, alignment_pass,
			raw_xml, processed_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		report.Domain,
		report.ReportID,
		report.BeginTime,
		report.EndTime,
		report.Organization,
		report.TotalMessages,
		report.SPFPass,
		report.DKIMPass,
		report.AlignmentPass,
		report.RawXML,
		report.ProcessedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create DMARC report: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	report.ID = id

	return nil
}

// GetByID retrieves a report by ID
func (r *dmarcReportsRepository) GetByID(ctx context.Context, id int64) (*domain.DMARCReport, error) {
	query := `
		SELECT
			id, domain, report_id, begin_time, end_time, organization,
			total_messages, spf_pass, dkim_pass, alignment_pass,
			raw_xml, processed_at
		FROM dmarc_reports
		WHERE id = ?
	`

	report := &domain.DMARCReport{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&report.ID,
		&report.Domain,
		&report.ReportID,
		&report.BeginTime,
		&report.EndTime,
		&report.Organization,
		&report.TotalMessages,
		&report.SPFPass,
		&report.DKIMPass,
		&report.AlignmentPass,
		&report.RawXML,
		&report.ProcessedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("DMARC report not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get DMARC report: %w", err)
	}

	return report, nil
}

// GetByReportID retrieves a report by its external report_id
func (r *dmarcReportsRepository) GetByReportID(ctx context.Context, reportID string) (*domain.DMARCReport, error) {
	query := `
		SELECT
			id, domain, report_id, begin_time, end_time, organization,
			total_messages, spf_pass, dkim_pass, alignment_pass,
			raw_xml, processed_at
		FROM dmarc_reports
		WHERE report_id = ?
	`

	report := &domain.DMARCReport{}
	err := r.db.QueryRowContext(ctx, query, reportID).Scan(
		&report.ID,
		&report.Domain,
		&report.ReportID,
		&report.BeginTime,
		&report.EndTime,
		&report.Organization,
		&report.TotalMessages,
		&report.SPFPass,
		&report.DKIMPass,
		&report.AlignmentPass,
		&report.RawXML,
		&report.ProcessedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("DMARC report not found with report_id %s", reportID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get DMARC report: %w", err)
	}

	return report, nil
}

// ListByDomain returns all reports for a domain
func (r *dmarcReportsRepository) ListByDomain(ctx context.Context, domainName string, limit, offset int) ([]*domain.DMARCReport, error) {
	query := `
		SELECT
			id, domain, report_id, begin_time, end_time, organization,
			total_messages, spf_pass, dkim_pass, alignment_pass,
			raw_xml, processed_at
		FROM dmarc_reports
		WHERE domain = ?
		ORDER BY begin_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list DMARC reports: %w", err)
	}
	defer rows.Close()

	var reports []*domain.DMARCReport
	for rows.Next() {
		report := &domain.DMARCReport{}
		err := rows.Scan(
			&report.ID,
			&report.Domain,
			&report.ReportID,
			&report.BeginTime,
			&report.EndTime,
			&report.Organization,
			&report.TotalMessages,
			&report.SPFPass,
			&report.DKIMPass,
			&report.AlignmentPass,
			&report.RawXML,
			&report.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DMARC report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DMARC reports: %w", err)
	}

	return reports, nil
}

// ListByTimeRange returns reports within a time range
func (r *dmarcReportsRepository) ListByTimeRange(ctx context.Context, startTime, endTime int64, limit, offset int) ([]*domain.DMARCReport, error) {
	query := `
		SELECT
			id, domain, report_id, begin_time, end_time, organization,
			total_messages, spf_pass, dkim_pass, alignment_pass,
			raw_xml, processed_at
		FROM dmarc_reports
		WHERE begin_time >= ? AND end_time <= ?
		ORDER BY begin_time DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list DMARC reports by time: %w", err)
	}
	defer rows.Close()

	var reports []*domain.DMARCReport
	for rows.Next() {
		report := &domain.DMARCReport{}
		err := rows.Scan(
			&report.ID,
			&report.Domain,
			&report.ReportID,
			&report.BeginTime,
			&report.EndTime,
			&report.Organization,
			&report.TotalMessages,
			&report.SPFPass,
			&report.DKIMPass,
			&report.AlignmentPass,
			&report.RawXML,
			&report.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DMARC report: %w", err)
		}
		reports = append(reports, report)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DMARC reports: %w", err)
	}

	return reports, nil
}

// GetDomainStats returns aggregated statistics for a domain
func (r *dmarcReportsRepository) GetDomainStats(ctx context.Context, domainName string, days int) (*domain.AlignmentAnalysis, error) {
	query := `
		SELECT
			COALESCE(SUM(total_messages), 0) as total,
			COALESCE(SUM(spf_pass), 0) as spf_pass,
			COALESCE(SUM(dkim_pass), 0) as dkim_pass,
			COALESCE(SUM(alignment_pass), 0) as alignment_pass
		FROM dmarc_reports
		WHERE domain = ? AND begin_time >= ?
	`

	startTime := int64(0)
	if days > 0 {
		startTime = time.Now().Unix() - int64(days*24*3600)
	}

	var total, spfPass, dkimPass, alignmentPass int64
	err := r.db.QueryRowContext(ctx, query, domainName, startTime).Scan(
		&total,
		&spfPass,
		&dkimPass,
		&alignmentPass,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get domain stats: %w", err)
	}

	stats := &domain.AlignmentAnalysis{
		Domain:          domainName,
		TotalMessages:   int(total),
		SPFAligned:      int(spfPass),
		DKIMAligned:     int(dkimPass),
		BothAligned:     int(alignmentPass),
		SPFMisaligned:   int(total - spfPass),
		DKIMMisaligned:  int(total - dkimPass),
		BothMisaligned:  int(total - alignmentPass),
	}

	if total > 0 {
		stats.SPFAlignmentRate = float64(spfPass) / float64(total)
		stats.DKIMAlignmentRate = float64(dkimPass) / float64(total)
		stats.OverallAlignmentRate = float64(alignmentPass) / float64(total)
	}

	return stats, nil
}

// GetRecentReports returns the most recent DMARC reports across all domains
func (r *dmarcReportsRepository) GetRecentReports(ctx context.Context, limit int) ([]*domain.DMARCReport, error) {
	query := `
		SELECT id, domain, report_id, begin_time, end_time, organization,
		       total_messages, spf_pass, dkim_pass, alignment_pass,
		       raw_xml, processed_at
		FROM dmarc_reports
		ORDER BY processed_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent reports: %w", err)
	}
	defer rows.Close()

	var reports []*domain.DMARCReport
	for rows.Next() {
		report := &domain.DMARCReport{}
		err := rows.Scan(
			&report.ID,
			&report.Domain,
			&report.ReportID,
			&report.BeginTime,
			&report.EndTime,
			&report.Organization,
			&report.TotalMessages,
			&report.SPFPass,
			&report.DKIMPass,
			&report.AlignmentPass,
			&report.RawXML,
			&report.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan report: %w", err)
		}
		reports = append(reports, report)
	}

	return reports, rows.Err()
}

// GetRecentActions returns recent DMARC auto-actions
func (r *dmarcReportsRepository) GetRecentActions(ctx context.Context, limit int) ([]*domain.DMARCAutoAction, error) {
	query := `
		SELECT id, domain, issue_type, description, action_taken,
		       taken_at, success, error_message
		FROM dmarc_auto_actions
		ORDER BY taken_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent actions: %w", err)
	}
	defer rows.Close()

	var actions []*domain.DMARCAutoAction
	for rows.Next() {
		action := &domain.DMARCAutoAction{}
		var issueType string
		var errorMsg sql.NullString

		err := rows.Scan(
			&action.ID,
			&action.Domain,
			&issueType,
			&action.Description,
			&action.ActionTaken,
			&action.TakenAt,
			&action.Success,
			&errorMsg,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan action: %w", err)
		}

		action.IssueType = domain.DMARCIssueType(issueType)
		if errorMsg.Valid {
			action.ErrorMessage = errorMsg.String
		}

		actions = append(actions, action)
	}

	return actions, rows.Err()
}

// CreateRecord stores a DMARC report record
func (r *dmarcReportsRepository) CreateRecord(ctx context.Context, record *domain.DMARCReportRecord) error {
	query := `
		INSERT INTO dmarc_report_records (
			report_id, source_ip, count, disposition,
			spf_result, dkim_result, spf_aligned, dkim_aligned,
			header_from, envelope_from
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		record.ReportID,
		record.SourceIP,
		record.Count,
		record.Disposition,
		record.SPFResult,
		record.DKIMResult,
		record.SPFAligned,
		record.DKIMAligned,
		record.HeaderFrom,
		record.EnvelopeFrom,
	)
	if err != nil {
		return fmt.Errorf("failed to create DMARC record: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	record.ID = id

	return nil
}

// GetRecordsByReportID retrieves all records for a report
func (r *dmarcReportsRepository) GetRecordsByReportID(ctx context.Context, reportID int64) ([]*domain.DMARCReportRecord, error) {
	query := `
		SELECT
			id, report_id, source_ip, count, disposition,
			spf_result, dkim_result, spf_aligned, dkim_aligned,
			header_from, envelope_from
		FROM dmarc_report_records
		WHERE report_id = ?
		ORDER BY count DESC
	`

	rows, err := r.db.QueryContext(ctx, query, reportID)
	if err != nil {
		return nil, fmt.Errorf("failed to get DMARC records: %w", err)
	}
	defer rows.Close()

	var records []*domain.DMARCReportRecord
	for rows.Next() {
		record := &domain.DMARCReportRecord{}
		err := rows.Scan(
			&record.ID,
			&record.ReportID,
			&record.SourceIP,
			&record.Count,
			&record.Disposition,
			&record.SPFResult,
			&record.DKIMResult,
			&record.SPFAligned,
			&record.DKIMAligned,
			&record.HeaderFrom,
			&record.EnvelopeFrom,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan DMARC record: %w", err)
		}
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating DMARC records: %w", err)
	}

	return records, nil
}
