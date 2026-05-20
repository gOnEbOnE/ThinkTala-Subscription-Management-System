# SETUP DEPLOYMENT DI UBUNTU SERVER

## Current Status
✅ Docker Compose sudah running dengan operational & subscription services  
✅ GitHub remote sudah di-add  
✅ server-yovan branch sudah di-push ke GitHub  
✅ Deployment config sudah di-setup  

---

## Langkah 1: Setup di Ubuntu Server (One-time)

### 1.1 Clone/Pull deployment config dari GitHub
```bash
cd ~/propensi

# Add GitHub remote (jika belum)
git remote add github https://github.com/gOnEbOnE/ThinkTala-Subscription-Management-System.git

# Fetch server-yovan branch
git fetch github server-yovan

# Lihat deployment config
cat deploy/server-yovan/operational.env
```

### 1.2 Setup permission untuk deployment script
```bash
chmod +x ~/propensi/deploy/server-yovan/up.sh
chmod +x ~/propensi/deploy/server-yovan/deploy.sh
```

### 1.3 Verify Docker setup
```bash
# Check services running
docker ps | grep thinknalyze

# Check network
docker network inspect propensi
```

---

## Langkah 2: Deploy di Ubuntu Server

### Option A: Menggunakan deployment script (Recommended)
```bash
cd ~/propensi
./deploy/server-yovan/deploy.sh
```

Script ini akan:
1. ✅ Pull dari GitHub server-yovan branch
2. ✅ Verify deployment configuration
3. ✅ Build Docker images
4. ✅ Start semua services
5. ✅ Show logs untuk verification

### Option B: Manual deployment
```bash
cd ~/propensi

# Pull dari GitLab main (development)
git pull origin main

# ATAU pull dari GitHub server-yovan (deployment)
git fetch github server-yovan
git checkout github/server-yovan

# Rebuild dan restart services
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

### Option C: Update tertentu services
```bash
cd ~/propensi

# Hanya update operational service
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build operational

# Atau subscription service
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build subscription
```

---

## Langkah 3: Automated Deployment (Cron Job)

### 3.1 Setup daily deployment
```bash
# Edit crontab
crontab -e

# Add line (deploy setiap hari jam 2 pagi):
0 2 * * * cd ~/propensi && ./deploy/server-yovan/deploy.sh >> ~/deployment.log 2>&1

# Atau setiap jam (every hour):
0 * * * * cd ~/propensi && ./deploy/server-yovan/deploy.sh >> ~/deployment.log 2>&1
```

### 3.2 Check cron logs
```bash
# View deployment logs
tail -f ~/deployment.log

# Or recent only
tail -20 ~/deployment.log
```

---

## Langkah 4: Workflow Integration

### For Development Team (GitLab)
```bash
# Development branch di GitLab
git pull origin main          # Pull latest development
git push origin branch-name   # Push development code
```

### For Server Deployment (GitHub)
```bash
# Server deployment dari GitHub server-yovan
git fetch github server-yovan
git merge github/server-yovan  # Or checkout

# Deploy ke server
./deploy/server-yovan/deploy.sh
```

### Keeping Both Repos in Sync
```bash
# Pull dari GitHub
git fetch github server-yovan

# Merge ke current development branch
git merge github/server-yovan

# Push ke GitLab
git push origin current-branch
```

---

## Langkah 5: Verification & Monitoring

### Check services status
```bash
# All containers
docker ps | grep thinknalyze

# Check specific services
docker ps | grep -E "operational|subscription|gateway"
```

### View service logs
```bash
# All services
docker compose logs -f

# Operational service
docker compose logs -f operational

# Subscription service  
docker compose logs -f subscription

# Gateway service
docker compose logs -f gateway

# Last 20 lines
docker compose logs --tail=20 operational
```

### Health checks
```bash
# Check operational API
curl http://localhost:5005/api/operational/stats

# Check if services can connect
docker exec thinknalyze-operational ping subscription
docker exec thinknalyze-operational curl http://subscription:5004/health 2>/dev/null || true
```

### Database checks
```bash
# Connect ke PostgreSQL
docker exec -it thinknalyze-postgres psql -U postgres -d thinknalyze -c "SELECT version();"

# Check tables
docker exec -it thinknalyze-postgres psql -U postgres -d thinknalyze -c "\dt"

# Check Redis
docker exec -it thinknalyze-redis redis-cli PING
```

---

## Troubleshooting

### Container won't start
```bash
# Check logs
docker compose logs operational

# Force rebuild
docker compose up -d --build operational

# Check image
docker images | grep operational
```

### Connection refused errors
```bash
# Verify network
docker network inspect propensi

# Ping from one container to another
docker exec thinknalyze-operational ping subscription

# Check port availability
netstat -tuln | grep 5004
```

### Configuration issues
```bash
# Verify env file
cat deploy/server-yovan/operational.env

# Update config
echo "SUBSCRIPTION_SERVICE_URL=http://subscription:5004" >> deploy/server-yovan/operational.env

# Force recreate with new config
docker compose up -d --force-recreate operational
```

### Clean rebuild
```bash
# Stop all
docker compose down

# Remove images (optional)
docker rmi -f $(docker images propensi-* -q)

# Rebuild from scratch
docker compose -f docker-compose.yml -f docker-compose.server.yml up -d --build
```

---

## Important Notes

### Service Discovery
- Services komunikasi via Docker network, bukan localhost
- `operational` → `subscription` via `http://subscription:5004`
- `operational` → `users` via `http://users:8080`

### Environment Variables
- Disimpan di `deploy/server-yovan/*.env`
- Diload otomatis oleh Docker Compose
- Perubahan butuh restart service

### Database
- PostgreSQL: `postgres:5432` (internal)
- Redis: `redis:6379` (internal)
- Backup: Check docker volumes

### Deployment Branch Strategy
```
GitLab (origin)        GitHub (github)
  main                   server-yovan ← Deployment branch
   ↓                          ↓
Development            Production Deploy
```

---

## Quick Commands Reference

```bash
# Deploy
cd ~/propensi && ./deploy/server-yovan/deploy.sh

# Check status
docker ps | grep thinknalyze

# View logs
docker compose logs -f operational

# Restart service
docker compose restart operational

# Pull latest
git fetch github server-yovan && git merge github/server-yovan

# Update single config
echo "NEW_VAR=value" >> deploy/server-yovan/operational.env
docker compose restart operational
```

---

## Support & Logs

All deployment logs stored in:
- Container logs: `docker compose logs [service-name]`
- System logs: `/var/log/` (if configured)
- Deployment script logs: `~/deployment_*.log`

For issues, check:
1. `docker compose logs` output
2. Container environment: `docker inspect <container> | grep Env`
3. Network: `docker network inspect propensi`
4. File permissions: `ls -la ~/propensi/deploy/server-yovan/`
