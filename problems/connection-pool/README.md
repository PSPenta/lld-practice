# Connection pool — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌ · **Code:** `Go/README.md §13 (snippet)`

## Step 1 — Clarify

### Questions (ask 6–8)
1. Max connections?
2. Idle timeout?
3. Block or fail when pool exhausted?
4. Health-check on checkout?
5. Per-DB or shared pool?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Application threads, ConnectionFactory, DB |
| **Use cases (v1)** | 1. Acquire connection 2. Use and release 3. Pool shutdown |
| **In scope** | Acquire/Release, max/min idle, factory |
| **Out of scope** | Read replica routing, sharding |
| **Assumptions** | Single DB; blocking acquire with timeout |

### Confirm understanding
> "Threads borrow a connection from the pool, use it, and return it instead of opening TCP each time."

---

## Step 2 — Entities & classes

`Pool`, `PooledConnection`, `Factory`, `PoolConfig` (max, min, idle TTL)

---

## Step 3 — Flows

Acquire → take idle or create if under max → use → release to idle or destroy if stale

---

## Step 4 — APIs

`Acquire(ctx)`, `Release(conn)`, `Close()` on pool

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Mutex on free list; context timeout on acquire; don't leak on panic (defer release)

---

## Step 6 — Evolve

Separate read/write pools; metrics on wait time; dynamic max under load

---

