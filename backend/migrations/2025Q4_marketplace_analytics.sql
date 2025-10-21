-- Marketplace analytics schema (US3)
BEGIN;

CREATE TABLE IF NOT EXISTS marketplace_usage_envelopes (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id        TEXT NOT NULL,
    license_id       UUID NOT NULL,
    plugin_id        TEXT NOT NULL,
    metrics          JSONB NOT NULL DEFAULT '[]'::jsonb,
    timestamp_start  TIMESTAMPTZ NOT NULL,
    timestamp_end    TIMESTAMPTZ NOT NULL,
    signature        TEXT NOT NULL,
    checksum         TEXT NOT NULL,
    ingest_status    TEXT NOT NULL DEFAULT 'processed',
    ingested_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_usage_checksum
    ON marketplace_usage_envelopes (checksum);

CREATE INDEX IF NOT EXISTS idx_marketplace_usage_license
    ON marketplace_usage_envelopes (tenant_id, license_id);

CREATE TABLE IF NOT EXISTS marketplace_usage_aggregates (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id   TEXT NOT NULL,
    license_id  UUID NOT NULL,
    metric      TEXT NOT NULL,
    window      TEXT NOT NULL,
    time_bucket TIMESTAMPTZ NOT NULL,
    total       NUMERIC(20,4) NOT NULL DEFAULT 0,
    delta       NUMERIC(20,4) NOT NULL DEFAULT 0,
    currency    TEXT,
    revenue     NUMERIC(18,4) NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_usage_aggregate_bucket
    ON marketplace_usage_aggregates (tenant_id, license_id, metric, window, time_bucket);

CREATE INDEX IF NOT EXISTS idx_marketplace_usage_aggregate_metric
    ON marketplace_usage_aggregates (tenant_id, metric, window, time_bucket);

CREATE TABLE IF NOT EXISTS marketplace_revenue_share_reports (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id       TEXT NOT NULL,
    vendor_id       TEXT NOT NULL,
    period_start    TIMESTAMPTZ NOT NULL,
    period_end      TIMESTAMPTZ NOT NULL,
    gross_amount    NUMERIC(18,4) NOT NULL DEFAULT 0,
    vendor_share    NUMERIC(18,4) NOT NULL DEFAULT 0,
    platform_share  NUMERIC(18,4) NOT NULL DEFAULT 0,
    fees            NUMERIC(18,4) NOT NULL DEFAULT 0,
    currency        TEXT NOT NULL,
    status          TEXT NOT NULL DEFAULT 'draft',
    generated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    export_uri      TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_revenue_report_period
    ON marketplace_revenue_share_reports (tenant_id, vendor_id, period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_marketplace_revenue_vendor
    ON marketplace_revenue_share_reports (tenant_id, vendor_id, status);

CREATE TABLE IF NOT EXISTS marketplace_notifications (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id      TEXT NOT NULL,
    recipient_type TEXT NOT NULL,
    recipient_id   TEXT NOT NULL,
    channel        TEXT NOT NULL,
    template_code  TEXT NOT NULL,
    payload        JSONB NOT NULL DEFAULT '{}'::jsonb,
    scheduled_at   TIMESTAMPTZ,
    sent_at        TIMESTAMPTZ,
    status         TEXT NOT NULL DEFAULT 'pending',
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_notifications_recipient
    ON marketplace_notifications (tenant_id, recipient_type, recipient_id, status);

COMMIT;

BEGIN;

ALTER TABLE marketplace_usage_envelopes FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_usage_envelopes_tenant_isolation ON marketplace_usage_envelopes
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_usage_aggregates FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_usage_aggregates_tenant_isolation ON marketplace_usage_aggregates
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_revenue_share_reports FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_revenue_reports_tenant_isolation ON marketplace_revenue_share_reports
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_notifications FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_notifications_tenant_isolation ON marketplace_notifications
    USING (tenant_id = current_tenant());

COMMIT;
