# Wallet / ledger — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Balance-only or double-entry?
2. Overdraft?
3. Multi-currency?
4. Audit trail?
5. Transfer idempotency?
6. Hold / authorize then capture?
7. Integer minor units (paise/cents)?
8. Concurrent transfers same wallet?

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

## Step 2 — Entities & classes

```text
Wallet { id, userId, currency }
LedgerEntry {
  id, walletId, amountSigned, type: CREDIT|DEBIT
  transferId?, createdAt
}
Transfer {
  id (idempotency), fromWallet, toWallet, amount, status
}

WalletService
  - Credit(walletId, amount, ref)
  - Debit(walletId, amount, ref)
  - Transfer(from, to, amount, transferId)
  - GetBalance(walletId)

Ledger
  - append(entries...)   // atomic for a transfer
  - sum(walletId) → balance
```

**Patterns:** Append-only ledger · Idempotency key on transfer · Double-entry (debit+credit pair)

## Step 3 — Flows

**Happy path — transfer**
1. Validate amount > 0; check idempotency key → return prior result if seen  
2. Lock wallets in stable id order (avoid deadlock)  
3. Ensure from balance ≥ amount (no overdraft)  
4. Append DEBIT on from + CREDIT on to in one transaction  
5. Mark transfer SUCCEEDED  

**Edge cases**
1. Insufficient funds → reject; no partial entries  
2. Retry same transferId → no double move

## Step 4 — APIs

```text
Transfer(from, to, amount, transferId) error
Credit(walletId, amount, ref)
GetBalance(userId|walletId) → amount
ListEntries(walletId, limit) → []LedgerEntry
```

```http
POST /wallets/{id}/transfers
GET  /wallets/{id}/balance
```

## Step 5 — Deepen

- Serializable / row locks; always lock wallet ids in sorted order  
- Idempotent transfer id stored uniquely  
- Immutable entries — never update balance without an entry  
- Integer minor units only (no float)  
- Balance = sum(entries) or cached balance updated in same txn

## Step 6 — Evolve

- Related: Splitwise pairwise balances; event sourcing for audit  
- Holds / authorize-capture; multi-currency books  
- Related: [order-management](../order-management/README.md), [subscription-manager](../subscription-manager/README.md)
