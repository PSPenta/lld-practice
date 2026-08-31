# Vibe Coding / AI-Assisted LLD Round — Preparation Guide

> **Based on:** 2025–2026 hiring trends where interviews allow **Cursor, Copilot, Claude, Windsurf** and score **prompt quality, verification, ownership, and orchestration** — not syntax recall.  
> **For whiteboard LLD (no AI):** [../README.md](../README.md)  
> **For manual code review (no AI generation):** [../ai-code-review-round/README.md](../ai-code-review-round/README.md)  
> **For LLD topics not in repo code:** [../lld-gaps/README.md](../lld-gaps/README.md)  
> **References:** [Underdog.io — Vibe Coding Interviews](https://underdog.io/blog/vibe-coding-interviews-for-engineers) · [Viblo — AI-Assisted Coding Assessments](https://viblo.asia/p/vibe-coding-interview-guide-ace-ai-assisted-coding-assessments-ZoJjeGDA4Y7)

---

## Table of contents

1. [What this round is (and is not)](#1-what-this-round-is-and-is-not)
2. [What interviewers score](#2-what-interviewers-score)
3. [Interview formats you'll see](#3-interview-formats-youll-see)
4. [Vibe coding + LLD — how they combine](#4-vibe-coding--lld--how-they-combine)
5. [The SCOPE workflow (memorize this)](#5-the-scope-workflow-memorize-this)
6. [Prompting for LLD problems](#6-prompting-for-lld-problems)
7. [Verification — catch AI mistakes](#7-verification--catch-ai-mistakes)
8. [Live session minute-by-minute](#8-live-session-minute-by-minute)
9. [Company patterns (2025–2026)](#9-company-patterns-20252026)
10. [Practice plan](#10-practice-plan)
11. [Speak templates](#11-speak-templates)
12. [Failure modes to avoid](#12-failure-modes-to-avoid)
13. [Rapid revision cheat sheet](#13-rapid-revision-cheat-sheet)

---

## 1. What this round is (and is not)

### Definition (interview meaning)

**Vibe coding** in interviews = **AI-assisted engineering where you remain engineer of record**. You describe intent, AI generates drafts, **you verify, steer, and own** the result.

It is **not** “accept whatever the model prints” (Karpathy’s provocative original phrase). Interviewers reject candidates who don't read AI output.

```text
Traditional          Vibe coding (interview)        Fail mode
───────────          ───────────────────────        ─────────
Write every line  →  Prompt → Review → Steer   →   Approve without reading
```

### Round type map (this repo)

| Round | AI allowed? | What you do | Prep doc |
|-------|-------------|-------------|----------|
| **Discussion LLD** | Usually **no** | Classes, APIs on whiteboard | [../README.md](../README.md) |
| **Machine coding** | Sometimes **yes** | Build working slice in 90–120 min | This doc + [../README.md](../README.md) |
| **Vibe coding / AI-assisted build** | **Yes** | Feature or LLD spike with Cursor/Copilot | **This doc** |
| **AI code review** | N/A | Find bugs in existing repo **manually** | [../ai-code-review-round](../ai-code-review-round/README.md) |
| **Review AI-generated code** | N/A | Interviewer gives AI slop; you critique | This doc §7 + ai-code-review |

### Why companies use it

- Most developers use AI tools daily; interviews should reflect **real work**.
- Signal shifts from **syntax recall** → **decomposition, verification, tradeoffs**.
- Especially common at **startups**, **Meta/Shopify pilots**, **AI-native companies**.

---

## 2. What interviewers score

From hiring-manager rubrics ([Underdog.io](https://underdog.io/blog/vibe-coding-interviews-for-engineers), [KnowledgeHut 2026](https://www.knowledgehut.com/blog/artificial-intelligence/vibe-coding-interview-questions)):

| Dimension | Strong signal | Weak signal |
|-----------|---------------|-------------|
| **Verification** | Run tests, catch wrong imports, fix edge cases aloud | “Looks good” → ship |
| **Prompt quality** | Scoped, iterative, references constraints | One giant “build everything” prompt |
| **Ownership** | “I chose Strategy because…” | Can't explain AI's design |
| **Orchestration** | Know when to prompt vs type yourself vs ask clarifying Q | Fight the model for 20 min |
| **Decomposition** | Entities + milestones **before** first prompt | Open AI immediately |
| **Communication** | Narrate while reviewing | Silent paste |

**Interview line:** “I'm using AI for boilerplate; **I'm responsible for architecture, correctness, and what we merge.**”

---

## 3. Interview formats you'll see

| Format | Duration | Task | LLD overlap |
|--------|----------|------|-------------|
| **Live AI-paired build** | 60–90 min | REST API, CLI, component | Implement Rate Limiter / Parking API |
| **Take-home + review** | 2–8 hr + call | Mini product slice | “Add caching to this service” |
| **Hybrid** | DS&A then AI round | Algo + feature build | Meta-style loops |
| **System design + AI spike** | 60 min | Design + **proof-of-concept** | Design notifier + spike one channel |
| **Review AI output** | 45 min | Find bugs in generated code | Same skills as code review round |
| **Extend existing repo** | 60–120 min | Add endpoint in codebase | Match patterns in `JavaScript/RateLimiter2/` |
| **Agentic round** | 60+ min | Direct Claude Code / Cursor agent | Senior+ — task spec + intervene |

---

## 4. Vibe coding + LLD — how they combine

**LLD knowledge still required.** AI does not replace:

- Clarifying requirements
- Choosing entities and relationships
- Picking Strategy vs State vs Observer
- API contracts, concurrency, idempotency

### Format A — “Design then build” (most common LLD + vibe)

```text
0–10 min   Clarify + entities + class sketch (YOU — whiteboard or comments)
10–15 min  Milestone plan (5 steps max)
15–45 min  AI assists implementation per milestone; YOU review each
45–55 min  Tests + edge cases + explain tradeoffs
55–60 min  Extensibility (“add PayPal later”)
```

**Use [../README.md](../README.md) §5 for design; use AI only after the model exists.**

### Format B — “Build and explain design”

Interviewer gives: *“Build a rate limiter with token bucket in the next hour.”*

You must **still speak LLD** while coding:

> “I'll separate `RateLimiter` from `RateLimiterStrategy` so we can swap algorithms — I'll have AI scaffold the interface, but I'll verify `allow()` handles burst correctly.”

**Repo reference after session:** compare with [../JavaScript/RateLimiter2/](../JavaScript/RateLimiter2/).

### Format C — “AI design doc + human review”

You prompt for class diagram or API spec, then **critique** it:

> “The model suggested a god-class `NotificationManager` — I'd split into `Channel`, `TemplateEngine`, and `DeliveryService` per SRP.”

---

## 5. The SCOPE workflow (memorize this)

Adapted from common interview frameworks ([Viblo guide](https://viblo.asia/p/vibe-coding-interview-guide-ace-ai-assisted-coding-assessments-ZoJjeGDA4Y7)):

| Step | Action | LLD tie-in |
|------|--------|------------|
| **S — State the problem** | Repeat requirements; list assumptions | Same as [../README.md](../README.md) §6 |
| **C — Cut scope** | MVP for 60 min: “create + redirect only, no analytics” | YAGNI |
| **O — Outline architecture** | Entities, interfaces, 5 milestones | §5 method — **before AI** |
| **P — Prompt with constraints** | Language, patterns, files, tests required | §6 below |
| **E — Evaluate & iterate** | Run, test, diff, fix; next milestone | §7 verification |

```text
     ┌─────────────┐
     │ S: Clarify  │
     └──────┬──────┘
            ▼
     ┌─────────────┐
     │ C: Cut scope│
     └──────┬──────┘
            ▼
     ┌─────────────┐     NO AI until here
     │ O: Outline  │ ─────────────────────
     └──────┬──────┘
            ▼
     ┌─────────────┐
     │ P: Prompt   │◄──┐
     └──────┬──────┘   │
            ▼          │
     ┌─────────────┐   │
     │ E: Verify   │───┘ iterate
     └─────────────┘
```

---

## 6. Prompting for LLD problems

### Bad vs good prompts

```text
BAD:  "Build a complete parking lot system with payment and SMS"

GOOD: "Node.js ES modules. Milestone 1 only:
       - Classes: ParkingLot, Floor, Slot, Vehicle enum (CAR, BIKE)
       - Method: park(vehicle) returns ticketId or throws NoSlot
       - No payment yet. Match style: one class per file.
       - Add 3 unit tests for park when full."
```

### Prompt template (copy for interviews)

```text
Context: [language, framework, existing repo path if any]
Milestone: [number] of [total] — [name only]
Architecture (fixed — do not change):
  - [Interface X with methods ...]
  - [Class Y composes Z]
Requirements:
  - [functional]
  - [edge cases]
Out of scope: [...]
Deliver: [files], [tests], no extra features
After generation I will review [specific risk: threading, auth, etc.]
```

### LLD-specific prompts (examples)

**Rate limiter:**

> “Implement `RateLimiterStrategy` interface with `allow(key): boolean`. First strategy: token bucket. Separate files. Do NOT put algorithm inside RateLimiter — composition only. JavaScript, no external libs.”

**Parking lot:**

> “Scaffold Floor/Sslot/ParkingLot only. Strategy for slot assignment: first-available. I'll review LSP on Vehicle types myself.”

**Notification (from [../lld-gaps/README.md](../lld-gaps/README.md)):**

> “Interface `NotificationChannel { send(msg) }`. Implement `EmailChannel` only. Stub `SmsChannel`. Template: `Hello {{name}}`. No queue yet — milestone 2.”

---

## 7. Verification — catch AI mistakes

AI **will** get these wrong — say you check for them aloud:

| Category | Common AI bug |
|----------|---------------|
| **Security** | SQL injection, hardcoded secrets, missing auth |
| **Concurrency** | Race on shared map; no lock on check-then-act |
| **LLD** | God class; no interface; Strategy as inheritance |
| **Edge cases** | Empty input, capacity full, duplicate idempotency key |
| **API** | Wrong HTTP status; no error body |
| **Dependencies** | Hallucinated package or wrong import path |
| **Tests** | Tests that mock nothing; always pass |

### Verification checklist (after each milestone)

- [ ] Run / compile / test
- [ ] Read diff line-by-line for **scope creep**
- [ ] Trace one request manually
- [ ] Name one thing you'd fix before production
- [ ] Explain every class **you** didn't write — can you own it?

**Overlap with code review:** same muscle as [../ai-code-review-round/README.md](../ai-code-review-round/README.md) — vibe coding **generates** the code to review; code review round **finds** planted bugs without AI.

---

## 8. Live session minute-by-minute

### 60-minute AI-assisted LLD build

| Min | You do | Say aloud |
|-----|--------|-----------|
| 0–5 | Clarify actors, MVP scope | “I'll skip payment for v1…” |
| 5–12 | Entity list + sketch diagram | “ParkingLot composes Floor, not inherits…” |
| 12–15 | Milestone plan (3–5 steps) | “M1: domain model + park()…” |
| 15–35 | AI milestone 1–2; **you review each** | “Model used sync file read — I'll fix…” |
| 35–45 | AI milestone 3 or hand-write critical path | “I'll write allow() myself — core logic” |
| 45–52 | Tests + one edge case | “Full lot → error, not hang” |
| 52–60 | Extensibility + tradeoffs | “Swap strategy via interface…” |

### When to **stop using AI**

- Core algorithm (token bucket math, LRU list splice)
- Concurrency-sensitive section
- When AI already failed twice on same task — **take the keyboard**

---

## 9. Company patterns (2025–2026)

| Company | Pattern | Prep note |
|---------|---------|-----------|
| **Shopify** | Multiple AI-enabled rounds; fix “garbage” live | Strong verification narrative |
| **Meta** | CoderPad + models; E7+ may replace one coding round | Hybrid DS&A + AI |
| **Google** | Pilot: human-led, AI-assisted comprehension (2026) | DS&A still AI-free |
| **Stripe** | **AI prohibited** in interviews | Train **without** AI too |
| **Amazon** | Classic OOP/LD + LP (as of 2026 reports) | [../README.md](../README.md) whiteboard |
| **Startups** | Take-home + “how did you use AI?” call | Document your prompts honestly |

**Always ask recruiter:** Which tools are allowed? Screen share? Internet? Copy-paste from your own notes?

---

## 10. Practice plan

| Day | Activity |
|-----|----------|
| **1** | Read this doc + [../README.md](../README.md) §5. No AI. Paper-design Rate Limiter. |
| **2** | **SCOPE** build: Parking Lot M1 only in Cursor (45 min timed). Compare [../JavaScript/ParkingLot2/](../JavaScript/ParkingLot2/). |
| **3** | **Verification drill:** Prompt AI for LRU cache; find 3 bugs yourself before running. |
| **4** | **Extend repo:** Add one strategy to RateLimiter2 with AI — must match existing pattern. |
| **5** | **Mock:** 60 min live build + record yourself narrating. |
| **6** | **Review round cross-train:** [../ai-code-review-round](../ai-code-review-round/) — Format 5 prep. |
| **7** | **Gap topic:** Notification or Elevator — design on paper ([../lld-gaps](../lld-gaps/)), then AI milestone 1 only. |

---

## 11. Speak templates

### Opening

> “I'll clarify scope first, sketch entities and interfaces, then use AI for scaffolding while I review every milestone and keep core logic under my control.”

### While reviewing AI output

> “The structure is close but this violates SRP — I'll split `Handler` from `Service` before we continue.”

### When AI is wrong

> “This misses idempotency on retry — I'll add a dedup key before calling send.”

### Closing

> “With more time I'd add a second channel behind the same interface and persistence behind a repository — the orchestration wouldn't change.”

---

## 12. Failure modes to avoid

| Failure | Why it fails |
|---------|--------------|
| **Prompt-first, no design** | Looks like AI architected it |
| **Blind acceptance** | Verification score zero |
| **Ignoring tool when disallowed** | Immediate fail at Stripe-like shops |
| **Can't explain generated code** | Ownership score zero |
| **Scope creep cleanup** | Wastes time; breaks tests |
| **No tests** | “Working” unverified |
| **Fighting agent 30 min** | Bad orchestration — reset scope |
| **Pure vibe, zero LLD vocabulary** | Misses OOP/LD round expectations |

---

## 13. Rapid revision cheat sheet

```
Vibe coding (interview) → AI drafts; YOU verify, steer, own
NOT → accept output blindly (Karpathy literal "vibes")

Scored on              → verification, prompts, ownership, orchestration
SCOPE                  → State, Cut scope, Outline, Prompt, Evaluate
Outline BEFORE AI      → entities + interfaces + milestones

LLD + vibe             → you design; AI scaffolds; you review each step
Format A               → design 10 min → build with AI → tests → extend
Format review          → same as ai-code-review-round

Prompt                 → one milestone, constraints, out-of-scope, tests
Verify                 → run, read diff, security, concurrency, edges
Stop AI when           → core algo, 2nd failure, concurrency

Stripe                 → AI often banned — practice without
Shopify/Meta/startups  → AI expected — narrate verification

Repo practice          → JavaScript/RateLimiter2, ParkingLot2
Gap topics             → lld-gaps/ (Elevator, Notification, …)
Whiteboard LLD         → ../README.md
Manual code review     → ../ai-code-review-round/
```

---

## Final note

**Vibe coding does not replace LLD.** It changes **how you produce code**, not **what you must know**. The strongest candidates use [../README.md](../README.md) for **design thinking** and this doc for **AI workflow under observation**.

Good luck — and always **read the diff**.
