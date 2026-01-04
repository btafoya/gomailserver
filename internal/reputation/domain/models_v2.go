package domain

import "time"

// DMARC Report Models

type DMARCReport struct {
	ID              int64
	Domain          string
	ReportID        string
	BeginTime       int64
	EndTime         int64
	Organization    string
	TotalMessages   int
	SPFPass         int
	DKIMPass        int
	AlignmentPass   int
	RawXML          string
	ProcessedAt     int64
	Records         []*DMARCReportRecord
}

type DMARCReportRecord struct {
	ID            int64
	ReportID      int64
	SourceIP      string
	Count         int
	Disposition   string // none, quarantine, reject
	SPFResult     string // pass, fail, neutral, softfail, temperror, permerror
	DKIMResult    string // pass, fail, neutral, temperror, permerror
	SPFAligned    bool
	DKIMAligned   bool
	HeaderFrom    string
	EnvelopeFrom  string
}

type DMARCAutoAction struct {
	ID           int64
	Domain       string
	IssueType    DMARCIssueType
	Description  string
	ActionTaken  string
	TakenAt      int64
	Success      bool
	ErrorMessage string
}

type DMARCIssueType string

const (
	IssueTypeSPFMisalign  DMARCIssueType = "spf_misalign"
	IssueTypeDKIMMisalign DMARCIssueType = "dkim_misalign"
	IssueTypeDKIMFail     DMARCIssueType = "dkim_fail"
	IssueTypeSPFFail      DMARCIssueType = "spf_fail"
)

type AlignmentAnalysis struct {
	Domain             string
	ReportID           int64
	TotalMessages      int
	AlignmentPassRate  float64
	SPFPassRate        float64
	DKIMPassRate       float64
	SPFAlignmentRate   float64
	DKIMAlignmentRate  float64
	Issues             []*AlignmentIssue
	RecommendedActions []string
}

type AlignmentIssue struct {
	IssueType   DMARCIssueType
	SourceIP    string
	Count       int
	Description string
	Severity    string // low, medium, high, critical
}

// ARF Complaint Models

type ARFReport struct {
	ID                     int64
	ReceivedAt             int64
	FeedbackType           string // abuse, fraud, virus, other
	UserAgent              string
	Version                string
	OriginalRcptTo         string
	ArrivalDate            int64
	ReportingMTA           string
	SourceIP               string
	AuthenticationResults  string
	MessageID              string
	Subject                string
	RawReport              string
	Processed              bool
	SuppressedRecipient    string
}

// External Metrics Models

type PostmasterMetrics struct {
	ID                 int64
	Domain             string
	FetchedAt          int64
	MetricDate         int64
	DomainReputation   string // HIGH, MEDIUM, LOW, BAD
	SpamRate           float64
	IPReputation       string // HIGH, MEDIUM, LOW, BAD
	AuthenticationRate float64
	EncryptionRate     float64
	UserSpamReports    int
	RawResponse        string
}

type SNDSMetrics struct {
	ID            int64
	IPAddress     string
	FetchedAt     int64
	MetricDate    int64
	SpamTrapHits  int
	ComplaintRate float64
	FilterLevel   string // GREEN, YELLOW, RED
	MessageCount  int
	RawResponse   string
}

// Provider-Specific Rate Limiting Models

type ProviderRateLimit struct {
	ID                   int64
	Domain               string
	Provider             MailProvider
	MaxHourlyRate        int
	MaxDailyRate         int
	CurrentHourCount     int
	CurrentDayCount      int
	HourResetAt          int64
	DayResetAt           int64
	CircuitBreakerActive bool
	LastUpdated          int64
}

type MailProvider string

const (
	ProviderGmail   MailProvider = "gmail"
	ProviderOutlook MailProvider = "outlook"
	ProviderYahoo   MailProvider = "yahoo"
	ProviderGeneric MailProvider = "generic"
)

// Custom Warm-up Models

type CustomWarmupSchedule struct {
	ID           int64
	Domain       string
	ScheduleName string
	Day          int
	MaxVolume    int
	CreatedAt    int64
	CreatedBy    string
	IsActive     bool
}

// Prediction Models

type ReputationPrediction struct {
	ID                    int64
	Domain                string
	PredictedAt           int64
	PredictionHorizon     int // hours ahead
	PredictedScore        int
	PredictedComplaintRate float64
	PredictedBounceRate   float64
	ConfidenceLevel       float64
	ModelVersion          string
	FeaturesUsed          map[string]interface{}
}

// Alert Models

type ReputationAlert struct {
	ID             int64
	Domain         string
	AlertType      AlertType
	Severity       AlertSeverity
	Title          string
	Message        string
	Details        map[string]interface{}
	CreatedAt      int64
	Acknowledged   bool
	AcknowledgedAt int64
	AcknowledgedBy string
	Resolved       bool
	ResolvedAt     int64
}

type AlertType string

const (
	AlertTypeDNSFailure        AlertType = "dns_failure"
	AlertTypeScoreDrop         AlertType = "score_drop"
	AlertTypeCircuitBreaker    AlertType = "circuit_breaker"
	AlertTypeExternalFeedback  AlertType = "external_feedback"
	AlertTypeDMARCIssue        AlertType = "dmarc_issue"
)

type AlertSeverity string

const (
	SeverityLow      AlertSeverity = "low"
	SeverityMedium   AlertSeverity = "medium"
	SeverityHigh     AlertSeverity = "high"
	SeverityCritical AlertSeverity = "critical"
)

// Helper Functions

func (d *DMARCReport) GetAlignmentRate() float64 {
	if d.TotalMessages == 0 {
		return 0.0
	}
	return float64(d.AlignmentPass) / float64(d.TotalMessages)
}

func (d *DMARCReport) GetSPFPassRate() float64 {
	if d.TotalMessages == 0 {
		return 0.0
	}
	return float64(d.SPFPass) / float64(d.TotalMessages)
}

func (d *DMARCReport) GetDKIMPassRate() float64 {
	if d.TotalMessages == 0 {
		return 0.0
	}
	return float64(d.DKIMPass) / float64(d.TotalMessages)
}

func (p *ProviderRateLimit) ShouldResetHour(now time.Time) bool {
	return now.Unix() >= p.HourResetAt
}

func (p *ProviderRateLimit) ShouldResetDay(now time.Time) bool {
	return now.Unix() >= p.DayResetAt
}

func (p *ProviderRateLimit) IsAtHourlyLimit() bool {
	return p.CurrentHourCount >= p.MaxHourlyRate
}

func (p *ProviderRateLimit) IsAtDailyLimit() bool {
	if p.MaxDailyRate == 0 {
		return false
	}
	return p.CurrentDayCount >= p.MaxDailyRate
}

func (a *ReputationAlert) IsUnacknowledged() bool {
	return !a.Acknowledged
}

func (a *ReputationAlert) IsUnresolved() bool {
	return !a.Resolved
}

func (p *ReputationPrediction) IsHighConfidence() bool {
	return p.ConfidenceLevel >= 0.7
}
