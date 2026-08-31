# Low-Level Design (LLD) Interview Preparation Guide

> A practical, detailed guide to prepare for **LLD / machine-coding / OOD** interview rounds at product companies.  
> This repository also contains working **JavaScript** implementations under **`JavaScript/`** and **Golang** ports under **`Go/`** — use them for practice after you learn the method.  
> For **AI code review** rounds (clone a repo, find production/RAG bugs), see **[ai-code-review-round/README.md](ai-code-review-round/README.md)**.

> **Start here (language & runtime first):** Before LLD design and repo problems, revise **JavaScript/Node.js** and **Go** fundamentals using the dedicated guides below. Most full-stack interviews test language/runtime knowledge alongside or before LLD.

| Prep guide | Path | Covers |
|------------|------|--------|
| **JavaScript & Node.js** | **[JavaScript/README.md](JavaScript/README.md)** | `this`, closures, event loop, Promises, streams, auth, gotchas |
| **Golang** | **[Go/README.md](Go/README.md)** | GMP, goroutines, channels, context, GC, interfaces, gotchas |
| **LLD (this doc)** | Below §1 | Design method, SOLID, patterns, worked examples |
| **LLD gaps (breadth)** | **[lld-gaps/README.md](lld-gaps/README.md)** | UML, missing problems, patterns not in repo, 12-week plan |
| **AI code review** | **[ai-code-review-round/README.md](ai-code-review-round/README.md)** | RAG repo review, production bugs (separate round type) |
| **Vibe coding / AI-assisted build** | **[vibe-coding-round/README.md](vibe-coding-round/README.md)** | Cursor/Copilot rounds — prompt, verify, own the design |

---

## Table of contents

