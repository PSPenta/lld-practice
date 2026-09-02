# Order management (OMS) — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Order states?
2. Cancel before ship?
3. Partial fulfillment?
4. Payment integration?
5. Inventory link?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Customer, OrderService, Inventory, Payment |
| **Use cases (v1)** | 1. Place order 2. Pay 3. Ship or cancel |
| **In scope** | State machine, cancel rules |
| **Out of scope** | Returns, split shipments |
| **Assumptions** | CREATED→PAID→SHIPPED; cancel before packed |

### Confirm understanding
> "Customer places order; after payment it moves through fulfillment states."

---

## Step 2 — Entities & classes

`Order`, `OrderItem`, `OrderState` enum, `OrderService`, `InventoryClient`, `PaymentClient`

---

## Step 3 — Flows

Create → pay → allocate inventory → pack → ship → deliver; cancel if before packed

---

## Step 4 — APIs

`CreateOrder`, `Cancel(orderId)`, `GetStatus`, webhook on state change

---

## Step 5 — Deepen (concurrency, failure, idempotency)

State machine enforces legal transitions; saga for payment+inventory rollback

---

## Step 6 — Evolve

Add refund path (Open/Closed on states); outbox for events

---

