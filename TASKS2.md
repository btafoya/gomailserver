# Phase 2: Security Foundation (Weeks 5-7)

**Status**: Not Started
**Priority**: MVP - Required
**Estimated Duration**: 2-3 weeks
**Dependencies**: Phase 1 (Core Mail Server)

---

## Overview

Implement email security standards (DKIM, SPF, DMARC), anti-virus/anti-spam integration (ClamAV, SpamAssassin), greylisting, rate limiting, and authentication security (2FA, brute force protection).

---

## 2.1 DKIM (DomainKeys Identified Mail) [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| DK-001 | RSA-2048/4096 key generation | [ ] | F-002 | `crypto/rsa` |
| DK-002 | Ed25519 key generation (RFC 8463) | [ ] | F-002 | `crypto/ed25519` |
| DK-003 | DKIM signing for outbound mail | [ ] | DK-001, S-002 |
| DK-004 | DKIM verification for inbound mail | [ ] | DK-001, S-003 |
| DK-005 | Multiple selector support per domain | [ ] | DK-001 |
| DK-006 | Key rotation mechanism | [ ] | DK-005 |
| DK-007 | DKIM keys storage (per domain) | [ ] | F-022, U-003 |

### DK-001: Key Generation

```go
// internal/security/dkim/keygen.go
package dkim

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
)

type KeyPair struct {
    PrivateKey string
    PublicKey  string
    Selector   string
}

func GenerateRSAKeyPair(bits int) (*KeyPair, error) {
    privateKey, err := rsa.GenerateKey(rand.Reader, bits)
    if err != nil {
        return nil, err
    }

    privateKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "RSA PRIVATE KEY",
        Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
    })

    publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
    if err != nil {
        return nil, err
    }

    publicKeyPEM := pem.EncodeToMemory(&pem.Block{
        Type:  "PUBLIC KEY",
        Bytes: publicKeyDER,
    })

    return &KeyPair{
        PrivateKey: string(privateKeyPEM),
        PublicKey:  string(publicKeyPEM),
        Selector:   generateSelector(),
    }, nil
}

func generateSelector() string {
    return fmt.Sprintf("s%d", time.Now().Unix())
}

// Generate DNS TXT record content
func (kp *KeyPair) DNSRecord() string {
    // Extract base64 public key
    block, _ := pem.Decode([]byte(kp.PublicKey))
    pubKeyB64 := base64.StdEncoding.EncodeToString(block.Bytes)
    return fmt.Sprintf("v=DKIM1; k=rsa; p=%s", pubKeyB64)
}
```

### DK-003: DKIM Signing

```go
// internal/security/dkim/signer.go
package dkim

import (
    "github.com/emersion/go-msgauth/dkim"
)

type Signer struct {
    domainService *service.DomainService
}

func (s *Signer) Sign(domain string, message []byte) ([]byte, error) {
    domainCfg, err := s.domainService.GetDKIMConfig(domain)
    if err != nil {
        return message, nil // Return unsigned if no DKIM config
    }

    options := &dkim.SignOptions{
        Domain:   domain,
        Selector: domainCfg.Selector,
        Signer:   domainCfg.PrivateKey,
        HeaderKeys: []string{
            "From", "To", "Subject", "Date", "Message-ID",
            "MIME-Version", "Content-Type",
        },
    }

    r := bytes.NewReader(message)
    var buf bytes.Buffer
    if err := dkim.Sign(&buf, r, options); err != nil {
        return nil, err
    }

    return buf.Bytes(), nil
}
```

### DK-004: DKIM Verification

```go
// internal/security/dkim/verifier.go
package dkim

import (
    "github.com/emersion/go-msgauth/dkim"
)

type VerificationResult struct {
    Valid       bool
    Domain      string
    Selector    string
    Error       string
    HeaderField string
}

func (v *Verifier) Verify(message []byte) (*VerificationResult, error) {
    r := bytes.NewReader(message)
    verifications, err := dkim.Verify(r)
    if err != nil {
        return &VerificationResult{Valid: false, Error: err.Error()}, nil
    }

    for _, v := range verifications {
        if v.Err == nil {
            return &VerificationResult{
                Valid:       true,
                Domain:      v.Domain,
                Selector:    v.Identifier,
                HeaderField: v.HeaderKeys[0],
            }, nil
        }
    }

    return &VerificationResult{Valid: false, Error: "No valid signature"}, nil
}
```

---

## 2.2 SPF (Sender Policy Framework) [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| SPF-001 | SPF record parsing | [ ] | F-002 | `miekg/dns` |
| SPF-002 | SPF validation for inbound mail | [ ] | SPF-001, S-003 |
| SPF-003 | SPF result headers | [ ] | SPF-002 |
| SPF-004 | IPv4/IPv6 support | [ ] | SPF-001 |
| SPF-005 | Configurable handling (none/softfail/fail) | [ ] | SPF-002 |