0. [Before you start — JavaScript & Go prep](#0-before-you-start--javascript--go-prep)
1. [What is LLD?](#1-what-is-lld)
2. [What interviewers evaluate](#2-what-interviewers-evaluate)
3. [LLD vs HLD vs DSA](#3-lld-vs-hld-vs-dsa)
4. [How a typical LLD round runs](#4-how-a-typical-lld-round-runs)
5. [The standard approach (memorize this)](#5-the-standard-approach-memorize-this)
6. [Clarifying questions checklist](#6-clarifying-questions-checklist)
7. [OOP building blocks](#7-oop-building-blocks)
7A. [Repository map — OOP, principles & patterns by LLD](#7a-repository-map--oop-principles--patterns-by-lld)
8. [Design principles (SOLID + DRY / KISS / YAGNI / PoLK)](#8-design-principles-solid--dry--kiss--yagni--polk)
9. [Design patterns (Creational · Behavioural · Structural)](#9-design-patterns-creational--behavioural--structural)
10. [API design for LLD](#10-api-design-for-lld)
11. [Data modeling & state](#11-data-modeling--state)
12. [Concurrency, idempotency & failure](#12-concurrency-idempotency--failure)
13. [Extensibility & evolution](#13-extensibility--evolution)
14. [Common LLD problems — how to think](#14-common-lld-problems--how-to-think)
15. [AI / LLM LLD for beginners](#15-ai--llm-lld-for-beginners)
16. [Worked example: Cache Client ⭐](#16-worked-example-cache-client-)
17. [Worked example: LRU Cache](#17-worked-example-lru-cache)
18. [Worked example: Rate Limiter](#18-worked-example-rate-limiter)
19. [Worked example: Parking Lot](#19-worked-example-parking-lot)
20. [Worked example: AI Suggest Reply](#20-worked-example-ai-suggest-reply-copilot-draft)
21. [Worked example: Ticket assign + notify](#21-worked-example-ticket-assign--notify)
22. [Redis eviction policies (cache / Redis LLD)](#22-redis-eviction-policies-cache--redis-lld)
23. [HLD topics that bleed into LLD](#23-hld-topics-that-bleed-into-lld)
24. [Timed mock + self-score](#24-timed-mock--self-score)
25. [How to practice (4-week plan)](#25-how-to-practice-4-week-plan)
26. [Interview day checklist](#26-interview-day-checklist)
27. [Problems in this repository](#27-problems-in-this-repository)
28. [Cheat sheet](#28-cheat-sheet)

---

## 0. Before you start — JavaScript & Go prep

**Do not jump straight into LLD.** Product interviews (especially full-stack / backend) usually expect solid **JavaScript/Node.js** and/or **Go** knowledge in addition to design. Complete these guides first, then return here for LLD method and repo problems.

### Recommended order

```text
1. JavaScript/README.md   → language + Node event loop + async (2–3 days)
2. Go/README.md           → if role is Go / polyglot backend (2–3 days)
3. This README (§1+)      → LLD method, SOLID, patterns, worked examples
4. Repo folders           → `JavaScript/*/` (JS) · `Go/*-go/` (Go) — after paper design (§27)
5. lld-gaps/              → paper-design missing problems (Elevator, Notification, UML, …)
6. ai-code-review-round/  → if round is manual code review (not whiteboard LLD)
7. vibe-coding-round/     → if round allows Cursor/Copilot — design first, AI second
```

### Which guide to prioritize?

| Your interview focus | Read first |
|----------------------|------------|
| Node / full-stack JS | **[JavaScript/README.md](JavaScript/README.md)** — then this doc |
| Go backend | **[Go/README.md](Go/README.md)** — then this doc |
| Both / unclear | **JavaScript** (§0–§12) + **Node** (§13–§26), then **Go** §1–§5 if time |
| LLD-only discussion (rare) | Skim §29 cheat sheets; still know event loop / goroutines basics |

### What each guide gives you

| Guide | Quick win sections |
|-------|-------------------|
| **[JavaScript/README.md](JavaScript/README.md)** | §15 Event loop · §19–§20 Promises · §27 gotchas · §29 cheat sheet |
| **[Go/README.md](Go/README.md)** | §1 GMP · §2–§4 concurrency · §18 gotchas · §20 cheat sheet |

**Extended JS Q&A:** [sudheerj/javascript-interview-questions](https://github.com/sudheerj/javascript-interview-questions) — use after this repo’s JavaScript guide, not instead of it.

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

### Interview round types (know which doc to use)

| Round type | What you do | Prep doc |
|------------|-------------|----------|
| **JavaScript / Node.js Q&A** | Language, event loop, Promises, closures, Node APIs | **[JavaScript/README.md](JavaScript/README.md)** — do **before** LLD if full-stack |
| **Golang Q&A** | GMP, goroutines, channels, context, GC | **[Go/README.md](Go/README.md)** — do **before** LLD if Go role |
| **Discussion LLD** (~60 min) | Classes, APIs, flows, trade-offs — no full coding | **This README** (§5–24) — after §0 language prep |
| **Machine coding** (90–120 min) | Working code for a small system | §0 guides + this README + **`JavaScript/*/`** + **`Go/*-go/`** |
| **Vibe coding / AI-assisted build** | Cursor/Copilot allowed; build + verify + narrate | **[vibe-coding-round/README.md](vibe-coding-round/README.md)** + this README §5 for design |
| **AI code review** | Clone repo, manual review — security, RAG, production gaps | **[ai-code-review-round/README.md](ai-code-review-round/README.md)** |
| **LLD breadth (paper only)** | Elevator, Chess, Logger, UML, extra patterns | **[lld-gaps/README.md](lld-gaps/README.md)** |

Same fundamentals (OOD, patterns, RAG vocabulary). **Language/runtime first (§0), then design.** Design rounds = speak structure; review rounds = find bugs with file names.

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

### Encapsulation
Hide internals; expose a small API (`Park`, `Unpark`), not raw fields.

**In this repo:** `JavaScript/LRU/index.js` — private `#removeNode`, `#addNode`; `JavaScript/Redis/index.js` — private `#evict`; `JavaScript/Database/Table.js` — rows/indexes behind methods; `Go/Redis-go/main.go` — `sync.Mutex` + unexported helpers; `Go/Splitwise-go/balance_sheet.go` — unexported balances map.

### Abstraction
Contract without implementation. Callers depend on the contract, not concrete types.

```text
Notifier
  + Send(message) error

EmailNotifier implements Notifier
SlackNotifier implements Notifier
```

**In this repo:** `JavaScript/RateLimiter2/RateLimiterStrategy.js` · `Go/RateLimiter2-go/strategy.go` · `Go/PaymentGateway-go/bank_gateway.go` · abstract `Expense.validate()` in `JavaScript/Splitwise/Expense.js` · `JavaScript/Parkinglot/Vehicle.js` — base class throws if instantiated directly.

### Polymorphism
Same interface, different behavior at runtime — via inheritance overrides or interface implementations.

**In this repo:** `RateLimiter.isAllowed()` → `strategy.isAllowed()` (`JavaScript/RateLimiter2/RateLimiter.js`) · `Slot.canFit(vehicle)` with `instanceof` (`JavaScript/Parkinglot/Slot.js`) · `ExactExpense` / `EqualExpense` / `PercentageExpense.validate()` (`JavaScript/Splitwise/Expense.js`) · `BankGateway.ProcessPayment` (`Go/PaymentGateway-go/bank_gateway.go`).

See **[§7A](#7a-repository-map--oop-principles--patterns-by-lld)** for the full cross-reference map.

### Inheritance vs composition (read this carefully)

Interviewers often ask for **composition over inheritance**. Both are OOP tools; composition is the default for most LLD designs.

#### Inheritance — “is-a”

One type **is a specialized kind of** another. The child inherits fields and methods from the parent.

```text
Vehicle
  ├── Car
  ├── Bike
  └── Truck
```

`Car` **is a** `Vehicle`. Code that accepts `Vehicle` can use a `Car`.

**Works when:** the relationship is stable, shallow, and truly “is-a” (e.g. `Square` is a `Shape` in a geometry exercise).

**Breaks down when:** you inherit to **reuse code** from unrelated concepts:

```text
Email
  └── Ticket
        └── SupportTicket
              └── AISupportTicket   ← fragile chain
```

Change `Email.send()` → every subclass may break, even if it only needed a ticket ID.

#### Composition — “has-a”

Build objects by **combining parts**. Behavior comes from collaborators, not from a parent class.

```text
Ticket
  - id, subject, status
  - assignee: Agent           ← has-a
  - messages: []Message       ← has-a
  - tags: []Tag                ← has-a

TicketService
  - repo: TicketRepository    ← has-a
  - notifier: Notifier        ← has-a (interface)
  - suggest: SuggestService   ← has-a
```

`Ticket` **has** messages and an assignee. It is **not** a subclass of `Email`.

In Go (no classical inheritance):

```go
type Ticket struct {
    ID       string
    Subject  string
    Assignee *Agent
    Messages []Message
}

type TicketService struct {
    Repo     TicketRepository
    Notifier Notifier      // SlackNotifier | EmailNotifier at runtime
    Suggest  SuggestService
}
```

**Interface implementation** (`SlackNotifier implements Notifier`) is a **contract**, not inheritance of state — still composition at the design level.

#### Side-by-side

| | Inheritance | Composition |
|---|-------------|-------------|
| Relationship | **is-a** | **has-a** / **uses-a** |
| Coupling | Child tied to parent changes | Parts swappable behind interfaces |
| Flexibility | Type fixed early | Swap implementations at runtime |
| Testing | Often need subclasses to mock | Inject mocks via interfaces |
| LLD default | Rare for business domains | **Preferred** |

#### Why composition is usually better

1. **Single Responsibility (SOLID)** — Each collaborator has one job. `TicketService` creates tickets; `Notifier` sends Slack; `LLMClient` calls AI. Not one mega-class.
2. **Open/Closed** — Add `SlackNotifier` as a new class implementing `Notifier`. Don’t edit `TicketService` with new `if channel == slack`.
3. **Runtime flexibility** — Tenant A gets Slack, Tenant B gets email — same `TicketService`, different injected `Notifier`.
4. **No fragile base class** — Parent changes don’t ripple through a deep tree.
5. **Easier testing** — Mock `Notifier` or `LLMClient`; no subclass hacks.

**Bad (inheritance thinking):**

```text
class SupportTicket extends Email extends MessageThread {
  assign() { ... }
  suggestReply() { ... }   // AI mixed into email hierarchy
}
```

**Good (composition):**

```text
TicketService
  - repo: TicketRepository
  - notifier: Notifier
  - suggest: SuggestService

Ticket
  - messages: []Message
  - assignee: Agent
```

Add AI later → inject `SuggestService`; don’t subclass `Ticket`.

#### AI / support inbox example (composition)

```text
SuggestService                         ← orchestrator
  ├── ticketRepo: TicketRepository
  ├── retriever: Retriever
  ├── promptBuilder: PromptBuilder
  ├── llm: LLMClient                   ← interface, not "extends OpenAI"
  ├── credits: CreditMeter
  └── guardrails: Guardrails
```

Not: `AISuggestTicket extends Ticket extends Email`.

Same pattern as **`PaymentGateway`** in this repo — business logic depends on `BankGateway` interface; concrete gateways are composed in, not inherited.

#### When inheritance is still OK

| OK | Avoid |
|----|--------|
| Small, stable **is-a** (Shape → Square) | `Ticket extends Email` because both contain text |
| Framework lifecycle hooks (rare in LLD) | Deep trees for changing business rules |
| Implementing a small interface in Go | Inheriting only to reuse unrelated methods |

**Rule for LLD:** If you ask “X **uses** Y” or “X **has** Y” → **composition**. If “X **is a** Y” and it’s stable → inheritance *might* fit; still prefer composition + interface in most product designs.

#### In this repository

**Node.js and Go LLDs that demonstrate composition (study these):**

| Style | Has-a relationship | Node.js (primary) | Go port |
|-------|------------------|-------------------|---------|
| **Composition + Strategy** | `RateLimiter` **has-a** `RateLimiterStrategy` | `JavaScript/RateLimiter2/RateLimiter.js` injects strategy in constructor · strategies: `TokenBucket.js`, `FixedWindowCounter.js`, `LeakyBucket.js`, `SlidingWindowLog.js`, `SlidingWindowCounter.js` · base: `RateLimiterStrategy.js` | `Go/RateLimiter2-go/rate_limiter.go`, `strategy.go`, `token_bucket.go`, … |
| **Composition + Adapter-style interface** | `PaymentGateway` **has-a** `BankGateway` | *(Go only in this repo)* | `Go/PaymentGateway-go/payment_gateway.go`, `bank_gateway.go`, `payment.go` |
| **Composition (ownership chain)** | `ParkingLot` **has-a** `Floor[]` **has-a** `Slot[]` | `JavaScript/ParkingLot2/ParkingLot.js`, `Floor.js`, `Slot.js`, `Ticket.js` | `Go/ParkingLot2-go/parkinglot.go`, `floor.go`, `slot.go` |
| **Composition (pipeline)** | `SearchEngine` **has-a** `Tokenizer`, `Trie`, `InvertedIndex`, `Ranker` | `JavaScript/SearchEngine/SearchEngine.js` · parts: `Tokenizer.js`, `Trie.js`, `InvertedIndex.js`, `Ranker.js` | — |
| **Composition (aggregate)** | `Database` **has-a** `Table` map | `JavaScript/Database/Database.js`, `JavaScript/Database/Table.js` | `Go/Database-go/database.go`, `table.go` |
| **Composition (internal structure)** | `Redis` / `LRUCache` **has-a** map + linked `Node` list | `JavaScript/Redis/index.js`, `JavaScript/LRU/index.js` | `Go/Redis-go/main.go`, `Go/LRU-go/main.go` |
| **Composition (event map)** | `PubSub` **has-a** `Map<event, callbacks[]>` | `JavaScript/Pub-Sub/index.js` | `Go/Pub-Sub-go/pubsub.go`, `pubsub_core.go` |
| **Service + repositories** | `PollService` orchestrates; `User`/`Poll`/`Vote` stay thin; `*Repository` stores | `JavaScript/PollingSystem/PollService.js`, `UserRepository.js`, `PollRepository.js`, `VoteRepository.js` | `Go/PollingSystem-go/poll_service.go`, `repository.go` |
| **Hybrid: composition + inheritance** | `ExpenseFactory` creates types; `Expense` **has-a** splits · subclasses for exact/equal/percent | `JavaScript/Splitwise/Expense.js`, `JavaScript/Splitwise/index.js` | `Go/Splitwise-go/expense.go` |
| **Inheritance (is-a, teaching)** | `Car` / `Bike` / `Truck` **extends** `Vehicle` | `JavaScript/Parkinglot/Vehicle.js`, `JavaScript/ParkingLot2/Vehicle.js` | `Go/Parkinglot-go/vehicle.go`, `Go/ParkingLot2-go/vehicle.go` |

**How to read the hybrid rows:** inheritance is used for **polymorphism** (vehicle type, expense algorithm). The **orchestrator** still **composes** parts — e.g. `RateLimiter` does not extend `TokenBucket`; it holds a strategy. Same for `ParkingLot2`: the lot composes floors; only `Vehicle` uses is-a.

**Prefer in interviews:** `RateLimiter2`, `ParkingLot2`, `SearchEngine`, `Database`, **`PollingSystem`** (service + repos) — composition-first designs. Mention `Splitwise` when discussing Factory + when inheritance is acceptable for variant algorithms.

**Quick code reference — RateLimiter2 composition:**

```1:4:c:\Users\user\Downloads\lld-practice\JavaScript\RateLimiter2\RateLimiter.js
class RateLimiter {
  constructor(strategy) {
    this.strategy = strategy;
  }
```

**Quick code reference — SearchEngine composition:**

```6:11:c:\Users\user\Downloads\lld-practice\JavaScript\SearchEngine\SearchEngine.js
class SearchEngine {
  constructor() {
    this.tokenizer = new Tokenizer();
    this.trie = new Trie();
    this.index = new InvertedIndex();
    this.ranker = new Ranker(this.index);
```

**Interview line with repo proof:**

> “In our repo, RateLimiter2 composes a strategy — the limiter class doesn’t inherit from TokenBucket. ParkingLot2 composes floors and slots. SearchEngine composes tokenizer, index, trie, and ranker. We use inheritance only for true is-a cases like Car extends Vehicle, not for the whole system.”

Prefer the **ParkingLot2 / RateLimiter2 / SearchEngine / Database** style in interviews.

#### Interview one-liner

> “Inheritance models is-a and couples subclasses to parent changes. Composition models has-a — I build a Ticket from Agent, Messages, and a Notifier interface. Responsibilities stay separate, testing is easy, and I can swap Slack, email, or AI without rewriting the core model.”

### Association types (useful vocabulary)

These refine **how** objects relate in composition diagrams:

- **Association** — uses (loose; either can exist alone)  
- **Aggregation** — has (parts can outlive the whole)  
- **Composition** — owns (lifecycle tied; part dies with whole)

Example: `ParkingLot` **owns** `Floor` **owns** `Slot` (composition). `TicketService` **uses** `Notifier` (association via interface).

Not: `OrderService extends DatabaseHelper extends Logger...`

---

## 7A. Repository map — OOP, principles & patterns by LLD

Use this section to **point at real code** in interviews. Paths are relative to the repo root. **Node.js** folders are the primary JS reference; **Go/** contains ports (often with interfaces + mutexes).

### OOP pillars — where to see each one

| OOP concept | What it means | Best repo examples (Node.js) | Go port |
|-------------|---------------|------------------------------|---------|
| **Encapsulation** | Hide state; small public API | `JavaScript/LRU/index.js` (`#` private methods) · `JavaScript/Redis/index.js` · `JavaScript/Database/Table.js` · `JavaScript/Splitwise/BalanceSheet.js` | `Go/Redis-go/main.go` · `Go/Splitwise-go/balance_sheet.go` |
| **Abstraction** | Contract without implementation | `JavaScript/RateLimiter2/RateLimiterStrategy.js` · `JavaScript/Splitwise/Expense.js` (`validate()` hook) · `JavaScript/Parkinglot/Vehicle.js` | `Go/RateLimiter2-go/strategy.go` · `Go/PaymentGateway-go/bank_gateway.go` |
| **Inheritance (is-a)** | Specialized subtype extends base | `JavaScript/RateLimiter2/TokenBucket.js` extends `RateLimiterStrategy` · `JavaScript/Splitwise/Expense.js` subclasses · `JavaScript/ParkingLot2/Vehicle.js` · `JavaScript/Parkinglot/Vehicle.js` | Go uses **interfaces / embedding** instead (e.g. `Go/Splitwise-go/expense.go`) |
| **Composition (has-a)** | Build from parts | `JavaScript/RateLimiter2/RateLimiter.js` · `JavaScript/ParkingLot2/ParkingLot.js` → `Floor.js` → `Slot.js` · `JavaScript/SearchEngine/SearchEngine.js` · `JavaScript/Database/Database.js` | Same folders under `Go/*-go/` |
| **Polymorphism** | Same call, different behavior | Strategy dispatch in `JavaScript/RateLimiter2/RateLimiter.js` · `JavaScript/Parkinglot/Slot.js` `canFit()` + `instanceof` · `JavaScript/Splitwise/Expense.js` overridden `validate()` | `Go/PaymentGateway-go/payment_gateway.go` → `BankGateway` |

**Interview line:** “We demonstrate composition in RateLimiter2 and ParkingLot2, inheritance only where is-a is real (vehicles, expense types), and abstraction via Strategy interfaces.”

---

### SOLID — where each principle appears

| Principle | Meaning (one line) | Repo examples |
|-----------|-------------------|---------------|
| **S — Single Responsibility** | One reason to change per class | `RateLimiter2` — context vs algorithm files · `JavaScript/SearchEngine/` — `Tokenizer.js`, `InvertedIndex.js`, `Ranker.js` each one job · `JavaScript/ParkingLot2/` — lot / floor / slot / ticket split · **Gap:** `JavaScript/UrlShortener/index.js` mixes HTTP + storage (anti-example) |
| **O — Open/Closed** | Extend with new classes, not editing core | `JavaScript/RateLimiter2/` — new strategy file, `RateLimiter.js` unchanged · `JavaScript/Splitwise/Expense.js` (`ExpenseFactory`) + new expense subclass · `Go/PaymentGateway-go/` — new `BankGateway` impl in map |
| **L — Liskov Substitution** | Subtypes honor base contract | All `JavaScript/RateLimiter2/*` strategies implement `isAllowed()` · `JavaScript/Parkinglot/Slot.js` — `Bike`/`Car`/`Truck` used via `instanceof` · **Bad example to avoid:** notifier that silently drops messages |
| **I — Interface Segregation** | Small focused interfaces | `RateLimiterStrategy` — single method · `BankGateway` — `ProcessPayment` only · **Gap:** `SearchEngine` uses concrete deps (no slim interfaces) |
| **D — Dependency Inversion** | High-level depends on abstraction | `JavaScript/RateLimiter2/RateLimiter.js` → `RateLimiterStrategy` · `Go/PaymentGateway-go/payment_gateway.go` → `BankGateway` · **Gap:** `JavaScript/SearchEngine/SearchEngine.js` → concrete `Ranker`, `InvertedIndex` |

---

### DRY · KISS · YAGNI · PoLK — repo examples

| Principle | Repo — good example | Repo — gap / lesson |
|-----------|---------------------|---------------------|
| **DRY** | `JavaScript/Splitwise/Expense.js` (`ExpenseFactory`) · `JavaScript/ParkingLot2/Vehicle.js` (`VehicleFactory`) — one place to create variants | Per-IP bucket logic repeated across `JavaScript/RateLimiter2/*` strategies; parking fallback loops in `JavaScript/ParkingLot2/ParkingLot.js` |
| **KISS** | `JavaScript/Queue/index.js` — single `FIFOQueue` · `JavaScript/LRU/index.js` — map + list only · `JavaScript/Ratelimiter/` v1 — teaching algorithms without Strategy ceremony | Prefer simple v1 before RateLimiter2-level abstraction in interviews |
| **YAGNI** | No Singleton/Builder/Decorator coded — patterns added only where needed | README lists patterns **not** in repo; don’t invent them in design unless requirement asks |
| **PoLK** | `JavaScript/ParkingLot2/ParkingLot.js` calls `floor.findAvailableSlot()` — not `floor.slots[0].vehicle...` | Avoid handlers reaching deep into another module’s internals |

---

### Design patterns — master reference table

| Pattern | Category | LLD / files (Node.js) | Go port | Notes |
|---------|----------|----------------------|---------|-------|
| **Strategy** | Behavioural | `JavaScript/RateLimiter2/RateLimiter.js` + `RateLimiterStrategy.js`, `TokenBucket.js`, `LeakyBucket.js`, `FixedWindowCounter.js`, `SlidingWindowLog.js`, `SlidingWindowCounter.js` | `Go/RateLimiter2-go/` | **Best demo in repo.** Context composes strategy. |
| **Strategy + interface** | Behavioural | — | `Go/PaymentGateway-go/bank_gateway.go`, `payment_gateway.go` | `PaymentGateway` maps method → `BankGateway` |
| **Factory** | Creational | `JavaScript/Splitwise/Expense.js` (`ExpenseFactory`) · `JavaScript/ParkingLot2/Vehicle.js` (`VehicleFactory`) | `Go/Splitwise-go/expense.go` · `Go/ParkingLot2-go/vehicle.go` | Creates concrete type without caller knowing class |
| **Observer / Pub-Sub** | Behavioural | `JavaScript/Pub-Sub/index.js` | `Go/Pub-Sub-go/pubsub.go`, `pubsub_core.go`, `pubsub_event.go` | Go has **3 variants** — good “evolution” discussion |
| **Composition / ownership** | Structural | `JavaScript/ParkingLot2/`, `JavaScript/Parkinglot/`, `JavaScript/Database/`, `JavaScript/SearchEngine/` | Matching `Go/*-go/` | has-a chains, not deep inheritance |
| **Facade / pipeline** | Structural | `JavaScript/SearchEngine/SearchEngine.js` orchestrates tokenizer, index, trie, ranker | `Go/SearchEngine-go/search_engine.go` | Facade-like; not named “Facade” in code |
| **Template Method** | Behavioural | `JavaScript/Splitwise/Expense.js` — shared `apply()`, subclasses override `validate()` | `Go/Splitwise-go/expense.go` | Informal hook method on base |
| **Middleware** | (web) | `JavaScript/Ratelimiter/leakyBucket.js` · `JavaScript/UrlShortener/index.js` (Express validators) | `Go/Ratelimiter-go/` | Request pipeline, not GoF |
| **Worker pool / concurrency limit** | (concurrency) | `JavaScript/Ratelimiter/serverRequestThrottler.js` | `Go/Ratelimiter-go/server_request_throttler.go` | Limits parallel work, not rate algorithms |
| **Hybrid inheritance + composition** | Mixed | `JavaScript/Splitwise/` · `JavaScript/ParkingLot2/` · `JavaScript/RateLimiter2/` | Go ports use interfaces/embedding | Orchestrator composes; variants may inherit |
| **Repository-like store** | (data) | `JavaScript/Database/Database.js`, `Table.js` · `JavaScript/PollingSystem/UserRepository.js`, `PollRepository.js`, `VoteRepository.js` | `Go/Database-go/` · `Go/PollingSystem-go/repository.go` | In-memory store; PollingSystem is the clearest **named** repository demo |
| **Application service** | (layering) | `JavaScript/PollingSystem/PollService.js` — create poll, assign voter, submit vote, stats | `Go/PollingSystem-go/poll_service.go` | Use-cases + auth rules; entities stay thin |
| **Inheritance (teaching)** | OOP | `JavaScript/Parkinglot/Vehicle.js` · `JavaScript/ParkingLot2/Vehicle.js` · `JavaScript/RateLimiter2/RateLimiterStrategy.js` subclasses | Go: enum/constructors instead in several ports | Use when is-a is genuine |

**Not implemented as named patterns in this repo:** Singleton, Builder, Abstract Factory, Prototype, Iterator (use language built-ins), Command, State, Proxy, Decorator, Adapter (PaymentGateway is Strategy-style interface, not classic Adapter), Bridge, Composite (except data-structure composites in JavaScript/LRU/Redis).

**Compare:** `JavaScript/Ratelimiter/` (v1) = algorithms **without** Strategy · `JavaScript/RateLimiter2/` = Strategy + composition — **prefer v2 in interviews.**

---

### Per-LLD quick matrix

| Folder | OOP highlights | Principles | Patterns |
|--------|----------------|------------|----------|
| **RateLimiter2** | Abstraction, composition, polymorphism | S, O, D, I, L | **Strategy** + composition |
| **Ratelimiter** | Encapsulation per algorithm | S, KISS, YAGNI | Middleware; worker pool in `serverRequestThrottler.js` |
| **ParkingLot2** | Composition chain + vehicle inheritance | S | **Factory**, composition, hybrid |
| **Parkinglot** | Inheritance-heavy vehicles + composition levels | S, L | Inheritance + composition |
| **Splitwise** | Inheritance + polymorphic `validate()` | S, O | **Factory**, Template Method–like |
| **Pub-Sub** | Encapsulated event map | S | **Observer** |
| **SearchEngine** | Composition pipeline | S ( ⚠️ no DIP interfaces) | Facade-like pipeline, Trie + inverted index |
| **Database** | Aggregate composition | S | Aggregate / index maps |
| **Redis** / **LRU** | Encapsulation, internal composition | KISS | Map + doubly linked list |
| **UrlShortener** | Minimal OOP | YAGNI ( ⚠️ SRP gap) | Base-62 util; Express middleware |
| **PollingSystem** | Thin `User`/`Poll`/`Vote`; no Admin class | S, layered design | **`PollService`** + **repositories**; creator cannot vote on own poll |
| **Queue** | Single class | KISS, YAGNI | Circular buffer FIFO |
| **PaymentGateway** | `JavaScript/PaymentGateway/` stub · **`Go/PaymentGateway-go/`** full impl | O, D, S | **Strategy** via `BankGateway` |

---

### How to cite this in an interview

| Topic | Say this + point to |
|-------|---------------------|
| **Composition** | “RateLimiter composes a strategy; ParkingLot2 composes floors and slots.” → `JavaScript/RateLimiter2/RateLimiter.js`, `JavaScript/ParkingLot2/ParkingLot.js` |
| **Strategy / OCP** | “New token bucket without editing RateLimiter.” → `JavaScript/RateLimiter2/TokenBucket.js` |
| **Factory** | “ExpenseFactory picks exact/equal/percent.” → `JavaScript/Splitwise/Expense.js` |
| **Observer** | “PubSub map of event → callbacks.” → `JavaScript/Pub-Sub/index.js` |
| **DIP** | “PaymentGateway depends on BankGateway interface.” → `Go/PaymentGateway-go/` |
| **SRP** | “SearchEngine splits tokenizer, index, ranker.” → `JavaScript/SearchEngine/SearchEngine.js` · “PollingSystem: entities vs `PollService` vs repos.” → `JavaScript/PollingSystem/` |
| **Service + Repository** | “Use-cases on PollService; User/Poll/Vote stay data + small invariants; stores behind repositories.” → `JavaScript/PollingSystem/PollService.js` |
| **Encapsulation** | “LRU hides list mutations in private methods.” → `JavaScript/LRU/index.js` |
| **Polymorphism** | “Slot.canFit uses vehicle subtype.” → `JavaScript/Parkinglot/Slot.js` |
| **KISS / YAGNI** | “Queue is one FIFO class — no extra patterns until needed.” → `JavaScript/Queue/index.js` |
| **Anti-pattern** | “UrlShortener mixes routes and storage — I’d split in production.” → `JavaScript/UrlShortener/index.js` |

---

## 8. Design principles (SOLID + DRY / KISS / YAGNI / PoLK)

Principles guide **how you structure code**. Patterns are reusable **shapes**. Learn principles first.

---

### SOLID

#### S — Single Responsibility Principle (SRP)

**Meaning:** A class / function / module should have **one reason to change** (one focused responsibility).

| Bad | Good |
|-----|------|
| `UserService` registers users, sends email, charges cards | `UserService`, `EmailNotifier`, `BillingService` |

**Interview line:** “If Slack notification logic changes, I shouldn’t have to touch ticket-creation code.”

**In this repo:** RateLimiter2 keeps algorithms in strategy classes; `RateLimiter` only delegates. Also: `JavaScript/SearchEngine/` splits tokenizer, index, ranker; `JavaScript/ParkingLot2/` splits lot, floor, slot, ticket; **`JavaScript/PollingSystem/`** splits thin entities (`User`/`Poll`/`Vote`), `PollService` (use-cases), and `*Repository` stores. Full map → **[§7A](#7a-repository-map--oop-principles--patterns-by-lld)**.

#### O — Open/Closed Principle (OCP)

**Meaning:** **Open for extension**, **closed for modification**.

Add behavior by adding new classes, not by editing a giant `if/else` in existing code.

```text
// Bad: edit Allow() every time you add an algorithm
if kind == "token" { ... } else if kind == "leaky" { ... }

// Good: new Strategy class; RateLimiter unchanged
RateLimiter { strategy.Allow(key) }
```

**In this repo:** `JavaScript/RateLimiter2/` — add Sliding Window without rewriting `RateLimiter.js` · `JavaScript/Splitwise/Expense.js` + new expense class · `Go/PaymentGateway-go/` — register new `BankGateway`. Full map → **§7A**.

#### L — Liskov Substitution Principle (LSP)

**Meaning:** Subclasses (or interface implementations) must be **substitutable** for the base type without breaking callers.

If code expects `Notifier.Send(msg) error`, every notifier must:
- honor that contract  
- not panic unexpectedly  
- not silently no-op when the caller expects delivery or a clear error  

**Bad:** `NullNotifier` that pretends success but drops all messages when the product requires delivery guarantees.

**In this repo:** Every `JavaScript/RateLimiter2/*` strategy must implement `isAllowed()` correctly · `JavaScript/Parkinglot/Slot.js` — `canFit()` behavior for `Bike`/`Car`/`Truck` subtypes · all `BankGateway` impls must honor `ProcessPayment` contract (`Go/PaymentGateway-go/`).

#### I — Interface Segregation Principle (ISP)

**Meaning:** No class should be forced to implement methods it doesn’t use. Prefer **small interfaces**.

```text
// Bad fat interface
Device { Read(); Write(); Fax(); Print(); BrewCoffee(); }

// Good
Reader { Read() }
Writer { Write() }
```

**In Go interviews:** “Accept interfaces with 1–2 methods” (`io.Reader`, `LLMClient.Stream`).

**In this repo:** `RateLimiterStrategy.isAllowed()` only · `BankGateway.ProcessPayment()` only · **Contrast (ISP gap):** `SearchEngine` depends on concrete collaborators — in an interview you’d introduce slim interfaces if multiple search backends were planned.

#### D — Dependency Inversion Principle (DIP)

**Meaning:** Depend on **abstractions**, not concrete classes.

```text
High-level: SuggestService / PaymentProcessor
                ↓ depends on
Abstract:     LLMClient / BankGateway / PaymentStrategy
                ↑ implemented by
Low-level:    OpenAIAdapter / UPIGateway / RazorpayPayment
```

**In this repo:** `JavaScript/RateLimiter2/RateLimiter.js` → `RateLimiterStrategy` · `Go/PaymentGateway-go/payment_gateway.go` → `BankGateway` interface. **Gap to mention:** `JavaScript/SearchEngine/SearchEngine.js` uses concrete `Ranker`/`InvertedIndex` — good trade-off discussion (YAGNI vs DIP).

---

### DRY — Don’t Repeat Yourself

**Meaning:** Avoid duplication of **logic, configuration, or behavior**. One source of truth.

| Duplicated | Better |
|------------|--------|
| Same validation copy-pasted in 4 handlers | Shared `ValidateTicket()` / middleware |
| Same credit-debit math in 3 services | One `CreditMeter` |
| Vehicle creation `switch` in many places | `VehicleFactory` / `ExpenseFactory` once (`JavaScript/ParkingLot2/Vehicle.js`, `JavaScript/Splitwise/Expense.js`) |

**Caution:** Don’t force unrelated things into one “util” god-object. DRY is about **knowledge**, not merging every two similar lines. **Gap in repo:** per-IP logic duplicated across `JavaScript/RateLimiter2/*` strategy files — acceptable for teaching, would extract in production.

### KISS — Keep It Simple Stupid

**Meaning:** Choose the **simplest solution** that solves the problem. Avoid unnecessary complexity.

- v1: one in-memory cache  
- Later: Redis, only when multi-instance forces it  

**In this repo:** `JavaScript/Queue/index.js` — one FIFO class, no extra layers · `JavaScript/Ratelimiter/` v1 before `JavaScript/RateLimiter2/` Strategy · `JavaScript/LRU/` — map + list only.

**Interview line:** “I’d start simple and introduce a JavaScript/Queue/cache when a requirement justifies it.”

### YAGNI — You Ain’t Gonna Need It

**Meaning:** Don’t build features until they are **actually needed**.

Don’t add Abstract Factory + 5 interfaces for one payment method “just in case.”  
When the **second** provider arrives → introduce the abstraction.

**In this repo:** No Singleton/Builder/Decorator in code — patterns appear only where variation exists (`RateLimiter2`, `Splitwise`, `PaymentGateway-go`). Compare `JavaScript/Ratelimiter/` (no Strategy) vs `JavaScript/RateLimiter2/` (Strategy when algorithms multiply).

### PoLK — Principle of Least Knowledge (Law of Demeter)

**Meaning:** Objects talk only to **direct collaborators** — not to “friends of friends.”

```text
// Bad — reach through internals
order.customer.address.zipCode

// Better — ask a direct collaborator
order.ShippingZip()
// or customer.PrimaryZip()
```

Reduces coupling: if `Address` structure changes, `Order` callers don’t all break.

**In this repo:** `JavaScript/ParkingLot2/ParkingLot.js` asks `floor.findAvailableSlot()` — does not reach into `floor.slots[i].vehicle` directly. Prefer module-level APIs over deep field chains.

---

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

---

## 9. Design patterns (Creational · Behavioural · Structural)

**Rule:** Use a pattern when it solves a real problem (variation, decoupling, reuse). Don’t sprinkle patterns to sound senior.

GoF groups:

| Category | Question it answers |
|----------|---------------------|
| **Creational** | How do we **create** objects cleanly? |
| **Behavioural** | How do objects **communicate / vary behavior**? |
| **Structural** | How do we **compose** objects for flexibility/scale? |

`*` = useful · `**` = very common in LLD interviews

---

### A) Creational patterns — how to create objects?

#### Singleton `**`

**Intent:** Exactly **one** shared instance (DB pool handle, logger, config).

**Idea:**
- private constructor (or unexported in Go)
- static / package-level instance
- `getInstance()` returns the same object every time

```text
Database.getInstance() → always same connection manager
```

**Pros:** Controlled global access.  
**Cons:** Hidden dependencies, hard to test; in Go prefer **DI** (pass `*sql.DB`) over classic Singleton.

**Interview tip:** Mention thread-safe init (`sync.Once` in Go).

#### Factory `*`

**Intent:** Create objects **without** exposing concrete classes to the caller.

```text
Vehicle (interface) → Car, Bike, Truck
car = VehicleFactory.createVehicle("car")  // returns Car
```

**When:** Many related types; caller shouldn’t `switch` everywhere.

**In this repo:**
- `JavaScript/ParkingLot2/Vehicle.js`, `Go/ParkingLot2-go/vehicle.go` (`CreateVehicle`)
- `JavaScript/Splitwise/Expense.js`, `Go/Splitwise-go/expense.go` (`CreateExpense` / `ExpenseFactory`)

#### Builder `*`

**Intent:** Construct a **complex object step-by-step** (optional fields, readable fluent API).

```text
QueryBuilder
  .select("id, name")
  .from("users")
  .where("active = true")
  .build()
```

Default/empty constructor + setters / chained methods for each part.  
**When:** Many optional params; telescoping constructors get ugly.

#### Abstract Factory

**Intent:** Factory of **families** of related products.

```text
CarAbstractFactory
  SportsCarFactory  → Bugatti, Ferrari
  BudgetCarFactory  → Tata, Mahindra
```

**When:** Multiple product families must stay consistent (UI theme: LightButton + LightDialog vs Dark*).  
**YAGNI warning:** Overkill if you only have one family — use a simple Factory first.

#### Prototype

**Intent:** Create new objects by **cloning** an existing one (avoid expensive `new` + setup).

```js
const obj2 = Object.create(obj1) // JS prototype-style share/clone idea
// or structuredClone(obj1) for a data copy
```

**When:** Object setup is costly; many similar instances with small differences (game pieces, document templates).

---

### B) Behavioural patterns — how objects communicate / vary behavior?

#### Observer `**` (Pub-Sub)

**Intent:** When the **Publisher** changes, **notify all Subscribers**.

```text
Publisher
  - subscribers[]
  + subscribe(sub)
  + unsubscribe(sub)
  + notify()        // call when state changes

Subscriber
  + update(event)
```

**Classic LLD:** Notification service — ticket assigned → email + Slack + analytics.

**In this repo:** `JavaScript/Pub-Sub/`, `Go/Pub-Sub-go/pubsub.go` (and `pubsub_core.go`, `pubsub_event.go`).

#### Strategy `**`

**Intent:** Same job, **different algorithms**, interchangeable at runtime.

```text
PaymentStrategy (interface)
  RazorpayPayment implements PaymentStrategy
  StripePayment implements PaymentStrategy

processor = new PaymentProcessor(new RazorpayPayment(...))
processor.pay(amount)
```

**Also:** Rate limit algorithms; expense split types (Equal/Exact/Percentage as strategies + factory).

**In this repo:**
- **Best example:** `JavaScript/RateLimiter2/` + `Go/RateLimiter2-go/strategy.go`
- Payment-style: `Go/PaymentGateway-go/bank_gateway.go`

#### Iterator `*`

**Intent:** Traverse a collection **without exposing** internal structure.

Built into JS (`for...of`, generators) and Go (`range` over slices — not a formal Iterator type).  
**LLD example:** Playlist iterator — `hasNext()` / `next()` over songs.

#### Command `*`

**Intent:** Encapsulate a request as an **object** (execute / undo / queue).

```text
MigrationCommand
  up()   → create table, run ops
  down() → rollback
```

One command can map to multiple steps: open DB → verify table → run operation.  
**Also:** remote controls, job queues, undo stacks.

#### State (know for interviews)

Behavior changes when **internal state** changes (ticket Open → Closed; vending machine). Often cleaner than huge `switch(status)`.

#### Template Method (know for interviews)

Skeleton algorithm in a base class; subclasses fill in steps (import CSV vs JSON with same pipeline).

---

### C) Structural patterns — how to compose objects for flexibility?

#### Adapter `**`

**Intent:** **Bridge** between two incompatible interfaces so they work together **without rewriting** either side.

```text
Your app expects: BankGateway.ProcessPayment(p)
Vendor SDK has:   razorpay.Charge(card, paise)

Adapter implements BankGateway and translates calls → vendor SDK
```

**Classic story:** Smart home hub speaking one protocol; each device adapter speaks Zigbee/Wi‑Fi/etc.  
**In this repo:** payment gateways behind `BankGateway` (`Go/PaymentGateway-go/`).

#### Proxy `**`

**Intent:** A **stand-in** object that controls access to another (caching, auth, rate limit, remote call).

```text
Client → Nginx reverse proxy → microservice
Client → RepositoryProxy (cache) → real Repository
```

**Microservices:** API gateway / reverse proxy in front of services.  
**App-level:** lazy-loading proxy, protection proxy.

#### Decorator `*`

**Intent:** Add behavior **dynamically** by wrapping objects (same interface).

```text
Coffee = WhippedCream(Sugar(Milk(BasicCoffee)))
cost() and description() stack
```

**Also:** HTTP middleware chain — logging(auth(handler))).  
Flexible alternative to deep inheritance.

#### Facade `*`

**Intent:** A **simple front door** over a messy subsystem.

```text
VideoConversionFacade.convert(file, "mp4")
  // internally: decoder, encoder, bitrate, filesystem...
```

Callers don’t need to know 15 internal classes.  
**LLD:** `OrderFacade.placeOrder()` orchestrates inventory + payment + notify.

#### Composite / Bridge (optional extras)

- **Composite:** tree of objects treated uniformly (file/folder UI).  
- **Bridge:** split abstraction from implementation (Remote ↔ TVBrand) — rarer in short LLD rounds.

---

### Patterns used in this repository (reference files)

| Pattern | LLD in this repo | How it shows up | Reference files |
|---------|------------------|-----------------|-----------------|
| **Strategy + composition** | **RateLimiter2** | `RateLimiter` **has-a** strategy; algorithms swappable at runtime | JS: `JavaScript/RateLimiter2/RateLimiter.js`, `RateLimiterStrategy.js`, `TokenBucket.js`, `LeakyBucket.js`, `FixedWindowCounter.js`, `SlidingWindowLog.js`, `SlidingWindowCounter.js` · Go: `Go/RateLimiter2-go/rate_limiter.go`, `strategy.go`, `token_bucket.go`, `leaky_bucket.go`, `fixed_window_counter.go`, `sliding_window_log.go`, `sliding_window_counter.go` |
| **Strategy + Adapter + composition** | **PaymentGateway** | `PaymentGateway` **has-a** `BankGateway` interface | Go: `Go/PaymentGateway-go/payment_gateway.go`, `bank_gateway.go`, `payment.go`, `main.go` |
| **Composition (ownership)** | **ParkingLot2** | `ParkingLot` → `Floor` → `Slot`; uses `Ticket` | JS: `JavaScript/ParkingLot2/ParkingLot.js`, `Floor.js`, `Slot.js`, `Ticket.js` · Go: `Go/ParkingLot2-go/parkinglot.go`, `floor.go`, `slot.go` |
| **Composition (pipeline)** | **SearchEngine** | `SearchEngine` **has-a** tokenizer, trie, index, ranker | JS: `JavaScript/SearchEngine/SearchEngine.js`, `Tokenizer.js`, `Trie.js`, `InvertedIndex.js`, `Ranker.js` |
| **Composition (aggregate)** | **Database** | `Database` **has-a** map of `Table` | JS: `JavaScript/Database/Database.js`, `Table.js` · Go: `Go/Database-go/database.go`, `table.go` |
| **Factory + inheritance (hybrid)** | **Splitwise** | `ExpenseFactory` + `ExactExpense`/`EqualExpense`/`PercentageExpense` extend `Expense` | JS: `JavaScript/Splitwise/Expense.js`, `JavaScript/Splitwise/index.js` · Go: `Go/Splitwise-go/expense.go`, `main.go` |
| **Factory + inheritance (hybrid)** | **ParkingLot2** | Vehicle factory pattern + `Car`/`Bike`/`Truck` extend `Vehicle` | JS: `JavaScript/ParkingLot2/Vehicle.js`, `index.js` · Go: `Go/ParkingLot2-go/vehicle.go`, `main.go` |
| **Observer / Pub-Sub + composition** | **Pub-Sub** | `PubSub` **has-a** event → callbacks map | JS: `JavaScript/Pub-Sub/index.js` · Go: `Go/Pub-Sub-go/pubsub.go`, `pubsub_core.go`, `pubsub_event.go` |
| **Composition (cache internals)** | **Redis**, **LRU** | map + doubly linked `Node` list | JS: `JavaScript/Redis/index.js`, `JavaScript/LRU/index.js` · Go: `Go/Redis-go/main.go`, `Go/LRU-go/main.go` |
| **Inheritance only (v1 teaching)** | **Parkinglot** | Car/Bike/Truck extend Vehicle — lot itself is simpler | JS: `JavaScript/Parkinglot/Vehicle.js` · Go: `Go/Parkinglot-go/vehicle.go` |

**Know for interviews but not clearly coded as named patterns here:** Singleton, Builder, Abstract Factory, Prototype, Iterator, Command, State, Template Method, Proxy, Decorator, Facade.

**Note:** `Ratelimiter` / `Go/Ratelimiter-go` = standalone algorithm types, **not** Strategy. Prefer **RateLimiter2** to demo Strategy.

### How to mention a pattern in interview

> “I’ll use Strategy because we have interchangeable algorithms — same idea as **RateLimiter2**: `RateLimiter` **composes** a strategy interface; `TokenBucket` is injected, not inherited.”

> “For structure, I prefer composition like **ParkingLot2** (lot owns floors owns slots) or **SearchEngine** (orchestrator owns tokenizer, index, ranker) — see our repo’s Node implementations.”

If only one implementation exists and none is planned → **YAGNI**: skip the pattern.

### Pattern ↔ principle quick links

| You say… | You’re applying… |
|----------|------------------|
| New class instead of editing switch | OCP + often Strategy/Factory |
| Depend on `LLMClient` interface | DIP + Adapter |
| Split God class | SRP |
| Don’t build Abstract Factory yet | YAGNI + KISS |
| Middleware wraps handler | Decorator |

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

Classic LLD skills still apply. AI rounds add **provider abstraction, RAG, credits, streaming, and unsafe model output**. See **[§15](#15-ai--llm-lld-for-beginners)**, **[§16 Cache Client](#16-worked-example-cache-client-)**, and **[§20 AI Suggest Reply](#20-worked-example-ai-suggest-reply-copilot-draft)**.

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
- Strategy: `JavaScript/RateLimiter2/`, `Go/RateLimiter2-go/`
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
4. **Suggest-reply / copilot** — full orchestration (see §20)  
5. **Tool-calling agent (basic)** — model returns tool name → your code runs tool → continue  

### What “good enough for a noob” looks like in the interview

You can draw the boxes, name interfaces, walk the happy path, then say:

- credits under concurrency  
- LLM timeout  
- tenant isolation in RAG  
- how you’d add a new model vendor (new Adapter only)

You do **not** need to implement a vector DB.

### Five AI LLD problems — expanded sketches

Practice **Suggest Reply (§20)** and **LLM abstraction** deeply; sketch the rest.

**A — LLM provider abstraction:** `LLMClient` interface with `Complete` / `Stream`; `OpenAIAdapter`, `AnthropicAdapter`; `LLMRouter` (Strategy) by plan/cost; failover on timeout.

**B — RAG retriever:** Offline: chunk → embed → store with `{ tenant_id, doc_id }`. Online: embed query → vector search **filtered by tenant** → top-K → prompt. Never cross-tenant retrieval.

**C — AI credit meter:** `Reserve` → `Commit` / `Release` (wallet-like, not just `Allow()`). DB ledger is source of truth; not cache alone.

**D — Suggest reply:** Full flow in §20.

**E — Tool-calling agent:** Loop: LLM → tool call → your code runs tool → append result → LLM again; cap max iterations; AuthZ per tool.

### Evaluation & AI prep drill

**Offline eval:** Golden tickets + expected facts. **Online:** thumbs, edit distance to sent reply, latency, cost. Store `prompt_version` + model on each suggestion.

**15-min drill:** Draw architecture → Suggest happy path → timeout/release credits → RAG + tenant filter → new vendor = new Adapter only.

For **RAG/code review depth** (chunking, retrieval, vector DB, reviewing repos like ragbot), see **[ai-code-review-round/README.md](ai-code-review-round/README.md)**.

---

## 16. Worked example: Cache Client ⭐

> **Very common at product companies** (including helpdesk/SaaS interviews): *“Design a Cache Client that caches frequent queries.”*

### Step 1 — Clarifying questions (ask 6–8)

1. Library or HTTP service?
2. In-memory only, or JavaScript/Redis/shared cache?
3. Single machine or many servers?
4. TTL (expire after time)?
5. Max capacity? Eviction policy (LRU)?
6. Thread-safe? (multiple goroutines)
7. On miss — return error, or load from DB (`GetOrLoad`)?
8. Cache stampede: 100 concurrent misses on same key — all hit DB?
9. Metrics needed (hit/miss rate)?

**State assumptions if they say “your call”:** in-memory v1, TTL + LRU, thread-safe, `GetOrLoad` with loader.

### Step 2 — Entities & classes

```text
CacheEntry: key, value, expiresAt

interface Cache {
  Get(key) (value, found)
  Set(key, value, ttl)
  Delete(key)
}

CacheClient implements Cache
  - store: map[string]*CacheEntry
  - order: doubly linked list (LRU)
  - capacity, defaultTTL
  - mu: Mutex

  GetOrLoad(key, loader func() (value, error))
```

Repo references: `JavaScript/LRU/`, `JavaScript/Redis/`, `Go/LRU-go/`, `Go/Redis-go/`.

### Step 3 — Flows

**Get:** lock → miss/expired → return miss → else move to MRU → unlock → return value  

**Set:** lock → update or insert → evict LRU tail if over capacity → unlock  

**GetOrLoad:** on miss, call loader — use **singleflight** so 50 concurrent misses → one DB call.

### Step 4 — API (if HTTP service)

```http
GET    /v1/cache/{key}
PUT    /v1/cache/{key}   body: { "value": "...", "ttl_sec": 60 }
DELETE /v1/cache/{key}
```

Often a **library** — ask first.

### Step 5 — Trade-offs / evolve

| Question | Answer |
|----------|--------|
| Traffic ↑ | Redis L2; local L1 |
| Memory ↑ | Capacity cap + TTL |
| Redis down | Fallback to DB (degraded) |
| Extend | `Cache` interface → Memory / Redis (Adapter) |
| Eviction | Strategy: LRU, LFU (see §22 for Redis policies) |
| Monitor | hit_rate, miss_rate, eviction_count, p99 |

**Patterns:** Strategy/Adapter for backends; SRP — cache does not know SQL (loader injected).

**Practice:** Explain in 20 minutes out loud without notes.

---

## 17. Worked example: LRU Cache

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

## 18. Worked example: Rate Limiter

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

## 19. Worked example: Parking Lot

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

## 20. Worked example: AI Suggest Reply (copilot draft)

Use this as your **AI LLD template**. Same structure as JavaScript/LRU/Parking Lot interviews.

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

## 21. Worked example: Ticket assign + notify

Helpdesk / shared-inbox classic (assign owner, notify team).

### Clarify

- Single assignee or multiple?
- Notify which channels (Slack, email, in-app)?
- Concurrent assign allowed?

### Classes

```text
TicketService
  - Assign(ticketID, agentID) error
  - Get(ticketID)

TicketRepository
Notifier (interface) → SlackNotifier | EmailNotifier | InAppNotifier
AssignmentPolicy → CanAssign(ticket, agent)
```

### Assign flow

```text
1. Load ticket (tenant-scoped)
2. Policy: agent in same tenant/inbox?
3. Optimistic lock: UPDATE ... WHERE version=?
   → 0 rows → 409 Conflict ("already assigned")
4. Save ticket, version++
5. Publish TicketAssigned event → Notifiers (Observer)
```

### Patterns

- **Observer:** one event, many notifiers  
- **Adapter:** Slack API wrapper  
- **Repository:** hide DB  

### Backup: Webhook ingest + connector sync (idempotent / async)

**Webhook (inbound events):**
```text
Verify signature → store event_id (dedupe) → enqueue → return 200 fast
Worker upserts ticket/message — never call LLM inline on webhook path
```

**Connector sync (pull files from Slack/Notion/Drive):**
```text
Bad:  POST /sync → fetch all files → chunk → embed in request (client waits)
Good: POST /sync → publish job to SQS/Kafka → 202 { job_id }
      Worker → fetch → ingest each doc → update last_synced_at (retries, DLQ)
```

Same rule: **heavy I/O and indexing never block the HTTP handler.**

---

## 22. Redis eviction policies (cache / Redis LLD)

When Redis hits `maxmemory`, it evicts keys per `maxmemory-policy`:

| Policy | Evicts |
|--------|--------|
| `noeviction` | Nothing — writes fail when full |
| `allkeys-lru` | Any key — least recently used (approximate) |
| `allkeys-lfu` | Any key — least frequently used |
| `volatile-lru` | Only keys **with TTL** — LRU among those |
| `volatile-ttl` | Keys with TTL — shortest remaining TTL first |
| `allkeys-random` / `volatile-random` | Random key |

**Interview lines:**
- Pure cache → `allkeys-lru` or `allkeys-lfu`
- Mix permanent + cache keys → `volatile-lru` + TTL on cache only
- Real Redis uses **approximate** LRU (samples keys), not exact linked-list LRU like `JavaScript/LRU/` in this repo
- **Eviction** (memory full) ≠ **expiration** (TTL) — related but different

---

## 23. HLD topics that bleed into LLD

Some **HLD** appears in LLD as *“how would you evolve?”* — not full multi-region design.

### Evolution template

```text
v1: Sync API → service → DB → LLM
v2: Add Redis (tenant config, hot retrieval)
v3: Queue (webhooks, bulk embed KB)
v4: Horizontal scale stateless API + workers
```

### Common bleed-in topics

| Topic | One-liner |
|-------|-----------|
| Caching | L1 in-process + L2 Redis; invalidate on KB update |
| Queue | Webhook + **connector file sync**: enqueue job (SQS/Kafka), 202 fast, worker ingests; DLQ |
| Rate limit / credits | Per-tenant 429 + reserve/commit |
| DB | Index `(tenant_id, inbox_id, status)`; cursor pagination |
| Scale out | Stateless API; shared DB/Redis/queue |
| Observability | `trace_id` across services; p99 per step |
| Multi-tenant | Every query/filter includes `tenant_id` |

**3-step answer:** (1) “I’d evolve without rewriting core classes.” (2) One concrete step. (3) Trade-off. Then stop.

**Mini-scenario:** Gmail webhook storm → verify, dedupe `event_id`, enqueue, 200 fast. **Connector sync** → don’t fetch Slack/GDrive in HTTP thread; enqueue sync job, return 202, workers chunk+embed. LLM suggest 15s timeout → release credits; `trace_id` ties flow.

---

## 24. Timed mock + self-score

**Timer: 40 minutes.** Pick **Cache Client (§16)** or **Rate Limiter (§18)**. Speak out loud; no notes for first 30 min.

After time, score 0–2 each (target ≥ 16/20):

1. Clarified requirements  
2. Stated assumptions  
3. Named entities and classes  
4. Drew simple diagram  
5. Walked main flows (Get/Set or Allow)  
6. Mentioned mutex / thread safety  
7. Mentioned failure case (DB/Redis/LLM down)  
8. Showed evolution (interface, Redis, queue)  
9. Named 1 pattern with reason  
10. Explained trade-offs plainly  

**Opening (memorize):** “I’ll clarify requirements first, then entities, classes, flows, APIs, concurrency, failures, and how the design evolves.”

**Closing:** “I kept v1 simple and showed how it evolves under load and failures. Happy to deep-dive any part.”

---

## 25. How to practice (4-week plan)

> **Prerequisite:** Finish **[JavaScript/README.md](JavaScript/README.md)** (and **[Go/README.md](Go/README.md)** if applicable) before Week 2 repo coding — see **[§0](#0-before-you-start--javascript--go-prep)**.

### Week 0 — Language & runtime (before LLD coding)
- **[JavaScript/README.md](JavaScript/README.md):** event loop, Promises, closures, §27 gotchas, §29 cheat sheet  
- **[Go/README.md](Go/README.md)** (if Go role): GMP, goroutines, channels, context, §18 gotchas  
- Answer §28 checklists in each guide aloud  

### Week 1 — Foundations
- SOLID + 5 patterns with your own examples  
- Draw class diagrams for LRU + Parking Lot on paper  
- Practice clarifying questions aloud  

### Week 2 — Core problems
Implement or redesign:
- LRU, Rate Limiter, Parking Lot, Pub-Sub  
- For each: write APIs + failure notes  

### Week 3 — More problems + AI + concurrency
- Splitwise, **PollingSystem** (service + repos), Notification system, Cache client  
- **AI:** LLM interface + RAG sketch + Suggest Reply (§20) once on paper  
- Add mutex/idempotency discussion every time  
- One machine-coding simulation (90 min timer)  

### Week 4 — Interview simulation
- 4 timed discussion LLDs (45 min each) — include **one AI** problem  
- Record yourself; check structure adherence  
- Prepare 2 stories linking designs to past work  
- Optional: one **vibe coding** mock ([vibe-coding-round/README.md](vibe-coding-round/README.md))  
- Paper-design 2 **gap** problems from [lld-gaps/README.md](lld-gaps/README.md)

### Practice method (every problem)
1. 5 min questions  
2. 10 min entities  
3. 20 min classes + flow  
4. 10 min API + deepen  
5. Compare with a reference solution (this repo / articles)  
6. Note 3 improvements  

---

## 26. Interview day checklist

**Before**
- [ ] **Language prep done:** [JavaScript/README.md](JavaScript/README.md) · [Go/README.md](Go/README.md) if Go role — see **§0**  
- [ ] Paper / shared editor ready  
- [ ] Know opening + closing lines (§24)  
- [ ] Know round type: **LLD design** (this doc) vs **code review** ([ai-code-review-round](ai-code-review-round/README.md))  
- [ ] Skim company product (helpdesk/SaaS: shared inbox, AI copilot, multi-tenant) — honest if you didn’t use the product  
- [ ] Sleep; don’t cram 10 new patterns  

**During (LLD design)**
- [ ] Ask questions first  
- [ ] State assumptions  
- [ ] Drive use-case by use-case  
- [ ] Prefer composition + interfaces  
- [ ] Mention concurrency & failure at least once  
- [ ] Show how design evolves (§23)  

**During (code review)**
- [ ] Trace main happy path, then untested routes  
- [ ] Prioritize: security → RAG correctness → production ops  
- [ ] Cite file names; group findings P0/P1/P2  

**Avoid**
- Jumping to Kafka/microservices immediately  
- Pattern-dropping without need  
- Coding before the model is clear (discussion rounds)  
- Ignoring edge cases (null, empty, capacity full)  
- Fake deep product usage  

### Company research template (helpdesk / AI SaaS)

If they ask *“What do you know about us?”*:

1. **One sentence:** AI-powered customer support — shared inbox + workflows + copilot/agents.  
2. **Problem:** Teams sharing support@ without chaos; knowledge for agents; draft replies grounded in docs.  
3. **Your fit:** Multi-tenant APIs, RAG, streaming, credits, reliability — same design space.  
4. **Honest product use:** Website/docs reviewed; trial if offered.

---

## 27. Problems in this repository

Use these as **hands-on practice** after designing on paper. **JavaScript implementations** live under **`JavaScript/`**; Go ports under **`Go/`**.

### JavaScript (`JavaScript/*/` — machine-coding style)

| Path | Focus | OOP | Principles | Patterns |
|------|--------|-----|------------|----------|
| **`JavaScript/RateLimiter2/`** | Strategy, swappable algorithms | Abstraction, composition, polymorphism | S, O, D, I, L | **Strategy** |
| **`JavaScript/ParkingLot2/`** | Multi-floor parking | Composition chain + vehicle is-a | S | **Factory**, composition |
| **`JavaScript/Splitwise/`** | Expense splitting | Inheritance + polymorphic `validate()` | S, O | **Factory**, Template Method–like |
| **`JavaScript/SearchEngine/`** | Index, rank, autocomplete | Composition pipeline | S ( ⚠️ DIP) | Facade-like pipeline |
| **`JavaScript/Pub-Sub/`** | Eventing | Encapsulated map | S | **Observer** |
| **`JavaScript/Database/`** | Schema, CRUD | Aggregate composition | S | Aggregate store |
| **`JavaScript/Redis/`** / **`JavaScript/LRU/`** | Cache eviction | Encapsulation | KISS | Map + DLL |
| **`JavaScript/Parkinglot/`** | Simple parking v1 | Inheritance-heavy vehicles | S, L | is-a + composition |
| **`JavaScript/Ratelimiter/`** | Standalone algorithms | Per-file encapsulation | KISS, YAGNI | Middleware; worker pool |
| **`JavaScript/PollingSystem/`** | Polls & votes (assign, one-vote, no self-vote) | Thin `User`/`Poll`/`Vote` | S + service layer | **`PollService`** + `UserRepository` / `PollRepository` / `VoteRepository` · Go: `Go/PollingSystem-go/` |
| **`JavaScript/Queue/`** | FIFO mechanics | Single class | KISS, YAGNI | Circular buffer |
| **`JavaScript/UrlShortener/`** | HTTP + short codes | Minimal | YAGNI ( ⚠️ SRP) | Encoding util |
| **`JavaScript/PaymentGateway/`** | Stub only (empty files) | — | — | **Full impl:** `Go/PaymentGateway-go/` — Strategy via `BankGateway` |

Full cross-reference → **[§7A](#7a-repository-map--oop-principles--patterns-by-lld)**.

### Golang (`Go/*-go`)

Same problems ported to Go — good for Go interviews (interfaces, mutexes, errors).  
**Language prep (do first):** [Go/README.md](Go/README.md) · **JS/Node prep:** [JavaScript/README.md](JavaScript/README.md) · **Code review round:** [ai-code-review-round/README.md](ai-code-review-round/README.md) · **LLD method:** this doc §1+.

**Note:** This repo’s coded LLDs are classic systems (cache, parking, rate limit, etc.). **AI LLD is covered in this README (§15–§20)** and **code review in ai-code-review-round/** — practice on paper/whiteboard; there is no separate `AI-Suggest-go` folder yet.

### Suggested workflow with this repo

1. Complete **[JavaScript/README.md](JavaScript/README.md)** and/or **[Go/README.md](Go/README.md)** (§0).  
2. Hide the code.  
3. Design on paper for 30–40 minutes (this doc’s method, §5).  
4. Only then open the folder and compare.  
5. Re-implement a smaller version from scratch timed.

---

## 28. Cheat sheet

### Structure
```text
Clarify → Entities → Classes/Interfaces → Flows → APIs →
Concurrency/Failure → Extend/Trade-offs
```

### SOLID + other principles one-liners
- **S:** one reason to change  
- **O:** extend without editing core  
- **L:** subtypes don’t break contracts  
- **I:** small interfaces  
- **D:** depend on abstractions  
- **Composition:** has-a + interfaces over deep is-a inheritance  
- **DRY:** don’t duplicate knowledge  
- **KISS / YAGNI / PoLK:** simplest design; add patterns when needed — repo map **§7A**  
- **KISS:** simplest workable design  
- **YAGNI:** don’t build unused features  
- **PoLK:** talk only to direct collaborators  

### Patterns one-liners
**Creational:** Singleton (one instance) · Factory (create by type) · Builder (step-by-step) · Abstract Factory (product families) · Prototype (clone)  
**Behavioural:** Observer (notify subscribers) · Strategy (swap algorithm) · Iterator (traverse) · Command (request as object) · State · Template Method  
**Structural:** Adapter (bridge interfaces) · Proxy (stand-in / gateway) · Decorator (wrap behavior) · Facade (simple front door)  
**Repo demos:** Full map → **§7A**. Quick: Strategy → RateLimiter2 · Factory → Splitwise / ParkingLot2 · Observer → Pub-Sub · DIP → PaymentGateway-go · Pipeline → SearchEngine · KISS → JavaScript/Queue/LRU

### AI LLD one-liners
- LLM behind `LLMClient` interface (Adapter)  
- RAG = retrieve tenant-scoped chunks → prompt → generate  
- Credits = Reserve before call, Commit/Release after  
- Never trust model output — validate  
- Timeout LLM calls; ticket must work if AI is down  
- Stream (SSE) for UX; queue for bulk/offline jobs  
- **Code review:** secrets, SQL injection, grounded RAG → [ai-code-review-round](ai-code-review-round/README.md)

### HLD bleed-in (when asked “scale?”)
- v1 simple → Redis → queue → stateless scale → indexes + `trace_id`  
- Multi-tenant: filter every query by `tenant_id`

### Must-practice designs
- **Cache Client (§16)** · Rate Limiter (§18) · AI Suggest (§20)

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

Start with **Cache Client (§16) → Rate Limiter (§18) → Parking Lot → Splitwise → Pub-Sub**, then **AI Suggest Reply (§20)**, then branch into domain problems. For **code review** prep, use **ai-code-review-round/README.md**.

Good luck.
