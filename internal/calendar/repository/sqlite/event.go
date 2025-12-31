package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
)

// EventRepository implements domain.EventRepository for SQLite
type EventRepository struct {
	db *sql.DB
}

// NewEventRepository creates a new SQLite event repository
func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{db: db}
}

// Create creates a new event
func (r *EventRepository) Create(event *domain.Event) error {
	query := `
		INSERT INTO events (calendar_id, uid, summary, description, location, start_time, end_time, all_day, rrule, attendees, organizer, status, sequence, etag, ical_data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		event.CalendarID,
		event.UID,
		event.Summary,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.AllDay,
		event.RRule,
		event.Attendees,
		event.Organizer,
		event.Status,
		event.Sequence,
		event.ETag,
		event.ICalData,
		event.CreatedAt,
		event.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = id

	return nil
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(id int64) (*domain.Event, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, start_time, end_time, all_day, rrule, attendees, organizer, status, sequence, etag, ical_data, created_at, updated_at
		FROM events
		WHERE id = ?
	`
	event := &domain.Event{}
	err := r.db.QueryRow(query, id).Scan(
		&event.ID,
		&event.CalendarID,
		&event.UID,
		&event.Summary,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.AllDay,
		&event.RRule,
		&event.Attendees,
		&event.Organizer,
		&event.Status,
		&event.Sequence,
		&event.ETag,
		&event.ICalData,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetByUID retrieves an event by UID and calendar ID
func (r *EventRepository) GetByUID(calendarID int64, uid string) (*domain.Event, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, start_time, end_time, all_day, rrule, attendees, organizer, status, sequence, etag, ical_data, created_at, updated_at
		FROM events
		WHERE calendar_id = ? AND uid = ?
	`
	event := &domain.Event{}
	err := r.db.QueryRow(query, calendarID, uid).Scan(
		&event.ID,
		&event.CalendarID,
		&event.UID,
		&event.Summary,
		&event.Description,
		&event.Location,
		&event.StartTime,
		&event.EndTime,
		&event.AllDay,
		&event.RRule,
		&event.Attendees,
		&event.Organizer,
		&event.Status,
		&event.Sequence,
		&event.ETag,
		&event.ICalData,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return event, nil
}

// GetByCalendar retrieves all events for a calendar
func (r *EventRepository) GetByCalendar(calendarID int64) ([]*domain.Event, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, start_time, end_time, all_day, rrule, attendees, organizer, status, sequence, etag, ical_data, created_at, updated_at
		FROM events
		WHERE calendar_id = ?
		ORDER BY start_time ASC
	`
	rows, err := r.db.Query(query, calendarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		event := &domain.Event{}
		err := rows.Scan(
			&event.ID,
			&event.CalendarID,
			&event.UID,
			&event.Summary,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.AllDay,
			&event.RRule,
			&event.Attendees,
			&event.Organizer,
			&event.Status,
			&event.Sequence,
			&event.ETag,
			&event.ICalData,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// GetByTimeRange retrieves events within a time range
func (r *EventRepository) GetByTimeRange(calendarID int64, start, end time.Time) ([]*domain.Event, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, start_time, end_time, all_day, rrule, attendees, organizer, status, sequence, etag, ical_data, created_at, updated_at
		FROM events
		WHERE calendar_id = ? AND (
			(start_time >= ? AND start_time < ?) OR
			(end_time > ? AND end_time <= ?) OR
			(start_time < ? AND end_time > ?)
		)
		ORDER BY start_time ASC
	`
	rows, err := r.db.Query(query, calendarID, start, end, start, end, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*domain.Event
	for rows.Next() {
		event := &domain.Event{}
		err := rows.Scan(
			&event.ID,
			&event.CalendarID,
			&event.UID,
			&event.Summary,
			&event.Description,
			&event.Location,
			&event.StartTime,
			&event.EndTime,
			&event.AllDay,
			&event.RRule,
			&event.Attendees,
			&event.Organizer,
			&event.Status,
			&event.Sequence,
			&event.ETag,
			&event.ICalData,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update updates an existing event
func (r *EventRepository) Update(event *domain.Event) error {
	query := `
		UPDATE events
		SET summary = ?, description = ?, location = ?, start_time = ?, end_time = ?, all_day = ?, rrule = ?, attendees = ?, organizer = ?, status = ?, sequence = ?, etag = ?, ical_data = ?, updated_at = ?
		WHERE id = ?
	`
	event.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		event.Summary,
		event.Description,
		event.Location,
		event.StartTime,
		event.EndTime,
		event.AllDay,
		event.RRule,
		event.Attendees,
		event.Organizer,
		event.Status,
		event.Sequence,
		event.ETag,
		event.ICalData,
		event.UpdatedAt,
		event.ID,
	)
	return err
}

// Delete deletes an event
func (r *EventRepository) Delete(id int64) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateETag updates the ETag for an event
func (r *EventRepository) UpdateETag(id int64, etag string) error {
	query := `UPDATE events SET etag = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, etag, time.Now(), id)
	return err
}
