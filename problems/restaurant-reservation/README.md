# Restaurant table reservation — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Common Razorpay-style machine coding prompt: check availability, book slot, avoid double booking.

## Code in this repo

No dedicated implementation yet. Same booking patterns as [hotel-booking](../hotel-booking/README.md) and [movie-booking](../movie-booking/README.md) — compare after your design.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Fixed table inventory or flexible party size → table assignment?
2. Reservation slot granularity (15 min, 30 min, 1 hour)?
3. Walk-ins reduce available tables?
4. Cancellation window?
5. Hold table before confirm (TTL hold)?
6. Same customer two bookings same slot?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Guest, `Restaurant`, `ReservationService` |
| **Use cases (v1)** | Check availability · create reservation · cancel |
| **In scope** | Tables, time slots, party size, conflict detection |
| **Out of scope** | Waitlist, menu, payment |
| **Assumptions** | Fixed tables; one reservation per table per slot; no overbooking |

### Confirm understanding
> "Guest picks date, time, and party size; system assigns an available table or rejects if full."

---

## Step 2 — Entities & classes

```text
Restaurant { id, name, openingHours }
Table { id, restaurantId, capacity, zone? }
TimeSlot { date, startTime, endTime }   // or store as instant range

Reservation { id, restaurantId, tableId, guestName, partySize, slot, status }
  status: HELD | CONFIRMED | CANCELLED

AvailabilityService
  - findAvailableTables(restaurantId, slot, partySize)
ReservationService
  - createReservation(...)   // atomic assign table
  - cancel(reservationId)
```

**Pattern:** **Repository** + domain service; slot overlap as core invariant

---

## Step 3 — Flows

**Check availability:** load tables with `capacity >= partySize` → filter out tables with overlapping CONFIRMED/HELD reservation for slot  

**Book:** within transaction — re-check availability → insert reservation → commit (fail if race)  

**Cancel:** mark CANCELLED; table freed for slot  

**Optional hold:** HELD + expiry job releases table if not confirmed

---

## Step 4 — APIs

```http
GET  /restaurants/{id}/availability?date=&time=&partySize=
POST /restaurants/{id}/reservations   { guestName, partySize, date, time }
DELETE /reservations/{id}
```

---

## Step 5 — Deepen

- **Double booking:** DB unique index on `(tableId, slotStart)` or transactional re-check
- **Overbooking:** reject when no table fits; optional waitlist in v2
- **Idempotency:** same client token on POST returns same reservation
- **Edge cases:** party larger than any single table → v2 combine tables; closing time overlap

---

## Step 6 — Evolve

- Dynamic table merging for large parties
- SMS reminder → [ticket-notify](../ticket-notify/README.md)
- Peak-hour slot limits → [rate-limiter](../rate-limiter/README.md) on booking API
