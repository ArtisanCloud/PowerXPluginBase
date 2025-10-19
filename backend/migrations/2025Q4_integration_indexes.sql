BEGIN;

CREATE INDEX IF NOT EXISTS idx_integration_webhook_subscriptions_tenant_event
  ON public.integration_webhook_subscriptions (tenant_id, event_type);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_subscriptions_status
  ON public.integration_webhook_subscriptions (status);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_subscriptions_tenant_status
  ON public.integration_webhook_subscriptions (tenant_id, status);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_attempts_status_next
  ON public.integration_webhook_attempts (status, next_delivery_at);

CREATE INDEX IF NOT EXISTS idx_integration_webhook_attempts_subscription
  ON public.integration_webhook_attempts (subscription_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_integration_secrets_tenant_type
  ON public.integration_secrets (tenant_id, integration_type);

CREATE INDEX IF NOT EXISTS idx_integration_secrets_next_due
  ON public.integration_secrets (next_rotation_due_at);

CREATE INDEX IF NOT EXISTS idx_integration_secrets_status_due
  ON public.integration_secrets (status, next_rotation_due_at);

CREATE INDEX IF NOT EXISTS idx_integration_secrets_tenant_status
  ON public.integration_secrets (tenant_id, status);

CREATE INDEX IF NOT EXISTS idx_integration_idempotency_tenant
  ON public.integration_idempotency_records (tenant_id, expires_at);

CREATE INDEX IF NOT EXISTS idx_integration_idempotency_scope_operation
  ON public.integration_idempotency_records (scope, operation);

COMMIT;
