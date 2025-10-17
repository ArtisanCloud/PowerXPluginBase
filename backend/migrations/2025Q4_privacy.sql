-- Privacy domain schema objects for powerx_plugin_base
BEGIN;

CREATE TABLE IF NOT EXISTS powerx_plugin_base.privacy_data_classifications (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         TEXT NOT NULL,
    asset_key         TEXT NOT NULL,
    category          TEXT NOT NULL,
    lawful_basis      TEXT NOT NULL,
    retention_policy  JSONB,
    purpose           TEXT NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_privacy_classification_tenant_asset
    ON powerx_plugin_base.privacy_data_classifications (tenant_id, asset_key);
CREATE INDEX IF NOT EXISTS idx_privacy_classification_category
    ON powerx_plugin_base.privacy_data_classifications (tenant_id, category);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.privacy_consent_tokens (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         TEXT NOT NULL,
    consent_token     TEXT NOT NULL,
    scope             JSONB NOT NULL,
    expires_at        TIMESTAMPTZ,
    issued_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    issued_by         TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'ACTIVE',
    revoked_at        TIMESTAMPTZ,
    revoked_reason    TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_privacy_consent_token
    ON powerx_plugin_base.privacy_consent_tokens (tenant_id, consent_token);
CREATE INDEX IF NOT EXISTS idx_privacy_consent_status
    ON powerx_plugin_base.privacy_consent_tokens (tenant_id, status);
CREATE INDEX IF NOT EXISTS idx_privacy_consent_expires
    ON powerx_plugin_base.privacy_consent_tokens (tenant_id, expires_at DESC);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.privacy_lifecycle_events (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id         TEXT NOT NULL,
    event_type        TEXT NOT NULL,
    asset_key         TEXT NOT NULL,
    payload           JSONB,
    occurred_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    recorded_by       TEXT NOT NULL,
    status            TEXT NOT NULL DEFAULT 'PENDING',
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_privacy_lifecycle_tenant_event
    ON powerx_plugin_base.privacy_lifecycle_events (tenant_id, event_type, occurred_at DESC);

COMMIT;

ALTER TABLE powerx_plugin_base.privacy_data_classifications ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.privacy_consent_tokens ENABLE ROW LEVEL SECURITY;
ALTER TABLE powerx_plugin_base.privacy_lifecycle_events ENABLE ROW LEVEL SECURITY;
