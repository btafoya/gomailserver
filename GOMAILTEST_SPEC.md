# GoMailTest - Production Mail Server Verification Tool

**Project**: gomailserver
**Tool**: gomailtest
**Date**: 2026-01-08
**Status**: Implementation Specification
**Version**: 1.0

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Requirements Overview](#requirements-overview)
3. [Architecture](#architecture)
4. [CLI Interface](#cli-interface)
5. [Configuration & Profiles](#configuration--profiles)
6. [Check System](#check-system)
7. [Reporting & Output](#reporting--output)
8. [Production Safety](#production-safety)
9. [Implementation Phases](#implementation-phases)
10. [File Structure](#file-structure)
11. [Usage Examples](#usage-examples)

---

## Executive Summary

**gomailtest** is a standalone CLI tool for verifying gomailserver installations in production environments. It performs comprehensive testing of mail server functionality, configuration validation, and security chain verification without requiring CI/CD integration or development test frameworks.

### Key Differentiators

- **Production-focused**: Designed for live server verification, not development testing
- **Standalone binary**: Separate `gomailtest` command, maximum code reuse from gomailserver
- **Dual operation modes**: Local (read gomailserver.conf) and remote (profile-based) testing
- **Comprehensive safety**: Test message marking, auto-cleanup, dry-run mode, rate limiting
- **Full security audit**: DKIM/SPF/DMARC configuration analysis and validation
- **Flexible output**: Console (quiet/summary/verbose), HTML reports, JSON export

---

## Requirements Overview

### Priority 1: Core Mail Flow (Must Have)

- SMTP connectivity and authentication testing
- IMAP connectivity and retrieval testing
- End-to-end mail flow verification (send → store → retrieve)
- Message integrity validation

### Priority 2: Configuration Validation (Must Have)

- Parse and validate gomailserver.conf syntax
- TLS certificate validation (existence, expiration)
- Port availability checks
- Database connectivity verification
- DNS record validation

### Priority 3: Security Chain Verification (Important)

**DKIM Configuration Audit:**
- Private key existence and permissions (0600)
- Key size validation (recommend 2048+ bits)
- DNS TXT record verification (selector._domainkey.domain)
- Public/private key matching

**SPF Policy Validation:**
- DNS TXT record existence
- Policy syntax validation
- IP/mechanism matching
- Policy strictness analysis (-all vs ~all vs +all)

**DMARC Policy Analysis:**
- DNS TXT record existence (_dmarc.domain)
- Policy strength validation (reject > quarantine > none)
- Alignment configuration (relaxed/strict for DKIM/SPF)
- Reporting address configuration (rua/ruf)

### Priority 4: Performance Diagnostics (Nice to Have)

- Connection timing measurements
- Throughput testing
- Queue processing speed
- Storage performance metrics

### Priority 5: Error Scenarios (Future)

- Invalid credential handling
- Network timeout behavior
- Malformed message handling
- Chaos testing capabilities

---

## Architecture

### Component Overview

```
┌────────────────────────────────────────────────────────┐
│                    gomailtest                        │
├────────────────────────────────────────────────────────┤
│                                                        │
│  ┌──────────────┐    ┌──────────────┐                │
│  │  CLI Parser  │───▶│  Config      │                │
│  │  (cobra)     │    │  Loader      │                │
│  └──────────────┘    └──────────────┘                │
│                            │                          │
│                            ▼                          │
│                     ┌──────────────┐                 │
│                     │  Verifier    │                 │
│                     │  Orchestrator│                 │
│                     └──────────────┘                 │
│                            │                          │
│              ┌─────────────┼─────────────┐           │
│              ▼             ▼             ▼           │
│         ┌────────┐   ┌─────────┐  ┌──────────┐     │
│         │ Config │   │MailFlow │  │ Security │     │
│         │ Checks │   │ Checks  │  │ Checks   │     │
│         └────────┘   └─────────┘  └──────────┘     │
│              │             │             │           │
│              └─────────────┴─────────────┘           │
│                            │                          │
│                            ▼                          │
│                     ┌──────────────┐                 │
│                     │  Reporter    │                 │
│                     │  (Console/   │                 │
│                     │   HTML/JSON) │                 │
│                     └──────────────┘                 │
│                                                        │
└────────────────────────────────────────────────────────┘

Code Reuse from gomailserver:
├─ internal/config (configuration parsing)
├─ internal/security/dkim (DKIM verification)
├─ internal/security/spf (SPF checking)
├─ internal/security/dmarc (DMARC validation)
├─ internal/smtp (SMTP client operations)
└─ internal/imap (IMAP client operations)
```

### Execution Flow

```
1. Parse CLI arguments
   ├─ Load configuration (file or profile)
   ├─ Apply CLI flag overrides
   └─ Validate configuration

2. Initialize verifier
   ├─ Register all checks based on config
   ├─ Set up reporters (console, HTML, JSON)
   └─ Configure safety settings

3. Run checks sequentially
   ├─ Config validation checks (errors fail immediately)
   ├─ Mail flow checks (skip if dry-run)
   └─ Security audit checks (warnings reported)

4. Generate reports
   ├─ Console output (quiet/summary/verbose)
   ├─ Save HTML report (if --report-html)
   └─ Save JSON export (if --report-json)

5. Exit with appropriate code
   ├─ 0: All checks passed (or warnings-only with --warnings-ok)
   └─ 1: One or more error-level checks failed
```

---

## CLI Interface

### Command Structure

```
gomailtest
├── verify           # Run all verification checks
├── test             # Run specific check categories
│   ├── config      # Configuration validation only
│   ├── smtp        # SMTP connectivity only
│   ├── imap        # IMAP connectivity only
│   ├── mailflow    # End-to-end mail flow only
│   └── security    # Security audit only
├── profile          # Manage profiles
│   ├── list        # List available profiles
│   ├── show <name> # Show profile details
│   └── validate    # Validate profile syntax
└── version          # Show version information
```

### Primary Commands

#### `gomailtest verify`

Run all verification checks.

**Flags:**

```
--config <path>       Path to gomailserver.conf (local mode)
--profile <name>      Profile name to use (remote mode)
--dry-run             Non-invasive checks only (no test messages)
--quiet, -q           Quiet mode (exit code only, no output)
--verbose, -v         Verbose output with detailed logging
--report-html <path>  Save HTML report to file
--report-json <path>  Save JSON report to file
--warnings-ok         Don't fail on warnings (exit 0 even with warnings)
--rate-limit <n>      Rate limit (operations per second, default: 10)
--no-cleanup          Don't auto-delete test messages
```

**Examples:**

```bash
# Local verification
gomailtest verify --config /etc/gomailserver/gomailserver.conf

# Remote verification with profile
gomailtest verify --profile production --report-html report.html

# Dry-run mode (non-invasive)
gomailtest verify --dry-run --profile staging

# Verbose output with JSON export
gomailtest verify --verbose --report-json results.json
```

#### `gomailtest test <category>`

Run specific check categories.

**Examples:**

```bash
# Only configuration checks
gomailtest test config --config /etc/gomailserver/gomailserver.conf

# Only SMTP connectivity
gomailtest test smtp --profile production

# Only security audit
gomailtest test security --verbose
```

### Auto-Discovery

When no `--config` or `--profile` is specified, gomailtest automatically:

1. Checks for gomailserver running locally
2. Attempts to read `/etc/gomailserver/gomailserver.conf`
3. Falls back to `./gomailserver.conf`
4. Errors if no configuration found

```bash
# Auto-discovery mode
gomailtest verify
# Equivalent to:
# gomailtest verify --config /etc/gomailserver/gomailserver.conf
```

---

## Configuration & Profiles

### Local Mode (Configuration File)

When testing a local gomailserver instance, use `--config`:

```bash
gomailtest verify --config /etc/gomailserver/gomailserver.conf
```

The tool will:
- Parse the gomailserver configuration
- Extract SMTP/IMAP ports, domains, TLS settings
- Use local credentials if available
- Prompt for test account credentials if needed

**Override with CLI flags:**

```bash
gomailtest verify \
  --config /etc/gomailserver/gomailserver.conf \
  --smtp-port 2587 \
  --imap-port 2993 \
  --test-user monitor@example.com
```

### Remote Mode (Profiles)

For remote testing, create profile files:

**Profile Location:**
- User profiles: `~/.gomailserver/profiles/<name>.yaml`
- System profiles: `/etc/gomailserver/profiles/<name>.yaml`

**Profile Format:**

```yaml
# ~/.gomailserver/profiles/production.yaml
name: production
smtp_host: mail.example.com
smtp_port: 587
imap_host: mail.example.com
imap_port: 993
domains:
  - example.com
  - mail.example.com
test_user: healthcheck@example.com
password_env: PROD_MAIL_PASSWORD  # Read from environment variable

# Optional settings
options:
  tls: true
  starttls: true
  timeout: 30s
```

**Usage:**

```bash
export PROD_MAIL_PASSWORD="your_test_password"
gomailtest verify --profile production
```

### Profile Management Commands

```bash
# List available profiles
gomailtest profile list

# Output:
# Available profiles:
#   production     /home/user/.gomailserver/profiles/production.yaml
#   staging        /home/user/.gomailserver/profiles/staging.yaml
#   local          /etc/gomailserver/profiles/local.yaml

# Show profile details
gomailtest profile show production

# Output:
# Profile: production
#   SMTP: mail.example.com:587
#   IMAP: mail.example.com:993
#   Domains: example.com, mail.example.com
#   Test User: healthcheck@example.com

# Validate profile syntax
gomailtest profile validate production
```

### Configuration Priority

1. CLI flags (highest priority)
2. Profile configuration
3. Config file (gomailserver.conf)
4. Default values (lowest priority)

**Example:**

```bash
gomailtest verify \
  --profile production \
  --smtp-port 2587  # Overrides profile's smtp_port
```

---

## Check System

### Check Interface

All checks implement this interface:

```go
type Check interface {
    Name() string
    Description() string
    Category() Category
    Severity() Severity
    Run(ctx context.Context, cfg *verifier.Config) *Result
}

type Result struct {
    Check    string                 // Check name
    Status   Status                 // Pass, Fail, Warning, Skip
    Severity Severity               // Error, Warning, Info
    Message  string                 // Human-readable result
    Details  map[string]interface{} // Structured details
    Duration time.Duration          // Execution time
    Error    error                  // Error (if failed)
}
```

### Check Categories

#### Category: Config (Priority 2)

| Check Name | Severity | Description |
|------------|----------|-------------|
| Config Syntax | Error | Parse and validate gomailserver.conf |
| TLS Certificates | Error | Validate cert existence and expiration |
| Port Availability | Error | Check SMTP/IMAP ports accessible |
| Database Connectivity | Error | Verify database connection |
| Domain Configuration | Warning | Validate domain settings |

#### Category: MailFlow (Priority 1)

| Check Name | Severity | Description |
|------------|----------|-------------|
| SMTP Connectivity | Error | Connect to SMTP server |
| SMTP Authentication | Error | Authenticate with test credentials |
| IMAP Connectivity | Error | Connect to IMAP server |
| IMAP Authentication | Error | Authenticate with test credentials |
| End-to-End Mail Flow | Error | Send → Store → Retrieve test message |
| Message Integrity | Error | Verify message content matches |

#### Category: Security (Priority 3)

| Check Name | Severity | Description |
|------------|----------|-------------|
| DKIM Config Audit | Warning | Validate DKIM keys, permissions, DNS |
| DKIM Signature Test | Warning | Send and verify DKIM signed message |
| SPF Policy Check | Warning | Validate SPF DNS records and policy |
| DMARC Policy Check | Warning | Validate DMARC DNS records and policy |
| Security Chain Test | Warning | Full DKIM+SPF+DMARC verification |

#### Category: Performance (Priority 4)

| Check Name | Severity | Description |
|------------|----------|-------------|
| SMTP Response Time | Info | Measure SMTP connection time |
| IMAP Response Time | Info | Measure IMAP connection time |
| Delivery Latency | Info | Measure end-to-end delivery time |
| Throughput Test | Info | Measure message throughput |

### Check Registry

Checks are registered at startup:

```go
func InitializeChecks(registry *Registry) {
    // Config checks (Priority 2)
    registry.Register(&ConfigSyntaxCheck{})
    registry.Register(&TLSCertificateCheck{})
    registry.Register(&PortAvailabilityCheck{})
    registry.Register(&DatabaseCheck{})

    // Mail flow checks (Priority 1)
    registry.Register(&SMTPConnectivityCheck{})
    registry.Register(&SMTPAuthenticationCheck{})
    registry.Register(&IMAPConnectivityCheck{})
    registry.Register(&IMAPAuthenticationCheck{})
    registry.Register(&MailFlowEndToEndCheck{})

    // Security checks (Priority 3)
    registry.Register(&DKIMConfigAudit{})
    registry.Register(&DKIMSignatureTest{})
    registry.Register(&SPFPolicyCheck{})
    registry.Register(&DMARCPolicyCheck{})
    registry.Register(&SecurityChainCheck{})

    // Performance checks (Priority 4)
    registry.Register(&SMTPResponseTimeCheck{})
    registry.Register(&IMAPResponseTimeCheck{})
}
```

### Severity Handling

**Error Severity:**
- **Always fails verification**
- Exit code 1
- Critical issues that prevent mail server operation
- Examples: Config parse error, TLS cert expired, SMTP unreachable

**Warning Severity:**
- **Fails verification UNLESS --warnings-ok**
- Exit code 1 (default) or 0 (with --warnings-ok)
- Important issues that should be addressed
- Examples: DKIM key < 2048 bits, SPF too permissive, DMARC policy=none

**Info Severity:**
- **Never fails verification**
- Always exit code 0 (unless errors present)
- Informational metrics and recommendations
- Examples: Performance metrics, optimization suggestions

---

## Reporting & Output

### Console Output Modes

#### Quiet Mode (`--quiet`)

Only exit code, no output:

```bash
gomailtest verify --quiet
echo $?  # 0 = pass, 1 = fail
```

#### Summary Mode (default)

Concise pass/fail summary:

```
Running verification for mail.example.com...

Configuration Validation:
  ✓ Config Syntax                    [  12ms]
  ✓ TLS Certificates                 [  45ms]
  ✓ Port Availability                [  23ms]
  ✓ Database Connectivity            [ 156ms]

Core Mail Flow:
  ✓ SMTP Connectivity                [  34ms]
  ✓ SMTP Authentication              [  67ms]
  ✓ IMAP Connectivity                [  28ms]
  ✓ IMAP Authentication              [  52ms]
  ✓ End-to-End Mail Flow             [2345ms]

Security Chain:
  ✓ DKIM Config Audit                [ 234ms]
  ⚠ SPF Policy Check                 [  89ms]
    Warning: SPF policy too permissive for example.com: v=spf1 ~all
  ✓ DMARC Policy Check               [  76ms]

Results: 11 passed, 1 warning, 0 failed
Duration: 3.161s

⚠ VERIFICATION PASSED WITH WARNINGS

To view details: gomailtest verify --verbose --report-html report.html
```

#### Verbose Mode (`--verbose`)

Detailed output with check details:

```
[2026-01-08 14:23:01] Starting verification...
[2026-01-08 14:23:01] Configuration: /etc/gomailserver/gomailserver.conf
[2026-01-08 14:23:01] Domains: example.com, mail.example.com
[2026-01-08 14:23:01] Test user: healthcheck@example.com

[2026-01-08 14:23:01] Running Config Syntax check...
[2026-01-08 14:23:01]   ✓ Config parsed successfully
[2026-01-08 14:23:01]   ✓ Found 2 domains: example.com, mail.example.com
[2026-01-08 14:23:01]   ✓ SMTP ports: [25, 587, 465]
[2026-01-08 14:23:01]   ✓ IMAP ports: [143, 993]
[2026-01-08 14:23:01]   Duration: 12ms

[2026-01-08 14:23:02] Running TLS Certificates check...
[2026-01-08 14:23:02]   ✓ Certificate: /etc/gomailserver/certs/mail.example.com.crt
[2026-01-08 14:23:02]   ✓ Expires: 2026-12-15 (341 days remaining)
[2026-01-08 14:23:02]   Duration: 45ms

[2026-01-08 14:23:05] Running DKIM Config Audit...
[2026-01-08 14:23:05]   Domain: example.com
[2026-01-08 14:23:05]     ✓ Private key exists: /etc/gomailserver/dkim/example.com.key
[2026-01-08 14:23:05]     ✓ Key permissions: 0600
[2026-01-08 14:23:05]     ✓ Key size: 2048 bits
[2026-01-08 14:23:05]     ✓ Selector: mail
[2026-01-08 14:23:05]     ✓ DNS record found: mail._domainkey.example.com
[2026-01-08 14:23:05]   Duration: 234ms

[2026-01-08 14:23:05] Running SPF Policy Check...
[2026-01-08 14:23:05]   Domain: example.com
[2026-01-08 14:23:05]     ✓ SPF record found: v=spf1 mx ~all
[2026-01-08 14:23:05]     ⚠ Warning: Policy uses ~all (softfail), recommend -all (fail)
[2026-01-08 14:23:05]   Duration: 89ms

...

[2026-01-08 14:23:08] Verification complete
[2026-01-08 14:23:08] Results: 11 passed, 1 warning, 0 failed
[2026-01-08 14:23:08] Total duration: 3.161s

⚠ VERIFICATION PASSED WITH WARNINGS
```

### HTML Report

Generated with `--report-html report.html`:

```html
<!DOCTYPE html>
<html>
<head>
    <title>GoMailTest Verification Report</title>
    <style>
        /* Modern, clean styling with status colors */
        .pass { color: #4CAF50; }
        .warning { color: #FF9800; }
        .fail { color: #F44336; }
        .check-category { margin: 20px 0; }
        .check-result { padding: 10px; margin: 5px 0; }
        .details { font-family: monospace; font-size: 0.9em; }
    </style>
</head>
<body>
    <h1>GoMailTest Verification Report</h1>

    <div class="summary">
        <h2>Summary</h2>
        <p><strong>Server:</strong> mail.example.com</p>
        <p><strong>Time:</strong> 2026-01-08 14:23:01</p>
        <p><strong>Duration:</strong> 3.161s</p>
        <p><strong>Results:</strong> 11 passed, 1 warning, 0 failed</p>
        <p class="warning"><strong>Status:</strong> PASSED WITH WARNINGS</p>
    </div>

    <div class="check-category">
        <h2>Configuration Validation</h2>

        <div class="check-result pass">
            <h3>✓ Config Syntax (12ms)</h3>
            <p>Configuration syntax valid</p>
            <div class="details">
                <pre>domains: 2
smtp_ports: [25, 587, 465]
imap_ports: [143, 993]</pre>
            </div>
        </div>

        <!-- More checks... -->
    </div>

    <!-- Timeline visualization -->
    <div class="timeline">
        <h2>Execution Timeline</h2>
        <svg><!-- Timeline chart --></svg>
    </div>
</body>
</html>
```

### JSON Export

Generated with `--report-json results.json`:

```json
{
  "timestamp": "2026-01-08T14:23:01Z",
  "server": "mail.example.com",
  "config_file": "/etc/gomailserver/gomailserver.conf",
  "duration_ms": 3161,
  "summary": {
    "total": 12,
    "passed": 11,
    "warnings": 1,
    "failed": 0
  },
  "status": "passed_with_warnings",
  "exit_code": 1,
  "checks": [
    {
      "name": "Config Syntax",
      "category": "config",
      "severity": "error",
      "status": "pass",
      "message": "Configuration syntax valid",
      "duration_ms": 12,
      "details": {
        "domains": 2,
        "smtp_ports": [25, 587, 465],
        "imap_ports": [143, 993]
      }
    },
    {
      "name": "SPF Policy Check",
      "category": "security",
      "severity": "warning",
      "status": "warning",
      "message": "SPF policy too permissive for example.com: v=spf1 ~all",
      "duration_ms": 89,
      "details": {
        "example.com": "v=spf1 mx ~all"
      }
    }
  ]
}
```

---

## Production Safety

### Test Message Identification

All test messages include special markers:

**Headers:**
```
X-GoMailTest: true
X-GoMailTest-ID: verify-2026-01-08-142301-abc123
X-GoMailTest-Timestamp: 2026-01-08T14:23:01Z
Subject: [MAILTEST] Health Check 2026-01-08-142301
```

**Benefits:**
- Easy identification in logs and mail queues
- Simple filtering rules for monitoring
- Clear separation from production traffic
- Enables automated cleanup

### Auto-Cleanup

With `--auto-cleanup` (default), test messages are automatically deleted after verification:

```bash
gomailtest verify --auto-cleanup  # Default behavior
gomailtest verify --no-cleanup    # Keep test messages
```

**Cleanup Process:**
1. Test message sent and retrieved successfully
2. After verification complete, connect to IMAP
3. Search for messages with X-GoMailTest-ID header
4. Delete matching messages
5. Expunge mailbox

### Dry-Run Mode

Non-invasive verification without sending test messages:

```bash
gomailtest verify --dry-run
```

**Dry-run checks:**
- ✓ Configuration validation
- ✓ TLS certificate checks
- ✓ DNS record verification
- ✓ SMTP/IMAP connectivity tests
- ✗ End-to-end mail flow (skipped)
- ✓ Security configuration audit

**Use cases:**
- Initial validation before full testing
- Monitoring scripts that run frequently
- Pre-deployment configuration checks

### Rate Limiting

Prevent overwhelming production servers:

```bash
gomailtest verify --rate-limit 5  # Max 5 operations/second
```

**Rate-limited operations:**
- SMTP connections
- IMAP connections
- DNS lookups
- Test message sends

**Default:** 10 operations/second

### Dedicated Test Accounts

Recommended practice: Create dedicated test accounts

**Example configuration:**

```yaml
# gomailserver.conf
accounts:
  - email: healthcheck@example.com
    password: $HEALTHCHECK_PASSWORD
    quota: 10MB
    auto_cleanup: true
    mailbox_retention: 1h
```

**Benefits:**
- Isolated from production users
- Can be monitored separately
- Easy to clean up
- No impact on user quotas

### Read-Only Mode (Future)

Check recent mail flow from logs/metrics without sending new messages:

```bash
gomailtest verify --read-only
```

**Read-only checks:**
- Parse recent SMTP logs for successful deliveries
- Check IMAP activity logs for recent retrieval
- Verify queue processing metrics
- Analyze recent DKIM/SPF/DMARC results

---

## Implementation Phases

### Phase 1: CLI Foundation & Configuration

**Goal:** Build standalone binary with config/profile support

**Tasks:**
1. Create `cmd/gomailtest` CLI structure with cobra
2. Implement config file parser (reuse gomailserver config package)
3. Implement profile management (load, list, validate)
4. Create check interface and registry
5. Implement basic console reporter

**Deliverables:**
- `gomailtest` binary that can be built
- `gomailtest verify --help` works
- `gomailtest profile list` works
- Configuration loading from file and profile

**Estimated Effort:** 2-3 days

---

### Phase 2: Core Checks Implementation

**Goal:** Implement all Priority 1 and 2 checks

**Config Checks:**
- ConfigSyntaxCheck
- TLSCertificateCheck
- PortAvailabilityCheck
- DatabaseCheck

**Mail Flow Checks:**
- SMTPConnectivityCheck
- SMTPAuthenticationCheck
- IMAPConnectivityCheck
- IMAPAuthenticationCheck
- MailFlowEndToEndCheck

**Deliverables:**
- All config and mail flow checks implemented
- Test coverage for each check
- Working end-to-end verification

**Estimated Effort:** 3-4 days

---

### Phase 3: Security Audit Implementation

**Goal:** Implement all Priority 3 security checks

**Security Checks:**
- DKIMConfigAudit (key permissions, size, DNS)
- DKIMSignatureTest (send and verify signed message)
- SPFPolicyCheck (DNS records, policy validation)
- DMARCPolicyCheck (DNS records, policy analysis)
- SecurityChainCheck (full DKIM+SPF+DMARC test)

**Deliverables:**
- All security checks implemented
- Full DKIM/SPF/DMARC verification working
- Test coverage for security checks

**Estimated Effort:** 3-4 days

---

### Phase 4: Reporting & Output

**Goal:** Implement all output formats

**Reporters:**
- Console reporter (quiet/summary/verbose)
- HTML reporter (with timeline visualization)
- JSON reporter (machine-readable export)

**Deliverables:**
- All three output formats working
- Beautiful HTML reports
- JSON suitable for automation

**Estimated Effort:** 2-3 days

---

### Phase 5: Production Safety Features

**Goal:** Implement all safety mechanisms

**Safety Features:**
- Test message identification (X-GoMailTest headers)
- Auto-cleanup functionality
- Dry-run mode
- Rate limiting
- Dedicated test account support

**Deliverables:**
- Production-safe testing
- Cleanup working reliably
- Rate limiting functional

**Estimated Effort:** 2 days

---

### Phase 6: Polish & Documentation

**Goal:** Production-ready release

**Tasks:**
- Comprehensive testing (integration tests)
- Documentation (README, usage guide)
- Build scripts and distribution
- Example profiles
- Error message polish

**Deliverables:**
- Stable v1.0 release
- Complete documentation
- Installation instructions

**Estimated Effort:** 2-3 days

---

## File Structure

```
gomailserver/
├── cmd/
│   └── gomailtest/
│       ├── main.go                  # CLI entry point
│       ├── commands/
│       │   ├── verify.go            # Verify command
│       │   ├── test.go              # Test command
│       │   └── profile.go           # Profile management
│       └── version.go               # Version info
│
├── internal/testing/
│   ├── verifier/
│   │   ├── verifier.go             # Main orchestration
│   │   ├── config.go               # Configuration types
│   │   ├── profile.go              # Profile management
│   │   └── result.go               # Result types
│   │
│   ├── checks/
│   │   ├── check.go                # Check interface
│   │   ├── registry.go             # Check registry
│   │   │
│   │   ├── config/
│   │   │   ├── syntax.go           # Config syntax check
│   │   │   ├── tls.go              # TLS certificate check
│   │   │   ├── ports.go            # Port availability check
│   │   │   └── database.go         # Database connectivity check
│   │   │
│   │   ├── mailflow/
│   │   │   ├── smtp_connectivity.go    # SMTP connectivity check
│   │   │   ├── smtp_auth.go            # SMTP authentication check
│   │   │   ├── imap_connectivity.go    # IMAP connectivity check
│   │   │   ├── imap_auth.go            # IMAP authentication check
│   │   │   └── end_to_end.go           # End-to-end mail flow check
│   │   │
│   │   ├── security/
│   │   │   ├── dkim_audit.go           # DKIM config audit
│   │   │   ├── dkim_signature.go       # DKIM signature test
│   │   │   ├── spf.go                  # SPF policy check
│   │   │   ├── dmarc.go                # DMARC policy check
│   │   │   └── chain.go                # Security chain test
│   │   │
│   │   └── performance/
│   │       ├── smtp_timing.go          # SMTP response time
│   │       ├── imap_timing.go          # IMAP response time
│   │       └── throughput.go           # Throughput test
│   │
│   ├── reporters/
│   │   ├── reporter.go             # Reporter interface
│   │   ├── console.go              # Console reporter
│   │   ├── html.go                 # HTML reporter
│   │   └── json.go                 # JSON reporter
│   │
│   └── safety/
│       ├── cleanup.go              # Test message cleanup
│       ├── ratelimit.go            # Rate limiting
│       └── markers.go              # Test message markers
│
├── testdata/
│   └── profiles/
│       ├── example-local.yaml      # Example local profile
│       ├── example-remote.yaml     # Example remote profile
│       └── example-production.yaml # Example production profile
│
└── docs/
    ├── gomailtest-guide.md        # User guide
    ├── profile-format.md          # Profile format specification
    └── check-reference.md         # Check reference documentation
```

---

## Usage Examples

### Example 1: Local Development Testing

```bash
# Test local gomailserver instance
gomailtest verify --config ./gomailserver.conf --verbose

# Output:
# Running verification for localhost...
# ✓ Configuration valid
# ✓ SMTP reachable (localhost:587)
# ✓ IMAP reachable (localhost:993)
# ✓ Mail flow working
# ✓ DKIM configured
# ⚠ SPF not configured (development mode)
# PASSED WITH WARNINGS (6.2s)
```

### Example 2: Production Server Verification

```bash
# Create production profile
cat > ~/.gomailserver/profiles/production.yaml <<EOF
name: production
smtp_host: mail.example.com
smtp_port: 587
imap_host: mail.example.com
imap_port: 993
domains:
  - example.com
test_user: healthcheck@example.com
password_env: PROD_MAIL_PASSWORD
EOF

# Set password
export PROD_MAIL_PASSWORD="your_secure_password"

# Run verification
gomailtest verify --profile production --report-html report.html

# Output:
# Running verification for mail.example.com...
# ✓ Configuration valid
# ✓ TLS certificates valid (expires in 341 days)
# ✓ SMTP connectivity (mail.example.com:587)
# ✓ IMAP connectivity (mail.example.com:993)
# ✓ Mail flow working (2.1s)
# ✓ DKIM configured correctly (2048-bit key)
# ✓ SPF policy valid (v=spf1 mx -all)
# ✓ DMARC policy valid (p=reject)
# ✓ Security chain verified
# PASSED (8.7s)
#
# Report saved: report.html
```

### Example 3: Monitoring Script Integration

```bash
#!/bin/bash
# /usr/local/bin/check-mailserver.sh
#
# Monitoring script for cron/systemd timer

REPORT_DIR="/var/log/gomailtest"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

# Run verification with JSON output
gomailtest verify \
  --profile production \
  --quiet \
  --report-json "$REPORT_DIR/report-$TIMESTAMP.json"

EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
  echo "Mail server verification: PASSED"
  # Send success metric to monitoring system
  curl -X POST https://monitoring.example.com/metrics \
    -d "mailserver.verification=1"
else
  echo "Mail server verification: FAILED"
  # Send alert to monitoring system
  curl -X POST https://monitoring.example.com/alerts \
    -d "mailserver.verification=0"
  # Send alert email
  mail -s "Mail Server Verification Failed" ops@example.com < "$REPORT_DIR/report-$TIMESTAMP.json"
fi

# Cleanup old reports (keep last 30 days)
find "$REPORT_DIR" -name "report-*.json" -mtime +30 -delete

exit $EXIT_CODE
```

**Crontab:**
```
# Check mail server every 15 minutes
*/15 * * * * /usr/local/bin/check-mailserver.sh
```

### Example 4: Pre-Deployment Validation

```bash
#!/bin/bash
# Pre-deployment configuration check

echo "Validating mail server configuration before deployment..."

# Dry-run verification (no test messages)
gomailtest verify \
  --config ./deploy/gomailserver.conf \
  --dry-run \
  --verbose

if [ $? -ne 0 ]; then
  echo "❌ Configuration validation failed!"
  echo "Fix errors before deploying."
  exit 1
fi

echo "✓ Configuration valid"
echo "Safe to deploy."
exit 0
```

### Example 5: Security Audit Only

```bash
# Run only security checks
gomailtest test security \
  --profile production \
  --report-html security-audit.html

# Output:
# Running security audit for mail.example.com...
#
# DKIM Configuration:
#   example.com:
#     ✓ Private key exists and readable
#     ✓ Key permissions secure (0600)
#     ✓ Key size: 2048 bits (recommended)
#     ✓ DNS record found: mail._domainkey.example.com
#     ✓ Public key matches private key
#
# SPF Configuration:
#   example.com:
#     ✓ SPF record found: v=spf1 mx ip4:203.0.113.10 -all
#     ✓ Policy includes server IP: 203.0.113.10
#     ✓ Policy strict: -all (reject unauthorized)
#
# DMARC Configuration:
#   example.com:
#     ✓ DMARC record found: v=DMARC1; p=reject; rua=mailto:dmarc@example.com
#     ✓ Policy: reject (strict)
#     ✓ DKIM alignment: relaxed
#     ✓ SPF alignment: relaxed
#     ✓ Aggregate reporting configured
#
# Security Chain Test:
#   ✓ Test message sent with DKIM signature
#   ✓ DKIM signature valid
#   ✓ SPF check passed
#   ✓ DMARC alignment verified
#
# SECURITY AUDIT PASSED (4.2s)
```

### Example 6: Specific Check Testing

```bash
# Test only SMTP connectivity
gomailtest test smtp --profile production

# Test only configuration
gomailtest test config --config /etc/gomailserver/gomailserver.conf

# Test only mail flow
gomailtest test mailflow --profile staging --verbose
```

---

## Dependencies

### Go Modules

```go
module github.com/btafoya/gomailserver

go 1.23.5

require (
    // CLI framework
    github.com/spf13/cobra v1.8.0

    // Configuration
    gopkg.in/yaml.v3 v3.0.1

    // Reused from gomailserver
    github.com/emersion/go-smtp v0.24.0
    github.com/emersion/go-imap/v2 v2.0.0-beta.7

    // Testing
    github.com/stretchr/testify v1.9.0
)
```

### External Dependencies

None required. All functionality built with:
- Go standard library
- Existing gomailserver packages
- Minimal external dependencies (cobra, yaml)

---

## Next Steps

1. **Review this specification** with stakeholders
2. **Create GitHub issue** tracking implementation phases
3. **Set up project structure** (cmd/gomailtest/, internal/testing/)
4. **Begin Phase 1 implementation** (CLI foundation)
5. **Iterate on feedback** from initial testing

---

**Document Status**: ✅ Complete
**Version**: 1.0
**Last Updated**: 2026-01-08
**Author**: btafoya
