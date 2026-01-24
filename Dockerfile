# gomailserver Production Docker Image
# Multi-stage build for security and size optimization

# Build stage
FROM golang:1.23.5-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    && update-ca-certificates

# Set working directory
WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with security flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a \
    -installsuffix cgo \
    -o gomailserver \
    ./cmd/gomailserver

# Verify the binary
RUN ./gomailserver version

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    clamav \
    clamav-daemon \
    spamassassin \
    && update-ca-certificates \
    && mkdir -p /var/lib/gomailserver /var/log/gomailserver /etc/gomailserver \
    && addgroup -g 1000 -S gomailserver \
    && adduser -u 1000 -S gomailserver -G gomailserver \
    && chown -R gomailserver:gomailserver /var/lib/gomailserver /var/log/gomailserver

# Copy the binary from builder stage
COPY --from=builder /build/gomailserver /usr/local/bin/gomailserver

# Copy default configuration
COPY --from=builder /build/gomailserver.example.yaml /etc/gomailserver/gomailserver.yaml

# Copy scripts for backup and management
COPY --from=builder /build/scripts/ /opt/gomailserver/scripts/

# Make scripts executable
RUN chmod +x /opt/gomailserver/scripts/*.sh

# Set proper permissions
RUN chmod 755 /usr/local/bin/gomailserver \
    && chmod 644 /etc/gomailserver/gomailserver.yaml

# Create volume mount points
VOLUME ["/var/lib/gomailserver", "/var/log/gomailserver", "/etc/gomailserver"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8980/api/v1/health || exit 1

# Expose ports
EXPOSE 25/tcp   # SMTP
EXPOSE 587/tcp  # SMTP Submission
EXPOSE 465/tcp  # SMTPS
EXPOSE 143/tcp  # IMAP
EXPOSE 993/tcp  # IMAPS
EXPOSE 8980/tcp # HTTP Admin/Webmail

# Switch to non-root user
USER gomailserver:gomailserver

# Set working directory
WORKDIR /var/lib/gomailserver

# Default command
CMD ["/usr/local/bin/gomailserver", "run", "--config", "/etc/gomailserver/gomailserver.yaml"]