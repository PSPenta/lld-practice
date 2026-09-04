# API Gateway — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Library in-process or standalone HTTP service?
2. Auth: JWT, API key, or mTLS?
3. Rate limit per tenant, user, or route?
4. Routing static config or admin API?
5. TLS termination here or at load balancer?
6. Request/response transform or pass-through only?
7. Circuit breaker / timeout per upstream?
8. Multi-region or single region for v1?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | API clients, gateway admin, upstream microservices |
| **Use cases (v1)** | 1. Proxy request to upstream 2. Enforce auth + rate limit 3. Register/update routes (admin) |
| **In scope** | Route table, filter chain (auth → rate limit → log), reverse proxy |
| **Out of scope** | Service mesh, full WAF, GraphQL gateway |
| **Assumptions** | Single region v1; in-memory route table; JWT validation |

### Confirm understanding
> "For v1, clients hit the gateway; it matches a route, runs filters, and proxies to one upstream URL."

## Step 2 — Entities & classes

```text
Route { id, pathPrefix, methods[], upstreamId, filters[] }
Upstream { id, baseURL, timeoutMs, health }

Filter (interface) → AuthFilter | RateLimitFilter | LoggingFilter
  - pre(request) → error?
  - post(response)?

Gateway
  - matchRoute(request) → Route
  - handle(request) → response   // run filter chain then proxy

RateLimiter, AuthValidator, AdminConfigStore
```

**Patterns:** Chain of Responsibility (filters) · Strategy (rate limit) · Reverse proxy facade

## Step 3 — Flows

**Happy path**
1. Request arrives → match path/method against route table  
2. Run filter chain in order (auth → rate limit → log)  
3. Proxy to upstream with timeout  
4. Map status/body back to client  

**Edge cases**
1. No route match → 404; auth fail → 401; rate limited → 429  
2. Upstream timeout / circuit open → 503 (no silent hang)

## Step 4 — APIs

```http
POST /admin/routes          { pathPrefix, methods, upstreamId }
GET  /admin/routes
ANY  /{prefix}/**           # proxied after filters
GET  /health
```

Library-style: `Gateway.Handle(ctx, req) (resp, error)`

## Step 5 — Deepen

- Timeout every upstream call; never block forever on a hung dependency  
- Circuit open → fail fast with 503; retry only on safe/idempotent methods  
- Per-tenant rate-limit keys; concurrent route-table reads vs admin writes need locking or copy-on-write  
- Idempotent admin upserts by route id  
- Propagate correlation/request id for observability

## Step 6 — Evolve

- New filter → implement `Filter` and insert in chain (**OCP**)  
- New upstream → register route without touching proxy core  
- Redis for distributed rate limit; service discovery for dynamic backends  
- Related: [load-balancer](../load-balancer/README.md), [circuit-breaker](../circuit-breaker/README.md), [distributed-rate-limiter](../distributed-rate-limiter/README.md)
