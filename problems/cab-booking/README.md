# Cab booking / Ride sharing — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Match rider to nearest driver?
2. Surge pricing?
3. Cancel before pickup?
4. Live GPS in v1?
5. Payment in scope?
6. One active trip per rider?
7. Driver reject / re-match?
8. Geo radius and max wait for matching?

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

## Step 2 — Entities & classes

```text
Rider { id, name, location? }
Driver { id, location, status: AVAILABLE | BUSY | OFFLINE }
Location { lat, lng }
Trip {
  id, riderId, driverId?, pickup, dropoff
  status: REQUESTED | ASSIGNED | IN_PROGRESS | COMPLETED | CANCELLED
  fare?
}

MatchingService
  - findNearbyDrivers(pickup, radius) → []Driver
  - assign(tripId, driverId)
PricingStrategy (interface) → FlatPricing | DistancePricing
TripRepository, TripService
```

**Patterns:** Strategy (pricing) · State machine on `Trip` · Repository

## Step 3 — Flows

**Happy path**
1. Rider `POST /trips` with pickup/dropoff → status REQUESTED  
2. MatchingService finds nearby AVAILABLE drivers → offer/assign  
3. Driver accepts → ASSIGNED → starts trip → IN_PROGRESS  
4. Complete → COMPLETED → compute fare → (payment if in scope)  

**Edge cases**
1. No drivers in radius → reject or queue; driver declines → re-match  
2. Rider cancels before pickup → free driver; ignore late complete

## Step 4 — APIs

```http
POST   /trips                 { pickup, dropoff }
PATCH  /trips/{id}/status     { status }   # accept | start | complete | cancel
GET    /trips/{id}
POST   /drivers/{id}/location { lat, lng }
```

## Step 5 — Deepen

- Double-book driver → lock driver row / CAS on `status=AVAILABLE`  
- Cancel must be idempotent (repeat cancel = same terminal state)  
- Driver location is eventually consistent — matching uses last-known geo  
- One active trip per rider/driver enforced in service layer  
- Payment failure after COMPLETED → separate payment retry, don’t rewind trip state blindly

## Step 6 — Evolve

- Surge as `PricingStrategy` (**OCP**)  
- Async matching queue at peak; ETA / live GPS stream  
- Split payment into its own service  
- Related: [hotel-booking](../hotel-booking/README.md) (hold/confirm), [order-management](../order-management/README.md) (state machine)
