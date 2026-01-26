# gomailserver Production Docker Image
# Multi-stage build for security and size optimization

# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies including CGO requirements for go-sqlite3
RUN apk add --no-cache \
    git \
    ca-certificates \
    tzdata \
    gcc \
    musl-dev \
    sqlite-dev \
    && update-ca-certificates

# Set working directory
WORKDIR /build

# Copy source code first (needed for local replace directives)
COPY . .

RUN go mod download

# Build the application with CGO enabled for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags='-w -s -linkmode external -extldflags "-static"' \
    -a \
    -o gomailserver \
    ./cmd/gomailserver

# Verify the binary
RUN ./gomailserver version

# Runtime stage
FROM alpine:latest

# Install runtime dependencies (including sqlite for go-sqlite3)
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite-libs \
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
EXPOSE 25/tcp
EXPOSE 587/tcp
EXPOSE 465/tcp
EXPOSE 143/tcp
EXPOSE 993/tcp
EXPOSE 8980/tcp

# Switch to non-root user
USER gomailserver:gomailserver

# Set working directory
WORKDIR /var/lib/gomailserver

# Default command
CMD ["/usr/local/bin/gomailserver", "run", "--config", "/etc/gomailserver/gomailserver.yaml"]