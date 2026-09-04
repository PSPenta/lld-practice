# Metrics collector — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Counter vs gauge vs histogram?
2. Push or pull scrape?
3. Label cardinality limits?
4. Aggregation window?
5. Exemplars / exemplars for traces?
6. Multi-process or single process registry?
7. Reset on scrape (rare) vs cumulative?
8. Naming conventions / units?

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

## Step 2 — Entities & classes

```text
Labels map[string]string

Metric (interface)
  - Name(), Type(), Collect() → samples

Counter { name, labels → atomic value; Inc / Add }
Gauge   { Set / Inc / Dec }
Histogram {
  buckets[], counts[], sum
  - Observe(value)
}

Registry
  - Register(metric)
  - Counter(name, labels) / Gauge / Histogram
  - Gather() → []Family

Exporter
  - Export() → Prometheus text / JSON
```

**Patterns:** Registry singleton-per-process (injectable) · Strategy per metric type

## Step 3 — Flows

**Happy path**
1. App registers or lazily gets Counter/Histogram  
2. Hot path calls `Inc` / `Observe`  
3. Scrape hits Export → Registry.Gather → text snapshot  

**Edge cases**
1. Unbounded label values (userId) → cardinality explosion — reject or bound  
2. Concurrent Inc during Gather → atomic reads; consistent-enough snapshot for v1

## Step 4 — APIs

```text
Increment(name, labels...)
SetGauge(name, value, labels...)
Observe(name, value, labels...)
Export() → string   // Prometheus text exposition
```

```http
GET /metrics
```

## Step 5 — Deepen

- Thread-safe registry and per-series atomics  
- Bound label cardinality; deny high-cardinality keys  
- Gather must not block writers for long — copy or lock briefly  
- Histograms: fixed buckets; observe is O(buckets)  
- Fail closed on register duplicate name+type mismatch

## Step 6 — Evolve

- Remote write; sampling for high-cardinality  
- Exemplars linking to traces  
- Related: [logging-framework](../logging-framework/README.md), [api-gateway](../api-gateway/README.md) (instrument filters)
