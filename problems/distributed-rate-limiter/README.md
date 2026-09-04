# Distributed rate limiter — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Local algorithms first: `JavaScript/RateLimiter2/` and [rate-limiter](../rate-limiter/README.md).

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Redis central store?
2. Approximate OK?
3. Per-tenant keys?
4. Fail open or closed if Redis down?
5. Sync with local cache?
6. Fixed window, sliding window, or token bucket?
7. Limit per API key / IP / user?
8. Soft vs hard limit (throttle vs reject)?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | API nodes, Redis, RateLimiter client |
| **Use cases (v1)** | 1. Allow(key) across all nodes consistently enough |
| **In scope** | Redis INCR + TTL or Lua script |
| **Out of scope** | Global precise quota billing |
| **Assumptions** | Redis; slight race OK; compare local RateLimiter first |

### Confirm understanding
> "All API instances share counter in Redis so limits apply cluster-wide."

## Step 2 — Entities & classes

```text
RateLimitConfig { limit, window, algorithm }

RedisCounter
  - IncrWithExpire(key, window) → count   // INCR + EXPIRE or Lua

DistributedRateLimiter
  - Allow(key) (allowed bool, remaining int)
  - config, redis, optional localCache

LocalTokenBucket?   // optional first-line shed before Redis
```

**Patterns:** Shared counter · Strategy (window vs bucket) same interface as local limiter · Fail-open/closed policy

## Step 3 — Flows

**Happy path**
1. `Allow(key)` → build Redis key `rl:{tenant}:{windowBucket}`  
2. Atomic INCR (+ set TTL on first hit) via Lua or pipeline  
3. If count ≤ limit → allow; else deny  

**Edge cases**
1. Redis timeout → fail open (allow) or fail closed (deny) — state aloud  
2. Clock skew across nodes: fixed window boundary burst — mention sliding/token bucket evolve

## Step 4 — APIs

```text
Allow(key) (allowed bool, remaining int, resetAt)
// Same surface as local limiter; backend is Redis
```

Middleware: `if !limiter.Allow(tenantId) { return 429 }`

## Step 5 — Deepen

- Race on INCR+EXPIRE → use Lua for atomicity  
- Redis down: choose fail-open vs fail-closed with product impact  
- Key cardinality: don’t put unbounded user ids without TTL hygiene  
- Local cache can approximate but must not undermine global limit badly  
- Idempotent reads of remaining are best-effort under concurrency

## Step 6 — Evolve

- Compare [rate-limiter](../rate-limiter/README.md) local Strategy first  
- Sliding window / GCRA; per-route limits; Redis Cluster  
- Related: [api-gateway](../api-gateway/README.md) filter chain
