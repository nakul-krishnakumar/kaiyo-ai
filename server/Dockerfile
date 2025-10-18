# syntax=docker/dockerfile:1

# ------------ builder (prod binary) ------------
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY cmd/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out/server ./cmd/server/main.go

# ------------ dev (hot-reload) ------------
FROM golang:1.24-alpine AS dev
RUN apk add --no-cache git ca-certificates tzdata
# Install Air and Goose into the Go bin
RUN go install github.com/air-verse/air@v1.62.0 \
 && go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app
# Bind-mount your source & migrations at runtime:
#   docker run -v "$(pwd)":/app ...
# so no COPY here

ENV PATH=$PATH:/go/bin
EXPOSE 8081
# wait for DB, run migrations, then Air
CMD ["sh","-c", "\
  until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do sleep 1; done && \
  goose -dir internal/database/migrations $GOOSE_DRIVER \"$GOOSE_DBSTRING\" up && \
  air --port 8081\
"]

# ------------ prod (minimal runtime) ------------
FROM gcr.io/distroless/static-debian12:nonroot AS prod
WORKDIR /app
COPY --from=builder /out/server /app/server
EXPOSE 8081
USER nonroot
ENTRYPOINT ["/app/server"]
