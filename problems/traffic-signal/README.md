# Traffic light / signal — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single intersection?
2. Pedestrian crossing?
3. Emergency override?
4. Fixed timing or sensor?
5. Coordination with neighbors?

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

---

## Step 2 — Entities & classes

`TrafficLight`, `Intersection`, **State** (NS_GREEN, EW_GREEN, ALL_RED)

---

## Step 3 — Flows

Timer tick → transition state → allow/block directions

---

## Step 4 — APIs

`Tick()`, `GetSignal(direction)`, `EmergencyOverride()`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Thread-safe state; fail-safe ALL_RED

---

## Step 6 — Evolve

Coordinator for green wave; strategy for timing plans

---

