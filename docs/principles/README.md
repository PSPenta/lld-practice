# Design Principles

← [Back to hub §8](../../README.md#8-design-principles-solid--dry--kiss--yagni--polk)

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

**In this repo:** RateLimiter2 keeps algorithms in strategy classes; `RateLimiter` only delegates. Also: `JavaScript/SearchEngine/` splits tokenizer, index, ranker; `JavaScript/ParkingLot2/` splits lot, floor, slot, ticket; **`JavaScript/PollingSystem2/PollingService/`** splits `models/`, `PollingService` (use-cases), and `repositories/`. Full map → **[../repo-map/README.md](../repo-map/README.md)**.

#### O — Open/Closed Principle (OCP)

**Meaning:** **Open for extension**, **closed for modification**.

Add behavior by adding new classes, not by editing a giant `if/else` in existing code.

```text
// Bad: edit Allow() every time you add an algorithm
if kind == "token" { ... } else if kind == "leaky" { ... }

// Good: new Strategy class; RateLimiter unchanged
RateLimiter { strategy.Allow(key) }
```

**In this repo:** `JavaScript/RateLimiter2/` — add Sliding Window without rewriting `RateLimiter.js` · `JavaScript/Splitwise/Expense.js` + new expense class · `Go/PaymentGateway-go/` — register new `BankGateway`. Full map → **[../repo-map/README.md](../repo-map/README.md)**.

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

