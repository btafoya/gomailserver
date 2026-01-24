# gomailserver Automated Backup System

## Overview

The gomailserver automated backup system provides comprehensive backup and recovery capabilities for all gomailserver data, including:

- **SQLite Database**: All metadata, user accounts, domains, and small messages (<1MB)
- **Filesystem Messages**: Large message attachments and content (≥1MB)
- **DKIM Keys**: Private cryptographic keys for email signing
- **Configuration**: All system configuration files
- **Backup Verification**: Automated integrity checking
- **Retention Management**: Automatic cleanup of old backups

## Quick Start

### Basic Backup
```bash
# Run manual backup
sudo /opt/gomailserver/scripts/backup.sh

# Run with custom options
sudo /opt/gomailserver/scripts/backup.sh --retention 30 --backup-dir /mnt/backups
```

### Automated Backups
```bash
# Install cron configuration
sudo cp /opt/gomailserver/scripts/cron-backup /etc/cron.d/gomailserver-backup
sudo systemctl restart cron

# Verify backups are running
sudo /opt/gomailserver/scripts/verify-backup.sh
```

## Backup Contents

### Database Backup
- **File**: `gomailserver_db_YYYYMMDD_HHMMSS.db`
- **Contents**: SQLite database with all metadata and small messages
- **Integrity**: Verified with `PRAGMA integrity_check`
- **Compression**: None (SQLite WAL mode for consistency)

### Messages Backup
- **File**: `gomailserver_messages_YYYYMMDD_HHMMSS.tar.gz`
- **Contents**: Large message files from `/var/lib/gomailserver/messages/`
- **Compression**: gzip
- **Permissions**: Preserved during backup/restore

### DKIM Keys Backup
- **File**: `gomailserver_dkim_YYYYMMDD_HHMMSS.tar.gz`
- **Contents**: DKIM private keys from `/var/lib/gomailserver/dkim/`
- **Security**: Restricted permissions (600)
- **Encryption**: Not encrypted (store on encrypted filesystem)

### Configuration Backup
- **File**: `gomailserver_config_YYYYMMDD_HHMMSS.tar.gz`
- **Contents**: Configuration files from `/etc/gomailserver/`
- **Includes**: YAML configs, certificates, keys

### Manifest File
- **File**: `gomailserver_manifest_YYYYMMDD_HHMMSS.txt`
- **Contents**: Backup metadata, file lists, verification commands

## Configuration

### Environment Variables
```bash
# Backup directory
export BACKUP_DIR="/var/lib/gomailserver/backups"

# Source directories
export DATABASE_PATH="/var/lib/gomailserver/mailserver.db"
export MESSAGE_DIR="/var/lib/gomailserver/messages"
export DKIM_DIR="/var/lib/gomailserver/dkim"
export CONFIG_DIR="/etc/gomailserver"

# Retention policy
export RETENTION_DAYS=7

# Options
export COMPRESSION="gzip"
export QUIET="false"
export DRY_RUN="false"
```

### Command Line Options
```bash
backup.sh [OPTIONS]

Options:
  -d, --backup-dir DIR       Backup directory (default: /var/lib/gomailserver/backups)
  --database PATH            Database file path
  --messages DIR             Messages directory
  --dkim DIR                 DKIM keys directory
  --config DIR               Config directory
  --retention DAYS           Days to keep backups (default: 7)
  --compression TYPE         Compression: gzip, bzip2, xz
  -q, --quiet                Suppress console output
  --dry-run                  Show what would be done
  -h, --help                 Show help
  -v, --version              Show version
```

## Scheduling

### Daily Backups
```bash
# /etc/cron.d/gomailserver-backup
0 2 * * * gomailserver /opt/gomailserver/scripts/backup.sh --quiet
```

### Weekly Full Backups
```bash
# /etc/cron.d/gomailserver-backup
0 3 * * 0 gomailserver /opt/gomailserver/scripts/backup.sh --quiet --retention 30
```

### Monthly Archives
```bash
# /etc/cron.d/gomailserver-backup
0 4 1 * * gomailserver /opt/gomailserver/scripts/backup.sh --quiet --retention 365 --backup-dir /var/lib/gomailserver/monthly-backups
```

## Verification

