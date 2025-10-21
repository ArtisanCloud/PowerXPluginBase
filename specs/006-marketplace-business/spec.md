# Feature Specification: Marketplace & Business (Listing, Pricing, Licensing, Analytics)

**Feature Branch**: `006-marketplace-business`  
**Created**: 2025-10-19  
**Status**: Draft  
**Input**: Title: Marketplace & Business (Listing, Pricing, Licensing, Analytics); WHAT/WHY: 建立 PowerX 插件 Marketplace 的商业分发与经济闭环，包括插件上架展示、价格计划、License 授权机制与使用分析报表，让 Vendor、Tenant、Platform 三方拥有可交易、可追踪、可计费的商业化能力。Scope: Listing & Branding 资产、定价与授权模型、Usage 分析与报表体系；Out-of-Scope: 01 Lifecycle、02 Capabilities、03 Runtime/Ops、04 Security & Compliance；Dependencies: Marketplace 具备 License Server、Billing Engine、Analytics Pipeline，插件产物已声明 pricing/license/usage，通信统一走 HTTPS + JWT/ToolGrant，Vendor 已完成 KYC 校验，RBAC `marketplace.usage.view` 控制租户报表访问。

## Clarifications

### Session 2025-10-19

- Q: Marketplace 推荐算法需要与现有 Discovery 系统如何集成，以支持插件 Listing 的曝光与排序？ → A: 建立独立服务按小时同步到 Discovery
- Q: Usage 分析管线是否需要支持离线重放与 GDPR“删除我”请求？ → A: 支持自动化离线重放与 GDPR 删除链路
- Q: Billing Engine 在多货币与跨区域税费方面需要支持到什么程度？ → A: 完整多货币结算并自动计算地域税费
- Q: Billing Engine 在跨区域税费计算上如何实现？ → A: 通过集成税务 SaaS 实时获取税率并自动开票
- Q: On-Prem 环境的 License 离线容错策略应如何设计？ → A: 允许 72 小时离线并提供一次性临时续期令牌
- Q: Vendor 品牌素材是否允许动画/视频资产？ → A: 允许静态图片和短视频（≤15 秒, ≤50 MB, MP4/WebM）
- Q: Listing Ready Checklist 是否需要开放 API 供自动化校验？ → A: 提供 GraphQL API，CI/CD 与控制台共用

## User Scenarios & Testing *(mandatory)*

### User Story 1 – Vendor 提交并上架插件 (Priority: P0)

作为 Vendor，我希望通过标准化的提交流程把 `.pxp` 包提交给 Marketplace，完成品牌素材和能力信息配置，在通过安全与合规审核后快速上架，便于潜在租户发现并试用。

**Why this priority**: 没有可用的提交流程和 Listing 渠道，就无法构建 Marketplace 的供给侧生态。  
**Independent Test**: 在 Sandbox 中提交一个包含必备 assets、docs 与 pricing 占位的 `.pxp` 包，经审核后 Listing 卡片可在 Marketplace 前台展示并带有 Vendor 认证徽章。

**Acceptance Scenarios**

1. **Given** Vendor 在控制台上传 `.pxp` 包并填写 Listing 表单，**When** 自动校验结构与资产通过，**Then** 系统生成草稿 Listing 并提示审核状态。
2. **Given** 审核人员对 Listing 进行安全、合规、兼容性检查，**When** 所有必填字段与素材符合品牌规范，**Then** Listing 状态切换为 “上架” 并可被搜索与推荐。
3. **Given** Vendor 更新封面或版本描述，**When** 变更在后台保存，**Then** Marketplace 触发增量审核并保留历史版本。

---

### User Story 2 – 租户购买并激活 License (Priority: P1)

作为租户管理员，我需要在 Marketplace 中查看插件的价格计划、试用政策与授权条款，购买后收到有效的 License 并在租户环境中自动激活，同时能在额度耗尽或到期前获得提示。

**Why this priority**: 定价与授权机制是商业闭环的核心，直接影响营收与合规。  
**Independent Test**: 选择 Subscription 计划下单 → Billing 引擎扣费 → License Server 签发 JWT → 插件启动时验证并缓存 License → License 即将到期时触发续费提醒。

**Acceptance Scenarios**

1. **Given** 租户选择某个付费计划并完成支付，**When** License Server 接收到 `issue` 请求，**Then** 返回包含 `license_id / plugin_id / tenant_id / plan_id / signature / expiry` 的 JWT，并写入授权审计记录。
2. **Given** License 临近到期，**When** 平台发送续费提醒并允许租户升级计划，**Then** 成功续费后 License Server 生成新的 `renewed` 事件并更新缓存。
3. **Given** Usage 上报显示租户超出额度，**When** License Server 校验发现违反计划限制，**Then** 返回错误码并触发停用/降级流程，同时记录异常事件。

---

### User Story 3 – Vendor 分析使用与收入 (Priority: P1)

作为 Vendor，我需要快速了解插件的安装量、调用次数、License 构成与地域分布，掌握订阅续费趋势，并在异常警报（如 usage spike 或 quota exceeded）出现时收到通知，以便优化产品和商业策略。

