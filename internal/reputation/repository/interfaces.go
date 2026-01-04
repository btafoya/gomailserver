package repository

import (
	"context"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
)

// EventsRepository handles storage and retrieval of sending events
type EventsRepository interface {
	// RecordEvent stores a new sending event
	RecordEvent(ctx context.Context, event *domain.SendingEvent) error

	// GetEventsInWindow retrieves events for a domain within a time window
	GetEventsInWindow(ctx context.Context, domain string, startTime, endTime int64) ([]*domain.SendingEvent, error)

	// GetEventCountsByType returns counts of each event type for a domain in a time window
	GetEventCountsByType(ctx context.Context, domain string, startTime, endTime int64) (map[string]int64, error)

	// CleanupOldEvents removes events older than the specified timestamp
	CleanupOldEvents(ctx context.Context, olderThan int64) error
}

// ScoresRepository handles domain reputation scores
type ScoresRepository interface {
	// GetReputationScore retrieves the reputation score for a domain
	GetReputationScore(ctx context.Context, domain string) (*domain.ReputationScore, error)

	// UpdateReputationScore updates or creates a reputation score for a domain
	UpdateReputationScore(ctx context.Context, score *domain.ReputationScore) error

	// ListAllScores retrieves all domain reputation scores
	ListAllScores(ctx context.Context) ([]*domain.ReputationScore, error)
}

// WarmUpRepository handles warm-up schedule management
type WarmUpRepository interface {
	// GetSchedule retrieves the warm-up schedule for a domain
	GetSchedule(ctx context.Context, domain string) ([]*domain.WarmUpDay, error)

	// CreateSchedule creates a new warm-up schedule for a domain
	CreateSchedule(ctx context.Context, domain string, schedule []*domain.WarmUpDay) error

	// UpdateDayVolume sets the actual volume sent for a specific day
	UpdateDayVolume(ctx context.Context, domain string, day int, volume int) error

	// IncrementDayVolume increments the actual volume sent for a specific day
	IncrementDayVolume(ctx context.Context, domain string, day int, increment int) error

	// DeleteSchedule removes the warm-up schedule for a domain
	DeleteSchedule(ctx context.Context, domain string) error
}

// CircuitBreakerRepository handles circuit breaker event tracking
type CircuitBreakerRepository interface {
	// RecordPause creates a new circuit breaker pause event
	RecordPause(ctx context.Context, event *domain.CircuitBreakerEvent) error

	// RecordResume marks a circuit breaker as resumed
	RecordResume(ctx context.Context, domain string, autoResumed bool, notes string) error

	// GetActiveBreakers retrieves all currently active circuit breakers
	GetActiveBreakers(ctx context.Context) ([]*domain.CircuitBreakerEvent, error)

	// GetBreakerHistory retrieves circuit breaker history for a domain
	GetBreakerHistory(ctx context.Context, domain string, limit int) ([]*domain.CircuitBreakerEvent, error)
}

// RetentionPolicyRepository handles retention policy management
type RetentionPolicyRepository interface {
	// GetPolicy retrieves the current retention policy
	GetPolicy(ctx context.Context) (*domain.RetentionPolicy, error)

	// UpdatePolicy updates the retention policy
	UpdatePolicy(ctx context.Context, policy *domain.RetentionPolicy) error
}