### Manual Verification
```bash
# Verify latest backup
sudo /opt/gomailserver/scripts/verify-backup.sh

# Verify specific backup directory
sudo /opt/gomailserver/scripts/verify-backup.sh /mnt/backups

# Check database integrity
sqlite3 /var/lib/gomailserver/backups/gomailserver_db_*.db 'PRAGMA integrity_check;'

# List message archive contents
tar -tzf /var/lib/gomailserver/backups/gomailserver_messages_*.tar.gz
```

### Automated Verification
```bash
# Add to cron for daily verification
# /etc/cron.d/gomailserver-backup
30 2 * * * gomailserver /opt/gomailserver/scripts/verify-backup.sh
```

## Recovery

### Database Recovery
```bash
# Stop gomailserver
sudo systemctl stop gomailserver

# Backup current database
sudo cp /var/lib/gomailserver/mailserver.db /var/lib/gomailserver/mailserver.db.broken

# Restore from backup
sudo cp /var/lib/gomailserver/backups/gomailserver_db_20240123_020000.db /var/lib/gomailserver/mailserver.db

# Verify integrity
sudo sqlite3 /var/lib/gomailserver/mailserver.db 'PRAGMA integrity_check;'

# Start gomailserver
sudo systemctl start gomailserver
```

### Messages Recovery
```bash
# Extract messages to temporary location
sudo mkdir -p /tmp/messages-recovery
sudo tar -xzf /var/lib/gomailserver/backups/gomailserver_messages_20240123_020000.tar.gz -C /tmp/messages-recovery

# Verify contents
sudo ls -la /tmp/messages-recovery/var/lib/gomailserver/messages/

# Replace current messages (if needed)
sudo rsync -av /tmp/messages-recovery/var/lib/gomailserver/messages/ /var/lib/gomailserver/messages/

# Clean up
sudo rm -rf /tmp/messages-recovery
```

### DKIM Keys Recovery
```bash
# Extract DKIM keys
sudo mkdir -p /tmp/dkim-recovery
sudo tar -xzf /var/lib/gomailserver/backups/gomailserver_dkim_20240123_020000.tar.gz -C /tmp/dkim-recovery

# Restore keys
sudo rsync -av /tmp/dkim-recovery/var/lib/gomailserver/dkim/ /var/lib/gomailserver/dkim/

# Set correct permissions
sudo chown -R gomailserver:gomailserver /var/lib/gomailserver/dkim/
sudo chmod 700 /var/lib/gomailserver/dkim/
sudo chmod 600 /var/lib/gomailserver/dkim/*/*.key

# Clean up
sudo rm -rf /tmp/dkim-recovery

# Restart gomailserver to reload keys
sudo systemctl restart gomailserver
```

## Monitoring

### Backup Logs
```bash
# View backup logs
sudo tail -f /var/log/gomailserver/backup.log

# Check for errors
sudo grep "ERROR" /var/log/gomailserver/backup.log

# View backup summary
sudo grep "Backup completed" /var/log/gomailserver/backup.log
```

### Disk Space Monitoring
```bash
# Check backup directory usage
sudo du -sh /var/lib/gomailserver/backups/

# Monitor available space
sudo df -h /var/lib/gomailserver/backups/

# Alert if low space
BACKUP_SPACE=$(df /var/lib/gomailserver/backups/ | tail -1 | awk '{print $5}' | sed 's/%//')
if [ "$BACKUP_SPACE" -gt 90 ]; then
    echo "WARNING: Backup disk usage is ${BACKUP_SPACE}%"
fi
```

### Backup Success Monitoring
```bash
# Check last backup success
LAST_BACKUP=$(find /var/lib/gomailserver/backups/ -name "gomailserver_db_*.db" -printf '%T@ %p\n' | sort -n | tail -1 | cut -d' ' -f2-)
LAST_BACKUP_TIME=$(stat -c %Y "$LAST_BACKUP")
CURRENT_TIME=$(date +%s)
AGE_HOURS=$(( (CURRENT_TIME - LAST_BACKUP_TIME) / 3600 ))

if [ "$AGE_HOURS" -gt 25 ]; then
    echo "WARNING: Last backup is ${AGE_HOURS} hours old"
fi
```

## Security Considerations

### Permissions
- Backup files: `640` (owner: gomailserver, group: gomailserver)
- DKIM keys: `600` (owner: gomailserver, group: gomailserver)
- Backup directory: `750` (owner: gomailserver, group: gomailserver)

