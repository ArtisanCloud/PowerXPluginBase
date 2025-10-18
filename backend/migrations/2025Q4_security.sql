-- Security baseline schema objects
BEGIN;

CREATE TABLE IF NOT EXISTS public.security_baseline_checklists (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    version      TEXT NOT NULL,
    controls     JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    retired_at   TIMESTAMPTZ
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_security_baseline_version
    ON public.security_baseline_checklists (version);

CREATE TABLE IF NOT EXISTS public.security_audit_reports (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    baseline_id        UUID NOT NULL REFERENCES public.security_baseline_checklists(id) ON DELETE RESTRICT,
    initiated_by       TEXT NOT NULL,
    status             TEXT NOT NULL,
    findings           JSONB,
    artifact_path      TEXT,
    sarif_path         TEXT,
    report_hash        TEXT,
    checklist_version  TEXT NOT NULL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_security_audit_status ON public.security_audit_reports (status, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_security_audit_baseline ON public.security_audit_reports (baseline_id);

COMMIT;
