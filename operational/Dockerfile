FROM golang:alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o operational main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata ca-certificates
COPY --from=builder /app/operational .
COPY --from=builder /app/public ./public
RUN chmod +x /app/operational

EXPOSE 5005
CMD ["/app/operational"]