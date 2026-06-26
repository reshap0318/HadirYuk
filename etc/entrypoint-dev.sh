#!/bin/sh
set -e

PRIVATE_KEY="${JWT_PRIVATE_KEY_PATH:-storage/keys/private.pem}"
PUBLIC_KEY="${JWT_PUBLIC_KEY_PATH:-storage/keys/public.pem}"
PASSPHRASE="${JWT_PASSPHRASE_PATH:-storage/keys/passphrase}"

if [ ! -f "$PRIVATE_KEY" ] || [ ! -f "$PUBLIC_KEY" ] || [ ! -f "$PASSPHRASE" ]; then
    echo "JWT keys not found, generating..."
    go run ./cmd/genkey -f
fi

# Start busybox crond in background
crond -l 8
echo "crond started"

exec air -c .air.toml
