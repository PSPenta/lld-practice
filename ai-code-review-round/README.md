# AI Code Review Round — Preparation Guide

> **Based on:** AI SDE / full-stack interviews where you clone a repo (e.g. [pranavgupta96/ragbot](https://github.com/pranavgupta96/ragbot)), read the code **manually**, and find what is wrong or not production-ready.  
> **For design / LLD discussion rounds** (classes, APIs, Cache Client, AI Suggest Reply on whiteboard), see **[../README.md](../README.md)** (§4 round types, [§16 design docs](../README.md#16-worked-examples--design-docs), §17–§22).  
> **For AI-assisted build rounds** (Cursor/Copilot allowed), see **[../vibe-coding-round/README.md](../vibe-coding-round/README.md)**.  
> **Audience:** Beginners — every concept explained in plain language first, then interview application.

---

## Table of contents

1. [What this round is (not LLD on a whiteboard)](#1-what-this-round-is-not-lld-on-a-whiteboard)
2. [What interviewers expect](#2-what-interviewers-expect)
3. [How to review a repo — step by step](#3-how-to-review-a-repo--step-by-step)
4. [The ragbot example — walkthrough](#4-the-ragbot-example--walkthrough)
5. [What to highlight (priority order)](#5-what-to-highlight-priority-order)
5A. [Sync vs async — connector file sync (SQS/Kafka)](#sync-vs-async--connector-file-sync-interviewer-highlight)
6. [RAG from zero — layman's guide](#6-rag-from-zero--laymans-guide)
7. [Semantic search — explained simply](#7-semantic-search--explained-simply)
8. [Vector DB — explained simply](#8-vector-db--explained-simply)
8A. [LangChain & LangGraph — vs hand-rolled RAG](#8a-langchain--langgraph--vs-hand-rolled-rag)
9. [Chunking — overlap, strategies, trade-offs](#9-chunking--overlap-strategies-trade-offs)
10. [Retrieval strategies in RAG](#10-retrieval-strategies-in-rag)
11. [Review checklist (print this)](#11-review-checklist-print-this)
12. [How to speak your findings (template)](#12-how-to-speak-your-findings-template)
13. [Practice plan](#13-practice-plan)

---

## 1. What this round is (not LLD on a whiteboard)

### Two different “AI rounds”

| Round type | What you do | Example |
|------------|-------------|---------|
| **Design LLD** | Draw classes, APIs, flows on board | “Design Cache Client” / “Design AI Suggest Reply” |
| **Vibe coding / AI build** | AI tools allowed; you design, verify, own output | “Build rate limiter with Cursor in 60 min” |
| **Code review** (this doc) | Clone repo, read code, list bugs & production gaps | “Review this RAG helpdesk service” |

> **Other round types:** whiteboard LLD → [../README.md](../README.md) · AI-assisted build → [../vibe-coding-round/README.md](../vibe-coding-round/README.md) · LLD breadth gaps → [../lld-gaps/README.md](../lld-gaps/README.md)

Companies like Hiver (AI SDE roles) increasingly use **code review** because:

- AI tools generate code fast — they want to see if **you** can judge quality.
- Real job = read PRs, spot security issues, say what’s not production-grade.
- Tests whether you understand **RAG, DB, APIs** in real code — not only theory.

### What “manual review” means

- Open files yourself; follow imports; trace one request end-to-end.
- Don’t only run the app — **read** `connector_routes.py` even if tests pass.
- Compare “clean” modules vs “messy” ones — bugs are often in the second path.

---

## 2. What interviewers expect

### They score you on

| Skill | Good signal |
|-------|-------------|
| **Prioritization** | Security first, then correctness, then performance |
| **Specificity** | File name + line + why it’s wrong + how to fix |
| **Production lens** | Pooling, secrets, timeouts, idempotency, observability |
| **RAG literacy** | Chunk vs retrieve vs generate; grounded vs hallucination |
| **Architecture eye** | Same app doing two different things inconsistently |
| **Communication** | Short summary + grouped findings, not 50 random nitpicks |

### They do NOT expect

- Fixing every file in the interview
- Memorizing LangChain API
- Perfect knowledge of every vector DB product name

### They DO expect

- “Hardcoded API keys in source” → **P0**
- “SQL built with f-strings from user input” → **P0**
- “This prompt tells the model to guess when context is thin” → **P0 for RAG product**
- “No connection pool on Postgres” → **valid P1/P2** depending on depth
- “Connector sync / bulk ingest runs in the HTTP handler” → **valid P1/P2** — should enqueue (SQS/Kafka), return **202 + job_id**, workers do fetch → chunk → embed

---

## 3. How to review a repo — step by step

Use **45–60 minutes** like this:

| Minutes | Action |
|---------|--------|
| 0–5 | Read README — what does it claim? (architecture, secrets, offline mode) |
| 5–10 | Trace **happy path**: ingest doc → ask question → answer |
| 10–20 | Core modules: `services/`, `retrieval/`, `llm/`, `db.py`, `config.py` |
| 20–35 | **All routes** — especially files not covered by tests |
| 35–45 | Grep: `password`, `api_key`, `sk-`, `f"SELECT`, `except: pass`, `requests.` |
| 45–50 | Note positives (clean patterns) — interviewers like balance |
| 50–60 | Prepare 2-min spoken summary |

### One request to trace (ragbot)

```text
POST /api/v1/ask
  → routes.py
  → AnswerService.ask()
  → Retriever.top_chunks()   ← retrieval
  → LLMClient.complete()     ← generation
  → AskResponse + sources
```

Then trace **alternate paths** (same file — where planted bugs live):

```text
POST /api/v1/connectors/{id}/ask
  → connector_routes.py
  → loads ALL chunks into prompt (not Retriever.top_k)
  → raw requests.post to Anthropic     ← bypasses LLMClient adapter
  → build_connector_prompt()           ← “always confident” / not grounded

POST /api/v1/connectors/{id}/sync
  → connector_routes.py
  → delete existing docs (sync)
  → fetch Slack / Notion / GDrive (sync, no queue)
  → for each item: ingest → chunk → embed → DB (sync)
  → return { documents_synced: N }     ← client waited for entire pipeline
```

**Why trace `/sync` separately:** It is the **interviewer highlight** — heavy I/O and embedding run **inside the HTTP request** instead of on a worker behind SQS/Kafka.

---

## 4. The ragbot example — walkthrough

Repo: [github.com/pranavgupta96/ragbot](https://github.com/pranavgupta96/ragbot) — “Acme Helpdesk RAG service.”

### Claimed architecture (README)

```text
api/           → thin HTTP
services/      → orchestration
repositories/  → DB only
ingestion/     → loaders + chunker
retrieval/     → embeddings + top-k
llm/           → adapter (Anthropic + fake)
```

### What is **good** (say this — shows fair review)

| Area | Why it’s good |
|------|----------------|
| `AnswerService` | Refuses to call LLM if retrieval empty — **no blind guessing** |
| `prompts.py` | “Use ONLY context”; treats passages as data not instructions — **anti prompt-injection** |
| `llm/anthropic_client.py` | Timeout, retries, adapter pattern |
| `ingestion/base.py` | Loader registry — Open/Closed for new content types |
| `chunker.py` | Keeps final partial chunk — avoids silent data loss |
| `retriever.py` (main) | Heap top-k + `yield_per` — doesn’t load all chunk bodies for scoring |
| Tests | Ask + ingest + pagination (main path) |

### What is **bad** (the planted / real issues)

**File to focus on:** `app/api/connector_routes.py`

| Issue | Layman explanation |
|-------|-------------------|
| Hardcoded Slack/Notion/GDrive/Anthropic keys | Like writing your house key on the front door — anyone with repo access owns your accounts |
| SQL `f"… WHERE title LIKE '%{q}%'"` | Attacker can inject SQL via search query |
| Prompt: “use your own knowledge, always confident answer” | Opposite of helpdesk RAG — **invites hallucination** |
| `ask_connector` calls Anthropic via raw `requests.post` | **Bypasses `LLMClient` adapter** — no shared timeout/retry/offline mode; duplicate HTTP path |
| Loads **all** chunks into one prompt | Like pasting 500 pages into ChatGPT — slow, expensive, misses focus |
| `while status == 429: retry` with no sleep | Computer spins forever when Google rate-limits you |
| `except Exception: pass` | Errors disappear — ops team sees “0 synced” with no clue why |
| `list_connectors` returns `access_token` | Leaking secrets in API response |
| `/search` re-embeds text, ignores stored vectors | Wasted work; inconsistent with main retriever |
| `/search` returns random chunks if no match | User thinks bad results are “relevant” |
| **`POST /connectors/{id}/sync` runs sync work synchronously** | HTTP request waits for Slack/Notion/GDrive fetch + ingest + embed — slow, timeouts, no retries; should use **SQS/Kafka + workers** (see below) |

**Other files:**

| Issue | File |
|-------|------|
| No DB pool / `pool_pre_ping` | `app/db.py` |
| `create_all()` on startup, no migrations | `app/main.py` |
| `content_hash` saved but never used for dedup | `ingestion_service.py` |
| Demo embedder (bag-of-words) | `retrieval/embedder.py` — OK if documented |

---

## 5. What to highlight (priority order)

Use this order when speaking:

### P0 — Security & data safety

1. Secrets in source control  
2. SQL injection  
3. Exposing tokens in API responses  
4. Prompt that bypasses grounding / encourages guessing  
5. No auth on ingest or ask (multi-tenant risk)

### P1 — Correctness & RAG quality

1. Retrieval inconsistent between endpoints  
2. **Connector ask bypasses `LLMClient`** — duplicate raw HTTP, no shared retries/timeouts  
3. No dedup on ingest (`content_hash` unused)  
4. Chunking loses structure (headers stripped, fixed word count)  
5. Weak / non-semantic embeddings in production path without guardrails  
6. Sync wipes all docs — no incremental update  
7. **Connector sync is synchronous** — file/API fetch + chunk + embed in the HTTP handler

### P2 — Production operations

1. DB connection pooling, pre-ping, pool size  
2. Migrations (Alembic) vs `create_all`  
3. Timeouts on **all** outbound HTTP (`requests` without timeout in connectors)  
4. Structured logging, trace_id, metrics  
5. Sync LLM in request thread — need JavaScript/Queue/workers at scale  
6. **Connector/document sync should be async** — SQS, Kafka, or similar; not blocking `POST /sync`  
7. Idempotency on POST ingest  

### P3 — Code quality

1. Raw `dict` bodies instead of Pydantic on some routes  
2. Duplicate logic (two retrieval implementations)  
3. Missing tests for connector/search paths  
4. Naive datetime vs UTC elsewhere  

### Sync vs async — connector file sync (interviewer highlight)

In `connector_routes.py`, `POST /connectors/{id}/sync` does **everything inside the HTTP request**:

```text
Current (bad for production):
  Client calls POST /sync
    → delete old docs (sync)
    → call Slack / Notion / Google Drive API (sync)
    → for each file/message: ingest → chunk → embed → DB write (sync)
    → return { documents_synced: N }   ← client waited for ALL of this
```

**Layman:** The user clicks “Sync” and the browser/API **holds the connection open** while the server downloads files from Slack, chunks them, and writes to the DB. If there are 500 files or Google is slow, the request **times out** or the server thread is **blocked** — no other syncs can run well.

**What production should do:**

```text
Better:
  Client calls POST /sync
    → validate connector, set sync_in_progress
    → publish job to queue (SQS / Kafka / Redis stream)
    → return 202 Accepted { job_id } immediately

  Worker(s) (separate process):
    → consume job
    → fetch from Slack/Notion/GDrive
    → ingest each doc (chunk, embed, store)
    → update connector.last_synced_at, clear sync_in_progress
    → on failure: retry with backoff, DLQ for poison messages
```

| | Synchronous (ragbot) | Async queue (SQS / Kafka) |
|---|---------------------|----------------------------|
| **API response** | Waits until all files synced | Returns fast with job id |
| **Timeouts** | Gateway/client may timeout | Workers can run minutes |
| **Retries** | Whole request fails | Per-message retry |
| **Scale** | One thread per sync request | Many workers in parallel |
| **Backpressure** | Overloads API server | Queue buffers spikes |

**When to say SQS vs Kafka (interview one-liner):**
- **SQS** — simple job queue, “sync this connector”, few consumers, AWS-native, at-least-once delivery.
- **Kafka** — high throughput, event log, many consumers, replay indexing pipeline, ordering per partition (e.g. per `tenant_id`).

**Same pattern applies to:**
- Gmail/webhook ingest (ack fast, process async)
- Bulk KB re-embedding after doc upload
- Never run **LLM calls** or **heavy embedding** on the webhook/sync HTTP path

**Interview line (memorize):**

> “The connector sync endpoint does file fetch and ingestion synchronously in the request thread. For production I’d enqueue a sync job to SQS or Kafka, return 202 immediately, and let workers handle fetch-chunk-embed with retries and a DLQ — same pattern as webhook ingest.”

---

## 6. RAG from zero — layman's guide

### The problem RAG solves

**LLM alone** = very smart intern who read the internet once but **doesn’t have your company’s latest docs** in front of them. They might **make up** refund policies — that’s **hallucination**.

**RAG** = Retrieval-Augmented Generation:

```text
Before answering, LOOK UP your company docs, THEN answer using what you found.
```

Like an open-book exam instead of memorization.

### RAG in one sentence

> **Find the right paragraphs from your knowledge base, put them in the prompt, then let the LLM write the answer using only (or mainly) that text.**

### The full pipeline (two phases)

#### Phase A — OFFLINE (when docs are added or updated)

Think: **building a searchable library index.**

```text
1. INGEST     — get documents (PDF, markdown, tickets, Slack)
2. PARSE      — extract clean text
3. CHUNK      — split long docs into small pieces
4. EMBED      — turn each chunk into numbers (a "vector")
5. STORE      — save chunk text + vector + metadata (tenant_id, doc_id)
```

You do this once per doc (or on update). Not on every user question.

#### Phase B — ONLINE (when user asks a question)

Think: ** librarian finds pages, then writer summarizes.**

```text
1. USER QUESTION  — "How long do refunds take?"
2. EMBED QUESTION — same kind of numbers as chunks
3. RETRIEVE       — find chunks most similar to the question (top-K)
4. AUGMENT        — paste chunks into the prompt as "context"
5. GENERATE       — LLM writes answer citing that context
6. VALIDATE       — check length, PII, format; return sources
```

**RAG is NOT just chunking.** Chunking is step A3. Retrieve + augment + generate are the rest.

### Simple analogy

| Step | Analogy |
|------|---------|
| Chunk | Cut a textbook into index cards |
| Embed | Give each card a fingerprint of its meaning |
| Vector DB | Filing cabinet sorted by fingerprint |
| Retrieve | Find cards closest to the question’s fingerprint |
| LLM | Teacher reads those cards and explains in plain English |

### What ragbot does well vs badly

| Good (`/api/v1/ask`) | Bad (`/connectors/.../ask`) |
|----------------------|----------------------------|
| Retrieve top-K chunks | Dumps **all** chunks |
| Grounded system prompt | “Use your own knowledge” |
| Uses LLM adapter | Raw HTTP + hardcoded key |
| Returns sources + scores | No proper citations |

---

## 7. Semantic search — explained simply

### Keyword search (old way)

Search for exact words: `"refund policy"`

- Misses: “money back within five days” (same meaning, different words)
- Matches junk: “no refund on refund fees” (word appears but wrong context)

### Semantic search (meaning-based)

Search by **meaning**, not exact letters.

- Question: “When do I get my money back?”
- Can match chunk: “Refunds processed within five business days”

**How?** Both sentences are converted to **vectors** (lists of numbers). Similar meaning → vectors are **close** in math space.

### Embedding (the magic step)

**Embedding** = a model that turns text into a fixed-size list of numbers (e.g. 768 or 1536 floats).

```text
"refund in 5 days"  → [0.02, -0.11, 0.45, ...]
"money back fast"   → [0.03, -0.09, 0.43, ...]   ← similar numbers!
"API rate limits"   → [-0.5, 0.8, 0.1, ...]      ← far away
```

You don’t hand-write these numbers — an **embedding model** (OpenAI, Cohere, sentence-transformers, etc.) produces them.

### Similarity score

Usually **cosine similarity** (angle between two vectors):

- **1.0** = very similar meaning  
- **0.0** = unrelated  
- ragbot demo uses bag-of-words hash — **keyword-ish**, not true semantic (fine for demo, weak for paraphrases)

### Semantic search in RAG

```text
Question → embed → compare to all chunk vectors → pick top-K highest scores
```

That’s **semantic retrieval** — the “R” in RAG.

---

## 8. Vector DB — explained simply

### Do you need a special database?

**Small scale:** You can store vectors in PostgreSQL (JSON column) and scan all rows — ragbot does this with brute-force + heap (OK for thousands of chunks).

**Large scale:** You need a **vector database** (or pgvector index) built for fast “nearest neighbor” search among **millions** of vectors.

### What a vector DB stores

Each row is roughly:

```text
chunk_id: 42
text: "Refunds are processed within five business days..."
vector: [0.02, -0.11, ...]
metadata: { tenant_id: "acme", doc_id: 7, source: "refund-policy.md" }
```

### What it does on query

```text
Input: query vector
Output: top-K chunk IDs with highest similarity
```

Algorithms: HNSW, IVF, etc. — interview depth: “indexed approximate nearest neighbor, sub-linear search.”

### Popular options (names only)

| Tool | One line |
|------|----------|
| **pgvector** | Postgres extension — good if you already use Postgres |
| **Qdrant / Pinecone / Weaviate** | Dedicated vector stores |
| **OpenSearch / Elasticsearch** | Can do kNN on vectors too |

ragbot README: *“Retriever is the seam where pgvector/Qdrant would slot in.”*

### Multi-tenant rule (critical for Hiver-style products)

**Always filter by `tenant_id` before or during vector search.**

Never search all customers’ docs and hope the top result is right — that’s a **data leak**.

---

## 8A. LangChain & LangGraph — vs hand-rolled RAG

You may hear these names in an AI SDE interview. You **do not** need to memorize their APIs — you need to know **what problem they solve** and how that maps to code like ragbot.

### LangChain (library)

**Layman:** A **toolkit** for building LLM apps — pre-built pieces for loading documents, splitting text, embedding, retrieving, and calling models, wired together in a “chain.”

| Concept in LangChain | ragbot equivalent |
|----------------------|-------------------|
| Document loaders | `ingestion/loaders/` (Slack, Notion, file) |
| Text splitter / chunker | `ingestion/chunker.py` |
| Embeddings | `retrieval/embedder.py` |
| Vector store | Postgres + `Chunk.embedding` (pgvector would slot in) |
| Retriever | `retrieval/retriever.py` → `top_chunks()` |
| LLM call | `llm/anthropic_client.py` via `LLMClient` interface |
| Prompt template | `prompts.py` |

**Interview line:** “This repo is a **hand-rolled LangChain-shaped pipeline** — same stages, explicit modules, no framework magic. That’s fine for a small service; LangChain helps when you want standard loaders/retrievers and faster iteration.”

**What to criticize in ragbot:** The **main path** follows the adapter pattern; **`connector_routes.py` breaks it** — raw `requests` to Anthropic instead of `LLMClient`, custom prompt instead of `prompts.py`, no `Retriever`.

### LangGraph (library)

**Layman:** Builds **multi-step AI workflows as a graph** — nodes (steps), edges (flow), optional **cycles** (agent loops: retrieve → think → tool call → retrieve again).

```text
Simple RAG (ragbot /ask):     ingest → retrieve → generate   (one straight line)

LangGraph-style agent:        retrieve → LLM → call tool → retrieve again → …
                              (state machine; can loop until done)
```

ragbot does **not** use LangGraph — `/ask` is a **linear pipeline** (retrieve once, generate once). That is appropriate for a helpdesk Q&A bot. LangGraph shines when you need **tool use, routing, or multi-hop reasoning** (e.g. “search KB → if empty search web → summarize → draft reply”).

**Interview line:** “LangGraph is for **stateful agent flows** with branching and loops. This helpdesk RAG is a single retrieve-then-answer path — a service class is enough; I’d reach for LangGraph if we added tool-calling or multi-step escalation workflows.”

### When frameworks help vs hurt

| Hand-rolled (ragbot core) | Framework (LangChain / LangGraph) |
|---------------------------|-----------------------------------|
| Clear files, easy to review in interviews | Faster to prototype |
| You own every line | Abstraction can hide bugs |
| Good when team wants minimal deps | Good when pipeline grows (10+ steps, many sources) |

**Red flag either way:** Two paths doing the same thing differently — ragbot’s clean `AnswerService` **and** messy `connector_routes` is exactly that anti-pattern.

---

## 9. Chunking — overlap, strategies, trade-offs

### Why chunk at all?

- Embedding models have **token limits** (can’t embed 200-page PDF at once).
- Retrieval works best on **focused** paragraphs — not whole books.
- Smaller chunks = more precise match; larger chunks = more context per hit.

### What is **overlap**?

**Overlap** = neighboring chunks **share some text** at the boundaries.

```text
Document words:  [1 … 200 | 201 … 400 | 401 … 500]

Without overlap:
  Chunk 1: words 1–200
  Chunk 2: words 201–400
  → Sentence split at 200/201 might break meaning

With overlap 40 (ragbot default: 200 size, 40 overlap):
  Chunk 1: words 1–200
  Chunk 2: words 161–360    ← words 161–200 appear in BOTH chunks
  Chunk 3: words 321–500
```

**Layman:** Like reading a book with **reread paragraphs** at each chapter boundary so you don’t lose a sentence cut in half.

**Why overlap helps retrieval:** The answer might sit on a boundary; overlap puts that sentence fully inside at least one chunk.

**Cost:** More chunks → more storage, more embed API calls, more rows to search.

ragbot: `chunk_size=200` **words**, `overlap=40` words, step = 160 words per slide.

### Chunking strategies (common in production RAG)

| Strategy | How it works | Pros | Cons |
|----------|--------------|------|------|
| **Fixed size (characters/words/tokens)** | Split every N tokens with overlap | Simple, predictable | Cuts mid-sentence, ignores structure |
| **Token-based** | Use tokenizer (e.g. 512 tokens) | Matches model limits | Needs tokenizer lib |
| **Recursive character** | Split on `\n\n`, then `\n`, then space | Respects paragraphs | Common in LangChain |
| **Semantic / paragraph** | Split on headings, blank lines, sections | Keeps ideas whole | Harder for messy HTML/PDF |
| **Document-structure** | Markdown headers, HTML tags, PDF pages | Good for docs | Parsers needed |
| **Parent–child** | Small chunks for search, link to big “parent” for LLM context | Best of both | More complex storage |
| **Sentence-based** | One or more sentences per chunk | Clean units | Many tiny chunks |
| **No chunk (whole doc)** | Embed entire short FAQ | OK for 1-page docs | Fails on long docs |

### What ragbot uses

- **Fixed word window + overlap** (`chunker.py`) — fine for demo.
- Markdown loader strips `#` headers before chunking — **structure lost** before split.
- Good: keeps **final partial chunk** (comment in code explains this bug if dropped).

### Interview line on chunking

> “I’d use token-based chunking with 10–20% overlap, prefer splitting on paragraph boundaries first, and store metadata (section title, doc_id, tenant_id). For long policies, parent-child chunks so search is precise but the LLM gets wider context.”

---

## 10. Retrieval strategies in RAG

**Retrieval** = given a question, **which chunks** do we fetch?

### 1. Dense retrieval (vector / semantic)

- Embed question + compare to chunk vectors (cosine).
- **Default for RAG today.**
- ragbot main path: this (with demo embedder).

### 2. Sparse retrieval (keyword)

- BM25, TF-IDF, inverted index — like Google 1990s.
- Great for exact product codes, error codes, IDs.
- **SearchEngine** in your lld-practice repo is this style (inverted index).

### 3. Hybrid retrieval

- Run **both** dense + sparse → merge scores (RRF = reciprocal rank fusion).
- Common in production: semantic + exact keyword match.

### 4. Top-K vs Top-P vs threshold

| Approach | Meaning |
|----------|---------|
| **Top-K** | Always take best K chunks (e.g. K=4) |
| **Score threshold** | Only chunks with similarity ≥ 0.15 (ragbot `min_score`) |
| **Top-K + threshold** | Best K, then drop any below threshold — ragbot does this; may return empty → refusal (good) |

### 5. Reranking

- First pass: cheap retrieval (top 20).
- Second pass: **cross-encoder** reranker scores (question, chunk) pairs precisely.
- Keep top 4 for LLM — better quality, extra latency/cost.

### 6. Metadata filtering (pre-filter)

- Before vector search: `WHERE tenant_id = X AND doc_type = 'policy'`
- **Mandatory** in SaaS — not optional.

### 7. Query transformation

- **Multi-query:** LLM rewrites question 3 ways → retrieve for each → merge.
- **HyDE:** LLM generates hypothetical answer → embed that → retrieve (advanced).

### 8. Full-context dump (anti-pattern)

- Put **all** documents in prompt — ragbot connector path does this.
- Breaks at scale; dilutes attention; expensive — **flag this in review.**

### 9. Brute force vs indexed

| | ragbot main | Production |
|---|-------------|------------|
| Scan all embeddings | Yes, with heap | HNSW / pgvector index |
| Complexity | O(n) per query | ~O(log n) |

### ragbot `/search` endpoint — retrieval anti-patterns

1. Re-embeds `chunk.text` instead of using stored `embedding` column.  
2. Loads **every** chunk into memory.  
3. Hard threshold `0.9` — arbitrary.  
4. If nothing matches, returns **random** chunks with score 0 — misleading.

**Interview:** “Main `Retriever` is reasonable for a demo; `/search` contradicts the storage model and UX.”

---

## 11. Review checklist (print this)

### Security

- [ ] API keys / tokens in code or git  
- [ ] SQL injection (f-strings in queries)  
- [ ] Secrets in API responses  
- [ ] No authentication / authorization  
- [ ] Prompt injection defenses on ingested content  

### RAG correctness

- [ ] Grounded answers (context-only) vs “use your knowledge”  
- [ ] Top-K retrieval vs dump-all-context  
- [ ] Uses stored embeddings at query time  
- [ ] Citations / sources returned  
- [ ] Refusal when no relevant chunks  
- [ ] Tenant isolation in retrieval  

### Chunking & indexing

- [ ] Sensible chunk size (tokens vs words)  
- [ ] Overlap at boundaries  
- [ ] Final partial chunk not dropped  
- [ ] Dedup / content hash on re-ingest  
- [ ] Async indexing for large docs (queue)  

### DB & infra

- [ ] Connection pool, pre-ping, recycle  
- [ ] Migrations vs create_all  
- [ ] Transactions on multi-row writes  
- [ ] N+1 queries (list loads all chunks?)  
- [ ] **Heavy sync/ingest async** — not in HTTP handler (SQS/Kafka workers)  

### HTTP / reliability

- [ ] Timeouts on all external calls  
- [ ] Retries with backoff (not infinite 429 loop)  
- [ ] Errors logged, not `except: pass`  
- [ ] Idempotency on writes  

### Architecture

- [ ] Single LLM abstraction (no duplicate raw HTTP)  
- [ ] Thin routes, logic in services  
- [ ] Tests cover risky paths (not only happy path)  

### Grep commands (quick hunt)

```bash
rg -i "api_key|secret|password|sk-|xoxb-|token\s*=" .
rg "f\"SELECT|f'SELECT|\+.*SELECT" .
rg "except.*pass" .
rg "requests\.(get|post)" .    # check for timeout=
```

---

## 12. How to speak your findings (template)

### Opening (30 sec)

> “I traced the main flow: ingest → chunk → embed → store, then ask → retrieve top-K → LLM with a grounded prompt. The core in `services/` and `llm/` is structured well. Most serious issues are in `connector_routes.py`, which diverges from that design.”

### Body (group by severity)

**P0 — Security:** …  
**P1 — RAG / correctness:** connector ask **bypasses `LLMClient`**, loads all chunks instead of `Retriever.top_k`, hallucination-friendly prompt …  
**P2 — Production:** connector **sync should be async** (SQS/Kafka workers, 202 + job id, not blocking HTTP) …  

### Positive close (10 sec)

> “I’d keep the AnswerService/Retriever/LLM adapter pattern and rewrite connectors to use the same grounding, config, and retrieval stack.”

### If you missed something

> “I’d want a second pass on auth and deployment config, and load test retrieval at chunk counts we expect in prod.”

---

## 13. Practice plan

| Day | Activity |
|-----|----------|
| 1 | Read §6–10 (RAG, semantic, vector DB, chunking, retrieval) — draw pipeline on paper |
| 2 | Clone ragbot; trace `/ask`, `/connectors/.../ask`, and **`/connectors/.../sync`**; write 10 findings with file names |
| 3 | Review another small RAG repo or your SuperStocks RAG code with §11 checklist |
| 4 | Mock: 40 min silent review + 5 min present findings to a friend or recorder |

### Concepts to know cold (one-liners)

| Term | One-liner |
|------|-----------|
| **RAG** | Retrieve docs → add to prompt → generate answer |
| **Chunk** | Small piece of a long document for search |
| **Overlap** | Shared text between adjacent chunks so boundaries aren’t lost |
| **Embedding** | Text → vector of numbers capturing meaning |
| **Semantic search** | Find by meaning similarity, not exact keywords |
| **Vector DB** | Database optimized for nearest-neighbor vector search |
| **Top-K** | Take the K best matching chunks |
| **Hallucination** | Model inventing facts — RAG + grounding reduces this |
| **Grounding** | Force answer to use retrieved context only |
| **LangChain** | Toolkit library for loaders, chunkers, retrievers, LLM chains |
| **LangGraph** | Graph/state-machine library for multi-step agent workflows with loops |

---

## Quick reference — ragbot file map

| File | Role | Review focus |
|------|------|--------------|
| `app/api/routes.py` | Main API | Thin? validation? |
| `app/api/connector_routes.py` | Connectors, search | **Primary bug farm** |
| `app/services/answer_service.py` | Q&A orchestration | Grounding, refusal |
| `app/services/ingestion_service.py` | Ingest pipeline | Dedup hash? |
| `app/retrieval/retriever.py` | Top-K search | Scale, min_score |
| `app/retrieval/embedder.py` | Vectors | Production model? |
| `app/ingestion/chunker.py` | Split text | Overlap, strategy |
| `app/llm/` | LLM adapter | Timeout, retries |
| `app/db.py` | Engine | Pooling, pre_ping |
| `app/prompts.py` | Safety rules | Injection, grounding |
| `app/config.py` | Env settings | vs hardcoded keys elsewhere |

---

## Final note

This round tests **engineering judgment**, not memorization. The ragbot repo is deliberately **split**: a teachable core and a flawed extension. Your job is to **notice the split**, explain **why it matters for a helpdesk AI product**, and prioritize fixes a team would ship first.

Good luck on the next round.
