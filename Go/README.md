# Golang Interview Preparation Guide
> Senior / Lead Golang Developer — Global Edition  
> Covers: Runtime, Concurrency, Memory, GC, Interfaces, Testing, HTTP, DB, Patterns, System Design

---

## Table of Contents
0. [Why Go, fmt & Interview Topic Map](#0-why-go-fmt--interview-topic-map)
1. [Go Runtime & Scheduler](#1-go-runtime--scheduler)
2. [Goroutines](#2-goroutines)
3. [Channels](#3-channels)
4. [Mutex, RWMutex, Atomics](#4-mutex-rwmutex-atomics)
5. [Context](#5-context)
6. [Memory & Garbage Collector](#6-memory--garbage-collector)
7. [Error Handling](#7-error-handling)
8. [Go Keywords — Complete Reference](#8-go-keywords--complete-reference)
9. [Interfaces & Type System](#9-interfaces--type-system)
   - [Reflection](#reflection)
10. [Data Structures — Map vs Slice](#10-data-structures--map-vs-slice)
11. [Testing](#11-testing)
12. [HTTP & Networking](#12-http--networking)
13. [Database](#13-database)
14. [Go Modules & Dependency Management](#14-go-modules--dependency-management)
15. [Design Patterns in Go](#15-design-patterns-in-go)
16. [Distributed Systems Concepts](#16-distributed-systems-concepts)
17. [LLD / System Design](#17-lld--system-design)
18. [Classic Interview Gotchas — Predict the Output](#18-classic-interview-gotchas--predict-the-output)
19. [Practice Questions](#19-practice-questions)
20. [Rapid Revision Cheat Sheet](#20-rapid-revision-cheat-sheet)

---

## 0. Why Go, fmt & Interview Topic Map

### Why Go?

Go tries to combine the **ease of dynamic languages** (simple syntax, fast iteration, garbage collection) with the **efficiency and safety of statically typed compiled languages** (compile-time checks, native performance, single binary).

| Strength | Interview one-liner |
|----------|---------------------|
| **Simple language** | Small spec, fast compile, easy onboarding |
| **Concurrency built-in** | Goroutines + channels, not bolt-on threads |
| **Static typing** | Catch bugs at compile time; great tooling |
| **GC + runtime** | No manual memory management; sub-ms GC pauses |
| **Deployment** | One static binary, cross-compile, containers |
| **Standard library** | HTTP, JSON, testing, crypto — production-ready |

Trade-offs to mention fairly: no generics until 1.18 (now has them), verbose error handling, less OOP than Java/C#, GC tuning for extreme latency workloads.

### fmt package — formatting & printing

```go
import "fmt"

name := "Go"
n := 42

fmt.Print("hello")                    // print with no newline
fmt.Println(name, n)                  // print values + newline
fmt.Printf("lang=%s count=%d\n", name, n)

// Verbs you'll use in interviews:
fmt.Printf("%v\n", n)    // default format → 42
fmt.Printf("%T\n", n)    // type of variable → int
fmt.Printf("%d\n", n)    // integer
fmt.Printf("%s\n", name) // string
fmt.Printf("%+v\n", u)   // struct with field names
fmt.Printf("%#v\n", u)   // Go-syntax representation

// Sprintf — same as Printf but RETURNS a string (does not print)
msg := fmt.Sprintf("user %s has %d items", name, n)

// Error wrapping (see §7)
err := fmt.Errorf("fetch user %d: %w", id, originalErr)
```

| Function | Prints? | Returns |
|----------|---------|---------|
| `Print` / `Println` / `Printf` | Yes → stdout | — |
| `Sprint` / `Sprintln` / `Sprintf` | No | `string` |
| `Fprint` / `Fprintf` | Yes → `io.Writer` | bytes written |

### Interview topic map (your checklist → this doc)

| Topic | Where in this doc |
|-------|-------------------|
| Why Go? | **§0** (above) |
| fmt (`%v`, `%T`, `Sprintf`) | **§0** |
| Goroutines vs threads, creating/managing | **§1** GMP, **§2** |
| Channels, communication | **§3** |
| Mutex vs RWMutex, `sync` package | **§4** |
| `select`, coordinating goroutines | **§3**, **§8** |
| Concurrency pitfalls, safe shared data | **§4**, **§18** |
| Error handling, panic/recover | **§7**, **§8** |
| Interfaces, polymorphism | **§9** |
| Garbage collection, memory leaks | **§6** |
| Reflection | **§9 → Reflection** |
| Error values vs error types | **§7** |
| Unit tests, mocking | **§11** |
| Go Modules vs GOPATH | **§14** |
| Performance optimization, profiler | **§6** |
| `context` — timeouts, cancellation | **§5** |
| Error wrapping | **§7** |
| Middleware & middleware chains | **§12** |

### How to revise with this doc (interview-ready path)

| Day | Focus | Sections |
|-----|-------|----------|
| **1 — Core concurrency** | GMP, goroutines, channels, sync, context | §1–§5, then **§18** gotchas #1–2, #8, #11–12 |
| **2 — Language + memory** | Errors, interfaces, GC, maps/slices, keywords | §6–§10, **§18** rest of gotchas |
| **3 — Production + drill** | HTTP, DB, modules, patterns; practice aloud | §11–§17, **§19** Qs, **§20** cheat sheet |

**Before the interview:** run through **§18** (predict output without looking), skim **§20** (30 min), pick 5 **§19** questions and answer out loud.

**What this doc covers well:** concurrency (GMP, channels, mutex, atomics, context), error handling, GC/profiling, interfaces, gotchas, HTTP/DB basics, modules, middleware.

**Pair with hands-on:** run `go test -race` on a small snippet, trace one goroutine leak fix, write a table-driven test — interviewers often follow theory with “what’s wrong with this code?”

---

## 1. Go Runtime & Scheduler

### GMP Model
- **G** — goroutine (user-space thread, ~2KB initial stack)
- **M** — OS thread (managed by the OS, ~1–8MB stack)
- **P** — logical processor (holds a local run queue of goroutines, max = `GOMAXPROCS`)

```
P1 [G1 G2 G3 ...]  →  M1 (OS thread running on CPU core)
P2 [G4 G5 G6 ...]  →  M2 (OS thread running on CPU core)
```

### Scheduling
- Go uses **M:N scheduling** — many goroutines multiplexed onto few OS threads
- When a goroutine blocks, it **leaves** P's runnable queue (see below); M usually keeps running other Gs from that P
- **Work stealing** — idle P steals goroutines from another P's local queue
- **Preemption** — since Go 1.14, goroutines are preempted at safe points even in CPU-bound loops (no longer cooperative-only)

### When G blocks — lifecycle (common interview question)

**P does not put a blocked G back on its run queue while it is waiting.** The G moves to a wait structure; when unblocked, the runtime marks it **Runnable** and enqueues it again (local or global queue).

| Block reason | Where G waits | What happens to M |
|--------------|---------------|-------------------|
| **Network I/O** (netpoller) | Netpoller wait set | Same M immediately runs the **next G** from P's queue |
| **Channel / mutex / timer** | That object's wait list | Same — M runs another G |
| **Blocking syscall** (e.g. file I/O) | G tied to syscall; not runnable | M may block in kernel → runtime **detaches P** and binds P to another M |

```text
Running on P1/M1
  → G blocks on network read
  → G parked in netpoller (NOT in P1's [G2 G3 …] queue)
  → M1 runs G2, G3, … from P1
  → socket ready → G marked Runnable → run queue → scheduled on some P/M again
```

**Syscall handoff (why total M can exceed GOMAXPROCS):**

```text
1. P1 + M7 running G1
2. G1 enters blocking syscall → M7 stuck in kernel
3. Runtime detaches P1 from M7, attaches P1 to M8 (new or from pool)
4. P1 + M8 runs G2, G3, …
5. Syscall returns on M7 → G1 runnable → queued → runs later
```

### Multi-core & P–M binding

On an **N-core** machine, default `GOMAXPROCS = N`:

```text
Core 0:  P0 [G G G] ──M0──► running G_a    ← parallel
Core 1:  P1 [G G  ] ──M1──► running G_b    ← parallel
Core 2:  P2 [G    ] ──M2──► running G_c
...
```

| Question | Answer |
|----------|--------|
| How many Gs run Go code in parallel? | At most **`GOMAXPROCS`** (one active G per P per core) |
| Multiple M on one P **at the same instant**? | **No** — 1:1 while executing user Go code |
| Multiple M serving one P **over time**? | **Yes** — when an M blocks in a syscall, P rebinds to another M |
| Who configures M count? | **Runtime only** — you do not set it; blocked syscall threads + active threads can make **total M > GOMAXPROCS** |
| What you configure | **`GOMAXPROCS`** only (= number of Ps) |

Idle P with an empty local queue **steals** Gs from another P so cores stay busy.

### GOMAXPROCS
```go
import "runtime"

runtime.GOMAXPROCS(4)          // use 4 OS threads = 4 goroutines truly parallel
runtime.GOMAXPROCS(0)          // query current value
// default: number of CPU cores on the machine
```
- Setting it higher than CPU cores does NOT improve CPU-bound performance
- For I/O-bound workloads, higher values can help

### Goroutine states
- **Runnable** — ready to run, sitting in a run queue (local to P or global), waiting for an M
- **Running** — currently executing on an M bound to a P
- **Waiting** — blocked on channel, netpoller, syscall, mutex, or timer — **not** in P's runnable queue

---

## 2. Goroutines

### Fundamentals

**Goroutine** = lightweight concurrent task scheduled by the Go runtime (not a 1:1 OS thread).

| | Goroutine | OS thread |
|---|-----------|-----------|
| Stack | ~2KB, grows dynamically | ~1–8MB fixed |
| Cost | Cheap — spawn millions | Expensive — thousands max |
| Scheduling | Go runtime (GMP) | OS kernel |
| Creation | `go func()` | `pthread` / platform API |

**Managing goroutines:** you cannot force-kill — use **`context` cancellation**, **closed channels**, or **`WaitGroup`** to coordinate lifecycle. Always give goroutines an exit path.

- Stack grows and shrinks dynamically (up to 1GB max)
- You cannot force-kill a goroutine — you signal it

### Starting goroutines
```go
go func() { doWork() }()           // fire and forget — bad, no exit path
go worker(ctx, jobs)               // correct — passes context for cancellation
```

**Classic trap:** `for i := 0; i < n; i++ { go func() { use(i) }() }` — loop variable capture (see §18). Prefer `go func(i int){...}(i)`.

### WaitGroup — wait for all to finish
```go
var wg sync.WaitGroup
for _, task := range tasks {
    wg.Add(1)
    go func(t Task) {
        defer wg.Done() // prevents wg.Wait() deadlock if panic occurs
        process(t)
    }(task)
}
wg.Wait() // blocked until EVERY goroutine calls Done
```

### `wg.Wait()` vs first-success (`Promise.race` equivalent)

| Goal | What to wait on | Early return? |
|---|---|---|
| Need **all** results / all done | `wg.Wait()` | No |
| Need **first success**, cancel rest | `<-ch` (channel receive) | Yes |
| Start goroutines and wait on nothing | — | Broken — `main` may exit and kill them |

**Key point:** skipping `WaitGroup` alone does not give first-success. You still must wait on something — usually a channel. `wg.Wait()` waits for everyone; `<-ch` wakes as soon as the first result arrives.

Node.js mapping:
- `Promise.race([...])` ≈ first receive from a shared results channel
- `AbortController.abort()` ≈ `context.WithCancel` + `cancel()`
- Note: raw `Promise.race` does **not** auto-cancel losers; same in Go — you must call `cancel()` yourself

```go
// First successful API response wins; cancel remaining calls
func fetchFromAny(parent context.Context, urls []string) ([]byte, error) {
    ctx, cancel := context.WithCancel(parent)
    defer cancel()

    type result struct {
        data []byte
        err  error
    }
    ch := make(chan result, len(urls)) // buffered: losers can send and exit

    client := &http.Client{Timeout: 5 * time.Second}
    for _, u := range urls {
        go func(url string) {
            req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
            if err != nil {
                ch <- result{err: err}
                return
            }
            resp, err := client.Do(req)
            if err != nil {
                ch <- result{err: err}
                return
            }
            defer resp.Body.Close()
            if resp.StatusCode != http.StatusOK {
                ch <- result{err: fmt.Errorf("status %d", resp.StatusCode)}
                return
            }
            body, err := io.ReadAll(resp.Body)
            ch <- result{data: body, err: err}
        }(u)
    }

    var errs []error
    for i := 0; i < len(urls); i++ {
        r := <-ch // unblocks on FIRST finished call (not after all)
        if r.err == nil && len(r.data) > 0 {
            cancel() // abort remaining HTTP requests via ctx
            return r.data, nil
        }
        errs = append(errs, r.err)
    }
    return nil, fmt.Errorf("all sources failed: %v", errs)
}
```

Two variants:
1. **First to finish** (like raw `Promise.race`) — take first settle, success or fail
2. **First successful** (usually what you want) — ignore failures until one succeeds or all fail (code above)

Do **not** put `wg.Wait()` before reading results if you want early return — that forces waiting for all.

### Error handling — one buffered channel per goroutine
```go
type Result struct {
    ID  int
    Err error
}

func runAndTrack(tasks []Task) []Task {
    errChs := make([]chan error, len(tasks))
    for i := range errChs {
        errChs[i] = make(chan error, 1) // buffered: goroutine never blocks on send
    }

    var wg sync.WaitGroup
    for i, task := range tasks {
        wg.Add(1)
        go func(id int, fn func() error, ch chan error) {
            defer wg.Done()
            defer func() {
                if r := recover(); r != nil {
                    ch <- fmt.Errorf("panic in task %d: %v", id, r)
                }
            }()
            ch <- fn()
        }(task.ID, task.Fn, errChs[i])
    }

    wg.Wait() // here we intentionally wait for ALL — then retry failures

    var failed []Task
    for i, ch := range errChs {
        close(ch)
        if err := <-ch; err != nil {
            failed = append(failed, tasks[i])
        }
    }
    return failed
}
```

### Retry only failed goroutines
```go
func runWithRetry(tasks []Task, maxRetries int) {
    remaining := tasks
    for attempt := 1; attempt <= maxRetries && len(remaining) > 0; attempt++ {
        remaining = runAndTrack(remaining)
        time.Sleep(time.Duration(attempt) * time.Second) // exponential backoff
    }
}
```

### Auto-restart goroutine after panic
```go
func runWithRestart(ctx context.Context, name string, fn func(context.Context) error) {
    go func() {
        for {
            func() {
                defer func() {
                    if r := recover(); r != nil {
                        log.Printf("[%s] panicked: %v", name, r)
                    }
                }()
                if err := fn(ctx); err != nil {
                    log.Printf("[%s] error: %v", name, err)
                }
            }()
            select {
            case <-ctx.Done():
                return
            case <-time.After(time.Second): // backoff before restart
            }
        }
    }()
}
```

### Goroutine leak — causes and fixes

| Cause | Fix |
|---|---|
| Blocked on channel send, reader gone | `make(chan T, 1)` buffered |
| Blocked on channel receive, sender gone | `close(ch)` to unblock |
| No exit condition | `<-ctx.Done()` + `return` |
| `time.After` in loop | `time.NewTicker` + `defer Stop()` |
| HTTP body not closed | `defer resp.Body.Close()` |
| DB rows not closed | `defer rows.Close()` |

### Detecting leaks
```go
// monitor goroutine count
log.Printf("goroutines: %d", runtime.NumGoroutine())

// expose pprof
import _ "net/http/pprof"
go http.ListenAndServe(":6060", nil)
// curl http://localhost:6060/debug/pprof/goroutine?debug=2

// in tests
import "go.uber.org/goleak"
defer goleak.VerifyNone(t)
```

---

## 3. Channels

**Channels** = typed conduits for **communication between goroutines** (CSP model: *"Don't communicate by sharing memory; share memory by communicating."*).

- **Send** `ch <- v` — pass a value (ownership transfer)
- **Receive** `v := <-ch` — wait for a value
- **Close** `close(ch)` — signal no more sends; receivers drain and exit

Use channels to pipeline work, fan-out/fan-in, signal completion, and limit concurrency (semaphore).

### Buffered vs Unbuffered
- **Unbuffered** `make(chan T)` — send blocks until receiver reads (synchronous)
- **Buffered** `make(chan T, n)` — send blocks only when buffer is full (decoupled)
- Use buffered size 1 when goroutine sends exactly once and must exit immediately
- **Classic trap:** unbuffered `ch <- x` then `<-ch` in the **same** goroutine → **deadlock** (see §18)

### Core patterns

```go
// Producer-consumer
jobs := make(chan Job, 10)
go func() {
    defer close(jobs) // signals consumer: no more items
    for _, j := range allJobs { jobs <- j }
}()
for job := range jobs { process(job) } // exits when channel closed + drained

// Fan-out: distribute work to N workers
for i := 0; i < numWorkers; i++ {
    go worker(ctx, jobs)
}

// Fan-in: collect from multiple goroutines (wait for all)
results := make(chan Result, len(tasks))
// each goroutine: results <- result
wg.Wait()
close(results)
for r := range results { collect(r) }

// First-success race (Promise.race + cancel losers) — NO wg.Wait() as gate
// wait on <-ch; first success → cancel(); return early
// see "wg.Wait() vs first-success" under Goroutines

// Semaphore: limit concurrency to N
sem := make(chan struct{}, N)
go func() {
    sem <- struct{}{}  // acquire
    defer func() { <-sem }() // release
    doWork()
}()

// Broadcast stop to all goroutines
done := make(chan struct{})
close(done) // every <-done unblocks simultaneously
```

### select — multiplex channels

**`select`** waits on **multiple channel operations** — whichever is ready first wins. Essential for timeouts, cancellation, and merging event sources.

```go
select {
case msg := <-jobs:
    process(msg)
case err := <-errs:
    handle(err)
case <-ctx.Done():
    return ctx.Err()
case <-time.After(5 * time.Second): // only outside loops — use ticker inside loops
    return errTimeout
default:
    // non-blocking branch — optional
}
```

Without `select`, you'd block on one channel only — hard to combine work + timeout + cancel.

### Channels vs Mutex
| | Channel | Mutex |
|---|---|---|
| Purpose | transfer data / ownership | protect shared memory |
| Analogy | SQS — read removes data | lock on a variable |
| Use when | ownership passes between goroutines | multiple goroutines share same variable |

---

## 4. Mutex, RWMutex, Atomics

### sync package — overview

The standard library `sync` package provides **low-level synchronisation primitives**. Use them when channels are awkward (shared memory, counters, one-time init).

| Primitive | Purpose |
|-----------|---------|
| `sync.Mutex` | Exclusive lock — one goroutine in critical section |
| `sync.RWMutex` | Many readers **or** one writer |
| `sync.WaitGroup` | Wait until N goroutines finish |
| `sync.Once` | Run a function exactly once (lazy init) |
| `sync.Pool` | Reuse objects; reduce GC pressure |
| `sync.Cond` | Wait for a condition (rare; prefer channels) |
| `sync/atomic` | Lock-free ops on single words (counters, flags) |

**Mutex vs RWMutex — when to use which:**

| Use **Mutex** | Use **RWMutex** |
|---------------|-----------------|
| Reads and writes both mutate state | **Many reads, rare writes** (config cache, registry) |
| Critical section is short | Read-heavy path must not block other readers |
| Simpler mental model | Writer excludes everyone; readers share lock |

If writes are frequent, RWMutex offers little benefit — writers still block everyone.

### Concurrency pitfalls & safe shared access

| Pitfall | Symptom | Fix |
|---------|---------|-----|
| **Data race** | Wrong counts, flaky tests | Mutex, atomic, or channels; **`go test -race`** |
| **Goroutine leak** | Memory grows; `NumGoroutine` rises | Exit on `ctx.Done()`; buffered err ch; close channels |
| **Deadlock** | All goroutines asleep | Unbuffered send/receive mismatch; lock ordering |
| **Copying mutex** | `go vet` warning; broken exclusion | Pointer receivers; don't copy structs with locks |
| **Concurrent map write** | Runtime fatal | `Mutex` + map or `sync.Map` |
| **Lock during I/O** | Throughput collapse | Release lock before network/DB calls |
| **`WaitGroup` misuse** | Early `Wait`, hang | `Add` before `go`; `defer Done()` in goroutine |

**How to ensure safe concurrent access:**

1. **Prefer message passing** — send ownership via channels when natural (§3).
2. **Protect shared memory** — `Mutex` / `RWMutex` for maps, slices, struct fields.
3. **Single-word updates** — `atomic` for metrics counters (§4 atomics).
4. **Immutability** — copy-on-write, read-only after publish (no lock on read path).
5. **Detect races in CI** — `go test -race ./...` on every PR.
6. **Design for cancellation** — `context.Context` so blocked goroutines can exit.

See **§18** for predict-the-output traps (loop variable, channel deadlock, etc.).

### sync.Mutex
```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
// only one goroutine executes this section at a time
```

### sync.RWMutex
```go
var rw sync.RWMutex

// multiple readers simultaneously
rw.RLock()
defer rw.RUnlock()

// exclusive write — blocks all readers and writers
rw.Lock()
defer rw.Unlock()
```

| | Mutex | RWMutex |
|---|---|---|
| `Lock()` | blocks everyone | blocks all |
| `RLock()` | N/A | allows concurrent readers |
| Best for | mixed read/write | many reads, rare writes |

**3 RLocks + 1 Lock scenario:**
> Writer calls `Lock()` → waits for all 3 RLocks to release → acquires exclusive access → new RLocks are blocked until writer unlocks.

**Rules:**
- Never hold a lock during I/O — other goroutines starve
- Always pair Lock with Unlock via `defer`

### sync/atomic — lock-free operations (multi-goroutine safe)

**Atomics are not single-threaded.** They exist precisely for **many goroutines** updating the **same memory location** at once — without a mutex.

#### The problem atomics solve

A normal read-modify-write is **not** one indivisible step:

```go
var counter int64

// WRONG — data race; two goroutines can interleave:
counter++   // actually: LOAD → ADD 1 → STORE (3 CPU steps)
```

```text
G1: LOAD counter (100)
G2: LOAD counter (100)    ← both read same value
G1: STORE 101
G2: STORE 101             ← lost update; expected 102, got 101
```

`atomic.AddInt64(&counter, 1)` makes that whole update **indivisible** at the hardware level — other goroutines cannot see a half-updated value.

#### How the CPU makes it safe (layman)

The Go runtime delegates to **CPU atomic instructions** (e.g. `LOCK`-prefixed ops on x86, LL/SC or CAS on ARM). The chip guarantees:

- One goroutine's atomic read-modify-write completes **before** another's is visible as a whole
- No torn reads/writes on supported sizes (`int32`, `int64`, pointers, etc.)

You do **not** get a Go mutex or a dedicated thread. All goroutines still run in parallel on many cores; the **memory operation itself** is serialized for that one word.

```text
Many goroutines on many cores
        │
        ▼
  atomic.AddInt64(&counter, 1)
        │
        ▼
  CPU ensures each +1 is one atomic step on &counter
        │
        ▼
  Final count = number of successful Add calls (no lost updates)
```

#### API (Go 1.19+ typed wrappers preferred)

```go
import "sync/atomic"

// Typed wrapper (preferred — harder to misuse)
var counter atomic.Int64
counter.Add(1)           // atomic increment — safe from any goroutine
counter.Load()           // atomic read
counter.Store(0)         // atomic write
counter.CompareAndSwap(old, new) // CAS: swap only if current == old

// Legacy pointer form — still common in older code
var n int64
atomic.AddInt64(&n, 1)
atomic.LoadInt64(&n)
atomic.StoreInt64(&n, 0)
atomic.CompareAndSwapInt64(&n, old, new)
```

**Compare-and-swap (CAS)** — the building block for lock-free algorithms:

```go
// Spin until we successfully increment (illustrative pattern)
for {
    old := counter.Load()
    if counter.CompareAndSwap(old, old+1) {
        break
    }
    // another goroutine changed counter between Load and CAS — retry
}
```

#### Atomic vs mutex

| | `sync/atomic` | `sync.Mutex` |
|---|---------------|--------------|
| **Protects** | One integer / pointer / flag | A **critical section** (many lines, many vars) |
| **Mechanism** | CPU atomic instruction | OS-level lock; goroutines may block/sleep |
| **Parallelism** | All cores keep running; only that one memory op is serialized | Only one goroutine in the section at a time |
| **Best for** | Counters, metrics, `done` flags, simple state | Maps, slices, multi-field invariants, I/O under lock |
| **Compound logic** | ❌ `if counter > 10 { counter++ }` still needs mutex or careful CAS loop | ✅ natural fit |

#### What atomics do NOT fix

```go
// WRONG — Load and Store are each atomic, but the pair is not:
if atomic.LoadInt64(&n) > 0 {
    atomic.AddInt64(&n, -1)  // two goroutines can both pass the if
}

// WRONG — two variables; atomic can't make them consistent together:
atomic.AddInt64(&sent, 1)
atomic.AddInt64(&failed, 1)  // need mutex if both must move as one unit

// WRONG — protecting a map/slice:
// maps and slices need mutex or channels, not atomic
```

**Rule:** atomic = **one memory location**, **one operation**. Multiple steps or multiple fields → **mutex**.

#### When to use which (interview line)

> Use **atomics** for simple counters and flags where goroutines only touch one `int64`/pointer. Use a **mutex** when you need a multi-line critical section or consistent view of several fields. Atomics are multi-goroutine safe; they are not single-threaded.

`sync/atomic` is also used internally by `sync.Map`, `sync.Once`, and parts of the runtime.

### sync.Once — run exactly once
```go
var once sync.Once
var db *sql.DB

func GetDB() *sql.DB {
    once.Do(func() { db, _ = sql.Open("postgres", dsn) })
    return db
}
```

### sync.Pool — reuse objects, reduce GC pressure
```go
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

func handle() {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset()
    defer bufPool.Put(buf) // returned to pool, not GC'd
    // use buf
}
```

---

## 5. Context

The **`context`** package carries **deadlines, cancellation signals, and request-scoped values** across API boundaries and goroutines. Every incoming HTTP request should get a context; pass it to DB calls, downstream HTTP, and worker goroutines.

**Interview uses:** timeout an external call, cancel sibling goroutines when one fails, propagate `request_id` for logging.

### What it provides
- **Cancellation** — signal goroutines to stop via `ctx.Done()` channel
- **Deadline/Timeout** — auto-cancel after duration
- **Values** — request-scoped data: request ID, auth token, trace ID

### How cancellation works
```
context.WithCancel(parent)
    └── returns (ctx, cancel)
            ctx.Done() → a channel (chan struct{})
            cancel()   → closes that channel
                         → all goroutines selecting on <-ctx.Done() unblock
                         → goroutine must return on its own (cooperative)
```

### Hierarchy — parent cancels all children
```go
root, cancel := context.WithCancel(context.Background())
defer cancel()

child1, c1 := context.WithTimeout(root, 5*time.Second)
defer c1()
child2, c2 := context.WithCancel(root)
defer c2()

cancel() // closes root.Done() → child1 and child2 also cancelled automatically
```

### Patterns
```go
// timeout on external call
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
resp, err := http.Get(url) // passes ctx implicitly via request
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)

// detect cancellation vs real error
if errors.Is(err, context.Canceled) {
    // caller cancelled — requeue, not a failure
}
if errors.Is(err, context.DeadlineExceeded) {
    // timed out — retry with backoff
}

// pass request-scoped values
type ctxKey string
ctx = context.WithValue(ctx, ctxKey("requestID"), "abc-123")
reqID := ctx.Value(ctxKey("requestID")).(string)

// detach — job must survive HTTP request cancellation
jobCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
go criticalJob(jobCtx) // independent of r.Context()
```

---

## 6. Memory & Garbage Collector

### Stack vs Heap — mental model

```
┌─────────────────────────────────────────────────────────┐
│  GOROUTINE 1            GOROUTINE 2            GOROUTINE 3
│  ┌──────────┐           ┌──────────┐           ┌──────────┐
│  │  Stack   │           │  Stack   │           │  Stack   │
│  │  ~2KB    │           │  ~2KB    │           │  ~2KB    │
│  │  grows   │           │  grows   │           │  grows   │
│  │  shrinks │           │  shrinks │           │  shrinks │
│  └──────────┘           └──────────┘           └──────────┘
│                                                            
│  ┌────────────────────────────────────────────────────┐   
│  │                     HEAP                          │   
│  │         shared across ALL goroutines              │   
│  │         managed by Garbage Collector              │   
│  └────────────────────────────────────────────────────┘   
└─────────────────────────────────────────────────────────┘
```

### Stack — how it behaves

- Each goroutine gets its **own stack**, starting at ~2KB
- Stack **grows dynamically** as functions are called (Go copies the stack to a larger space)
- Stack **shrinks** when functions return — memory immediately reclaimed, no GC involved
- Stack is **private** — no other goroutine can touch it
- Stack allocation is essentially **free** — just moving a pointer

```go
func add(a, b int) int {
    result := a + b  // result lives on the stack
    return result    // stack frame popped — result memory gone instantly
}
// no GC needed, no cost
```

**Call stack frames:**
```
main()          ← frame 3 (bottom)
  → fetchUser() ← frame 2
    → queryDB() ← frame 1 (top — currently executing)

When queryDB() returns → frame 1 popped → memory freed instantly
When fetchUser() returns → frame 2 popped → memory freed instantly
```

### Heap — how it behaves

- Shared memory pool across all goroutines
- Objects that need to **outlive a function call** go here
- GC periodically scans and frees unreachable objects
- Allocation is slower than stack — requires bookkeeping
- More heap allocations = more GC pressure = potential latency spikes

```go
func createUser() *User {
    u := &User{Name: "Alice"}  // u escapes to heap — address returned
    return u                   // caller holds reference → must survive this function
}
// u cannot be on the stack — it needs to outlive createUser()
// so the compiler puts it on the heap
// GC will collect it when no reference remains
```

### Escape analysis — who decides stack vs heap?

The **compiler** decides at build time — not at runtime. This is called escape analysis.

```bash
go build -gcflags='-m' .
# output shows:
# ./main.go:5:6: u escapes to heap     ← goes to heap
# ./main.go:9:6: x does not escape     ← stays on stack
```

**When does a variable escape to heap?**

```go
// 1. Address is returned — caller holds a pointer, value must outlive function
func newUser() *User {
    u := User{}  // escapes — address returned
    return &u
}

// 2. Stored in an interface — compiler cannot know the concrete size at compile time
var i interface{} = 42  // 42 escapes to heap (boxing)

// 3. Captured by a goroutine — goroutine may outlive the function
x := 10
go func() {
    fmt.Println(x)  // x escapes — goroutine may run after enclosing function returns
}()

// 4. Too large for stack — Go has a size limit per frame
big := [1_000_000]int{}  // large arrays escape to heap

// 5. Stored in a map or slice — container lives on heap, so elements do too
m["key"] = &user  // user escapes

// Does NOT escape — stays on stack:
func compute() int {
    x := 42       // stays on stack — never leaves this function
    return x * 2  // value returned, not pointer
}
```

### Practical impact on performance

```go
// BAD — allocates on heap every call — GC pressure
func handler(w http.ResponseWriter, r *http.Request) {
    buf := make([]byte, 4096) // may escape to heap
    // use buf
}

// GOOD — reuse from pool — bypasses GC entirely
var bufPool = sync.Pool{New: func() any { return make([]byte, 4096) }}
func handler(w http.ResponseWriter, r *http.Request) {
    buf := bufPool.Get().([]byte)
    defer bufPool.Put(buf)
    // use buf — never GC'd, recycled across requests
}

// BAD — interface boxing in hot loop — every value escapes to heap
func sum(nums []any) int { // interface{} forces heap allocation
    total := 0
    for _, n := range nums { total += n.(int) }
    return total
}

// GOOD — concrete type stays on stack
func sum(nums []int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}
```

### Memory layout in a running Go program

```
┌─────────────────────────────────────┐
│  Text segment                       │  compiled code (read-only)
├─────────────────────────────────────┤
│  Data segment                       │  global variables, constants
├─────────────────────────────────────┤
│  Heap  ↑ grows upward               │  dynamic allocations (GC managed)
│                                     │
│  (free space)                       │
│                                     │
│  Stack ↓ grows downward (per goroutine, not OS stack)
└─────────────────────────────────────┘
```

**Key numbers to remember:**
- Goroutine stack starts: ~2KB
- Goroutine stack max: 1GB (configurable)
- OS thread stack: 1–8MB (fixed)
- Stack allocation cost: ~nanoseconds
- Heap allocation cost: ~100ns + future GC cost
- GC STW pause in modern Go: sub-millisecond

### Escape analysis — inspect it

```bash
# see every escape decision
go build -gcflags='-m' .

# more verbose
go build -gcflags='-m -m' .

# on a specific file
go build -gcflags='-m' ./path/to/file.go
```

**What to look for:**
```
./server.go:12:6: user escapes to heap    ← heap allocation — GC cost
./server.go:20:6: x does not escape      ← stack — free
./server.go:35:13: inlining call to ...  ← inlined — stack, faster
```

### How GC works
- **Tri-color mark-and-sweep** running concurrently with your code
- Very short stop-the-world pauses (sub-millisecond in modern Go)
- Trade-off: low latency, slightly lower throughput vs Java's STW GC

### Default — let GC run automatically
For most services, no tuning needed. GC runs when heap doubles (GOGC=100).

### Deliberately tuning GC

```go
import "runtime/debug"

// GOGC — controls GC frequency
// default 100 = GC when heap doubles
// GOGC=200 → GC when heap grows 3x → less GC, more memory used
// GOGC=50  → GC when heap grows 1.5x → more GC, less memory
debug.SetGCPercent(200)

// GOMEMLIMIT — hard memory ceiling (Go 1.19+)
// set to ~90% of container memory limit
// prevents OOM kills in Kubernetes / Docker
debug.SetMemoryLimit(900 * 1024 * 1024) // 900 MB
// when heap approaches limit, GC ignores GOGC and runs aggressively
```

### Deliberately triggering GC
```go
import "runtime"

// after processing a large batch — temp objects allocated, now release them
func processBigFile(path string) []Result {
    data := loadEntireFile(path) // allocates hundreds of MB
    results := transform(data)
    data = nil          // release reference
    runtime.GC()        // force immediate collection instead of waiting
    return results
}

// before a latency-sensitive window (e.g. before traffic spike)
runtime.GC() // clean heap now → fewer mid-request GC pauses later
```

### Reduce GC pressure
```go
// pre-allocate slices
s := make([]Item, 0, 1000) // no reallocations during append

// reuse with sync.Pool
var pool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// avoid interface boxing in hot paths
// var i interface{} = x  ← x escapes to heap
// keep concrete types in tight loops

// use value types over pointers for small structs
// pointer → heap allocation; value → stack
```

### Profiling & performance optimization

**Workflow:** measure → profile → fix hotspot → benchmark again. Don't optimize without data.

| Tool / technique | What it finds |
|------------------|---------------|
| `go test -bench=. -benchmem` | Hot functions, allocations per op |
| `go test -race ./...` | Data races (correctness, not speed) |
| `go build -gcflags='-m'` | Escape analysis — heap vs stack |
| `net/http/pprof` + `go tool pprof` | CPU, heap, goroutine profiles at runtime |
| `go tool trace` | Timeline: goroutines, GC, syscalls |
| `GODEBUG=gctrace=1` | GC frequency and pause times |
| `GOGC` / `GOMEMLIMIT` | Tune GC vs memory (§6) |
| `sync.Pool`, preallocate slices | Reduce allocation rate |
| Avoid `interface{}` boxing in hot loops | Fewer heap allocs |

```bash
# CPU profile (30s sample while load test runs)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Heap — what's using memory
go tool pprof http://localhost:6060/debug/pprof/heap

# Inside pprof interactive mode:
# top10          — hottest functions
# list MyFunc      — line-level cost
# web              — call graph (needs graphviz)
```

**Interview line:** "I'd use benchmarks to set a baseline, pprof CPU/heap to find the top allocators or hot paths, fix those, then re-benchmark. Race detector stays on in CI."

### Avoiding memory leaks in Go

Go has GC, but **leaks still happen** when you hold references to objects that are no longer needed, or goroutines never exit.

| Leak source | Fix |
|-------------|-----|
| Goroutine blocked forever | `ctx` cancellation; close channels; timeouts |
| HTTP `resp.Body` not closed | `defer resp.Body.Close()` |
| DB `rows` not closed | `defer rows.Close()` |
| `time.After` in a loop | `time.NewTicker` + `defer Stop()` |
| Global map/cache growing unbounded | TTL, max size, eviction |
| Finalizers | Rare; don't rely on them — explicit `Close()` |

Monitor: `runtime.NumGoroutine()`, heap pprof, `goleak` in tests.

### Profiling (quick reference)
```bash
GODEBUG=gctrace=1 ./server                        # GC log in real time
go tool pprof http://localhost:6060/debug/pprof/heap          # heap profile
go tool pprof -alloc_objects http://.../heap      # allocation hotspots
go tool trace trace.out                           # goroutine/GC timeline
```

---

## 7. Error Handling

### Basics
```go
// errors are values — check explicitly
result, err := doSomething()
if err != nil {
    return fmt.Errorf("doSomething failed: %w", err) // wrap with context
}
```

### Wrapping and unwrapping
```go
// wrap — preserves original error, adds context
err = fmt.Errorf("user service: %w", originalErr)

// errors.Is — checks if any error in chain matches a target (sentinel)
if errors.Is(err, sql.ErrNoRows) { ... }

// errors.As — checks if any error in chain is of a specific type
var netErr *net.OpError
if errors.As(err, &netErr) {
    fmt.Println(netErr.Op)
}

// errors.Unwrap — get the next error in the chain
inner := errors.Unwrap(err)
```

### Sentinel errors vs custom types

```go
// sentinel — compare with errors.Is
var ErrNotFound = errors.New("not found")

// custom type — extract data with errors.As
type ValidationError struct {
    Field   string
    Message string
}
func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}
```

### Error values vs error interfaces — pros & cons

In Go, **`error` is an interface** with one method: `Error() string`. Everything below **implements** that interface.

| Approach | Pros | Cons | Use when |
|----------|------|------|----------|
| **Plain `error` value** (`errors.New`, `fmt.Errorf`) | Simple, no types | Hard to branch on details | Generic failures, wrapping |
| **Sentinel errors** (`var ErrX = errors.New(...)`) | Compare with `errors.Is` | Package API surface grows | Stable, known conditions (`ErrNotFound`) |
| **Custom error types** (structs with fields) | Rich data via `errors.As` | More boilerplate | Validation errors, HTTP status + body |
| **Wrapped errors** (`fmt.Errorf("%w")`) | Context chain preserved | Must use `Is`/`As`, not `==` | Almost always in production |

```go
// BAD — breaks wrapped errors
if err == sql.ErrNoRows { ... }

// GOOD
if errors.Is(err, sql.ErrNoRows) { ... }

var ve *ValidationError
if errors.As(err, &ve) {
    fmt.Println(ve.Field)
}
```

**Interview line:** "Return `error` interface values; use sentinels for stable identity, custom types when callers need fields, and `%w` wrapping for context — always inspect with `errors.Is` / `errors.As`."

### panic vs error
- `error` — expected, recoverable: DB not found, invalid input, timeout
- `panic` — unexpected, programming bug: nil pointer, index out of range
- Never use panic for control flow

### recover
```go
func safeExecute(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v\n%s", r, debug.Stack())
        }
    }()
    fn()
    return nil
}
```

### errgroup — structured goroutine error handling
```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(context.Background())
g.Go(func() error { return fetchUser(ctx) })
g.Go(func() error { return fetchOrder(ctx) })
if err := g.Wait(); err != nil {
    // first error; ctx cancelled for all other goroutines
}
```

---

## 8. Go Keywords — Complete Reference

### defer — guaranteed execution, even on panic

```go
// defer runs LIFO when the surrounding function returns (normally or via panic)
func example() {
    defer fmt.Println("third")  // runs last
    defer fmt.Println("second") // runs second
    defer fmt.Println("first")  // runs first
    fmt.Println("body")
}
// Output: body → first → second → third
```

**Does defer run on panic? Yes — same goroutine only.**
```go
func riskyFunc() {
    defer fmt.Println("defer ran") // ← WILL run even on panic

    panic("something went wrong")
    fmt.Println("this never runs")
}
// Output: "defer ran" then program crashes

// With recover — panic caught, program continues
func safeFunc() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
    panic("boom") // caught by defer above
}

// IMPORTANT: defer in main does NOT catch panic in another goroutine
func main() {
    defer func() { recover() }() // ← does NOT catch this:
    go func() {
        panic("goroutine panic") // crashes entire program
    }()
    // Fix: put defer+recover INSIDE the goroutine itself
}
```

**Common uses:**
```go
defer file.Close()      // always release file
defer mu.Unlock()       // always unlock — even if function panics
defer cancel()          // always release context resources
defer wg.Done()         // always signal WaitGroup
defer rows.Close()      // always release DB cursor
defer resp.Body.Close() // always release HTTP connection
```

---

### panic & recover
```go
panic("critical error")  // stops current goroutine, runs deferred funcs, unwinds stack
panic(err)               // can panic with any value

// recover — only useful inside a deferred function
defer func() {
    if r := recover(); r != nil {
        fmt.Println("caught:", r) // r is whatever was passed to panic()
    }
}()
```

---

### make vs new vs var
```go
// make — only for slice, map, channel — returns initialized value (not pointer)
s := make([]int, 0, 10)
m := make(map[string]int)
ch := make(chan int, 5)

// new — any type — returns *T pointing to zero value
p := new(int)    // *int → 0
p := new(MyStruct) // *MyStruct → zero struct

// var — zero value for any type
var i int        // 0
var s string     // ""
var p *MyStruct  // nil
var m map[string]int // nil (panics on write — use make)
```

---

### go — start a goroutine
```go
go doWork()            // fire and forget
go func() { ... }()   // anonymous goroutine
go func(x int) { ... }(val) // pass value at launch time — avoid closure capture bug
```

**Closure capture bug:**
```go
// BUG — all goroutines capture same variable i
for i := 0; i < 3; i++ {
    go func() { fmt.Println(i) }() // may print 3,3,3
}

// FIX — pass i as argument
for i := 0; i < 3; i++ {
    go func(n int) { fmt.Println(n) }(i) // prints 0,1,2
}
```

---

### select — multiplex channels
```go
select {
case v := <-ch1:         // receive from ch1
case ch2 <- x:           // send to ch2
case <-ctx.Done():       // context cancelled
case <-time.After(5*time.Second): // timeout
default:                 // non-blocking — runs immediately if no case ready
}
```

---

### range — iterate over collections
```go
// slice
for i, v := range slice { }
for _, v := range slice { }  // skip index
for i := range slice { }     // index only

// map
for k, v := range m { }
for k := range m { }         // keys only

// string — iterates runes (Unicode), not bytes
for i, ch := range "hello" { } // ch is rune (int32)

// channel — reads until closed
for v := range ch { }  // exits when ch is closed

// range with integer (Go 1.22+)
for i := range 5 { fmt.Println(i) } // 0,1,2,3,4
```

---

### init — package initialization
```go
func init() {
    // runs automatically before main()
    // one per file, multiple init() allowed
    // use for: setup, validation, registering drivers
    sql.Register("postgres", &pq.Driver{})
}
// order: imported packages init first, then current package
// avoid heavy logic in init — hard to test and debug
```

---

### type — define new types
```go
type UserID int           // new type (not alias) — methods can be added
type Handler func(http.ResponseWriter, *http.Request) // function type
type Celsius float64
type Fahrenheit float64
// Celsius and Fahrenheit are different types — cannot mix without conversion

type MyError struct {     // custom error type
    Code    int
    Message string
}
func (e *MyError) Error() string { return e.Message }
```

---

### const & iota
```go
const Pi = 3.14159

// iota — auto-increment in const block
type Status int
const (
    Pending Status = iota // 0
    Active                // 1
    Closed                // 2
)

type ByteSize float64
const (
    _           = iota // ignore first value
    KB ByteSize = 1 << (10 * iota) // 1024
    MB                              // 1048576
    GB                              // 1073741824
)
```

---

### goto, fallthrough, break, continue
```go
// fallthrough — in switch, continues to next case (rare in Go)
switch x {
case 1:
    fmt.Println("one")
    fallthrough     // also runs case 2
case 2:
    fmt.Println("two")
}

// labeled break — break out of nested loop
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 { break outer } // breaks both loops
    }
}

// labeled continue — continue outer loop from inner loop
outer:
for i := 0; i < 3; i++ {
    for j := 0; j < 3; j++ {
        if j == 1 { continue outer }
    }
}

// goto — jumps to label (avoid in production code)
goto cleanup
cleanup:
    closeDB()
```

---

### blank identifier _
```go
_, err := doSomething()    // ignore first return value
for _, v := range slice {} // ignore index
_ = expensiveCall()       // explicitly discard return value
import _ "net/http/pprof"  // import for side effects only (runs init())
```

---

## 9. Interfaces & Type System

### How interfaces work
- Satisfied **implicitly** — no `implements` keyword
- Interface value = two pointers: **type pointer** + **data pointer**
- Empty interface `any` (alias for `interface{}`) — holds any value

**Polymorphism:** functions accept an **interface type**; any concrete type with matching methods works at runtime (duck typing).

```go
type Logger interface {
    Log(msg string)
}

type ConsoleLogger struct{}
func (c ConsoleLogger) Log(msg string) { fmt.Println(msg) }

type FileLogger struct{ path string }
func (f FileLogger) Log(msg string) { /* write to file */ }

func Process(l Logger) {  // accepts ANY type with Log(string)
    l.Log("done")
}

Process(ConsoleLogger{})           // polymorphic call
Process(FileLogger{path: "a.log"})
```

### One type implementing multiple interfaces

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
type Writer interface {
    Write(p []byte) (n int, err error)
}
type Closer interface {
    Close() error
}

// Compose interfaces by embedding
type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

type File struct{ /* ... */ }

func (f *File) Read(p []byte) (int, error)  { /* ... */ }
func (f *File) Write(p []byte) (int, error) { /* ... */ }
func (f *File) Close() error                { /* ... */ }

// *File satisfies Reader, Writer, Closer, and ReadWriteCloser — no explicit declaration

func Copy(dst Writer, src Reader) (int64, error) {
    // both args accept any matching type — e.g. *File, *bytes.Buffer, net.Conn
    return io.Copy(dst, src)
}
```

Same pattern as `io` package: `io.Reader`, `io.Writer`, `io.ReadWriter`.

### Nil interface gotcha
```go
var p *MyStruct = nil
var i interface{} = p
fmt.Println(i == nil) // false — type pointer is set, data pointer is nil
// always return untyped nil for interface return types
func getUser() error { return nil } // correct
```

### Struct vs Interface
| | Struct | Interface |
|---|---|---|
| What | concrete data + methods | behavior contract |
| Testing | harder to mock | easy — swap with fake impl |
| Design | composition via embedding | dependency injection |

**Rule: accept interfaces, return concrete types.**
```go
func NewService(repo UserRepository) *Service { // interface in, struct out
    return &Service{repo: repo}
}
```

### Embedding — composition over inheritance
```go
type Animal struct{ Name string }
func (a Animal) Speak() string { return a.Name }

type Dog struct {
    Animal              // embedded — inherits methods
    Breed string
}
d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Lab"}
d.Speak() // promoted method
```

### Type assertion and type switch
```go
// type assertion
val, ok := i.(string) // ok = false if not string (safe)
val := i.(string)     // panics if not string (unsafe)

// type switch
switch v := i.(type) {
case string:  fmt.Println("string:", v)
case int:     fmt.Println("int:", v)
default:      fmt.Println("unknown")
}
```

### Generics (Go 1.18+)
```go
func Map[T, U any](s []T, f func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s { result[i] = f(v) }
    return result
}

// type constraints
type Number interface { int | int64 | float64 }
func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums { total += n }
    return total
}
```

---

### Reflection

**Reflection** = inspect or modify types and values **at runtime** via the `reflect` package.

```go
import "reflect"

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

u := User{Name: "Ada", Age: 36}
v := reflect.ValueOf(u)
t := reflect.TypeOf(u)

fmt.Println(t.Name())              // User
fmt.Println(v.FieldByName("Age"))  // 36

// Iterate struct fields (e.g. build JSON, ORM, validators)
for i := 0; i < t.NumField(); i++ {
    f := t.Field(i)
    fmt.Println(f.Name, f.Tag.Get("json"))
}
```

| Use reflection for | Avoid reflection for |
|--------------------|----------------------|
| JSON/XML encode-decode (stdlib uses it) | Hot paths — **10–100× slower** than direct calls |
| Generic containers before Go 1.18 generics | Simple dispatch — use interfaces instead |
| ORMs, validators, debug tools | Anything you can express with interfaces + codegen |

**Limitations:**
- Needs **exported** fields for reflection from other packages
- Loses compile-time type safety — errors at runtime
- Harder to read and refactor
- `reflect.Value` must be **addressable** to modify fields (`reflect.ValueOf(&x).Elem()`)

**Interview line:** "Reflection is powerful for frameworks and serialization, but I avoid it in request handlers — interfaces and generics give compile-time checks and better performance."

---

## 10. Data Structures — Map vs Slice

### Slice — full syntax reference

```go
// --- Declaration ---
var s []int                        // nil slice — safe to append, len=0
s := []int{}                       // empty slice — initialized, len=0
s := []int{1, 2, 3}               // literal
s := make([]int, 5)               // len=5, cap=5, all zeros
s := make([]int, 0, 100)          // len=0, cap=100 — pre-allocated, no realloc

// --- Length & Capacity ---
len(s)   // number of elements currently in slice
cap(s)   // total allocated capacity of backing array

// --- Read / Write ---
s[0]         // read first element
s[len(s)-1]  // read last element
s[1] = 99    // write

// --- Append ---
s = append(s, 4)           // always assign back — may reallocate backing array
s = append(s, 5, 6, 7)    // append multiple
s = append(s, other...)   // append entire other slice

// --- Slicing (shares backing array) ---
s2 := s[1:4]   // elements at index 1,2,3 — modifying s2 modifies s
s2 := s[:3]    // first 3 elements
s2 := s[2:]    // from index 2 to end

// --- Deep copy (independent) ---
dst := make([]int, len(s))
copy(dst, s)

// --- Check if element exists (no built-in — must loop) ---
func contains(s []int, target int) bool {
    for _, v := range s {
        if v == target { return true }
    }
    return false
}
// Go 1.21+: slices.Contains(s, target)

// --- Traversal ---
for i, v := range s { fmt.Println(i, v) }  // index + value
for _, v := range s { fmt.Println(v) }      // value only
for i := range s   { fmt.Println(s[i]) }   // index only

// --- Delete element at index i (order preserved) ---
s = append(s[:i], s[i+1:]...)

// --- Delete element at index i (order NOT preserved — faster) ---
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// --- Stack (LIFO) ---
stack = append(stack, v)           // push
top := stack[len(stack)-1]         // peek
stack = stack[:len(stack)-1]       // pop

// --- Queue (FIFO) ---
queue = append(queue, v)   // enqueue
front := queue[0]          // peek
queue = queue[1:]          // dequeue (use circular buffer for production)

// --- Sort ---
import "sort"
sort.Ints(s)
sort.Strings(s)
sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
```

| Operation | Time |
|---|---|
| Access by index | O(1) |
| Append (amortized) | O(1) |
| Insert / Delete at middle | O(n) |
| Search (unsorted) | O(n) |
| Search (sorted + binary) | O(log n) |

---

### Map — full syntax reference

```go
// --- Declaration ---
var m map[string]int               // nil map — panics on write
m := map[string]int{}             // empty map — safe to write
m := map[string]int{"a": 1, "b": 2} // literal
m := make(map[string]int)         // initialized — safe to write
m := make(map[string]int, 100)    // with initial capacity hint

// --- Write ---
m["key"] = 42

// --- Read ---
val := m["key"]          // returns zero value if key absent — never panics
val, ok := m["key"]      // ok=true if key exists, ok=false if absent

// --- Check if key exists ---
if _, ok := m["key"]; ok {
    fmt.Println("exists")
}

// --- Delete ---
delete(m, "key")          // safe even if key doesn't exist

// --- Length ---
len(m)   // number of key-value pairs

// --- Traversal (order NOT guaranteed) ---
for k, v := range m { fmt.Println(k, v) }
for k := range m   { fmt.Println(k) }   // keys only

// --- Traversal in sorted order ---
keys := make([]string, 0, len(m))
for k := range m { keys = append(keys, k) }
sort.Strings(keys)
for _, k := range keys { fmt.Println(k, m[k]) }

// --- Nested map ---
m := map[string]map[string]int{}
m["outer"] = map[string]int{"inner": 1}

// --- Map of slices ---
m := map[string][]int{}
m["evens"] = append(m["evens"], 2, 4, 6)

// --- Frequency count ---
words := []string{"go", "is", "go", "fast"}
freq := make(map[string]int)
for _, w := range words { freq[w]++ }
// freq = {"go":2, "is":1, "fast":1}

// --- Deduplication ---
seen := make(map[string]struct{})
for _, v := range items {
    if _, ok := seen[v]; ok { continue }
    seen[v] = struct{}{}
    unique = append(unique, v)
}
```

| Operation | Time |
|---|---|
| Get / Set / Delete | O(1) average |
| Iteration order | **not guaranteed** |
| Concurrent read+write | **panics** — use Mutex or sync.Map |
| nil map read | safe (returns zero value) |
| nil map write | **panics** |

---

### Map concurrency
```go
// option 1: RWMutex — general purpose
var mu sync.RWMutex
mu.RLock(); val := m[k]; mu.RUnlock()
mu.Lock(); m[k] = v; mu.Unlock()

// option 2: sync.Map — best for write-once-read-many, or disjoint keys
var sm sync.Map
sm.Store("key", value)
val, ok := sm.Load("key")
sm.LoadOrStore("key", defaultVal) // atomic get-or-set
sm.Delete("key")
sm.Range(func(k, v any) bool {
    fmt.Println(k, v)
    return true // return false to stop iteration
})
```

---

### Side-by-side comparison

| | Slice | Map |
|---|---|---|
| Declaration | `make([]T, 0, cap)` | `make(map[K]V)` |
| Access | `s[i]` | `m[key]` or `m[key], ok` |
| Length | `len(s)` | `len(m)` |
| Check exists | loop or `slices.Contains` | `_, ok := m[key]` |
| Traverse | `for i, v := range s` | `for k, v := range m` |
| Order | maintained | **not guaranteed** |
| Zero value | nil — safe to append | nil — panics on write |
| Concurrent | not safe | not safe |
| Use when | ordered list, stack, queue | lookup, dedup, counting |

---

## 11. Testing

Go's standard library **`testing`** package — no third-party framework required. Tests live in `*_test.go` files, same package (or `package foo_test` for black-box tests).

```bash
go test ./...              # run all tests
go test -v -run TestAdd    # verbose, single test
go test -cover ./...       # coverage
go test -race ./...        # race detector
```

### Table-driven tests (Go standard)
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.expected {
                t.Errorf("Add(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.expected)
            }
        })
    }
}
```

### Mock interfaces with fake structs

**Idiomatic mocking:** depend on **interfaces**; in tests, pass a **fake struct** with the same methods (no code generation required).

```go
type UserRepository interface {
    GetUser(id int) (*User, error)
}

// fake implementation for tests
type mockUserRepo struct {
    user *User
    err  error
}
func (m *mockUserRepo) GetUser(id int) (*User, error) { return m.user, m.err }

func TestGetUser_NotFound(t *testing.T) {
    svc := NewService(&mockUserRepo{err: ErrNotFound})
    _, err := svc.GetUser(99)
    if !errors.Is(err, ErrNotFound) {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}
```

### HTTP handler tests
```go
func TestHandler(t *testing.T) {
    req := httptest.NewRequest("GET", "/user/1", nil)
    w := httptest.NewRecorder()

    MyHandler(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("expected 200, got %d", resp.StatusCode)
    }
}
```

### Benchmark tests
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
// go test -bench=. -benchmem
```

### Race detector — always run in CI
```bash
go test -race ./...
go run -race .
```

### Subtests and parallel
```go
t.Run("subtest", func(t *testing.T) {
    t.Parallel() // runs in parallel with other t.Parallel() subtests
    // test code
})
```

---

## 12. HTTP & Networking

### http.Client — configure timeouts
```go
// never use http.DefaultClient in production — no timeouts
client := &http.Client{
    Timeout: 10 * time.Second, // end-to-end timeout
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   3 * time.Second, // TCP connect timeout
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout:   5 * time.Second,
        ResponseHeaderTimeout: 5 * time.Second,
        MaxIdleConns:          100,
        MaxIdleConnsPerHost:   10,
        IdleConnTimeout:       90 * time.Second,
    },
}
```

**Never create a new `http.Client` per request** — each one creates a new connection pool. Create once and reuse.

### Request with context
```go
req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close() // always close — or connection leaks
body, _ := io.ReadAll(resp.Body)
```

### Graceful shutdown
```go
srv := &http.Server{Addr: ":8080", Handler: router}

go func() {
    if err := srv.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(ctx) // waits for in-flight requests to complete
```

### Middleware pattern

**Middleware** = wrapper around an `http.Handler` that runs **before/after** the next handler (logging, auth, metrics, recovery).

```go
type Middleware func(http.Handler) http.Handler

func Logger(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func Auth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Header.Get("Authorization") == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### Middleware chains

Chain middleware so the **first listed runs outermost** (like onion layers):

```go
// Chain applies m1, then m2, then m3 around final handler
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        h = middlewares[i](h)
    }
    return h
}

mux := http.NewServeMux()
mux.Handle("/api/", Chain(apiHandler,
    Logger,      // outer — logs entire request lifecycle
    Auth,        // middle — reject before handler
    RateLimit,   // inner — closest to handler
))

// Equivalent manual nesting:
// Logger(Auth(RateLimit(apiHandler)))
```

**Request flow:** `Logger` → `Auth` → `RateLimit` → `apiHandler` → response unwinds back through wrappers.

Common middleware: request ID, panic recovery, CORS, timeout via `context`, Prometheus metrics.

**With routers:** Chi, Echo, Gin expose `Use(mw)` — same idea, different API.

---

## 13. Database

### JavaScript/Database/sql — connection pool
```go
db, err := sql.Open("postgres", dsn)

// always configure pool — default is no limit = connection exhaustion
db.SetMaxOpenConns(25)                  // max connections to DB
db.SetMaxIdleConns(10)                  // kept open when idle
db.SetConnMaxLifetime(5 * time.Minute)  // recycle connections
db.SetConnMaxIdleTime(1 * time.Minute)  // close idle connections sooner

// pool exhaustion = requests queue up = latency spikes
// monitor with db.Stats()
stats := db.Stats()
log.Printf("open: %d, idle: %d, in-use: %d", stats.OpenConnections, stats.Idle, stats.InUse)
```

### Always use context
```go
// propagates cancellation and timeouts into DB calls
rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil { return err }
defer rows.Close() // always close rows

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil { return err }
}
return rows.Err() // check for iteration error
```

### Transactions
```go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // no-op if already committed

_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amount, fromID)
if err != nil { return err }

_, err = tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amount, toID)
if err != nil { return err }

return tx.Commit()
```

---

## 14. Go Modules & Dependency Management

### Go Modules (current standard — Go 1.16+ default)

**Modules** = dependency management with **`go.mod`** at project root. Versioned, reproducible builds.

```go
// go.mod
module github.com/you/myapp

go 1.22

require (
    github.com/lib/pq v1.10.9
    golang.org/x/sync v0.7.0
)
```

| Command | Purpose |
|---------|---------|
| `go mod init module/path` | Create new module |
| `go get pkg@v1.2.3` | Add or upgrade dependency |
| `go mod tidy` | Add missing, remove unused deps |
| `go mod vendor` | Copy deps into `vendor/` (offline/air-gapped builds) |
| `go list -m all` | List all module versions |
| `go mod download` | Cache modules locally |

**`go.sum`** — cryptographic checksums for reproducible builds. Commit both `go.mod` and `go.sum`.

**Semantic import versioning:** `v2+` modules require `/v2` in import path (e.g. `github.com/foo/bar/v2`).

### GOPATH (legacy)

Before modules, all code lived under **`$GOPATH/src`**. Dependencies were cloned into GOPATH with **no per-project version pinning** — "latest on disk" chaos.

| | **GOPATH** | **Go Modules** |
|---|------------|----------------|
| Project location | Must be under `$GOPATH/src` | Anywhere on disk |
| Version pinning | No (vendor or manual) | `go.mod` with semver |
| Reproducible builds | Hard | `go.sum` + module cache |
| Status | Legacy / deprecated for apps | **Use this** |

**Interview line:** "Modules replaced GOPATH with versioned dependencies in `go.mod`, a global module cache, and `go.sum` for integrity. I use `go mod tidy` and commit `go.mod`/`go.sum`."

---

## 15. Design Patterns in Go

### Functional options — flexible constructors
```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithPort(p int) Option       { return func(s *Server) { s.port = p } }
func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }

func NewServer(host string, opts ...Option) *Server {
    s := &Server{host: host, port: 8080, timeout: 30 * time.Second} // defaults
    for _, o := range opts { o(s) }
    return s
}

// usage
srv := NewServer("localhost", WithPort(9090), WithTimeout(10*time.Second))
```

### Worker pool
```go
func WorkerPool(ctx context.Context, jobs <-chan Job, numWorkers int) <-chan Result {
    results := make(chan Result, numWorkers)
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                select {
                case <-ctx.Done():
                    return
                case results <- process(job):
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

### Circuit breaker — prevent cascading failure
```go
// States: Closed (normal) → Open (failing, reject all) → Half-Open (test recovery)
type CircuitBreaker struct {
    failures  int
    threshold int
    state     string // "closed", "open", "half-open"
    mu        sync.Mutex
    resetAt   time.Time
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    if cb.state == "open" && time.Now().Before(cb.resetAt) {
        cb.mu.Unlock()
        return errors.New("circuit open")
    }
    cb.mu.Unlock()

    err := fn()
    cb.mu.Lock()
    defer cb.mu.Unlock()
    if err != nil {
        cb.failures++
        if cb.failures >= cb.threshold {
            cb.state = "open"
            cb.resetAt = time.Now().Add(30 * time.Second)
        }
    } else {
        cb.failures = 0
        cb.state = "closed"
    }
    return err
}
```

### Singleflight — prevent cache stampede
```go
import "golang.org/x/sync/singleflight"

var group singleflight.Group

func getUser(id string) (*User, error) {
    // only ONE call to DB even if 1000 goroutines request same id simultaneously
    v, err, _ := group.Do(id, func() (any, error) {
        return db.GetUser(id) // DB call runs once, result shared with all callers
    })
    if err != nil { return nil, err }
    return v.(*User), nil
}
```

### Pipeline pattern
```go
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums { out <- n }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in { out <- n * n }
    }()
    return out
}

// pipeline: generate → square → consume
for v := range square(generate(1, 2, 3, 4)) {
    fmt.Println(v)
}
```

---

## 16. Distributed Systems Concepts

### Idempotency
> The same operation applied multiple times produces the same result.  
> Critical for retries — without it, retrying a failed payment charges twice.

Implementation: idempotency key in request → check DB if already processed → skip if yes.

### Rate limiting strategies (from LLD)
- **Fixed window** — simple, allows burst at window boundary
- **Sliding window log** — precise, memory-heavy
- **Sliding window counter** — approximation, memory-efficient
- **Token bucket** — allows bursts up to bucket size
- **Leaky bucket** — smooths traffic, constant output rate

### Caching strategies
- **Cache-aside** (lazy): read → miss → fetch DB → store cache
- **Write-through**: write to cache + DB simultaneously
- **Write-behind**: write to cache → async write to DB
- **TTL**: expire after duration — simplest invalidation
- **Event-driven**: invalidate on change event (pub-sub)

### Singleflight for thundering herd
> 1000 requests hit expired cache → 1000 DB queries simultaneously = stampede.  
> Singleflight: deduplicate — 1 DB query, 999 waiters share the result.

### Distributed tracing
```go
// pass trace context via context.Context
ctx = context.WithValue(ctx, "traceID", uuid.New())

// OpenTelemetry: instrument at service boundaries
span := tracer.Start(ctx, "fetchUser")
defer span.End()
```

### CAP theorem (for system design discussions)
- **C**onsistency — every read gets latest write
- **A**vailability — every request gets a response
- **P**artition tolerance — system works despite network splits
- You can only fully guarantee 2 of 3. Most distributed systems choose AP or CP.

---

## 17. LLD / System Design

### LLD problems to know (all implemented in this repo as Go)
| Problem | Key concepts |
|---|---|
| Rate Limiter | token bucket, sliding window, fixed window, leaky bucket |
| LRU Cache | doubly linked list + map, O(1) get/put |
| Pub-Sub | observer pattern, goroutines for async delivery |
| Parking Lot | OOP design, strategy pattern for slot selection |
| Splitwise | graph-based debt simplification |
| URL Shortener | base62 encoding, hash collision, redirect |
| Database (in-memory) | schema, indexing, CRUD |
| Payment Gateway | strategy pattern for payment providers |
| Search Engine | inverted index, TF-IDF ranking, trie autocomplete |

### System design — how to approach
1. **Clarify** scope: scale, read/write ratio, consistency needs
2. **Estimate**: QPS, storage/day, bandwidth, memory
3. **High-level design**: clients, load balancer, services, DB, cache
4. **Deep dive**: bottlenecks, DB schema, cache strategy, failure handling
5. **Trade-offs**: always discuss what you gave up

### Latency debugging checklist
- Is it p50 or p99? (tail latency is often GC pauses or lock contention)
- Check goroutine count — leak = resource exhaustion
- Check DB pool stats — exhaustion = queued requests
- Check GC trace — `GODEBUG=gctrace=1`
- Check mutex profile — contention hotspots
- Check for missing timeouts on HTTP clients or DB calls

---

## 18. Classic Interview Gotchas — Predict the Output

> Interviewers love short snippets. Cover these cold. Say: **what happens**, **why**, **how to fix**.

### 1. Unbuffered send + receive in the same goroutine → deadlock

```go
func main() {
    ch := make(chan int)
    ch <- 10              // blocks forever — no other goroutine to receive
    fmt.Println(<-ch)     // never reached
}
```

**Output:** `fatal error: all goroutines are asleep - deadlock`

**Why:** Unbuffered channel needs a **rendezvous**. `main` is stuck on send, so it never reaches the receive.

**Fix:** another goroutine, or buffer: `make(chan int, 1)`.

---

### 2. Loop variable captured by goroutine

```go
func main() {
    for i := 0; i < 5; i++ {
        go func() {
            fmt.Println(i)
        }()
    }
    time.Sleep(time.Second)
}
```

| Go version | Typical output |
|------------|----------------|
| **&lt; 1.22** (classic interview answer) | `5 5 5 5 5` (order may vary) |
| **≥ 1.22** | `0 1 2 3 4` (any order) — each iteration has its own `i` |

**Why (pre-1.22):** all closures share **one** `i`; they usually run after the loop, when `i == 5`.

**Fix (always safe — say this in interviews):**

```go
for i := 0; i < 5; i++ {
    i := i                    // shadow: new var per iteration
    go func() { fmt.Println(i) }()
}
// or: go func(i int) { fmt.Println(i) }(i)
```

Also: prefer `WaitGroup` over `time.Sleep` to wait for goroutines.

---

### 3. Nil interface vs nil pointer

```go
func returnsError() error {
    var p *MyError = nil
    return p              // returns a typed nil
}
fmt.Println(returnsError() == nil) // false!
```

**Why:** Interface = `(type, value)`. Here type is `*MyError`, value is `nil` → interface is **not** nil.

**Fix:** `return nil` (untyped), or check with `errors.As` / typed helpers carefully.

---

### 4. Send / close on closed channel → panic

```go
ch := make(chan int, 1)
close(ch)
ch <- 1        // panic: send on closed channel
close(ch)      // panic: close of closed channel
```

Receive on closed channel is **safe** → zero value, `ok == false`.

**Rule:** only the **sender** closes; never close from the receiver; close **once**.

---

### 5. Nil map write → panic

```go
var m map[string]int   // nil map
fmt.Println(m["a"])    // 0 — read OK
m["a"] = 1             // panic: assignment to entry in nil map
```

**Fix:** `m = make(map[string]int)` or `m := map[string]int{}`.

---

### 6. Slice append shares backing array (unexpected overwrite)

```go
a := []int{1, 2, 3, 4}
b := a[:2]           // b = [1 2], same array as a
b = append(b, 99)    // may overwrite a[2]!
fmt.Println(a)       // often [1 2 99 4] — surprising
```

**Why:** `append` reuses capacity when there is room.

**Fix:** `b = append([]int(nil), a[:2]...)` or `slices.Clone`, or `a[:2:2]` (full slice expression — cap=2 forces realloc on append).

---

### 7. `defer` args evaluated immediately; defers run LIFO

```go
func main() {
    x := 1
    defer fmt.Println(x)  // prints 1 — arg snapshotted now
    x = 2
    defer fmt.Println("second")
    defer fmt.Println("first")
}
// Output:
// first
// second
// 1
```

**LIFO:** last `defer` runs first. **Args** of `defer f(x)` are evaluated when `defer` is registered, not when it runs. For late evaluation use `defer func() { fmt.Println(x) }()`.

---

### 8. `WaitGroup.Add` inside the goroutine → race / undercount

```go
var wg sync.WaitGroup
for i := 0; i < 5; i++ {
    go func() {
        wg.Add(1)       // WRONG — may race with Wait
        defer wg.Done()
        work()
    }()
}
wg.Wait()               // may return before all Add(1) calls
```

**Fix:** `wg.Add(1)` **before** `go`, or `wg.Add(n)` once before the loop.

---

### 9. Copying a mutex (or struct containing one)

```go
type Counter struct {
    mu sync.Mutex
    n  int
}
func (c Counter) Inc() {  // value receiver COPIES the mutex
    c.mu.Lock()
    c.n++
    c.mu.Unlock()
}
```

**Why:** each call locks a **different** mutex copy → no real exclusion; also `go vet` warns.

**Fix:** pointer receiver `func (c *Counter) Inc()`, never copy mutexes.

---

### 10. Concurrent map write → fatal (not a recover-able panic)

```go
m := map[int]int{}
go func() { m[1] = 1 }()
go func() { m[2] = 2 }()
// fatal error: concurrent map writes
```

**Fix:** `sync.Mutex` / `RWMutex`, or `sync.Map`, or shard maps. Always run `go test -race`.

---

### 11. Range over channel without close → hang forever

```go
ch := make(chan int)
go func() {
    ch <- 1
    ch <- 2
    // forgot close(ch)
}()
for v := range ch {  // waits forever for next value after 2
    fmt.Println(v)
}
```

**Fix:** sender `defer close(ch)` when done (only if no more sends).

---

### 12. `select` on nil channel never fires that case

```go
var ch chan int  // nil
select {
case <-ch:           // never selected — nil chan blocks forever
default:
    fmt.Println("ok")
}
```

**Useful trick:** set `ch = nil` to **disable** a `select` case dynamically.

---

### 13. Method set — pointer vs value receiver

```go
type T struct{}
func (t *T) M() {}

var t T
var i interface{ M() } = t   // compile error — T does not have M
var i2 interface{ M() } = &t // OK
```

**Rule:** `*T` methods are not in `T`'s method set. Prefer pointer receivers for mutating methods.

---

### 14. Shadowing with `:=` hides outer `err`

```go
err := doA()
if err != nil { return err }
if err := doB(); err != nil {  // new err in if scope — OK
    return err
}
// but this is a bug:
result, err := doC()           // fine if err already declared
if data, err := doD(); err != nil {
    return nil                 // returned WITHOUT the doD err if you return wrong var
}
_ = result
```

Prefer clear scopes; don’t reuse `:=` in a way that ignores the outer error you meant to return.

---

### 15. String / range — runes, not bytes

```go
s := "Go语言"
fmt.Println(len(s))           // 8 — bytes, not characters
for i, r := range s {
    fmt.Println(i, string(r)) // i is byte index; r is rune
}
```

Use `utf8.RuneCountInString` or `[]rune(s)` when you mean characters.

---

### 16. Array vs slice — value vs reference-ish

```go
a := [3]int{1, 2, 3}
b := a
b[0] = 99
fmt.Println(a[0]) // 1 — arrays are copied

s1 := []int{1, 2, 3}
s2 := s1
s2[0] = 99
fmt.Println(s1[0]) // 99 — same backing array
```

---

### Quick “predict output” drill table

| Snippet gist | Result |
|--------------|--------|
| Unbuffered `ch<-` then `<-ch` in `main` | Deadlock |
| Loop `go func(){ print(i) }` (Go &lt; 1.22) | All `len` / final `i` |
| `var e error = (*T)(nil)` then `e == nil` | `false` |
| `close` then send | Panic |
| Nil map write | Panic |
| `append` into subslice with spare cap | May mutate original |
| Multiple `defer` prints | LIFO order |
| `wg.Add` inside `go` | Race / early `Wait` |
| Value-receiver method + mutex | Broken locking |
| Concurrent map write | Fatal crash |
| `for range ch` without `close` | Hang |
| `len("世界")` | Byte length (6), not 2 |

---

## 19. Practice Questions

### Goroutines & Concurrency
1. Run 10 goroutines, track which failed, retry only the failed ones
2. Write a goroutine that restarts itself after a panic
3. How do you stop a goroutine from outside? Why can't you force-kill it?
4. What is the difference between `go func(){}()` and `go func(x int){}(x)`?
5. Implement a worker pool with N workers and graceful shutdown via context
6. What is a goroutine leak? How do you detect and fix it?
7. Call 3 APIs in parallel; return the first successful response and cancel the rest (`Promise.race` equivalent). Why can't you use `wg.Wait()` for this?
7a. Predict the output of the loop+goroutine `print(i)` snippet — and the unbuffered channel in `main` only (§18)

### Channels
8. What is the difference between buffered and unbuffered channels?
9. What happens if you send to a closed channel?
10. Implement a fan-in that merges results from 3 goroutines into one channel
11. How do you broadcast a stop signal to 10 goroutines simultaneously?
11a. Why does `for range ch` hang if the sender never closes?

### Memory & GC
12. What is escape analysis? How do you inspect it?
13. When would you deliberately call `runtime.GC()`?
14. What is `GOMEMLIMIT` and when is it essential?
15. What is `sync.Pool` and when do you use it?
16. What is a GC assist and how does it cause latency spikes?

### Error Handling
17. What is the difference between `errors.Is` and `errors.As`?
18. How do you wrap an error with context in Go?
19. When do you use a sentinel error vs a custom error type?

### Interfaces & Types
20. What is the nil interface gotcha?
21. What does "accept interfaces, return structs" mean?
22. Write a mock for a `UserRepository` interface for use in a unit test
23. What are generics in Go and when do you use them?

### Data Structures
24. Concurrent map writes — why does it panic and how do you fix it?
25. When do you use `sync.Map` vs map + Mutex?
26. What is the difference between `len` and `cap` on a slice?

### HTTP & DB
27. Why should you not use `http.DefaultClient` in production?
28. How do you configure connection pool for `JavaScript/Database/sql`?
29. How does graceful HTTP server shutdown work in Go?

### Design Patterns
30. What is the functional options pattern? Why is it preferred over config structs?
31. What is the singleflight pattern and what problem does it solve?
32. Explain the circuit breaker pattern and its 3 states

### Language & tooling (checklist topics)
33. Why Go? What problems does it solve?
34. Explain `fmt.Printf` vs `fmt.Sprintf` and common verbs (`%v`, `%T`)
35. What is reflection? When would you use it vs interfaces?
36. Go Modules vs GOPATH — how do you version dependencies?
37. What middleware is and how do you chain it in HTTP handlers?
38. Error values vs sentinel vs custom error types — pros and cons?
39. How do you profile and optimize a Go service (pprof, benchmarks)?
40. How do you avoid memory leaks despite GC?

---

## 20. Rapid Revision Cheat Sheet

```
Why Go                  → simple, fast compile, goroutines, static types, single binary, GC
fmt.Printf / Sprintf    → Printf prints; Sprintf returns string; %v value, %T type
Topic map               → §0 table maps interview checklist to sections

GMP model               → G=goroutine, M=OS thread, P=logical processor
Blocked G               → NOT on P's run queue; netpoller/channel/syscall wait; re-queued when ready
P–M binding             → 1 M per P at a time; P rebinds to new M if syscall blocks old M
Goroutine vs thread     → 2KB vs 1–8MB stack; millions vs thousands
Work stealing           → idle P steals goroutines from busy P's queue
GOMAXPROCS              → number of Ps (parallel Go code); default = CPU cores; you don't set M count

defer wg.Done()         → prevents wg.Wait() deadlock; NOT leak prevention
wg.Wait()               → wait for ALL goroutines; cannot early-return on first success
First-success race      → wait on <-ch (not wg.Wait); cancel() losers via context
Promise.race (Node)     → Go: parallel goroutines + buffered results channel + cancel()
ctx.Done()              → channel closed by cancel(); select on it to exit
Buffered channel        → send blocks only when full; size 1 lets single-send goroutine exit
Unbuffered in main only → ch<- then <-ch in same G = deadlock
close(ch)               → unblocks all receivers; range ch loop exits
Loop + go print(i)      → Go <1.22: often all final i; pass i or i:=i; Go 1.22+ per-iter i
Goroutine leak          → goroutine blocked forever; memory never GC'd
Detect leak             → runtime.NumGoroutine(); pprof goroutine dump; goleak
wg.Add                  → BEFORE go; never Add inside the goroutine vs Wait race

time.After in loop      → new timer every iteration; use NewTicker instead
defer resp.Body.Close() → HTTP connection returned to pool
defer rows.Close()      → DB connection returned to pool
defer order             → LIFO; defer f(x) snapshots x now
Mutex.Lock()            → one goroutine at a time (readers + writers)
Never copy mutex        → pointer receiver; go vet catches copies
RWMutex.RLock()         → many concurrent readers; writers blocked
RWMutex.Lock()          → exclusive; blocks all readers and writers
atomic ops              → multi-goroutine safe; CPU makes single read/write indivisible; NOT single-threaded
atomic vs mutex         → atomic = one word/one op; mutex = critical section / many fields
counter++ without atomic → data race; lost updates under parallel goroutines
CompareAndSwap (CAS)    → update only if value unchanged; retry loops for lock-free code
sync.Once               → exactly one execution across all goroutines
sync.Pool               → reuse objects across requests; reduces GC pressure

cancel()                → closes ctx.Done() channel; signals all goroutines
ctx.Err()               → context.Canceled or context.DeadlineExceeded
errgroup                → structured goroutines; returns first error; cancels rest

sync package            → Mutex, RWMutex, WaitGroup, Once, Pool, Cond, atomic
Concurrency pitfalls    → races, leaks, deadlock; -race in CI; §4 table + §18 gotchas

Stack                   → per-goroutine; reclaimed on return; zero GC cost
Heap                    → shared; GC managed; minimize allocations here
Memory leaks            → goroutine leaks, unclosed Body/rows, timers; not "GC broken"
Escape analysis         → go build -gcflags='-m'
pprof                   → CPU/heap/goroutine profiles; benchmark before/after
GOGC=100                → default; GC when heap doubles
GOMEMLIMIT              → hard ceiling; ~90% of container limit (Go 1.19+)
runtime.GC()            → force GC after large batch or before latency window
sync.Pool               → bypass GC for reusable objects (buffers, structs)

Interface / polymorphism → implicit satisfaction; accept interface, return struct
nil interface gotcha    → interface{nilConcretePtr} != nil
Multiple interfaces     → one type implements many; compose with embedding
Reflection              → reflect package; slow; JSON/ORM; prefer interfaces/generics
Go Modules              → go.mod + go.sum; go mod tidy; replaces GOPATH
GOPATH                  → legacy; no version pinning; code under $GOPATH/src

Middleware              → func(http.Handler) http.Handler; Chain wraps outer→inner
errors.Is / errors.As   → sentinel error match anywhere in chain
fmt.Errorf("%w", err)   → wrap error with context

Map concurrent          → not safe; use RWMutex or sync.Map
nil map write           → panic; make() first
sync.Map                → best for write-once-read-many pattern
Slice append            → always assign back; may reallocate
Slice sharing           → s[i:j] shares backing array; append may overwrite
len(string)             → bytes, not runes; range gives runes + byte index
send on closed ch       → panic; receive on closed → zero, ok=false

http.DefaultClient      → no timeouts; never use in production
Connection pool (DB)    → SetMaxOpenConns + SetMaxIdleConns + SetConnMaxLifetime
Graceful shutdown       → srv.Shutdown(ctx); waits for in-flight requests

Functional options      → flexible, backward-compatible constructors
Singleflight            → deduplicate concurrent calls; prevents stampede
Circuit breaker         → Closed → Open → Half-Open; prevents cascade failure
Worker pool             → N goroutines, shared jobs channel, ctx cancellation

Race detector           → go test -race ./... — always in CI
Benchmark               → go test -bench=. -benchmem
Table-driven tests      → Go standard; use t.Run for subtests
Gotchas drill           → §18 predict-the-output table
```
