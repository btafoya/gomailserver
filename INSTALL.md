# Installation Guide

Welcome to gomailserver! This guide walks you through installing and configuring your mail server.

---

## Table of Contents

- [Quick Start](#quick-start)
- [System Requirements](#system-requirements)
- [Installation Methods](#installation-methods)
  - [Method 1: Build from Source](#method-1-build-from-source)
  - [Method 2: Docker](#method-2-docker)
  - [Method 3: Docker Compose](#method-3-docker-compose)
- [Configuration](#configuration)
- [DNS Setup](#dns-setup)
- [TLS Certificates](#tls-certificates)
- [Running as a Service](#running-as-a-service)
- [Firewall Configuration](#firewall-configuration)
- [First Steps After Installation](#first-steps-after-installation)
- [Backup and Recovery](#backup-and-recovery)
- [Upgrading](#upgrading)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

**Already familiar with mail servers?** Here's the fast track:

```bash
# Clone and build
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver
make build

# Configure
cp .env.example .env
# Edit .env with your domain settings

# Run
./build/gomailserver run
```

Need more guidance? Read on!

---

## System Requirements

### Minimum Requirements

| Component | Requirement |
|-----------|-------------|
| **OS** | Linux (Ubuntu 20.04+, Debian 11+, RHEL 8+) |
| **CPU** | 1 core |
| **RAM** | 512 MB |
| **Disk** | 10 GB (depends on mail volume) |
| **Go** | 1.25.1+ (for building from source) |

### Recommended for Production

| Component | Requirement |
|-----------|-------------|
| **CPU** | 2+ cores |
| **RAM** | 2 GB+ |
| **Disk** | SSD with 50+ GB |
| **Network** | Static IP address |

### Network Ports

Your server needs these ports open:

| Port | Protocol | Purpose |
|------|----------|---------|
| 25 | TCP | SMTP (receiving mail) |
| 587 | TCP | SMTP Submission (sending mail) |
| 465 | TCP | SMTPS (secure sending) |
| 143 | TCP | IMAP (mail access) |
| 993 | TCP | IMAPS (secure mail access) |
| 8080 | TCP | Web UI (Admin + Webmail) |
| 8980 | TCP | Admin API |
| 8800 | TCP | CalDAV/CardDAV (optional) |

---

## Installation Methods

### Method 1: Build from Source

This is the recommended approach for production deployments.

#### Step 1: Install Prerequisites

**Ubuntu/Debian:**

```bash
sudo apt update
sudo apt install -y git make gcc
```

**RHEL/CentOS/Fedora:**

```bash
sudo dnf install -y git make gcc
```

**macOS:**

```bash
xcode-select --install
brew install go
```

#### Step 2: Install Go

Download and install Go 1.25.1 or later:

```bash
# Download Go
wget https://go.dev/dl/go1.25.1.linux-amd64.tar.gz

# Remove old version (if any) and install
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.25.1.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc for permanence)
export PATH=$PATH:/usr/local/go/bin

# Verify installation
go version
```

#### Step 3: Clone the Repository

```bash
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver
```

#### Step 4: Build

**Standard build:**

```bash
make build
```

**Static build (recommended for production):**

```bash
make build-static
```

The binary appears at `./build/gomailserver`.

#### Step 5: Verify the Build

```bash
./build/gomailserver version
```

---

### Method 2: Docker

The simplest way to get started.

#### Step 1: Build the Image

```bash
docker build -t gomailserver:latest .
```

#### Step 2: Create Data Directories

```bash
mkdir -p /var/lib/gomailserver
mkdir -p /etc/gomailserver
mkdir -p /var/log/gomailserver
```

#### Step 3: Create Configuration

```bash
cp .env.example /etc/gomailserver/.env
# Edit the configuration file
nano /etc/gomailserver/.env
```

#### Step 4: Run the Container

```bash
docker run -d \
  --name gomailserver \
  --restart unless-stopped \
  --env-file /etc/gomailserver/.env \
  -p 25:25 \
  -p 587:587 \
  -p 465:465 \
  -p 143:143 \
  -p 993:993 \
  -p 8080:8080 \
  -p 8980:8980 \
  -v /var/lib/gomailserver:/var/lib/gomailserver \
  -v /var/log/gomailserver:/var/log/gomailserver \
  gomailserver:latest
```

---

### Method 3: Docker Compose

Best for managing the full stack including optional services.

#### Step 1: Create docker-compose.yml

```bash
cp .doc_archive/docker-compose.yml ./docker-compose.yml
```

Or create a minimal version:

```yaml
version: '3.8'

services:
  gomailserver:
    build: .
    image: gomailserver:latest
    container_name: gomailserver
    restart: unless-stopped
    env_file:
      - .env
    ports:
      - "25:25"
      - "587:587"
      - "465:465"
      - "143:143"
      - "993:993"
      - "8080:8080"
      - "8980:8980"
    volumes:
      - gomailserver_data:/var/lib/gomailserver
      - gomailserver_logs:/var/log/gomailserver
    environment:
      - TZ=UTC
      - GOMEMLIMIT=512MiB
    healthcheck:
      test: ["CMD", "wget", "--spider", "-q", "http://localhost:8980/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3

volumes:
  gomailserver_data:
  gomailserver_logs:
```

#### Step 2: Start the Stack

```bash
# Build and start
docker-compose up -d

# View logs
docker-compose logs -f

# Stop
docker-compose down
```

---

## Configuration

gomailserver uses environment variables for all configuration. Settings can be defined in a `.env` file.

### Create Your Configuration

```bash
cp .env.example .env
nano .env
```

### Essential Settings

Update these settings in your `.env` file:

```bash
# ===========================================
# Server Identity (REQUIRED)
# ===========================================
SERVER_HOSTNAME=mail.example.com    # Your mail server's FQDN
SERVER_DOMAIN=example.com           # Your primary domain

# ===========================================
# Web UI (Admin Panel + Webmail)
# ===========================================
WEBUI_ENABLED=true
WEBUI_PORT=8080

# ===========================================
# Admin API
# ===========================================
API_PORT=8980
API_READ_TIMEOUT=15
API_WRITE_TIMEOUT=15
# API_JWT_SECRET=your-secret-key-here   # Set this for production!

# ===========================================
# SMTP Settings
# ===========================================
SMTP_HOSTNAME=mail.example.com
SMTP_SUBMISSION_PORT=587
SMTP_RELAY_PORT=25
SMTPS_PORT=465
SMTP_MAX_MESSAGE_SIZE=52428800      # 50MB

# ===========================================
# IMAP Settings
# ===========================================
IMAP_PORT=143
IMAPS_PORT=993
IMAP_IDLE_TIMEOUT=1800              # 30 minutes

# ===========================================
# Database
# ===========================================
DB_DRIVER=sqlite3
DB_SQLITE_PATH=./data/mailserver.db
DB_SQLITE_WAL_ENABLED=true

# ===========================================
# Logging
# ===========================================
LOG_LEVEL=info                      # debug, info, warn, error
LOG_FORMAT=json                     # json or text
```

### Complete Environment Variable Reference

#### Server Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_HOSTNAME` | localhost | Mail server FQDN |
| `SERVER_DOMAIN` | (empty) | Primary domain |

#### Web UI Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBUI_ENABLED` | true | Enable web interface |
| `WEBUI_PORT` | 8080 | Web UI port |

#### API Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `API_PORT` | 8980 | Admin API port |
| `API_READ_TIMEOUT` | 15 | Read timeout (seconds) |
| `API_WRITE_TIMEOUT` | 15 | Write timeout (seconds) |
| `API_JWT_SECRET` | (auto) | JWT signing secret |
| `API_CORS_ORIGINS` | (none) | Allowed CORS origins |

#### WebDAV Settings (CalDAV/CardDAV)

| Variable | Default | Description |
|----------|---------|-------------|
| `WEBDAV_ENABLED` | true | Enable CalDAV/CardDAV |
| `WEBDAV_PORT` | 8800 | WebDAV port |
| `WEBDAV_READ_TIMEOUT` | 30 | Read timeout (seconds) |
| `WEBDAV_WRITE_TIMEOUT` | 30 | Write timeout (seconds) |

#### SMTP Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `SMTP_HOSTNAME` | localhost | SMTP hostname |
| `SMTP_SUBMISSION_PORT` | 587 | Submission port |
| `SMTP_RELAY_PORT` | 25 | Relay port |
| `SMTPS_PORT` | 465 | Secure SMTP port |
| `SMTP_MAX_MESSAGE_SIZE` | 52428800 | Max message size (bytes) |

#### IMAP Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `IMAP_PORT` | 143 | IMAP port |
| `IMAPS_PORT` | 993 | Secure IMAP port |
| `IMAP_IDLE_TIMEOUT` | 1800 | IDLE timeout (seconds) |

#### Database Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_DRIVER` | sqlite3 | Database driver |
| `DB_SQLITE_PATH` | ./data/mailserver.db | SQLite file path |
| `DB_SQLITE_WAL_ENABLED` | true | Enable WAL mode |

**PostgreSQL (alternative):**

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_POSTGRES_HOST` | localhost | PostgreSQL host |
| `DB_POSTGRES_PORT` | 5432 | PostgreSQL port |
| `DB_POSTGRES_DATABASE` | gomailserver | Database name |
| `DB_POSTGRES_USER` | gomailserver | Username |
| `DB_POSTGRES_PASSWORD` | (empty) | Password |
| `DB_POSTGRES_SSL_MODE` | disable | SSL mode |

#### TLS Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `TLS_CERT_FILE` | (empty) | Path to certificate |
| `TLS_KEY_FILE` | (empty) | Path to private key |
| `ACME_ENABLED` | false | Enable Let's Encrypt |
| `ACME_EMAIL` | (empty) | ACME account email |
| `ACME_PROVIDER` | cloudflare | DNS provider |
| `CLOUDFLARE_API_TOKEN` | (empty) | Cloudflare API token |

#### Logging Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | info | Log level |
| `LOG_FORMAT` | json | Log format |

#### DNS Settings

| Variable | Default | Description |
|----------|---------|-------------|
| `DNS_RESOLVER` | 1.1.1.1:53 | DNS resolver address |
| `DNS_TIMEOUT` | 5 | DNS timeout (seconds) |
| `DNS_USE_TCP` | false | Use TCP for DNS |

#### Security Services (Optional)

| Variable | Default | Description |
|----------|---------|-------------|
| `CLAMAV_SOCKET_PATH` | /var/run/clamav/clamd.ctl | ClamAV socket |
| `CLAMAV_TIMEOUT` | 60 | ClamAV timeout (seconds) |
| `SPAMASSASSIN_HOST` | localhost | SpamAssassin host |
| `SPAMASSASSIN_PORT` | 783 | SpamAssassin port |
| `SPAMASSASSIN_TIMEOUT` | 30 | SpamAssassin timeout |

---

## DNS Setup

Proper DNS configuration is critical for email delivery.

### Required DNS Records

Replace `mail.example.com` with your actual hostname and `203.0.113.1` with your IP.

#### A Record (Mail Server Address)

```
mail.example.com.    IN A    203.0.113.1
```

#### MX Record (Mail Exchanger)

```
example.com.         IN MX   10 mail.example.com.
```

#### SPF Record (Sender Policy Framework)

```
example.com.         IN TXT  "v=spf1 mx a:mail.example.com -all"
```

#### DKIM Record

After starting gomailserver, retrieve your DKIM public key:

```bash
curl http://localhost:8980/api/v1/domains/example.com/dkim
```

Add the returned key as a DNS record:

```
default._domainkey.example.com.  IN TXT  "v=DKIM1; k=rsa; p=YOUR_PUBLIC_KEY_HERE"
```

#### DMARC Record

```
_dmarc.example.com.  IN TXT  "v=DMARC1; p=quarantine; rua=mailto:dmarc@example.com"
```

### Reverse DNS (PTR Record)

Contact your hosting provider to set up reverse DNS:

```
1.113.0.203.in-addr.arpa.  IN PTR  mail.example.com.
```

### Verify DNS Configuration

```bash
# Check MX record
dig MX example.com

# Check SPF
dig TXT example.com

# Check DKIM
dig TXT default._domainkey.example.com

# Check DMARC
dig TXT _dmarc.example.com
```

---

## TLS Certificates

gomailserver needs TLS certificates for secure connections.

### Option 1: Let's Encrypt (Recommended)

Automatic certificate management with ACME. Add to your `.env`:

```bash
ACME_ENABLED=true
ACME_EMAIL=admin@example.com
ACME_PROVIDER=cloudflare
CLOUDFLARE_API_TOKEN=your_api_token
```

Supported ACME providers:
- Cloudflare
- Route53
- Google Cloud DNS
- And many more...

### Option 2: Manual Certificates

If you have existing certificates, add to your `.env`:

```bash
TLS_CERT_FILE=/etc/ssl/certs/mail.example.com.crt
TLS_KEY_FILE=/etc/ssl/private/mail.example.com.key
```

### Option 3: Self-Signed (Development Only)

Leave `TLS_CERT_FILE` and `TLS_KEY_FILE` empty. gomailserver generates self-signed certificates automatically.

**Warning:** Self-signed certificates are **not recommended for production**.

---

## Running as a Service

### systemd Installation

For production Linux servers, use the included systemd installer:

```bash
# Build first
make build

# Install as systemd service
sudo ./scripts/install-systemd.sh --start
```

Options:

| Flag | Description |
|------|-------------|
| `--start` | Start immediately after install |
| `--enable` | Enable on boot (default) |
| `--no-enable` | Don't enable on boot |
| `--user USER` | Run as user (default: gomailserver) |
| `--prefix PATH` | Install prefix (default: /usr/local) |

### Managing the Service

```bash
# Start
sudo systemctl start gomailserver

# Stop
sudo systemctl stop gomailserver

# Restart
sudo systemctl restart gomailserver

# Check status
sudo systemctl status gomailserver

# View logs
sudo journalctl -u gomailserver -f

# Enable on boot
sudo systemctl enable gomailserver
```

### File Locations (systemd install)

| Path | Purpose |
|------|---------|
| `/usr/local/bin/gomailserver` | Binary |
| `/etc/gomailserver/.env` | Configuration |
| `/var/lib/gomailserver/` | Data and database |
| `/var/log/gomailserver/` | Logs |

---

## Firewall Configuration

### UFW (Ubuntu)

```bash
# Allow mail ports
sudo ufw allow 25/tcp    # SMTP
sudo ufw allow 587/tcp   # Submission
sudo ufw allow 465/tcp   # SMTPS
sudo ufw allow 143/tcp   # IMAP
sudo ufw allow 993/tcp   # IMAPS
sudo ufw allow 8080/tcp  # Web UI
sudo ufw allow 8980/tcp  # Admin API

# Apply
sudo ufw enable
```

### firewalld (RHEL/CentOS)

```bash
# Allow mail ports
sudo firewall-cmd --permanent --add-service=smtp
sudo firewall-cmd --permanent --add-service=smtps
sudo firewall-cmd --permanent --add-service=imap
sudo firewall-cmd --permanent --add-service=imaps
sudo firewall-cmd --permanent --add-port=587/tcp
sudo firewall-cmd --permanent --add-port=8080/tcp
sudo firewall-cmd --permanent --add-port=8980/tcp

# Reload
sudo firewall-cmd --reload
```

### iptables

```bash
# SMTP
sudo iptables -A INPUT -p tcp --dport 25 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 587 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 465 -j ACCEPT

# IMAP
sudo iptables -A INPUT -p tcp --dport 143 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 993 -j ACCEPT

# Web UI and Admin API
sudo iptables -A INPUT -p tcp --dport 8080 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8980 -j ACCEPT

# Save rules
sudo iptables-save > /etc/iptables/rules.v4
```

---

## First Steps After Installation

### 1. Access the Web Interface

Open your browser and navigate to:

```
http://your-server:8080
```

### 2. Access the Admin API

The admin API is available at:

```
http://your-server:8980
```

### 3. Create an Admin User

```bash
curl -X POST http://localhost:8980/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "secure-password", "role": "admin"}'
```

### 4. Add Your Domain

```bash
curl -X POST http://localhost:8980/api/v1/domains \
  -H "Content-Type: application/json" \
  -d '{"name": "example.com"}'
```

### 5. Create a Mailbox

```bash
curl -X POST http://localhost:8980/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "user-password",
    "display_name": "First User"
  }'
```

### 6. Test Email Delivery

Send a test email:

```bash
# Using swaks (install with: apt install swaks)
swaks --to user@example.com \
      --from test@example.com \
      --server localhost \
      --port 587 \
      --auth LOGIN \
      --auth-user user@example.com \
      --auth-password user-password \
      --tls
```

### 7. Configure an Email Client

Use these settings in your email client:

**Incoming (IMAP):**
- Server: mail.example.com
- Port: 993
- Security: SSL/TLS
- Username: user@example.com

**Outgoing (SMTP):**
- Server: mail.example.com
- Port: 587
- Security: STARTTLS
- Username: user@example.com

---

## Backup and Recovery

### Automated Backups

Use the included backup script:

```bash
# Basic backup
./scripts/backup.sh

# Custom backup location
./scripts/backup.sh --backup-dir /mnt/backups

# Keep 30 days of backups
./scripts/backup.sh --retention 30

# Dry run (see what would be backed up)
./scripts/backup.sh --dry-run
```

### Schedule Daily Backups

Add to crontab:

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /path/to/gomailserver/scripts/backup.sh --quiet
```

### What Gets Backed Up

| Component | Description |
|-----------|-------------|
| Database | SQLite database with all metadata |
| Messages | Email messages (large files) |
| DKIM Keys | Domain signing keys |
| Configuration | `.env` and other config files |

### Restore from Backup

```bash
# Stop the server
sudo systemctl stop gomailserver

# Restore database
sqlite3 /var/lib/gomailserver/mailserver.db ".restore '/path/to/backup/gomailserver_db_TIMESTAMP.db'"

# Restore messages
tar -xzf /path/to/backup/gomailserver_messages_TIMESTAMP.tar.gz -C /

# Restore DKIM keys
tar -xzf /path/to/backup/gomailserver_dkim_TIMESTAMP.tar.gz -C /

# Start the server
sudo systemctl start gomailserver
```

---

## Upgrading

### From Source

```bash
# Stop the server
sudo systemctl stop gomailserver

# Pull latest changes
cd /path/to/gomailserver
git pull origin main

# Rebuild
make clean
make build

# Reinstall
sudo cp build/gomailserver /usr/local/bin/

# Start the server
sudo systemctl start gomailserver
```

### Docker

```bash
# Pull latest
docker pull gomailserver:latest

# Or rebuild
docker-compose build --no-cache

# Restart
docker-compose down
docker-compose up -d
```

### Database Migrations

Migrations run automatically when the server starts. For manual control:

```bash
./scripts/migrate.sh
```

---

## Troubleshooting

### Server Won't Start

**Check configuration:**

```bash
# Verify .env file exists
cat .env

# Try running manually
./build/gomailserver run
```

Look for error messages about missing or invalid settings.

**Check permissions:**

```bash
# Ensure data directory is writable
ls -la /var/lib/gomailserver/
```

**Check ports:**

```bash
# See if ports are already in use
sudo ss -tlnp | grep -E '(25|587|465|143|993|8080|8980)'
```

### Can't Send Mail

**Check DNS:**

```bash
dig MX example.com
dig TXT example.com  # SPF
```

**Check connectivity:**

```bash
telnet mail.example.com 25
```

**Check logs:**

```bash
sudo journalctl -u gomailserver -f
# Or
tail -f /var/log/gomailserver/gomailserver.log
```

### Can't Receive Mail

**Verify MX record points to your server:**

```bash
dig MX example.com
```

**Verify port 25 is open:**

```bash
# From another machine
telnet mail.example.com 25
```

**Check firewall:**

```bash
sudo ufw status
# or
sudo firewall-cmd --list-all
```

### TLS Certificate Issues

**Check certificate validity:**

```bash
openssl s_client -connect mail.example.com:993 -showcerts
```

**Check certificate expiry:**

```bash
echo | openssl s_client -connect mail.example.com:993 2>/dev/null | openssl x509 -noout -dates
```

### Email Marked as Spam

Common causes and fixes:

1. **Missing SPF record** - Add SPF TXT record
2. **Missing DKIM signature** - Configure DKIM in admin panel
3. **No DMARC policy** - Add DMARC TXT record
4. **No reverse DNS** - Contact hosting provider
5. **IP on blacklist** - Check mxtoolbox.com and request removal

### Health Check

```bash
# API health endpoint
curl http://localhost:8980/api/v1/health

# Docker health
docker inspect --format='{{.State.Health.Status}}' gomailserver
```

---

## Getting Help

- **Documentation:** [GitHub Wiki](https://github.com/btafoya/gomailserver/wiki)
- **Issues:** [GitHub Issues](https://github.com/btafoya/gomailserver/issues)
- **API Reference:** Access at `http://localhost:8980/api/docs` when running

---

## Next Steps

Now that you have gomailserver running:

1. **Secure your setup** - Enable DKIM, SPF, and DMARC
2. **Set up monitoring** - Configure alerts for service health
3. **Plan backups** - Schedule automated backups
4. **Test deliverability** - Use mail-tester.com to check your score
5. **Review security** - Enable ClamAV and SpamAssassin if needed

Happy mailing!
