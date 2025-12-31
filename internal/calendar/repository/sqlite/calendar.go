package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
)

// CalendarRepository implements domain.CalendarRepository for SQLite
type CalendarRepository struct {
	db *sql.DB
}

// NewCalendarRepository creates a new SQLite calendar repository
func NewCalendarRepository(db *sql.DB) *CalendarRepository {
	return &CalendarRepository{db: db}
}

// Create creates a new calendar
func (r *CalendarRepository) Create(calendar *domain.Calendar) error {
	query := `
		INSERT INTO calendars (user_id, name, display_name, color, description, timezone, sync_token, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		calendar.UserID,
		calendar.Name,
		calendar.DisplayName,
		calendar.Color,
		calendar.Description,
		calendar.Timezone,
		calendar.SyncToken,
		calendar.CreatedAt,
		calendar.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	calendar.ID = id

	return nil
}

// GetByID retrieves a calendar by ID
func (r *CalendarRepository) GetByID(id int64) (*domain.Calendar, error) {
	query := `
		SELECT id, user_id, name, display_name, color, description, timezone, sync_token, created_at, updated_at
		FROM calendars
		WHERE id = ?
	`
	calendar := &domain.Calendar{}
	err := r.db.QueryRow(query, id).Scan(
		&calendar.ID,
		&calendar.UserID,
		&calendar.Name,
		&calendar.DisplayName,
		&calendar.Color,
		&calendar.Description,
		&calendar.Timezone,
		&calendar.SyncToken,
		&calendar.CreatedAt,
		&calendar.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return calendar, nil
}

// GetByUserID retrieves all calendars for a user
func (r *CalendarRepository) GetByUserID(userID int64) ([]*domain.Calendar, error) {
	query := `
		SELECT id, user_id, name, display_name, color, description, timezone, sync_token, created_at, updated_at
		FROM calendars
		WHERE user_id = ?
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*domain.Calendar
	for rows.Next() {
		calendar := &domain.Calendar{}
		err := rows.Scan(
			&calendar.ID,
			&calendar.UserID,
			&calendar.Name,
			&calendar.DisplayName,
			&calendar.Color,
			&calendar.Description,
			&calendar.Timezone,
			&calendar.SyncToken,
			&calendar.CreatedAt,
			&calendar.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		calendars = append(calendars, calendar)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return calendars, nil
}

// GetByUserAndName retrieves a calendar by user ID and name
func (r *CalendarRepository) GetByUserAndName(userID int64, name string) (*domain.Calendar, error) {
	query := `
		SELECT id, user_id, name, display_name, color, description, timezone, sync_token, created_at, updated_at
		FROM calendars
		WHERE user_id = ? AND name = ?
	`
	calendar := &domain.Calendar{}
	err := r.db.QueryRow(query, userID, name).Scan(
		&calendar.ID,
		&calendar.UserID,
		&calendar.Name,
		&calendar.DisplayName,
		&calendar.Color,
		&calendar.Description,
		&calendar.Timezone,
		&calendar.SyncToken,
		&calendar.CreatedAt,
		&calendar.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return calendar, nil
}

// Update updates an existing calendar
func (r *CalendarRepository) Update(calendar *domain.Calendar) error {
	query := `
		UPDATE calendars
		SET display_name = ?, color = ?, description = ?, timezone = ?, sync_token = ?, updated_at = ?
		WHERE id = ?
	`
	calendar.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		calendar.DisplayName,
		calendar.Color,
		calendar.Description,
		calendar.Timezone,
		calendar.SyncToken,
		calendar.UpdatedAt,
		calendar.ID,
	)
	return err
}

// Delete deletes a calendar
func (r *CalendarRepository) Delete(id int64) error {
	query := `DELETE FROM calendars WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateSyncToken updates the sync token for a calendar
func (r *CalendarRepository) UpdateSyncToken(id int64, token string) error {
	query := `UPDATE calendars SET sync_token = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, token, time.Now(), id)
	return err
}
