# Data Model – Protocols & Integrations

## IntegrationEnvelope
- **Purpose**: 标准化插件与宿主之间的消息封装。
- **Fields**:
  - `message_id` (UUID, PK)
  - `trace_id` (UUID)
  - `correlation_id` (UUID)
  - `tenant_id` (string)
  - `tool_scope` (string)
  - `issued_at` (timestamp)
  - `idempotency_key` (string, nullable)
  - `payload_ref` (string, 预签名 URL 或内联 JSON)
  - `signature` (string, base64 HMAC)
  - `metadata` (JSONB, 包含 adapter、版本、重试次数等)
- **Relationships**: 与 DeliveryAttempt (1:N)
- **Constraints**: `message_id` 唯一；`payload_ref` 若为 URL 必须包含过期时间；`idempotency_key`+`tenant_id` 组合唯一。

## GrantMatrix
- **Purpose**: 维护 ToolScope → 资源/操作/通道 的映射及策略覆盖。
- **Fields**:
  - `id` (UUID, PK)
  - `scope` (string)
  - `channel` (enum: HTTP, GRPC, MCP, WEBHOOK)
  - `resource` (string, e.g., `/api/v1/admin/security/audit-reports`)
  - `action` (string, e.g., GET, POST, CALL)
  - `constraints` (JSONB, 包含速率、租户限制等)
  - `source` (enum: STATIC, OVERRIDE)
  - `version` (integer)
  - `approved_by` (string)
  - `approved_at` (timestamp)
- **Relationships**: 无直接外键，但 Envelope 校验时引用。
- **Constraints**: `scope + channel + resource + action` 唯一；`source=STATIC` 只读。

## WebhookSubscription
- **Purpose**: 描述 Webhook 订阅与交付策略。
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (string)
  - `event_type` (string)
  - `target_url` (string)
  - `secret` (string, encrypted)
  - `retry_policy` (JSONB，默认 1m→5m→15m)
  - `status` (enum: ACTIVE, PAUSED, DISABLED)
  - `created_at` / `updated_at`
- **Relationships**: 与 DeliveryAttempt (1:N)；与 SecretCredential (可选 1:1)
- **Constraints**: `tenant_id + event_type + target_url` 唯一；`target_url` 必须为 HTTPS。

## DeliveryAttempt
- **Purpose**: 跟踪每次 Webhook/事件投递的状态。
- **Fields**:
  - `id` (UUID, PK)
  - `subscription_id` (UUID, FK -> WebhookSubscription)
  - `envelope_id` (UUID, FK -> IntegrationEnvelope)
  - `status` (enum: PENDING, RETRYING, SUCCEEDED, FAILED, DLQ)
  - `retry_count` (integer)
  - `last_error` (text)
  - `next_delivery_at` (timestamp)
  - `updated_at`
- **Constraints**: `subscription_id` 必须存在；`retry_count` ≥0；DLQ 状态需记录 `last_error`。

## SecretCredential
- **Purpose**: 管理外部 API/Webhook 凭证的生命周期。
- **Fields**:
  - `id` (UUID, PK)
  - `tenant_id` (string)
  - `integration_type` (enum: WEBHOOK, API, OTHER)
  - `current_secret_ref` (string，指向 Secrets Manager)
  - `pending_secret_ref` (string, nullable)
  - `rotation_interval_days` (integer, default 30)
  - `last_rotated_at` (timestamp, nullable)
  - `next_rotation_due_at` (timestamp)
  - `status` (enum: ACTIVE, ROTATING, REVOKED)
  - `audit_log` (JSONB, 记录操作历史)
- **Relationships**: 与 WebhookSubscription 可选关联；与 GrantMatrix 无直接关系。
- **Constraints**: `tenant_id + integration_type` 可建唯一索引；`status=ACTIVE` 时 `current_secret_ref` 必填。
