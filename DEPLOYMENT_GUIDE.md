# Deployment Guide - Server Yovan

## Current Setup
- **GitLab Remote** (origin): Main development repository
- **GitHub Remote** (github): Backup/deployment branch repository
- **Deploy Directory**: `/deploy/server-yovan/`

## Deployment Architecture

### Environment Files
```
deploy/server-yovan/
├── account.env
├── gateway.env
├── management.env
├── notification.env
├── operational.env          ← Contains SUBSCRIPTION_SERVICE_URL
├── subscription.env
├── tickets.env
├── users.env
├── routes.json
└── up.sh                    ← Deployment script
```

### Docker Compose Setup
Uses TWO compose files:
- `docker-compose.yml` - Base services (postgres, redis, networks)
- `docker-compose.server.yml` - Server-specific overrides

## Deployment Workflow

### Option 1: Deploy from GitHub (server-yovan branch)
```bash
cd ~/propensi

# Pull latest from GitHub server-yovan
git fetch github server-yovan
git checkout github/server-yovan

# Or merge into your local main
git fetch github
git merge github/server-yovan

# Deploy dengan up.sh script
./deploy/server-yovan/up.sh
```

### Option 2: Deploy from GitLab + Update Config
```bash
cd ~/propensi

# Pull dari GitLab
git pull origin main

# Update deployment config dari GitHub server-yovan
git show github/server-yovan:deploy/server-yovan/operational.env > /tmp/operational.env
# Review and merge manually if needed

# Deploy
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

### Option 3: Manual Deployment (Current workflow)
```bash
cd ~/propensi

# Ensure deployment config is up to date
# deploy/server-yovan/operational.env should have:
#   SUBSCRIPTION_SERVICE_URL=http://subscription:5004

# Force recreate operational service if needed
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --force-recreate operational

# Or full rebuild
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

## Key Configuration Notes

### Critical Environment Variables
- `SUBSCRIPTION_SERVICE_URL=http://subscription:5004` - Must be set in operational.env
- All `.env` files in `deploy/server-yovan/` are loaded by Docker Compose

### Service Network
- Services communicate via Docker network (not localhost:port)
- Internal URLs: `http://servicename:port`
- Example: operational → subscription via `http://subscription:5004`

### Database Connection
- DB_HOST=postgres (not localhost)
- Uses internal Docker network

## Automated Deployment (Optional Setup)

### Using up.sh script
```bash
#!/bin/sh
set -e
cd "$(dirname "$0")/../.."
git pull origin server-yovan
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

Modify for your needs:
- Change `origin server-yovan` to your deployment branch
- Add error handling/notifications

### Using Cron Job
```bash
0 2 * * * cd ~/propensi && git pull origin main && docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build >> /var/log/deployment.log 2>&1
```

### Using GitHub Actions (For CD Pipeline)
Create `.github/workflows/deploy.yml` to auto-deploy on push to server-yovan branch.

## Troubleshooting

### Container not connecting
```bash
# Check network
docker network ls
docker network inspect propensi

# Check service connectivity
docker exec thinknalyze-operational ping subscription
```

### Port conflicts
```bash
# Show listening ports
netstat -tuln | grep LISTEN
docker ps
```

### View logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f operational
docker compose logs -f subscription
```

### Restart services
```bash
# Single service
docker compose restart operational

# All services
docker compose restart

# Rebuild and restart
docker compose up -d --build
```

## Quick Reference Commands

```bash
# Status check
docker ps | grep thinknalyze

# Start deployment
cd ~/propensi && ./deploy/server-yovan/up.sh

# Pull updates
git fetch origin && git pull origin main

# View deployment config
cat deploy/server-yovan/operational.env

# Check service health
curl http://localhost:5005/api/operational/stats
```

## File Sync Between Repos

### To update deployment config from GitHub to local:
```bash
git show github/server-yovan:deploy/server-yovan/operational.env > deploy/server-yovan/operational.env.new
# Review differences
diff deploy/server-yovan/operational.env deploy/server-yovan/operational.env.new
# Apply if needed
```

### To push deployment changes to GitHub:
```bash
git add deploy/server-yovan/
git commit -m "Update deployment configuration"
git push github server-yovan
```
