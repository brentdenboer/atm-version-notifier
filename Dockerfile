# Use the official Go image as the base image
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o atm10-version-monitor

# Use a minimal alpine image for the final container
FROM alpine:latest

# Add metadata labels
LABEL maintainer="ATM10 Version Monitor Team" \
      description="Monitors ATM10 modpack version changes and sends Discord notifications" \
      version="1.0.0"

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create user with UID 1000 and GID 1000 to match itzg/minecraft-server
RUN addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/atm10-version-monitor .

RUN chmod +x atm10-version-monitor

# Create directories for volumes and set permissions
RUN mkdir -p /data /reference-data && \
    chown -R appuser:appgroup /app /data /reference-data

# Document environment variables
ENV DISCORD_WEBHOOK_URL="" \
    REFERENCE_FILE_PATH="/reference-data/version_reference.json" \
    CONFIG_FILE_PATH="/data/config/bcc-common.toml" \
    FILE_CHECK_INTERVAL_SECONDS="60" \
    RECHECK_TIMEOUT_SECONDS="30"

# Declare volumes
VOLUME ["/data", "/reference-data"]

# Switch to non-root user
USER appuser

# Add healthcheck
HEALTHCHECK --interval=60s --timeout=3s --start-period=5s --retries=3 \
    CMD pgrep atm10-version-monitor || exit 1

# Set the entry point
ENTRYPOINT ["/app/atm10-version-monitor"]
