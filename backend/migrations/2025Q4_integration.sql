-- Integration foundational schema objects
BEGIN;

CREATE TABLE IF NOT EXISTS public.integration_grant_matrix_overrides (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scope         TEXT NOT NULL,
    channel       TEXT NOT NULL,
    resource      TEXT NOT NULL,
    action        TEXT NOT NULL,
    constraints   JSONB NOT NULL DEFAULT '{}'::jsonb,
    status        TEXT NOT NULL DEFAULT 'PENDING',
    version       INTEGER NOT NULL DEFAULT 1,
    approved_by   TEXT,
    approved_at   TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_integration_grant_matrix_override
    ON public.integration_grant_matrix_overrides (scope, channel, resource, action);

CREATE TABLE IF NOT EXISTS public.integration_webhook_subscriptions (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id     TEXT NOT NULL,
    event_type    TEXT NOT NULL,
    target_url    TEXT NOT NULL,
    secret        TEXT,
    retry_policy  JSONB NULL,
    status        TEXT NOT NULL DEFAULT 'ACTIVE',
    metadata      JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_integration_webhook_subscription
    ON public.integration_webhook_subscriptions (tenant_id, event_type, target_url);

CREATE TABLE IF NOT EXISTS public.integration_webhook_attempts (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id   UUID NOT NULL REFERENCES public.integration_webhook_subscriptions(id) ON DELETE CASCADE,
    envelope_id       UUID,
    status            TEXT NOT NULL,
    retry_count       INTEGER NOT NULL DEFAULT 0,
    last_error        TEXT,
    next_delivery_at  TIMESTAMPTZ,
    payload_snapshot  JSONB,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_attempt_subscription
    ON public.integration_webhook_attempts (subscription_id);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_attempt_status
    ON public.integration_webhook_attempts (status, next_delivery_at);

CREATE TABLE IF NOT EXISTS public.integration_webhook_dlq (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id     UUID NOT NULL REFERENCES public.integration_webhook_attempts(id) ON DELETE CASCADE,
    failure_reason TEXT,
    payload        JSONB,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_dlq_attempt
    ON public.integration_webhook_dlq (attempt_id);

CREATE TABLE IF NOT EXISTS public.integration_secrets (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id              TEXT NOT NULL,
    integration_type       TEXT NOT NULL,
    current_secret_ref     TEXT,
    pending_secret_ref     TEXT,
    rotation_interval_days INTEGER NOT NULL DEFAULT 30,
    last_rotated_at        TIMESTAMPTZ,
    next_rotation_due_at   TIMESTAMPTZ,
    status                 TEXT NOT NULL DEFAULT 'ACTIVE',
    audit_log              JSONB NOT NULL DEFAULT '[]'::jsonb,
    metadata               JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS uniq_integration_secrets_tenant_type
    ON public.integration_secrets (tenant_id, integration_type);

CREATE TABLE IF NOT EXISTS public.integration_idempotency_records (
    key            TEXT PRIMARY KEY,
    tenant_id      TEXT NOT NULL,
    scope          TEXT,
    operation      TEXT,
    payload_hash   TEXT,
    response_data  JSONB,
    metadata       JSONB NOT NULL DEFAULT '{}'::jsonb,
    expires_at     TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_integration_idempotency_expires
    ON public.integration_idempotency_records (expires_at);

COMMIT;
