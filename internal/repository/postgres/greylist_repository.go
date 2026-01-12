package postgres

import (
	"database/sql"
)

type greylistRepository struct {
	db *database.DB
}

func NewGreylistRepository(db *database.DB) repository.GreylistRepository {
	return &greylistRepository{db: db}
}

func (r *greylistRepository) Create(entry *domain.Greylist) error {
	panic("postgres repository not implemented yet")
}

func (r *greylistRepository) GetTriplet(senderIP, senderEmail, recipientEmail string) (*domain.Greylist, error) {
	panic("postgres repository not implemented yet")
}

func (r *greylistRepository) Update(entry *domain.Greylist) error {
	panic("postgres repository not implemented yet")
}

func (r *greylistRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *greylistRepository) Cleanup(expiredAt time.Time) error {
	panic("postgres repository not implemented yet")
}