### SPF-001: SPF Record Lookup

```go
// internal/security/spf/resolver.go
package spf

import (
    "github.com/miekg/dns"
    "net"
)

type Resolver struct {
    client *dns.Client
}

func (r *Resolver) LookupSPF(domain string) (string, error) {
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(domain), dns.TypeTXT)

    resp, _, err := r.client.Exchange(m, "8.8.8.8:53")
    if err != nil {
        return "", err
    }

    for _, ans := range resp.Answer {
        if txt, ok := ans.(*dns.TXT); ok {
            record := strings.Join(txt.Txt, "")
            if strings.HasPrefix(record, "v=spf1") {
                return record, nil
            }
        }
    }

    return "", ErrNoSPFRecord
}
```

### SPF-002: SPF Validation

```go
// internal/security/spf/validator.go
package spf

type Result string

const (
    ResultNone      Result = "none"
    ResultNeutral   Result = "neutral"
    ResultPass      Result = "pass"
    ResultFail      Result = "fail"
    ResultSoftFail  Result = "softfail"
    ResultTempError Result = "temperror"
    ResultPermError Result = "permerror"
)

type Validator struct {
    resolver *Resolver
}

func (v *Validator) Check(ip net.IP, domain, sender string) (Result, error) {
    record, err := v.resolver.LookupSPF(domain)
    if err != nil {
        if err == ErrNoSPFRecord {
            return ResultNone, nil
        }
        return ResultTempError, err
    }

    return v.evaluate(record, ip, domain, sender)
}

func (v *Validator) evaluate(record string, ip net.IP, domain, sender string) (Result, error) {
    mechanisms := parseSPF(record)

    for _, mech := range mechanisms {
        match, err := v.matchMechanism(mech, ip, domain)
        if err != nil {
            continue
        }
        if match {
            return mech.Qualifier, nil
        }
    }

    return ResultNeutral, nil
}
```

**Acceptance Criteria**:
- [ ] Support all SPF mechanisms: ip4, ip6, a, mx, include, all
- [ ] Support qualifiers: + (pass), - (fail), ~ (softfail), ? (neutral)
- [ ] IPv4 and IPv6 address matching
- [ ] DNS lookup timeout: 5 seconds per query
- [ ] Maximum 10 DNS lookups per SPF check (prevent DoS)
- [ ] Results added to Authentication-Results header
- [ ] Softfail treated as pass with warning logged
- [ ] Fail results in SMTP rejection (configurable)
- [ ] Cache SPF records for 1 hour (TTL-based)
- [ ] PermError for malformed SPF records

**Structured Logging (slog)**:
- [ ] **INFO**: SPF check passed (domain, sending_ip, spf_result="pass", mechanism_matched, auth_results_header)
- [ ] **WARN**: SPF softfail (domain, sending_ip, spf_result="softfail", spf_record, auth_results_header)
- [ ] **ERROR**: SPF check failed (domain, sending_ip, spf_result="fail", spf_record, rejection_code="550", message_id)
- [ ] **DEBUG**: SPF DNS lookup (domain, spf_record, lookup_duration_ms, ttl, cached)
- [ ] **WARN**: SPF lookup limit exceeded (domain, lookup_count=11, max_lookups=10, spf_result="permerror")
- [ ] **TRACE**: SPF mechanism evaluation (domain, mechanism, qualifier, ip_match_result, evaluation_order)
- [ ] **Fields**: domain, sending_ip, spf_result, spf_record, mechanism_matched, lookup_count, duration_ms, message_id

**Given/When/Then Scenarios**:
```
Given domain "example.com" has SPF record "v=spf1 ip4:192.0.2.0/24 -all"
And sending IP is 192.0.2.50
When SPF check is performed
Then result is "pass"
And Authentication-Results header contains "spf=pass"
And message is accepted

Given domain "example.org" has SPF record "v=spf1 ip4:198.51.100.0/24 -all"
And sending IP is 203.0.113.10 (not in range)
When SPF check is performed
Then result is "fail"
And SMTP returns "550 5.7.1 SPF check failed"
And message is rejected
And sender receives bounce notification

Given domain "test.com" has SPF record "v=spf1 ~all" (softfail)
And sending IP is 10.0.0.1
When SPF check is performed
Then result is "softfail"
And warning is logged
And message is accepted (softfail = accept with warning)
And Authentication-Results header contains "spf=softfail"

Given domain "nodomain.invalid" has no SPF record
When SPF check is performed
Then DNS query returns NXDOMAIN
And result is "none"
And message is accepted (no SPF = neutral)
And Authentication-Results header contains "spf=none"

Given domain "include.example" has SPF "v=spf1 include:sendgrid.net -all"
And sendgrid.net has SPF "v=spf1 ip4:149.72.0.0/16 -all"
And sending IP is 149.72.100.50
When SPF check is performed with include resolution
Then include mechanism is followed
And nested SPF record is evaluated
And result is "pass"
And DNS lookup count is 2 (includes both lookups)

Given SPF check requires 11 DNS lookups (exceeds limit)
When SPF validation is performed
Then check terminates after 10 lookups
And result is "permerror"
And error is logged: "SPF lookup limit exceeded"
And message handling follows permerror policy
```

