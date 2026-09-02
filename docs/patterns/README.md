# Design Patterns

← [Back to hub §9](../../README.md#9-design-patterns-creational--behavioural--structural)

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

