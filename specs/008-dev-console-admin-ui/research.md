# Research Log — Dev Console & Admin UI

## Topic: Admin Audit & Configuration History Persistence
- **Decision**: Introduce dedicated tables `admin_console_audit_events` and `admin_console_config_changes` in schema `powerx_plugin_base`, modelled with `tenant_id`, `plugin_id`, `actor_id`, `actor_name`, `action`, `resource`, `diff_snapshot`, and indexed `occurred_at` to support filtering and export.
- **Rationale**: Existing `runtime_ops.RuntimeAuditEvent` lacks actor/context fields needed for compliance exports. A dedicated table lets us capture rich metadata (permission code, before/after summary) while aligning with repository patterns (`models.S(...)`) and maintaining RLS guarantees.
- **Alternatives considered**: Reuse `runtime_ops` audit table (insufficient columns); rely solely on external log aggregation (breaks offline export requirement); store changes as append-only JSON in configuration tables (hurts filtering and auditing fidelity).

## Topic: Job/Task History Surface
- **Decision**: Persist job execution snapshots in new table `admin_console_job_runs` populated by safe-op service callbacks, capturing `job_id`, `job_type`, `scope_type`, `scope_ref` (tenant/environment), `trigger_source`, `status`, `duration_ms`, and retry metadata.
- **Rationale**: No existing repository tracks admin-level job runs. Persisting snapshots enables consistent troubleshooting UI, retry eligibility checks, and RBAC-scoped history without querying live queues.
- **Alternatives considered**: Pull directly from task queue backend (would require privileged access and lacks historical durability); encode history in observability metrics (difficult to reconstruct per-run context); piggyback on runtime_ops quota records (does not represent arbitrary jobs).

## Topic: Troubleshooting Metrics & Webhook Diagnostics
- **Decision**: Expose a consolidated `/admin/dev-console/troubleshooting` endpoint that aggregates health checks (via existing runtime_ops health service), quota usage (runtime_ops quota service), and webhook delivery stats from `integration_webhook_attempts`, using STS-backed clients where remote calls are needed.
- **Rationale**: Reusing vetted services avoids duplicating integration logic; we can enrich results with timestamps and friendly labels for console display while respecting permission scopes.
- **Alternatives considered**: Query Prometheus directly from the console (adds infra coupling and auth complexity); build bespoke webhook metrics service (duplicates integration domain logic); rely purely on frontend to orchestrate multiple endpoints (harder to guard RBAC and consistency).
