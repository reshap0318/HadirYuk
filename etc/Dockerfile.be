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

# -ldflags="-s -w"     strip debug info → smaller binary
# -tags timetzdata     embed timezone DB into binary (no need OS tzdata)
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/server ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /opt/genkey ./cmd/genkey

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /opt/cleartmp ./cmd/cleartmp

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/markabsent ./cmd/markabsent

# DEMO/DEV ONLY — fabricates attendance rows for the portfolio demo
# employees. Built into the image so it can be run manually (`docker exec
# <container> /dummy`) on a demo/staging deployment, but deliberately NOT
# added to etc/crontab — it must never auto-run against a real database.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/dummy ./cmd/dummy

# Migration CLI (up/down/seed/dummy/refresh). Built into the image for manual
# use (`docker exec <container> /migration up`), but deliberately NOT called
# from entrypoint.sh — its "up" command unconditionally drops the
# attendances table on every run (see cmd/migration/main.go), so auto-running
# it on every container start would wipe attendance data on every restart.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -tags timetzdata \
    -o /opt/migration ./cmd/migration

# ── Runtime ───────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Set system timezone so busybox crond uses local time correctly
RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && \
    echo "Asia/Jakarta" > /etc/timezone

COPY --from=builder /opt/server     /server
COPY --from=builder /opt/genkey     /genkey
COPY --from=builder /opt/cleartmp   /cleartmp
COPY --from=builder /opt/markabsent /markabsent
COPY --from=builder /opt/dummy      /dummy
COPY --from=builder /opt/migration  /migration
COPY etc/entrypoint.sh              /entrypoint.sh
COPY etc/crontab                    /etc/crontabs/root

RUN chmod +x /entrypoint.sh && mkdir -p /var/log

WORKDIR /app
EXPOSE 8080

ENTRYPOINT ["/entrypoint.sh"]
