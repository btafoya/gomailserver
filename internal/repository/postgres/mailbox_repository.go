package postgres

import (
	"database/sql"
)

type mailBoxRepository struct {
	db *database.DB
}

func NewMailboxRepository(db *database.DB) repository.MailboxRepository {
	return &mailBoxRepository{db: db}
}

func (r *mailBoxRepository) Create(mailbox *domain.Mailbox) error {
	panic("postgres repository not implemented yet")
}

func (r *mailBoxRepository) GetByID(id int64) (*domain.Mailbox, error) {
	panic("postgres repository not implemented yet")
}

func (r *mailBoxRepository) GetByUserID(userID int64) ([]*domain.Mailbox, error) {
	panic("postgres repository not implemented yet")
}

func (r *mailBoxRepository) Update(mailbox *domain.Mailbox) error {
	panic("postgres repository not implemented yet")
}

func (r *mailBoxRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *mailBoxRepository) List(userID int64, offset, limit int) ([]*domain.Mailbox, error) {
	panic("postgres repository not implemented yet")
}
