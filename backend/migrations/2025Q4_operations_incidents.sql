-- Operations Incident schema objects
BEGIN;

CREATE TABLE IF NOT EXISTS operations_incidents (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id         TEXT NOT NULL,
    tenant_id         TEXT,
    severity          TEXT NOT NULL,
    status            TEXT NOT NULL,
    detection_source  TEXT NOT NULL,
    summary           TEXT NOT NULL,
    impact            JSONB,
    mitigation        TEXT,
    root_cause        TEXT,
    next_update_at    TIMESTAMPTZ,
    labels            JSONB,
    confidentiality   TEXT,
    detected_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    acknowledged_at   TIMESTAMPTZ,
    mitigated_at      TIMESTAMPTZ,
    resolved_at       TIMESTAMPTZ,
    closed_at         TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_incidents_scope
    ON operations_incidents (plugin_id, tenant_id);
CREATE INDEX IF NOT EXISTS idx_operations_incidents_status
    ON operations_incidents (status, severity);

CREATE TABLE IF NOT EXISTS operations_incident_updates (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id        UUID NOT NULL REFERENCES operations_incidents(id) ON DELETE CASCADE,
    entry_type         TEXT NOT NULL,
    message            TEXT NOT NULL,
    stakeholder_channel TEXT,
    author_role        TEXT,
    metadata           JSONB,
    posted_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_incident_updates_incident
    ON operations_incident_updates (incident_id, posted_at);

CREATE TABLE IF NOT EXISTS operations_incident_checklist (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    incident_id UUID NOT NULL REFERENCES operations_incidents(id) ON DELETE CASCADE,
    item_key    TEXT NOT NULL,
    description TEXT,
    status      TEXT NOT NULL,
    completed_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (incident_id, item_key)
);

CREATE INDEX IF NOT EXISTS idx_operations_incident_checklist_incident
    ON operations_incident_checklist (incident_id);

COMMIT;

ALTER TABLE operations_incidents ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_incident_updates ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_incident_checklist ENABLE ROW LEVEL SECURITY;
