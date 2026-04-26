# Stage 1: Build
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# -ldflags="-w -s" reduces binary size
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o kshanik-search ./cmd/server/

# Stage 2: Final image
FROM alpine:latest

# Install CA certificates for HTTPS requests to search providers
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/kshanik-search .

# Copy environment example (user will provide .env via volume or env_file)
COPY .env.example .env

# Expose port
EXPOSE 8080

# Run the application
CMD ["./kshanik-search"]
