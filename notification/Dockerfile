FROM golang:alpine AS builder
WORKDIR /app
COPY . .
# Download dependensi jika ada go.mod
RUN if [ -f go.mod ]; then go mod download; fi
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add tzdata
COPY --from=builder /app/app .

ENV PORT=2001
EXPOSE 2001
CMD ["./app"]