# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (opsional, untuk private dependency)
RUN apk add --no-cache git

# Copy go mod files first untuk caching layer
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Build dengan optimization
# - CGO_ENABLED=0: Static binary, tidak bergantung C library
# - GOOS=linux: Target OS Linux
# - GOARCH=amd64: Target architecture
# - -ldflags="-w -s": Hilangkan DWARF debug info dan symbol table
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o main cmd/server/main.go

# Run stage - Final image yang akan digunakan
FROM alpine:latest

# Install packages yang diperlukan:
# - ca-certificates: Untuk HTTPS request
# - tzdata: Timezone data
# - wget: Untuk healthcheck command
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /app

# Buat user non-root untuk keamanan
RUN adduser -D -g '' appuser

# Copy binary dari builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/config/local.yml ./config/local.yml

# Change ownership ke appuser
RUN chown -R appuser:appuser /app

# Switch ke non-root user
USER appuser

# Expose port 3020
EXPOSE 3020

# Health check - check endpoint /health setiap 30 detik
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3020/health || exit 1

# Run aplikasi
CMD ["./main"]
