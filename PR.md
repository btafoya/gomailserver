# Mail Server Implementation - Project Requirements

## Project Overview
Build a production-ready, standards-compliant mail server in Go 1.23.5+ with full SMTP, IMAP4, WebDAV (CalDAV/CardDAV) support, comprehensive email security features, SQLite-based storage, and a Gmail-like webmail interface. The system prioritizes ease of installation and management while maintaining enterprise-grade features.

## Design Philosophy
- **Simplicity First**: Single binary deployment with minimal external dependencies
- **SQLite-Powered**: All configuration and metadata in SQLite for easy backup and portability
- **User-Friendly**: Web-based admin interface and user self-service portal from day one
- **Production-Ready**: Built-in security, monitoring, and operational tools

## Core Requirements

### 1. Protocol Support

#### SMTP (Simple Mail Transfer Protocol)
- **Library**: https://github.com/emersion/go-smtp
- Full RFC 5321 compliance
- SMTP submission (port 587) with STARTTLS
- SMTP relay (port 25) with opportunistic TLS
- SMTPS (port 465) support
- Authentication mechanisms: PLAIN, LOGIN, CRAM-MD5
- PIPELINING support (RFC 2920)
- SIZE extension support (RFC 1870)
- 8BITMIME support (RFC 6152)
- CHUNKING support (RFC 3030)
- Message size limits (configurable per domain/user)
- Rate limiting and connection throttling
- Queue management with retry logic
- Bounce handling and DSN (Delivery Status Notifications)

#### IMAP4 (Internet Message Access Protocol v4)
- **Library**: https://github.com/emersion/go-imap
- Full RFC 3501 compliance
- IMAP4rev1 support
- STARTTLS support
- Authentication mechanisms: PLAIN, LOGIN, CRAM-MD5
- IDLE support (RFC 2177) for push notifications
- UIDPLUS extension (RFC 4315)
- QUOTA extension (RFC 2087)
- SORT and THREAD extensions (RFC 5256, RFC 5267)
- NAMESPACE extension (RFC 2342)
- ACL support for shared folders (RFC 4314)
- Special-use mailboxes (RFC 6154): Drafts, Sent, Trash, Spam, etc.
- Mailbox subscriptions
- Message flags: \Seen, \Answered, \Flagged, \Deleted, \Draft, custom flags
- Server-side search capabilities

#### Message Handling
- **Library**: https://github.com/emersion/go-message
- MIME parsing and composition
- Attachment handling
- HTML and plain text message support
- Character encoding support (UTF-8, ISO-8859-1, etc.)
- Message threading
- Message deduplication

#### WebDAV/CalDAV/CardDAV
- CalDAV (RFC 4791) for calendar synchronization
- CardDAV (RFC 6352) for contact synchronization
- WebDAV (RFC 4918) base protocol
- **Client Support (Critical)**: Thunderbird, Apple Mail/Calendar/Contacts, iOS, Android, Microsoft Outlook, Evolution
- Event reminders and invitations
- Calendar sharing and permissions
- Contact groups and distribution lists
- Recurring events support (RFC 5545)
- Free/busy information
- Resource booking (rooms, equipment)
- Event invitations and RSVP handling

#### Sieve Server-Side Filtering (RFC 5228)
- **Critical Feature**: Full Sieve script support
- Rule-based message filtering
- Automatic folder filing
- Conditional forwarding
- Auto-reply rules
- Spam classification rules
- ManageSieve protocol (RFC 5804) for script management
- Web-based rule editor in admin interface

### 2. Email Security Features

#### Antivirus Integration (ClamAV)
- Real-time scanning of incoming and outgoing messages
- Attachment scanning
- Configurable actions per domain/user: reject, quarantine, tag
- Regular signature updates
- Performance optimization with caching
- Scan result logging
- Run on same server as mail server

#### Anti-Spam Integration
- **SpamAssassin Integration**: https://github.com/teamwork/spamc
- Spam scoring and classification
- Bayesian filtering
- Network-based blacklist checking
- User-accessible spam quarantine
- Per-user spam threshold configuration
- Learn from user actions (spam/not spam)
- Spam report generation

#### DKIM (DomainKeys Identified Mail) - RFC 6376
- Outbound email signing
- Inbound signature verification
- Multiple selector support per domain
- Key rotation capabilities
- Configurable signing policies
- Support for RSA-2048 and RSA-4096 keys
- Ed25519 support (RFC 8463)

