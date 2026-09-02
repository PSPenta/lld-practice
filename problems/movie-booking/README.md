# Movie / seat booking — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Hold seat TTL?
2. Payment in scope?
3. Same seat concurrent book?
4. Show = screen + time?
5. Refund?

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

---

## Step 2 — Entities & classes

`Show`, `Screen`, `Seat`, `Booking`, `SeatLock`, `BookingService`

---

## Step 3 — Flows

Select seats → hold lock TTL → pay → confirm → release lock on timeout

---

## Step 4 — APIs

`GetAvailable(showId)`, `HoldSeats(seatIds)`, `Confirm(bookingId)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Idempotency on confirm; row-level lock or compare-and-set on seat status

---

## Step 6 — Evolve

Waitlist; dynamic pricing Strategy

---

