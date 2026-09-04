# LFU cache — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Compare after design with solved LRU: `JavaScript/LRU/` (✅).

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Capacity?
2. TTL?
3. Thread-safe?
4. Evict least frequently used — tie-break LRU within same freq?
5. get/put O(1) required?
6. Null values allowed?
7. Delete / clear APIs?
8. Eviction callback / stats?

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

## Step 2 — Entities & classes

```text
Node { key, value, freq, prev, next }   // in a freq bucket DLL

FreqBucket { freq, head, tail }         // doubly linked list of Nodes
  - addToHead(node) / remove(node) / popTail()

LFUCache {
  capacity, size, minFreq
  keyToNode  map[key]*Node
  freqToList map[freq]*FreqBucket
  - Get(key) (val, ok)
  - Put(key, val)
  - Delete(key)
}
```

**Patterns:** Hash map + doubly linked lists per frequency · same structural idea as LRU, keyed by freq

## Step 3 — Flows

**Happy path — Get**
1. Lookup key → miss returns not found  
2. Hit → remove from current freq list → freq++ → insert into new list  
3. If old freq list empty and freq==minFreq → minFreq++  

**Happy path — Put**
1. Update existing → same as Get + set value  
2. If full → evict tail of `freqToList[minFreq]` → delete key  
3. Insert new node at freq=1; minFreq=1  

**Edge cases**
1. capacity 0 → no-op / error; duplicate Put same key doesn’t grow size  
2. Tie at minFreq → evict LRU within that bucket (list tail)

## Step 4 — APIs

```text
Get(key) (value, ok)
Put(key, value)
Delete(key) bool
Len() / Cap()
```

## Step 5 — Deepen

- Lock entire structure (v1) or shard by key hash for less contention  
- Keep get/put amortized O(1) with map + DLL; avoid scanning all keys  
- Thread-safe minFreq updates when a bucket drains  
- Optional TTL: separate expiry heap or lazy expire on access  
- Put of existing key must bump frequency (policy — state aloud)

## Step 6 — Evolve

- Compare with LRU walkthrough; hybrid LRFU / TinyLFU  
- Soft capacity / weighted entries  
- Related: [rate-limiter](../rate-limiter/README.md) only for “hot key” discussion — cache is local structure
