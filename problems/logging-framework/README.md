# Logging framework — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Levels DEBUG–ERROR?
2. Multiple appenders?
3. Async vs sync?
4. Structured JSON format?
5. Per-logger config?
6. MDC / correlation id?
7. Log sampling / rate limit?
8. Rotation for file appender?

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

## Step 2 — Entities & classes

```text
LogLevel: DEBUG < INFO < WARN < ERROR

LogRecord { timestamp, level, message, loggerName, fields? }

Formatter (interface) → TextFormatter | JSONFormatter
  - Format(record) → string

Appender (interface) → ConsoleAppender | FileAppender
  - Append(record)

Logger
  - level threshold
  - appenders []
  - Debug/Info/Warn/Error(msg, fields...)
  - log(level, msg)  // filter → format → each appender

LoggerFactory / Config
```

**Patterns:** Strategy (formatter) · Composite / fan-out appenders · optional Chain of Responsibility for filters

## Step 3 — Flows

**Happy path**
1. App calls `logger.Info("paid", {orderId})`  
2. Level below threshold → return early  
3. Build LogRecord → Formatter  
4. Fan-out to each Appender  

**Edge cases**
1. File appender disk full → isolate failure; don’t crash app (policy)  
2. Async mode: enqueue; overflow drop or block — state aloud

## Step 4 — APIs

```text
Logger.Debug(msg) / Info / Warn / Error
Logger.With(fields) → child logger
Configure(level, appenders[], formatter)
Flush() / Close()     // shutdown
```

## Step 5 — Deepen

- Sync v1 is simple; async queue + worker so logging doesn’t block hot path  
- Flush / Close on shutdown so buffered logs aren’t lost  
- Thread-safe appenders (file write lock)  
- Avoid Singleton — inject logger for testability  
- Bound queue size; never let logging OOM the process

## Step 6 — Evolve

- Chain of Responsibility for filters (PII redact, sampling)  
- Remote appender (HTTP/syslog); hierarchical loggers  
- Related: [metrics-collector](../metrics-collector/README.md), [webhook-delivery](../webhook-delivery/README.md)
