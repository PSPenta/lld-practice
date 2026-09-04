# Repository Map

← [Back to hub §7A](../../README.md#7a-repository-map--oop-principles--patterns-by-lld)

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
| **Factory** | Creational | `JavaScript/Splitwise/Expense.js` (`ExpenseFactory`) · `JavaScript/ParkingLot2/Vehicle.js` (`VehicleFactory`) | `Go/Splitwise-go/expense.go` · `Go/ParkingLot2-go/vehicle.go` | Service + pairwise ledger + paise |
| **Observer / Pub-Sub** | Behavioural | `JavaScript/Pub-Sub/index.js` | `Go/Pub-Sub-go/pubsub.go`, `pubsub_core.go`, `pubsub_event.go` | Go has **3 variants** — good “evolution” discussion |
| **Composition / ownership** | Structural | `JavaScript/ParkingLot2/`, `JavaScript/Parkinglot/`, `JavaScript/Database/`, `JavaScript/SearchEngine/` | Matching `Go/*-go/` | has-a chains, not deep inheritance |
| **Facade / pipeline** | Structural | `JavaScript/SearchEngine/SearchEngine.js` orchestrates tokenizer, index, trie, ranker | `Go/SearchEngine-go/search_engine.go` | Facade-like; not named “Facade” in code |
| **Template Method** | Behavioural | `JavaScript/Splitwise/Expense.js` — shared `apply()`, subclasses override `validate()` | `Go/Splitwise-go/expense.go` | Informal hook method on base |
| **Middleware** | (web) | `JavaScript/Ratelimiter/leakyBucket.js` · `JavaScript/UrlShortener/index.js` (Express validators) | `Go/Ratelimiter-go/` | Request pipeline, not GoF |
| **Worker pool / concurrency limit** | (concurrency) | `JavaScript/Ratelimiter/serverRequestThrottler.js` | `Go/Ratelimiter-go/server_request_throttler.go` | Limits parallel work, not rate algorithms |
| **Hybrid inheritance + composition** | Mixed | `JavaScript/Splitwise/` · `JavaScript/ParkingLot2/` · `JavaScript/RateLimiter2/` | Go ports use interfaces/embedding | Orchestrator composes; variants may inherit |
| **Repository-like store** | (data) | `JavaScript/PollingSystem2/PollingService/repositories/` | `Go/PollingSystem2-go/repositories/` | In-memory store; **PollingSystem2** is the clearest **named** repository demo |
| **Application service** | (layering) | `JavaScript/PollingSystem2/PollingService/index.js` — **PollingService** | `Go/PollingSystem2-go/pollingservice/polling_service.go` | Use-cases + auth rules; models stay thin |
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
| **Splitwise** | Service + pairwise ledger + Factory + paise | S, O, KISS | **Factory**, Template Method, application service |
| **Pub-Sub** | Encapsulated event map | S | **Observer** |
| **SearchEngine** | Composition pipeline | S ( ⚠️ no DIP interfaces) | Facade-like pipeline, Trie + inverted index |
| **Database** | Aggregate composition | S | Aggregate / index maps |
| **Redis** / **LRU** | Encapsulation, internal composition | KISS | Map + doubly linked list |
| **UrlShortener** | Minimal OOP | YAGNI ( ⚠️ SRP gap) | Base-62 util; Express middleware |
| **PollingSystem** | `Admin` + static `Polls`/`Results` | S gaps | v1 teaching — spot bugs, refactor to v2 |
| **PollingSystem2** | `PollingService/models/` + `repositories/` | S, layered design | **PollingService** + repos; public/private; creator ≠ voter; `isClosed` |
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
| **SRP** | “SearchEngine splits tokenizer, index, ranker.” → `JavaScript/SearchEngine/SearchEngine.js` · “PollingSystem2: models vs `PollingService` vs repos.” → `JavaScript/PollingSystem2/PollingService/` |
| **Service + Repository** | “Use-cases on PollingService; User/Poll/Vote in `models/`; stores in `repositories/`.” → `JavaScript/PollingSystem2/PollingService/index.js` |
| **Encapsulation** | “LRU hides list mutations in private methods.” → `JavaScript/LRU/index.js` |
| **Polymorphism** | “Slot.canFit uses vehicle subtype.” → `JavaScript/Parkinglot/Slot.js` |
| **KISS / YAGNI** | “Queue is one FIFO class — no extra patterns until needed.” → `JavaScript/Queue/index.js` |
| **Anti-pattern** | “UrlShortener mixes routes and storage — I’d split in production.” → `JavaScript/UrlShortener/index.js` |

---

