# Restaurant table reservation — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Common Razorpay-style machine coding prompt: check availability, book slot, avoid double booking.

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Fixed table inventory or flexible party size → table assignment?
2. Reservation slot granularity (15 min, 30 min, 1 hour)?
3. Walk-ins reduce available tables?
4. Cancellation window?
5. Hold table before confirm (TTL hold)?
6. Same customer two bookings same slot?
7. Overbooking policy?
8. Timezone / closing-time overlap?

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

## Step 3 — Flows

**Happy path — check availability**
1. Load tables with `capacity >= partySize`  
2. Filter out tables with overlapping CONFIRMED/HELD reservation for slot  
3. Return remaining tables / boolean available  

**Happy path — book**
1. Within transaction — re-check availability  
2. Insert reservation (assign table) → commit  
3. On race, fail and ask client to retry  

**Edge cases**
1. Cancel → mark CANCELLED; table freed for slot  
2. Optional hold: HELD + expiry job releases if not confirmed; closing-time overlap rejected

## Step 4 — APIs

```http
GET  /restaurants/{id}/availability?date=&time=&partySize=
POST /restaurants/{id}/reservations   { guestName, partySize, date, time }
DELETE /reservations/{id}
```

## Step 5 — Deepen

- **Double booking:** DB unique index on `(tableId, slotStart)` or transactional re-check  
- **Overbooking:** reject when no table fits; optional waitlist in v2  
- **Idempotency:** same client token on POST returns same reservation  
- **Edge cases:** party larger than any single table → v2 combine tables  
- Concurrent book of last table → one commit wins, other gets conflict

## Step 6 — Evolve

- Dynamic table merging for large parties  
- SMS reminder → [ticket-notify](../ticket-notify/README.md)  
- Peak-hour slot limits → [rate-limiter](../rate-limiter/README.md) on booking API  
- Same hold patterns as [hotel-booking](../hotel-booking/README.md) and [movie-booking](../movie-booking/README.md)
