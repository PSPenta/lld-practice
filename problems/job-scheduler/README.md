# Job scheduler — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Cron or fixed delay?
2. Persist jobs?
3. Exactly-once or at-least-once?
4. Priority?
5. Worker count?
6. Max retries / DLQ?
7. Cancel in-flight or only pending?
8. Single process or multi-node?

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

## Step 2 — Entities & classes

```text
Job {
  id, payload, cronExpr? | runAt
  priority, status: PENDING|RUNNING|SUCCEEDED|FAILED|CANCELLED
  attempts, maxAttempts
}

JobQueue          // priority heap by (runAt, priority)
JobStore          // optional persistence
WorkerPool
  - start(n), stop()
  - pullDue() → Job; execute handler

Scheduler
  - Schedule(spec, payload) → jobId
  - Cancel(jobId)
  - GetStatus(jobId)
```

**Patterns:** Priority queue · Worker pool · Strategy for cron vs one-shot

## Step 3 — Flows

**Happy path**
1. Client Schedule → validate → push Job onto queue with `runAt`  
2. Worker loop wakes / polls due jobs  
3. Mark RUNNING → execute handler  
4. Success → SUCCEEDED; failure → schedule retry with backoff or FAIL  

**Edge cases**
1. Cancel while PENDING → remove; while RUNNING → cooperative cancel if supported  
2. Poison job exceeds maxAttempts → DLQ / FAILED terminal; don’t block workers

## Step 4 — APIs

```text
Schedule(cronOrRunAt, payload, opts?) → jobId
Cancel(jobId) error
GetStatus(jobId) → Job
ListDue(now) → []Job          // internal
```

## Step 5 — Deepen

- Mutex / heap lock around queue; workers don’t steal same job twice  
- At-least-once: handler must be idempotent  
- DLQ for poison jobs; backoff with jitter  
- Persist jobs if process can restart mid-flight  
- Multi-node later needs leader election — call out as out of scope for v1

## Step 6 — Evolve

- Horizontal workers; shard queue by tenant  
- Persistent store + lease/heartbeat for RUNNING jobs  
- Related: [retry-scheduler](../retry-scheduler/README.md), [message-queue](../message-queue/README.md)
