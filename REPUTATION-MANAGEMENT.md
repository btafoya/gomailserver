# Automated Reputation Management for gomailserver Mail Server

This document describes what is required to **build, maintain, and automate sender reputation** for a self‑hosted mail server operating in cloud environments. It is written with the assumption that the reader is a mail‑server developer and operator, and that reputation management should be integrated directly into the server application itself rather than treated as an external afterthought.

The focus is on **engineering controls**, **feedback loops**, and **automated remediation**, suitable for an open‑source MTA such as `gomailserver`.

---

## Implementation Status

### ✅ Phase 1: Telemetry Foundation (COMPLETE)
**Status**: Production-ready
**Documentation**: `ISSUE1-PHASE1-COMPLETE.md`

- Reputation score calculation (0-100 scale)
- Event tracking (sent, delivered, bounce, complaint, defer)
- SQLite-based metrics storage
- Automated score calculation (every 5 minutes)
- Data retention policies (90-day rolling window)

### ✅ Phase 2: Deliverability Readiness Auditor (COMPLETE)
**Status**: Production-ready
**Documentation**: `ISSUE2-PHASE2-COMPLETE.md`

- DNS and authentication validation (SPF, DKIM, DMARC)
- rDNS and FCrDNS verification
- TLS certificate validation
- Operational mailbox checks (postmaster, abuse)
- RESTful API endpoints for reputation monitoring
- Real-time alert system

### ✅ Phase 3: Adaptive Sending Policy Engine (COMPLETE)
**Status**: Production-ready
**Documentation**: `ISSUE3-PHASE3-COMPLETE.md`

- Reputation-aware rate limiting (0-100 score → 0.0-1.0 multiplier)
- Circuit breaker with 3 trigger types:
  - High complaint rate (>0.1%)
  - High bounce rate (>10%)
  - Major provider blocks
- Auto-resume with exponential backoff (1h → 2h → 4h → 8h)
- Progressive warm-up (14-day schedule: 100 → 80,000 msgs/day)
- Auto-detection of new domains/IPs requiring warm-up
- SMTP integration with real-time enforcement
- Automated scheduler jobs:
  - Circuit breaker checks (every 15 minutes)
  - Auto-resume attempts (hourly)
  - Warm-up advancement (daily at midnight)
  - New domain detection (daily at 1 AM)

### 🚧 Phase 4: Dashboard UI (PLANNED)
- Real-time reputation visualization
- Circuit breaker status monitoring
- Warm-up progress tracking
- Manual override controls
- Alert configuration interface

### 🚧 Phase 5: Advanced Automation (PLANNED)
- DMARC report processing and visualization
- Automated complaint handling (ARF ingestion)
- Provider-specific rate limiting
- Custom warm-up schedules
- Machine learning for predictive scoring

---

## 1. Reputation Is a System, Not a Setting

Mail reputation is not controlled by a single configuration value or DNS record. It is an *emergent property* produced by the interaction of:

* Identity correctness (IP, domain, authentication)
* Consistent, predictable sending behavior
* Recipient feedback (complaints, bounces, engagement proxies)
* Infrastructure hygiene and protocol correctness

For a mail server, this naturally maps to a **closed‑loop control system**:

> **Measure → Classify → Adapt → Remediate**

Automated reputation management means embedding this loop directly into the MTA’s outbound pipeline.

---

## 2. Core Reputation Surfaces You Must Control

### 2.1 IP Reputation

Reputation is evaluated per sending IP and per receiving ecosystem. New IPs are untrusted by default and must instruct receivers through slow, predictable behavior. Sudden spikes, erratic sending, or persistent failures rapidly degrade trust.

Key factors:

* Volume stability
* Error rate
* Complaint rate
* Recipient validity

### 2.2 Domain Reputation

Modern filtering is increasingly domain‑centric. Domains used in:

* RFC5322.From
* DKIM `d=`
* SMTP MAIL FROM / Return‑Path

