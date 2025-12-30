package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config holds the entire application configuration
type Config struct {
	Server   ServerConfig   `mapstructure:"server" yaml:"server"`
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Logger   LoggerConfig   `mapstructure:"logger" yaml:"logger"`
	TLS      TLSConfig      `mapstructure:"tls" yaml:"tls"`
	SMTP     SMTPConfig     `mapstructure:"smtp" yaml:"smtp"`
	IMAP     IMAPConfig     `mapstructure:"imap" yaml:"imap"`
	Security SecurityConfig `mapstructure:"security" yaml:"security"`
}

// ServerConfig holds general server configuration
type ServerConfig struct {
	Hostname string `mapstructure:"hostname" yaml:"hostname" env:"SERVER_HOSTNAME"`
	Domain   string `mapstructure:"domain" yaml:"domain" env:"SERVER_DOMAIN"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Path       string `mapstructure:"path" yaml:"path" env:"DB_PATH" default:"./data/mailserver.db"`
	WALEnabled bool   `mapstructure:"wal_enabled" yaml:"wal_enabled" env:"DB_WAL_ENABLED" default:"true"`
}

// TLSConfig holds TLS/certificate configuration
type TLSConfig struct {
	CertFile string     `yaml:"cert_file" env:"TLS_CERT_FILE"`
	KeyFile  string     `yaml:"key_file" env:"TLS_KEY_FILE"`
	ACME     ACMEConfig `yaml:"acme"`
}

// ACMEConfig holds Let's Encrypt ACME configuration
type ACMEConfig struct {
	Enabled  bool   `yaml:"enabled" env:"ACME_ENABLED" default:"false"`
	Email    string `yaml:"email" env:"ACME_EMAIL"`
	Provider string `yaml:"provider" env:"ACME_PROVIDER" default:"cloudflare"`
	APIToken string `yaml:"api_token" env:"CLOUDFLARE_API_TOKEN"`
	CacheDir string `yaml:"cache_dir" env:"ACME_CACHE_DIR" default:"./data/acme"`
}

// SMTPConfig holds SMTP server configuration
type SMTPConfig struct {
	SubmissionPort int    `mapstructure:"submission_port" yaml:"submission_port" env:"SMTP_SUBMISSION_PORT" default:"587"`
	RelayPort      int    `mapstructure:"relay_port" yaml:"relay_port" env:"SMTP_RELAY_PORT" default:"25"`
	SMTPSPort      int    `mapstructure:"smtps_port" yaml:"smtps_port" env:"SMTPS_PORT" default:"465"`
	MaxMessageSize int64  `mapstructure:"max_message_size" yaml:"max_message_size" env:"SMTP_MAX_MESSAGE_SIZE" default:"52428800"` // 50MB
	Hostname       string `mapstructure:"hostname" yaml:"hostname" env:"SMTP_HOSTNAME"`
}

// IMAPConfig holds IMAP server configuration
type IMAPConfig struct {
	Port        int `mapstructure:"port" yaml:"port" env:"IMAP_PORT" default:"143"`
	IMAPSPort   int `mapstructure:"imaps_port" yaml:"imaps_port" env:"IMAPS_PORT" default:"993"`
	IdleTimeout int `mapstructure:"idle_timeout" yaml:"idle_timeout" env:"IMAP_IDLE_TIMEOUT" default:"1800"` // 30 minutes
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	ClamAV       ClamAVConfig       `yaml:"clamav"`
	SpamAssassin SpamAssassinConfig `yaml:"spamassassin"`
	Greylisting  GreylistingConfig  `yaml:"greylisting"`
}

// ClamAVConfig holds ClamAV configuration
type ClamAVConfig struct {
	Enabled    bool   `yaml:"enabled" env:"CLAMAV_ENABLED" default:"true"`
	SocketPath string `yaml:"socket_path" env:"CLAMAV_SOCKET_PATH" default:"/var/run/clamav/clamd.ctl"`
}

// SpamAssassinConfig holds SpamAssassin configuration
type SpamAssassinConfig struct {
	Enabled bool   `yaml:"enabled" env:"SPAMASSASSIN_ENABLED" default:"true"`
	Host    string `yaml:"host" env:"SPAMASSASSIN_HOST" default:"localhost"`
	Port    int    `yaml:"port" env:"SPAMASSASSIN_PORT" default:"783"`
}

// GreylistingConfig holds greylisting configuration
type GreylistingConfig struct {
	Enabled      bool `yaml:"enabled" env:"GREYLISTING_ENABLED" default:"true"`
	DelayMinutes int  `yaml:"delay_minutes" env:"GREYLISTING_DELAY_MINUTES" default:"5"`
	ExpiryDays   int  `yaml:"expiry_days" env:"GREYLISTING_EXPIRY_DAYS" default:"30"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set config file path
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("gomailserver")
		v.SetConfigType("yaml")
		v.AddConfigPath("/etc/gomailserver")
		v.AddConfigPath("$HOME/.gomailserver")
		v.AddConfigPath(".")
	}

	// Enable environment variables
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Set defaults
	setDefaults(v)

	// Unmarshal into config struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Database directory creation will be handled by database package
	return &cfg, nil
}

//nolint:mnd // Default configuration values
func setDefaults(v *viper.Viper) {
	// Database
	v.SetDefault("database.path", "./data/mailserver.db")
	v.SetDefault("database.wal_enabled", true)

	// Logger
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.output_path", "stdout")

	// SMTP
	v.SetDefault("smtp.submission_port", 587)
	v.SetDefault("smtp.relay_port", 25)
	v.SetDefault("smtp.smtps_port", 465)
	v.SetDefault("smtp.max_message_size", 52428800) // 50MB

	// IMAP
	v.SetDefault("imap.port", 143)
	v.SetDefault("imap.imaps_port", 993)
	v.SetDefault("imap.idle_timeout", 1800) // 30 minutes

	// Security
	v.SetDefault("security.clamav.enabled", true)
	v.SetDefault("security.clamav.socket_path", "/var/run/clamav/clamd.ctl")
	v.SetDefault("security.spamassassin.enabled", true)
	v.SetDefault("security.spamassassin.host", "localhost")
	v.SetDefault("security.spamassassin.port", 783)
	v.SetDefault("security.greylisting.enabled", true)
	v.SetDefault("security.greylisting.delay_minutes", 5)
	v.SetDefault("security.greylisting.expiry_days", 30)

	// TLS/ACME
	v.SetDefault("tls.acme.enabled", false)
	v.SetDefault("tls.acme.provider", "cloudflare")
	v.SetDefault("tls.acme.cache_dir", "./data/acme")
}
