package imap

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
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
			zap.Error(err))
		return nil, err
	}

	// Calculate message counts
	totalMessages := len(messages)
	recentMessages := 0
	unseenMessages := 0
	firstUnseenSeqNum := uint32(0)

	// Count messages by flags
	for i, msg := range messages {
		// Check if recent (received within last 5 minutes)
		if time.Since(msg.ReceivedAt) < 5*time.Minute {
			recentMessages++
		}

		// Check if unseen
		if !strings.Contains(msg.Flags, "\\Seen") {
			unseenMessages++
			// Track first unseen sequence number (1-indexed)
			if firstUnseenSeqNum == 0 {
				firstUnseenSeqNum = uint32(i + 1)
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
		zap.Int("items_count", len(items)),
	)

	// Fetch all messages from the mailbox
	messages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		m.logger.Error("failed to get messages",
			zap.Error(err))
		return fmt.Errorf("failed to get messages: %w", err)
	}

	// Expand fetch items (e.g., ALL -> FLAGS, INTERNALDATE, RFC822.SIZE, ENVELOPE)
	expandedItems := make([]imap.FetchItem, 0, len(items))
	for _, item := range items {
		expandedItems = append(expandedItems, item.Expand()...)
	}

	// Process each message
	for seqNum, msg := range messages {
		seqNumVal := uint32(seqNum + 1) // 1-indexed

		// Check if message matches the sequence set
		if seqSet != nil {
			var matches bool
			if uid {
				matches = seqSet.Contains(msg.UID)
			} else {
				matches = seqSet.Contains(seqNumVal)
			}
			if !matches {
				continue
			}
		}

		// Create IMAP message with requested items
		imapMsg := imap.NewMessage(seqNumVal, expandedItems)
		imapMsg.Uid = msg.UID

		// Populate requested items
		for _, item := range expandedItems {
			switch item {
			case imap.FetchEnvelope:
				imapMsg.Envelope = m.buildEnvelope(msg)

			case imap.FetchFlags:
				imapMsg.Flags = m.parseFlags(msg.Flags)

			case imap.FetchInternalDate:
				imapMsg.InternalDate = msg.InternalDate

			case imap.FetchRFC822Size:
				imapMsg.Size = uint32(msg.Size)

			case imap.FetchUid:
				// Already set above

			case imap.FetchBodyStructure, imap.FetchBody:
				imapMsg.BodyStructure = m.buildBodyStructure(msg)

			default:
				// Check if it's a body section request (e.g., BODY[], BODY[HEADER])
				if strings.HasPrefix(string(item), "BODY[") || strings.HasPrefix(string(item), "BODY.PEEK[") {
					section, err := imap.ParseBodySectionName(item)
					if err != nil {
						m.logger.Warn("failed to parse body section",
							zap.String("item", string(item)),
							zap.Error(err))
						continue
					}

					// Get the full message content
					fullMsg, err := m.messageService.GetByID(msg.ID)
					if err != nil {
						m.logger.Error("failed to get message content",
							zap.Error(err))
						continue
					}

					// Initialize body map if needed
					if imapMsg.Body == nil {
						imapMsg.Body = make(map[*imap.BodySectionName]imap.Literal)
					}

					// Extract the requested section
					content := m.extractBodySection(fullMsg, section)
					imapMsg.Body[section] = content
				}
			}
		}

		// Send message to channel
		ch <- imapMsg
	}

	return nil
}

// parseFlags converts a comma-separated flags string to a slice
func (m *Mailbox) parseFlags(flags string) []string {
	if flags == "" {
		return nil
	}
	var result []string
	for _, flag := range strings.Split(flags, ",") {
		flag = strings.TrimSpace(flag)
		if flag != "" {
			result = append(result, flag)
		}
	}
	return result
}

// bodyStructureJSON represents the stored body structure format
type bodyStructureJSON struct {
	Parts []string `json:"parts"`
}

