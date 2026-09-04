# Chess / board game — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Full rules or simplified?
2. Two players hot-seat?
3. Move validation only?
4. Undo?
5. AI opponent?
6. Castling / en passant / promotion in v1?
7. Clocks / timed games?
8. Persist game for resume?

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

## Step 2 — Entities & classes

```text
Color: WHITE | BLACK
Square { file, rank }  // or 0..63

Piece (abstract)
  King | Queen | Rook | Bishop | Knight | Pawn
  - color, possibleMoves(board, from) → []Square

Board {
  grid[8][8] Piece?
  - get / set / clone
}

Move { from, to, promoPiece? }
MoveValidator
  - isLegal(board, move, side)  // includes check filter

Game {
  board, turn, status: ONGOING|CHECKMATE|STALEMATE|DRAW
  history []Move
  - makeMove(from, to) error
  - getStatus() / undo()?
}
```

**Patterns:** Piece hierarchy (polymorphism) · Command for undo · Strategy later for AI

## Step 3 — Flows

**Happy path**
1. Start game → initial board setup; turn = WHITE  
2. Player proposes move → generate pseudo-legal moves for piece  
3. Filter moves that leave own king in check  
4. Apply move (captures, special rules) → switch turn  
5. If opponent has no legal moves → checkmate or stalemate  

**Edge cases**
1. Illegal move / wrong turn → reject, board unchanged  
2. Undo → pop history and restore prior board snapshot

## Step 4 — APIs

```text
Game.New() → *Game
Game.MakeMove(from, to, promo?) error
Game.GetStatus() → status, turn
Game.LegalMoves(from?) → []Move
Game.Undo() error                 // optional
```

## Step 5 — Deepen

- Immutable board or copy-on-write snapshots for undo / search  
- Legal-move generation must consider pins, checks, castling rights  
- Single-threaded v1; online multiplayer later needs turn auth  
- Validate promotion piece type; reject moves into check  
- Draw rules (threefold, fifty-move) as evolve if time permits

## Step 6 — Evolve

- Strategy for AI evaluation; Command pattern undo stack  
- Networked play + clocks  
- Related breadth: [lld-gaps](../../docs/lld-gaps/README.md) chess notes if present
