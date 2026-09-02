# Cache Client — LLD design walkthrough

> **Paper design** (no dedicated code folder). Overlaps: `JavaScript/LRU/`, `JavaScript/Redis/`, `Go/LRU-go/`, `Go/Redis-go/`.  
> **Method:** [../../README.md §5](../../README.md) · **Redis eviction policies:** [../../JavaScript/Redis/README.md](../../JavaScript/Redis/README.md)

> **Very common at product companies** (including helpdesk/SaaS interviews): *“Design a Cache Client that caches frequent queries.”*

---

## Step 1 — Clarifying questions (ask 6–8)

1. Library or HTTP service?
2. In-memory only, or JavaScript/Redis/shared cache?
3. Single machine or many servers?
4. TTL (expire after time)?
5. Max capacity? Eviction policy (LRU)?
6. Thread-safe? (multiple goroutines)
7. On miss — return error, or load from DB (`GetOrLoad`)?
8. Cache stampede: 100 concurrent misses on same key — all hit DB?
9. Metrics needed (hit/miss rate)?

**State assumptions if they say “your call”:** in-memory v1, TTL + LRU, thread-safe, `GetOrLoad` with loader.

---

## Step 2 — Entities & classes

```text
CacheEntry: key, value, expiresAt

interface Cache {
  Get(key) (value, found)
  Set(key, value, ttl)
  Delete(key)
}

CacheClient implements Cache
  - store: map[string]*CacheEntry
  - order: doubly linked list (LRU)
  - capacity, defaultTTL
  - mu: Mutex

  GetOrLoad(key, loader func() (value, error))
```

**Code:** [JavaScript/LRU/](../../JavaScript/LRU/) · [JavaScript/Redis/](../../JavaScript/Redis/) · [Go/LRU-go/](../../Go/LRU-go/) · [Go/Redis-go/](../../Go/Redis-go/)

---

## Step 3 — Flows

**Get:** lock → miss/expired → return miss → else move to MRU → unlock → return value  

**Set:** lock → update or insert → evict LRU tail if over capacity → unlock  

**GetOrLoad:** on miss, call loader — use **singleflight** so 50 concurrent misses → one DB call.

---

## Step 4 — API (if HTTP service)

```http
GET    /v1/cache/{key}
PUT    /v1/cache/{key}   body: { "value": "...", "ttl_sec": 60 }
DELETE /v1/cache/{key}
```

Often a **library** — ask first.

---

## Step 5 — Trade-offs / evolve

| Question | Answer |
|----------|--------|
| Traffic ↑ | Redis L2; local L1 |
| Memory ↑ | Capacity cap + TTL |
| Redis down | Fallback to DB (degraded) |
| Extend | `Cache` interface → Memory / Redis (Adapter) |
| Eviction | Strategy: LRU, LFU — see [Redis README](../../JavaScript/Redis/README.md) |
| Monitor | hit_rate, miss_rate, eviction_count, p99 |

**Patterns:** Strategy/Adapter for backends; SRP — cache does not know SQL (loader injected).

**Practice:** Explain in 20 minutes out loud without notes.
