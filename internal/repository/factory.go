package repository

import (
	"github.com/btafoya/gomailserver/internal/database"
)

type Repositories struct {
	User         UserRepository
	Domain       DomainRepository
	Message      MessageRepository
	Mailbox      MailboxRepository
	Alias        AliasRepository
	Queue        QueueRepository
	LoginAttempt LoginAttemptRepository
	IPBlacklist  IPBlacklistRepository
	Greylist     GreylistRepository
	RateLimit    RateLimitRepository
	Webhook      WebhookRepository
}

func NewRepositories(db *database.DB) *Repositories {
	switch db.DriverType() {
	case database.DriverPostgres:
		return newPostgresRepositories(db)
	case database.DriverSQLite:
		return newSQLiteRepositories(db)
	default:
		panic("unsupported database driver")
	}
}

func newSQLiteRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:         sqlite.NewUserRepository(db),
		Domain:       sqlite.NewDomainRepository(db),
		Message:      sqlite.NewMessageRepository(db),
		Mailbox:      sqlite.NewMailboxRepository(db),
		Alias:        sqlite.NewAliasRepository(db),
		Queue:        sqlite.NewQueueRepository(db),
		LoginAttempt: sqlite.NewLoginAttemptRepository(db),
		IPBlacklist:  sqlite.NewIPBlacklistRepository(db),
		Greylist:     sqlite.NewGreylistRepository(db),
		RateLimit:    sqlite.NewRateLimitRepository(db),
		Webhook:      sqlite.NewWebhookRepository(db),
	}
}

func newPostgresRepositories(db *database.DB) *Repositories {
	return &Repositories{
		User:         postgres.NewUserRepository(db),
		Domain:       postgres.NewDomainRepository(db),
		Message:      postgres.NewMessageRepository(db),
		Mailbox:      postgres.NewMailboxRepository(db),
		Alias:        postgres.NewAliasRepository(db),
		Queue:        postgres.NewQueueRepository(db),
		LoginAttempt: postgres.NewLoginAttemptRepository(db),
		IPBlacklist:  postgres.NewIPBlacklistRepository(db),
		Greylist:     postgres.NewGreylistRepository(db),
		RateLimit:    postgres.NewRateLimitRepository(db),
		Webhook:      postgres.NewWebhookRepository(db),
	}
}
