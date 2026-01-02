package domain

import "time"

// AuditLog represents an audit log entry for tracking admin actions and security events
type AuditLog struct {
	ID           int64     `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	UserID       *int64    `json:"user_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   string    `json:"resource_id,omitempty"`
	Details      string    `json:"details,omitempty"`
	IPAddress    string    `json:"ip_address,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Severity     string    `json:"severity"`
	Success      bool      `json:"success"`
}

// Severity levels for audit logs
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityError    = "error"
	SeverityCritical = "critical"
)

// Common action types
const (
	ActionUserCreated        = "user.created"
	ActionUserUpdated        = "user.updated"
	ActionUserDeleted        = "user.deleted"
	ActionUserPasswordChange = "user.password_changed"
	ActionUserLogin          = "user.login"
	ActionUserLoginFailed    = "user.login_failed"
	ActionUserLogout         = "user.logout"

	ActionDomainCreated = "domain.created"
	ActionDomainUpdated = "domain.updated"
	ActionDomainDeleted = "domain.deleted"

	ActionAliasCreated = "alias.created"
	ActionAliasUpdated = "alias.updated"
	ActionAliasDeleted = "alias.deleted"

	ActionConfigUpdated = "config.updated"

	ActionSecurityDKIMEnabled  = "security.dkim_enabled"
	ActionSecuritySPFEnabled   = "security.spf_enabled"
	ActionSecurityDMARCEnabled = "security.dmarc_enabled"

	ActionPGPKeyImported = "pgp.key_imported"
	ActionPGPKeyDeleted  = "pgp.key_deleted"

	ActionMailSent     = "mail.sent"
	ActionMailReceived = "mail.received"
	ActionMailBlocked  = "mail.blocked"

	ActionSystemStartup  = "system.startup"
	ActionSystemShutdown = "system.shutdown"
)

// Resource types
const (
	ResourceTypeUser   = "user"
	ResourceTypeDomain = "domain"
	ResourceTypeAlias  = "alias"
	ResourceTypeMail   = "mail"
	ResourceTypeConfig = "config"
	ResourceTypePGP    = "pgp"
	ResourceTypeSystem = "system"
)
