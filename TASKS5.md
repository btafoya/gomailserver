# Phase 5: Advanced Security (Weeks 16-17)

**Status**: Not Started
**Priority**: Full Feature (Post-MVP)
**Estimated Duration**: 2 weeks
**Dependencies**: Phase 2 (Security Foundation)

---

## Overview

Implement advanced email security features including DANE (DNS-based Authentication of Named Entities), MTA-STS (Mail Transfer Agent Strict Transport Security), PGP/GPG key management, and comprehensive audit logging.

---

## 5.1 DANE (DNS-based Authentication of Named Entities) [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| DA-001 | TLSA record lookup | [ ] | F-002 |
| DA-002 | DNSSEC validation | [ ] | DA-001 |
| DA-003 | DANE-TA support | [ ] | DA-001 |
| DA-004 | DANE-EE support | [ ] | DA-001 |
| DA-005 | Fallback mechanisms | [ ] | DA-001 |

### DA-001: TLSA Record Lookup

```go
// internal/security/dane/resolver.go
package dane

import (
    "github.com/miekg/dns"
)

type TLSARecord struct {
    Usage        uint8  // 0: CA, 1: Service cert, 2: Trust anchor, 3: Domain-issued
    Selector     uint8  // 0: Full cert, 1: SubjectPublicKeyInfo
    MatchingType uint8  // 0: No hash, 1: SHA-256, 2: SHA-512
    Certificate  string // Hex-encoded
}

type Resolver struct {
    client *dns.Client
}

func (r *Resolver) LookupTLSA(hostname string, port int) ([]TLSARecord, error) {
    // Construct TLSA record name: _port._tcp.hostname
    name := fmt.Sprintf("_%d._tcp.%s", port, dns.Fqdn(hostname))

    m := new(dns.Msg)
    m.SetQuestion(name, dns.TypeTLSA)
    m.SetEdns0(4096, true) // Enable DNSSEC

    resp, _, err := r.client.Exchange(m, "8.8.8.8:53")
    if err != nil {
        return nil, err
    }

    records := []TLSARecord{}
    for _, ans := range resp.Answer {
        if tlsa, ok := ans.(*dns.TLSA); ok {
            records = append(records, TLSARecord{
                Usage:        tlsa.Usage,
                Selector:     tlsa.Selector,
                MatchingType: tlsa.MatchingType,
                Certificate:  tlsa.Certificate,
            })
        }
    }

    return records, nil
}
```

### DA-002: DNSSEC Validation

```go
// internal/security/dane/dnssec.go
package dane

type DNSSECValidator struct {
    resolver *Resolver
}

func (v *DNSSECValidator) Validate(domain string) (*ValidationResult, error) {
    m := new(dns.Msg)
    m.SetQuestion(dns.Fqdn(domain), dns.TypeTLSA)
    m.SetEdns0(4096, true)
    m.CheckingDisabled = false
    m.RecursionDesired = true
    m.AuthenticatedData = true

    resp, _, err := v.resolver.client.Exchange(m, "8.8.8.8:53")
    if err != nil {
        return nil, err
    }

    result := &ValidationResult{
        Secure:         resp.AuthenticatedData,
        Bogus:          false,
        Indeterminate:  false,
    }

    // Check for bogus responses
    if resp.Rcode == dns.RcodeServerFailure {
        result.Bogus = true
    }

    return result, nil
}

type ValidationResult struct {
    Secure        bool
    Bogus         bool
    Indeterminate bool
    Reason        string
}
```

### DA-003 & DA-004: DANE Verification

