package imap

import (
	"time"

	"github.com/emersion/go-imap"
	"go.uber.org/zap"

	"github.com/btafoya/gomailserver/internal/domain"
	"github.com/btafoya/gomailserver/internal/service"
)

// Mailbox implements IMAP mailbox interface
type Mailbox struct {
	mailbox        *domain.Mailbox
	user           *domain.User
	messageService service.MessageServiceInterface
	mailboxService service.MailboxServiceInterface
	logger         *zap.Logger
}

// Name returns the mailbox name
func (m *Mailbox) Name() string {
	return m.mailbox.Name
}

// Info returns mailbox information
func (m *Mailbox) Info() (*imap.MailboxInfo, error) {
	info := &imap.MailboxInfo{
		Attributes: []string{},
		Delimiter:  "/",
		Name:       m.mailbox.Name,
	}

	// Add special-use attributes
	switch m.mailbox.SpecialUse {
	case "\\Drafts":
		info.Attributes = append(info.Attributes, imap.DraftsAttr)
	case "\\Sent":
		info.Attributes = append(info.Attributes, imap.SentAttr)
	case "\\Trash":
		info.Attributes = append(info.Attributes, imap.TrashAttr)
	case "\\Junk":
		info.Attributes = append(info.Attributes, imap.JunkAttr)
	case "\\Archive":
		info.Attributes = append(info.Attributes, imap.ArchiveAttr)
	}

	return info, nil
}

// Status returns mailbox status
func (m *Mailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(m.mailbox.Name, items)
	status.UidValidity = uint32(m.mailbox.UIDValidity)
	status.UidNext = uint32(m.mailbox.UIDNext)

	// TODO: Get actual message counts from database
	status.Messages = 0      // Total messages
	status.Recent = 0        // Recent messages
	status.Unseen = 0        // First unseen message sequence number
	status.UnseenSeqNum = 0  // Sequence number of first unseen

	return status, nil
}

// SetSubscribed sets the subscription status
func (m *Mailbox) SetSubscribed(subscribed bool) error {
	m.logger.Debug("setting mailbox subscription",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.Bool("subscribed", subscribed),
	)

	m.mailbox.Subscribed = subscribed
	return m.mailboxService.UpdateSubscription(m.mailbox.ID, subscribed)
}

// Check requests a checkpoint of the currently selected mailbox
func (m *Mailbox) Check() error {
	m.logger.Debug("mailbox check",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
	)
	// Checkpoint is a no-op for SQLite (WAL mode handles this)
	return nil
}

// ListMessages lists messages in the mailbox
func (m *Mailbox) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)

	m.logger.Debug("listing messages",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.Bool("uid", uid),
	)

	// TODO: Fetch messages from database
	// TODO: Apply sequence set filter
	// TODO: Fetch requested items
	// For now, return empty result

	return nil
}

// SearchMessages searches for messages matching criteria
func (m *Mailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	m.logger.Debug("searching messages",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.Bool("uid", uid),
	)

	// TODO: Implement search
	// TODO: Support various search criteria (FROM, TO, SUBJECT, etc.)
	return []uint32{}, nil
}

// CreateMessage appends a new message to the mailbox
func (m *Mailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	m.logger.Debug("creating message",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.Strings("flags", flags),
	)

	// TODO: Read message from body
	// TODO: Store message using messageService
	// TODO: Set flags
	// TODO: Increment UIDNext

	return nil
}

// UpdateMessagesFlags updates message flags
func (m *Mailbox) UpdateMessagesFlags(uid bool, seqSet *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	m.logger.Debug("updating message flags",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.Bool("uid", uid),
		zap.String("operation", string(operation)),
		zap.Strings("flags", flags),
	)

	// TODO: Fetch messages by sequence set
	// TODO: Update flags based on operation (SET, ADD, REMOVE)
	// TODO: Handle \Seen, \Deleted, \Flagged, \Answered, \Draft

	return nil
}

// CopyMessages copies messages to another mailbox
func (m *Mailbox) CopyMessages(uid bool, seqSet *imap.SeqSet, dest string) error {
	m.logger.Debug("copying messages",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
		zap.String("destination", dest),
		zap.Bool("uid", uid),
	)

	// TODO: Get destination mailbox
	// TODO: Fetch messages by sequence set
	// TODO: Copy message data to destination mailbox
	// TODO: Preserve flags and date

	return nil
}

// Expunge permanently removes messages marked for deletion
func (m *Mailbox) Expunge() error {
	m.logger.Debug("expunging messages",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
	)

	// TODO: Find messages with \Deleted flag
	// TODO: Permanently delete them from storage
	// TODO: Update sequence numbers

	return nil
}
