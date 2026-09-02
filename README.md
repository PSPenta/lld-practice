# Low-Level Design (LLD) Interview Preparation Guide

> A practical guide for **LLD / machine-coding / OOD** rounds at product companies.  
> Working **JavaScript** implementations under **`JavaScript/`**; **Golang** ports under **`Go/`** — use after you learn the method.  
> **AI code review** rounds → **[ai-code-review-round/README.md](ai-code-review-round/README.md)**.

**Do not jump straight into LLD.** Revise language/runtime and REST first, then use the sections below. Depth for each topic lives in **`docs/`** (linked throughout) — this page is the **table of contents + cheat cards + master problem list**.

| Prep guide | Path | Covers |
|------------|------|--------|
| **JavaScript & Node.js** | **[JavaScript/README.md](JavaScript/README.md)** | `this`, closures, event loop, Promises, streams, auth, gotchas |
| **TypeScript** | **[TypeScript/README.md](TypeScript/README.md)** | type system, generics, unions, `strict` tsconfig, LLD interfaces |
| **Golang** | **[Go/README.md](Go/README.md)** | GMP, goroutines, channels, context, GC, interfaces, high-throughput HTTP |
| **REST API & Backend** | **[Backend/README.md](Backend/README.md)** | REST, idempotency, payments, pagination, JWT/OAuth, status codes |
| **LLD method & theory** | **This doc §1–§22** · **[docs/](docs/README.md)** | Method, OOP, SOLID, patterns, walkthroughs in **`problems/`** |
| **LLD gaps (breadth)** | **[lld-gaps/README.md](lld-gaps/README.md)** | UML, extra patterns, 12-week plan |
| **AI code review** | **[ai-code-review-round/README.md](ai-code-review-round/README.md)** | RAG repo review, production bugs |
| **Vibe coding** | **[vibe-coding-round/README.md](vibe-coding-round/README.md)** | Cursor/Copilot rounds — design first, AI second |
| **External** | [awesome-low-level-design](https://github.com/ashishps1/awesome-low-level-design) · [Coder Army LLD playlist](https://youtube.com/playlist?list=PLQEaRBV9gAFvzp6XhcNFpk1WdOcyVo9qT) | Problem index · video supplement (Java/C++) |

| Your focus | Read first |
|------------|------------|
| Node / full-stack | **JavaScript** → **Backend** → §1+ |
| TypeScript / Nest | **JavaScript** → **TypeScript** → **Backend** |
| Go backend | **Go** → **Backend** → §1+ |
| Fintech / payments | **Backend** + **Go** + `Go/PaymentGateway-go/` |
| LLD-only discussion | Skim JS/Go/Backend cheat sheets; still know event loop + idempotency |

**Recommended order:** JavaScript (2–3d) → TypeScript if TS role → Go if Go role → Backend (1–2d) → **§4 round type** + **§5 method** → **[§21 walkthrough](#21-lld-problem-checklist)** → code if ✅ → **lld-gaps/** for UML & drills.

---

## Table of contents

1. [What is LLD?](#1-what-is-lld) · 2. [What interviewers evaluate](#2-what-interviewers-evaluate) · 3. [LLD vs HLD vs DSA](#3-lld-vs-hld-vs-dsa)  
4. [How a typical LLD round runs](#4-how-a-typical-lld-round-runs) · 5. [Standard approach](#5-the-standard-approach-memorize-this) · 6. [Clarifying questions](#6-clarifying-questions-checklist)  
7. [OOP building blocks](#7-oop-building-blocks) · [7A. Repository map](#7a-repository-map--oop-principles--patterns-by-lld)  
8. [Design principles](#8-design-principles-solid--dry--kiss--yagni--polk) · 9. [Design patterns](#9-design-patterns-creational--behavioural--structural)  
10. [API design](#10-api-design-for-lld) · 11. [Data modeling & state](#11-data-modeling--state) · 12. [Concurrency, idempotency & failure](#12-concurrency-idempotency--failure)  
13. [Extensibility & evolution](#13-extensibility--evolution) · 14. [Common LLD problems — how to think](#14-common-lld-problems--how-to-think)  
15. [AI / LLM LLD](#15-ai--llm-lld-for-beginners) · 16. [Worked examples — design docs](#16-worked-examples--design-docs)  
17. [HLD bleed-in](#17-hld-topics-that-bleed-into-lld) · 18. [Timed mock](#18-timed-mock--self-score) · 19. [4-week plan](#19-how-to-practice-4-week-plan) · 20. [Interview day](#20-interview-day-checklist)  
21. [LLD problem checklist (master list)](#21-lld-problem-checklist) · 22. [Cheat sheet](#22-cheat-sheet)

---

## 1. What is LLD?

**Low-Level Design** = classes, modules, APIs — not AWS topology (HLD).

| Question you answer | |
|---------------------|---|
| Main **entities**? | |
| **Classes / interfaces** and methods? | |
| **Collaboration** per use case? | |
| **API** shape? | |
| **Clean, testable, extensible** structure? | |
| **Concurrency, failures, growth**? | |

**Discussion LLD** — whiteboard, 60–75 min, trade-offs. **Machine coding** — working code, 90–120 min. Same fundamentals.

→ Full detail: **[docs/method/README.md §1](docs/method/README.md#1-what-is-lld)**

---

## 2. What interviewers evaluate

Problem breakdown · Modeling · Abstraction · SOLID · Patterns (purposeful) · APIs · Data & state · Concurrency · Failure handling · Extensibility · Communication.

→ **[docs/method/README.md §2](docs/method/README.md#2-what-interviewers-evaluate)**

---

## 3. LLD vs HLD vs DSA

| | DSA | LLD | HLD |
|---|-----|-----|-----|
| Focus | Algorithms, Big-O | Classes, APIs | Scale, infra |
| Example | Shortest path | Rate Limiter classes | URL shortener at 100M QPS |

→ **[docs/method/README.md §3](docs/method/README.md#3-lld-vs-hld-vs-dsa)**

---

## 4. How a typical LLD round runs

Same **six steps** ([§5](#5-the-standard-approach-memorize-this)) — different **time budget**. See both tables below.

### Discussion LLD (~60 min)

| Min | Activity |
|-----|----------|
| 0–5 | Clarify requirements & assumptions |
| 5–12 | Entities + relationships |
| 12–30 | Classes/interfaces + main flows |
| 30–40 | APIs (+ light schema if asked) |
| 40–50 | Concurrency, failures, extensibility |
| 50–60 | Trade-offs, evolve design, Q&A |

### Machine coding (~90–120 min)

| Min | Activity |
|-----|----------|
| 0–10 | Clarify + lock v1 scope (what you skip) |
| 10–25 | Entities + classes + one flow (design only) |
| 25–85 | Code happy path + core classes |
| 85–105 | Public API, 2–3 edge cases, basic thread-safety |
| 105–120 | Trade-offs + one evolution out loud |

**Other rounds:** Vibe coding · AI code review · Language Q&A — **[docs/method/ §4](docs/method/README.md#4-how-a-typical-lld-round-runs)**.

**Problem walkthroughs:** one doc per checklist row under **`problems/<name>/`** — see [§21](#21-lld-problem-checklist).

---

## 5. The standard approach (memorize this)

```text
Understand → Model → Design → APIs → Deepen → Evolve
```

1. **Understand** — actors, use cases, constraints  
2. **Model** — entities, relationships, state machines  
3. **Design** — classes, interfaces, responsibilities  
4. **APIs** — endpoints or public methods  
5. **Deepen** — concurrency, idempotency, errors  
6. **Evolve** — “If scale 10× or requirement changes…”

→ **[docs/method/README.md §5–§6](docs/method/README.md#5-the-standard-approach-memorize-this)**

---

## 6. Clarifying questions checklist

**Product:** actors, top 3 use cases, out of scope · **Scale:** users, QPS, single machine vs distributed · **Functional:** auth, sync/async, persistence, TTL · **Quality:** consistency, failures, observability — then **state assumptions aloud**.

→ Full list: **[docs/method/README.md §6](docs/method/README.md#6-clarifying-questions-checklist)**

---

## 7. OOP building blocks

| Concept | One line | Deep dive |
|---------|----------|-----------|
| **Class** | Blueprint — data + behavior | [docs/oop/](docs/oop/README.md) |
| **Encapsulation** | Hide internals; small public API | [docs/oop/](docs/oop/README.md) · `JavaScript/LRU/`, `Database/Table.js` |
| **Abstraction** | Contract without implementation | [docs/oop/](docs/oop/README.md) · `RateLimiterStrategy`, `BankGateway` |
| **Polymorphism** | Same interface, different runtime behavior | [docs/oop/](docs/oop/README.md) · Strategy dispatch, expense `validate()` |
| **Inheritance (is-a)** | Specialized subtype extends base | [docs/oop/](docs/oop/README.md) · `Vehicle` → Car/Bike/Truck |
| **Composition (has-a)** | Build from parts + interfaces — **LLD default** | [docs/oop/](docs/oop/README.md) · `RateLimiter2`, `ParkingLot2`, `SearchEngine` |
| **Association / Aggregation / Composition** | uses · has · owns (lifecycle) | [docs/oop/](docs/oop/README.md) |

**Interview line:** Composition + interfaces over deep inheritance. **Repo proof:** RateLimiter2, ParkingLot2, SearchEngine, PollingSystem2.

→ Full treatment + repo tables: **[docs/oop/README.md](docs/oop/README.md)**

---

## 7A. Repository map — OOP, principles & patterns by LLD

Point at real code in interviews — **Strategy** → `RateLimiter2` · **Factory** → `Splitwise` / `ParkingLot2` · **Observer** → `Pub-Sub` · **DIP** → `PaymentGateway-go` · **Service + repos** → `PollingSystem2`.

→ Full matrix (OOP pillars, SOLID per folder, pattern table): **[docs/repo-map/README.md](docs/repo-map/README.md)**

---

## 8. Design principles (SOLID + DRY / KISS / YAGNI / PoLK)

| Principle | One line | Deep dive |
|-----------|----------|-----------|
| **S — Single Responsibility** | One reason to change | [docs/principles/](docs/principles/README.md) |
| **O — Open/Closed** | Extend with new classes, not editing core | [docs/principles/](docs/principles/README.md) |
| **L — Liskov Substitution** | Subtypes honor the contract | [docs/principles/](docs/principles/README.md) |
| **I — Interface Segregation** | Small interfaces | [docs/principles/](docs/principles/README.md) |
| **D — Dependency Inversion** | Depend on abstractions | [docs/principles/](docs/principles/README.md) |
| **DRY** | Don’t duplicate knowledge | [docs/principles/](docs/principles/README.md) |
| **KISS** | Simplest workable design | [docs/principles/](docs/principles/README.md) |
| **YAGNI** | Don’t build unused features | [docs/principles/](docs/principles/README.md) |
| **PoLK** | Talk only to direct collaborators | [docs/principles/](docs/principles/README.md) |

### Principles cheat card

| Principle | One line |
|-----------|----------|
| SRP | One reason to change |
| OCP | Extend without editing core |
| LSP | Impls must honor the contract |
| ISP | Small interfaces |
| DIP | Depend on abstractions |
| DRY | Don’t duplicate knowledge |
| KISS | Simplest workable design |
| YAGNI | Don’t build unused features |
| PoLK | Talk only to close friends |

→ Interview phrasing + repo examples: **[docs/principles/README.md](docs/principles/README.md)**

---

## 9. Design patterns (Creational · Behavioural · Structural)

**Rule:** Use when it solves variation or decoupling — not decoration. `*` = useful · `**` = very common in LLD.

### GoF groups

| Category | Question it answers |
|----------|---------------------|
| **Creational** | How do we **create** objects cleanly? |
| **Behavioural** | How do objects **communicate / vary behavior**? |
| **Structural** | How do we **compose** objects for flexibility? |

### Pattern index (detail in doc)

| Creational | Behavioural | Structural |
|------------|-------------|------------|
| Singleton `**` | Observer `**` (Pub-Sub) | Adapter `**` |
| Factory `*` | Strategy `**` | Proxy `**` |
| Builder `*` | Iterator `*` | Decorator `*` |
| Abstract Factory | Command `*` | Facade `*` |
| Prototype | State · Template Method | Composite · Bridge |

**In this repo (named):** Strategy → `RateLimiter2`, `PaymentGateway-go` · Factory → `Splitwise`, `ParkingLot2` · Observer → `Pub-Sub` · Facade-like → `SearchEngine`.  
**Paper only:** Singleton, Builder, Command, State, Proxy, Decorator, Adapter — **[lld-gaps/ §4](lld-gaps/README.md)**.

### Pattern ↔ principle quick links

| You say… | You’re applying… |
|----------|------------------|
| New class instead of editing switch | OCP + often Strategy/Factory |
| Depend on `LLMClient` interface | DIP + Adapter |
| Split God class | SRP |
| Don’t build Abstract Factory yet | YAGNI + KISS |
| Middleware wraps handler | Decorator |

→ Full pattern write-ups + repo files: **[docs/patterns/README.md](docs/patterns/README.md)**

---

## 10. API design for LLD

Resource names · HTTP verbs (POST/PUT/PATCH) · versioning · error model + `trace_id` · cursor pagination · idempotency keys · validation (422) · rate limit (429).

→ REST depth: **[Backend/README.md](Backend/README.md)** · LLD slice: **[docs/design-topics/ §10](docs/design-topics/README.md#10-api-design-for-lld)**

---

## 11. Data modeling & state

Entities & relationships (1—*, composition chains) · **state machines** (ticket, order, payment) · persistence choices (in-memory vs DB, indexes, unique keys) · derived vs stored fields.

→ **[docs/design-topics/ §11](docs/design-topics/README.md#11-data-modeling--state)**

---

## 12. Concurrency, idempotency & failure

**Concurrency:** locks / channels · optimistic locking · don’t hold locks during I/O · **Idempotency:** `Idempotency-Key` · **Failure:** timeouts, backoff, partial failure + queue/outbox.

→ **[docs/design-topics/ §12](docs/design-topics/README.md#12-concurrency-idempotency--failure)**

---

## 13. Extensibility & evolution

New payment provider → Adapter + Factory · new rate limit algo → Strategy · new notify channel → Notifier interface · 10× traffic → cache, queue, partition · multi-tenant → `tenant_id` everywhere.

→ **[docs/design-topics/ §13](docs/design-topics/README.md#13-extensibility--evolution)**

---

## 14. Common LLD problems — how to think

Actors → use cases → entities → classes → APIs → concurrency → extend.

### Must-practice set

| Problem | Core ideas |
|---------|------------|
| **LRU Cache** | Hash map + DLL; O(1) get/put |
| **Rate Limiter** | Strategy; per-key state; thread safety |
| **Parking Lot** | Floors/slots/vehicles; ticket |
| **Splitwise** | Split strategies; balances |
| **Pub-Sub** | Topics; fan-out |
| **URL Shortener** | encode; mappings; APIs |
| **Elevator / Traffic** | State machine; scheduling |
| **Chess / Snake & Ladder** | Board; rules |
| **Logging** | Levels; appenders; chain of responsibility |
| **ATM / Vending** | State + chain for dispense |
| **Notification** | Templates; channels; retry |
| **Booking** | Seats; locks; idempotency |
| **Cache client** | TTL/LRU; loader; stampede |
| **Job scheduler** | Priority queue; workers; retries |

→ AI-era variants: **[docs/ai-lld/](docs/ai-lld/README.md)** · **[docs/design-topics/ §14](docs/design-topics/README.md#14-common-lld-problems--how-to-think)**

---

## 15. AI / LLM LLD for beginners

`LLMClient` interface · RAG (tenant-scoped retrieve → prompt → generate) · credits (reserve/commit) · guardrails · streaming · ticket works if AI is down.

→ **[docs/ai-lld/README.md](docs/ai-lld/README.md)** · walkthroughs: **[Cache Client](problems/cache-client/README.md)**, **[AI Suggest Reply](problems/ai-suggest-reply/README.md)**

---

## 16. Worked examples — design docs

Every checklist row has a **walkthrough** under **`problems/<name>/`** — Steps 1–6, timed per [§4](#4-how-a-typical-lld-round-runs). Hide code → walkthrough → compare with **`JavaScript/*/`** if ✅.

---

## 17. HLD topics that bleed into LLD

v1 sync → Redis → queue → stateless scale · caching · multi-tenant `tenant_id` · `trace_id` · rate limits / credits.

→ **[docs/interview-prep/ §17](docs/interview-prep/README.md#17-hld-topics-that-bleed-into-lld)**

---

## 18. Timed mock + self-score

40 min · Cache Client or Rate Limiter · score 10 criteria (0–2 each, target ≥16/20).

→ **[docs/interview-prep/ §18](docs/interview-prep/README.md#18-timed-mock--self-score)**

---

## 19. How to practice (4-week plan)

Week 0 language · Week 1 SOLID + diagrams · Week 2 core coded problems · Week 3 concurrency + AI · Week 4 mocks + ❌ paper gaps.

→ **[docs/interview-prep/ §19](docs/interview-prep/README.md#19-how-to-practice-4-week-plan)**

---

## 20. Interview day checklist

Before / during / avoid lists · company research template · LLD vs code-review round.

→ **[docs/interview-prep/ §20](docs/interview-prep/README.md#20-interview-day-checklist)**

---

## 21. LLD problem checklist

**How to use:** Open **Walkthrough** → follow six steps + [§4](#4-how-a-typical-lld-round-runs) time boxes for your round type → code only if ✅.

| Legend | Meaning |
|--------|---------|
| ✅ | Working code in `JavaScript/` and/or `Go/` (compare after design) |
| ❌ | Paper / discussion round — walkthrough only |

### Master list

| Problem | Solved | Walkthrough (§4 pattern) |
|---------|:------:|--------------------------|
| **API gateway** | ❌ | [problems/api-gateway/](problems/api-gateway/README.md) |
| **ATM** | ❌ | [problems/atm/](problems/atm/README.md) |
| **Cab booking** | ❌ | [problems/cab-booking/](problems/cab-booking/README.md) |
| **Cache (TTL + eviction)** | ✅ | [problems/cache-client/](problems/cache-client/README.md) |
| **Circuit breaker** | ❌ | [problems/circuit-breaker/](problems/circuit-breaker/README.md) |
| **Connection pool** | ❌ | [problems/connection-pool/](problems/connection-pool/README.md) |
| **Elevator** | ❌ | [problems/elevator/](problems/elevator/README.md) |
| **Expense splitter (Splitwise)** | ✅ | [problems/splitwise/](problems/splitwise/README.md) |
| **File storage (in-memory FS)** | ❌ | [problems/file-storage/](problems/file-storage/README.md) |
| **In-memory database** | ✅ | [problems/in-memory-database/](problems/in-memory-database/README.md) |
| **Inventory management** | ❌ | [problems/inventory-management/](problems/inventory-management/README.md) |
| **Job scheduler** | ❌ | [problems/job-scheduler/](problems/job-scheduler/README.md) |
| **LFU cache** | ❌ | [problems/lfu-cache/](problems/lfu-cache/README.md) |
| **Load balancer** | ❌ | [problems/load-balancer/](problems/load-balancer/README.md) |
| **Logging framework** | ❌ | [problems/logging-framework/](problems/logging-framework/README.md) |
| **LRU cache** | ✅ | [problems/lru-cache/](problems/lru-cache/README.md) |
| **Message queue (ack, DLQ)** | ❌ | [problems/message-queue/](problems/message-queue/README.md) |
| **Metrics collector** | ❌ | [problems/metrics-collector/](problems/metrics-collector/README.md) |
| **Movie / seat booking** | ❌ | [problems/movie-booking/](problems/movie-booking/README.md) |
| **Notification service** | ❌ | [problems/ticket-notify/](problems/ticket-notify/README.md) |
| **Order management (OMS)** | ❌ | [problems/order-management/](problems/order-management/README.md) |
| **Parking lot** | ✅ | [problems/parking-lot/](problems/parking-lot/README.md) |
| **Payment gateway** | ✅ | [problems/payment-gateway/](problems/payment-gateway/README.md) |
| **Polling / voting system** | ✅ | [problems/polling-system/](problems/polling-system/README.md) |
| **Pub/Sub system** | ✅ | [problems/pub-sub/](problems/pub-sub/README.md) |
| **Rate limiter** | ✅ | [problems/rate-limiter/](problems/rate-limiter/README.md) |
| **Restaurant table reservation** | ❌ | [problems/restaurant-reservation/](problems/restaurant-reservation/README.md) |
| **Retry scheduler** | ❌ | [problems/retry-scheduler/](problems/retry-scheduler/README.md) |
| **Search engine** | ✅ | [problems/search-engine/](problems/search-engine/README.md) |
| **Subscription manager** | ❌ | [problems/subscription-manager/](problems/subscription-manager/README.md) |
| **Task board (Trello-like)** | ❌ | [problems/task-board/](problems/task-board/README.md) |
| **Task queue (FIFO)** | ✅ | [problems/task-queue/](problems/task-queue/README.md) |
| **URL shortener** | ✅ | [problems/url-shortener/](problems/url-shortener/README.md) |
| **Vending machine** | ❌ | [problems/vending-machine/](problems/vending-machine/README.md) |
| **Wallet / ledger** | ❌ | [problems/wallet-ledger/](problems/wallet-ledger/README.md) |
| **Webhook delivery system** | ❌ | [problems/webhook-delivery/](problems/webhook-delivery/README.md) |
| **AI suggest reply / copilot** | ❌ | [problems/ai-suggest-reply/](problems/ai-suggest-reply/README.md) |
| **Distributed rate limiter** | ❌ | [problems/distributed-rate-limiter/](problems/distributed-rate-limiter/README.md) |
| **Chess / board game** | ❌ | [problems/chess/](problems/chess/README.md) |
| **Hotel booking** | ❌ | [problems/hotel-booking/](problems/hotel-booking/README.md) |
| **Traffic light / signal** | ❌ | [problems/traffic-signal/](problems/traffic-signal/README.md) |

**Scorecard:** **12 ✅** · **29 ❌** (41 total). **Coded impl map:** [docs/repo-map/](docs/repo-map/README.md). **Browse by difficulty:** [awesome-low-level-design](https://github.com/ashishps1/awesome-low-level-design).

---

## 22. Cheat sheet

### Structure

```text
Clarify → Entities → Classes/Interfaces → Flows → APIs →
Concurrency/Failure → Extend/Trade-offs
```

### Principles (one-liners)

SRP · OCP · LSP · ISP · DIP · **Composition** over deep is-a · DRY · KISS · YAGNI · PoLK — full card in [§8](#8-design-principles-solid--dry--kiss--yagni--polk).

### Patterns (one-liners)

**Creational:** Singleton · Factory · Builder · Abstract Factory · Prototype  
**Behavioural:** Observer · Strategy · Iterator · Command · State · Template Method  
**Structural:** Adapter · Proxy · Decorator · Facade · Composite · Bridge  
**Repo demos:** [§7A](#7a-repository-map--oop-principles--patterns-by-lld) · [docs/repo-map/](docs/repo-map/README.md)

### AI LLD

`LLMClient` · RAG · credits reserve/commit · validate output · timeout · stream vs queue · [ai-code-review-round/](ai-code-review-round/README.md)

### Must-practice designs

[Cache Client](problems/cache-client/README.md) · [Rate Limiter](problems/rate-limiter/README.md) · [AI Suggest Reply](problems/ai-suggest-reply/README.md)

### Closing line

> “I kept v1 simple, made extension points where variation is real, and called out concurrency and failure modes. Happy to go deeper on any part.”

→ Full cheat sheet + production questions: **[docs/interview-prep/ §22](docs/interview-prep/README.md#22-cheat-sheet)**

---

## Final note

LLD mastery is a **repeatable process**, not 50 memorized diagrams. Practice 10–15 problems out loud.

**Start:** [Cache Client](problems/cache-client/README.md) → [Rate Limiter](problems/rate-limiter/README.md) → [Parking Lot](problems/parking-lot/README.md) → [Splitwise](problems/splitwise/README.md) → [Pub-Sub](problems/pub-sub/README.md) → [AI Suggest Reply](problems/ai-suggest-reply/README.md). **Code review:** [ai-code-review-round/README.md](ai-code-review-round/README.md).

Good luck.
