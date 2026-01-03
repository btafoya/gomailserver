package service

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/repository/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

func setupTelemetryTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create schema
	schema := `
	CREATE TABLE sending_events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		timestamp INTEGER NOT NULL,
		domain TEXT NOT NULL,
		recipient_domain TEXT NOT NULL,
		event_type TEXT NOT NULL,
		bounce_type TEXT,
		enhanced_status_code TEXT,
		smtp_response TEXT,
		ip_address TEXT NOT NULL,
		metadata TEXT
	);
	CREATE INDEX idx_sending_events_timestamp ON sending_events(timestamp);
	CREATE INDEX idx_sending_events_domain ON sending_events(domain);

	CREATE TABLE domain_reputation_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL UNIQUE,
		reputation_score INTEGER NOT NULL,
		complaint_rate REAL NOT NULL,
		bounce_rate REAL NOT NULL,
		delivery_rate REAL NOT NULL,
		circuit_breaker_active BOOLEAN DEFAULT 0,
		circuit_breaker_reason TEXT,
		warm_up_active BOOLEAN DEFAULT 0,
		warm_up_day INTEGER DEFAULT 0,
		last_updated INTEGER NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestTelemetryService_RecordDelivery(t *testing.T) {
	db := setupTelemetryTestDB(t)
	defer db.Close()

	eventsRepo := sqlite.NewEventsRepository(db)
	scoresRepo := sqlite.NewScoresRepository(db)
	logger := zap.NewNop()
	service := NewTelemetryService(eventsRepo, scoresRepo, logger)

	ctx := context.Background()

	tests := []struct {
		name            string
		domain          string
		recipientDomain string
		ip              string
		wantErr         bool
	}{
		{
			name:            "record successful delivery",
			domain:          "example.com",
			recipientDomain: "recipient.com",
			ip:              "192.168.1.1",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RecordDelivery(ctx, tt.domain, tt.recipientDomain, tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordDelivery() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTelemetryService_RecordBounce(t *testing.T) {
	db := setupTelemetryTestDB(t)
	defer db.Close()

	eventsRepo := sqlite.NewEventsRepository(db)
	scoresRepo := sqlite.NewScoresRepository(db)
	logger := zap.NewNop()
	service := NewTelemetryService(eventsRepo, scoresRepo, logger)

	ctx := context.Background()

	tests := []struct {
		name            string
		domain          string
		recipientDomain string
		ip              string
		bounceType      string
		statusCode      string
		response        string
		wantErr         bool
	}{
		{
			name:            "record hard bounce",
			domain:          "example.com",
			recipientDomain: "recipient.com",
			ip:              "192.168.1.1",
			bounceType:      "hard",
			statusCode:      "5.1.1",
			response:        "User unknown",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.RecordBounce(ctx, tt.domain, tt.recipientDomain, tt.ip, tt.bounceType, tt.statusCode, tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordBounce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTelemetryService_CalculateReputationScore(t *testing.T) {
	db := setupTelemetryTestDB(t)
	defer db.Close()

	eventsRepo := sqlite.NewEventsRepository(db)
	scoresRepo := sqlite.NewScoresRepository(db)
	logger := zap.NewNop()
	service := NewTelemetryService(eventsRepo, scoresRepo, logger)

	ctx := context.Background()
	now := time.Now().Unix()

	tests := []struct {
		name              string
		setupEvents       []*domain.SendingEvent
		domain            string
		expectedMinScore  int
		expectedMaxScore  int
		expectedBounceMin float64
		expectedBounceMax float64
	}{
		{
			name: "excellent reputation - all deliveries",
			setupEvents: []*domain.SendingEvent{
				{Timestamp: now - 3600, Domain: "example.com", RecipientDomain: "r1.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "example.com", RecipientDomain: "r2.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "example.com", RecipientDomain: "r3.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "example.com", RecipientDomain: "r4.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "example.com", RecipientDomain: "r5.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
			},
			domain:            "example.com",
			expectedMinScore:  95,
			expectedMaxScore:  100,
			expectedBounceMin: 0.0,
			expectedBounceMax: 0.1,
		},
		{
			name: "poor reputation - high bounce rate",
			setupEvents: []*domain.SendingEvent{
				{Timestamp: now - 3600, Domain: "bad.com", RecipientDomain: "r1.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "bad.com", RecipientDomain: "r2.com", EventType: domain.EventBounce, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "bad.com", RecipientDomain: "r3.com", EventType: domain.EventBounce, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "bad.com", RecipientDomain: "r4.com", EventType: domain.EventBounce, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
			},
			domain:            "bad.com",
			expectedMinScore:  0,
			expectedMaxScore:  60,
			expectedBounceMin: 70.0,
			expectedBounceMax: 80.0,
		},
		{
			name: "very poor reputation - high complaint rate",
			setupEvents: []*domain.SendingEvent{
				{Timestamp: now - 3600, Domain: "spammer.com", RecipientDomain: "r1.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "spammer.com", RecipientDomain: "r2.com", EventType: domain.EventDelivery, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "spammer.com", RecipientDomain: "r3.com", EventType: domain.EventComplaint, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
				{Timestamp: now - 3600, Domain: "spammer.com", RecipientDomain: "r4.com", EventType: domain.EventComplaint, IPAddress: "1.1.1.1", Metadata: make(map[string]interface{})},
			},
			domain:           "spammer.com",
			expectedMinScore: 0,
			expectedMaxScore: 40,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Insert test events
			for _, event := range tt.setupEvents {
				if err := eventsRepo.RecordEvent(ctx, event); err != nil {
					t.Fatalf("failed to insert test event: %v", err)
				}
			}

			// Calculate reputation score
			score, err := service.CalculateReputationScore(ctx, tt.domain)
			if err != nil {
				t.Errorf("CalculateReputationScore() error = %v", err)
				return
			}

			if score.ReputationScore < tt.expectedMinScore || score.ReputationScore > tt.expectedMaxScore {
				t.Errorf("CalculateReputationScore() score = %v, want between %v and %v",
					score.ReputationScore, tt.expectedMinScore, tt.expectedMaxScore)
			}

			if tt.expectedBounceMin > 0 || tt.expectedBounceMax > 0 {
				if score.BounceRate < tt.expectedBounceMin || score.BounceRate > tt.expectedBounceMax {
					t.Errorf("CalculateReputationScore() bounce rate = %v, want between %v and %v",
						score.BounceRate, tt.expectedBounceMin, tt.expectedBounceMax)
				}
			}
		})
	}
}

func TestTelemetryService_CleanupOldData(t *testing.T) {
	db := setupTelemetryTestDB(t)
	defer db.Close()

	eventsRepo := sqlite.NewEventsRepository(db)
	scoresRepo := sqlite.NewScoresRepository(db)
	logger := zap.NewNop()
	service := NewTelemetryService(eventsRepo, scoresRepo, logger)

	ctx := context.Background()
	now := time.Now().Unix()

	// Insert old and recent events
	events := []*domain.SendingEvent{
		{
			Timestamp:       now - (100 * 86400), // 100 days ago
			Domain:          "example.com",
			RecipientDomain: "r1.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "1.1.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 3600, // 1 hour ago
			Domain:          "example.com",
			RecipientDomain: "r2.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "1.1.1.1",
			Metadata:        make(map[string]interface{}),
		},
	}

	for _, event := range events {
		if err := eventsRepo.RecordEvent(ctx, event); err != nil {
			t.Fatalf("failed to insert test event: %v", err)
		}
	}

	// Run cleanup
	if err := service.CleanupOldData(ctx); err != nil {
		t.Errorf("CleanupOldData() error = %v", err)
	}

	// Verify only recent events remain
	allEvents, err := eventsRepo.GetEventsInWindow(ctx, "example.com", 0, now)
	if err != nil {
		t.Fatalf("failed to get events: %v", err)
	}

	if len(allEvents) != 1 {
		t.Errorf("CleanupOldData() should leave 1 event, got %d", len(allEvents))
	}
}
