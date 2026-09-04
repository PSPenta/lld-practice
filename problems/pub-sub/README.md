# Pub/Sub system — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/Pub-Sub/`](../../JavaScript/Pub-Sub/) | in-memory Observer |
| **Go** | [`Go/Pub-Sub-go/`](../../Go/Pub-Sub-go/) | 3 variants — good evolution discussion |

### Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `pubsubCoreJS.js` | Core topic → callback list; subscribe / publish / unsubscribe |
| `pubsubEventEmitter.js` | Variant using EventEmitter-style API |
| `index.js` | Demo fan-out |
| `Go/Pub-Sub-go/` | Multiple variants for “how would you evolve?” |

**Read order:** `subscribe` / `publish` — copy listener list before invoke.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Topics or direct subscriber list?
2. Sync or async delivery?
3. At-least-once?
4. Ordering per topic?
5. Persistence?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Publishers, subscribers, `PubSub` bus |
| **Use cases (v1)** | Subscribe · publish · unsubscribe |
| **In scope** | topic → callbacks map, fan-out |
| **Out of scope** | Persistent log, consumer groups |
| **Assumptions** | In-memory; sync callbacks |

### Confirm understanding
> "Publisher emits on topic; all subscribers receive callback."

---

## Step 2 — Entities & classes

```text
PubSub
  - subscribers: Map<event, callback[]>
  subscribe(event, cb), unsubscribe(event, cb), publish(event, payload)
```

**Pattern:** **Observer** / Pub-Sub

---

## Step 3 — Flows

**Subscribe:** append callback to event bucket  
**Publish:** copy callback list → invoke each with payload (avoid mutation during iterate)  
**Unsubscribe:** remove callback reference

---

## Step 4 — APIs

`subscribe(topic, handler)`, `unsubscribe(topic, handler)`, `publish(topic, data)`

---

## Step 5 — Deepen

- Copy subscriber list before iteration
- Handler throw → don't break other subscribers (policy choice)

---

## Step 6 — Evolve

- Async queue per subscriber; evolve to [message-queue](../message-queue/README.md)
