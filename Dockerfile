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
# - nodejs 20 + typescript:   Node.js 20 (NodeSource) plus the official
#                              `typescript` package providing `tsc`. We use
#                              tsc (pure JavaScript) instead of tsx /
#                              esbuild / swc / Node's --experimental-strip-types,
#                              because every "fast" TS transpiler ships a
#                              WASM module and V8's WASM sandbox fails inside
#                              our docker runner with `WebAssembly.instantiate():
#                              Out of memory`. tsc is slower to start (~1s)
#                              but has zero WASM dependencies and just works.
# - ca-certificates + curl + gnupg: needed for the NodeSource setup script.
RUN apt-get update && apt-get install -y --no-install-recommends \
    default-jdk-headless \
    python3 \
    golang \
    ca-certificates \
    curl \
    gnupg \
    && curl -fsSL https://deb.nodesource.com/setup_20.x | bash - \
    && apt-get install -y --no-install-recommends nodejs \
    && npm install -g typescript@5.6.3 \
    && apt-get purge -y curl gnupg \
    && apt-get autoremove -y \
    && rm -rf /var/lib/apt/lists/* /root/.npm

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
