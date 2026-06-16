# ── Production Backend ────────────────────────────────────────────
# Multi-stage build: compile Go binary, run FROM scratch (~15 MB)
# ──────────────────────────────────────────────────────────────────

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN mkdir -p /src/storage

# -ldflags="-s -w"     strip debug info → smaller binary
# -tags timetzdata     embed timezone DB into binary (no need OS tzdata)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/server ./cmd/api

# ── Runtime (scratch = 0 MB, total image ~15 MB) ──────────────────
FROM scratch

# CA certificates for HTTPS (SMTP, external APIs)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Binary
COPY --from=builder /opt/server /server

# Storage directories (attendance evidence, face photos, keys, logs, tmp)
COPY --from=builder /src/storage /app/storage

WORKDIR /app
EXPOSE 8080

CMD ["/server"]