---

## 2.3 DMARC (Domain-based Message Authentication) [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DM-001 | DMARC record parsing | [ ] | SPF-001 |
| DM-002 | DMARC policy enforcement | [ ] | DM-001, DK-004, SPF-002 |
| DM-003 | Alignment checking (relaxed/strict) | [ ] | DM-002 |
| DM-004 | Aggregate report generation | [ ] | DM-002 |
| DM-005 | Forensic report generation | [ ] | DM-002 |
| DM-006 | Report sending scheduler | [ ] | DM-004, Q-002 |

### DM-001: DMARC Record Lookup

```go
// internal/security/dmarc/resolver.go
package dmarc

type Policy struct {
    Version          string  // v=DMARC1
    Policy           string  // p=none|quarantine|reject
    SubdomainPolicy  string  // sp=
    Percentage       int     // pct=
    DKIM             string  // adkim=r|s
    SPF              string  // aspf=r|s
    ReportAggregate  string  // rua=
    ReportForensic   string  // ruf=
    ReportInterval   int     // ri=
    FailureReporting string  // fo=
}

func (r *Resolver) LookupDMARC(domain string) (*Policy, error) {
    dmarcDomain := "_dmarc." + domain
    record, err := r.lookupTXT(dmarcDomain)
    if err != nil {
        return nil, err
    }

    return parseDMARCRecord(record)
}
```

### DM-002: DMARC Enforcement

```go
// internal/security/dmarc/enforcer.go
package dmarc

type EnforcementResult struct {
    Policy      string  // none, quarantine, reject
    Action      string  // none, quarantine, reject
    SPFResult   spf.Result
    DKIMResult  bool
    SPFAligned  bool
    DKIMAligned bool
    Reason      string
}

func (e *Enforcer) Enforce(msg *Message) (*EnforcementResult, error) {
    fromDomain := extractDomain(msg.From)

    policy, err := e.resolver.LookupDMARC(fromDomain)
    if err != nil {
        return &EnforcementResult{Policy: "none", Action: "none"}, nil
    }

    result := &EnforcementResult{
        Policy:     policy.Policy,
        SPFResult:  msg.SPFResult,
        DKIMResult: msg.DKIMValid,
    }

    // Check alignment
    result.SPFAligned = e.checkSPFAlignment(msg, policy)
    result.DKIMAligned = e.checkDKIMAlignment(msg, policy)

    // Determine action
    if result.SPFAligned || result.DKIMAligned {
        result.Action = "none" // Pass
    } else {
        result.Action = policy.Policy
        result.Reason = "DMARC alignment failed"
    }

    return result, nil
}

func (e *Enforcer) checkDKIMAlignment(msg *Message, policy *Policy) bool {
    fromDomain := extractDomain(msg.From)
    dkimDomain := msg.DKIMDomain

    if policy.DKIM == "s" { // Strict
        return fromDomain == dkimDomain
    }
    // Relaxed (default)
    return hasOrgDomain(fromDomain, dkimDomain)
}
```

---

## 2.4 Anti-Virus (ClamAV) [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| AV-001 | ClamAV socket connection | [ ] | F-002 |
| AV-002 | Message scanning integration | [ ] | AV-001, M-001 |
| AV-003 | Configurable actions (reject/quarantine/tag) | [ ] | AV-002 |
| AV-004 | Per-domain/user configuration | [ ] | AV-003, U-003 |
| AV-005 | Scan result logging | [ ] | AV-002, F-010 |

### AV-001: ClamAV Client

