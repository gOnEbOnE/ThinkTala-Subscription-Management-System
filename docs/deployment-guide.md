# Deployment Guide

Step-by-step deployment guide for reproducible Railway setup.

## 1. Create Railway Services
Create or confirm services in one Railway project:
- Gateway
- Users
- Account
- Redis
- Postgres

## 2. Add Redis Service
If Redis does not exist, add Railway Redis service first.

## 3. Set Environment Variables
Set required variables exactly as below.

```env
redis=true
redis_host=redis.railway.internal
redis_port=6379
redis_pass=<Railway Redis password>
redis_db=0

SESSION_KEY=<required>
SESSION_LIFETIME=86400
JWT_EXPIRED=3600

postgres=true
read_db_host=<...>
read_db_user=<...>
read_db_pass=<...>
```

Apply on:
- Users: all required variables
- Gateway: all `redis*` variables and auth-relevant session vars

## 4. Deploy Service (`railway up`)
From each service root as needed:

```bash
railway up
```

Or via linked project/service context:

```bash
railway service <service-name>
railway up
```

## 5. Validate Endpoints
Use a gateway public URL and verify both login and protected route.

### 5.1 Set gateway URL
```bash
export GATEWAY_BASE_URL="https://<gateway-domain>"
```

### 5.2 Login test (expect 200)
```bash
curl -i -c cookies.txt \
  -H "Content-Type: application/json" \
  -X POST "$GATEWAY_BASE_URL/account/login/auth" \
  -d '{"email":"superadmin@thinktala.com","password":"Super123"}'
```

Expected:
- HTTP status: `200`
- Response body contains successful auth message
- Cookie jar receives auth cookies (`za_session`, `token`)

### 5.3 Protected route test (expect 200)
```bash
curl -i -b cookies.txt \
  "$GATEWAY_BASE_URL/api/admin/users"
```

Expected:
- HTTP status: `200`
- Protected data returned

## Verification Checklist
- `POST /account/login/auth` returns `200`
- `GET /api/admin/users` returns `200`
- Users logs show Redis connected
- Gateway auth path can read Redis-backed session/token

## Critical Warnings
- Do not use `localhost` for service-to-service communication.
- Do not use `users-service` naming.
- Use `users.railway.internal:8080`.
- Railway runtime uses `PORT=8080` in this deployment.
- Redis is required for auth and protected routes.
