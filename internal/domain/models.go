package domain

import "time"

// Domain represents an email domain
type Domain struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	MaxUsers       int       `json:"max_users"`
	MaxMailboxSize int64     `json:"max_mailbox_size"`
	DefaultQuota   int64     `json:"default_quota"`
	CatchallEmail  string    `json:"catchall_email,omitempty"`
	BackupMX       bool      `json:"backup_mx"`
	DKIMSelector   string    `json:"dkim_selector,omitempty"`
	DKIMPrivateKey string    `json:"-"`
	DKIMPublicKey  string    `json:"dkim_public_key,omitempty"`
	SPFRecord      string    `json:"spf_record,omitempty"`
	DMARCPolicy    string    `json:"dmarc_policy,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// User represents a mail user
type User struct {
	ID               int64      `json:"id"`
	Email            string     `json:"email"`
	DomainID         int64      `json:"domain_id"`
	PasswordHash     string     `json:"-"`
	FullName         string     `json:"full_name,omitempty"`
	DisplayName      string     `json:"display_name,omitempty"`
	Quota            int64      `json:"quota"`
	UsedQuota        int64      `json:"used_quota"`
	Status           string     `json:"status"`
	AuthMethod       string     `json:"auth_method"`
	TOTPSecret       string     `json:"-"`
	TOTPEnabled      bool       `json:"totp_enabled"`
	ForwardTo        string     `json:"forward_to,omitempty"`
	AutoReplyEnabled bool       `json:"auto_reply_enabled"`
	AutoReplySubject string     `json:"auto_reply_subject,omitempty"`
	AutoReplyBody    string     `json:"auto_reply_body,omitempty"`
	SpamThreshold    float64    `json:"spam_threshold"`
	Language         string     `json:"language"`
	LastLogin        *time.Time `json:"last_login,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Alias represents an email alias
type Alias struct {
	ID                int64     `json:"id"`
	AliasEmail        string    `json:"alias_email"`
	DomainID          int64     `json:"domain_id"`
	DestinationEmails string    `json:"destination_emails"` // JSON array
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
}

// Mailbox represents a mail folder
type Mailbox struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	Name        string    `json:"name"`
	ParentID    *int64    `json:"parent_id,omitempty"`
	Subscribed  bool      `json:"subscribed"`
	SpecialUse  string    `json:"special_use,omitempty"`
	UIDValidity int64     `json:"uid_validity"`
	UIDNext     int64     `json:"uid_next"`
	CreatedAt   time.Time `json:"created_at"`
}

// Message represents an email message
type Message struct {
	ID            int64     `json:"id"`
	UserID        int64     `json:"user_id"`
	MailboxID     int64     `json:"mailbox_id"`
	UID           uint32    `json:"uid"`
	Size          int64     `json:"size"`
	Flags         string    `json:"flags"`
	Categories    string    `json:"categories"`
	ThreadID      string    `json:"thread_id,omitempty"`
	ReceivedAt    time.Time `json:"received_at"`
	InternalDate  time.Time `json:"internal_date"`
	Subject       string    `json:"subject,omitempty"`
	From          string    `json:"from,omitempty"`
	To            string    `json:"to,omitempty"`
	CC            string    `json:"cc,omitempty"`
	BCC           string    `json:"bcc,omitempty"`
	ReplyTo       string    `json:"reply_to,omitempty"`
	MessageID     string    `json:"message_id,omitempty"`
	InReplyTo     string    `json:"in_reply_to,omitempty"`
	Refs          string    `json:"refs,omitempty"`
	Headers       string    `json:"headers,omitempty"`
	BodyStructure string    `json:"body_structure,omitempty"`
	StorageType   string    `json:"storage_type"`
	Content       []byte    `json:"-"`
	ContentPath   string    `json:"content_path,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// QueueItem represents a queued message for delivery
type QueueItem struct {
	ID           int64      `json:"id"`
	Sender       string     `json:"sender"`
	Recipients   string     `json:"recipients"` // JSON array
	MessageID    string     `json:"message_id,omitempty"`
	MessagePath  string     `json:"message_path"`
	RetryCount   int        `json:"retry_count"`
	MaxRetries   int        `json:"max_retries"`
	NextRetry    *time.Time `json:"next_retry,omitempty"`
	Status       string     `json:"status"`
	ErrorMessage string     `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
