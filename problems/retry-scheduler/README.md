# Retry scheduler — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Max retries?
2. Exponential backoff + jitter?
3. Dead letter after N?
4. Persist scheduled tasks?
5. Clock skew?
6. Cancel pending retries?
7. Which errors are retryable vs permanent?
8. Single worker or pool?

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

## Step 2 — Entities & classes

```text
RetryPolicy {
  maxAttempts, baseDelay, maxDelay, jitter
  - nextDelay(attempt) → duration
}

ScheduledTask {
  id, payload, attempt, runAt, status
}

BackoffCalculator   // exp backoff + full/equal jitter
RetryQueue          // min-heap by runAt
Worker
  - pollDue(now) → tasks; execute; on fail reschedule

RetryScheduler
  - Schedule(task, policy)
  - Cancel(taskId)
```

**Patterns:** Strategy (backoff) · Delay queue · DLQ for exhausted attempts

## Step 3 — Flows

**Happy path**
1. Work fails (retryable) → Schedule with attempt=1, runAt=now+backoff  
2. Worker picks due tasks → execute  
3. Success → done; fail → attempt++ → reschedule or DLQ if max  

**Edge cases**
1. Cancel before run → remove from queue; ignore late execution if cancelled  
2. Clock skew / delayed worker → still safe if task idempotent

## Step 4 — APIs

```text
Schedule(task, policy) → taskId
Cancel(taskId) error
GetStatus(taskId)
```

## Step 5 — Deepen

- Idempotent task execution required (at-least-once wakeups)  
- Jitter avoids thundering herd; cap maxDelay  
- Thread-safe heap; only one worker claims a due task  
- Distinguish retryable vs permanent errors  
- Persist queue if process restart must not lose retries

## Step 6 — Evolve

- Integrate with [message-queue](../message-queue/README.md) DLQ pattern  
- Shared design with [job-scheduler](../job-scheduler/README.md) and [webhook-delivery](../webhook-delivery/README.md) backoff  
- Distributed delay via Redis ZSET / SQS delay
