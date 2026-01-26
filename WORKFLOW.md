# Release Workflow Documentation

This document describes the release process and CI/CD pipeline for gomailserver.

## Table of Contents

- [Release Process](#release-process)
- [CI/CD Pipeline](#cicd-pipeline)
- [Local Testing](#local-testing)
- [Repository Setup](#repository-setup)
- [Troubleshooting](#troubleshooting)

---

## Release Process

### Version Numbering

gomailserver follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes to configuration, API, or behavior
- **MINOR**: New features, backwards-compatible
- **PATCH**: Bug fixes, backwards-compatible

Examples:
- `v1.0.0` - First stable release
- `v1.1.0` - New feature added
- `v1.1.1` - Bug fix
- `v2.0.0` - Breaking configuration change

### Creating a Release

1. **Ensure all changes are merged to main**
   ```bash
   git checkout main
   git pull origin main
   ```

2. **Verify tests pass**
   ```bash
   make test
   make lint
   ```

3. **Check GoReleaser configuration**
   ```bash
   make release-check
   ```

4. **Create and push a version tag**
   ```bash
   # For a release
   git tag v1.0.0
   git push origin v1.0.0

   # For a pre-release
   git tag v1.0.0-rc.1
   git push origin v1.0.0-rc.1
   ```

5. **Monitor the release workflow**
   - Go to GitHub Actions
   - Watch the "Release" workflow
   - Verify all jobs complete successfully

6. **Verify the release**
   - Check GitHub Releases page
   - Test installation (see [Verification](#verification))

### Pre-release Checklist

Before creating a release tag:

- [ ] All tests pass (`make test`)
- [ ] Linter passes (`make lint`)
- [ ] GoReleaser config is valid (`make release-check`)
- [ ] Snapshot build works (`make release-snapshot`)
- [ ] CHANGELOG is updated (or rely on auto-generation)
- [ ] Documentation is current
- [ ] Breaking changes are documented

### Release Artifacts

Each release produces:

| Artifact | Description |
|----------|-------------|
| `gomailserver_VERSION_linux_amd64.tar.gz` | Linux x86_64 binary |
| `gomailserver_VERSION_linux_arm64.tar.gz` | Linux ARM64 binary |
| `gomailserver_VERSION_darwin_amd64.tar.gz` | macOS x86_64 binary |
| `gomailserver_VERSION_darwin_arm64.tar.gz` | macOS ARM64 (Apple Silicon) |
| `gomailserver_VERSION_windows_amd64.zip` | Windows x86_64 binary |
| `gomailserver_VERSION_amd64.deb` | Debian/Ubuntu package |
| `gomailserver_VERSION_arm64.deb` | Debian/Ubuntu ARM64 package |
| `gomailserver_VERSION_x86_64.rpm` | RHEL/Fedora package |
| `gomailserver_VERSION_aarch64.rpm` | RHEL/Fedora ARM64 package |
| `gomailserver_VERSION_x86_64.apk` | Alpine package |
| `gomailserver_VERSION_aarch64.apk` | Alpine ARM64 package |
| `checksums.txt` | SHA256 checksums for all artifacts |

### Container Images

Multi-arch images are pushed to GitHub Container Registry:

```bash
# Specific version
docker pull ghcr.io/btafoya/gomailserver:1.0.0

# Latest release
docker pull ghcr.io/btafoya/gomailserver:latest
```

---

## CI/CD Pipeline

### Workflow Overview

```
push tag v* ─────► Release Workflow
                        │
                        ├── Build binaries (linux/darwin/windows)
                        ├── Create archives (tar.gz/zip)
                        ├── Build Linux packages (deb/rpm/apk)
                        ├── Build Docker images (amd64/arm64)
                        ├── Push to GHCR
                        └── Create GitHub Release
```

### Triggers

| Event | Workflow | Description |
|-------|----------|-------------|
| Push tag `v*` | release.yml | Full release with all artifacts |
| Push to main | ci.yml | Tests and linting |
| Pull request | ci.yml | Tests and linting |

### Build Matrix

| OS | Architectures | Notes |
|----|--------------|-------|
| Linux | amd64, arm64 | Full support |
| macOS | amd64, arm64 | Full support |
| Windows | amd64 | Binary only |

### Required Secrets

Configure these in repository Settings > Secrets and variables > Actions:

| Secret | Purpose | Required |
|--------|---------|----------|
| `GITHUB_TOKEN` | Release creation, GHCR push | Auto-provided |

### Workflow Jobs

#### release.yml

1. **Build and Release**
   - Checkout code with full history
   - Setup Go and Docker Buildx
   - Login to GitHub Container Registry
   - Run GoReleaser (builds, packages, releases)
   - Upload artifacts

2. **Verify Release**
   - Download artifacts
   - Verify checksums
   - List all artifacts

---

## Local Testing

### Prerequisites

Install GoReleaser:

```bash
# macOS
brew install goreleaser

# Linux (snap)
sudo snap install --classic goreleaser

# Go install
go install github.com/goreleaser/goreleaser/v2@latest
```

### Check Configuration

Validate the GoReleaser configuration:

```bash
make release-check
# or
goreleaser check
```

### Build Snapshot

Build a snapshot release without publishing:

```bash
make release-snapshot
# or
goreleaser release --snapshot --clean
```

This creates all artifacts in `./dist/`:

```bash
ls -la dist/
# gomailserver_linux_amd64.tar.gz
# gomailserver_darwin_arm64.tar.gz
# gomailserver_1.0.0-next_amd64.deb
# etc.
```

### Test Packages Locally

#### DEB Package

```bash
# Install
sudo dpkg -i dist/gomailserver_*_amd64.deb

# Verify
gomailserver version

# Check service file
cat /lib/systemd/system/gomailserver.service

# Uninstall
sudo dpkg -r gomailserver
```

#### RPM Package

```bash
# Install
sudo rpm -i dist/gomailserver_*_x86_64.rpm

# Verify
gomailserver version

# Uninstall
sudo rpm -e gomailserver
```

#### Docker Image

```bash
# Build snapshot (includes Docker images)
make release-snapshot

# Load local image
docker load < dist/gomailserver_linux_amd64.tar

# Run
docker run --rm ghcr.io/btafoya/gomailserver:*-next version
```

### Test Release Workflow

Create a release candidate tag:

```bash
git tag v0.1.0-rc.1
git push origin v0.1.0-rc.1
```

This triggers the full release workflow but marks it as a pre-release.

---

## Repository Setup

### Initial Setup Checklist

1. **Verify workflows**
   - Check `.github/workflows/release.yml` exists
   - Run `make release-check`

2. **Test with snapshot**
   - Run `make release-snapshot`
   - Verify all artifacts build

3. **Create initial release**
   - Tag `v0.1.0` or `v1.0.0`
   - Monitor workflow
   - Verify artifacts

### File Structure

```
.
├── .github/
│   └── workflows/
│       ├── release.yml      # GoReleaser release workflow
│       ├── ci.yml           # Tests and linting
│       └── build-deb.yml    # Distro-specific DEBs (kept)
├── .goreleaser.yaml         # GoReleaser configuration
├── Dockerfile               # Development Docker build
├── Dockerfile.goreleaser    # Release Docker build
├── debian/
│   ├── postinst             # Post-install script
│   ├── prerm                # Pre-remove script
│   └── postrm               # Post-remove script
├── scripts/
│   └── gomailserver.service # Systemd service file
└── Makefile                 # Build targets
```

---

## Troubleshooting

### Common Issues

#### GoReleaser fails to build

```
Error: build for linux_amd64 failed
```

**Solution**: Check Go version compatibility and module dependencies:

```bash
go mod tidy
go build ./cmd/gomailserver
```

#### Docker push fails

```
Error: unauthorized: authentication required
```

**Solution**: Verify `GITHUB_TOKEN` has `packages:write` permission.

#### nfpm packaging fails

```
Error: open debian/postinst: no such file or directory
```

**Solution**: Ensure all files referenced in `.goreleaser.yaml` exist:

```bash
ls -la debian/
ls -la scripts/gomailserver.service
```

#### Version not injected

```
gomailserver version
# Shows: dev
```

**Solution**: Build with ldflags:

```bash
go build -ldflags="-X main.version=1.0.0" ./cmd/gomailserver
```

### Debug Commands

```bash
# Verbose GoReleaser output
goreleaser release --snapshot --clean --verbose

# Check specific package
rpm -qpl dist/gomailserver_*.rpm
dpkg-deb -c dist/gomailserver_*.deb

# Inspect Docker image
docker inspect ghcr.io/btafoya/gomailserver:latest
```

### Getting Help

- [GoReleaser Documentation](https://goreleaser.com/)
- [nfpm Documentation](https://nfpm.goreleaser.com/)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Project Issues](https://github.com/btafoya/gomailserver/issues)

---

## Quick Reference

### Common Commands

```bash
# Check configuration
make release-check

# Build snapshot (local testing)
make release-snapshot

# Create release
git tag v1.0.0 && git push origin v1.0.0

# View release artifacts
ls -la dist/

# Install local DEB
sudo dpkg -i dist/gomailserver_*_amd64.deb

# Install local RPM
sudo rpm -i dist/gomailserver_*_x86_64.rpm

# Test Docker image
docker run --rm ghcr.io/btafoya/gomailserver:latest version
```

### Useful Links

- **GoReleaser**: https://goreleaser.com/
- **nfpm**: https://nfpm.goreleaser.com/
- **Container Registry**: https://ghcr.io/btafoya/gomailserver
