package postgres

import (
	"database/sql"
)

type rateLimitRepository struct {
	db *database.DB
}

func NewRateLimitRepository(db *database.DB) repository.RateLimitRepository {
	return &rateLimitRepository{db: db}
}

func (r *rateLimitRepository) Create(entry *domain.RateLimit) error {
	panic("postgres repository not implemented yet")
}

func (r *rateLimitRepository) Get(entityType, entityValue, actionType string) (*domain.RateLimit, error) {
	panic("postgres repository not implemented yet")
}

func (r *rateLimitRepository) Increment(entityType, entityValue, actionType string) error {
	panic("postgres repository not implemented yet")
}

func (r *rateLimitRepository) IsExceeded(entityType, entityValue, actionType string) (bool, error) {
	panic("postgres repository not implemented yet")
}

func (r *rateLimitRepository) ResetWindowStart(entityType, entityValue, actionType string) error {
	panic("postgres repository not implemented yet")
}

func (r *rateLimitRepository) Cleanup() error {
	panic("postgres repository not implemented yet")
}
