# Retry scheduler — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Max retries?
2. Exponential backoff + jitter?
3. Dead letter after N?
4. Persist scheduled tasks?
5. Clock skew?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Task submitter, RetryQueue, Worker |
| **Use cases (v1)** | 1. Fail task 2. Schedule retry 3. Execute at due time |
| **In scope** | Backoff policy, due queue, cancel |
| **Out of scope** | Cross-region scheduling |
| **Assumptions** | In-memory delay queue; idempotent tasks |

### Confirm understanding
> "Failed work is re-queued with increasing delay until success or DLQ."

---

## Step 2 — Entities & classes

`RetryPolicy`, `ScheduledTask`, `RetryQueue`, `BackoffCalculator`

---

## Step 3 — Flows

Fail → schedule retry at now+backoff → worker picks due tasks → execute

---

## Step 4 — APIs

`Schedule(task, policy)`, `Cancel(taskId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Idempotent task execution; clock skew tolerance

---

## Step 6 — Evolve

Integrate with message queue DLQ pattern

---