// buildBodyStructure creates an IMAP body structure from a message
func (m *Mailbox) buildBodyStructure(msg *domain.Message) *imap.BodyStructure {
	// Parse the stored body structure if available
	if msg.BodyStructure != "" {
		var stored bodyStructureJSON
		if err := json.Unmarshal([]byte(msg.BodyStructure), &stored); err == nil {
			// Single part message
			if len(stored.Parts) == 1 {
				mimeType, mimeSubType, params := parseContentType(stored.Parts[0])
				return &imap.BodyStructure{
					MIMEType:    mimeType,
					MIMESubType: mimeSubType,
					Params:      params,
					Size:        uint32(msg.Size),
					Extended:    true,
				}
			}

			// Multipart message
			if len(stored.Parts) > 1 {
				bs := &imap.BodyStructure{
					MIMEType:    "multipart",
					MIMESubType: "mixed",
					Parts:       make([]*imap.BodyStructure, 0, len(stored.Parts)),
					Extended:    true,
				}

				for _, part := range stored.Parts {
					mimeType, mimeSubType, params := parseContentType(part)

					// Check if this is a multipart boundary definition
					if mimeType == "multipart" {
						bs.MIMESubType = mimeSubType
						bs.Params = params
						continue
					}

					partBS := &imap.BodyStructure{
						MIMEType:    mimeType,
						MIMESubType: mimeSubType,
						Params:      params,
						Extended:    true,
					}
					bs.Parts = append(bs.Parts, partBS)
				}

				return bs
			}
		}
	}

	// Default to simple text/plain structure
	return &imap.BodyStructure{
		MIMEType:    "text",
		MIMESubType: "plain",
		Params:      map[string]string{"charset": "utf-8"},
		Size:        uint32(msg.Size),
		Extended:    true,
	}
}

// parseContentType parses a Content-Type header value into MIME type, subtype, and parameters
func parseContentType(contentType string) (mimeType, mimeSubType string, params map[string]string) {
	params = make(map[string]string)

	// Default values
	mimeType = "text"
	mimeSubType = "plain"

	if contentType == "" {
		return
	}

	// Split by semicolon to separate type from parameters
	parts := strings.Split(contentType, ";")
	if len(parts) == 0 {
		return
	}

	// Parse the MIME type (e.g., "text/plain")
	typeParts := strings.SplitN(strings.TrimSpace(parts[0]), "/", 2)
	if len(typeParts) >= 1 {
		mimeType = strings.ToLower(strings.TrimSpace(typeParts[0]))
	}
	if len(typeParts) >= 2 {
		mimeSubType = strings.ToLower(strings.TrimSpace(typeParts[1]))
	}

	// Parse parameters (e.g., charset=utf-8, boundary=xxx)
	for i := 1; i < len(parts); i++ {
		param := strings.TrimSpace(parts[i])
		if param == "" {
			continue
		}

		kv := strings.SplitN(param, "=", 2)
		if len(kv) == 2 {
			key := strings.ToLower(strings.TrimSpace(kv[0]))
			value := strings.TrimSpace(kv[1])
			// Remove quotes if present
			value = strings.Trim(value, "\"'")
			params[key] = value
		}
	}

	return
}

