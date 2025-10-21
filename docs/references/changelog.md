# Changelog

> 仅记录与 PowerX Marketplace 商业闭环相关的可交付版本。更早的安全、协议迭代请参考 `docs/releases/` 历史记录。

## v0.5.0 · Marketplace Analytics Loop

发布日期：2025-10-21

- **Usage 聚合与告警**：新增 `marketplace_usage_envelopes` / `marketplace_usage_aggregates` / `marketplace_notifications` 表，提供批量幂等上报、配额/Spike 告警以及 GDPR 删除链路。  
- **Revenue Share 报表**：自动计算 Vendor / Platform / Fee 分润，Expose `/marketplace/revenue-share/reports` API 与 Dashboard 导出。  
- **Admin Dashboard**：在 `/_p/com.powerx.plugins.base/admin/integration/marketplace/dashboard` 展示调用趋势、剩余额度、异常告警与分润结果，并上报首屏性能指标。  
- **文档更新**：补充 [Marketplace 商业闭环指南](../overview/marketplace_business_loop.md)，README 与 Quickstart 指引完整闭环演练。

升级提示：

1. 执行 `make migrate` 以创建新的 Usage/Revenue 表结构及 RLS 策略。  
2. 更新 `backend/etc/config.yaml` 中的 `integration.billing.reconciliation` 与 `marketplace` 节点，确认分润比例、提醒渠道与 Redis 配置。  
3. 重新构建后端与前端包：`make build && make frontend-build && make dist`。  
4. 若生产环境需要 Usage 历史迁移，请在升级前导出旧数据或利用 `PrivacyService.PurgeUsageData` 工具清理不需要的 Envelope。
