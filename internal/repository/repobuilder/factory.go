package repobuilder

import (
	"github.com/btafoya/gomailserver/internal/database"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/repository/postgres"
	"github.com/btafoya/gomailserver/internal/repository/sqlite"
)

// NewRepositories creates a Repositories struct based on the database driver type
func NewRepositories(db *database.DB) *repository.Repositories {
	switch db.DriverType() {
	case database.DriverPostgres:
		return newPostgresRepositories(db)
	case database.DriverSQLite:
		return newSQLiteRepositories(db)
	default:
		panic("unsupported database driver")
	}
}

func newSQLiteRepositories(db *database.DB) *repository.Repositories {
	return &repository.Repositories{
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

func newPostgresRepositories(db *database.DB) *repository.Repositories {
	return &repository.Repositories{
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
