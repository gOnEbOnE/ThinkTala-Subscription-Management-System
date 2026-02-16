FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN if [ -f go.mod ]; then go mod download; fi
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata
COPY --from=builder /app/app .

ENV PORT=2003
EXPOSE 2003
CMD ["./app"]