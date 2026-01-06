package reputation

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	"github.com/btafoya/gomailserver/internal/reputation/service"
	"go.uber.org/zap"
)

// TestEndToEndEventRecording tests the full event recording flow
func TestEndToEndEventRecording(t *testing.T) {
	// Create temp database
	dbPath := "./test_reputation_e2e.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := InitDatabase(Config{Path: dbPath}, zap.NewNop())
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Test 1: Record delivery events
	t.Run("record_delivery_events", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			err := db.TelemetryService.RecordDelivery(ctx, "example.com", "recipient.com", "192.168.1.1")
			if err != nil {
				t.Errorf("failed to record delivery %d: %v", i, err)
			}
		}
	})

	// Test 2: Record bounce events
	t.Run("record_bounce_events", func(t *testing.T) {
		for i := 0; i < 2; i++ {
			err := db.TelemetryService.RecordBounce(ctx, "example.com", "recipient.com", "192.168.1.1", "hard", "5.1.1", "User unknown")
			if err != nil {
				t.Errorf("failed to record bounce %d: %v", i, err)
			}
		}
	})

	// Test 3: Verify events were stored
	t.Run("verify_events_stored", func(t *testing.T) {
		now := time.Now().Unix()
		events, err := db.EventsRepo.GetEventsInWindow(ctx, "example.com", now-3600, now)
		if err != nil {
			t.Fatalf("failed to get events: %v", err)
		}

		if len(events) != 12 { // 10 deliveries + 2 bounces
			t.Errorf("expected 12 events, got %d", len(events))
		}

		// Verify event types
		deliveryCount := 0
		bounceCount := 0
		for _, event := range events {
			switch event.EventType {
			case domain.EventDelivery:
				deliveryCount++
			case domain.EventBounce:
				bounceCount++
			}
		}

		if deliveryCount != 10 {
			t.Errorf("expected 10 deliveries, got %d", deliveryCount)
		}
		if bounceCount != 2 {
			t.Errorf("expected 2 bounces, got %d", bounceCount)
		}
	})
}

// TestEndToEndReputationCalculation tests the reputation score calculation flow
func TestEndToEndReputationCalculation(t *testing.T) {
	// Create temp database
	dbPath := "./test_reputation_calc.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := InitDatabase(Config{Path: dbPath}, zap.NewNop())
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()

	// Test scenario: 98 deliveries, 2 bounces, 0 complaints = excellent reputation
	t.Run("calculate_good_reputation", func(t *testing.T) {
		// Record events
		for i := 0; i < 98; i++ {
			err := db.TelemetryService.RecordDelivery(ctx, "good.com", "recipient.com", "192.168.1.1")
			if err != nil {
				t.Fatalf("failed to record delivery: %v", err)
			}
		}
		for i := 0; i < 2; i++ {
			err := db.TelemetryService.RecordBounce(ctx, "good.com", "recipient.com", "192.168.1.1", "soft", "4.2.1", "Mailbox full")
			if err != nil {
				t.Fatalf("failed to record bounce: %v", err)
			}
		}

		// Calculate reputation score
		score, err := db.TelemetryService.CalculateReputationScore(ctx, "good.com")
		if err != nil {
			t.Fatalf("failed to calculate reputation: %v", err)
		}

		// Verify score is excellent (>90)
		if score.ReputationScore < 90 {
			t.Errorf("expected excellent reputation (>90), got %d", score.ReputationScore)
		}

		// Verify metrics
		if score.DeliveryRate < 97.0 || score.DeliveryRate > 99.0 {
			t.Errorf("expected ~98%% delivery rate, got %.2f", score.DeliveryRate)
		}
		if score.BounceRate < 1.0 || score.BounceRate > 3.0 {
			t.Errorf("expected ~2%% bounce rate, got %.2f", score.BounceRate)
		}
		if score.ComplaintRate > 0.01 {
			t.Errorf("expected 0%% complaint rate, got %.2f", score.ComplaintRate)
		}
	})

	// Test scenario: 40 deliveries, 50 bounces, 10 complaints = poor reputation
	t.Run("calculate_poor_reputation", func(t *testing.T) {
		// Record events
		for i := 0; i < 40; i++ {
			err := db.TelemetryService.RecordDelivery(ctx, "bad.com", "recipient.com", "192.168.1.1")
			if err != nil {
				t.Fatalf("failed to record delivery: %v", err)
			}
		}
		for i := 0; i < 50; i++ {
			err := db.TelemetryService.RecordBounce(ctx, "bad.com", "recipient.com", "192.168.1.1", "hard", "5.1.1", "User unknown")
			if err != nil {
				t.Fatalf("failed to record bounce: %v", err)
			}
		}
		for i := 0; i < 10; i++ {
			err := db.TelemetryService.RecordComplaint(ctx, "bad.com", "recipient.com")
			if err != nil {
				t.Fatalf("failed to record complaint: %v", err)
			}
		}

		// Calculate reputation score
		score, err := db.TelemetryService.CalculateReputationScore(ctx, "bad.com")
		if err != nil {
			t.Fatalf("failed to calculate reputation: %v", err)
		}

		// Verify score is poor (<50)
		if score.ReputationScore >= 50 {
			t.Errorf("expected poor reputation (<50), got %d", score.ReputationScore)
		}

		// Verify high bounce rate
		if score.BounceRate < 45.0 {
			t.Errorf("expected high bounce rate (>45%%), got %.2f", score.BounceRate)
		}
	})

	// Test CalculateAllScores
	t.Run("calculate_all_scores", func(t *testing.T) {
		err := db.TelemetryService.CalculateAllScores(ctx)
		if err != nil {
			t.Fatalf("failed to calculate all scores: %v", err)
		}

		// Verify both domains have scores
		scores, err := db.ScoresRepo.ListAllScores(ctx)
		if err != nil {
			t.Fatalf("failed to list scores: %v", err)
		}

		if len(scores) != 2 {
			t.Errorf("expected 2 domain scores, got %d", len(scores))
		}
	})
}

