package verifier

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type Config struct {
	Mode        Mode
	ConfigFile  string
	ProfileName string
	OutputMode  types.OutputMode
	DryRun      bool
	WarningsOk  bool
	RateLimit   int
	NoCleanup   bool
	ReportHTML  string
	ReportJSON  string
}

func DefaultConfig() *Config {
	return &Config{
		Mode:        ModeLocal,
		ConfigFile:  "",
		ProfileName: "",
		OutputMode:  types.OutputSummary,
		DryRun:      false,
		WarningsOk:  false,
		RateLimit:   10,
		NoCleanup:   false,
		ReportHTML:  "",
		ReportJSON:  "",
	}
}

type Mode string

const (
	ModeLocal  Mode = "local"
	ModeRemote Mode = "remote"
)

type ProfileOptions struct {
	TLS      bool          `yaml:"tls"`
	StartTLS bool          `yaml:"starttls"`
	Timeout  time.Duration `yaml:"timeout"`
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

func LoadProfile(name string) (*Profile, error) {
	paths := []string{
		"~/.gomailserver/profiles/" + name + ".yaml",
		"/etc/gomailserver/profiles/" + name + ".yaml",
	}

	for _, path := range paths {
		expandedPath := os.ExpandEnv(path)
		data, err := os.ReadFile(expandedPath)
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
