# OOP Building Blocks

← [Back to hub §7](../../README.md#7-oop-building-blocks)

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

See **[../repo-map/README.md](../repo-map/README.md)** for the full cross-reference map.

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
| **Service + repositories** | `PollingService` orchestrates; models stay thin; `*Repository` stores | `JavaScript/PollingSystem2/PollingService/` (`models/`, `repositories/`, `index.js`) | `Go/PollingSystem2-go/pollingservice/`, `models/`, `repositories/` |
| **Hybrid: composition + inheritance** | `ExpenseFactory` creates types; `Expense` **has-a** splits · subclasses for exact/equal/percent | `JavaScript/Splitwise/Expense.js`, `JavaScript/Splitwise/index.js` | `Go/Splitwise-go/expense.go` |
| **Inheritance (is-a, teaching)** | `Car` / `Bike` / `Truck` **extends** `Vehicle` | `JavaScript/Parkinglot/Vehicle.js`, `JavaScript/ParkingLot2/Vehicle.js` | `Go/Parkinglot-go/vehicle.go`, `Go/ParkingLot2-go/vehicle.go` |

**How to read the hybrid rows:** inheritance is used for **polymorphism** (vehicle type, expense algorithm). The **orchestrator** still **composes** parts — e.g. `RateLimiter` does not extend `TokenBucket`; it holds a strategy. Same for `ParkingLot2`: the lot composes floors; only `Vehicle` uses is-a.

**Prefer in interviews:** `RateLimiter2`, `ParkingLot2`, `SearchEngine`, `Database`, **`PollingSystem2`** (service + repos) — composition-first designs. Mention `Splitwise` when discussing Factory + when inheritance is acceptable for variant algorithms. Use **`PollingSystem`** v1 only as a refactor drill.

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

