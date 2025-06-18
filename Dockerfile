# Build Stage (Using Golang image)
FROM golang:1.24.3-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Copy go mod files and download dependencies early for better cache use
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go binary (statically linked by default in Go)
RUN CGO_ENABLED=0 GOOS=linux go build -o tazapay-mcp-server ./cmd/server

# Runtime Stage (Minimal Image)
FROM debian:stable-slim

# Set working directory
WORKDIR /app

# Install required packages
RUN apt-get update && apt-get install -y --no-install-recommends \
    openssl \
    ca-certificates \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Fetch and store the certificate
RUN echo | openssl s_client -showcerts -connect service.tazapay.com:443 2>/dev/null \
    | awk '/-----BEGIN CERTIFICATE-----/,/-----END CERTIFICATE-----/ { print }' \
    > /usr/local/share/ca-certificates/tazapay.crt

# Update CA trust store
RUN update-ca-certificates

# Set default log file path (can be overridden during runtime with -e LOG_FILE_PATH=/your/path.log)
ENV LOG_FILE_PATH=/app/logs/app.log
# Set default server type (can be overridden at runtime)
ENV TRANSPORT_TYPE=streamablehttp

# Ensure the log directory exists (default, but if overridden, user must ensure directory exists)
RUN mkdir -p /app/logs

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/tazapay-mcp-server .

# Entrypoint (can be overridden to pass env vars)
CMD ["/app/tazapay-mcp-server"]
