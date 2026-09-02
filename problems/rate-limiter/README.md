# Rate limiter — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/RateLimiter2/`](../../JavaScript/RateLimiter2/) | **Strategy** pattern — preferred interview demo |
| **Go** | [`Go/RateLimiter2-go/`](../../Go/RateLimiter2-go/) | interfaces + mutexes |
| Teaching v1 | [`JavaScript/Ratelimiter/`](../../JavaScript/Ratelimiter/) | standalone algorithms, no Strategy |

Distributed variant → [distributed-rate-limiter](../distributed-rate-limiter/README.md)

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Limit per user, IP, or API key?
2. Algorithm preference?
3. Burst allowed?
4. Single node or distributed?
5. Return boolean or throw / HTTP 429?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | HTTP middleware or API callers, `RateLimiter` |
| **Use cases (v1)** | 1. `Allow(key)` before request 2. Reject when exceeded |
| **In scope** | Strategy pattern, per-key state, swappable algorithms |
| **Out of scope** | Redis cluster (see distributed walkthrough) |
| **Assumptions** | Token bucket v1; in-memory; per-IP |

### Confirm understanding
> "Each request checks `Allow(key)`; over limit returns false or HTTP 429."

---

## Step 2 — Entities & classes

```text
RateLimiterStrategy { isAllowed(key): boolean }

RateLimiter
  - strategy: RateLimiterStrategy   // composition, not inheritance

TokenBucketStrategy | FixedWindowCounterStrategy | SlidingWindowLogStrategy | ...
  - per-key buckets / counters
  - mu: Mutex (shared state)
```

### Strategies (pick one for v1; know others for trade-offs)

| Strategy | Idea | Pros | Cons |
|----------|------|------|------|
| Fixed window | Count in current minute | Simple | Burst at window edge |
| Sliding window log | Store timestamps | Accurate | Memory heavy |
| Token bucket | Tokens refill over time | Smooth + burst | Slightly more logic |
| Leaky bucket | Steady outflow | Smooth egress | Less burst friendly |

**Pattern:** **Strategy** + **composition** — `RateLimiter` does not extend `TokenBucket`.

---

## Step 3 — Flows

**Allow(key):**
1. Delegate to `strategy.isAllowed(key)`
2. Strategy updates per-key state under lock
3. Return true (proceed) or false (reject → 429)

**HTTP middleware:** wrap handler → if `!Allow(clientIP)` return 429 + `Retry-After`

---

## Step 4 — APIs

- Library: `allow(key: string): boolean`
- HTTP: middleware returns **429** when false

---

## Step 5 — Deepen (concurrency, failure, idempotency)

- Mutex on per-key maps; avoid holding lock across I/O
- Sliding window: bound memory (trim old timestamps)
- Distributed: race on Redis counters — Lua or accept approximate limit

---

## Step 6 — Evolve

- [distributed-rate-limiter](../distributed-rate-limiter/README.md) — Redis-backed counters
- New algorithm → new strategy class; `RateLimiter.js` unchanged (**OCP**)
