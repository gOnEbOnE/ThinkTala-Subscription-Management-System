#!/bin/sh
# entrypoint.sh — use baked certs, or decode from env as fallback

set -e

mkdir -p certs

# Use baked certs if they exist (preferred)
if [ -f "certs/private.pem" ] && [ -f "certs/public.pem" ]; then
    echo "[entrypoint] Using baked JWT keys from certs/"
else
    # Decode RSA private key from env
    if [ -n "$JWT_PRIVATE_KEY_B64" ]; then
        printf '%s' "$JWT_PRIVATE_KEY_B64" | tr -d '"\n\r ' | base64 -d > certs/private.pem
        echo "[entrypoint] JWT private key loaded from env"
    elif [ ! -f "certs/private.pem" ]; then
        echo "[entrypoint] FATAL: JWT_PRIVATE_KEY_B64 not set and certs/private.pem not found"
        exit 1
    fi

    # Decode RSA public key from env
    if [ -n "$JWT_PUBLIC_KEY_B64" ]; then
        printf '%s' "$JWT_PUBLIC_KEY_B64" | tr -d '"\n\r ' | base64 -d > certs/public.pem
        echo "[entrypoint] JWT public key loaded from env"
    elif [ ! -f "certs/public.pem" ]; then
        echo "[entrypoint] FATAL: JWT_PUBLIC_KEY_B64 not set and certs/public.pem not found"
        exit 1
    fi
fi

exec /app/subscription "$@"
