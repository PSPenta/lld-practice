# Logging framework — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Levels DEBUG–ERROR?
2. Multiple appenders?
3. Async vs sync?
4. Structured JSON format?
5. Per-logger config?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Application code, Logger, Appenders |
| **Use cases (v1)** | 1. Log at level 2. Fan-out to console/file 3. Filter below threshold |
| **In scope** | Logger, level filter, appender chain |
| **Out of scope** | Remote log aggregation |
| **Assumptions** | Sync v1; console + file appenders |

### Confirm understanding
> "Code logs at a level; framework formats and writes to configured appenders."

---

## Step 2 — Entities & classes

`Logger`, `LogLevel`, `Appender`, `Formatter`, `LogRecord`

---

## Step 3 — Flows

Log message → level filter → format → fan-out to each appender

---

## Step 4 — APIs

`Debug(msg)`, `Info`, `Warn`, `Error` · configure appenders at startup

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Async queue + worker so logging doesn't block hot path; flush on shutdown

---

## Step 6 — Evolve

Chain of Responsibility for filters; avoid Singleton — inject logger

---

