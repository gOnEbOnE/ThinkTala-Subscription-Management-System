# Build Stage - context harus root repo agar ../account tersedia
FROM golang:alpine AS builder
WORKDIR /repo

RUN apk add --no-cache git

# Copy seluruh repo agar replace ../account bisa resolve
COPY . .

WORKDIR /repo/gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o thinknalyze-gateway main.go

# Run Stage
FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata ca-certificates

COPY --from=builder /repo/gateway/thinknalyze-gateway .
COPY --from=builder /repo/gateway/routes.json .
COPY --from=builder /repo/gateway/entrypoint.sh .
RUN chmod +x entrypoint.sh && mkdir -p certs
COPY --from=builder /repo/gateway/certs/*.pem certs/

ENV port=2000
EXPOSE 2000

ENTRYPOINT ["./entrypoint.sh"]
