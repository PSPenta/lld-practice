# Circuit breaker — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Per dependency or global?
2. Failure threshold to open?
3. Half-open probe count?
4. Library wrapper or sidecar?
5. Which failures count (timeout, 5xx)?
6. Sliding window vs consecutive failures?
7. Fallback when open (cached value / default / error)?
8. Metrics / events on state change?

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

## Step 2 — Entities & classes

```text
State: CLOSED | OPEN | HALF_OPEN

CircuitBreaker
  - failureThreshold, openTimeout, halfOpenMaxProbes
  - failureCount, lastFailureAt, state
  - Execute(fn) (result, error)
  - onSuccess() / onFailure()

CallExecutor   // runs fn with timeout
FailureCounter // consecutive or windowed
Clock          // injectable for tests
```

**Patterns:** State machine · Decorator / wrapper around outbound calls

## Step 3 — Flows

**Happy path**
1. `Execute(fn)` while CLOSED → run call  
2. Success → reset failure count; stay CLOSED  
3. Failures exceed threshold → transition OPEN; start cooldown timer  
4. After timeout → HALF_OPEN; allow limited probes  
5. Probe success → CLOSED; probe failure → OPEN again  

**Edge cases**
1. Call while OPEN before cooldown → immediate error (fail fast)  
2. Concurrent probes in HALF_OPEN → limit to N in-flight probes

## Step 4 — APIs

```text
NewCircuitBreaker(config) → *CircuitBreaker
Execute(fn) (result, error)   // wraps any downstream call
State() → CLOSED|OPEN|HALF_OPEN
Reset()                       // admin / tests
```

## Step 5 — Deepen

- Thread-safe counters and state transitions (mutex or atomics)  
- Don’t count timeouts as success; decide whether 4xx counts as failure  
- HALF_OPEN allows limited traffic so recovery doesn’t stampede the dependency  
- Fallback hook when open (optional) vs hard error  
- Emit metrics on open/close for ops

## Step 6 — Evolve

- Per-tenant / per-endpoint breakers  
- Sliding-window failure rate instead of consecutive count  
- Bulkhead: separate thread pools / semaphores per dependency  
- Related: [retry-scheduler](../retry-scheduler/README.md), [api-gateway](../api-gateway/README.md)
