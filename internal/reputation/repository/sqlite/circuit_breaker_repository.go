package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type circuitBreakerRepository struct {
	db *sql.DB
}

// NewCircuitBreakerRepository creates a new SQLite circuit breaker repository
func NewCircuitBreakerRepository(db *sql.DB) repository.CircuitBreakerRepository {
	return &circuitBreakerRepository{db: db}
}

// RecordPause creates a new circuit breaker pause event
func (r *circuitBreakerRepository) RecordPause(ctx context.Context, event *domain.CircuitBreakerEvent) error {
	query := `
		INSERT INTO circuit_breaker_events (
			domain, trigger_type, trigger_value, threshold, paused_at, admin_notes
		) VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		event.Domain,
		event.TriggerType,
		event.TriggerValue,
		event.Threshold,
		event.PausedAt,
		event.AdminNotes,
	)
	if err != nil {
		return fmt.Errorf("failed to record pause event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get event ID: %w", err)
	}

	event.ID = id
	return nil
}

// RecordResume marks a circuit breaker as resumed
func (r *circuitBreakerRepository) RecordResume(ctx context.Context, domainName string, autoResumed bool, notes string) error {
	query := `
		UPDATE circuit_breaker_events
		SET resumed_at = ?, auto_resumed = ?, admin_notes = ?
		WHERE domain = ? AND resumed_at IS NULL
	`

	now := sql.NullInt64{Int64: nowUnix(), Valid: true}
	_, err := r.db.ExecContext(ctx, query, now, autoResumed, notes, domainName)
	if err != nil {
		return fmt.Errorf("failed to record resume: %w", err)
	}

	return nil
}

// GetActiveBreakers retrieves all currently active circuit breakers
func (r *circuitBreakerRepository) GetActiveBreakers(ctx context.Context) ([]*domain.CircuitBreakerEvent, error) {
	query := `
		SELECT
			id, domain, trigger_type, trigger_value, threshold,
			paused_at, resumed_at, auto_resumed, admin_notes
		FROM circuit_breaker_events
		WHERE resumed_at IS NULL
		ORDER BY paused_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active breakers: %w", err)
	}
	defer rows.Close()

	var events []*domain.CircuitBreakerEvent
	for rows.Next() {
		event := &domain.CircuitBreakerEvent{}
		var resumedAt sql.NullInt64

		err := rows.Scan(
			&event.ID,
			&event.Domain,
			&event.TriggerType,
			&event.TriggerValue,
			&event.Threshold,
			&event.PausedAt,
			&resumedAt,
			&event.AutoResumed,
			&event.AdminNotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan breaker event: %w", err)
		}

		if resumedAt.Valid {
			event.ResumedAt = &resumedAt.Int64
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating active breakers: %w", err)
	}

	return events, nil
}

// GetBreakerHistory retrieves circuit breaker history for a domain
func (r *circuitBreakerRepository) GetBreakerHistory(ctx context.Context, domainName string, limit int) ([]*domain.CircuitBreakerEvent, error) {
	query := `
		SELECT
			id, domain, trigger_type, trigger_value, threshold,
			paused_at, resumed_at, auto_resumed, admin_notes
		FROM circuit_breaker_events
		WHERE domain = ?
		ORDER BY paused_at DESC
		LIMIT ?
	`

	rows, err := r.db.QueryContext(ctx, query, domainName, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get breaker history: %w", err)
	}
	defer rows.Close()

	var events []*domain.CircuitBreakerEvent
	for rows.Next() {
		event := &domain.CircuitBreakerEvent{}
		var resumedAt sql.NullInt64

		err := rows.Scan(
			&event.ID,
			&event.Domain,
			&event.TriggerType,
			&event.TriggerValue,
			&event.Threshold,
			&event.PausedAt,
			&resumedAt,
			&event.AutoResumed,
			&event.AdminNotes,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan breaker event: %w", err)
		}

		if resumedAt.Valid {
			event.ResumedAt = &resumedAt.Int64
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating breaker history: %w", err)
	}

	return events, nil
}

// nowUnix returns the current Unix timestamp
func nowUnix() int64 {
	return time.Now().Unix()
}
