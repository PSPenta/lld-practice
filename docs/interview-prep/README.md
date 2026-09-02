# Interview Prep

← [Back to hub §17–§22](../../README.md#17-hld-topics-that-bleed-into-lld)

---

## 17. HLD topics that bleed into LLD

Some **HLD** appears in LLD as *“how would you evolve?”* — not full multi-region design.

### Evolution template

```text
v1: Sync API → service → DB → LLM
v2: Add Redis (tenant config, hot retrieval)
v3: Queue (webhooks, bulk embed KB)
v4: Horizontal scale stateless API + workers
```

### Common bleed-in topics

| Topic | One-liner |
|-------|-----------|
| Caching | L1 in-process + L2 Redis; invalidate on KB update |
| Queue | Webhook + **connector file sync**: enqueue job (SQS/Kafka), 202 fast, worker ingests; DLQ |
| Rate limit / credits | Per-tenant 429 + reserve/commit |
| DB | Index `(tenant_id, inbox_id, status)`; cursor pagination |
| Scale out | Stateless API; shared DB/Redis/queue |
| Observability | `trace_id` across services; p99 per step |
| Multi-tenant | Every query/filter includes `tenant_id` |

**3-step answer:** (1) “I’d evolve without rewriting core classes.” (2) One concrete step. (3) Trade-off. Then stop.

**Mini-scenario:** Gmail webhook storm → verify, dedupe `event_id`, enqueue, 200 fast. **Connector sync** → don’t fetch Slack/GDrive in HTTP thread; enqueue sync job, return 202, workers chunk+embed. LLM suggest 15s timeout → release credits; `trace_id` ties flow.

---

## 18. Timed mock + self-score

**Timer: 40 minutes.** Pick **[Cache Client](../../problems/cache-client/README.md)** or **[Rate Limiter](../../problems/rate-limiter/README.md)**. Speak out loud; no notes for first 30 min.

After time, score 0–2 each (target ≥ 16/20):

1. Clarified requirements  
2. Stated assumptions  
3. Named entities and classes  
4. Drew simple diagram  
5. Walked main flows (Get/Set or Allow)  
6. Mentioned mutex / thread safety  
7. Mentioned failure case (DB/Redis/LLM down)  
8. Showed evolution (interface, Redis, queue)  
9. Named 1 pattern with reason  
10. Explained trade-offs plainly  

**Opening (memorize):** “I’ll clarify requirements first, then entities, classes, flows, APIs, concurrency, failures, and how the design evolves.”

**Closing:** “I kept v1 simple and showed how it evolves under load and failures. Happy to deep-dive any part.”

---

## 19. How to practice (4-week plan)

> **Prerequisite:** Finish **[JavaScript/README.md](../../JavaScript/README.md)** (and **[Go/README.md](../../Go/README.md)** if applicable) before Week 2 — see the **prep table at the top of [README.md](../../README.md)**.

### Week 0 — Language & runtime (before LLD coding)
- **[JavaScript/README.md](../../JavaScript/README.md):** event loop, Promises, closures, §27 gotchas, §29 cheat sheet  
- **[Go/README.md](../../Go/README.md)** (if Go role): GMP, goroutines, channels, context, §18 gotchas  
- Answer §28 checklists in each guide aloud  

### Week 1 — Foundations
- SOLID + 5 patterns with your own examples  
- Draw class diagrams for LRU + Parking Lot on paper  
- Practice clarifying questions aloud  

### Week 2 — Core problems
Implement or redesign:
- LRU, Rate Limiter, Parking Lot, Pub-Sub  
- For each: write APIs + failure notes  

### Week 3 — More problems + AI + concurrency
- Splitwise, **PollingSystem2** (service + repos), Notification system, Cache client  
- **AI:** LLM interface + RAG sketch + [Suggest Reply](../../problems/ai-suggest-reply/README.md) once on paper  
- Add mutex/idempotency discussion every time  
- One machine-coding simulation (90 min timer)  

### Week 4 — Interview simulation
- 4 timed discussion LLDs (45 min each) — include **one AI** problem  
- Record yourself; check structure adherence  
- Prepare 2 stories linking designs to past work  
- Optional: one **vibe coding** mock ([vibe-coding-round/README.md](../../vibe-coding-round/README.md))  
- Paper-design 2 **gap** problems from [lld-gaps/README.md](../../lld-gaps/README.md)

### Practice method (every problem)
1. 5 min questions  
2. 10 min entities  
3. 20 min classes + flow  
4. 10 min API + deepen  
5. Compare with a reference solution (this repo / articles)  
6. Note 3 improvements  

---

## 20. Interview day checklist

**Before**
- [ ] **Language prep done:** [JavaScript/README.md](../../JavaScript/README.md) · [Go/README.md](../../Go/README.md) if Go role — see **prep table on [README.md](../../README.md)**  
- [ ] Paper / shared editor ready  
- [ ] Know opening + closing lines ([§18](../interview-prep/README.md#18-timed-mock--self-score))  
- [ ] Know round type: **LLD design** ([hub](../../README.md)) vs **code review** ([ai-code-review-round](../../ai-code-review-round/README.md))  
- [ ] Skim company product (helpdesk/SaaS: shared inbox, AI copilot, multi-tenant) — honest if you didn’t use the product  
- [ ] Sleep; don’t cram 10 new patterns  

**During (LLD design)**
- [ ] Ask questions first  
- [ ] State assumptions  
- [ ] Drive use-case by use-case  
- [ ] Prefer composition + interfaces  
- [ ] Mention concurrency & failure at least once  
- [ ] Show how design evolves ([§17](../interview-prep/README.md#17-hld-topics-that-bleed-into-lld))  

**During (code review)**
- [ ] Trace main happy path, then untested routes  
- [ ] Prioritize: security → RAG correctness → production ops  
- [ ] Cite file names; group findings P0/P1/P2  

**Avoid**
- Jumping to Kafka/microservices immediately  
- Pattern-dropping without need  
- Coding before the model is clear (discussion rounds)  
- Ignoring edge cases (null, empty, capacity full)  
- Fake deep product usage  

### Company research template (helpdesk / AI SaaS)

If they ask *“What do you know about us?”*:

1. **One sentence:** AI-powered customer support — shared inbox + workflows + copilot/agents.  
2. **Problem:** Teams sharing support@ without chaos; knowledge for agents; draft replies grounded in docs.  
3. **Your fit:** Multi-tenant APIs, RAG, streaming, credits, reliability — same design space.  
4. **Honest product use:** Website/docs reviewed; trial if offered.

---

## 22. Cheat sheet

### Structure
```text
Clarify → Entities → Classes/Interfaces → Flows → APIs →
Concurrency/Failure → Extend/Trade-offs
```

### SOLID + other principles one-liners
- **S:** one reason to change  
- **O:** extend without editing core  
- **L:** subtypes don’t break contracts  
- **I:** small interfaces  
- **D:** depend on abstractions  
- **Composition:** has-a + interfaces over deep is-a inheritance  
- **DRY:** don’t duplicate knowledge  
- **KISS / YAGNI / PoLK:** simplest design; add patterns when needed — repo map **§7A**  
- **KISS:** simplest workable design  
- **YAGNI:** don’t build unused features  
- **PoLK:** talk only to direct collaborators  

### Patterns one-liners
**Creational:** Singleton (one instance) · Factory (create by type) · Builder (step-by-step) · Abstract Factory (product families) · Prototype (clone)  
**Behavioural:** Observer (notify subscribers) · Strategy (swap algorithm) · Iterator (traverse) · Command (request as object) · State · Template Method  
**Structural:** Adapter (bridge interfaces) · Proxy (stand-in / gateway) · Decorator (wrap behavior) · Facade (simple front door)  
**Repo demos:** Full map → **[../repo-map/README.md](../repo-map/README.md)**. Quick: Strategy → RateLimiter2 · Factory → Splitwise / ParkingLot2 · Observer → Pub-Sub · DIP → PaymentGateway-go · Pipeline → SearchEngine · KISS → JavaScript/Queue/LRU

### AI LLD one-liners
- LLM behind `LLMClient` interface (Adapter)  
- RAG = retrieve tenant-scoped chunks → prompt → generate  
- Credits = Reserve before call, Commit/Release after  
- Never trust model output — validate  
- Timeout LLM calls; ticket must work if AI is down  
- Stream (SSE) for UX; queue for bulk/offline jobs  
- **Code review:** secrets, SQL injection, grounded RAG → [ai-code-review-round](../../ai-code-review-round/README.md)

### HLD bleed-in (when asked “scale?”)
- v1 simple → Redis → queue → stateless scale → indexes + `trace_id`  
- Multi-tenant: filter every query by `tenant_id`

### Must-practice designs
- **[Cache Client](../../problems/cache-client/README.md)** · [Rate Limiter](../../problems/rate-limiter/README.md) · [AI Suggest](../../problems/ai-suggest-reply/README.md)

### Production questions (ask yourself)
- What if traffic ×10?  
- What if data ×10?  
- What if dependency fails?  
- What if two requests hit same row?  
- What if we add a new variant?  
- How do we observe this in prod?  

### Closing line
> “I kept v1 simple, made extension points where variation is real, and called out concurrency and failure modes. Happy to go deeper on any part.”

---

## Final note

→ **[Hub final note](../../README.md#final-note)** — start path and code-review pointer live on the root README only.

