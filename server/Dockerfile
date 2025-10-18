# syntax=docker/dockerfile:1

# -------- Build Stage --------
FROM golang:1.24 AS builder

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Build the Go binary into /app/build/
RUN mkdir -p build && go build -o build/server ./cmd/server/main.go

# -------- Runtime Stage --------
FROM golang:1.24

WORKDIR /app

# Copy binary and migration files from builder
COPY --from=builder /app/build/server ./build/server
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy .env file
COPY .env .env

# Load environment variables from .env
ENV $(cat .env | xargs)

# Expose app port
EXPOSE 8081

# Run migrations before starting the server
CMD goose -dir $GOOSE_MIGRATION_DIR $GOOSE_DRIVER "$GOOSE_DBSTRING" up && ./build/server
