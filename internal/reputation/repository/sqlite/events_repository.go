package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type eventsRepository struct {
	db *sql.DB
}

// NewEventsRepository creates a new SQLite events repository
func NewEventsRepository(db *sql.DB) repository.EventsRepository {
	return &eventsRepository{db: db}
}

// RecordEvent stores a new sending event
func (r *eventsRepository) RecordEvent(ctx context.Context, event *domain.SendingEvent) error {
	var metadata *string
	if event.Metadata != nil {
		jsonData, err := json.Marshal(event.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		metadataStr := string(jsonData)
		metadata = &metadataStr
	}

	query := `
		INSERT INTO sending_events (
			timestamp, domain, recipient_domain, event_type,
			bounce_type, enhanced_status_code, smtp_response, ip_address, metadata
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		event.Timestamp,
		event.Domain,
		event.RecipientDomain,
		event.EventType,
		event.BounceType,
		event.EnhancedStatusCode,
		event.SMTPResponse,
		event.IPAddress,
		metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to record event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get event ID: %w", err)
	}

	event.ID = id
	return nil
}

// GetEventsInWindow retrieves events for a domain within a time window
func (r *eventsRepository) GetEventsInWindow(ctx context.Context, domainName string, startTime, endTime int64) ([]*domain.SendingEvent, error) {
	query := `
		SELECT
			id, timestamp, domain, recipient_domain, event_type,
			bounce_type, enhanced_status_code, smtp_response, ip_address, metadata
		FROM sending_events
		WHERE domain = ? AND timestamp >= ? AND timestamp <= ?
		ORDER BY timestamp DESC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []*domain.SendingEvent
	for rows.Next() {
		event := &domain.SendingEvent{}
		var metadata *string

		err := rows.Scan(
			&event.ID,
			&event.Timestamp,
			&event.Domain,
			&event.RecipientDomain,
			&event.EventType,
			&event.BounceType,
			&event.EnhancedStatusCode,
			&event.SMTPResponse,
			&event.IPAddress,
			&metadata,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}

		if metadata != nil {
			if err := json.Unmarshal([]byte(*metadata), &event.Metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}

// GetEventCountsByType returns counts of each event type for a domain in a time window
func (r *eventsRepository) GetEventCountsByType(ctx context.Context, domainName string, startTime, endTime int64) (map[string]int64, error) {
	query := `
		SELECT event_type, COUNT(*) as count
		FROM sending_events
		WHERE domain = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY event_type
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query event counts: %w", err)
	}
	defer rows.Close()

	counts := make(map[string]int64)
	for rows.Next() {
		var eventType string
		var count int64

		if err := rows.Scan(&eventType, &count); err != nil {
			return nil, fmt.Errorf("failed to scan event count: %w", err)
		}

		counts[eventType] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating event counts: %w", err)
	}

	return counts, nil
}

// CleanupOldEvents removes events older than the specified timestamp
func (r *eventsRepository) CleanupOldEvents(ctx context.Context, olderThan int64) error {
	query := `DELETE FROM sending_events WHERE timestamp < ?`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return fmt.Errorf("failed to cleanup old events: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected > 0 {
		// Log cleanup for debugging (could be enhanced with proper logging)
		_ = rowsAffected
	}

	return nil
}
