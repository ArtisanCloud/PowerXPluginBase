# Research Log – Marketplace & Business

## Context

基于 `spec.md` 的功能目标，需要补充税务 SaaS 对接策略、离线 License 兜底机制以及 Listing Ready Checklist GraphQL API 的最佳实践，以指导后续设计与实现。

## Findings

### Decision: 采用税务 SaaS（Stripe Tax/Avalara）进行跨区域税费计算
- **Rationale**: 现有 Billing Engine 无内建税则库，与成熟 SaaS 集成可获得自动更新的税率、发票与合规支持，满足多货币、多地域的实时结算需求，并降低财会团队维护成本。
- **Alternatives considered**:
  - **本地维护税率表**：需要持续维护各国税法与发票格式，更新滞后风险高，不符合“自动计算地域税费”的 SLA。
  - **仅统一货币人工处理税费**：违背商业闭环自动化目标，且无法满足 P0 场景的合规与体验要求。

### Decision: On-Prem License 支持 72 小时离线+一次性续期令牌
- **Rationale**: 72 小时窗口覆盖周末与短期网络波动；一次性续期令牌保证管理员在断网时可人工介入，同时避免长期离线绕过计费。该策略简化 SDK 逻辑（同一路径处理线上/离线刷新），也便于审计事件追踪。
- **Alternatives considered**:
  - **24 小时离线**：窗口过短，企业客户在变更维护或跨时区假期期间可能频繁触发停用，影响体验。
  - **完全在线验证**：不切实际，On-Prem 场景常见隔离网络，违背 Spec 中的离线容错需求。

### Decision: Listing Ready Checklist 通过 GraphQL API 暴露给控制台与 CI
- **Rationale**: GraphQL 允许细粒度查询 Checklist 项目、状态与审核备注，可被前端与自动化流程复用；结合现有 admin API 前缀即可统一鉴权。同时可在 Schema 中开放 mutation 触发重新校验或标记完成。
- **Alternatives considered**:
  - **仅控制台交互**：阻断 CI/CD 自动校验，无法在提交前发现缺失资产。
  - **CLI 专用接口**：增加维护成本，且与前端逻辑重复。
