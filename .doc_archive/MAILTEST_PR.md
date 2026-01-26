# GoMailTest - Production Mail Server Verification Tool

**Project**: gomailserver
**Date**: 2026-01-08
**Status**: Implementation Specification

---

## Table of Contents

1. [Overview](#overview)
2. [Research Findings](#research-findings)
3. [Tool Architecture](#tool-architecture)
4. [Implementation Plan](#implementation-plan)
5. [File Structure](#file-structure)
6. [CLI Interface](#cli-interface)
7. [Configuration & Profiles](#configuration--profiles)
8. [Production Safety](#production-safety)
9. [Next Steps](#next-steps)

---

## Overview

This document specifies **gomailtest**, a standalone CLI tool for verifying gomailserver installations in production environments. The tool performs comprehensive end-to-end testing of mail server functionality, configuration validation, and security chain verification without requiring CI/CD integration or development test frameworks.

### Key Requirements

- ✅ **Standalone CLI tool** (`gomailtest` binary) for production verification
- ✅ **Local and remote testing** - verify servers from any location
- ✅ **Core mail flow testing** - SMTP send/receive, IMAP retrieval, storage verification
- ✅ **Configuration validation** - parse configs, verify TLS, check DNS records
- ✅ **Security chain audit** - comprehensive DKIM/SPF/DMARC verification
- ✅ **Multiple output formats** - console (quiet/summary/verbose), HTML reports, JSON export
- ✅ **Production safety** - test accounts, message marking, auto-cleanup, dry-run mode
- ✅ **Profile management** - simple profiles for different environments

---

## Research Findings

Based on extensive research of open-source repositories, testing tools, and best practices:

### Existing Tools Analyzed

| Tool | Type | Language | Key Features | Relevance |
|--------|-------|-----------|----------------|-------------|
| **smtptest** (k1LoW/smtptest) | SMTP Test Server | Go | In-process SMTP server for testing, auth support, message retrieval | ✅ High - Go-native, perfect for unit tests |
| **go-smtp-mock** (mocktools) | SMTP Mock Server | Go | Configurable behavior, error injection, multi-recipient support, response delays | ✅ High - Perfect for chaos testing |
| **MailHog** | Fake SMTP + Web UI | Go | Email capture, web viewer, REST API, WebSocket updates | ✅ Medium - Web UI useful for visual testing |
| **Inbucket** | SMTP/POP3 + REST | Go | Multi-protocol, swaks integration, MIME handling, message retention | ✅ Medium - Good for protocol testing |
| **Mailpit** | Fake SMTP + Web UI | Go | Modern MailHog alternative, faster, actively maintained | ✅ Medium - Good replacement for MailHog |
| **Swaks** | CLI Testing Tool | Perl | RFC compliance testing, TLS support, custom headers, attachments | ✅ High - Script automation, RFC verification |
| **smtpbench** | Load Testing | Python | Multi-threaded load testing, MX failover, detailed JSON logging | ✅ High - Performance benchmarking |
| **MailSlurp** | Email Testing Service | Multiple SDKs | CI/CD integration, email waiting, verification code extraction | ✅ High - Perfect for automated tests |
| **Imitate Email** | Email Sandbox | Multiple | Full SMTP/IMAP/POP, programmatic API, WebSocket support | ✅ Medium - Good for E2E testing |
| **Greenmail** | Integration Testing Server | Java | JUnit integration, SSL/TLS, rule-based filtering | ⚠️ Low - Java-based, less relevant for Go project |

### Key Insights

1. **Go-Native Tools Preferred**: smtptest and go-smtp-mock are ideal for Go-based testing
2. **Mock-Based Testing Works**: Both Inbucket and go-smtp-mock demonstrate success with mock dependencies
3. **Error Injection Critical**: go-smtp-mock's configurable behavior enables chaos testing
4. **CI/CD Integration Common**: MailSlurp and similar services show CI/CD patterns
5. **Load Testing Essential**: smtpbench provides production-ready load testing
6. **Your MailSandbox Exists**: You already have btafoya/mailsandbox - leverage it!
7. **Security Testing Missing**: Most tools don't include DKIM/SPF/DMARC verification - major opportunity

### Your Current Testing State

**Strengths:**
- ✅ Good unit tests with mocks (smtp/backend_test.go, imap/backend_test.go)
- ✅ Service layer tests (message_service_test.go, queue_service_test.go)
- ✅ Reputation integration tests (integration_test.go)
- ✅ Security modules implemented (DKIM, SPF, DMARC)
- ✅ Reusable packages for SMTP, IMAP, configuration parsing

**Gaps for Production Verification:**
- ❌ No standalone tool for verifying live server deployments
- ❌ No end-to-end verification from external perspective (remote testing)
- ❌ No comprehensive security chain audit tool
- ❌ No production-safe testing framework with cleanup
- ❌ No profile-based testing for multiple environments

---

## Tool Architecture

### GoMailTest - Production Verification CLI

**"Standalone CLI tool for comprehensive mail server verification"**

#### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                   GoMailTest CLI Tool                     │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Verification Phases (Priority Order):                      │
│                                                             │
│  1. Configuration Validation                                │
│     ├─ Parse gomailserver.conf                             │
│     ├─ Validate TLS certificates                           │
│     ├─ Check port availability                             │
│     └─ Verify database connectivity                        │
│                                                             │
│  2. Core Mail Flow Testing                                  │
│     ├─ SMTP send (local/remote)                            │
│     ├─ Message storage verification                         │
│     ├─ IMAP retrieval                                       │
│     └─ Content integrity check                             │
│                                                             │
│  3. Security Chain Audit                                    │
│     ├─ DKIM: config audit, DNS verification, crypto check  │
│     ├─ SPF: policy validation, IP matching                 │
│     └─ DMARC: policy analysis, alignment verification      │
│                                                             │
│  4. Performance Diagnostics (optional)                      │
│     ├─ Connection timing                                    │
│     ├─ Throughput measurement                              │
│     └─ Queue processing speed                              │
│                                                             │
│  Output Modes:                                              │
│  • Console: quiet/summary/verbose                          │
│  • HTML: detailed audit reports                            │
│  • JSON: machine-readable for automation                   │
│  • Structured logs: for aggregation systems                │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

#### Key Components

**1. Verifier Interface**

```go
type Check interface {
    Name() string
    Run(ctx context.Context, cfg *Config) *CheckResult
    Severity() Severity  // Error, Warning, Info
}

type CheckResult struct {
    Name      string
    Status    Status  // Pass, Fail, Warning
    Severity  Severity
    Message   string
    Details   map[string]interface{}
    Duration  time.Duration
    Error     error
}
```

**2. Verification Orchestrator**

```go
type Verifier struct {
    Config    *Config
    Profile   *Profile
    Checks    []Check
    Reporter  Reporter

    // Reused gomailserver components
    DKIMVerifier *dkim.Verifier
    SPFChecker   *spf.Checker
    DMARCChecker *dmarc.Checker
}

func (v *Verifier) Run(ctx context.Context) (*Report, error) {
    results := []CheckResult{}

    for _, check := range v.Checks {
        result := check.Run(ctx, v.Config)
        results = append(results, result)

        // Stop on critical errors
        if result.Severity == SeverityError && result.Status == StatusFail {
            return nil, fmt.Errorf("critical check failed: %s", check.Name())
        }
    }

    return v.Reporter.Generate(results), nil
}
```

**3. Check Implementations**

```go
// Configuration checks
type ConfigSyntaxCheck struct{}
type TLSCertificateCheck struct{}
type PortAvailabilityCheck struct{}
type DatabaseCheck struct{}

// Mail flow checks
type SMTPConnectivityCheck struct{}
type SMTPAuthenticationCheck struct{}
type IMAPConnectivityCheck struct{}
type MailFlowEndToEndCheck struct{}

// Security checks
type DKIMConfigAudit struct{}
type DKIMSignatureCheck struct{}
type SPFPolicyCheck struct{}
type DMARCPolicyCheck struct{}
type SecurityChainCheck struct{}
```

**4. Profile & Configuration**

```go
type Config struct {
    // Local mode: read from gomailserver.conf
    ConfigFile string

    // Remote mode: connection parameters
    SMTPHost string
    SMTPPort int
    IMAPHost string
    IMAPPort int

    // Test credentials
    TestUser     string
    TestPassword string

    // Safety settings
    DryRun       bool
    AutoCleanup  bool
    RateLimit    int

    // Output control
    Verbose      bool
    Quiet        bool
    OutputFormat OutputFormat
}

type Profile struct {
    Name         string
    SMTPHost     string
    SMTPPort     int
    IMAPHost     string
    IMAPPort     int
    Domains      []string
    TestUser     string
    PasswordEnv  string
}
```

---

## Implementation Plan

### Phase 1: CLI Foundation & Configuration (Priority 1)

**Goal:** Build standalone `gomailtest` binary with config parsing and profile support

#### 1.1 Create CLI Structure

**Files to Create:**

```bash
cmd/gomailtest/
├── main.go              # CLI entry point
├── commands/
│   ├── verify.go       # Main verification command
│   ├── test.go         # Individual test commands
│   └── profile.go      # Profile management
└── version.go          # Version info

internal/testing/
├── verifier/
│   ├── verifier.go     # Main orchestration
│   ├── config.go       # Configuration parsing
│   ├── profile.go      # Profile management
│   └── result.go       # Result types
├── checks/
│   ├── check.go        # Check interface
│   ├── config/         # Configuration validation checks
│   ├── mailflow/       # Mail flow checks
│   └── security/       # Security audit checks
└── reporters/
    ├── console.go      # Console output
    ├── html.go         # HTML reports
    └── json.go         # JSON export
```

**Key Code:**

```go
// cmd/gomailtest/main.go
package main

import (
    "github.com/spf13/cobra"
    "github.com/btafoya/gomailserver/internal/testing/verifier"
)

func main() {
    rootCmd := &cobra.Command{
        Use:   "gomailtest",
        Short: "Production mail server verification tool",
    }

    rootCmd.AddCommand(
        newVerifyCommand(),
        newTestCommand(),
        newProfileCommand(),
    )

    rootCmd.Execute()
}

func newVerifyCommand() *cobra.Command {
    var (
        configFile  string
        profile     string
        quiet       bool
        verbose     bool
        dryRun      bool
        reportHTML  string
        reportJSON  string
    )

    cmd := &cobra.Command{
        Use:   "verify",
        Short: "Run all verification checks",
        RunE: func(cmd *cobra.Command, args []string) error {
            cfg := verifier.NewConfig()

            // Load config/profile
            if configFile != "" {
                cfg.LoadFromFile(configFile)
            }
            if profile != "" {
                cfg.LoadProfile(profile)
            }

            // Apply flags
            cfg.Quiet = quiet
            cfg.Verbose = verbose
            cfg.DryRun = dryRun

            v := verifier.New(cfg)
            report, err := v.RunAll(cmd.Context())
            if err != nil {
                return err
            }

            // Output reports
            if reportHTML != "" {
                report.SaveHTML(reportHTML)
            }
            if reportJSON != "" {
                report.SaveJSON(reportJSON)
            }

            return report.ExitError()
        },
    }

    cmd.Flags().StringVar(&configFile, "config", "", "Path to gomailserver.conf")
    cmd.Flags().StringVar(&profile, "profile", "", "Profile name to use")
    cmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Quiet mode (exit code only)")
    cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
    cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Non-invasive checks only")
    cmd.Flags().StringVar(&reportHTML, "report-html", "", "Save HTML report")
    cmd.Flags().StringVar(&reportJSON, "report-json", "", "Save JSON report")

    return cmd
}
```

#### 1.2 Configuration & Profile System

```go
// internal/testing/verifier/config.go
type Config struct {
    // Source
    ConfigFile string
    Profile    *Profile

    // Connection (local or remote)
    SMTPHost string
    SMTPPort int
    IMAPHost string
    IMAPPort int

    // Credentials
    TestUser     string
    TestPassword string

    // Domains to verify
    Domains []string

    // Safety
    DryRun      bool
    AutoCleanup bool
    RateLimit   int

    // Output
    Quiet        bool
    Verbose      bool
    OutputFormat OutputFormat
}

type Profile struct {
    Name        string            `yaml:"name"`
    SMTPHost    string            `yaml:"smtp_host"`
    SMTPPort    int               `yaml:"smtp_port"`
    IMAPHost    string            `yaml:"imap_host"`
    IMAPPort    int               `yaml:"imap_port"`
    Domains     []string          `yaml:"domains"`
    TestUser    string            `yaml:"test_user"`
    PasswordEnv string            `yaml:"password_env"`
    Options     map[string]string `yaml:"options"`
}

func LoadProfile(name string) (*Profile, error) {
    // Check ~/.gomailserver/profiles/<name>.yaml
    // Fall back to /etc/gomailserver/profiles/<name>.yaml
    path := filepath.Join(os.Getenv("HOME"), ".gomailserver", "profiles", name+".yaml")

    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var profile Profile
    if err := yaml.Unmarshal(data, &profile); err != nil {
        return nil, err
    }

    return &profile, nil
}
```

#### 1.3 Check Interface & Registry

```go
// internal/testing/checks/check.go
type Check interface {
    Name() string
    Description() string
    Category() Category
    Severity() Severity
    Run(ctx context.Context, cfg *verifier.Config) *Result
}

type Category string
const (
    CategoryConfig   Category = "config"
    CategoryMailFlow Category = "mailflow"
    CategorySecurity Category = "security"
    CategoryPerf     Category = "performance"
)

type Severity string
const (
    SeverityError   Severity = "error"   // Fails verification
    SeverityWarning Severity = "warning" // Reports but doesn't fail
    SeverityInfo    Severity = "info"    // Informational only
)

type Result struct {
    Check    string
    Status   Status
    Severity Severity
    Message  string
    Details  map[string]interface{}
    Duration time.Duration
    Error    error
}

type Registry struct {
    checks map[Category][]Check
}

func (r *Registry) Register(check Check) {
    cat := check.Category()
    r.checks[cat] = append(r.checks[cat], check)
}

func (r *Registry) GetByCategory(cat Category) []Check {
    return r.checks[cat]
}
```

---

### Phase 2: Core Checks Implementation (Priority 2)

**Goal:** Implement configuration, mail flow, and security checks

#### 2.1 Configuration Validation Checks

```go
// internal/testing/checks/config/syntax.go
type ConfigSyntaxCheck struct{}

func (c *ConfigSyntaxCheck) Name() string { return "Config Syntax" }
func (c *ConfigSyntaxCheck) Category() Category { return CategoryConfig }
func (c *ConfigSyntaxCheck) Severity() Severity { return SeverityError }

func (c *ConfigSyntaxCheck) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    // Parse gomailserver.conf using existing parser
    serverCfg, err := config.LoadFile(cfg.ConfigFile)
    if err != nil {
        return &Result{
            Status:   StatusFail,
            Message:  fmt.Sprintf("Failed to parse config: %v", err),
            Duration: time.Since(start),
            Error:    err,
        }
    }

    details := map[string]interface{}{
        "domains":    len(serverCfg.Domains),
        "smtp_ports": serverCfg.SMTP.Ports,
        "imap_ports": serverCfg.IMAP.Ports,
    }

    return &Result{
        Status:   StatusPass,
        Message:  "Configuration syntax valid",
        Details:  details,
        Duration: time.Since(start),
    }
}

// internal/testing/checks/config/tls.go
type TLSCertificateCheck struct{}

func (c *TLSCertificateCheck) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    serverCfg, _ := config.LoadFile(cfg.ConfigFile)

    // Check each TLS certificate
    for _, certPath := range serverCfg.TLS.Certificates {
        cert, err := tls.LoadX509KeyPair(certPath.Cert, certPath.Key)
        if err != nil {
            return &Result{
                Status:   StatusFail,
                Severity: SeverityError,
                Message:  fmt.Sprintf("Invalid TLS certificate: %s", certPath.Cert),
                Error:    err,
                Duration: time.Since(start),
            }
        }

        // Parse certificate to check expiration
        x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
        if err != nil {
            continue
        }

        daysUntilExpiry := int(time.Until(x509Cert.NotAfter).Hours() / 24)

        if daysUntilExpiry < 0 {
            return &Result{
                Status:   StatusFail,
                Severity: SeverityError,
                Message:  "TLS certificate expired",
                Details:  map[string]interface{}{"expired_days_ago": -daysUntilExpiry},
                Duration: time.Since(start),
            }
        } else if daysUntilExpiry < 30 {
            return &Result{
                Status:   StatusWarning,
                Severity: SeverityWarning,
                Message:  "TLS certificate expires soon",
                Details:  map[string]interface{}{"days_until_expiry": daysUntilExpiry},
                Duration: time.Since(start),
            }
        }
    }

    return &Result{
        Status:   StatusPass,
        Message:  "TLS certificates valid",
        Duration: time.Since(start),
    }
}
```

#### 2.2 Mail Flow Checks

```go
// internal/testing/checks/mailflow/smtp_connectivity.go
type SMTPConnectivityCheck struct{}

func (c *SMTPConnectivityCheck) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    addr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)

    conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
    if err != nil {
        return &Result{
            Status:   StatusFail,
            Severity: SeverityError,
            Message:  fmt.Sprintf("Cannot connect to SMTP server: %s", addr),
            Error:    err,
            Duration: time.Since(start),
        }
    }
    defer conn.Close()

    // Try SMTP handshake
    smtpConn := smtp.NewClientConn(conn)
    _, err = smtpConn.Text.ReadLine()  // Read greeting
    if err != nil {
        return &Result{
            Status:   StatusFail,
            Message:  "SMTP handshake failed",
            Error:    err,
            Duration: time.Since(start),
        }
    }

    return &Result{
        Status:   StatusPass,
        Message:  fmt.Sprintf("SMTP server reachable: %s", addr),
        Duration: time.Since(start),
    }
}

// internal/testing/checks/mailflow/end_to_end.go
type MailFlowEndToEndCheck struct{}

func (c *MailFlowEndToEndCheck) Run(ctx context.Context, cfg *Config) *Result {
    if cfg.DryRun {
        return &Result{
            Status:  StatusSkip,
            Message: "Skipped (dry-run mode)",
        }
    }

    start := time.Now()

    // Generate unique test message
    testID := generateTestID()
    msg := buildTestMessage(testID, cfg.TestUser)

    // Phase 1: Send via SMTP
    smtpAddr := fmt.Sprintf("%s:%d", cfg.SMTPHost, cfg.SMTPPort)
    err := sendSMTP(smtpAddr, cfg.TestUser, cfg.TestPassword, msg)
    if err != nil {
        return &Result{
            Status:   StatusFail,
            Message:  "Failed to send test message",
            Error:    err,
            Duration: time.Since(start),
        }
    }

    // Phase 2: Wait for delivery
    time.Sleep(2 * time.Second)

    // Phase 3: Retrieve via IMAP
    imapAddr := fmt.Sprintf("%s:%d", cfg.IMAPHost, cfg.IMAPPort)
    retrieved, err := fetchIMAP(imapAddr, cfg.TestUser, cfg.TestPassword, testID)
    if err != nil {
        return &Result{
            Status:   StatusFail,
            Message:  "Failed to retrieve test message via IMAP",
            Error:    err,
            Duration: time.Since(start),
        }
    }

    // Phase 4: Cleanup
    if cfg.AutoCleanup {
        deleteIMAP(imapAddr, cfg.TestUser, cfg.TestPassword, testID)
    }

    return &Result{
        Status:   StatusPass,
        Message:  "Mail flow working: SMTP → Storage → IMAP",
        Details: map[string]interface{}{
            "test_id":   testID,
            "retrieved": retrieved,
        },
        Duration: time.Since(start),
    }
}
```

#### 2.3 Security Audit Checks

```go
// internal/testing/checks/security/dkim.go
type DKIMConfigAudit struct{}

func (c *DKIMConfigAudit) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    serverCfg, _ := config.LoadFile(cfg.ConfigFile)

    issues := []string{}
    details := map[string]interface{}{}

    for _, domain := range serverCfg.Domains {
        dkimCfg := serverCfg.DKIM[domain]

        // Check private key exists
        if _, err := os.Stat(dkimCfg.PrivateKeyPath); err != nil {
            issues = append(issues, fmt.Sprintf("DKIM private key not found: %s", domain))
            continue
        }

        // Check key permissions
        info, _ := os.Stat(dkimCfg.PrivateKeyPath)
        if info.Mode().Perm() != 0600 {
            issues = append(issues, fmt.Sprintf("DKIM key permissions too open for %s: %o (should be 0600)",
                domain, info.Mode().Perm()))
        }

        // Load and check key size
        keyData, _ := os.ReadFile(dkimCfg.PrivateKeyPath)
        block, _ := pem.Decode(keyData)
        privKey, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

        keySize := privKey.N.BitLen()
        details[domain] = map[string]interface{}{
            "selector": dkimCfg.Selector,
            "key_size": keySize,
        }

        if keySize < 2048 {
            issues = append(issues, fmt.Sprintf("DKIM key size too small for %s: %d bits (recommend 2048+)",
                domain, keySize))
        }

        // Check DNS record
        dnsName := fmt.Sprintf("%s._domainkey.%s", dkimCfg.Selector, domain)
        txtRecords, err := net.LookupTXT(dnsName)
        if err != nil || len(txtRecords) == 0 {
            issues = append(issues, fmt.Sprintf("DKIM DNS record not found: %s", dnsName))
        }
    }

    if len(issues) > 0 {
        return &Result{
            Status:   StatusWarning,
            Severity: SeverityWarning,
            Message:  fmt.Sprintf("%d DKIM issues found", len(issues)),
            Details:  details,
            Duration: time.Since(start),
        }
    }

    return &Result{
        Status:   StatusPass,
        Message:  "DKIM configuration valid",
        Details:  details,
        Duration: time.Since(start),
    }
}

// internal/testing/checks/security/spf.go
type SPFPolicyCheck struct{}

func (c *SPFPolicyCheck) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    serverCfg, _ := config.LoadFile(cfg.ConfigFile)

    issues := []string{}
    details := map[string]interface{}{}

    for _, domain := range serverCfg.Domains {
        txtRecords, err := net.LookupTXT(domain)
        if err != nil {
            issues = append(issues, fmt.Sprintf("DNS lookup failed for %s", domain))
            continue
        }

        var spfRecord string
        for _, record := range txtRecords {
            if strings.HasPrefix(record, "v=spf1") {
                spfRecord = record
                break
            }
        }

        if spfRecord == "" {
            issues = append(issues, fmt.Sprintf("SPF record not found for %s", domain))
            continue
        }

        details[domain] = spfRecord

        // Check policy strictness
        if strings.Contains(spfRecord, "~all") || strings.Contains(spfRecord, "+all") {
            issues = append(issues, fmt.Sprintf("SPF policy too permissive for %s: %s", domain, spfRecord))
        }
    }

    if len(issues) > 0 {
        return &Result{
            Status:   StatusWarning,
            Message:  fmt.Sprintf("%d SPF issues found", len(issues)),
            Details:  details,
            Duration: time.Since(start),
        }
    }

    return &Result{
        Status:   StatusPass,
        Message:  "SPF policies valid",
        Details:  details,
        Duration: time.Since(start),
    }
}

// internal/testing/checks/security/dmarc.go
type DMARCPolicyCheck struct{}

func (c *DMARCPolicyCheck) Run(ctx context.Context, cfg *Config) *Result {
    start := time.Now()

    serverCfg, _ := config.LoadFile(cfg.ConfigFile)

    issues := []string{}
    details := map[string]interface{}{}

    for _, domain := range serverCfg.Domains {
        dmarcDomain := "_dmarc." + domain
        txtRecords, err := net.LookupTXT(dmarcDomain)
        if err != nil || len(txtRecords) == 0 {
            issues = append(issues, fmt.Sprintf("DMARC record not found for %s", domain))
            continue
        }

        dmarcRecord := txtRecords[0]
        details[domain] = dmarcRecord

        // Parse policy
        if !strings.Contains(dmarcRecord, "p=quarantine") && !strings.Contains(dmarcRecord, "p=reject") {
            issues = append(issues, fmt.Sprintf("DMARC policy weak for %s: should be quarantine or reject", domain))
        }

        // Check for reporting
        if !strings.Contains(dmarcRecord, "rua=") {
            issues = append(issues, fmt.Sprintf("DMARC aggregate reporting not configured for %s", domain))
        }
    }

    if len(issues) > 0 {
        return &Result{
            Status:   StatusWarning,
            Message:  fmt.Sprintf("%d DMARC issues found", len(issues)),
            Details:  details,
            Duration: time.Since(start),
        }
    }

    return &Result{
        Status:   StatusPass,
        Message:  "DMARC policies valid",
        Details:  details,
        Duration: time.Since(start),
    }
}
```

---

### Phase 3: Reporting & Output Formats (Priority 3)

#### 3.1 Integrate go-smtp-mock for Error Injection

```go
// internal/testing/chaos/smtp_chaos_test.go
import smtpmock "github.com/mocktools/go-smtp-mock/v2"

type ChaosTest struct {
    Name          string
    ErrorScenario func(*smtpmock.Server)
    ExpectRecovery bool
}

func (t *ChaosTest) Execute(ctx context.Context) error {
    // Configure mock server with error injection
    config := smtpmock.ConfigurationAttr{
        HostAddress:      "127.0.0.1",
        PortNumber:        2525,
        LogToStdout:      false,
        LogServerActivity: true,
        BlacklistedDomains: []string{},
        ResponseDelay:     0,
    }

    mockServer := smtpmock.New(config)
    defer mockServer.Stop()

    // Apply error scenario
    t.ErrorScenario(mockServer)

    // Start mock server
    if err := mockServer.Start(); err != nil {
        return err
    }

    // Try to send email
    client, err := smtp.Dial("127.0.0.1:2525", "", nil)
    if err != nil {
        return err
    }
    defer client.Quit()

    err = client.Mail("test@example.com", []string{"recipient@example.com"},
        []byte("Subject: Test\r\n\r\nTest body"))

    if t.ExpectRecovery {
        // Server should handle error gracefully
        if err == nil {
            return fmt.Errorf("Expected error but got none")
        }
        // Check server logs for proper error handling
        logs := mockServer.GetLogs()
        if !containsGracefulHandling(logs) {
            return fmt.Errorf("Server didn't handle error gracefully")
        }
    } else {
        // Expect complete failure
        if err == nil {
            return fmt.Errorf("Expected failure but message succeeded")
        }
    }

    return nil
}
```

#### 3.2 Load Testing with smtpbench

```bash
# scripts/load/smtp_load_test.sh
#!/bin/bash

# Install smtpbench
pip install smtpbench

# Test configuration
SMTP_HOST="localhost"
SMTP_PORT="587"
USER="test@example.com"
PASS="password"

echo "=== SMTP Load Testing ==="
echo

# Test 1: Basic throughput (100 msgs/sec for 60s)
echo "Test 1: Basic throughput (6,000 messages)..."
smtpbench \
    --host $SMTP_HOST \
    --port $SMTP_PORT \
    --username $USER \
    --password $PASS \
    --rate 100 \
    --duration 60 \
    --output load-test-1.json

# Analyze results
echo "Results:"
jq -r '.summary' load-test-1.json

# Test 2: Sustained load (1,000 msgs/min for 10 min)
echo "Test 2: Sustained load (10,000 messages)..."
smtpbench \
    --host $SMTP_HOST \
    --port $SMTP_PORT \
    --username $USER \
    --password $PASS \
    --rate 16 \
    --duration 600 \
    --output load-test-2.json

echo
echo "=== Load Testing Complete ==="
echo "Reports saved in load-test-*.json"
```

---

### Phase 4: Reporting + CI/CD Integration (Week 4)

#### 4.1 HTML Report with Security Visualization

```go
// internal/testing/reporters/security_html.go
type SecurityReport struct {
    TestID       string
    StartTime    time.Time
    EndTime      time.Time

    // Mail Flow
    SMTPSent     bool
    Queued       bool
    Stored       bool
    IMAPFetched  bool

    // Security Chain
    DKIM         DKIMSummary
    SPF          SPFSummary
    DMARC        DMARCSummary

    // Reputation
    Reputation   ReputationScore

    // Overall
    Status       string  // "pass", "fail", "partial"
    Trace        []TraceEvent
}

type DKIMSummary struct {
    Signed       bool
    Verified     bool
    Selector     string
    KeySize      int
    Algorithm    string
    ValidFrom    time.Time
}

type SPFSummary struct {
    Result       string  // "pass", "fail", "neutral"
    Record       string
    IPRange      string
    Aligned      bool
}

type DMARCSummary struct {
    Result       string
    Policy       string
    Aligned      bool
    Alignment    string  // "spf", "dkim", "both"
    Pct          int
}

func GenerateHTMLReport(report *SecurityReport) (string, error) {
    template := `<!DOCTYPE html>
<html>
<head>
    <title>GoMailServer Security Test Report - {{.TestID}}</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .pass { color: #4CAF50; font-weight: bold; }
        .fail { color: #F44336; font-weight: bold; }
        .warning { color: #FF9800; font-weight: bold; }
        .phase-box {
            border: 1px solid #ddd;
            border-radius: 5px;
            padding: 15px;
            margin: 10px 0;
        }
        .security-box {
            border: 2px solid #333;
            padding: 20px;
            margin: 20px 0;
        }
        .metric { margin: 5px 0; }
        .header { background: #f5f5f5; padding: 10px; font-weight: bold; }
    </style>
</head>
<body>
    <h1>GoMailServer Security Test Report</h1>
    <p>Test ID: <code>{{.TestID}}</code></p>
    <p>Duration: {{.Duration}}</p>

    <div class="security-box">
        <h2>Mail Flow</h2>
        {{if .SMTPSent}}<p>✅ SMTP Send: Message accepted</p>{{end}}
        {{if .Queued}}<p>✅ Queue Processing: Message queued</p>{{end}}
        {{if .Stored}}<p>✅ Storage: Message stored in database</p>{{end}}
        {{if .IMAPFetched}}<p>✅ IMAP Fetch: Message retrieved</p>{{end}}
    </div>

    <div class="security-box">
        <h2>Security Chain</h2>

        <div class="phase-box">
            <h3>DKIM</h3>
            {{if .DKIM.Signed}}<p>✅ Signing: {{.DKIM.Algorithm}}, selector={{.DKIM.Selector}}</p>{{end}}
            {{if .DKIM.Verified}}<p>✅ Verification: Signature valid</p>{{end}}
            <p>Key: <code>{{.DKIM.Selector}}._domainkey.example.com</code></p>
            <p class="metric">Key Size: {{.DKIM.KeySize}} bits</p>
        </div>

        <div class="phase-box">
            <h3>SPF</h3>
            <p>{{if .SPF.Aligned}}✅{{else}}❌{{end}} Result: {{.SPF.Result}}</p>
            <p>Record: <code>{{.SPF.Record}}</code></p>
            {{if .SPF.Aligned}}<p>✅ Alignment: SPF matches MAIL FROM</p>{{end}}
        </div>

        <div class="phase-box">
            <h3>DMARC</h3>
            <p>{{if .DMARC.Aligned}}✅{{else}}❌{{end}} Result: {{.DMARC.Result}}</p>
            <p>Policy: <code>{{.DMARC.Policy}}</code></p>
            <p>Alignment: {{.DMARC.Alignment}}</p>
            <p class="metric">Pct: {{.DMARC.Pct}}</p>
        </div>
    </div>

    <div class="security-box">
        <h2>Reputation Impact</h2>
        <p class="metric">Delivery Count: {{.Reputation.DeliveryCount}}</p>
        <p class="metric">Reputation Score: {{.Reputation.Score}}</p>
        <p class="metric">Status: {{if gt .Reputation.Score 90}}Excellent{{else if gt .Reputation.Score 70}}Good{{else}}Poor{{end}}</p>
    </div>

    <div class="security-box">
        <h2>Execution Timeline</h2>
        <ul>
        {{range .Trace}}
        <li>[{{.Timestamp.Format "15:04:05"}}] {{.Phase}}: {{.Action}}
            <span class="{{.Status}}">{{.Status}}</span> ({{.Duration}})</li>
        {{end}}
        </ul>
    </div>

    <div class="security-box">
        <h2>Overall Status</h2>
        <p class="{{.Status}}">{{if eq .Status "pass"}}✅ ALL TESTS PASSED{{else}}❌ TESTS FAILED{{end}}</p>
    </div>
</body>
</html>`

    // Execute template
    return executeTemplate(template, report)
}
```

#### 4.2 CI/CD Integration with MailSlurp

```go
// tests/ci/email_flow_test.go
package ci

import (
    "context"
    "testing"
    "time"

    "github.com/mailslurp/mailslurp-go"
)

func TestUserSignupEmailFlow(t *testing.T) {
    // Skip if no API key
    apiKey := os.Getenv("MAILSLURP_API_KEY")
    if apiKey == "" {
        t.Skip("MAILSLURP_API_KEY not set")
    }

    // Create MailSlurp client
    client := mailslurp.NewClient(apiKey)

    // Create test inbox
    inbox, err := client.CreateInbox(context.Background())
    if err != nil {
        t.Fatalf("Failed to create inbox: %v", err)
    }
    defer inbox.Close()

    t.Logf("Created test inbox: %s", inbox.EmailAddress)

    // Start GoMailServer (or use existing instance)
    // ...

    // Simulate user signup via API
    signupResponse := sendSignupRequest(t, inbox.EmailAddress)
    if signupResponse.StatusCode != 200 {
        t.Fatalf("Signup failed: %d", signupResponse.StatusCode)
    }

    // Wait for welcome email (timeout 30s)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    welcomeEmail, err := inbox.WaitForEmailBySubject(
        ctx,
        "Welcome to GoMailServer",
        nil,  // no body filter
    )
    if err != nil {
        t.Fatalf("Welcome email not received: %v", err)
    }

    t.Logf("Received welcome email with ID: %s", welcomeEmail.ID)

    // Extract verification link
    link, err := extractLink(welcomeEmail.Body, `/verify/`)
    if err != nil {
        t.Fatalf("Failed to extract verification link: %v", err)
    }

    t.Logf("Extracted verification link: %s", link)

    // Verify account
    verifyResponse := sendVerificationRequest(t, link)
    if verifyResponse.StatusCode != 200 {
        t.Fatalf("Verification failed: %d", verifyResponse.StatusCode)
    }

    // Check for confirmation email
    confirmationEmail, err := inbox.WaitForEmailBySubject(
        ctx,
        "Account Verified",
        nil,
    )
    if err != nil {
        t.Fatalf("Confirmation email not received: %v", err)
    }

    t.Logf("Received confirmation email with ID: %s", confirmationEmail.ID)

    // Test reply flow
    replyBody := "Thanks for the welcome email!"
    err = sendEmailViaGoMailServer(t,
        "user@example.com",
        inbox.EmailAddress,
        "Re: Welcome to GoMailServer",
        replyBody,
    )
    if err != nil {
        t.Fatalf("Failed to send reply: %v", err)
    }

    // Wait for reply in test inbox
    replyEmail, err := inbox.WaitForEmailBySubject(
        ctx,
        "Re: Welcome to GoMailServer",
        nil,
    )
    if err != nil {
        t.Fatalf("Reply not received: %v", err)
    }

    // Verify reply content
    if !strings.Contains(replyEmail.Body, replyBody) {
        t.Errorf("Reply body mismatch: expected '%s', got '%s'",
            replyBody, replyEmail.Body)
    }

    t.Logf("✅ Full email flow test passed")
}
```

#### 4.3 GitHub Actions Workflow

```yaml
# .github/workflows/mail-test.yml
name: Mail Delivery Test

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 6 * * *'  # Daily at 6 AM UTC

jobs:
  basic-delivery:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23.5'

      - name: Install dependencies
        run: |
          go install github.com/k1LoW/smtptest@latest
          go install github.com/mocktools/go-smtp-mock/v2@latest

      - name: Run basic delivery tests
        run: |
          go test -v ./internal/testing/scenarios/... -run TestBasicDelivery
        env:
          GOMAILSERVER_CONFIG: ./testdata/test-config.yaml

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: test-results
          path: test-reports/

  security-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23.5'

      - name: Run DKIM tests
        run: |
          go test -v ./internal/testing/scenarios/... -run TestDKIM
        env:
          GOMAILSERVER_CONFIG: ./testdata/test-config.yaml

      - name: Run SPF tests
        run: |
          go test -v ./internal/testing/scenarios/... -run TestSPF

      - name: Run DMARC tests
        run: |
          go test -v ./internal/testing/scenarios/... -run TestDMARC

      - name: Upload security reports
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: security-reports
          path: test-reports/security/

  integration-tests:
    runs-on: ubuntu-latest
    services:
      gomailserver:
        image: btafoya/gomailserver:latest
        ports:
          - 25:25
          - 587:587
          - 143:143
          - 993:993
          - 8980:8980
        env:
          GOMAILSERVER_CONFIG: /etc/gomailserver/gomailserver.yaml

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5

      - name: Wait for server to start
        run: |
          for i in {1..30}; do
            if nc -z localhost 587; then
              echo "Server is ready"
              break
            fi
            sleep 2
          done

      - name: Run integration tests
        run: |
          go test -v ./tests/integration/...
        env:
          MAILSLURP_API_KEY: ${{ secrets.MAILSLURP_API_KEY }}

      - name: Upload integration reports
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: integration-reports
          path: test-reports/integration/

  load-test:
    runs-on: ubuntu-latest
    if: github.event_name == 'schedule'
    steps:
      - uses: actions/checkout@v4

      - name: Set up Python
        uses: actions/setup-python@v5
        with:
          python-version: '3.x'

      - name: Install smtpbench
        run: pip install smtpbench

      - name: Run load tests
        run: |
          ./scripts/load/smtp_load_test.sh

      - name: Upload load test results
        uses: actions/upload-artifact@v4
        with:
          name: load-test-results
          path: load-test-*.json
```

---

## File Structure

```
gomailserver/
├── internal/testing/
│   ├── runner.go              # Test orchestration
│   ├── tracer.go              # Event collection
│   ├── config.go              # Test configuration
│   ├── reporters/
│   │   ├── html.go           # HTML report generator
│   │   ├── json.go           # JSON export
│   │   └── security_html.go  # Security visualization
│   ├── scenarios/
│   │   ├── base.go           # Base test interface
│   │   ├── basic.go          # Basic mail flow
│   │   ├── thread.go         # Threaded conversation
│   │   ├── multi.go          # Multi-recipient
│   │   ├── large.go          # Large attachment
│   │   ├── dkim.go           # DKIM scenarios
│   │   ├── spf.go            # SPF scenarios
│   │   ├── dmarc.go          # DMARC scenarios
│   │   └── security_chain.go # Full chain test
│   ├── chaos/
│   │   ├── smtp_chaos.go    # SMTP error injection
│   │   └── imap_chaos.go    # IMAP error injection
│   ├── servers/
│   │   ├── smtp_test.go      # Test SMTP server wrapper
│   │   └── imap_test.go      # Test IMAP server wrapper
│   └── security/
│       ├── dkim_tracer.go     # DKIM verification
│       ├── spf_tracer.go      # SPF verification
│       └── dmarc_tracer.go    # DMARC verification
├── cmd/mailtest/
│   └── main.go               # CLI entry point
├── tests/
│   ├── integration/
│   │   ├── mail_flow_test.go  # End-to-end mail flow
│   │   ├── imap_client_test.go  # IMAP operations
│   │   └── security_test.go   # Security chain
│   ├── ci/
│   │   ├── email_flow_test.go   # CI with MailSlurp
│   │   └── signup_flow_test.go  # User signup flows
│   └── load/
│       └── smtp_load_test.sh   # Load testing script
├── scripts/
│   ├── test-scenarios.sh      # Run predefined scenarios
│   ├── dkim-rotation.sh     # DKIM key rotation test
│   ├── dmarc-audit.sh       # DMARC audit
│   └── generate-reports.sh   # Generate test reports
└── testdata/
    ├── test-config.yaml       # Test configuration
    └── messages/            # Test email messages
```

---

## Integration Strategy

### 1. MailSandbox Integration

Leverage your existing MailSandbox tool (btafoya/mailsandbox):

```bash
# Start MailSandbox
mailsandbox start --port 2525

# Configure GoMailServer to relay to MailSandbox
# gomailserver.yaml:
smtp:
  relay_server: localhost:2525

# Run tests
./gomailserver test delivery

# View captured emails
open http://localhost:8080
```

**Benefits:**
- Zero-configuration email capture
- Web UI for viewing received emails
- Postmark API emulation (you already support this!)
- Perfect for development testing

### 2. Use Existing Security Modules

```go
// internal/testing/security/dkim_tracer.go
import (
    "github.com/btafoya/gomailserver/internal/security/dkim"
)

type DKIMTracer struct {
    verifier *dkim.Verifier
    tracer   *TraceCollector
}

func (t *DKIMTracer) Verify(rawMsg []byte) (*DKIMResult, error) {
    trace := t.tracer.Start("dkim_verify")
    defer trace.End()

    // Use your existing DKIM verifier
    result, err := t.verifier.Verify(rawMsg)

    trace.WithDetails(map[string]interface{}{
        "dkim_valid":   result.Valid,
        "selector":      result.Selector,
        "domain":        result.Domain,
        "key_size":      result.KeySize,
        "algorithm":     result.Algorithm,
    })

    return result, err
}
```

### 3. smtptest for In-Process Testing

```go
// Use smtptest for Go-native SMTP testing
import smtptest "github.com/k1LoW/smtptest"

func TestSMTPTransaction(t *testing.T) {
    ts, auth, err := smtptest.NewServerWithAuth()
    if err != nil {
        t.Fatal(err)
    }
    defer ts.Close()

    addr := ts.Addr()
    if err := smtp.SendMail(addr, auth, "sender@example.org",
        []string{"alice@example.net"}, []byte(testMsg)); err != nil {
        t.Fatal(err)
    }

    if len(ts.Messages()) != 1 {
        t.Errorf("got %v\nwant %v", len(ts.Messages()), 1)
    }
}
```

### 4. go-smtp-mock for Chaos Testing

```go
import smtpmock "github.com/mocktools/go-smtp-mock/v2"

func TestErrorInjection(t *testing.T) {
    // Configure mock server with error injection
    configuration := createConfiguration()
    server := newServer(configuration)
    server.SetResponseDelay(30 * time.Second)

    // Test scenario
    // ...
}
```

---

## CLI Usage Examples

### Basic Delivery Test

```bash
# Run basic delivery test
./gomailserver test delivery \
    --sender test@example.com \
    --recipient user@example.com \
    --subject "Test Email" \
    --body "Test message" \
    --report-html ./test-reports/delivery.html

# Output:
Running mail delivery test...
✓ SMTP send successful (45ms)
✓ Message queued (12ms)
✓ Message stored (98ms)
✓ IMAP fetch successful (145ms)
✓ Content verified (12ms)

✓ ALL TESTS PASSED (312ms)

Report: ./test-reports/delivery.html
```

### Security Chain Test

```bash
# Run full security chain test
./gomailserver test security-chain \
    --domain example.com \
    --selector k1 \
    --verify-dkim \
    --verify-spf \
    --verify-dmarc \
    --report-html ./test-reports/security.html

# Output:
Running security chain test...
Checking DNS records...
✓ SPF record found
✓ DMARC record found
✓ DKIM key found

Running message tests...
✓ Message 1: DKIM pass, SPF pass, DMARC pass
✓ Message 2: DKIM pass, SPF pass, DMARC pass
✓ Message 3: DKIM pass, SPF pass, DMARC pass

✓ SECURITY CHAIN PASSED

Report: ./test-reports/security.html
```

### DKIM Rotation Test

```bash
# Run DKIM rotation test
./gomailserver test dkim-rotate \
    --domain example.com \
    --old-selector k1 \
    --new-selector k2 \
    --report-html ./test-reports/dkim-rotation.html

# Output:
Running DKIM rotation test...
✓ Old selector (k1) signing works
✓ Rotated to new selector (k2)
✓ New selector (k2) signing works
✓ Overlap period: Old signatures still valid

✓ DKIM ROTATION PASSED

Report: ./test-reports/dkim-rotation.html
```

### DMARC Report Test

```bash
# Run DMARC report test
./gomailserver test dmarc-report \
    --report-file ./test-report.xml \
    --report-email dmarc@example.com

# Output:
Running DMARC report test...
✓ Report submitted successfully
✓ Report parsed (50 events)
✓ Events recorded in database
✓ Reputation score updated (95 → 96)

✓ DMARC REPORT TEST PASSED

View report details: http://localhost:8980/admin/reputation/dmarc
```

### Security Audit

```bash
# Run comprehensive security audit
./gomailserver test security-audit \
    --domain example.com

# Output:
Running security audit...
Checking DNS records...
✓ SPF record found
✓ DMARC record found
✓ DKIM key found (k1)
✓ DKIM key found (k2) - rotation in progress

Testing DKIM signatures...
✓ k1 selector: Valid (2048-bit RSA)
✓ k2 selector: Valid (2048-bit RSA)

Testing DMARC processing...
✓ Report submission works
✓ Report parsing works
✓ Reputation updates correctly

✓ SECURITY AUDIT PASSED

Report: ./test-reports/security-audit-20260108.html
```

### Interactive Mode

```bash
# Start interactive test runner
./gomailserver test interactive

MailFlow Tracer v1.0
====================

[✓] SMTP server listening on :2525
[✓] IMAP server listening on :9143
[✓] Database connected to ./test.db

Select test scenario:
  [1] Basic delivery (send → retrieve)
  [2] Threaded conversation (send → reply → retrieve)
  [3] Multi-recipient test
  [4] Large attachment test
  [5] TLS handshake test
  [6] Error injection suite
  [7] DKIM rotation test
  [8] DMARC report processing
  [9] Security chain audit
  [10] Load testing

> 1

Running scenario: Basic delivery
=================================
Phase 1: SMTP Send          [██████████] 0ms
Phase 2: Queue Verification   [██████████] 45ms ✓
Phase 3: Storage            [██████████] 120ms ✓
Phase 4: IMAP Fetch         [██████████] 200ms ✓
Phase 5: Content Verify      [██████████] 235ms ✓

✓ Test PASSED in 235ms

View detailed trace? [y/N] > y
Opening browser at http://localhost:8080/trace/abc-123...
```

---

## Next Steps

### Week 1: Foundation
1. ✅ Create `internal/testing` package
2. ✅ Implement `TestRunner` framework
3. ✅ Add `TraceCollector`
4. ✅ Implement `BasicDeliveryTest`
5. ✅ Create HTML reporter
6. ✅ Add CLI command

### Week 2: Security
1. ✅ Implement `DKIMRotationTest`
2. ✅ Implement `DMARCReportTest`
3. ✅ Implement `SecurityChainTest`
4. ✅ Add security visualization
5. ✅ Integrate with your DKIM/SPF/DMARC modules

### Week 3: Advanced
1. ✅ Implement error injection with go-smtp-mock
2. ✅ Add load testing with smtpbench
3. ✅ Implement `ThreadedConversationTest`
4. ✅ Implement `MultiRecipientTest`
5. ✅ Implement `LargeAttachmentTest`

### Week 4: Integration
1. ✅ Integrate MailSandbox
2. ✅ Add CI/CD with GitHub Actions
3. ✅ Add MailSlurp integration
4. ✅ Document testing approach
5. ✅ Create test scenarios library

---

## Questions for Implementation

1. **Starting point**: Should I start with Phase 1 (Foundation) or a specific component?
2. **Priority**: Which test scenario is most critical for your current needs?
3. **Integration**: Should this integrate tightly with MailSandbox from the start?
4. **CI/CD**: Should I set up GitHub Actions immediately or after the framework is complete?
5. **Deliverables**: Do you want just the framework, or also a set of predefined test scenarios?
6. **Testing**: Should I implement this as a separate `mailtest` command or integrate into existing `gomailserver` CLI?

---

## Appendix: Dependencies

### External Tools to Integrate

| Tool | Purpose | Install Command |
|-------|---------|-----------------|
| smtptest | Go SMTP test server | `go install github.com/k1LoW/smtptest@latest` |
| go-smtp-mock | Configurable mock SMTP server | `go install github.com/mocktools/go-smtp-mock/v2@latest` |
| smtpbench | SMTP load testing | `pip install smtpbench` |
| MailSlurp SDK | CI/CD email testing | `go get github.com/mailslurp/mailslurp-go` |

### Go Modules

```go
module github.com/btafoya/gomailserver

go 1.23.5

require (
    github.com/k1LoW/smtptest v0.10.2
    github.com/mocktools/go-smtp-mock/v2 v2.0.0
    github.com/mailslurp/mailslurp-go v1.0.0
    github.com/emersion/go-smtp v0.24.0
    github.com/emersion/go-imap v2.0.0-beta.7
)
```

---

**Document Status**: ✅ Complete
**Version**: 1.0
**Last Updated**: 2026-01-08