must remain stable, aligned, and authentic. Disposable or rotating domains are interpreted as evasive behavior.

### 2.3 Content and Behavioral Signals

Mailbox providers infer whether mail is wanted using indirect signals:

* User spam complaints
* Delete‑without‑read behavior
* Low engagement across cohorts

Your server cannot see engagement directly, but it *can* detect complaint feedback and bounce semantics that correlate strongly with poor engagement.

### 2.4 Infrastructure Trust Signals

Poor SMTP hygiene is indistinguishable from low‑effort spam infrastructure. This includes:

* Broken or missing reverse DNS
* Inconsistent HELO/EHLO identity
* Weak or missing TLS
* Protocol violations

---

## 3. Cloud Baseline Requirements (Non‑Negotiable)

### 3.1 Stable Identity and Reverse DNS

Every sending IP must have:

* A PTR record pointing to a hostname
* A/AAAA records resolving that hostname back to the same IP

This forward‑confirmed reverse DNS (FCrDNS) is table‑stakes for inbox delivery.

### 3.2 Standards‑Compliant SMTP Behavior

The server must strictly adhere to SMTP norms:

* Correct HELO/EHLO identity
* Proper response codes
* Predictable retry semantics

Operational addresses such as `postmaster@` (and preferably `abuse@`) must exist and be deliverable.

### 3.3 TLS Everywhere

At minimum:

* Opportunistic TLS for all SMTP sessions

Recommended:

* MTA‑STS (policy + HTTPS endpoint)
* TLS reporting (TLS‑RPT)
* DANE where DNSSEC is available

TLS failures are increasingly treated as delivery failures rather than soft warnings.

---

## 4. Authentication and Alignment

Authentication is no longer optional for real‑world delivery.

### Required

* **SPF** (valid envelope sender authorization)
* **DKIM** (cryptographic message signing)
* **DMARC** with reporting enabled

Alignment matters. SPF and DKIM must align with the visible From domain or DMARC will fail. DMARC reporting is not only defensive; it is a critical observability channel for reputation health.

Your server should:

* Always DKIM‑sign authenticated outbound mail
* Support per‑domain keys and rotation
* Provide clear diagnostics for alignment failures

---

## 5. Sending Behavior That Preserves Reputation

### 5.1 Complaint Rate Control

Complaint rate is one of the strongest negative signals. Even fractions of a percent can cause filtering. Complaints must result in **immediate suppression** of the complaining recipient.

### 5.2 Unsubscribe Mechanics

Bulk or recurring mail must support:

* RFC‑compliant `List‑Unsubscribe`
* One‑click unsubscribe where possible

Failure here directly increases complaints.

### 5.3 Bounce and Recipient Hygiene

Your server must rapidly suppress:

* Non‑existent recipients
* Repeated permanent failures

Hard bounces must not be retried beyond minimal confirmation.

### 5.4 Warm‑Up and Traffic Shaping

New IPs and domains must ramp slowly:

* Low initial volume
* Gradual increases
* Stable sending windows

Volume cliffs are interpreted as abuse or compromise.

---

## 6. Observability: Measuring Reputation in Practice

### 6.1 First‑Party Telemetry

Your MTA should emit structured events for:

* Delivery attempts
* Deferrals (with enhanced status codes)
* Hard vs soft bounces
* Policy blocks
* Complaints

Metrics should be aggregated by:

* Tenant / account
* Sending domain
* Sending IP
* Recipient domain

### 6.2 External Feedback Sources

Reputation cannot be inferred from SMTP alone. Integration with mailbox‑provider feedback systems is essential:

* Gmail Postmaster metrics
* Microsoft SNDS / complaint reporting
* Abuse Reporting Format (ARF) ingestion

These signals close the loop between sending behavior and recipient perception.

---

## 7. Designing Automated Reputation Management Inside the MTA

### 7.1 Deliverability Readiness Auditor ✅ IMPLEMENTED

