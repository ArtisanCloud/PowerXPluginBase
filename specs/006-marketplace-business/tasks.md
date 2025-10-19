# Tasks: Marketplace & Business (Listing, Pricing, Licensing, Analytics)

**Input**: Design docs from `/specs/006-marketplace-business/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

> Format: `[ID] [P?] [Story] Description`

---

## Phase 1 — Setup (Shared Infrastructure)

- [x] T001 运行 `make dev-setup`（仓库根目录）确保 Go 依赖与 `golangci-lint` 就绪。  
- [ ] T002 [P] 在 `web-admin/` 执行 `pnpm install` 安装 Nuxt 依赖与测试工具。  
- [ ] T003 [P] 启动本地依赖：`docker compose -f config/docker-compose.integration.yml up -d` 供 Redis/Webhook mock 使用。

---

## Phase 2 — Foundational (Blocking Prerequisites)

- [x] T010 更新 `plugin.yaml` 以及 `backend/internal/router/router.go`、`backend/internal/transport/http/registry.go`，预留 `marketplace` 路由分组与占位 RBAC。  
- [x] T011 [P] 扩展配置：修改 `backend/etc/config.example.yaml` 与 `backend/internal/config/config.go`，新增 `integration.billing.tax_provider`、Stripe Tax/Avalara 凭据等字段并加载到 `config.Config`.  
- [x] T012 [P] 在 `backend/internal/observability/marketplace/metrics.go` 创建指标埋点骨架（license 验证、usage ingest、tax provider error）。  
- [x] T013 [P] 初始化税务 SaaS 客户端封装：在 `backend/internal/services/marketplace/tax_provider_client.go` 与 `backend/internal/shared/app/deps.go` 注册依赖，暴露重试/回放接口。
- [x] T014 [P] 加固零信任栈：在 `backend/internal/transport/http/middleware`、`backend/internal/router/router.go` 与 `backend/internal/transport/http/registry.go` 校验 ToolGrant JWT/HTTPS 配置、同步 `marketplace.*` RBAC，并补充安全回归测试。

**Checkpoint**：基础能力准备完毕，可进入各用户故事实现。

---

## Phase 3 — User Story 1 — Vendor 提交并上架插件 (Priority: P0) 🎯 MVP

**Goal**: Vendor 完成 `.pxp` 包上传、Listing 配置与审核流，Checklist 可在控制台与 CI 校验。  
**Independent Test**: 在 Sandbox 上传带 assets/pricing 占位的 `.pxp` → 审核通过 → Listing 出现在前台并带 Vendor 认证。

### Tests for US1

- [ ] T101 [P] [US1] 在 `backend/tests/contract/marketplace/listings_test.go` 编写 OpenAPI 合同测试，覆盖 `POST/GET/PATCH /marketplace/listings` 与 `/status`.  
- [ ] T102 [P] [US1] 在 `backend/tests/contract/marketplace/checklist_graphql_test.go` 添加 GraphQL 合同测试，验证 checklist 查询/触发 mutation。  
- [ ] T103 [P] [US1] 在 `backend/internal/services/marketplace/listing_service_test.go` 编写服务层单测，覆盖草稿创建、审核通过、资产更新。

### Implementation for US1

- [ ] T104 [US1] 新增迁移文件 `backend/migrations/2025Q4_marketplace_listings.sql` 定义 `marketplace_listings`、`listing_assets`、`listing_versions`、`checklist_runs`、`checklist_items` 表及 RLS。  
- [ ] T105 [P] [US1] 在 `backend/internal/domain/models/marketplace/` 创建 `listing.go`、`asset.go`、`checklist.go`，映射上述表并声明枚举。  
- [ ] T106 [P] [US1] 实现仓储：`backend/internal/domain/repository/marketplace/listing_repository.go`、`checklist_repository.go`，提供草稿创建、资产管理、Checklist 读写接口。  
- [ ] T107 [US1] 完成服务层 `backend/internal/services/marketplace/listing_service.go`，封装提交流程、审核状态转换、资产存储调用与事件发布。  
- [ ] T108 [US1] 构建 GraphQL Resolver：在 `backend/internal/transport/http/admin/marketplace/checklist_resolver.go` 实现 Query/Mutation，复用服务层并校验 RBAC。  
- [ ] T109 [US1] 开发 Admin HTTP Handler：`backend/internal/transport/http/admin/marketplace/listings_handler.go` 暴露列表、详情、草稿更新、审核提交 API。  
- [ ] T110 [US1] 更新路由与权限：在 `backend/internal/transport/http/admin/routes.go` 注册 handler，在 `backend/internal/transport/http/registry.go` 合并 RBAC，补充 `marketplace.listings.{read,write,review}`。  
- [ ] T111 [US1] 更新 `plugin.yaml` 与 `backend/etc/config.example.yaml`，同步 Checklist GraphQL 端点、RBAC 权限描述与说明文档链接。  
- [ ] T112 [P] [US1] 前端页面：在 `web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/listings.vue` 实现 Listing 控制台，含上传表单与审核状态视图。  
- [ ] T113 [P] [US1] 创建 Checklist 组件与 composable：`web-admin/app/components/marketplace/ChecklistRunner.vue`、`web-admin/app/composables/useMarketplaceChecklist.ts`，对接 GraphQL。  
- [ ] T114 [P] [US1] 编写前端测试：`web-admin/tests/marketplace/listings.spec.ts` 验证表单校验、GraphQL 交互与列表渲染。
- [ ] T115 [P] [US1] 在 `backend/tests/integration/marketplace/listing_edge_cases_test.go` 覆盖 `.pxp` 缺资产、KYC 撤销时的自动阻断与复审流程。  
- [ ] T116 [US1] 构建推荐服务：`backend/internal/services/recommendation/engine.go` 计算排序权重、A/B 实验指标，并落地 `recommended_weight`.  
- [ ] T117 [P] [US1] 实现 Discovery 同步任务：`backend/internal/jobs/marketplace/recommendation_sync.go` 与调度配置，按小时推送候选数据。  
- [ ] T118 [US1] 扩展 Admin Handler/前端，提供推荐配置与实验面板：`backend/internal/transport/http/admin/marketplace/recommendation_handler.go`、`web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/recommendation.vue`。  
- [ ] T119 [P] [US1] 在 `backend/internal/services/marketplace/listing_service.go` 增加品牌素材尺寸/视频时长校验（NFR-005），并为前端提供即时反馈。  
- [ ] T120 [US1] 添加提交流程 SLA 监控：`backend/internal/observability/marketplace/metrics.go` & `backend/tests/perf/marketplace/listing_submission_test.go`，确保 3 分钟内返回校验结果。  
- [ ] T121 [P] [US1] 在 `web-admin/tests/marketplace/recommendation.spec.ts` 验证推荐配置、A/B 切换与曝光指标呈现。

**Checkpoint**：US1 独立可演示（提交→审核→展示）。

---

## Phase 4 — User Story 2 — 租户购买并激活 License (Priority: P1)

**Goal**: 租户选购价格计划完成支付，License Server 签发 JWT，支持离线 72h 与续费提醒。  
**Independent Test**: Subscription 下单 → Billing 引擎扣费 → License Server `issue` → 插件验证缓存 → 临期提醒。

### Tests for US2

- [ ] T201 [P] [US2] 在 `backend/tests/contract/marketplace/licenses_test.go` 添加 REST 合同测试（创建、查询、续订、离线扩展）。  
- [ ] T202 [P] [US2] 在 `backend/tests/integration/marketplace/license_flow_test.go` 编写端到端测试，模拟支付成功→License 发放→续费。  
- [ ] T203 [P] [US2] 在 `backend/internal/services/marketplace/license_service_test.go` 编写单测，覆盖额度校验、离线续期令牌生成、事件记录。

### Implementation for US2

- [ ] T204 [US2] 新增迁移 `backend/migrations/2025Q4_marketplace_licensing.sql`，创建 `pricing_plans`、`plan_tiers`、`licenses`、`license_events`、`tax_transactions` 并配置 RLS。  
- [ ] T205 [P] [US2] 增补模型：`backend/internal/domain/models/marketplace/pricing.go`、`license.go` 定义计划/License 结构与约束。  
- [ ] T206 [P] [US2] 实现仓储：`backend/internal/domain/repository/marketplace/pricing_repository.go`、`license_repository.go`，支持计划查询、License 发放与事件写入。  
- [ ] T207 [US2] 编写 `backend/internal/services/marketplace/license_service.go`，实现购买→税费计算→调用 Billing Engine→请求 License Server→缓存 JWT→记录事件。  
- [ ] T208 [US2] 构建税务 SaaS 适配：在 `backend/internal/services/marketplace/tax_provider_client.go` 扩展多货币结算/重试逻辑并持久化 `tax_transactions`。  
- [ ] T209 [US2] 开发租户侧 HTTP Handler：`backend/internal/transport/http/tenant/marketplace/licenses_handler.go` 支持计划浏览、下单、续期、离线续期令牌。  
- [ ] T210 [US2] 更新 gRPC 入口：在 `backend/internal/transport/grpc/marketplace/license_server.go` 提供 License Server 回调/验证接口并接入现有 `server.go`。  
- [ ] T211 [US2] 添加缓存与通知：在 `backend/internal/services/marketplace/license_service.go` 整合 Redis 缓存（离线 72h）与 `backend/internal/observability/marketplace/events.go` 发布续费提醒。  
- [ ] T212 [P] [US2] 管理前端价格计划编排：`web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/plans.vue` 提供计划编辑与发布。  
- [ ] T213 [P] [US2] 租户购买界面：`web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/purchase.vue` 实现价格展示、支付触发、License 状态提醒。  
- [ ] T214 [P] [US2] 前端测试：`web-admin/tests/marketplace/license_flow.spec.ts` 覆盖计划选择、支付回执、离线续期提示。  
- [ ] T215 [US2] 更新 `plugin.yaml` 与 `backend/etc/config.example.yaml`，声明 `marketplace.license.*` RBAC、税务配置示例与 License Server 事件订阅。
- [ ] T216 [P] [US2] 压测与指标：在 `backend/tests/perf/marketplace/license_latency_test.go` 与 `backend/internal/observability/marketplace/metrics.go` 校验 License 验证 p95 <200ms，并输出告警。  
- [ ] T217 [US2] 处理支付/签发延迟：实现 `backend/internal/services/marketplace/license_recovery.go`，在 Billing 与 License Server 不一致时补偿；新增测试 `backend/tests/integration/marketplace/license_delay_test.go`。  
- [ ] T218 [P] [US2] 构建续费提醒任务：`backend/internal/jobs/marketplace/license_renewal_notifier.go` 定期扫描 `offline_until`/到期记录并触发通知渠道。

**Checkpoint**：US1 + US2 可独立运行（提交流程 + 购买激活）。

---

## Phase 5 — User Story 3 — Vendor 分析使用与收入 (Priority: P1)

**Goal**: Vendor 可查看安装/调用/收入趋势，系统支持 Usage 聚合、异常告警、分润报表。  
**Independent Test**: SDK 批量上报 Usage → Analytics Pipeline 聚合 → Dashboard 展示趋势与告警 → 月末导出分润报表。

### Tests for US3

- [ ] T301 [P] [US3] 在 `backend/tests/contract/marketplace/usage_test.go` 添加 REST 合同测试，覆盖 `/marketplace/usage` ingest 与 metrics 查询、报表列表。  
- [ ] T302 [P] [US3] 在 `backend/tests/integration/marketplace/analytics_flow_test.go` 编写集成测试，模拟 Usage 上报→聚合→告警→报表生成。  
- [ ] T303 [P] [US3] 在 `backend/internal/services/marketplace/analytics_service_test.go` 编写单测，验证配额判定、spike detection、分润计算。

### Implementation for US3

- [ ] T304 [US3] 创建迁移 `backend/migrations/2025Q4_marketplace_analytics.sql` 定义 `usage_envelopes`、`usage_aggregates`、`revenue_share_reports`、`notifications` 表与 RLS 策略。  
- [ ] T305 [P] [US3] 定义模型：`backend/internal/domain/models/marketplace/usage.go`、`revenue.go`、`notification.go`。  
- [ ] T306 [P] [US3] 实现仓储：`backend/internal/domain/repository/marketplace/usage_repository.go`、`revenue_repository.go`、`notification_repository.go`。  
- [ ] T307 [US3] 编写 `backend/internal/services/marketplace/usage_ingest_service.go`，处理批量上报、重放、幂等校验，落地 `usage_envelopes`。  
- [ ] T308 [US3] 编写 `backend/internal/services/marketplace/analytics_service.go`，聚合 Usage→生成趋势、quota 检查、分润报表，并触发告警通知。  
- [ ] T309 [US3] 增强观测：在 `backend/internal/observability/marketplace/metrics.go` 补充 usage lag、revenue generation 指标，新增 `backend/internal/observability/marketplace/events.go` 发布 `usage.spike.detected`。  
- [ ] T310 [US3] 实现 HTTP Handler：`backend/internal/transport/http/admin/marketplace/analytics_handler.go` 暴露 Usage metrics、报表导出 API；注册至路由与 RBAC。  
- [ ] T311 [US3] 同步前端仪表盘：`web-admin/app/pages/_p/com.powerx.plugins.base/admin/integration/marketplace/dashboard.vue` 显示趋势图、异常告警、分润报表导出。  
- [ ] T312 [P] [US3] 构建前端 stores/composables：`web-admin/app/stores/marketplaceAnalytics.ts`、`web-admin/app/composables/useUsageMetrics.ts`。  
- [ ] T313 [P] [US3] 添加前端测试：`web-admin/tests/marketplace/analytics_dashboard.spec.ts` 验证趋势渲染、告警提示、报表下载。  
- [ ] T314 [US3] 更新 `quickstart.md` 与文档，加入 Usage 上报、Dashboard 操作指南与告警说明。
- [ ] T315 [P] [US3] 在 `backend/tests/perf/marketplace/usage_load_test.go` 模拟 10K req/s 吞吐并验证 `usage_ingest_lag` 指标（NFR-003）。  
- [ ] T316 [US3] 实现 GDPR 删除链路：`backend/internal/services/marketplace/privacy_service.go`、`backend/tests/integration/marketplace/gdpr_delete_test.go`，确保 24h 内删除 Usage 数据（NFR-004）。  
- [ ] T317 [P] [US3] Dashboard 性能监控：在 `web-admin/tests/perf/marketplace/dashboard_performance.spec.ts` 与 `web-admin/app/plugins/metrics.client.ts` 记录首屏加载时间（NFR-006）。  
- [ ] T318 [US3] 增补异常告警测试：`backend/tests/integration/marketplace/usage_spike_test.go` 确认 `usage.spike.detected` 与配额兜底逻辑。

**Checkpoint**：三大用户故事均独立通过测试 → 功能完整。

---

## Phase 6 — Polish & Cross-Cutting

- [ ] T901 [P] 整理 `docs/` 与 `README.md`，新增 Marketplace 商业闭环章节。  
- [ ] T902 清理临时代码/提高覆盖率：运行 `make fmt && make lint && make test && pnpm lint`.  
- [ ] T903 [P] 校验打包流程：执行 `make build && make frontend-build && make dist`，确认产物包含新增 API/UI。  
- [ ] T904 验证配置样例：审查 `backend/etc/config.example.yaml`、`config/docker-compose.integration.yml` 是否包含 Marketplace 说明。  
- [ ] T905 更新 `plugin.yaml` 版本号与 `docs/references/changelog.md`（若存在）记录新功能。

---

## Dependencies & Execution Order

1. **Phase 1 → Phase 2 → Phase 3/4/5 → Phase 6**  
2. User story顺序：US1 (P0) → US2 (P1) → US3 (P1)。US3 依赖 License/计划数据，因此需在 US2 完成后进行。  
3. 通过 `[P]` 标记的任务可在确保无文件冲突时并行执行。

### Story-Level Dependency Graph

- US1 → US2 → US3  
  - US1 交付 Listing + Checklist  
  - US2 在 US1 基础上提供定价与 License 闭环  
  - US3 消费 License/Usage 数据生成分析

### Parallel Execution Examples

- Phase 1：T002 与 T003 可与 T001 并行。  
- US1：T105/T106/T112/T113/T114 涵盖不同文件，可多人并行。  
- US2：T205/T206/T212/T213/T214 可与 T207/T208/T209 并行推进。  
- US3：T305/T306/T312/T313 可与 T307–T310 并行。

---

## Implementation Strategy

1. **MVP（US1）**：优先实现 Listing + Checklist 流程，确保 Marketplace 基础供给侧可演示。  
2. **商业闭环（US2）**：在 MVP 稳定后补完定价、License、税费与续费提醒，形成完整交易链路。  
3. **分析洞察（US3）**：最后实现 Usage 聚合与分润报表，为 Vendor/Platform 提供运营视角。  
4. 每个阶段结束运行 `make test`、`pnpm test` 与合同测试，保持可部署增量。

---

## Task Tally

- **Total tasks**: 70  
  - Setup: 3  
  - Foundational: 5  
  - US1: 21  
  - US2: 18  
  - US3: 18  
  - Polish: 5
- **Parallelizable tasks** (`[P]`): 39  
  - US1: 12  
  - US2: 10  
  - US3: 9  
  - Setup/Foundational/Polish: 8
- **Independent Test Criteria**:  
  - US1：Listing 提交→审核→展示与推荐同步（T101–T121）。  
  - US2：购买→税费→License 签发→续费提醒（T201–T218）。  
  - US3：Usage 上报→聚合→告警→分润/报表（T301–T318）。
- **Suggested MVP Scope**: 完成 Phase 1–3（含 US1），交付可上架插件的供应链闭环。
