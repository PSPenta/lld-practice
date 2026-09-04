# LRU cache — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Capacity fixed?
2. TTL needed?
3. Thread-safe?
4. Library or HTTP service?
5. `get`/`put` must be O(1)?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Cache callers |
| **Use cases (v1)** | 1. Get (refresh recency) 2. Put (evict LRU if full) |
| **In scope** | Hash map + DLL, capacity eviction |
| **Out of scope** | Distributed cache, GetOrLoad (see cache-client) |
| **Assumptions** | In-memory; mutex; fixed capacity |

### Confirm understanding
> "Get/Put are O(1); least recently used entry evicted when over capacity."

---

## Step 2 — Entities & classes

**Core idea:** **Map** for O(1) lookup key → node · **Doubly linked list** for recency (head = MRU, tail = LRU)

```text
Node { key, value, prev, next }

LRUCache {
  capacity, map, head, tail, mu
  get(key) → value | miss
  put(key, value)
  // private: moveToFront, evictTail, removeNode
}
```

---

## Step 3 — Flows

**Get hit:** find node → move to front → return value  
**Get miss:** return not found  

**Put existing key:** update value → move to front  

**Put new key:** insert at front → if size > capacity → remove tail (LRU)

---

## Step 4 — APIs

- `get(key)`, `put(key, value)`, optional `delete(key)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

- Wrap `get`/`put` with **mutex** if multi-threaded
- Missing key on get — don't throw; return miss sentinel

---

## Step 6 — Evolve

- TTL on node + lazy expiry on get
- `Cache` interface → memory LRU / Redis adapter
- Full client: [cache-client](../cache-client/README.md)


---

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/LRU/`](../../JavaScript/LRU/) | map + doubly linked list |
| **Go** | [`Go/LRU-go/`](../../Go/LRU-go/) | |
| Broader client | [cache-client](../cache-client/README.md) | TTL, GetOrLoad, stampede |
| Redis-style TTL | [`JavaScript/Redis/`](../../JavaScript/Redis/) · [`Go/Redis-go/`](../../Go/Redis-go/) | eviction policies |

## Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `LRU/index.js` | `LRUCache` — `Map` + doubly linked list; `get` / `put`; private move/evict helpers |
| `Go/LRU-go/` | Same structure with mutex if concurrent |

**Read order:** `get` / `put` → `moveToFront` → `evictTail` when over capacity.

