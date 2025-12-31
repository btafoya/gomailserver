# gomailserver Deployment Guide

## Overview

gomailserver uses a **SQLite-first** architecture where all configuration and metadata are stored in SQLite, with only bootstrap settings (database path, ports, external service connections) in YAML configuration files.

## Architecture Philosophy

### Bootstrap Configuration (YAML)
- Database connection settings
- Server ports (SMTP, IMAP, API)
- Logger configuration
- External service connections (ClamAV, SpamAssassin)
- TLS/ACME settings

### Per-Domain Configuration (SQLite)
- All security policies (DKIM, SPF, DMARC, antivirus, antispam)
- Rate limiting rules
- Authentication security settings
- Greylisting configuration
- User management
- Mailbox configuration

### Benefits
- **Hot-Reload**: Domain settings changeable without server restart
- **Multi-Tenant**: Different security policies per domain
- **API-Driven**: All settings manageable via REST API
- **No Config Sprawl**: Single database file vs. dozens of config files
- **Scalable**: Easy to manage hundreds of domains programmatically

---

## Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu 22.04+ recommended)
- **Go**: 1.23.5 or later (for building from source)
- **Database**: SQLite 3.x (included)
- **Memory**: Minimum 512MB RAM (2GB+ recommended for production)
- **Disk**: 10GB+ for mail storage

### External Services (Optional but Recommended)
- **ClamAV**: Virus scanning (`clamd` daemon)
- **SpamAssassin**: Spam filtering (`spamd` daemon)

---

## Installation

### Option 1: Build from Source

```bash
# Clone repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Build
go build -o gomailserver ./cmd/gomailserver

# Or use the build script
./build.sh

# Install to system (requires sudo)
./build.sh install
```

### Option 2: Docker

```bash
# Build Docker image
docker build -t gomailserver:latest .

# Run with volume mounts
docker run -d \
  --name gomailserver \
  -p 25:25 \
  -p 143:143 \
  -p 465:465 \
  -p 587:587 \
  -p 993:993 \
  -p 8080:8080 \
  -v /path/to/data:/data \
  -v /path/to/config:/etc/gomailserver \
  gomailserver:latest
```

---

## Configuration

### Step 1: Bootstrap Configuration File

Create `gomailserver.yaml`:

```yaml
# Server identity
server:
  hostname: mail.example.com
  domain: example.com

# Database
database:
  path: ./data/mailserver.db
  wal_enabled: true

# Logging
logger:
  level: info
  format: json
  output_path: stdout

# SMTP ports
smtp:
  submission_port: 587    # Authenticated submission
  relay_port: 25          # Receiving from other servers
  smtps_port: 465         # SMTP over TLS
  max_message_size: 52428800  # 50MB

# IMAP ports
imap:
  port: 143               # Standard IMAP
  imaps_port: 993         # IMAP over TLS
  idle_timeout: 1800      # 30 minutes

# Admin API
api:
  port: 8080
  admin_token: "your-secure-random-token-here"  # IMPORTANT: Change this!

# External security services
security:
  clamav:
    socket_path: /var/run/clamav/clamd.ctl
    timeout: 60

  spamassassin:
    host: localhost
    port: 783
    timeout: 30

# TLS/ACME (Let's Encrypt)
tls:
  # Option 1: Manual certificates
  # cert_file: /path/to/cert.pem
  # key_file: /path/to/key.pem

  # Option 2: Automatic ACME (Let's Encrypt)
  acme:
    enabled: false  # Set to true for Let's Encrypt
    email: admin@example.com
    provider: cloudflare
    api_token: your-cloudflare-api-token
    cache_dir: ./data/acme
```

### Step 2: Initialize Database

On first run, the database will be automatically created and migrations applied:

```bash
./gomailserver run --config gomailserver.yaml
```

This will:
1. Create SQLite database at the configured path
2. Run all schema migrations
3. Create the `_default` domain template with sensible security defaults
4. Start all services (SMTP, IMAP, Admin API)

