package domain

import "time"

// MTASTSPolicy represents a cached MTA-STS policy
type MTASTSPolicy struct {
	ID         int64     `json:"id"`
	Domain     string    `json:"domain"`
	Version    string    `json:"version"`
	Mode       string    `json:"mode"`
	MaxAge     int       `json:"max_age"`
	MXPatterns string    `json:"mx_patterns"` // JSON array of MX patterns
	FetchedAt  time.Time `json:"fetched_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	PolicyText string    `json:"policy_text"`
}

// MTA-STS Policy modes
const (
	MTASTSModeNone    = "none"
	MTASTSModeTesting = "testing"
	MTASTSModeEnforce = "enforce"
)

// TLSReport represents a TLS reporting entry (TLSRPT - RFC 8460)
type TLSReport struct {
	ID             int64     `json:"id"`
	ReportID       string    `json:"report_id"`
	Domain         string    `json:"domain"`
	DateRangeStart time.Time `json:"date_range_start"`
	DateRangeEnd   time.Time `json:"date_range_end"`
	ContactInfo    string    `json:"contact_info,omitempty"`
	ReportJSON     string    `json:"report_json"`
	CreatedAt      time.Time `json:"created_at"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
}
