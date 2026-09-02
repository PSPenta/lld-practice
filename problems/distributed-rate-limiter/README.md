# Distributed rate limiter — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌ · **Code:** `JavaScript/RateLimiter2/ (local algo)`

## Step 1 — Clarify

### Questions (ask 6–8)
1. Redis central store?
2. Approximate OK?
3. Per-tenant keys?
4. Fail open or closed if Redis down?
5. Sync with local cache?

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

---

## Step 2 — Entities & classes

`DistributedRateLimiter`, `RedisCounter`, local cache optional

---

## Step 3 — Flows

Allow → INCR Redis key with TTL → compare to limit → return

---

## Step 4 — APIs

Same `Allow(key)` as local; backend is Redis

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Race tolerance or Lua script atomicity; Redis down → fail open or closed?

---

## Step 6 — Evolve

Compare [rate-limiter](../rate-limiter/README.md) local Strategy first

---

