# Binary Search

← [DSA index](README.md) · [docs/](../README.md) · [Hub](../../README.md)

Crisp revision. Add new BS problems here, a row in the index, and a **Practice** link.

---

## Index

| # | Topic | Practice |
|---:|--------|----------|
| 1 | [Safe mid](#safe-mid) | — |
| 2 | [Total occurrences in a sorted array](#total-occurrences-in-a-sorted-array) | [34](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/) |
| 3 | [Search in rotated sorted array](#search-in-rotated-sorted-array) | [33](https://leetcode.com/problems/search-in-rotated-sorted-array/) |
| 4 | [Lowest (minimum) in rotated sorted array](#lowest-minimum-in-rotated-sorted-array) | [153](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/) |
| 5 | [Duplicates (`start == mid == end`)](#duplicates-start--mid--end) | [81](https://leetcode.com/problems/search-in-rotated-sorted-array-ii/) · [154](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array-ii/) |
| 6 | [Is the array rotated-sorted? — O(n)](#is-the-array-rotated-sorted--on) | [1752](https://leetcode.com/problems/check-if-array-is-sorted-and-rotated/) |
| 7 | [Single element (pairs of 2, one unique)](#single-element-pairs-of-2-one-unique) | [540](https://leetcode.com/problems/single-element-in-a-sorted-array/) |
| 8 | [Find Peak Element](#find-peak-element) | [162](https://leetcode.com/problems/find-peak-element/) |
| 9 | [Square root of `n` (floor)](#square-root-of-n-floor--binary-search) | [69](https://leetcode.com/problems/sqrtx/) |
| 10 | [nth root of `m` (exact)](#nth-root-of-m-exact--binary-search) | [GFG](https://www.geeksforgeeks.org/problems/find-nth-root-of-m5843/1) |
| 11 | [Koko Eating Bananas](#koko-eating-bananas--binary-search-on-answer) | [875](https://leetcode.com/problems/koko-eating-bananas/) |
| 12 | [Minimum Days to Make m Bouquets](#minimum-days-to-make-m-bouquets--binary-search-on-answer) | [1482](https://leetcode.com/problems/minimum-number-of-days-to-make-m-bouquets/) |
| 13 | [Smallest Divisor Given a Threshold](#smallest-divisor-given-a-threshold--binary-search-on-answer) | [1283](https://leetcode.com/problems/find-the-smallest-divisor-given-a-threshold/) |
| 14 | [Capacity to Ship Packages Within D Days](#capacity-to-ship-packages-within-d-days--binary-search-on-answer) | [1011](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/) |
| 15 | [Kth Missing Positive Number](#kth-missing-positive-number) | [1539](https://leetcode.com/problems/kth-missing-positive-number/) |

1–8 classic / array. 9–10 search on numeric range. 11–14 **binary search on answer** (monotonic feasibility, return `start`). 15 missing-count on a sorted index. nth root has no LeetCode twin — use GFG.

---

## Safe mid

```text
// overflow-prone
mid = (low + high) / 2

// safe
mid = low + ((high - low) / 2)
```

## Total occurrences in a sorted array

**Practice:** [34. Find First and Last Position of Element in Sorted Array](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/) — count = last − first + 1.

Use **lower bound** + **upper bound**.

When `mid == target`:

| Bound | Move |
|-------|------|
| **Lower** (first index) | `end = mid` |
| **Upper** (last index) | `start = mid` |

```text
count = upper - lower + 1   // if found; else 0
```

## Search in rotated sorted array

**Practice:** [33. Search in Rotated Sorted Array](https://leetcode.com/problems/search-in-rotated-sorted-array/)

1. Find which half is **sorted**.
2. Check if target lies in that sorted half.
3. Search that half; else search the other.

## Lowest (minimum) in rotated sorted array

**Practice:** [153. Find Minimum in Rotated Sorted Array](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/)

```text
if mid > end  →  lowest is in the right half
else          →  lowest is in the left half (incl. mid)
```

- **Index of lowest** = **number of times** the array was rotated.
- Fully sorted → lowest at index `0` → rotated `0` times.

## Duplicates (`start == mid == end`)

**Practice:** [81. Search in Rotated Sorted Array II](https://leetcode.com/problems/search-in-rotated-sorted-array-ii/) · [154. Find Minimum in Rotated Sorted Array II](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array-ii/)

```text
if start == mid == end:
  start++
  end--
else:
  normal binary-search branch
```

## Is the array rotated-sorted? — O(n)

**Practice:** [1752. Check if Array Is Sorted and Rotated](https://leetcode.com/problems/check-if-array-is-sorted-and-rotated/)

Count **drops** (`discrepancy`) where `a[i] > next`.

**Next index (wrap) — equivalent:**

```text
next = (i == n - 1) ? 0 : i + 1
next = (i + 1) % n                 // prefer — same thing
```

### Linear (no wrap) — `i = 0 .. n-2`

```text
if a[i] > a[i + 1]: discrepancy++
```

| `discrepancy` | Meaning |
|---------------|---------|
| `0` | Sorted (0 rotations) |
| `1` | Rotated-sorted |
| `> 1` | **Not** rotated-sorted |

### Circular (with wrap) — `i = 0 .. n-1`

```text
if a[i] > a[(i + 1) % n]: discrepancy++
// or: a[i] > a[i == n - 1 ? 0 : i + 1]
```

| `discrepancy` | Meaning |
|---------------|---------|
| `1` | Valid circularly sorted (includes **fully sorted** — drop is last→first) |
| `≠ 1` | Invalid (usually `> 1`; all-equal edge cases aside) |

## Single element (pairs of 2, one unique)

**Practice:** [540. Single Element in a Sorted Array](https://leetcode.com/problems/single-element-in-a-sorted-array/)

All elements appear **twice** except one. After the single, pair parity flips.

| Side of single | Index pattern of each pair |
|----------------|----------------------------|
| **Left** | `(even, odd)` |
| **Right** | `(odd, even)` |

```text
start = 1
end   = n - 2          // 2nd-last

// check a[0] and a[n-1] separately
// BS on [1 .. n-2] using even/odd pair pattern to pick half
```

## Find Peak Element

**Practice:** [162. Find Peak Element](https://leetcode.com/problems/find-peak-element/)

Peak = index `i` where `a[i] > neighbors` (treat outsides as −∞).

```text
// edges first
if n == 1                    → 0
if a[0] > a[1]               → 0
if a[n - 1] > a[n - 2]       → n - 1

start = 1
end   = n - 2

mid = start + (end - start) / 2

if a[mid] > a[mid - 1] && a[mid] > a[mid + 1]
  → return mid                    // peak

if a[mid] > a[mid + 1]
  → end = mid - 1                 // peak on left (descending slope)
else
  → start = mid + 1               // peak on right (ascending slope)
```

## Square root of `n` (floor) — Binary Search

**Practice:** [69. Sqrt(x)](https://leetcode.com/problems/sqrtx/)

```text
// handle n == 0 or n == 1 → return n

start = 1
end   = ceil(n / 2)                        // tighter than end = n (√n ≤ n/2 for n ≥ 4)
// end = n  also works

while start <= end:
  mid = start + (end - start) / 2

  if mid * mid == n  → return mid          // perfect square
  if mid * mid > n   → end = mid - 1       // too big (not end = mid)
  else               → start = mid + 1     // mid * mid < n

return end                                 // floor(sqrt(n))
```

- Prefer `mid > n / mid` instead of `mid * mid` when overflow is a risk.
- `end = mid` with `return end` does **not** give floor for non-squares; use `end = mid - 1`.

## nth root of `m` (exact) — Binary Search

**Practice:** [GFG — Find Nth root of M](https://www.geeksforgeeks.org/problems/find-nth-root-of-m5843/1) (no LeetCode equivalent)

Find integer `x` such that `x^n == m`. If none → `-1`.

```text
// power(n, k) = k^n  — multiply k, n times
power(n, k):
  ans = 1
  for i in 1 .. n:
    ans *= k
  return ans

start = 1
end   = m

while start <= end:
  mid = start + (end - start) / 2
  val = power(n, mid)                 // mid^n

  if val == m  → return mid           // exact root
  if val > m   → end = mid - 1
  else         → start = mid + 1      // val < m

return -1                             // no exact integer nth root
```

- Same mid overflow care: stop early in `power` if `ans > m / k`.
- Only exact roots (unlike floor √n above).

## Koko Eating Bananas — Binary Search on answer

**Practice:** [875. Koko Eating Bananas](https://leetcode.com/problems/koko-eating-bananas/)

Given piles `array` and hours `h`. Find **min** eating speed `k` (bananas/hour) so all piles finish in `≤ h` hours.

```text
largest = max(array)

start = 1
end   = largest

totalTime(array, mid):                  // hours needed at speed mid
  totalHours = 0
  for x in array:
    totalHours += ceil(x / mid)         // Math.ceil(array[i] / mid)
  return totalHours

while start <= end:
  mid = start + (end - start) / 2

  if totalTime(array, mid) <= h:        // mid works — try slower
    end = mid - 1
  else:                                 // too slow — need faster
    start = mid + 1

return start                            // min feasible speed (low lands on answer)
```

- Search space = speed, not index. Feasibility is monotonic → BS.
- `ceil(x / mid)` per pile (one pile per hour slot; can’t split across parallel piles).

## Minimum Days to Make m Bouquets — Binary Search on answer

**Practice:** [1482. Minimum Number of Days to Make m Bouquets](https://leetcode.com/problems/minimum-number-of-days-to-make-m-bouquets/)

Unsorted `array[i]` = day the `i`th flower blooms. Need `m` bouquets, each of `k` **adjacent** blooms. One flower → one bouquet.

```text
if m * k > n  →  return -1              // not enough flowers

start = min(array)                      // earliest bloom  (not max)
end   = max(array)                      // latest bloom

canMake(mid):                           // how many bouquets if we wait `mid` days
  flowersCount = 0
  bouquets = 0
  for x in array:
    if x <= mid:                        // bloomed by day mid
      flowersCount++
      if flowersCount == k:             // k adjacent → one bouquet
        bouquets++
        flowersCount = 0                // flowers not reused
    else:
      flowersCount = 0                  // adjacency broken
  return bouquets >= m

while start <= end:
  mid = start + (end - start) / 2

  if canMake(mid):                      // mid works — try fewer days
    end = mid - 1
  else:
    start = mid + 1

return start                            // min feasible day (low lands on answer)
```

- Search space = **wait days**, not index. More days → never fewer bouquets → monotonic → BS.
- Adjacent only: a gap (`x > mid`) resets `flowersCount`.
- Same pattern as Koko: no extra `days` — feasible → `end = mid - 1`; else `start = mid + 1`; return `start`.

## Smallest Divisor Given a Threshold — Binary Search on answer

**Practice:** [1283. Find the Smallest Divisor Given a Threshold](https://leetcode.com/problems/find-the-smallest-divisor-given-a-threshold/)

Unsorted `array`. Find **min** divisor `d` such that `sum(ceil(array[i] / d)) ≤ threshold`.

```text
largest = max(array)

start = 1
end   = largest

sumDiv(array, mid):                     // result of dividing every element by mid
  total = 0
  for x in array:
    total += ceil(x / mid)
  return total

while start <= end:
  mid = start + (end - start) / 2

  if sumDiv(array, mid) <= threshold:   // mid works — try smaller divisor
    end = mid - 1
  else:                                 // sum too big — need larger divisor
    start = mid + 1

return start                            // min feasible divisor (low lands on answer)
```

- Search space = **divisor**, not index. Larger `d` → smaller (or equal) sum → monotonic → BS.
- Same pattern as Koko: `ceil(x / mid)` in a linear scan; feasible → `end = mid - 1`; else `start = mid + 1`; return `start`.

## Capacity to Ship Packages Within D Days — Binary Search on answer

**Practice:** [1011. Capacity To Ship Packages Within D Days](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/)

Packages `array` in order (cannot reorder). Find **min** ship capacity so all packages go out in `≤ D` days. One ship per day; a day’s load cannot exceed capacity.

```text
start = max(array)                      // at least the heaviest package
end   = sum(array)                      // one day, whole cargo

daysNeeded(array, mid):                 // days used at capacity mid
  totalWeight = 0
  totalDays = 1
  for x in array:
    if totalWeight + x > mid:           // doesn't fit today
      totalDays++
      totalWeight = x                   // start a new day with x
    else:
      totalWeight += x
  return totalDays

while start <= end:
  mid = start + (end - start) / 2

  if daysNeeded(array, mid) > D:        // too many days — need more capacity
    start = mid + 1
  else:                                 // mid works — try smaller capacity
    end = mid - 1

return start                            // min feasible capacity (low lands on answer)
```

- Search space = **capacity**, not index. Larger capacity → fewer (or equal) days → monotonic → BS.
- `start = max`, not `1`: a package heavier than capacity can never ship.
- Same pattern: too slow/too small → `start = mid + 1`; feasible → `end = mid - 1`; return `start`.

## Kth Missing Positive Number

**Practice:** [1539. Kth Missing Positive Number](https://leetcode.com/problems/kth-missing-positive-number/)

Sorted `arr` of unique positives. Some positives are missing. Return the `k`th missing positive.

```text
// Bruteforce — bump k past every present number that "eats" a missing slot
for x in arr:
  if x <= k: k++                        // x is present, so kth missing shifts right
  else: break                           // remaining values are all > k
return k


// Binary search — missing count before index mid
// arr[mid] should be (mid + 1) if nothing was missing (1-indexed values)
// missing = arr[mid] - (mid + 1)       // +1 skips 0-based index

start = 0
end   = n - 1

while start <= end:
  mid = start + (end - start) / 2
  missing = arr[mid] - (mid + 1)

  if missing < k:                       // not enough missing yet — go right
    start = mid + 1
  else:                                 // kth missing is at or before mid — go left
    end = mid - 1

return start + k                        // same as end + k + 1  (loop exits with start = end + 1)
```

- Brute is `x <= k` then `k++`, not `<`. If `arr` contains `k`, that value is not missing.
- BS branch is `missing < k → start = mid + 1`, not `>`. Too few missing → search right.
- After the loop, `end` is the rightmost index whose prefix still has `< k` missing. `start = end + 1`, so **`start + k` == `end + k + 1`**.
- Empty / all-missing-before-array: `start` stays `0` → `k` itself.