### Step 3: Configure Domains via Admin API

#### Create Admin API Token

Set a secure random token in `gomailserver.yaml`:

```bash
# Generate secure token
openssl rand -base64 32
```

#### Create Your First Domain

```bash
# Create domain from default template
curl -X POST http://localhost:8080/api/domains \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "use_default_template": true
  }'
```

#### Customize Domain Security Settings

```bash
# Update security configuration
curl -X PUT http://localhost:8080/api/domains/example.com/security \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "dkim": {
      "signing_enabled": true,
      "key_size": 4096
    },
    "spam": {
      "reject_score": 12.0
    },
    "auth": {
      "totp_enforced": true
    }
  }'
```

---

## Production Deployment

### System Service Setup

Create systemd service file `/etc/systemd/system/gomailserver.service`:

```ini
[Unit]
Description=gomailserver Mail Server
After=network.target

[Service]
Type=simple
User=mail
Group=mail
WorkingDirectory=/opt/gomailserver
ExecStart=/usr/local/bin/gomailserver run --config /etc/gomailserver/gomailserver.yaml
Restart=always
RestartSec=10

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/gomailserver/data

[Install]
WantedBy=multi-user.target
```

Start and enable the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable gomailserver
sudo systemctl start gomailserver
sudo systemctl status gomailserver
```

### ClamAV Setup

```bash
# Install ClamAV
sudo apt-get install clamav clamav-daemon

# Update virus definitions
sudo freshclam

# Start ClamAV daemon
sudo systemctl enable clamav-daemon
sudo systemctl start clamav-daemon

# Verify socket exists
ls -la /var/run/clamav/clamd.ctl
```

### SpamAssassin Setup

```bash
# Install SpamAssassin
sudo apt-get install spamassassin

# Enable and start spamd
sudo systemctl enable spamassassin
sudo systemctl start spamassassin

# Verify it's listening
netstat -tlnp | grep 783
```

### TLS with Let's Encrypt

#### Option 1: Automatic ACME (Recommended)

Enable ACME in `gomailserver.yaml`:

```yaml
tls:
  acme:
    enabled: true
    email: admin@example.com
    provider: cloudflare  # or other DNS provider
    api_token: your-dns-provider-api-token
    cache_dir: ./data/acme
```

Supported DNS providers:
- Cloudflare
- AWS Route53
- Google Cloud DNS
- DigitalOcean
- Others via DNS-01 challenge

#### Option 2: Manual Certificates

```yaml
tls:
  cert_file: /etc/letsencrypt/live/mail.example.com/fullchain.pem
  key_file: /etc/letsencrypt/live/mail.example.com/privkey.pem
```

Get certificates with certbot:

```bash
sudo certbot certonly --standalone -d mail.example.com
```

### Firewall Configuration

```bash
# Allow mail ports
sudo ufw allow 25/tcp    # SMTP relay
sudo ufw allow 143/tcp   # IMAP
sudo ufw allow 465/tcp   # SMTPS
sudo ufw allow 587/tcp   # SMTP submission
sudo ufw allow 993/tcp   # IMAPS

# Allow admin API (restrict to trusted IPs in production)
sudo ufw allow from 10.0.0.0/8 to any port 8080

sudo ufw enable
```

---

## DNS Configuration

### Required DNS Records

```dns
; MX record
example.com.           IN  MX  10 mail.example.com.

; A record for mail server
mail.example.com.      IN  A   203.0.113.10

; SPF record
example.com.           IN  TXT "v=spf1 mx -all"

; DMARC record
_dmarc.example.com.    IN  TXT "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"

; DKIM record (get public key from API)
default._domainkey.example.com. IN TXT "v=DKIM1; k=rsa; p=MIIBIj..."
```

### Get DKIM Public Key

```bash
curl -X GET http://localhost:8080/api/domains/example.com \
  -H "Authorization: Bearer your-admin-token" \
  | jq -r '.dkim_public_key'