```go
// internal/security/antivirus/clamav.go
package antivirus

import (
    "net"
    "bufio"
)

type ClamAV struct {
    socketPath string
}

func NewClamAV(socketPath string) *ClamAV {
    return &ClamAV{socketPath: socketPath}
}

func (c *ClamAV) Scan(data []byte) (*ScanResult, error) {
    conn, err := net.Dial("unix", c.socketPath)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to clamd: %w", err)
    }
    defer conn.Close()

    // Send INSTREAM command
    fmt.Fprintf(conn, "zINSTREAM\x00")

    // Send data in chunks
    chunkSize := 2048
    for i := 0; i < len(data); i += chunkSize {
        end := i + chunkSize
        if end > len(data) {
            end = len(data)
        }
        chunk := data[i:end]

        // Write chunk size (4 bytes, network byte order)
        binary.Write(conn, binary.BigEndian, uint32(len(chunk)))
        conn.Write(chunk)
    }

    // End stream
    binary.Write(conn, binary.BigEndian, uint32(0))

    // Read response
    reader := bufio.NewReader(conn)
    response, _ := reader.ReadString('\x00')

    return parseResponse(response), nil
}

type ScanResult struct {
    Clean    bool
    Virus    string
    Error    string
}

func parseResponse(response string) *ScanResult {
    response = strings.TrimSuffix(response, "\x00")

    if strings.HasSuffix(response, "OK") {
        return &ScanResult{Clean: true}
    }

    if strings.Contains(response, "FOUND") {
        parts := strings.Split(response, ":")
        virus := strings.TrimSpace(strings.TrimSuffix(parts[len(parts)-1], "FOUND"))
        return &ScanResult{Clean: false, Virus: virus}
    }

    return &ScanResult{Clean: false, Error: response}
}
```

### AV-002: Scanning Integration

```go
// internal/security/antivirus/scanner.go
package antivirus

type Scanner struct {
    clamav        *ClamAV
    domainService *service.DomainService
    logger        logger.Logger
}

type ScanAction string

const (
    ActionReject     ScanAction = "reject"
    ActionQuarantine ScanAction = "quarantine"
    ActionTag        ScanAction = "tag"
)

func (s *Scanner) ScanMessage(domain string, message []byte) (*ScanResult, ScanAction, error) {
    result, err := s.clamav.Scan(message)
    if err != nil {
        s.logger.Error("ClamAV scan failed", zap.Error(err))
        return nil, ActionTag, err // Fail open with tag
    }

    if result.Clean {
        return result, "", nil
    }

    // Get domain configuration
    cfg, _ := s.domainService.GetAntivirusConfig(domain)
    action := cfg.VirusAction
    if action == "" {
        action = ActionQuarantine // Default
    }

    s.logger.Warn("Virus detected",
        zap.String("virus", result.Virus),
        zap.String("domain", domain),
        zap.String("action", string(action)))

    return result, action, nil
}
```

---

## 2.5 Anti-Spam (SpamAssassin) [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| AS-001 | Integrate spamc client | [ ] | F-002 | `teamwork/spamc` |
| AS-002 | Message scoring integration | [ ] | AS-001, M-001 |
| AS-003 | Per-user spam threshold | [ ] | AS-002, U-001 |
| AS-004 | Spam quarantine system | [ ] | AS-002, F-022 |
| AS-005 | Learn from user actions (spam/ham) | [ ] | AS-001 |
| AS-006 | Spam report generation | [ ] | AS-002 |

### AS-001: SpamAssassin Client

```go
// internal/security/antispam/spamassassin.go
package antispam

import (
    "github.com/teamwork/spamc"
)

type SpamAssassin struct {
    client *spamc.Client
}

func NewSpamAssassin(host string, port int) *SpamAssassin {
    return &SpamAssassin{
        client: spamc.New(fmt.Sprintf("%s:%d", host, port), 10*time.Second),
    }
}

func (s *SpamAssassin) Check(message []byte) (*SpamResult, error) {
    reply, err := s.client.Check(context.Background(), bytes.NewReader(message), nil)
    if err != nil {
        return nil, err
    }

    return &SpamResult{
        Score:     reply.Score,
        Threshold: reply.BaseScore,
        IsSpam:    reply.IsSpam,
        Rules:     parseRules(reply.Headers),
    }, nil
}

func (s *SpamAssassin) Learn(message []byte, isSpam bool) error {
    var cmd string
    if isSpam {
        cmd = "LEARN_SPAM"
    } else {
        cmd = "LEARN_HAM"
    }

    _, err := s.client.Tell(context.Background(), bytes.NewReader(message), cmd, nil)
    return err
}

type SpamResult struct {
    Score     float64
    Threshold float64
    IsSpam    bool
    Rules     []SpamRule
}

type SpamRule struct {
    Name        string
    Score       float64
    Description string
}
```

### AS-004: Quarantine System

