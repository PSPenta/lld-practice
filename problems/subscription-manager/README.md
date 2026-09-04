# Subscription manager — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Plans free/pro?
2. Billing cycle?
3. Trial?
4. Proration?
5. Payment webhooks?
6. Cancel immediate vs end of period?
7. One active subscription per user?
8. Grace period / PAST_DUE?

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

## Step 2 — Entities & classes

```text
Plan { id, name, price, interval: MONTHLY|YEARLY }
BillingCycle { start, end }

Subscription {
  id, userId, planId
  status: TRIALING | ACTIVE | PAST_DUE | CANCELLED
  currentPeriodStart, currentPeriodEnd
  cancelAtPeriodEnd bool
}

SubscriptionService
  - Subscribe(userId, planId)
  - Cancel(subId, immediate?)
  - ChangePlan(subId, newPlanId)
  - HandleInvoicePaid / HandlePaymentFailed  // webhooks

PaymentClient
```

**Patterns:** State machine on subscription · Idempotent webhook handler · Adapter to payment provider

## Step 3 — Flows

**Happy path**
1. Subscribe → create ACTIVE (or TRIALING) with period dates; charge if needed  
2. Period end → renew charge → extend period  
3. Cancel at period end → flag `cancelAtPeriodEnd`; status CANCELLED when period ends  
4. ChangePlan → upgrade/downgrade (proration policy stated aloud)  

**Edge cases**
1. Payment failed → PAST_DUE; dunning retries; eventually CANCELLED  
2. Duplicate webhook → idempotent by event id

## Step 4 — APIs

```text
Subscribe(userId, planId) → subId
Cancel(subId, opts)
ChangePlan(subId, newPlanId)
GetSubscription(userId)
```

```http
POST /subscriptions
POST /subscriptions/{id}/cancel
POST /subscriptions/{id}/change-plan
POST /webhooks/payment          # provider events
```

## Step 5 — Deepen

- Webhook idempotency from payment provider (store event ids)  
- State ACTIVE / PAST_DUE / CANCELLED with legal transitions  
- One active sub per user enforced; concurrent ChangePlan → version lock  
- Don’t trust client for “paid” — trust signed webhooks  
- Renew job must be safe if run twice (idempotent period extension)

## Step 6 — Evolve

- Usage-based add-on meter; dunning retries  
- Proration strategies; tax  
- Related: [wallet-ledger](../wallet-ledger/README.md), [webhook-delivery](../webhook-delivery/README.md)
