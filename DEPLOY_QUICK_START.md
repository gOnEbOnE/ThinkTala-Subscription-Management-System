# QUICK START - Deploy di Ubuntu Server

## TL;DR - Langkah Cepat

### Pertama kali di server
```bash
cd ~/propensi

# Add GitHub remote
git remote add github https://github.com/gOnEbOnE/ThinkTala-Subscription-Management-System.git

# Fetch deployment config
git fetch github server-yovan

# Setup permission
chmod +x ~/propensi/deploy/server-yovan/*.sh
```

### Deploy setiap kali ada update
```bash
cd ~/propensi

# Cara 1: Menggunakan script (Recommended)
./deploy/server-yovan/deploy.sh

# Cara 2: Manual
git fetch github server-yovan
git merge github/server-yovan
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

### Setup automated deployment (optional)
```bash
# Edit crontab
crontab -e

# Add one of these:
# Deploy every day at 2 AM
0 2 * * * cd ~/propensi && ./deploy/server-yovan/deploy.sh >> ~/deployment.log 2>&1

# Deploy every hour
0 * * * * cd ~/propensi && ./deploy/server-yovan/deploy.sh >> ~/deployment.log 2>&1
```

### Monitor services
```bash
# Check status
docker ps | grep thinknalyze

# View logs
docker compose logs -f operational

# Test API
curl http://localhost:5005/api/operational/stats
```

---

## File Structure

```
~/propensi/
├── DEPLOYMENT_GUIDE.md          ← Read this for detailed info
├── SETUP_DEPLOYMENT_SERVER.md   ← Read this for step-by-step setup
├── docker-compose.yml           ← Base configuration
├── docker-compose.server.yml    ← Server-specific overrides
├── deploy/server-yovan/
│   ├── up.sh                    ← Simple deploy script
│   ├── deploy.sh                ← Advanced deploy script with logging
│   ├── operational.env          ← Operational service config
│   ├── subscription.env         ← Subscription service config
│   ├── account.env              ← Account service config
│   ├── gateway.env              ← Gateway service config
│   ├── notification.env         ← Notification service config
│   ├── management.env           ← Management service config
│   ├── users.env                ← Users service config
│   ├── tickets.env              ← Tickets service config
│   └── routes.json              ← Gateway routes
```

---

## Critical Config

### operational.env
```
SUBSCRIPTION_SERVICE_URL=http://subscription:5004
```
✅ Already configured

### Service Network
- All services communicate via `docker network propensi`
- Internal hostnames: `operational`, `subscription`, `gateway`, etc.
- No localhost - use service names

---

## Troubleshooting Quick Links

| Issue | Command |
|-------|---------|
| Services not connecting | `docker exec thinknalyze-operational ping subscription` |
| Check service logs | `docker compose logs -f operational` |
| Rebuild services | `docker compose up -d --build` |
| Update config | Edit `deploy/server-yovan/*.env` then `docker compose restart` |
| Check network | `docker network inspect propensi` |
| Database status | `docker exec -it thinknalyze-postgres psql -U postgres -d thinknalyze` |

---

## GitHub & GitLab Workflow

### Development (GitLab)
```bash
git pull origin main          # Pull development code
git push origin branch-name   # Push code changes
```

### Server Deployment (GitHub)
```bash
git fetch github server-yovan      # Get deployment config
./deploy/server-yovan/deploy.sh    # Deploy to server
```

### Keep in sync
```bash
git fetch github server-yovan
git merge github/server-yovan     # Merge deployment config to current branch
git push origin current-branch    # Push back to GitLab if needed
```

---

## Next Steps

1. ✅ **Copy deployment setup** - Already done, GitHub remote added
2. ✅ **Config verified** - SUBSCRIPTION_SERVICE_URL is set
3. ✅ **Services running** - Docker compose up and healthy
4. 📋 **TODO**: Setup cron for automated deployment (optional)
5. 📋 **TODO**: Configure monitoring/alerts for service failures

---

## Reference Docs

- [DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md) - Full deployment architecture
- [SETUP_DEPLOYMENT_SERVER.md](SETUP_DEPLOYMENT_SERVER.md) - Detailed setup steps

---

**Need help?** Check:
1. Docker logs: `docker compose logs service-name`
2. Config files: `cat deploy/server-yovan/*.env`
3. Network: `docker network inspect propensi`
4. Documentation: See files above
