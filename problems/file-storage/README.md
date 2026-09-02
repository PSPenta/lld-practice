# File storage (in-memory FS) — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Files + folders only?
2. Max file size?
3. Symlinks?
4. Permissions model?
5. Path separator rules?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Client code, FileSystem API |
| **Use cases (v1)** | 1. Create/read file 2. mkdir -p 3. List/delete by path |
| **In scope** | Composite File/Directory, path resolve |
| **Out of scope** | Disk persistence, ACLs |
| **Assumptions** | In-memory tree; Unix-style paths |

### Confirm understanding
> "Clients create paths like /a/b.txt, read/write files, and list directories."

---

## Step 2 — Entities & classes

`FileSystemNode` interface, `File`, `Directory` (Composite), `PathResolver`

---

## Step 3 — Flows

Create path → mkdir -p → write file → read by path → list dir → delete

---

## Step 4 — APIs

`CreateFile(path, data)`, `ReadFile(path)`, `ListDir(path)`, `Delete(path)`

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Path traversal guard; concurrent writes to same file → lock

---

## Step 6 — Evolve

Add metadata, permissions; persist to disk adapter

---