```go
// internal/security/dane/verifier.go
package dane

import (
    "crypto/sha256"
    "crypto/sha512"
    "crypto/x509"
)

type Verifier struct {
    resolver *Resolver
    dnssec   *DNSSECValidator
}

type VerificationResult struct {
    Valid    bool
    Secure   bool
    UsedType string // DANE-TA, DANE-EE
    Error    string
}

func (v *Verifier) Verify(hostname string, port int, cert *x509.Certificate) (*VerificationResult, error) {
    // Lookup TLSA records
    records, err := v.resolver.LookupTLSA(hostname, port)
    if err != nil {
        return &VerificationResult{Valid: false, Error: err.Error()}, nil
    }

    if len(records) == 0 {
        return &VerificationResult{Valid: false, Error: "No TLSA records found"}, nil
    }

    // Validate DNSSEC
    dnssecResult, _ := v.dnssec.Validate(hostname)

    for _, tlsa := range records {
        match, usedType := v.matchCertificate(tlsa, cert)
        if match {
            return &VerificationResult{
                Valid:    true,
                Secure:   dnssecResult.Secure,
                UsedType: usedType,
            }, nil
        }
    }

    return &VerificationResult{Valid: false, Error: "No matching TLSA record"}, nil
}

func (v *Verifier) matchCertificate(tlsa TLSARecord, cert *x509.Certificate) (bool, string) {
    var data []byte

    // Select data based on selector
    switch tlsa.Selector {
    case 0: // Full certificate
        data = cert.Raw
    case 1: // SubjectPublicKeyInfo
        data = cert.RawSubjectPublicKeyInfo
    }

    // Hash data based on matching type
    var hash string
    switch tlsa.MatchingType {
    case 0: // No hash
        hash = fmt.Sprintf("%x", data)
    case 1: // SHA-256
        h := sha256.Sum256(data)
        hash = fmt.Sprintf("%x", h)
    case 2: // SHA-512
        h := sha512.Sum512(data)
        hash = fmt.Sprintf("%x", h)
    }

    // Determine usage type
    usedType := ""
    switch tlsa.Usage {
    case 2:
        usedType = "DANE-TA"
    case 3:
        usedType = "DANE-EE"
    }

    return strings.EqualFold(hash, tlsa.Certificate), usedType
}
```

---

## 5.2 MTA-STS (Mail Transfer Agent Strict Transport Security) [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| MS-001 | MTA-STS policy fetching | [ ] | F-002 |
| MS-002 | Policy caching | [ ] | MS-001, F-022 |
| MS-003 | TLS enforcement based on policy | [ ] | MS-002 |
| MS-004 | TLSRPT reporting (RFC 8460) | [ ] | MS-001 |

### MS-001: MTA-STS Policy Fetching

```go
// internal/security/mtasts/policy.go
package mtasts

import (
    "net/http"
)

type Policy struct {
    Version   string   `json:"version"`
    Mode      string   `json:"mode"` // enforce, testing, none
    MaxAge    int64    `json:"max_age"`
    MX        []string `json:"mx"`
    FetchedAt time.Time
    ExpiresAt time.Time
}

type PolicyFetcher struct {
    httpClient *http.Client
}

func (f *PolicyFetcher) Fetch(domain string) (*Policy, error) {
    // First, check for _mta-sts TXT record
    txtRecord, err := f.lookupMTASTSRecord(domain)
    if err != nil {
        return nil, err
    }

    if txtRecord == "" {
        return nil, ErrNoMTASTS
    }

    // Fetch policy from well-known URL
    policyURL := fmt.Sprintf("https://mta-sts.%s/.well-known/mta-sts.txt", domain)

    resp, err := f.httpClient.Get(policyURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, ErrPolicyNotFound
    }

    return f.parsePolicy(resp.Body)
}

func (f *PolicyFetcher) parsePolicy(r io.Reader) (*Policy, error) {
    policy := &Policy{
        FetchedAt: time.Now(),
    }

    scanner := bufio.NewScanner(r)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }

        parts := strings.SplitN(line, ":", 2)
        if len(parts) != 2 {
            continue
        }

        key := strings.TrimSpace(parts[0])
        value := strings.TrimSpace(parts[1])

        switch key {
        case "version":
            policy.Version = value
        case "mode":
            policy.Mode = value
        case "max_age":
            policy.MaxAge, _ = strconv.ParseInt(value, 10, 64)
        case "mx":
            policy.MX = append(policy.MX, value)
        }
    }

    policy.ExpiresAt = policy.FetchedAt.Add(time.Duration(policy.MaxAge) * time.Second)

    return policy, nil
}
```

