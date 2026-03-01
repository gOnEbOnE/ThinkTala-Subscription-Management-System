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
    go build -ldflags="-s -w" -o notification .

# =====================
# RUNTIME STAGE
# =====================
FROM alpine:latest

WORKDIR /app

# Install timezone data
RUN apk --no-cache add tzdata

# Copy binary dari builder
COPY --from=builder /app/notification /app/notification
RUN chmod +x /app/notification

# Copy env file (opsional)
# Copy env file
COPY .env .env

# ---> TAMBAHKAN BARIS INI UNTUK MENGCOPY FILE PEM <---
COPY private.pem public.pem ./

COPY public/ ./public/

# Expose port aplikasi
EXPOSE 2002

# Jalankan aplikasi
CMD ["/app/notification"]