# Task board (Trello-like) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Common Razorpay-style machine coding prompt: boards, lists, cards, drag between columns.

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single user or multi-user boards?
2. Roles (owner, member, viewer)?
3. Move card between lists only, or reorder within a list?
4. Card fields beyond title (description, assignee, due date)?
5. Archive vs delete?
6. Concurrent edits on same board?
7. Default lists on board create?
8. Soft delete / undo move?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | User, `BoardService`, repositories |
| **Use cases (v1)** | Create board · add lists · add cards · move card to another list |
| **In scope** | Board → lists (ordered) → cards (ordered per list) |
| **Out of scope** | Comments, attachments, notifications, real-time sync |
| **Assumptions** | One owner; members can edit; move = remove from old list + append to new |

### Confirm understanding
> "Users create a board with columns like To Do / In Progress / Done; cards move between columns."

## Step 2 — Entities & classes

```text
User { id, name }
Board { id, title, ownerId, memberIds[] }
ListColumn { id, boardId, title, position }
Card { id, listId, title, position, assigneeId? }

BoardRepository, ListRepository, CardRepository
BoardService
  - createBoard, addList, addCard
  - moveCard(cardId, toListId, position?)
```

**Patterns:** **Repository** · application **service** for use-cases · ordered collections (position index)

## Step 3 — Flows

**Happy path — create board**
1. Validate user → persist board  
2. Optionally seed default lists (To Do / Doing / Done)  

**Happy path — add / move card**
1. Add card: verify list exists → `position = max+1` in that list  
2. Move card: load card → remove from source (recompact positions) → insert at target position → persist atomically  
3. Permission check: only board members may mutate  

**Edge cases**
1. Move to list on another board → reject  
2. Concurrent moves → lock board/card; conflict on stale position

## Step 4 — APIs

```http
POST   /boards                    { title }
POST   /boards/{id}/lists         { title }
POST   /lists/{id}/cards          { title }
PATCH  /cards/{id}/move           { toListId, position }
GET    /boards/{id}               → board + lists + cards
```

## Step 5 — Deepen

- **Concurrency:** lock board or card row on move; reject stale `position` with conflict error  
- **Idempotency:** duplicate move with same target is no-op  
- **Validation:** card cannot belong to two lists; list must belong to same board as card  
- Positions: integer recompact vs fractional indices — state trade-off  
- AuthZ on every mutating API

## Step 6 — Evolve

- Drag reorder within list → `reorderCard(cardId, newPosition)`  
- Activity log / audit trail; role-based permissions (viewer read-only)  
- WebSocket fan-out for live board updates → [pub-sub](../pub-sub/README.md)  
- Compare FIFO-only [task-queue](../task-queue/README.md) after your design
