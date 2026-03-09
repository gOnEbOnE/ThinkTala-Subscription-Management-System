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

# Copy entrypoint (decodes JWT keys from Railway env vars at startup)
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

COPY public/ ./public/

RUN mkdir -p certs

# Expose port aplikasi
EXPOSE 2006

# Gunakan entrypoint untuk decode JWT keys dari env vars Railway
ENTRYPOINT ["/app/entrypoint.sh"]