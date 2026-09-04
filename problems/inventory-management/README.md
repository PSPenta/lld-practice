# Inventory management — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. SKUs per warehouse?
2. Reserve on order or deduct on ship?
3. Allow oversell?
4. Returns restock?
5. Soft hold TTL if payment never completes?
6. Partial reservation when stock low?
7. Idempotency key on reserve?
8. Multi-warehouse in v1?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | OMS, InventoryService, Warehouse stock |
| **Use cases (v1)** | 1. Reserve stock 2. Release on cancel 3. Confirm on ship |
| **In scope** | Available vs reserved qty, reserve/release |
| **Out of scope** | Multi-warehouse optimization |
| **Assumptions** | Single warehouse; reserve then ship |

### Confirm understanding
> "Order reserves inventory; cancel releases; ship confirms deduction."

## Step 2 — Entities & classes

```text
Product { sku, name }
Warehouse { id, name }
Stock {
  sku, warehouseId
  available, reserved
  version                 // optimistic lock
}

Reservation {
  id, orderId, sku, qty, status: HELD|RELEASED|CONSUMED
  expiresAt?
}

InventoryService
  - Reserve(sku, qty, orderId) → reservationId
  - Release(reservationId)
  - ConfirmShip(orderId)      // reserved → consumed
```

**Patterns:** Domain service · Optimistic locking · Reservation as first-class entity

## Step 3 — Flows

**Happy path**
1. OMS calls Reserve → if `available >= qty`, move qty to `reserved`  
2. Payment succeeds → reservation stays HELD until ship  
3. Ship Confirm → decrement reserved (and total on-hand)  
4. Cancel before ship → Release → reserved back to available  

**Edge cases**
1. Concurrent reserves for last unit → one wins via version check; other gets insufficient stock  
2. Hold TTL expires → background release; idempotent Release if already released

## Step 4 — APIs

```text
Reserve(sku, qty, orderId, idempotencyKey?) → reservationId
Release(reservationId)
ConfirmShip(orderId)
GetAvailability(sku) → { available, reserved }
```

```http
POST /inventory/reservations
POST /inventory/reservations/{id}/release
POST /orders/{id}/ship
```

## Step 5 — Deepen

- Optimistic locking (`version`) or row lock on stock to prevent oversell  
- Idempotent reserve by `(orderId, sku)` or client idempotency key  
- Never allow `available < 0`; reject oversell unless policy says otherwise  
- Release and Confirm must be safe under retries  
- Expired holds need a sweeper job (related to [job-scheduler](../job-scheduler/README.md))

## Step 6 — Evolve

- Multi-warehouse allocation Strategy  
- Event sourcing / stock movement audit trail  
- Related: [order-management](../order-management/README.md), [movie-booking](../movie-booking/README.md) (hold TTL)
