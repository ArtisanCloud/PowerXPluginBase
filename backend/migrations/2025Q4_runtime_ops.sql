-- Runtime Ops schema objects for powerx_plugin_base
BEGIN;

-- Runtime assignments track each plugin instance and resource envelope
CREATE TABLE IF NOT EXISTS powerx_plugin_base.runtime_assignments (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id         TEXT NOT NULL,
    tenant_scope      TEXT,
    runtime_mode      TEXT NOT NULL,
    host_id           TEXT NOT NULL,
    port              INTEGER,
    status            TEXT NOT NULL,
    cpu_limit         TEXT,
    memory_limit      TEXT,
    network_profile   TEXT,
    restart_count     INTEGER DEFAULT 0,
    ready_at          TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE powerx_plugin_base.runtime_assignments
    ADD COLUMN IF NOT EXISTS tenant_id TEXT;

CREATE INDEX IF NOT EXISTS idx_runtime_assignments_plugin ON powerx_plugin_base.runtime_assignments (plugin_id);
CREATE INDEX IF NOT EXISTS idx_runtime_assignments_tenant ON powerx_plugin_base.runtime_assignments (tenant_id);
CREATE INDEX IF NOT EXISTS idx_runtime_assignments_status_host ON powerx_plugin_base.runtime_assignments (status, host_id);
CREATE INDEX IF NOT EXISTS idx_runtime_assignments_ready_at ON powerx_plugin_base.runtime_assignments (ready_at DESC);

-- Port reservations enforce unique port usage per instance
CREATE TABLE IF NOT EXISTS powerx_plugin_base.port_reservations (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    runtime_assignment_id UUID NOT NULL REFERENCES powerx_plugin_base.runtime_assignments(id) ON DELETE CASCADE,
    port              INTEGER NOT NULL,
    host_id           TEXT NOT NULL,
    state             TEXT NOT NULL DEFAULT 'active',
    reserved_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    released_at       TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_port_reservation_active ON powerx_plugin_base.port_reservations (host_id, port) WHERE state = 'active';

-- MCP sessions table captures lifecycle state and heartbeat metrics
CREATE TABLE IF NOT EXISTS powerx_plugin_base.mcp_sessions (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    runtime_assignment_id UUID NOT NULL REFERENCES powerx_plugin_base.runtime_assignments(id) ON DELETE CASCADE,
    tenant_id         TEXT NOT NULL,
    state             TEXT NOT NULL,
    jwt_id            TEXT,
    capabilities_hash TEXT,
    missed_heartbeats INTEGER DEFAULT 0,
    last_ping_at      TIMESTAMPTZ,
    closed_at         TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mcp_sessions_state ON powerx_plugin_base.mcp_sessions (state);
CREATE INDEX IF NOT EXISTS idx_mcp_sessions_tenant ON powerx_plugin_base.mcp_sessions (tenant_id);
CREATE INDEX IF NOT EXISTS idx_mcp_sessions_assignment ON powerx_plugin_base.mcp_sessions (runtime_assignment_id);

-- Runtime audit log for immutable events
CREATE TABLE IF NOT EXISTS powerx_plugin_base.runtime_audit_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id         TEXT NOT NULL,
    tenant_id         TEXT,
    event_type        TEXT NOT NULL,
    payload           JSONB,
    occurred_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_runtime_audit_plugin ON powerx_plugin_base.runtime_audit_events (plugin_id);
CREATE INDEX IF NOT EXISTS idx_runtime_audit_tenant ON powerx_plugin_base.runtime_audit_events (tenant_id);
CREATE INDEX IF NOT EXISTS idx_runtime_audit_type_time ON powerx_plugin_base.runtime_audit_events (event_type, occurred_at DESC);

-- Quota ledger aggregates consumption by scope
CREATE TABLE IF NOT EXISTS powerx_plugin_base.quota_ledger (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope_type        TEXT NOT NULL,
    scope_ref         TEXT NOT NULL,
    window_start      TIMESTAMPTZ NOT NULL,
    window_end        TIMESTAMPTZ NOT NULL,
    tokens_consumed   NUMERIC DEFAULT 0,
    cpu_seconds       NUMERIC DEFAULT 0,
    bandwidth_mb      NUMERIC DEFAULT 0,
    invocations       NUMERIC DEFAULT 0,
    over_limit_action TEXT,
    reported_at       TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_quota_ledger_scope ON powerx_plugin_base.quota_ledger (scope_type, scope_ref);
CREATE INDEX IF NOT EXISTS idx_quota_ledger_window ON powerx_plugin_base.quota_ledger (window_start, window_end);
CREATE INDEX IF NOT EXISTS idx_quota_ledger_reported_at ON powerx_plugin_base.quota_ledger (reported_at);

-- Marketplace hourly overage summaries
CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_overages (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id         TEXT NOT NULL,
    tenant_id         TEXT,
    hour_window       TIMESTAMPTZ NOT NULL,
    quota_metric      TEXT NOT NULL,
    breach_count      INTEGER NOT NULL DEFAULT 0,
    last_breach_at    TIMESTAMPTZ,
    reported          BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_overages_window ON powerx_plugin_base.marketplace_overages (hour_window);
CREATE INDEX IF NOT EXISTS idx_marketplace_overages_plugin_tenant ON powerx_plugin_base.marketplace_overages (plugin_id, tenant_id);
CREATE INDEX IF NOT EXISTS idx_marketplace_overages_unreported ON powerx_plugin_base.marketplace_overages (plugin_id, tenant_id) WHERE reported = false;

COMMIT;

-- Apply row-level security and tenant guardrails
ALTER TABLE powerx_plugin_base.runtime_assignments ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.port_reservations ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.mcp_sessions ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.runtime_audit_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.quota_ledger ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.marketplace_overages ENABLE ROW LEVEL SECURITY;
