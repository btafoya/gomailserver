package testing

import (
	"time"
)

// TestConfig holds configuration for test execution
type TestConfig struct {
	// Server configuration
	SMTPAddr     string `yaml:"smtp_addr" json:"smtp_addr"`
	IMAPAddr     string `yaml:"imap_addr" json:"imap_addr"`
	HTTPAddr     string `yaml:"http_addr" json:"http_addr"`
	DatabasePath string `yaml:"database_path" json:"database_path"`

	// Test settings
	Timeout time.Duration `yaml:"timeout" json:"timeout"`
	Verbose bool          `yaml:"verbose" json:"verbose"`
	Debug   bool          `yaml:"debug" json:"debug"`

	// Authentication
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`

	// Security settings
	VerifyDKIM  bool `yaml:"verify_dkim" json:"verify_dkim"`
	VerifySPF   bool `yaml:"verify_spf" json:"verify_spf"`
	VerifyDMARC bool `yaml:"verify_dmarc" json:"verify_dmarc"`

	// Report settings
	OutputDir  string `yaml:"output_dir" json:"output_dir"`
	HTMLReport bool   `yaml:"html_report" json:"html_report"`
	JSONReport bool   `yaml:"json_report" json:"json_report"`

	// Advanced settings
	MaxRetries int           `yaml:"max_retries" json:"max_retries"`
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() TestConfig {
	return TestConfig{
		SMTPAddr:     "localhost:587",
		IMAPAddr:     "localhost:143",
		HTTPAddr:     "http://localhost:8980",
		DatabasePath: "./data/mailserver.db",
		Timeout:      30 * time.Second,
		Verbose:      false,
		Debug:        false,
		Username:     "test@example.com",
		Password:     "password",
		VerifyDKIM:   true,
		VerifySPF:    true,
		VerifyDMARC:  true,
		OutputDir:    "./test-reports",
		HTMLReport:   true,
		JSONReport:   true,
		MaxRetries:   3,
		RetryDelay:   1 * time.Second,
	}
}

// TestResult represents the result of a test execution
type TestResult struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Passed      bool          `json:"passed"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Duration    time.Duration `json:"duration"`
	Trace       []TraceEvent  `json:"trace"`
	Errors      []error       `json:"errors"`
	Summary     string        `json:"summary"`
}

// SecuritySummary contains security verification results
type SecuritySummary struct {
	DKIM          DKIMSummary     `json:"dkim"`
	SPF           SPFSummary      `json:"spf"`
	DMARC         DMARCSummary    `json:"dmarc"`
	Reputation    ReputationScore `json:"reputation"`
	OverallStatus string          `json:"overall_status"`
}

// DKIMSummary contains DKIM verification results
type DKIMSummary struct {
	Signed    bool      `json:"signed"`
	Verified  bool      `json:"verified"`
	Selector  string    `json:"selector"`
	KeySize   int       `json:"key_size"`
	Algorithm string    `json:"algorithm"`
	ValidFrom time.Time `json:"valid_from"`
	Error     string    `json:"error,omitempty"`
}

// SPFSummary contains SPF verification results
type SPFSummary struct {
	Result  string `json:"result"`
	Record  string `json:"record"`
	IPRange string `json:"ip_range"`
	Aligned bool   `json:"aligned"`
	Error   string `json:"error,omitempty"`
}

// DMARCSummary contains DMARC verification results
type DMARCSummary struct {
	Result    string `json:"result"`
	Policy    string `json:"policy"`
	Aligned   bool   `json:"aligned"`
	Alignment string `json:"alignment"`
	Pct       int    `json:"pct"`
	Error     string `json:"error,omitempty"`
}

// ReputationScore contains reputation information
type ReputationScore struct {
	Score         int     `json:"score"`
	DeliveryRate  float64 `json:"delivery_rate"`
	BounceRate    float64 `json:"bounce_rate"`
	ComplaintRate float64 `json:"complaint_rate"`
	Status        string  `json:"status"`
}

// TraceEvent represents a single trace event
type TraceEvent struct {
	Timestamp time.Time              `json:"timestamp"`
	Phase     string                 `json:"phase"`
	Component string                 `json:"component"`
	Action    string                 `json:"action"`
	Duration  time.Duration          `json:"duration"`
	Status    string                 `json:"status"`
	Details   map[string]interface{} `json:"details,omitempty"`
	Error     string                 `json:"error,omitempty"`
}

// NewTestResult creates a new test result
func NewTestResult(name, description string) *TestResult {
	return &TestResult{
		ID:          generateID(),
		Name:        name,
		Description: description,
		StartTime:   time.Now(),
		Trace:       []TraceEvent{},
		Errors:      []error{},
	}
}

// Complete marks the test result as completed
func (r *TestResult) Complete(passed bool, summary string) {
	r.EndTime = time.Now()
	r.Duration = r.EndTime.Sub(r.StartTime)
	r.Passed = passed
	r.Summary = summary
}

// AddError adds an error to the test result
func (r *TestResult) AddError(err error) {
	r.Errors = append(r.Errors, err)
}

// AddTrace adds a trace event to the result
func (r *TestResult) AddTrace(event TraceEvent) {
	r.Trace = append(r.Trace, event)
}

// generateID generates a unique test ID
func generateID() string {
	return time.Now().Format("20060102-150405-") + randomString(6)
}

// randomString generates a random string of given length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
		time.Sleep(1) // Ensure different values
	}
	return string(b)
}
