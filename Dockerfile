# Build Stage (Using Golang image)
FROM golang:1.24.2-alpine AS builder

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
FROM debian:bookworm-slim

# Set working directory in the minimal image
WORKDIR /app

# Install ca-certificates to allow trust of custom root certs
RUN apt-get update && apt-get install -y ca-certificates

# Copy the custom certificate (tazapay-chain.crt) into the container
COPY tazapay.crt /usr/local/share/ca-certificates/tazapay.crt

# Update the CA certificates (this will include the custom CA)
RUN update-ca-certificates

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/tazapay-mcp-server .

# Run the binary on container start
ENTRYPOINT ["/app/tazapay-mcp-server"]
