package sqlite

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/btafoya/gomailserver/internal/reputation/domain"
	_ "github.com/mattn/go-sqlite3"
)

func setupScoresTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Create schema
	schema := `
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

func TestScoresRepository_GetReputationScore(t *testing.T) {
	db := setupScoresTestDB(t)
	defer db.Close()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	// Insert test score
	score := &domain.ReputationScore{
		Domain:          "example.com",
		ReputationScore: 85,
		ComplaintRate:   0.1,
		BounceRate:      2.5,
		DeliveryRate:    97.4,
		LastUpdated:     time.Now().Unix(),
	}

	if err := repo.UpdateReputationScore(ctx, score); err != nil {
		t.Fatalf("failed to insert test score: %v", err)
	}

	tests := []struct {
		name    string
		domain  string
		wantErr bool
	}{
		{
			name:    "get existing score",
			domain:  "example.com",
			wantErr: false,
		},
		{
			name:    "get non-existent score",
			domain:  "nonexistent.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetReputationScore(ctx, tt.domain)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReputationScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Domain != tt.domain {
					t.Errorf("GetReputationScore() domain = %v, want %v", got.Domain, tt.domain)
				}
				if got.ReputationScore != 85 {
					t.Errorf("GetReputationScore() score = %v, want 85", got.ReputationScore)
				}
			}
		})
	}
}

func TestScoresRepository_UpdateReputationScore(t *testing.T) {
	db := setupScoresTestDB(t)
	defer db.Close()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	tests := []struct {
		name    string
		score   *domain.ReputationScore
		wantErr bool
	}{
		{
			name: "insert new score",
			score: &domain.ReputationScore{
				Domain:          "example.com",
				ReputationScore: 90,
				ComplaintRate:   0.05,
				BounceRate:      1.5,
				DeliveryRate:    98.45,
				LastUpdated:     time.Now().Unix(),
			},
			wantErr: false,
		},
		{
			name: "update existing score",
			score: &domain.ReputationScore{
				Domain:          "example.com",
				ReputationScore: 95,
				ComplaintRate:   0.02,
				BounceRate:      1.0,
				DeliveryRate:    98.98,
				LastUpdated:     time.Now().Unix(),
			},
			wantErr: false,
		},
		{
			name: "insert score with circuit breaker",
			score: &domain.ReputationScore{
				Domain:               "blocked.com",
				ReputationScore:      30,
				ComplaintRate:        0.5,
				BounceRate:           10.0,
				DeliveryRate:         89.5,
				CircuitBreakerActive: true,
				CircuitBreakerReason: "High complaint rate",
				LastUpdated:          time.Now().Unix(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdateReputationScore(ctx, tt.score)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateReputationScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the score was updated
				got, err := repo.GetReputationScore(ctx, tt.score.Domain)
				if err != nil {
					t.Errorf("failed to verify score: %v", err)
					return
				}

				if got.ReputationScore != tt.score.ReputationScore {
					t.Errorf("UpdateReputationScore() score = %v, want %v", got.ReputationScore, tt.score.ReputationScore)
				}
			}
		})
	}
}

func TestScoresRepository_ListAllScores(t *testing.T) {
	db := setupScoresTestDB(t)
	defer db.Close()

	repo := NewScoresRepository(db)
	ctx := context.Background()

	// Insert test scores
	scores := []*domain.ReputationScore{
		{
			Domain:          "example1.com",
			ReputationScore: 90,
			ComplaintRate:   0.05,
			BounceRate:      1.5,
			DeliveryRate:    98.45,
			LastUpdated:     time.Now().Unix(),
		},
		{
			Domain:          "example2.com",
			ReputationScore: 85,
			ComplaintRate:   0.1,
			BounceRate:      2.0,
			DeliveryRate:    97.9,
			LastUpdated:     time.Now().Unix(),
		},
		{
			Domain:          "example3.com",
			ReputationScore: 95,
			ComplaintRate:   0.02,
			BounceRate:      1.0,
			DeliveryRate:    98.98,
			LastUpdated:     time.Now().Unix(),
		},
	}

	for _, score := range scores {
		if err := repo.UpdateReputationScore(ctx, score); err != nil {
			t.Fatalf("failed to insert test score: %v", err)
		}
	}

	t.Run("list all scores", func(t *testing.T) {
		got, err := repo.ListAllScores(ctx)
		if err != nil {
			t.Errorf("ListAllScores() error = %v", err)
			return
		}

		if len(got) != 3 {
			t.Errorf("ListAllScores() got %d scores, want 3", len(got))
		}

		// Verify scores are sorted by domain
		if got[0].Domain > got[1].Domain || got[1].Domain > got[2].Domain {
			t.Error("ListAllScores() scores are not sorted by domain")
		}
	})
}
