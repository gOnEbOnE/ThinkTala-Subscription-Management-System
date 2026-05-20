#!/bin/bash
# Deploy script untuk Ubuntu Server
# Tempat: ~/propensi/deploy/server-yovan/deploy.sh

set -e

echo "=== Starting Deployment Process ==="
echo "Timestamp: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Variables
REPO_PATH="$HOME/propensi"
DEPLOY_DIR="$REPO_PATH/deploy/server-yovan"
LOG_FILE="$HOME/deployment_$(date +%Y%m%d_%H%M%S).log"

# capture all output to log
exec > >(tee -a "$LOG_FILE") 2>&1

# Change to repo directory
cd "$REPO_PATH"

echo "[1/4] Pulling latest changes from GitHub..."
git fetch github server-yovan
git checkout server-yovan
# avoid pull conflicts: ensure working tree exactly matches remote branch
git reset --hard github/server-yovan
git clean -fd

echo "[2/4] Verifying deployment configuration..."
if [ ! -f "$DEPLOY_DIR/operational.env" ]; then
    echo "ERROR: operational.env not found!"
    exit 1
fi

# Verify critical config
if ! grep -q "SUBSCRIPTION_SERVICE_URL" "$DEPLOY_DIR/operational.env"; then
    echo "WARNING: SUBSCRIPTION_SERVICE_URL not found in operational.env"
    echo "Adding: SUBSCRIPTION_SERVICE_URL=http://subscription:5004"
    echo "SUBSCRIPTION_SERVICE_URL=http://subscription:5004" >> "$DEPLOY_DIR/operational.env"
fi

echo "[3/4] Building Docker images..."
docker compose -f docker-compose.yml -f docker-compose.server.yml build --no-cache

echo "[4/4] Starting services..."
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --remove-orphans --force-recreate --build

echo "[5/5] Recreating notification service to ensure updated envs are applied..."
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --no-deps --force-recreate --build notification || true

echo "Verifying notification container environment (redis & smtp vars):"
if docker ps --format '{{.Names}}' | grep -q thinknalyze-notification; then
    docker exec thinknalyze-notification env | grep -i redis || true
    docker exec thinknalyze-notification env | grep -i smtp || true
else
    echo "Notification container not running yet."
fi

echo ""
echo "=== Deployment Completed Successfully ==="
echo "Log: $LOG_FILE"
echo ""
echo "Checking service status..."
docker compose ps

echo ""
echo "Service health check:"
sleep 2
docker compose logs --tail=5 operational
echo "Recent notification logs (tail 200):"
docker compose logs --tail=200 notification || true
