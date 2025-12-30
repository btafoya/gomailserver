package repository

import (
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
)

// UserRepository defines user data access interface
type UserRepository interface {
	Create(user *domain.User) error
	GetByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	UpdateLastLogin(id int64) error
	UpdatePassword(userID int64, passwordHash string) error
	Delete(id int64) error
	List(domainID int64, offset, limit int) ([]*domain.User, error)
}

// MessageRepository defines message data access interface
type MessageRepository interface {
	Create(message *domain.Message) error
	GetByID(id int64) (*domain.Message, error)
	GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error)
	Update(message *domain.Message) error
	Delete(id int64) error
}

// MailboxRepository defines mailbox data access interface
type MailboxRepository interface {
	Create(mailbox *domain.Mailbox) error
	GetByID(id int64) (*domain.Mailbox, error)
	GetByUser(userID int64) ([]*domain.Mailbox, error)
	GetByName(userID int64, name string) (*domain.Mailbox, error)
	Update(mailbox *domain.Mailbox) error
	Delete(id int64) error
}

// DomainRepository defines domain data access interface
type DomainRepository interface {
	Create(domain *domain.Domain) error
	GetByID(id int64) (*domain.Domain, error)
	GetByName(name string) (*domain.Domain, error)
	Update(domain *domain.Domain) error
	Delete(id int64) error
	List(offset, limit int) ([]*domain.Domain, error)
}

// QueueRepository defines queue data access interface
type QueueRepository interface {
	Enqueue(item *domain.QueueItem) error
	GetPending() ([]*domain.QueueItem, error)
	GetByID(id int64) (*domain.QueueItem, error)
	UpdateStatus(id int64, status string, errorMsg string) error
	UpdateRetry(id int64, retryCount int, nextRetry time.Time) error
	Delete(id int64) error
}
