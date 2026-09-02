# Rate Limiter — LLD design walkthrough

> **Code:** this folder (Strategy pattern — **preferred demo**) · v1 teaching: [../Ratelimiter/](../Ratelimiter/) · **Go:** [../../Go/RateLimiter2-go/](../../Go/RateLimiter2-go/)  
> **Method:** [../../README.md §5](../../README.md) · **Distributed variant:** [../../lld-gaps/README.md §10](../../lld-gaps/README.md)

---

## Clarify

- Per user / IP / API key?  
- Limit & window?  
- Burst allowed?  
- Single node or distributed?

---

## Strategies (Strategy pattern)

| Strategy | Idea | Pros | Cons |
|----------|------|------|------|
| Fixed window | Count in current minute | Simple | Burst at window edge |
| Sliding window log | Store timestamps | Accurate | Memory heavy |
| Token bucket | Tokens refill over time | Smooth + burst | Slightly more logic |
| Leaky bucket | Steady outflow | Smooth egress | Less burst friendly |

---

## Classes

```text
RateLimiterStrategy { Allow(key) bool }
RateLimiter { strategy; Allow(key) }
TokenBucketStrategy { capacity, refillRate, buckets map, mu }
```

---

## API

- Library: `Allow(key string) bool`  
- HTTP middleware: return `429` when false  

---

## Evolve to distributed

Store counters/tokens in Redis; accept approximate limits under race, or use Lua for atomicity.
