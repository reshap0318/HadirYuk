#!/bin/sh
set -e

PRIVATE_KEY="${JWT_PRIVATE_KEY_PATH:-storage/keys/private.pem}"
PUBLIC_KEY="${JWT_PUBLIC_KEY_PATH:-storage/keys/public.pem}"

if [ ! -f "$PRIVATE_KEY" ] || [ ! -f "$PUBLIC_KEY" ]; then
    echo "JWT keys not found, generating..."
    go run ./cmd/genkey -f
fi

# Start busybox crond in background
crond -l 8
echo "crond started"

exec air -c .air.toml