**Why this priority**: 没有可观测的数据洞察，Vendor 无法迭代定价或证明价值，也无法与平台协同运营。  
**Independent Test**: 插件通过 SDK 批量上报 Usage Envelope → Analytics Pipeline 聚合 → Vendor Dashboard 显示租户/版本维度的趋势图与异常提示。

**Acceptance Scenarios**

1. **Given** 插件按小时批量上报 Usage Envelope，**When** Analytics Pipeline 将数据入库，**Then** Vendor Dashboard 展示按租户、计划、版本分组的调用量与收入趋势。
2. **Given** 某租户调用量异常飙升，**When** `spike.detected` 事件被触发，**Then** Vendor 与 Platform 共同收到告警并可下钻到具体 License 与 API 指标。
3. **Given** Vendor 需要对账结算，**When** 月结周期结束，**Then** 系统生成 Vendor 80% / Platform 15% / Fee 5% 的分润报表并提供导出。

---

### User Story 4 – 租户监控使用与费用 (Priority: P2)

作为租户运营人员，我希望在控制台查看当前订阅的额度消耗、费用预估、续费入口及异常事件（如 quota 超限、License 锁定），以便及时调整计划或与 Vendor 联系。

**Why this priority**: 提升租户对商业条款的透明度，可以降低支持成本并促进续费。  
**Independent Test**: 租户访问 “My Licenses” → 查看 Usage Dashboard → 接收 quota 预警 → 通过入口续费成功。

**Acceptance Scenarios**

1. **Given** 租户用户具备 `marketplace.usage.view` 权限，**When** 访问 Usage 页面，**Then** 仪表盘展示额度使用率、近期预估费用以及 License 状态。
2. **Given** 额度达到 80%，**When** 系统触发阈值事件，**Then** 租户收到站内信/邮件提醒并可一键升级或购买附加包。
3. **Given** License 被平台因违规暂挂，**When** 租户请求重新激活，**Then** 平台记录审计并在合规通过后恢复授权。

---

### Edge Cases

- `.pxp` 包缺失 `plugin.yaml` 中的 pricing/licensing 配置或包含不兼容的 assets，如何在提交阶段立即阻断并给出修复指引？  
- Vendor 资料 KYC 状态变更（例如被撤销认证）时，现有 Listing 与 License 是否需要自动下架或冻结？  
- Usage 上报断流或网络异常导致的重复/延迟报送，如何实现幂等与补偿并保持账务准确？  
- License Server 与 Billing Engine 之间出现延迟（已付款但未签发 License），如何确保租户体验与财务对账一致？  
- 推荐算法可能偏好高装机量插件，如何保证新插件或细分垂直能获得曝光（探索 vs. 利益）？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Marketplace 必须提供基于 `.pxp` 包的 Vendor 提交流程，对包结构（`plugin.yaml`, `backend/`, `web-admin/`, `assets/`, `docs/`）执行自动校验，并生成草稿 Listing。  
- **FR-002**: 系统必须记录 Listing Metadata（名称、简介、版本、分类、能力标签、封面、截图、文档链接、品牌主题色），并支持多语言本地化。  
- **FR-003**: 审核工作流必须覆盖安全、合规、兼容性、文档完整性检查，并允许 Reviewer 留存备注及版本回溯。  
- **FR-004**: Marketplace 推荐引擎需以独立服务按小时刷新候选列表，将结果同步至现有 Discovery 排序管线，排序权重基于安装量、评分、更新频率、响应时间、品牌完整度，并支持配置化 A/B 测试。  
- **FR-005**: `plugin.yaml` 必须扩展 `pricing.plans`、`pricing.trial`、`pricing.usage` 节点，描述 Free / One-Time / Subscription / Usage-Based 计划、计费周期、试用时长与用途限制。  
- **FR-006**: License Server 必须支持 `issue`, `verify`, `renew`, `usage`, `revoke` 接口，返回 JWT 结构化 License（`license_id`, `plugin_id`, `tenant_id`, `plan_id`, `features`, `quota`, `signature`, `expires_at`）。  
- **FR-007**: 插件启动时必须验证 License 并缓存到期时间，缓存更新后需要在后台刷新，并在 On-Prem 环境允许最长 72 小时离线，同时支持一次性临时续期令牌兜底。  
- **FR-008**: License 生命周期事件（`issued`, `renewed`, `revoked`, `usage.reported`, `trial.expired`）必须写入审计日志并发往 Marketplace EventBus。  
- **FR-009**: Billing Engine 必须支持分润与结算（Vendor 80% / Platform 15% / Fee 5%），通过集成税务 SaaS（如 Stripe Tax / Avalara）实时获取税率并自动开票，提供月度对账报表和异常记录（退款、chargeback），并实现完整的多货币结算，依据租户地域自动计算增值税/销售税。  
- **FR-010**: Usage 上报需定义 Envelope（`plugin_id`, `tenant_id`, `license_id`, `metrics`, `timestamp_range`, `signature`），插件 SDK 需支持批量上报、重试及 HMAC 签名。  
- **FR-011**: Analytics Pipeline 必须支持按租户、版本、计划、地域、时间窗口、License 状态等维度的聚合，提供 Vendor、Tenant、Platform 不同视图，并具备自动化离线重放能力补偿延迟或失败上报。  
- **FR-012**: Vendor Dashboard 必须提供安装量趋势、收入曲线、地域地图、异常告警列表；Tenant Dashboard 必须展示使用率、费用预估、续费入口与异常状态。  
- **FR-013**: 系统必须提供异常策略处理（试用结束、额度超限、License 过期、违规停用、恢复激活），并确保前台展示与后端状态一致。  
- **FR-014**: 所有 Marketplace 接口和报表必须通过 HTTPS，使用 ToolGrant / JWT 鉴权并尊重 RBAC（包含 `marketplace.usage.view`）。  
- **FR-015**: Listing Ready Checklist 必须对 Vendor 可见，提供提交前自检指引并与审核结果关联，同时暴露 GraphQL API 以供 CI/CD 与控制台共用自动化校验。

