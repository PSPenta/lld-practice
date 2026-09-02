# Parking Lot — LLD design walkthrough

> **Code:** this folder (**preferred**) · v1 teaching: [../Parkinglot/](../Parkinglot/) · **Go:** [../../Go/ParkingLot2-go/](../../Go/ParkingLot2-go/)  
> **Method:** [../../README.md §5](../../README.md)

---

## Clarify

- Multiple floors?  
- Vehicle types & slot types?  
- Pricing?  
- Entry/exit gates count?

---

## Entities

`ParkingLot`, `Floor`, `Slot`, `Vehicle`, `Ticket`

---

## Responsibilities

- Find suitable free slot (Strategy: first-fit, type-fit, nearest)  
- Issue ticket on park  
- Calculate fee on unpark  
- Maintain availability counts  

---

## Classes (sketch)

```text
Vehicle { number, type }
Slot { id, type, isFree, vehicle }
Floor { id, slots[]; FindSlot(vehicleType) }
ParkingLot { floors[]; Park(vehicle); Unpark(ticketId) }
Ticket { id, slotId, entryTime }
PricingStrategy { Calculate(ticket, exitTime) }
```

---

## Extensibility

New vehicle type → mapping to slot types.  
New pricing → new `PricingStrategy`.
