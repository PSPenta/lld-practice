# Ticket assign + notify — LLD walkthrough

> **Timed steps:** [Hub §4](../../README.md#4-how-a-typical-lld-round-runs) · **Solved:** ❌  
> Broader notification patterns: [logging-framework](../logging-framework/README.md) · [webhook-delivery](../webhook-delivery/README.md)

**Round opening (say aloud):**
> "I'll clarify requirements and v1 scope, outline entities and classes, walk the main flows, define APIs, then cover concurrency/failures, and how I'd evolve the design."

## Step 1 — Clarify

### Questions (ask 6–8)
1. Single assignee?
2. Notify channels: Slack, email, in-app?
3. Concurrent assign allowed?
4. Multi-tenant?
5. Template per event?
6. Reassign / unassign in v1?
7. Notify on create only or on every status change?
8. Sync notify in request path or async queue?

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

## Step 2 — Entities & classes

```text
Ticket {
  id, tenantId, subject, status
  assigneeId?, version
}

TicketService
  - Assign(ticketID, agentID) error
  - Get(ticketID)

TicketRepository
AssignmentPolicy
  - CanAssign(ticket, agent) → bool

Notifier (interface)
  - Notify(event TicketAssigned) error
  SlackNotifier | EmailNotifier | InAppNotifier

EventPublisher / Observer list of Notifiers
```

**Patterns:** Observer (fan-out notify) · Adapter (Slack API) · Repository · Optimistic lock

## Step 3 — Flows

**Happy path**
1. Load ticket (tenant-scoped)  
2. Policy: agent in same tenant/inbox?  
3. Optimistic lock: `UPDATE ... WHERE version=?` → 0 rows → 409 Conflict  
4. Save ticket, version++  
5. Publish `TicketAssigned` → each Notifier (Observer)  

**Edge cases**
1. Already assigned → conflict; do not notify twice for no-op  
2. Notifier failure → don’t roll back assign (async retry / outbox)

## Step 4 — APIs

```http
POST /tickets/{id}/assign   { agentId }
GET  /tickets/{id}
```

```text
TicketService.Assign(ticketID, agentID) error
Notifier.Notify(TicketAssigned)
```

## Step 5 — Deepen

- Optimistic lock prevents double-assign under concurrency  
- Observer: one event, many notifiers; Adapter wraps Slack/email APIs  
- Prefer outbox / queue so HTTP assign returns after persist, not after Slack  
- Idempotent assign of same agent = no-op; different agent = conflict or reassign policy  
- Multi-tenant: every query filtered by tenantId

## Step 6 — Evolve

- Reassign, SLA timers, template engine per event  
- Webhook ingest + connector sync (same async rule):

```text
Webhook inbound:
  Verify signature → store event_id (dedupe) → enqueue → return 200 fast
  Worker upserts ticket — never call LLM inline on webhook path

Connector sync:
  POST /sync → 202 { job_id } → worker fetch/ingest (retries, DLQ)
  Heavy I/O never blocks the HTTP handler
```

- Related: [webhook-delivery](../webhook-delivery/README.md), [message-queue](../message-queue/README.md)
