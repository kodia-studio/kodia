# Multi-stage Dockerfile for Kodia Backend
# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary (statically linked, Linux target)
RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-s -w" \
  -o server \
  cmd/server/main.go

---

# Stage 2: Runtime stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata curl

# Create non-root user for security
RUN addgroup -g 1001 kodia && \
    adduser -D -u 1001 -G kodia kodia

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .

# Copy config template for reference (optional)
COPY --from=builder /app/configs ./configs

# Set ownership to kodia user
RUN chown -R kodia:kodia /app

# Switch to non-root user
USER kodia

# Expose application port (app) and metrics port (Prometheus)
EXPOSE 8080 9090

# Health check - probe the /health endpoint
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]
