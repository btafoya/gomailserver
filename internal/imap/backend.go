package imap

import (
	"errors"
	"net"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/repository"
	"github.com/btafoya/gomailserver/internal/security/bruteforce"
	"github.com/btafoya/gomailserver/internal/security/ratelimit"
	"github.com/btafoya/gomailserver/internal/security/totp"
	"github.com/btafoya/gomailserver/internal/service"
)

// Backend implements IMAP backend interface
type Backend struct {
	userService    service.UserServiceInterface
	mailboxService service.MailboxServiceInterface
	messageService service.MessageServiceInterface
	domainRepo     repository.DomainRepository
	logger         *zap.Logger

	// Security services
	rateLimiter *ratelimit.Limiter
	bruteForce  *bruteforce.Protection
	totpService *totp.TOTPService
}

// NewBackend creates a new IMAP backend with all dependencies
func NewBackend(
	userService service.UserServiceInterface,
	mailboxService service.MailboxServiceInterface,
	messageService service.MessageServiceInterface,
	domainRepo repository.DomainRepository,
	rateLimiter *ratelimit.Limiter,
	bruteForce *bruteforce.Protection,
	totpService *totp.TOTPService,
	logger *zap.Logger,
) *Backend {
	return &Backend{
		userService:    userService,
		mailboxService: mailboxService,
		messageService: messageService,
		domainRepo:     domainRepo,
		logger:         logger,
		rateLimiter:    rateLimiter,
		bruteForce:     bruteForce,
		totpService:    totpService,
	}
}

// Login authenticates a user
func (b *Backend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	b.logger.Info("IMAP authentication attempt",
		zap.String("username", username),
		zap.String("remote_addr", connInfo.RemoteAddr.String()),
	)

	// Extract domain from username
	domain := extractDomain(username)
	if domain == "" {
		return nil, backend.ErrInvalidCredentials
	}

	// Load domain configuration
	domainConfig, err := b.domainRepo.GetByName(domain)
	if err != nil {
		b.logger.Error("failed to load domain config",
			zap.String("domain", domain),
			zap.Error(err),
		)
		// Continue even if domain config fails
		domainConfig = nil
	}

	// Extract IP address
	remoteIP := extractIP(connInfo.RemoteAddr.String())

	// Check brute force protection if enabled
	if domainConfig != nil && b.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
		blocked, err := b.bruteForce.IsBlocked(remoteIP)
		if err != nil {
			b.logger.Error("brute force check failed", zap.Error(err))
		} else if blocked {
			b.logger.Warn("IMAP authentication blocked - brute force protection",
				zap.String("username", username),
				zap.String("remote_ip", remoteIP),
			)
			return nil, backend.ErrInvalidCredentials
		}
	}

	// Check IMAP rate limiting if enabled
	if domainConfig != nil && b.rateLimiter != nil && domainConfig.RateLimitEnabled {
		allowed, err := b.rateLimiter.Check("imap_per_user", username)
		if err != nil {
			b.logger.Error("rate limit check failed", zap.Error(err))
		} else if !allowed {
			b.logger.Warn("IMAP rate limited",
				zap.String("username", username),
				zap.String("remote_ip", remoteIP),
				zap.String("domain", domain),
			)
			return nil, backend.ErrInvalidCredentials
		}
	}

	user, err := b.userService.Authenticate(username, password)
	if err != nil {
		// Record failed login attempt for brute force protection
		if domainConfig != nil && b.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
			if err := b.bruteForce.RecordFailure(remoteIP, username); err != nil {
				b.logger.Error("failed to record login failure", zap.Error(err))
			}
		}

		b.logger.Warn("IMAP authentication failed",
			zap.String("username", username),
			zap.String("remote_addr", connInfo.RemoteAddr.String()),
			zap.Error(err),
		)
		return nil, backend.ErrInvalidCredentials
	}

	if user.Status != "active" {
		b.logger.Warn("IMAP authentication failed - user disabled",
			zap.String("username", username),
			zap.String("status", user.Status),
		)
		return nil, backend.ErrInvalidCredentials
	}

	// Check TOTP if enforced
	if domainConfig != nil && b.totpService != nil && domainConfig.AuthTOTPEnforced {
		if !user.TOTPEnabled {
			b.logger.Warn("IMAP authentication failed - TOTP required but not enabled",
				zap.String("username", username),
			)
			return nil, errors.New("TOTP required but not configured")
		}
		// Note: TOTP validation would require a separate authentication step
		// For now, we just check if TOTP is enabled
	}

	// Record successful login for brute force protection
	if domainConfig != nil && b.bruteForce != nil && domainConfig.AuthBruteForceEnabled {
		if err := b.bruteForce.RecordSuccess(remoteIP, username); err != nil {
			b.logger.Error("failed to record successful login", zap.Error(err))
		}
	}

	b.logger.Info("IMAP authentication successful",
		zap.String("username", username),
		zap.Int64("user_id", user.ID),
	)

	return &User{
		user:           user,
		backend:        b,
		mailboxService: b.mailboxService,
		messageService: b.messageService,
		logger:         b.logger,
	}, nil
}

