# Phase 5: Automation & User Onboarding

**Status**: Planning
**Priority**: High
**Dependencies**: Phase 3 (Settings Management), Phase 4 (CalDAV/CardDAV), Let's Encrypt Integration

## Overview

Phase 5 eliminates manual configuration steps through automated DNS management and client autoconfiguration. This phase focuses on reducing onboarding friction and improving user experience by automating two critical setup tasks:

1. **DNS Management for Cloudflare** - Automated DNS record configuration with admin approval workflow
2. **Mail Client Autoconfiguration** - Zero-configuration client setup for Thunderbird, Outlook, and Apple Mail

## Success Criteria

✅ Admins can propose and approve DNS changes through Existing Admin UI developed in TASKS3.md and manage CLOUDFLARE API key stored in sqlite.
✅ DNS records automatically created/updated via Cloudflare API
✅ DNS propagation monitoring shows real-time status
✅ Mail clients auto-discover server settings without manual input
✅ Complete audit trail of all DNS and configuration changes
✅ Security review passed with no critical vulnerabilities
✅ Client compatibility tested across major platforms (Thunderbird, Outlook, Apple Mail)

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────┐
│                   Web UI Layer                      │
├─────────────────────────────────────────────────────┤
│  DNS Management UI  │  Approval Queue  │  Status    │
└─────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────┐
│                  API Layer                          │
├─────────────────────────────────────────────────────┤
│  /api/dns/*  │  /mail/config-v1.1.xml  │  /autodiscover  │
└─────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────┐
│               Business Logic Layer                  │
├──────────────────────┬──────────────────────────────┤
│  internal/dns        │  internal/autoconfig         │
│  - Cloudflare API    │  - Mozilla Autoconfig        │
│  - Record Templates  │  - MS Autodiscover           │
│  - Validation        │  - Apple Mobileconfig        │
│  - Propagation Check │  - go-autoconfig integration │
└──────────────────────┴──────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────┐
│                 Data Layer (SQLite)                 │
├─────────────────────────────────────────────────────┤
│  dns_records  │  dns_checks  │  domains  │  users   │
└─────────────────────────────────────────────────────┘
```

### Integration Points

**Leverages Existing Components:**
- Settings Management API (DNS provider credentials)
- Admin Authentication System (approval workflow)
- Cloudflare API Integration (Let's Encrypt/DNS)
- SQLite Schema (domain/user data)
- Web UI Framework (admin dashboard)
- TLS Certificate Infrastructure (secure endpoints)

## Milestones

### Milestone 1: DNS Management Foundation (2-3 weeks)

**Goal**: Core DNS management infrastructure with Cloudflare integration

**Tasks:**
- **1.1**: Extend Cloudflare API client for DNS record management
  - CRUD operations for SRV, MX, TXT, DKIM records
  - Batch record creation/updates
  - Error handling and retry logic

- **1.2**: Design and implement DNS record templates
  - Mail service SRV records (_imap, _smtp, _submission)
  - MX records with priority configuration
  - DKIM, DMARC, SPF record generation
  - Autodiscover SRV record support

- **1.3**: Build admin approval workflow
  - Propose → Review → Approve → Apply pipeline
  - Role-based access control integration
  - Approval notification system

- **1.4**: DNS validation and propagation checking
  - Pre-deployment validation (syntax, conflicts)
  - Post-deployment propagation monitoring
  - Multi-nameserver validation
  - Rollback capability

**Database Schema:**
```sql
CREATE TABLE dns_records (
    id INTEGER PRIMARY KEY,
    domain_id INTEGER REFERENCES domains(id),
    record_type TEXT NOT NULL, -- SRV, MX, DKIM, DMARC, SPF, TXT
    name TEXT NOT NULL,
    value TEXT NOT NULL,
    ttl INTEGER DEFAULT 3600,
    priority INTEGER,
    status TEXT DEFAULT 'pending', -- pending, approved, active, failed
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    approved_by INTEGER REFERENCES users(id),
    approved_at DATETIME,
    cloudflare_id TEXT -- External record ID for updates/deletes
);

CREATE TABLE dns_checks (
    id INTEGER PRIMARY KEY,
    record_id INTEGER REFERENCES dns_records(id),
    check_time DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT, -- propagated, pending, failed
    nameserver TEXT,
    response TEXT
);
```

**API Endpoints:**
- `POST /api/domains/{id}/dns/propose` - Propose DNS changes
- `GET /api/domains/{id}/dns/preview` - Preview generated records
- `POST /api/admin/dns/approve/{record_id}` - Approve pending change
- `DELETE /api/admin/dns/reject/{record_id}` - Reject proposed change
- `GET /api/domains/{id}/dns/status` - Check propagation status

**Deliverables:**
- `internal/dns` package with Cloudflare integration
- DNS record template generator
- Admin approval workflow engine
- DNS propagation monitoring service

---

### Milestone 2: DNS Management UI (1-2 weeks)

**Goal**: User-facing interface for DNS configuration and monitoring

**Tasks:**
- **2.1**: Domain DNS configuration interface
  - Visual domain DNS overview
  - Record type selection and configuration
  - Bulk record generation (all mail records at once)

- **2.2**: DNS record preview and validation
  - Show generated records before submission
  - Syntax validation with error highlighting
  - Conflict detection with existing records

- **2.3**: Admin approval queue UI
  - Pending changes dashboard
  - Side-by-side diff view (current vs. proposed)
  - Batch approval capability

- **2.4**: DNS propagation status monitoring
  - Real-time propagation status display
  - Per-nameserver check results
  - Visual propagation timeline

**UI Components:**
- DNS Configuration Dashboard (Vue.js)
- Record Editor with validation
- Approval Queue with filtering/sorting
- Propagation Status Widget with live updates

**Deliverables:**
- Vue.js components for DNS management
- Real-time WebSocket updates for propagation status
- Admin approval dashboard
- User documentation for DNS features

---

### Milestone 3: Autoconfiguration Infrastructure (2-3 weeks)

**Goal**: Zero-configuration client setup for major email platforms

**Tasks:**
- **3.1**: Integrate go-autoconfig library
  - Dependency management and integration
  - Configuration mapping from SQLite domain data
  - Protocol-specific response generation

- **3.2**: Design autoconfig data model
  - Domain → client config mapping
  - Server capabilities detection (IMAP/SMTP/CalDAV)
  - TLS/SSL configuration detection
  - Port and hostname resolution

- **3.3**: Mozilla Autoconfig endpoint implementation
  - `/mail/config-v1.1.xml?emailaddress={email}` endpoint
  - `/.well-known/autoconfig/mail/config-v1.1.xml` alternative
  - XML response generation with proper MIME types
  - Email address validation and domain extraction

- **3.4**: Microsoft Autodiscover endpoint implementation
  - `/autodiscover/autodiscover.xml` POST endpoint
  - XML request parsing and validation
  - Autodiscover response generation
  - Exchange protocol compatibility

- **3.5**: Apple Mobileconfig generation
  - `/mobileconfig/{email}` endpoint with authentication
  - .mobileconfig XML profile generation
  - Email, CalDAV, CardDAV configuration
  - Signed profile support (optional)

**Protocol Requirements:**

**Mozilla Autoconfig XML:**
```xml
<clientConfig version="1.1">
  <emailProvider id="example.com">
    <domain>example.com</domain>
    <displayName>Example Mail</displayName>
    <incomingServer type="imap">
      <hostname>mail.example.com</hostname>
      <port>993</port>
      <socketType>SSL</socketType>
      <authentication>password-cleartext</authentication>
      <username>%EMAILADDRESS%</username>
    </incomingServer>
    <outgoingServer type="smtp">
      <hostname>mail.example.com</hostname>
      <port>587</port>
      <socketType>STARTTLS</socketType>
      <authentication>password-cleartext</authentication>
      <username>%EMAILADDRESS%</username>
    </outgoingServer>
  </emailProvider>
</clientConfig>
```

**Microsoft Autodiscover Response:**
```xml
<Autodiscover>
  <Response>
    <Account>
      <AccountType>email</AccountType>
      <Action>settings</Action>
      <Protocol>
        <Type>IMAP</Type>
        <Server>mail.example.com</Server>
        <Port>993</Port>
        <SSL>on</SSL>
        <AuthRequired>on</AuthRequired>
      </Protocol>
    </Account>
  </Response>
</Autodiscover>
```

**API Endpoints:**
- `GET /mail/config-v1.1.xml?emailaddress={email}` - Mozilla Autoconfig
- `GET /.well-known/autoconfig/mail/config-v1.1.xml?emailaddress={email}` - Mozilla alt
- `POST /autodiscover/autodiscover.xml` - Microsoft Autodiscover
- `GET /mobileconfig/{email}?token={auth_token}` - Apple Mobileconfig

**Deliverables:**
- `internal/autoconfig` package
- Mozilla Autoconfig endpoint handler
- Microsoft Autodiscover endpoint handler
- Apple Mobileconfig generator
- go-autoconfig library integration

---

### Milestone 4: Security & Testing (2 weeks)

**Goal**: Comprehensive security audit and multi-client compatibility validation

**Tasks:**
- **4.1**: Security audit for autoconfig endpoints
  - Authentication and authorization review
  - Rate limiting implementation (prevent enumeration)
  - Information disclosure assessment
  - TLS enforcement validation
  - Input sanitization and validation

- **4.2**: Client compatibility testing matrix
  - Thunderbird (Linux, Windows, macOS)
  - Microsoft Outlook (Windows, macOS)
  - Apple Mail (macOS, iOS)
  - K-9 Mail (Android)
  - Nine (iOS/Android)
  - Manual fallback documentation

- **4.3**: DNS validation and monitoring tests
  - Record creation/update/delete tests
  - Propagation monitoring accuracy
  - Cloudflare API error handling
  - Rollback functionality validation

- **4.4**: Integration testing for complete onboarding flow
  - End-to-end domain setup with DNS
  - Client autoconfiguration verification
  - Error handling and recovery paths
  - Performance and load testing

**Security Checklist:**
- [ ] Rate limiting on all autoconfig endpoints (10 req/min per IP)
- [ ] Email address validation prevents enumeration
- [ ] TLS 1.2+ required for all autoconfig endpoints
- [ ] Authentication tokens for mobileconfig downloads
- [ ] Audit logging for all DNS changes and autoconfig requests
- [ ] No internal hostname disclosure in responses
- [ ] CORS headers properly configured
- [ ] Input sanitization for all user-provided data

**Testing Matrix:**

| Client | Platform | Protocol | Status |
|--------|----------|----------|--------|
| Thunderbird 115+ | Linux | Mozilla Autoconfig | ⏳ Pending |
| Thunderbird 115+ | Windows | Mozilla Autoconfig | ⏳ Pending |
| Thunderbird 115+ | macOS | Mozilla Autoconfig | ⏳ Pending |
| Outlook 2021 | Windows | Autodiscover | ⏳ Pending |
| Outlook | macOS | Autodiscover | ⏳ Pending |
| Apple Mail | macOS 14+ | Mobileconfig | ⏳ Pending |
| Apple Mail | iOS 17+ | Mobileconfig | ⏳ Pending |
| K-9 Mail | Android | Mozilla Autoconfig | ⏳ Pending |
| Nine | iOS/Android | Autodiscover | ⏳ Pending |

**Deliverables:**
- Security audit report with findings
- Client compatibility test results
- Automated test suite for DNS and autoconfig
- Performance benchmarks and optimization recommendations

---

## Technical Specifications

### DNS Record Templates

**Example Domain: example.com**

```
# SRV Records for Service Discovery
_imap._tcp.example.com          SRV  10  20  143  mail.example.com.
_imaps._tcp.example.com         SRV  0   1   993  mail.example.com.
_submission._tcp.example.com    SRV  10  20  587  mail.example.com.
_autodiscover._tcp.example.com  SRV  0   0   443  mail.example.com.

# MX Record
example.com                     MX   10  mail.example.com.

# SPF Record
example.com                     TXT  "v=spf1 mx ~all"

# DMARC Record
_dmarc.example.com              TXT  "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"

# DKIM Record (generated per domain)
default._domainkey.example.com  TXT  "v=DKIM1; k=rsa; p=MIGfMA0GCS..."
```

### Cloudflare API Integration

**Required Operations:**
- `POST /zones/{zone_id}/dns_records` - Create DNS record
- `GET /zones/{zone_id}/dns_records` - List DNS records
- `PUT /zones/{zone_id}/dns_records/{record_id}` - Update DNS record
- `DELETE /zones/{zone_id}/dns_records/{record_id}` - Delete DNS record

**Authentication:**
- API Token (preferred) or Global API Key
- Stored in settings database (encrypted)
- Scoped to DNS edit permissions only

### go-autoconfig Integration

**Library Reference**: https://github.com/philband/go-autoconfig

**Key Functions:**
- Generate Mozilla Autoconfig XML responses
- Generate Microsoft Autodiscover XML responses
- Generate Apple Mobileconfig profiles
- Parse and validate email addresses
- Detect domain from email address

**Configuration Mapping:**
```go
type AutoconfigSettings struct {
    Domain           string
    DisplayName      string
    IMAPHost         string
    IMAPPort         int
    IMAPSecurity     string // SSL, STARTTLS
    SMTPHost         string
    SMTPPort         int
    SMTPSecurity     string
    CalDAVHost       string
    CalDAVPort       int
    CardDAVHost      string
    CardDAVPort      int
    Authentication   string // password-cleartext, OAuth2
}
```

---

## Security Considerations

### DNS Management Security

**Authorization:**
- Only admins can propose DNS changes
- Separate role for DNS approval (configurable)
- Audit trail of all DNS operations

**Validation:**
- Pre-deployment syntax validation
- Conflict detection with existing records
- Cloudflare API error handling and retry logic

**Rollback:**
- Previous DNS state stored before changes
- One-click rollback for failed deployments
- Automatic rollback on propagation failure (optional)

**Rate Limiting:**
- Max 100 DNS operations per hour per admin
- Cloudflare API rate limit handling (1200 req/5min)

### Autoconfiguration Security

**Authentication:**
- Email ownership validation (optional token-based)
- Rate limiting: 10 requests per minute per IP
- No authentication required for Mozilla/Autodiscover (per spec)
- Token-based authentication for Mobileconfig downloads

**Information Disclosure:**
- Only return config for valid domains
- Don't reveal internal hostnames or IPs
- Generic error messages (don't confirm user existence)

**TLS Requirements:**
- All autoconfig endpoints require HTTPS
- TLS 1.2+ minimum
- Valid certificate required (Let's Encrypt)

**Privacy:**
- Log autoconfig requests for abuse detection
- No PII stored beyond audit logs
- GDPR compliance for EU users

### Compliance

**GDPR:**
- User consent for autoconfiguration logging
- Data retention policy (30-day audit logs)
- Right to deletion of autoconfig logs

**Security Best Practices:**
- OWASP top 10 review
- Input validation on all endpoints
- SQL injection prevention (parameterized queries)
- XSS prevention in UI

---

## Testing Strategy

### Unit Tests

**DNS Management:**
- Cloudflare API client (mocked responses)
- Record template generation
- Validation logic
- Propagation checker

**Autoconfiguration:**
- Protocol response generation (XML validation)
- Email address parsing
- Configuration mapping from database
- go-autoconfig integration

### Integration Tests

**DNS Workflow:**
- End-to-end DNS record creation
- Approval workflow with multiple admins
- Propagation monitoring
- Rollback scenarios

**Autoconfiguration Workflow:**
- Mozilla Autoconfig endpoint with Thunderbird simulation
- Autodiscover endpoint with Outlook simulation
- Mobileconfig generation and validation

### Client Testing

**Real Client Testing:**
- Thunderbird (latest stable)
- Outlook (Windows/macOS)
- Apple Mail (macOS/iOS)
- K-9 Mail (Android)

**Test Scenarios:**
- New account setup from scratch
- Existing account reconfiguration
- Error handling (invalid email, wrong domain)
- Multiple account configurations

### Performance Testing

**DNS Operations:**
- Bulk record creation (100+ records)
- Concurrent approval operations
- Propagation check performance

**Autoconfiguration:**
- Concurrent autoconfig requests (100 req/sec)
- Response time under load
- Database query optimization

---

## Documentation

### Admin Documentation

**DNS Management:**
- Cloudflare API setup and authentication
- DNS record template customization
- Approval workflow configuration
- Troubleshooting DNS propagation issues

**Autoconfiguration:**
- Enabling autoconfiguration features
- Security configuration (rate limits, tokens)
- Client compatibility notes
- Debugging autoconfig requests

### User Documentation

**DNS Setup:**
- Understanding DNS record requirements
- Requesting DNS changes
- Monitoring DNS propagation

**Client Setup:**
- Supported email clients
- Autoconfiguration URLs
- Manual fallback instructions
- Troubleshooting common issues

### Developer Documentation

**API Documentation:**
- DNS management endpoints
- Autoconfiguration endpoints
- Request/response examples
- Error codes and handling

**Architecture Documentation:**
- System component diagram
- Database schema
- Integration points
- Extension points for multi-provider support

---

## Implementation Order

**Recommended Sequence:**

1. **Week 1-2**: Milestone 1.1-1.2 (DNS foundation, templates)
2. **Week 3**: Milestone 1.3-1.4 (approval workflow, validation)
3. **Week 4-5**: Milestone 2 (DNS UI)
4. **Week 6-7**: Milestone 3.1-3.3 (autoconfig infrastructure, Mozilla)
5. **Week 8**: Milestone 3.4-3.5 (Autodiscover, Mobileconfig)
6. **Week 9-10**: Milestone 4 (security audit, testing)

**Total Estimated Duration:** 10-12 weeks

---

## Future Enhancements

**Multi-Provider Support:**
- Route 53 (AWS)
- Google Cloud DNS
- Azure DNS
- Generic DNS API adapter interface

**Advanced Features:**
- DNSSEC support
- Automated DKIM key rotation
- Let's Encrypt DNS-01 challenge integration
- Terraform/IaC export for DNS configuration

**Enhanced Autoconfiguration:**
- OAuth2 configuration support
- S/MIME certificate distribution
- Enterprise policy enforcement
- Custom branding for mobileconfig profiles

---

## References

- **go-autoconfig**: https://github.com/philband/go-autoconfig
- **automx2 Implementation**: https://rseichter.github.io/automx2/
- **Mozilla Autoconfig Spec**: https://wiki.mozilla.org/Thunderbird:Autoconfiguration
- **Microsoft Autodiscover**: https://docs.microsoft.com/en-us/exchange/client-developer/web-service-reference/autodiscover-web-service-reference-for-exchange
- **Apple Configuration Profile**: https://developer.apple.com/documentation/devicemanagement/mail
- **Cloudflare DNS API**: https://developers.cloudflare.com/api/operations/dns-records-for-a-zone-list-dns-records

---

## Status Tracking

| Milestone | Status | Started | Completed | Notes |
|-----------|--------|---------|-----------|-------|
| M1: DNS Foundation | 📋 Planned | - | - | - |
| M2: DNS UI | 📋 Planned | - | - | - |
| M3: Autoconfig | 📋 Planned | - | - | - |
| M4: Security/Testing | 📋 Planned | - | - | - |

**Legend:**
- 📋 Planned
- 🔄 In Progress
- ✅ Completed
- ⚠️ Blocked
- ❌ Cancelled
