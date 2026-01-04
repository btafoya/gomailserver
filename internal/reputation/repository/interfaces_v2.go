package repository

import (
	"context"
	"github.com/btafoya/gomailserver/internal/reputation/domain"
)

// DMARC Reports Repository

type DMARCReportsRepository interface {
	// Create stores a new DMARC report
	Create(ctx context.Context, report *domain.DMARCReport) error

	// GetByID retrieves a report by ID
	GetByID(ctx context.Context, id int64) (*domain.DMARCReport, error)

	// GetByReportID retrieves a report by its external report_id
	GetByReportID(ctx context.Context, reportID string) (*domain.DMARCReport, error)

	// ListByDomain returns all reports for a domain
	ListByDomain(ctx context.Context, domain string, limit, offset int) ([]*domain.DMARCReport, error)

	// ListByTimeRange returns reports within a time range
	ListByTimeRange(ctx context.Context, startTime, endTime int64, limit, offset int) ([]*domain.DMARCReport, error)

	// GetDomainStats returns aggregated statistics for a domain
	GetDomainStats(ctx context.Context, domain string, days int) (*domain.AlignmentAnalysis, error)

	// CreateRecord stores a DMARC report record
	CreateRecord(ctx context.Context, record *domain.DMARCReportRecord) error

	// GetRecordsByReportID retrieves all records for a report
	GetRecordsByReportID(ctx context.Context, reportID int64) ([]*domain.DMARCReportRecord, error)
}

type DMARCActionsRepository interface {
	// RecordAction logs an automated action taken
	RecordAction(ctx context.Context, action *domain.DMARCAutoAction) error

	// ListActions returns recent automated actions
	ListActions(ctx context.Context, domain string, limit int) ([]*domain.DMARCAutoAction, error)

	// ListAllActions returns all actions with pagination
	ListAllActions(ctx context.Context, limit, offset int) ([]*domain.DMARCAutoAction, error)
}

// ARF Reports Repository

type ARFReportsRepository interface {
	// Create stores a new ARF complaint report
	Create(ctx context.Context, report *domain.ARFReport) error

	// GetByID retrieves a report by ID
	GetByID(ctx context.Context, id int64) (*domain.ARFReport, error)

	// ListUnprocessed returns unprocessed complaints
	ListUnprocessed(ctx context.Context, limit int) ([]*domain.ARFReport, error)

	// MarkProcessed marks a report as processed
	MarkProcessed(ctx context.Context, id int64, suppressedRecipient string) error

	// ListByTimeRange returns complaints within a time range
	ListByTimeRange(ctx context.Context, startTime, endTime int64, limit, offset int) ([]*domain.ARFReport, error)

	// GetComplaintRate calculates complaint rate for a domain/IP
	GetComplaintRate(ctx context.Context, domain string, hours int) (float64, error)
}

// External Metrics Repositories

type PostmasterMetricsRepository interface {
	// Create stores new Gmail Postmaster metrics
	Create(ctx context.Context, metrics *domain.PostmasterMetrics) error

	// GetLatest returns the latest metrics for a domain
	GetLatest(ctx context.Context, domain string) (*domain.PostmasterMetrics, error)

	// ListByDomain returns metrics history for a domain
	ListByDomain(ctx context.Context, domain string, days int) ([]*domain.PostmasterMetrics, error)

	// GetReputationTrend returns domain reputation trend
	GetReputationTrend(ctx context.Context, domain string, days int) ([]string, error)
}

type SNDSMetricsRepository interface {
	// Create stores new Microsoft SNDS metrics
	Create(ctx context.Context, metrics *domain.SNDSMetrics) error

	// GetLatest returns the latest metrics for an IP
	GetLatest(ctx context.Context, ipAddress string) (*domain.SNDSMetrics, error)

	// ListByIP returns metrics history for an IP
	ListByIP(ctx context.Context, ipAddress string, days int) ([]*domain.SNDSMetrics, error)

	// GetFilterLevelTrend returns filter level trend
	GetFilterLevelTrend(ctx context.Context, ipAddress string, days int) ([]string, error)
}

// Provider Rate Limits Repository

