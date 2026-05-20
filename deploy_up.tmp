#!/bin/sh
set -e
cd "$(dirname "$0")/../.."
git pull origin server-yovan
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
