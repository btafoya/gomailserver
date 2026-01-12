package postgres

import (
	"database/sql"
)

type queueRepository struct {
	db *database.DB
}

func NewQueueRepository(db *database.DB) repository.QueueRepository {
	return &queueRepository{db: db}
}

func (r *queueRepository) Create(queue *domain.Queue) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) GetByID(id int64) (*domain.Queue, error) {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) List(offset, limit int) ([]*domain.Queue, error) {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) Update(queue *domain.Queue) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) Delete(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) GetPending(limit int) ([]*domain.Queue, error) {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) UpdateStatus(id int64, status string) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) IncrementRetry(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) MarkDelivered(id int64) error {
	panic("postgres repository not implemented yet")
}

func (r *queueRepository) MarkFailed(id int64, errorMessage string) error {
	panic("postgres repository not implemented yet")
}
