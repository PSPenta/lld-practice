# TypeScript — Interview Preparation Guide

> **SDE-3 / Staff** — type system, generics, API design, LLD patterns, production `tsconfig`, and revision cheat sheet.  
> **Prerequisite:** [../JavaScript/README.md](../JavaScript/README.md) — TS is a **type layer on JS**; runtime behaviour (event loop, `this`, Promises) lives there.  
> **LLD method & problems:** [../README.md](../README.md) · JS implementations: [../JavaScript/](../JavaScript/) · Go ports: [../Go/](../Go/)

---

## Table of Contents

0. [How to use this doc & interview topic map](#0-how-to-use-this-doc--interview-topic-map)
1. [What TypeScript is (and is not)](#1-what-typescript-is-and-is-not)
2. [`type` vs `interface` — Staff decision guide](#2-type-vs-interface--staff-decision-guide)
3. [Structural typing & assignability](#3-structural-typing--assignability)
4. [Narrowing & type guards](#4-narrowing--type-guards)
5. [Discriminated unions & exhaustive `switch`](#5-discriminated-unions--exhaustive-switch)
6. [Generics — interview depth](#6-generics--interview-depth)
7. [Utility & mapped types](#7-utility--mapped-types)
8. [OOP in TypeScript — LLD patterns](#8-oop-in-typescript--lld-patterns)
9. [Functions, overloads & `this`](#9-functions-overloads--this)
10. [`any`, `unknown`, `never`, `void`](#10-any-unknown-never-void)
11. [`as const`, `satisfies`, branded types](#11-as-const-satisfies-branded-types)
12. [Modules, `.d.ts`, and project layout](#12-modules-dts-and-project-layout)
13. [`tsconfig` — what Staff must know](#13-tsconfig--what-staff-must-know)
14. [Async typing — Promises & Result types](#14-async-typing--promises--result-types)
15. [Error handling at scale](#15-error-handling-at-scale)
16. [Testing & mocking with interfaces](#16-testing--mocking-with-interfaces)
17. [Node / Express / NestJS typing](#17-node--express--nestjs-typing)
18. [Variance — covariance & contravariance](#18-variance--covariance--contravariance)
19. [TypeScript for LLD interviews](#19-typescript-for-lld-interviews)
20. [Classic gotchas — predict compile/runtime](#20-classic-gotchas--predict-compileruntime)
21. [Staff / system-design talking points](#21-staff--system-design-talking-points)
22. [Practice checklist](#22-practice-checklist)
23. [Rapid revision cheat sheet](#23-rapid-revision-cheat-sheet)

---

## 0. How to use this doc & interview topic map

### Revision path (2–3 days)

| Day | Focus | Sections |
|-----|-------|----------|
| **1** | Type system core | §1–§5, **§20** gotchas |
| **2** | Generics, utilities, OOP | §6–§9, **§19** LLD |
| **3** | Production TS + drill | §10–§18, **§22** checklist, **§23** cheat sheet |

**Before interview:** answer **§22** aloud; skim **§23** (30 min); whiteboard one LLD with interfaces from **§19**.

### External resources (don't duplicate)

| Use **this doc** | Use elsewhere |
|------------------|---------------|
| Interview one-liners, LLD + types | [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/) — full reference |
| `tsconfig` trade-offs | [Total TypeScript](https://www.totaltypescript.com/) — advanced puzzles |
| Staff API / module boundaries | Your team's style guide + `strict` baseline |

### Interview topic map

| Topic | Section |
|-------|---------|
| `type` vs `interface`, declaration merging | **§2** |
| Structural typing, excess property checks | **§3** |
| Narrowing, user-defined type guards | **§4** |
| Discriminated unions, `never` exhaustiveness | **§5** |
| Generics, constraints, `infer` | **§6** |
| `Partial`, `Pick`, `Record`, mapped types | **§7** |
| Strategy / Repository as `interface` | **§8**, **§19** |
| `unknown` vs `any`, Result types | **§10**, **§15** |
| `strict`, `strictNullChecks`, project references | **§13** |
| Mocking repos in tests | **§16** |
| Function parameter variance | **§18** |
| Nest decorators (high level) | **§17** |

### Repo cross-links

| Language | LLD code | Prep doc |
|----------|----------|----------|
| JavaScript | [../JavaScript/](../JavaScript/) | [../JavaScript/README.md](../JavaScript/README.md) |
| TypeScript | Patterns in **§19** (port `.ts` when ready) | **This doc** |
| Go | [../Go/](../Go/) | [../Go/README.md](../Go/README.md) |
| REST / Backend | — | [../Backend/README.md](../Backend/README.md) |

---

## 1. What TypeScript is (and is not)

TypeScript = **JavaScript + static types**, erased at compile time. **No runtime type checking** unless you add libraries (zod, io-ts).

| | Compile time | Runtime |
|---|--------------|---------|
| `interface User { id: string }` | Yes | **Gone** — not in JS output |
| `enum Status { Active }` | Yes | **Emitted** as JS object (unless `const enum`) |
| `as User` type assertion | Checks (loose) | **No check** — trust the developer |
| `zod.parse(data)` | Inferred types | **Validates** at runtime |

**Interview line (Staff):** “TS catches contract violations before deploy; boundary code (HTTP, queues, third-party JSON) still needs runtime validation. Types are documentation the compiler enforces.”

### Three roles of types in large codebases

1. **Correctness** — null safety, exhaustive unions, wrong field access  
2. **API design** — public surfaces are explicit (`export type CreateUserDto`)  
3. **Refactoring** — rename/move with compiler as safety net  

---

## 2. `type` vs `interface` — Staff decision guide

| | **`interface`** | **`type`** |
|---|-------------------|------------|
| **Extends** | `extends` (merge-friendly) | `&` intersection |
| **Declaration merging** | ✅ Same name merges | ❌ Duplicate name = error |
| **Unions / primitives** | ❌ | ✅ `type Id = string \| number` |
| **Mapped / conditional** | ❌ | ✅ `type Readonly<T> = ...` |
| **implements** | ✅ Classes implement | ✅ (via intersection shape) |
| **Performance (large projects)** | Slightly faster for IDE checks on huge hierarchies | Fine for most cases |

### When to use which (interview answer)

```typescript
// interface — public OOP contract, extension across files, library augmentation
interface PaymentGateway {
  charge(amount: Money): Promise<PaymentResult>;
}

// type — unions, utilities, composition, DTOs
type PaymentStatus = 'pending' | 'succeeded' | 'failed';

type ApiResponse<T> =
  | { ok: true; data: T }
  | { ok: false; error: string };
```

**Staff rule:** Use **`interface`** for **object shapes you expect to extend** (services, repositories). Use **`type`** for **unions, tuples, mapped types, and DTO aliases**.

### Declaration merging (libraries)

```typescript
// express augments Request — only works with interface
declare global {
  namespace Express {
    interface Request {
      userId?: string;
    }
  }
}
```

---

## 3. Structural typing & assignability

TypeScript uses **structural** (duck) typing: if it has the shape, it fits.

```typescript
interface Point { x: number; y: number; }
const p = { x: 1, y: 2, z: 3 };
const pt: Point = p; // OK — extra z allowed when assigning object literal to variable first

const pt2: Point = { x: 1, y: 2, z: 3 }; // ERROR — excess property check on fresh literal
```

| Concept | Meaning |
|---------|---------|
| **Excess property check** | Fresh object literals can't have extra keys vs target type |
| **Width subtyping** | More properties → assignable to fewer **when** not a fresh literal |
| **Readonly** | Shallow freeze on properties; nested objects still mutable |

**Interview line:** “TS is not Java — no `implements` required for structural match; that's why duck typing works but also why accidental matches happen.”

---

## 4. Narrowing & type guards

| Mechanism | Example |
|-----------|---------|
| `typeof` | `typeof x === 'string'` |
| `instanceof` | `err instanceof Error` |
| `in` | `'code' in err` |
| Equality | `status === 'active'` |
| Truthiness | `if (user)` — narrows null/undefined (with strictNullChecks) |
| **User-defined guard** | `function isPoll(x): x is Poll` |

```typescript
interface Admin { role: 'admin'; permissions: string[] }
interface Member { role: 'member' }
type User = Admin | Member;

function canDelete(u: User): boolean {
  if (u.role === 'admin') {
    return u.permissions.includes('delete'); // narrowed to Admin
  }
  return false;
}

function isError(value: unknown): value is Error {
  return value instanceof Error;
}
```

**Staff pattern:** Prefer **discriminant + guard** over `as` casts at system boundaries.

---

## 5. Discriminated unions & exhaustive `switch`

```typescript
type VoteResult =
  | { kind: 'success'; pollId: number }
  | { kind: 'not_assigned'; pollId: number }
  | { kind: 'poll_closed' };

function message(r: VoteResult): string {
  switch (r.kind) {
    case 'success':
      return `Voted on poll ${r.pollId}`;
    case 'not_assigned':
      return 'You are not assigned to this poll';
    case 'poll_closed':
      return 'Poll is closed';
    default: {
      const _exhaustive: never = r;
      return _exhaustive;
    }
  }
}
```

Adding a new variant without updating `switch` → **compile error** on `never`.

**Interview line:** “Discriminated unions model domain outcomes explicitly — better than optional fields or error codes as magic strings.”

---

## 6. Generics — interview depth

### Basics

```typescript
function first<T>(items: T[]): T | undefined {
  return items[0];
}

interface Repository<T, ID = string> {
  getById(id: ID): Promise<T | null>;
  save(entity: T): Promise<void>;
}
```

### Constraints

```typescript
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
  return obj[key];
}
```

### `infer` (Staff favourite)

```typescript
// Extract return type without ReturnType utility
type MyReturn<T> = T extends (...args: never[]) => infer R ? R : never;

// Unwrap Promise
type Awaited<T> = T extends Promise<infer U> ? U : T;
```

### Defaults & multiple type params

```typescript
type Paginated<T, Cursor = string> = {
  items: T[];
  nextCursor: Cursor | null;
};
```

### Anti-patterns

| Bad | Better |
|-----|--------|
| `function foo<T>(x: T): T` everywhere | Concrete types at leaves; generics at boundaries |
| `any` in generic default `T = any` | `unknown` or omit default |
| 5+ type parameters on one function | Break into object type param |

---

## 7. Utility & mapped types

| Utility | Effect | LLD use |
|---------|--------|---------|
| `Partial<T>` | All optional | PATCH updates |
| `Required<T>` | All required | After validation |
| `Readonly<T>` | Shallow readonly | Immutable DTO |
| `Pick<T, K>` | Subset keys | Public API slice |
| `Omit<T, K>` | Exclude keys | Hide internal fields |
| `Record<K, V>` | Key → value map | `Record<PollId, Poll>` |
| `ReturnType<F>` | Function return | Wrap service methods |
| `Parameters<F>` | Function args tuple | Decorators / wrappers |

```typescript
type CreatePollInput = Pick<Poll, 'question' | 'options' | 'isPrivate'>;
type PollUpdate = Partial<CreatePollInput>;

type Flags<T> = {
  [K in keyof T]: boolean;
};
```

**Interview line:** “Utilities encode DTO transformations without duplicating field lists — single source of truth from domain model.”

---

## 8. OOP in TypeScript — LLD patterns

### Strategy (preferred in LLD)

```typescript
interface RateLimiterStrategy {
  isAllowed(key: string): boolean;
}

class TokenBucket implements RateLimiterStrategy {
  constructor(private capacity: number, private refillRate: number) {}
  isAllowed(key: string): boolean { /* ... */ return true; }
}

class RateLimiter {
  constructor(private strategy: RateLimiterStrategy) {}
  isAllowed(key: string): boolean {
    return this.strategy.isAllowed(key);
  }
}
```

### Repository (interface + in-memory impl)

```typescript
interface PollRepository {
  add(poll: Poll): void;
  getById(id: PollId): Poll | undefined;
  update(poll: Poll): void;
}

class InMemoryPollRepository implements PollRepository {
  private polls = new Map<PollId, Poll>();
  add(poll: Poll): void { this.polls.set(poll.id, poll); }
  getById(id: PollId): Poll | undefined { return this.polls.get(id); }
  update(poll: Poll): void { this.polls.set(poll.id, poll); }
}
```

### Abstract class vs interface

| Use **interface** | Use **abstract class** |
|-------------------|------------------------|
| Multiple implementations, DI, testing | Shared concrete helper methods |
| Strategy, Repository, Notifier | Template Method with partial impl |

```typescript
abstract class Expense {
  constructor(protected amount: number) {}
  abstract validate(): boolean;
  apply(): void { if (!this.validate()) throw new Error('invalid'); }
}
```

### Access modifiers

| Modifier | Meaning |
|----------|---------|
| `public` | Default — anywhere |
| `private` | Class body only (runtime enforced in TS 4.3+ with `#` optional) |
| `protected` | Class + subclasses |
| `readonly` | Assign once (shallow) |

**Interview line:** “I use interfaces at module boundaries; classes for entities with behaviour; `readonly` on IDs and value objects.”

---

## 9. Functions, overloads & `this`

### Overloads (API ergonomics)

```typescript
function createPoll(question: string, options: string[]): Poll;
function createPoll(input: CreatePollInput): Poll;
function createPoll(a: string | CreatePollInput, b?: string[]): Poll {
  if (typeof a === 'string') {
    return new Poll(a, b ?? []);
  }
  return new Poll(a.question, a.options);
}
```

### `this` parameter (not runtime)

```typescript
interface Clickable {
  label: string;
  onClick(this: Clickable, ev: Event): void;
}
```

Arrow functions **don't** have their own `this` — same as JS (see [JavaScript README §1](../JavaScript/README.md)).

---

## 10. `any`, `unknown`, `never`, `void`

| Type | Meaning | Staff usage |
|------|---------|-------------|
| **`any`** | Opt out of checking | **Ban** in prod code (`noImplicitAny`); legacy migration only |
| **`unknown`** | Safe top type — must narrow before use | JSON parse, catch blocks, plugin boundaries |
| **`never`** | No values — unreachable | Exhaustive switch, throw helpers |
| **`void`** | No useful return | Callbacks, `Promise<void>` |

```typescript
function parseJson(raw: string): unknown {
  return JSON.parse(raw);
}

function assertNever(x: never): never {
  throw new Error(`Unexpected: ${JSON.stringify(x)}`);
}
```

**Interview line:** “`unknown` forces handling; `any` hides bugs. Staff bar: zero explicit `any` in new code.”

---

## 11. `as const`, `satisfies`, branded types

### `as const` — literal inference

```typescript
const OPTIONS = ['Delhi', 'Mumbai'] as const;
type Option = (typeof OPTIONS)[number]; // 'Delhi' | 'Mumbai'
```

### `satisfies` — validate without widening

```typescript
const config = {
  privateDefault: true,
  maxOptions: 10,
} satisfies PollConfig;
// config.privateDefault stays literal true, but keys are checked against PollConfig
```

### Branded types (nominal-ish)

```typescript
type PollId = string & { readonly __brand: 'PollId' };
type UserId = string & { readonly __brand: 'UserId' };

function pollId(id: string): PollId {
  return id as PollId;
}

// pollId and userId can't be swapped accidentally
```

**Staff use:** IDs, currency, email — prevent primitive obsession bugs.

---

## 12. Modules, `.d.ts`, and project layout

```text
src/
  domain/
    models/Poll.ts
    repositories/PollRepository.ts
  application/PollingService.ts
  infrastructure/InMemoryPollRepository.ts
  index.ts
```

| File | Purpose |
|------|---------|
| `.ts` | Implementation |
| `.d.ts` | Type-only declarations (ambient types, `@types/*`) |
| `export type` | Re-export types without value (tree-shake friendly) |

```typescript
// types-only import — erased at compile
import type { Poll } from './models/Poll';
```

**Barrel files (`index.ts`):** convenient but can hurt tree-shaking and circular deps — Staff teams often ban deep barrels or use `export type` only.

---

## 13. `tsconfig` — what Staff must know

### Non-negotiable strict baseline

```json
{
  "compilerOptions": {
    "strict": true,
    "noImplicitAny": true,
    "strictNullChecks": true,
    "strictFunctionTypes": true,
    "noUncheckedIndexedAccess": true,
    "exactOptionalPropertyTypes": true,
    "noImplicitOverride": true,
    "verbatimModuleSyntax": true
  }
}
```

| Flag | Why it matters |
|------|----------------|
| `strictNullChecks` | `null`/`undefined` explicit — kills NPE class |
| `noUncheckedIndexedAccess` | `arr[i]` → `T \| undefined` |
| `exactOptionalPropertyTypes` | `{ x?: number }` can't assign `undefined` unless allowed |
| `verbatimModuleSyntax` | `import type` required — clear value vs type imports |
| `isolatedModules` | Each file transpile-safe (Babel/esbuild) |

### Project references (monorepos)

```json
{ "references": [{ "path": "./packages/domain" }] }
```

**Interview line:** “I'd rather tighten `tsconfig` than add `@ts-ignore`. Project references keep large repos incremental.”

---

## 14. Async typing — Promises & Result types

```typescript
type Result<T, E = Error> =
  | { ok: true; value: T }
  | { ok: false; error: E };

async function submitVote(
  voterId: UserId,
  pollId: PollId,
  option: string,
): Promise<Result<Vote, 'NOT_ASSIGNED' | 'POLL_CLOSED' | 'INVALID_OPTION'>> {
  // ...
  return { ok: true, value: vote };
}
```

| Pattern | When |
|---------|------|
| `throw` on failure | Express middleware, familiar JS |
| `Result<T, E>` | Explicit control flow, Staff APIs |
| Discriminated union errors | Typed error channels without exceptions |

`Promise.all` vs `allSettled` — same semantics as JS ([JavaScript README §20](../JavaScript/README.md)).

---

## 15. Error handling at scale

```typescript
try {
  await service.submitVote(user, poll, option);
} catch (err: unknown) {
  if (err instanceof PollClosedError) { /* ... */ }
  else if (isDomainError(err)) { /* ... */ }
  else {
    logger.error('unexpected', { err });
    throw err;
  }
}
```

| Approach | Pros | Cons |
|----------|------|------|
| **Custom error classes** | `instanceof`, stack traces | Hierarchy sprawl |
| **Tagged errors** | Serializable, exhaustive | More boilerplate |
| **HTTP status only** | Simple | Loses type info inside app |

**Staff:** Define **bounded error unions** per use-case layer; never catch `any`.

---

## 16. Testing & mocking with interfaces

```typescript
interface VoteRepository {
  add(vote: Vote): void;
  getByPollId(pollId: PollId): Vote[];
}

class FakeVoteRepository implements VoteRepository {
  votes: Vote[] = [];
  add(vote: Vote): void { this.votes.push(vote); }
  getByPollId(pollId: PollId): Vote[] {
    return this.votes.filter((v) => v.pollId === pollId);
  }
}

// jest
const votes = new FakeVoteRepository();
const service = new PollingService(users, polls, votes);
```

| Tool | Role |
|------|------|
| **Fake** | Working in-memory impl |
| **Mock** | `jest.fn()` — assert calls |
| **Stub** | Fixed return — `getById: () => poll` |

**Interview line:** “Interfaces make DI testable — I inject fakes in unit tests, mocks only when verifying collaboration.”

---

## 17. Node / Express / NestJS typing

### Express (minimal)

```typescript
import type { Request, Response, NextFunction } from 'express';

interface AuthRequest extends Request {
  userId: UserId;
}

app.post('/polls', (req: Request, res: Response) => {
  const body = req.body as CreatePollInput; // prefer zod parse → typed
  res.json(poll);
});
```

### NestJS (Staff summary)

| Concept | TS angle |
|---------|----------|
| **Modules** | DI container boundaries |
| **Providers** | `implements` repository interfaces |
| **DTOs** | `class-validator` + `class-transformer` |
| **Decorators** | Experimental metadata — know `@Injectable()`, `@Controller()` |

**Interview line:** “Nest leverages decorators + DI; I'd still keep domain logic free of framework imports.”

### Runtime validation (production)

```typescript
import { z } from 'zod';

const CreatePollSchema = z.object({
  question: z.string().min(1),
  options: z.array(z.string()).min(2),
  isPrivate: z.boolean().default(false),
});

type CreatePollInput = z.infer<typeof CreatePollSchema>;
```

**Staff bar:** Types at compile time + schema at boundary — not one or the other.

---

## 18. Variance — covariance & contravariance

With `strictFunctionTypes`:

| Position | Arrays/objects | Function params | Function returns |
|----------|----------------|-----------------|------------------|
| **Covariant** | ✅ `Dog[]` ⊄ `Animal[]` (mutable array — unsound) | — | Return type can be subtype |
| **Contravariant** | — | Param can be **wider** | — |

```typescript
type Handler = (animal: Animal) => void;
let dogHandler: (dog: Dog) => void = (d) => console.log(d.bark());
let handler: Handler = dogHandler; // OK — param contravariant

// Return types covariant
type Getter = () => Dog;
let getAnimal: () => Animal = (): Dog => new Dog();
```

**Interview one-liner:** “Function args are contravariant (accept wider), returns covariant (return narrower).”

---

## 19. TypeScript for LLD interviews

### Whiteboard structure (45 min)

1. Clarify actors & rules (same as [../README.md](../README.md) §5)  
2. **Entities** — `class` or `type` + behaviour methods  
3. **Interfaces** — `Repository`, `Strategy`, `Notifier`  
4. **Service** — orchestration, authorization  
5. **Types** — `PollId`, discriminated results, DTOs  

### Map to repo (JavaScript → how you'd write TS)

| JS folder | TS upgrade |
|-----------|------------|
| [../JavaScript/RateLimiter2/](../JavaScript/RateLimiter2/) | `interface RateLimiterStrategy` |
| [../JavaScript/PollingSystem2/](../JavaScript/PollingSystem2/) | `PollingService` + `interface *Repository` |
| [../JavaScript/Splitwise/](../JavaScript/Splitwise/) | `abstract class Expense` + `implements` |
| [../JavaScript/PaymentGateway/](../JavaScript/PaymentGateway/) | `interface BankGateway` (see Go port) |

### Sample LLD snippet (Polling — Staff)

```typescript
interface PollingService {
  createPoll(creator: User, input: CreatePollInput): Poll;
  assignVoter(creator: User, pollId: PollId, voter: User): void;
  submitVote(voter: User, pollId: PollId, option: string): VoteResult;
  getStatistics(creator: User, pollId: PollId): PollStatistics;
}

type VoteResult =
  | { status: 'ok'; vote: Vote }
  | { status: 'not_assigned' }
  | { status: 'own_poll' }
  | { status: 'poll_closed' };
```

**Say aloud:** “Interfaces for ports, classes for domain, service for use-cases, branded IDs, union for outcomes.”

---

## 20. Classic gotchas — predict compile/runtime

### 1. `enum` vs union

```typescript
enum Status { Active, Inactive } // runtime object emitted
type Status = 'active' | 'inactive'; // erased — prefer for new code
```

### 2. `keyof` and index signature

```typescript
interface Scores { [userId: string]: number }
type K = keyof Scores; // string | number (not just string!)
```

### 3. Optional vs `undefined`

With `exactOptionalPropertyTypes`, `{ x?: number }` ≠ `{ x: number | undefined }`.

### 4. Array covariance (mutable)

```typescript
let dogs: Dog[] = [];
let animals: Animal[] = dogs; // allowed historically — unsound if you push Cat
```

### 5. `as` assertion lies

```typescript
const user = {} as User; // compiles; runtime disaster
```

### 6. `interface` merging surprise

```typescript
interface Box { height: number }
interface Box { width: number }
// Box = { height, width } — intentional in libs, confusing in app code
```

### 7. `void` vs `undefined` in callbacks

```typescript
type CB = () => void;
const cb: CB = () => undefined; // OK — void allows any return to be ignored
```

---

## 21. Staff / system-design talking points

| Topic | What to say |
|-------|-------------|
| **Adopt TS gradually** | `allowJs`, `checkJs`, strict island, expand coverage |
| **API versioning** | DTO types per version; don't leak domain entities |
| **Monorepo** | Project references, shared `domain` package |
| **Codegen** | OpenAPI → types; GraphQL codegen — single source |
| **Performance** | `tsc --build` incremental; avoid huge unions in hot paths |
| **any debt** | Measure `@ts-expect-error` count; burn down per squad |
| **LLD in TS** | Interfaces for OCP; generics for repositories; unions for state machines |

---

## 22. Practice checklist

Answer **without notes** (SDE-3 / Staff bar):

### Type system
- [ ] Explain structural typing vs Java nominal typing  
- [ ] When `interface` vs `type`? Declaration merging?  
- [ ] Implement discriminated union + exhaustive `never`  
- [ ] `unknown` vs `any` at HTTP boundary  
- [ ] Branded type for `UserId` / `PollId`  

### Generics & utilities
- [ ] Write `Repository<T, ID>` with constraint  
- [ ] Use `Pick` / `Omit` / `Partial` on a domain model  
- [ ] Explain `keyof` and `typeof` for config objects  

### LLD
- [ ] Whiteboard Rate Limiter with `interface Strategy`  
- [ ] Whiteboard Polling with service + repository interfaces  
- [ ] Explain how you'd test with fake repositories  

### Production
- [ ] Name 5 `strict` flags and why  
- [ ] Runtime validation strategy (zod) + compile-time types  
- [ ] Function variance one-liner  

### JS runtime (still required)
- [ ] Event loop, Promises — [JavaScript README](../JavaScript/README.md) §15–§20  

---

## 23. Rapid revision cheat sheet

```
TS = JS + erasable types; runtime needs zod/io-ts at boundaries

interface vs type     → interface: extend/merge/OOP contracts
                      → type: unions, mapped, utilities
structural typing     → shape matters, not declaration name
excess property       → fresh literal can't have extra keys

narrowing             → typeof, in, discriminant, type guards
discriminated union   → shared kind/tag field + switch + never

generics              → T, extends keyof, default type params
infer                 → unwrap Promise/return inside conditional type

utility types         → Partial Pick Omit Record Readonly ReturnType

LLD                   → interface Strategy/Repository; class entity
                      → service orchestrates; branded IDs; Result union

any vs unknown        → never any in new code; unknown + narrow

as const / satisfies  → literal unions; validate without widen
branded types         → PollId vs UserId on string

strict                → strictNullChecks, noUncheckedIndexedAccess
variance              → params contravariant, returns covariant

testing               → interface + fake impl; mock for call verify
boundaries            → zod parse → z.infer<typeof Schema>

JS runtime            → still JavaScript/README.md for event loop etc.
repo LLD              → JavaScript/* + design in TS with §19 patterns
```

---

*Last aligned with repo: `JavaScript/PollingSystem2/`, `JavaScript/RateLimiter2/`, [../README.md](../README.md) LLD method.*