// extractBodySection extracts a specific body section from a message
func (m *Mailbox) extractBodySection(msg *domain.Message, section *imap.BodySectionName) imap.Literal {
	content := msg.Content
	if len(content) == 0 {
		return bytes.NewReader(nil)
	}

	// Determine what section is requested
	specifier := section.Specifier

	switch specifier {
	case imap.HeaderSpecifier:
		// Return only headers
		if idx := bytes.Index(content, []byte("\r\n\r\n")); idx != -1 {
			content = content[:idx+2]
		} else if idx := bytes.Index(content, []byte("\n\n")); idx != -1 {
			content = content[:idx+1]
		}

	case imap.TextSpecifier:
		// Return only body (after headers)
		if idx := bytes.Index(content, []byte("\r\n\r\n")); idx != -1 {
			content = content[idx+4:]
		} else if idx := bytes.Index(content, []byte("\n\n")); idx != -1 {
			content = content[idx+2:]
		}

	case imap.MIMESpecifier:
		// Return MIME headers for a specific part
		// For simple messages, this is similar to header
		if idx := bytes.Index(content, []byte("\r\n\r\n")); idx != -1 {
			content = content[:idx+2]
		}

	default:
		// EntireSpecifier or no specifier - return full content
	}

	// Apply partial extraction if specified
	if len(section.Partial) >= 2 {
		start := section.Partial[0]
		length := section.Partial[1]
		if start < len(content) {
			end := start + length
			if end > len(content) {
				end = len(content)
			}
			content = content[start:end]
		} else {
			content = nil
		}
	}

	return bytes.NewReader(content)
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
			zap.Error(err))
		return nil, err
	}

	var matchingUIDs []uint32

	// Apply search criteria
	for seqNum, msg := range messages {
		matches := true

		// Search by header fields (criteria.Header is textproto.MIMEHeader)
		if criteria.Header != nil && len(criteria.Header) > 0 {
			// Check subject
			if subjects, ok := criteria.Header["Subject"]; ok {
				for _, searchSubject := range subjects {
					if !containsFold(msg.Subject, searchSubject) {
						matches = false
						break
					}
				}
			}

			// Check from
			if froms, ok := criteria.Header["From"]; ok && matches {
				for _, searchFrom := range froms {
					if !containsFold(msg.From, searchFrom) {
						matches = false
						break
					}
				}
			}

			// Check to
			if tos, ok := criteria.Header["To"]; ok && matches {
				for _, searchTo := range tos {
					if !containsFold(msg.To, searchTo) {
						matches = false
						break
					}
				}
			}
		}

		// Search by body content (criteria.Body is []string)
		if criteria.Body != nil && matches {
			for _, searchText := range criteria.Body {
				// Search in message content
				content := string(msg.Content)
				if !containsFold(content, searchText) {
					matches = false
					break
				}
			}
		}

		// Search by text (header + body)
		if criteria.Text != nil && matches {
			for _, searchText := range criteria.Text {
				// Search in subject, from, to, and content
				found := containsFold(msg.Subject, searchText) ||
					containsFold(msg.From, searchText) ||
					containsFold(msg.To, searchText) ||
					containsFold(string(msg.Content), searchText)
				if !found {
					matches = false
					break
				}
			}
		}

		// Search by flags - WithFlags (all must be present)
		if criteria.WithFlags != nil && matches {
			for _, flag := range criteria.WithFlags {
				if !strings.Contains(msg.Flags, flag) {
					matches = false
					break
				}
			}
		}

		// Search by flags - WithoutFlags (none must be present)
		if criteria.WithoutFlags != nil && matches {
			for _, flag := range criteria.WithoutFlags {
				if strings.Contains(msg.Flags, flag) {
					matches = false
					break
				}
			}
		}

		// Search by date criteria (time.Time values, use IsZero to check if set)
		if !criteria.Since.IsZero() && msg.ReceivedAt.Before(criteria.Since) {
			matches = false
		}
		if !criteria.Before.IsZero() && msg.ReceivedAt.After(criteria.Before) {
			matches = false
		}

		// Search by size (uint32 values, check against 0)
		if criteria.Larger > 0 && uint32(msg.Size) <= criteria.Larger {
			matches = false
		}
		if criteria.Smaller > 0 && uint32(msg.Size) >= criteria.Smaller {
			matches = false
		}

		// Search by UID (criteria.Uid is *SeqSet)
		if criteria.Uid != nil && !criteria.Uid.Contains(msg.UID) {
			matches = false
		}

		// Search by sequence number (criteria.SeqNum is *SeqSet)
		if criteria.SeqNum != nil && !criteria.SeqNum.Contains(uint32(seqNum+1)) {
			matches = false
		}

		// Add to results if matches
		if matches {
			if uid {
				matchingUIDs = append(matchingUIDs, msg.UID)
			} else {
				matchingUIDs = append(matchingUIDs, uint32(seqNum+1))
			}
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

	// Read message content from the literal body
	messageData, err := io.ReadAll(body)
	if err != nil {
		m.logger.Error("failed to read message body",
			zap.Error(err))
		return fmt.Errorf("failed to read message body: %w", err)
	}

	// Get the next UID for this mailbox
	uid := m.mailbox.UIDNext

	// Store the message
	msg, err := m.messageService.Store(m.user.ID, m.mailbox.ID, uid, messageData)
	if err != nil {
		m.logger.Error("failed to store message",
			zap.Error(err))
		return fmt.Errorf("failed to store message: %w", err)
	}

	// Set flags if provided
	if len(flags) > 0 {
		ctx := context.Background()
		if err := m.messageService.UpdateFlags(ctx, int(msg.ID), int(m.user.ID), flags, "set"); err != nil {
			m.logger.Error("failed to set message flags",
				zap.Error(err))
			return fmt.Errorf("failed to set message flags: %w", err)
		}
	}

	// Increment UIDNext for the mailbox
	m.mailbox.UIDNext++

	m.logger.Info("message created",
		zap.Int64("message_id", msg.ID),
		zap.Int64("uid", uid),
		zap.Int64("size", int64(len(messageData))),
	)

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
			zap.Error(err))
		return err
	}

	// Filter messages by sequence set
	var messagesToUpdate []*domain.Message
	for seqNum, msg := range messages {
		var shouldInclude bool
		if seqSet == nil {
			shouldInclude = true
		} else if uid {
			shouldInclude = seqSet.Contains(msg.UID)
		} else {
			shouldInclude = seqSet.Contains(uint32(seqNum + 1))
		}

		if shouldInclude {
			messagesToUpdate = append(messagesToUpdate, msg)
		}
	}

	// Determine the action string for the service
	var action string
	switch operation {
	case imap.SetFlags:
		action = "set"
	case imap.AddFlags:
		action = "add"
	case imap.RemoveFlags:
		action = "remove"
	default:
		action = "set"
	}

	// Update each message
	ctx := context.Background()
	for _, msg := range messagesToUpdate {
		if err := m.messageService.UpdateFlags(ctx, int(msg.ID), int(m.user.ID), flags, action); err != nil {
			m.logger.Error("failed to update message flags",
				zap.Error(err))
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
			zap.Error(err))
		return err
	}

	// Get all messages from source mailbox
	allMessages, err := m.messageService.GetByMailbox(m.mailbox.ID)
	if err != nil {
		return err
	}

	// Filter messages by sequence set
	var messagesToCopy []*domain.Message
	for seqNum, msg := range allMessages {
		var shouldInclude bool
		if seqSet == nil {
			shouldInclude = true
		} else if uid {
			shouldInclude = seqSet.Contains(msg.UID)
		} else {
			shouldInclude = seqSet.Contains(uint32(seqNum + 1))
		}

		if shouldInclude {
			messagesToCopy = append(messagesToCopy, msg)
		}
	}

	// Copy each message
	for _, msg := range messagesToCopy {
		// Get the full message content
		fullMsg, err := m.messageService.GetByID(msg.ID)
		if err != nil {
			m.logger.Error("failed to get full message for copy",
				zap.Error(err))
			return err
		}

		// Store copy in destination mailbox with same content
		_, err = m.messageService.Store(m.user.ID, destMailbox.ID, int64(destMailbox.UIDNext), fullMsg.Content)
		if err != nil {
			m.logger.Error("failed to copy message",
				zap.Error(err))
			return err
		}
	}

	m.logger.Info("messages copied",
		zap.String("destination", dest),
		zap.Int("count", len(messagesToCopy)),
	)

	return nil
}

