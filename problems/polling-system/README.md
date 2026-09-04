# Polling / voting system — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/PollingSystem2/`](../../JavaScript/PollingSystem2/) | `PollingService/` + `models/` + `repositories/` |
| **Go** | [`Go/PollingSystem2-go/`](../../Go/PollingSystem2-go/) | |
| Refactor drill | [`JavaScript/PollingSystem/`](../../JavaScript/PollingSystem/) | v1 — spot SRP gaps |

### Codebase map (how the code is organized)

| Path | Responsibility |
|------|----------------|
| `PollingService/models/` | Thin `User`, `Poll`, `Vote` entities |
| `PollingService/repositories/` | In-memory `User` / `Poll` / `Vote` stores |
| `PollingService/index.js` | Use-cases: create poll, vote, close, results + auth rules |
| `PollingSystem2/index.js` | Demo wiring |
| `PollingSystem/` (v1) | Teaching anti-example — compare layering |

**Read order:** `PollingService` methods → repositories → models.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Public vs private polls?
2. One vote per user?
3. Creator can close poll?
4. Anonymous votes?
5. Max options?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | User, `PollingService`, repositories |
| **Use cases (v1)** | Create poll · cast vote · view results |
| **In scope** | Dedupe voter, tally, close poll |
| **Out of scope** | Weighted voting |
| **Assumptions** | One vote per (poll, user); layered design |

### Confirm understanding
> "Creator opens poll; each user votes once; results tally live."

---

## Step 2 — Entities & classes

```text
models/: User, Poll, Option, Vote
repositories/: PollRepository, VoteRepository
PollingService
  - createPoll, vote, closePoll, getResults
  - rules: poll open, user != creator if required, one vote
```

**Patterns:** **Application service** + **Repository** · thin models

---

## Step 3 — Flows

**Create poll:** validate user → save poll + options  

**Vote:** load poll → if closed reject → if already voted reject → save vote → increment tally  

**Close:** owner marks `isClosed`

---

## Step 4 — APIs

`createPoll(...)`, `vote(pollId, optionId, userId)`, `getResults(pollId)`

---

## Step 5 — Deepen

- Unique constraint `(pollId, userId)`
- Concurrent votes → transaction or row lock

---

## Step 6 — Evolve

- Compare v1 `PollingSystem/` → v2 layering in `PollingSystem2/`
