# LFU cache — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌ · **Code:** `JavaScript/LRU/ (compare LRU ✅)`

## Step 1 — Clarify

### Questions (ask 6–8)
1. Capacity?
2. TTL?
3. Thread-safe?
4. Evict least frequently used — tie-break LRU?
5. get/put O(1) required?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Cache client callers |
| **Use cases (v1)** | 1. Get (promote freq) 2. Put (evict min freq if full) |
| **In scope** | LFU eviction, O(1) target |
| **Out of scope** | Distributed cache |
| **Assumptions** | In-memory; fixed capacity; mutex |

### Confirm understanding
> "Cache evicts items used least often when full, not just oldest."

---

## Step 2 — Entities & classes

`LFUCache`, `Node`, freq buckets map, `minFreq` pointer

---

## Step 3 — Flows

Get → increment freq, move bucket → Put → evict from minFreq bucket if full

---

## Step 4 — APIs

`Get(key)`, `Put(key, val)`, `Delete(key)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Lock entire structure or striping; O(1) goal for get/put

---

## Step 6 — Evolve

Compare with LRU walkthrough; hybrid LRFU policy

---

