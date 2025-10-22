# Tasks: Support & Operations (Support Playbook, Incident Handling, SLA/SLO)

**Input**: spec.md, plan.md, research.md, data-model.md, contracts/operations-openapi.yaml, quickstart.md  
**Prerequisites**: `make dev-setup`, `cd web-admin && npm install`, operations Redis/Webhook mocks running via `docker compose -f config/docker-compose.integration.yml up -d`

> Format: `[ID] [P?] [Story] Description`

---

## Phase 1 — Setup (Shared Infrastructure)

- [x] T001 [P] 运行 `(make dev-setup) && (cd web-admin && npm install)`，确认 Go/Nuxt 依赖、`golangci-lint` 与测试工具就绪。  
- [x] T002 [P] 初始化 `backend/internal/{domain,repository,services}/operations/` 与 `web-admin/app/{pages,components,stores}/operations/` 目录，添加 README 或 package 文档说明领域边界。  

---

## Phase 2 — Foundational (Blocking Prerequisites)

- [x] T010 更新 RBAC 与配置：在 `plugin.yaml`、`backend/internal/transport/http/registry.go`、`backend/internal/router/router.go` 中预留 `operations.support.*`、`operations.incident.*`、`operations.sla.*` 策略与路由骨架。  
- [x] T011 扩展配置与常量：在 `backend/etc/config.example.yaml`、`backend/internal/config/config.go` 增加 `operations` 节点（渠道、Webhook 签名、SLA 周期），同步 `docs/readme.md` 说明。  
- [x] T012 建立观测基线：在 `backend/internal/observability/operations/metrics.go` 定义票务/事故/SLA 指标占位，注册于 `internal/shared/app/deps.go`。  
- [x] T013 校准 Checklist 框架：扩展 `backend/internal/services/admin/runtime_ops/quota_service_test.go` 等 readiness 逻辑，支持 `support_ready`、`incident_ready`、`sla_ready` 三套清单骨架。  

---

## Phase 3 — User Story 1 — Support Playbook & Channels (Priority: P1) 🎯 MVP

**Goal**: 支持经理可在 Admin Console 配置多渠道支持体系，完成 Support Ready Checklist，并触发 `ticket.created` Webhook。  
**Independent Test**: 管理端完成渠道配置 → Checklist 全部通过 → 提交模拟 P2 工单 → 收到 `operations.support.ticket.created` Webhook → KPI 仪表盘展示 FRT/MTTR。  

### Implementation for US1

- [x] T020 [US1] 创建迁移：`backend/migrations/2025Q4_operations_support.sql` 定义 `operations_support_channels`、`operations_support_tickets`、`operations_support_ticket_events`、Checklist 相关表及 RLS。  
- [x] T021 [US1] 建模：在 `backend/internal/domain/models/operations/` 创建 `support_channel.go`、`support_ticket.go`，声明 GORM 标签与状态机。  
- [x] T022 [US1] 实现仓储：在 `backend/internal/domain/repository/operations/support_repository.go` 提供渠道 CRUD、票务状态流转、事件写入接口。  
- [x] T023 [US1] 服务层：在 `backend/internal/services/operations/support_service.go` 封装渠道配置校验、Support Ready Checklist 更新、工单事件出站逻辑。  
- [x] T024 [US1] Webhook 整合：复用 `internal/services/integration/webhook_service.go`，新增 `operations.support.ticket.*` 主题与负载签名。  
- [x] T025 [US1] Admin API：根据 contracts 在 `backend/internal/transport/http/admin/operations/support_handler.go` 实现 `/support/playbook`、`/support/metrics`、`/support/channels/test`。  
- [x] T026 [US1] 路由集成：更新 `backend/internal/transport/http/admin/routes.go`、`registry.go` 注册新 handler，补充中间件校验。  
- [x] T027 [US1] 就绪清单：在 `backend/internal/services/admin/runtime_ops` 与 `backend/internal/transport/http/admin/runtime_ops` 中增加 Support Ready 项与 API。  
- [x] T028 [US1] Go 测试：编写 `internal/services/operations/support_service_test.go`、`backend/tests/integration/operations/support_flow_test.go` 覆盖渠道配置、Webhook 发送、Checklist 完成。  
- [x] T029 [US1] Admin UI 页面：实现 `web-admin/app/pages/_p/com.powerx.plugins.base/admin/operations/support.vue`、相关组件与 store，支持渠道配置、知识库、校验。  
- [x] T030 [US1] 前端测试：在 `web-admin/tests/operations/support_playbook.spec.ts` 验证表单校验、Webhook 状态提示与 Checklist 渲染。  
- [x] T031 [US1] 审计日志：在 `backend/internal/services/operations/support_service.go` 写入渠道变更与工单状态审计事件，补充 `support_service_audit_test.go` 覆盖日志内容。  

