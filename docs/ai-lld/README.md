# AI / LLM LLD

← [Back to hub §15](../../README.md#15-ai--llm-lld-for-beginners)

---

## 15. AI / LLM LLD for beginners

This section is for people new to AI interviews. You are **not** asked to train models. You are asked to **design application components** that *use* LLMs safely.

### What interviewers mean by “AI LLD”

| They want | They usually do **not** want |
|-----------|------------------------------|
| Classes/interfaces around LLM calls | Deriving transformer math |
| How RAG pieces connect | Training a custom model from scratch |
| Credits, timeouts, retries, validation | Picking the “best” model brand only |
| Multi-tenant isolation of data | Building a vector DB from scratch |

Think: **same as PaymentGateway** — wrap a vendor behind an interface; orchestrate a workflow; handle failure and cost.

### Words you must understand (plain English)

| Term | Meaning |
|------|---------|
| **LLM** | A text model API (OpenAI, etc.) you call with a prompt; it returns text |
| **Prompt** | Instructions + context you send to the model |
| **Token** | Chunk of text the model bills/limits by (roughly pieces of words) |
| **Hallucination** | Confident but wrong answer — why we add retrieval + validation |
| **Embedding** | Turning text into a list of numbers so “similar meaning” is searchable |
| **Vector DB** | Store embeddings; find nearest neighbors for a query |
| **RAG** | Retrieve relevant docs → put them in the prompt → then generate |
| **Chunking** | Splitting long docs into smaller pieces before embedding |
| **Tool / function calling** | Model asks your app to run a function (search DB, create ticket) |
| **Guardrails** | Checks before/after model: size limits, PII, blocked topics, schema |
| **SSE / streaming** | Send answer token-by-token to the UI instead of one big wait |
| **Credit metering** | Limit how much AI usage a tenant/user can consume |

### Core building blocks (LEGO pieces)

Almost every AI feature LLD uses some of these:

```text
1. API / Handler          — auth, validate input
2. Orchestrator Service   — runs the steps in order
3. Guardrails             — input checks / PII redact
4. CreditMeter            — reserve / commit / release credits
5. Retriever (RAG)        — find relevant chunks for this tenant
6. PromptBuilder          — build the final prompt string
7. LLMClient (interface)  — talk to any model vendor
8. OutputValidator        — parse/validate model output (JSON/schema)
9. SuggestionStore / Repo — save drafts / history
10. StreamSession (optional) — SSE to browser
```

### Design patterns that show up in AI LLD

| Pattern | AI use |
|---------|--------|
| **Strategy** | Model routing by plan (cheap vs strong model); different retrievers |
| **Adapter** | `OpenAIClient` / `AnthropicClient` behind `LLMClient` |
| **Factory** | Create tool handlers by name |
| **Observer / Pub-Sub** | “Suggestion completed” → analytics / notify |
| **Decorator / Middleware** | Auth, tenant, credit check around the API |
| **Repository** | TicketStore, ChunkStore, SuggestionStore |

Repo references for the same ideas (non-AI, but same patterns):
- Strategy: `JavaScript/RateLimiter2/`, `Go/RateLimiter2-go/`
- Adapter-style interface: `Go/PaymentGateway-go/bank_gateway.go`
- Factory: `Go/Splitwise-go/expense.go`, `Go/ParkingLot2-go/vehicle.go`

### RAG — simplest correct story

```text
Offline (indexing):
  Document → split into Chunks → Embedding → save in Vector DB
             (always store tenant_id on each chunk)

Online (user asks / agent clicks Suggest):
  User text → Embedding → search Vector DB (filter tenant_id)
           → top-K chunks → add to Prompt → LLM → answer
```

**Why RAG?** Model alone may invent facts. Retrieved chunks ground the answer.

**Must say in interview:** retrieval is **tenant-scoped** — never mix Customer A’s docs into Customer B’s prompt.

### Sync vs async (when to choose)

| Approach | Use when |
|----------|----------|
| **Sync + SSE stream** | User is waiting for a draft reply in UI |
| **Async queue + worker** | Bulk re-embed docs, nightly summaries, non-interactive jobs |

Don’t call a slow LLM inside a DB transaction. Don’t do heavy embedding on the webhook request path — enqueue it.

### Reliability checklist for AI components

1. **Timeout** every LLM/retriever call (`context` deadline)  
2. **Retry** only transient failures; use idempotency key for billable calls  
3. **Reserve credits before** call; **commit** on success; **release** on failure  
4. **Validate output** — don’t trust raw model JSON  
5. **Partial failure** — ticket/email still works if AI is down  
6. **Observability** — `trace_id`, latency, cost/tokens, error rate, thumbs-down  

### Minimum AI LLD problems to practice

1. **LLM provider abstraction** — `LLMClient` + 2 adapters + router  
2. **RAG retriever** — chunk metadata, top-K, tenant filter  
3. **AI credit meter** — reserve/commit (like wallet) + rate limit  
4. **Suggest-reply / copilot** — full orchestration → [problems/ai-suggest-reply/README.md](../../problems/ai-suggest-reply/README.md)  
5. **Tool-calling agent (basic)** — model returns tool name → your code runs tool → continue  

### What “good enough for a noob” looks like in the interview

You can draw the boxes, name interfaces, walk the happy path, then say:

- credits under concurrency  
- LLM timeout  
- tenant isolation in RAG  
- how you’d add a new model vendor (new Adapter only)

You do **not** need to implement a vector DB.

### Five AI LLD problems — expanded sketches

Practice **[AI Suggest Reply](../../problems/ai-suggest-reply/README.md)** and **LLM abstraction** deeply; sketch the rest.

**A — LLM provider abstraction:** `LLMClient` interface with `Complete` / `Stream`; `OpenAIAdapter`, `AnthropicAdapter`; `LLMRouter` (Strategy) by plan/cost; failover on timeout.

**B — RAG retriever:** Offline: chunk → embed → store with `{ tenant_id, doc_id }`. Online: embed query → vector search **filtered by tenant** → top-K → prompt. Never cross-tenant retrieval.

**C — AI credit meter:** `Reserve` → `Commit` / `Release` (wallet-like, not just `Allow()`). DB ledger is source of truth; not cache alone.

**D — Suggest reply:** Full flow → [problems/ai-suggest-reply/README.md](../../problems/ai-suggest-reply/README.md).

**E — Tool-calling agent:** Loop: LLM → tool call → your code runs tool → append result → LLM again; cap max iterations; AuthZ per tool.

### Evaluation & AI prep drill

**Offline eval:** Golden tickets + expected facts. **Online:** thumbs, edit distance to sent reply, latency, cost. Store `prompt_version` + model on each suggestion.

**15-min drill:** Draw architecture → Suggest happy path → timeout/release credits → RAG + tenant filter → new vendor = new Adapter only.

For **RAG/code review depth** (chunking, retrieval, vector DB, reviewing repos like ragbot), see **[ai-code-review-round/README.md](../../ai-code-review-round/README.md)**.

---

