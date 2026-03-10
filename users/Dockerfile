# =====================
# BUILD STAGE
# =====================
FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency dulu (cache optimal)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -ldflags="-s -w" -o users .

# =====================
# RUNTIME STAGE
# =====================
FROM alpine:latest

WORKDIR /app

# Install timezone data
RUN apk --no-cache add tzdata openssl

# Copy binary dari builder
COPY --from=builder /app/users /app/users
RUN chmod +x /app/users

# Copy entrypoint
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

COPY public/ ./public/
# Copy pre-committed RSA key pair (baked into image)
COPY certs/ ./certs/

RUN mkdir -p certs

# Ensure uploads directory exists at build time
# (Railway Volume should be mounted at /app/public/uploads via Railway dashboard)
RUN mkdir -p /app/public/uploads/kyc

# Expose port aplikasi
EXPOSE 2006

# Gunakan entrypoint untuk decode JWT keys dari env vars Railway
ENTRYPOINT ["/app/entrypoint.sh"]