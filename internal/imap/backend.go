package imap

import (
	"errors"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
)

// Backend implements IMAP backend interface
type Backend struct {
	userService    service.UserServiceInterface
	mailboxService service.MailboxServiceInterface
	messageService service.MessageServiceInterface
	logger         *zap.Logger
}

// Login authenticates a user
func (b *Backend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	b.logger.Info("IMAP authentication attempt",
		zap.String("username", username),
		zap.String("remote_addr", connInfo.RemoteAddr.String()),
	)

	user, err := b.userService.Authenticate(username, password)
	if err != nil {
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
