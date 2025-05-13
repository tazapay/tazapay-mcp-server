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
RUN echo | openssl s_client -showcerts -connect api-orange.tazapay.com:443 2>/dev/null \
    | awk '/-----BEGIN CERTIFICATE-----/,/-----END CERTIFICATE-----/ { print }' \
    > /usr/local/share/ca-certificates/tazapay.crt

# Update CA trust store
RUN update-ca-certificates

# Set default log file path (can be overridden during runtime)
ENV LOG_FILE_PATH=/app/logs/app.log

# Ensure the log directory exists
RUN mkdir -p /app/logs

# Optional: Final CMD
CMD ["bash"]

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/tazapay-mcp-server .

# Run the binary on container start
ENTRYPOINT ["/app/tazapay-mcp-server"]
