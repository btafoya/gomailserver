package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository"
)

type customWarmupRepository struct {
	db *sql.DB
}

// NewCustomWarmupRepository creates a new SQLite custom warm-up repository
func NewCustomWarmupRepository(db *sql.DB) repository.CustomWarmupRepository {
	return &customWarmupRepository{db: db}
}

// CreateSchedule creates a custom warm-up schedule
func (r *customWarmupRepository) CreateSchedule(ctx context.Context, schedule []*domain.CustomWarmupSchedule) error {
	if len(schedule) == 0 {
		return fmt.Errorf("schedule cannot be empty")
	}

	// Start a transaction for batch insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `
		INSERT INTO custom_warmup_schedules (
			domain, schedule_name, day, max_volume, created_at, created_by, is_active
		) VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(domain, day) DO UPDATE SET
			schedule_name = excluded.schedule_name,
			max_volume = excluded.max_volume,
			created_at = excluded.created_at,
			created_by = excluded.created_by,
			is_active = excluded.is_active
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, day := range schedule {
		result, err := stmt.ExecContext(ctx,
			day.Domain,
			day.ScheduleName,
			day.Day,
			day.MaxVolume,
			day.CreatedAt,
			day.CreatedBy,
			day.IsActive,
		)
		if err != nil {
			return fmt.Errorf("failed to insert schedule day %d: %w", day.Day, err)
		}

		if day.ID == 0 {
			id, err := result.LastInsertId()
			if err != nil {
				return fmt.Errorf("failed to get last insert id: %w", err)
			}
			day.ID = id
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetSchedule retrieves custom warm-up schedule for a domain
func (r *customWarmupRepository) GetSchedule(ctx context.Context, domainName string) ([]*domain.CustomWarmupSchedule, error) {
	query := `
		SELECT
			id, domain, schedule_name, day, max_volume, created_at, created_by, is_active
		FROM custom_warmup_schedules
		WHERE domain = ?
		ORDER BY day ASC
	`

	rows, err := r.db.QueryContext(ctx, query, domainName)
	if err != nil {
		return nil, fmt.Errorf("failed to get custom warmup schedule: %w", err)
	}
	defer rows.Close()

	var schedule []*domain.CustomWarmupSchedule
	for rows.Next() {
		day := &domain.CustomWarmupSchedule{}
		err := rows.Scan(
			&day.ID,
			&day.Domain,
			&day.ScheduleName,
			&day.Day,
			&day.MaxVolume,
			&day.CreatedAt,
			&day.CreatedBy,
			&day.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan warmup schedule: %w", err)
		}
		schedule = append(schedule, day)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating warmup schedule: %w", err)
	}

	return schedule, nil
}

// UpdateSchedule updates an existing schedule
func (r *customWarmupRepository) UpdateSchedule(ctx context.Context, schedule *domain.CustomWarmupSchedule) error {
	query := `
		UPDATE custom_warmup_schedules
		SET schedule_name = ?,
		    max_volume = ?,
		    is_active = ?
		WHERE domain = ? AND day = ?
	`

	_, err := r.db.ExecContext(ctx, query,
		schedule.ScheduleName,
		schedule.MaxVolume,
		schedule.IsActive,
		schedule.Domain,
		schedule.Day,
	)
	if err != nil {
		return fmt.Errorf("failed to update warmup schedule: %w", err)
	}

	return nil
}

// DeleteSchedule deletes a custom schedule
func (r *customWarmupRepository) DeleteSchedule(ctx context.Context, domainName string) error {
	query := `
		DELETE FROM custom_warmup_schedules
		WHERE domain = ?
	`

	_, err := r.db.ExecContext(ctx, query, domainName)
	if err != nil {
		return fmt.Errorf("failed to delete warmup schedule: %w", err)
	}

	return nil
}

// ListActiveSchedules returns all active custom schedules
func (r *customWarmupRepository) ListActiveSchedules(ctx context.Context) (map[string][]*domain.CustomWarmupSchedule, error) {
	query := `
		SELECT
			id, domain, schedule_name, day, max_volume, created_at, created_by, is_active
		FROM custom_warmup_schedules
		WHERE is_active = 1
		ORDER BY domain ASC, day ASC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active warmup schedules: %w", err)
	}
	defer rows.Close()

	schedules := make(map[string][]*domain.CustomWarmupSchedule)
	for rows.Next() {
		day := &domain.CustomWarmupSchedule{}
		err := rows.Scan(
			&day.ID,
			&day.Domain,
			&day.ScheduleName,
			&day.Day,
			&day.MaxVolume,
			&day.CreatedAt,
			&day.CreatedBy,
			&day.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan warmup schedule: %w", err)
		}

		schedules[day.Domain] = append(schedules[day.Domain], day)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating warmup schedules: %w", err)
	}

	return schedules, nil
}

// SetActive activates/deactivates a schedule
func (r *customWarmupRepository) SetActive(ctx context.Context, domainName string, active bool) error {
	query := `
		UPDATE custom_warmup_schedules
		SET is_active = ?
		WHERE domain = ?
	`

	_, err := r.db.ExecContext(ctx, query, active, domainName)
	if err != nil {
		return fmt.Errorf("failed to set warmup schedule active state: %w", err)
	}

	return nil
}