```go
// internal/service/quarantine_service.go
package service

type QuarantineService struct {
    repo   repository.QuarantineRepository
    logger logger.Logger
}

func (s *QuarantineService) Quarantine(userID int64, messageID int64, score float64) error {
    item := &domain.QuarantineItem{
        MessageID:    messageID,
        UserID:       userID,
        SpamScore:    score,
        QuarantinedAt: time.Now(),
        AutoDeleteAt:  time.Now().AddDate(0, 0, 30), // 30 days
    }
    return s.repo.Create(item)
}

func (s *QuarantineService) Release(userID int64, itemID int64) error {
    item, err := s.repo.GetByID(itemID)
    if err != nil {
        return err
    }

    if item.UserID != userID {
        return ErrUnauthorized
    }

    // Move message back to inbox
    if err := s.messageService.MoveToInbox(item.MessageID); err != nil {
        return err
    }

    return s.repo.MarkReleased(itemID)
}

func (s *QuarantineService) ListForUser(userID int64) ([]*domain.QuarantineItem, error) {
    return s.repo.ListByUser(userID)
}

func (s *QuarantineService) CleanupExpired() error {
    return s.repo.DeleteExpired()
}
```

---

## 2.6 Greylisting [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| GL-001 | Greylisting table schema | [ ] | F-022 |
| GL-002 | Implement greylist check | [ ] | GL-001, S-003 |
| GL-003 | Auto-whitelist after pass | [ ] | GL-002 |
| GL-004 | Configurable delay and expiry | [ ] | GL-002 |

### GL-002: Greylist Implementation

```go
// internal/security/greylist/greylist.go
package greylist

type Greylister struct {
    repo   repository.GreylistRepository
    delay  time.Duration // Minimum delay before accepting
    expiry time.Duration // How long to remember triplets
}

func NewGreylister(repo repository.GreylistRepository) *Greylister {
    return &Greylister{
        repo:   repo,
        delay:  5 * time.Minute,
        expiry: 30 * 24 * time.Hour, // 30 days
    }
}

type CheckResult struct {
    Action   string // "accept", "defer", "pass"
    Message  string
    WaitTime time.Duration
}

func (g *Greylister) Check(ip, sender, recipient string) (*CheckResult, error) {
    triplet, err := g.repo.Get(ip, sender, recipient)
    if err != nil {
        // First time seeing this triplet
        g.repo.Create(ip, sender, recipient)
        return &CheckResult{
            Action:   "defer",
            Message:  "Greylisting in effect, please retry later",
            WaitTime: g.delay,
        }, nil
    }

    // Check if enough time has passed
    if time.Since(triplet.FirstSeen) < g.delay {
        remaining := g.delay - time.Since(triplet.FirstSeen)
        return &CheckResult{
            Action:   "defer",
            Message:  "Greylisting in effect, please retry later",
            WaitTime: remaining,
        }, nil
    }

    // Passed greylisting
    g.repo.IncrementPass(triplet.ID)

    return &CheckResult{
        Action:  "accept",
        Message: "Greylisting passed",
    }, nil
}

func (g *Greylister) Cleanup() error {
    return g.repo.DeleteOlderThan(g.expiry)
}
```

**Acceptance Criteria**:
- [ ] Minimum delay: 5 minutes before accepting retry
- [ ] Triplet tracking: (IP address, sender email, recipient email)
- [ ] First-time triplets: defer with SMTP 451 "Greylisted, retry later"
- [ ] Auto-whitelist: after 3 successful passes, bypass greylisting
- [ ] Expiry: 30 days for inactive triplet entries
- [ ] Cleanup job: run hourly to remove expired entries
- [ ] Whitelist bypass: trusted IPs skip greylisting
- [ ] Pass counter: track successful deliveries per triplet
- [ ] Retry timing: accept after delay, reject if before delay

**Structured Logging (slog)**:
- [ ] **INFO**: New triplet created (ip_address, sender, recipient, action="defer", delay="5m", triplet_id, first_seen)
- [ ] **WARN**: Retry too soon (ip_address, sender, recipient, elapsed_time, required_delay="5m", remaining_wait, action="defer")
- [ ] **INFO**: Greylisting passed (ip_address, sender, recipient, total_wait_time, pass_count, action="accept", triplet_id)
- [ ] **INFO**: Triplet auto-whitelisted (ip_address, sender, recipient, pass_count=3, whitelisted_at, bypass_enabled)
- [ ] **DEBUG**: Cleanup executed (expired_triplets_count, oldest_expiry_date, cleanup_duration_ms)
- [ ] **TRACE**: Triplet lookup (ip_address, sender, recipient, found, first_seen, last_seen, pass_count)
- [ ] **Fields**: ip_address, sender, recipient, triplet_id, action, delay, elapsed_time, pass_count, first_seen, last_seen