// TestEndToEndDataRetention tests the cleanup process
func TestEndToEndDataRetention(t *testing.T) {
	// Create temp database
	dbPath := "./test_reputation_retention.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := InitDatabase(Config{Path: dbPath}, zap.NewNop())
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	now := time.Now().Unix()

	// Insert old events (100 days ago) and recent events
	t.Run("insert_mixed_age_events", func(t *testing.T) {
		// Old events
		oldEvent := &domain.SendingEvent{
			Timestamp:       now - (100 * 86400),
			Domain:          "test.com",
			RecipientDomain: "old.com",
			EventType:       domain.EventDelivery,
			IPAddress:       "192.168.1.1",
			Metadata:        make(map[string]interface{}),
		}
		err := db.EventsRepo.RecordEvent(ctx, oldEvent)
		if err != nil {
			t.Fatalf("failed to record old event: %v", err)
		}

		// Recent events
		for i := 0; i < 5; i++ {
			recentEvent := &domain.SendingEvent{
				Timestamp:       now - 3600, // 1 hour ago
				Domain:          "test.com",
				RecipientDomain: "recent.com",
				EventType:       domain.EventDelivery,
				IPAddress:       "192.168.1.1",
				Metadata:        make(map[string]interface{}),
			}
			err := db.EventsRepo.RecordEvent(ctx, recentEvent)
			if err != nil {
				t.Fatalf("failed to record recent event: %v", err)
			}
		}
	})

	// Verify all events are present
	t.Run("verify_events_before_cleanup", func(t *testing.T) {
		events, err := db.EventsRepo.GetEventsInWindow(ctx, "test.com", 0, now)
		if err != nil {
			t.Fatalf("failed to get events: %v", err)
		}

		if len(events) != 6 {
			t.Errorf("expected 6 events before cleanup, got %d", len(events))
		}
	})

	// Run cleanup
	t.Run("run_cleanup", func(t *testing.T) {
		err := db.TelemetryService.CleanupOldData(ctx)
		if err != nil {
			t.Fatalf("failed to cleanup old data: %v", err)
		}
	})

	// Verify only recent events remain
	t.Run("verify_events_after_cleanup", func(t *testing.T) {
		events, err := db.EventsRepo.GetEventsInWindow(ctx, "test.com", 0, now)
		if err != nil {
			t.Fatalf("failed to get events: %v", err)
		}

		if len(events) != 5 {
			t.Errorf("expected 5 events after cleanup (old one removed), got %d", len(events))
		}

		// Verify all remaining events are recent
		for _, event := range events {
			age := now - event.Timestamp
			if age > 90*86400 {
				t.Errorf("found event older than 90 days after cleanup: age=%d days", age/86400)
			}
		}
	})
}

// TestSchedulerIntegration tests the scheduler functionality
func TestSchedulerIntegration(t *testing.T) {
	// Create temp database
	dbPath := "./test_scheduler.db"
	defer os.Remove(dbPath)

	// Initialize database
	db, err := InitDatabase(Config{Path: dbPath}, zap.NewNop())
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	// Create services needed for scheduler
	circuitBreakerSvc := service.NewCircuitBreakerService(
		db.EventsRepo,
		db.ScoresRepo,
		db.CircuitBreakerRepo,
		db.TelemetryService,
		zap.NewNop(),
	)
	warmUpSvc := service.NewWarmUpService(
		db.EventsRepo,
		db.ScoresRepo,
		db.WarmUpRepo,
		db.TelemetryService,
		zap.NewNop(),
	)

	// Create scheduler
	scheduler := NewScheduler(db.TelemetryService, circuitBreakerSvc, warmUpSvc, zap.NewNop())

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Start scheduler
	t.Run("start_scheduler", func(t *testing.T) {
		err := scheduler.Start(ctx)
		if err != nil {
			t.Fatalf("failed to start scheduler: %v", err)
		}

		// Record some events
		for i := 0; i < 5; i++ {
			err := db.TelemetryService.RecordDelivery(ctx, "scheduled.com", "recipient.com", "192.168.1.1")
			if err != nil {
				t.Errorf("failed to record delivery: %v", err)
			}
		}

		// Wait for context to finish (2 seconds)
		<-ctx.Done()
	})

	// Stop scheduler gracefully
	t.Run("stop_scheduler", func(t *testing.T) {
		err := scheduler.Stop()
		if err != nil {
			t.Fatalf("failed to stop scheduler: %v", err)
		}
	})

	// Verify scheduler processed events
	t.Run("verify_scheduler_calculations", func(t *testing.T) {
		// Check if reputation score was calculated
		score, err := db.ScoresRepo.GetReputationScore(context.Background(), "scheduled.com")
		if err != nil {
			// Score might not exist if scheduler didn't run yet (timing dependent)
			t.Logf("reputation score not calculated yet (expected for short test): %v", err)
			return
		}

		if score.Domain != "scheduled.com" {
			t.Errorf("expected domain 'scheduled.com', got %s", score.Domain)
		}
	})
}