type ProviderRateLimitsRepository interface {
	// Get retrieves rate limit for a domain and provider
	Get(ctx context.Context, domain string, provider domain.MailProvider) (*domain.ProviderRateLimit, error)

	// CreateOrUpdate creates or updates rate limit
	CreateOrUpdate(ctx context.Context, limit *domain.ProviderRateLimit) error

	// IncrementHourly increments hourly counter
	IncrementHourly(ctx context.Context, domain string, provider domain.MailProvider, count int) error

	// IncrementDaily increments daily counter
	IncrementDaily(ctx context.Context, domain string, provider domain.MailProvider, count int) error

	// ResetHourly resets hourly counter
	ResetHourly(ctx context.Context, domain string, provider domain.MailProvider, newResetTime int64) error

	// ResetDaily resets daily counter
	ResetDaily(ctx context.Context, domain string, provider domain.MailProvider, newResetTime int64) error

	// ListByDomain returns all provider limits for a domain
	ListByDomain(ctx context.Context, domain string) ([]*domain.ProviderRateLimit, error)

	// SetCircuitBreaker activates/deactivates circuit breaker
	SetCircuitBreaker(ctx context.Context, domain string, provider domain.MailProvider, active bool) error
}

// Custom Warm-up Repository

type CustomWarmupRepository interface {
	// CreateSchedule creates a custom warm-up schedule
	CreateSchedule(ctx context.Context, schedule []*domain.CustomWarmupSchedule) error

	// GetSchedule retrieves custom warm-up schedule for a domain
	GetSchedule(ctx context.Context, domain string) ([]*domain.CustomWarmupSchedule, error)

	// UpdateSchedule updates an existing schedule
	UpdateSchedule(ctx context.Context, schedule *domain.CustomWarmupSchedule) error

	// DeleteSchedule deletes a custom schedule
	DeleteSchedule(ctx context.Context, domain string) error

	// ListActiveSchedules returns all active custom schedules
	ListActiveSchedules(ctx context.Context) (map[string][]*domain.CustomWarmupSchedule, error)

	// SetActive activates/deactivates a schedule
	SetActive(ctx context.Context, domain string, active bool) error
}

// Predictions Repository

type PredictionsRepository interface {
	// Create stores a new prediction
	Create(ctx context.Context, prediction *domain.ReputationPrediction) error

	// GetLatest returns the latest prediction for a domain
	GetLatest(ctx context.Context, domain string) (*domain.ReputationPrediction, error)

	// ListByDomain returns prediction history for a domain
	ListByDomain(ctx context.Context, domain string, limit int) ([]*domain.ReputationPrediction, error)

	// GetByHorizon returns predictions for a specific time horizon
	GetByHorizon(ctx context.Context, domain string, hours int) (*domain.ReputationPrediction, error)
}

// Alerts Repository

type AlertsRepository interface {
	// Create stores a new alert
	Create(ctx context.Context, alert *domain.ReputationAlert) error

	// GetByID retrieves an alert by ID
	GetByID(ctx context.Context, id int64) (*domain.ReputationAlert, error)

	// ListUnacknowledged returns unacknowledged alerts
	ListUnacknowledged(ctx context.Context, limit int) ([]*domain.ReputationAlert, error)

	// ListByDomain returns alerts for a domain
	ListByDomain(ctx context.Context, domain string, limit, offset int) ([]*domain.ReputationAlert, error)

	// ListBySeverity returns alerts by severity level
	ListBySeverity(ctx context.Context, severity domain.AlertSeverity, limit int) ([]*domain.ReputationAlert, error)

	// Acknowledge marks an alert as acknowledged
	Acknowledge(ctx context.Context, id int64, acknowledgedBy string) error

	// Resolve marks an alert as resolved
	Resolve(ctx context.Context, id int64) error

	// GetUnacknowledgedCount returns count of unacknowledged alerts
	GetUnacknowledgedCount(ctx context.Context) (int, error)

	// GetUnacknowledgedCountByDomain returns count of unacknowledged alerts for a domain
	GetUnacknowledgedCountByDomain(ctx context.Context, domain string) (int, error)
}
