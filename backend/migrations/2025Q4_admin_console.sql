-- Admin Console schema objects for powerx_plugin_base
BEGIN;

CREATE TABLE IF NOT EXISTS admin_console_audit_events (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id        TEXT NOT NULL,
    tenant_id        TEXT,
    actor_id         TEXT NOT NULL,
    actor_name       TEXT,
    actor_email      TEXT,
    permission_code  TEXT NOT NULL,
    action           TEXT NOT NULL,
    resource_type    TEXT NOT NULL,
    resource_ref     TEXT,
    summary          TEXT,
    diff             JSONB,
    occurred_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_console_audit_tenant_time
    ON admin_console_audit_events (tenant_id, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_console_audit_plugin_action_time
    ON admin_console_audit_events (plugin_id, action, occurred_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_console_audit_diff
    ON admin_console_audit_events USING gin (diff);

CREATE TABLE IF NOT EXISTS admin_console_config_changes (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id          TEXT NOT NULL,
    tenant_id          TEXT,
    section_key        TEXT NOT NULL,
    change_type        TEXT NOT NULL,
    previous_snapshot  JSONB,
    next_snapshot      JSONB,
    validation_summary JSONB,
    audit_event_id     UUID NOT NULL REFERENCES admin_console_audit_events(id) ON DELETE CASCADE,
    applied_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_console_config_section_time
    ON admin_console_config_changes (tenant_id, section_key, applied_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_console_config_plugin_time
    ON admin_console_config_changes (plugin_id, applied_at DESC);

CREATE TABLE IF NOT EXISTS admin_console_job_runs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id      TEXT NOT NULL,
    tenant_id      TEXT,
    environment    TEXT,
    job_type       TEXT NOT NULL,
    trigger_source TEXT NOT NULL,
    status         TEXT NOT NULL,
    started_at     TIMESTAMPTZ,
    finished_at    TIMESTAMPTZ,
    duration_ms    BIGINT GENERATED ALWAYS AS (
        COALESCE(
            (EXTRACT(EPOCH FROM (finished_at - started_at)) * 1000)::BIGINT,
            0
        )
    ) STORED,
    message        TEXT,
    retry_of       UUID REFERENCES admin_console_job_runs(id),
    audit_event_id UUID REFERENCES admin_console_audit_events(id),
    created_by     TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_admin_console_job_type_time
    ON admin_console_job_runs (tenant_id, job_type, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_console_job_status_time
    ON admin_console_job_runs (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_console_job_environment_time
    ON admin_console_job_runs (plugin_id, environment, created_at DESC);

COMMIT;

CREATE OR REPLACE FUNCTION current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE admin_console_audit_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE admin_console_config_changes ENABLE ROW LEVEL SECURITY;
ALTER TABLE admin_console_job_runs ENABLE ROW LEVEL SECURITY;

ALTER TABLE admin_console_audit_events FORCE ROW LEVEL SECURITY;
CREATE POLICY admin_console_audit_events_tenant_isolation
    ON admin_console_audit_events
    USING (tenant_id IS NULL OR tenant_id = current_tenant());

ALTER TABLE admin_console_config_changes FORCE ROW LEVEL SECURITY;
CREATE POLICY admin_console_config_changes_tenant_isolation
    ON admin_console_config_changes
    USING (tenant_id IS NULL OR tenant_id = current_tenant());

ALTER TABLE admin_console_job_runs FORCE ROW LEVEL SECURITY;
CREATE POLICY admin_console_job_runs_tenant_isolation
    ON admin_console_job_runs
    USING (tenant_id IS NULL OR tenant_id = current_tenant());
