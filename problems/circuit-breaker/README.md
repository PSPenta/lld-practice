# Circuit breaker — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Per dependency or global?
2. Failure threshold to open?
3. Half-open probe count?
4. Library wrapper or sidecar?
5. Which failures count (timeout, 5xx)?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Application code, downstream dependency |
| **Use cases (v1)** | 1. Execute call through breaker 2. Fail fast when OPEN 3. Probe recovery in HALF_OPEN |
| **In scope** | CLOSED/OPEN/HALF_OPEN states, failure counter, Execute(fn) |
| **Out of scope** | Distributed breaker sync across pods |
| **Assumptions** | In-process library; count consecutive failures |

### Confirm understanding
> "Callers wrap downstream calls; after N failures we fail fast until a cooldown probe succeeds."

---

## Step 2 — Entities & classes

`CircuitBreaker` states CLOSED/OPEN/HALF_OPEN, `CallExecutor`, `FailureCounter`, `Clock`

---

## Step 3 — Flows

Call → if OPEN fail fast → if CLOSED run → on failures exceed threshold → OPEN → after timeout HALF_OPEN → probe

---

## Step 4 — APIs

`Execute(fn) (result, error)` wraps any downstream call

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Thread-safe counters; don't count timeouts as success; half-open allows limited traffic

---

## Step 6 — Evolve

Per-tenant breakers; metrics export; bulkhead separate thread pools

---

