# Feature Specification: Support & Operations (Support Playbook, Incident Handling, SLA/SLO)

**Feature Branch**: `007-support-and-operations`  
**Created**: 2025-10-21  
**Status**: Draft  
**Input**: User description: "Title: Support & Operations (Support Playbook, Incident Handling, SLA/SLO)\nWHAT/WHY: 建立 PowerX 插件在 Marketplace 上线后的标准化支持与运维体系，\n确保插件运行质量、事故可控、SLA 可量化、支持体验可追踪。\n目标是形成从客户支持 → 事故响应 → SLA 量化 → 改进复盘的闭环机制，为租户与 Vendor 提供可预期的服务可靠性。\nScope:\n  - **客户支持体系（Customer Support Playbook）**\n    - 统一的支持渠道：Marketplace 工单系统、Vendor 邮箱、In-App Feedback、紧急通道；\n    - 工单生命周期（Created → Assigned → In-Progress → Resolved → Closed）；\n    - 工单字段与 API/Webhook 事件（ticket.created / updated / resolved / closed）；\n    - 工单分级标准（P0–P4：响应时间与解决目标）；\n    - 知识库与自助文档要求（README, FAQ, Troubleshooting, Support_Policy）；\n    - 绩效指标：FRT ≤ 4h、MTTR ≤ 24h、CSAT ≥ 4.5、Resolution ≥ 95%；\n    - SLA 不达标处罚机制与季度审核；\n    - 安全控制与脱敏规则（JWT 访问、附件隔离、审计日志）；\n    - RCA 模板与支持团队角色（Agent / Engineer / Manager / Liaison / QA）；\n    - 支持准备清单（Support Ready Checklist）:contentReference[oaicite:0]{index=0}。\n  - **事故处理与应急响应（Incident Handling）**\n    - SEV-0~SEV-4 分级矩阵与响应时间（15 min – 24 h）；\n    - 事件生命周期：检测 → 确认 → 通报 → 缓解 → RCA → 修复 → 验证 → 复盘；\n    - 通报模板（incident_id / plugin_id / severity / mitigation / next_update）；\n    - 安全事件流程（检测 → 隔离 → 通报 → 修复 → 通告 → 复盘）；\n    - 常见事件类型及应对策略（功能中断 / 性能下降 / 外部依赖失效 / 配置错误 / 安全漏洞）；\n    - 监控指标集成（response_time_ms, error_rate, cpu_usage, dependency.failure_rate）；\n    - 通信机制（Support Hub / Hotline / security@powerx.io / 状态页）；\n    - 事件标签与季度复盘统计（#availability #security #performance #dependency）；\n    - 自检清单（Incident Ready Checklist）:contentReference[oaicite:1]{index=1}。\n  - **服务等级与目标（SLA / SLO / SLI）**\n    - 概念定义：SLA = 承诺，SLO = 目标，SLI = 指标；\n    - 插件分类与推荐 SLA（Real-time 99.9% / Transactional 99.5% / Utility 99.0%）；\n    - 核心指标：Uptime ≥ 99.5%，Avg Response < 800 ms，成功率 ≥ 99%，FRT ≤ 4 h；\n    - 公式：\n      - Uptime % = (Total − Downtime)/Total × 100；\n      - MTTR = 修复耗时总和 / 事件数；\n      - SLA Score = 0.4 × Uptime + 0.3 × Support + 0.3 × Reliability；\n    - 监控来源：Health Gateway / APM / Observability Hub / Support Hub；\n    - 处罚与激励机制（违规 → 降级，下架；连续 99.9% → 推荐位加权）；\n    - SLA 公示与 API：\n      `GET /api/v1/marketplace/sla/{plugin_id}`；\n    - manifest 内 SLO 定义 （availability / response_time / error_rate / support_frt）；\n    - SLA 与 License/Pricing 的联动（高 SLA 计划 → 更高佣金）；\n    - 自检清单（SLA Ready Checklist）:contentReference[oaicite:2]{index=2}。\nOut-of-Scope:\n  - 插件生命周期与构建（归属 01_plugin_lifecycle）；\n  - 能力与 Schema 契约（归属 02_capabilities_and_schema）；\n  - 安全与 ToolGrant 机制（归属 04_security_and_compliance）；\n  - Marketplace 商业策略与 License 分润（归属 06_marketplace_and_business）。\nDependencies/Assumptions:\n  - PowerX Support Hub 与 Incident Center 已部署；\n  - Vendor 已注册支持邮箱与 Webhook 回调；\n  - Observability Stack (Logs + Metrics + APM) 接入；\n  - SLA 数据每日采样、月度汇总、季度复盘；\n  - 所有支持/事件操作均记录 Audit Log；\n  - 插件 manifest 包含 support 与 incident 配置；\n  - Marketplace Dashboard 实时展示 SLA 等级与 支持绩效。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Stand Up Support Playbook & Channels (Priority: P1)