#### SPF (Sender Policy Framework) - RFC 7208
- Inbound SPF validation
- SPF record checking for all domains
- Configurable handling: none, softfail, fail
- SPF result headers
- IPv4 and IPv6 support

#### DMARC (Domain-based Message Authentication, Reporting & Conformance) - RFC 7489
- DMARC policy enforcement
- Alignment checking (relaxed/strict)
- Aggregate report generation
- Forensic report generation
- Report sending scheduling
- DMARC record parsing

#### DANE (DNS-based Authentication of Named Entities) - RFC 7672
- TLSA record validation
- DANE-TA and DANE-EE support
- DNSSEC validation
- Fallback mechanisms

#### MTA-STS (Mail Transfer Agent Strict Transport Security) - RFC 8461
- Policy enforcement
- Policy caching
- TLS version and cipher enforcement
- TLSRPT (TLS Reporting) support (RFC 8460)

#### PGP/GPG End-to-End Encryption
- Built-in PGP/GPG support
- Key management per user
- Automatic encryption when recipient key available
- Signature verification
- Key import/export via web interface
- Integration with webmail client

#### Additional Security
- TLS 1.2+ enforcement (configurable)
- Modern cipher suite support
- Certificate validation
- STARTTLS everywhere
- Connection encryption logging
- Brute force protection
- IP-based blacklisting/whitelisting
- **Greylisting enabled by default**
- Real-time Blackhole List (RBL) checking
- Rate limiting per user, per IP, per domain
- TOTP-based two-factor authentication (2FA)
- Bcrypt password hashing
- Full audit trail for admin actions

### 3. Database Architecture (SQLite 3)

#### Storage Strategy
**SQLite for Everything**:
- All configuration stored in SQLite
- All metadata (users, domains, folders) in SQLite
- Message metadata in SQLite
- **Hybrid Storage for Large Messages**:
  - Small messages (< 1MB): Store in SQLite BLOB
  - Large messages (≥ 1MB): Store on filesystem, path in SQLite
  - Configurable threshold per domain
- Automatic schema management and migrations
- Single file for easy backup: `mailserver.db`
- WAL (Write-Ahead Logging) mode for better concurrency

#### Schema Requirements

**Domains Table**
- Domain name (primary identifier)
- Status (active/inactive)
- Creation date, modification date
- DKIM keys and selectors
- SPF configuration
- DMARC policy
- Catchall email address
- Max users per domain
- Max mailbox size per domain
- Default quota
- Backup MX support

**Subdomains Table**
- Subdomain name
- Parent domain reference
- Independent or inherited settings

**Users Table**
- Email address (primary key)
- Domain reference (foreign key)
- Password hash (bcrypt/argon2)
- Full name
- Display name
- Mailbox quota
- Current mailbox usage
- Status (active/disabled/suspended)
- Creation date, last login
- Authentication method
- Two-factor authentication settings
- Email forwarding rules
- Auto-reply/vacation settings
- Spam threshold
- Language preference

**Aliases Table**
- Alias address
- Destination addresses (comma-separated or JSON array)
- Domain reference
- Status (active/inactive)
- Creation date

**Messages Table**
- Message ID (UUID or auto-increment)
- User reference (foreign key)
- Mailbox/folder
- Flags (\Seen, \Deleted, etc.)
- Size
- Received date
- Internal date
- Subject (indexed)
- Sender (indexed)
- Recipients
- Message headers (JSON or text)
- Body structure
- Storage type (blob/file)
- Message content (BLOB for small) or filesystem path (TEXT for large)
- Gmail-like categories/labels (JSON array)
- Thread ID for conversation view

**Mailboxes/Folders Table**
- Mailbox ID
- User reference
- Name (INBOX, Sent, Drafts, etc.)
- Parent folder (for hierarchy)
- Subscribed status
- Special use flag
- UIDVALIDITY
- UIDNEXT

**SMTP Queue Table**
- Queue ID
- Sender
- Recipients
- Message content reference
- Retry count
- Next retry time
- Status (pending/processing/failed)
- Error message
- Created timestamp

**Logs Table**
- Timestamp
- Log level
- Service (SMTP/IMAP/CalDAV/etc.)
- User/IP address
- Action
- Result
- Message

**Sessions Table**
- Session ID
- User reference
- Protocol (SMTP/IMAP/WebDAV)
- IP address
- Start time
- Last activity
- Connection state

