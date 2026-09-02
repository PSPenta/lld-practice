# LRU Cache — LLD design walkthrough

> **Code:** this folder · **Go port:** [../../Go/LRU-go/](../../Go/LRU-go/)  
> **Method:** [../../README.md §5](../../README.md) · **Cache client (broader):** [../../problems/cache-client/README.md](../../problems/cache-client/README.md)

---

## Clarify

- Capacity?  
- TTL needed?  
- Thread-safe?  
- `get` updates recency? (yes for LRU)

---

## Core idea

- **Map** for O(1) lookup key → node  
- **Doubly linked list** for recency order (head = most recent, tail = least)

---

## Classes

```text
Node { key, value, prev, next }
LRUCache {
  capacity
  map
  head, tail
  Get(key) → value
  Put(key, value)
  // private: moveToFront, evictTail, removeNode
}
```

---

## Flows

**Get hit:** move node to front; return value  
**Put existing:** update value; move to front  
**Put new:** insert front; if over capacity, remove tail  

---

## Concurrency

Wrap `Get`/`Put` with mutex if multi-threaded.

---

## Extend

- TTL field on node + lazy expiry on get  
- `Cache` interface → MemoryLRU / Redis adapter  
