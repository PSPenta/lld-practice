# Wallet / ledger — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Balance-only or double-entry?
2. Overdraft?
3. Multi-currency?
4. Audit trail?
5. Transfer idempotency?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Users, WalletService, Ledger |
| **Use cases (v1)** | 1. Credit/debit 2. Transfer 3. Query balance |
| **In scope** | Append entries, compute balance, idempotent transfer id |
| **Out of scope** | FX, interest |
| **Assumptions** | Single currency; serializable transfers |

### Confirm understanding
> "Transfers move balance between wallets with immutable ledger entries."

---

## Step 2 — Entities & classes

`Wallet`, `LedgerEntry`, `Transaction`, `WalletService`

---

## Step 3 — Flows

Credit/debit → append entry → compute balance from entries or cached balance

---

## Step 4 — APIs

`Transfer(from, to, amount)`, `GetBalance(userId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Serializable transactions; idempotent transfer id

---

## Step 6 — Evolve

Related: Splitwise balances; event sourcing for audit

---

