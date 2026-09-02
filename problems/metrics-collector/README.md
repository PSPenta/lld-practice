# Metrics collector — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · [Method §5](../../README.md#5-the-standard-approach-memorize-this)  
> **Solved in repo:** ❌

## Step 1 — Clarify

### Questions (ask 6–8)
1. Counter vs gauge vs histogram?
2. Push or pull scrape?
3. Label cardinality limits?
4. Aggregation window?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Instrumented services, Registry, Exporter |
| **Use cases (v1)** | 1. Record metric 2. Export snapshot for scrape |
| **In scope** | Register metric, increment/observe, export text |
| **Out of scope** | Long-term TSDB storage |
| **Assumptions** | In-process registry; Prometheus-style export |

### Confirm understanding
> "Services record counters/gauges; scraper pulls a text snapshot."

---

## Step 2 — Entities & classes

`Metric`, `Counter`, `Gauge`, `Histogram`, `Registry`, `Exporter`

---

## Step 3 — Flows

Record sample → aggregate in registry → scrape/export snapshot

---

## Step 4 — APIs

`Increment(name, labels)`, `Observe(name, value)`, `Export()` Prometheus text

---

## Step 5 — Deepen (concurrency, failure, idempotency)

Thread-safe registry; bound label cardinality

---

## Step 6 — Evolve

Remote write; sampling for high-cardinality

---

