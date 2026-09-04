# Task queue (FIFO) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Bounded capacity?
2. Block or error when full?
3. Multiple consumers?
4. Thread-safe?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Producers, consumers, `FIFOQueue` |
| **Use cases (v1)** | Enqueue · dequeue in order |
| **In scope** | Circular buffer, O(1) enqueue/dequeue |
| **Out of scope** | Priority, persistence, workers |
| **Assumptions** | Fixed capacity; single process |

### Confirm understanding
> "FIFO queue: enqueue at tail, dequeue from head."

---

## Step 2 — Entities & classes

```text
FIFOQueue
  - buffer[], head, tail, size, capacity
  enqueue(item), dequeue() → item | error, isEmpty(), isFull()
```

**Principles:** **KISS**, **YAGNI**

---

## Step 3 — Flows

**Enqueue:** if full → error → else write at tail, advance tail  

**Dequeue:** if empty → error → else read at head, advance head

---

## Step 4 — APIs

`enqueue(item)`, `dequeue()`, `size()`

---

## Step 5 — Deepen

- Mutex if multi-threaded
- Full vs empty distinguish (size counter or sacrifice one slot)

---

## Step 6 — Evolve

- [message-queue](../message-queue/README.md) — ack, retry, DLQ


---

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/Queue/`](../../JavaScript/Queue/) | circular buffer FIFO |
| **Go** | [`Go/Queue-go/`](../../Go/Queue-go/) | |

Full worker/DLQ → [message-queue](../message-queue/README.md)

## Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `Queue/index.js` | `FIFOQueue` — circular buffer, `enqueue` / `dequeue`, capacity |
| `Go/Queue-go/` | Same FIFO structure |

**Read order:** `enqueue` / `dequeue` — how full vs empty is distinguished.

