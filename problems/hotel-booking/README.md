# Hotel booking — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Room types?
2. Date-range search?
3. Overbooking?
4. Cancellation policy?
5. Payment when?

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

---

## Step 2 — Entities & classes

`Hotel`, `Room`, `Booking`, `AvailabilityService`, `DateRange`

---

## Step 3 — Flows

Search available rooms → hold → pay → confirm booking

---

## Step 4 — APIs

`Search(checkIn, checkOut, guests)`, `Book(roomId, dates)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Prevent double book same room+dates; idempotent book API

---

## Step 6 — Evolve

Similar patterns to [movie-booking](../movie-booking/README.md)

---

