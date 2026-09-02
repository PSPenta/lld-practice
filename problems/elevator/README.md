# Elevator — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. How many floors and elevators?
2. Hall call vs destination dispatch?
3. Morning rush priority?
4. Weight capacity?
5. Maintenance mode?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Building users, ElevatorController, Elevator cars |
| **Use cases (v1)** | 1. Request up/down at floor 2. Move elevator 3. Open/close doors |
| **In scope** | Request queue, scheduler strategy, elevator state machine |
| **Out of scope** | Fire evacuation algorithm, smart building IoT |
| **Assumptions** | 3 elevators, 20 floors, FCFS or SCAN scheduler |

### Confirm understanding
> "Users press hall buttons; controller assigns an elevator and moves it to serve requests."

---

## Step 2 — Entities & classes

`ElevatorController`, `Elevator` (IDLE/MOVING/DOORS_OPEN), `Request`, `SchedulerStrategy` (FCFS/SCAN)

---

## Step 3 — Flows

Request arrives → controller assigns elevator → move → open doors → complete request

---

## Step 4 — APIs

`RequestFloor(elevatorId, floor, direction)`, `GetStatus()`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Thread-safe request queues; don't starve opposite direction (SCAN helps)

---

## Step 6 — Evolve

Swap scheduler Strategy; weight capacity; express elevator

---

