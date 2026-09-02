# Vending machine — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Coin types?
2. Change available?
3. Multiple products?
4. Admin refill?
5. Out-of-stock?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Customer, VendingMachine (state machine) |
| **Use cases (v1)** | 1. Insert coins 2. Select product 3. Dispense + change |
| **In scope** | State pattern, inventory, change calc |
| **Out of scope** | Card payment |
| **Assumptions** | Coins only; finite change reservoir |

### Confirm understanding
> "Customer pays, selects slot; machine dispenses if stock and change OK."

---

## Step 2 — Entities & classes

`VendingMachine`, `Product`, `Inventory`, **State** (Idle, HasMoney, Dispensing)

---

## Step 3 — Flows

Insert coin → select product → check stock & change → dispense → return change

---

## Step 4 — APIs

`InsertCoin`, `SelectSlot`, `Dispense`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

State pattern avoids giant switch; insufficient change → refund

---

## Step 6 — Evolve

New payment method → extend state/transitions

---