PowerX 支持经理需要为新上线插件配置统一支持渠道、定义工单字段与分级标准，并发布 Support Ready Checklist，确保租户和 Vendor 可以通过明确定义的入口获得响应。

**Why this priority**: 没有完整的支持规范就无法为 Marketplace 提供基础服务承诺，会直接影响上线审核与租户满意度。

**Independent Test**: 仅通过 Support Hub 校验新的支持配置、完成 Support Ready Checklist，以及触发 `ticket.created` Webhook 便可验证核心价值。

**Acceptance Scenarios**:

1. **Given** 插件尚未启用任何支持渠道，**When** 支持经理在 Support Hub 完成渠道配置并发布 Support Playbook，**Then** 所有渠道在 10 分钟内对租户与 Vendor 可见，并触发 Support Ready Checklist 勾选。
2. **Given** Vendor 通过 In-App Feedback 提交 P2 工单，**When** 工单进入队列并自动套用模板字段，**Then** 系统在 4 小时内记录首响时间并发出 `ticket.created` 通知。

---

### User Story 2 - Coordinate Incident Response Lifecycle (Priority: P2)

当出现重大故障时，Incident Commander 需要按 SEV 矩阵触发应急流程，完成通报、缓解与 RCA，并通过 Incident Ready Checklist 追踪进度。

**Why this priority**: 事故响应直接影响 SLA 合规与信任度，必须保证流程可执行且信息闭环。

**Independent Test**: 通过模拟 SEV-1 事件，从检测、通报到 RCA 归档的全流程演练即可独立验证该能力。

**Acceptance Scenarios**:

1. **Given** 监控检测到响应时间异常并触发 SEV-1 告警，**When** Incident Commander 接手并在 15 分钟内发布首次通报，**Then** 事件时间线包含检测、确认、缓解动作且 Incident Ready Checklist 标记为进行中。
2. **Given** 事故完成修复，**When** 团队在 48 小时内提交 RCA 模板并归档标签，**Then** 系统自动计算 MTTR 并在季度复盘报告中纳入统计。

---

### User Story 3 - Publish SLA & Reliability Insights (Priority: P3)

产品运营经理需要面向租户公开 SLA 指标、计划差异和处罚/激励结果，确保 Marketplace Dashboard 与 API 展示一致，并驱动持续改进。

**Why this priority**: 透明的 SLA 与绩效报告是 Marketplace 引导 Vendor 提升服务质量、租户评估可靠性的核心。

**Independent Test**: 通过月度 SLA 汇总任务生成报告、调用 `GET /api/v1/marketplace/sla/{plugin_id}` 并核对 Dashboard 展示即可独立验证价值。

**Acceptance Scenarios**:

1. **Given** 月度 SLA 汇总任务执行完成，**When** 运营经理下载报告并核对 Dashboard 数据，**Then** 所有核心 SLI（Uptime、Response、Success、FRT）与 API 返回值保持一致且标记激励或处罚结果。
2. **Given** 某插件连续两个月达成 Real-time 等级目标，**When** 系统自动更新推荐权重并记录在 SLA Ready Checklist，**Then** Marketplace Dashboard 显示“推荐”标识且通知 Vendor。

---

### Edge Cases

