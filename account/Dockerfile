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
    go build -ldflags="-s -w" -o account .

# =====================
# RUNTIME STAGE
# =====================
FROM alpine:latest

WORKDIR /app

# Install timezone data
RUN apk --no-cache add tzdata

# Copy binary dari builder
COPY --from=builder /app/account /app/account
RUN chmod +x /app/account

# Copy key pair and TLS assets baked in repository
COPY certs/ ./certs/

COPY public/ ./public/

COPY entrypoint.sh /app/entrypoint.sh
RUN tr -d '\r' < /app/entrypoint.sh > /tmp/ep && mv /tmp/ep /app/entrypoint.sh && chmod +x /app/entrypoint.sh

# Expose port aplikasi
EXPOSE 2001

ENTRYPOINT ["/app/entrypoint.sh"]