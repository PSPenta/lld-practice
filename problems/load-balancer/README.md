# Load balancer — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. L4 or L7?
2. Algorithm: RR, least conn, weighted?
3. Health checks?
4. Sticky sessions?
5. How many backends?
6. Passive vs active health probes?
7. Connection draining on deregister?
8. TLS terminate here?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Clients, LoadBalancer, Backend pool |
| **Use cases (v1)** | 1. Forward request to healthy backend 2. Register/deregister backend |
| **In scope** | Selection strategy, health flag, forward |
| **Out of scope** | Global anycast, autoscaling |
| **Assumptions** | Round-robin; HTTP pass-through; 3 backends |

### Confirm understanding
> "Client sends request; LB picks a healthy backend and proxies the response."

## Step 2 — Entities & classes

```text
Backend {
  id, host, weight?
  healthy bool
  activeConns int
}

SelectionStrategy (interface)
  - Pick(backends []Backend) → Backend
  RoundRobin | LeastConnections | Weighted

HealthChecker
  - Probe(backend) → healthy
  - Start(interval)

LoadBalancer
  - Register(backend) / Deregister(id)
  - Handle(request) → response
```

**Patterns:** Strategy (selection) · Health observer updating pool

## Step 3 — Flows

**Happy path**
1. Request arrives → filter to healthy backends  
2. SelectionStrategy picks one → increment activeConns  
3. Forward / proxy → return response → decrement conns  

**Edge cases**
1. No healthy backends → 503  
2. Passive failures mark unhealthy; active probe restores; deregister drains in-flight

## Step 4 — APIs

```text
Handle(request) (response, error)
RegisterBackend(host, opts?)
DeregisterBackend(id)
SetStrategy(strategy)
GetBackends() → []Backend
```

```http
GET /proxy/**          # data plane
POST /admin/backends
DELETE /admin/backends/{id}
```

## Step 5 — Deepen

- Thread-safe backend list; RR index under mutex  
- Mark unhealthy after N consecutive failures; avoid flapping with hysteresis  
- Sticky sessions via consistent hash if required  
- Timeout on backend call; don’t leave activeConns incremented on panic  
- Drain: stop new traffic, wait for in-flight before remove

## Step 6 — Evolve

- Consistent hash for sticky / cache-friendly routing  
- Active health probes; weighted least-conn  
- Related: [api-gateway](../api-gateway/README.md), [circuit-breaker](../circuit-breaker/README.md)
