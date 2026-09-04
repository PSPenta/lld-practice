# In-memory database — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Fixed schema or dynamic tables?
2. Indexes beyond primary key?
3. Transactions?
4. SQL subset or CRUD API?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Application, `Database`, `Table` |
| **Use cases (v1)** | Create table · insert/select/update/delete |
| **In scope** | Table map, PK, simple WHERE |
| **Out of scope** | JOIN, WAL, replication |
| **Assumptions** | Dynamic tables; in-memory rows |

### Confirm understanding
> "App creates tables and CRUDs rows keyed by id."

---

## Step 2 — Entities & classes

```text
Database { tables: Map<name, Table> }
Table { schema, rows[], indexes }
  insert(row), select(where), update(where, patch), delete(where)
```

**Pattern:** aggregate — `Database` owns `Table` instances

---

## Step 3 — Flows

**Insert:** validate schema → append row → update indexes  
**Select:** scan or index lookup by PK / indexed column  
**Delete/Update:** find rows → mutate → refresh indexes

---

## Step 4 — APIs

`createTable(name, schema)`, `insert(table, row)`, `select(table, where)`, …

---

## Step 5 — Deepen

- Lock per table for concurrent writers
- Schema validation on insert

---

## Step 6 — Evolve

- Secondary indexes; WAL for durability


---

## Code in this repo

| Language | Path |
|----------|------|
| **JavaScript** | [`JavaScript/Database/`](../../JavaScript/Database/) |
| **Go** | [`Go/Database-go/`](../../Go/Database-go/) |

## Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `Database.js` | Registry of tables — `createTable`, route CRUD to a `Table` |
| `Table.js` | Schema, rows, indexes; `insert` / `select` / `update` / `delete` |
| `index.js` | Demo create + CRUD |

**Read order:** `Database.createTable` → `Table.insert` / `select`.

