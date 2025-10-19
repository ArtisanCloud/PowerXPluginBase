# Tasks: Protocols & Integrations (A2A, HTTP/gRPC/MCP, Webhooks & Events, Secrets)

**Input**: Design docs from `/specs/005-protocols-integrations/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: 按《PowerX Plugin Constitution》之 *Observable & Testable Delivery*，在完成各故事实现后执行 `go test ./...`、Webhook 重试演练、Nuxt e2e smoke，以验证端到端场景。规范未要求 TDD，因此测试任务嵌入在各 phase 的实施说明中，而非独立 TDD 步骤。

**Organization**: Tasks 按用户故事分组，确保每个故事都能独立交付、独立验收。

> **格式**: `[ID] [P?] [Story] 描述`
>
> - **[P]**: 可并行（不同文件、无依赖）
> - **[Story]**: 对应用户故事（US1/US2/US3）
> - 必须给出精确文件路径（符合 `plan.md` 的结构）

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: 创建 integration 目录结构、安装依赖、准备配置和本地工具。

- [X] T001 [Setup] 根据 `plan.md` 创建后端/前端目录骨架  
      `backend/internal/services/integration/`, `backend/internal/domain/repository/integration/`, `backend/internal/observability/integration/`, `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/`
- [X] T002 [P] [Setup] 安装/更新 Go module 与 Nuxt 依赖（含脚本）  
      `make dev-setup`, `npm install --prefix web-admin`
- [X] T003 [Setup] 初始化 integration 配置项（1 MB 阈值、重试策略、轮换周期）  
      `backend/etc/config.yaml`, `backend/internal/config/config.go`
- [X] T004 [P] [Setup] 准备幂等缓存后端（优先 Redis，PostgreSQL 回退）并配置 docker-compose  
      `config/docker-compose.integration.yml`
- [X] T005 [Setup] 更新 `AGENTS.md` / 团队 README 说明新的语言、依赖及测试入口  
      `AGENTS.md`

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: 所有用户故事共享的核心基础能力，必须完成后才可进入故事开发。

- [X] T010 [Foundational] 创建数据库迁移（GrantMatrix override、Webhook attempts/DLQ、Secrets）  
      `backend/migrations/2025Q4_integration.sql`
- [X] T011 [Foundational] 加载 integration 配置结构体，暴露 `payload_threshold_bytes`、`retry_policy`、`rotation_days`  
      `backend/internal/config/integration.go`
- [X] T012 [Foundational] 抽象幂等存储接口与 Redis/Postgres 实现  
      `backend/internal/domain/repository/integration/idempotency_provider.go`
- [X] T013 [Foundational] 实现 GrantMatrix 静态 YAML + DB override 加载器及缓存失效机制  
      `backend/internal/services/integration/grant_matrix_loader.go`
- [X] T014 [P] [Foundational] 初始化 integration 观测指标骨架（counter/gauge/events）  
      `backend/internal/observability/integration/metrics.go`
- [X] T015 [P] [Foundational] 注册 integration HTTP/GRPC 路由、RBAC 权限占位  
      `backend/internal/transport/http/routes.go`, `backend/internal/transport/grpc/server.go`, `backend/internal/transport/http/rbac.go`
- [X] T016 [Foundational] 构建后台调度框架（用于 Webhook 重试、Secrets 轮换提醒）  
      `backend/internal/jobs/integration/scheduler.go`
- [X] T017 [Foundational] 更新 Quickstart 配置样例与脚本占位（Redis、webhook mock）  
      `scripts/mock-webhook-target.sh`, `specs/005-protocols-integrations/quickstart.md`
- [X] T018 [Foundational] 实现配置审批/双人复核工作流（仓储、服务、审批记录迁移）  
      `backend/migrations/2025Q4_integration_approvals.sql`, `backend/internal/domain/repository/integration/approval_repository.go`, `backend/internal/services/integration/approval_service.go`

**Checkpoint**: 基础设施就绪，用户故事可独立启动。

---

## Phase 3 — User Story 1 — 统一 A2A 协议投递 (Priority: P1) 🎯 MVP

**Goal**: 提供统一 Envelope、ToolScope 校验、HTTP/gRPC/MCP 适配器，打造首次可演示的 A2A 请求/响应闭环。  
**Independent Test**: 使用 `/dispatch` 发送带 ToolScope 的 Envelope，验证幂等、GrantMatrix 拒绝路径，以及 HTTP/gRPC/MCP 三通道可互换返回同一 trace。

### Implementation for US1

- [X] T100 [US1] 建模 IntegrationEnvelope、DTO 与验证逻辑  
      `backend/internal/domain/models/integration/envelope.go`, `backend/internal/transport/http/integration/dto_envelope.go`
- [X] T101 [P] [US1] 实现幂等记录仓储（Redis + Postgres 回退）  
      `backend/internal/domain/repository/integration/idempotency_repository.go`
- [X] T102 [P] [US1] 实现 GrantMatrix 查询与 scope 校验服务  
      `backend/internal/services/integration/grant_matrix_service.go`
- [X] T103 [US1] 实现 DispatchService（Envelope 校验 → 幂等检测 → 调用宿主）  
      `backend/internal/services/integration/dispatch_service.go`
- [X] T104 [US1] HTTP Handler：解析请求、调用服务、写入观测指标  
      `backend/internal/transport/http/integration/dispatch_handler.go`
- [X] T105 [P] [US1] gRPC Server：映射同一服务层  
      `backend/internal/transport/grpc/integration/dispatch_server.go`
- [X] T106 [P] [US1] MCP Session 适配（握手验证 ToolGrant，建立流式上下文）  
      `backend/internal/mcp/integration/session_adapter.go`
- [X] T107 [US1] 将 dispatch 路由/RBAC/GrantMatrix 元数据写入 OpenAPI 与配置  
      `specs/005-protocols-integrations/contracts/integration-openapi.yaml`, `backend/internal/transport/http/routes.go`, `backend/internal/shared/app/rbac.go`
- [X] T108 [US1] 观测指标与日志增强：请求计数、错误分类、幂等重放标记  
      `backend/internal/observability/integration/metrics.go`, `backend/internal/services/integration/dispatch_service.go`
- [X] T109 [US1] 更新 Quickstart：添加 Envelope 调用与校验步骤  
      `specs/005-protocols-integrations/quickstart.md`

**Checkpoint**: `/dispatch` + ToolScope 校验可独立演示；即完成 MVP。

---

## Phase 4 — User Story 2 — Webhook 与事件可靠交付 (Priority: P2)

**Goal**: 统一 Webhook 订阅、重试、DLQ 协作流程，并提供运营可见性。  
**Independent Test**: 创建订阅→模拟目标不可用→验证退避重试/最终送达；触发连续失败进入 DLQ→联合处理→成功补发。

### Implementation for US2

- [ ] T200 [US2] 建模 WebhookSubscription、DeliveryAttempt 实体与映射  
      `backend/internal/domain/models/integration/webhook_subscription.go`, `backend/internal/domain/models/integration/delivery_attempt.go`
- [ ] T201 [P] [US2] 实现订阅仓储 CRUD（含签名密钥加密）  
      `backend/internal/domain/repository/integration/webhook_subscription_repository.go`
- [ ] T202 [P] [US2] 实现投递尝试仓储（重试状态、DLQ 标记）  
      `backend/internal/domain/repository/integration/delivery_attempt_repository.go`
- [ ] T203 [US2] 实现 WebhookService：创建/更新订阅、生成签名、记录投递结果  
      `backend/internal/services/integration/webhook_service.go`
- [ ] T204 [US2] 实现重试/调度器（引用 Phase 2 scheduler），支持 1m→5m→15m 退避、DLQ 入列  
      `backend/internal/jobs/integration/webhook_retry_worker.go`
- [ ] T205 [US2] 管理端 HTTP API：订阅 CRUD、DLQ replay、统计视图  
      `backend/internal/transport/http/admin/integration/webhook_handler.go`
- [ ] T206 [P] [US2] 观测指标与告警：成功率、重试次数、DLQ 积压  
      `backend/internal/observability/integration/metrics.go`, `backend/internal/observability/integration/alerts.md`
- [ ] T207 [P] [US2] 前端管理页面（列表、详情、DLQ 操作）  
      `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/webhooks.vue`, `web-admin/app/types/integration.ts`
- [ ] T208 [US2] 更新 EventBus 集成（投递失败/成功事件）、Quickstart 演练文档  
      `backend/internal/observability/integration/events.go`, `specs/005-protocols-integrations/quickstart.md`
- [ ] T209 [US2] 将订阅创建/变更纳入审批流程（提交/审批/审计 UI）  
      `backend/internal/services/integration/webhook_service.go`, `backend/internal/transport/http/admin/integration/webhook_handler.go`, `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/webhooks.vue`

**Checkpoint**: Webhook 生命周期可独立运行；运营可恢复失败通知。

---

## Phase 5 — User Story 3 — 外部 Secrets 生命周期治理 (Priority: P3)

**Goal**: 管理外部 API 凭证的创建、轮换、吊销与审计，支持租户隔离。  
**Independent Test**: 创建凭证→触发轮换→验证双密钥宽限→审计日志可查询；吊销后相关调用被拒绝并有告警。

### Implementation for US3

- [ ] T300 [US3] 建模 SecretCredential、审计记录结构  
      `backend/internal/domain/models/integration/secret_credential.go`
- [ ] T301 [P] [US3] 实现 Secrets 仓储（CRUD、轮换计划、挂账状态）  
      `backend/internal/domain/repository/integration/secret_repository.go`
- [ ] T302 [US3] 实现 SecretsService：创建/轮换/吊销、触发审计事件  
      `backend/internal/services/integration/secret_service.go`
- [ ] T303 [P] [US3] 集成宿主 Secrets Manager（STS 凭证交换、双密钥切换）  
      `backend/internal/services/integration/secret_provider.go`
- [ ] T304 [US3] HTTP Admin API：Secrets 管理、轮换触发、审计查询  
      `backend/internal/transport/http/admin/integration/secret_handler.go`
- [ ] T305 [P] [US3] 定时轮换提醒任务 + 告警（复用 scheduler）  
      `backend/internal/jobs/integration/secret_rotation_worker.go`
- [ ] T306 [US3] 管理前端页面（凭证列表、轮换、吊销操作、历史记录）  
      `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/secrets.vue`, `web-admin/app/types/integration.ts`
- [ ] T307 [US3] 更新 Quickstart 与运维 runbook（轮换流程、告警确认）  
      `specs/005-protocols-integrations/quickstart.md`, `docs/security/vulnerability-response.md`
- [ ] T308 [US3] 将凭证创建/轮换/吊销操作纳入审批流程并记录审计  
      `backend/internal/services/integration/secret_service.go`, `backend/internal/transport/http/admin/integration/secret_handler.go`, `backend/internal/domain/repository/integration/approval_repository.go`

**Checkpoint**: Secrets 生命周期闭环可独立验收。

---

## Phase 6 — Polish & Cross-Cutting

- [ ] T900 [P] Polish 更新文档：README、docs/security、plan/research 补充最终决策  
      `README.md`, `docs/security/integration.md`
- [ ] T901 Polish 性能调优与索引（幂等键、DeliveryAttempt、Secrets 查询）  
      `backend/migrations/2025Q4_integration_indexes.sql`
- [ ] T902 Polish 发布产物 & release notes（OpenAPI、Nuxt 构建、manifest 更新）  
      `docs/releases/2025-10-integrations.md`, `plugin.yaml`
- [ ] T903 Polish 全量测试与演练（go test、webhook replay、Nuxt build）  
      `Makefile`, `scripts/ci/integration.sh`
- [ ] T904 Polish 安全与合规复核（GrantMatrix 覆盖、RBAC、审计日志抽样）  
      `docs/security/audit-logs.md`, `backend/internal/shared/app/rbac.go`
- [ ] T905 [P] Polish 构建仪表盘/报表（Envelope 采用率、Webhook 成功率、Secrets 轮换进度等）  
      `docs/observability/integration-dashboard.json`, `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/insights.vue`, `specs/005-protocols-integrations/quickstart.md`
- [ ] T906 Polish 编写成功指标验证脚本与发布前检查（SC-001~SC-005）  
      `scripts/ci/verify_integration_metrics.sh`, `docs/observability/integration-checklist.md`

---

## Dependencies & Execution Order

- **Phase Order**: Setup → Foundational → US1 (P1) → US2 (P2) → US3 (P3) → Polish。  
  US2、US3 依赖 GrantMatrix/幂等基础，因此需等待 Foundational，同时建议在 US1 完成后再继续，以复用 Envelope 验证与指标。
- **Story Dependencies**:
  - US2 ← Foundational + US1（重用 Envelope、GrantMatrix 校验）
  - US3 ← Foundational（Secrets 生命周期独立，但需现有配置与 observer）

---

## Parallel Execution Examples

- **Phase 1**: T002 ↔ T004 可并行（不同文件/服务）。
- **Foundational**: T013 ↔ T014 ↔ T015 可由不同工程师同步推进。
- **US1**:  
  - T101 ↔ T102 ↔ T106 可并行（仓储/GrantMatrix/MCP 属不同文件）。  
  - T104 与 T105 在 Service (T103) 完成后可并行。
- **US2**:  
  - T201 ↔ T202 ↔ T207 可并行。  
  - T204 需待仓储准备好。
- **US3**:  
  - T301 ↔ T303 ↔ T306 可并行；T304 依赖服务 (T302)。
- **Polish**: T900、T901、T902、T903 大多互不依赖，可按职责拆分。

---

## Implementation Strategy

1. **MVP** = 完成 US1：统一 Envelope + ToolScope 校验 + 三通道适配，确保最早可演示插件与宿主互通。  
2. 在 MVP 稳定后，扩展 Webhook 重试（US2），提供可靠事件通知。  
3. 最后实现 Secrets 生命周期（US3），补足合规要求。  
4. Polish 阶段合并所有文档、release notes、性能调优，准备发布。
