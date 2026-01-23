package repository

// Repositories is a container for all repository implementations
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
