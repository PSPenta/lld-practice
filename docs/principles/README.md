# Design Principles

← [Back to hub §8](../../README.md#8-design-principles-solid--dry--kiss--yagni--polk)

Principles guide **how you structure code**. Patterns are reusable **shapes**. Learn principles first.

| Jump | |
|------|--|
| [SRP](#srp) · [OCP](#ocp) · [LSP](#lsp) · [ISP](#isp) · [DIP](#dip) | SOLID |
| [DRY](#dry) · [KISS](#kiss) · [YAGNI](#yagni) · [PoLK](#polk) | Everyday |

Hub one-liner index: [§8](../../README.md#8-design-principles-solid--dry--kiss--yagni--polk)

---

## SOLID

<a id="srp"></a>

### S — Single Responsibility Principle (SRP)

**Meaning:** A class / module should have **one reason to change** — one focused job.

**Why:** Mixed responsibilities force unrelated edits into the same file and make tests brittle.

**Example:** A ticket API that creates tickets *and* sends Slack *and* writes billing will change for three unrelated product reasons. Split into `TicketService`, `SlackNotifier`, `BillingClient`.

| Bad | Good |
|-----|------|
| `UserService` registers users, sends email, charges cards | `UserService`, `EmailNotifier`, `BillingService` |
| One `SearchEngine.js` that tokenizes, indexes, and ranks inline | Separate `Tokenizer`, `InvertedIndex`, `Ranker` |

**Interview line:** “If Slack notification logic changes, I shouldn’t have to touch ticket-creation code.”

**In this repo:** `RateLimiter` only delegates; algorithms live in strategy classes. Also: `JavaScript/SearchEngine/` · `JavaScript/ParkingLot2/` (lot / floor / slot / ticket) · `JavaScript/PollingSystem2/PollingService/` (`models/` · service · `repositories/`). Map → [../repo-map/README.md](../repo-map/README.md).

---

<a id="ocp"></a>

### O — Open/Closed Principle (OCP)

**Meaning:** **Open for extension**, **closed for modification** — add behavior with new types, don’t keep editing a core `switch`.

**Why:** Every new algorithm shouldn’t risk regressing the orchestrator that already works.

**Example:**

```text
// Bad: edit Allow() every time you add an algorithm
if kind == "token" { ... } else if kind == "leaky" { ... }

// Good: new Strategy class; RateLimiter unchanged
RateLimiter { strategy.Allow(key) }
```

Same idea for Splitwise: new expense type → new class + factory entry, not a growing `if` in the service.

**Interview line:** “I’d add a new class that implements the strategy interface instead of editing the core allow path.”

**In this repo:** `JavaScript/RateLimiter2/` · `JavaScript/Splitwise/Expense.js` · `Go/PaymentGateway-go/` new `BankGateway`. Map → [../repo-map/README.md](../repo-map/README.md).

---

<a id="lsp"></a>

### L — Liskov Substitution Principle (LSP)

**Meaning:** Any subtype / interface impl must be **safe to plug in** wherever the base type is used — same contract, no surprises.

**Why:** Callers program to the abstraction; a lying implementation breaks them silently.

**Example:** Code expects `Notifier.send(msg)` to deliver or return a clear error. A `NullNotifier` that always returns OK but drops messages violates LSP if the product requires delivery guarantees.

| Honor the contract | Break the contract |
|--------------------|--------------------|
| Every rate-limit strategy returns a real allow/deny | Strategy that always returns `true` “for tests” left in prod |
| Every `BankGateway` either charges or fails explicitly | Impl that no-ops on failure and reports success |

**Interview line:** “Subtypes must honor the base contract — no weaker preconditions, no surprising side effects.”

**In this repo:** All `JavaScript/RateLimiter2/*` strategies · `BankGateway` in `Go/PaymentGateway-go/` · vehicle/`canFit` behavior in parking.

---

<a id="isp"></a>

### I — Interface Segregation Principle (ISP)

**Meaning:** Prefer **small interfaces**. Don’t force clients to depend on methods they don’t use.

**Why:** Fat interfaces create empty stubs and couple unrelated capabilities.

**Example:**

```text
// Bad fat interface
Device { Read(); Write(); Fax(); Print(); BrewCoffee(); }

// Good
Reader { Read() }
Writer { Write() }
```

In Go interviews: accept interfaces with 1–2 methods (`io.Reader`, `LLMClient.Stream`).

**Interview line:** “I’d keep the strategy interface to a single `isAllowed(key)` — not a kitchen-sink rate-limit admin API.”

**In this repo:** `RateLimiterStrategy.isAllowed()` · `BankGateway.ProcessPayment()` · **Gap:** `SearchEngine` uses concrete deps — you’d introduce slim interfaces only if multiple backends appear (YAGNI vs ISP/DIP).

---

<a id="dip"></a>

### D — Dependency Inversion Principle (DIP)

**Meaning:** High-level policy depends on **abstractions**, not concrete low-level classes.

**Why:** You can swap OpenAI ↔ Anthropic or UPI ↔ card without rewriting the use-case.

**Example:**

```text
High-level: SuggestService / PaymentGateway
                ↓ depends on
Abstract:     LLMClient / BankGateway
                ↑ implemented by
Low-level:    OpenAIAdapter / UPIGateway
```

**Interview line:** “The payment orchestrator depends on a `BankGateway` interface; providers plug in underneath.”

**In this repo:** `RateLimiter` → `RateLimiterStrategy` · `Go/PaymentGateway-go` → `BankGateway`. **Trade-off:** `SearchEngine` wires concrete `Ranker`/`InvertedIndex` — fine for v1 (YAGNI); invert when a second backend shows up.

---

## Everyday principles

<a id="dry"></a>

### DRY — Don’t Repeat Yourself

**Meaning:** Don’t duplicate **knowledge** (rules, formulas, config). One source of truth.

**Why:** Copied validation drifts — one path gets a bugfix, the other doesn’t.

**Example:**

| Duplicated | Better |
|------------|--------|
| Same ticket validation in 4 handlers | Shared `ValidateTicket()` / middleware |
| Credit-debit math in 3 services | One `CreditMeter` / ledger helper |
| Vehicle `switch` in many call sites | One `VehicleFactory` / `ExpenseFactory` |

**Caution:** Don’t smash unrelated snippets into a god `utils.js`. DRY is about shared *meaning*, not merging every similar line.

**In this repo:** `ExpenseFactory` · `VehicleFactory`. **Teaching gap:** per-key bucket helpers repeated across RateLimiter2 strategies — extract in production.

---

<a id="kiss"></a>

### KISS — Keep It Simple Stupid

**Meaning:** Prefer the **simplest design that works** for the stated v1.

**Why:** Extra layers (queues, Abstract Factory, microservices) cost time and interview clarity when requirements don’t need them.

**Example:**
- v1 cache → in-memory map + LRU  
- Only later → Redis when multi-instance or shared state is required  

**Interview line:** “I’d ship a single FIFO / in-memory sheet first, then add complexity when a requirement forces it.”

**In this repo:** `JavaScript/Queue/` · `JavaScript/LRU/` · `Ratelimiter/` v1 before `RateLimiter2/` Strategy ceremony.

---

<a id="yagni"></a>

### YAGNI — You Ain’t Gonna Need It

**Meaning:** Don’t build features or abstractions **until a real need appears**.

**Why:** Speculative “flexibility” is often unused — and harder to explain under time pressure.

**Example:** Don’t invent Abstract Factory + five interfaces for one payment provider. When the **second** provider arrives → introduce `BankGateway` / Strategy.

**Interview line:** “I won’t add an interface until I have a second implementation or a clear seam.”

**In this repo:** Patterns only where variation exists (`RateLimiter2`, `Splitwise`, `PaymentGateway-go`). Contrast `Ratelimiter/` (no Strategy) vs `RateLimiter2/`.

---

<a id="polk"></a>

### PoLK — Principle of Least Knowledge (Law of Demeter)

**Meaning:** Talk only to **direct collaborators** — not friends-of-friends.

**Why:** Deep chains (`a.b.c.d`) couple you to internals; refactors cascade.

**Example:**

```text
// Bad — reach through internals
order.customer.address.zipCode

// Better — ask a direct collaborator
order.shippingZip()
// or customer.primaryZip()
```

**Interview line:** “ParkingLot asks Floor for a free slot — it doesn’t index into `floor.slots[i].vehicle`.”

**In this repo:** `JavaScript/ParkingLot2/ParkingLot.js` → `floor.findAvailableSlot()`.
