# gomailserver Production Docker Setup

## Overview

This Docker setup provides a production-ready gomailserver deployment with security hardening, multi-stage builds, and optimized performance.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/btafoya/gomailserver.git
cd gomailserver

# Start the services
make docker-run

# Access admin interface
open http://localhost:8980/admin/
```

### Manual Docker Commands

```bash
# Build the image
make docker-build

# Run the container
docker run -d \
  --name gomailserver \
  -p 25:25 -p 587:587 -p 465:465 \
  -p 143:143 -p 993:993 -p 8980:8980 \
  -v gomailserver_data:/var/lib/gomailserver \
  -v gomailserver_logs:/var/log/gomailserver \
  -v gomailserver_config:/etc/gomailserver \
  --security-opt no-new-privileges:true \
  --cap-drop ALL \
  --cap-add NET_BIND_SERVICE \
  --read-only \
  --tmpfs /tmp:noexec,nosuid,size=100m \
  gomailserver:latest
```

## Docker Image Features

### Security Hardening
- **Non-root user**: Runs as `gomailserver` user (UID 1000)
- **Minimal base image**: Alpine Linux for reduced attack surface
- **No new privileges**: Security option prevents privilege escalation
- **Capability dropping**: Drops all capabilities except NET_BIND_SERVICE
- **Read-only filesystem**: Prevents unauthorized modifications
- **tmpfs for temp files**: Secure temporary file handling

### Multi-Stage Build
- **Builder stage**: Go compilation with security flags
- **Runtime stage**: Minimal Alpine with only necessary runtime dependencies
- **Static binary**: No dynamic linking dependencies
- **Small image size**: ~50MB compressed

### Production Optimizations
- **Health checks**: HTTP-based health monitoring
- **Resource limits**: Memory and CPU limits configured
- **Logging**: JSON structured logging with size limits
- **Volumes**: Persistent data storage for database, logs, and config

## Configuration

### Environment Variables

```bash
# Timezone
TZ=UTC

# Go runtime limits
GOMEMLIMIT=512MiB
GOMAXPROCS=2
```

### Volume Mounts

| Host Path | Container Path | Purpose |
|-----------|----------------|---------|
| `gomailserver_data` | `/var/lib/gomailserver` | SQLite database and message storage |
| `gomailserver_logs` | `/var/log/gomailserver` | Application logs |
| `gomailserver_config` | `/etc/gomailserver` | Configuration files |

### Ports

| Port | Protocol | Purpose |
|------|----------|---------|
| 25 | TCP | SMTP (mail relay) |
| 587 | TCP | SMTP Submission (authenticated) |
| 465 | TCP | SMTPS (implicit TLS) |
| 143 | TCP | IMAP (STARTTLS) |
| 993 | TCP | IMAPS (implicit TLS) |
| 8980 | TCP | HTTP (Admin/Webmail UI) |

## Docker Compose Services

### Main Service (gomailserver)
- Production-optimized gomailserver container
- All mail protocols and web interface
- Integrated ClamAV and SpamAssassin support

### Optional Services

#### ClamAV (Antivirus)
```yaml
clamav:
  image: clamav/clamav:latest
  volumes:
    - clamav_data:/var/lib/clamav
```

#### Redis (Caching)
```yaml
redis:
  image: redis:7-alpine
  command: redis-server --maxmemory 128mb
  volumes:
    - redis_data:/data
```

## Backup and Recovery

### Automated Backups in Docker

```bash
# Run backup inside container
docker exec gomailserver /opt/gomailserver/scripts/backup.sh

# Copy backups to host
docker cp gomailserver:/var/lib/gomailserver/backups ./backups

# Restore from backup
docker cp ./backups/gomailserver_db_*.db gomailserver:/var/lib/gomailserver/mailserver.db
```

### Volume Backups

```bash
# Backup volumes
docker run --rm -v gomailserver_data:/data -v $(pwd):/backup alpine tar czf /backup/data.tar.gz -C /data .

# Restore volumes
docker run --rm -v gomailserver_data:/data -v $(pwd):/backup alpine tar xzf /backup/data.tar.gz -C /data
```

## Monitoring and Troubleshooting

### Container Logs

```bash
# View container logs
docker-compose logs -f gomailserver

