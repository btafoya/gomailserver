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
	ListAll() ([]*domain.User, error)
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

// AliasRepository defines alias data access interface
type AliasRepository interface {
	Create(alias *domain.Alias) error
	GetByID(id int64) (*domain.Alias, error)
	GetByEmail(email string) (*domain.Alias, error)
	Update(alias *domain.Alias) error
	Delete(id int64) error
	ListAll() ([]*domain.Alias, error)
	ListByDomain(domainID int64) ([]*domain.Alias, error)
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

// GreylistRepository defines greylist data access interface
type GreylistRepository interface {
	Create(ip, sender, recipient string) (*domain.GreylistTriplet, error)
	Get(ip, sender, recipient string) (*domain.GreylistTriplet, error)
	IncrementPass(id int64) error
	DeleteOlderThan(age time.Duration) error
}

// RateLimitRepository defines rate limit data access interface
type RateLimitRepository interface {
	Get(key string, limitType string) (*domain.RateLimitEntry, error)
	CreateOrUpdate(entry *domain.RateLimitEntry) error
	Cleanup(windowDuration time.Duration) error
}

// LoginAttemptRepository defines login attempt tracking interface
type LoginAttemptRepository interface {
	Record(ip, email string, success bool) error
	GetRecentFailures(ip string, duration time.Duration) (int, error)
	GetRecentUserFailures(email string, duration time.Duration) (int, error)
	Cleanup(age time.Duration) error
}

// IPBlacklistRepository defines IP blacklist interface
type IPBlacklistRepository interface {
	Add(ip, reason string, expiresAt *time.Time) error
	IsBlacklisted(ip string) (bool, error)
	Remove(ip string) error
	RemoveExpired() error
}

// QuarantineRepository defines quarantine data access interface
type QuarantineRepository interface {
	Create(message *domain.QuarantineMessage) error
	GetByID(id int64) (*domain.QuarantineMessage, error)
	List(offset, limit int) ([]*domain.QuarantineMessage, error)
	UpdateAction(id int64, action string) error
	Delete(id int64) error
	DeleteOlderThan(age time.Duration) error
}

// APIKeyRepository defines API key data access interface
type APIKeyRepository interface {
	Create(apiKey *domain.APIKey) error
	GetByKeyHash(keyHash string) (*domain.APIKey, error)
	GetByID(id int64) (*domain.APIKey, error)
	ListByUser(userID int64) ([]*domain.APIKey, error)
	UpdateLastUsed(id int64, ip string) error
	Delete(id int64) error
}
