package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/domain"
	"go.uber.org/zap"
)

// AuditService handles audit logging for admin actions and security events
type AuditService struct {
	db     *database.DB
	logger *zap.Logger
}

// NewAuditService creates a new audit service
func NewAuditService(db *database.DB, logger *zap.Logger) *AuditService {
	return &AuditService{
		db:     db,
		logger: logger,
	}
}

// Log creates an audit log entry
func (s *AuditService) Log(ctx context.Context, log *domain.AuditLog) error {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}

	if log.Severity == "" {
		log.Severity = domain.SeverityInfo
	}

	query := `
		INSERT INTO audit_logs (
			timestamp, user_id, username, action, resource_type, resource_id,
			details, ip_address, user_agent, severity, success
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(ctx, query,
		log.Timestamp,
		log.UserID,
		log.Username,
		log.Action,
		log.ResourceType,
		log.ResourceID,
		log.Details,
		log.IPAddress,
		log.UserAgent,
		log.Severity,
		log.Success,
	)
	if err != nil {
		s.logger.Error("failed to create audit log",
			zap.Error(err),
			zap.String("action", log.Action),
		)
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	id, err := result.LastInsertId()
	if err == nil {
		log.ID = id
	}

	// Also log to structured logger for real-time monitoring
	fields := []zap.Field{
		zap.String("action", log.Action),
		zap.String("resource_type", log.ResourceType),
		zap.String("severity", log.Severity),
		zap.Bool("success", log.Success),
	}
	if log.UserID != nil {
		fields = append(fields, zap.Int64("user_id", *log.UserID))
	}
	if log.Username != "" {
		fields = append(fields, zap.String("username", log.Username))
	}
	if log.ResourceID != "" {
		fields = append(fields, zap.String("resource_id", log.ResourceID))
	}
	if log.IPAddress != "" {
		fields = append(fields, zap.String("ip_address", log.IPAddress))
	}

	s.logger.Info("audit_event", fields...)

	return nil
}

// LogAction is a convenience method for logging admin actions
func (s *AuditService) LogAction(ctx context.Context, userID *int64, username, action, resourceType, resourceID string, details interface{}) error {
	detailsJSON := ""
	if details != nil {
		b, err := json.Marshal(details)
		if err == nil {
			detailsJSON = string(b)
		}
	}

	return s.Log(ctx, &domain.AuditLog{
		UserID:       userID,
		Username:     username,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Details:      detailsJSON,
		Success:      true,
	})
}

// LogSecurityEvent logs a security event with appropriate severity
func (s *AuditService) LogSecurityEvent(ctx context.Context, action, resourceType string, success bool, severity string, details interface{}) error {
	detailsJSON := ""
	if details != nil {
		b, err := json.Marshal(details)
		if err == nil {
			detailsJSON = string(b)
		}
	}

	return s.Log(ctx, &domain.AuditLog{
		Action:       action,
		ResourceType: resourceType,
		Details:      detailsJSON,
		Severity:     severity,
		Success:      success,
	})
}

// GetLogs retrieves audit logs with filtering and pagination
func (s *AuditService) GetLogs(ctx context.Context, filter AuditLogFilter) ([]*domain.AuditLog, error) {
	query := `
		SELECT id, timestamp, user_id, username, action, resource_type, resource_id,
		       details, ip_address, user_agent, severity, success
		FROM audit_logs
		WHERE 1=1
	`
	args := []interface{}{}

	if filter.UserID != nil {
		query += " AND user_id = ?"
		args = append(args, *filter.UserID)
	}

	if filter.Action != "" {
		query += " AND action = ?"
		args = append(args, filter.Action)
	}

	if filter.ResourceType != "" {
		query += " AND resource_type = ?"
		args = append(args, filter.ResourceType)
	}

	if filter.Severity != "" {
		query += " AND severity = ?"
		args = append(args, filter.Severity)
	}

	if !filter.StartTime.IsZero() {
		query += " AND timestamp >= ?"
		args = append(args, filter.StartTime)
	}

	if !filter.EndTime.IsZero() {
		query += " AND timestamp <= ?"
		args = append(args, filter.EndTime)
	}

	query += " ORDER BY timestamp DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
		if filter.Offset > 0 {
			query += " OFFSET ?"
			args = append(args, filter.Offset)
		}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*domain.AuditLog
	for rows.Next() {
		log := &domain.AuditLog{}
		var userID sql.NullInt64
		var username, resourceID, details, ipAddress, userAgent sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.Timestamp,
			&userID,
			&username,
			&log.Action,
			&log.ResourceType,
			&resourceID,
			&details,
			&ipAddress,
			&userAgent,
			&log.Severity,
			&log.Success,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}

		if userID.Valid {
			log.UserID = &userID.Int64
		}
		log.Username = username.String
		log.ResourceID = resourceID.String
		log.Details = details.String
		log.IPAddress = ipAddress.String
		log.UserAgent = userAgent.String

		logs = append(logs, log)
	}

	return logs, rows.Err()
}

// AuditLogFilter defines filter criteria for retrieving audit logs
type AuditLogFilter struct {
	UserID       *int64
	Action       string
	ResourceType string
	Severity     string
	StartTime    time.Time
	EndTime      time.Time
	Limit        int
	Offset       int
}

// DeleteOldLogs removes audit logs older than the specified retention period
func (s *AuditService) DeleteOldLogs(ctx context.Context, retentionDays int) (int64, error) {
	cutoff := time.Now().AddDate(0, 0, -retentionDays)

	result, err := s.db.ExecContext(ctx,
		"DELETE FROM audit_logs WHERE timestamp < ?",
		cutoff,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	deleted, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	s.logger.Info("deleted old audit logs",
		zap.Int64("deleted", deleted),
		zap.Int("retention_days", retentionDays),
		zap.Time("cutoff", cutoff),
	)

	return deleted, nil
}
