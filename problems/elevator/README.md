# Elevator — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. How many floors and elevators?
2. Hall call vs destination dispatch?
3. Morning rush priority?
4. Weight capacity?
5. Maintenance mode?
6. Internal cabin buttons vs external hall buttons?
7. Door obstruction / reopen?
8. Single controller or distributed per car?

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

## Step 2 — Entities & classes

```text
Elevator {
  id, currentFloor, direction: UP|DOWN|IDLE
  state: IDLE | MOVING | DOORS_OPEN
  pendingStops: set/queue of floors
  - step() / openDoors() / closeDoors()
}

Request { floor, direction?, elevatorId? }  // hall vs cabin

SchedulerStrategy (interface)
  - assign(request, elevators[]) → elevatorId
  FCFSScheduler | SCANScheduler

ElevatorController
  - requestHall(floor, direction)
  - requestCabin(elevatorId, floor)
  - tick()   // advance simulation / hardware loop
```

**Patterns:** Strategy (scheduler) · State machine per car · Controller orchestrates

## Step 3 — Flows

**Happy path**
1. Hall request arrives → controller asks SchedulerStrategy for a car  
2. Car adds floor to pending stops  
3. On tick: move toward next stop → open doors → complete request → close  
4. If more stops in same direction, continue (SCAN)  

**Edge cases**
1. All cars busy / maintenance → queue request; don’t drop silently  
2. Door obstruction → reopen; opposite-direction starvation mitigated by SCAN

## Step 4 — APIs

```text
RequestHall(floor, direction)
RequestCabin(elevatorId, floor)
GetStatus() → []{ id, floor, state, direction }
Tick()                    // or event-driven MoveComplete
```

## Step 5 — Deepen

- Thread-safe request queues; single writer for car state or lock per elevator  
- Don’t starve opposite direction — SCAN / LOOK scheduling helps  
- Idempotent duplicate hall presses (same floor+direction already pending)  
- Fail-safe: on controller crash, cars stop or go to nearest floor (policy)  
- Capacity / weight: reject new cabin entry when over limit

## Step 6 — Evolve

- Swap scheduler Strategy without rewriting cars (**OCP**)  
- Weight capacity, express elevator, destination dispatch  
- Related: [traffic-signal](../traffic-signal/README.md) (timed state machines)
