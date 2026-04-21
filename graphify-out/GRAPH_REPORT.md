# Graph Report - gateway  (2026-04-21)

## Corpus Check
- Corpus is ~3,773 words - fits in a single context window. You may not need a graph.

## Summary
- 64 nodes · 96 edges · 10 communities detected
- Extraction: 78% EXTRACTED · 22% INFERRED · 0% AMBIGUOUS · INFERRED: 21 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Community Hubs (Navigation)
- [[_COMMUNITY_Frontend & Auth System|Frontend & Auth System]]
- [[_COMMUNITY_Config Management|Config Management]]
- [[_COMMUNITY_Microservice Routing|Microservice Routing]]
- [[_COMMUNITY_Route Handlers|Route Handlers]]
- [[_COMMUNITY_Auth Middleware|Auth Middleware]]
- [[_COMMUNITY_Proxy Management|Proxy Management]]
- [[_COMMUNITY_Entry Point & RBAC|Entry Point & RBAC]]
- [[_COMMUNITY_Docker Deployment|Docker Deployment]]
- [[_COMMUNITY_Environment Config|Environment Config]]
- [[_COMMUNITY_Redis Session Store|Redis Session Store]]

## God Nodes (most connected - your core abstractions)
1. `main()` - 10 edges
2. `Gateway Main Entry Point` - 10 edges
3. `Authentication Middleware` - 8 edges
4. `Reverse Proxy Pool Manager` - 6 edges
5. `withRoleAuth()` - 5 edges
6. `createProxyHandler()` - 5 edges
7. `SetupRoutes()` - 5 edges
8. `Environment Configuration System` - 5 edges
9. `withRolesAuth()` - 4 edges
10. `GetUserFromToken()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `main()` --calls--> `InitRedis()`  [INFERRED]
  gateway\main.go → gateway\auth\redis_init.go
- `normalizeTargetForEnv()` --calls--> `getEnv()`  [INFERRED]
  gateway\main.go → gateway\core\routes\routes.go
- `withRoleAuth()` --calls--> `CheckRoleAccess()`  [INFERRED]
  gateway\main.go → gateway\auth\auth.go
- `main()` --calls--> `Env()`  [INFERRED]
  gateway\main.go → gateway\system\env.go
- `Env()` --calls--> `getEnv()`  [INFERRED]
  gateway\system\env.go → gateway\core\routes\routes.go

## Hyperedges (group relationships)
- **Authentication & Authorization Flow** — auth_middleware, token_validation, redis_session, role_authorization [EXTRACTED 0.95]
- **Microservice Request Routing** — proxy_pool, users_service, notification_service, subscription_service, operational_service [EXTRACTED 0.95]
- **Frontend Application Serving** — static_assets, account_routes, client_routes, ops_routes, compliance_routes [EXTRACTED 0.95]
- **Docker Multi-Stage Deployment Pipeline** — docker_builder, docker_runtime, docker_config_copy, gateway_port_2000 [EXTRACTED 0.90]

## Communities

### Community 0 - "Frontend & Auth System"
Cohesion: 0.2
Nodes (14): Account Pages (Login/Register), Authentication Middleware, Client Dashboard Pages, Compliance Dashboard Pages, Environment Configuration System, .env File Fallback, Gateway Main Entry Point, Operations Dashboard Pages (+6 more)

### Community 1 - "Config Management"
Cohesion: 0.29
Nodes (6): Config, responseWriter, RouteConfig, loadConfig(), normalizeTargetForEnv(), serveFrontendPage()

### Community 2 - "Microservice Routing"
Cohesion: 0.25
Nodes (8): Auth API Routes → Users Service, Notification API Routes → Notification Service, Notification Service (Port 5003), Operational Service (Port 5005), Reverse Proxy Pool Manager, Subscription API Routes → Subscription Service, Subscription Service (Port 5004), Users Service (Port 5002/2006)

### Community 3 - "Route Handlers"
Cohesion: 0.57
Nodes (5): getEnv(), reverseProxy(), reverseProxyRewrite(), reverseProxyRewriteWild(), SetupRoutes()

### Community 4 - "Auth Middleware"
Cohesion: 0.53
Nodes (5): AuthMiddleware(), CheckRoleAccess(), isPublicPath(), isRoleAllowed(), TokenUser

### Community 5 - "Proxy Management"
Cohesion: 0.4
Nodes (4): ProxyPool, createProxyHandler(), withCORS(), withLogging()

### Community 6 - "Entry Point & RBAC"
Cohesion: 0.6
Nodes (5): GetUserFromToken(), main(), redirectByRole(), withRoleAuth(), withRolesAuth()

### Community 7 - "Docker Deployment"
Cohesion: 0.4
Nodes (5): Go Builder Stage (Alpine), Configuration File Deployment, Multi-Stage Docker Build, Runtime Stage (Alpine), Gateway Listening Port (2000)

### Community 8 - "Environment Config"
Cohesion: 0.5
Nodes (2): Env(), init()

### Community 9 - "Redis Session Store"
Cohesion: 1.0
Nodes (1): InitRedis()

## Knowledge Gaps
- **13 isolated node(s):** `RouteConfig`, `Config`, `TokenUser`, `Operational Service (Port 5005)`, `Static Assets Serving` (+8 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **Thin community `Redis Session Store`** (2 nodes): `redis_init.go`, `InitRedis()`
  Too small to be a meaningful cluster - may be noise or needs more connections extracted.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Gateway Main Entry Point` connect `Frontend & Auth System` to `Microservice Routing`, `Docker Deployment`?**
  _High betweenness centrality (0.117) - this node is a cross-community bridge._
- **Why does `main()` connect `Entry Point & RBAC` to `Environment Config`, `Config Management`, `Proxy Management`, `Redis Session Store`?**
  _High betweenness centrality (0.101) - this node is a cross-community bridge._
- **Why does `Reverse Proxy Pool Manager` connect `Microservice Routing` to `Frontend & Auth System`?**
  _High betweenness centrality (0.077) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `main()` (e.g. with `InitRedis()` and `Env()`) actually correct?**
  _`main()` has 3 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `Gateway Main Entry Point` (e.g. with `Authentication Middleware` and `Multi-Stage Docker Build`) actually correct?**
  _`Gateway Main Entry Point` has 2 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `Authentication Middleware` (e.g. with `Gateway Main Entry Point` and `Account Pages (Login/Register)`) actually correct?**
  _`Authentication Middleware` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 2 inferred relationships involving `withRoleAuth()` (e.g. with `GetUserFromToken()` and `CheckRoleAccess()`) actually correct?**
  _`withRoleAuth()` has 2 INFERRED edges - model-reasoned connections that need verification._