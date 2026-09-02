# Cab booking / Ride sharing — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Match rider to nearest driver?
2. Surge pricing?
3. Cancel before pickup?
4. Live GPS in v1?
5. Payment in scope?
6. One active trip per rider?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Rider, Driver, MatchingService, TripRepository |
| **Use cases (v1)** | 1. Request ride 2. Driver accepts 3. Complete trip |
| **In scope** | Create trip, assign driver, status transitions |
| **Out of scope** | Surge, multi-stop, chat |
| **Assumptions** | Simple nearest-driver match; sync accept |

### Confirm understanding
> "Rider requests ride; system assigns a nearby driver; driver accepts and trip completes."

---

## Step 2 — Entities & classes

`Rider`, `Driver`, `Trip`, `Location`, `PricingStrategy`, `MatchingService`, `TripRepository`

---

## Step 3 — Flows

Request ride → find nearby drivers → assign → driver accepts → trip in progress → complete → payment

---

## Step 4 — APIs

`POST /trips`, `PATCH /trips/{id}/status`, `GET /trips/{id}`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Double-book driver → lock driver row; cancel idempotency; driver location eventual consistency

---

## Step 6 — Evolve

Surge as Strategy; queue for matching at peak; split payment service

---

