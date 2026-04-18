FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o subscription main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata ca-certificates
COPY --from=builder /app/subscription .
COPY entrypoint.sh .
COPY certs/ ./certs/
RUN chmod +x /app/subscription && chmod +x entrypoint.sh

EXPOSE 5004
ENTRYPOINT ["/app/entrypoint.sh"]