### Non-Functional & Quality Attributes

- **NFR-001**: Listing 提交流程应在 3 分钟内完成自动校验反馈；审核 SLA 默认 2 个工作日，支持加速。  
- **NFR-002**: License 验证 API p95 延迟 < 200ms，失败率 < 0.1%；License 缓存过期前 5 分钟必须主动刷新。  
- **NFR-003**: Usage 上报系统支持至少 10K req/s，保障 180 天数据保留并可按租户隔离，同时提供自动化离线重放以恢复数据缺口。  
- **NFR-004**: Usage Analytics 必须提供 GDPR “删除我” 处理流程并在 24 小时内响应 Data Subject Request，且分润对账要求财务级准确性，任何错账需在 24 小时内可追溯并回滚。  
- **NFR-005**: 品牌素材需符合尺寸/格式限制（Logo SVG/PNG, Cover 16:9, Font & 颜色遵守 Marketplace 指南），允许上传短视频（≤15 秒, ≤50 MB, MP4/WebM）用于展示动画效果。  
- **NFR-006**: 所有仪表盘需在 5 秒内加载最近 30 天数据，可分页或按时间窗口下钻。

### Key Entities

- **MarketplaceListing**: 汇总插件 Listing 信息、状态、审核记录、发布历史与推荐权重。  
- **VendorProfile**: Vendor 的身份认证、KYC 状态、品牌素材、支持联系与 SLA 指标。  
- **PricingPlan**: 描述 Free/One-Time/Subscription/Usage-Based 计划的定价、额度、试用策略、取消政策。  
- **LicenseToken**: License Server 签发的 JWT，包含授权上下文、额度、特性与签名信息。  
- **LicenseEvent**: 记录 License 生命周期中的重要事件（issue/renew/revoke/usage/trial），用于审计与通知。  
- **UsageEnvelope**: 插件上报的使用信息，支持批量、签名、时间窗口与指标集合。  
- **RevenueShareReport**: 归档分润信息，按月度生成 Vendor/Platform/Fee 的对账视图。  
- **AnalyticsDashboardConfig**: 定义 Vendor/Tenant/Platform 三方的指标视图、维度过滤与告警阈值。

## Success Criteria *(mandatory)*

- **SC-001**: 90% 的新提交插件在首次审核中通过 Listing Ready Checklist，审核平均用时 ≤ 2 个工作日。  
- **SC-002**: 100% 正式上架插件的 License 验证成功率 ≥ 99.5%，试用结束或额度超限的停用响应时间 ≤ 5 分钟。  
- **SC-003**: Usage 上报延迟（插件发送到 Marketplace Analytics 可查询）在 95% 情况下 ≤ 2 分钟，数据保留 180 天无丢失。  
- **SC-004**: Vendor Dashboard 与结算报表在每月 3 个工作日内出具，分润误差率 < 0.5%。  
- **SC-005**: Tenant 续费预警触达率 ≥ 95%，Quota 超限后平均修复时间（MTTR） ≤ 1 个工作日。  
- **SC-006**: 推荐算法 A/B 实验显示，新插件曝光量在不降低转化率的前提下提升 ≥ 25%，且无明显品牌违规。

## Constraints, Assumptions & Dependencies

- Marketplace 依赖现有的 License Server、Billing Engine、Analytics Pipeline，可扩展相关 API。  
- 插件 `.pxp` 包需事先包含 pricing/licensing/usage 声明；SDK 已支持 License 验证与 Usage 上报。  
- 所有通信使用 HTTPS + ToolGrant/JWT，Audit & Telemetry 已启用。  
- Vendor 完成 KYC 并绑定收款账户；Marketplace 审核团队维护版本与认证。  
- 租户访问 Usage 报表必须具备 `marketplace.usage.view` 权限。  
- 安全、合规细则由 04 Security & Compliance 项目负责，此处仅引用审核结果。

## Open Questions & Follow-ups

暂无新的开放问题，后续规划阶段根据实现细节再补充。
