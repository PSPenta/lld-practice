# AI Suggest Reply (copilot draft) — LLD design walkthrough

> **Paper design** — no coded folder yet.  
> **Method:** [../../README.md §5](../../README.md) · **AI primer:** [../../README.md §15](../../README.md)

Use this as your **AI LLD template**. Same structure as LRU / Parking Lot interviews.

---

## Clarify (ask first!)

1. Draft only or auto-send to customer?  
2. Streaming to UI (SSE) required?  
3. Multi-tenant SaaS? (almost always yes)  
4. One LLM vendor or many?  
5. Do we have a knowledge base / past tickets for RAG?  
6. Credits / quotas per plan?  
7. Approximate QPS of “Suggest” clicks?

**Assumptions to state if they don’t specify:**
- Draft only (human sends)
- SSE streaming
- Multi-tenant
- RAG over KB + ticket thread
- Credits required

---

## Actors & use cases

- **Agent** clicks Suggest on a ticket  
- **System** builds grounded draft, streams it, stores it  
- **Admin** (optional) uploads KB docs (indexing path — mention briefly)

---

## Entities

```text
Tenant, Ticket, Message
Suggestion { id, ticket_id, status, model, text, created_by }
Chunk { id, tenant_id, doc_id, text, embedding }
CreditAccount / CreditReservation
IdempotencyRecord
```

---

## Classes / interfaces

```text
SuggestHandler          // HTTP boundary: auth + validation
SuggestService          // orchestrator (main flow)

TicketRepository        // load ticket + messages
Retriever               // vector search with tenant filter
PromptBuilder           // system rules + chunks + thread
LLMClient (interface)   // Complete / Stream
  OpenAIAdapter
  AnthropicAdapter
OutputValidator         // length / JSON / basic safety
CreditMeter             // Reserve / Commit / Release
Guardrails              // max input size, PII redact
SuggestionRepository    // save draft + status
```

---

## Happy path (say this in order)

```text
1. AuthZ: agent can access ticket (same tenant)
2. Idempotency-Key → return existing suggestion if replay
3. Guardrails on input size
4. CreditMeter.Reserve(units)
5. Load ticket thread
6. Retriever.topK(query, tenant_id)
7. PromptBuilder.build(...)
8. Save Suggestion(status=streaming)
9. LLMClient.Stream(ctx) → SSE to client
10. Validate final text; save; CreditMeter.Commit
```

---

## Failure path

```text
LLM timeout / error
  → cancel context
  → CreditMeter.Release
  → Suggestion status=failed
  → UI shows retry (same Idempotency-Key OK)
Ticket itself is unchanged (partial failure is fine)
```

---

## APIs

```http
POST /v1/tickets/{id}/ai/suggestions
Headers: Authorization, Idempotency-Key
Body: { "tone": "friendly" }
→ 201 { "suggestion_id", "stream_url" }

GET /v1/tickets/{id}/ai/suggestions/{sid}/stream
→ SSE tokens, then done|error
```

---

## Concurrency

- Two Suggest clicks: two reservations if credits allow; or limit one in-flight suggestion per ticket  
- Credit reserve must be atomic (DB row lock / conditional update)  
- Never hold DB transaction open during LLM HTTP call  

---

## Patterns used here

| Pattern | Where |
|---------|--------|
| Adapter | LLM vendors behind `LLMClient` |
| Strategy | Model pick by plan / latency |
| Repository | Ticket / Suggestion / Chunk stores |
| Middleware | Auth + credit checks |

---

## Evolve

| Change | Design move |
|--------|-------------|
| New model vendor | New Adapter only |
| 10× traffic | Cache retrieval; queue non-interactive jobs; scale workers |
| Auto-send | Stronger validator + approval flag |
| Tool calling | Agent loop: model → tool → observe → model (cap max steps) |
| Eval | Store prompt version; thumbs-up/down API; golden set offline |

---

## What to practice aloud (20–30 min)

Close the doc. Draw boxes from memory. Hit: tenant RAG filter, credits, timeout, SSE, new vendor.
