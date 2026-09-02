# ATM / Cash dispenser — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single ATM or fleet?
2. Ops: withdraw, deposit, balance?
3. PIN validation in scope?
4. Multi-denomination cash dispense?
5. Concurrent sessions on one ATM?
6. Link to central bank or local ledger?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Customer, ATM hardware, AccountService |
| **Use cases (v1)** | 1. Withdraw cash 2. Check balance 3. Deposit (optional v1) |
| **In scope** | Card/PIN session, dispense chain, debit account |
| **Out of scope** | Fraud ML, cross-bank settlement |
| **Assumptions** | Single ATM; withdraw + balance; central AccountService |

### Confirm understanding
> "Customer inserts card, enters PIN, withdraws amount; ATM dispenses notes and debits account."

---

## Step 2 — Entities & classes

`ATM`, `CardReader`, `CashDispenser`, `DispenseChain` (100→50→20), `AccountService`, `Transaction`

---

## Step 3 — Flows

Insert card → PIN → select withdraw → check balance → dispense chain → commit debit → eject card

---

## Step 4 — APIs

`Withdraw(amount)`, `GetBalance()`, `Deposit(cash)` on `ATM` facade

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Concurrent sessions serialized; rollback if dispense fails mid-way; insufficient funds → clear error

---

## Step 6 — Evolve

State pattern for ATM session; add deposit envelope reader without rewriting withdraw flow

---

