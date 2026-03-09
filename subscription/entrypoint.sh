#!/bin/sh
# entrypoint.sh — decode JWT keys from Railway env vars before starting subscription service

set -e

mkdir -p certs

# Decode RSA private key (strip ALL whitespace + quotes in case Railway line-wraps long vars)
if [ -n "$JWT_PRIVATE_KEY_B64" ]; then
    printf '%s' "$JWT_PRIVATE_KEY_B64" | tr -d '"\n\r ' | base64 -d > certs/private.pem
    echo "[entrypoint] JWT private key loaded from env"
elif [ ! -f "certs/private.pem" ]; then
    echo "[entrypoint] FATAL: JWT_PRIVATE_KEY_B64 not set and certs/private.pem not found"
    exit 1
fi

# Decode RSA public key (strip ALL whitespace + quotes in case Railway line-wraps long vars)
if [ -n "$JWT_PUBLIC_KEY_B64" ]; then
    printf '%s' "$JWT_PUBLIC_KEY_B64" | tr -d '"\n\r ' | base64 -d > certs/public.pem
    echo "[entrypoint] JWT public key loaded from env"
elif [ ! -f "certs/public.pem" ]; then
    echo "[entrypoint] FATAL: JWT_PUBLIC_KEY_B64 not set and certs/public.pem not found"
    exit 1
fi

exec /app/subscription "$@"
