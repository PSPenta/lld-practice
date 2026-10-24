# AI-Native Open Challenge — Preparation Guide

> **Based on:** Open-ended take-homes like **Billeasy’s “Build anything” AI challenge** (AI Native Tech Lead / AI Native Engineer tracks) and the broader 2025–2026 shift to **portfolio + process** over prescribed coding tests.  
> **Related:** [vibe-coding-round](../vibe-coding-round/README.md) (prescribed problem + AI IDE) · [ai-code-review-round](../ai-code-review-round/README.md) · hub LLD [../README.md](../README.md)  
> **This is not classic LLD.** Same skills (judgment, architecture, ownership) — different brief: **you choose the problem**; they score **how you use AI**, not the idea alone.

---

## Table of contents

1. [What this round is](#1-what-this-round-is)
2. [What they score](#2-what-they-score)
3. [Billeasy / AI-Native company signals](#3-billeasy--ai-native-company-signals)
4. [How it differs from vibe coding / Chakra ADE](#4-how-it-differs-from-vibe-coding--chakra-ade)
5. [Playbook — pick, build, document](#5-playbook--pick-build-document)
6. [Write-up template (match their bullets)](#6-write-up-template-match-their-bullets)
7. [Idea filters (personality + craft)](#7-idea-filters-personality--craft)
8. [Failure modes](#8-failure-modes)
9. [Rapid revision cheat sheet](#9-rapid-revision-cheat-sheet)

---

## 1. What this round is

Typical brief (paraphrased from invites like Billeasy’s):

```text
Build anything.
No prescribed problem, format, length, or tool ban.
Submit: artifact + short write-up on HOW you used AI.
Deadline: fixed (e.g. a few days).
Note: personalities — not “just another techie.”
```

| They say | They mean |
|----------|-----------|
| “Not the idea itself” | Idea is a vehicle; **process & AI depth** are the product |
| “Any form” (repo, link, video, zip, doc) | Pick a medium that **shows** the work — don’t hide behind a PDF only |
| “Personality” | Taste, judgment, voice, what you care about — not a sterile CRUD clone |
| “How AI shaped decisions” | Moments of **steering, surprise, rejection** of AI — not a tool name-drop list |

**Interview line:** “I’m submitting a small, opinionated build plus a process diary — proof I lead with AI without surrendering ownership.”

---

## 2. What they score

| Dimension | Strong signal | Weak signal |
|-----------|---------------|-------------|
| **AI integration depth** | Tools chosen for roles; multi-step orchestration; you verify | One ChatGPT paste; no critique |
| **Deliberate tool choice** | “Claude Code for impl, Cursor for nav, Whisper for X…” | “I used AI” with no roles |
| **Process from start→finish** | Clarify → cut scope → milestones → verify → ship | Idea → dump → polish write-up only |
| **Decision moments** | Architecture / scope pivots attributed to AI or to rejecting AI | No forks; linear fairy tale |
| **Surprise & adaptation** | Honest “model hallucinated X → I did Y” | Only praise; no friction |
| **Personality / taste** | Domain you care about; constraints that sound like you | Generic todo app with no voice |
| **Ownership** | You can defend every major choice in a follow-up call | Can’t explain generated parts |

Overlaps vibe-coding rubrics ([Underdog / AI-assisted assessments](../vibe-coding-round/README.md#2-what-interviewers-score)) — here **scope selection** is itself a scored skill.

---

## 3. Billeasy / AI-Native company signals

Public hiring language (AI Native Tech Lead / AI Native Software Engineer):

- **AI-first culture** — Claude Code called out as a primary daily tool; Copilot/Cursor also welcome.
- **“AI drafts; you own the quality”** — same ownership bar as vibe coding.
- **TypeScript / full-stack / Postgres–Mongo** appear on eng JDs — not required for the open challenge, but a **billing / ops / retail** flavoured prototype can resonate with [billeasy.in](https://billeasy.in)-style product space without copying their product.
- Open round = filter for **AI-native personality** before (or beside) classic LLD / leadership loops.

There is little public Glassdoor-style write-up of this exact challenge; treat the **email brief as the rubric**.

---

## 4. How it differs from vibe coding / Chakra ADE

| | **Open AI challenge** (this doc) | **Vibe coding / ADE** |
|--|----------------------------------|------------------------|
| Problem | **You invent it** | Given (README + tests / feature) |
| Time | Days (async take-home) | Live 60–120 min |
| Artifact | Build **+ process write-up** | Working code under observation |
| Personality | Explicitly scored | Secondary to verification |
| LLD | Optional — use if the build needs structure | Often required (entities first) |

Use [vibe-coding SCOPE](../vibe-coding-round/README.md#5-the-scope-workflow-memorize-this) **inside** the build once you’ve chosen the problem. Use [hub §5](../README.md#5-the-standard-approach-memorize-this) if you ship a small system design.

---

## 5. Playbook — pick, build, document

```text
1. Pick a sharp problem (personal + finishable in the deadline)
2. Name tools + role for each (before heavy coding)
3. Keep a running log: decisions, surprises, rejects
4. Build a thin vertical slice that actually runs / demos
5. Write-up last — but sourced from the log (not rewritten fiction)
6. Package: README + link/repo/video + write-up
```

| Step | Do | Don’t |
|------|----|-------|
| **Pick** | One user, one pain, one demo path | “Platform for everything” |
| **Tools** | 2–4 tools with clear jobs | Ten logos, no narrative |
| **Build** | Shippable MVP; tests or scripted demo | Half-finished mega-architecture |
| **Document** | Answer their four bullets exactly | Generic “I love AI” essay |
| **Personality** | Constraints, humour, domain taste | Anonymous boilerplate |

---

## 6. Write-up template (match their bullets)

Copy structure; keep voice yours.

```text
## What I built and why
- One paragraph: artifact + who it’s for
- Why this (personal / work taste) — not “because it was easy”

## AI tools, models, platforms — and roles
| Tool / model | Role in this project |
|--------------|----------------------|
| …            | Spec / research / code / review / design / voice |

## How AI influenced decisions (key moments)
1. Moment — what I asked — what changed (architecture / scope / API)
2. …
3. …

## Surprises / unexpected / approach changes
- Hallucination or wrong API → how I caught it
- Suggestion I rejected and why
- Pivot that improved the build

## How to run / watch
- Link, clone steps, or video timestamp
```

---

## 7. Idea filters (personality + craft)

Good open-challenge ideas tend to be:

- **Personal** — something you’d actually use (interview prep harness, billing reconciler toy, meeting-notes → tickets, DSA drill coach tied to *your* notes).
- **AI-shaped** — AI is in the loop of the product *or* clearly in the build process (better: both).
- **Finite** — demo in &lt;5 minutes; README has one happy path.
- **Opinionated** — a constraint that reveals taste (offline-first, no backend, TS-strict, cost-capped LLM calls).

Weak: another auth CRUD, “ChatGPT wrapper” with no craft, copy of a famous tutorial with renamed vars.

---

## 8. Failure modes

| Failure | Why it hurts |
|---------|--------------|
| Generic take-home energy | They asked for personality |
| Tool dump, no roles | Looks like checkbox AI |
| No surprise / no rejection of AI | Unrealistic; weak ownership |
| Idea so big it doesn’t run | Process without proof |
| Perfect polished lie | Follow-up call exposes it — keep the log honest |
| Zero architecture when the artifact needs it | Tech Lead bar still expects judgment |

---

## 9. Rapid revision cheat sheet

```
Open AI challenge     → YOU pick problem; score = AI process + personality
Submit                → artifact + write-up (4 bullets from invite)
Not classic LLD       → but use SCOPE / §5 if you build a system
Billeasy signal       → AI-native, Claude Code culture, own the quality
vs vibe/ADE           → prescribed problem + live IDE vs open + diary
Win                   → deliberate tools, decision moments, honest surprises, ships
Fail                  → sterile CRUD, tool name-drop, unfinished mega-scope
Craft docs            → vibe-coding-round · hub LLD method
```

---

## Final note

Treat this as a **proof-of-work portfolio piece**, not a school assignment. The write-up is half the grade. If a follow-up call happens, rehearse walking the **decision log** — same muscle as narrating vibe-coding verification.
