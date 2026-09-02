# LLD Method

← [Back to hub §1–§6](../../README.md#1-what-is-lld)

---

## 1. What is LLD?

**Low-Level Design** is designing software at the **class / module / API** level.

You answer questions like:

- What are the main **entities**?
- What are the **classes/interfaces** and their **methods**?
- How do objects **collaborate** for each use case?
- What does the **API** look like?
- How do we keep the design **clean, testable, and extensible**?
- What happens under **concurrency**, **failures**, and **growth**?

You are **not** primarily designing AWS topology (that is HLD).  
You **are** designing the code structure someone would implement in a sprint.

### Machine coding vs discussion LLD

| Style | What happens |
|-------|----------------|
| **Discussion LLD** | Whiteboard / shared doc: classes, flows, APIs, trade-offs (60–75 min) |
| **Machine coding** | Write working code in 90–120 min for a small system |

Same fundamentals. Machine coding needs speed + compiling code. Discussion needs clearer verbal trade-offs.

---

## 2. What interviewers evaluate

At senior / mid levels they score:

| Skill | What “good” looks like |
|-------|-------------------------|
| Problem breakdown | Clarifies ambiguity; states assumptions |
| Modeling | Correct entities & relationships |
| Abstraction | Interfaces where variation exists |
| SOLID / clean design | Responsibilities are separated |
| Patterns | Used purposefully, not decoration |
| APIs | Intuitive, consistent, versionable |
| Data & state | Clear persistence & transitions |
| Concurrency | Mentions races and how to prevent them |
| Failure handling | Timeouts, retries, partial failure |
| Extensibility | Can add a feature without rewriting core |
| Communication | Structured; trade-offs explained |

There is rarely one “correct” design. **Sound decisions + clear reasoning** wins.

---

## 3. LLD vs HLD vs DSA

| | DSA | LLD | HLD |
|---|-----|-----|-----|
| Focus | Algorithms, complexity | Classes, APIs, modules | Systems, scale, infra |
| Example | Shortest path in graph | Design Rate Limiter classes | Design URL shortener at 100M QPS |
| Output | Code + Big-O | Class diagram + APIs + flows | Boxes: LB, DB, cache, queues |

Many companies mix a little scale talk into LLD (“what if traffic grows?”). Answer at **component** level first, then briefly say how you’d evolve.

---

## 4. How a typical LLD round runs

Same **six content steps** (§5) — different **time budget** by round type (tables below). Every **`problems/<name>/README.md`** uses Steps 1–6 mapped to these clocks.

### Discussion LLD (~60 min)

Whiteboard / shared doc — classes, flows, APIs, trade-offs. **No full implementation.**

| Minutes | Activity | Maps to step |
|---------|----------|--------------|
| 0–5 | Clarify requirements & assumptions | Step 1 — Clarify |
| 5–12 | Entities + relationships | Step 2 — Entities & classes |
| 12–30 | Classes/interfaces + main flows | Step 3 — Flows |
| 30–40 | APIs (+ light DB schema if asked) | Step 4 — APIs |
| 40–50 | Concurrency, failures, extensibility | Step 5 — Deepen |
| 50–60 | Trade-offs, evolve design, Q&A | Step 6 — Evolve |

**Opening line:**

> “I’ll clarify the requirements first, then outline entities and responsibilities, design classes and key workflows, define APIs, and finally discuss concurrency, failure modes, and how the design evolves.”

### Machine coding (~90–120 min)

Working code for a small system — **same six steps**, compressed design phase, long coding phase.

| Minutes | Activity | Maps to step |
|---------|----------|--------------|
| 0–10 | Clarify + state v1 scope & non-goals | Step 1 — Clarify |
| 10–25 | Entities, classes, one happy-path flow (no coding yet) | Steps 2–3 |
| 25–85 | Implement core classes + happy path; compile often | Step 3 — Flows (code) |
| 85–105 | Public API, 2–3 edge cases, mutex/channel if shared state | Steps 4–5 |
| 105–120 | Narrate trade-offs + one evolution (cache, queue, strategy) | Step 6 — Evolve |

**Machine-coding opening line:**

> “I’ll clarify scope for v1, sketch entities and interfaces, then implement the main flow first, add edge cases and thread-safety where needed, and leave time to explain how I’d evolve this under load.”

### Interview round types (know which doc to use)

| Round type | What you do | Prep doc |
|------------|-------------|----------|
| **JavaScript / Node.js Q&A** | Language, event loop, Promises, closures, Node APIs | **[JavaScript/README.md](../../JavaScript/README.md)** — do **before** LLD if full-stack |
| **Golang Q&A** | GMP, goroutines, channels, context, GC | **[Go/README.md](../../Go/README.md)** — do **before** LLD if Go role |
| **Discussion LLD** (~60 min) | Classes, APIs, flows, trade-offs — no full coding | **[hub §1–§6](../../README.md)** — after language prep (top table) |
| **Machine coding** (90–120 min) | Working code for a small system | Prep guides + [hub §5](../../README.md#5-the-standard-approach-memorize-this) + **`JavaScript/*/`** + **`Go/*-go/`** |
| **Vibe coding / AI-assisted build** | Cursor/Copilot allowed; build + verify + narrate | **[vibe-coding-round/README.md](../../vibe-coding-round/README.md)** + [method doc](../method/README.md) §5 for design |
| **AI code review** | Clone repo, manual review — security, RAG, production gaps | **[ai-code-review-round/README.md](../../ai-code-review-round/README.md)** |
| **LLD breadth (paper only)** | Elevator, Chess, Logger, UML, extra patterns | **[lld-gaps/README.md](../../lld-gaps/README.md)** |

Same fundamentals (OOD, patterns, RAG vocabulary). **Language/runtime first (hub prep table), then design.** Design rounds = speak structure; review rounds = find bugs with file names.

---

## 5. The standard approach (memorize this)

```text
Understand → Model → Design → APIs → Deepen → Evolve
```

### Step-by-step

1. **Understand** — actors, use cases, constraints, non-goals  
2. **Model** — entities, relationships, state machines  
3. **Design** — classes, interfaces, responsibilities, collaborations  
4. **APIs** — endpoints or public methods of the library  
5. **Deepen** — concurrency, idempotency, caching, errors  
6. **Evolve** — “If requirements change / scale 10×, I would…”

Never jump to patterns or microservices before entities and use cases are clear.

---

## 6. Clarifying questions checklist

Ask a subset relevant to the problem:

**Product / scope**
- Who are the users (actors)?
- What are the top 3 use cases for v1?
- What is explicitly out of scope?

**Scale (order of magnitude)**
- Users? Requests/sec? Data size?
- Single machine library or distributed service?

**Functional**
- Auth / multi-tenant?
- Sync or async?
- Persistence needed? TTL?
- Exact vs approximate correctness?

**Quality**
- Consistency needs?
- Failure expectations?
- Observability needs?

Then **state assumptions out loud** so the interviewer can correct you.

### What Step 1 must include (every `problems/*/README.md`)

Step 1 is **Understand** — not the full design flow (that is Steps 2–3).

| Subsection | Purpose |
|------------|---------|
| **Questions (6–8)** | Ambiguity you ask the interviewer |
| **v1 expectations** | Actors, top use cases, in scope, out of scope, stated assumptions |
| **Confirm understanding** | One sentence happy path — proves you got the prompt |

**Do not put** class diagrams or step-by-step flows in Step 1 — those belong in **Step 2** (model) and **Step 3** (flows).

**Template:**

```markdown
## Step 1 — Clarify

### Questions (ask 6–8)
1. ...

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | ... |
| **Use cases (v1)** | 1. ... 2. ... |
| **In scope** | ... |
| **Out of scope** | ... |
| **Assumptions** | ... |

### Confirm understanding
> "For v1, [actor] can [main action]; we [won't build X yet]."
```

---

