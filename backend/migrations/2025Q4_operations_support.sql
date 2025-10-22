-- Operations Support schema objects
BEGIN;

CREATE TABLE IF NOT EXISTS operations_support_channels (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id       TEXT NOT NULL,
    tenant_id       TEXT,
    channel         TEXT NOT NULL,
    is_enabled      BOOLEAN NOT NULL DEFAULT TRUE,
    service_window  JSONB,
    escalation_path JSONB,
    metadata        JSONB,
    version         INTEGER NOT NULL DEFAULT 1,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_support_channels_scope
    ON operations_support_channels (plugin_id, tenant_id);

CREATE TABLE IF NOT EXISTS operations_support_tickets (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id          TEXT NOT NULL,
    tenant_id          TEXT NOT NULL,
    channel_id         UUID REFERENCES operations_support_channels(id) ON DELETE SET NULL,
    external_ref       TEXT,
    subject            TEXT NOT NULL,
    description        TEXT,
    priority           TEXT NOT NULL,
    status             TEXT NOT NULL,
    requested_by       JSONB,
    assigned_team      TEXT,
    assigned_to        UUID,
    knowledge_base_refs TEXT[],
    first_response_at  TIMESTAMPTZ,
    resolved_at        TIMESTAMPTZ,
    closed_at          TIMESTAMPTZ,
    csat_score         NUMERIC(2,1),
    resolution_code    TEXT,
    reopen_count       INTEGER NOT NULL DEFAULT 0,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_support_tickets_scope
    ON operations_support_tickets (plugin_id, tenant_id);

CREATE INDEX IF NOT EXISTS idx_operations_support_tickets_status
    ON operations_support_tickets (status, priority);

CREATE TABLE IF NOT EXISTS operations_support_ticket_events (
    id             BIGSERIAL PRIMARY KEY,
    ticket_id      UUID NOT NULL REFERENCES operations_support_tickets(id) ON DELETE CASCADE,
    event_type     TEXT NOT NULL,
    payload        JSONB,
    emitted_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    webhook_status TEXT DEFAULT 'pending',
    retry_count    INTEGER NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_operations_support_ticket_events_ticket
    ON operations_support_ticket_events (ticket_id, emitted_at DESC);

CREATE TABLE IF NOT EXISTS operations_readiness_checklist_items (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plugin_id   TEXT NOT NULL,
    type        TEXT NOT NULL,
    item_key    TEXT NOT NULL,
    description TEXT,
    status      TEXT NOT NULL,
    owner_role  TEXT,
    due_date    DATE,
    completed_at TIMESTAMPTZ,
    notes       TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (plugin_id, type, item_key)
);

CREATE INDEX IF NOT EXISTS idx_operations_readiness_type
    ON operations_readiness_checklist_items (plugin_id, type);

COMMIT;

ALTER TABLE operations_support_channels ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_support_tickets ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_support_ticket_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE operations_readiness_checklist_items ENABLE ROW LEVEL SECURITY;
