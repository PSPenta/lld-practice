# Cache (TTL + eviction) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅  
> Very common at product companies: *“Design a cache client for frequent queries.”*

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Library or HTTP service?
2. In-memory only, or shared Redis?
3. Single machine or many servers?
4. TTL (expire after time)?
5. Max capacity? Eviction policy (LRU)?
6. Thread-safe?
7. On miss — error or `GetOrLoad` from DB?
8. Cache stampede on concurrent misses?
9. Metrics (hit/miss rate)?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Application code or HTTP client |
| **Use cases (v1)** | Get/Set/Delete · TTL expiry · evict at capacity |
| **In scope** | In-memory cache, LRU, optional loader |
| **Out of scope** | Multi-node consistency (Step 6) |
| **Assumptions** | Library; TTL + LRU; thread-safe; `GetOrLoad` |

### Confirm understanding
> "Bounded in-memory cache with LRU eviction and optional loader on miss."

---

## Step 2 — Entities & classes

```text
CacheEntry { key, value, expiresAt }

interface Cache {
  get(key), set(key, value, ttl), delete(key)
}

CacheClient
  - store: map[string]*CacheEntry
  - order: doubly linked list (LRU)
  - capacity, defaultTTL, mu
  getOrLoad(key, loader)
```

---

## Step 3 — Flows

**Get:** lock → miss/expired → return miss → else move to MRU → unlock → return  

**Set:** lock → update/insert → evict LRU tail if over capacity → unlock  

**GetOrLoad:** on miss → **singleflight** so N concurrent misses → one DB call

---

## Step 4 — APIs

Library methods above, or if HTTP:

```http
GET    /v1/cache/{key}
PUT    /v1/cache/{key}   body: { "value", "ttl_sec" }
DELETE /v1/cache/{key}
```

---

## Step 5 — Deepen (concurrency, failure, idempotency)

- Mutex on map + list
- Loader failure → don't cache poison values
- Stampede protection via singleflight / per-key lock

---

## Step 6 — Evolve

| Topic | Answer |
|-------|--------|
| Traffic ↑ | Redis L2; local L1 |
| Redis down | Fallback to DB (degraded) |
| Backends | `Cache` interface → Memory / Redis (**Adapter**) |
| Monitor | hit_rate, miss_rate, eviction_count |

### Redis eviction policies (when interviewer asks “like Redis”)

| Policy | Evicts |
|--------|--------|
| `noeviction` | Nothing — writes fail when full |
| `allkeys-lru` | Any key — LRU (approximate in real Redis) |
| `allkeys-lfu` | Any key — LFU |
| `volatile-lru` | Keys **with TTL** only — LRU among those |
| `volatile-ttl` | Shortest TTL first |

**Interview lines:** pure cache → `allkeys-lru`; mix permanent + cache keys → `volatile-lru` + TTL on cache only. **Eviction** (memory full) ≠ **expiration** (TTL).

**Patterns:** Strategy/Adapter for backends; SRP — cache does not know SQL (loader injected).


---

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript LRU** | [`JavaScript/LRU/`](../../JavaScript/LRU/) | exact LRU (map + DLL) |
| **JavaScript Redis-style** | [`JavaScript/Redis/`](../../JavaScript/Redis/) | TTL + policy-based eviction |
| **Go** | [`Go/LRU-go/`](../../Go/LRU-go/) · [`Go/Redis-go/`](../../Go/Redis-go/) | |
| Pure LRU only | [lru-cache](../lru-cache/README.md) | subset of this problem |

## Codebase map (how the code is organized)

| Path | Responsibility |
|------|----------------|
| `JavaScript/LRU/index.js` | Exact LRU — O(1) get/put via map + doubly linked list |
| `JavaScript/Redis/index.js` | TTL expiry + `maxmemory`-style eviction policy switch |
| `Go/LRU-go/`, `Go/Redis-go/` | Ports of the above |

**Interview tip:** start design from LRU; add TTL + eviction policies when asked “like Redis.”

