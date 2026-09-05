# Expense splitter (Splitwise) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅  
> Service + pairwise ledger + integer paise + Factory (Equal / Exact / Percentage).

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Split types: equal, exact, percent?
2. Settle up between two users (full + partial)?
3. Groups in v1?
4. Multi-currency?
5. Who can add expenses?
6. Money precision (rupees vs paise)?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Users, `SplitwiseService`, `ExpenseFactory`, `BalanceSheet` |
| **Use cases (v1)** | Add expense · view pairwise balances · settle (full/partial) |
| **In scope** | Equal/exact/percent, pairwise debts, integer minor units |
| **Out of scope** | Real payments, multi-currency FX, simplify (optional evolve) |
| **Assumptions** | Single currency; API accepts rupees, ledger stores paise; Equal remainder→last, Percentage remainder→smallest |

### Confirm understanding
> "Users add shared expenses; we track who owes whom in paise; settle clears or reduces an edge."

---

## Step 2 — Entities & classes

```text
User { id, name, email }
Split { userId, amount?, percentage? }
Balance { debtorId, creditorId, amountPaise }

Expense (abstract / interface)
  ExactExpense | EqualExpense | PercentageExpense
  validate() → fills/checks paise amounts; apply(sheet) writes debts
  Equal: floor + remainder → last split
  Percentage: floor + remainder → smallest split (tie → first)
ExpenseFactory / CreateExpense(type, …)

BalanceSheet
  addDebt(debtor, creditor, paise)  // merge same edge; net opposite
  getBalance / getPairwiseBalances

SplitwiseService
  addUser, addExpense, settleUp, getPairwiseBalances
  toAmount(rupees) at API boundary
```

**Patterns:** **Factory** + Template Method (`validate`) · **Application service** · pairwise ledger

---

## Step 3 — Flows

1. `addExpense` → `toAmount(rupees)` → factory → `validate()` (fills/checks **paise** shares) → `apply` pairwise debts  
2. **Equal:** `floor(total/n)` each; **remainder on last** so Σ === total  
3. **Percentage:** percents must sum to 100; `floor` each share; **remainder on smallest share** (first index on ties)  
4. **Exact:** split amounts must already sum to expense (else throw)  
5. `settleUp(from, to, amount?)` → reverse `addDebt` to reduce/clear edge  
6. Opposite edges net inside `addDebt`; display uses `fromAmount(paise)`

---

## Step 4 — APIs

```text
addUser(name, email)
addExpense({ type, paidBy, amountRupees, splits })
getPairwiseBalances() → display strings (fromAmount)
settleUp(payerId, payeeId, amountRupees?)   // null/omit = full edge
```

---

## Step 5 — Deepen

- Integer **paise** in the ledger only — convert at API boundary; never `toFixed` as source of truth  
- Equal: remainder on **last**; Percentage: remainder on **smallest** — both keep Σ shares === total exactly  
- Exact validates sum of amounts === expense (demo rejects mismatches)  
- Reject non-positive settle / `addDebt` amounts  
- Concurrent adds → lock `BalanceSheet`

---

## Step 6 — Evolve

- New split type → new expense type + factory entry (**OCP**)  
- `simplify()` min transactions · groups · multi-currency


---

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/Splitwise/`](../../JavaScript/Splitwise/) | Preferred — `SplitwiseService`, pairwise `BalanceSheet`, `money.js` |
| **Go** | [`Go/Splitwise-go/`](../../Go/Splitwise-go/) | Same design (interfaces + embedding) |

## Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `money.js` / `money.go` | `toAmount(rupees)` → integer paise; `fromAmount` for display |
| `User.js` | User entity (`id`, `name`, `email`) |
| `Split.js` | `EqualSplit` / `ExactSplit` / `PercentageSplit` |
| `Balance.js` | Pairwise edge `{ debtorId, creditorId, amount }` |
| `BalanceSheet.js` | `addDebt` (merge same edge, net opposite), `getBalance` |
| `Expense.js` | Factory + Equal (remainder→**last**) / Percentage (remainder→**smallest**) / Exact |
| `SplitwiseService.js` | Façade: `addUser`, `addExpense`, `settleUp`, `getPairwiseBalances` |
| `index.js` / `main.go` | Demo: equal, percent, settle, Exact mismatch, final clear |

**Read order:** `SplitwiseService` → `Expense` factory → `BalanceSheet.addDebt` → `money.toAmount`.

