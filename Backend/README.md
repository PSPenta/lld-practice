# REST API & Backend — Interview Preparation Guide

> **SDE-3 / Staff / Lead** — REST design, HTTP semantics, idempotency, payments, pagination, auth, platform APIs.  
> **Pair with:** [../Go/README.md](../Go/README.md) (Go HTTP/production) · [../JavaScript/README.md](../JavaScript/README.md) (Node runtime, JWT §26) · [../README.md](../README.md) (LLD method, [Rate Limiter design](../JavaScript/RateLimiter2/README.md), idempotency §12) · [../Go/PaymentGateway-go/](../Go/PaymentGateway-go/) (Strategy + gateway port)

---

## Table of Contents

0. [How to use this doc & topic map](#0-how-to-use-this-doc--topic-map)
1. [REST principles](#1-rest-principles)
2. [HTTP methods — POST vs PUT vs PATCH vs DELETE](#2-http-methods--post-vs-put-vs-patch-vs-delete)
3. [Idempotency & safe retries](#3-idempotency--safe-retries)
4. [Designing an idempotent payment API](#4-designing-an-idempotent-payment-api)
5. [API versioning](#5-api-versioning)
6. [Pagination — offset vs cursor](#6-pagination--offset-vs-cursor)
7. [Rate limiting](#7-rate-limiting)
8. [Authentication vs authorization](#8-authentication-vs-authorization)
9. [JWT & OAuth (interview depth)](#9-jwt--oauth-interview-depth)
10. [API gateway](#10-api-gateway)
11. [Request validation](#11-request-validation)
12. [Error response design & HTTP status codes](#12-error-response-design--http-status-codes)
13. [Backward compatibility & evolving APIs](#13-backward-compatibility--evolving-apis)
14. [APIs consumed by millions of developers](#14-apis-consumed-by-millions-of-developers)
15. [Practice checklist](#15-practice-checklist)
16. [Rapid revision cheat sheet](#16-rapid-revision-cheat-sheet)

---

## 0. How to use this doc & topic map

### Revision path (1–2 days)

| Day | Focus | Sections |
|-----|-------|----------|
| **1** | HTTP semantics + idempotency + payments | §1–§4, **§12** status codes |
| **2** | Platform APIs + auth + drill | §5–§11, §13–§14, **§16** cheat sheet |

**Before interview:** whiteboard **§4** payment flow; answer **§15** aloud; skim **§16** (20 min).

### Interview topic map

| Topic | Section |
|-------|---------|
| REST principles | **§1** |
| PUT vs PATCH, POST vs PUT | **§2** |
| Idempotency, retries | **§3** |
| Idempotent payment API | **§4** |
| Versioning | **§5** |
| Cursor vs offset pagination | **§6** |
| Rate limiting | **§7** |
| Auth vs authz, JWT, OAuth | **§8–§9** |
| API gateway | **§10** |
| Request validation | **§11** |
| Error shape, status codes | **§12** |
| Backward compat, API evolution | **§13** |
| Developer platform at scale | **§14** |

---

## 1. REST principles

**REST** (Representational State Transfer) = architectural style for **resource-oriented** HTTP APIs.

| Principle | Interview one-liner |
|-----------|---------------------|
| **Resources** | Nouns in URLs: `/payments/{id}`, not `/createPayment` |
| **Representations** | JSON/XML of resource state; same resource, multiple formats via `Accept` |
| **Stateless** | Server stores no client session in memory between requests — auth via token/key each call |
| **Uniform interface** | Standard HTTP methods + status codes + hypermedia optional (HATEOAS) |
| **Cacheable** | GET responses cacheable when safe; use `Cache-Control`, `ETag` |
| **Layered** | Client → gateway → service → DB; client doesn't know internal topology |

**Not REST purism in interviews:** pragmatic JSON APIs with clear resources + verbs are enough. Mention **RPC-style** (`POST /payments/{id}/capture`) when action doesn't map cleanly to PUT on a resource.

**Interview line:** “I design around **resources and state transitions**, use HTTP semantics correctly, keep handlers stateless, and push session state to tokens or DB.”

---

## 2. HTTP methods — POST vs PUT vs PATCH vs DELETE

| Method | Purpose | Idempotent? | Safe? | Typical use |
|--------|---------|-------------|-------|-------------|
| **GET** | Read resource | Yes | Yes | `GET /v1/payments/{id}` |
| **POST** | Create sub-resource / action | **No** (usually) | No | `POST /v1/payments` create |
| **PUT** | Replace entire resource at URI | **Yes** | No | `PUT /v1/customers/{id}` full replace |
| **PATCH** | Partial update | Often yes if designed | No | `PATCH /v1/customers/{id}` `{ "email": "..." }` |
| **DELETE** | Remove resource | Yes | No | `DELETE /v1/webhooks/{id}` |

### POST vs PUT

| | **POST** | **PUT** |
|---|----------|---------|
| **URI** | Collection: `POST /payments` | Known ID: `PUT /payments/pay_123` |
| **Create** | Server assigns ID | Client or server assigns ID at fixed URI |
| **Repeat call** | May create **duplicate** unless idempotency key | Same body → same final state |
| **Example** | “Charge this order” (new payment attempt) | “Upsert config at this key” |

### PUT vs PATCH

| | **PUT** | **PATCH** |
|---|---------|-----------|
| **Body** | **Full** representation required | **Partial** fields only |
| **Missing fields** | Often interpreted as null/cleared | Unmentioned fields unchanged |
| **Idempotency** | Yes by definition | Yes if patches are absolute (not `{ "increment": 1 }`) |
| **Interview tip** | “Replace document” | “Merge patch” — JSON Merge Patch (RFC 7396) or JSON Patch (RFC 6902) |

**Interview line:** “POST for creates and non-idempotent actions; PUT for full replace at a known URI; PATCH for partial updates — document which fields are required.”

---

## 3. Idempotency & safe retries

**Idempotent operation:** calling it **once or many times** has the **same effect** as calling it once (for server state that matters).

| Layer | Mechanism |
|-------|-----------|
| **HTTP** | GET, PUT, DELETE idempotent by spec; POST is not |
| **Application** | `Idempotency-Key` header on POST |
| **Database** | Unique constraint on business key; upsert |
| **Distributed** | Dedup table: `(tenant_id, idempotency_key) → response` |

### Retry rules

| Retry? | When |
|--------|------|
| **Yes** | 408, 429 (with backoff), 502, 503, 504 — **transient** |
| **No** | 400, 401, 403, 404, 409, 422 — **client/logic** error |
| **Careful** | 500 — unknown if side effect happened → use idempotency key before retry |

### Exponential backoff + jitter

```text
delay = min(cap, base * 2^attempt) + random_jitter
```

**Interview line:** “Retries without idempotency double-charge. I retry only transient failures with backoff+jitter and the same idempotency key.”

---

## 4. Designing an idempotent payment API

### Create payment (POST)

```http
POST /v1/payments
Authorization: Bearer sk_live_...
Idempotency-Key: 550e8400-e29b-41d4-a716-446655440000
Content-Type: application/json

{
  "amount": 50000,
  "currency": "INR",
  "order_id": "order_abc",
  "customer_id": "cust_xyz"
}
```

### Server flow

```text
1. Validate auth + body (amount > 0, currency supported)
2. Lookup idempotency_key in store (Redis/DB, TTL 24–72h)
   → if hit: return stored response (same status + body)
3. Begin business logic:
   a. Check order not already paid (unique on order_id)
   b. Call payment provider (Strategy — see PaymentGateway-go)
   c. Persist Payment row + idempotency record in same transaction
4. Return 201 { payment_id, status, ... }
```

### States (state machine)

```text
created → processing → succeeded | failed
                    ↘ requires_action (3DS)
```

| Concern | Design |
|---------|--------|
| **Double submit** | Idempotency-Key + unique `order_id` |
| **Provider timeout** | Mark `processing`; webhook or poll completes; retry with **same key** |
| **Webhook replay** | Store `event_id`; ignore duplicate deliveries |
| **Partial failure** | Core write in DB txn; provider call outside or saga with compensating action |

**Repo tie-in:** [../Go/PaymentGateway-go/](../Go/PaymentGateway-go/) — pluggable `BankGateway`; LLD idempotency in [../README.md](../README.md) §12.

**Interview line:** “Idempotency key is the contract with the client; unique business keys are the contract inside the domain; webhooks get their own dedup.”

---

## 5. API versioning

| Strategy | Pros | Cons |
|----------|------|------|
| **URL path** `/v1/`, `/v2/` | Obvious, easy routing | URL churn |
| **Header** `Accept: application/vnd.company.v2+json` | Clean URLs | Harder to test in browser |
| **Query** `?version=2` | Rare | Cache-unfriendly |

**Staff default:** **URL path** for public developer APIs; **header** acceptable for internal services.

### Rules

- **Never** break v1 silently — ship v2, deprecate v1 with sunset date
- Document migration guide; support overlap period (6–12 months for public APIs)
- Version **DTOs**, not internal domain entities leaked to JSON

---

## 6. Pagination — offset vs cursor

| | **Offset** `?offset=100&limit=20` | **Cursor** `?after=pay_xyz&limit=20` |
|---|-----------------------------------|--------------------------------------|
| **How** | `SKIP offset LIMIT limit` | `WHERE id > cursor ORDER BY id LIMIT n` |
| **Pros** | Jump to page N; simple | Stable under concurrent inserts/deletes |
| **Cons** | Slow at large offset; **duplicate/skip** rows if data changes | No random page jump; cursor opaque |
| **Use when** | Small lists, admin UI | **High-volume lists**, feeds, webhooks, payments |

### Cursor design

```json
{
  "data": [ ... ],
  "has_more": true,
  "next_cursor": "pay_9f3a..."
}
```

- Cursor = opaque base64 of `(sort_key, id)` or signed blob
- **Always** index `(tenant_id, created_at, id)` for stable sort

**Interview line:** “Offset is fine for page 1–10; cursor for production lists at scale and payment history.”

---

## 7. Rate limiting

| Algorithm | Behaviour | LLD in repo |
|-----------|-----------|-------------|
| **Token bucket** | Allows bursts | [../JavaScript/RateLimiter2/](../JavaScript/RateLimiter2/) |
| **Fixed window** | Simple; boundary burst | Same |
| **Sliding window** | Smoother | Same |

### HTTP response

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 60
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1710000000

{ "code": "RATE_LIMITED", "message": "Too many requests", "trace_id": "req_abc" }
```

**Layers:** edge (API gateway / CDN) → service middleware → per-tenant DB quota.

**Distributed:** Redis `INCR` + TTL or token bucket in Redis; clock skew awareness.

---

## 8. Authentication vs authorization

| | **Authentication (AuthN)** | **Authorization (Authz)** |
|---|---------------------------|---------------------------|
| **Question** | Who are you? | What may you do? |
| **Mechanisms** | API key, JWT, mTLS, session cookie | RBAC, scopes, ABAC, resource ownership |
| **Example** | Verify `Bearer` token signature | Check `payments:write` scope or `tenant_id` match |

**Interview line:** “AuthN at the edge; Authz in the service layer against the resource being touched.”

See also [../JavaScript/README.md](../JavaScript/README.md) §26 (Session vs JWT).

---

## 9. JWT & OAuth (interview depth)

### JWT (JSON Web Token)

| Part | Content |
|------|---------|
| **Header** | alg (`RS256`), typ |
| **Payload** | `sub`, `exp`, `iss`, custom claims / scopes |
| **Signature** | Verify with public key (RS256) or secret (HS256 — avoid for multi-service) |

| Pros | Cons |
|------|------|
| Stateless verification | Hard to revoke before `exp` |
| Horizontal scale | Token size grows with claims |

**Production pattern:** short-lived access JWT (15 min) + refresh token (rotating, stored server-side or in secure cookie).

### OAuth 2.0 (roles)

| Role | Role |
|------|------|
| **Resource owner** | User |
| **Client** | App requesting access |
| **Authorization server** | Issues tokens |
| **Resource server** | Your API |

| Flow | Use |
|------|-----|
| **Authorization Code + PKCE** | User-facing apps (SPAs, mobile) |
| **Client Credentials** | Machine-to-machine, merchant server → payment API |
| **Refresh token** | Renew access without re-login |

**Interview line:** “OAuth is **delegation** — client gets scoped access without sharing user password; Client Credentials for server-side integrations.”

---

## 10. API gateway

**API gateway** = single entry point in front of microservices.

| Responsibility | Examples |
|----------------|----------|
| **Routing** | `/v1/payments` → payment-service |
| **TLS termination** | HTTPS at edge |
| **AuthN** | Validate API key / JWT |
| **Rate limiting** | Per key, per IP, per tenant |
| **Request/response transform** | Header injection, legacy adapter |
| **Observability** | Access logs, trace_id propagation |

**Products:** Kong, AWS API Gateway, Envoy, Nginx, cloud provider gateways.

**vs Load balancer:** LB = L4/L7 distribute; gateway = **API policy** (auth, quota, routing rules).

**Interview line:** “Gateway enforces cross-cutting policies so services stay thin; business authz still lives in the service.”

---

## 11. Request validation

**Validate at the boundary** (handler/middleware) before domain logic.

| Layer | What |
|-------|------|
| **Syntax** | JSON parseable, required fields present |
| **Semantic** | amount > 0, enum values, string lengths |
| **Business** | order exists, currency matches — often in service layer |

```json
// 422 Unprocessable Entity
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid request",
  "errors": [
    { "field": "amount", "code": "must_be_positive" }
  ],
  "trace_id": "req_abc"
}
```

**Node/TS:** zod, joi, class-validator (Nest). **Go:** `go-playground/validator`, manual checks, or codegen from OpenAPI.

**Interview line:** “Fail fast at the edge with structured field errors; never trust client JSON.”

---

## 12. Error response design & HTTP status codes

### Consistent error shape

```json
{
  "code": "PAYMENT_ALREADY_CAPTURED",
  "message": "This payment has already been captured",
  "trace_id": "req_abc123",
  "doc_url": "https://docs.example.com/errors#payment_already_captured"
}
```

| Field | Purpose |
|-------|---------|
| `code` | Machine-readable, stable across versions |
| `message` | Human-readable (not for branching logic) |
| `trace_id` | Support / log correlation |
| `doc_url` | Developer platform — link to fix |

### Status code cheat sheet

| Code | Meaning | When to use |
|------|---------|-------------|
| **200** | OK | GET/PUT/PATCH success |
| **201** | Created | POST created resource |
| **202** | Accepted | Async job queued (`{ "job_id" }`) |
| **204** | No Content | DELETE success |
| **400** | Bad Request | Malformed JSON, wrong types |
| **401** | Unauthorized | Missing/invalid credentials |
| **403** | Forbidden | Valid auth, insufficient permission |
| **404** | Not Found | Resource doesn't exist (or hidden) |
| **409** | Conflict | Duplicate, state conflict |
| **422** | Unprocessable | Valid JSON, semantic validation failed |
| **429** | Too Many Requests | Rate limited |
| **500** | Internal Error | Unhandled bug — log + alert |
| **502/503/504** | Gateway/Unavailable/Timeout | Upstream/down/overloaded — **retry candidate** |

**Interview line:** “401 vs 403 is the classic trap — unauthenticated vs authenticated but not allowed.”

---

## 13. Backward compatibility & evolving APIs

### Safe changes (non-breaking)

- Add **optional** JSON fields
- Add new endpoints
- Add new enum values **only if** clients ignore unknown values
- Add new query params with defaults

### Breaking changes (require new version)

- Remove or rename fields
- Change field type (`string` → `number`)
- Change URL path or auth scheme
- Tighten validation on existing fields

### Evolution playbook

```text
1. Add new field/endpoint in v1 (optional) — document
2. Deprecation header: Sunset: Sat, 01 Jan 2027 00:00:00 GMT
3. Ship v2 with breaking change
4. Migration guide + SDK update
5. Monitor v1 traffic; turn down after sunset
```

**Interview line:** “Be additive by default; version when you must break; never change semantics under the same field name silently.”

---

## 14. APIs consumed by millions of developers

Payment/fintech platform APIs (Razorpay-style) need **product + engineering** discipline:

| Pillar | What to mention |
|--------|-----------------|
| **Documentation** | OpenAPI spec, quickstart, idempotency guide, error catalog |
| **SDKs** | Official Java, Node, Python — generated from OpenAPI where possible |
| **Sandbox** | Test keys, fake payments, webhook simulator |
| **Consistency** | Same error shape, pagination, auth across all resources |
| **Versioning** | Clear v1/v2 policy, long deprecation windows |
| **Reliability** | SLOs, status page, idempotent webhooks |
| **Security** | Key rotation, IP allowlist, HMAC webhook signatures |
| **Rate limits** | Documented tiers; 429 + Retry-After |
| **Supportability** | `trace_id` in every response; request ID in logs |

### Webhook design (brief)

```http
POST https://merchant.com/webhooks
X-Webhook-Signature: sha256=...
X-Webhook-Id: evt_unique_123

{ "event": "payment.captured", "payload": { ... } }
```

- Verify signature; dedup on `event_id`; respond **200** quickly; process async if heavy.

**Interview line:** “Developer experience is part of the API — predictable errors, sandbox, SDKs, and never breaking silently.”

---

## 15. Practice checklist

Answer **without notes**:

### HTTP & REST
- [ ] State REST constraints in one minute
- [ ] POST vs PUT vs PATCH with idempotency column
- [ ] When 401 vs 403 vs 409 vs 422

### Payments & reliability
- [ ] Draw idempotent payment flow with Idempotency-Key
- [ ] What happens on provider timeout + client retry
- [ ] Webhook deduplication strategy

### Scale & platform
- [ ] Cursor vs offset — trade-offs
- [ ] Rate limit response headers
- [ ] How to evolve API without breaking mobile clients
- [ ] OAuth Client Credentials vs Auth Code + PKCE

### Cross-stack
- [ ] JWT vs session ([JS §26](../JavaScript/README.md))
- [ ] Go graceful shutdown + worker pool ([Go §12, §15](../Go/README.md))
- [ ] Rate limit algorithms ([RateLimiter2 design](../JavaScript/RateLimiter2/README.md))

---

## 16. Rapid revision cheat sheet

```
REST                  → resources (nouns), stateless, HTTP semantics, cacheable GET
POST vs PUT           → POST create/action (not idempotent); PUT replace at URI (idempotent)
PUT vs PATCH          → PUT full body; PATCH partial; document required fields
Idempotency           → HTTP: GET/PUT/DELETE; app: Idempotency-Key on POST
Retries               → transient 5xx/429 only; same key; exponential backoff + jitter
Payment API           → validate → key lookup → txn → provider → persist → 201
                      → unique order_id + idempotency record + webhook event_id dedup
Versioning            → /v1/ path; additive changes; Sunset header; v2 for breaks
Pagination            → offset simple/slow; cursor stable at scale + index (tenant, ts, id)
Rate limit            → 429 + Retry-After; token bucket; edge + service; Redis distributed
AuthN vs Authz        → who vs what allowed; JWT verify at edge, scopes in service
JWT                   → short access + refresh; RS256; revoke via blocklist/TTL
OAuth                 → delegation; Client Credentials M2M; Auth Code+PKCE for users
API gateway           → TLS, auth, rate limit, route, trace_id; not business logic
Validation            → boundary fail-fast; 422 + field errors; zod/validator
Errors                → code + message + trace_id; stable machine codes
Status codes          → 201 create · 202 async · 409 conflict · 429 limit · 502 retry
Backward compat       → add optional fields; never rename/remove in same version
Platform API          → docs, sandbox, SDKs, webhooks+HMAC, SLO, status page
```

---

*Aligned with [../README.md](../README.md) §10–§12, [../Go/PaymentGateway-go/](../Go/PaymentGateway-go/), [../JavaScript/RateLimiter2/](../JavaScript/RateLimiter2/).*
