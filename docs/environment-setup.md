# Environment Setup

This file contains only required environment variables for stable auth and database connectivity.

## Required Variables

### Required for Users and Gateway
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

## Service Scope
- Set `redis*` variables on both Users and Gateway services.
- Set `SESSION_KEY`, `SESSION_LIFETIME`, `JWT_EXPIRED` where auth/session logic is evaluated.
- Set `postgres` and `read_db_*` where service reads from Postgres.

## Critical Warnings
- Do not use `localhost` for inter-service communication on Railway.
- Do not use `users-service` naming for Users route target.
- Use internal DNS: `users.railway.internal:8080`
- Railway runtime binds app process on `PORT=8080` in this setup.
- Redis is required for auth and protected route authorization.
