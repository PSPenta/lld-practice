# Vending machine — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Coin types?
2. Change available?
3. Multiple products?
4. Admin refill?
5. Out-of-stock?
6. Cancel / return coins mid-flow?
7. Exact change only mode?
8. One customer transaction at a time?

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

## Step 2 — Entities & classes

```text
Product { slotId, name, price, quantity }
Inventory { map[slotId]Product; refill(slot, qty) }
CoinInventory { denominations → counts; makeChange(amount) }

State (interface)
  - insertCoin(m, coin)
  - selectProduct(m, slot)
  - cancel(m)
IdleState | HasMoneyState | DispensingState

VendingMachine {
  state, balance, inventory, coinInventory
  - InsertCoin / SelectSlot / Cancel / Dispense
}
```

**Patterns:** State (Idle / HasMoney / Dispensing) · Inventory + change calculator

## Step 3 — Flows

**Happy path**
1. Idle → InsertCoin → HasMoney (accumulate balance)  
2. SelectSlot → check price ≤ balance, stock > 0, change possible  
3. Dispense product → return change → decrement stock → Idle  

**Edge cases**
1. Insufficient funds / out of stock / cannot make change → stay HasMoney or refund  
2. Cancel → return inserted coins → Idle

## Step 4 — APIs

```text
InsertCoin(denomination)
SelectSlot(slotId)
Cancel()
AdminRefill(slotId, qty)
AdminAddCoins(denomination, count)
GetDisplay() → { balance, message }
```

## Step 5 — Deepen

- State pattern avoids a giant switch; illegal ops rejected per state  
- Insufficient change → refund, don’t dispense (atomic decision)  
- Single transaction serialized on one machine  
- Coin inventory must update with both inserted coins and change given  
- Admin refill concurrent with purchase → lock inventory

## Step 6 — Evolve

- New payment method (card) → extend states/transitions (**OCP**)  
- Related: [atm](../atm/README.md) (dispense chain), [wallet-ledger](../wallet-ledger/README.md) for cash tracking analogies