# View specific service logs
docker logs -f gomailserver
```

### Health Checks

```bash
# Check container health
docker ps

# Manual health check
curl http://localhost:8980/api/v1/health
```

### Debugging

```bash
# Access container shell
docker exec -it gomailserver sh

# View configuration
docker exec gomailserver cat /etc/gomailserver/gomailserver.yaml

# Check database
docker exec gomailserver sqlite3 /var/lib/gomailserver/mailserver.db ".tables"
```

### Common Issues

#### Port Conflicts
```bash
# Check if ports are in use
netstat -tulpn | grep :25

# Stop conflicting services
sudo systemctl stop postfix
sudo systemctl disable postfix
```

#### Permission Issues
```bash
# Fix volume permissions
docker exec gomailserver chown -R gomailserver:gomailserver /var/lib/gomailserver
```

#### Memory Issues
```bash
# Check memory usage
docker stats gomailserver

# Increase memory limit in docker-compose.yml
deploy:
  resources:
    limits:
      memory: 2G
```

## Production Deployment

### System Requirements
- Docker Engine 20.10+
- Docker Compose 2.0+
- 2GB RAM minimum
- 10GB storage minimum
- Linux kernel 4.0+

### Security Considerations
- Use Docker secrets for sensitive configuration
- Implement proper firewall rules
- Regular security updates
- Monitor container logs
- Backup data regularly

### High Availability
- Use Docker Swarm or Kubernetes for clustering
- Implement load balancing for multiple instances
- Use external database for multi-instance deployments
- Configure proper health checks and auto-healing

### Performance Tuning

#### Resource Limits
```yaml
deploy:
  resources:
    limits:
      memory: 1G
      cpus: '2.0'
    reservations:
      memory: 512M
      cpus: '0.5'
```

#### Network Optimization
- Use host networking mode for better performance
- Configure proper MTU settings
- Implement connection pooling

#### Storage Optimization
- Use SSD storage for volumes
- Implement proper backup strategies
- Monitor disk space usage

## Updating

### Update Procedure

```bash
# Pull latest image
docker-compose pull

# Stop services
docker-compose down

# Start with new version
docker-compose up -d

# Check logs
docker-compose logs -f
```

### Rolling Updates

```bash
# Zero-downtime updates
docker-compose up -d --scale gomailserver=2

# Wait for new container to be healthy
sleep 30

# Scale down old container
docker-compose up -d --scale gomailserver=1
```

## Customization

### Custom Configuration

```bash
# Copy custom config
docker cp gomailserver.yaml gomailserver:/etc/gomailserver/gomailserver.yaml

# Restart container
docker-compose restart gomailserver
```

### Custom Dockerfile

```dockerfile
FROM gomailserver:latest

# Add custom SSL certificates
COPY cert.pem /etc/ssl/certs/
COPY key.pem /etc/ssl/private/

# Add custom configuration
COPY gomailserver.yaml /etc/gomailserver/

# Custom entrypoint
COPY docker-entrypoint.sh /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
```

### Environment-Specific Builds

```bash
# Build for different architectures
docker build --platform linux/amd64 -t gomailserver:amd64 .
docker build --platform linux/arm64 -t gomailserver:arm64 .

# Multi-arch manifest
docker manifest create gomailserver:latest \
  gomailserver:amd64 \
  gomailserver:arm64
```

## Integration Examples

### Reverse Proxy (nginx)

```nginx
server {
    listen 80;
    server_name mail.example.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name mail.example.com;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/mail.example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/mail.example.com/privkey.pem;

    # Proxy to gomailserver
    location / {
        proxy_pass http://127.0.0.1:8980;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### Monitoring (Prometheus)

```yaml
scrape_configs:
  - job_name: 'gomailserver'
    static_configs:
      - targets: ['localhost:8980']
    metrics_path: '/api/v1/metrics'
    scrape_interval: 30s
```

### Log Aggregation (ELK Stack)

```yaml
# Filebeat configuration
filebeat.inputs:
- type: docker
  containers.ids:
    - gomailserver
  processors:
  - add_docker_metadata:
      host: "unix:///var/run/docker.sock"

output.elasticsearch:
  hosts: ["elasticsearch:9200"]
```

This Docker setup provides a secure, scalable, and maintainable deployment of gomailserver for production use.