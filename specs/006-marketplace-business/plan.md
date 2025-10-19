# Implementation Plan: Marketplace & Business (Listing, Pricing, Licensing, Analytics)

**Branch**: `006-marketplace-business` | **Date**: 2025-10-19 | **Spec**: specs/006-marketplace-business/spec.md  
**Input**: Feature specification from `/specs/006-marketplace-business/spec.md`

---

## Summary

构建插件 Marketplace 的商业闭环：提供 Vendor 上架、定价/License 发放、Usage 分析及分润报表。后端基于现有 Gin + GORM 服务扩展 Marketplace 模块、对接税务 SaaS（Stripe Tax/Avalara）与 License Server，前端 Nuxt 4 Admin 增补 Marketplace Console 与 Vendor Dashboard。核心决策见 `research.md`（税务 SaaS、72h 离线许可兜底、GraphQL Checklist）。

---

## Technical Context

**Language/Version**: Go 1.24 (backend), Node 20 + Nuxt 4 (admin UI)  
**Primary Dependencies**: Gin, GORM, Redis (幂等/缓存), Stripe Tax SDK (HTTP 客户端封装), PowerX EventBus & ToolGrant  
**Storage**: PostgreSQL schema `powerx_plugin_base`（新表详见 `data-model.md`），Redis 用于 License 缓存与幂等记录  
**Testing**: `go test ./...`（含 service/repository 单元测试 + integration tenant tx）、`make test`、前端 `pnpm test` + `@nuxt/test-utils`、合同测试（OpenAPI + GraphQL）  
**Target Platform**: Linux 宿主（PowerX 插件运行时），前端打包至宿主反向代理  
**Project Type**: 插件后端 + Admin 前端（Web，全栈）  
**Performance Goals**: License 验证 API p95 < 200ms，Usage Ingest ≥10K req/s，Dashboard 首屏 <5s  
**Constraints**: 多租户 RLS、72 小时离线 License 限制、税务自动化、GDPR 删除 24h 内完成  
**Scale/Scope**: 100+ 发布 Listing、每月 1k 交易、Usage 数据 180 天留存（~数十 GB）

### Platform / Hosting Integration

- **Reverse Proxy & Routes**: `/api/v1/marketplace/**` （宿主代理 `/_p/com.powerx.plugin.base/api/v1`），GraphQL Checklist 走 `/api/v1/admin/marketplace/graphql`。  
- **Context Signing**: 统一使用 ToolGrant JWT，Handler 注入 `tenant_id/user_id/permissions`，SDK 侧 HMAC 校验 Usage。  
- **Tenant/RBAC**: `middleware.RBAC` 映射新资源（`marketplace.listings.*`, `marketplace.usage.view`, `marketplace.revenue.export`）；Postgres `SET LOCAL app.tenant_id`.  
- **Outbound Access**: Stripe Tax/Avalara 调用通过 STS 获取短期凭证或配置化 HTTPS。  
- **Observability**: `/healthz`、结构化日志（tenant_id/request_id/license_id）、Metrics（license_verify_latency、usage_ingest_lag）。  
- **Packaging**: 更新 `plugin.yaml`（routes、RBAC、schema）、`make build && make frontend-build`，随发布打包 `.output/`。

---

## Constitution Check

- [x] **Host Contract First** {PX-HOST-001}  
  所有 API 聚合在 `/api/v1/marketplace`，GraphQL 置于 admin 前缀；`plugin.yaml` 更新 manifest 与 RBAC。
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  ToolGrant JWT → middleware 验签；Repository 内使用 `BeginTenantTx`；Redis 缓存以 `tenant_id` 命名空间隔离。
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  新增 `internal/services/marketplace` 统一复用，HTTP/gRPC Handler 仅做请求转换与授权；Repository 层抽象。
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  定义资源动作矩阵（Vendor 投稿、Reviewer 审核、Tenant 购买、Platform 对账）；UI 只控制可见性。
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  增加指标、审计事件、EventBus；Service + Repository 单测、迁移冒烟、合同/前端测试覆盖核心流程。
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  依赖新增仅限税务 SDK；SemVer 次版本提升，release 前执行 `make dist`，文档/manifest 同步。

---

## Project Structure

### Documentation (this feature)

```
specs/006-marketplace-business/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
└── contracts/
    ├── marketplace-openapi.yaml
    └── ready-checklist.graphql
```

### Source Code (repository root)

```
backend/
├── internal/
│   ├── domain/models/marketplace/
│   ├── domain/repository/marketplace/
│   ├── services/{marketplace,billing,analytics}/
│   ├── transport/http/admin/marketplace/
│   ├── transport/http/tenant/marketplace/
│   ├── transport/grpc/marketplace/            # License server hooks
│   └── observability/marketplace/
├── internal/router/                           # 注册新路由组
├── cmd/database/migrate/                      # 注册迁移
└── pkg/ (若需税务客户端封装)

web-admin/
├── app/pages/_p/com.powerx.plugin.base/admin/marketplace/
├── app/components/marketplace/
├── app/types/marketplace.ts
├── app/stores/marketplace/
└── tests/marketplace/
```

**Structure Decision**: 保持现有后端分层（transport → services → repository → models），新增 `marketplace` 子域；前端沿 Nuxt 页面/组件/类型目录扩展 Marketplace 管理界面与仪表盘。

---

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| None | - | - |