### MS-002: Policy Cache

```go
// internal/security/mtasts/cache.go
package mtasts

type PolicyCache struct {
    repo repository.MTASTSRepository
}

func (c *PolicyCache) Get(domain string) (*Policy, error) {
    cached, err := c.repo.Get(domain)
    if err != nil {
        return nil, err
    }

    if cached.ExpiresAt.Before(time.Now()) {
        return nil, ErrPolicyExpired
    }

    return cached, nil
}

func (c *PolicyCache) Set(domain string, policy *Policy) error {
    return c.repo.Upsert(domain, policy)
}

// SQL Schema
/*
CREATE TABLE mta_sts_policies (
    domain TEXT PRIMARY KEY,
    version TEXT,
    mode TEXT,
    max_age INTEGER,
    mx TEXT,
    fetched_at DATETIME,
    expires_at DATETIME
);
*/
```

### MS-003: TLS Enforcement

```go
// internal/security/mtasts/enforcer.go
package mtasts

type Enforcer struct {
    fetcher *PolicyFetcher
    cache   *PolicyCache
    logger  logger.Logger
}

type EnforcementResult struct {
    Enforce   bool
    Policy    *Policy
    MXMatch   bool
    TLSResult string
}

func (e *Enforcer) CheckDelivery(domain string, mxHost string) (*EnforcementResult, error) {
    // Try cache first
    policy, err := e.cache.Get(domain)
    if err != nil {
        // Fetch fresh policy
        policy, err = e.fetcher.Fetch(domain)
        if err != nil {
            // No MTA-STS policy - proceed without enforcement
            return &EnforcementResult{Enforce: false}, nil
        }
        e.cache.Set(domain, policy)
    }

    result := &EnforcementResult{
        Policy:  policy,
        Enforce: policy.Mode == "enforce",
    }

    // Check MX match
    for _, mx := range policy.MX {
        if matchMX(mxHost, mx) {
            result.MXMatch = true
            break
        }
    }

    if !result.MXMatch && result.Enforce {
        return result, ErrMXMismatch
    }

    return result, nil
}

func matchMX(host, pattern string) bool {
    if strings.HasPrefix(pattern, "*.") {
        // Wildcard match
        suffix := pattern[1:]
        return strings.HasSuffix(host, suffix)
    }
    return host == pattern
}
```

### MS-004: TLSRPT Reporting

```go
// internal/security/mtasts/tlsrpt.go
package mtasts

type TLSReport struct {
    OrganizationName string        `json:"organization-name"`
    DateRange        DateRange     `json:"date-range"`
    ContactInfo      string        `json:"contact-info"`
    ReportID         string        `json:"report-id"`
    Policies         []PolicyStats `json:"policies"`
}

type DateRange struct {
    StartDateTime time.Time `json:"start-datetime"`
    EndDateTime   time.Time `json:"end-datetime"`
}

type PolicyStats struct {
    Policy           PolicyInfo     `json:"policy"`
    Summary          Summary        `json:"summary"`
    FailureDetails   []FailureDetail `json:"failure-details,omitempty"`
}

type Summary struct {
    TotalSuccessfulCount int `json:"total-successful-session-count"`
    TotalFailureCount    int `json:"total-failure-session-count"`
}

type TLSReportGenerator struct {
    repo   repository.TLSReportRepository
    domain string
    email  string
}

func (g *TLSReportGenerator) Generate(start, end time.Time) (*TLSReport, error) {
    stats, err := g.repo.GetStats(start, end)
    if err != nil {
        return nil, err
    }

    report := &TLSReport{
        OrganizationName: g.domain,
        DateRange: DateRange{
            StartDateTime: start,
            EndDateTime:   end,
        },
        ContactInfo: g.email,
        ReportID:    uuid.New().String(),
        Policies:    stats,
    }

    return report, nil
}

func (g *TLSReportGenerator) Send(report *TLSReport, recipientDomain string) error {
    // Lookup TLSRPT record for recipient
    reportURL, err := g.lookupTLSRPT(recipientDomain)
    if err != nil {
        return err
    }

    // Compress and send report
    reportJSON, _ := json.Marshal(report)
    compressed := gzip.Compress(reportJSON)

    if strings.HasPrefix(reportURL, "mailto:") {
        return g.sendByEmail(reportURL[7:], compressed)
    } else if strings.HasPrefix(reportURL, "https://") {
        return g.sendByHTTPS(reportURL, compressed)
    }

    return ErrInvalidReportURL
}
```

