# Research Notes: Plugin Runtime & Ops

## Decision: Heartbeat Cadence & Session Timeout
- **Rationale**: Aligns with distributed system norms (Redis Sentinel, Kubernetes liveness) where 45s timeout balances detection speed and network jitter across multi-tenant clusters.
- **Alternatives Considered**:
  - 10s total timeout (too aggressive; amplifies transient network jitter).
  - 60s+ timeout (too lax for stale sessions; delays failover).

## Decision: Tenant Quota Refill Window (5 minutes)
- **Rationale**: Mirrors existing PowerX rate-limit tiers; fits Marketplace SLA windows and keeps billing aggregation aligned with existing 5-minute buckets.
- **Alternatives Considered**:
  - 1-minute buckets (higher churn, risk of sudden oscillations).
  - 15-minute buckets (slower recovery for tenants after burst).

## Decision: Restart Backoff Strategy (Exponential up to 2 minutes)
- **Rationale**: Matches host auto-healing expectations and avoids CPU thrash in crash loops; exponential growth ensures quick retries initially, then cool-off.
- **Alternatives Considered**:
  - Fixed 10s interval (risk repeated thundering herd).
  - Linear 15s increments (slower to dampen repeated failures).

## Decision: Log Retention (7-day local window with archival)
- **Rationale**: Keeps local disk usage manageable while meeting observability retention policies; 7-day rotation matches Loki ingestion SLA and on-call triage window.
- **Alternatives Considered**:
  - 3-day retention (insufficient for weekly regression analysis).
  - 14-day retention (excessive disk use on constrained hosts).

## Decision: Marketplace Over-limit Notifications (Hourly Summary)
- **Rationale**: Prevents alert fatigue while giving operators timely insight; hourly cadence aligns with billing aggregation cycle.
- **Alternatives Considered**:
  - Per-event notifications (too noisy in spikes).
  - First-event only (insufficient for prolonged overages).
