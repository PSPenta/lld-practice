# Movie / seat booking — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Hold seat TTL?
2. Payment in scope?
3. Same seat concurrent book?
4. Show = screen + time?
5. Refund?
6. Partial confirm if some seats fail?
7. Seat types / pricing tiers?
8. Idempotency key on confirm?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Customer, BookingService, SeatLock |
| **Use cases (v1)** | 1. Browse available seats 2. Hold 3. Pay and confirm |
| **In scope** | Seat state, hold lock, confirm idempotent |
| **Out of scope** | Dynamic pricing, waitlist |
| **Assumptions** | Hold 10 min; row-level seat lock |

### Confirm understanding
> "User selects seats, holds them, pays, and booking confirms atomically."

## Step 2 — Entities & classes

```text
Movie { id, title }
Screen { id, theaterId, seatLayout }
Show { id, screenId, movieId, startTime }
Seat { id, showId, row, number, status: FREE|HELD|BOOKED }

SeatLock { seatIds[], userId, expiresAt, bookingId }
Booking {
  id, showId, userId, seatIds[], status: HELD|CONFIRMED|EXPIRED|CANCELLED
}

BookingService
  - getAvailable(showId)
  - holdSeats(showId, seatIds, userId) → bookingId
  - confirm(bookingId, paymentRef)
  - releaseExpired()
```

**Patterns:** Optimistic/pessimistic seat lock · TTL hold · State on Booking

## Step 3 — Flows

**Happy path**
1. GetAvailable → seats with status FREE  
2. HoldSeats → CAS FREE→HELD with expiresAt; create Booking HELD  
3. Pay (external) → Confirm → HELD→BOOKED seats; Booking CONFIRMED  
4. Sweeper releases HELD past TTL back to FREE  

**Edge cases**
1. Two users hold same seat → one wins lock; other gets conflict  
2. Confirm after expiry → reject; payment may need refund path

## Step 4 — APIs

```http
GET  /shows/{id}/seats
POST /shows/{id}/holds       { seatIds } → bookingId
POST /bookings/{id}/confirm  { paymentRef, Idempotency-Key }
DELETE /bookings/{id}        # cancel hold
```

## Step 5 — Deepen

- Row-level lock or compare-and-set on seat status prevents double book  
- Confirm is idempotent by booking id / Idempotency-Key  
- Never hold DB transaction open across payment HTTP call  
- TTL sweeper must be safe under concurrent confirm (CAS)  
- Partial hold failure → release any seats already held in that request

## Step 6 — Evolve

- Waitlist; dynamic pricing Strategy  
- Related: [hotel-booking](../hotel-booking/README.md), [restaurant-reservation](../restaurant-reservation/README.md), [inventory-management](../inventory-management/README.md)
