package postgres

import (
	"database/sql"
)

type ipBlacklistRepository struct {
	db *database.DB
}

func NewIPBlacklistRepository(db *database.DB) repository.IPBlacklistRepository {
	return &ipBlacklistRepository{db: db}
}

func (r *ipBlacklistRepository) Create(entry *domain.IPBlacklist) error {
	panic("postgres repository not implemented yet")
}

func (r *ipBlacklistRepository) GetByIP(ipAddress string) (*domain.IPBlacklist, error) {
	panic("postgres repository not implemented yet")
}

func (r *ipBlacklistRepository) List(offset, limit int) ([]*domain.IPBlacklist, error) {
	panic("postgres repository not implemented yet")
}

func (r *ipBlacklistRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *ipBlacklistRepository) Cleanup() error {
	panic("postgres repository not implemented yet")
}
