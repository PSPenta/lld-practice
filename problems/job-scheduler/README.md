# Job scheduler — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Cron or fixed delay?
2. Persist jobs?
3. Exactly-once or at-least-once?
4. Priority?
5. Worker count?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Scheduler, WorkerPool, JobStore |
| **Use cases (v1)** | 1. Schedule job 2. Worker executes 3. Retry on failure |
| **In scope** | Priority queue, schedule API, worker loop |
| **Out of scope** | Distributed cron leader election |
| **Assumptions** | In-memory queue; at-least-once; exponential backoff |

### Confirm understanding
> "Jobs are scheduled with a run time; workers pull and execute with retries."

---

## Step 2 — Entities & classes

`Job`, `Scheduler`, `WorkerPool`, `JobQueue` (priority heap), `JobStore`

---

## Step 3 — Flows

Schedule job → push queue → worker pulls → execute → mark success/fail → retry with backoff

---

## Step 4 — APIs

`Schedule(cron, payload)`, `Cancel(jobId)`, `GetStatus(jobId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Leader election for scheduler; DLQ for poison jobs; mutex on job state

---

## Step 6 — Evolve

Horizontal workers; shard queue by tenant

---

