package service

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/btafoya/gomailserver/internal/calendar/domain"
	"github.com/emersion/go-ical"
)

// TaskService implements domain.TaskService
type TaskService struct {
	taskRepo     domain.TaskRepository
	calendarRepo domain.CalendarRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo domain.TaskRepository, calendarRepo domain.CalendarRepository) *TaskService {
	return &TaskService{
		taskRepo:     taskRepo,
		calendarRepo: calendarRepo,
	}
}

// CreateTask creates a new task from iCalendar data
func (s *TaskService) CreateTask(calendarID int64, icalData string) (*domain.Task, error) {
	// Verify calendar exists
	calendar, err := s.calendarRepo.GetByID(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar: %w", err)
	}
	if calendar == nil {
		return nil, fmt.Errorf("calendar not found")
	}

	// Parse iCalendar data
	dec := ical.NewDecoder(strings.NewReader(icalData))
	cal, err := dec.Decode()
	if err != nil {
		return nil, fmt.Errorf("failed to parse iCalendar data: %w", err)
	}

	// Extract VTODO component
	var vtodo *ical.Component
	for _, comp := range cal.Children {
		if comp.Name == "VTODO" {
			vtodo = comp
			break
		}
	}
	if vtodo == nil {
		return nil, fmt.Errorf("no VTODO component found in iCalendar data")
	}

	// Extract task properties
	task := &domain.Task{
		CalendarID: calendarID,
		ICalData:   icalData,
		Status:     domain.TaskStatusNeedsAction,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Extract UID (required)
	if prop := vtodo.Props.Get("UID"); prop != nil {
		task.UID = prop.Value
	} else {
		return nil, fmt.Errorf("UID property is required")
	}

	// Check if task with same UID already exists
	existing, err := s.taskRepo.GetByUID(calendarID, task.UID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing task: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("task with UID %s already exists", task.UID)
	}

	// Extract other properties
	s.extractTaskProperties(vtodo, task)

	// Generate ETag
	task.ETag = s.GenerateETag(task)

	// Create task in repository
	if err := s.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetTask retrieves a task by ID
func (s *TaskService) GetTask(id int64) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

// GetCalendarTasks retrieves all tasks for a calendar
func (s *TaskService) GetCalendarTasks(calendarID int64) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetByCalendar(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get calendar tasks: %w", err)
	}
	return tasks, nil
}

// GetTasksByStatus retrieves tasks by status
func (s *TaskService) GetTasksByStatus(calendarID int64, status string) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetByStatus(calendarID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	return tasks, nil
}

// GetOverdueTasks retrieves overdue tasks
func (s *TaskService) GetOverdueTasks(calendarID int64) ([]*domain.Task, error) {
	tasks, err := s.taskRepo.GetOverdue(calendarID)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	return tasks, nil
}

// UpdateTask updates a task from iCalendar data
func (s *TaskService) UpdateTask(id int64, icalData string) error {
	// Get existing task
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	// Parse new iCalendar data
	dec := ical.NewDecoder(strings.NewReader(icalData))
	cal, err := dec.Decode()
	if err != nil {
		return fmt.Errorf("failed to parse iCalendar data: %w", err)
	}

	// Extract VTODO component
	var vtodo *ical.Component
	for _, comp := range cal.Children {
		if comp.Name == "VTODO" {
			vtodo = comp
			break
		}
	}
	if vtodo == nil {
		return fmt.Errorf("no VTODO component found in iCalendar data")
	}

	// Update task properties
	task.ICalData = icalData
	task.UpdatedAt = time.Now()
	s.extractTaskProperties(vtodo, task)

	// Update ETag
	task.ETag = s.GenerateETag(task)

	// Update in repository
	if err := s.taskRepo.Update(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(id int64) error {
	task, err := s.taskRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return fmt.Errorf("task not found")
	}

	now := time.Now()
	task.Status = domain.TaskStatusCompleted
	task.Completed = &now
	task.Percent = 100
	task.UpdatedAt = now

	// Regenerate iCalendar data with updated status
	task.ICalData = s.generateTaskICalData(task)
	task.ETag = s.GenerateETag(task)

	if err := s.taskRepo.Update(task); err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id int64) error {
	if err := s.taskRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}

// GenerateETag generates a new ETag for a task
func (s *TaskService) GenerateETag(task *domain.Task) string {
	data := fmt.Sprintf("%s-%s-%d", task.UID, task.UpdatedAt.Format(time.RFC3339), task.Sequence)
	hash := sha256.Sum256([]byte(data))
	return `"` + hex.EncodeToString(hash[:]) + `"`
}

// extractTaskProperties extracts VTODO properties into a Task struct
func (s *TaskService) extractTaskProperties(vtodo *ical.Component, task *domain.Task) {
	if prop := vtodo.Props.Get("SUMMARY"); prop != nil {
		task.Summary = prop.Value
	}
	if prop := vtodo.Props.Get("DESCRIPTION"); prop != nil {
		task.Description = prop.Value
	}
	if prop := vtodo.Props.Get("LOCATION"); prop != nil {
		task.Location = prop.Value
	}
	if prop := vtodo.Props.Get("STATUS"); prop != nil {
		task.Status = prop.Value
	}
	if prop := vtodo.Props.Get("PRIORITY"); prop != nil {
		fmt.Sscanf(prop.Value, "%d", &task.Priority)
	}
	if prop := vtodo.Props.Get("PERCENT-COMPLETE"); prop != nil {
		fmt.Sscanf(prop.Value, "%d", &task.Percent)
	}
	if prop := vtodo.Props.Get("ORGANIZER"); prop != nil {
		task.Organizer = prop.Value
	}
	if prop := vtodo.Props.Get("SEQUENCE"); prop != nil {
		fmt.Sscanf(prop.Value, "%d", &task.Sequence)
	}
	if prop := vtodo.Props.Get("CATEGORIES"); prop != nil {
		task.Categories = prop.Value
	}

	// Extract DTSTART
	if prop := vtodo.Props.Get("DTSTART"); prop != nil {
		dtstart, err := prop.DateTime(time.UTC)
		if err == nil {
			task.Start = &dtstart
		}
	}

	// Extract DUE
	if prop := vtodo.Props.Get("DUE"); prop != nil {
		due, err := prop.DateTime(time.UTC)
		if err == nil {
			task.Due = &due
		}
	}

	// Extract COMPLETED
	if prop := vtodo.Props.Get("COMPLETED"); prop != nil {
		completed, err := prop.DateTime(time.UTC)
		if err == nil {
			task.Completed = &completed
		}
	}
}

// generateTaskICalData generates iCalendar data for a task
func (s *TaskService) generateTaskICalData(task *domain.Task) string {
	var sb strings.Builder
	sb.WriteString("BEGIN:VCALENDAR\r\n")
	sb.WriteString("VERSION:2.0\r\n")
	sb.WriteString("PRODID:-//gomailserver//CalDAV Server//EN\r\n")
	sb.WriteString("BEGIN:VTODO\r\n")
	sb.WriteString(fmt.Sprintf("UID:%s\r\n", task.UID))
	sb.WriteString(fmt.Sprintf("DTSTAMP:%s\r\n", time.Now().UTC().Format("20060102T150405Z")))

	if task.Summary != "" {
		sb.WriteString(fmt.Sprintf("SUMMARY:%s\r\n", escapeICalText(task.Summary)))
	}
	if task.Description != "" {
		sb.WriteString(fmt.Sprintf("DESCRIPTION:%s\r\n", escapeICalText(task.Description)))
	}
	if task.Location != "" {
		sb.WriteString(fmt.Sprintf("LOCATION:%s\r\n", escapeICalText(task.Location)))
	}
	if task.Status != "" {
		sb.WriteString(fmt.Sprintf("STATUS:%s\r\n", task.Status))
	}
	if task.Priority > 0 {
		sb.WriteString(fmt.Sprintf("PRIORITY:%d\r\n", task.Priority))
	}
	if task.Percent > 0 {
		sb.WriteString(fmt.Sprintf("PERCENT-COMPLETE:%d\r\n", task.Percent))
	}
	if task.Organizer != "" {
		sb.WriteString(fmt.Sprintf("ORGANIZER:%s\r\n", task.Organizer))
	}
	if task.Sequence > 0 {
		sb.WriteString(fmt.Sprintf("SEQUENCE:%d\r\n", task.Sequence))
	}
	if task.Categories != "" {
		sb.WriteString(fmt.Sprintf("CATEGORIES:%s\r\n", task.Categories))
	}
	if task.Start != nil {
		sb.WriteString(fmt.Sprintf("DTSTART:%s\r\n", task.Start.UTC().Format("20060102T150405Z")))
	}
	if task.Due != nil {
		sb.WriteString(fmt.Sprintf("DUE:%s\r\n", task.Due.UTC().Format("20060102T150405Z")))
	}
	if task.Completed != nil {
		sb.WriteString(fmt.Sprintf("COMPLETED:%s\r\n", task.Completed.UTC().Format("20060102T150405Z")))
	}

	sb.WriteString("END:VTODO\r\n")
	sb.WriteString("END:VCALENDAR\r\n")

	return sb.String()
}
