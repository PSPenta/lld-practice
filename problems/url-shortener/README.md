# URL shortener — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ✅

## Code in this repo

| Language | Path | Notes |
|----------|------|--------|
| **JavaScript** | [`JavaScript/UrlShortener/`](../../JavaScript/UrlShortener/) | Express + in-memory map |
| **Go** | [`Go/UrlShortener-go/`](../../Go/UrlShortener-go/) | |

Scale at 100M QPS → HLD (not LLD focus)

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Custom aliases allowed?
2. Expiry TTL?
3. Click analytics?
4. Base62 counter vs hash?
5. Collision handling?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Client, `UrlService`, redirect handler |
| **Use cases (v1)** | Shorten long URL · redirect by code |
| **In scope** | Encode, store mapping, HTTP redirect |
| **Out of scope** | Global scale, CDN |
| **Assumptions** | Auto-generated codes; in-memory store |

### Confirm understanding
> "POST long URL → short code; GET /{code} → 302 to original."

---

## Step 2 — Entities & classes

```text
UrlService
  - repository: Map<shortCode, longUrl>
  - encoder: ShortCodeGenerator (base62 counter or hash)

createShortUrl(longUrl) → code
resolve(code) → longUrl
```

---

## Step 3 — Flows

**Shorten:** validate URL → generate unique code → store → return code  

**Redirect:** lookup code → if missing 404 → else 302 to long URL

---

## Step 4 — APIs

```http
POST /shorten   body: { "url": "https://..." }
GET  /{code}    → 302 Location: long url
```

---

## Step 5 — Deepen

- Collision on hash → retry or increment
- Idempotent shorten for same URL (optional same code)

---

## Step 6 — Evolve

- DB + unique index on code; cache hot codes; HLD for read scale
