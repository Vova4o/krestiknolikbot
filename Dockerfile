### Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build static binary for linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

### Runtime stage
FROM alpine:3.19
WORKDIR /app

RUN apk add --no-cache ca-certificates

# Copy binary and assets
COPY --from=builder /app/server /app/server
COPY --from=builder /app/web /app/web

ENV GIN_MODE=release
ENV PORT=8080

EXPOSE 8080

CMD ["/app/server"]