---

## 5.3 PGP/GPG Integration [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| PGP-001 | User PGP key storage | [ ] | F-022 |
| PGP-002 | Key import/export API | [ ] | PGP-001, API-002 |
| PGP-003 | Automatic encryption when key available | [ ] | PGP-001 |
| PGP-004 | Signature verification | [ ] | PGP-001 |
| PGP-005 | Web UI for key management | [ ] | UP-001, PGP-002 |

### PGP-001: Key Storage

```go
// internal/domain/pgp.go
package domain

type PGPKey struct {
    ID          int64     `json:"id"`
    UserID      int64     `json:"user_id"`
    KeyID       string    `json:"key_id"`       // Short key ID
    Fingerprint string    `json:"fingerprint"`  // Full fingerprint
    PublicKey   string    `json:"public_key"`   // ASCII armored
    Algorithm   string    `json:"algorithm"`    // RSA, Ed25519, etc.
    KeySize     int       `json:"key_size"`
    ExpiresAt   time.Time `json:"expires_at"`
    Primary     bool      `json:"primary"`
    Revoked     bool      `json:"revoked"`
    CreatedAt   time.Time `json:"created_at"`
}

// SQL Schema
/*
CREATE TABLE pgp_keys (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    key_id TEXT NOT NULL,
    fingerprint TEXT NOT NULL UNIQUE,
    public_key TEXT NOT NULL,
    algorithm TEXT,
    key_size INTEGER,
    expires_at DATETIME,
    is_primary INTEGER DEFAULT 0,
    revoked INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pgp_keys_user ON pgp_keys(user_id);
CREATE INDEX idx_pgp_keys_keyid ON pgp_keys(key_id);
*/
```

### PGP-002: Key Import/Export

```go
// internal/service/pgp_service.go
package service

import (
    "golang.org/x/crypto/openpgp"
    "golang.org/x/crypto/openpgp/armor"
)

type PGPService struct {
    repo   repository.PGPKeyRepository
    logger logger.Logger
}

func (s *PGPService) ImportKey(userID int64, armoredKey string) (*domain.PGPKey, error) {
    // Parse armored key
    block, err := armor.Decode(strings.NewReader(armoredKey))
    if err != nil {
        return nil, ErrInvalidKey
    }

    entityList, err := openpgp.ReadKeyRing(block.Body)
    if err != nil {
        return nil, ErrInvalidKey
    }

    if len(entityList) == 0 {
        return nil, ErrNoKeysFound
    }

    entity := entityList[0]
    primaryKey := entity.PrimaryKey

    key := &domain.PGPKey{
        UserID:      userID,
        KeyID:       fmt.Sprintf("%X", primaryKey.KeyId),
        Fingerprint: fmt.Sprintf("%X", primaryKey.Fingerprint),
        PublicKey:   armoredKey,
        Algorithm:   algorithmName(primaryKey.PubKeyAlgo),
        KeySize:     primaryKey.BitLength,
    }

    // Check expiry
    for _, identity := range entity.Identities {
        if identity.SelfSignature != nil && identity.SelfSignature.KeyExpires != nil {
            key.ExpiresAt = primaryKey.CreationTime.Add(*identity.SelfSignature.KeyExpires)
        }
    }

    if err := s.repo.Create(key); err != nil {
        return nil, err
    }

    return key, nil
}

func (s *PGPService) ExportKey(userID int64, keyID string) (string, error) {
    key, err := s.repo.GetByKeyID(userID, keyID)
    if err != nil {
        return "", err
    }
    return key.PublicKey, nil
}

func (s *PGPService) LookupByEmail(email string) (*domain.PGPKey, error) {
    return s.repo.GetPrimaryByEmail(email)
}
```

