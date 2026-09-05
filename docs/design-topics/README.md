# Design Topics

← [Back to hub §10–§14](../../README.md#10-api-design-for-lld)

---

## 10. API design for LLD

Whether HTTP or in-process library APIs. **Full REST/backend interview prep:** **[Backend/README.md](../../Backend/README.md)** (HTTP methods, idempotent payments, pagination, status codes, platform APIs).

### Good practices
- Resource-oriented names (`/tickets`, not `/doCreateTicket`)
- Correct HTTP verbs — **POST** create, **PUT** full replace, **PATCH** partial update ([Backend §2](../../Backend/README.md))
- Versioning for public HTTP (`/v1/...`)
- Explicit error model with stable `code` + `trace_id`
- Pagination for lists — **cursor** preferred for large data ([Backend §6](../../Backend/README.md))
- Idempotency for unsafe retries (payments, webhooks, creates) — `Idempotency-Key` ([Backend §3–§4](../../Backend/README.md))
- Validation at the boundary (422 + field errors)
- Rate limiting → 429 + `Retry-After` ([Backend §7](../../Backend/README.md); [Rate Limiter walkthrough](../../problems/rate-limiter/README.md))

### HTTP status codes (quick reference)

| Code | Use |
|------|-----|
| 200 / 201 / 202 / 204 | Success (read / created / async accepted / delete) |
| 400 / 422 | Bad JSON / semantic validation |
| 401 / 403 | Not authenticated / not authorized |
| 404 / 409 | Not found / conflict (duplicate, wrong state) |
| 429 | Rate limited |
| 502 / 503 / 504 | Upstream failure — retry with backoff |

Full table: **[Backend/README.md §12](../../Backend/README.md)**.

### Example error shape

```json
{
  "code": "RATE_LIMITED",
  "message": "Too many requests",
  "trace_id": "req_abc123"
}
```

### Library vs service
Ask: “Is this a library used in-process, or a networked service?”  
Design public methods vs HTTP endpoints accordingly.

---

## 11. Data modeling & state

### Entities & relationships
Draw:
- User 1—* Order  
- Order 1—* OrderItem  
- ParkingLot 1—* Floor 1—* Slot  

### State machines
For lifecycle objects (ticket, order, payment):

```text
CREATED → PAID → SHIPPED → DELIVERED
              ↘ CANCELLED
```

Enforce legal transitions in domain methods (`MarkPaid()`), not scattered if-else in controllers.

### Persistence choices in LLD
You often stay in-memory for machine coding. Still mention:
- What would be a row/document in DB
- Indexes for lookup keys (`user_id`, `short_code`)
- Unique constraints for idempotency keys

### Derived vs stored
- Stored: balance ledger entries  
- Derived: available slots count (or stored with careful updates)

---

## 12. Concurrency, idempotency & failure

Even in LLD discussion, seniors are expected to mention these.

### Concurrency
- Shared maps/counters need locks (`Mutex` / `RWMutex`) or single-threaded ownership via channels  
- Optimistic locking: `UPDATE ... WHERE version=?`  
- Avoid holding locks during slow I/O (HTTP, LLM calls)

### Idempotency
Same request retried twice should not double-charge / double-create.

Mechanism: client sends `Idempotency-Key`; server stores key → response.

### Failure handling
- Timeouts on external calls  
- Retries with backoff for transient errors  
- Partial failure: core write succeeds; side effects retry via queue  
- Clear error types for callers  

### Minimal reliability story

> “I’ll keep the critical write in a short DB transaction, perform the external call outside the transaction, and use an outbox/queue for retries if the side effect fails.”

---

## 13. Extensibility & evolution

Interviewers often ask: “What if we need X tomorrow?”

Show you can evolve without a rewrite:

| Change | Design move |
|--------|-------------|
| New payment provider | Adapter + Factory |
| New rate limit algorithm | Strategy |
| New notification channel | Notifier interface |
| 10× traffic | Add cache, queue async work, partition data |
| Multi-tenant | Thread `tenant_id` through APIs & storage keys |

**Rule:** start simple (one implementation), introduce abstractions when the second variant appears or is clearly required.

---

## 14. Common LLD problems — how to think

For each problem: actors → use cases → entities → classes → APIs → concurrency → extend.

Must-practice set (core ideas per problem): **[hub §14](../../README.md#14-common-lld-problems--how-to-think)** — not repeated here.

### AI-era LLD variants (increasingly asked)

Classic LLD skills still apply. AI rounds add **provider abstraction, RAG, credits, streaming, and unsafe model output**. See **[../ai-lld/README.md](../ai-lld/README.md)**, **[Cache Client design](../../problems/cache-client/README.md)**, and **[AI Suggest Reply design](../../problems/ai-suggest-reply/README.md)**.

---