**Given/When/Then Scenarios**:
```
Given IP 192.0.2.100 sends from "user@example.com" to "alice@mydomain.com" (new triplet)
When greylisting check is performed
Then triplet is created with FirstSeen = now
And action is "defer"
And SMTP returns "451 4.7.1 Greylisting in effect, please retry in 5 minutes"
And message is temporarily rejected

Given triplet (192.0.2.100, user@example.com, alice@mydomain.com) exists
And FirstSeen was 3 minutes ago
When retry is attempted
Then time elapsed (3 min) < delay (5 min)
And remaining wait time is 2 minutes
And action is "defer"
And SMTP returns "451 4.7.1 Greylisting in effect, retry in 2 minutes"

Given triplet (192.0.2.100, user@example.com, alice@mydomain.com) exists
And FirstSeen was 6 minutes ago (>5 min delay)
When retry is attempted
Then greylisting delay has passed
And action is "accept"
And pass counter increments to 1
And message is delivered normally
And sender is one step closer to whitelist

Given triplet has 3 successful passes
When 4th message arrives from same triplet
Then triplet is auto-whitelisted
And greylisting is bypassed
And message is accepted immediately
And no delay is imposed

Given triplet was last seen 31 days ago
When cleanup job runs
Then triplet is identified as expired (>30 days)
And triplet record is deleted from database
And space is reclaimed for active triplets

Given legitimate mail server retries after 5 minutes
And spambot does not retry
When delivery statistics are analyzed
Then legitimate server achieves 95%+ delivery rate
And spambot achieves 0% delivery rate (no retries)
And greylisting blocks spam without affecting real mail
```

---

## 2.7 Rate Limiting [MVP]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| RL-001 | Rate limit table schema | [ ] | F-022 |
| RL-002 | Per-IP rate limiting | [ ] | RL-001 |
| RL-003 | Per-user rate limiting | [ ] | RL-001 |
| RL-004 | Per-domain rate limiting | [ ] | RL-001 |
| RL-005 | Configurable limits | [ ] | RL-002 |

### RL-002: Rate Limiter

```go
// internal/security/ratelimit/limiter.go
package ratelimit

type Limiter struct {
    repo repository.RateLimitRepository
}

type Limit struct {
    Count    int           // Max requests
    Window   time.Duration // Time window
}

var DefaultLimits = map[string]Limit{
    "smtp_per_ip":      {Count: 100, Window: time.Hour},
    "smtp_per_user":    {Count: 500, Window: time.Hour},
    "smtp_per_domain":  {Count: 1000, Window: time.Hour},
    "auth_per_ip":      {Count: 10, Window: 15 * time.Minute},
    "imap_per_user":    {Count: 1000, Window: time.Hour},
}

func (l *Limiter) Check(limitType, key string) (bool, error) {
    limit := DefaultLimits[limitType]

    count, err := l.repo.GetCount(limitType, key, limit.Window)
    if err != nil {
        return true, err // Fail open
    }

    if count >= limit.Count {
        return false, nil // Rate limited
    }

    l.repo.Increment(limitType, key)
    return true, nil
}

func (l *Limiter) CheckIP(ip string) (bool, error) {
    return l.Check("smtp_per_ip", ip)
}

func (l *Limiter) CheckUser(userID string) (bool, error) {
    return l.Check("smtp_per_user", userID)
}

func (l *Limiter) CheckAuth(ip string) (bool, error) {
    return l.Check("auth_per_ip", ip)
}
```

**Acceptance Criteria**:
- [ ] SMTP per IP: 100 requests/hour (burst protection)
- [ ] SMTP per user: 500 requests/hour (spam prevention)
- [ ] SMTP per domain: 1000 requests/hour (domain-level limiting)
- [ ] Auth per IP: 10 attempts/15 minutes (brute force prevention)
- [ ] IMAP per user: 1000 requests/hour (DoS prevention)
- [ ] Sliding window algorithm for accurate rate tracking
- [ ] Fail-open on repository errors (availability over strict limiting)
- [ ] Redis-backed counters for distributed deployments
- [ ] Rate limit headers in SMTP responses (X-RateLimit-*)

**Structured Logging (slog)**:
- [ ] **INFO**: Rate limit check passed (limit_type="smtp_per_ip|smtp_per_user|auth_per_ip", key, current_count, limit, window, remaining)
- [ ] **WARN**: Rate limit threshold approaching (limit_type, key, current_count, limit, percentage="80%|90%", window)
- [ ] **ERROR**: Rate limit exceeded (limit_type, key, ip_address, user_id, current_count, limit, window, exceeded_by)
- [ ] **DEBUG**: Rate limit counter increment (limit_type, key, old_count, new_count, window_start, window_end)
- [ ] **TRACE**: Rate limit window reset (limit_type, keys_reset, old_window, new_window)
- [ ] **FATAL**: Rate limit repository failure (error="redis_connection_lost", fail_open=true, limit_type)
- [ ] **Fields**: limit_type, key, ip_address, user_id, current_count, limit, window, percentage, exceeded_by, duration_ms