### Encryption
- Store backups on encrypted filesystems
- Use encrypted backup destinations for off-site storage
- Consider encrypting sensitive configuration files

### Access Control
- Limit access to backup files
- Use separate backup storage for sensitive data
- Implement backup file integrity verification

## Troubleshooting

### Common Issues

#### Backup Fails with Permission Denied
```bash
# Check directory permissions
ls -la /var/lib/gomailserver/backups/

# Fix permissions
sudo chown gomailserver:gomailserver /var/lib/gomailserver/backups/
sudo chmod 750 /var/lib/gomailserver/backups/
```

#### Database Backup Shows Corruption
```bash
# Check database integrity before backup
sudo sqlite3 /var/lib/gomailserver/mailserver.db 'PRAGMA integrity_check;'

# If corrupted, restore from previous backup
sudo cp /var/lib/gomailserver/backups/gomailserver_db_previous.db /var/lib/gomailserver/mailserver.db
```

#### Insufficient Disk Space
```bash
# Check available space
df -h /var/lib/gomailserver/backups/

# Clean old backups manually
sudo find /var/lib/gomailserver/backups/ -name "gomailserver_*" -mtime +7 -delete

# Or change retention policy
sudo /opt/gomailserver/scripts/backup.sh --retention 3
```

#### Backup Verification Fails
```bash
# Check specific backup
sudo /opt/gomailserver/scripts/verify-backup.sh /path/to/backup/dir

# Manual verification
sqlite3 /var/lib/gomailserver/backups/gomailserver_db_*.db 'PRAGMA integrity_check;'
tar -tzf /var/lib/gomailserver/backups/gomailserver_messages_*.tar.gz | head -10
```

## Performance Tuning

### Backup Speed Optimization
```bash
# Use faster compression (if CPU is bottleneck)
sudo /opt/gomailserver/scripts/backup.sh --compression gzip

# Disable compression for local backups
sudo /opt/gomailserver/scripts/backup.sh --compression none

# Parallel compression for large message archives
# (automatically handled by tar with pigz if available)
```

### Storage Optimization
```bash
# Use compression for large message backups
sudo /opt/gomailserver/scripts/backup.sh --compression xz  # Best compression
sudo /opt/gomailserver/scripts/backup.sh --compression bzip2  # Good compression
sudo /opt/gomailserver/scripts/backup.sh --compression gzip  # Fast compression
```

### Retention Policy Tuning
```bash
# Short retention for development
sudo /opt/gomailserver/scripts/backup.sh --retention 3

# Standard retention for production
sudo /opt/gomailserver/scripts/backup.sh --retention 7

# Long retention for compliance
sudo /opt/gomailserver/scripts/backup.sh --retention 365
```

## Integration with Monitoring Systems

### Prometheus Metrics
```bash
# Backup success/failure metrics
backup_last_success_timestamp{job="gomailserver"} <timestamp>
backup_duration_seconds{job="gomailserver"} <duration>
backup_size_bytes{job="gomailserver"} <size>

# Disk usage alerts
backup_disk_usage_percent{job="gomailserver"} <percentage>
backup_disk_available_bytes{job="gomailserver"} <bytes>
```

### Webhook Notifications
```bash
# Configure webhook for backup notifications
curl -X POST http://localhost:8980/api/v1/webhooks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Backup Monitor",
    "url": "https://monitoring.example.com/webhook",
    "events": ["system.backup_completed", "system.backup_failed"],
    "active": true
  }'
```

## Best Practices

1. **Test Backups Regularly**: Always verify backup integrity
2. **Monitor Backup Success**: Set up alerts for backup failures
3. **Use Multiple Locations**: Keep backups in multiple locations
4. **Encrypt Sensitive Data**: Use encrypted storage for DKIM keys
5. **Document Recovery Procedures**: Keep recovery steps updated
6. **Automate Everything**: Use cron jobs and monitoring
7. **Monitor Disk Space**: Prevent backup failures due to full disks
8. **Test Recovery**: Regularly test full recovery procedures

## Support

For issues with the backup system:
1. Check `/var/log/gomailserver/backup.log`
2. Run verification: `sudo /opt/gomailserver/scripts/verify-backup.sh`
3. Review backup manifest for details
4. Check GitHub issues: https://github.com/btafoya/gomailserver/issues