# Research Notes — Support & Operations

## Decision 1: Establish dedicated `operations` bounded context in backend and web-admin
- **Decision**: 新增 `backend/internal/domain/models/operations/**`、`backend/internal/domain/repository/operations/**`、`backend/internal/services/operations/**` 与 `web-admin/app/pages/_p/com.powerx.plugins.base/admin/operations/**` 目录，集中实现支持与运维逻辑。
- **Rationale**: 将 Support/Incident/SLA 能力从 marketplace 与 integration 领域拆分，便于权限与发布节奏独立治理，符合宪章的 Service-Centric 分层。
- **Alternatives considered**:
  - 复用 `integration` 目录：会导致领域职责混淆，测试覆盖难以界定。
  - 拆散到现有 marketplace 服务：运维动作跨 marketplace 之外的租户/插件，耦合不合理。

## Decision 2: Persist support tickets and incidents in normalized tables with checklist snapshots
- **Decision**: 使用三组主表 `operations_support_tickets`、`operations_incidents`、`operations_readiness_checklist_items`，辅以 `operations_incident_updates` 和 `operations_ticket_metrics` 记录事件时间线与 KPI。
- **Rationale**: 结构化存储支持精确计算 FRT/MTTR/SLA Score，并允许审计与历史追溯，满足宪章要求的可观测与合规。
- **Alternatives considered**:
  - 仅写审计日志：难以执行实时 KPI 与激励/处罚；查询复杂。
  - 采用文档存储：与现有 Postgres + RLS 体系不匹配。

## Decision 3: Generate SLA metrics via scheduled aggregation job plus real-time cache
- **Decision**: 使用定时任务（每日/月度/季度）汇总 SLA 指标，写入 `operations_sla_profiles` 表，同时在请求 `GET /api/v1/marketplace/sla/{plugin_id}` 时读取最近快照。
- **Rationale**: 满足 FR-010 的批量采样，同时避免每次请求重算；与既有 job runner 结合简单。
- **Alternatives considered**:
  - 实时查询所有 ticket/incident：高开销，难以满足 API latency。
  - 单纯依赖外部 Support Hub：无法保证插件侧的审计一致性。

## Decision 4: Emit ticket Webhook events via existing event bus abstraction
- **Decision**: 复用 `backend/internal/services/integration/webhook_service.go` 的事件投递管线，新增 `operations.support.ticket.*` 主题，提供 Webhook 与审计双写。
- **Rationale**: 减少重复实现，继承已有重试/签名机制，符合零信任与审计要求。
- **Alternatives considered**:
  - 直接 HTTP 调用 Vendor：缺少重试/签名；实现成本高。
  - 引入新队列服务：超出范围并增加维护成本。

## Decision 5: Admin UI uses existing Nuxt operations layout with checklist components
- **Decision**: 扩展 `web-admin` 的 Admin 菜单，新增 “Operations > Support Playbook”、“Incident Center”、“SLA Dashboard” 页面，复用 ChecklistRunner 组件。
- **Rationale**: 保持 UI/UX 一致，减少新组件开发，降低培训成本。
- **Alternatives considered**:
  - 独立前端应用：破坏宿主路由约束。
  - CLI-only 管理：不符合 Playbook 的可视化需求与审核流程。