**Checkpoint**：Support Playbook 生效，Checklist 绿色，Webhook 可靠。

---

## Phase 4 — User Story 2 — Incident Response Lifecycle (Priority: P2)

**Goal**: Incident Commander 可按 SEV 矩阵声明/同步事故，维持时间线及 RCA，Incident Ready Checklist 通过。  
**Independent Test**: 模拟 SEV-1 事件 → 15 分钟内通报 → 按阶段更新状态 → 48 小时内提交 RCA → Incident Ready Checklist 完成。  

### Implementation for US2

- [ ] T040 [US2] 迁移：扩展同批 SQL 或新增 `operations_incidents.sql`，建立 `operations_incidents`、`operations_incident_updates`、`operations_incident_checklist`。  
- [ ] T041 [US2] 建模：在 `domain/operations/models/incident.go` 表达 SEV/状态流转、时间戳、标签、保密级别。  
- [ ] T042 [US2] 仓储：`repository/incident_repository.go` 支持创建、状态更新、时间线、标签检索与分页过滤。  
- [ ] T043 [US2] 服务层：`services/operations/incident_service.go` 实现 SEV 驱动 SLA 时钟、通报计划、事件复发处理。  
- [ ] T044 [US2] 通知整合：新增 incident 通知通道，向 Support Hub / Hotline / security 邮箱推送（扩展 webhook service & email client）。  
- [ ] T045 [US2] Admin API：在 `transport/http/admin/operations/incident_handler.go` 实现 `/incidents` CRUD、`/{incidentId}/timeline`，校验 SEV 响应时间。  
- [ ] T046 [US2] 状态页同步：扩展 `transport/http/admin/runtime_ops`/外部 SDK，确保状态页数据写入（若仅触发 webhook，记录说明）。  
- [ ] T047 [US2] Checklist：实现 Incident Ready Checklist 项与阻断逻辑（复用 readiness 服务）。  
- [ ] T048 [US2] Go 测试：编写 `incident_service_test.go`、`tests/integration/operations/incident_flow_test.go` 覆盖 SEV 响应、时间线、RCA。  
- [ ] T049 [US2] Admin UI：实现 `web-admin/app/pages/_p/.../operations/incidents.vue`、`IncidentTimeline.vue` 等组件，支持标签、通报计划、RCA 上传。  
- [ ] T050 [US2] 前端测试：`web-admin/tests/operations/incident_flow.spec.ts` 验证 SEV 流程、时间线、Checklist 状态。  
- [ ] T051 [US2] 审计日志：在 `services/operations/incident_service.go` 与通知通道记录 SEV 升级、时间线更新、RCA 提交的审计事件，新增 `incident_audit_test.go`。  

**Checkpoint**：Incident 流程闭环，通报与 RCA 达标。

---

## Phase 5 — User Story 3 — SLA Transparency & Incentives (Priority: P3)

**Goal**: 运营经理可查看/调整 SLA 目标，发布公共 SLA API，并自动执行激励/处罚。  
**Independent Test**: 月度作业生成 SLA 快照 → Admin Dashboard 与 API 返回一致 → Score≥85 触发推荐位，Score<70 触发处罚。  

### Implementation for US3

