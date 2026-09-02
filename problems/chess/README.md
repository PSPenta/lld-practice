# Chess / board game — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Full rules or simplified?
2. Two players hot-seat?
3. Move validation only?
4. Undo?
5. AI opponent?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Players, Game, Board, MoveValidator |
| **Use cases (v1)** | 1. Start game 2. Make legal move 3. Detect checkmate/stalemate |
| **In scope** | Board, pieces, legal moves, turn order |
| **Out of scope** | AI, online multiplayer |
| **Assumptions** | Full rules; two human players |

### Confirm understanding
> "Players alternate legal moves until win/draw condition."

---

## Step 2 — Entities & classes

`Board`, `Piece` hierarchy, `Move`, `Game`, `MoveValidator`

---

## Step 3 — Flows

Start → loop turns → generate legal moves → apply move → check/checkmate

---

## Step 4 — APIs

`Game.makeMove(from, to)`, `Game.getStatus()`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Immutable board or copy-on-write for undo

---

## Step 6 — Evolve

Strategy for AI; command pattern for undo stack

---