```

---

## Backup and Recovery

### Backup Strategy

```bash
#!/bin/bash
# backup.sh - Daily backup script

BACKUP_DIR=/backup/gomailserver
DATE=$(date +%Y%m%d)

# Backup SQLite database
sqlite3 /opt/gomailserver/data/mailserver.db ".backup '$BACKUP_DIR/mailserver-$DATE.db'"

# Backup mail storage
tar -czf $BACKUP_DIR/mail-$DATE.tar.gz /opt/gomailserver/data/mail

# Backup ACME certificates
tar -czf $BACKUP_DIR/acme-$DATE.tar.gz /opt/gomailserver/data/acme

# Keep last 30 days
find $BACKUP_DIR -name "*.db" -mtime +30 -delete
find $BACKUP_DIR -name "*.tar.gz" -mtime +30 -delete
```

Schedule with cron:

```cron
0 2 * * * /opt/gomailserver/backup.sh
```

### Recovery

```bash
# Stop server
sudo systemctl stop gomailserver

# Restore database
cp /backup/gomailserver/mailserver-20250101.db /opt/gomailserver/data/mailserver.db

# Restore mail storage
tar -xzf /backup/gomailserver/mail-20250101.tar.gz -C /

# Start server
sudo systemctl start gomailserver
```

---

## Monitoring

### Health Checks

```bash
# API health check
curl http://localhost:8080/health

# Check service status
sudo systemctl status gomailserver

# View logs
sudo journalctl -u gomailserver -f
```

### Log Monitoring

Key log events to monitor:

```bash
# Authentication failures
journalctl -u gomailserver | grep "authentication failed"

# Rate limit hits
journalctl -u gomailserver | grep "rate_limit_exceeded"

# Virus detections
journalctl -u gomailserver | grep "virus_detected"

# Spam rejections
journalctl -u gomailserver | grep "spam"

# DKIM/SPF/DMARC failures
journalctl -u gomailserver | grep -E "dkim|spf|dmarc"
```

### Metrics to Track

- Message throughput (messages/hour)
- Authentication success/failure rates
- SPF/DKIM/DMARC pass rates
- Greylisting deferral rates
- Spam/virus detection rates
- API response times
- Database size growth

---

## Security Hardening

### File Permissions

```bash
# Create dedicated user
sudo useradd -r -s /bin/false mail

# Set ownership
sudo chown -R mail:mail /opt/gomailserver/data
sudo chmod 700 /opt/gomailserver/data

