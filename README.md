![gomailserver](website-assets/gomailserver.png)

# gomailserver

[![CI](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml/badge.svg)](https://github.com/btafoya/gomailserver/actions/workflows/ci.yml)
[![Release](https://github.com/btafoya/gomailserver/actions/workflows/release.yml/badge.svg)](https://github.com/btafoya/gomailserver/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/btafoya/gomailserver)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/btafoya/gomailserver)](https://goreportcard.com/report/github.com/btafoya/gomailserver)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL_v3-blue.svg)](LICENSE.txt)

A modern, composable, all-in-one mail server written in Go designed to replace complex mail server stacks (Postfix, Dovecot, OpenDKIM, etc.) with a single daemon.

---

## Installation

### Quick Install (Linux)

#### Debian / Ubuntu (.deb)

```bash
# Download the latest release
curl -LO https://github.com/btafoya/gomailserver/releases/latest/download/gomailserver_*_amd64.deb

# Install
sudo dpkg -i gomailserver_*_amd64.deb

# Start the service
sudo systemctl enable --now gomailserver
```

#### RHEL / Fedora / CentOS (.rpm)

```bash
# Download the latest release
curl -LO https://github.com/btafoya/gomailserver/releases/latest/download/gomailserver_*_x86_64.rpm

# Install
sudo rpm -i gomailserver_*_x86_64.rpm

# Start the service
sudo systemctl enable --now gomailserver
```

#### Alpine Linux (.apk)

```bash
# Download the latest release
curl -LO https://github.com/btafoya/gomailserver/releases/latest/download/gomailserver_*_x86_64.apk

# Install
sudo apk add --allow-untrusted gomailserver_*_x86_64.apk
```

### Docker

```bash
# Pull the latest image
docker pull ghcr.io/btafoya/gomailserver:latest

# Run with default configuration
docker run -d \
  --name gomailserver \
  -p 25:25 \
  -p 587:587 \
  -p 465:465 \
  -p 143:143 \
  -p 993:993 \
  -p 8980:8980 \
  -v gomailserver-data:/var/lib/gomailserver \
  -v gomailserver-config:/etc/gomailserver \
  ghcr.io/btafoya/gomailserver:latest
```

### Binary Download

Download pre-built binaries for your platform from the [Releases](https://github.com/btafoya/gomailserver/releases) page:

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux | x86_64 | `gomailserver_VERSION_linux_amd64.tar.gz` |
| Linux | ARM64 | `gomailserver_VERSION_linux_arm64.tar.gz` |
| macOS | x86_64 | `gomailserver_VERSION_darwin_amd64.tar.gz` |
| macOS | ARM64 (Apple Silicon) | `gomailserver_VERSION_darwin_arm64.tar.gz` |
| Windows | x86_64 | `gomailserver_VERSION_windows_amd64.zip` |

```bash
# Example: Linux x86_64
curl -LO https://github.com/btafoya/gomailserver/releases/latest/download/gomailserver_*_linux_amd64.tar.gz
tar -xzf gomailserver_*_linux_amd64.tar.gz
sudo mv gomailserver /usr/local/bin/
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Build
make build

# Or build a static binary
make build-static

# Install system-wide
sudo make install
```

---

## Quick Start

1. **Edit the configuration file**

   ```bash
   sudo cp /usr/share/gomailserver/gomailserver.example.yaml /etc/gomailserver/gomailserver.yaml
   sudo nano /etc/gomailserver/gomailserver.yaml
   ```

2. **Configure your domain**

   ```yaml
   server:
     hostname: mail.example.com
     domain: example.com
   ```

3. **Start the service**

   ```bash
   sudo systemctl start gomailserver
   sudo systemctl status gomailserver
   ```

4. **Access the web interface**

   Open `http://your-server:8980` in your browser.

---

## Ports

| Port | Protocol | Description |
|------|----------|-------------|
| 25 | SMTP | Mail Transfer (MX) |
| 587 | SMTP | Mail Submission (STARTTLS) |
| 465 | SMTPS | Mail Submission (Implicit TLS) |
| 143 | IMAP | Mail Access |
| 993 | IMAPS | Mail Access (Implicit TLS) |
| 8980 | HTTP | Web Interface / API |

---

## Features

### Core Mail Services
- **SMTP** - Send and receive mail with full MTA/MX support
- **IMAP** - Access mail from any client
- **CalDAV** - Calendar synchronization with ACL support
- **CardDAV** - Contact synchronization

### Security
- **DKIM** - DomainKeys Identified Mail signing
- **SPF** - Sender Policy Framework validation
- **DMARC** - Domain-based Message Authentication
- **TLS** - Automatic certificate management via Let's Encrypt
- **AI-Powered Phishing Detection** - Real-time threat analysis

### Administration
- **Web Interface** - Modern, responsive admin panel
- **RESTful API** - Full programmatic access
- **Systemd Integration** - Service management
- **SQLite/PostgreSQL** - Flexible database options

---

## Configuration

See [`gomailserver.example.yaml`](gomailserver.example.yaml) for a complete configuration reference.

Key configuration sections:

```yaml
server:
  hostname: mail.example.com
  domain: example.com

database:
  driver: sqlite3  # or "postgres"
  sqlite:
    path: /var/lib/gomailserver/mailserver.db

smtp:
  submission_port: 587
  relay_port: 25
  smtps_port: 465

imap:
  port: 143
  imaps_port: 993

tls:
  acme:
    enabled: true
    email: admin@example.com
```

---

## Documentation

- [Installation Guide](INSTALL.md) - Detailed installation instructions
- [Release Workflow](WORKFLOW.md) - CI/CD and release process
- [API Documentation](api-docs/) - REST API reference
- [TODO](TODO.md) - Roadmap and planned features

---

## Contributing

Contributions are welcome! Please read our contributing guidelines before submitting pull requests.

```bash
# Run tests
make test

# Run linter
make lint

# Build and test locally
make build
./build/gomailserver version
```

---

## License

This project is licensed under the [GNU Affero General Public License v3.0](LICENSE.txt).

---

## Support

- **Issues**: [GitHub Issues](https://github.com/btafoya/gomailserver/issues)
- **Discussions**: [GitHub Discussions](https://github.com/btafoya/gomailserver/discussions)
