package postgres

import (
	"database/sql"
)

type aliasRepository struct {
	db *database.DB
}

func NewAliasRepository(db *database.DB) repository.AliasRepository {
	return &aliasRepository{db: db}
}

func (r *aliasRepository) Create(alias *domain.Alias) error {
	panic("postgres repository not implemented yet")
}

func (r *aliasRepository) GetByID(id int64) (*domain.Alias, error) {
	panic("postgres repository not implemented yet")
}

func (r *aliasRepository) GetByDomainID(domainID int64) ([]*domain.Alias, error) {
	panic("postgres repository not implemented yet")
}

func (r *aliasRepository) Update(alias *domain.Alias) error {
	panic("postgres repository not implemented yet")
}

func (r *aliasRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *aliasRepository) List(domainID int64, offset, limit int) ([]*domain.Alias, error) {
	panic("postgres repository not implemented yet")
}