### PGP-003: Automatic Encryption

```go
// internal/security/pgp/encryptor.go
package pgp

type Encryptor struct {
    keyService *service.PGPService
    logger     logger.Logger
}

func (e *Encryptor) EncryptIfPossible(recipient string, message []byte) ([]byte, bool, error) {
    // Look up recipient's public key
    key, err := e.keyService.LookupByEmail(recipient)
    if err != nil {
        // No key found - return original message
        return message, false, nil
    }

    if key.Revoked || (key.ExpiresAt.Before(time.Now()) && !key.ExpiresAt.IsZero()) {
        return message, false, nil
    }

    // Encrypt message
    encrypted, err := e.encrypt(key.PublicKey, message)
    if err != nil {
        e.logger.Warn("PGP encryption failed", zap.Error(err))
        return message, false, err
    }

    return encrypted, true, nil
}

func (e *Encryptor) encrypt(armoredKey string, message []byte) ([]byte, error) {
    block, _ := armor.Decode(strings.NewReader(armoredKey))
    entityList, _ := openpgp.ReadKeyRing(block.Body)

    var buf bytes.Buffer
    armorWriter, _ := armor.Encode(&buf, "PGP MESSAGE", nil)

    plainWriter, err := openpgp.Encrypt(armorWriter, entityList, nil, nil, nil)
    if err != nil {
        return nil, err
    }

    plainWriter.Write(message)
    plainWriter.Close()
    armorWriter.Close()

    return buf.Bytes(), nil
}
```

### PGP-004: Signature Verification

```go
// internal/security/pgp/verifier.go
package pgp

type SignatureResult struct {
    Valid       bool
    SignerKeyID string
    SignerEmail string
    SignedAt    time.Time
    Error       string
}

func (v *Verifier) VerifySignature(message []byte, signature []byte) (*SignatureResult, error) {
    sigBlock, _ := armor.Decode(bytes.NewReader(signature))

    // Try to find signer's key
    keyring, err := v.buildKeyring()
    if err != nil {
        return nil, err
    }

    entity, err := openpgp.CheckDetachedSignature(keyring, bytes.NewReader(message), sigBlock.Body)
    if err != nil {
        return &SignatureResult{
            Valid: false,
            Error: err.Error(),
        }, nil
    }

    result := &SignatureResult{
        Valid:       true,
        SignerKeyID: fmt.Sprintf("%X", entity.PrimaryKey.KeyId),
    }

    for name := range entity.Identities {
        result.SignerEmail = extractEmail(name)
        break
    }

    return result, nil
}
```

---

## 5.4 Audit Logging [FULL]

| ID | Task | Status | Dependencies |
|----|------|--------|--------------|
| AL-001 | Admin action audit trail | [ ] | F-022 |
| AL-002 | Security event logging | [ ] | AL-001 |
| AL-003 | Audit log viewer in admin UI | [ ] | AUI-001, AL-001 |

### AL-001: Audit Trail

