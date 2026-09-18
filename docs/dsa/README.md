# DSA Notes

← [docs/](../README.md) · [Hub](../../README.md)

Crisp revision points.

---

## Binary Search

### Safe mid

```text
// overflow-prone
mid = (low + high) / 2

// safe
mid = low + ((high - low) / 2)
```

### Total occurrences in a sorted array

Use **lower bound** + **upper bound**.

When `mid == target`:

| Bound | Move |
|-------|------|
| **Lower** (first index) | `end = mid` |
| **Upper** (last index) | `start = mid` |

```text
count = upper - lower + 1   // if found; else 0
```

### Search in rotated sorted array

1. Find which half is **sorted**.
2. Check if target lies in that sorted half.
3. Search that half; else search the other.

### Lowest (minimum) in rotated sorted array

```text
if mid > end  →  lowest is in the right half
else          →  lowest is in the left half (incl. mid)
```

- **Index of lowest** = **number of times** the array was rotated.
- Fully sorted → lowest at index `0` → rotated `0` times.

### Duplicates (`start == mid == end`)

```text
if start == mid == end:
  start++
  end--
else:
  normal binary-search branch
```

### Is the array rotated-sorted? — O(n)

Count **drops** (`discrepancy`) where `a[i] > next`.

**Next index (wrap) — equivalent:**

```text
next = (i == n - 1) ? 0 : i + 1
next = (i + 1) % n                 // prefer — same thing
```

#### Linear (no wrap) — `i = 0 .. n-2`

```text
if a[i] > a[i + 1]: discrepancy++
```

| `discrepancy` | Meaning |
|---------------|---------|
| `0` | Sorted (0 rotations) |
| `1` | Rotated-sorted |
| `> 1` | **Not** rotated-sorted |

#### Circular (with wrap) — `i = 0 .. n-1`

```text
if a[i] > a[(i + 1) % n]: discrepancy++
// or: a[i] > a[i == n - 1 ? 0 : i + 1]
```

| `discrepancy` | Meaning |
|---------------|---------|
| `1` | Valid circularly sorted (includes **fully sorted** — drop is last→first) |
| `≠ 1` | Invalid (usually `> 1`; all-equal edge cases aside) |

### Single element (pairs of 2, one unique)

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

### Find Peak Element

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

### Square root of `n` (floor) — Binary Search

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

### nth root of `m` (exact) — Binary Search

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

### Koko Eating Bananas — Binary Search on answer

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

### Minimum Days to Make m Bouquets — Binary Search on answer

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
