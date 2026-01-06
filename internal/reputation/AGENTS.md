# REPUTATION MANAGEMENT

**Generated:** 2026-01-06
**Commit:** $(git rev-parse HEAD 2>/dev/null | cut -c1-8 || echo "unknown")
**Branch:** $(git branch --show-current 2>/dev/null || echo "unknown")

## OVERVIEW
Automated sender reputation management with real-time scoring, circuit breakers, and adaptive sending.

## STRUCTURE
```
internal/reputation/
├── domain/              # Core models (events, scores, warmup, circuit breakers, DMARC)
├── repository/          # Data access interfaces + SQLite implementations
│   └── sqlite/          # Concrete repositories (events, scores, warmup, breakers, etc.)
├── service/             # Business logic (telemetry, limiter, auditor, DMARC, warmup)
├── database.go          # Database connection and initialization
└── scheduler.go         # Periodic reputation scoring and cleanup
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Track sending events | repository/sqlite/events_repository.go | RecordEvent() called on each email action |
| Calculate reputation scores | service/telemetry_service.go | CalculateReputationScore() - 90-day rolling window |
| Adaptive rate limiting | service/adaptive_limiter.go | GetLimit() adjusts limits based on score |
| Circuit breaker logic | service/circuit_breaker_service.go | Auto-pause on complaints >0.1%, bounces >10% |
| Warm-up schedules | service/warmup_service.go | Progressive volume ramping (100 → 80,000 msgs) |
| DMARC report analysis | service/dmarc_analyzer.go | Parse RUA reports, detect SPF/DKIM failures |
| External feedback | service/gmail_postmaster.go, microsoft_snds.go | Provider-specific telemetry ingestion |
| Audit DNS/deliverability | service/auditor_service.go | Checks SPF/DKIM/DMARC/rDNS/TLS/MTA-STS |

## CONVENTIONS
- **Event-driven architecture**: All reputation data flows through SendingEvent (sent/delivered/bounce/complaint/defer)
- **Repository pattern**: All data access through interfaces in repository/ package
- **Domain models separate**: Core types in domain/ package, independent of storage
- **Score-based decisions**: Reputation scores (0-100) drive rate limiting, circuit breakers, warm-up
- **90-day rolling window**: Reputation calculated from last 90 days of events
- **Threshold-driven triggers**: Circuit breakers activate on exceeding configured thresholds
- **Context everywhere**: All repository and service methods accept context for cancellation
- **Structured logging**: zap logger with domain, trigger type, thresholds for debugging

## ANTI-PATTERNS (THIS MODULE)
- NEVER bypass circuit breaker - always check reputation before sending
- NEVER use fixed rate limits - always call AdaptiveLimiter.GetLimit() for current limit
- NEVER ignore feedback - complaints/bounces must trigger reputation score updates
- NEVER hardcode provider thresholds - use configurable thresholds per provider
- NEVER skip warm-up for new domains - progressive ramping required for new IPs/domains
- NEVER calculate reputation on every send - use cached scores, refresh periodically via scheduler
- NEVER resume circuit breaker manually without investigation - auto-resume with exponential backoff
