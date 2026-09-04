# Connection pool — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Related snippet: `Go/README.md` §13 (illustrative only — not a full solution).

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Max connections?
2. Idle timeout?
3. Block or fail when pool exhausted?
4. Health-check on checkout?
5. Per-DB or shared pool?
6. Min idle warm connections?
7. Validation query before reuse?
8. Fairness under contention (FIFO waiters)?

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

## Step 2 — Entities & classes

```text
PoolConfig { maxOpen, maxIdle, idleTTL, acquireTimeout }

Connection (interface)
  - Query / Exec / Close

PooledConnection { conn, lastUsedAt, inUse }
ConnectionFactory
  - Open() → Connection

Pool
  - Acquire(ctx) → Connection
  - Release(conn)
  - Close()                 // drain & destroy all
  - idle[] , openCount
```

**Patterns:** Object pool · Factory for real connections · RAII / defer-release in callers

## Step 3 — Flows

**Happy path**
1. `Acquire(ctx)` → take from idle list if healthy  
2. Else if `openCount < maxOpen` → factory.Open(), return  
3. Else wait until ctx timeout or a connection is released  
4. Caller uses connection → `Release` returns to idle (or destroys if idle TTL exceeded)  

**Edge cases**
1. Acquire timeout → error; do not leak waiter forever  
2. Broken connection on checkout → discard and create/retry once

## Step 4 — APIs

```text
NewPool(factory, config) → *Pool
Acquire(ctx) (conn, error)
Release(conn)
Close() error
Stats() → { open, idle, waiters }
```

## Step 5 — Deepen

- Mutex (or channel) protecting free list and openCount  
- Context timeout on acquire; cancel removes waiter  
- Always `defer Release` so panic / early return doesn’t leak  
- Idle reaper closes stale connections; Close() is idempotent  
- Don’t return closed/broken connections to idle

## Step 6 — Evolve

- Separate read/write pools; metrics on wait time and create latency  
- Dynamic max under load; health-check on borrow  
- Related: [circuit-breaker](../circuit-breaker/README.md) wrapping factory.Open
