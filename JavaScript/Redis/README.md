# Redis-style cache & eviction policies

> **Code:** this folder · **Go port:** [../../Go/Redis-go/](../../Go/Redis-go/) · **Exact LRU (linked list):** [../LRU/README.md](../LRU/README.md)  
> **Cache client design:** [../../problems/cache-client/README.md](../../problems/cache-client/README.md)

---

## Redis eviction policies (cache / Redis LLD)

When Redis hits `maxmemory`, it evicts keys per `maxmemory-policy`:

| Policy | Evicts |
|--------|--------|
| `noeviction` | Nothing — writes fail when full |
| `allkeys-lru` | Any key — least recently used (approximate) |
| `allkeys-lfu` | Any key — least frequently used |
| `volatile-lru` | Only keys **with TTL** — LRU among those |
| `volatile-ttl` | Keys with TTL — shortest remaining TTL first |
| `allkeys-random` / `volatile-random` | Random key |

---

## Interview lines

- Pure cache → `allkeys-lru` or `allkeys-lfu`
- Mix permanent + cache keys → `volatile-lru` + TTL on cache only
- Real Redis uses **approximate** LRU (samples keys), not exact linked-list LRU like [../LRU/](../LRU/)
- **Eviction** (memory full) ≠ **expiration** (TTL) — related but different
