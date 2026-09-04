# Traffic light / signal — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single intersection?
2. Pedestrian crossing?
3. Emergency override?
4. Fixed timing or sensor?
5. Coordination with neighbors?
6. Left-turn protected phases?
7. Yellow / all-red clearance times?
8. Fail-safe behavior on controller fault?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Traffic controller, Intersection |
| **Use cases (v1)** | 1. Cycle NS/EW green 2. All-red clearance 3. Emergency all-stop |
| **In scope** | State machine per approach |
| **Out of scope** | City-wide coordination |
| **Assumptions** | One intersection; timed phases |

### Confirm understanding
> "Lights cycle safely NS then EW with all-red between."

## Step 2 — Entities & classes

```text
Direction: NS | EW
SignalColor: RED | YELLOW | GREEN

TrafficLight { direction, color }

Phase {
  name, allowedGreens[], duration
  // e.g. NS_GREEN, NS_YELLOW, ALL_RED, EW_GREEN, ...
}

Intersection {
  lights map[Direction]TrafficLight
  currentPhase, phaseStartedAt
  - Tick(now)
  - GetSignal(direction) → color
  - EmergencyOverride()      // all red / flash
}

TimingPlan / PhaseSchedule
```

**Patterns:** State machine (phases) · Strategy later for timing plans / sensor-actuated

## Step 3 — Flows

**Happy path**
1. Controller starts in ALL_RED or first phase  
2. Tick: if phase duration elapsed → transition to next phase  
3. Update each TrafficLight color from phase  
4. Repeat cycle NS green → yellow → all-red → EW green → …  

**Edge cases**
1. EmergencyOverride → force ALL_RED (or flash); resume plan when cleared  
2. Clock jump / missed ticks → recompute phase from elapsed time, don’t skip clearance

## Step 4 — APIs

```text
Tick(now)
GetSignal(direction) → RED|YELLOW|GREEN
EmergencyOverride(on bool)
SetTimingPlan(plan)           // evolve
```

## Step 5 — Deepen

- Thread-safe phase updates; single ticker goroutine/thread preferred  
- Fail-safe ALL_RED on unknown state / hardware fault  
- Never allow NS and EW green simultaneously — invariant test  
- Yellow + all-red clearance mandatory between opposing greens  
- Idempotent emergency on/off

## Step 6 — Evolve

- Coordinator for green wave across intersections  
- Strategy for fixed vs sensor-actuated timing plans  
- Pedestrian phases; protected left turns  
- Related: [elevator](../elevator/README.md) (timed state machines)
