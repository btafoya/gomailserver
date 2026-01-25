package sqlite

import (
	"database/sql"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
)

// TaskRepository implements domain.TaskRepository for SQLite
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new SQLite task repository
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create creates a new task
func (r *TaskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query,
		task.CalendarID,
		task.UID,
		task.Summary,
		task.Description,
		task.Location,
		task.Due,
		task.Start,
		task.Completed,
		task.Status,
		task.Priority,
		task.Percent,
		task.Categories,
		task.Organizer,
		task.Attendees,
		task.Sequence,
		task.ETag,
		task.ICalData,
		task.CreatedAt,
		task.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = id

	return nil
}

// GetByID retrieves a task by ID
func (r *TaskRepository) GetByID(id int64) (*domain.Task, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	task := &domain.Task{}
	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.CalendarID,
		&task.UID,
		&task.Summary,
		&task.Description,
		&task.Location,
		&task.Due,
		&task.Start,
		&task.Completed,
		&task.Status,
		&task.Priority,
		&task.Percent,
		&task.Categories,
		&task.Organizer,
		&task.Attendees,
		&task.Sequence,
		&task.ETag,
		&task.ICalData,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetByUID retrieves a task by UID and calendar ID
func (r *TaskRepository) GetByUID(calendarID int64, uid string) (*domain.Task, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at
		FROM tasks
		WHERE calendar_id = ? AND uid = ?
	`
	task := &domain.Task{}
	err := r.db.QueryRow(query, calendarID, uid).Scan(
		&task.ID,
		&task.CalendarID,
		&task.UID,
		&task.Summary,
		&task.Description,
		&task.Location,
		&task.Due,
		&task.Start,
		&task.Completed,
		&task.Status,
		&task.Priority,
		&task.Percent,
		&task.Categories,
		&task.Organizer,
		&task.Attendees,
		&task.Sequence,
		&task.ETag,
		&task.ICalData,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetByCalendar retrieves all tasks for a calendar
func (r *TaskRepository) GetByCalendar(calendarID int64) ([]*domain.Task, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at
		FROM tasks
		WHERE calendar_id = ?
		ORDER BY COALESCE(due, '9999-12-31') ASC, priority ASC
	`
	rows, err := r.db.Query(query, calendarID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// GetByStatus retrieves tasks by status
func (r *TaskRepository) GetByStatus(calendarID int64, status string) ([]*domain.Task, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at
		FROM tasks
		WHERE calendar_id = ? AND status = ?
		ORDER BY COALESCE(due, '9999-12-31') ASC, priority ASC
	`
	rows, err := r.db.Query(query, calendarID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// GetOverdue retrieves overdue tasks
func (r *TaskRepository) GetOverdue(calendarID int64) ([]*domain.Task, error) {
	query := `
		SELECT id, calendar_id, uid, summary, description, location, due, start, completed, status, priority, percent, categories, organizer, attendees, sequence, etag, ical_data, created_at, updated_at
		FROM tasks
		WHERE calendar_id = ? AND due IS NOT NULL AND due < ? AND status != 'COMPLETED' AND status != 'CANCELLED'
		ORDER BY due ASC, priority ASC
	`
	rows, err := r.db.Query(query, calendarID, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTasks(rows)
}

// Update updates an existing task
func (r *TaskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET summary = ?, description = ?, location = ?, due = ?, start = ?, completed = ?, status = ?, priority = ?, percent = ?, categories = ?, organizer = ?, attendees = ?, sequence = ?, etag = ?, ical_data = ?, updated_at = ?
		WHERE id = ?
	`
	task.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		task.Summary,
		task.Description,
		task.Location,
		task.Due,
		task.Start,
		task.Completed,
		task.Status,
		task.Priority,
		task.Percent,
		task.Categories,
		task.Organizer,
		task.Attendees,
		task.Sequence,
		task.ETag,
		task.ICalData,
		task.UpdatedAt,
		task.ID,
	)
	return err
}

// Delete deletes a task
func (r *TaskRepository) Delete(id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// UpdateETag updates the ETag for a task
func (r *TaskRepository) UpdateETag(id int64, etag string) error {
	query := `UPDATE tasks SET etag = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, etag, time.Now(), id)
	return err
}

// scanTasks scans rows into task slice
func (r *TaskRepository) scanTasks(rows *sql.Rows) ([]*domain.Task, error) {
	var tasks []*domain.Task
	for rows.Next() {
		task := &domain.Task{}
		err := rows.Scan(
			&task.ID,
			&task.CalendarID,
			&task.UID,
			&task.Summary,
			&task.Description,
			&task.Location,
			&task.Due,
			&task.Start,
			&task.Completed,
			&task.Status,
			&task.Priority,
			&task.Percent,
			&task.Categories,
			&task.Organizer,
			&task.Attendees,
			&task.Sequence,
			&task.ETag,
			&task.ICalData,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
