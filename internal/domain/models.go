package domain

import "time"

// Domain represents an email domain with per-domain security configuration
type Domain struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	MaxUsers       int       `json:"max_users"`
	MaxMailboxSize int64     `json:"max_mailbox_size"`
	DefaultQuota   int64     `json:"default_quota"`
	CatchallEmail  string    `json:"catchall_email,omitempty"`
	BackupMX       bool      `json:"backup_mx"`

	// DKIM configuration
	DKIMSelector        string `json:"dkim_selector,omitempty"`
	DKIMPrivateKey      string `json:"-"`
	DKIMPublicKey       string `json:"dkim_public_key,omitempty"`
	DKIMSigningEnabled  bool   `json:"dkim_signing_enabled"`
	DKIMVerifyEnabled   bool   `json:"dkim_verify_enabled"`
	DKIMKeySize         int    `json:"dkim_key_size"`
	DKIMKeyType         string `json:"dkim_key_type"`
	DKIMHeadersToSign   string `json:"dkim_headers_to_sign"` // JSON array

	// SPF configuration
	SPFRecord          string `json:"spf_record,omitempty"`
	SPFEnabled         bool   `json:"spf_enabled"`
	SPFDNSServer       string `json:"spf_dns_server"`
	SPFDNSTimeout      int    `json:"spf_dns_timeout"`
	SPFMaxLookups      int    `json:"spf_max_lookups"`
	SPFFailAction      string `json:"spf_fail_action"`
	SPFSoftFailAction  string `json:"spf_softfail_action"`

	// DMARC configuration
	DMARCPolicy        string `json:"dmarc_policy,omitempty"`
	DMARCEnabled       bool   `json:"dmarc_enabled"`
	DMARCDNSServer     string `json:"dmarc_dns_server"`
	DMARCDNSTimeout    int    `json:"dmarc_dns_timeout"`
	DMARCReportEnabled bool   `json:"dmarc_report_enabled"`
	DMARCReportEmail   string `json:"dmarc_report_email,omitempty"`

	// ClamAV antivirus configuration
	ClamAVEnabled      bool   `json:"clamav_enabled"`
	ClamAVMaxScanSize  int64  `json:"clamav_max_scan_size"`
	ClamAVVirusAction  string `json:"clamav_virus_action"`
	ClamAVFailAction   string `json:"clamav_fail_action"`

	// SpamAssassin configuration
	SpamEnabled           bool    `json:"spam_enabled"`
	SpamRejectScore       float64 `json:"spam_reject_score"`
	SpamQuarantineScore   float64 `json:"spam_quarantine_score"`
	SpamLearningEnabled   bool    `json:"spam_learning_enabled"`

	// Greylisting configuration
	GreylistEnabled         bool `json:"greylist_enabled"`
	GreylistDelayMinutes    int  `json:"greylist_delay_minutes"`
	GreylistExpiryDays      int  `json:"greylist_expiry_days"`
	GreylistCleanupInterval int  `json:"greylist_cleanup_interval"`
	GreylistWhitelistAfter  int  `json:"greylist_whitelist_after"`

	// Rate limiting configuration (JSON objects)
	RateLimitEnabled           bool   `json:"ratelimit_enabled"`
	RateLimitSMTPPerIP         string `json:"ratelimit_smtp_per_ip"`         // JSON: {"count":100,"window_minutes":60}
	RateLimitSMTPPerUser       string `json:"ratelimit_smtp_per_user"`       // JSON
	RateLimitSMTPPerDomain     string `json:"ratelimit_smtp_per_domain"`     // JSON
	RateLimitAuthPerIP         string `json:"ratelimit_auth_per_ip"`         // JSON
	RateLimitIMAPPerUser       string `json:"ratelimit_imap_per_user"`       // JSON
	RateLimitCleanupInterval   int    `json:"ratelimit_cleanup_interval"`

	// Authentication security configuration
	AuthTOTPEnforced             bool `json:"auth_totp_enforced"`
	AuthBruteForceEnabled        bool `json:"auth_brute_force_enabled"`
	AuthBruteForceThreshold      int  `json:"auth_brute_force_threshold"`
	AuthBruteForceWindowMinutes  int  `json:"auth_brute_force_window_minutes"`
	AuthBruteForceBlockMinutes   int  `json:"auth_brute_force_block_minutes"`
	AuthIPBlacklistEnabled       bool `json:"auth_ip_blacklist_enabled"`
	AuthCleanupInterval          int  `json:"auth_cleanup_interval"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// DKIMConfig represents DKIM signing configuration
type DKIMConfig struct {
	Domain     string `json:"domain"`
	Selector   string `json:"selector"`
	PrivateKey []byte `json:"-"`
	PublicKey  string `json:"public_key,omitempty"`
}

// AntivirusConfig represents antivirus configuration
type AntivirusConfig struct {
	VirusAction string `json:"virus_action"` // reject, quarantine, tag
}

// GreylistTriplet represents a greylisting entry
type GreylistTriplet struct {
	ID        int64     `json:"id"`
	IP        string    `json:"ip"`
	Sender    string    `json:"sender"`
	Recipient string    `json:"recipient"`
	FirstSeen time.Time `json:"first_seen"`
	PassCount int       `json:"pass_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RateLimitEntry represents a rate limit tracking entry
type RateLimitEntry struct {
	ID         int64     `json:"id"`
	Key        string    `json:"key"` // IP or user identifier
	Type       string    `json:"type"` // "ip" or "user"
	Count      int       `json:"count"`
	WindowStart time.Time `json:"window_start"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// LoginAttempt represents a login attempt for brute force tracking
type LoginAttempt struct {
	ID        int64     `json:"id"`
	IP        string    `json:"ip"`
	Email     string    `json:"email,omitempty"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

// IPBlacklist represents a blacklisted IP address
type IPBlacklist struct {
	ID        int64      `json:"id"`
	IP        string     `json:"ip"`
	Reason    string     `json:"reason"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// QuarantineMessage represents a quarantined message
type QuarantineMessage struct {
	ID          int64     `json:"id"`
	MessageID   string    `json:"message_id"`
	Sender      string    `json:"sender"`
	Recipient   string    `json:"recipient"`
	Subject     string    `json:"subject,omitempty"`
	Reason      string    `json:"reason"` // virus, spam
	Score       float64   `json:"score,omitempty"`
	MessagePath string    `json:"message_path"`
	Action      string    `json:"action"` // quarantined, deleted, released
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
