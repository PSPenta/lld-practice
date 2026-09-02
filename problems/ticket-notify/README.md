# Ticket assign + notify — LLD walkthrough

> **Round pattern:** [Discussion 60 min · Machine coding 90–120 min](../../docs/method/README.md#4-how-a-typical-lld-round-runs) · [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Broader **notification platform** patterns: [logging-framework](../logging-framework/README.md) · [webhook-delivery](../webhook-delivery/README.md)

---

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single assignee?
2. Notify channels: Slack, email, in-app?
3. Concurrent assign allowed?
4. Multi-tenant?
5. Template per event?

### v1 expectations (state aloud)
| | |
|---|---|
| **Actors** | Agent, TicketService, Notifier implementations |
| **Use cases (v1)** | 1. Assign ticket 2. Notify assignee/team 3. Prevent double-assign |
| **In scope** | Assign + Observer notifiers, tenant scope |
| **Out of scope** | Full campaign builder, SMS provider billing |
| **Assumptions** | Single assignee; Slack + email; optimistic lock on assign |

### Confirm understanding
> "Agent assigns ticket; system persists and fan-out notifies on TicketAssigned."

---


---

## Step 2 — Entities & classes

```text
TicketService
  - Assign(ticketID, agentID) error
  - Get(ticketID)

TicketRepository
Notifier (interface) → SlackNotifier | EmailNotifier | InAppNotifier
AssignmentPolicy → CanAssign(ticket, agent)
```

---

## Step 3 — Flows

```text
1. Load ticket (tenant-scoped)
2. Policy: agent in same tenant/inbox?
3. Optimistic lock: UPDATE ... WHERE version=?
   → 0 rows → 409 Conflict ("already assigned")
4. Save ticket, version++
5. Publish TicketAssigned event → Notifiers (Observer)
```

---

## Step 5 — Deepen (patterns)

- **Observer:** one event, many notifiers  
- **Adapter:** Slack API wrapper  
- **Repository:** hide DB  

---

## Backup: Webhook ingest + connector sync (idempotent / async)

**Webhook (inbound events):**
```text
Verify signature → store event_id (dedupe) → enqueue → return 200 fast
Worker upserts ticket/message — never call LLM inline on webhook path
```

**Connector sync (pull files from Slack/Notion/Drive):**
```text
Bad:  POST /sync → fetch all files → chunk → embed in request (client waits)
Good: POST /sync → publish job to SQS/Kafka → 202 { job_id }
      Worker → fetch → ingest each doc → update last_synced_at (retries, DLQ)
```

Same rule: **heavy I/O and indexing never block the HTTP handler.**
