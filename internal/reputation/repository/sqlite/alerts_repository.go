package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type alertsRepository struct {
	db *sql.DB
}

// NewAlertsRepository creates a new SQLite alerts repository
func NewAlertsRepository(db *sql.DB) repository.AlertsRepository {
	return &alertsRepository{db: db}
}

// Create stores a new alert
func (r *alertsRepository) Create(ctx context.Context, alert *domain.ReputationAlert) error {
	query := `
		INSERT INTO reputation_alerts (
			domain, alert_type, severity, title, message, details,
			created_at, acknowledged, acknowledged_at, acknowledged_by,
			resolved, resolved_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		alert.Domain,
		alert.AlertType,
		alert.Severity,
		alert.Title,
		alert.Message,
		alert.DetailsJSON(),
		alert.CreatedAt,
		alert.Acknowledged,
		alert.AcknowledgedAt,
		alert.AcknowledgedBy,
		alert.Resolved,
		alert.ResolvedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create alert: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	alert.ID = id

	return nil
}

// GetByID retrieves an alert by ID
func (r *alertsRepository) GetByID(ctx context.Context, id int64) (*domain.ReputationAlert, error) {
	query := `
		SELECT
			id, domain, alert_type, severity, title, message, details,
			created_at, acknowledged, acknowledged_at, acknowledged_by,
			resolved, resolved_at
		FROM reputation_alerts
		WHERE id = ?
	`

	alert := &domain.ReputationAlert{}
	var detailsJSON string
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&alert.ID,
		&alert.Domain,
		&alert.AlertType,
		&alert.Severity,
		&alert.Title,
		&alert.Message,
		&detailsJSON,
		&alert.CreatedAt,
		&alert.Acknowledged,
		&alert.AcknowledgedAt,
		&alert.AcknowledgedBy,
		&alert.Resolved,
		&alert.ResolvedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("alert not found with id %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get alert: %w", err)
	}

	// Parse details JSON
	if err := alert.ParseDetailsJSON(detailsJSON); err != nil {
		return nil, fmt.Errorf("failed to parse details JSON: %w", err)
	}

	return alert, nil
}

// ListUnacknowledged returns unacknowledged alerts
func (r *alertsRepository) ListUnacknowledged(ctx context.Context, limit int) ([]*domain.ReputationAlert, error) {
	query := `
		SELECT
			id, domain, alert_type, severity, title, message, details,
			created_at, acknowledged, acknowledged_at, acknowledged_by,
			resolved, resolved_at
		FROM reputation_alerts
		WHERE acknowledged = 0
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list unacknowledged alerts: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

// ListByDomain returns alerts for a domain
func (r *alertsRepository) ListByDomain(ctx context.Context, domainName string, limit, offset int) ([]*domain.ReputationAlert, error) {
	query := `
		SELECT
			id, domain, alert_type, severity, title, message, details,
			created_at, acknowledged, acknowledged_at, acknowledged_by,
			resolved, resolved_at
		FROM reputation_alerts
		WHERE domain = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by domain: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

// ListBySeverity returns alerts by severity level
func (r *alertsRepository) ListBySeverity(ctx context.Context, severity domain.AlertSeverity, limit int) ([]*domain.ReputationAlert, error) {
	query := `
		SELECT
			id, domain, alert_type, severity, title, message, details,
			created_at, acknowledged, acknowledged_at, acknowledged_by,
			resolved, resolved_at
		FROM reputation_alerts
		WHERE severity = ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, string(severity), limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list alerts by severity: %w", err)
	}
	defer rows.Close()

	return r.scanAlerts(rows)
}

// Acknowledge marks an alert as acknowledged
func (r *alertsRepository) Acknowledge(ctx context.Context, id int64, acknowledgedBy string) error {
	query := `
		UPDATE reputation_alerts
		SET acknowledged = 1,
		    acknowledged_at = ?,
		    acknowledged_by = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, domain.Now(), acknowledgedBy, id)
	if err != nil {
		return fmt.Errorf("failed to acknowledge alert: %w", err)
	}

	return nil
}

// Resolve marks an alert as resolved
func (r *alertsRepository) Resolve(ctx context.Context, id int64) error {
	query := `
		UPDATE reputation_alerts
		SET resolved = 1,
		    resolved_at = ?
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, domain.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to resolve alert: %w", err)
	}

	return nil
}

// GetUnacknowledgedCount returns count of unacknowledged alerts
func (r *alertsRepository) GetUnacknowledgedCount(ctx context.Context) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM reputation_alerts
		WHERE acknowledged = 0
	`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unacknowledged count: %w", err)
	}

	return count, nil
}

// GetUnacknowledgedCountByDomain returns count of unacknowledged alerts for a domain
func (r *alertsRepository) GetUnacknowledgedCountByDomain(ctx context.Context, domainName string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM reputation_alerts
		WHERE acknowledged = 0 AND domain = ?
	`

	var count int
	err := r.db.QueryRowContext(ctx, query, domainName).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unacknowledged count by domain: %w", err)
	}

	return count, nil
}

// scanAlerts is a helper function to scan alert rows
func (r *alertsRepository) scanAlerts(rows *sql.Rows) ([]*domain.ReputationAlert, error) {
	var alerts []*domain.ReputationAlert
	for rows.Next() {
		alert := &domain.ReputationAlert{}
		var detailsJSON string
		err := rows.Scan(
			&alert.ID,
			&alert.Domain,
			&alert.AlertType,
			&alert.Severity,
			&alert.Title,
			&alert.Message,
			&detailsJSON,
			&alert.CreatedAt,
			&alert.Acknowledged,
			&alert.AcknowledgedAt,
			&alert.AcknowledgedBy,
			&alert.Resolved,
			&alert.ResolvedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan alert: %w", err)
		}

		// Parse details JSON
		if err := alert.ParseDetailsJSON(detailsJSON); err != nil {
			return nil, fmt.Errorf("failed to parse details JSON: %w", err)
		}

		alerts = append(alerts, alert)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating alerts: %w", err)
	}

	return alerts, nil
}
