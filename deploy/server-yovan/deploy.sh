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

# Change to repo directory
cd "$REPO_PATH"

echo "[1/4] Pulling latest changes from GitHub..."
git fetch github server-yovan
git checkout server-yovan
git pull github server-yovan --rebase

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
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d

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
