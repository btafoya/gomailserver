package service

import (
	"time"

	"github.com/btafoya/gomailserver/internal/domain"
)

// UserServiceInterface defines the user service interface
type UserServiceInterface interface {
	Create(user *domain.User, password string) error
	Authenticate(email, password string) (*domain.User, error)
	GetByID(id int64) (*domain.User, error)
	GetByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	UpdatePassword(userID int64, newPassword string) error
	Delete(id int64) error
}

// MessageServiceInterface defines the message service interface
type MessageServiceInterface interface {
	Store(userID, mailboxID, uid int64, messageData []byte) (*domain.Message, error)
	GetByID(id int64) (*domain.Message, error)
	GetByMailbox(mailboxID int64, offset, limit int) ([]*domain.Message, error)
	Delete(id int64) error
}

// MailboxServiceInterface defines the mailbox service interface
type MailboxServiceInterface interface {
	Create(userID int64, name, specialUse string) error
	GetByName(userID int64, name string) (*domain.Mailbox, error)
	List(userID int64, subscribedOnly bool) ([]*domain.Mailbox, error)
	Delete(mailboxID int64) error
	Rename(mailboxID int64, newName string) error
	UpdateSubscription(id int64, subscribed bool) error
}

// QueueServiceInterface defines the queue service interface
type QueueServiceInterface interface {
	Enqueue(from string, to []string, message []byte) (string, error)
	GetPending() ([]*domain.QueueItem, error)
	MarkDelivered(id int64) error
	MarkFailed(id int64, errorMsg string) error
	IncrementRetry(id int64, currentRetryCount int, failedAt time.Time) error
	CalculateNextRetry(retryCount int, failedAt time.Time) time.Time
}