**Security Tables**
- Failed login attempts
- Blacklisted IPs
- Whitelisted IPs
- DKIM keys per domain
- TLS certificates (or Let's Encrypt metadata)
- User 2FA tokens
- User PGP/GPG public keys

**Sieve Scripts Table**
- Script ID
- User reference
- Script name
- Script content
- Active status
- Creation/modification dates

**Spam Quarantine Table**
- Message ID reference
- User reference
- Quarantine date
- Spam score
- Release status
- Auto-delete date

**Webhooks Table**
- Webhook ID
- Domain reference
- URL endpoint
- Events to trigger on
- Authentication token
- Active status
- Retry configuration

#### Database Performance
- Proper indexing on frequently queried fields
- Foreign key constraints
- Prepared statements
- Connection pooling (limited for SQLite)
- Query optimization
- WAL mode for better concurrent access
- PRAGMA optimizations (synchronous=NORMAL, cache_size, etc.)
- Regular VACUUM and ANALYZE
- Automatic cleanup of old logs and deleted messages

#### Backup Strategy
- Simple file-based backup (copy `mailserver.db`)
- Built-in backup command with hot backup support
- Incremental backup using SQLite backup API
- Export to SQL for portability
- Daily automated backups (configurable)
- 30-day retention by default

### 4. Scalability Requirements

#### Unlimited Support
- **Domains**: No hard-coded limits on domain count
- **Subdomains**: Unlimited subdomains per domain
- **Users**: Scalable user management (tested up to 10,000 users)
- **Aliases**: Unlimited aliases per domain/user
- **Storage**: Configurable per user/domain, no global limits

#### Performance Targets
- Handle hundreds of concurrent connections
- Process 100,000+ emails per day
- Sub-second IMAP response times
- Efficient memory usage (< 512MB for typical workloads)
- Single-server architecture (horizontal scaling in future versions)
- Optimized for VPS/bare metal deployment

#### SQLite Scalability Considerations
- SQLite performs well for small-to-medium deployments (< 100 domains, < 1000 users)
- For larger deployments, hybrid storage mitigates blob size issues
- Read-heavy workloads (IMAP) benefit from SQLite's performance
- Write bottlenecks addressed through WAL mode and connection pooling
- Can migrate to PostgreSQL in future if needed (modular design)

### 5. Configuration Management

#### Configuration Storage
- All configuration in SQLite (no config files for domains/users)
- Minimal system configuration file (YAML) for:
  - SQLite database path
  - Port bindings
  - TLS mode (Let's Encrypt or manual)
  - ClamAV socket path
  - SpamAssassin connection
- Environment variable support for Docker deployment
- Hot-reload capabilities for certain settings via API

#### Configurable Items
- Port bindings
- TLS certificates (Let's Encrypt ACME + Cloudflare DNS by default)
- Automatic certificate renewal
- Per-domain certificate support (SNI)
- Timeout values
- Resource limits
- Feature toggles
- Logging levels
- SQLite database settings
- ClamAV socket path
- SpamAssassin connection
- Backup settings
- Message size thresholds (blob vs file)
- Quota warning thresholds

### 6. Web Interface Requirements

#### Admin Web Interface (Critical - Phase 1)
- Modern, responsive web UI
- Domain management (CRUD)
- User management (CRUD)
- Alias management (CRUD)
- Quota management and visualization
- Real-time statistics dashboard
- Log viewer with filtering
- Queue management interface
- Security settings (DKIM, SPF, DMARC per domain)
- TLS certificate status and management
- Backup/restore interface
- System health monitoring
- Role-based access control (admin/read-only)

#### User Self-Service Portal (Critical - Phase 1)
- Password change
- 2FA setup (TOTP)
- Alias management (create/delete own aliases)
- Quota usage display
- Forwarding rules
- Sieve filter management (visual editor)
- PGP key management
- Spam quarantine review
- Session management
- Email signature settings

#### Webmail Client (Critical - Phase 2)
- **Gmail-like interface** with modern design
- **Categories/Labels system** (Primary, Social, Promotions, etc.)
- Conversation/thread view
- Rich text editor (HTML email composition)
- Attachment handling (drag-and-drop, inline images)
- Contact management (integrated with CardDAV)
- Calendar view (integrated with CalDAV)
- Search with advanced filters
- Keyboard shortcuts
- Multiple account support
- Dark mode
- Mobile-responsive design
- Offline capability (Progressive Web App)
- PGP/GPG integration (encrypt/decrypt/sign)
- Spam reporting
- Message templates
- Auto-save drafts

### 7. API Requirements

#### Management REST API
- JSON-based REST API
- Domain management (CRUD)
- User management (CRUD)
- Alias management (CRUD)
- Quota management
- Statistics and monitoring
- Log retrieval
- Queue management
- Authentication (API keys, JWT tokens)
- Rate limiting
- OpenAPI/Swagger documentation

#### Webhook Support (Critical)
- Email received notifications
- Email sent notifications
- Delivery status notifications (success/failure)
- Security event notifications (failed login, blocked IP)
- Quota warnings (80%, 90%, 100%)
- Configurable per domain
- Retry logic with exponential backoff
- Webhook testing interface in admin UI

### 8. Monitoring and Logging

#### Metrics
- Connection counts (SMTP/IMAP/WebDAV)
- Message throughput
- Queue depth
- Database query performance
- Error rates
- Authentication success/failure rates
- Storage usage per user/domain
- TLS usage statistics

#### Logging
- Structured logging (JSON format)
- Log levels: DEBUG, INFO, WARN, ERROR
- Separate logs per service
- Log rotation
- Integration with log aggregation tools (syslog, ELK, etc.)

#### Health Checks
- Service health endpoints
- Database connectivity
- External service checks (ClamAV, DNS)
- Queue status
- Disk space monitoring

#### Sieve Filtering (RFC 5228)
- RFC 5228 (Sieve base)
- RFC 5229 (Variables)
- RFC 5230 (Vacation)
- RFC 5231 (Relational)
- RFC 5233 (Subaddress)
- RFC 5235 (Spamtest)

### 9. Standards Compliance

#### Email RFCs
- RFC 5321 (SMTP)
- RFC 5322 (Internet Message Format)
- RFC 3501 (IMAP4rev1)
- RFC 2045-2049 (MIME)
- RFC 6376 (DKIM)
- RFC 7208 (SPF)
- RFC 7489 (DMARC)
- RFC 7672 (DANE)
- RFC 8461 (MTA-STS)
- RFC 2142 (Mailbox Names)
- RFC 6409 (Message Submission)

#### CalDAV/CardDAV RFCs
- RFC 4791 (CalDAV)
- RFC 6352 (CardDAV)
- RFC 4918 (WebDAV)
- RFC 5545 (iCalendar)
- RFC 6350 (vCard)

### 10. Code Quality Requirements

#### Go Best Practices
- Go 1.23.5+ features
- Idiomatic Go code
- Comprehensive error handling
- Context usage for cancellation
- Graceful shutdown
- Unit tests (80%+ coverage)
- Integration tests
- Benchmark tests for critical paths
- Race condition detection

#### Code Organization
- Clean architecture/hexagonal architecture
- Dependency injection
- Interface-based design
- Repository pattern for data access
- Service layer for business logic
- Clear separation of concerns

#### Documentation
- Godoc comments for all exported types/functions
- README with setup instructions
- Architecture documentation
- API documentation
- Deployment guide
- Troubleshooting guide

### 11. Deployment Requirements

#### Build
- Single binary compilation (Go 1.23.5+)
- Embed web UI assets in binary
- Cross-platform support (Linux primary, macOS, Windows)
- Docker container with Alpine base
- Docker Compose for easy deployment
- Minimal external dependencies (ClamAV, DNS server)
- SQLite included (no external database server needed)

#### Installation
- **One-line installation script** for Ubuntu/Debian
- Automatic SQLite database creation
- Built-in database migration tool
- **Web-based initial setup wizard**:
  - First domain configuration
  - First admin user creation
  - TLS certificate setup (Let's Encrypt auto-config)
  - DKIM key generation
  - DNS record suggestions
- Configuration validation
- Pre-flight checks (ports, ClamAV, DNS)
- Systemd service file generation

#### Backup and Recovery
- Built-in backup command: `mailserver backup`
- Automatic daily backups (configurable)
- Single file backup (SQLite + message directory)
- Restore command: `mailserver restore <backup-file>`
- Export/import functionality
- 30-day backup retention
- Backup integrity verification

## Success Criteria

1. **Email Functionality**: Send and receive emails reliably
2. **Security**: DKIM, SPF, DMARC, DANE, MTA-STS all functional
3. **CalDAV/CardDAV**: Sync working with Thunderbird, Apple, iOS, Android, Outlook, Evolution
4. **Anti-Virus/Spam**: ClamAV and SpamAssassin scanning all messages
5. **Sieve Filtering**: Server-side rules working correctly
6. **PGP/GPG**: Encryption and signing functional in webmail
7. **Scalability**: Tested with 100 domains and 200 users
8. **Performance**: Sub-second IMAP response times, handle 100,000 emails/day
9. **Web Interfaces**: Admin UI, user portal, and webmail all functional
10. **TLS**: Let's Encrypt automatic certificate management working
11. **Testing**: All unit and integration tests passing
12. **Documentation**: Complete installation, admin, and user guides
13. **Deliverability**: Passes mail-tester.com score 10/10
14. **Installation**: Can be installed from scratch in < 30 minutes

## Revised Timeline Estimate

### Phase 1: Core Mail Server (3-4 weeks)
- SMTP server (send/receive)
- IMAP server (read emails)
- SQLite database schema
- User authentication
- Basic folder management
- TLS support
- Message storage (hybrid blob/file)
- Queue management

### Phase 2: Security Foundation (2-3 weeks)
- DKIM signing and verification
- SPF validation
- DMARC enforcement
- ClamAV integration
- SpamAssassin integration
- Greylisting
- Rate limiting
- 2FA (TOTP)

### Phase 3: Web Interfaces - Admin & Portal (4-5 weeks)
- REST API foundation
- Admin web UI (domain/user management)
- User self-service portal
- Dashboard and statistics
- Sieve filter visual editor
- System monitoring interface
- Let's Encrypt ACME integration
- Initial setup wizard

### Phase 4: CalDAV/CardDAV (2-3 weeks)
- CalDAV server (RFC 4791)
- CardDAV server (RFC 6352)
- Event and contact storage
- Client compatibility testing
- Calendar sharing
- Recurring events
- Free/busy information

### Phase 5: Advanced Security (2 weeks)
- DANE (TLSA validation)
- MTA-STS policy enforcement
- PGP/GPG key management
- End-to-end encryption
- Advanced audit logging

### Phase 6: Sieve & Filtering (1-2 weeks)
- Sieve interpreter
- ManageSieve protocol
- Common filter templates
- Integration with spam detection
- Testing with various rules

### Phase 7: Webmail Client (4-6 weeks)
- Modern responsive UI framework
- Mailbox view (list/thread)
- Message composition (rich text)
- Gmail-like categories/labels
- Attachment handling
- Search functionality
- Contact integration
- Calendar integration
- PGP integration in webmail
- Keyboard shortcuts
- Mobile responsiveness

### Phase 8: Webhooks & Integration (1 week)
- Webhook framework
- Event triggers
- Retry logic
- Testing interface

### Phase 9: Polish & Documentation (2-3 weeks)
- Integration testing
- Performance optimization
- Installation scripts
- Comprehensive documentation
- User guides
- Video tutorials
- Docker images
- Example configurations

### Phase 10: Real-World Testing (1-2 weeks)
- Deliverability testing (mail-tester.com)
- Client compatibility testing
- Load testing
- Security audit
- Bug fixes

**Total Estimated Time**: 22-31 weeks (5-7 months)

**Minimum Viable Product (MVP)**: Phases 1-3 (9-12 weeks)
- Can send/receive email securely
- Basic spam/virus protection
- Admin web interface
- Suitable for small personal/business use

## Dependencies

### Required
- Go 1.23.5 or higher (build time)
- SQLite 3 (embedded, no separate installation)
- ClamAV daemon (clamd)
- SpamAssassin (spamd)
- DNS server access (for DKIM, SPF, DMARC, DANE setup)
- Cloudflare account (for Let's Encrypt DNS validation)

### Optional
- Redis (for session storage in clustered deployments - future)
- Reverse proxy (nginx/caddy) for additional security

## Out of Scope

### Definitely Out of Scope
- POP3 protocol support (use IMAP)
- Built-in DNS server
- Clustering/high availability (single instance architecture)
- S3 or object storage backend (filesystem only)
- LDAP/Active Directory integration
- Multi-language support (English only initially)

### Future Considerations (Post-MVP)
- Mobile apps (iOS/Android)
- Desktop apps
- Migration tools from other mail servers
- PostgreSQL backend option (for larger deployments)
- Advanced spam filtering (AI/ML-based)
- Mailing list management
- Auto-reply/vacation messages
- Email archiving/compliance features
- SSO integration (OAuth, SAML)