**Given/When/Then Scenarios**:
```
Given IP 192.168.1.100 has made 50 SMTP connections
When connection 51 is attempted within the same hour
Then rate limit check passes (under 100/hour limit)
And connection is accepted
And counter increments to 51

Given IP 192.168.1.100 has made 100 SMTP connections in current hour
When connection 101 is attempted
Then rate limit check fails
And SMTP returns "421 4.7.0 Rate limit exceeded. Try again later."
And connection is rejected
And counter does not increment

Given IP 10.0.0.5 has failed authentication 9 times in 15 minutes
When 10th authentication attempt is made
Then auth rate limit check fails
And login is rejected with "Too many login attempts"
And IP is temporarily blocked for remainder of window
And security event is logged

Given user "user@example.com" has sent 499 messages in current hour
When 500th message is submitted
Then user rate limit check passes (at limit)
And message is accepted
And counter increments to 500

Given user "spammer@example.com" has sent 500 messages in current hour
When 501st message is submitted
Then user rate limit check fails
And SMTP returns "452 4.3.1 User sending rate exceeded"
And message is rejected
And admin alert is triggered

Given rate limit window started at 10:00:00
When current time is 11:00:01 (>1 hour later)
Then previous window counters expire
And new window begins
And rate limits reset to 0
```

---

## 2.8 Authentication Security [MVP]

| ID | Task | Status | Dependencies | Library |
|----|------|--------|--------------|---------|
| AU-001 | TOTP 2FA implementation | [ ] | U-001 | `pquerna/otp` |
| AU-002 | Failed login tracking | [ ] | F-022 |
| AU-003 | IP blacklisting | [ ] | AU-002 |
| AU-004 | IP whitelisting | [ ] | F-022 |
| AU-005 | Brute force protection | [ ] | AU-002 |

### AU-001: TOTP 2FA

```go
// internal/security/totp/totp.go
package totp

import (
    "github.com/pquerna/otp/totp"
)

type TOTPService struct {
    issuer string
}

func NewTOTPService(issuer string) *TOTPService {
    return &TOTPService{issuer: issuer}
}

func (s *TOTPService) GenerateSecret(email string) (*TOTPSetup, error) {
    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      s.issuer,
        AccountName: email,
    })
    if err != nil {
        return nil, err
    }

    return &TOTPSetup{
        Secret:     key.Secret(),
        URL:        key.URL(),
        QRCodeData: generateQRCode(key.URL()),
    }, nil
}

func (s *TOTPService) Validate(secret, code string) bool {
    return totp.Validate(code, secret)
}

type TOTPSetup struct {
    Secret     string
    URL        string
    QRCodeData []byte // PNG image data
}
```

### AU-002 & AU-005: Brute Force Protection

```go
// internal/security/bruteforce/protection.go
package bruteforce

type Protection struct {
    repo      repository.FailedLoginRepository
    blacklist repository.BlacklistRepository
    threshold int           // Max failures before block
    window    time.Duration // Time window for failures
    blockTime time.Duration // How long to block
}

func NewProtection(repo repository.FailedLoginRepository, bl repository.BlacklistRepository) *Protection {
    return &Protection{
        repo:      repo,
        blacklist: bl,
        threshold: 5,
        window:    15 * time.Minute,
        blockTime: 1 * time.Hour,
    }
}

func (p *Protection) RecordFailure(ip, username string) error {
    p.repo.Create(ip, username)

    // Check if threshold exceeded
    count, _ := p.repo.CountByIP(ip, p.window)
    if count >= p.threshold {
        p.blacklist.Add(ip, "Brute force protection", time.Now().Add(p.blockTime))
    }

    return nil
}

func (p *Protection) IsBlocked(ip string) bool {
    return p.blacklist.Exists(ip)
}

func (p *Protection) ClearOnSuccess(ip string) {
    p.repo.DeleteByIP(ip)
}
```

**Acceptance Criteria**:
- [ ] Threshold: 5 failed login attempts within 15-minute window
- [ ] Block duration: 1 hour after threshold exceeded
- [ ] Failed login tracking per IP address
- [ ] Failed login tracking per username (separate counter)
- [ ] Automatic blacklist addition at threshold
- [ ] Failed login counter resets on successful authentication
- [ ] Blacklist entries expire automatically after block duration
- [ ] Security event logging for all brute force blocks
- [ ] Admin notification on brute force detection
- [ ] Whitelist bypass for trusted IPs

