# Dynamic Programming

← [DSA index](README.md) · [docs/](../README.md) · [Hub](../../README.md)

Crisp revision. Add new DP problems here, a row in the index, and a **Practice** link.

---
 
## Index

| # | Topic | Practice |
|---:|--------|----------|
| 1 | [DP shortcut](#dp-shortcut) | — |
| 2 | [Memoization vs Tabulation](#memoization-vs-tabulation) | [509](https://leetcode.com/problems/fibonacci-number/) |
| 3 | [Climbing Stairs](#climbing-stairs) | [70](https://leetcode.com/problems/climbing-stairs/) |
| 4 | [Frog Jump (min energy)](#frog-jump-min-energy) | [GFG](https://www.geeksforgeeks.org/problems/geek-jump/1) |
| 5 | [Frog Jump with K distance](#frog-jump-with-k-distance) | [GFG](https://www.geeksforgeeks.org/problems/minimal-cost/1) |
| 6 | [House Robber (max sum non-adjacent)](#house-robber-max-sum-non-adjacent) | [198](https://leetcode.com/problems/house-robber/) |

---
 
## DP shortcut

Three steps before you write code:

```text
1. Represent the problem in terms of an index (or state).
2. Do all possible stuff on that index.
3. Combine:
     count ways  → sum of all stuffs
     minimise    → min(stuffs)
     maximise    → max(stuffs)
```

| Step | Meaning | Climbing Stairs |
|------|---------|-----------------|
| **1. Index** | What is `f(i)`? | `f(i)` = ways to reach stair `i` |
| **2. Stuff** | Transitions from `i` | Come from `i-1` (1-step) or `i-2` (2-step) |
| **3. Combine** | How answers merge | **Count ways** → `f(i-1) + f(i-2)` |

- State can be more than one index (`f(i, j)`, knapsack capacity, …) — still “represent as state.”
- “All possible stuff” = every legal choice at that state (pick/skip, 1-step/2-step, …).
- Same skeleton for memo or tabulation — only fill order differs ([§2](#memoization-vs-tabulation)).

**Interview line:** “Index → try all transitions → sum / min / max.”

---

## Memoization vs Tabulation

DP = store answers to overlapping subproblems so you never recompute them. Two ways to fill that store:

| | **Memoization** | **Tabulation** |
|--|-----------------|----------------|
| Direction | **Top-down** | **Bottom-up** |
| Shape | Recursion + cache | Iterative table / array |
| Start | Current / large state (`fib(5)`) | Base cases (`0`, `1`, …) |
| Flow | Call smaller states as needed | Build up until the answer index |
| When | Natural recursive recurrence; sparse states | Clear order of dependencies; want O(1) stack |

### Memoization (top-down)

Start from the value you need. Recurse into smaller subproblems; cache each result.

```text
// fib(5) → fib(4) + fib(3) → …
memo = {}                              // or array filled with sentinel

fib(n):
  if n <= 1: return n
  if n in memo: return memo[n]
  memo[n] = fib(n - 1) + fib(n - 2)
  return memo[n]
```

- First call is `fib(n)` (the “current” problem).
- Without `memo`, exponential recomputation; with it, each `n` once → `O(n)`.

### Tabulation (bottom-up)

Start from the bottom-most bases (`0` / `1`), gradually fill until the current answer.

```text
// dp[i] = fib(i)
dp = array of size n + 1
dp[0] = 0
dp[1] = 1                              // bases first

for i in 2 .. n:
  dp[i] = dp[i - 1] + dp[i - 2]

return dp[n]
```

- No call stack; order of the loop must respect dependencies (`i` needs `i-1`, `i-2`).
- Space can often shrink to a few variables when only recent states matter.

**Interview line:** Memo = top-down recursion + cache; Tabulation = bottom-up fill from bases. Same recurrence — different evaluation order.

---

## Climbing Stairs

**Practice:** [70. Climbing Stairs](https://leetcode.com/problems/climbing-stairs/)

`n` stairs; each step take **1** or **2**. Ways to reach the top = `ways(n - 1) + ways(n - 2)` (same recurrence as Fibonacci).

### Tabulation — bottom-up

Fill ways for stair `1`, then `2`, then up to `n`.

```text
if n <= 2: return n

dp = array of size n + 1               // sentinel fill optional
dp[1] = 1                              // one way: single 1-step
dp[2] = 2                              // 1+1 or 2

for i in 3 .. n:
  dp[i] = dp[i - 1] + dp[i - 2]        // last step was 1 or 2

return dp[n]
```

- Start from the bottom (`1`, `2`), not from `n` — classic tabulation.
- Guard `n <= 2` (or size the array carefully) so you never read unset `dp[2]` when `n == 1`.

### Space O(1) — rolling variables

Same tabulation; only the last two values matter.

```text
if n <= 2: return n                     // REQUIRED — else n=1 returns 2

prev2 = 1                               // ways(1)
prev  = 2                               // ways(2)

for i in 3 .. n:
  curr  = prev + prev2
  prev2 = prev
  prev  = curr

return prev                             // ways(n)
```

- Without `if n <= 2`, the loop never runs and you always return `2`.
- `curr = prev` before the loop is unused — drop it.

```text
if n <= 2: return n                     // REQUIRED — else n=1 returns 2

prev2 = 1                               // ways(1)
prev  = 2                               // ways(2)

for i in 3 .. n:
  curr  = prev + prev2
  prev2 = prev
  prev  = curr

return prev                             // ways(n)
```

- `prev2` / `prev` = last two answers; no `dp[]` array.
- Without `if n <= 2`, the loop never runs and you always return `2`.

### Memoization — top-down

Start from `n`, recurse to smaller stairs; cache each `dp[i]`.

```text
climbStairs(n):
  dp = array(n + 1).fill(-1)           // -1 = unknown
  return climb(n, dp)

climb(n, dp):
  if n <= 2: return n                  // NOT n <= 1 — else ways(2) becomes 1
  if dp[n] != -1: return dp[n]         // hit cache (check sentinel, not truthiness)
  dp[n] = climb(n - 1, dp) + climb(n - 2, dp)
  return dp[n]
```

**Bugs to avoid** (common JS mistakes):

| Mistake | Why it fails |
|---------|----------------|
| `if (n <= 1) return n` | `ways(2) = ways(1)+ways(0) = 1+0 = 1`, but answer is `2` |
| `if (!dp[k])` with `fill(-1)` | `-1` is truthy → never fills; or if `0` were valid, `!0` wrongly recomputes |
| Only write `dp[n-1]` / `dp[n-2]`, never `dp[n]` | Store the **current** answer at `dp[n]` |
| Return `dp[n-1] + dp[n-2]` without caching sum | Still cache `dp[n] = …` |

- Same recurrence as tabulation; evaluation order is top-down.
- Cache check must be **`dp[n] != -1`** (or `undefined`), not `!dp[n]`.

---

## Frog Jump (min energy)

**Practice:** [GFG — Geek Jump / Frog Jump](https://www.geeksforgeeks.org/problems/geek-jump/1)

Heights `height[0 .. n-1]`. From `i` jump to `i+1` or `i+2`. Cost of a jump = `|height[i] - height[j]|`. Return **min total energy** to reach `n-1`.

**Why `abs`:** energy is the **magnitude** of the height change. Jumping down (`height[j] > height[i]`) makes a raw subtract **negative** — that would wrongly shrink (or negate) total cost. `abs` keeps every jump cost ≥ 0 whether you go up or down.

Same skeleton as Climbing Stairs, but combine with **min** (shortcut step 3) and add jump cost.

### DP shortcut

| Step | Here |
|------|------|
| **1. Index** | `f(i)` = min energy to reach index `i` from `0` |
| **2. Stuff** | Arrive via `i-1` (1-jump) or `i-2` (2-jump) |
| **3. Combine** | **Minimise** → `min(left, right)` |

### Memoization — top-down

```text
frogJump(height):
  n = height.length
  dp = array(n).fill(-1)
  return f(n - 1, height, dp)          // ask for LAST index — not 0

f(i, height, dp):
  if i == 0: return 0                  // already at start
  if dp[i] != -1: return dp[i]

  left = f(i - 1, height, dp) + abs(height[i] - height[i - 1])   // abs: cost ≥ 0 up or down
  right = INF
  if i > 1:
    right = f(i - 2, height, dp) + abs(height[i] - height[i - 2])

  dp[i] = min(left, right)
  return dp[i]
```

### Tabulation — bottom-up

```text
dp[0] = 0
dp[1] = abs(height[1] - height[0])     // only 1-jump if n > 1

for i in 2 .. n - 1:
  one = dp[i - 1] + abs(height[i] - height[i - 1])
  two = dp[i - 2] + abs(height[i] - height[i - 2])
  dp[i] = min(one, two)

return dp[n - 1]
```

### Common bugs

| Issue | Fix |
|-------|-----|
| `return f(0, …)` | Always `0` — call **`f(n - 1)`** |
| Diff without `abs` | Downhill jump → negative “cost”; use **`abs(...)`** so energy is always ≥ 0 |
| `dp` size | `Array(n)` is enough (`0 .. n-1`) |

**Interview line:** “Stairs with cost — `f(i) = min(f(i-1)+|Δ|, f(i-2)+|Δ|)`.”

---

## Frog Jump with K distance

**Practice:** [GFG — Minimal Cost](https://www.geeksforgeeks.org/problems/minimal-cost/1)

Same as [Frog Jump (min energy)](#frog-jump-min-energy), but from `i` the frog may jump **any** `1 .. k` steps (not only 1 or 2).

### DP shortcut

| Step | Here |
|------|------|
| **1. Index** | `f(i)` = min energy to reach `i` |
| **2. Stuff** | Arrive from `i-1`, `i-2`, …, `i-k` (if valid) |
| **3. Combine** | **Minimise** → `min` over those jumps |

### Tabulation — O(n · k)

```text
dp[0] = 0

for i in 1 .. n - 1:
  minCost = INF
  for j in 1 .. k:                       // j = jump length — NOT from 0
    if i - j >= 0:
      minCost = min(minCost, dp[i - j] + abs(height[i] - height[i - j]))
  dp[i] = minCost

return dp[n - 1]
```

### Memoization — top-down

```text
f(i):
  if i == 0: return 0
  if dp[i] != -1: return dp[i]

  minCost = INF
  for j in 1 .. k:
    if i - j >= 0:
      minCost = min(minCost, f(i - j) + abs(height[i] - height[i - j]))
  dp[i] = minCost
  return dp[i]

// answer = f(n - 1)
```

### Bugs in the pasted code

```text
// WRONG pieces
dp[1] = abs(...)                        // optional; loop can set it
for i in 0 .. n-1:
  for j in 0 .. k:                      // j=0 is the killer
    cost = dp[i - j] + abs(...)
```

| Bug | Why |
|-----|-----|
| **`j` from `0`** | `j=0` → `dp[i] + 0` while `dp[i]` is still `-1` → `minCost` becomes `-1` for every `i` |
| Jump of 0 | Not a legal move — start **`j = 1`** |
| Loop `i` from `0` | Harmless if `j≥1`, but clearer to start **`i = 1`** and leave `dp[0] = 0` |

- `k = 2` reduces to classic frog jump (min energy).
- Time `O(n·k)`, space `O(n)` (or rolling window of size `k` if you want).

---

## House Robber (max sum non-adjacent)

**Practice:** [198. House Robber](https://leetcode.com/problems/house-robber/)

You are a professional robber planning to rob houses along a street. Each house has a certain amount of money stashed, the only constraint stopping you from robbing each of them is that adjacent houses have security systems connected and it will automatically contact the police if two adjacent houses were broken into on the same night.

Given an integer array `nums` representing the amount of money in each house, return the maximum amount of money you can rob tonight without alerting the police.

### DP shortcut

| Step | Here |
|------|------|
| **1. Index** | `f(i)` = max money that can be robbed from houses `[0..i]` |
| **2. Stuff** | At house `i`, you can either rob it (then you can't rob `i-1`) or skip it |
| **3. Combine** | **Maximise** → `max(rob current + f(i-2), skip current + f(i-1))` |

### Tabulation — bottom-up

```text
if n == 0: return 0
if n == 1: return nums[0]

dp = array of size n
dp[0] = nums[0]
dp[1] = max(nums[0], nums[1])

for i in 2 .. n-1:
  dp[i] = max(dp[i-1], nums[i] + dp[i-2])

return dp[n-1]
```

### Space O(1) — rolling variables

```text
if n == 0: return 0
if n == 1: return nums[0]

prev2 = nums[0]                    // dp[i-2]
prev  = max(nums[0], nums[1])      // dp[i-1]

for i in 2 .. n-1:
  curr  = max(prev, nums[i] + prev2)
  prev2 = prev
  prev  = curr

return prev
```

### Memoization — top-down

```text
houseRobber(nums):
  n = nums.length
  if n == 0: return 0
  memo = array(n).fill(-1)
  return rob(nums, n-1, memo)

rob(nums, i, memo):
  if i < 0: return 0
  if i == 0: return nums[0]
  if memo[i] != -1: return memo[i]
  
  // Rob current house + skip previous
  robCurrent = nums[i] + rob(nums, i-2, memo)
  // Skip current house
  robSkip = rob(nums, i-1, memo)
  
  memo[i] = max(robCurrent, robSkip)
  return memo[i]
```

### Common bugs to avoid

| Mistake | Why it fails |
|---------|--------------|
| `for i in 0 .. n-1` without base cases | Will access `dp[-1]` or `nums[-1]` |
| Forgetting `n == 0` and `n == 1` edge cases | Returns incorrect values for small arrays |
| Using `dp[i] = max(dp[i-1], nums[i] + dp[i-1])` | Doesn't skip adjacent house (should be `i-2`) |
| Not initializing `dp[1]` correctly | Should be `max(nums[0], nums[1])` not `nums[1]` |

**Interview line:** "House robber is a classic DP where `f(i) = max(f(i-1), nums[i] + f(i-2))` — either skip current house or rob it and skip the previous one."