# Payment gateway — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅ (Go)

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **Go** | [`Go/PaymentGateway-go/`](../../Go/PaymentGateway-go/) | **full impl** — `BankGateway` Strategy |
| JavaScript | [`JavaScript/PaymentGateway/`](../../JavaScript/PaymentGateway/) | stub only (empty files) |

REST/idempotency depth → [Backend/README.md §4](../../Backend/README.md)

### Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `bank_gateway.go` | `BankGateway` interface — `ProcessPayment` |
| `payment_gateway.go` | Maps payment method → gateway; idempotency store |
| `payment.go` | Payment model + status |
| `main.go` | Demo process + retry with same key |

**Read order:** `PaymentGateway.ProcessPayment` → idempotency check → `BankGateway` impl.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Providers (UPI, card, Razorpay)?
2. Idempotent pay on retry?
3. Refunds in v1?
4. Webhook callbacks?
5. Sync confirm or async?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Merchant, `PaymentGateway`, `BankGateway` providers |
| **Use cases (v1)** | Process payment · route to provider · safe retry |
| **In scope** | Strategy per provider, idempotency store |
| **Out of scope** | PCI vault, full 3DS |
| **Assumptions** | `Idempotency-Key`; map method → gateway |

### Confirm understanding
> "Same idempotency key on retry returns same result — never double-charge."

---

## Step 2 — Entities & classes

```text
BankGateway (interface) { processPayment(payment): Result }

PaymentGateway
  - gateways: Map<PaymentMethod, BankGateway>

Payment { id, amount, method, idempotencyKey, status }
```

**Pattern:** **Strategy** + **DIP** — gateway depends on `BankGateway` interface

---

## Step 3 — Flows

1. `ProcessPayment(idempotencyKey, amount, method)`
2. If key seen → return stored response
3. Else pick `BankGateway` from map → call provider
4. Persist result keyed by idempotency key

---

## Step 4 — APIs

- `processPayment(idempotencyKey, amount, method)`
- HTTP: `POST /payments` + header `Idempotency-Key`

---

## Step 5 — Deepen

- Timeout → unknown state; client retries with **same** key
- Provider failure → structured error; no partial double charge

---

## Step 6 — Evolve

- New provider → new `BankGateway` impl registered in map (**OCP**)