// User implements IMAP user interface
type User struct {
	user           *domain.User
	backend        *Backend
	mailboxService service.MailboxServiceInterface
	messageService service.MessageServiceInterface
	logger         *zap.Logger
}

// Username returns the user's email
func (u *User) Username() string {
	return u.user.Email
}

// ListMailboxes lists user's mailboxes
func (u *User) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	u.logger.Debug("listing mailboxes",
		zap.Int64("user_id", u.user.ID),
		zap.Bool("subscribed", subscribed),
	)

	mailboxes, err := u.mailboxService.List(u.user.ID, subscribed)
	if err != nil {
		u.logger.Error("failed to list mailboxes",
			zap.Error(err),
			zap.Int64("user_id", u.user.ID),
		)
		return nil, err
	}

	result := make([]backend.Mailbox, len(mailboxes))
	for i, mb := range mailboxes {
		result[i] = &Mailbox{
			mailbox:        mb,
			user:           u.user,
			messageService: u.messageService,
			mailboxService: u.mailboxService,
			logger:         u.logger,
		}
	}

	return result, nil
}

// GetMailbox retrieves a specific mailbox
func (u *User) GetMailbox(name string) (backend.Mailbox, error) {
	u.logger.Debug("getting mailbox",
		zap.Int64("user_id", u.user.ID),
		zap.String("mailbox", name),
	)

	mb, err := u.mailboxService.GetByName(u.user.ID, name)
	if err != nil {
		u.logger.Error("failed to get mailbox",
			zap.Error(err),
			zap.Int64("user_id", u.user.ID),
			zap.String("mailbox", name),
		)
		return nil, backend.ErrNoSuchMailbox
	}

	return &Mailbox{
		mailbox:        mb,
		user:           u.user,
		messageService: u.messageService,
		mailboxService: u.mailboxService,
		logger:         u.logger,
	}, nil
}

// CreateMailbox creates a new mailbox
func (u *User) CreateMailbox(name string) error {
	u.logger.Info("creating mailbox",
		zap.Int64("user_id", u.user.ID),
		zap.String("mailbox", name),
	)

	err := u.mailboxService.Create(u.user.ID, name, "")
	if err != nil {
		u.logger.Error("failed to create mailbox",
			zap.Error(err),
			zap.Int64("user_id", u.user.ID),
			zap.String("mailbox", name),
		)
		return err
	}

	return nil
}

// DeleteMailbox deletes a mailbox
func (u *User) DeleteMailbox(name string) error {
	u.logger.Info("deleting mailbox",
		zap.Int64("user_id", u.user.ID),
		zap.String("mailbox", name),
	)

	// Prevent deletion of INBOX
	if name == "INBOX" {
		return errors.New("cannot delete INBOX")
	}

	mb, err := u.mailboxService.GetByName(u.user.ID, name)
	if err != nil {
		return backend.ErrNoSuchMailbox
	}

	err = u.mailboxService.Delete(mb.ID)
	if err != nil {
		u.logger.Error("failed to delete mailbox",
			zap.Error(err),
			zap.Int64("user_id", u.user.ID),
			zap.String("mailbox", name),
		)
		return err
	}

	return nil
}

// RenameMailbox renames a mailbox
func (u *User) RenameMailbox(oldName, newName string) error {
	u.logger.Info("renaming mailbox",
		zap.Int64("user_id", u.user.ID),
		zap.String("old_name", oldName),
		zap.String("new_name", newName),
	)

	// Prevent renaming INBOX
	if oldName == "INBOX" {
		return errors.New("cannot rename INBOX")
	}

	mb, err := u.mailboxService.GetByName(u.user.ID, oldName)
	if err != nil {
		return backend.ErrNoSuchMailbox
	}

	err = u.mailboxService.Rename(mb.ID, newName)
	if err != nil {
		u.logger.Error("failed to rename mailbox",
			zap.Error(err),
			zap.Int64("user_id", u.user.ID),
			zap.String("old_name", oldName),
			zap.String("new_name", newName),
		)
		return err
	}

	return nil
}

// Logout ends the user session
func (u *User) Logout() error {
	u.logger.Debug("IMAP session ended",
		zap.String("username", u.user.Email),
		zap.Int64("user_id", u.user.ID),
	)
	return nil
}

// extractDomain extracts domain from email address
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// extractIP extracts IP address from remote address string
func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}
