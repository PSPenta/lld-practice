# Webhook delivery system — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Retry policy?
2. Signing secret?
3. Ordering per subscriber?
4. DLQ?
5. Timeout?
6. Fan-out to many subscribers per event?
7. At-least-once to subscriber?
8. Admin replay from DLQ?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Event source, DeliveryWorker, Subscriber URL |
| **Use cases (v1)** | 1. Register webhook 2. Deliver event POST 3. Retry with backoff |
| **In scope** | Sign payload, retry, idempotency key |
| **Out of scope** | Subscriber management UI |
| **Assumptions** | HMAC signature; exponential backoff |

### Confirm understanding
> "Events POST to subscriber URLs with retries until ack or DLQ."

## Step 2 — Entities & classes

```text
WebhookSubscription {
  id, url, secret, eventTypes[], active
}

DeliveryJob {
  id, subscriptionId, eventId, payload
  attempt, nextRunAt, status: PENDING|SUCCESS|DLQ
}

SignatureService
  - Sign(payload, secret) → header   // e.g. HMAC-SHA256

DeliveryWorker
  - Enqueue(event)
  - Deliver(job)   // HTTP POST + timeout
  - scheduleRetry / moveToDLQ

WebhookService
  - RegisterWebhook(url, events, secret)
```

**Patterns:** Outbox / queue · Retry with backoff · HMAC signing

## Step 3 — Flows

**Happy path**
1. Domain event occurs → enqueue DeliveryJob per matching subscription  
2. Worker POSTs signed payload with Idempotency-Key / event id  
3. 2xx → SUCCESS  
4. Failure → backoff nextRunAt; after N → DLQ  

**Edge cases**
1. Subscriber slow → hard timeout; don’t block other deliveries  
2. Duplicate event publish → same event id / idempotency so subscriber can dedupe

## Step 4 — APIs

```text
RegisterWebhook(url, events, secret) → subscriptionId
DisableWebhook(id)
EnqueueDelivery(event)          # internal
ReplayDLQ(jobId)                # ops
```

```http
POST /webhooks
DELETE /webhooks/{id}
# outbound: POST {subscriber.url}  Headers: X-Signature, Idempotency-Key
```

## Step 5 — Deepen

- Idempotency-Key / event id for subscriber dedupe  
- Per-request timeout; exponential backoff + jitter  
- Concurrent workers: claim job with lease so two workers don’t POST same job  
- Verify we only retry on network/5xx — not on 410 Gone (disable sub)  
- Outbox from main service so DB commit and enqueue stay consistent

## Step 6 — Evolve

- See Backend docs for production webhook patterns; transactional outbox  
- Ordering per subscription via partition key  
- Related: [retry-scheduler](../retry-scheduler/README.md), [message-queue](../message-queue/README.md), [ticket-notify](../ticket-notify/README.md)
