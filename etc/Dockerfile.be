# ── Production Backend ─────────────────────────────────────────────
# Multi-stage build: compile Go binary, run FROM alpine (~20 MB)
# Alpine dipilih karena:
#   - Punya shell (sh) untuk entrypoint script
#   - crond (busybox) sudah built-in — siap pakai jika butuh cron job
# ──────────────────────────────────────────────────────────────────

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /src

COPY be/go.mod be/go.sum ./
RUN go mod download

COPY be/ .

RUN mkdir -p /src/storage/keys

# -ldflags="-s -w"     strip debug info → smaller binary
# -tags timetzdata     embed timezone DB into binary (no need OS tzdata)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/server ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /opt/genkey ./cmd/genkey

# ── Runtime ───────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    tzdata

COPY --from=builder /opt/server /server
COPY --from=builder /opt/genkey /genkey
COPY --from=builder /src/storage /app/storage
COPY etc/entrypoint-prod.sh /entrypoint.sh

RUN chmod +x /entrypoint.sh

WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
