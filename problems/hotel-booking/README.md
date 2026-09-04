# Hotel booking — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Room types?
2. Date-range search?
3. Overbooking?
4. Cancellation policy?
5. Payment when?
6. Hold TTL before confirm?
7. Inventory per room instance or per room type?
8. Idempotency on book?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Guest, AvailabilityService, Booking |
| **Use cases (v1)** | 1. Search rooms 2. Hold/book dates 3. Cancel |
| **In scope** | Room-night availability, book idempotent |
| **Out of scope** | Dynamic pricing |
| **Assumptions** | No overbooking; pay on book |

### Confirm understanding
> "Guest searches dates, books a room type if available."

## Step 2 — Entities & classes

```text
Hotel { id, name }
RoomType { id, hotelId, name, capacity, totalRooms }
Room { id, hotelId, roomTypeId }           // optional concrete rooms

DateRange { checkIn, checkOut }            // nights [checkIn, checkOut)

Booking {
  id, hotelId, roomTypeId, roomId?
  guestId, range, status: HELD|CONFIRMED|CANCELLED
  expiresAt?
}

AvailabilityService
  - search(hotelId, range, guests) → []RoomTypeAvailability
BookingService
  - hold / book(range, roomTypeId, idempotencyKey)
  - cancel(bookingId)
```

**Patterns:** Date-range inventory · Hold TTL · Idempotent book

## Step 3 — Flows

**Happy path**
1. Search → for each room type, count booked nights overlapping range → available = total − booked  
2. Hold/Book → transactionally reserve capacity for each night (or assign room)  
3. Pay → CONFIRM  
4. Cancel → free nights per policy  

**Edge cases**
1. Two guests book last room-night → one wins unique constraint / version  
2. Hold expires before pay → release; confirm after expiry rejected

## Step 4 — APIs

```http
GET  /hotels/{id}/availability?checkIn=&checkOut=&guests=
POST /hotels/{id}/bookings     { roomTypeId, checkIn, checkOut, Idempotency-Key }
DELETE /bookings/{id}
GET  /bookings/{id}
```

```text
Search(checkIn, checkOut, guests)
Book(roomTypeId, dates, idempotencyKey)
Cancel(bookingId)
```

## Step 5 — Deepen

- Prevent double book same room+dates (unique per room-night or atomic counter)  
- Idempotent book API by client key  
- Inclusive/exclusive date math — state night boundaries aloud  
- Don’t hold DB txn across payment; use HELD + sweeper  
- Cancellation policy windows enforced in service

## Step 6 — Evolve

- Similar patterns to [movie-booking](../movie-booking/README.md) and [restaurant-reservation](../restaurant-reservation/README.md)  
- Dynamic pricing Strategy; overbooking policy  
- Multi-property search / allotment
