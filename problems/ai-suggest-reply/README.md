# AI Suggest Reply (copilot draft) — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> **AI primer:** [Hub §15](../../README.md#15-ai--llm-lld-for-beginners) · [docs/ai-lld/](../../docs/ai-lld/README.md)

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Draft only or auto-send to customer?
2. Streaming to UI (SSE) required?
3. Multi-tenant SaaS? (almost always yes)
4. One LLM vendor or many?
5. Do we have a knowledge base / past tickets for RAG?
6. Credits / quotas per plan?
7. Approximate QPS of “Suggest” clicks?
8. PII redaction / retention of prompts?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Support agent, SuggestService, LLM provider |
| **Use cases (v1)** | 1. Agent clicks Suggest 2. System streams grounded draft 3. Agent edits and sends |
| **In scope** | RAG retrieve, prompt build, stream, credit check |
| **Out of scope** | Auto-send to customer, fine-tuning, multi-model routing |
| **Assumptions** | Draft-only; SSE streaming; multi-tenant; RAG over KB + ticket thread; credits required; ticket works if AI down |

### Confirm understanding
> "Agent triggers suggest; we retrieve tenant KB chunks, stream a draft into the composer; human sends."

## Step 2 — Entities & classes

```text
Tenant, Ticket, Message
Suggestion { id, ticket_id, status, model, text, created_by }
Chunk { id, tenant_id, doc_id, text, embedding }
CreditAccount / CreditReservation
IdempotencyRecord

SuggestHandler          // HTTP boundary: auth + validation
SuggestService          // orchestrator (main flow)

TicketRepository        // load ticket + messages
Retriever               // vector search with tenant filter
PromptBuilder           // system rules + chunks + thread
LLMClient (interface)   // Complete / Stream
  OpenAIAdapter | AnthropicAdapter
OutputValidator         // length / JSON / basic safety
CreditMeter             // Reserve / Commit / Release
Guardrails              // max input size, PII redact
SuggestionRepository    // save draft + status
```

**Patterns:** Adapter (LLM vendors) · Strategy (model by plan) · Repository · Middleware (auth + credits)

## Step 3 — Flows

**Happy path**
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

**Edge cases / failure**
1. LLM timeout/error → cancel context → CreditMeter.Release → Suggestion failed → UI retry (same Idempotency-Key OK)  
2. Ticket unchanged on AI failure (partial failure is fine)

## Step 4 — APIs

```http
POST /v1/tickets/{id}/ai/suggestions
Headers: Authorization, Idempotency-Key
Body: { "tone": "friendly" }
→ 201 { "suggestion_id", "stream_url" }

GET /v1/tickets/{id}/ai/suggestions/{sid}/stream
→ SSE tokens, then done|error
```

## Step 5 — Deepen

- Two Suggest clicks: two reservations if credits allow, or limit one in-flight suggestion per ticket  
- Credit reserve must be atomic (DB row lock / conditional update)  
- Never hold DB transaction open during LLM HTTP call  
- Always filter retrieval by `tenant_id` (no cross-tenant RAG leakage)  
- Cancel context on client disconnect; Release credits on failure

## Step 6 — Evolve

| Change | Design move |
|--------|-------------|
| New model vendor | New Adapter only (**OCP**) |
| 10× traffic | Cache retrieval; queue non-interactive jobs; scale workers |
| Auto-send | Stronger validator + approval flag |
| Tool calling | Agent loop: model → tool → observe → model (cap max steps) |
| Eval | Store prompt version; thumbs-up/down API; golden set offline |

**Practice aloud (20–30 min):** Close the doc. Draw boxes from memory. Hit: tenant RAG filter, credits, timeout, SSE, new vendor.
