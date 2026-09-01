# LLD Gaps — What This Repo Does Not Cover (Yet)

> **Use with:** [../README.md](../README.md) (method + worked examples + implemented problems) · [../JavaScript/README.md](../JavaScript/README.md) · [../Go/README.md](../Go/README.md)  
> **Based on:** common FAANG-style LLD syllabi (e.g. [ChatGPT LLD Roadmap share](https://chatgpt.com/share/6a949bd1-856c-83ee-b35b-b1b5badaaa91)) — this doc fills **gaps**, not duplicates what the main README already teaches.  
> **Audience:** SDE-2+ switching — you know Parking Lot and Rate Limiter from this repo; you still need breadth on paper.

---

## Table of contents

1. [How to use this doc](#1-how-to-use-this-doc)
2. [Coverage map — main README vs gaps](#2-coverage-map--main-readme-vs-gaps)
3. [UML — minimum for interviews](#3-uml--minimum-for-interviews)
4. [Patterns — study on paper (not in repo code)](#4-patterns--study-on-paper-not-in-repo-code)
5. [Problem ladder — missing implementations](#5-problem-ladder--missing-implementations)
6. [Paper-design prompts (copy-paste practice)](#6-paper-design-prompts-copy-paste-practice)
7. [Testing & testable design](#7-testing--testable-design)
8. [Refactoring drills](#8-refactoring-drills)
9. [Persistence & layered architecture](#9-persistence--layered-architecture)
10. [FAANG / senior combo topics](#10-faang--senior-combo-topics)
11. [Extended 12-week plan (optional)](#11-extended-12-week-plan-optional)
12. [Gap checklist — print before interview](#12-gap-checklist--print-before-interview)

---

## 1. How to use this doc

```text
Main README     → HOW to design (method, SOLID, patterns in repo, worked examples)
JavaScript/Go   → Language + runtime
Repo code       → ~12 implemented LLDs (depth)
This doc        → BREADTH you still need on paper / whiteboard
```

**Rule:** For each gap item, spend **30–45 min on paper** using [../README.md](../README.md) §5 method — you do **not** need working code for every problem.

---

## 2. Coverage map — main README vs gaps

| Topic | In main README / repo? | This doc |
|-------|------------------------|----------|
| LLD method (clarify → entities → APIs) | ✅ §4–§6 | — |
| OOP, SOLID, DRY/KISS/YAGNI/PoLK | ✅ §7–§8 | — |
| Strategy, Factory, Observer, Template | ✅ coded + §7A | More patterns §4 |
| Parking Lot, Rate Limiter, LRU, Splitwise, Pub-Sub, DB, URL, Search, **PollingSystem2** | ✅ `JavaScript/PollingSystem2/PollingService/` · models + repos | — |
| Payment pluggable gateway | ✅ `Go/PaymentGateway-go/` | — |
| AI LLD (Cache Client, Suggest Reply, RAG) | ✅ §15–§21 | — |
| **UML** (class, sequence, state) | ⚠️ light | **§3** |
| **Beginner games** (Chess, TTT, Snake) | ⚠️ §14 mention | **§5 L1** |
| **Elevator, Notification, Logger, File system** | ⚠️ mention only | **§5 L2** |
| **Workflow, Rule engine, Scheduler, LFU** | ❌ | **§5 L3–§10** |
| **Singleton, Builder, State, Decorator** (hands-on) | ❌ coded | **§4** |
| **Testing / mocking for LLD** | ⚠️ Go README | **§7** |
| **Refactoring god-class drills** | ❌ | **§8** |
| **Repository / DAO / Unit of Work** | ✅ PollingSystem2 (`PollingService/repositories/` + service) | **§9** (generalize + Unit of Work) |
| **12-week breadth timeline** | 4-week in §25 | **§11** |

---

## 3. UML — minimum for interviews

You do **not** need full UML certification. You need **three diagrams fast**.

### 3.1 Class diagram (most common)

```text
┌─────────────┐       ┌──────────────┐
│ ParkingLot  │1    * │    Floor     │
├─────────────┤───────├──────────────┤
│ +park(v)    │       │ +findSlot()  │
│ +unpark(t)  │       └──────┬───────┘
└─────────────┘              │1
                             │*
                      ┌──────▼───────┐
                      │     Slot     │
                      └──────────────┘

<<interface>>
RateLimiterStrategy
      △
      │ implements
 TokenBucket   FixedWindow
```

**Draw when:** after listing entities (minute 10–15 of LLD round).

**Relationships to label:**

| Symbol | Meaning | Example |
|--------|---------|---------|
| ──▷ dashed | implements | `TokenBucket` implements `RateLimiterStrategy` |
| ──▶ solid | inherits | `Car` extends `Vehicle` |
| ◆── composition | lifecycle bound | `ParkingLot` owns `Floor` |
| ◇── aggregation | has-a, weaker | `Order` has `LineItems` |
| ── uses | dependency | `OrderService` uses `PaymentGateway` |

### 3.2 Sequence diagram (one happy path)

```text
Client    ParkingLot    Floor    Slot    Ticket
  │           │           │        │        │
  │──park()──►│           │        │        │
  │           │─findSlot()►│        │        │
  │           │           │─free?─►│        │
  │           │◄──────────│◄───────│        │
  │           │──createTicket────────────────►│
  │◄──ticket──│           │        │        │
```

**Draw when:** interviewer asks “walk through park flow” or multi-actor (notify, payment).

### 3.3 State diagram (Elevator, ATM, Order, Ticket)

```text
[Created] ──pay──► [Paid] ──ship──► [Shipped] ──deliver──► [Delivered]
    │                  │
    └──cancel──────────┴──cancel──► [Cancelled]
```

**Use for:** Elevator, vending machine, workflow engine, notification retry states.

**Interview line:** “I'll use a state diagram here because behavior depends on **state**, not only class type.”

---

## 4. Patterns — study on paper (not in repo code)

Main README §9 lists what **is** in the repo. Master these **additionally on paper**:

| Pattern | Problem → naive fail → pattern | When NOT to use |
|---------|----------------------------------|-----------------|
| **Singleton** | One config/connection — multiple instances break | Prefer DI in Go/Node; hard to test |
| **Builder** | HTTP request / pizza / query with 10 optional fields | Small constructors (< 4 params) |
| **Abstract Factory** | UI themes / DB families (MySQL vs Postgres widgets) | One product line only |
| **State** | Order/Elevator — giant `if (status)` | Few states, stable |
| **Command** | Undo/redo, job queue, macro actions | Simple CRUD only |
| **Chain of Responsibility** | Logging levels, auth middleware chain | Fixed pipeline known upfront |
| **Decorator** | Coffee + milk + whip; middleware stack | Deep nesting hurts debug |
| **Adapter** | Third-party payment API → your `PaymentProcessor` | Same interface already |
| **Proxy** | Lazy load, access control, caching wrapper | Premature — start direct |
| **Facade** | SearchEngine orchestrates tokenizer+index+ranker | Already have Facade-like in repo |

**Drill (15 min each):** Pick pattern → draw bad design → refactor diagram → say tradeoff aloud.

---

## 5. Problem ladder — missing implementations

### Level 1 — Beginner (warmup, entity modeling)

| # | Problem | Core skill | In repo? |
|---|---------|------------|----------|
| 1 | **Tic Tac Toe** | Board, win check, 2 players | ❌ paper |
| 2 | **Snake & Ladder** | Board, dice, jump rules | ❌ paper |
| 3 | **Deck of Cards** | Suit, rank, shuffle, deal | ❌ paper |
| 4 | **Vending Machine** | State + coin inventory + dispense | ❌ paper |
| 5 | **Coffee Machine** | Recipe, ingredients, State | ❌ paper |
| 6 | **ATM** | Card, PIN, cash dispenser, State | ❌ paper |
| 7 | **Library Management** | Book, member, loan, fine | ❌ paper |
| 8 | **Movie ticket booking** | Seat lock, show, payment idempotency | ❌ paper |
| 9 | **Restaurant / food order** | Menu, order, kitchen queue | ❌ paper |
| 10 | **Parking Lot** | — | ✅ `JavaScript/ParkingLot2/` |

### Level 2 — Intermediate

| # | Problem | Core skill | In repo? |
|---|---------|------------|----------|
| 1 | **Elevator** | State, scheduling (SCAN/FCFS), multiple cars | ❌ paper |
| 2 | **Notification system** | Channels, templates, retry, idempotency | ⚠️ §21 only |
| 3 | **Logging framework** | Levels, appenders, Chain of Responsibility | ❌ paper |
| 4 | **File system (in-memory)** | Tree, path resolve, permissions | ❌ paper |
| 5 | **Hotel / Cab booking** | Inventory, locks, cancellation | ❌ paper |
| 6 | **Inventory management** | SKU, reservations, concurrent updates | ❌ paper |
| 7 | **Order management** | State machine, payment, fulfillment | ❌ paper |
| 8 | **Splitwise** | — | ✅ coded |
| 9 | **Rate Limiter** | — | ✅ coded |
| 10 | **Car rental** | Vehicle pool, booking, late return | ❌ paper |

### Level 3 — Advanced

| # | Problem | Core skill | In repo? |
|---|---------|------------|----------|
| 1 | **Chess** | Board, pieces, move validation, check | ❌ paper |
| 2 | **LFU cache** | Freq map + min heap / dual structure | ❌ (LRU ✅) |
| 3 | **Task / Job scheduler** | Priority queue, workers, retry | ❌ paper |
| 4 | **Message queue** | Producers, consumers, ack, DLQ | ❌ paper |
| 5 | **Pub/Sub** | — | ✅ coded |
| 6 | **Payment system** | — | ✅ Go only |
| 7 | **Connection pool** | Acquire, release, max, timeout | ❌ paper |
| 8 | **In-memory DB** | — | ✅ coded |
| 9 | **Event processing** | ingest → process → store | ❌ paper |
| 10 | **Cache (TTL + eviction)** | — | ✅ Redis/LRU + §16 Cache Client |

### Level 4 — Senior / FAANG combo

| Problem | Combines |
|---------|----------|
| **Extensible notification platform** | Observer + Strategy (channel) + retry + templates |
| **Pluggable payment** | Strategy + Factory + idempotency |
| **Distributed rate limiter** | Rate limit + shared store + clock skew |
| **Workflow engine** | State + Command + DAG |
| **Rule engine** | Strategy + chain + config-driven rules |
| **Plugin architecture** | Open/Closed + registry + interface |

See **§10** for how to practice these without coding everything.

---

## 6. Paper-design prompts (copy-paste practice)

Use [../README.md](../README.md) §5 for each. **45 min timer.**

### Elevator (45 min)

> Design an elevator system for a building: multiple elevators, multiple floors, up/down requests. Support morning rush. Discuss scheduling strategy.

**Hit:** `ElevatorController`, `Elevator` (state: IDLE/MOVING/DOORS_OPEN), `Request` queue, strategy interface (FCFS vs SCAN).

### Notification system (45 min)

> Design in-app + email + SMS notifications for a SaaS product. Templates per event. Retry on failure. User opt-out.

**Hit:** `NotificationService`, `Channel` interface, `TemplateRenderer`, `Outbox` or queue, idempotency key.

### Logging framework (30 min)

> Design a logger with levels DEBUG–ERROR, multiple appenders (console, file), formatters.

**Hit:** Chain of Responsibility or composite appender; don't over-engineer Singleton.

### Vending machine (30 min)

> Accept coins, select product, dispense, return change. Handle out-of-stock.

**Hit:** **State pattern** — Idle, HasMoney, Dispensing, ReturnChange.

### Workflow engine (60 min — senior)

> Orders go through: created → paid → packed → shipped → delivered. Allow cancel before ship. Add refund path later without rewriting core.

**Hit:** State per order; transitions table; Open/Closed for new states.

---

## 7. Testing & testable design

LLD interviews ask: *“How would you test this?”*

| Design choice | Makes testing easy |
|---------------|-------------------|
| **DIP** — depend on interfaces | Inject fake `PaymentGateway`, fake `Notifier` |
| **Pure functions** for rules | Test split logic without DB |
| **No static singletons** | Pass dependencies in constructor |
| **Small public API** | Few integration tests cover surface |

### Mock vs stub vs fake (one-liners)

| Term | Meaning |
|------|---------|
| **Stub** | Returns canned data (`fakeRepo.getUser()` → fixed user) |
| **Mock** | Verifies calls were made (`expect(notifier).toHaveBeenCalledWith(...)`) |
| **Fake** | Working in-memory impl (`FakePaymentGateway` records charges) |

**Interview line:** “I'd unit-test `Expense.validate()` and split math with fakes; integration-test `ParkingLot.park()` end-to-end with an in-memory store.”

**Hands-on:** see [../Go/README.md](../Go/README.md) §11 (table-driven tests, mock interfaces).

---

## 8. Refactoring drills

Take **intentionally bad** designs and refactor on paper (15 min each).

### God class

```text
BEFORE: ApplicationManager — park(), sendEmail(), chargeCard(), log(), saveToDb()

AFTER:  ParkingLot + NotificationService + PaymentService + Logger
        (each injected into orchestrator)
```

### Deep inheritance

```text
BEFORE: Animal → Mammal → Dog → WorkingDog

AFTER:  Dog has-a TrainableBehaviour, WorkerRole (composition)
```

### Feature envy

```text
BEFORE: OrderService reaches into order.items[i].product.price

AFTER:  order.getTotal() or OrderLine.getSubtotal()
```

**Drill:** Open [../JavaScript/UrlShortener/](../JavaScript/UrlShortener/) — README calls it SRP gap; sketch how you'd split HTTP vs storage.

---

## 9. Persistence & layered architecture

Standard layers (say aloud in LLD):

```text
Handler / Controller  → HTTP, validation, auth
Service             → use cases, orchestration
Repository          → persistence interface
Entity / Domain     → business objects
```

| Pattern | LLD use |
|---------|---------|
| **Repository** | `UserRepository.findById()` — hide SQL/ORM |
| **DTO** | API shape ≠ domain shape |
| **Unit of Work** | One transaction across multiple repos |

**Interview line:** “Persistence is behind `Repository`; service doesn't know if it's Postgres or in-memory — see **`JavaScript/PollingSystem2/PollingService/`** (`repositories/` + **`PollingService`**) or the `JavaScript/Database/` aggregate.”

**Repo sketch (PollingSystem2):**

```text
index.js               ← demo
PollingService/
  index.js             ← PollingService (use-cases)
  models/              ← User, Poll, Vote
  repositories/        ← UserRepository, PollRepository, VoteRepository
```

Go mirror: `Go/PollingSystem2-go/main.go` · `pollingservice/` · `models/` · `repositories/`

| Poll type | Vote rule | Assign |
|-----------|-----------|--------|
| **Private** (`isPrivate: true`) | Must be on creator's invite list | Creator assigns voters |
| **Public** (default) | Any user except creator | Assign blocked |

No separate `Admin` class — any user may create polls and vote on **others’** polls, never their own. **`isClosed`** (or past `validTill` via `poll.isCompleted()`) stops new votes.

---

## 10. FAANG / senior combo topics

Practice **speaking** these without full code:

| Topic | 3-minute structure |
|-------|-------------------|
| **Distributed rate limiter** | Local counter vs Redis; token bucket; key per user/IP; clock skew |
| **Distributed cache** | Consistent hashing; TTL; stampede → singleflight |
| **Workflow engine** | States + transitions; idempotent handlers; DLQ |
| **Rule engine** | Rules as data; Strategy per rule type; eval pipeline |
| **Plugin architecture** | Registry; interface; load at runtime; version compatibility |

Link to [../README.md](../README.md) §23 (HLD bleed-in) for queues, CAP, idempotency.

---

## 11. Extended 12-week plan (optional)

Use if you want **ChatGPT-roadmap breadth** on top of this repo's **4-week depth**.

| Weeks | Focus |
|-------|-------|
| **1–2** | §0 language prep + main README §7–§8 (OOP/SOLID) + **§3 UML** this doc |
| **3–4** | Main README §9 patterns + **§4 gaps patterns** on paper |
| **5–6** | **§5 Level 1–2** paper designs (2 per week) |
| **7–8** | Repo code compare: paper design → open `JavaScript/*/` |
| **9–10** | **§5 Level 3–4** + concurrency (Go README §1–§4) |
| **11** | **§7–§8** testing + refactoring drills |
| **12** | Timed mocks: main README §24 + [vibe-coding-round](../vibe-coding-round/README.md) if AI-assisted |

---

## 12. Gap checklist — print before interview

**UML**
- [ ] Can draw class diagram with composition vs inheritance
- [ ] Can draw one sequence diagram for happy path
- [ ] Know when to use state diagram

**Patterns (paper)**
- [ ] State, Builder, Command, Decorator, Chain of Responsibility — one example each

**Problems (paper-designed at least once)**
- [ ] One **game** (TTT or Snake)
- [ ] One **state machine** (Elevator or Vending or ATM)
- [ ] **Notification** or **Logger**
- [ ] One **senior combo** (workflow or pluggable payments — beyond repo)

**Engineering**
- [ ] Explain Repository + Service layers
- [ ] Explain how you'd test with mocks/fakes
- [ ] Refactor god-class example aloud

**Already strong from main repo** (don't re-read from zero)
- [ ] Parking Lot, Rate Limiter, LRU, Splitwise, Pub-Sub, Cache Client, AI Suggest Reply

---

## Quick links

| Doc | Purpose |
|-----|---------|
| [../README.md](../README.md) | LLD method + implemented problems |
| [../vibe-coding-round/README.md](../vibe-coding-round/README.md) | AI-assisted coding / build rounds |
| [../ai-code-review-round/README.md](../ai-code-review-round/README.md) | Manual code review rounds |
| [ChatGPT LLD Roadmap](https://chatgpt.com/share/6a949bd1-856c-83ee-b35b-b1b5badaaa91) | Full syllabus reference |
