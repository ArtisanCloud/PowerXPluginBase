-- Operations SLA schema objects
BEGIN;

CREATE TABLE IF NOT EXISTS operations_sla_profiles (
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id                  TEXT NOT NULL,
    plan_type                  TEXT NOT NULL,
    uptime_target              NUMERIC(5,2) NOT NULL DEFAULT 0,
    uptime_actual              NUMERIC(5,2) NOT NULL DEFAULT 0,
    response_target_ms         INTEGER NOT NULL DEFAULT 0,
    response_actual_ms         INTEGER NOT NULL DEFAULT 0,
    success_target_pct         NUMERIC(5,2) NOT NULL DEFAULT 0,
    success_actual_pct         NUMERIC(5,2) NOT NULL DEFAULT 0,
    support_frt_target_hours   NUMERIC(4,2) NOT NULL DEFAULT 0,
    support_frt_actual_hours   NUMERIC(4,2) NOT NULL DEFAULT 0,
    sla_score                  NUMERIC(5,2) NOT NULL DEFAULT 0,
    incentive_applied_at       TIMESTAMPTZ,
    penalty_applied_at         TIMESTAMPTZ,
    notes                      TEXT,
    computed_at                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (plugin_id, plan_type)
);

CREATE INDEX IF NOT EXISTS idx_operations_sla_profiles_plugin
    ON operations_sla_profiles (plugin_id, plan_type);

CREATE TABLE IF NOT EXISTS operations_sla_adjustments (
    id            BIGSERIAL PRIMARY KEY,
    plugin_id     TEXT NOT NULL,
    plan_type     TEXT NOT NULL,
    period_start  DATE NOT NULL,
    period_end    DATE NOT NULL,
    score_before  NUMERIC(5,2) NOT NULL,
    score_after   NUMERIC(5,2) NOT NULL,
    action        TEXT NOT NULL,
    details       TEXT,
    applied_by    UUID,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_sla_adjustments_plugin
    ON operations_sla_adjustments (plugin_id, plan_type, period_start DESC);

COMMIT;

ALTER TABLE operations_sla_profiles ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_sla_adjustments ENABLE ROW LEVEL SECURITY;
