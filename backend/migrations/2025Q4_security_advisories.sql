-- Vulnerability advisory and distribution tables
BEGIN;

CREATE TABLE IF NOT EXISTS public.security_vulnerability_advisories (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reference          TEXT NOT NULL,
    severity           TEXT NOT NULL,
    status             TEXT NOT NULL,
    affected_versions  JSONB NOT NULL DEFAULT '[]'::jsonb,
    patched_in_version TEXT,
    summary            TEXT NOT NULL,
    details_markdown   TEXT,
    published_at       TIMESTAMPTZ,
    patched_at         TIMESTAMPTZ,
    closed_at          TIMESTAMPTZ,
    sla_deadline       TIMESTAMPTZ,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_security_advisory_reference
    ON public.security_vulnerability_advisories (reference);

CREATE INDEX IF NOT EXISTS idx_security_advisory_severity_status
    ON public.security_vulnerability_advisories (severity, status);

CREATE INDEX IF NOT EXISTS gin_security_advisory_versions
    ON public.security_vulnerability_advisories USING gin (affected_versions jsonb_path_ops);

CREATE TABLE IF NOT EXISTS public.security_advisory_distributions (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    advisory_id  UUID NOT NULL REFERENCES public.security_vulnerability_advisories(id) ON DELETE CASCADE,
    tenant_id    TEXT NOT NULL,
    channel      TEXT NOT NULL,
    delivered_at TIMESTAMPTZ,
    status       TEXT NOT NULL,
    metadata     JSONB,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_security_advisory_delivery
    ON public.security_advisory_distributions (advisory_id, tenant_id, channel);

CREATE INDEX IF NOT EXISTS idx_security_advisory_delivery_status
    ON public.security_advisory_distributions (status, delivered_at);

COMMIT;