**Structured Logging (slog)**:
- [ ] **WARN**: Failed login attempt recorded (ip_address, username, attempt_count, threshold=5, window="15m", session_id)
- [ ] **ERROR**: Brute force threshold exceeded (ip_address, username, attempt_count=5, block_duration="1h", event_type="brute_force_detected")
- [ ] **ERROR**: Blocked IP authentication attempt (ip_address, block_reason="brute_force", remaining_block_time, attempted_username)
- [ ] **INFO**: Failed login counter cleared (ip_address, reason="successful_auth", previous_attempts)
- [ ] **INFO**: Blacklist entry expired (ip_address, block_reason, blocked_duration, total_attempts_during_block)
- [ ] **DEBUG**: Blacklist check (ip_address, is_blocked, block_expires_at, username)
- [ ] **FATAL**: Admin notification failed (ip_address, event_type="brute_force", error_msg, notification_channel)
- [ ] **Fields**: ip_address, username, user_id, attempt_count, threshold, window, block_duration, event_type, session_id, trace_id

**Given/When/Then Scenarios**:
```
Given IP 192.168.1.50 has 0 failed login attempts
When authentication fails
Then failed login is recorded with IP and username
And counter increments to 1
And IP is not blocked (under threshold)

Given IP 10.0.0.100 has 4 failed login attempts in last 15 minutes
When 5th authentication failure occurs
Then brute force threshold is exceeded
And IP is added to blacklist for 1 hour
And security event is logged with "Brute force detected from 10.0.0.100"
And admin notification is sent

Given IP 172.16.0.50 is blacklisted for brute force
When authentication is attempted from this IP
Then request is rejected immediately before password check
And response is "Too many failed login attempts. Try again later."
And no database queries are executed (fail fast)

Given IP 192.168.1.100 has 3 failed login attempts
When successful authentication occurs from this IP
Then failed login counter is cleared for this IP
And IP is not blocked
And normal access is restored

Given IP 203.0.113.50 was blacklisted at 10:00 AM
When current time is 11:01 AM (>1 hour later)
Then blacklist entry expires automatically
And IP can attempt authentication again
And counter resets to 0

Given IP 10.1.1.1 has failed login for user "alice@example.com" 3 times
And same IP has failed login for user "bob@example.com" 2 times
When brute force check runs
Then combined IP counter is 5 (threshold)
And IP is blocked for targeting multiple accounts
And both usernames are logged in security event
```

---

## Acceptance Criteria

### DKIM
- [ ] RSA-2048 and RSA-4096 keys can be generated
- [ ] Ed25519 keys can be generated
- [ ] Outbound mail is signed correctly
- [ ] Inbound signatures are verified
- [ ] DNS record format is correct
- [ ] Multiple selectors work

### SPF
- [ ] SPF records are parsed correctly
- [ ] All mechanisms evaluated (ip4, ip6, a, mx, include, etc.)
- [ ] Results added to headers
- [ ] IPv4 and IPv6 supported

### DMARC
- [ ] DMARC records parsed correctly
- [ ] Policy enforcement works (none, quarantine, reject)
- [ ] Alignment checking works (relaxed and strict)
- [ ] Aggregate reports generated

### ClamAV
- [ ] Socket connection works
- [ ] Messages scanned successfully
- [ ] Viruses detected and handled
- [ ] Actions configurable per domain

### SpamAssassin
- [ ] spamc connection works
- [ ] Spam scoring returns results
- [ ] Per-user thresholds applied
- [ ] Quarantine system works
- [ ] Learn spam/ham works

### Greylisting
- [ ] First-time senders deferred
- [ ] Retry after delay accepted
- [ ] Auto-whitelist after pass
- [ ] Cleanup removes old entries

### Rate Limiting
- [ ] Per-IP limits enforced
- [ ] Per-user limits enforced
- [ ] Per-domain limits enforced
- [ ] Limits configurable

### Auth Security
- [ ] TOTP setup works
- [ ] TOTP validation works
- [ ] Failed logins tracked
- [ ] Brute force protection blocks attackers
- [ ] Whitelist bypasses checks

---

## Go Dependencies for Phase 2

```go
// Additional go.mod entries
require (
    github.com/emersion/go-msgauth v0.6.8
    github.com/miekg/dns v1.1.58
    github.com/teamwork/spamc v0.0.0-20200109085853-a4e0c5c3f7a0
    github.com/pquerna/otp v1.4.0
)
```

---

## Testing

```bash
# Test DKIM signing
echo "Subject: Test" | ./gomailserver dkim sign example.com

# Test SPF
./gomailserver spf check 192.0.2.1 example.com

# Test ClamAV
echo "X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*" | \
  ./gomailserver scan

# Test SpamAssassin
./gomailserver spam check < test_message.eml
```

---

## Next Phase

After completing Phase 2, proceed to [TASKS3.md](TASKS3.md) - Web Interfaces.
