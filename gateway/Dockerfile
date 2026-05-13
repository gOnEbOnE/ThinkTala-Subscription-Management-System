# Build Stage — go.mod is 1.25.x; default golang:alpine can be older, which breaks -mod=vendor parsing.
FROM golang:alpine AS builder
ENV GOTOOLCHAIN=go1.25.5
WORKDIR /app
# Install git jika package kamu ada di private repo atau github
RUN apk add --no-cache git 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -ldflags="-s -w" -o thinknalyze-gateway main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata ca-certificates

# Copy binary
COPY --from=builder /app/thinknalyze-gateway .

# JWT keys (see main.go InitJWTLoadKeys paths)
COPY certs/ ./certs/

# Copy semua config wajib
# Always use routes.railway.json for Docker (production deployments)
COPY routes.railway.json routes.json

# Metadata
COPY entrypoint.sh .
RUN tr -d '\r' < /app/entrypoint.sh > /tmp/ep && mv /tmp/ep /app/entrypoint.sh && chmod +x /app/entrypoint.sh

ENV port=2000
EXPOSE 2000

ENTRYPOINT ["/app/entrypoint.sh"]