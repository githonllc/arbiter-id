# Build stage
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Copy module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o arbiter-id-service .

# Runtime stage
FROM alpine:latest

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh arbiter

# Set working directory
WORKDIR /home/arbiter

# Copy binary from build stage
COPY --from=builder /app/arbiter-id-service .

# Change file ownership
RUN chown arbiter:arbiter arbiter-id-service

# Switch to non-root user
USER arbiter

# Expose port
EXPOSE 8080

# Set default environment variables
ENV NODE_ID=0
ENV PORT=8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT}/health || exit 1

# Start service
CMD ["./arbiter-id-service"]