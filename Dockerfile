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
# - nodejs 22:            Node.js 22 (NodeSource) for TypeScript via the
#                          built-in --experimental-strip-types flag. We do
#                          NOT use tsx because every tsx version depends on
#                          the es-module-lexer WASM build, which fails inside
#                          docker with `WebAssembly.instantiate(): Out of
#                          memory` due to seccomp / cgroup memory limits on
#                          WASM allocations. Node 22's strip-types is pure
#                          C++/JS, no WASM, and handles the type-erasure-only
#                          syntax our student solutions use.
# - ca-certificates + curl + gnupg: needed for the NodeSource setup script.
RUN apt-get update && apt-get install -y --no-install-recommends \
    default-jdk-headless \
    python3 \
    golang \
    ca-certificates \
    curl \
    gnupg \
    && curl -fsSL https://deb.nodesource.com/setup_22.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && apt-get purge -y curl gnupg \
    && apt-get autoremove -y \
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
