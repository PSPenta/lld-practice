# Message queue (ack, DLQ) — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌ · **Code:** `JavaScript/Queue/ (FIFO only)`

## Step 1 — Clarify

### Questions (ask 6–8)
1. At-least-once delivery?
2. Consumer groups?
3. DLQ after N retries?
4. Ordering per partition?
5. Persistence?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Producers, Broker, Consumers |
| **Use cases (v1)** | 1. Publish 2. Subscribe/consume 3. Ack or retry → DLQ |
| **In scope** | Topic, ack, retry, DLQ |
| **Out of scope** | Kafka-level partitioning |
| **Assumptions** | Single broker in-memory; at-least-once |

### Confirm understanding
> "Producers publish; consumers process and ack; failures retry then dead-letter."

---

## Step 2 — Entities & classes

`Broker`, `Topic`, `Message`, `Consumer`, `AckTracker`, `DeadLetterQueue`

---

## Step 3 — Flows

Publish → persist → consumer poll → process → ack → on fail retry → DLQ

---

## Step 4 — APIs

`Publish(topic, body)`, `Subscribe(topic, handler)`, `Ack(msgId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Visibility timeout; idempotent consumers; poison message isolation

---

## Step 6 — Evolve

Compare basic FIFO in `JavaScript/Queue/`; add persistence adapter

---

