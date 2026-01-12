package postgres

import (
	"database/sql"
)

type loginAttemptRepository struct {
	db *database.DB
}

func NewLoginAttemptRepository(db *database.DB) repository.LoginAttemptRepository {
	return &loginAttemptRepository{db: db}
}

func (r *loginAttemptRepository) Create(attempt *domain.LoginAttempt) error {
	panic("postgres repository not implemented yet")
}

func (r *loginAttemptRepository) GetByIP(ipAddress string) ([]*domain.LoginAttempt, error) {
	panic("postgres repository not implemented yet")
}

func (r *loginAttemptRepository) List(offset, limit int) ([]*domain.LoginAttempt, error) {
	panic("postgres repository not implemented yet")
}

func (r *loginAttemptRepository) Cleanup(timestamp time.Time) error {
	panic("postgres repository not implemented yet")
}
