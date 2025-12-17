FROM golang:1.24.10-alpine AS builder

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o service-a ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/service-a .
EXPOSE 8080
ENTRYPOINT ["/app/service-a"]
