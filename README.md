# Low-Level Design (LLD) Interview Preparation Guide

> A practical, detailed guide to prepare for **LLD / machine-coding / OOD** interview rounds at product companies.  
> This repository also contains working implementations (JavaScript + Golang) of common LLD problems — use them for practice after you learn the method.

---

## Table of contents

1. [What is LLD?](#1-what-is-lld)
2. [What interviewers evaluate](#2-what-interviewers-evaluate)
3. [LLD vs HLD vs DSA](#3-lld-vs-hld-vs-dsa)
4. [How a typical LLD round runs](#4-how-a-typical-lld-round-runs)
5. [The standard approach (memorize this)](#5-the-standard-approach-memorize-this)
6. [Clarifying questions checklist](#6-clarifying-questions-checklist)
7. [OOP building blocks](#7-oop-building-blocks)
8. [SOLID principles (with examples)](#8-solid-principles-with-examples)
9. [Design patterns that actually matter](#9-design-patterns-that-actually-matter)
10. [API design for LLD](#10-api-design-for-lld)
11. [Data modeling & state](#11-data-modeling--state)
12. [Concurrency, idempotency & failure](#12-concurrency-idempotency--failure)
13. [Extensibility & evolution](#13-extensibility--evolution)
14. [Common LLD problems — how to think](#14-common-lld-problems--how-to-think)
15. [AI / LLM LLD for beginners](#15-ai--llm-lld-for-beginners)
16. [Worked example: LRU Cache](#16-worked-example-lru-cache)
17. [Worked example: Rate Limiter](#17-worked-example-rate-limiter)
18. [Worked example: Parking Lot](#18-worked-example-parking-lot)
19. [Worked example: AI Suggest Reply](#19-worked-example-ai-suggest-reply-copilot-draft)
20. [How to practice (4-week plan)](#20-how-to-practice-4-week-plan)
21. [Interview day checklist](#21-interview-day-checklist)
22. [Problems in this repository](#22-problems-in-this-repository)
23. [Cheat sheet](#23-cheat-sheet)

---

## 1. What is LLD?

**Low-Level Design** is designing software at the **class / module / API** level.

You answer questions like:

- What are the main **entities**?
- What are the **classes/interfaces** and their **methods**?
- How do objects **collaborate** for each use case?
- What does the **API** look like?
- How do we keep the design **clean, testable, and extensible**?
- What happens under **concurrency**, **failures**, and **growth**?

You are **not** primarily designing AWS topology (that is HLD).  
You **are** designing the code structure someone would implement in a sprint.

### Machine coding vs discussion LLD

| Style | What happens |
|-------|----------------|
| **Discussion LLD** | Whiteboard / shared doc: classes, flows, APIs, trade-offs (60–75 min) |
| **Machine coding** | Write working code in 90–120 min for a small system |

Same fundamentals. Machine coding needs speed + compiling code. Discussion needs clearer verbal trade-offs.

---

## 2. What interviewers evaluate

At senior / mid levels they score:

| Skill | What “good” looks like |
|-------|-------------------------|
| Problem breakdown | Clarifies ambiguity; states assumptions |
| Modeling | Correct entities & relationships |
| Abstraction | Interfaces where variation exists |
| SOLID / clean design | Responsibilities are separated |
| Patterns | Used purposefully, not decoration |
| APIs | Intuitive, consistent, versionable |
| Data & state | Clear persistence & transitions |
| Concurrency | Mentions races and how to prevent them |
| Failure handling | Timeouts, retries, partial failure |
| Extensibility | Can add a feature without rewriting core |
| Communication | Structured; trade-offs explained |

There is rarely one “correct” design. **Sound decisions + clear reasoning** wins.

---

## 3. LLD vs HLD vs DSA

| | DSA | LLD | HLD |
|---|-----|-----|-----|
| Focus | Algorithms, complexity | Classes, APIs, modules | Systems, scale, infra |
| Example | Shortest path in graph | Design Rate Limiter classes | Design URL shortener at 100M QPS |
| Output | Code + Big-O | Class diagram + APIs + flows | Boxes: LB, DB, cache, queues |

Many companies mix a little scale talk into LLD (“what if traffic grows?”). Answer at **component** level first, then briefly say how you’d evolve.

---

## 4. How a typical LLD round runs

**Duration:** 45–60 minutes (sometimes 90 for machine coding).

Suggested time split:

| Minutes | Activity |
|---------|----------|
| 0–5 | Clarify requirements & assumptions |
| 5–12 | Entities + relationships |
| 12–30 | Classes/interfaces + main flows |
| 30–40 | APIs (+ light DB schema if asked) |
| 40–50 | Concurrency, failures, extensibility |
| 50–60 | Trade-offs, evolve design, Q&A |

**Opening line you can use:**

> “I’ll clarify the requirements first, then outline entities and responsibilities, design classes and key workflows, define APIs, and finally discuss concurrency, failure modes, and how the design evolves.”

---

## 5. The standard approach (memorize this)

```text
Understand → Model → Design → APIs → Deepen → Evolve
```

### Step-by-step

1. **Understand** — actors, use cases, constraints, non-goals  
2. **Model** — entities, relationships, state machines  
3. **Design** — classes, interfaces, responsibilities, collaborations  
4. **APIs** — endpoints or public methods of the library  
5. **Deepen** — concurrency, idempotency, caching, errors  
6. **Evolve** — “If requirements change / scale 10×, I would…”

Never jump to patterns or microservices before entities and use cases are clear.

---

## 6. Clarifying questions checklist

Ask a subset relevant to the problem:

**Product / scope**
- Who are the users (actors)?
- What are the top 3 use cases for v1?
- What is explicitly out of scope?

**Scale (order of magnitude)**
- Users? Requests/sec? Data size?
- Single machine library or distributed service?

**Functional**
- Auth / multi-tenant?
- Sync or async?
- Persistence needed? TTL?
- Exact vs approximate correctness?

**Quality**
- Consistency needs?
- Failure expectations?
- Observability needs?

Then **state assumptions out loud** so the interviewer can correct you.

---

## 7. OOP building blocks

### Class
Blueprint for objects — data + behavior.

### Interface / abstract behavior
Contract without implementation. Callers depend on the contract.

```text
Notifier
  + Send(message) error

EmailNotifier implements Notifier
SlackNotifier implements Notifier
```

### Encapsulation
Hide internals; expose a small API (`Park`, `Unpark`), not raw fields.

### Composition over inheritance
Prefer “has-a” over deep “is-a” trees.

```text
OrderService has PaymentGateway
OrderService has InventoryService
```
Not: `OrderService extends DatabaseHelper extends Logger...`

### Association types (useful vocabulary)
- **Association** — uses  
- **Aggregation** — has (can exist independently)  
- **Composition** — owns (lifecycle tied)

---

## 8. SOLID principles (with examples)

### S — Single Responsibility Principle
A class should have **one reason to change**.

Bad: `UserService` that registers users, sends email, and writes billing.

Good:
- `UserService` — registration rules  
- `EmailNotifier` — emails  
- `BillingService` — charges  

### O — Open/Closed Principle
Open for **extension**, closed for **modification**.

Add `SlackNotifier` without rewriting `NotificationService` — program to `Notifier` interface.

### L — Liskov Substitution Principle
Subtypes must be usable wherever the parent/interface is expected, without breaking callers.

If `Notifier.Send` is documented to return errors (not panic), every implementation must honor that.

### I — Interface Segregation Principle
Prefer small, focused interfaces over fat ones.

Bad: `Worker` with `Read()`, `Write()`, `AdminFlush()`, `RenderPDF()`.  
Good: separate interfaces per capability.

### D — Dependency Inversion Principle
High-level modules depend on **abstractions**, not concrete vendors.

```text
SuggestService → LLMClient (interface)
                    ↑
           OpenAIClient / AnthropicClient
```

This is critical for payments, LLMs, email providers, SMS gateways, etc.

---

## 9. Design patterns that actually matter

Use patterns to solve a real variation problem — not to impress.

| Pattern | Intent | Classic LLD use |
|---------|--------|-----------------|
| **Strategy** | Swap algorithms | Rate limit algorithms; fare calculation; payment methods |
| **Factory / Abstract Factory** | Create objects without exposing construction | VehicleFactory; ExpenseFactory |
| **Adapter** | Make incompatible APIs work together | Wrap bank/UPI/card behind one gateway interface |
| **Observer / Pub-Sub** | Notify many listeners | Topic publish → multiple subscribers |
| **Singleton** | One shared instance | Use carefully (config); often DI is better |
| **Decorator** | Add behavior dynamically | Auth middleware, logging wrapper |
| **Repository** | Abstract persistence | `TicketRepository` |
| **Builder** | Construct complex objects step-by-step | Query builder; HTTP request builder |
| **Command** | Encapsulate an action | Undo/redo; job queue tasks |
| **State** | Behavior changes with state | Ticket Open/Pending/Closed transitions |
| **Template Method** | Shared algorithm skeleton | Import pipelines with customizable steps |

### Patterns used in this repository (with reference files)

| Pattern | LLD in this repo | How it shows up | Reference files |
|---------|------------------|-----------------|-----------------|
| **Strategy** | **RateLimiter2** | `RateLimiter` holds a `RateLimiterStrategy`; swap Fixed Window / Token Bucket / Leaky Bucket / Sliding Window without changing the caller | JS: `RateLimiter2/RateLimiterStrategy.js`, `RateLimiter2/RateLimiter.js`, `RateLimiter2/TokenBucket.js`, `RateLimiter2/LeakyBucket.js`, `RateLimiter2/FixedWindowCounter.js`, `RateLimiter2/SlidingWindowLog.js`, `RateLimiter2/SlidingWindowCounter.js` · Go: `Go/RateLimiter2-go/strategy.go`, `Go/RateLimiter2-go/token_bucket.go`, `Go/RateLimiter2-go/leaky_bucket.go`, `Go/RateLimiter2-go/fixed_window_counter.go`, `Go/RateLimiter2-go/sliding_window_log.go`, `Go/RateLimiter2-go/sliding_window_counter.go`, `Go/RateLimiter2-go/main.go` |
| **Strategy + Adapter-style interface** | **PaymentGateway** | `BankGateway` interface; `UPIGateway` / `CardGateway` / `NetBankingGateway` are interchangeable processors | Go: `Go/PaymentGateway-go/bank_gateway.go`, `Go/PaymentGateway-go/payment_gateway.go`, `Go/PaymentGateway-go/payment.go`, `Go/PaymentGateway-go/main.go` |
| **Factory** | **Splitwise** | `ExpenseFactory` / `CreateExpense` builds Exact / Equal / Percentage expenses | JS: `Splitwise/Expense.js`, `Splitwise/index.js` · Go: `Go/Splitwise-go/expense.go`, `Go/Splitwise-go/main.go` |
| **Factory** | **ParkingLot2** | `VehicleFactory` / `CreateVehicle` creates Bike / Car / Truck | JS: `ParkingLot2/Vehicle.js`, `ParkingLot2/index.js` · Go: `Go/ParkingLot2-go/vehicle.go`, `Go/ParkingLot2-go/main.go` |
| **Observer / Pub-Sub** | **Pub-Sub** | Subscribe handlers to a topic/event; publish fans out to listeners | JS: `Pub-Sub/` · Go: `Go/Pub-Sub-go/pubsub.go`, `Go/Pub-Sub-go/pubsub_core.go`, `Go/Pub-Sub-go/pubsub_event.go`, `Go/Pub-Sub-go/main.go` |
| **Inheritance / specialization** (not full Strategy) | **Parkinglot** (v1) | `Car` / `Bike` / `Truck` extend `Vehicle` | JS: `Parkinglot/Vehicle.js` · Go: `Go/Parkinglot-go/vehicle.go` |

**Not clearly implemented as named patterns in this repo (know them for interviews anyway):** Singleton, Decorator, Repository, Builder, Command, State, Template Method. You can still *discuss* them (e.g. URL shortener HTTP handlers as middleware/decorator-style; ticket status as State) even when the folder does not label them.

**Note:** `Ratelimiter` / `Go/Ratelimiter-go` implement multiple algorithms as **standalone types**, not behind a shared Strategy interface. For Strategy pattern demos, prefer **RateLimiter2**.

### How to mention a pattern in interview

> “I’ll use Strategy here because we already have two interchangeable algorithms and may add a third — same idea as `RateLimiter2` in my practice repo.”

If only one implementation exists and none planned, skip the pattern.

---

## 10. API design for LLD

Whether HTTP or in-process library APIs:

### Good practices
- Resource-oriented names (`/tickets`, not `/doCreateTicket`)
- Consistent verbs / method names
- Versioning for public HTTP (`/v1/...`)
- Explicit error model
- Pagination for lists (cursor preferred for large data)
- Idempotency for unsafe retries (payments, webhooks, creates)
- Validation at the boundary

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

### Must-practice set

| Problem | Core ideas |
|---------|------------|
| **LRU Cache** | Hash map + doubly linked list; O(1) get/put; capacity eviction |
| **Rate Limiter** | Strategy (fixed window, sliding, token bucket); per-key state; thread safety |
| **Parking Lot** | Floors/slots/vehicles; strategy for slot assignment; ticket |
| **Splitwise** | Users, expenses, split strategies; balances; settle-up |
| **Pub-Sub** | Topics, subscribers, publish fan-out; sync vs async |
| **URL Shortener** (LLD slice) | encode ID, mappings, APIs (full scale is HLD) |
| **Elevator / Traffic lights** | State machine; scheduling strategy |
| **Chess / Snake & Ladder** | Board, pieces/moves, game rules |
| **Logging framework** | Levels, appenders, chain of responsibility |
| **ATM / Vending** | State + chain for cash dispense |
| **Notification system** | Templates, channels, retry |
| **Booking (movie/hotel)** | Seats, locks, payment, idempotency |
| **Cache client** | get/set/TTL/LRU; optional loader; stampede |
| **Task / Job scheduler** | Priority queue; workers; retries |

### AI-era LLD variants (increasingly asked)

Classic LLD skills still apply. AI rounds add **provider abstraction, RAG, credits, streaming, and unsafe model output**. See **[§15. AI / LLM LLD for beginners](#15-ai--llm-lld-for-beginners)** and **[§19. Worked example: AI Suggest Reply](#19-worked-example-ai-suggest-reply)**.

---

## 15. AI / LLM LLD for beginners

This section is for people new to AI interviews. You are **not** asked to train models. You are asked to **design application components** that *use* LLMs safely.

### What interviewers mean by “AI LLD”

| They want | They usually do **not** want |
|-----------|------------------------------|
| Classes/interfaces around LLM calls | Deriving transformer math |
| How RAG pieces connect | Training a custom model from scratch |
| Credits, timeouts, retries, validation | Picking the “best” model brand only |
| Multi-tenant isolation of data | Building a vector DB from scratch |

Think: **same as PaymentGateway** — wrap a vendor behind an interface; orchestrate a workflow; handle failure and cost.

### Words you must understand (plain English)

| Term | Meaning |
|------|---------|
| **LLM** | A text model API (OpenAI, etc.) you call with a prompt; it returns text |
| **Prompt** | Instructions + context you send to the model |
| **Token** | Chunk of text the model bills/limits by (roughly pieces of words) |
| **Hallucination** | Confident but wrong answer — why we add retrieval + validation |
| **Embedding** | Turning text into a list of numbers so “similar meaning” is searchable |
| **Vector DB** | Store embeddings; find nearest neighbors for a query |
| **RAG** | Retrieve relevant docs → put them in the prompt → then generate |
| **Chunking** | Splitting long docs into smaller pieces before embedding |
| **Tool / function calling** | Model asks your app to run a function (search DB, create ticket) |
| **Guardrails** | Checks before/after model: size limits, PII, blocked topics, schema |
| **SSE / streaming** | Send answer token-by-token to the UI instead of one big wait |
| **Credit metering** | Limit how much AI usage a tenant/user can consume |

### Core building blocks (LEGO pieces)

Almost every AI feature LLD uses some of these:

```text
1. API / Handler          — auth, validate input
2. Orchestrator Service   — runs the steps in order
3. Guardrails             — input checks / PII redact
4. CreditMeter            — reserve / commit / release credits
5. Retriever (RAG)        — find relevant chunks for this tenant
6. PromptBuilder          — build the final prompt string
7. LLMClient (interface)  — talk to any model vendor
8. OutputValidator        — parse/validate model output (JSON/schema)
9. SuggestionStore / Repo — save drafts / history
10. StreamSession (optional) — SSE to browser
```

### Design patterns that show up in AI LLD

| Pattern | AI use |
|---------|--------|
| **Strategy** | Model routing by plan (cheap vs strong model); different retrievers |
| **Adapter** | `OpenAIClient` / `AnthropicClient` behind `LLMClient` |
| **Factory** | Create tool handlers by name |
| **Observer / Pub-Sub** | “Suggestion completed” → analytics / notify |
| **Decorator / Middleware** | Auth, tenant, credit check around the API |
| **Repository** | TicketStore, ChunkStore, SuggestionStore |

Repo references for the same ideas (non-AI, but same patterns):
- Strategy: `RateLimiter2/`, `Go/RateLimiter2-go/`
- Adapter-style interface: `Go/PaymentGateway-go/bank_gateway.go`
- Factory: `Go/Splitwise-go/expense.go`, `Go/ParkingLot2-go/vehicle.go`

### RAG — simplest correct story

```text
Offline (indexing):
  Document → split into Chunks → Embedding → save in Vector DB
             (always store tenant_id on each chunk)

Online (user asks / agent clicks Suggest):
  User text → Embedding → search Vector DB (filter tenant_id)
           → top-K chunks → add to Prompt → LLM → answer
```

**Why RAG?** Model alone may invent facts. Retrieved chunks ground the answer.

**Must say in interview:** retrieval is **tenant-scoped** — never mix Customer A’s docs into Customer B’s prompt.

### Sync vs async (when to choose)

| Approach | Use when |
|----------|----------|
| **Sync + SSE stream** | User is waiting for a draft reply in UI |
| **Async queue + worker** | Bulk re-embed docs, nightly summaries, non-interactive jobs |

Don’t call a slow LLM inside a DB transaction. Don’t do heavy embedding on the webhook request path — enqueue it.

### Reliability checklist for AI components

1. **Timeout** every LLM/retriever call (`context` deadline)  
2. **Retry** only transient failures; use idempotency key for billable calls  
3. **Reserve credits before** call; **commit** on success; **release** on failure  
4. **Validate output** — don’t trust raw model JSON  
5. **Partial failure** — ticket/email still works if AI is down  
6. **Observability** — `trace_id`, latency, cost/tokens, error rate, thumbs-down  

### Minimum AI LLD problems to practice

1. **LLM provider abstraction** — `LLMClient` + 2 adapters + router  
2. **RAG retriever** — chunk metadata, top-K, tenant filter  
3. **AI credit meter** — reserve/commit (like wallet) + rate limit  
4. **Suggest-reply / copilot** — full orchestration (see §19)  
5. **Tool-calling agent (basic)** — model returns tool name → your code runs tool → continue  

### What “good enough for a noob” looks like in the interview

You can draw the boxes, name interfaces, walk the happy path, then say:

- credits under concurrency  
- LLM timeout  
- tenant isolation in RAG  
- how you’d add a new model vendor (new Adapter only)

You do **not** need to implement a vector DB.

---

## 16. Worked example: LRU Cache

### Clarify
- Capacity?  
- TTL needed?  
- Thread-safe?  
- `get` updates recency? (yes for LRU)

### Core idea
- **Map** for O(1) lookup key → node  
- **Doubly linked list** for recency order (head = most recent, tail = least)

### Classes
```text
Node { key, value, prev, next }
LRUCache {
  capacity
  map
  head, tail
  Get(key) → value
  Put(key, value)
  // private: moveToFront, evictTail, removeNode
}
```

### Flows
**Get hit:** move node to front; return value  
**Put existing:** update value; move to front  
**Put new:** insert front; if over capacity, remove tail  

### Concurrency
Wrap `Get`/`Put` with mutex if multi-threaded.

### Extend
- TTL field on node + lazy expiry on get  
- `Cache` interface → MemoryLRU / Redis adapter  

---

## 17. Worked example: Rate Limiter

### Clarify
- Per user / IP / API key?  
- Limit & window?  
- Burst allowed?  
- Single node or distributed?

### Strategies (Strategy pattern)

| Strategy | Idea | Pros | Cons |
|----------|------|------|------|
| Fixed window | Count in current minute | Simple | Burst at window edge |
| Sliding window log | Store timestamps | Accurate | Memory heavy |
| Token bucket | Tokens refill over time | Smooth + burst | Slightly more logic |
| Leaky bucket | Steady outflow | Smooth egress | Less burst friendly |

### Classes
```text
RateLimiterStrategy { Allow(key) bool }
RateLimiter { strategy; Allow(key) }
TokenBucketStrategy { capacity, refillRate, buckets map, mu }
```

### API
- Library: `Allow(key string) bool`  
- HTTP middleware: return `429` when false  

### Evolve to distributed
Store counters/tokens in Redis; accept approximate limits under race, or use Lua for atomicity.

---

## 18. Worked example: Parking Lot

### Clarify
- Multiple floors?  
- Vehicle types & slot types?  
- Pricing?  
- Entry/exit gates count?

### Entities
`ParkingLot`, `Floor`, `Slot`, `Vehicle`, `Ticket`

### Responsibilities
- Find suitable free slot (Strategy: first-fit, type-fit, nearest)  
- Issue ticket on park  
- Calculate fee on unpark  
- Maintain availability counts  

### Classes (sketch)
```text
Vehicle { number, type }
Slot { id, type, isFree, vehicle }
Floor { id, slots[]; FindSlot(vehicleType) }
ParkingLot { floors[]; Park(vehicle); Unpark(ticketId) }
Ticket { id, slotId, entryTime }
PricingStrategy { Calculate(ticket, exitTime) }
```

### Extensibility
New vehicle type → mapping to slot types.  
New pricing → new `PricingStrategy`.

---

## 19. Worked example: AI Suggest Reply (copilot draft)

Use this as your **AI LLD template**. Same structure as LRU/Parking Lot interviews.

### Clarify (ask first!)

1. Draft only or auto-send to customer?  
2. Streaming to UI (SSE) required?  
3. Multi-tenant SaaS? (almost always yes)  
4. One LLM vendor or many?  
5. Do we have a knowledge base / past tickets for RAG?  
6. Credits / quotas per plan?  
7. Approximate QPS of “Suggest” clicks?

**Assumptions to state if they don’t specify:**
- Draft only (human sends)
- SSE streaming
- Multi-tenant
- RAG over KB + ticket thread
- Credits required

### Actors & use cases

- **Agent** clicks Suggest on a ticket  
- **System** builds grounded draft, streams it, stores it  
- **Admin** (optional) uploads KB docs (indexing path — mention briefly)

### Entities

```text
Tenant, Ticket, Message
Suggestion { id, ticket_id, status, model, text, created_by }
Chunk { id, tenant_id, doc_id, text, embedding }
CreditAccount / CreditReservation
IdempotencyRecord
```

### Classes / interfaces

```text
SuggestHandler          // HTTP boundary: auth + validation
SuggestService          // orchestrator (main flow)

TicketRepository        // load ticket + messages
Retriever               // vector search with tenant filter
PromptBuilder           // system rules + chunks + thread
LLMClient (interface)   // Complete / Stream
  OpenAIAdapter
  AnthropicAdapter
OutputValidator         // length / JSON / basic safety
CreditMeter             // Reserve / Commit / Release
Guardrails              // max input size, PII redact
SuggestionRepository    // save draft + status
```

### Happy path (say this in order)

```text
1. AuthZ: agent can access ticket (same tenant)
2. Idempotency-Key → return existing suggestion if replay
3. Guardrails on input size
4. CreditMeter.Reserve(units)
5. Load ticket thread
6. Retriever.topK(query, tenant_id)
7. PromptBuilder.build(...)
8. Save Suggestion(status=streaming)
9. LLMClient.Stream(ctx) → SSE to client
10. Validate final text; save; CreditMeter.Commit
```

### Failure path

```text
LLM timeout / error
  → cancel context
  → CreditMeter.Release
  → Suggestion status=failed
  → UI shows retry (same Idempotency-Key OK)
Ticket itself is unchanged (partial failure is fine)
```

### APIs

```http
POST /v1/tickets/{id}/ai/suggestions
Headers: Authorization, Idempotency-Key
Body: { "tone": "friendly" }
→ 201 { "suggestion_id", "stream_url" }

GET /v1/tickets/{id}/ai/suggestions/{sid}/stream
→ SSE tokens, then done|error
```

### Concurrency

- Two Suggest clicks: two reservations if credits allow; or limit one in-flight suggestion per ticket  
- Credit reserve must be atomic (DB row lock / conditional update)  
- Never hold DB transaction open during LLM HTTP call  

### Patterns used here

| Pattern | Where |
|---------|--------|
| Adapter | LLM vendors behind `LLMClient` |
| Strategy | Model pick by plan / latency |
| Repository | Ticket / Suggestion / Chunk stores |
| Middleware | Auth + credit checks |

### Evolve

| Change | Design move |
|--------|-------------|
| New model vendor | New Adapter only |
| 10× traffic | Cache retrieval; queue non-interactive jobs; scale workers |
| Auto-send | Stronger validator + approval flag |
| Tool calling | Agent loop: model → tool → observe → model (cap max steps) |
| Eval | Store prompt version; thumbs-up/down API; golden set offline |

### What to practice aloud (20–30 min)

Close the doc. Draw boxes from memory. Hit: tenant RAG filter, credits, timeout, SSE, new vendor.

---

## 20. How to practice (4-week plan)

### Week 1 — Foundations
- SOLID + 5 patterns with your own examples  
- Draw class diagrams for LRU + Parking Lot on paper  
- Practice clarifying questions aloud  

### Week 2 — Core problems
Implement or redesign:
- LRU, Rate Limiter, Parking Lot, Pub-Sub  
- For each: write APIs + failure notes  

### Week 3 — More problems + AI + concurrency
- Splitwise, Notification system, Cache client  
- **AI:** LLM interface + RAG sketch + Suggest Reply (§19) once on paper  
- Add mutex/idempotency discussion every time  
- One machine-coding simulation (90 min timer)  

### Week 4 — Interview simulation
- 4 timed discussion LLDs (45 min each) — include **one AI** problem  
- Record yourself; check structure adherence  
- Prepare 2 stories linking designs to past work  

### Practice method (every problem)
1. 5 min questions  
2. 10 min entities  
3. 20 min classes + flow  
4. 10 min API + deepen  
5. Compare with a reference solution (this repo / articles)  
6. Note 3 improvements  

---

## 21. Interview day checklist

**Before**
- [ ] Paper / shared editor ready  
- [ ] Know your opening script  
- [ ] Sleep; don’t cram 10 new patterns  

**During**
- [ ] Ask questions first  
- [ ] State assumptions  
- [ ] Drive use-case by use-case  
- [ ] Prefer composition + interfaces  
- [ ] Mention concurrency & failure at least once  
- [ ] Show how design evolves  

**Avoid**
- Jumping to Kafka/microservices immediately  
- Pattern-dropping without need  
- Coding before the model is clear (discussion rounds)  
- Ignoring edge cases (null, empty, capacity full)  

---

## 22. Problems in this repository

Use these as **hands-on practice** after designing on paper.

### JavaScript (machine-coding style)

| Folder | Focus |
|--------|--------|
| `Database` | Schema, indexing, CRUD |
| `LRU` | Cache eviction |
| `Parkinglot` / `ParkingLot2` | OO parking design |
| `PollingSystem` | Polls & votes |
| `Pub-Sub` | Eventing |
| `Queue` | Queue mechanics |
| `Ratelimiter` / `RateLimiter2` | Limiting strategies |
| `Redis` | Cache + TTL flavor |
| `SearchEngine` | Indexing & ranking |
| `Splitwise` | Expense splitting |
| `UrlShortener` | Short codes + HTTP |
| `PaymentGateway` | Provider abstraction |

### Golang (`Go/*-go`)

Same problems ported to Go — good for Go interviews (interfaces, mutexes, errors).  
Also see `Go/README.md` for Go concurrency interview prep. For a slower, ultra-beginner AI walkthrough used in company-specific prep, see `Go/hiver-lld-prep-simple.md` (same ideas as §15–§19).

**Note:** This repo’s coded LLDs are classic systems (cache, parking, rate limit, etc.). **AI LLD is covered in this README as design guidance (§15, §19)** — practice it on paper/whiteboard; there is no separate `AI-Suggest-go` folder yet.

### Suggested workflow with this repo

1. Hide the code.  
2. Design on paper for 30–40 minutes.  
3. Only then open the folder and compare.  
4. Re-implement a smaller version from scratch timed.

---

## 23. Cheat sheet

### Structure
```text
Clarify → Entities → Classes/Interfaces → Flows → APIs →
Concurrency/Failure → Extend/Trade-offs
```

### SOLID one-liners
- **S:** one reason to change  
- **O:** extend without editing core  
- **L:** subtypes don’t break contracts  
- **I:** small interfaces  
- **D:** depend on abstractions  

### Patterns one-liners
- Strategy = swap algorithm  
- Factory = create cleanly  
- Adapter = wrap vendor  
- Observer = fan-out events  
- Repository = hide DB  

### AI LLD one-liners
- LLM behind `LLMClient` interface (Adapter)  
- RAG = retrieve tenant-scoped chunks → prompt → generate  
- Credits = Reserve before call, Commit/Release after  
- Never trust model output — validate  
- Timeout LLM calls; ticket must work if AI is down  
- Stream (SSE) for UX; queue for bulk/offline jobs  

### Production questions (ask yourself)
- What if traffic ×10?  
- What if data ×10?  
- What if dependency fails?  
- What if two requests hit same row?  
- What if we add a new variant?  
- How do we observe this in prod?  

### Closing line
> “I kept v1 simple, made extension points where variation is real, and called out concurrency and failure modes. Happy to go deeper on any part.”

---

## Final note

LLD mastery is not memorizing 50 class diagrams. It is building a **repeatable thinking process**, practicing it on 10–15 problems out loud, and explaining **why** your design is good enough for v1 and how it will evolve.

Start with **LRU → Rate Limiter → Parking Lot → Splitwise → Pub-Sub**, then **AI Suggest Reply (§19)**, then branch into domain problems for your target company (payments, inbox, copilots, etc.).

Good luck.
