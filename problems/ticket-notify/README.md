# Ticket assign + notify — LLD design walkthrough

> **Paper design** — helpdesk / shared-inbox classic (assign owner, notify team).  
> **Method:** [../../README.md §5](../../README.md) · **Full notification platform:** [../../lld-gaps/README.md §6](../../lld-gaps/README.md)

---

## Clarify

- Single assignee or multiple?
- Notify which channels (Slack, email, in-app)?
- Concurrent assign allowed?

---

## Classes

```text
TicketService
  - Assign(ticketID, agentID) error
  - Get(ticketID)

TicketRepository
Notifier (interface) → SlackNotifier | EmailNotifier | InAppNotifier
AssignmentPolicy → CanAssign(ticket, agent)
```

---

## Assign flow

```text
1. Load ticket (tenant-scoped)
2. Policy: agent in same tenant/inbox?
3. Optimistic lock: UPDATE ... WHERE version=?
   → 0 rows → 409 Conflict ("already assigned")
4. Save ticket, version++
5. Publish TicketAssigned event → Notifiers (Observer)
```

---

## Patterns

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
