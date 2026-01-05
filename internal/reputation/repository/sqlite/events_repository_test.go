package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	_ "github.com/mattn/go-sqlite3"
)

func setupTestDB(t *testing.T) *sql.DB {
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
	CREATE INDEX idx_sending_events_event_type ON sending_events(event_type);
	`

	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func TestEventsRepository_RecordEvent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEventsRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		event   *domain.SendingEvent
		wantErr bool
	}{
		{
			name: "record delivery event",
			event: &domain.SendingEvent{
				Timestamp:       time.Now().Unix(),
				Domain:          "example.com",
				RecipientDomain: "recipient.com",
				EventType:       domain.EventDelivery,
				IPAddress:       "192.168.1.1",
				Metadata:        map[string]interface{}{"test": "value"},
			},
			wantErr: false,
		},
		{
			name: "record bounce event with details",
			event: &domain.SendingEvent{
				Timestamp:          time.Now().Unix(),
				Domain:             "example.com",
				RecipientDomain:    "recipient.com",
				EventType:          domain.EventBounce,
				BounceType:         stringPtr("hard"),
				EnhancedStatusCode: stringPtr("5.1.1"),
				SMTPResponse:       stringPtr("User unknown"),
				IPAddress:          "192.168.1.1",
				Metadata:           map[string]interface{}{"retry": 1},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.RecordEvent(ctx, tt.event)
			if (err != nil) != tt.wantErr {
				t.Errorf("RecordEvent() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.event.ID == 0 {
				t.Error("RecordEvent() did not set event ID")
			}
		})
	}
}

func TestEventsRepository_GetEventsInWindow(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEventsRepository(db)
	ctx := context.Background()

	// Insert test events
	now := time.Now().Unix()
	events := []*domain.SendingEvent{
		{
			Timestamp:       now - 3600, // 1 hour ago
			Domain:          "example.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 7200, // 2 hours ago
			Domain:          "example.com",
			RecipientDomain: "recipient2.com",
			EventType:       domain.EventBounce,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 10800, // 3 hours ago
			Domain:          "other.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.2",
			Metadata:        make(map[string]interface{}),
		},
	}

	for _, event := range events {
		if err := repo.RecordEvent(ctx, event); err != nil {
			t.Fatalf("failed to insert test event: %v", err)
		}
	}

	tests := []struct {
		name      string
		domain    string
		startTime int64
		endTime   int64
		wantCount int
	}{
		{
			name:      "get all events for example.com",
			domain:    "example.com",
			startTime: now - 86400, // 24 hours ago
			endTime:   now,
			wantCount: 2,
		},
		{
			name:      "get events in narrow window",
			domain:    "example.com",
			startTime: now - 4000,
			endTime:   now - 3000,
			wantCount: 1,
		},
		{
			name:      "no events for non-existent domain",
			domain:    "nonexistent.com",
			startTime: now - 86400,
			endTime:   now,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetEventsInWindow(ctx, tt.domain, tt.startTime, tt.endTime)
			if err != nil {
				t.Errorf("GetEventsInWindow() error = %v", err)
				return
			}

			if len(got) != tt.wantCount {
				t.Errorf("GetEventsInWindow() got %d events, want %d", len(got), tt.wantCount)
			}
		})
	}
}

func TestEventsRepository_GetEventCountsByType(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEventsRepository(db)
	ctx := context.Background()

	// Insert test events
	now := time.Now().Unix()
	events := []*domain.SendingEvent{
		{
			Timestamp:       now - 3600,
			Domain:          "example.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 3600,
			Domain:          "example.com",
			RecipientDomain: "recipient2.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 3600,
			Domain:          "example.com",
			RecipientDomain: "recipient3.com",
			EventType:       domain.EventBounce,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
	}

	for _, event := range events {
		if err := repo.RecordEvent(ctx, event); err != nil {
			t.Fatalf("failed to insert test event: %v", err)
		}
	}

	tests := []struct {
		name              string
		domain            string
		startTime         int64
		endTime           int64
		wantDeliveryCount int64
		wantBounceCount   int64
	}{
		{
			name:              "count events for example.com",
			domain:            "example.com",
			startTime:         now - 86400,
			endTime:           now,
			wantDeliveryCount: 2,
			wantBounceCount:   1,
		},
		{
			name:              "no events in window",
			domain:            "example.com",
			startTime:         now - 86400,
			endTime:           now - 7200,
			wantDeliveryCount: 0,
			wantBounceCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetEventCountsByType(ctx, tt.domain, tt.startTime, tt.endTime)
			if err != nil {
				t.Errorf("GetEventCountsByType() error = %v", err)
				return
			}

			if got[string(domain.EventDelivery)] != tt.wantDeliveryCount {
				t.Errorf("GetEventCountsByType() delivery count = %d, want %d", got[string(domain.EventDelivery)], tt.wantDeliveryCount)
			}

			if got[string(domain.EventBounce)] != tt.wantBounceCount {
				t.Errorf("GetEventCountsByType() bounce count = %d, want %d", got[string(domain.EventBounce)], tt.wantBounceCount)
			}
		})
	}
}

func TestEventsRepository_CleanupOldEvents(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewEventsRepository(db)
	ctx := context.Background()

	// Insert test events
	now := time.Now().Unix()
	events := []*domain.SendingEvent{
		{
			Timestamp:       now - (100 * 86400), // 100 days ago
			Domain:          "example.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - (50 * 86400), // 50 days ago
			Domain:          "example.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
		{
			Timestamp:       now - 3600, // 1 hour ago
			Domain:          "example.com",
			RecipientDomain: "recipient.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		},
	}

	for _, event := range events {
		if err := repo.RecordEvent(ctx, event); err != nil {
			t.Fatalf("failed to insert test event: %v", err)
		}
	}

	// Cleanup events older than 90 days
	cutoff := now - (90 * 86400)
	if err := repo.CleanupOldEvents(ctx, cutoff); err != nil {
		t.Errorf("CleanupOldEvents() error = %v", err)
	}

	// Verify only recent events remain
	allEvents, err := repo.GetEventsInWindow(ctx, "example.com", 0, now)
	if err != nil {
		t.Fatalf("failed to get events: %v", err)
	}

	if len(allEvents) != 2 {
		t.Errorf("CleanupOldEvents() should leave 2 events, got %d", len(allEvents))
	}
}

func stringPtr(s string) *string {
	return &s
}