- 多渠道提交的重复工单需要合并，防止统计重复计入 FRT/MTTR。
- Incident 在缓解后再次复发时需维持同一 incident_id，并重置沟通频率。
- Vendor 未完成支持或 SLA 准备清单时禁止发布 Marketplace Listing。
- SLA 指标缺失或采样失败时，应使用最近的有效数据并向运营发出数据质量告警。
- 安全事件包含敏感附件时，下载权限必须由 Incident Manager 审批并自动脱敏日志。

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须支持为插件配置 Marketplace 工单、Vendor 邮箱、In-App Feedback、紧急通道等渠道，并允许为每个渠道定义可用时间与升级路径。
- **FR-002**: 系统必须提供标准工单生命周期（Created、Assigned、In-Progress、Resolved、Closed）及状态转换规则，阻止跳过必要状态。
- **FR-003**: 系统必须针对工单事件推送 `ticket.created`、`ticket.updated`、`ticket.resolved`、`ticket.closed` Webhook，并提供带签名的 API 查询接口以供 Vendor 对账。
- **FR-004**: 系统必须内置 P0–P4 工单分级标准，自动计算并展示 FRT、MTTR、CSAT、Resolution Rate，支持按租户/Vendor/渠道过滤。
- **FR-005**: 系统必须要求插件维护 Support Playbook 文档集（README、FAQ、Troubleshooting、Support_Policy），并在发布前通过 Support Ready Checklist 验证完整性。
- **FR-006**: 系统必须为 Incident Lifecycle 提供 SEV-0~SEV-4 模板、通报模板和沟通渠道，并强制记录检测、确认、缓解、修复、验证、复盘各阶段时间戳。
- **FR-007**: 系统必须将安全事件流程与通用 Incident 流程区分，提供隔离步骤、保密标签以及安全联系人通知，确保 24 小时内完成对外通告草稿。
- **FR-008**: 系统必须提供 SLA/SLO 配置界面，可为插件类型选择推荐等级，记录 Uptime、Avg Response、Success Rate、Support FRT 等 SLI 的目标值及当前表现。
- **FR-009**: 系统必须计算 SLA Score = 0.4 × Uptime + 0.3 × Support + 0.3 × Reliability，并在得分 ≥ 85 时自动授予推荐位加权激励，在得分 < 70 时触发处罚（警告或下架流程），并实时更新 Marketplace Dashboard。
- **FR-010**: 系统必须每日采样 SLA 数据、月度生成报告、季度输出复盘，所有操作自动记录审计日志并可导出。
- **FR-011**: 系统必须在插件 manifest 中启用支持与事件配置校验，缺失必填字段时阻止发布。
- **FR-012**: 系统必须向 Audit Log 写入与支持、事件相关的所有操作，包括附件访问、权限变更和 SLA 调整，并可由合规团队检索。

### Key Entities *(include if feature involves data)*

- **Support Ticket**: 表示租户或 Vendor 提出的支持请求，包含渠道、优先级、状态、首次响应时间、解决时间、满意度评分及关联的知识库链接。
- **Incident Record**: 表示按 SEV 分级的故障或安全事件，记录时间线节点、责任人、沟通渠道、标签（#availability/#security/#performance/#dependency）以及 RCA 文档。
- **SLA Profile**: 表示插件针对不同计划或类型的 SLA 承诺，包含 SLI 目标、实际表现、SLA Score、激励/处罚结果与对外公示内容。
- **Readiness Checklist**: 表示 Support Ready、Incident Ready、SLA Ready 三套自检清单，跟踪项状态、负责人、完成日期及阻断发布策略。

### Assumptions & Dependencies

- PowerX Support Hub、Incident Center 与 Observability Hub 已可用，并与 Marketplace 核心系统集成。
- Vendor 已完成支持邮箱、Webhook 回调的注册与验真。
- SLA 数据采集依赖现有监控探针，采样频率为每日一次，月度与季度任务由调度系统驱动。
- 全部支持/事件操作需写入既有 Audit Log 管道，用于合规与追责。
- Marketplace Dashboard 与 SLA API 共用同一指标来源，确保数据一致性。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% 新插件在 Marketplace 上线前完成 Support Ready、Incident Ready、SLA Ready 三套清单，并通过审核记录。
- **SC-002**: 上线后首个季度内，95% 工单的首次响应时间 ≤ 4 小时，平均解决时间 ≤ 24 小时。
- **SC-003**: Incident 演练或真实事件中，SEV-1 及以上事件首次通报时间 ≤ 15 分钟，RCA 提交率在 48 小时内达到 100%。
- **SC-004**: SLA Dashboard 与 `GET /api/v1/marketplace/sla/{plugin_id}` 返回的数据一致率达到 100%，并覆盖所有 Real-time、Transactional、Utility 分类插件。
- **SC-005**: 连续两个季度 SLA Score ≥ 85 的插件占比达到 70%，并记录至少一次激励或处罚执行案例。

## Clarifications

### Session 2025-10-21

- Q: SLA Score 触发激励或处罚的具体分界线应如何定义？ → A: 激励 ≥85，处罚 <70
