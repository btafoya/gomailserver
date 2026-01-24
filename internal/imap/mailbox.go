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

		messages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		m.logger.Error("failed to get messages for status",
			zap.Error(err),
			return nil, err
	}

	// Calculate message counts
	totalMessages := len(messages)
	recentMessages := 0
	unseenMessages := 0
	firstUnseenSeqNum := uint32(0)

	// Count messages by flags
	for _, msg := range messages {
		// Check if recent (received within last 5 minutes)
		if time.Since(msg.ReceivedAt) < 5*time.Minute {
			recentMessages++
		}

		// Check if unseen
		if !strings.Contains(msg.Flags, "\\Seen") {
			unseenMessages++
			// Track first unseen sequence number
			if firstUnseenSeqNum == 0 && msg.UID < firstUnseenSeqNum {
				firstUnseenSeqNum = uint32(msg.UID)
			}
		}
	}

	status.Messages = uint32(totalMessages)
	status.Recent = uint32(recentMessages)
	status.Unseen = uint32(unseenMessages)
	status.UnseenSeqNum = firstUnseenSeqNum

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

	// Get all messages for searching
	messages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		m.logger.Error("failed to get messages for search",
			zap.Error(err),
			return nil, err
	}

	var matchingUIDs []uint32

	// Apply search criteria
	for _, msg := range messages {
		matches := true

		// Search by subject
		if criteria.Header != nil {
			for _, field := range criteria.Header {
				if strings.EqualFold(field.Key, "subject") && strings.ContainsFold(msg.Subject, field.Value) {
					continue // Matches, don't exclude
				}
				if strings.EqualFold(field.Key, "subject") && !strings.ContainsFold(msg.Subject, field.Value) {
					matches = false
					break
				}
			}
		}

		// Search by from
		if criteria.Header != nil && matches {
			for _, field := range criteria.Header {
				if strings.EqualFold(field.Key, "from") && strings.ContainsFold(msg.Sender, field.Value) {
					continue // Matches, don't exclude
				}
				if strings.EqualFold(field.Key, "from") && !strings.ContainsFold(msg.Sender, field.Value) {
					matches = false
					break
				}
			}
		}

		// Search by body content
		if criteria.Body != nil && matches {
			for _, field := range criteria.Body {
				if strings.EqualFold(field.Key, "text") && strings.ContainsFold(msg.Body, field.Value) {
					continue // Matches, don't exclude
				}
				if strings.EqualFold(field.Key, "text") && !strings.ContainsFold(msg.Body, field.Value) {
					matches = false
					break
				}
			}
		}

		// Search by flags
		if criteria.Flag != nil && matches {
			for _, flag := range criteria.Flag {
				hasFlag := strings.Contains(msg.Flags, "\\"+flag)
				switch flag.Not {
				case true:
					if hasFlag {
						matches = false
						break
					}
				case false:
					if !hasFlag {
						matches = false
						break
					}
				}
			}
		}

		// Search by date criteria
		if criteria.Since != nil && msg.ReceivedAt.Before(*criteria.Since) {
			matches = false
		}
		if criteria.Before != nil && msg.ReceivedAt.After(*criteria.Before) {
			matches = false
		}

		// Search by size
		if criteria.Larger != nil && uint32(len(msg.Body)) < *criteria.Larger {
			matches = false
		}
		if criteria.Smaller != nil && uint32(len(msg.Body)) > *criteria.Smaller {
			matches = false
		}

		// Search by UID
		if criteria.UID != nil {
			for _, uidSet := range criteria.UID.Set() {
				if uint32(msg.UID) == uidSet {
					matches = true
					break
				}
			}
		}

		// Search by sequence number
		if criteria.SeqNum != nil && msg.SequenceNum == *criteria.SeqNum {
			matches = true
		}

		// Add to results if matches
		if matches {
			matchingUIDs = append(matchingUIDs, uint32(msg.UID))
		}
	}

	return matchingUIDs, nil
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

	// Get messages to update
	messages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		m.logger.Error("failed to get messages for flag update",
			zap.Error(err),
			return err
	}

	// Filter by sequence set if provided
	var messagesToUpdate []*domain.Message
	if seqSet != nil {
		// Convert sequence set to message indices
		indices := seqSet.Set()
		for _, index := range indices {
			// Convert to 0-based for slice access
			msgIndex := int(index) - 1
			if msgIndex >= 0 && msgIndex < len(messages) {
				messagesToUpdate = append(messagesToUpdate, messages[msgIndex])
			}
		}
	} else {
		// No sequence set, update all messages
		messagesToUpdate = messages
	}

	// Convert IMAP flags to domain flags
	domainFlags := m.convertIMAPToDomainFlags(flags)

	// Update each message
	for _, msg := range messagesToUpdate {
		var newFlags string
		switch operation {
		case imap.FlagsSet:
			newFlags = domainFlags
		case imap.FlagsAdd:
			// Add new flags to existing ones
			existingFlags := m.convertDomainToMessageFlags(msg.Flags)
			mergedFlags := append(existingFlags, domainFlags...)
			// Remove duplicates
			uniqueFlags := []string{}
			flagSet := make(map[string]bool)
			for _, flag := range mergedFlags {
				if !flagSet[flag] {
					flagSet[flag] = true
					uniqueFlags = append(uniqueFlags, flag)
				}
			}
			newFlags = strings.Join(uniqueFlags, ",")
		case imap.FlagsRemove:
			// Remove specified flags from existing ones
			existingFlags := m.convertDomainToMessageFlags(msg.Flags)
			remainingFlags := []string{}
			flagSet := make(map[string]bool)
			for _, flag := range existingFlags {
				flagSet[flag] = true
			}
			for _, flag := range domainFlags {
				if !flagSet[flag] {
					remainingFlags = append(remainingFlags, flag)
				}
			}
			newFlags = strings.Join(remainingFlags, ",")
		}

		// Update message in database
		if err := m.messageService.UpdateFlags(msg.ID, newFlags); err != nil {
			m.logger.Error("failed to update message flags",
				zap.Error(err),
				return err
		}
	}

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

	// Get destination mailbox
	destMailbox, err := m.mailboxService.GetByName(m.user.ID, dest)
	if err != nil {
		m.logger.Error("failed to get destination mailbox",
			zap.Error(err),
			return err
	}

	// Get messages to copy
	var messagesToCopy []*domain.Message
	if seqSet != nil {
		// Convert sequence set to message indices
		indices := seqSet.Set()
		for _, index := range indices {
			// Convert to 0-based for slice access
			msgIndex := int(index) - 1
			allMessages, err := m.messageService.GetByMailbox(m.mailbox.ID)
			if err != nil {
				return err
			}
			if msgIndex >= 0 && msgIndex < len(allMessages) {
				messagesToCopy = append(messagesToCopy, allMessages[msgIndex])
			}
		}
	} else {
		// No sequence set, get all messages
		allMessages, err := m.messageService.GetByMailbox(m.mailbox.ID)
		if err != nil {
			return err
		}
		messagesToCopy = allMessages
	}

	// Copy each message
	for _, msg := range messagesToCopy {
		// Create message copy with same content and flags
		copyMsg := &domain.Message{
			Subject:      msg.Subject,
			Body:         msg.Body,
			Flags:        msg.Flags,
			ReceivedAt:    msg.ReceivedAt,
			Sender:       msg.Sender,
			Recipients:  msg.Recipients,
			ThreadID:     msg.ThreadID,
			MessageID:    msg.MessageID,
		}

		// Store copy in destination mailbox
		_, err := m.messageService.Store(copyMsg, destMailbox.ID, m.convertDomainToMessageFlags(msg.Flags))
		if err != nil {
			m.logger.Error("failed to copy message",
				zap.Error(err),
				return err
		}
	}

	m.logger.Info("messages copied",
		zap.String("destination", dest),
		zap.Int("count", len(messagesToCopy)),
	)

	return nil
}

// Expunge permanently removes messages marked for deletion
func (m *Mailbox) Expunge() error {
	m.logger.Debug("expunging messages",
		zap.Int64("mailbox_id", m.mailbox.ID),
		zap.String("mailbox", m.mailbox.Name),
	)

	// Get messages to delete
	messages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		m.logger.Error("failed to get messages for expunge",
			zap.Error(err),
			return err
	}

	// Find messages with \Deleted flag
	var messagesToDelete []*domain.Message
	for _, msg := range messages {
		if strings.Contains(msg.Flags, "\\Deleted") {
			messagesToDelete = append(messagesToDelete, msg)
		}
	}

	// Delete messages permanently
	for _, msg := range messagesToDelete {
		if err := m.messageService.Delete(msg.ID); err != nil {
			m.logger.Error("failed to delete message",
				zap.Error(err),
				return err
		}
	}

	m.logger.Info("messages expunged",
		zap.Int("count", len(messagesToDelete)),
	)

	return nil
}
