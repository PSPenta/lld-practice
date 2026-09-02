# API Gateway — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Library in-process or standalone HTTP service?
2. Auth: JWT, API key, or mTLS?
3. Rate limit per tenant, user, or route?
4. Routing static config or admin API?
5. TLS termination here or at load balancer?
6. Request/response transform or pass-through only?

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

---

## Step 2 — Entities & classes

`Gateway`, `Route`, `Upstream`, `Filter`/`Middleware` chain, `RateLimiter`, `AuthValidator`

---

## Step 3 — Flows

Request in → match route → run filter chain (auth → rate limit → log) → proxy to upstream → map response

---

## Step 4 — APIs

Admin: `POST /routes` · Data: `ANY /{prefix}/**` proxied · Health: `GET /health`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Timeout upstream; circuit open → 503; idempotent retries only on safe methods; per-tenant rate limit keys

---

## Step 6 — Evolve

New filter → chain of responsibility; new upstream → register route (OCP); Redis for distributed rate limit

---