```go
// internal/audit/logger.go
package audit

type AuditLogger struct {
    repo   repository.AuditLogRepository
    logger logger.Logger
}

type AuditEvent struct {
    ID         int64     `json:"id"`
    UserID     int64     `json:"user_id"`
    Action     string    `json:"action"`
    TargetType string    `json:"target_type"`
    TargetID   int64     `json:"target_id"`
    OldValue   string    `json:"old_value"`
    NewValue   string    `json:"new_value"`
    IPAddress  string    `json:"ip_address"`
    UserAgent  string    `json:"user_agent"`
    Timestamp  time.Time `json:"timestamp"`
}

func (a *AuditLogger) Log(ctx context.Context, action string, target interface{}, oldValue, newValue interface{}) {
    userID := ctx.Value("user_id").(int64)
    ipAddress := ctx.Value("ip_address").(string)

    event := &AuditEvent{
        UserID:     userID,
        Action:     action,
        TargetType: reflect.TypeOf(target).Name(),
        TargetID:   getID(target),
        OldValue:   toJSON(oldValue),
        NewValue:   toJSON(newValue),
        IPAddress:  ipAddress,
        Timestamp:  time.Now(),
    }

    if err := a.repo.Create(event); err != nil {
        a.logger.Error("Failed to log audit event", zap.Error(err))
    }
}

// Audit actions
const (
    ActionDomainCreate = "domain.create"
    ActionDomainUpdate = "domain.update"
    ActionDomainDelete = "domain.delete"
    ActionUserCreate   = "user.create"
    ActionUserUpdate   = "user.update"
    ActionUserDelete   = "user.delete"
    ActionUserLogin    = "user.login"
    ActionUserLogout   = "user.logout"
    ActionDKIMGenerate = "dkim.generate"
    ActionDKIMRotate   = "dkim.rotate"
    ActionConfigChange = "config.change"
    ActionBackupCreate = "backup.create"
    ActionBackupRestore = "backup.restore"
)
```

### AL-002: Security Events

```go
// internal/audit/security.go
package audit

type SecurityEventLogger struct {
    repo   repository.SecurityEventRepository
    logger logger.Logger
}

type SecurityEvent struct {
    ID        int64     `json:"id"`
    Type      string    `json:"type"`
    Severity  string    `json:"severity"` // info, warning, critical
    Source    string    `json:"source"`   // IP address
    Details   string    `json:"details"`
    Timestamp time.Time `json:"timestamp"`
}

const (
    EventLoginSuccess     = "login.success"
    EventLoginFailed      = "login.failed"
    EventBruteForceBlock  = "bruteforce.block"
    EventSpamDetected     = "spam.detected"
    EventVirusDetected    = "virus.detected"
    EventDMARCFail        = "dmarc.fail"
    EventDANEFail         = "dane.fail"
    EventRateLimitHit     = "ratelimit.hit"
    EventCertExpiring     = "cert.expiring"
    EventUnauthorizedAccess = "unauthorized.access"
)

func (s *SecurityEventLogger) Log(eventType, severity, source string, details map[string]interface{}) {
    event := &SecurityEvent{
        Type:      eventType,
        Severity:  severity,
        Source:    source,
        Details:   toJSON(details),
        Timestamp: time.Now(),
    }

    s.repo.Create(event)

    // Alert on critical events
    if severity == "critical" {
        s.sendAlert(event)
    }
}
```

---

## Acceptance Criteria

### DANE
- [ ] TLSA records looked up correctly
- [ ] DNSSEC validation works
- [ ] DANE-TA verification works
- [ ] DANE-EE verification works
- [ ] Graceful fallback when DANE unavailable

### MTA-STS
- [ ] Policy fetched from well-known URL
- [ ] Policy cached correctly
- [ ] TLS enforced when mode=enforce
- [ ] MX patterns matched correctly
- [ ] TLSRPT reports generated and sent

### PGP/GPG
- [ ] Keys can be imported
- [ ] Keys can be exported
- [ ] Automatic encryption works
- [ ] Signature verification works
- [ ] Key management UI works

### Audit Logging
- [ ] All admin actions logged
- [ ] Security events logged
- [ ] Audit log searchable
- [ ] Retention policy applied

---

## Go Dependencies for Phase 5

```go
// Additional go.mod entries
require (
    golang.org/x/crypto v0.19.0 // includes openpgp
    github.com/google/uuid v1.6.0
)
```

---

## Next Phase

After completing Phase 5, proceed to [TASKS6.md](TASKS6.md) - Sieve Filtering.