// buildEnvelope creates an IMAP envelope from a message
func (m *Mailbox) buildEnvelope(msg *domain.Message) *imap.Envelope {
	return &imap.Envelope{
		Date:      msg.InternalDate,
		Subject:   msg.Subject,
		From:      m.parseAddresses(msg.From),
		Sender:    m.parseAddresses(msg.From),
		ReplyTo:   m.parseAddresses(msg.ReplyTo),
		To:        m.parseAddresses(msg.To),
		Cc:        m.parseAddresses(msg.CC),
		Bcc:       m.parseAddresses(msg.BCC),
		InReplyTo: msg.InReplyTo,
		MessageId: msg.MessageID,
	}
}

// parseAddresses parses a comma-separated address string into IMAP addresses
func (m *Mailbox) parseAddresses(addrs string) []*imap.Address {
	if addrs == "" {
		return nil
	}

	var result []*imap.Address
	for _, addr := range strings.Split(addrs, ",") {
		addr = strings.TrimSpace(addr)
		if addr == "" {
			continue
		}

		imapAddr := &imap.Address{}

		// Parse "Name <email@host>" or "email@host" format
		if idx := strings.Index(addr, "<"); idx != -1 {
			imapAddr.PersonalName = strings.TrimSpace(addr[:idx])
			addr = strings.TrimSuffix(strings.TrimPrefix(addr[idx:], "<"), ">")
		}

		// Split email into mailbox@host
		if atIdx := strings.LastIndex(addr, "@"); atIdx != -1 {
			imapAddr.MailboxName = addr[:atIdx]
			imapAddr.HostName = addr[atIdx+1:]
		} else {
			imapAddr.MailboxName = addr
		}

		result = append(result, imapAddr)
	}
	return result
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
			zap.Error(err))
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
				zap.Error(err))
			return err
		}
	}

	m.logger.Info("messages expunged",
		zap.Int("count", len(messagesToDelete)),
	)

	return nil
}

// containsFold performs a case-insensitive substring search
func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
