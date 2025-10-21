-- Marketplace licensing & billing schema (US2)
BEGIN;

CREATE TABLE IF NOT EXISTS marketplace_pricing_plans (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id           TEXT NOT NULL,
    listing_id          UUID NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    plan_code           TEXT NOT NULL,
    plan_type           TEXT NOT NULL,
    currency            TEXT NOT NULL,
    amount              NUMERIC(18,4),
    billing_period      TEXT,
    trial_period_days   INTEGER,
    quota_limit         NUMERIC(18,4),
    overage_policy      TEXT,
    feature_matrix      JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_default          BOOLEAN NOT NULL DEFAULT FALSE,
    status              TEXT NOT NULL DEFAULT 'active',
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_pricing_plan_code
    ON marketplace_pricing_plans (listing_id, plan_code);

CREATE TABLE IF NOT EXISTS marketplace_plan_tiers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id        UUID NOT NULL REFERENCES marketplace_pricing_plans(id) ON DELETE CASCADE,
    tenant_id      TEXT NOT NULL,
    metric         TEXT NOT NULL,
    range_from     NUMERIC(18,4) NOT NULL,
    range_to       NUMERIC(18,4),
    unit_amount    NUMERIC(18,4) NOT NULL,
    unit_name      TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_plan_tiers_plan
    ON marketplace_plan_tiers (plan_id, range_from);

CREATE TABLE IF NOT EXISTS marketplace_licenses (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id          TEXT NOT NULL,
    listing_id         UUID NOT NULL REFERENCES marketplace_listings(id) ON DELETE CASCADE,
    plan_id            UUID NOT NULL REFERENCES marketplace_pricing_plans(id) ON DELETE RESTRICT,
    license_token      TEXT NOT NULL,
    status             TEXT NOT NULL DEFAULT 'active',
    issued_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at         TIMESTAMPTZ NOT NULL,
    renewal_token      TEXT,
    offline_until      TIMESTAMPTZ,
    last_validated_at  TIMESTAMPTZ,
    issued_by          TEXT,
    metadata           JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_offline_within_72h CHECK (
        offline_until IS NULL OR issued_at IS NULL OR offline_until <= issued_at + INTERVAL '72 hours'
    )
);

CREATE INDEX IF NOT EXISTS idx_marketplace_license_tenant_listing
    ON marketplace_licenses (tenant_id, listing_id);

CREATE INDEX IF NOT EXISTS idx_marketplace_license_plan_status
    ON marketplace_licenses (plan_id, status);

CREATE TABLE IF NOT EXISTS marketplace_license_events (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     TEXT NOT NULL,
    license_id    UUID NOT NULL REFERENCES marketplace_licenses(id) ON DELETE CASCADE,
    event_type    TEXT NOT NULL,
    event_payload JSONB NOT NULL DEFAULT '{}'::jsonb,
    emitted_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actor_id      TEXT,
    trace_id      TEXT,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_license_events
    ON marketplace_license_events (license_id, event_type, emitted_at);

CREATE TABLE IF NOT EXISTS marketplace_tax_transactions (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              TEXT NOT NULL,
    billing_id             TEXT NOT NULL,
    external_provider      TEXT NOT NULL,
    external_transaction_id TEXT,
    jurisdiction           TEXT,
    tax_amount             NUMERIC(18,4) NOT NULL,
    currency               TEXT NOT NULL,
    settlement_currency    TEXT,
    exchange_rate          NUMERIC(18,6),
    tax_amount_settlement  NUMERIC(18,4),
    raw_payload            JSONB NOT NULL DEFAULT '{}'::jsonb,
    status                 TEXT NOT NULL DEFAULT 'pending',
    synced_at              TIMESTAMPTZ,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_tax_transactions_status
    ON marketplace_tax_transactions (status, external_provider);

COMMIT;

BEGIN;

CREATE OR REPLACE FUNCTION current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE marketplace_pricing_plans FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_pricing_plans_tenant ON marketplace_pricing_plans
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_plan_tiers FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_plan_tiers_tenant ON marketplace_plan_tiers
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_licenses FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_licenses_tenant ON marketplace_licenses
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_license_events FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_license_events_tenant ON marketplace_license_events
    USING (tenant_id = current_tenant());

ALTER TABLE marketplace_tax_transactions FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_tax_transactions_tenant ON marketplace_tax_transactions
    USING (tenant_id = current_tenant());

COMMIT;
