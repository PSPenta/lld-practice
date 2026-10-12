# Staff / SDE-3 Interview Q&A — Node · Backend · Security

← [docs/](../README.md) · [Hub](../../README.md)

> **Bar:** Staff / SDE-3 — answers include **trade-offs**, **failure modes**, and **what you’d ship**.  
> **Depth elsewhere:** [JavaScript](../../JavaScript/README.md) · [TypeScript](../../TypeScript/README.md) · [backend](../backend/README.md) · [ai-lld](../ai-lld/README.md) · [principles](../principles/README.md) · [patterns](../patterns/README.md)

---

## How to use

1. Cover a section out loud in **2–4 minutes**.
2. Always add: **“What breaks in production?”** and **“How would you measure it?”**
3. Prefer **one concrete design** over listing buzzwords.

---

## Table of contents

1. [Node.js architecture](#1-nodejs-architecture)
2. [Event loop flow](#2-event-loop-flow)
3. [Microservice architecture](#3-microservice-architecture)
4. [Microservice design patterns](#4-microservice-design-patterns)
5. [SOLID principles](#5-solid-principles)
6. [MongoDB vs Postgres](#6-mongodb-vs-postgres)
7. [Indexes](#7-indexes)
8. [Kafka & Docker](#8-kafka--docker)
9. [JavaScript & TypeScript](#9-javascript--typescript)
10. [JWT & OAuth2](#10-jwt--oauth2)
11. [DevOps fundamentals](#11-devops-fundamentals)
12. [Agile](#12-agile)
13. [Performance tuning](#13-performance-tuning)
14. [Exception / error handling](#14-exception--error-handling)
15. [12-factor app](#15-12-factor-app)
16. [HTTP vs HTTPS](#16-http-vs-https)
17. [Design patterns](#17-design-patterns)
18. [Scenario — bank upload 50k records](#18-scenario--bank-upload-50k-records)
19. [`package.json`](#19-packagejson)
20. [Production deployment](#20-production-deployment)
21. [Automation testing](#21-automation-testing)
22. [Load testing](#22-load-testing)
23. [Memory leaks](#23-memory-leaks)
24. [GenAI / LLM](#24-genai--llm)
25. [Security & concurrency (Staff extras)](#25-security--concurrency-staff-extras)
26. [Rapid revision cheat sheet](#26-rapid-revision-cheat-sheet)

---

## 1. Node.js architecture

**Q: What is Node.js architecture?**

Node = **V8** (JS execution) + **libuv** (async I/O, thread pool) + **event loop** on one main JS thread.

```text
Your JS  →  V8 heap/stack  →  Event loop
                ↑
         libuv / OS (epoll, IOCP)
         · network, timers
         · thread pool (fs, dns, crypto) — default 4
```

**Staff add:** Node is great for **I/O-bound** APIs. CPU-bound work blocks the loop → use **worker_threads**, **child_process**, or a separate worker service.

**Q: Single-threaded — how does it scale?**

Concurrency via non-blocking I/O + many in-flight requests. Horizontal scale: **stateless** processes behind a load balancer (`cluster` / K8s replicas). Shared state in Redis/DB — not in process memory.

**Interview line:** “One JS thread per process; scale out processes; offload CPU; never block the loop.”

---

## 2. Event loop flow

**Q: Walk the event loop.**

Simplified Node phases (libuv):

```text
timers → pending callbacks → idle/prepare → poll (I/O) → check (setImmediate) → close
         microtasks (Promises, queueMicrotask, process.nextTick*) run between / after callbacks
```

\* `process.nextTick` drains **before** other microtasks — can starve the loop if recursive.

**Q: `setTimeout(0)` vs `Promise.then` vs `setImmediate`?**

| API | Queue | Order vs sync |
|-----|-------|----------------|
| Sync code | Call stack | First |
| `Promise.then` / `queueMicrotask` | Microtask | After current stack, before timers |
| `setTimeout(0)` | Timers (macrotask) | Later turn |
| `setImmediate` | Check phase | After poll (Node) |

**Output classic:** `console.log(2); Promise.resolve().then(() => console.log(1))` → `2` then `1`.

**Staff trap:** Long sync work (JSON.parse huge payload, crypto sync, tight loops) = **latency cliff** for all connections on that process.

---

## 3. Microservice architecture

**Q: When microservices vs monolith?**

| Prefer monolith | Prefer microservices |
|-----------------|----------------------|
| Small team, unclear domains | Clear bounded contexts, independent scale/deploy |
| Strong consistency, simple ops | Different SLAs, polyglot, team ownership |

**Staff:** Start modular monolith; split when **deploy cadence**, **scale**, or **failure isolation** demand it — not for résumé patterns.

**Q: What must every service have?**

- Clear **API contract** (versioned)
- **Own data** (no shared DB writes across services)
- **Idempotency** on writes
- **Observability:** `trace_id`, metrics, structured logs
- **Authn/z** at edge + service identity (mTLS / JWT)
- **Backpressure** (queues, rate limits, circuit breakers)

**Q: How do services communicate?**

Sync: HTTP/gRPC (latency, coupling). Async: Kafka/SQS (decouple, eventual consistency). Staff default: **commands sync when user waits; events async for fan-out.**

---

## 4. Microservice design patterns

### Saga

Long business transaction across services **without** distributed 2PC.

- **Choreography:** each service emits events; others react (simple, harder to see global state).
- **Orchestration:** a Saga coordinator drives steps (clearer, single point).

**Compensate** on failure (e.g. payment success → booking fail → refund).

### CQRS

Separate **write model** (commands) from **read model** (queries/projections). Use when read shapes ≠ write aggregates or read QPS ≫ writes. Cost: eventual consistency, dual schemas.

### Circuit breaker

Stop calling a failing dependency after threshold → **open** (fail fast) → **half-open** (probe) → **closed**. Prevents cascade. Pair with timeouts + bulkheads.

### API gateway

Single edge: auth, TLS, rate limit, routing, request aggregation. **Not** a place for heavy business logic. Staff: keep gateway thin; domain stays in services.

### Service discovery

Clients or sidecars find instances (K8s DNS/Service, Consul, Eureka). In Kubernetes, **Services + endpoints** usually suffice.

### ACID (and where it applies)

| Property | Meaning |
|----------|---------|
| Atomicity | All-or-nothing in one DB transaction |
| Consistency | Valid state transitions (constraints) |
| Isolation | Concurrent tx visibility (levels) |
| Durability | Committed data survives crash |

**Staff:** ACID inside **one** database. Across services use **Saga + idempotency + outbox**, not distributed XA.

---

## 5. SOLID principles

| Letter | Staff one-liner | Smell if violated |
|--------|-----------------|-------------------|
| **S** | One reason to change | God class / fat service |
| **O** | Extend via new types, not editing core | Endless `if (type === …)` |
| **L** | Subtypes honor base contract | Override that weakens guarantees |
| **I** | Small interfaces | Fat “Utils” / god interface |
| **D** | Depend on abstractions | `new ConcreteDb()` deep in domain |

**Interview:** Point to Strategy (RateLimiter), Factory (Splitwise), DIP (PaymentGateway) — see [principles](../principles/README.md) · [repo-map](../repo-map/README.md).

---

## 6. MongoDB vs Postgres

| | **Postgres** | **MongoDB** |
|--|--------------|-------------|
| Model | Relational, strong schema | Document (BSON), flexible |
| Transactions | Mature multi-row ACID | Multi-doc tx (since 4.x) — use carefully |
| Query | SQL, joins, rich analytics | Aggregation pipeline; joins via `$lookup` (costly) |
| Consistency | Strong default | Tunable; beware secondary reads |
| Best for | Money, relations, constraints | Variable docs, high write ingest, hierarchical JSON |

**Mongo architecture (Staff):**

- **mongod** processes; **replica set** (primary + secondaries) for HA
- **WiredTiger** storage engine; journals for durability
- **Sharding:** mongos → config servers → shards (shard key choice is critical — avoid hot shards)

**When to pick:** Ledger/payments → Postgres. Content/catalog with evolving shape → Mongo (or Postgres JSONB). Staff often: **Postgres primary + Redis cache**.

---

## 7. Indexes

**Q: What does an index do?**

Trade **write cost + disk** for **read speed**. B-tree (Postgres/Mongo default) for equality/range; hash for equality; GIN/GiST for full-text/JSONB (Postgres).

**Q: Design indexes for Staff.**

1. Index **filter + sort** columns used in hot queries.
2. **Composite** order: equality cols first, then range/sort (`(tenant_id, status, created_at)`).
3. Covering index / `INCLUDE` to avoid heap lookups when useful.
4. Partial indexes for hot subsets (`WHERE status = 'open'`).
5. Measure with `EXPLAIN (ANALYZE)` / Mongo `explain("executionStats")`.

**Anti-patterns:** Index every column; low-cardinality alone (`boolean`); unused indexes; functions on column preventing index use (`WHERE LOWER(email) = …` → use expression index).

**Q: Why is my query slow with an index?**

Wrong order in composite; selectivity too low; stale stats; sequential scan cheaper for large % of table; OR across non-indexed paths; hidden type casts.

---

## 8. Kafka & Docker

### Kafka

**What:** Distributed **commit log**. Producers append; consumers read by offset.

| Concept | Staff meaning |
|---------|----------------|
| Topic / partition | Parallelism unit; ordering **per partition** |
| Consumer group | Each partition → one consumer in group |
| Offset | Position; commit after processing (at-least-once typical) |
| Replication | ISR; durability vs lag trade-off |
| Exactly-once | Idempotent producer + transactions — complex; often at-least-once + idempotent consumers |

**vs SQS:** Kafka = replay, high throughput, ordering keys. SQS = simpler ops, competing consumers, less “event streaming platform.”

### Docker

Image = layered filesystem + metadata. Container = running instance. **Dockerfile** → build; **compose** for local multi-service.

**Staff production:** Containers on **K8s** — probes (liveness/readiness), resource limits, secrets not in image, non-root user, small base images, immutable tags (digest).

---

## 9. JavaScript & TypeScript

**JS Staff hits:** `this`/closures, event loop, Promises vs async, `==` vs `===`, prototypes vs `class`, module systems, memory (closures retaining DOM/caches).

**TS Staff hits:**

- `interface` vs `type` (extendable contracts vs unions/mapped)
- `unknown` > `any` at boundaries
- Discriminated unions + exhaustive `never`
- `strict` tsconfig; validate HTTP/JSON with Zod at edges
- Prefer composition + interfaces for LLD

Depth: [JavaScript/README.md](../../JavaScript/README.md) · [TypeScript/README.md](../../TypeScript/README.md).

---

## 10. JWT & OAuth2

### JWT

Signed (often HS256/RS256) claims: `sub`, `exp`, `iat`, optional `roles`/`tenant_id`. **Stateless** ⇒ hard revoke → short TTL + refresh tokens or blocklist.

**Staff security:** Validate `iss`/`aud`/`exp`; never trust payload without verify; don’t store secrets in JWT; prefer **httpOnly Secure** cookies for browsers over `localStorage`.

### OAuth2

Delegation: client obtains **access token** to call APIs on user’s behalf.

| Flow | Use |
|------|-----|
| Authorization code (+ PKCE) | Web/SPA/mobile — default |
| Client credentials | Service-to-service |
| Implicit | Legacy — avoid |

**OIDC** = identity layer on OAuth (ID token).

**Interview line:** “OAuth authorizes; OIDC authenticates; JWT is a token format — not a protocol by itself.”

---

## 11. DevOps fundamentals

### Git

Trunk-based or short-lived PRs; protected main; conventional commits optional. Staff cares about: **reviewable diffs**, bisectable history, no force-push to main, signed commits if required.

### Jenkins (and CI generally)

Pipeline as code (`Jenkinsfile`): build → test → lint → security scan → artifact → deploy. **Idempotent** deploys; secrets in vault/credentials store; separate build vs deploy stages; fail fast on tests.

**Staff:** Prefer GitHub Actions/GitLab CI/Tekton equally — concepts matter more than brand.

---

## 12. Agile

**Q: Scrum vs Kanban?** Scrum = timeboxed sprints, roles, ceremonies. Kanban = continuous flow, WIP limits. Staff: pick for **predictability vs continuous delivery**.

**Q: What’s your role as Staff in Agile?**

- Slice stories for **vertical value**
- Kill scope creep with clear NFRs
- Unblock teams (tech spikes, ADRs)
- Definition of Done includes **observability + rollback**

Expect 1–2 behavioral questions — answer with **situation → action → measurable result**.

---

## 13. Performance tuning

**Method (Staff):**

1. **Measure** — p50/p95/p99, CPU, memory, DB time, external calls (`trace_id`)
2. Find bottleneck (app vs DB vs network)
3. Fix **highest impact** first
4. Re-measure

**Node app levers:** avoid sync CPU; pool DB connections; cache hot reads; paginate; compress; HTTP keep-alive; stream large payloads; `UV_THREADPOOL_SIZE` if fs/crypto saturated.

**DB levers:** indexes, query rewrite, connection pooling, read replicas, batch writes.

**Never:** Optimize without numbers; premature micro-opts.

---

## 14. Exception / error handling

**Operational vs programmer errors:** retry timeouts/5xx from deps; crash-loop on bugs after log (or central handler).

**Node patterns:**

- `async/await` + `try/catch`; never swallow
- `unhandledRejection` / `uncaughtException` → log + **graceful shutdown** (don’t keep serving)
- Typed AppError with `statusCode` + stable `code`
- Map to HTTP: 400 validation, 401/403 auth, 404, 409 conflict, 429, 5xx

**Staff API errors:** problem+json or `{ error: { code, message, trace_id } }` — no stack traces to clients.

---

## 15. 12-factor app

Classic Heroku factors Staff still recite:

| # | Idea | Example |
|---|------|---------|
| I–III | Codebase, deps, config | One repo; lockfile; env/config service — not hardcoded secrets |
| IV–VI | Backing services, build/release/run, processes | DB/Redis as attached resources; immutable releases; stateless processes |
| VII–IX | Port binding, concurrency, disposability | Self-hosted HTTP; scale via process count; fast start/stop |
| X–XII | Dev/prod parity, logs, admin | Same artifacts; logs to stdout; one-off migrations as tasks |

**Interview line:** “Stateless app processes + externalized config + disposable deploys.”

---

## 16. HTTP vs HTTPS

| | HTTP | HTTPS |
|--|------|-------|
| Transport | Cleartext | TLS encrypted |
| Integrity | None | Protected against tampering |
| Auth of server | No | Certificate chain |
| SEO/browser | Mixed content blocked | Required for cookies `Secure`, HTTP/2 often |

**TLS handshake** (simplified): client hello → cert verify → key exchange → symmetric session keys.

**Staff:** HSTS, TLS 1.2+, strong ciphers, terminate TLS at LB or mesh; mTLS for service-to-service when needed.

---

## 17. Design patterns

Staff expects **when + trade-off**, not memorizing UML.

| Pattern | Use when | Caution |
|---------|----------|---------|
| Strategy | Swap algorithms (rate limit) | Don’t over-abstract one impl |
| Factory | Create family variants | Not a replacement for `new` everywhere |
| Observer/Pub-Sub | Fan-out events | Ordering, failure, backpressure |
| Circuit breaker | Protect from cascading failure | Fallback UX |
| Repository | Persist behind interface | Don’t make “repo” a second god class |
| Adapter | Wrap foreign API | Keep mapping thin |

Repo demos: [patterns](../patterns/README.md) · [repo-map](../repo-map/README.md).

---

## 18. Scenario — bank upload 50k records

**Prompt:** Design a banking app feature to upload ~50,000 records — how to process, what infra?

### Clarify

- File format? Max size? Sync UX or async job?
- Idempotency (re-upload)? Validation rules? Compliance (PII)?
- Latency expectation (minutes OK)?

### Staff design

```text
Client → API Gateway (auth, size limit)
      → Upload service: store object (S3) + enqueue job_id (202 Accepted)
Worker pool:
  · stream/parse file in chunks (CSV parser stream)
  · validate row batch
  · write via COPY / bulk insert (Postgres) in transactions of N rows
  · idempotency key = (upload_id, row_hash) or natural business key
  · DLQ + poison row quarantine
Progress: Redis/DB job status; webhook/SSE optional
```

**Infra:** Object storage, API (stateless), queue (SQS/Kafka), workers (autoscaled), Postgres (pooled), Redis (progress), observability, WAF/virus scan on upload.

**Why not parse 50k in HTTP request:** blocks event loop / timeouts; no retry story; memory spike.

**Interview line:** “Async job + chunked bulk load + idempotent rows + observable progress.”

---

## 19. `package.json`

**Key fields Staff mention:**

| Field | Why |
|-------|-----|
| `name` / `version` | Identity; semver |
| `type: "module"` | ESM vs CJS |
| `engines` | Node version pin |
| `dependencies` vs `devDependencies` | Prod image size / security |
| `scripts` | `start`, `test`, `lint`, `build` |
| `exports` / `main` / `types` | Package entrypoints |
| `package-lock.json` / `npm ci` | Reproducible installs |

**Staff:** Pin Node in CI and Docker; audit (`npm audit` / Snyk); avoid `*` versions; private registry for internal packages.

---

## 20. Production deployment

**Happy path:**

```text
CI green → build immutable image (git SHA tag)
     → push registry → deploy (K8s rolling / blue-green / canary)
     → migrate DB (expand/contract; never break old app mid-rollout)
     → health checks pass → shift traffic
     → monitor error rate / latency → auto or manual rollback
```

**Staff checklist:** secrets from vault; config via env; readiness ≠ liveness; DB migrations backward compatible; feature flags; runbooks; on-call.

---

## 21. Automation testing

| Layer | Tools | Staff use |
|-------|-------|-----------|
| Unit | Jest / Vitest / Mocha | Domain logic, pure functions |
| Integration | Supertest + testcontainers | API + DB |
| E2E UI | **Cypress** / Playwright / Selenium | Critical user journeys only |
| Contract | Pact | Provider/consumer APIs |

**Cypress vs Selenium:** Cypress = modern, fast feedback, JS-native; Selenium = multi-browser/language, heavier. Staff: **few E2E**, many unit/integration — pyramid.

---

## 22. Load testing

**Tools:** k6, Locust, JMeter, Artillery.

**What to measure:** RPS, latency percentiles, error %, saturation (CPU, DB connections, queue lag).

**Staff method:** Define SLO → baseline → step load → soak (memory leaks) → spike. Test **prod-like** data volume. Never load-test prod without approval.

---

## 23. Memory leaks

**Node symptoms:** Rising `process.memoryUsage().heapUsed` / RSS; eventual OOM kill.

**Common causes:**

- Global caches without bounds/TTL
- EventEmitter listeners not removed
- Closures retaining large objects
- Uncleared `setInterval` / timers
- Growing arrays/maps keyed by request id
- Native buffers retained

**Debug:** `node --inspect`, Chrome DevTools heap snapshots, `clinic.js`, compare snapshots; watch `heapUsed` under soak test.

**Fix:** Bound caches (LRU); weak refs where appropriate; clear listeners; stream don’t buffer entire files.

Also: [JavaScript §22](../../JavaScript/README.md) GC notes.

---

## 24. GenAI / LLM

**Staff LLD angles** (not model training):

| Topic | Answer |
|-------|--------|
| Architecture | Client → API → **LLMClient** interface → provider (OpenAI/Bedrock) |
| RAG | Retrieve **tenant-scoped** chunks → prompt → generate; cite sources |
| Cost/abuse | Credits reserve/commit; rate limit; max tokens |
| Safety | Timeout, output validation, PII redaction, prompt injection awareness |
| Reliability | Retries with jitter; fallback model; degrade to non-AI path |
| Latency | Stream tokens; async jobs for long work |

Depth: [ai-lld](../ai-lld/README.md) · [ai-code-review-round](../../ai-code-review-round/README.md).

---

## 25. Security & concurrency (Staff extras)

### Concurrency (Node)

- JS is single-threaded — **no parallel mutation of JS memory** across requests on one process, but **async interleaving** causes races on shared state (in-memory counters, caches).
- Cross-process: Redis `INCR` / Lua / DB row locks / optimistic versioning.
- `worker_threads` / `child_process` for CPU (see [JS child process notes](../../JavaScript/README.md)).

**Classic race:** read-modify-write balance without transaction → lost updates. Fix: `UPDATE … WHERE version = ?` or serializable tx / Redis lock.

### Security checklist (expect questions)

| Area | Staff answer |
|------|----------------|
| Authn | OIDC/OAuth; short-lived access tokens; secure session cookies |
| Authz | RBAC/ABAC; check **every** request; never trust client roles alone |
| Injection | Parameterized SQL; escape/ORM; sanitize HTML; no `eval` |
| SSRF | Allowlist URLs; block link-local/metadata IPs |
| Upload | Type/size limits; store outside web root; virus scan; signed download URLs |
| Secrets | Vault/KMS; rotate; never commit `.env` |
| Dependency | Lockfile + audit + SBOM |
| Transport | HTTPS everywhere; HSTS; mTLS internal |
| Logging | No passwords/tokens/PII in logs; redact |
| Multi-tenant | `tenant_id` on every query — IDOR is a Staff killer question |

### OWASP Top 10 (name a few)

Broken access control · cryptographic failures · injection · insecure design · security misconfig · vulnerable components · auth failures · integrity failures · logging failures · SSRF.

---

## 26. Rapid revision cheat sheet

```text
Node              → V8 + libuv + one JS thread; don’t block the loop
Event loop        → sync → microtasks → timers/I/O macrotasks
Microservices     → modular monolith first; own data; Saga across services
SAGA/CQRS/CB/GW   → compensate; split read/write; fail fast; thin gateway
SOLID             → Strategy/Factory/DIP in this repo
Postgres vs Mongo → relations/ACID vs flexible docs/sharding discipline
Indexes           → composite order; EXPLAIN; don’t over-index
Kafka             → log + partitions; at-least-once + idempotent consumer
Docker/K8s        → immutable image; probes; limits; secrets out of image
JWT/OAuth         → short JWT + refresh; auth code+PKCE; OIDC for login
12-factor         → config in env; stateless processes; logs stdout
HTTPS             → TLS confidentiality + integrity + server auth
50k upload        → S3 + queue + chunked workers + bulk insert + idempotency
Perf              → measure p99 → fix bottleneck → re-measure
Errors            → typed AppError; graceful shutdown on fatal
Memory leak       → bound caches; remove listeners; soak + heap snapshot
Security          → authz every request; parameterized queries; tenant isolation
LLM               → LLMClient + RAG tenant scope + credits + timeouts
```

---

## Final note

Staff interviews reward **judgment under constraints**. For every answer: **default design → failure mode → metric → evolution**. Cross-link deeper docs when they ask you to go deeper — don’t dump the encyclopedia unprompted.
