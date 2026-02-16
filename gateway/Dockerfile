# Build Stage
FROM golang:alpine AS builder
WORKDIR /app
# Install git jika package kamu ada di private repo atau github
RUN apk add --no-cache git 
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o thinknalyze-gateway main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata ca-certificates

# Copy binary
COPY --from=builder /app/thinknalyze-gateway .

# Copy semua config wajib
COPY routes.json .
COPY .env .

# Metadata
ENV port=2000
EXPOSE 2000

CMD ["./thinknalyze-gateway"]