# ATM / Cash dispenser — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single ATM or fleet?
2. Ops: withdraw, deposit, balance?
3. PIN validation in scope?
4. Multi-denomination cash dispense?
5. Concurrent sessions on one ATM?
6. Link to central bank or local ledger?
7. Daily withdrawal limits?
8. Receipt / mini-statement in v1?

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

## Step 2 — Entities & classes

```text
ATM
  - insertCard(card) → Session
  - enterPin(pin) → ok/fail
  - withdraw(amount) / getBalance() / ejectCard()

Session { cardNumber, authenticated, state }
CardReader, CashDispenser, Keypad, Screen

DispenseChain (Chain of Responsibility)
  NoteHandler(100) → NoteHandler(50) → NoteHandler(20)
  - dispense(amount) → note counts or error

AccountService (remote)
  - getBalance(accountId)
  - debit(accountId, amount, txnId)   // idempotent txnId

Transaction { id, type, amount, status }
```

**Patterns:** State (session lifecycle) · Chain of Responsibility (denominations) · Facade (`ATM`)

## Step 3 — Flows

**Happy path (withdraw)**
1. Insert card → create session  
2. Enter PIN → authenticate with AccountService  
3. Select withdraw → enter amount  
4. Check balance / limits → run dispense chain  
5. Commit debit with txn id → eject card  

**Edge cases**
1. Insufficient funds or cannot make exact change with available notes → abort, no debit  
2. Dispense hardware fails after debit attempt → compensate/rollback using txn id; keep card if bank says so

## Step 4 — APIs

```text
ATM facade:
  InsertCard(cardNumber)
  EnterPin(pin)
  Withdraw(amount)
  GetBalance()
  Deposit(cash)          // optional v1
  EjectCard() / Cancel()
```

## Step 5 — Deepen

- Serialize sessions on one ATM — one authenticated session at a time  
- Dispense-then-debit vs debit-then-dispense: prefer debit with compensating credit if dispense fails  
- Idempotent debit via transaction id so retries don’t double-charge  
- Insufficient funds / wrong PIN → clear errors; lock card after N PIN failures (policy)  
- Cash cassette inventory must stay consistent under concurrent admin refill (ops lock)

## Step 6 — Evolve

- State pattern for Idle → CardInserted → Authenticated → Dispensing  
- Add deposit / transfer without rewriting withdraw (**OCP** on session states)  
- Fleet: central monitoring, cassette low alerts  
- Related: [wallet-ledger](../wallet-ledger/README.md) for ledger semantics
