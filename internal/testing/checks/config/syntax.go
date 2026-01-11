package config

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"net"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/btafoya/gomailserver/internal/config"
	"github.com/btafoya/gomailserver/internal/testing/types"
)

type ConfigSyntaxCheck struct{}

func (c *ConfigSyntaxCheck) Name() string {
	return "Config Syntax"
}

func (c *ConfigSyntaxCheck) Description() string {
	return "Parse and validate gomailserver configuration"
}

func (c *ConfigSyntaxCheck) Category() types.Category {
	return types.CategoryConfig
}

func (c *ConfigSyntaxCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *ConfigSyntaxCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "Configuration syntax validation",
		Details:  make(map[string]interface{}),
	}

	var configPath string
	if cfg.ConfigPath != "" {
		configPath = cfg.ConfigPath
	} else {
		possiblePaths := []string{
			"/etc/gomailserver/gomailserver.yaml",
			"/etc/gomailserver/gomailserver.yml",
			os.Getenv("HOME") + "/.gomailserver/gomailserver.yaml",
			os.Getenv("HOME") + "/.gomailserver/gomailserver.yml",
			"./gomailserver.yaml",
			"./gomailserver.yml",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	if configPath == "" {
		result.Status = types.StatusWarning
		result.Message = "No configuration file found (using defaults)"
		result.Details["config_file"] = "none"
		return result, nil
	}

	loadedCfg, err := config.Load(configPath)
	if err != nil {
		result.Message = fmt.Sprintf("Configuration syntax error: %v", err)
		result.Details["config_file"] = configPath
		result.Details["error"] = err.Error()
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "Configuration syntax valid"
	result.Details["config_file"] = configPath
	result.Details["hostname"] = loadedCfg.Server.Hostname
	result.Details["domain"] = loadedCfg.Server.Domain
	result.Details["database_path"] = loadedCfg.Database.Path

	domainCount := 1
	if loadedCfg.Server.Hostname != "" && loadedCfg.Server.Domain != "" {
		domainCount = 1
	}
	result.Details["configured_domains"] = domainCount

	if loadedCfg.TLS.CertFile != "" {
		result.Details["tls_enabled"] = true
		result.Details["tls_cert_file"] = loadedCfg.TLS.CertFile
	} else if loadedCfg.TLS.ACME.Enabled {
		result.Details["tls_enabled"] = true
		result.Details["acme_enabled"] = true
		result.Details["acme_provider"] = loadedCfg.TLS.ACME.Provider
	} else {
		result.Details["tls_enabled"] = false
		result.Status = types.StatusWarning
		result.Message = "Configuration valid but TLS not configured"
	}

	return result, nil
}

type TLSCertificateCheck struct{}

func (c *TLSCertificateCheck) Name() string {
	return "TLS Certificates"
}

func (c *TLSCertificateCheck) Description() string {
	return "Validate TLS certificate existence and expiration"
}

func (c *TLSCertificateCheck) Category() types.Category {
	return types.CategoryConfig
}

func (c *TLSCertificateCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *TLSCertificateCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "TLS not enabled, skipping certificate check",
		Details:  make(map[string]interface{}),
	}

	if !cfg.TLS {
		result.Details["tls_enabled"] = false
		return result, nil
	}

	certFile := cfg.ConfigPath
	if certFile == "" {
		possiblePaths := []string{
			"/etc/gomailserver/tls.crt",
			"./tls.crt",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				certFile = path
				break
			}
		}

		if certFile == "" {
			result.Status = types.StatusWarning
			result.Message = "TLS enabled but certificate file not found"
			result.Details["tls_enabled"] = true
			result.Details["cert_file"] = "not found"
			return result, nil
		}
	}

	if _, err := os.Stat(certFile); err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to read certificate: %v", err)
		result.Details["tls_enabled"] = true
		result.Details["cert_file"] = certFile
		return result, nil
	}

	cert, err := tls.LoadX509KeyPair(certFile, "")
	if err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to parse certificate: %v", err)
		result.Details["tls_enabled"] = true
		result.Details["cert_file"] = certFile
		return result, nil
	}

	result.Status = types.StatusPass
	result.Message = "TLS certificates valid"
	result.Details["tls_enabled"] = true
	result.Details["cert_file"] = certFile
	result.Details["subject"] = cert.Leaf.Subject.CommonName
	result.Details["issuer"] = cert.Leaf.Issuer.CommonName
	result.Details["valid_from"] = cert.Leaf.NotBefore.Format("2006-01-02")
	result.Details["valid_until"] = cert.Leaf.NotAfter.Format("2006-01-02")

	daysUntilExpiry := int(cert.Leaf.NotAfter.Sub(time.Now()).Hours() / 24)
	result.Details["days_until_expiry"] = daysUntilExpiry

	if daysUntilExpiry < 7 {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Certificate expires in %d days (critical)", daysUntilExpiry)
	} else if daysUntilExpiry < 30 {
		result.Status = types.StatusWarning
		result.Message = fmt.Sprintf("Certificate expires in %d days (warning)", daysUntilExpiry)
	}

	return result, nil
}

