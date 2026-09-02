# Inventory management — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. SKUs per warehouse?
2. Reserve on order or deduct on ship?
3. Allow oversell?
4. Returns restock?

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

---

## Step 2 — Entities & classes

`Product`, `Warehouse`, `Stock`, `Reservation`, `InventoryService`

---

## Step 3 — Flows

Reserve qty → decrement available → on cancel release → on ship confirm deduct

---

## Step 4 — APIs

`Reserve(sku, qty)`, `Release(reservationId)`, `Ship(orderId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Optimistic locking on stock row; idempotent reserve key

---

## Step 6 — Evolve

Multi-warehouse allocation strategy; event sourcing for audit

---

