-- Integration configuration approval workflow
BEGIN;

CREATE TABLE IF NOT EXISTS integration_change_approvals (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    target_type    TEXT NOT NULL,
    target_id      TEXT NOT NULL,
    payload        JSONB NOT NULL,
    status         TEXT NOT NULL DEFAULT 'PENDING',
    submitted_by   TEXT NOT NULL,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    reviewed_by    TEXT,
    reviewed_at    TIMESTAMPTZ,
    reason         TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integration_change_approvals_status
    ON integration_change_approvals (status, submitted_at DESC);

CREATE INDEX IF NOT EXISTS idx_integration_change_approvals_target
    ON integration_change_approvals (target_type, target_id);

COMMIT;
