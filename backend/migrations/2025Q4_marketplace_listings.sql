-- Marketplace listings & checklist schema (US1)
BEGIN;

CREATE SCHEMA IF NOT EXISTS powerx_plugin_base;

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_listings (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              TEXT NOT NULL,
    plugin_id              TEXT NOT NULL,
    vendor_id              TEXT NOT NULL,
    status                 TEXT NOT NULL DEFAULT 'draft',
    title                  TEXT NOT NULL,
    slug                   TEXT NOT NULL,
    summary                TEXT,
    description            TEXT,
    cover_asset_id         UUID,
    hero_video_asset_id    UUID,
    categories             JSONB NOT NULL DEFAULT '[]'::jsonb,
    tags                   JSONB NOT NULL DEFAULT '[]'::jsonb,
    locale                 TEXT NOT NULL DEFAULT 'en',
    version                TEXT,
    ready_checklist_score  INTEGER NOT NULL DEFAULT 0,
    recommended_weight     NUMERIC(10,4) NOT NULL DEFAULT 0,
    published_at           TIMESTAMPTZ,
    reviewed_at            TIMESTAMPTZ,
    reviewer_id            TEXT,
    audit_notes            TEXT,
    branding_theme         JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at             TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_listing_slug
    ON powerx_plugin_base.marketplace_listings (tenant_id, slug, locale)
    WHERE deleted_at IS NULL;

CREATE INDEX IF NOT EXISTS idx_marketplace_listing_status
    ON powerx_plugin_base.marketplace_listings (tenant_id, plugin_id, status);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_listing_assets (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id     UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_listings(id) ON DELETE CASCADE,
    tenant_id      TEXT NOT NULL,
    asset_type     TEXT NOT NULL,
    storage_uri    TEXT NOT NULL,
    checksum       TEXT,
    is_primary     BOOLEAN NOT NULL DEFAULT FALSE,
    locale         TEXT NOT NULL DEFAULT 'en',
    weight         INTEGER NOT NULL DEFAULT 0,
    metadata       JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_marketplace_listing_asset_type
    ON powerx_plugin_base.marketplace_listing_assets (listing_id, asset_type);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_listing_versions (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id     UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_listings(id) ON DELETE CASCADE,
    tenant_id      TEXT NOT NULL,
    version        TEXT NOT NULL,
    changelog      TEXT,
    metadata       JSONB NOT NULL DEFAULT '{}'::jsonb,
    submitted_by   TEXT NOT NULL,
    review_state   TEXT NOT NULL DEFAULT 'draft',
    reviewer_id    TEXT,
    reviewed_at    TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_listing_version
    ON powerx_plugin_base.marketplace_listing_versions (listing_id, version);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_pricing_plans (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id        UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_listings(id) ON DELETE CASCADE,
    tenant_id         TEXT NOT NULL,
    plan_code         TEXT NOT NULL,
    plan_type         TEXT NOT NULL,
    currency          TEXT NOT NULL,
    amount            NUMERIC(18,4),
    billing_period    TEXT,
    trial_period_days INTEGER,
    quota_limit       NUMERIC(18,4),
    overage_policy    TEXT,
    feature_matrix    JSONB NOT NULL DEFAULT '{}'::jsonb,
    is_default        BOOLEAN NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_pricing_plan_code
    ON powerx_plugin_base.marketplace_pricing_plans (listing_id, plan_code);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_plan_tiers (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id        UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_pricing_plans(id) ON DELETE CASCADE,
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
    ON powerx_plugin_base.marketplace_plan_tiers (plan_id, range_from);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_checklist_runs (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    listing_id     UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_listings(id) ON DELETE CASCADE,
    tenant_id      TEXT NOT NULL,
    trigger_source TEXT NOT NULL,
    run_number     INTEGER NOT NULL DEFAULT 1,
    status         TEXT NOT NULL DEFAULT 'pending',
    started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at   TIMESTAMPTZ,
    summary        TEXT,
    ci_pipeline_id TEXT,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_checklist_run_num
    ON powerx_plugin_base.marketplace_checklist_runs (listing_id, run_number);

CREATE TABLE IF NOT EXISTS powerx_plugin_base.marketplace_checklist_items (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    checklist_run_id  UUID NOT NULL REFERENCES powerx_plugin_base.marketplace_checklist_runs(id) ON DELETE CASCADE,
    tenant_id         TEXT NOT NULL,
    code              TEXT NOT NULL,
    description       TEXT NOT NULL,
    result            TEXT NOT NULL,
    evidence_uri      TEXT,
    notes             TEXT,
    auto_fix_link     TEXT,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_marketplace_checklist_item_code
    ON powerx_plugin_base.marketplace_checklist_items (checklist_run_id, code);

COMMIT;

BEGIN;

CREATE OR REPLACE FUNCTION powerx_plugin_base.current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE powerx_plugin_base.marketplace_listings FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_listings_tenant_isolation ON powerx_plugin_base.marketplace_listings
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_listing_assets FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_listing_assets_tenant_isolation ON powerx_plugin_base.marketplace_listing_assets
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_listing_versions FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_listing_versions_tenant_isolation ON powerx_plugin_base.marketplace_listing_versions
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_pricing_plans FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_pricing_plans_tenant_isolation ON powerx_plugin_base.marketplace_pricing_plans
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_plan_tiers FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_plan_tiers_tenant_isolation ON powerx_plugin_base.marketplace_plan_tiers
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_checklist_runs FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_checklist_runs_tenant_isolation ON powerx_plugin_base.marketplace_checklist_runs
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_checklist_items FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_checklist_items_tenant_isolation ON powerx_plugin_base.marketplace_checklist_items
    USING (tenant_id = powerx_plugin_base.current_tenant());

COMMIT;
