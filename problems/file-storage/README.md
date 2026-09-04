# File storage (in-memory FS) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Files + folders only?
2. Max file size?
3. Symlinks?
4. Permissions model?
5. Path separator rules?
6. Overwrite vs create-exclusive?
7. Move/rename in v1?
8. Case sensitivity?

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

## Step 2 — Entities & classes

```text
FileSystemNode (interface)
  - Name(), IsDir(), Children()? 

File { name, data []byte, updatedAt }
Directory { name, children map[string]FileSystemNode }

PathResolver
  - resolve(path) → (parent Dir, name, node)
  - mkdirAll(path)

FileSystem
  - CreateFile(path, data)
  - ReadFile(path) → data
  - Mkdir(path)
  - ListDir(path) → names
  - Delete(path)
```

**Patterns:** Composite (File/Directory) · Facade (`FileSystem`)

## Step 3 — Flows

**Happy path**
1. Normalize path (`/a/b/c.txt`)  
2. `mkdir -p` parents as needed  
3. Create/write file node under parent directory  
4. Read by resolving path → return bytes  
5. ListDir / Delete by path  

**Edge cases**
1. Path traversal (`/../etc`) → reject; file where dir expected → error  
2. Delete non-empty dir → fail or recursive flag; concurrent write → lock file node

## Step 4 — APIs

```text
CreateFile(path, data) error
ReadFile(path) ([]byte, error)
Mkdir(path) error
ListDir(path) ([]string, error)
Delete(path) error
Exists(path) bool
```

## Step 5 — Deepen

- Path traversal / absolute-path guards before mutate  
- Concurrent writes to same file → per-node mutex or FS-wide lock for v1  
- Idempotent mkdir; create-exclusive vs overwrite policy stated aloud  
- Bound max file size / total memory to avoid OOM  
- Delete file vs directory semantics must be explicit

## Step 6 — Evolve

- Metadata, permissions, move/rename  
- Persist via disk adapter (**OCP** on storage backend)  
- Related: [inventory-management](../inventory-management/README.md) only for locking analogies — FS is Composite-heavy