# Protect configuration
sudo chmod 600 /etc/gomailserver/gomailserver.yaml
```

### Admin API Security

1. **Use Strong Tokens**: Generate with `openssl rand -base64 32`
2. **Restrict Access**: Firewall API port to trusted IPs only
3. **Use HTTPS**: Reverse proxy with nginx/caddy for TLS
4. **Rotate Tokens**: Change admin token periodically

Example nginx reverse proxy:

```nginx
server {
    listen 443 ssl http2;
    server_name admin.example.com;

    ssl_certificate /etc/letsencrypt/live/admin.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/admin.example.com/privkey.pem;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### Rate Limiting

Adjust default rate limits via Admin API based on your needs:

```bash
curl -X PUT http://localhost:8080/api/domains/_default \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "ratelimit_smtp_per_ip": "{\"count\":50,\"window_minutes\":60}",
    "ratelimit_auth_per_ip": "{\"count\":5,\"window_minutes\":15}"
  }'
```

---

## Multi-Domain Hosting

### Add Additional Domains

```bash
# Domain 1
curl -X POST http://localhost:8080/api/domains \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{"name": "company1.com", "use_default_template": true}'

# Domain 2 with custom settings
curl -X POST http://localhost:8080/api/domains \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "company2.com",
    "use_default_template": false,
    "domain": {
      "status": "active",
      "max_users": 50,
      "default_quota": 2147483648,
      "dkim_signing_enabled": true,
      "spf_enabled": true,
      "auth_totp_enforced": true
    }
  }'
```

### Per-Domain Security Policies

Each domain can have unique security settings:

```bash
# Strict security for domain 1
curl -X PUT http://localhost:8080/api/domains/company1.com/security \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "spam": {"reject_score": 8.0},
    "auth": {"totp_enforced": true},
    "greylist": {"enabled": true}
  }'

# Relaxed security for domain 2
curl -X PUT http://localhost:8080/api/domains/company2.com/security \
  -H "Authorization: Bearer your-admin-token" \
  -H "Content-Type: application/json" \
  -d '{
    "spam": {"reject_score": 15.0},
    "auth": {"totp_enforced": false},
    "greylist": {"enabled": false}
  }'
```

---

## Performance Tuning

### SQLite Optimization

```bash
# Enable WAL mode (already default in config)
sqlite3 mailserver.db "PRAGMA journal_mode=WAL;"

# Optimize database periodically
sqlite3 mailserver.db "VACUUM;"
sqlite3 mailserver.db "ANALYZE;"
```

Add to weekly cron:

```cron
0 3 * * 0 sqlite3 /opt/gomailserver/data/mailserver.db "VACUUM; ANALYZE;"
```

### Resource Limits

Adjust in systemd service file:

```ini
[Service]
LimitNOFILE=65536       # Max open files
LimitNPROC=4096         # Max processes
MemoryMax=4G            # Max memory
CPUQuota=200%           # Max CPU (2 cores)
```

---

## Troubleshooting

### Common Issues

#### SMTP Connection Refused

```bash
# Check if SMTP server is running
ss -tlnp | grep :25

# Check firewall
sudo ufw status

# Check logs
journalctl -u gomailserver | grep smtp
```

#### Database Locked Errors

```bash
# Check WAL mode
sqlite3 mailserver.db "PRAGMA journal_mode;"

# Should return: wal

# If not, enable it
sqlite3 mailserver.db "PRAGMA journal_mode=WAL;"
```

#### ClamAV Not Scanning

```bash
# Check ClamAV daemon
sudo systemctl status clamav-daemon

# Check socket permissions
ls -la /var/run/clamav/clamd.ctl

# Test connection
clamdscan --version
```

#### Admin API 401 Unauthorized

```bash
# Verify token in config
grep admin_token /etc/gomailserver/gomailserver.yaml

# Check authorization header format
curl -H "Authorization: Bearer your-token" http://localhost:8080/api/domains
```

### Debug Mode

Enable debug logging:

```yaml
logger:
  level: debug
  format: json
  output_path: /var/log/gomailserver/debug.log
```

---

## Migration from Other Mail Servers

### From Postfix/Dovecot

1. Export user accounts:
```bash
# Via Admin API, create users for each domain
curl -X POST http://localhost:8080/api/users ...
```

2. Migrate mail storage (Maildir format compatible)
3. Update DNS records (MX, SPF, DKIM, DMARC)
4. Test with secondary MX before switching primary
5. Monitor delivery logs during transition

### Configuration Mapping

| Postfix/Dovecot | gomailserver (SQLite) |
|-----------------|----------------------|
| `/etc/postfix/main.cf` | `domains` table + `gomailserver.yaml` |
| `/etc/postfix/virtual` | `domains` table + `users` table |
| `/etc/dovecot/dovecot.conf` | `gomailserver.yaml` IMAP section |
| OpenDKIM config | `domains.dkim_*` columns |
| Amavis config | `domains.clamav_*`, `domains.spam_*` |
| Policy daemon rules | `domains.ratelimit_*`, `domains.greylist_*` |

---

## See Also

- [API.md](API.md) - Admin API documentation
- [INTEGRATION_GUIDE.md](INTEGRATION_GUIDE.md) - Security integration details
- [ISSUE003.md](ISSUE003.md) - Phase 2 security implementation
