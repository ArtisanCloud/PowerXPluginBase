# Data Model — Support & Operations

## 1. SupportChannelConfig (`operations_support_channels`)
- **id** (UUID, PK)
- **plugin_id** (string, FK → plugin manifest id)
- **tenant_id** (bigint, nullable; `NULL` 表示全局配置)
- **channel** (enum: `marketplace_ticket`, `vendor_email`, `in_app_feedback`, `emergency_hotline`)
- **is_enabled** (bool)
- **service_window** (jsonb: business_hours + timezone + escalation contacts)
- **escalation_path** (jsonb: ordered roles + SLA clocks)
- **metadata** (jsonb: webhook endpoints、邮箱地址、紧急电话等)
- **version** (int, 乐观锁)
- **created_at / updated_at** (timestamptz)
- **created_by / updated_by** (uuid → admin user)

> **Relationships**: `plugin_id` + `tenant_id` scoped by RLS。用于 Support Playbook 页面渲染。

## 2. SupportTicket (`operations_support_tickets`)
- **id** (UUID, PK)
- **plugin_id** (string)
- **tenant_id** (bigint)
- **channel_id** (UUID, FK → SupportChannelConfig)
- **external_ref** (string, optional，用于对接 Support Hub 工单编号)
- **subject / description** (text)
- **priority** (enum: `P0`–`P4`)
- **status** (enum: `created`, `assigned`, `in_progress`, `resolved`, `closed`)
- **requested_by** (jsonb: tenant_contact、vendor_contact)
- **assigned_team** (enum: `agent`, `engineer`, `manager`, `liaison`, `qa`)
- **assigned_to** (uuid)
- **knowledge_base_refs** (text[])
- **first_response_at / resolved_at / closed_at** (timestamptz)
- **csat_score** (numeric(2,1), nullable)
- **resolution_code** (enum: `fixed`, `workaround`, `known_issue`, `rejected`, `duplicate`)
- **reopen_count** (int)
- **created_at / updated_at** (timestamptz)

> **Derived Metrics**: FRT = `first_response_at - created_at`；MTTR = `resolved_at - created_at`。

## 3. SupportTicketEvent (`operations_support_ticket_events`)
- **id** (bigserial, PK)
- **ticket_id** (UUID, FK → SupportTicket)
- **event_type** (enum: `created`, `updated`, `assigned`, `status_changed`, `webhook_dispatched`)
- **payload** (jsonb)
- **emitted_at** (timestamptz)
- **webhook_status** (enum: `pending`, `delivered`, `failed`, `retrying`)
- **retry_count** (int)

> 驱动 Webhook 通知与审计日志。`event_type` 决定 `operations.support.ticket.*` 主题。

## 4. ReadinessChecklistItem (`operations_readiness_checklist_items`)
- **id** (UUID, PK)
- **plugin_id** (string)
- **type** (enum: `support_ready`, `incident_ready`, `sla_ready`)
- **item_key** (string)
- **description** (text)
- **status** (enum: `pending`, `in_progress`, `completed`, `blocked`)
- **owner_role** (enum: `agent`, `engineer`, `manager`, `liaison`, `operations`)
- **due_date** (date)
- **completed_at** (timestamptz)
- **notes** (text)
- **created_at / updated_at** (timestamptz)

> ChecklistRunner 重新使用，发布前需全部 `completed`。

## 5. IncidentRecord (`operations_incidents`)
- **id** (UUID, PK)
- **plugin_id** (string)
- **tenant_id** (bigint, nullable)
- **severity** (enum: `sev0`–`sev4`)
- **status** (enum: `detected`, `acknowledged`, `mitigated`, `monitoring`, `resolved`, `closed`)
- **detected_at / acknowledged_at / mitigated_at / resolved_at / closed_at** (timestamptz)
- **detection_source** (enum: `monitoring`, `support`, `vendor`, `security`, `dependency`)
- **summary** (text)
- **impact** (jsonb: affected tenants, features, revenue estimate)
- **mitigation** (text)
- **root_cause** (text)
- **next_update_at** (timestamptz)
- **labels** (text[]; enforced vocabulary: `#availability`, `#security`, `#performance`, `#dependency`)
- **confidentiality** (enum: `public`, `restricted`, `security-only`)
- **created_at / updated_at** (timestamptz)

## 6. IncidentTimelineEntry (`operations_incident_updates`)
- **id** (bigserial, PK)
- **incident_id** (UUID, FK → IncidentRecord)
- **entry_type** (enum: `announcement`, `update`, `mitigation`, `resolution`, `postmortem`)
- **message** (text)
- **author_role** (enum: `incident_commander`, `liaison`, `security`, `engineering`)
- **posted_at** (timestamptz)
- **stakeholder_channel** (enum: `support_hub`, `status_page`, `security_email`, `hotline`)

> 支持 15 分钟 / 1 小时节奏化通报。

## 7. SLAProfile (`operations_sla_profiles`)
- **id** (UUID, PK)
- **plugin_id** (string)
- **plan_type** (enum: `real_time`, `transactional`, `utility`)
- **uptime_target** (numeric(5,2))
- **uptime_actual** (numeric(5,2))
- **response_target_ms / response_actual_ms** (int)
- **success_target_pct / success_actual_pct** (numeric(5,2))
- **support_frt_target_hours / support_frt_actual_hours** (numeric(4,2))
- **sla_score** (numeric(5,2))
- **incentive_applied_at** (timestamptz, nullable)
- **penalty_applied_at** (timestamptz, nullable)
- **notes** (text)
- **computed_at** (timestamptz)

## 8. SLAAdjustmentHistory (`operations_sla_adjustments`)
- **id** (bigserial, PK)
- **plugin_id** (string)
- **period** (date range)
- **score_before / score_after** (numeric(5,2))
- **action** (enum: `incentive`, `penalty`)
- **details** (text)
- **applied_by** (uuid)
- **created_at** (timestamptz)

## State & Validation Notes

- 所有实体继承租户/插件 RLS 约束：`plugin_id` + `tenant_id` (如适用)。
- 状态机：
  - **SupportTicket.status**：必须遵循 `created → assigned → in_progress → resolved → closed` 或 `created → resolved → closed`，禁止跳级。
  - **IncidentRecord.status**：`detected → acknowledged → mitigated → monitoring → resolved → closed`。复发时保持同一 incident_id 并追加 `IncidentTimelineEntry`。
- SLA Score 计算：`0.4*uptime_actual + 0.3*(support_frt_target_hours / support_frt_actual_hours * 100)`（截断 100）+ `0.3*success_actual_pct`；评分用于激励/处罚。
- Checklist Items 需在发布前全部 `completed`；阻塞项在 Admin UI 显示警告，并阻止 `marketplace` Listing 发布。
- Webhook 负载包含 `ticket_id`, `event_type`, `occurred_at`, `priority`, `status`, `tenant_id`, `plugin_id`, `payload_checksum`。
