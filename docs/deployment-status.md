# Deployment Status

## System Overview
ThinkNalyze is deployed as Go-based microservices on Railway with Gateway as the public entrypoint and internal service-to-service routing over Railway private DNS.

## Current Stable State
- Gateway service: running
- Users service: running
- Account service: running
- Login endpoint verified: `POST /account/login/auth` returns HTTP 200
- Protected endpoint verified: `GET /api/admin/users` returns HTTP 200
- Redis integration verified for auth/session flow in both Gateway and Users
- Internal routing verified with Railway private DNS (`users.railway.internal:8080`)

## Deployed Services
- Gateway
- Users
- Account
- Redis
- Postgres

## Auth Flow (Gateway <-> Users <-> Redis)
1. Client sends login request to Gateway: `POST /account/login/auth`
2. Gateway proxies request to Users service (`users.railway.internal:8080`)
3. Users validates credentials and writes auth session/token reference into Redis
4. Gateway reads auth token/session from Redis for protected route authorization
5. Protected requests are allowed only when Redis-backed auth lookup succeeds

Redis is part of the auth critical path and is required for protected routes.

## Verified Endpoints
- `POST /account/login/auth` -> `200 OK`
- `GET /api/admin/users` -> `200 OK`

## Operational Notes
- Railway deployment runtime uses `PORT=8080` for internal service process binding in this environment.
- Keep service communication internal via `*.railway.internal` hostnames.
- Avoid routing/auth changes unless there is a planned migration with validation.