- [X] T060 [US3] 迁移：创建/扩展 SQL 加入 `operations_sla_profiles`、`operations_sla_adjustments`、聚合辅助表。  
- [X] T061 [US3] 建模：在 `models/sla_profile.go` 定义计划类型、目标/实际指标、score 字段。  
- [X] T062 [US3] 仓储：`repository/sla_repository.go` 负责快照查询、调整历史记录。  
- [X] T063 [US3] 服务层：`services/operations/sla_service.go` 计算 Score（≥85 激励，<70 处罚）、写入调整历史、驱动 Dashboard。  
- [X] T064 [US3] 定时任务：实现 `services/operations/jobs/sla_recompute_job.go` & 调度入口（cmd/cron 或 existing runner），每日/月度/季度聚合支持。  
- [X] T065 [US3] Admin API：在 `transport/http/admin/operations/sla_handler.go` 实现 `/sla/profiles`、`/sla/profiles/recompute`。  
- [X] T066 [US3] 公共 API：在 `transport/http/public/marketplace/sla_handler.go` 提供 `GET /api/v1/marketplace/sla/{plugin_id}`，含缓存与 404 处理。  
- [X] T067 [US3] Dashboard UI：实现 `operations/sla.vue`、`SlaScoreCard.vue`，展示激励/处罚与趋势。  
- [X] T068 [US3] Tests：`sla_service_test.go`、`tests/integration/operations/sla_refresh_test.go` 校验计算、API 一致性与激励/处罚触发。  
- [X] T069 [US3] 前端测试：`web-admin/tests/operations/sla_dashboard.spec.ts` 覆盖指标显示、筛选、Badge 更新。  
- [X] T070 [US3] 审计日志：在 `services/operations/sla_service.go` 与 `sla_recompute_job.go` 写入 SLA 调整、激励/处罚执行的审计事件，补充 `sla_audit_test.go`。  

**Checkpoint**：SLA 透明可查，激励/处罚自动落地。

---

## Phase 6 — Polish & Cross-Cutting

- [X] T090 更新文档：同步 `docs/overview/marketplace_business_loop.md`、`quickstart.md`、新增操作指南，说明 Support/Incident/SLA 演练。  
- [X] T091 清理临时代码并运行验证：`make fmt && make lint && make test && npm run lint && npm run test -- operations`。  
- [X] T092 校验打包：执行 `make build && make frontend-build && make dist`，检查 dist 包含新的 operations 资源。  
- [X] T093 审核审计与指标：确认 audit log、metrics、Webhook 重试配置均正确，更新 `docs/references/changelog.md` 与 `plugin.yaml` 版本号。  

---

## Dependencies & Execution Order

1. **Phase 1 → Phase 2 → Phase 3/4/5 → Phase 6**  
2. User story顺序：US1 (P1) → US2 (P2) → US3 (P3)。US2 依赖已存在的 Support Ready 骨架，US3 依赖 Incident/SLA 数据模型与渠道事件。  
3. 带 `[P]` 的任务可在保持文件独立时并行推进，若存在同文件修改需串行合并。

### Parallel Execution Examples

- Phase 3：T021/T022/T024 可与 T029/T030 并行；T025 完成后可并行推进 T026/T027。  
- Phase 4：T041/T042/T044 同步推进；T049 与 T050 可由前端团队并行。  
- Phase 5：T061/T062/T064 可并行；T067 与 T069 并行。  
- Phase 6：T090 与 T091 可并行（文档 vs 验证）。

---

## Implementation Strategy

1. **MVP（US1）**：优先交付 Support Playbook + Webhook 流程，确保上线审核与支持渠道可用。  
2. **Incident Lifecycle（US2）**：在支持体系稳定后补全事故响应，保障 SEV 流程与 RCA 追踪。  
3. **SLA Transparency（US3）**：最后实现 SLA 量化与激励/处罚，完成运营闭环。  
4. 每个阶段结束执行 `make test`、`npm run test -- operations` 与合约测试，确保回归稳定。  

---

## Task Tally

- **Total tasks**: 35  
  - Setup: 2  
  - Foundational: 4  
  - US1: 12  
  - US2: 12  
  - US3: 11  
  - Polish: 4  
- **Parallelizable tasks** (`[P]`): 6  
  - Setup/Foundational: 2  
  - US1: 2  
  - US2: 1  
  - US3: 1  
  - Polish: 0  
- **Independent Test Criteria**:  
  - US1：Support Playbook 完成 & Webhook 验证（T028–T031）。  
  - US2：SEV-1 演练与 Incident Ready Checklist（T048–T051）。  
  - US3：SLA API/Dashboard 一致与激励/处罚触发（T068–T070）。  
- **Suggested MVP Scope**: 完成 Phase 1–3，交付 Support Playbook & Channels。
