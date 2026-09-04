# Search engine — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/SearchEngine/`](../../JavaScript/SearchEngine/) | composition pipeline |
| **Go** | [`Go/SearchEngine-go/`](../../Go/SearchEngine-go/) | |

### Codebase map (how the code is organized)

| File | Responsibility |
|------|----------------|
| `SearchEngine.js` | Facade — wires tokenizer, index, trie, ranker |
| `Tokenizer.js` | Text → terms |
| `InvertedIndex.js` | term → posting list of doc ids |
| `Trie.js` | Prefix / autocomplete |
| `Ranker.js` | Score and order results |
| `index.js` | Demo index + search + suggest |

**Read order:** `SearchEngine.index` / `search` → each pipeline piece.

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. In-memory index size?
2. Autocomplete / prefix search?
3. Ranking signals?
4. Incremental index updates?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Indexer, searcher, query API |
| **Use cases (v1)** | Index doc · search ranked · prefix suggest |
| **In scope** | Tokenizer, inverted index, trie, ranker |
| **Out of scope** | Distributed sharding |
| **Assumptions** | Single node; simple term frequency ranker |

### Confirm understanding
> "Documents indexed by terms; queries return ranked document ids."

---

## Step 2 — Entities & classes

```text
SearchEngine (facade / orchestrator)
  - tokenizer: Tokenizer
  - trie: Trie                    // autocomplete
  - index: InvertedIndex
  - ranker: Ranker

Tokenizer.tokenize(text) → terms[]
InvertedIndex.add(docId, terms)
Ranker.rank(queryTerms) → docIds[]
```

**Pattern:** **composition** pipeline (Facade-like)

---

## Step 3 — Flows

**Index:** tokenize doc → add terms to inverted index → update trie  

**Search:** tokenize query → lookup posting lists → intersect/union → rank → top-k  

**Suggest:** walk trie prefix

---

## Step 4 — APIs

`index(docId, text)`, `search(query, limit)`, `suggest(prefix)`

---

## Step 5 — Deepen

- Lock index during write; concurrent reads (RW lock) or copy-on-write
- Empty query, stop words (optional)

---

## Step 6 — Evolve

- BM25 ranker as **Strategy**; DIP with interfaces if multiple backends planned