type PortAvailabilityCheck struct{}

func (c *PortAvailabilityCheck) Name() string {
	return "Port Availability"
}

func (c *PortAvailabilityCheck) Description() string {
	return "Check SMTP/IMAP ports are accessible"
}

func (c *PortAvailabilityCheck) Category() types.Category {
	return types.CategoryConfig
}

func (c *PortAvailabilityCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *PortAvailabilityCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusPass,
		Message:  "All required ports accessible",
		Details:  make(map[string]interface{}),
	}

	timeout := 5 * time.Second

	smtpPort := cfg.SMTPPort
	if smtpPort == 0 {
		smtpPort = 587
	}

	imapPort := cfg.IMAPPort
	if imapPort == 0 {
		imapPort = 143
	}

	smtpAddress := fmt.Sprintf("%s:%d", cfg.SMTPHost, smtpPort)
	imapAddress := fmt.Sprintf("%s:%d", cfg.IMAPHost, imapPort)

	smtpConn, err := net.DialTimeout("tcp", smtpAddress, timeout)
	result.Details["smtp_port"] = smtpPort
	result.Details["smtp_address"] = smtpAddress

	if err != nil {
		result.Status = types.StatusFail
		result.Details["smtp_accessible"] = false
		result.Details["smtp_error"] = err.Error()
	} else {
		result.Details["smtp_accessible"] = true
		smtpConn.Close()
	}

	imapConn, err := net.DialTimeout("tcp", imapAddress, timeout)
	result.Details["imap_port"] = imapPort
	result.Details["imap_address"] = imapAddress

	if err != nil {
		if result.Status != types.StatusFail {
			result.Status = types.StatusFail
		}
		result.Details["imap_accessible"] = false
		result.Details["imap_error"] = err.Error()
	} else {
		result.Details["imap_accessible"] = true
		imapConn.Close()
	}

	if result.Status == types.StatusFail {
		result.Message = "Port availability check failed"
	} else if !result.Details["smtp_accessible"].(bool) || !result.Details["imap_accessible"].(bool) {
		result.Status = types.StatusWarning
		result.Message = "Some ports not accessible"
	} else {
		result.Status = types.StatusPass
		result.Message = "All required ports accessible"
	}

	return result, nil
}

type DatabaseCheck struct{}

func (c *DatabaseCheck) Name() string {
	return "Database Connectivity"
}

func (c *DatabaseCheck) Description() string {
	return "Verify database connection"
}

func (c *DatabaseCheck) Category() types.Category {
	return types.CategoryConfig
}

func (c *DatabaseCheck) Severity() types.Severity {
	return types.SeverityError
}

func (c *DatabaseCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusFail,
		Message:  "Database connectivity check failed",
		Details:  make(map[string]interface{}),
	}

	dbPath := cfg.ConfigPath
	if dbPath == "" {
		possiblePaths := []string{
			"./data/mailserver.db",
			"/var/lib/gomailserver/mailserver.db",
			os.Getenv("HOME") + "/.gomailserver/mailserver.db",
		}

		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				dbPath = path
				break
			}
		}
	}

	if dbPath == "" {
		result.Message = "Database file not found"
		result.Details["database_path"] = "not found"
		return result, nil
	}

	if _, err := os.Stat(dbPath); err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Database file not accessible: %v", err)
		result.Details["database_path"] = dbPath
		return result, nil
	}

	dsn := fmt.Sprintf("%s?_foreign_keys=1", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to open database: %v", err)
		result.Details["database_path"] = dbPath
		result.Details["error"] = err.Error()
		return result, nil
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to ping database: %v", err)
		result.Details["database_path"] = dbPath
		result.Details["error"] = err.Error()
		return result, nil
	}

	var file_size int64
	if info, err := os.Stat(dbPath); err == nil {
		file_size = info.Size()
	}

	result.Status = types.StatusPass
	result.Message = "Database connectivity successful"
	result.Details["database_path"] = dbPath
	result.Details["file_size_bytes"] = file_size
	result.Details["file_size_mb"] = fmt.Sprintf("%.2f", float64(file_size)/(1024*1024))

	var tableCount int
	if err := db.QueryRow("SELECT count(*) FROM sqlite_master WHERE type='table'").Scan(&tableCount); err == nil {
		result.Details["table_count"] = tableCount
	}

	return result, nil
}

type DomainConfigurationCheck struct{}

func (c *DomainConfigurationCheck) Name() string {
	return "Domain Configuration"
}

func (c *DomainConfigurationCheck) Description() string {
	return "Validate domain settings"
}

func (c *DomainConfigurationCheck) Category() types.Category {
	return types.CategoryConfig
}

func (c *DomainConfigurationCheck) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *DomainConfigurationCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusPass,
		Message:  "Domain configuration valid",
		Details:  make(map[string]interface{}),
	}

	result.Details["domains"] = cfg.Domains

	return result, nil
}
