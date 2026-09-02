# Webhook delivery system — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Retry policy?
2. Signing secret?
3. Ordering per subscriber?
4. DLQ?
5. Timeout?

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

---

## Step 2 — Entities & classes

`WebhookSubscription`, `DeliveryJob`, `DeliveryWorker`, `SignatureService`

---

## Step 3 — Flows

Event → enqueue delivery → POST with signature → retry backoff → DLQ

---

## Step 4 — APIs

`RegisterWebhook(url, events)`, internal `EnqueueDelivery`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Idempotency-Key; timeout; verify subscriber response

---

## Step 6 — Evolve

See Backend §14; outbox from main service

---

