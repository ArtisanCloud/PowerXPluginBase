# Feature Specification: Protocols & Integrations (A2A, HTTP/gRPC/MCP, Webhooks & Events, Secrets)

**Feature Branch**: `005-protocols-integrations`  
**Created**: 2025-10-17  
**Status**: Draft  
**Input**: User description: "Title: Protocols & Integrations (A2A, HTTP/gRPC/MCP, Webhooks & Events, Secrets); WHAT/WHY: Provide unified integration patterns so plugins and PowerX interoperate via A2A protocol, transports (HTTP/gRPC/MCP), webhooks/event subscriptions, and external API secret handling; Scope: A2A message shapes & correlation; transport adapter interfaces (HTTP/gRPC/MCP) and handshake; ToolScopes & GrantMatrix mapping to endpoints; webhook event schema & retry/backoff; external API credential lifecycle & rotation; idempotency & retry semantics; Out-of-Scope: Business orchestration flows; UI command surfaces; marketplace routing; Dependencies/Assumptions: Capability contracts exist; EventBus topics & DLQ configured; Security & Tool Grants enforced; PowerX gRPC SDK available to plugins; gateway reverse-proxy supports per-plugin routing"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 统一 A2A 协议投递 (Priority: P1)

作为插件集成工程师，我需要使用统一的消息封装和校验规则，从而在各类宿主接口（HTTP、gRPC、MCP）之间保持一致的调用体验，无需重复设计安全字段或追踪信息。

**Why this priority**: 没有统一协议，将持续造成集成失败和审计缺口，是整个项目的根基。  
**Independent Test**: 通过仅实现 Envelope 定义与 ToolScope 校验即可驱动一次端到端互操作演示（插件调用宿主 API 并得到可追踪响应）。

**Acceptance Scenarios**:

1. **Given** 插件开发者提交带标准字段的 A2A 请求，**When** 宿主收到后验证 ToolScope 与幂等键，**Then** 请求被接受且响应包含对应的 trace 信息。
2. **Given** 请求缺少必填字段或携带无效 ToolScope，**When** 宿主校验，**Then** 返回清晰的拒绝原因和修复指引。

---

### User Story 2 - Webhook 与事件可靠交付 (Priority: P2)

作为平台 SRE，我需要统一的 Webhook/事件订阅、重试和 DLQ 协作流程，确保通知在客户端不可用时仍能追踪、补偿，并让插件团队能及时处理失败。

**Why this priority**: 事件通知直接影响安全告警和运营动作，必须在协议基础完成后尽快保障可靠性。  
**Independent Test**: 只需部署事件订阅与重试队列，即可模拟下游失败并验证重投与协同流程。

**Acceptance Scenarios**:

1. **Given** Webhook 目标短暂不可达，**When** 系统自动按照退避策略重试，**Then** 成功通知在 SLA 内送达且指标记录重试次数。
2. **Given** Webhook 多次失败进入 DLQ，**When** SRE 触发联合处理流程，**Then** 插件团队收到事件并完成补发确认。

---

### User Story 3 - 外部 Secrets 生命周期治理 (Priority: P3)

作为安全合规负责人，我需要集中管理外部 API 凭证的创建、轮换、吊销和审计记录，以避免长期凭证泄露风险，并支持租户隔离追踪。

**Why this priority**: 虽然依赖前两项基础能力，但及时治理 Secrets 能显著降低安全风险。  
**Independent Test**: 仅构建 Secrets 生命周期接口及审计记录，即可单独验证轮换提醒与审批链闭环。

**Acceptance Scenarios**:

1. **Given** 秘钥到达轮换窗口，**When** 系统提醒负责人并完成新旧密钥并行期，**Then** 旧密钥在宽限期后自动失效并记录审计事件。
2. **Given** 外部集成被撤销，**When** 安全负责人吊销凭证，**Then** 所有相关调用立即被拒绝且触发告警。

---

### Edge Cases

- 订阅方长时间不可用或返回非标准错误码时，系统如何控制重试次数与通知优先级？  
- 同一 ToolScope 在短时间内被多个插件实例同时调用时，如何确保幂等和速率限制不冲突？  
- 秘钥轮换期间新旧凭证都可能生效，如何避免出现双重计费或重复调用？  
- EventBus 或 DLQ 本身发生故障时，如何保留待处理通知并恢复处理顺序？

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 系统必须定义一个包含身份、追踪、幂等与签名字段的标准 A2A Envelope，并要求所有插件→宿主调用使用该封装。  
- **FR-002**: 系统必须提供针对 HTTP、gRPC、MCP 三类通道的统一适配器行为，包括请求验证、错误映射与超时约定。  
- **FR-003**: 系统必须维护 GrantMatrix，将 ToolScope 与允许访问的资源、操作动作和通道建立关系，并在请求进入服务层前执行校验。  
- **FR-004**: 系统必须为 Webhook/事件订阅提供签名校验、退避重试（至少 3 次）与 DLQ 协同处理能力，并记录交付状态。  
- **FR-005**: 系统必须允许安全负责人创建、轮换、吊销外部 API 凭证，并在每一步生成不可篡改的审计日志。  
- **FR-006**: 系统必须在消息或事件超过 1 MB 载荷阈值时，仅在 Envelope 中传递引用（如预签名链接）并强制验证有效期。  
- **FR-007**: 系统必须提供幂等键管理与缓存失效策略，确保重复请求返回一致结果且可追踪。  
- **FR-008**: 系统必须向平台运营与插件团队提供仪表盘或报表视图，展示请求成功率、重试率、轮换进度等关键指标。  
- **FR-009**: 系统必须为配置变更（GrantMatrix 更新、订阅新增、凭证轮换）提供审批或双人复核流程记录。  
- **FR-010**: 系统必须在可观测性层面输出统一的日志、指标与告警信号，支持按租户和 ToolScope 进行过滤分析。

### Key Entities

- **Integration Envelope**: 描述任何插件与宿主交互时的标准消息封装，包含身份、追踪、幂等、签名与 payload 引用信息。  
- **Grant Matrix**: 表示 ToolScope 与资源操作、通道、限制条件之间的映射关系，并记录审批历史与版本。  
- **Webhook Subscription**: 记录订阅目标、事件类型、签名密钥、重试策略及最近交付状态。  
- **Secret Credential**: 描述外部 API 或 Webhook 的凭证元数据、轮换计划、关联租户及审计日志。  
- **Delivery Attempt**: 跟踪每次 Webhook 或事件交付的状态、重试次数、最终结果和责任归属。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 在功能上线后 30 天内，95% 的插件→宿主调用都使用统一 Envelope，并通过平台仪表盘可验证。  
- **SC-002**: Webhook/事件通知的成功交付率在任意 30 天滚动窗口内不低于 99%，且任意单事件平均重试次数 ≤ 2。  
- **SC-003**: 外部 API 凭证轮换延迟（从提醒到完成）在 90% 情况下不超过 3 个工作日，相关审计记录可查询。  
- **SC-004**: 通过统一适配器后，集成相关的支持工单数量在三个月内同比下降 40% 以上。  
- **SC-005**: 关键集成指标（幂等冲突、权限拒绝、DLQ 积压）均可在 5 分钟内触发告警并告知责任团队。
