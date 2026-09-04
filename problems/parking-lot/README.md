# Parking lot — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/ParkingLot2/`](../../JavaScript/ParkingLot2/) | composition chain — **preferred** |
| **Go** | [`Go/ParkingLot2-go/`](../../Go/ParkingLot2-go/) | |
| Teaching v1 | [`JavaScript/Parkinglot/`](../../JavaScript/Parkinglot/) | inheritance-heavy vehicles |

### Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `ParkingLot.js` | Orchestrator — `park` / `unpark` across floors |
| `Floor.js` | Owns slots; `findAvailableSlot(vehicleType)` |
| `Slot.js` | Free/occupied + vehicle type fit |
| `Vehicle.js` | Car/Bike/Truck + `VehicleFactory` |
| `Ticket.js` | Entry time + slot id for fee calc |
| `index.js` | Demo park/unpark |

**Read order:** `ParkingLot.park` → `Floor.findAvailableSlot` → `Slot` / `Ticket`.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Multiple floors?
2. Vehicle types vs slot types?
3. Pricing model?
4. Single entry/exit gate?
5. Behavior when lot full?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Driver, `ParkingLot`, ticket system |
| **Use cases (v1)** | 1. Park vehicle 2. Unpark and pay fee |
| **In scope** | Find slot, issue ticket, fee on exit |
| **Out of scope** | Reservations, EV charging |
| **Assumptions** | 2 floors; car/bike/truck; first-fit slot |

### Confirm understanding
> "Vehicle enters, gets ticket and slot; on exit pays and slot frees."

---

## Step 2 — Entities & classes

```text
Vehicle { number, type }           // Car | Bike | Truck via Factory
Slot { id, type, isFree, vehicle }
Floor { id, slots[]; findAvailableSlot(vehicleType) }
ParkingLot { floors[]; park(vehicle); unpark(ticketId) }
Ticket { id, slotId, entryTime }
PricingStrategy { calculate(ticket, exitTime) }   // optional v1
```

**Responsibilities:** find suitable slot · issue ticket on park · compute fee on unpark · track availability

**Patterns:** **Factory** (`VehicleFactory`) · **composition** (lot → floors → slots)

---

## Step 3 — Flows

**Park:**
1. `findAvailableSlot(vehicleType)` across floors
2. If none → reject (lot full)
3. Assign vehicle to slot; create `Ticket` with entry time

**Unpark:**
1. Load ticket → compute fee via `PricingStrategy`
2. Free slot; return amount

---

## Step 4 — APIs

- `park(vehicle): Ticket`
- `unpark(ticketId): fee`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

- Thread-safe slot assignment (lock lot or slot row)
- Full lot → fail fast before ticket issued
- Invalid ticket on unpark → clear error

---

## Step 6 — Evolve

- New vehicle type → `VehicleFactory` + slot compatibility map
- New pricing → new `PricingStrategy` (**Strategy** / **OCP**)
