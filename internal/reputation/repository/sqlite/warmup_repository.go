package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type warmUpRepository struct {
	db *sql.DB
}

// NewWarmUpRepository creates a new SQLite warm-up repository
func NewWarmUpRepository(db *sql.DB) repository.WarmUpRepository {
	return &warmUpRepository{db: db}
}

// GetSchedule retrieves the warm-up schedule for a domain
func (r *warmUpRepository) GetSchedule(ctx context.Context, domainName string) ([]*domain.WarmUpDay, error) {
	query := `
		SELECT domain, day, max_volume, actual_volume, created_at
		FROM warm_up_schedules
		WHERE domain = ?
		ORDER BY day ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get warm-up schedule: %w", err)
	}
	defer rows.Close()

	var schedule []*domain.WarmUpDay
	for rows.Next() {
		day := &domain.WarmUpDay{}
		err := rows.Scan(
			&day.Domain,
			&day.Day,
			&day.MaxVolume,
			&day.ActualVolume,
			&day.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan warm-up day: %w", err)
		}
		schedule = append(schedule, day)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating warm-up schedule: %w", err)
	}

	return schedule, nil
}

// CreateSchedule creates a new warm-up schedule for a domain
func (r *warmUpRepository) CreateSchedule(ctx context.Context, domainName string, schedule []*domain.WarmUpDay) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing schedule if any
	_, err = tx.ExecContext(ctx, "DELETE FROM warm_up_schedules WHERE domain = ?", domainName)
	if err != nil {
		return fmt.Errorf("failed to delete existing schedule: %w", err)
	}

	// Insert new schedule
	query := `
		INSERT INTO warm_up_schedules (domain, day, max_volume, actual_volume, created_at)
		VALUES (?, ?, ?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, day := range schedule {
		_, err = stmt.ExecContext(ctx, domainName, day.Day, day.MaxVolume, day.ActualVolume, day.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert warm-up day: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// UpdateDayVolume updates the actual volume sent for a specific day
func (r *warmUpRepository) UpdateDayVolume(ctx context.Context, domainName string, day int, volume int) error {
	query := `
		UPDATE warm_up_schedules
		SET actual_volume = ?
		WHERE domain = ? AND day = ?
	`

	result, err := r.db.ExecContext(ctx, query, volume, domainName, day)
	if err != nil {
		return fmt.Errorf("failed to update day volume: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no warm-up day found for domain %s day %d", domainName, day)
	}

	return nil
}

// IncrementDayVolume increments the actual volume sent for a specific day
func (r *warmUpRepository) IncrementDayVolume(ctx context.Context, domainName string, day int, increment int) error {
	query := `
		UPDATE warm_up_schedules
		SET actual_volume = actual_volume + ?
		WHERE domain = ? AND day = ?
	`

	result, err := r.db.ExecContext(ctx, query, increment, domainName, day)
	if err != nil {
		return fmt.Errorf("failed to increment day volume: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no warm-up day found for domain %s day %d", domainName, day)
	}

	return nil
}

// DeleteSchedule removes the warm-up schedule for a domain
func (r *warmUpRepository) DeleteSchedule(ctx context.Context, domainName string) error {
	query := `DELETE FROM warm_up_schedules WHERE domain = ?`

	_, err := r.db.ExecContext(ctx, query, domainName)
	if err != nil {
		return fmt.Errorf("failed to delete warm-up schedule: %w", err)
	}

	return nil
}
