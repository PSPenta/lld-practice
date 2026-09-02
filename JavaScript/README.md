# JavaScript & Node.js — Interview Preparation Guide

> Quick-revision doc for JS language + Node.js runtime interviews.  
> **Extended Q&A bank:** use [sudheerj/javascript-interview-questions](https://github.com/sudheerj/javascript-interview-questions) for hundreds of additional questions — this doc avoids duplicating that list and focuses on **structured answers + one-liners** for fast recall.  
> **LLD implementations (JS):** working code for Rate Limiter, Parking Lot, LRU, etc. lives in **this folder** — see [LLD implementations](#lld-implementations-in-this-folder) below. **Design method:** [../README.md](../README.md).  
> **TypeScript prep (types, generics, LLD in TS):** [../TypeScript/README.md](../TypeScript/README.md) — read **after** JS language core (§1–§12); runtime behaviour stays in this doc.  
> **REST / backend / payments API prep:** [../Backend/README.md](../Backend/README.md) — HTTP semantics, idempotency, pagination, OAuth; complements §26 here (JWT).

---

## Table of Contents
0. [How to Use This Doc & External Resources](#0-how-to-use-this-doc--external-resources)
1. [`this`, Lexical `this` & Explicit Binding](#1-this-lexical-this--explicit-binding)
2. [map, filter, reduce](#2-map-filter-reduce)
3. [Hoisting & Temporal Dead Zone (TDZ)](#3-hoisting--temporal-dead-zone-tdz)
4. [Rest vs Spread & Destructuring](#4-rest-vs-spread--destructuring)
5. [Closures & Memoisation](#5-closures--memoisation)
6. [`use strict`](#6-use-strict)
7. [Shallow Copy vs Deep Copy](#7-shallow-copy-vs-deep-copy)
8. [Dynamic Typing — How JS Stores Values](#8-dynamic-typing--how-js-stores-values)
9. [Programming Paradigms in JS](#9-programming-paradigms-in-js)
10. [Prototypes & Prototypal Inheritance](#10-prototypes--prototypal-inheritance)
11. [setTimeout, setImmediate, setInterval](#11-settimeout-setimmediate-setinterval)
12. [NaN vs null vs undefined](#12-nan-vs-null-vs-undefined)
13. [What Is Node.js? JS vs Node](#13-what-is-nodejs-js-vs-node)
14. [Why Node.js?](#14-why-nodejs)
15. [Event Loop Architecture](#15-event-loop-architecture)
16. [Sync vs Async vs Concurrency vs Parallelism](#16-sync-vs-async-vs-concurrency-vs-parallelism)
17. [V8 & JavaScript Engines](#17-v8--javascript-engines)
18. [Execution Model — Browser JS vs Node](#18-execution-model--browser-js-vs-node)
19. [Callbacks, Promises, async/await & Generators](#19-callbacks-promises-asyncawait--generators)
20. [Promise Combinators — all, allSettled, race, any](#20-promise-combinators--all-allsettled-race-any)
21. [Error Handling in Node.js](#21-error-handling-in-nodejs)
22. [Memory Management & GC in Node](#22-memory-management--gc-in-node)
23. [Child Processes — spawn, exec, fork](#23-child-processes--spawn-exec-fork)
24. [Streams in Node.js](#24-streams-in-nodejs)
25. [Node.js & the DOM](#25-nodejs--the-dom)
26. [Authentication & Authorization — Session vs JWT](#26-authentication--authorization--session-vs-jwt)
27. [Classic Interview Gotchas — Predict the Output](#27-classic-interview-gotchas--predict-the-output)
28. [Practice Checklist](#28-practice-checklist)
29. [Rapid Revision Cheat Sheet](#29-rapid-revision-cheat-sheet)

---

## 0. How to Use This Doc & External Resources

### Revision path (2–3 days)

| Day | Focus | Sections |
|-----|-------|----------|
| **1** | Language core | §1–§12, **§27** gotchas |
| **2** | Node runtime & async | §13–§21, **§20** Promise combinators |
| **3** | Production topics + drill | §22–§26, **§28** checklist, **§29** cheat sheet |

**Before interview:** answer **§28** aloud; skim **§29** (30 min); run **§27** without looking.

### sudheerj repo — use for depth, not duplication

| Use **this doc** for | Use **[sudheerj/javascript-interview-questions](https://github.com/sudheerj/javascript-interview-questions)** for |
|----------------------|---------------------------------------------------------------------------------------------------------------------|
| Structured one-pagers, event loop, Node architecture | 1000+ individual Q&A (closures, ES6+, DOM, React-adjacent) |
| Promise combinator comparison tables | Niche trivia, version-history questions |
| Session vs JWT, streams, child processes | Extra browser/DOM/browser-API questions |

**Rule:** master concepts here first; use sudheerj for **mock interview volume** on topics you already understand.

### Backend / REST API roles

For **Lead backend**, **fintech**, or **platform API** interviews, read **[../Backend/README.md](../Backend/README.md)** after Node basics (§13–§21). It covers POST vs PUT vs PATCH, idempotent payments, cursor pagination, status codes, OAuth, and developer-platform APIs. **§26 below** stays the quick JWT vs session reference.

### TypeScript roles

If the job description says **TypeScript**, **NestJS**, or **typed React**, finish **§1–§12** here (runtime + language), then switch to **[../TypeScript/README.md](../TypeScript/README.md)** for the type system, generics, `strict` tsconfig, and LLD with `interface` / discriminated unions. Do **not** skip JS — TS interviews still test event loop, Promises, and closures.

### LLD implementations in this folder

JavaScript **machine-coding / LLD** solutions — **design on paper first** via **[problems/](../README.md#21-lld-problem-checklist)** walkthroughs, then compare code:

| Folder | Problem | Walkthrough (design) |
|--------|---------|----------------------|
| `RateLimiter2/` | Rate limiter — **Strategy** (preferred) | [problems/rate-limiter](../problems/rate-limiter/README.md) |
| `Ratelimiter/` | Rate algorithms v1 — teaching | [problems/rate-limiter](../problems/rate-limiter/README.md) |
| `ParkingLot2/` | Multi-floor parking | [problems/parking-lot](../problems/parking-lot/README.md) |
| `Parkinglot/` | Parking v1 — inheritance-heavy | [problems/parking-lot](../problems/parking-lot/README.md) |
| `Splitwise/` | Expense split | [problems/splitwise](../problems/splitwise/README.md) |
| `SearchEngine/` | Index, rank, trie | [problems/search-engine](../problems/search-engine/README.md) |
| `Pub-Sub/` | Observer / event bus | [problems/pub-sub](../problems/pub-sub/README.md) |
| `Database/` | In-memory DB | [problems/in-memory-database](../problems/in-memory-database/README.md) |
| `Redis/` / `LRU/` | Cache + eviction | [problems/cache-client](../problems/cache-client/README.md) · [problems/lru-cache](../problems/lru-cache/README.md) |
| `PollingSystem2/` | Polls v2 — **PollingService** + repos | [problems/polling-system](../problems/polling-system/README.md) |
| `PollingSystem/` | v1 teaching / spot-the-bugs | [problems/polling-system](../problems/polling-system/README.md) |
| `Queue/` | FIFO queue | [problems/task-queue](../problems/task-queue/README.md) |
| `UrlShortener/` | Short URLs + Express | [problems/url-shortener](../problems/url-shortener/README.md) |
| `PaymentGateway/` | Stub — full impl in **`../Go/PaymentGateway-go/`** | [problems/payment-gateway](../problems/payment-gateway/README.md) |

Paper-only (no code yet): [problems/ai-suggest-reply](../problems/ai-suggest-reply/README.md) · [problems/ticket-notify](../problems/ticket-notify/README.md)

Go ports of the same problems: **`../Go/*-go/`** (PaymentGateway is **complete in Go only**).

---

## 1. `this`, Lexical `this` & Explicit Binding

### What is `this`?

`this` is **not** defined by where a function is written — it is determined by **how the function is called** (except arrow functions — see below).

| Call site | `this` value |
|-----------|--------------|
| Top-level in browser module/script (non-strict) | `window` |
| Top-level in **Node.js** module | `exports` / module wrapper — **not** `global` for `this` at module scope; use `globalThis` for true global |
| `obj.method()` | `obj` |
| `func()` alone (non-strict) | `window` / `globalThis` |
| `func()` in **strict mode** | `undefined` |
| `new Func()` | new empty object, then returned if constructor doesn't return object |
| Arrow function | **Lexical `this`** — inherits from enclosing scope; **cannot** be rebound |

```javascript
// Browser: bare this at top level → window
console.log(this === window); // true (non-module script)

// Node: module scope
console.log(this === exports); // true in CommonJS module
```

### Lexical `this` (arrow functions)

Arrow functions **do not have their own `this`**. They capture `this` from the enclosing lexical scope at creation time.

```javascript
const obj = {
  name: "Ada",
  regular() {
    setTimeout(function () {
      console.log(this.name); // undefined — `this` is window/global
    }, 0);
  },
  arrow() {
    setTimeout(() => {
      console.log(this.name); // "Ada" — lexical this from obj
    }, 0);
  },
};
```

**Interview line:** Regular functions get `this` from the call site; arrow functions get **lexical `this`** from where they were defined.

### Explicit binding — `call`, `apply`, `bind`

Force a specific `this` when invoking a function.

| Method | Invokes now? | Args | Returns |
|--------|--------------|------|---------|
| `fn.call(ctx, a, b)` | Yes | Comma-separated | Return value of `fn` |
| `fn.apply(ctx, [a, b])` | Yes | Array of args | Return value of `fn` |
| `fn.bind(ctx, a)` | No | Partial args OK | **New function** with fixed `this` |

```javascript
function greet(greeting) {
  console.log(`${greeting}, ${this.name}`);
}
const user = { name: "Bob" };

greet.call(user, "Hello");       // Hello, Bob
greet.apply(user, ["Hi"]);       // Hi, Bob
const bound = greet.bind(user);
bound("Hey");                    // Hey, Bob
```

**Use cases:** borrow methods (`Array.prototype.slice.call(arguments)`), fix `this` in callbacks (`this.handleClick = this.handleClick.bind(this)`), partial application with `bind`.

---

## 2. map, filter, reduce

| Method | Purpose | Returns |
|--------|---------|---------|
| **`map`** | Transform each element | **New array** same length |
| **`filter`** | Keep elements matching predicate | **New array** (≤ length) |
| **`reduce`** | Fold array to single value (sum, group, flatten) | **Any type** (accumulator) |

```javascript
const nums = [1, 2, 3, 4];
nums.map((n) => n * 2);           // [2, 4, 6, 8]
nums.filter((n) => n % 2 === 0);  // [2, 4]
nums.reduce((acc, n) => acc + n, 0); // 10
```

### Sync vs async

**All three run synchronously** — they do not wait for Promises inside the callback.

```javascript
// WRONG — async map doesn't await
await [1, 2, 3].map(async (id) => fetch(`/api/${id}`)); // array of Pending Promises

// RIGHT — for async work use for...of or Promise.all
const results = [];
for (const id of ids) {
  results.push(await fetch(`/api/${id}`));
}
// or
const results = await Promise.all(ids.map((id) => fetch(`/api/${id}`)));
```

**Interview line:** `map`/`filter`/`reduce` are sync; for async per-item work, use **`for...of` + await** or **`Promise.all` + map**.

---

## 3. Hoisting & Temporal Dead Zone (TDZ)

### Hoisting

JS **moves declarations to the top of their scope** during compilation (conceptually — not literally moving lines).

| Declaration | Hoisted? | Initial value before line runs |
|-------------|----------|--------------------------------|
| `var x` | Yes | `undefined` |
| `function foo(){}` | Yes (full function) | Callable |
| `let x` / `const x` | Yes (but in TDZ) | **Cannot access** until declaration |
| `class C` | Yes (TDZ) | Cannot access until declaration |

```javascript
console.log(a); // undefined (var hoisted)
var a = 1;

console.log(b); // ReferenceError — TDZ
let b = 2;
```

### Hoisting inside functions

Hoisting is **per scope**. `var` and `function` declarations hoist to the **top of the function** (or block for `let`/`const` in block scope).

```javascript
function demo() {
  console.log(x); // undefined
  var x = 10;
}
```

### Temporal Dead Zone (TDZ)

The region **from the start of the block until the `let`/`const` declaration line** where the binding exists but **must not be accessed**.

```javascript
{
  // TDZ for `value` starts here
  // console.log(value); // ReferenceError
  let value = 42;       // TDZ ends
}
```

**Interview line:** `var` hoists as `undefined`; `let`/`const` hoist but live in **TDZ** until initialized.

---

## 4. Rest vs Spread & Destructuring

### Rest (`...`) — collect remaining

Gathers remaining items into an **array** (or object in destructuring).

```javascript
function sum(first, ...rest) {
  return rest.reduce((a, b) => a + b, first);
}
sum(1, 2, 3, 4); // 10

const [head, ...tail] = [1, 2, 3]; // head=1, tail=[2,3]
```

### Spread (`...`) — expand

Expands iterable/object into **individual elements** (copy, merge, pass args).

```javascript
const arr = [1, 2, 3];
Math.max(...arr);                    // 3
const copy = [...arr];
const merged = { ...defaults, ...overrides };
fn(...args);
```

| | **Rest** | **Spread** |
|---|----------|------------|
| **Role** | Collect many → one | Expand one → many |
| **Position** | Last param / left side destructure | Call args, array/object literals |

### Destructuring

Pull values from arrays/objects into variables.

```javascript
const { name, age = 18 } = user;           // default age = 18
const { address: { city } } = user;         // nested destructure
const { name: userName } = user;            // rename

const [a, , c] = [1, 2, 3];                 // skip element
const [first, ...rest] = [1, 2, 3];
```

**Nested object + defaults:**

```javascript
const config = {
  server: {
    host: "localhost",
    // port omitted
  },
};
const {
  server: { host, port = 3000 } = {},
} = config;
// host = "localhost", port = 3000
```

---

## 5. Closures & Memoisation

### Closure

A **closure** is a function that remembers variables from its **outer lexical scope** even after the outer function has returned.

```javascript
function makeCounter() {
  let count = 0;
  return function () {
    return ++count;
  };
}
const counter = makeCounter();
counter(); // 1
counter(); // 2 — still has access to `count`
```

**Uses:** data privacy, factories, event handlers, partial application, memoisation.

### Memoisation

Cache function results by input to **avoid recomputing** the same expensive logic.

```javascript
function memoize(fn) {
  const cache = new Map();
  return function (...args) {
    const key = JSON.stringify(args);
    if (cache.has(key)) return cache.get(key);
    const result = fn.apply(this, args);
    cache.set(key, result);
    return result;
  };
}

const fib = memoize(function fib(n) {
  if (n <= 1) return n;
  return fib(n - 1) + fib(n - 2);
});
```

**Interview line:** Closure keeps the cache alive; memoisation trades memory for CPU.

---

## 6. `use strict`

`"use strict";` (file top or function body) enables **Strict Mode** — catches silent errors and disables unsafe features.

| Behaviour | Sloppy mode | Strict mode |
|-----------|-------------|-------------|
| Undeclared assignment | Creates global | **ReferenceError** |
| Duplicate param names | Allowed | **SyntaxError** |
| `this` in plain function call | `window` / global | **`undefined`** |
| `delete` variable | Silent fail | **SyntaxError** |
| `with` statement | Allowed | **SyntaxError** |

```javascript
"use strict";
function f() {
  console.log(this); // undefined (not window)
}
f();
```

---

## 7. Shallow Copy vs Deep Copy

| | **Shallow copy** | **Deep copy** |
|---|------------------|---------------|
| **Top-level** | New container | New container |
| **Nested objects** | **Shared references** | **Fully duplicated** |
| **Methods** | `Object.assign`, spread `{...obj}`, `[...arr]`, `slice()` | `structuredClone()`, JSON hack, lodash `cloneDeep` |

```javascript
const original = { a: 1, nested: { b: 2 } };

// Shallow
const shallow = { ...original };
shallow.nested.b = 99;
console.log(original.nested.b); // 99 — shared

// Deep (modern)
const deep = structuredClone(original);
deep.nested.b = 0;
console.log(original.nested.b); // 99 — unchanged

// Deep (JSON — loses Date, Map, undefined, functions)
const jsonDeep = JSON.parse(JSON.stringify(original));
```

---

## 8. Dynamic Typing — How JS Stores Values

JavaScript is **dynamically typed** — a variable can hold any type; type is tied to the **value**, not the variable name.

```javascript
let x = 42;       // number
x = "hello";      // now string — legal
```

**Behind the scenes (conceptual):**

- Primitives (`number`, `string`, `boolean`, `null`, `undefined`, `symbol`, `bigint`) are stored **by value** (engine-optimised; often inline for small ints).
- Objects (including arrays, functions) are stored **by reference** — variable holds a pointer to heap object.
- Engines (V8) use **hidden classes**, inline caches, and tagging for fast property access — you don't manage types manually.

**Interview line:** JS variables are untyped bindings; values have types at runtime; objects are reference types on the heap.

---

## 9. Programming Paradigms in JS

JavaScript is **multi-paradigm**:

| Paradigm | JS support | Example |
|----------|------------|---------|
| **Procedural** | Functions, statements | Scripts, utilities |
| **Object-oriented** | Objects, prototypes, classes (syntactic sugar) | `class User { }` |
| **Functional** | First-class functions, map/filter/reduce, immutability | Pure functions, composition |
| **Event-driven** | Callbacks, events, Promises | DOM events, Node `EventEmitter` |
| **Imperative** | Loops, mutation | `for`, `while` |
| **Declarative** | React JSX, SQL-like chains | `.filter().map()` |

**Interview line:** JS isn't purely OOP like Java — it's prototype-based OOP + functional + event-driven, especially in Node.

---

## 10. Prototypes & Prototypal Inheritance

### `__proto__` vs `Object.prototype`

| | **`Object.prototype`** | **`__proto__` (deprecated name)** |
|---|------------------------|-----------------------------------|
| **What** | The **root prototype object** — shared methods like `.toString()` | **Link** on an instance pointing to its prototype |
| **Better API** | — | `Object.getPrototypeOf(obj)` / `Object.setPrototypeOf()` |

```javascript
const obj = { a: 1 };
Object.getPrototypeOf(obj) === Object.prototype; // true

function Person(name) {
  this.name = name;
}
Person.prototype.greet = function () {
  return `Hi, ${this.name}`;
};
const p = new Person("Ada");
p.greet(); // lookup: p → Person.prototype → Object.prototype
```

### Prototypal vs classical inheritance

| | **Prototypal (JS)** | **Classical (Java/C++)** |
|---|---------------------|---------------------------|
| **Mechanism** | Objects inherit from objects via prototype chain | Classes/instances, compile-time hierarchy |
| **Syntax** | `Object.create`, prototype, `class` sugar | `extends`, interfaces |
| **Flexibility** | Can change prototype at runtime | Fixed at compile time |

```javascript
// ES6 class = syntactic sugar over prototypes
class Dog extends Animal {
  bark() { return "woof"; }
}
// Still prototype under the hood: Dog.prototype.__proto__ === Animal.prototype
```

### Interfaces in JavaScript (LLD / Strategy pattern)

JavaScript **does not** have a compile-time `interface` keyword like Java or TypeScript. Common ways to express a contract:

| Approach | When to use |
|----------|-------------|
| **Abstract base class** | LLD demos in this repo — method throws if not overridden (`RateLimiterStrategy`, `Expense`, `Vehicle`) |
| **Duck typing** | Any object with `isAllowed(ip)` works; no base class required |
| **JSDoc `@interface` + `@implements`** | Documents intent in plain `.js` files (used in `RateLimiter2/`, `Splitwise/`, `Parkinglot/`) |
| **TypeScript `interface`** | Real compile-time checks when you use TS |

```javascript
// Plain JS — interface-like contract (see JavaScript/RateLimiter2/RateLimiterStrategy.js)
class RateLimiterStrategy {
  isAllowed(ip) { throw new Error("abstract"); }
}
class TokenBucket extends RateLimiterStrategy {
  isAllowed(ip) { /* ... */ }
}

// TypeScript — actual interface
interface RateLimiterStrategy { isAllowed(ip: string): boolean; }
```

**Interview line:** “In JS I use an abstract base class or duck typing for Strategy; in TS I'd use an `interface`.”

---

## 11. setTimeout, setImmediate, setInterval

| API | Environment | Behaviour |
|-----|-------------|-----------|
| **`setTimeout(fn, ms)`** | Browser + Node | Run `fn` once after ≥ `ms` (not guaranteed exact) |
| **`setInterval(fn, ms)`** | Browser + Node | Repeat `fn` every ~`ms` until cleared |
| **`setImmediate(fn)`** | **Node only** | Run `fn` after I/O callbacks in current event loop turn |
| **`requestAnimationFrame`** | Browser | Before next paint |

**Node ordering (simplified):**

```text
sync code → process.nextTick queue → microtasks (Promises)
→ macrotasks (setTimeout, setImmediate, I/O callbacks)
```

```javascript
setTimeout(() => console.log("timeout"), 0);
setImmediate(() => console.log("immediate")); // Node only
Promise.resolve().then(() => console.log("microtask"));
// Typical: microtask → timeout/immediate order varies by context
```

**Interview line:** `setTimeout`/`setInterval` are timers; `setImmediate` is Node-specific, runs after poll phase I/O.

---

## 12. NaN vs null vs undefined

| | **`undefined`** | **`null`** | **`NaN`** |
|---|-----------------|------------|-----------|
| **Meaning** | Declared but not assigned / missing property | Intentional **empty** value | Invalid **number** result |
| **Type** | `undefined` | `object` (historical bug) | `number` |
| **Equality** | `undefined == null` → **true** | same | `NaN === NaN` → **false** |

```javascript
let a;              // undefined
const user = { name: "x" };
user.age;           // undefined — no property

const empty = null; // explicit "no object"

0 / 0;              // NaN
Number("abc");      // NaN
Number.isNaN(NaN);  // true — prefer over global isNaN
```

---

## 13. What Is Node.js? JS vs Node

| | **JavaScript** | **Node.js** |
|---|----------------|-------------|
| **What** | Programming language (ECMAScript spec) | **Runtime** that runs JS **outside the browser** |
| **APIs** | Language + host (DOM in browser) | Language + **Node APIs** (fs, http, crypto, …) |
| **Engine** | V8 (Chrome), SpiderMonkey, etc. | **V8** + **libuv** + Node bindings |
| **Use** | Web pages, now also servers/tools | APIs, CLIs, microservices, tooling |

**Interview line:** JavaScript is the language; Node.js is a **V8 + libuv** platform that adds I/O and system APIs.

---

## 14. Why Node.js?

| Advantage vs Java/Python (typical interview angles) |
|-----------------------------------------------------|
| **Single language** full-stack (JS front + back) |
| **Non-blocking I/O** — good for I/O-heavy, many concurrent connections |
| **npm ecosystem** — largest package registry |
| **Fast startup** vs JVM for small services |
| **JSON-native** — APIs and web fit naturally |

| Trade-offs to mention |
|-----------------------|
| CPU-bound work blocks the main thread unless offloaded (worker threads, child process) |
| Callback complexity (mitigated by Promises/async-await) |
| Less opinionated than Spring/Django — need discipline on structure |

---

## 15. Event Loop Architecture

Node.js is **event-driven** and **single-threaded for JavaScript execution** (main thread), with **libuv** handling async I/O.

```text
┌─────────────────────────────────────────────────────────┐
│                    YOUR JS CODE                          │
│              (runs on single main thread)                │
└─────────────────────────┬───────────────────────────────┘
                          │
                   ┌──────▼──────┐
                   │  CALL STACK  │  ← sync functions LIFO
                   └──────┬──────┘
                          │ stack empty?
                   ┌──────▼──────────────────────────────┐
                   │           EVENT LOOP                 │
                   │  phases: timers → pending → poll →   │
                   │  check (setImmediate) → close        │
                   └──────┬──────────────────────────────┘
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
  MICROTASK QUEUE   CALLBACK QUEUE    LIBUV THREAD POOL
  (Promises,        (setTimeout,       (fs, crypto,
   queueMicrotask)    I/O callbacks)     dns — heavy work)
```

| Component | Role |
|-----------|------|
| **Call stack** | Executes sync JS; one frame per function call |
| **Event loop** | Picks tasks when stack is empty |
| **Microtask queue** | Promises, `process.nextTick` (Node — runs before other microtasks) |
| **Callback/macrotask queue** | `setTimeout`, I/O completion callbacks |
| **libuv thread pool** | Default **4 threads** for file/crypto/DNS — offload blocking work from main JS thread |

**Interview line:** JS runs on one thread; slow I/O goes to libuv/OS; when ready, callbacks enqueue and run when the stack clears.

---

## 16. Sync vs Async vs Concurrency vs Parallelism

| Term | Meaning in Node context |
|------|-------------------------|
| **Synchronous** | Blocks call stack until done (`JSON.parse`, tight loop) |
| **Asynchronous** | Start now, finish later via callback/Promise (I/O, timers) |
| **Concurrency** | **Many tasks in progress** — interleaved on one thread (event loop) + libuv pool |
| **Parallelism** | **True simultaneous CPU work** — `worker_threads`, `cluster`, child processes |

```text
Node default:
  Main thread     → async/concurrent (one JS thread, many in-flight I/O)
  libuv pool      → concurrent file/crypto (few threads)
  worker_threads  → parallelism on multiple CPU cores
```

### "Channelisation" (libuv / thread pool)

Heavy or blocking OS operations are **delegated to libuv's thread pool** (default size 4, configurable via `UV_THREADPOOL_SIZE`). When complete, the callback is **queued back** to the main thread — this is how Node **channels** work off the single JS thread without blocking it for disk/crypto.

**Interview line:** Node is async on the main thread but **concurrent** via the event loop and libuv pool; **parallelism** needs `worker_threads` or multiple processes.

---

## 17. V8 & JavaScript Engines

### Why V8 in Node?

Node was built on **V8** (Google, Chrome) because when Ryan Dahl created Node (2009), V8 was **open source**, **fast**, and **embeddable** — ideal for a standalone JS runtime. Node wasn't tied to a browser.

| Engine | Host |
|--------|------|
| **V8** | Chrome, Node.js, Deno |
| **SpiderMonkey** | Firefox |
| **JavaScriptCore (Nitro)** | Safari |
| **Chakra** | Legacy Edge (deprecated) |

### V8 pros/cons (interview level)

| V8 advantages | Trade-offs |
|---------------|------------|
| Mature, heavily optimised JIT | Memory hungry vs some engines |
| Large community, Node ecosystem | Tied to Google's release cycle |
| `--inspect` debugging, snapshots | Not the only fast engine today |

**Interview line:** Node chose V8 for speed + embeddability; engine choice is about ecosystem and tooling, not raw language features.

---

## 18. Execution Model — Browser JS vs Node

| | **Browser JS** | **Node.js** |
|---|----------------|-------------|
| **Host objects** | DOM, `window`, `fetch` | `fs`, `http`, `process` |
| **Entry** | HTML script, module | `node file.js` |
| **Module systems** | ES modules + (bundlers) | CommonJS + ESM |
| **Concurrency** | Web Workers, main thread + event loop | Event loop + worker_threads + cluster |
| **Security** | Same-origin, sandbox | Full OS access — trust your code |

Both share: **call stack, heap, event loop pattern, microtasks before macrotasks** (with Node's `process.nextTick` nuance).

---

## 19. Callbacks, Promises, async/await & Generators

### Evolution

```text
Callbacks → Promise chains → async/await (syntax sugar over Promises)
```

| Style | Pros | Cons |
|-------|------|------|
| **Callbacks** | Simple, Node-native | Callback hell, error handling inconsistent |
| **Promises** | Composable, `.then/.catch` | Verbose chains |
| **async/await** | Reads like sync | Easy to serialise by mistake (`await` in loop) |

### Callback hell vs promise chaining vs async/await

Beyond readability:
- **Promises** standardise **one error path** (`.catch`)
- **async/await** enables **`try/catch`**, **`finally`**, loops with `await`
- **Promises** start microtasks; plain callbacks use macrotasks

### Generators — async without async/await

Use **`function*`** + manual Promise driving to write async flow that **looks** sequential:

```javascript
function* gen() {
  const a = yield fetch("/a").then((r) => r.json());
  const b = yield fetch(`/b?id=${a.id}`).then((r) => r.json());
  return b;
}

function runGenerator(g) {
  return new Promise((resolve, reject) => {
    const it = g();
    function step(result) {
      if (result.done) return resolve(result.value);
      Promise.resolve(result.value).then(
        (val) => step(it.next(val)),
        (err) => step(it.throw(err))
      );
    }
    step(it.next());
  });
}
```

Libraries like **co** wrapped this pattern before async/await was standard.

---

## 20. Promise Combinators — all, allSettled, race, any

| Method | Resolves when | Rejects when | Use case |
|--------|---------------|--------------|----------|
| **`Promise.all`** | **All** fulfill | **First** reject | All must succeed (batch fetch) |
| **`Promise.allSettled`** | **All** settled (always) | Never — returns `{status, value/reason}` | Audit all outcomes |
| **`Promise.race`** | **First** settle (fulfill or reject) | First reject if that's first | Timeout, first response |
| **`Promise.any`** | **First fulfill** | **All** reject (AggregateError) | First success wins |
| **`Promise.finally`** | After fulfill or reject | — | Cleanup (like `finally` block) |

### all vs allSettled

```javascript
// all — one failure fails everything
await Promise.all([p1, p2, p3]);

// allSettled — get every result
const results = await Promise.allSettled([p1, p2, p3]);
// [{status:'fulfilled', value:...}, {status:'rejected', reason:...}]
```

### any vs race

```javascript
// race — first to settle wins (even if rejection)
Promise.race([fastFail, slowOk]); // may reject

// any — first SUCCESS wins; ignore failures until all fail
Promise.any([fail1, fail2, ok3]); // resolves with ok3
```

---

## 21. Error Handling in Node.js

| Technique | When |
|-----------|------|
| **Error-first callbacks** `(err, data)` | Legacy Node APIs (`fs.readFile`) |
| **`.then().catch()` / `.finally()`** | Promise chains |
| **`try/catch` + `async/await`** | Modern async handlers |
| **`process.on('uncaughtException')`** | Last resort — log, graceful shutdown |
| **`process.on('unhandledRejection')`** | Forgotten Promise rejections |
| **Custom error classes** | Operational vs programmer errors |

```javascript
class AppError extends Error {
  constructor(message, statusCode = 500) {
    super(message);
    this.statusCode = statusCode;
    this.name = "AppError";
  }
}

// Express-style
app.use((err, req, res, next) => {
  console.error(err);
  res.status(err.statusCode || 500).json({ error: err.message });
});
```

**Best practices:** fail fast on programmer errors; distinguish operational errors (404, validation); never swallow errors; always handle `unhandledRejection` in production.

---

## 22. Memory Management & GC in Node

Node uses **V8's garbage collector** — primarily **generational** collection with **mark-and-sweep** (and incremental/concurrent marking in modern V8).

| Concept | Layman |
|---------|--------|
| **Young generation** | Short-lived objects, scavenged often |
| **Old generation** | Long-lived objects, mark-sweep-compact |
| **Mark-and-sweep** | Mark reachable from roots; sweep unreachable |

### Memory leaks in Node (despite GC)

| Cause | Fix |
|-------|-----|
| Global arrays/maps growing | Bound cache, TTL |
| Forgotten timers/listeners | `clearInterval`, `removeListener` |
| Closures holding large refs | Null out when done |
| Unclosed DB/HTTP connections | `close()`, pool limits |

**Tools:** `node --inspect`, heap snapshots in Chrome DevTools, `process.memoryUsage()`.

---

## 23. Child Processes — spawn, exec, fork

| Method | Shell? | Buffers stdout? | IPC channel? | Use |
|--------|--------|-----------------|--------------|-----|
| **`spawn`** | No (by default) | Stream — you pipe | Optional | Long-running, streaming (ffmpeg) |
| **`exec`** | Yes (`sh -c`) | Buffers (max buffer limit) | No | Short shell commands |
| **`fork`** | No | Inherits stdio optional | **Yes** (`send/on('message')`) | Run another **Node** script |

```javascript
import { spawn, exec, fork } from "child_process";

spawn("ls", ["-la"]); // streams
exec("ls -la", (err, stdout) => { }); // buffered
const child = fork("./worker.js");
child.send({ job: 1 });
child.on("message", (msg) => { });
```

**Interview line:** `fork` = Node-to-Node with IPC; `spawn` = streaming; `exec` = shell one-liner with buffered output.

---

## 24. Streams in Node.js

**Streams** process data **chunk by chunk** instead of loading everything into memory.

| Type | Direction | Example |
|------|-----------|---------|
| **Readable** | Source | `fs.createReadStream`, HTTP request body |
| **Writable** | Destination | `fs.createWriteStream`, HTTP response |
| **Duplex** | Both | TCP socket |
| **Transform** | Read + modify + write | `zlib.createGzip()` |

```javascript
import { createReadStream, createWriteStream } from "fs";
import { pipeline } from "stream/promises";

await pipeline(
  createReadStream("input.txt"),
  createWriteStream("output.txt")
);
```

**Why:** constant memory for large files; backpressure handled automatically with `pipeline`.

---

## 25. Node.js & the DOM

**No.** Node.js has **no DOM** (`document`, `window`, HTML elements). It runs on the server/CLI.

- **Browser JS** → DOM APIs
- **Node.js** → `fs`, `http`, `path`, etc.
- **SSR / testing:** JSDOM or Puppeteer/Playwright provide DOM-like APIs when needed

---

## 26. Authentication & Authorization

| | **Authentication** | **Authorization** |
|---|---------------------|---------------------|
| **Question** | Who are you? | What can you do? |
| **Example** | Login, JWT verify | Role `admin`, scope `read:users` |

### Session vs JWT

| | **Session (server-side)** | **JWT (stateless token)** |
|---|---------------------------|----------------------------|
| **Storage** | Session ID in cookie; data on **server** (Redis) | Signed token on **client** |
| **Revocation** | Easy — delete session | Hard — need blocklist or short TTL |
| **Scale** | Needs shared session store | Easy horizontal scale |
| **Size** | Small cookie | Token can grow large |
| **Best for** | Traditional web apps, instant logout | Microservices, mobile APIs, SPAs |

**Interview line:** Sessions = server state, easy revoke; JWT = stateless, scale-friendly — often **refresh token + short-lived access JWT** hybrid.

Neither is universally "better" — trade-offs on scale, logout, and security.

---

## 27. Classic Interview Gotchas — Predict the Output

| # | Snippet gist | Answer |
|---|--------------|--------|
| 1 | `var x = 1; function f(){ console.log(x); var x = 2; } f();` | `undefined` (hoisting) |
| 2 | `let a = 1; { console.log(a); let a = 2; }` | ReferenceError (TDZ) |
| 3 | `[] + []` | `""` (array toString → empty + empty) |
| 4 | `[] + {}` vs `{} + []` | `"" + "[object Object]"` vs object literal confusion |
| 5 | `typeof null` | `"object"` |
| 6 | `0.1 + 0.2 === 0.3` | `false` (float precision) |
| 7 | `this` in arrow vs regular in object method | Arrow inherits outer `this` |
| 8 | `Promise.resolve().then(() => console.log(1)); console.log(2);` | `2` then `1` (microtask) |
| 9 | `setTimeout(() => console.log(1), 0); console.log(2);` | `2` then `1` |
| 10 | `async function f(){ return 1; } f().then(console.log)` | `1` (async returns Promise) |
| 11 | Closure in loop with `var` + `setTimeout` | All print same final `i` |
| 12 | `==` vs `===` for `null`/`undefined` | `null == undefined` true; `===` false |

---

## 28. Practice Checklist

Answer aloud without looking:

**JavaScript**
- [ ] Explain `this` binding rules + lexical `this` in arrows
- [ ] `call` vs `apply` vs `bind`
- [ ] Why `map` doesn't await async callbacks
- [ ] Hoisting + TDZ
- [ ] Rest vs spread
- [ ] Nested destructuring with defaults
- [ ] Closure + memoisation example
- [ ] Shallow vs deep copy methods
- [ ] Prototype chain vs `class`
- [ ] NaN / null / undefined differences

**Node.js**
- [ ] JS vs Node vs browser
- [ ] Event loop: call stack, microtasks, macrotasks, libuv pool
- [ ] Concurrency vs parallelism in Node
- [ ] Why V8
- [ ] Callback → Promise → async/await → when generators
- [ ] all vs allSettled, any vs race
- [ ] Error handling patterns + unhandledRejection
- [ ] GC + common Node memory leaks
- [ ] spawn vs exec vs fork
- [ ] Stream types + why pipeline
- [ ] Session vs JWT trade-offs

**Backend / REST** (full guide: [../Backend/README.md](../Backend/README.md))
- [ ] POST vs PUT vs PATCH + idempotency
- [ ] Idempotent payment with Idempotency-Key
- [ ] 401 vs 403 vs 409; cursor vs offset pagination

**Extended drill:** [sudheerj/javascript-interview-questions](https://github.com/sudheerj/javascript-interview-questions)

---

## 29. Rapid Revision Cheat Sheet

```
this (regular)          → from call site; strict mode plain call → undefined
this (arrow)            → lexical — from enclosing scope
call/apply/bind         → explicit this; bind returns new function
map/filter/reduce       → sync; async → for...of or Promise.all + map

hoisting                → var → undefined; function decl callable; let/const → TDZ
TDZ                     → let/const from block start until declaration line
rest                    → collect ...args; spread → expand iterable/object
closure                 → inner fn remembers outer scope
memoize                 → cache by input key

use strict              → no implicit globals; this undefined on plain call
shallow copy            → spread/assign — nested shared
deep copy               → structuredClone (modern)

dynamic typing          → variable untyped; value has type; objects by reference
prototypes              → chain lookup; class is sugar
__proto__               → use getPrototypeOf instead

setTimeout              → macrotask; 0 ≠ immediate
setImmediate            → Node only; after I/O phase
microtasks              → Promises before setTimeout

undefined               → missing; null → intentional empty; NaN → bad number

Node.js                 → V8 + libuv; JS on server
event loop              → stack → microtasks → macrotasks
libuv thread pool       → fs/crypto; default 4 threads; UV_THREADPOOL_SIZE
concurrency             → many in-flight on one JS thread
parallelism             → worker_threads / cluster / fork

Promise.all             → all succeed or first reject
Promise.allSettled      → wait for all; never rejects
Promise.race            → first settle (fail or success)
Promise.any             → first success; all fail → AggregateError

errors Node             → err-first cb; try/catch await; unhandledRejection
spawn/exec/fork         → stream / shell buffer / Node IPC
streams                 → readable writable duplex transform; pipeline
no DOM in Node          → use JSDOM or browser for DOM

session                 → server state; easy revoke
JWT                     → stateless; scale; short TTL + refresh pattern
REST/API depth          → Backend/README.md (POST/PUT/PATCH, payments, pagination)

sudheerj repo           → extra Q&A volume after this doc
```
