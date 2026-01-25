package domain

import (
	"time"
)

// Task represents a CalDAV VTODO component
type Task struct {
	ID          int64
	CalendarID  int64
	UID         string
	Summary     string
	Description string
	Location    string
	Due         *time.Time // DTDUE
	Start       *time.Time // DTSTART (optional)
	Completed   *time.Time // COMPLETED
	Status      string     // NEEDS-ACTION, IN-PROCESS, COMPLETED, CANCELLED
	Priority    int        // 0-9 (0 = undefined, 1 = highest, 9 = lowest)
	Percent     int        // PERCENT-COMPLETE (0-100)
	Categories  string     // JSON array of categories
	Organizer   string
	Attendees   string // JSON array of attendees
	Sequence    int
	ETag        string
	ICalData    string // Full iCalendar data
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	// Create creates a new task
	Create(task *Task) error

	// GetByID retrieves a task by ID
	GetByID(id int64) (*Task, error)

	// GetByUID retrieves a task by UID and calendar ID
	GetByUID(calendarID int64, uid string) (*Task, error)

	// GetByCalendar retrieves all tasks for a calendar
	GetByCalendar(calendarID int64) ([]*Task, error)

	// GetByStatus retrieves tasks by status
	GetByStatus(calendarID int64, status string) ([]*Task, error)

	// GetOverdue retrieves overdue tasks
	GetOverdue(calendarID int64) ([]*Task, error)

	// Update updates an existing task
	Update(task *Task) error

	// Delete deletes a task
	Delete(id int64) error

	// UpdateETag updates the ETag for a task
	UpdateETag(id int64, etag string) error
}

// TaskService defines business logic for tasks
type TaskService interface {
	// CreateTask creates a new task from iCalendar data
	CreateTask(calendarID int64, icalData string) (*Task, error)

	// GetTask retrieves a task by ID
	GetTask(id int64) (*Task, error)

	// GetCalendarTasks retrieves all tasks for a calendar
	GetCalendarTasks(calendarID int64) ([]*Task, error)

	// GetTasksByStatus retrieves tasks by status
	GetTasksByStatus(calendarID int64, status string) ([]*Task, error)

	// GetOverdueTasks retrieves overdue tasks
	GetOverdueTasks(calendarID int64) ([]*Task, error)

	// UpdateTask updates a task from iCalendar data
	UpdateTask(id int64, icalData string) error

	// CompleteTask marks a task as completed
	CompleteTask(id int64) error

	// DeleteTask deletes a task
	DeleteTask(id int64) error

	// GenerateETag generates a new ETag for a task
	GenerateETag(task *Task) string
}

// TaskStatus constants
const (
	TaskStatusNeedsAction = "NEEDS-ACTION"
	TaskStatusInProcess   = "IN-PROCESS"
	TaskStatusCompleted   = "COMPLETED"
	TaskStatusCancelled   = "CANCELLED"
)

// TaskPriority constants
const (
	TaskPriorityUndefined = 0
	TaskPriorityHighest   = 1
	TaskPriorityHigh      = 2
	TaskPriorityMedium    = 5
	TaskPriorityLow       = 8
	TaskPriorityLowest    = 9
)
