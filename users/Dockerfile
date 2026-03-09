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

# Generate RSA keys for JWT + self-signed cert (fallback jika app_secure=true)
RUN mkdir -p certs && \
    openssl genrsa -out certs/private.pem 2048 && \
    openssl rsa -in certs/private.pem -pubout -out certs/public.pem && \
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout certs/localhost.key -out certs/localhost.crt \
        -subj "/CN=localhost"

COPY public/ ./public/

# Expose port aplikasi
EXPOSE 2006

# Jalankan aplikasi
CMD ["/app/users"]