# gomailserver - Administrator Guide

**Version**: 1.0  
**Date**: January 11, 2026  
**gomailserver Version**: 0.10.0+

---

## Table of Contents

1. [Overview](#overview)
2. [Installation](#installation)
3. [Initial Setup](#initial-setup)
4. [System Configuration](#system-configuration)
5. [Domain Management](#domain-management)
6. [User Management](#user-management)
7. [Security Configuration](#security-configuration)
8. [TLS Certificates](#tls-certificates)
9. [Reputation Management](#reputation-management)
10. [Monitoring and Logging](#monitoring-and-logging)
11. [Backup and Recovery](#backup-and-recovery)
12. [Troubleshooting](#troubleshooting)
13. [Production Best Practices](#production-best-practices)

---

## Overview

gomailserver is a modern, all-in-one mail server written in Go that replaces complex mail server stacks (Postfix, Dovecot, OpenDKIM, etc.) with a single daemon. This guide covers system administration tasks, configuration, and maintenance procedures for production deployments.

### Key Features

- **Single Binary**: No complex dependencies, easy deployment
- **SQLite Backend**: Single database file for simple backup
- **Web Interface**: Unified admin/portal/webmail (Nuxt.js)
- **Security**: DKIM, SPF, DMARC, DANE, MTA-STS, ClamAV, SpamAssassin
- **Reputation**: Automated sender reputation with external feedback (Gmail, Microsoft SNDS)
- **Auto-Configuration**: Let's Encrypt ACME with Cloudflare DNS
- **Testing Tools**: gomailtest for production verification

---

## Installation

### Requirements

- **OS**: Debian 12 (bookworm), Ubuntu 24.04, or compatible
- **Architecture**: amd64 or arm64
- **Memory**: 2GB RAM minimum, 4GB recommended for production
- **Storage**: 10GB minimum for production
- **Network**: Static IP recommended for production mail servers
- **DNS**: Cloudflare account recommended for ACME (Let's Encrypt)
- **Software**:
  - Go 1.23.5+ (build time only)
  - ClamAV daemon (clamd)
  - SpamAssassin daemon (spamd)

### Option 1: APT Package (Recommended)

```bash
# Add GPG key
curl -fsSL https://btafoya.github.io/gomailserver/repo/public.key | \
  sudo gpg --dearmor -o /usr/share/keyrings/gomailserver-archive-keyring.gpg

# Add repository (replace 'jammy' with your distribution)
echo "deb [signed-by=/usr/share/keyrings/gomailserver-archive-keyring.gpg] https://btafoya.github.io/gomailserver/repo jammy main" | \
  sudo tee /etc/apt/sources.list.d/gomailserver.list

# Install
sudo apt update
sudo apt install gomailserver
```

### Option 2: Build from Source

```bash
# Clone repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Build
make build

# Install
sudo install -m 0755 build/gomailserver /usr/local/bin/
```

### Option 3: systemd Installation

```bash
# Build binary
make build

# Run installer
sudo ./scripts/install-systemd.sh --start

# The installer creates:
# - User and group: gomailserver
# - Binary: /usr/local/bin/gomailserver
# - Directories: /var/lib/gomailserver, /var/log/gomailserver
# - Config: /etc/gomailserver
# - Systemd service
```

### Service Management

```bash
# Start service
sudo systemctl start gomailserver

# Stop service
sudo systemctl stop gomailserver

# Restart service
sudo systemctl restart gomailserver

# Enable on boot
sudo systemctl enable gomailserver

# Check status
sudo systemctl status gomailserver

# View logs
sudo journalctl -u gomailserver -f
```

---

## Initial Setup

### 1. Access Setup Wizard

After first installation, access the setup wizard at:

```
http://your-server-ip:8980/admin/setup
```

The setup wizard will guide you through:
1. **System Configuration** - Hostname, ports, storage paths
2. **Domain Setup** - Primary domain configuration
3. **Admin User Creation** - First administrative account

### 2. Using gomailserver-control Script

```bash
# Start in development mode (uses ./gomailserver.yaml)
./scripts/gomailserver-control.sh start --dev

# Start in production mode (uses /etc/gomailserver/gomailserver.yaml)
./scripts/gomailserver-control.sh start

# Check status
./scripts/gomailserver-control.sh status

# Stop server
./scripts/gomailserver-control.sh stop
```

### 3. Creating First Admin User via CLI

```bash
./build/gomailserver create-admin
```

Follow prompts for:
- Email address (e.g., admin@example.com)
- Full name (e.g., System Administrator)
- Password (strong password required)

---

## System Configuration

### Main Configuration File

Primary configuration is stored in:
- **Production**: `/etc/gomailserver/gomailserver.yaml`
- **Development**: `./gomailserver.yaml`

Example configuration:

```yaml
server:
  hostname: mail.example.com
  domain: example.com
  listen_addr: 0.0.0.0:25

smtp:
  enabled: true
  port: 25
  submission_port: 587
  smtps_port: 465

imap:
  enabled: true
  port: 143
  imaps_port: 993

http:
  enabled: true
  port: 8980

database:
  path: /var/lib/gomailserver/mailserver.db
  backup_enabled: true
  backup_interval: 86400  # 24 hours

tls:
  acme:
    enabled: true
    email: admin@example.com
    provider: cloudflare
    api_token: your_cloudflare_api_token
  cert_file: /etc/ssl/certs/gomailserver.crt
  key_file: /etc/ssl/private/gomailserver.key

security:
  dkim_enabled: true
  spf_enabled: true
  dmarc_enabled: true
  clamav_enabled: true
  spamassassin_enabled: true
  totp_enabled: true
```

### Configuration Sections

#### Server Section
- `hostname`: Fully qualified domain name (FQDN)
- `domain`: Primary mail domain
- `listen_addr`: IP:port to bind SMTP

#### SMTP Section
- `port`: Port 25 (relay/inbound)
- `submission_port`: Port 587 (authenticated submission)
- `smtps_port`: Port 465 (implicit TLS)

#### IMAP Section
- `port`: Port 143 (STARTTLS)
- `imaps_port`: Port 993 (implicit TLS)

#### HTTP Section
- `port`: Admin/API/Portal/Webmail UI port (default: 8980)

#### TLS Section
- `acme.enabled`: Enable automatic Let's Encrypt certificates
- `acme.email`: Email for Let's Encrypt notifications
- `acme.provider`: DNS provider (cloudflare, digitalocean, etc.)
- `acme.api_token`: API token for DNS provider

---

## Domain Management

### Adding a Domain

**Via Web UI**:
1. Access Admin UI: `http://server:8980/admin/domains`
2. Click "Add Domain"
3. Enter domain name (e.g., example.com)
4. Configure limits and security settings
5. Click "Save"

**Via API**:
```bash
curl -X POST http://server:8980/api/v1/domains \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "description": "Primary mail domain"
  }'
```

### Domain Configuration

After adding a domain, configure:

#### DNS Records Required

```
# A Record (MX)
example.com.  IN MX 10 mail.example.com.

# A Record (Server)
mail.example.com.  IN A your-server-ip-address

# TXT Record (SPF)
example.com.  IN TXT "v=spf1 ip4:your-server-ip -all"

# TXT Record (DKIM)
selector1._domainkey.example.com.  IN TXT "v=DKIM1; k=rsa; p=MIGfMA0GCs..."

# TXT Record (DMARC)
_dmarc.example.com.  IN TXT "v=DMARC1; p=none; rua=mailto:dmarc@example.com"
```

#### Security Settings

- **DKIM**: Generate and configure DKIM keys
- **SPF**: Configure SPF policy
- **DMARC**: Set up DMARC reporting
- **DANE/MTA-STS**: Optional additional security layers

### Generating DKIM Keys

**Via Web UI**:
1. Access domain settings: `/admin/domains/example.com`
2. Click "Generate DKIM"
3. Select key size (2048 or 4096 bits)
4. Click "Generate"

**Via API**:
```bash
curl -X POST http://server:8980/api/v1/domains/example.com/dkim \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "selector": "selector1",
    "key_size": 2048
  }'
```

### Testing Domain Configuration

```bash
# Run domain security audit
./build/gomailtest verify --profile production \
  --category security \
  --report-html security-audit.html
```

---

## User Management

### Creating Users

**Via Web UI**:
1. Access: `http://server:8980/admin/users`
2. Click "Add User"
3. Enter user details:
   - Email address (e.g., user@example.com)
   - Full name
   - Password (strong password required)
   - Domain assignment
   - Quota (optional)
4. Click "Save"

**Via API**:
```bash
curl -X POST http://server:8980/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "full_name": "John Doe",
    "password": "SecurePassword123!",
    "domain_id": 1,
    "quota_mb": 1024
  }'
```

### User Roles

- **admin**: Full system access
- **user**: Normal mail access (default for users without role)

### User Quotas

Set quotas to manage disk usage:
- **Unlimited**: Set to 0 or leave blank
- **Size-based**: Set in megabytes (e.g., 1024 MB = 1 GB)
- **Mailbox quotas**: Limit per-mailbox storage

### Managing User Passwords

Users can reset their own passwords via Portal (`/portal/`):
1. User logs in to portal
2. Navigates to "Account Settings"
3. Requests password reset (email verification required)

Administrators can reset passwords via Admin UI:
1. Access user list
2. Click user email
3. Click "Reset Password"
4. Enter new password or send reset link

### Deleting Users

**Warning**: Deleting a user removes:
- All email data
- Mailboxes and messages
- Cannot be undone

Always archive user data before deletion.

---

## Security Configuration

### DKIM (DomainKeys Identified Mail)

DKIM provides email authentication by signing outbound messages.

**Key Configuration**:
- **Key Sizes**: 2048 bits (recommended) or 4096 bits (high security)
- **Rotation**: Rotate keys every 6-12 months
- **Selectors**: Use descriptive selectors (e.g., `mail1`, `mail2024`)

**DNS Record**:
```
selector._domainkey.example.com.  IN TXT "v=DKIM1; k=rsa; p=..." " DKIM public key
```

### SPF (Sender Policy Framework)

SPF prevents sender spoofing by listing authorized mail servers.

**Example SPF Records**:
```
# Simple (all mail from this server)
example.com.  IN TXT "v=spf1 ip4:your-server-ip -all"

# Include multiple servers
example.com.  IN TXT "v=spf1 ip4:your-server-ip include:_spf.google.com include:spf1.protection.outlook.com ~all"

# Third-party mail providers
example.com.  IN TXT "v=spf1 include:mailgun.org include:sendgrid.net ~all"
```

### DMARC

DMARC enforces SPF and DKIM policies and provides reporting.

**Example DMARC Record**:
```
_dmarc.example.com.  IN TXT "v=DMARC1; p=none; rua=mailto:dmarc@example.com; ruf=mailto:dmarc@example.com"
```

**Policy Levels**:
- `p=none`: No action (monitor only)
- `p=quarantine`: Quarantine suspicious messages
- `p=reject`: Reject failing messages (recommended)

### Two-Factor Authentication (2FA)

Enable TOTP for admin users:
1. Access user settings
2. Enable TOTP
3. Scan QR code with authenticator app (Google Authenticator, Authy, etc.)
4. Enter TOTP code when logging in

**Recovery**: Save TOTP recovery codes during setup for account recovery.

### Antivirus (ClamAV)

ClamAV scans email attachments and messages.

**Configuration**:
- Enable via settings: `/admin/settings/security`
- Update virus definitions: ClamAV automatically updates
- Quarantine infected emails: Messages moved to quarantine folder

### Anti-Spam (SpamAssassin)

SpamAssassin filters spam messages.

**Configuration**:
- Enable via settings: `/admin/settings/security`
- Configure thresholds: Adjust spam score levels
- Train filter: Mark messages as spam/not spam to improve accuracy

---

## TLS Certificates

### Automatic Certificates (Let's Encrypt via ACME)

**Cloudflare DNS Configuration**:
1. Log in to Cloudflare dashboard
2. Navigate to "API Tokens"
3. Create token with Zone:DNS and Zone:Read permissions
4. Add token to gomailserver config:

```yaml
tls:
  acme:
    enabled: true
    provider: cloudflare
    api_token: your_cloudflare_api_token
    email: admin@example.com
```

**Certificate Management**:
- Certificates auto-renew 30 days before expiry
- Certificates stored in database
- No manual intervention required after initial setup

### Manual Certificates

For environments without ACME or with specific CA certificates:

1. Upload certificate files to server
2. Update configuration:

```yaml
tls:
  acme:
    enabled: false
  cert_file: /path/to/certificate.crt
  key_file: /path/to/private.key
```

3. Restart service: `sudo systemctl restart gomailserver`

### Self-Signed Certificates (Testing Only)

For development/testing, gomailserver can generate self-signed certificates.

**Warning**: Self-signed certificates will cause warnings in email clients and are NOT for production use.

### Certificate Expiry Monitoring

Certificates are checked daily:
- 30-day warning: Logged but service continues
- 7-day warning: Additional warnings logged
- Expired: Certificate auto-renewal attempted

---

## Reputation Management

gomailserver includes comprehensive automated sender reputation management.

### Reputation Telemetry

Metrics collected per domain:
- **Deliveries**: Successful deliveries
- **Bounces**: Failed deliveries
- **Complaints**: User-reported spam (ARF)
- **Deferrals**: Temporary failures
- **90-day rolling window**: Metrics retained for 90 days

### Reputation Scoring

**Score Range**: 0-100
- **0-49**: Poor - sending blocked
- **50-69**: Fair - limited sending
- **70-89**: Good - normal sending
- **90-100**: Excellent - maximum allowed

**Score Impact**:
- Rate limiting based on reputation
- Circuit breaker activation on poor scores
- Warm-up schedule for new domains

### Circuit Breakers

Automatic sending pause on:
- **High complaint rate**: >0.1%
- **High bounce rate**: >10%
- **Provider blocks**: Gmail, Outlook, Yahoo blocks

**Auto-Resume**: Exponential backoff (1h → 2h → 4h → 8h)

### External Feedback Integration

**Gmail Postmaster Tools**:
1. Add Postmaster Tools account to domain
2. Configure API key in settings
3. Enable automated reporting

**Microsoft SNDS**:
1. Add domain to SNDS
2. Configure access
3. Data collection automatic

### Managing Reputation

**Via Web UI** (`/admin/reputation/`):
- View domain scores
- Check circuit breaker status
- Review DMARC reports
- Manage warm-up schedules
- Configure provider rate limits

**Monitoring Alerts**:
- Reputation score drops
- Circuit breaker activation
- Certificate expiry warnings
- High bounce/complaint rates

---

## Monitoring and Logging

### Log Levels

Configure log verbosity:

```yaml
# In gomailserver.yaml
logging:
  level: info  # debug, info, warn, error
  output: both  # stdout, file, both
```

**Log Levels**:
- `debug`: Detailed information (development only)
- `info`: Normal operations (production recommended)
- `warn`: Warning conditions
- `error`: Error conditions

### Viewing Logs

**Via systemd**:
```bash
# Follow logs in real-time
sudo journalctl -u gomailserver -f

# View last 100 lines
sudo journalctl -u gomailserver -n 100

# View logs for specific time range
sudo journalctl -u gomailserver --since "1 hour ago"
```

**Via Control Script**:
```bash
./scripts/gomailserver-control.sh status
```

### Log Files

Production logs stored in:
- `/var/log/gomailserver/gomailserver.log` (structured JSON)
- `/var/log/gomailserver/access.log` (HTTP access)
- `/var/log/gomailserver/error.log` (errors only)

### Monitoring Webhooks

Configure webhook notifications for events:
- Email deliveries/failures
- Security events
- Reputation alerts
- System errors

**Via API**:
```bash
curl -X POST http://server:8980/api/v1/webhooks \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://your-endpoint.com/webhook",
    "events": ["email.sent", "email.failed", "security.breach"],
    "active": true
  }'
```

---

## Backup and Recovery

### Database Backup

**Manual Backup**:
```bash
# Stop service (optional, for consistency)
sudo systemctl stop gomailserver

# Copy database
sudo cp /var/lib/gomailserver/mailserver.db \
  /backups/mailserver-$(date +%Y%m%d-%H%M%S).db

# Start service
sudo systemctl start gomailserver
```

**Automatic Backup**:
- Enabled by default (every 24 hours)
- Backups stored in `/var/lib/gomailserver/backups/`
- Retain 7 days of backups

### Message Storage

Messages are stored using hybrid strategy:
- **<1MB**: Stored in SQLite database
- **≥1MB**: Stored on filesystem at `/var/lib/gomailserver/messages/`

Backup strategy:
1. Database backup includes all metadata
2. Filesystem backup required for large messages

### Recovery Procedure

**Restoring Database**:
```bash
# Stop service
sudo systemctl stop gomailserver

# Backup current database
sudo cp /var/lib/gomailserver/mailserver.db \
  /var/lib/gomailserver/mailserver.db.broken

# Restore backup
sudo cp /backups/mailserver-20260111-120000.db \
  /var/lib/gomailserver/mailserver.db

# Start service
sudo systemctl start gomailserver
```

**Filesystem Recovery**:
- Restore from backup if `/var/lib/gomailserver/messages/` is lost
- Verify checksums if available

---

## Troubleshooting

### Common Issues

#### Server Won't Start

**Check 1**: Port conflicts
```bash
sudo netstat -tulpn | grep :25
sudo netstat -tulpn | grep :143
sudo netstat -tulpn | grep :587
```

**Check 2**: File permissions
```bash
ls -la /var/lib/gomailserver
ls -la /var/log/gomailserver
```

**Check 3**: Service status
```bash
sudo systemctl status gomailserver
sudo journalctl -xeu gomailserver
```

#### Email Not Delivering

**Check 1**: DNS records
```bash
# Check MX records
dig +short mx example.com

# Check SPF record
dig +short txt example.com

# Check DMARC record
dig +short txt _dmarc.example.com
```

**Check 2**: Blacklists
```bash
# Check if IP is blacklisted
./build/gomailtest verify --category smtp

# Manual check via online tools
# https://mxtoolbox.com/blacklists.aspx
# https://multirbl.valli.org/
```

**Check 3**: Reputation score
```bash
# Check domain reputation
curl http://server:8980/api/v1/reputation/scores/example.com \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Authentication Failures

**Common Causes**:
- Incorrect password
- Account disabled
- TOTP code expired/invalid
- JWT token expired

**Via Web UI**:
1. Access Portal: `/portal/`
2. Log in with email/password
3. If TOTP enabled, enter code from authenticator app
4. Reset password if needed via "Forgot Password?"

#### High Resource Usage

**Check Memory**:
```bash
free -h
```

**Check Disk Space**:
```bash
df -h
du -sh /var/lib/gomailserver
```

**Check CPU**:
```bash
top -p gomailserver
```

### Testing Tools

Use `gomailtest` for production verification:

```bash
# Run all checks
./build/gomailtest verify --profile production

# Run specific categories
./build/gomailtest test config
./build/gomailtest test smtp
./build/gomailtest test imap
./build/gomailtest test mailflow
./build/gomailtest test security

# Generate HTML report
./build/gomailtest verify --profile production --report-html report.html
```

---

## Production Best Practices

### Security

1. **Always use strong passwords** (16+ characters, mixed case, numbers, symbols)
2. **Enable 2FA** for all admin accounts
3. **Keep software updated**: `sudo apt update && sudo apt upgrade`
4. **Monitor logs**: Review logs regularly for suspicious activity
5. **Use HTTPS**: Access admin panel via HTTPS (TLS required)
6. **Firewall**: Restrict access to only necessary ports (25, 143, 587, 993, 8980)

### Reliability

1. **Backup regularly**: Daily automated backups
2. **Test backups**: Verify backups are restorable
3. **Monitor certificates**: Ensure certificates don't expire
4. **Load testing**: Test system capacity before production launch
5. **High availability**: Consider load balancer for multiple servers

### Performance

1. **Monitor reputation**: Keep reputation score high (70+)
2. **Optimize queues**: Don't let queue grow beyond system capacity
3. **Clean regularly**: Remove old logs and messages
4. **Resource limits**: Configure appropriate quotas for users/domains
5. **Monitor resources**: Use monitoring tools for CPU, memory, disk

### Maintenance

1. **Regular updates**: Check for gomailserver updates monthly
2. **Key rotation**: Rotate DKIM keys every 6-12 months
3. **Review reports**: Check DMARC reports weekly
4. **Audit access**: Review admin access logs periodically
5. **Disaster recovery**: Have documented recovery procedures

### Scaling Considerations

**When to Add Servers**:
- Queue size consistently >1000 messages
- Average delivery time >10 seconds
- CPU usage >80% during peak hours
- Single server resource limits reached

**Multi-Server Architecture**:
1. **Load balancer**: Distribute incoming SMTP traffic
2. **Shared database**: Use PostgreSQL instead of SQLite (future)
3. **Queue cluster**: Distribute queue processing
4. **Message storage**: NFS or object storage for large messages

---

## API Reference

### Authentication

**Login**:
```bash
curl -X POST http://server:8980/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "YourPassword"
  }'
```

Response:
```json
{
  "token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "admin@example.com",
    "full_name": "Admin User"
  },
  "expires_at": "2026-01-12T00:00:00Z"
}
```

**Refresh Token**:
```bash
curl -X POST http://server:8980/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_REFRESH_TOKEN" \
  -d '{}'
```

### Domains

**List Domains**:
```bash
curl -X GET http://server:8980/api/v1/domains \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Create Domain**:
```bash
curl -X POST http://server:8980/api/v1/domains \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "description": "Primary mail domain"
  }'
```

**Generate DKIM**:
```bash
curl -X POST http://server:8980/api/v1/domains/1/dkim \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "selector": "mail1",
    "key_size": 2048
  }'
```

### Users

**List Users**:
```bash
curl -X GET http://server:8980/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Create User**:
```bash
curl -X POST http://server:8980/api/v1/users \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "full_name": "John Doe",
    "password": "SecurePassword123!",
    "domain_id": 1,
    "quota_mb": 1024
  }'
```

### Settings

**Get Settings**:
```bash
curl -X GET http://server:8980/api/v1/settings/server \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Update Settings**:
```bash
curl -X PUT http://server:8980/api/v1/settings/server \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "hostname": "mail.example.com",
    "smtp_enabled": true,
    "imap_enabled": true
  }'
```

### Queue

**View Queue**:
```bash
curl -X GET http://server:8980/api/v1/queue \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Retry Message**:
```bash
curl -X POST http://server:8980/api/v1/queue/1/retry \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Statistics

**Dashboard Stats**:
```bash
curl -X GET http://server:8980/api/v1/stats/dashboard \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Domain Stats**:
```bash
curl -X GET http://server:8980/api/v1/stats/domains/1 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Reputation

**Get Domain Score**:
```bash
curl -X GET http://server:8980/api/v1/reputation/scores/example.com \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**List Alerts**:
```bash
curl -X GET http://server:8980/api/v1/reputation/alerts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Support Resources

### Documentation

- **README**: https://github.com/btafoya/gomailserver
- **Reputation Management**: See `REPUTATION-MANAGEMENT.md`
- **Configuration Examples**: See `gomailserver.example.yaml`

### Community

- **GitHub Issues**: https://github.com/btafoya/gomailserver/issues
- **Discussions**: https://github.com/btafoya/gomailserver/discussions
- **Wiki**: https://github.com/btafoya/gomailserver/wiki

### Testing

- **gomailtest**: Built-in production verification tool
- Run: `./build/gomailtest verify --help`

---

## Appendix

### Default Ports

| Service | Port | Protocol | Notes |
|---------|-------|----------|--------|
| SMTP (relay/inbound) | 25 | TCP | Standard SMTP |
| SMTP (submission) | 587 | TCP | Authenticated submission |
| SMTPS (implicit TLS) | 465 | TCP | Legacy SMTPS |
| IMAP (STARTTLS) | 143 | TCP | Standard IMAP |
| IMAPS (implicit TLS) | 993 | TCP | Legacy IMAPS |
| HTTP (Admin/API/Portal) | 8980 | TCP | Web interface |

### Default Paths

| Path | Description |
|------|-------------|
| `/etc/gomailserver/` | Configuration directory |
| `/var/lib/gomailserver/` | Data directory |
| `/var/log/gomailserver/` | Log directory |
| `/var/lib/gomailserver/messages/` | Large message storage (≥1MB) |
| `/var/lib/gomailserver/backups/` | Database backups |

### File Permissions

Recommended permissions:
- Configuration: `644` (rw-r--r--)
- Data directory: `750` (rwxr-x---)
- Log directory: `750` (rwxr-x---)
- Binary: `755` (rwxr-xr-x)
- Private key: `600` (rw-------)

### Security Checklist

Before going to production, ensure:

- [ ] Strong admin passwords (16+ characters)
- [ ] 2FA enabled for all admin accounts
- [ ] TLS certificates configured (not self-signed)
- [ ] Firewall rules in place
- [ ] DNS records properly configured (MX, SPF, DKIM, DMARC)
- [ ] DKIM keys generated and DNS configured
- [ ] ClamAV enabled and updated
- [ ] SpamAssassin enabled and trained
- [ ] Backup schedule configured
- [ ] Monitoring webhooks configured
- [ ] Rate limits configured appropriately
- [ ] Reputation scoring enabled
- [ ] Logging configured (info level)
- [ ] System tested with gomailtest
- [ ] Disaster recovery plan documented

---

**End of Administrator Guide**

For the latest version of this guide and additional documentation, visit:
https://github.com/btafoya/gomailserver
