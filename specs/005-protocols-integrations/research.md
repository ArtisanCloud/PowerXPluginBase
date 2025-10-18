# Research Notes – Protocols & Integrations

**Date**: 2025-10-17  
**Sources**: PowerX integration guidelines、行业最佳实践（webhook 重试、Secrets 管理）、内部安全基线

## Decision 1 – A2A Envelope 字段与大小阈值
- **Decision**: 使用单一 Envelope 结构，必含 `message_id`, `trace_id`, `correlation_id`, `tenant_id`, `tool_scope`, `issued_at`, `idempotency_key`, `payload_ref`, `signature`；负载超过 1 MB 时，`payload_ref` 仅携带预签名 URL。
- **Rationale**: 统一追踪与审计字段，1 MB 阈值可避免主链路传输大数据；预签名 URL 便于权限控制和到期。
- **Alternatives considered**:
  - 直接 Base64 内嵌二进制：消息体膨胀 + 传输重试成本高，被拒绝。
  - 动态协商字段：增加适配复杂度，无法保证安全审计覆盖。

## Decision 2 – GrantMatrix 存储模式
- **Decision**: 静态 YAML 管控基础映射，数据库表存储增量/覆盖配置，并通过缓存失效广播到各服务。
- **Rationale**: 承袭版本管控优势（代码评审）同时支持运营期快速调整；数据库记录审批与历史。
- **Alternatives considered**:
  - 全静态配置：无法满足临时调整和多租户差异化需求。
  - 全动态配置：缺乏默认策略和 Git 版本追溯，不利审计。

## Decision 3 – Webhook 重试与 DLQ 策略
- **Decision**: 采用 1m→5m→15m 指数退避，最大 3 次；失败进入 DLQ，由平台 SRE 监控、插件团队执行补发；每条事件带签名与到期时间。
- **Rationale**: 三次重试兼顾及时性与队列负载；联合处理符合运营职责分工；签名保障数据安全。
- **Alternatives considered**:
  - 固定间隔重试：无法适应下游冷启动场景。
  - 仅插件团队负责：平台缺少统一监控易造成遗漏。

## Decision 4 – Secrets 生命周期
- **Decision**: 所有外部 API 凭证必须登记轮换周期（默认 30 天提醒），支持双密钥过渡；吊销立即生效并记录审计事件。
- **Rationale**: 满足安全基线与合规要求；双密钥过渡避免服务中断。
- **Alternatives considered**:
  - 手动轮换无提醒：易遗漏导致凭证过期或泄露。
  - 单密钥立即切换：对合作方造成中断风险。

## Decision 5 – 幂等缓存与速率控制
- **Decision**: 使用 Redis 作为首选幂等缓存（24 小时 TTL），若不可用则回退 PostgreSQL；GrantMatrix 执行速率限制前检查幂等结果。
- **Rationale**: Redis 低延迟适合高频调用，PostgreSQL 回退保障部署简化；幂等先决可减轻重复调用压力。
- **Alternatives considered**:
  - 仅数据库：性能不足以支撑高并发。
  - 不做回退：Redis 故障时会导致请求失败。
