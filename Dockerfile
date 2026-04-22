# Docker Build Context Setup
#
# This Dockerfile builds tinydsa-tester using the published tester-utils from GitHub
# Build from the tinydsa-tester directory:
#   cd tinydsa-tester
#   docker build -t ghcr.io/bootcraft-cn/tinydsa-tester .

# Stage 1: Build the Go binary
FROM golang:1.24-bookworm AS builder

WORKDIR /app

# Copy go module files first for better caching
COPY go.mod go.sum ./

# Download dependencies from GitHub
RUN go mod download

# Copy the rest of the project
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -o tinydsa-tester \
    -ldflags="-s -w" \
    .

# Stage 2: Runtime image with Java, Python, Go and Node.js
FROM debian:bookworm-slim

# Install runtime dependencies for the four supported student languages:
# - default-jdk-headless: Java compiler and runtime (javac/java)
# - python3:               Python interpreter
# - golang:                Go toolchain (for go starter's run.sh)
# - nodejs + npm:          Node.js for TypeScript (tsx / ts-node)
# - ca-certificates:       HTTPS for any solutions hitting the network
RUN apt-get update && apt-get install -y \
    default-jdk-headless \
    python3 \
    golang \
    nodejs \
    npm \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Create a non-root user for running tests
RUN useradd -m -s /bin/bash tester

# Copy the binary from builder
COPY --from=builder /app/tinydsa-tester /usr/local/bin/tinydsa-tester

# Set working directory
WORKDIR /workspace

# Switch to non-root user
USER tester

# Default command shows help
ENTRYPOINT ["tinydsa-tester"]
CMD ["--help"]
