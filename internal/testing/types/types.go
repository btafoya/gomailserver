package types

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Category string

const (
	CategoryConfig      Category = "config"
	CategoryMailFlow    Category = "mailflow"
	CategorySecurity    Category = "security"
	CategoryPerformance Category = "performance"
)

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

type Status string

const (
	StatusPass    Status = "pass"
	StatusFail    Status = "fail"
	StatusWarning Status = "warning"
	StatusSkip    Status = "skip"
)

type OutputMode string

const (
	OutputQuiet   OutputMode = "quiet"
	OutputSummary OutputMode = "summary"
	OutputVerbose OutputMode = "verbose"
)

type Mode string

const (
	ModeLocal  Mode = "local"
	ModeRemote Mode = "remote"
)

type CheckResult struct {
	Check     string                 `json:"check"`
	Category  Category               `json:"category"`
	Severity  Severity               `json:"severity"`
	Status    Status                 `json:"status"`
	Message   string                 `json:"message"`
	Details   map[string]interface{} `json:"details"`
	Duration  int64                  `json:"duration_ms"`
	Error     string                 `json:"error,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

type ServerConfig struct {
	SMTPHost    string
	SMTPPort    int
	IMAPHost    string
	IMAPPort    int
	Domains     []string
	TestUser    string
	TestPass    string
	TLS         bool
	StartTLS    bool
	Timeout     time.Duration
	ConfigPath  string
	PasswordEnv string
	DryRun      bool
	AutoCleanup bool
}

type Profile struct {
	Name        string         `yaml:"name"`
	SMTPHost    string         `yaml:"smtp_host"`
	SMTPPort    int            `yaml:"smtp_port"`
	IMAPHost    string         `yaml:"imap_host"`
	IMAPPort    int            `yaml:"imap_port"`
	Domains     []string       `yaml:"domains"`
	TestUser    string         `yaml:"test_user"`
	PasswordEnv string         `yaml:"password_env"`
	Options     ProfileOptions `yaml:"options"`
}

type ProfileOptions struct {
	TLS      bool          `yaml:"tls"`
	StartTLS bool          `yaml:"starttls"`
	Timeout  time.Duration `yaml:"timeout"`
}

func LoadProfile(name string) (*Profile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	paths := []string{
		filepath.Join(homeDir, ".gomailserver/profiles", name+".yaml"),
		"/etc/gomailserver/profiles/" + name + ".yaml",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			var profile Profile
			if err := yaml.Unmarshal(data, &profile); err != nil {
				return nil, err
			}
			return &profile, nil
		}
	}

	return nil, os.ErrNotExist
}

func SaveProfile(profile *Profile, path string) error {
	data, err := yaml.Marshal(profile)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (cfg *ServerConfig) SMTPAddress() string {
	return net.JoinHostPort(cfg.SMTPHost, fmt.Sprintf("%d", cfg.SMTPPort))
}

func (cfg *ServerConfig) IMAPAddress() string {
	return net.JoinHostPort(cfg.IMAPHost, fmt.Sprintf("%d", cfg.IMAPPort))
}

func (cfg *ServerConfig) TLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.SMTPHost,
	}
}
