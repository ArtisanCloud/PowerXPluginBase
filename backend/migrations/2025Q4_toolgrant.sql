-- ToolGrant lifecycle tables
BEGIN;

CREATE TABLE IF NOT EXISTS public.tool_grant_revocations (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    TEXT NOT NULL,
    toolgrant_id TEXT NOT NULL,
    revoked_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_by   TEXT NOT NULL,
    reason       TEXT,
    ttl_expiry   TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_tool_grant_revocation ON public.tool_grant_revocations (tenant_id, toolgrant_id);
CREATE INDEX IF NOT EXISTS idx_tool_grant_revocation_expiry ON public.tool_grant_revocations (ttl_expiry);

CREATE TABLE IF NOT EXISTS public.tool_grant_usage_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    TEXT NOT NULL,
    toolgrant_id TEXT NOT NULL,
    event_type   TEXT NOT NULL,
    capability   TEXT NOT NULL,
    agent_id     TEXT NOT NULL,
    occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    metadata     JSONB
);

CREATE INDEX IF NOT EXISTS idx_tool_grant_usage_tenant ON public.tool_grant_usage_events (tenant_id, toolgrant_id);
CREATE INDEX IF NOT EXISTS idx_tool_grant_usage_event ON public.tool_grant_usage_events (event_type, occurred_at DESC);

COMMIT;
