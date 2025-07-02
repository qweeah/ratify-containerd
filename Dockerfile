FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY main.go .

# Build the watcher binary
RUN go build -o watcher main.go

# Use a minimal image for runtime
FROM alpine:3.19

WORKDIR /app

# Copy the watcher binary from builder
COPY --from=builder /app/watcher /app/watcher

# Create the shared-data directory (optional, for mount)
RUN mkdir -p /shared-data

# Set the entrypoint to the watcher binary
ENTRYPOINT ["/app/watcher"]