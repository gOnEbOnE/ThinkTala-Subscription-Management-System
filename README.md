# ThinkNalyze Deployment Handoff

This repository contains ThinkNalyze microservices and deployment assets for Railway.

## System Overview
ThinkNalyze runs as Go microservices behind a Gateway service. Authentication and protected route authorization depend on Redis-backed session/token lookup between Gateway and Users.

Current stable production state:
- Gateway running
- Users running
- Account running
- Login endpoint verified (`POST /account/login/auth` -> 200)
- Protected endpoint verified (`GET /api/admin/users` -> 200)
- Redis integrated for auth/session in Gateway and Users
- Internal routing uses Railway private DNS (`users.railway.internal:8080`)

## Architecture Summary
- Gateway: public entrypoint and reverse proxy
- Users: auth/session logic and protected user/admin APIs
- Account: account-related flow endpoints
- Redis: required auth/session state backend
- Postgres: persistent data store

Auth path (high-level):
1. Client logs in through Gateway
2. Gateway proxies login to Users
3. Users validates credentials and writes session/token reference to Redis
4. Gateway reads Redis-backed auth context for protected routes

## Documentation
- [docs/deployment-status.md](docs/deployment-status.md)
- [docs/environment-setup.md](docs/environment-setup.md)
- [docs/deployment-guide.md](docs/deployment-guide.md)

## Quick Start (Railway)
1. Create/confirm services: Gateway, Users, Account, Redis, Postgres.
2. Set required environment variables from [docs/environment-setup.md](docs/environment-setup.md).
3. Deploy each service with `railway up`.
4. Validate:
   - `POST /account/login/auth` returns 200
   - `GET /api/admin/users` returns 200 (with login cookie/session)

## Critical Warnings
- Do not use `localhost` for service-to-service communication on Railway.
- Do not use `users-service` as target host naming.
- Use `users.railway.internal:8080`.
- Railway runtime uses `PORT=8080` in this setup.
- Redis is required for auth and protected routes.