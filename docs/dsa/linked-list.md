# Linked List

← [DSA index](README.md) · [docs/](../README.md) · [Hub](../../README.md)

Crisp revision. Add new LL problems here, a row in the index, and a **Practice** link.

---

## Index

| # | Topic | Practice |
|---:|--------|----------|
| 1 | [Reverse a linked list](#reverse-a-linked-list) | [206](https://leetcode.com/problems/reverse-linked-list/) |
| 2 | [Fast & slow pointers](#fast--slow-pointers) | [141](https://leetcode.com/problems/linked-list-cycle/) · [142](https://leetcode.com/problems/linked-list-cycle-ii/) · [876](https://leetcode.com/problems/middle-of-the-linked-list/) |
| 3 | [Merge two sorted lists](#merge-two-sorted-lists) | [21](https://leetcode.com/problems/merge-two-sorted-lists/) |
| 4 | [Merge K sorted lists](#merge-k-sorted-lists) | [23](https://leetcode.com/problems/merge-k-sorted-lists/) |

---

## Reverse a linked list

**Practice:** [206. Reverse Linked List](https://leetcode.com/problems/reverse-linked-list/)

Flip every `next` pointer. New head is the old tail.

### Iterative — O(n) time, O(1) space

```text
prev = null

while head:
  next = head.next                     // save rest
  head.next = prev                     // reverse link
  prev = head                          // walk prev forward
  head = next                          // walk head forward

return prev                            // old tail is new head
```

- Three pointers: `prev` (reversed so far), `head` (current), `next` (rest).
- Empty list / single node → loop never runs or runs once → `prev` is correct.
- Do **not** return `head` — after the loop `head` is `null`.

### Recursive — O(n) time, O(n) stack

Same idea as iterative: reverse link, then recurse on the rest. Carry `prev`.

```text
reverse(head, prev = null):
  if head == null:
    return prev                        // empty → new head is prev (old tail)

  next = head.next
  head.next = prev
  return reverse(next, head)           // walk forward with updated prev
```

Classic post-order form (no `prev` arg) — reverse the suffix first, then attach current:

```text
reverse(head):
  if head == null || head.next == null:
    return head                        // empty / single → already reversed

  newHead = reverse(head.next)         // reverse rest; newHead = old tail
  head.next.next = head                // old next now points back to head
  head.next = null                     // head becomes new tail
  return newHead
```

- Base case with `prev`: return **`prev`**, not `head` (`head` is `null` at the end).
- Must **`return`** the recursive result — assigning to `next` and dropping it loses the new head.
- Pass **`head` as the new `prev`** into the recursive call (same order as the iterative loop). Mutating after the call with the old `prev` does not reverse.
- Stack depth = `n`. Prefer iterative in interviews unless asked for recursion.

---

## Fast & slow pointers

**Practice:** [141. Linked List Cycle](https://leetcode.com/problems/linked-list-cycle/) · [142. Linked List Cycle II](https://leetcode.com/problems/linked-list-cycle-ii/) · [876. Middle of the Linked List](https://leetcode.com/problems/middle-of-the-linked-list/)

`slow` moves one step, `fast` moves two. By the time `fast` reaches the end, `slow` is at the mid. On a cyclic list the pointers meet instead — used to detect a cycle and find its entry.

```text
slow = fast = head

while fast != null and fast.next != null:
  slow = slow.next
  fast = fast.next.next
```

- Guard `fast != null && fast.next != null` **before** using `fast.next.next` (avoids null-deref / TypeError on a terminating list).
- Always compare by **reference**: `slow === fast`.
- On a terminating list, loop ends with `slow` at the mid.
- On a cyclic list, loop never terminates; check `slow == fast` inside to detect the cycle.

### Find the middle — O(n) time, O(1) space — 876

Return `slow`. Halfway condition decides left vs right middle on even length.

```text
slow = fast = head

while fast != null and fast.next != null:
  slow = slow.next
  fast = fast.next.next

return slow
```

- `while fast != null and fast.next != null` → `slow` ends at the **right middle** (index `n/2` for even `n`).
- Use `while fast.next != null and fast.next.next != null` if you want the **left middle**.
- Use cases: split list for merge sort, palindrome check (reverse first half, compare).

### Detect a cycle — O(n) time, O(1) space — 141

```text
slow = fast = head

while fast != null and fast.next != null:
  slow = slow.next
  fast = fast.next.next
  if slow == fast:
    return true           // pointers met → cycle exists

return false
```

- If a cycle exists, both pointers get trapped in it; `fast` closes the gap by one per step, so capture is guaranteed.
- `fast` is checked before every double-step so a terminating list does not deref `null`.
- Compare **`slow === fast`** (reference), never by value.

### Find cycle entry — O(n) time, O(1) space — 142

Once `slow` and `fast` meet, reset one pointer to `head`; advance both one step per turn. They meet at the cycle's start.

```text
// Phase 1 — find meeting point (as above)
slow = fast = head
found = false
while fast != null and fast.next != null:
  slow = slow.next
  fast = fast.next.next
  if slow == fast:
    found = true
    break
if not found: return null

// Phase 2 — find entry
ptr = head
while ptr != slow:
  ptr = ptr.next
  slow = slow.next

return ptr
```

- Let head→entry = `a`, entry→meeting = `b`, cycle length = `c`. At meeting, `slow = a + b`, `fast = a + b + k*c`. With `fast = 2*slow`: `a = k*c - b`, so the pointer from head advances exactly `a` and lands on the entry.
- If no cycle was found in Phase 1, return `null`.

---

## Merge two sorted lists

**Practice:** [21. Merge Two Sorted Lists](https://leetcode.com/problems/merge-two-sorted-lists/)

Walk both lists from their heads; at each step pick the smaller current value, append it to a new list, and advance only that list's pointer. Once one list is fully drained, the other may still hold (larger) values — append its remainder as-is.

### Iterative (dummy node) — O(n + m) time, O(1) extra

Use a `dummy` head and a `tail` cursor so the very first append doesn't need a special case.

```text
dummy = {next: null}
tail  = dummy

while l1 != null and l2 != null:
  if l1.val <= l2.val:
    tail.next = l1
    l1 = l1.next
  else:
    tail.next = l2
    l2 = l2.next
  tail = tail.next

// one list drained; the other may still have values left
if l1 != null: tail.next = l1
if l2 != null: tail.next = l2

return dummy.next
```

- Reuses the existing nodes (no new nodes), so extra space is O(1).
- `<=` vs `<` only affects which list's node is chosen on ties; picking from `l1` keeps it stable.
- The leftover check (`tail.next = l1 || l2`) is required because the loop exits as soon as **either** pointer becomes null.

### Recursive — O(n + m) time, O(n + m) stack

```text
mergeTwoLists(l1, l2):
  if l1 == null: return l2
  if l2 == null: return l1

  if l1.val <= l2.val:
    l1.next = mergeTwoLists(l1.next, l2)
    return l1
  else:
    l2.next = mergeTwoLists(l1, l2.next)
    return l2
```

- Base cases are the leftover rule: an empty list contributes nothing.
- Prefer iterative in interviews (no stack cost).

---

## Merge K sorted lists

**Practice:** [23. Merge k Sorted Lists](https://leetcode.com/problems/merge-k-sorted-lists/)

Given `k` sorted linked lists, merge them into one sorted list. Reuse [`mergeTwoLists`](#merge-two-sorted-lists) as the primitive: pop two lists, merge them, and put the merged result back so it gets merged with the next one. Repeat until fewer than two lists remain.

### Pairwise with a working list — O(k^2 · n) worst case

```text
function mergeKLists(lists):
  if lists == null or lists.length == 0: return null

  while lists.length >= 2:
    l1 = lists.shift()                 // take first
    l2 = lists.shift()                 // take second (guaranteed here)
    merged = mergeTwoLists(l1, l2)
    lists.unshift(merged)              // merged joins the front for the next pass

  // 0 lists → null already handled; 1 list → it is the answer
  return lists.length == 1 ? lists[0] : null
```

- Each pass reduces the count by 1 (`k - 1` merges total), so `mergeTwoLists` runs `k - 1` times.
- On an odd count, the last unpaired list is never touched until later passes pair it with a merged result.
- `shift`/`unshift` on a JS array is `O(k)` per call — for large `k`, prefer a queue (deque) to get `O(1)` pop/push, or divide-and-conquer.

### Pitfalls

| Mistake | Why it fails |
|---------|--------------|
| Not handling `lists.length == 0` / empty array | Returns `undefined` instead of `null` |
| Popping only once when `length == 1` | Infinite/incorrect handling; check `length < 2` before popping a second |
| Re-pairing the same two lists | Ensure the merged list re-enters the pool, not that `l1`/`l2` stay around |
| Ignoring a leftover single list at the end | After the loop, return `lists[0]` (the sole survivor) |
| Mutating the input array | `shift`/`unshift` destroy caller's order — copy first if needed |

### Better asymptotics (for large `k`)

- **Divide & conquer:** pair lists and merge level by level → `O(N log k)` time, `O(1)` extra.
- **Min-heap:** keep heads in a priority queue; pop smallest, push its `next` → `O(N log k)` time, `O(k)` space.
- Use `mergeTwoLists` unchanged as the merge primitive either way; only the *scheduling* of merges changes.

**Interview line:** "K sorted lists = repeated 2-way merge. Min-heap or divide-and-conquer gets `O(N log k)`; naive front-of-array pairing is `O(k^2 · n)`."
