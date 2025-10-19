# Data Model – Marketplace & Business

## Schema

所有表位于 `powerx_plugin_base` schema，启用多租户隔离（`tenant_id`）与审计列（`created_at`, `updated_at`, `deleted_at`）。主键除特别说明外均为 `uuid`.

## Entities

### marketplace_listings
- **Fields**: `id`, `tenant_id`, `plugin_id`, `vendor_id`, `status` (`draft|in_review|published|suspended`), `title`, `slug`, `summary`, `description`, `cover_asset_id`, `hero_video_asset_id`, `categories` (string[]), `tags` (string[]), `locale`, `version`, `ready_checklist_score`, `recommended_weight`, `published_at`, `reviewed_at`, `reviewer_id`, `audit_notes`, `branding_theme` (jsonb).
- **Indexes**: `(tenant_id, plugin_id, status)`, `slug` unique per locale.
- **Relations**: 1:N → `listing_assets`, `pricing_plans`, `checklist_runs`, `listing_versions`.
- **Rules**: published listings必须具备完成的 checklist；状态切换会写入 `listing_versions`.
- **State transitions**: `draft → in_review → published → suspended`; `suspended → in_review` 复审。

### listing_assets
- **Fields**: `id`, `listing_id`, `tenant_id`, `asset_type` (`logo|cover|screenshot|video`), `storage_uri`, `checksum`, `is_primary`, `locale`, `weight`, `metadata` (jsonb: resolution, duration, size).
- **Indexes**: `(listing_id, asset_type)`.
- **Rules**: 视频资产（`asset_type=video`）约束为时长 ≤ 15s/50MB，格式 MP4/WebM。

### listing_versions
- **Fields**: `id`, `listing_id`, `tenant_id`, `version`, `changelog`, `metadata`, `submitted_by`, `review_state`, `reviewer_id`, `reviewed_at`.
- **Rules**: 保存每次版本提交与审核结果；`review_state` 关联主 Listing 状态。

### pricing_plans
- **Fields**: `id`, `listing_id`, `tenant_id`, `plan_code`, `plan_type` (`free|one_time|subscription|usage`), `currency`, `amount`, `billing_period` (`monthly|yearly|per_use|lifetime`), `trial_period_days`, `quota_limit`, `overage_policy`, `feature_matrix` (jsonb), `is_default`.
- **Indexes**: `(listing_id, plan_code)` unique.
- **Rules**: 与 Billing Engine 同步 `plan_id`; `currency` 需与税务 SaaS 支持列表匹配。

### plan_tiers (usage-based)
- **Fields**: `id`, `plan_id`, `tenant_id`, `metric`, `from`, `to`, `unit_amount`, `unit_name`.
- **Rules**: 与 `pricing_plans.plan_type=usage` 关联；有序不重叠。

### licenses
- **Fields**: `id`, `tenant_id`, `listing_id`, `plan_id`, `license_token` (JWT blob), `status` (`active|trial|revoked|expired|suspended`), `issued_at`, `expires_at`, `renewal_token`, `offline_until`, `last_validated_at`, `issued_by`.
- **Indexes**: `(tenant_id, listing_id)`, `(plan_id, status)`.
- **Rules**: `offline_until` ≤ `issued_at + 72h`；`renewal_token` 一次性使用。

### license_events
- **Fields**: `id`, `tenant_id`, `license_id`, `event_type` (`issued|renewed|revoked|usage_reported|trial_expired|offline_extend`), `event_payload` (jsonb), `emitted_at`, `actor_id`, `trace_id`.
- **Indexes**: `(license_id, event_type, emitted_at)`.
- **Rules**: 所有 License 状态变更与异常需写入；对接审计与通知。

### usage_envelopes
- **Fields**: `id`, `tenant_id`, `license_id`, `plugin_id`, `metrics` (jsonb array `{metric_name, unit, value}`), `timestamp_start`, `timestamp_end`, `ingested_at`, `signature`, `checksum`, `ingest_status` (`pending|processed|replayed`).
- **Indexes**: `(license_id, timestamp_start)`, `(ingest_status)`.
- **Rules**: 保留 180 天；幂等键由 `checksum` + `signature` 保证，断流补偿写入 `ingest_status=replayed`。

### usage_aggregates
- **Fields**: `id`, `tenant_id`, `license_id`, `time_bucket` (hour/day/month), `metric`, `total`, `delta`, `currency`, `revenue`, `aggregation_window`.
- **Rules**: Analytics Pipeline 产出，用于 Dashboard 与配额核对。

### revenue_share_reports
- **Fields**: `id`, `tenant_id`, `vendor_id`, `period_start`, `period_end`, `gross_amount`, `vendor_share`, `platform_share`, `fees`, `currency`, `status` (`draft|ready|exported|reconciled`), `generated_at`, `export_uri`.
- **Indexes**: `(vendor_id, period_start, period_end)`.
- **Rules**: Vendor 80% / Platform 15% / Fee 5%；出口 CSV/PDF。

### checklist_runs
- **Fields**: `id`, `listing_id`, `tenant_id`, `trigger_source` (`vendor|ci|auto`), `run_number`, `status` (`pending|passed|failed`), `started_at`, `completed_at`, `summary`, `ci_pipeline_id`.
- **Relations**: 1:N → `checklist_items`.

### checklist_items
- **Fields**: `id`, `checklist_run_id`, `tenant_id`, `code`, `description`, `result` (`passed|failed|warning`), `evidence_uri`, `notes`, `auto_fix_link`.
- **Rules**: GraphQL 查询/变更目标；`code` 对应 spec 中的 Ready Checklist 条目。

### tax_transactions
- **Fields**: `id`, `tenant_id`, `billing_id`, `external_provider` (`stripe_tax|avalara`), `external_transaction_id`, `jurisdiction`, `tax_amount`, `currency`, `raw_payload` (jsonb), `status` (`pending|completed|failed`), `synced_at`.
- **Rules**: 记录 SaaS 访问与 retries，失败需可重放。

### notifications (marketplace scoped)
- **Fields**: `id`, `tenant_id`, `recipient_type` (`vendor|tenant|platform`), `recipient_id`, `channel` (`email|webhook|in_app`), `template_code`, `payload`, `scheduled_at`, `sent_at`, `status`.
- **Purpose**: 支撑续费提醒、usage spike 通知。

## Shared Enumerations
- `status_listing`: `draft`, `in_review`, `published`, `suspended`.
- `license_status`: `trial`, `active`, `expired`, `revoked`, `suspended`.
- `checklist_trigger`: `vendor`, `ci`, `auto`.
- `aggregation_window`: `hour`, `day`, `month`.

## Data Retention
- `usage_envelopes`: 180 天后归档/删除，遵守 GDPR 删除请求 → 批量标记并触发清理任务。
- `license_events`, `revenue_share_reports`: 永久保留（合规审计），允许逻辑删除但不可物理移除。

## Tenancy & RLS
- 每个表包含 `tenant_id` 并通过 Repository `BeginTenantTx` 设置 `app.tenant_id`。
- 跨租户只允许 Platform 角色通过视图访问聚合数据，具体通过 `tenant_id in (requested|platform_overrides)` 的安全策略实现。

## Audit & Observability
- 关键表（`licenses`, `license_events`, `tax_transactions`, `revenue_share_reports`）需启用结构化审计日志，写入 `observability/integration` 的事件流。
- 数据变更触发 `EventBus` 推送，供通知与分析流程消费。
