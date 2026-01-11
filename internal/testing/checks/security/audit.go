package security

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/btafoya/gomailserver/internal/testing/types"
)

type DKIMConfigAudit struct{}

func (c *DKIMConfigAudit) Name() string {
	return "DKIM Config Audit"
}

func (c *DKIMConfigAudit) Description() string {
	return "Validate DKIM keys, permissions, DNS"
}

func (c *DKIMConfigAudit) Category() types.Category {
	return types.CategorySecurity
}

func (c *DKIMConfigAudit) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *DKIMConfigAudit) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusPass,
		Message:  "DKIM configuration audit",
		Details:  make(map[string]interface{}),
	}

	possibleKeyPaths := []string{
		"/etc/gomailserver/dkim/private.key",
		"./data/dkim/private.key",
		os.Getenv("HOME") + "/.gomailserver/dkim/private.key",
	}

	keyPath := ""
	for _, path := range possibleKeyPaths {
		if _, err := os.Stat(path); err == nil {
			keyPath = path
			break
		}
	}

	if keyPath == "" {
		result.Status = types.StatusWarning
		result.Message = "DKIM key file not found (DKIM disabled or not configured)"
		result.Details["key_file"] = "not found"
		return result, nil
	}

	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to read DKIM key: %v", err)
		result.Details["key_file"] = keyPath
		result.Details["error"] = err.Error()
		return result, nil
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		result.Status = types.StatusFail
		result.Message = "DKIM key file is not valid PEM format"
		result.Details["key_file"] = keyPath
		return result, nil
	}

	rsaKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		rsaKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			result.Status = types.StatusFail
			result.Message = "DKIM key is not valid RSA or Ed25519 private key"
			result.Details["key_file"] = keyPath
			result.Details["key_type"] = block.Type
			return result, nil
		}
	}

	result.Details["key_size_bytes"] = len(keyData)

	info, err := os.Stat(keyPath)
	if err != nil {
		result.Status = types.StatusFail
		result.Message = fmt.Sprintf("Failed to stat DKIM key file: %v", err)
		result.Details["key_file"] = keyPath
		return result, nil
	}

	perm := info.Mode().Perm()
	if perm&0077 != 0 {
		result.Status = types.StatusWarning
		result.Message = "DKIM key file has insecure permissions (should be 0600 or more restrictive)"
		result.Details["key_file"] = keyPath
		result.Details["permissions"] = fmt.Sprintf("%04o", perm)
		return result, nil
	}

	keyDir := filepath.Dir(keyPath)
	if keyDir != "." && keyDir != "/" {
		dirInfo, err := os.Stat(keyDir)
		if err == nil {
			dirPerm := dirInfo.Mode().Perm()
			if dirPerm&0077 != 0 {
				result.Status = types.StatusWarning
				result.Message = "DKIM key directory has insecure permissions"
				result.Details["key_file"] = keyPath
				result.Details["key_dir"] = keyDir
				result.Details["dir_permissions"] = fmt.Sprintf("%04o", dirPerm)
				return result, nil
			}
		}
	}

	result.Status = types.StatusPass
	result.Message = "DKIM configuration valid"
	result.Details["key_file"] = keyPath
	result.Details["key_type"] = block.Type
	result.Details["permissions"] = fmt.Sprintf("%04o", perm)
	result.Details["key_size_bytes"] = len(keyData)

	selector := filepath.Base(keyPath)
	selector = filepath.Base(filepath.Dir(keyPath)) + "." + selector
	result.Details["selector"] = "default"

	return result, nil
}

type DKIMSignatureTest struct{}

func (c *DKIMSignatureTest) Name() string {
	return "DKIM Signature Test"
}

func (c *DKIMSignatureTest) Description() string {
	return "Send and verify DKIM signed message"
}

func (c *DKIMSignatureTest) Category() types.Category {
	return types.CategorySecurity
}

func (c *DKIMSignatureTest) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *DKIMSignatureTest) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "DKIM signature test skipped (requires mail server connection)",
		Details:  make(map[string]interface{}),
	}

	result.Details["requires_server"] = true
	result.Details["test_message"] = "Send test message and verify DKIM signature"

	return result, nil
}

type SPFPolicyCheck struct{}

func (c *SPFPolicyCheck) Name() string {
	return "SPF Policy Check"
}

func (c *SPFPolicyCheck) Description() string {
	return "Validate SPF DNS records and policy"
}

func (c *SPFPolicyCheck) Category() types.Category {
	return types.CategorySecurity
}

func (c *SPFPolicyCheck) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *SPFPolicyCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "SPF policy check skipped (requires DNS lookup)",
		Details:  make(map[string]interface{}),
	}

	if len(cfg.Domains) == 0 {
		result.Status = types.StatusWarning
		result.Message = "No domains configured, skipping SPF check"
		result.Details["domains"] = []string{}
		return result, nil
	}

	domain := cfg.Domains[0]
	result.Details["test_domain"] = domain
	result.Details["expected_record"] = fmt.Sprintf("_spf.%s TXT record", domain)
	result.Details["lookup_method"] = "DNS TXT record query"

	return result, nil
}

type DMARCPolicyCheck struct{}

func (c *DMARCPolicyCheck) Name() string {
	return "DMARC Policy Check"
}

func (c *DMARCPolicyCheck) Description() string {
	return "Validate DMARC DNS records and policy"
}

func (c *DMARCPolicyCheck) Category() types.Category {
	return types.CategorySecurity
}

func (c *DMARCPolicyCheck) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *DMARCPolicyCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "DMARC policy check skipped (requires DNS lookup)",
		Details:  make(map[string]interface{}),
	}

	if len(cfg.Domains) == 0 {
		result.Status = types.StatusWarning
		result.Message = "No domains configured, skipping DMARC check"
		result.Details["domains"] = []string{}
		return result, nil
	}

	domain := cfg.Domains[0]
	result.Details["test_domain"] = domain
	result.Details["expected_record"] = fmt.Sprintf("_dmarc.%s TXT record", domain)
	result.Details["lookup_method"] = "DNS TXT record query"
	result.Details["valid_policies"] = []string{"none", "quarantine", "reject"}

	return result, nil
}

type SecurityChainCheck struct{}

func (c *SecurityChainCheck) Name() string {
	return "Security Chain Test"
}

func (c *SecurityChainCheck) Description() string {
	return "Full DKIM+SPF+DMARC verification"
}

func (c *SecurityChainCheck) Category() types.Category {
	return types.CategorySecurity
}

func (c *SecurityChainCheck) Severity() types.Severity {
	return types.SeverityWarning
}

func (c *SecurityChainCheck) Run(ctx context.Context, cfg *types.ServerConfig) (*types.CheckResult, error) {
	result := &types.CheckResult{
		Check:    c.Name(),
		Category: c.Category(),
		Severity: c.Severity(),
		Status:   types.StatusSkip,
		Message:  "Security chain test skipped (requires mail server connection and DNS lookup)",
		Details:  make(map[string]interface{}),
	}

	result.Details["requires_server"] = true
	result.Details["requires_dns"] = true
	result.Details["test_chain"] = []string{
		"Send DKIM-signed message",
		"Verify DKIM signature",
		"Check SPF validation",
		"Check DMARC validation",
	}
	result.Details["expected_outcome"] = "All three security mechanisms pass"

	return result, nil
}
