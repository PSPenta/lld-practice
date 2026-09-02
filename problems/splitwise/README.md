# Expense splitter (Splitwise) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/Splitwise/`](../../JavaScript/Splitwise/) | Factory + expense subclasses |
| **Go** | [`Go/Splitwise-go/`](../../Go/Splitwise-go/) | |

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Split types: equal, exact, percent?
2. Settle up between two users?
3. Groups?
4. Multi-currency?
5. Who can add expenses?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Users, `ExpenseFactory`, `BalanceSheet` |
| **Use cases (v1)** | Add expense · view balances · settle debt |
| **In scope** | Split strategies, balance tracking |
| **Out of scope** | Payment integration |
| **Assumptions** | Single currency; equal/exact/percent |

### Confirm understanding
> "Users add shared expenses; balances update; settle-up clears debt."

---

## Step 2 — Entities & classes

```text
User { id, name }
Expense (abstract) { amount, paidBy, splits[]; validate() }
  ExactExpense | EqualExpense | PercentageExpense
ExpenseFactory.create(type, ...)
BalanceSheet { balances: Map<userId, amount> }
```

**Patterns:** **Factory** · Template Method–like `validate()` on subclasses

---

## Step 3 — Flows

1. `AddExpense` → factory creates typed expense → `validate()` (splits sum to total)
2. Update `BalanceSheet` (who owes whom)
3. `Settle(from, to)` → transfer balance between two users

---

## Step 4 — APIs

- `addExpense(...)`, `getBalance(userId)`, `settle(from, to)`

---

## Step 5 — Deepen

- Validate splits sum to total before commit
- Lock balance sheet on concurrent adds

---

## Step 6 — Evolve

- New split type → new `Expense` subclass + factory entry (**OCP**)
