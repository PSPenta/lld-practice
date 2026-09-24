# Binary Search

← [DSA index](README.md) · [docs/](../README.md) · [Hub](../../README.md)

Crisp revision. Add new BS problems here, a row in the index, and a **Practice** link.

---

## Index

| # | Topic | Practice |
|---:|--------|----------|
| 1 | [Safe mid](#safe-mid) | — |
| 2 | [Total occurrences in a sorted array](#total-occurrences-in-a-sorted-array) | [34](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/) |
| 3 | [Search in rotated sorted array I / II](#search-in-rotated-sorted-array-i--ii) | [33](https://leetcode.com/problems/search-in-rotated-sorted-array/) · [81](https://leetcode.com/problems/search-in-rotated-sorted-array-ii/) |
| 4 | [Lowest (minimum) in rotated sorted array I / II](#lowest-minimum-in-rotated-sorted-array-i--ii) | [153](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/) · [154](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array-ii/) |
| 5 | [Is the array rotated-sorted? — O(n)](#is-the-array-rotated-sorted--on) | [1752](https://leetcode.com/problems/check-if-array-is-sorted-and-rotated/) |
| 6 | [Single element in a sorted array](#single-element-in-a-sorted-array) | [540](https://leetcode.com/problems/single-element-in-a-sorted-array/) |
| 7 | [Find Peak Element / Mountain Peak Index](#find-peak-element--mountain-peak-index) | [162](https://leetcode.com/problems/find-peak-element/) · [852](https://leetcode.com/problems/peak-index-in-a-mountain-array/) |
| 8 | [Square root of `n` (floor)](#square-root-of-n-floor--binary-search) | [69](https://leetcode.com/problems/sqrtx/) |
| 9 | [nth root of `m` (exact)](#nth-root-of-m-exact--binary-search) | [GFG](https://www.geeksforgeeks.org/problems/find-nth-root-of-m5843/1) |
| 10 | [Koko Eating Bananas](#koko-eating-bananas--binary-search-on-answer) | [875](https://leetcode.com/problems/koko-eating-bananas/) |
| 11 | [Minimum Days to Make m Bouquets](#minimum-days-to-make-m-bouquets--binary-search-on-answer) | [1482](https://leetcode.com/problems/minimum-number-of-days-to-make-m-bouquets/) |
| 12 | [Smallest Divisor Given a Threshold](#smallest-divisor-given-a-threshold--binary-search-on-answer) | [1283](https://leetcode.com/problems/find-the-smallest-divisor-given-a-threshold/) |
| 13 | [Capacity to Ship Packages Within D Days](#capacity-to-ship-packages-within-d-days--binary-search-on-answer) | [1011](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/) |
| 14 | [Kth Missing Positive Number](#kth-missing-positive-number) | [1539](https://leetcode.com/problems/kth-missing-positive-number/) |
| 15 | [Aggressive Cows](#aggressive-cows) | [GFG](https://www.geeksforgeeks.org/problems/aggressive-cows/1) · [1552](https://leetcode.com/problems/magnetic-force-between-two-balls/) |
| 16 | [Search in 2D Matrix](#search-in-2d-matrix) | [74](https://leetcode.com/problems/search-a-2d-matrix/) |
| 17 | [Search in 2D Matrix II](#search-in-2d-matrix-ii) | [240](https://leetcode.com/problems/search-a-2d-matrix-ii/) |
| 18 | [Row with Maximum Ones](#row-with-maximum-ones) | [2643](https://leetcode.com/problems/row-with-maximum-ones/) |
| 19 | [Row with Maximum Ones (row-wise sorted)](#row-with-maximum-ones-row-wise-sorted) | [GFG](https://www.geeksforgeeks.org/problems/row-with-max-1s0023/1) |
| 20 | [Allocate Minimum Pages](#allocate-minimum-pages) | [GFG](https://www.geeksforgeeks.org/problems/allocate-minimum-number-of-pages0937/1) · [410](https://leetcode.com/problems/split-array-largest-sum/) · [painters](https://www.geeksforgeeks.org/problems/the-painters-partition-problem1535/1) |
| 21 | [Minimize Max Distance of Adjacent Gas Stations](#minimize-max-distance-of-adjacent-gas-stations) | [GFG](https://www.geeksforgeeks.org/problems/minimize-max-distance-to-gas-station/1) · [774](https://leetcode.com/problems/minimize-max-distance-to-gas-station/) |
| 22 | [Median of Two Sorted Arrays](#median-of-two-sorted-arrays) | [4](https://leetcode.com/problems/median-of-two-sorted-arrays/) · [GFG kth](https://www.geeksforgeeks.org/problems/k-th-element-of-two-sorted-array1317/1) |

1–7 classic / array. 8–9 search on numeric range. 10–13, **20–21** **binary search on answer** (minimise). 14 missing-count on a sorted index. 15 **maximise** min-distance (return `end`). 16–17 2D search (LC 74 / 240). **18** unsorted max-ones (nested loops). **19** row-wise sorted → lower-bound BS per row. 21 float BS (`ε`). **22** partition BS on cut (same family as **k-th of two sorted arrays**). nth root has no LeetCode twin — use GFG.

---

## Safe mid

```text
// overflow-prone
mid = (low + high) / 2

// safe
mid = low + ((high - low) / 2)
```

## Total occurrences in a sorted array

**Practice:** [34. Find First and Last Position of Element in Sorted Array](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/)

Two separate binary searches — **lower** (first index) and **upper** (one past last). Rest of the branches are normal BS.

```text
// Lower — first occurrence
start = 0, end = n - 1, lower = -1
while start <= end:
  mid = start + (end - start) / 2
  if nums[mid] == target:
    lower = mid
    end = mid - 1                       // keep looking left
  else if nums[mid] < target:
    start = mid + 1
  else:
    end = mid - 1

if lower < 0: return 0                  // not found — skip upper BS

// Upper — last occurrence (then exclusive end = last + 1)
start = 0, end = n - 1, upper = -1
while start <= end:
  mid = start + (end - start) / 2
  if nums[mid] == target:
    upper = mid
    start = mid + 1                     // keep looking right
  else if nums[mid] < target:
    start = mid + 1
  else:
    end = mid - 1

return (upper + 1) - lower              // same as upper - lower + 1 on inclusive indices
```

| Bound | When `nums[mid] == target` |
|-------|----------------------------|
| **Lower** (first) | `end = mid - 1` |
| **Upper** (last) | `start = mid + 1` |

- Inclusive last index → count = `upper - lower + 1`.
- If you treat `upper` as exclusive (`last + 1`), count = `upper - lower` (same number).

## Search in rotated sorted array I / II

**Practice:** [33. Search in Rotated Sorted Array](https://leetcode.com/problems/search-in-rotated-sorted-array/) · [81. Search in Rotated Sorted Array II](https://leetcode.com/problems/search-in-rotated-sorted-array-ii/)

Find which half is **sorted**, then check if `target` lies in that half.

```text
start = 0
end   = n - 1

while start <= end:
  mid = start + (end - start) / 2

  if nums[mid] == target:
    return mid                          // or true (II)

  // II only — cannot decide half when all three collide
  if nums[start] == nums[mid] && nums[mid] == nums[end]:
    start++
    end--
    continue

  if nums[start] < nums[mid]:           // left half [start..mid] is sorted
    if nums[start] <= target && target < nums[mid]:
      end = mid - 1
    else:
      start = mid + 1
  else:                                 // right half [mid..end] is sorted
    if nums[mid] < target && target <= nums[end]:
      start = mid + 1
    else:
      end = mid - 1

return -1                               // or false (II)
```

- **I** (unique): skip the `start == mid == end` block — values never all collide.
- **II** (duplicates): when `nums[start] == nums[mid] == nums[end]`, shrink both ends; worst case O(n).
- Inclusive ends: `nums[start] <= target` and `target <= nums[end]`. Strict `<` on both sides misses a hit on `start` / `end`.
- Left-sorted test is `nums[start] < nums[mid]` (strict). If `==` and not the triple case, treat as right-sorted / keep going via the duplicate shrink.

## Lowest (minimum) in rotated sorted array I / II

**Practice:** [153. Find Minimum in Rotated Sorted Array](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/) · [154. Find Minimum in Rotated Sorted Array II](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array-ii/)

Use **`start < end`** (not `<=`). Answer lands on `nums[start]`. Compare mid to **`nums[end]`**, not `nums[start]`.

```text
start = 0
end   = n - 1

while start < end:                      // not start <= end
  mid = start + (end - start) / 2

  if nums[mid] > nums[end]:             // min is strictly right of mid
    start = mid + 1
  else:                                 // min is at mid or left  (incl. equals)
    end = mid                           // keep mid — not mid - 1

return nums[start]
```

- `nums[mid] > nums[end]` → right half is unsorted / has the drop → go right (`start = mid + 1`).
- Else min is in `[start .. mid]` → `end = mid` (mid can still be the answer).
- **Index of lowest** = rotation count. Fully sorted → lowest at `0`.

### II — duplicates

Only when **all three** are equal you cannot pick a half — shrink both ends:

```text
if nums[start] == nums[mid] && nums[mid] == nums[end]:
  start++
  end--                                 // drop duplicates; worst case O(n)
else if nums[mid] > nums[end]:
  start = mid + 1
else:
  end = mid
```

Rest of the loop is the same (`start < end`, `return nums[start]`).

## Is the array rotated-sorted? — O(n)

**Practice:** [1752. Check if Array Is Sorted and Rotated](https://leetcode.com/problems/check-if-array-is-sorted-and-rotated/)

Linear pass — count **drops** (`discrepancy`) where `a[i] > next`.

```text
discrepancy = 0
for i in 0 .. n-1:                      // wrap last→first with %
  if a[i] > a[(i + 1) % n]:
    discrepancy++

// same without wrap: only i = 0 .. n-2  (misses last→first drop)
```

| `discrepancy` (no wrap, `i = 0 .. n-2`) | Meaning |
|----------------------------------------|---------|
| `0` | Simple sorted (0 rotations) |
| `1` | Rotated and sorted |
| `> 1` | Unsorted |

| `discrepancy` (with wrap, LC 1752) | Meaning |
|------------------------------------|---------|
| `≤ 1` | Valid (sorted or rotated-sorted; all-equal → `0`) |
| `> 1` | Unsorted |

- Prefer `(i + 1) % n` for the circular/LC check.
- No-wrap `0 / 1 / >1` is the revision one-liner; wrap treats fully sorted as `1` (last→first drop) unless all equal.

## Single element in a sorted array

**Practice:** [540. Single Element in a Sorted Array](https://leetcode.com/problems/single-element-in-a-sorted-array/)

Sorted array — every value appears in a **pair** except **one** unique. After the single, pair parity flips.

| Side of single | Index pattern of each pair |
|----------------|----------------------------|
| **Left** | `(even, odd)` |
| **Right** | `(odd, even)` |

```text
// edges — unique can sit at ends
if n == 1 or a[0] != a[1]:     return a[0]
if a[n - 1] != a[n - 2]:       return a[n - 1]

start = 1
end   = n - 2

while start <= end:
  mid = start + (end - start) / 2

  if a[mid] != a[mid - 1] && a[mid] != a[mid + 1]:
    return a[mid]                       // unique at mid

  if mid % 2 == 0:                      // even index
    if a[mid] == a[mid + 1]:            // pair intact → single is right
      start = mid + 1
    else:                               // pair broken → single is left (or mid)
      end = mid - 1
  else:                                 // odd index
    if a[mid] == a[mid - 1]:            // pair intact → single is right
      start = mid + 1
    else:
      end = mid - 1

return -1                               // shouldn't happen if input is valid
```

- Left of single: pairs start on **even** indices. Even mid matching `mid+1` (or odd mid matching `mid-1`) → still left → go right.
- Right of single: pair starts on **odd** → same checks fail → go left.

## Find Peak Element / Mountain Peak Index

**Practice:** [162. Find Peak Element](https://leetcode.com/problems/find-peak-element/) · [852. Peak Index in a Mountain Array](https://leetcode.com/problems/peak-index-in-a-mountain-array/) (same code — always return **index**)

Peak = index `i` where `a[i] > neighbors` (treat outsides as −∞). Mountain array = one peak (strict rise then fall); same BS.

```text
// edges first (162; mountain peak is never at ends for n ≥ 3)
if n == 1                    → 0
if a[0] > a[1]               → 0
if a[n - 1] > a[n - 2]       → n - 1

start = 1
end   = n - 2

while start <= end:
  mid = start + (end - start) / 2

  if a[mid] > a[mid - 1] && a[mid] > a[mid + 1]
    → return mid                    // peak index

  if a[mid] > a[mid + 1]
    → end = mid - 1                 // descending slope — peak on left
  else
    → start = mid + 1               // ascending slope — peak on right
```

- Return **index**, not `a[mid]`. 852 is the same algorithm on a guaranteed single mountain.

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

## Aggressive Cows

**Practice:** [GFG — Aggressive Cows](https://www.geeksforgeeks.org/problems/aggressive-cows/1) · same as [1552. Magnetic Force Between Two Balls](https://leetcode.com/problems/magnetic-force-between-two-balls/)

Place `k` cows in stalls `arr` so the **minimum** distance between any two cows is **maximised**.

```text
sort arr                               // unsorted stalls — must sort first

start = 1                              // min gap (distance, not arr[0] position)
end   = arr[n - 1] - arr[0]            // max gap = max - min  (not max - 1)

canPlace(dist):                        // greedy: first cow at arr[0]
  cows = 1
  prev = arr[0]
  for i in 1 .. n-1:
    if arr[i] - prev >= dist:
      cows++
      prev = arr[i]
      if cows >= k: return true
  return false

while start <= end:
  mid = start + (end - start) / 2

  if canPlace(mid):                    // dist works — try a larger gap
    start = mid + 1
  else:
    end = mid - 1

return end                             // largest feasible distance
```

- Search space = **distance**, not stall coordinate. `arr[0]` is a position; `start` is a gap.
- Sort first. Greedy left-to-right only works on sorted stalls.
- This **maximises**, so feasible → `start = mid + 1` and **return `end`** (Koko minimises and returns `start`).

## Search in 2D Matrix

**Practice:** [74. Search a 2D Matrix](https://leetcode.com/problems/search-a-2d-matrix/)

Each row is sorted. Each row’s first value is **greater** than the previous row’s last — treat as one sorted stream. Find the **row** that could hold `target`, then BS in that row.

```text
// Naive — linear find row, then linear scan the row
for r in 0 .. m-1:
  if matrix[r][0] <= target <= matrix[r][n - 1]:
    for x in matrix[r]:
      if x == target: return true
    return false
return false


// BS — row, then column
rStart = 0
rEnd   = m - 1
cStart = 0
cEnd   = n - 1                          // last col, not n
row    = -1

while rStart <= rEnd:
  rMid = rStart + (rEnd - rStart) / 2
  if matrix[rMid][cStart] <= target <= matrix[rMid][cEnd]:
    row = rMid
    break                               // this row can hold target
  else if target > matrix[rMid][cEnd]:
    rStart = rMid + 1                   // target is in a later row
  else:
    rEnd = rMid - 1                     // target < row start → earlier row

if row < 0: return false

while cStart <= cEnd:
  cMid = cStart + (cEnd - cStart) / 2
  if matrix[row][cMid] == target: return true
  if target > matrix[row][cMid]:
    cStart = cMid + 1
  else:
    cEnd = cMid - 1

return false
```

- Inclusive range: `<=` on both ends. Strict `<` misses `target` on a row’s first or last cell.
- `cEnd = n - 1`, not `n` — `matrix[r][n]` is out of bounds.
- Not [LC 240](#search-in-2d-matrix-ii): next row can start **below** the previous row’s last — staircase, not row-then-col.
- Flatten option: one BS on `[0, m*n)`. `row = mid / n`, `col = mid % n`. Same problem.

## Search in 2D Matrix II

**Practice:** [240. Search a 2D Matrix II](https://leetcode.com/problems/search-a-2d-matrix-ii/)

Each **row** sorted left→right, each **column** sorted top→bottom. A row’s last value **can be larger** than the next row’s first — not one global stream. Do **not** use LC 74 row-then-col or flatten.

Start at **row 0, last column** (`r = 0`, `c = n - 1`). Compare with `target`:

- `mat[r][c] == target` → found.
- `target` **larger** than `mat[r][c]` → this **row** can’t hold it (row is sorted left→right; everything left is smaller). Discard the row → `r++`.
- `target` **smaller** than `mat[r][c]` → this **column** can’t hold it (column sorted top→bottom; everything below is larger). The row might still have it → discard the column → `c--`.

**Do not shrink both row and col from `matrix[rMid][cMid]`.** If mid `<` target, it may still sit **right in a higher row** or **below in a left column**. Moving `rStart` **and** `cStart` drops a valid quadrant.

```text
// WRONG — 2D mid is not 1D BS
// matrix[rMid][cMid] < target → rStart = rMid+1 AND cStart = cMid+1
// loses target in the top-right or bottom-left


// Correct — staircase from top-right  O(m + n)
r = 0
c = n - 1                               // 0th row, last column

while r <= m - 1 && c >= 0:
  if matrix[r][c] == target: return true
  if target > matrix[r][c]:
    r++                                 // current row too small — go down
  else:
    c--                                 // current col too big — go left

return false
```

- Top-right is a **saddle**: left is smaller, down is larger. One comparison discards a row **or** a column.
- Same idea from bottom-left: `r = m-1`, `c = 0`; too big → `r--`, too small → `c++`.
- Naive fallback: BS each row `O(m log n)` — correct but slower. The 2D-mid shrink is not a valid BS.

## Row with Maximum Ones

**Practice:** [2643. Row With Maximum Ones](https://leetcode.com/problems/row-with-maximum-ones/)

Binary matrix, **unsorted**. Return `[rowIndex, onesCount]` for the row with the most `1`s (tie → smallest row index).

### Nested linear — O(m · n)

```text
bestRow = 0, bestCount = 0

for i in 0 .. m - 1:
  ones = 0
  for j in 0 .. n - 1:
    ones += mat[i][j]                  // 0/1 matrix
  if ones > bestCount:                 // strict > → keeps smallest i on tie
    bestCount = ones
    bestRow = i

return [bestRow, bestCount]
```

- No BS — rows/cols have no order to exploit.
- Update only on `ones > bestCount` so ties keep the earlier row.

## Row with Maximum Ones (row-wise sorted)

**Practice:** [GFG — Row with max 1s](https://www.geeksforgeeks.org/problems/row-with-max-1s0023/1)

Each row is sorted (`0`s then `1`s). Different rows can have different ones-counts. Return the row index with the most `1`s (often also track the count).

### Linear rows + BS first `1` — O(m · log n)

```text
bestRow = -1, bestCount = 0             // or 0 / 0 per problem statement

for i in 0 .. m - 1:
  // lower bound of 1 — first index where row[mid] == 1
  start = 0, end = n - 1, firstOne = n  // n = "not found" → 0 ones
  while start <= end:
    mid = start + (end - start) / 2
    if mat[i][mid] == 1:
      firstOne = mid
      end = mid - 1                     // keep looking left
    else:
      start = mid + 1                   // still in 0s

  ones = n - firstOne                   // NOT n - 1 - firstOne
  if ones > bestCount:
    bestCount = ones
    bestRow = i

return bestRow                          // or [bestRow, bestCount]
```

- Ones live in a suffix → count = `n - firstOne`. Off-by-one if you use `n - 1 - firstOne`.
- Row with no `1`s → `firstOne` stays `n` → `ones = 0`.
- Same lower-bound pattern as first occurrence / occurrences notes.

## Allocate Minimum Pages

**Practice:** [GFG — Allocate minimum number of pages](https://www.geeksforgeeks.org/problems/allocate-minimum-number-of-pages0937/1) · [410. Split Array Largest Sum](https://leetcode.com/problems/split-array-largest-sum/) · [GFG — Painter's Partition](https://www.geeksforgeeks.org/problems/the-painters-partition-problem1535/1) (same code)

Give `k` students contiguous books from `books[]` (pages). **Minimise** the max pages any student gets. Same skeleton as [ship packages](#capacity-to-ship-packages-within-d-days--binary-search-on-answer). Same as painter’s partition: books → boards, students → painters; if each unit takes `T` time, `return start * T`.

```text
if n < k: return -1                    // not enough books (each student gets ≥ 1)

start = max(books)                     // one student must take the thickest book
end   = sum(books)                     // one student takes everything

studentsNeeded(mid):                   // students required if max load = mid
  students = 1
  pages = 0
  for p in books:
    if pages + p > mid:                // doesn't fit this student
      students++
      pages = p
    else:
      pages += p
  return students

while start <= end:
  mid = start + (end - start) / 2

  if studentsNeeded(mid) > k:          // too many students — need a larger cap
    start = mid + 1
  else:                                // mid works — try a smaller max
    end = mid - 1

return start                           // min feasible max-pages
```

- Search space = **max pages per student**, not an index. `start = max`, not `1`.
- Fit check is `pages + p > mid` (or `<= mid` to stay). Strict `< mid` rejects a load that **equals** `mid` (e.g. `[10,20,30]`, `k=2`, `mid=30`).
- Count from `students = 1`. Then `students > k` is “infeasible”. Starting at `0` and testing `students < k` only works if you keep that pairing — don’t mix.
- Contiguous assignment: do **not** sort the books.
- Same as ship-capacity: packages → books, days → students, capacity → page cap.
- Painter’s partition is this problem; no separate write-up.

## Minimize Max Distance of Adjacent Gas Stations

**Practice:** [GFG — Minimize Max Distance to Gas Station](https://www.geeksforgeeks.org/problems/minimize-max-distance-to-gas-station/1) · [774. Minimize Max Distance to Gas Station](https://leetcode.com/problems/minimize-max-distance-to-gas-station/)

Sorted `stations` on a line. Add **`k`** new stations (anywhere, not only integers). **Minimise** the maximum adjacent distance.

### Naive — place one-by-one into the worst gap — O(k · n)

```text
placed[0 .. n-2] = 0                   // how many new stations inside each gap

for _ in 1 .. k:                       // place k stations
  maxDist = -1
  maxIndex = -1
  for i in 0 .. n-2:
    currDist = (stations[i+1] - stations[i]) / (placed[i] + 1)
    if currDist > maxDist:
      maxDist = currDist
      maxIndex = i
  placed[maxIndex]++                   // put next station in the worst gap

maxDist = 0
for i in 0 .. n-2:
  dist = (stations[i+1] - stations[i]) / (placed[i] + 1)
  maxDist = max(maxDist, dist)

return maxDist
```

- Correct greedy; each new station splits the current worst gap into equal pieces.
- Too slow when `k` is large → prefer BS.

### BS on answer (float) — O(n log(range / ε))

```text
start = 0
end   = 0
for i in 0 .. n-2:
  end = max(end, stations[i+1] - stations[i])   // max existing gap

while end - start > 1e-6:              // float — not start <= end
  mid = start + (end - start) / 2      // no floor; keep decimal

  stationsNeeded = 0
  for i in 0 .. n-2:
    gap = stations[i+1] - stations[i]
    if gap > mid:
      stationsNeeded += ceil(gap / mid) - 1   // +=  accumulate all gaps

  if stationsNeeded <= k:              // mid works — try smaller max
    end = mid
  else:                                // too tight — need larger max
    start = mid

return start                           // or end — they converge within ε
```

- Search space = **max adjacent distance** (real number), not an index.
- Feasible → `end = mid`; infeasible → `start = mid`. Do **not** use `mid ± 1`.
- `ceil(gap / mid) - 1` = new stations needed inside that gap so every piece ≤ `mid`.
- Use `+=`, not `=`. Guard `mid == 0` (or keep `ε` so mid never stays 0 forever).
- Same monotonic pattern as pages / ship: larger allowed distance → fewer stations needed.

<a id="median-of-two-sorted-arrays"></a>

## Median of Two Sorted Arrays

**Practice:** [4. Median of Two Sorted Arrays](https://leetcode.com/problems/median-of-two-sorted-arrays/) · same cut as [GFG — K-th element of two sorted arrays](https://www.geeksforgeeks.org/problems/k-th-element-of-two-sorted-array1317/1)

Two sorted arrays `arr1`, `arr2` (sizes may differ). Return the median of the merged sorted view.

### Naive — merge then pick — O(n1 + n2)

```text
arr = [], i1 = 0, i2 = 0

while i1 < n1 && i2 < n2:
  if arr1[i1] < arr2[i2]:
    arr.push(arr1[i1]); i1++
  else:
    arr.push(arr2[i2]); i2++

while i1 < n1: arr.push(arr1[i1]); i1++
while i2 < n2: arr.push(arr2[i2]); i2++

n = arr.length
if n % 2 == 1:                         // length odd/even — not mid % 2
  return arr[floor(n / 2)]
return (arr[n / 2 - 1] + arr[n / 2]) / 2
```

### BS — partition cut on the smaller array — O(log min(n1, n2))

BS on **how many elements to take from `arr1`**, not on the median value. Left side of the cut must hold `half` elements; left max <= right min.

```text
if n1 > n2: swap arr1, arr2            // always BS on shorter

start = 0
end   = n1                             // take 0 .. n1 from arr1
half  = (n1 + n2 + 1) / 2              // left size (odd -> median on left)

while start <= end:
  mid = start + (end - start) / 2      // i = take mid from arr1
  j   = half - mid                     // take j from arr2

  maxL1 = (mid == 0)  ? -INF : arr1[mid - 1]
  minR1 = (mid == n1) ? +INF : arr1[mid]
  maxL2 = (j == 0)    ? -INF : arr2[j - 1]
  minR2 = (j == n2)   ? +INF : arr2[j]

  if maxL1 <= minR2 && maxL2 <= minR1: // valid cut
    if (n1 + n2) % 2 == 1:
      return max(maxL1, maxL2)
    return (max(maxL1, maxL2) + min(minR1, minR2)) / 2

  if maxL1 > minR2:                    // too many from arr1
    end = mid - 1
  else:                                // too few from arr1
    start = mid + 1
```

- Empty side -> +/- INF so edge cuts work.
- Odd total -> median is max of left border. Even -> average of left max and right min.
- Same pattern as other BS: `start` / `end` / `mid`, branch left or right until the condition holds.

### Same problem family — K-th of two sorted arrays

**Yes.** Median is a special case of **k-th element in two sorted arrays**:

| Need | `k` (1-based in merged order) |
|------|-------------------------------|
| Odd median | `k = (n1+n2+1) / 2` |
| Even median | average of `k` and `k+1` where `k = (n1+n2)/2` |

Same BS cut: take `mid` from the shorter array so left has exactly `k` elements, enforce `maxL1 <= minR2` and `maxL2 <= minR1`, answer is `max(maxL1, maxL2)`.

Interview line: "Median = find the middle k-th (and maybe next) via the same partition BS."