**Status**: Production-ready (Phase 2 complete)
**Documentation**: `ISSUE2-PHASE2-COMPLETE.md`

A built‑in auditor prevents broken configurations from damaging reputation.

Implemented checks:

* SPF / DKIM / DMARC presence and alignment
* rDNS and FCrDNS validation
* TLS posture and certificate validity
* MTA-STS policy validation
* Required operational mailboxes (postmaster@, abuse@)

Exposed via:

* RESTful API endpoints (`/api/v1/reputation/audit/:domain`)
* Concurrent DNS/TLS validation for sub-second audit performance
* Overall deliverability score (0-100) with detailed check results
* Real-time alert generation for critical issues

### 7.2 Reputation Telemetry Pipeline ✅ IMPLEMENTED

**Status**: Production-ready (Phase 1 complete)
**Documentation**: `ISSUE1-PHASE1-COMPLETE.md`

All outbound mail generates normalized events feeding a SQLite metrics store. Rolling windows (24h, 7d, 30d) allow detection of trend changes before blocks occur.

Implemented features:

* Event tracking: sent, delivered, bounce, complaint, defer
* Automated reputation score calculation (0-100 scale)
* Metrics aggregation by domain, IP, and time window
* Scheduled score recalculation (every 5 minutes)
* Data retention policies (90-day rolling cleanup)
* RESTful API for score retrieval and monitoring

### 7.3 Adaptive Sending Policy Engine ✅ IMPLEMENTED

**Status**: Production-ready (Phase 3 complete)
**Documentation**: `ISSUE3-PHASE3-COMPLETE.md`

This is the heart of automated reputation management.

The engine dynamically adjusts:

* Per‑domain rate limits based on reputation score (0-100 → 0.0-1.0 multiplier)
* Circuit breaker enforcement with automatic pausing
* Progressive warm-up volume caps for new domains/IPs

Implemented **circuit breakers** with:

* Automatic pause when complaints spike (>0.1% threshold)
* Automatic pause when bounces spike (>10% threshold)
* Major provider block detection and response
* Exponential backoff auto-resume (1h → 2h → 4h → 8h max)
* Manual operator override capability

**Integration**:
- Real-time SMTP enforcement (421 error codes when paused/limited)
- Automated scheduler jobs (circuit checks every 15min, auto-resume hourly)
- Daily warm-up progression at midnight
- Daily new domain detection at 1 AM

### 7.4 Automated Complaint Handling

Feedback loop messages must:

* Be parsed automatically
* Mapped to recipients or accounts
* Trigger immediate suppression

Failure to suppress complainers guarantees escalating reputation damage.

### 7.5 Guardrails for Multi‑Tenant Safety

If the server supports multiple users or tenants, reputation isolation is mandatory:

* Per‑account rate limits
* Anomaly detection (credential abuse)
* Mandatory warm‑up for new tenants
* Hard ceilings on outbound volume

Without these controls, one compromised account will poison the entire IP and domain set.

### 7.6 Automated Inbound DKIM monitoring and actions

Present the DMARC reports in the Admin WebUI (new page) for the administrators to review as a summary and details.

Automatically take actions to address issues being reported back.

---

## 8. What a Successful UX Looks Like

A reputation‑aware mail server should expose:

1. **Setup Wizard** – live validation of DNS, auth, TLS, and identity
2. **Deliverability Dashboard** – acceptance rates, deferrals, complaints, top failures
3. **Actionable Alerts** – early warnings before blocking occurs
4. **Automated Remediation Logs** – every throttle, pause, or suppression recorded
5. **Operator Overrides** – explicit, auditable controls

---

## 9. Key Takeaway

Reputation is not something you configure once; it is something you **operate continuously**.

By embedding observability, adaptive policy, and automated remediation directly into the mail server, you turn reputation from an opaque external judgment into a **managed engineering system**.

This approach is particularly well‑suited to modern open‑source MTAs that aim to replace fragile multi‑component stacks with a single, coherent control plane.
