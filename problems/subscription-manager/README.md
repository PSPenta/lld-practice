# Subscription manager — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Plans free/pro?
2. Billing cycle?
3. Trial?
4. Proration?
5. Payment webhooks?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | User, SubscriptionService, PaymentProvider |
| **Use cases (v1)** | 1. Subscribe 2. Renew/cancel 3. Change plan |
| **In scope** | Plan, period, status ACTIVE/CANCELLED |
| **Out of scope** | Usage metering |
| **Assumptions** | Monthly billing; webhook idempotent |

### Confirm understanding
> "User picks plan; subscription active until period end or cancel."

---

## Step 2 — Entities & classes

`Subscription`, `Plan`, `BillingCycle`, `SubscriptionService`, `PaymentClient`

---

## Step 3 — Flows

Subscribe → active until period end → renew or cancel → downgrade/upgrade

---

## Step 4 — APIs

`Subscribe(userId, planId)`, `Cancel`, `ChangePlan`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Webhook idempotency from payment provider; state ACTIVE/PAST_DUE/CANCELLED

---

## Step 6 — Evolve

Usage-based add-on meter; dunning retries

---

