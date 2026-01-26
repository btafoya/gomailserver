package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
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
	API      APIConfig      `mapstructure:"api" yaml:"api"`
	WebUI    WebUIConfig    `mapstructure:"webui" yaml:"webui"`
	WebDAV   WebDAVConfig   `mapstructure:"webdav" yaml:"webdav"`
	Security SecurityConfig `mapstructure:"security" yaml:"security"`
}

// ServerConfig holds general server configuration
type ServerConfig struct {
	Hostname string `mapstructure:"hostname" yaml:"hostname" env:"SERVER_HOSTNAME"`
	Domain   string `mapstructure:"domain" yaml:"domain" env:"SERVER_DOMAIN"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string         `mapstructure:"driver" yaml:"driver" env:"DB_DRIVER" default:"sqlite3"`
	SQLite   SQLiteConfig   `mapstructure:"sqlite" yaml:"sqlite"`
	Postgres PostgresConfig `mapstructure:"postgres" yaml:"postgres"`
}

// SQLiteConfig holds SQLite-specific configuration
type SQLiteConfig struct {
	Path       string `mapstructure:"path" yaml:"path" env:"DB_SQLITE_PATH" default:"./data/mailserver.db"`
	WALEnabled bool   `mapstructure:"wal_enabled" yaml:"wal_enabled" env:"DB_SQLITE_WAL_ENABLED" default:"true"`
}

// PostgresConfig holds PostgreSQL-specific configuration
type PostgresConfig struct {
	Host            string `mapstructure:"host" yaml:"host" env:"DB_POSTGRES_HOST" default:"localhost"`
	Port            int    `mapstructure:"port" yaml:"port" env:"DB_POSTGRES_PORT" default:"5432"`
	Database        string `mapstructure:"database" yaml:"database" env:"DB_POSTGRES_DATABASE" default:"gomailserver"`
	User            string `mapstructure:"user" yaml:"user" env:"DB_POSTGRES_USER" default:"gomailserver"`
	Password        string `mapstructure:"password" yaml:"password" env:"DB_POSTGRES_PASSWORD"`
	SSLMode         string `mapstructure:"ssl_mode" yaml:"ssl_mode" env:"DB_POSTGRES_SSL_MODE" default:"disable"`
	MaxOpenConns    int    `mapstructure:"max_open_conns" yaml:"max_open_conns" env:"DB_POSTGRES_MAX_OPEN_CONNS" default:"25"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns" env:"DB_POSTGRES_MAX_IDLE_CONNS" default:"5"`
	ConnMaxLifetime string `mapstructure:"conn_max_lifetime" yaml:"conn_max_lifetime" env:"DB_POSTGRES_CONN_MAX_LIFETIME" default:"1h"`
	ConnMaxIdleTime string `mapstructure:"conn_max_idle_time" yaml:"conn_max_idle_time" env:"DB_POSTGRES_CONN_MAX_IDLE_TIME" default:"30m"`
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

// APIConfig holds admin API server configuration
type APIConfig struct {
	Port           int      `mapstructure:"port" yaml:"port" env:"API_PORT" default:"8980"`
	ReadTimeout    int      `mapstructure:"read_timeout" yaml:"read_timeout" env:"API_READ_TIMEOUT" default:"15"`
	WriteTimeout   int      `mapstructure:"write_timeout" yaml:"write_timeout" env:"API_WRITE_TIMEOUT" default:"15"`
	MaxHeaderBytes int      `mapstructure:"max_header_bytes" yaml:"max_header_bytes" env:"API_MAX_HEADER_BYTES" default:"1048576"` // 1MB
	AdminToken     string   `mapstructure:"admin_token" yaml:"admin_token" env:"API_ADMIN_TOKEN"`                                  // Bearer token for admin endpoints (deprecated, use JWT)
	JWTSecret      string   `mapstructure:"jwt_secret" yaml:"jwt_secret" env:"API_JWT_SECRET"`                                     // Secret for signing JWT tokens
	CORSOrigins    []string `mapstructure:"cors_origins" yaml:"cors_origins" env:"API_CORS_ORIGINS"`                               // Allowed CORS origins
}

// WebUIConfig holds Web UI (admin and webmail) server configuration
type WebUIConfig struct {
	Enabled  bool `mapstructure:"enabled" yaml:"enabled" env:"WEBUI_ENABLED" default:"true"`
	Port     int  `mapstructure:"port" yaml:"port" env:"WEBUI_PORT" default:"8080"`
	VitePort int  `mapstructure:"vite_port" yaml:"vite_port" env:"WEBUI_VITE_PORT" default:"5173"` // Dev mode Vite server port
}

// WebDAVConfig holds WebDAV/CalDAV/CardDAV server configuration
type WebDAVConfig struct {
	Enabled      bool `mapstructure:"enabled" yaml:"enabled" env:"WEBDAV_ENABLED" default:"true"`
	Port         int  `mapstructure:"port" yaml:"port" env:"WEBDAV_PORT" default:"8800"`
	ReadTimeout  int  `mapstructure:"read_timeout" yaml:"read_timeout" env:"WEBDAV_READ_TIMEOUT" default:"30"`
	WriteTimeout int  `mapstructure:"write_timeout" yaml:"write_timeout" env:"WEBDAV_WRITE_TIMEOUT" default:"30"`
}

// SecurityConfig holds external security service connection configuration
// All security policies and settings are stored in SQLite per-domain
type SecurityConfig struct {
	ClamAV       ClamAVConfig       `mapstructure:"clamav" yaml:"clamav"`
	SpamAssassin SpamAssassinConfig `mapstructure:"spamassassin" yaml:"spamassassin"`
	DNS          DNSConfig          `mapstructure:"dns" yaml:"dns"`
}

// DNSConfig holds DNS resolver configuration for DANE, MTA-STS, and other DNS lookups
type DNSConfig struct {
	Resolver        string `mapstructure:"resolver" yaml:"resolver" env:"DNS_RESOLVER" default:"1.1.1.1:53"`
	FallbackServers []string `mapstructure:"fallback_servers" yaml:"fallback_servers" env:"DNS_FALLBACK_SERVERS"`
	Timeout         int    `mapstructure:"timeout" yaml:"timeout" env:"DNS_TIMEOUT" default:"5"`
	UseTCP          bool   `mapstructure:"use_tcp" yaml:"use_tcp" env:"DNS_USE_TCP" default:"false"`
}

// ClamAVConfig holds ClamAV connection configuration
// All scanning policies and actions are stored in SQLite per-domain
type ClamAVConfig struct {
	SocketPath string `mapstructure:"socket_path" yaml:"socket_path" env:"CLAMAV_SOCKET_PATH" default:"/var/run/clamav/clamd.ctl"`
	Timeout    int    `mapstructure:"timeout" yaml:"timeout" env:"CLAMAV_TIMEOUT" default:"60"`
}

// SpamAssassinConfig holds SpamAssassin connection configuration
// All spam scoring policies and thresholds are stored in SQLite per-domain/user
type SpamAssassinConfig struct {
	Host    string `mapstructure:"host" yaml:"host" env:"SPAMASSASSIN_HOST" default:"localhost"`
	Port    int    `mapstructure:"port" yaml:"port" env:"SPAMASSASSIN_PORT" default:"783"`
	Timeout int    `mapstructure:"timeout" yaml:"timeout" env:"SPAMASSASSIN_TIMEOUT" default:"30"`
}

// Load loads configuration from file and environment variables
func Load(configPath string) (*Config, error) {
	// Load .env file if it exists (from current directory or config directory)
	loadDotEnv(configPath)

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

	// Validate security configuration
	if err := cfg.ValidateSecurityConfig(); err != nil {
		return nil, fmt.Errorf("invalid security configuration: %w", err)
	}

	// Database directory creation will be handled by database package
	return &cfg, nil
}

//nolint:mnd // Default configuration values
func setDefaults(v *viper.Viper) {
	// Database
	v.SetDefault("database.driver", "sqlite3")
	v.SetDefault("database.sqlite.path", "./data/mailserver.db")
	v.SetDefault("database.sqlite.wal_enabled", true)
	v.SetDefault("database.postgres.host", "localhost")
	v.SetDefault("database.postgres.port", "5432")
	v.SetDefault("database.postgres.database", "gomailserver")
	v.SetDefault("database.postgres.user", "gomailserver")
	v.SetDefault("database.postgres.ssl_mode", "disable")
	v.SetDefault("database.postgres.max_open_conns", "25")
	v.SetDefault("database.postgres.max_idle_conns", "5")
	v.SetDefault("database.postgres.conn_max_lifetime", "1h")
	v.SetDefault("database.postgres.conn_max_idle_time", "30m")

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

	// Admin API
	v.SetDefault("api.port", 8980)
	v.SetDefault("api.read_timeout", 15)
	v.SetDefault("api.write_timeout", 15)
	v.SetDefault("api.max_header_bytes", 1048576) // 1MB

	// WebUI (Admin and Webmail interfaces)
	v.SetDefault("webui.enabled", true)
	v.SetDefault("webui.port", 8080)

	// WebDAV/CalDAV/CardDAV
	v.SetDefault("webdav.enabled", true)
	v.SetDefault("webdav.port", 8800)
	v.SetDefault("webdav.read_timeout", 30)
	v.SetDefault("webdav.write_timeout", 30)

	// Security - External service connections only
	// All security policies are stored in SQLite per-domain
	v.SetDefault("security.clamav.socket_path", "/var/run/clamav/clamd.ctl")
	v.SetDefault("security.clamav.timeout", 60)
	v.SetDefault("security.spamassassin.host", "localhost")
	v.SetDefault("security.spamassassin.port", 783)
	v.SetDefault("security.spamassassin.timeout", 30)

	// DNS resolver configuration for DANE, MTA-STS, etc.
	v.SetDefault("security.dns.resolver", "1.1.1.1:53")
	v.SetDefault("security.dns.fallback_servers", []string{"8.8.8.8:53", "9.9.9.9:53"})
	v.SetDefault("security.dns.timeout", 5)
	v.SetDefault("security.dns.use_tcp", false)

	// TLS/ACME
	v.SetDefault("tls.acme.enabled", false)
	v.SetDefault("tls.acme.provider", "cloudflare")
	v.SetDefault("tls.acme.cache_dir", "./data/acme")
}

// loadDotEnv attempts to load .env file from various locations
func loadDotEnv(configPath string) {
	// Try loading from config directory first
	if configPath != "" {
		configDir := filepath.Dir(configPath)
		envPath := filepath.Join(configDir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}
	}

	// Try current directory
	if _, err := os.Stat(".env"); err == nil {
		_ = godotenv.Load(".env")
		return
	}

	// Try /etc/gomailserver/.env
	if _, err := os.Stat("/etc/gomailserver/.env"); err == nil {
		_ = godotenv.Load("/etc/gomailserver/.env")
		return
	}

	// Try $HOME/.gomailserver/.env
	if home, err := os.UserHomeDir(); err == nil {
		envPath := filepath.Join(home, ".gomailserver", ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
		}
	}
}

// ValidateSecurityConfig validates external security service connection settings
// All security policies are validated at the domain level in the database
func (c *Config) ValidateSecurityConfig() error {
	// ClamAV connection validation
	if c.Security.ClamAV.SocketPath == "" {
		return fmt.Errorf("clamav.socket_path cannot be empty")
	}
	if c.Security.ClamAV.Timeout <= 0 {
		return fmt.Errorf("clamav.timeout must be positive, got %d", c.Security.ClamAV.Timeout)
	}

	// SpamAssassin connection validation
	if c.Security.SpamAssassin.Host == "" {
		return fmt.Errorf("spamassassin.host cannot be empty")
	}
	if c.Security.SpamAssassin.Port <= 0 || c.Security.SpamAssassin.Port > 65535 {
		return fmt.Errorf("spamassassin.port must be between 1 and 65535, got %d", c.Security.SpamAssassin.Port)
	}
	if c.Security.SpamAssassin.Timeout <= 0 {
		return fmt.Errorf("spamassassin.timeout must be positive, got %d", c.Security.SpamAssassin.Timeout)
	}

	return nil
}
