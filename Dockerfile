# Dockerfile
FROM golang:1.24.2-alpine as builder

WORKDIR /app

# Install Git for Go modules
RUN apk add --no-cache git

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o banking-ledger ./cmd/main.go

# Final minimal image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/banking-ledger .

# Expose the service port
EXPOSE 8080

# Start the app
CMD ["./banking-ledger"]
