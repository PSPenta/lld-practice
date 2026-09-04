# Order management (OMS) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Order states?
2. Cancel before ship?
3. Partial fulfillment?
4. Payment integration?
5. Inventory link?
6. Idempotency on create/pay?
7. Multi-item single order?
8. Notifications on state change?

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

## Step 2 — Entities & classes

```text
Order {
  id, userId, items[], total
  status: CREATED | PAID | PACKED | SHIPPED | DELIVERED | CANCELLED
  version
}
OrderItem { sku, qty, price }

OrderService
  - Create(items) → orderId
  - Pay(orderId, paymentRef)
  - Cancel(orderId)
  - MarkPacked / MarkShipped
  - GetStatus(orderId)

InventoryClient.Reserve / Release
PaymentClient.Charge
OrderStateMachine.canTransition(from, to)
```

**Patterns:** State machine · Saga / compensating actions · Repository

## Step 3 — Flows

**Happy path**
1. Create order (CREATED) → reserve inventory  
2. Pay → PAID  
3. Pack → PACKED → Ship → SHIPPED → Deliver  

**Edge cases**
1. Cancel if before PACKED → release inventory; illegal transition rejected  
2. Payment fails after reserve → release; pay retry with idempotency key

## Step 4 — APIs

```http
POST   /orders
POST   /orders/{id}/pay
POST   /orders/{id}/cancel
GET    /orders/{id}
POST   /orders/{id}/ship      # internal / warehouse
```

Webhooks/outbox on state change (optional evolve)

## Step 5 — Deepen

- State machine enforces legal transitions only  
- Saga: payment + inventory — compensate on partial failure  
- Idempotent Create/Pay with client keys  
- Concurrent cancel vs ship → version check; one wins  
- Don’t call payment inside a long DB transaction

## Step 6 — Evolve

- Refund path (extend states without rewriting core — Open/Closed)  
- Outbox for domain events; split shipments / returns  
- Related: [inventory-management](../inventory-management/README.md), [wallet-ledger](../wallet-ledger/README.md)
