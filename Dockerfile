# ----------- Stage 1: Build -----------
FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o api ./cmd/api

# Build Migration binary
RUN GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o migration ./cmd/migrate


# ----------- Stage 2: Runtime -----------
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    libc6 \
    wkhtmltopdf \
    libreoffice \
    fonts-dejavu \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binaries
COPY --from=builder /app/api .
COPY --from=builder /app/migration .

# Create required folders
RUN mkdir -p /app/storage /app/logs

# Create non-root user
RUN useradd --system --no-create-home --uid 1001 appuser && \
    chown -R appuser:appuser /app

USER appuser

HEALTHCHECK CMD curl -f http://localhost:6014/health || exit 1

# Application Port
EXPOSE 6014

ENTRYPOINT ["./api"]
