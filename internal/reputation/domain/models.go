package domain

// SendingEvent represents a single email sending event
type SendingEvent struct {
	ID                 int64
	Timestamp          int64
	Domain             string
	RecipientDomain    string
	EventType          EventType
	BounceType         *string
	EnhancedStatusCode *string
	SMTPResponse       *string
	IPAddress          string
	Metadata           map[string]interface{}
}

// EventType represents the type of sending event
type EventType string

const (
	// EventDelivery indicates successful message delivery
	EventDelivery EventType = "delivery"
	// EventBounce indicates message bounce
	EventBounce EventType = "bounce"
	// EventDefer indicates temporary delivery failure
	EventDefer EventType = "defer"
	// EventComplaint indicates spam complaint
	EventComplaint EventType = "complaint"
)

// BounceType represents the type of bounce
type BounceType string

const (
	// BounceHard indicates permanent delivery failure
	BounceHard BounceType = "hard"
	// BounceSoft indicates temporary delivery failure
	BounceSoft BounceType = "soft"
)

// ReputationScore represents the reputation metrics for a domain
type ReputationScore struct {
	Domain               string
	ReputationScore      int     // 0-100
	ComplaintRate        float64 // Percentage
	BounceRate           float64 // Percentage
	DeliveryRate         float64 // Percentage
	CircuitBreakerActive bool
	CircuitBreakerReason string
	WarmUpActive         bool
	WarmUpDay            int
	LastUpdated          int64
}

// WarmUpDay represents a day in the warm-up schedule
type WarmUpDay struct {
	Domain       string
	Day          int
	MaxVolume    int
	ActualVolume int
	CreatedAt    int64
}

// CircuitBreakerEvent represents a circuit breaker trigger event
type CircuitBreakerEvent struct {
	ID           int64
	Domain       string
	TriggerType  TriggerType
	TriggerValue float64
	Threshold    float64
	PausedAt     int64
	ResumedAt    *int64
	AutoResumed  bool
	AdminNotes   string
}

// TriggerType represents the type of circuit breaker trigger
type TriggerType string

const (
	// TriggerComplaint indicates high complaint rate
	TriggerComplaint TriggerType = "complaint"
	// TriggerBounce indicates high bounce rate
	TriggerBounce TriggerType = "bounce"
	// TriggerBlock indicates repeated blocks from major providers
	TriggerBlock TriggerType = "block"
)

// RetentionPolicy represents data retention settings
type RetentionPolicy struct {
	ID           int
	RetentionDays int
	LastCleanup  int64
}

// AuditResult represents the deliverability readiness audit for a domain
type AuditResult struct {
	Domain       string
	Timestamp    int64
	SPFStatus    CheckStatus
	DKIMStatus   CheckStatus
	DMARCStatus  CheckStatus
	RDNSStatus   CheckStatus
	FCrDNSStatus CheckStatus
	TLSStatus    CheckStatus
	MTASTSStatus CheckStatus
	PostmasterOK bool
	AbuseOK      bool
	OverallScore int // 0-100
	Issues       []string
}

// CheckStatus represents the status of a deliverability check
type CheckStatus struct {
	Passed  bool
	Message string
	Details map[string]interface{}
}
