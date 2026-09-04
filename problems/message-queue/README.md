# Message queue (ack, DLQ) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Basic FIFO practice exists at `JavaScript/Queue/` — this problem adds ack, retry, DLQ.

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. At-least-once delivery?
2. Consumer groups?
3. DLQ after N retries?
4. Ordering per partition?
5. Persistence?
6. Visibility timeout / ack deadline?
7. Fan-out (pub-sub) vs competing consumers?
8. Payload size limits?

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

## Step 2 — Entities & classes

```text
Message {
  id, topic, body, attempts
  visibleAt, acked
}

Topic { name, queue[], dlq[] }
AckTracker { inFlight map[msgId]deadline }

Consumer { id, topic, handler }
Broker
  - Publish(topic, body) → msgId
  - Subscribe(topic, handler)
  - Ack(msgId) / Nack(msgId)
```

**Patterns:** Broker as facade · visibility timeout · DLQ as separate queue

## Step 3 — Flows

**Happy path**
1. Publish → append to topic queue, assign id  
2. Consumer poll / push → mark in-flight with visibility deadline  
3. Process → Ack → remove from in-flight  
4. On Nack / timeout → attempts++ → requeue or after N → DLQ  

**Edge cases**
1. Consumer crash before ack → message reappears after visibility timeout  
2. Poison message always fails → isolated in DLQ after max attempts

## Step 4 — APIs

```text
Publish(topic, body) → msgId
Subscribe(topic, handler)
Ack(msgId)
Nack(msgId)                 // optional explicit
GetDLQ(topic) → []Message   // ops
```

## Step 5 — Deepen

- Visibility timeout prevents double-delivery while still allowing retry  
- Consumers must be idempotent (at-least-once)  
- Poison message isolation via DLQ  
- Concurrent pollers: only one consumer gets a message until timeout  
- Persist if broker restart must not lose unacked messages (call out evolve)

## Step 6 — Evolve

- Compare basic FIFO in `JavaScript/Queue/`; add persistence adapter  
- Consumer groups / partitions for ordering + scale  
- Related: [retry-scheduler](../retry-scheduler/README.md), [pub-sub](../pub-sub/README.md), [webhook-delivery](../webhook-delivery/README.md)
