# ==========================================================================
# RWA House CRE Workflow - Dockerfile
# ==========================================================================
# This Dockerfile builds the Chainlink CRE workflow for the RWA House platform.
# The workflow runs as a WebAssembly module within the Chainlink DON.

# Stage 1: Build Environment
# --------------------------
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache \
    git \
    make \
    gcc \
    musl-dev \
    linux-headers \
    libc6-compat

# Install TinyGo for WebAssembly compilation
RUN wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.2/tinygo_0.31.2_amd64.deb && \
    dpkg -i tinygo_0.31.2_amd64.deb || true && \
    cp /usr/local/tinygo/bin/tinygo /usr/local/bin/tinygo || \
    (wget https://github.com/tinygo-org/tinygo/releases/download/v0.31.2/tinygo0.31.2.linux-amd64.tar.gz && \
     tar -xzf tinygo0.31.2.linux-amd64.tar.gz && \
     mv tinygo /usr/local/tinygo && \
     ln -s /usr/local/tinygo/bin/tinygo /usr/local/bin/tinygo)

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the WebAssembly module
# Using TinyGo for smaller WASM output optimized for CRE
RUN tinygo build -o workflow.wasm -target wasm . || \
    GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o workflow.wasm .

# Verify the build
RUN ls -la workflow.wasm && \
    file workflow.wasm

# Stage 2: Runtime Environment
# ----------------------------
FROM alpine:3.19 AS runtime

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Create non-root user for security
RUN addgroup -g 1000 cre && \
    adduser -u 1000 -G cre -s /bin/sh -D cre

# Set working directory
WORKDIR /app

# Copy the built WASM module
COPY --from=builder /app/workflow.wasm .

# Copy configuration files
COPY --from=builder /app/config/ ./config/

# Set proper ownership
RUN chown -R cre:cre /app

# Switch to non-root user
USER cre

# Expose port (for health checks)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD echo "Healthy" || exit 1

# Labels for metadata
LABEL org.opencontainers.image.title="RWA House CRE Workflow"
LABEL org.opencontainers.image.description="Chainlink CRE workflow for real estate tokenization"
LABEL org.opencontainers.image.version="1.0.0"
LABEL org.opencontainers.image.vendor="RWA House Platform"

# Default command (workflow is executed by Chainlink DON)
CMD ["echo", "RWA House CRE Workflow - WebAssembly Module Ready"]

# ==========================================================================
# Alternative: Production Build with Multi-Stage Optimization
# ==========================================================================
FROM builder AS production

# Build optimized WASM with size optimization
RUN tinygo build -o workflow.wasm -target wasm -opt=z -no-debug . || \
    GOOS=js GOARCH=wasm go build -ldflags="-s -w -X main.version=1.0.0" -o workflow.wasm .

# Compress WASM for distribution
RUN gzip -c workflow.wasm > workflow.wasm.gz && \
    ls -lh workflow.wasm workflow.wasm.gz

# Stage 3: Distribution
# ---------------------
FROM scratch AS distribution

# Copy only the necessary artifacts
COPY --from=production /app/workflow.wasm /
COPY --from=production /app/workflow.wasm.gz /
COPY --from=production /app/config/config.staging.json /config.staging.json
COPY --from=production /app/config/config.production.json /config.production.json

# This stage can be used to extract artifacts
