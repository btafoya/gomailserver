# Using RabbitMQ in a Go-Based Mail Server

**Audience:** CLI coding agents, autonomous refactoring tools, and human developers working on Go mail servers

---

## Purpose

This document explains **how and why RabbitMQ should be integrated into a Go-based mail server** to handle asynchronous, failure-prone, and bursty workloads while keeping SMTP and IMAP paths fast, reliable, and scalable.

RabbitMQ acts as a **durable internal work queue** for mail processing stages that should *not* block protocol handlers.

---

## High-Level Concept

In a mail server, not all work is equal:

* **Latency-critical:** SMTP accept, IMAP commands
* **Compute-heavy / failure-prone:** spam scanning, AV, outbound retries, indexing

RabbitMQ allows you to:

* Accept mail quickly
* Persist messages safely
* Defer expensive processing to background workers
* Retry safely without losing state

---

## Core Architectural Shift

### Without RabbitMQ

```
SMTP → scan → route → deliver → index
```

Problems:

* Tight coupling
* SMTP stalls under load
* Hard-to-scale retry logic

### With RabbitMQ

```
SMTP → persist → enqueue → ACK client
                 ↓
           worker pipelines
```

Benefits:

* Fast protocol response
* Horizontal scaling
* Explicit backpressure and retry control

---

## Primary RabbitMQ Use Cases

### 1. SMTP Ingest Decoupling

**Goal:** Keep SMTP fast and reliable.

Flow:

1. SMTP handler receives message
2. Assign internal `message_id`
3. Store raw message (disk / object store / DB)
4. Publish event: `mail.received`

Worker responsibilities:

* SPF / DKIM / DMARC evaluation
* Spam filtering
* Virus scanning
* Routing decision

Why RabbitMQ:

* Burst tolerance
* Failure isolation
* Easy horizontal scaling

---

### 2. Outbound Mail Delivery & Retry Queue

**Goal:** Reliable delivery with retry and backoff.

Flow:

* Publish `outbound.deliver`
* Worker attempts SMTP delivery
* On result:

  * `ACK` → success
  * Tempfail → republish with delay
  * Permfail → publish `mail.bounced`

Why RabbitMQ:

* Durable retry queues
* Controlled concurrency
* Avoids complex scheduler logic in SQL

---

### 3. IMAP Indexing & Metadata Processing

**Goal:** Fast IMAP operations without blocking delivery.

Triggered by:

* `mail.delivered`

Worker tasks:

* Mailbox counters
* Seen/unseen flags
* Threading metadata
* Full-text search indexing

Why RabbitMQ:

* Keeps IMAP latency low
* Indexing can lag without user impact

---

### 4. Security Reporting & Telemetry

Use cases:

* DMARC aggregate reports
* TLS-RPT
* Delivery analytics

Pattern:

* Emit events (`dmarc.result`, `smtp.delivery.result`)
* Aggregate asynchronously

Why RabbitMQ:

* Clean separation of mail plane vs analytics plane

---

### 5. Multi-Tenant Configuration Propagation

Use cases:

* Domain added/removed
* DKIM key rotation
* Alias updates

Pattern:

* Control plane publishes `config.changed`
* Mail nodes refresh caches

Why RabbitMQ:

* Reliable fan-out
* No polling required

---

## Suggested Exchange & Routing Model

### Exchanges

* `mail.events` (topic)

  * `mail.received`
  * `mail.delivered`
  * `mail.quarantined`
  * `mail.bounced`

* `outbound.jobs` (topic)

  * `outbound.deliver`
  * `outbound.retry`

* `security.jobs`

  * `scan.spam`
  * `scan.virus`

---

## Message Payload Guidelines

**DO:**

* Include message IDs
* Include storage pointers (path/object key)
* Include tenant/domain identifiers

**DO NOT:**

* Put full raw email bodies in the queue

Reason:

* Large payloads degrade broker performance
* Queues should move *work references*, not data blobs

---

## Reliability & Delivery Semantics

### Required Patterns

* Durable queues
* Persistent messages
* Manual ACK/NACK
* Dead-letter queues (DLQ)

### Idempotency Rule (Critical)

All workers **must be safe to run more than once**.

Examples:

* Check if `message_id` already processed
* Use database constraints or state flags

RabbitMQ provides **at-least-once delivery**, not exactly-once.

---

## Operational Controls for Agents

CLI agents should monitor:

* Queue depth
* Consumer lag
* DLQ growth
* Broker disk usage

Backpressure strategy:

* If queues grow → slow SMTP accept or issue tempfail

---

## Minimal High-Value Implementation (Recommended First Step)

1. SMTP accepts mail → persists → publishes `mail.received`
2. Processing worker:

   * Auth checks
   * Spam/AV
   * Routing decision
3. Emit:

   * `mail.delivered` OR
   * `outbound.deliver`
4. IMAP indexer consumes `mail.delivered`

This delivers immediate gains with minimal complexity.

---

## Summary for Coding Agents

RabbitMQ is best used as:

* A **durable internal work scheduler**
* A **retry-safe delivery engine**
* A **decoupling layer** between protocols and processing

If a task:

* Is slow
* Can fail
* Can be retried
* Does not need to block SMTP/IMAP

→ It belongs in RabbitMQ.

---

**End of document**
