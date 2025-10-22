# Data Model — Dev Console & Admin UI

## Entities

### AdminConsoleAuditEvent
- **Purpose**: Canonical audit log for privileged admin actions performed through the console.
- **Fields**:
  - `id UUID PK`
  - `tenant_id TEXT NOT NULL` (RLS scope; nullable for global plugin actions)
  - `plugin_id TEXT NOT NULL`
  - `actor_id TEXT NOT NULL`
  - `actor_name TEXT`
  - `actor_email TEXT`
  - `permission_code TEXT NOT NULL`
  - `action TEXT NOT NULL` (enum: `CONFIG_UPDATE`, `SAFE_OP_REPLAY`, `SAFE_OP_RETRY`, `SAFE_OP_DRAIN`, `JOB_RETRY`, `EXPORT_AUDIT`, `VIEW_HEALTH`, etc.)
  - `resource_type TEXT NOT NULL` (e.g., `config.section`, `job.run`, `webhook.attempt`)
  - `resource_ref TEXT` (identifier of target resource)
  - `summary TEXT` (plain-language change summary)
  - `diff JSONB` (structured before/after deltas)
  - `occurred_at TIMESTAMPTZ DEFAULT now()`
  - `created_at TIMESTAMPTZ DEFAULT now()`
- **Indexes**:
  - `(tenant_id, occurred_at DESC)`
  - `(plugin_id, action, occurred_at DESC)`
  - GIN on `diff` for structured search
- **Relationships**:
  - Referenced by `AdminConsoleConfigChange.audit_event_id`
  - Referenced by `AdminConsoleJobRun.audit_event_id`

### AdminConsoleConfigChange
- **Purpose**: Persist a normalized history of configuration updates applied through the console.
- **Fields**:
  - `id UUID PK`
  - `tenant_id TEXT NOT NULL`
  - `plugin_id TEXT NOT NULL`
  - `section_key TEXT NOT NULL` (maps to form section identifier)
  - `change_type TEXT NOT NULL` (enum: `create`, `update`, `delete`)
  - `previous_snapshot JSONB` (state prior to change)
  - `next_snapshot JSONB` (state after change)
  - `validation_summary JSONB` (captured validation warnings/errors shown at submission time)
  - `audit_event_id UUID NOT NULL REFERENCES admin_console_audit_events(id) ON DELETE CASCADE`
  - `applied_at TIMESTAMPTZ DEFAULT now()`
- **Indexes**:
  - `(tenant_id, section_key, applied_at DESC)`
  - `(plugin_id, applied_at DESC)`

### AdminConsoleJobRun
- **Purpose**: Track safe-ops executions and background jobs initiated via the console to drive history tables and retry eligibility.
- **Fields**:
  - `id UUID PK`
  - `plugin_id TEXT NOT NULL`
  - `tenant_id TEXT` (nullable for global jobs)
  - `environment TEXT` (enum: `production`, `staging`, `sandbox`)
  - `job_type TEXT NOT NULL` (enum: `webhook_replay`, `task_retry`, `queue_drain`, `health_probe`, `custom`)
  - `trigger_source TEXT NOT NULL` (enum: `manual`, `schedule`, `alert`, `api`)
  - `status TEXT NOT NULL` (enum: `pending`, `running`, `succeeded`, `failed`, `cancelled`)
  - `started_at TIMESTAMPTZ`
  - `finished_at TIMESTAMPTZ`
  - `duration_ms BIGINT GENERATED ALWAYS AS (COALESCE(EXTRACT(EPOCH FROM (finished_at - started_at)) * 1000, 0)) STORED`
  - `message TEXT` (last status update)
  - `retry_of UUID REFERENCES admin_console_job_runs(id)`
  - `audit_event_id UUID REFERENCES admin_console_audit_events(id)`
  - `created_by TEXT NOT NULL` (operator id)
  - `created_at TIMESTAMPTZ DEFAULT now()`
  - `updated_at TIMESTAMPTZ DEFAULT now()`
- **Indexes**:
  - `(tenant_id, job_type, created_at DESC)`
  - `(status, created_at DESC)`
  - `(plugin_id, environment, created_at DESC)`

### AdminConsoleSafeOpLock (virtual)
- **Purpose**: Application-level concurrency guard (implemented via advisory lock or transactional flag) preventing overlapping safe operations on the same scope.
- **Fields**:
  - Not persisted as a table; realized via Postgres advisory locks keyed by `(plugin_id, scope_ref, action)` to avoid separate schema maintenance.
- **Reasoning**: We rely on advisory locks rather than a table to minimize write contention and because scope cardinality is low; lock acquisition failure returns a user-visible notice.

### Reused Entities
- `integration_webhook_attempts`: Source of delivery history; add composite index `(tenant_id, plugin_id, delivered_at DESC)` if missing.
- `runtime_ops.RuntimeAuditEvent`: Remains for runtime metrics but cross-linked via `audit_event_id` where safe-op actions already emit runtime audits.

## Data Retention & Governance
- Audit and config change tables retain 365 days of history, controlled by existing audit log retention config; archival job will purge older entries with export backup.
- Job run history keeps the most recent 45 days by default, configurable via admin console settings (backed by config key `admin_console.job_history_days`).

## Validation Rules
- `tenant_id` must align with JWT context; service layer enforces this before insert.
- `section_key` must match registered configuration schema descriptor; invalid keys rejected.
- `status` transitions validated: only `pending -> running -> {succeeded|failed}` or `pending -> cancelled`.
- `environment` restricted to enumerated values and must be consistent with tenant binding (production environment allowed only when manifest marks plugin as production-enabled).

## Migrations
- New tables created via `backend/migrations/2025Q4_admin_console.sql`.
- Migration adds indexes listed above and registers models in `backend/cmd/database/migrate/migrate.go`.
