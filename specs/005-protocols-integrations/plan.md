# Implementation Plan: Protocols & Integrations (A2A, HTTP/gRPC/MCP, Webhooks & Events, Secrets)

**Branch**: `005-protocols-integrations` | **Date**: 2025-10-17 | **Spec**: specs/005-protocols-integrations/spec.md  
**Input**: Feature specification from `/specs/005-protocols-integrations/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

---

## Summary

为 PowerX 插件提供统一的集成协议：标准化 A2A Envelope、HTTP/gRPC/MCP 适配器、GrantMatrix 权限校验、Webhook 重试与 DLQ 协作，以及外部 Secrets 生命周期治理。实现将以现有 Go/Nuxt 代码结构为基础，新增 integration 服务层、观测指标命名空间、静态 YAML + 数据库覆盖的 GrantMatrix 管理，并引入幂等缓存与事件重试流水线。

---

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24（后端）、Node 20 + Nuxt 4（管理前端）  
**Primary Dependencies**: Gin、PowerX gRPC SDK、MCP client 库、PostgreSQL、可选 Redis（幂等缓存/队列）、PowerX EventBus  
**Storage**: PostgreSQL `powerx_plugin_base` schema；新增表 `integration_grant_matrix_overrides`、`integration_webhook_attempts`、`integration_webhook_dlq`、`integration_secrets`  
**Testing**: `go test ./...`（unit/service/repository）、integration tests（模拟适配器和重试）、contract tests（OpenAPI + Webhook schema）、Nuxt e2e smoke  
**Target Platform**: PowerX 托管容器（Linux amd64）；反代路径 `/api/v1/**` 与 `/_p/<plugin-id>/admin/**`  
**Project Type**: Go 后端 + Nuxt 管理前端  
**Performance Goals**: Adapter 开销 <5ms p95；Webhook 成功率 ≥99%，平均重试≤2；Secrets 轮换响应 <5s；DLQ 消息 5 分钟内被监控发现  
**Constraints**: TLS 1.3 强制；ToolGrant TTL ≤24h；幂等键有效期 24 小时；DLQ 保留 30 天；Observability 统一在 `internal/observability/integration`  
**Scale/Scope**: 覆盖 1000+ 插件实例；Webhook 日发送 1M；GrantMatrix 管理 200+ ToolScope；Secrets 关联 500+ 外部系统

### Platform / Hosting Integration

- **Reverse Proxy & Routes**: 业务端点统一在 `/api/v1/integration/**`；管理 UI 页面位于 `/_p/<plugin-id>/admin/integration/**`；OpenAPI 与 Webhook 文档同步发布。  
- **Context Signing**: 所有入口验证 ToolGrant JWT（HS256）；Envelope 携带 `tenant_id`, `tool_scope`, `request_id`, `idempotency_key`；可使用 STS 交换服务凭证。  
- **Tenant/RBAC**: Service 层设置 `SET LOCAL app.tenant_id`；GrantMatrix + RBAC 资源 `integration.manage` / `integration.read` 控制访问；审批操作记录责任人。  
- **Outbound Access**: 调用宿主 Secrets Manager、EventBus 时使用短期 STS 凭证并限制作用域。  
- **Observability**: `/healthz` 增加 integration readiness 检查；日志包含 `tenant_id`, `trace_id`, `adapter`, `scope`；metrics 使用 `powerx_integration_*` 命名。  
- **Packaging**: 更新 `plugin.yaml` 的 `data_usage` 与 `security_baseline_version`；发布包包含 OpenAPI、Webhook schema、Nuxt 构建产物以及新的配置样例。

---

- [x] **Host Contract First** {PX-HOST-001}  
  `/api/v1/integration/**` 遵循 host contract；管理端路由纳入 `/_p/<plugin-id>/admin/**`；manifest 更新 data_usage/SLA 说明。
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  ToolGrant + STS 校验；GrantMatrix/RLS 保证 tenant 隔离；Secrets/队列按 tenant_id 分区；测试覆盖多租户场景。
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  新建 `internal/services/integration` 聚合逻辑；HTTP/gRPC/MCP handler 仅解析/转发；repository 继续使用 BaseRepository。
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  明确资源/动作；GrantMatrix 与 RBAC 双重控制；UI 仅暴露可见性，写操作由服务端鉴权。
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  指标/日志/Audit 统一放在 `internal/observability/integration`；`/healthz` 子检查；计划 unit + integration + contract + webhook smoke 测试。
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  评估 Redis 依赖（可降级至 PostgreSQL）；SemVer 次版本；发布需含 OpenAPI/Webhook schema、Nuxt 构建、release notes。

> Any unchecked item must be resolved or explicitly justified in **Complexity Tracking** below.

---

## Project Structure

### Documentation (this feature)

```

specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)

```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., internal/services/foo, internal/transport/http/bar).
  The delivered plan must not include Option labels.
-->

```

backend/
├── internal/
│   ├── domain/models/
│   ├── domain/repository/
│   ├── services/
│   └── transport/{http,grpc}/
└── tests/

frontend/
├── app/ (pages/layouts/components/stores)
└── tests/

```

**Structure Decision**: 采用现有 Go 后端 + Nuxt 管理前端结构。后端新增 `internal/services/integration`、`internal/domain/repository/integration`、`internal/transport/http/admin/integration`、`internal/transport/grpc/integration`、`internal/observability/integration`。数据库迁移位于 `backend/migrations`。前端新增 `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/**`、`web-admin/app/types/integration.ts` 与配套组件/测试。

---

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| [e.g., extra transport] | [current need] | [why single transport insufficient] |
| [e.g., custom repo pattern] | [specific problem] | [why direct DB access insufficient] |
