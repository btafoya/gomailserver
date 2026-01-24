# Installing gomailserver via APT

You can install gomailserver on Debian/Ubuntu systems using our APT repository hosted on GitHub Pages.

## Supported Distributions

- Ubuntu: focal (20.04), jammy (22.04), noble (24.04)
- Debian: bullseye (11), bookworm (12)

## Supported Architectures

- amd64 (x86_64)
- arm64 (aarch64)

## Installation Steps

### 1. Add the GPG Key

First, download and add the repository signing key:

```bash
curl -fsSL https://btafoya.github.io/gomailserver/repo/public.key | \
  sudo gpg --dearmor -o /usr/share/keyrings/gomailserver-archive-keyring.gpg
```

### 2. Add the Repository

Add the repository to your APT sources. Replace `DISTRO` with your distribution codename (e.g., `jammy`, `bookworm`):

```bash
echo "deb [signed-by=/usr/share/keyrings/gomailserver-archive-keyring.gpg] https://btafoya.github.io/gomailserver/repo DISTRO main" | \
  sudo tee /etc/apt/sources.list.d/gomailserver.list
```

**For Ubuntu 22.04 (jammy):**
```bash
echo "deb [signed-by=/usr/share/keyrings/gomailserver-archive-keyring.gpg] https://btafoya.github.io/gomailserver/repo jammy main" | \
  sudo tee /etc/apt/sources.list.d/gomailserver.list
```

**For Debian 12 (bookworm):**
```bash
echo "deb [signed-by=/usr/share/keyrings/gomailserver-archive-keyring.gpg] https://btafoya.github.io/gomailserver/repo bookworm main" | \
  sudo tee /etc/apt/sources.list.d/gomailserver.list
```

### 3. Update Package List

```bash
sudo apt update
```

### 4. Install gomailserver

```bash
sudo apt install gomailserver
```

## Configuration

After installation:

1. Edit the configuration file:
   ```bash
   sudo nano /etc/gomailserver/gomailserver.yaml
   ```

2. Enable and start the service:
   ```bash
   sudo systemctl enable gomailserver
   sudo systemctl start gomailserver
   ```

3. Check the service status:
   ```bash
   sudo systemctl status gomailserver
   ```

## Upgrading

To upgrade to the latest version:

```bash
sudo apt update
sudo apt upgrade gomailserver
```

## Uninstalling

To remove gomailserver but keep configuration:

```bash
sudo apt remove gomailserver
```

To remove everything including configuration and data:

```bash
sudo apt purge gomailserver
```

## File Locations

- Binary: `/usr/bin/gomailserver`
- Configuration: `/etc/gomailserver/gomailserver.yaml`
- Web assets: `/usr/share/gomailserver/`
- Data directory: `/var/lib/gomailserver/`
- Log directory: `/var/log/gomailserver/`
- Systemd service: `/lib/systemd/system/gomailserver.service`

## Troubleshooting

### Repository not found

Make sure you've replaced `DISTRO` with your actual distribution codename. You can find it with:

```bash
lsb_release -cs
```

### GPG key errors

If you get GPG errors, ensure the key was properly downloaded and added:

```bash
ls -l /usr/share/keyrings/gomailserver-archive-keyring.gpg
```

### Service won't start

Check the logs:

```bash
sudo journalctl -u gomailserver -n 50 --no-pager
```

Verify your configuration:

```bash
sudo gomailserver -config /etc/gomailserver/gomailserver.yaml -validate
```
