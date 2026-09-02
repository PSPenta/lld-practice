# Load balancer — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. L4 or L7?
2. Algorithm: RR, least conn, weighted?
3. Health checks?
4. Sticky sessions?
5. How many backends?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Clients, LoadBalancer, Backend pool |
| **Use cases (v1)** | 1. Forward request to healthy backend 2. Register/deregister backend |
| **In scope** | Selection strategy, health flag, forward |
| **Out of scope** | Global anycast, autoscaling |
| **Assumptions** | Round-robin; HTTP pass-through; 3 backends |

### Confirm understanding
> "Client sends request; LB picks a healthy backend and proxies the response."

---

## Step 2 — Entities & classes

`LoadBalancer`, `Backend`, `HealthChecker`, `SelectionStrategy`

---

## Step 3 — Flows

Request → pick healthy backend via strategy → forward → return response

---

## Step 4 — APIs

`Handle(request)` or `RegisterBackend(host)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Mark unhealthy after N failures; thread-safe backend list

---

## Step 6 — Evolve

Consistent hash for sticky; active health probes

---

