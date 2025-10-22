# Quickstart — Support & Operations Playbook

## Prerequisites
- 已执行 `make dev-setup && npm install`，本地环境运行 `make run` 与 `npm run dev`。
- 确保 `config/docker-compose.integration.yml` 中的 Redis / Webhook mock 已启动（`docker compose -f config/docker-compose.integration.yml up -d`）。
- `backend/etc/config.yaml` 追加 `operations` 配置节点（支持渠道、Webhook 签名密钥、SLA 采样周期）。

## 1. 配置 Support Playbook
1. 打开 Admin Console：`/_p/com.powerx.plugins.base/admin/operations/support`。
2. 编辑支持渠道：填充 Marketplace 工单系统、Vendor 邮箱、In-App Feedback、紧急热线，设置服务窗口与升级路径。
3. 上传或链接 README / FAQ / Troubleshooting / Support Policy。
4. 点击 “Validate Channels” 触发合成检测；在 `operations_support_ticket_events` 中确认 `webhook_status=delivered`。
5. Support Ready Checklist 全部勾选后，页面显示 “Support Ready✅”。

## 2. 演练 Incident 生命周期
1. 进入 `/_p/com.powerx.plugins.base/admin/operations/incidents`，点击 “Declare Incident”。
2. 选择 `SEV-1`、标签 `#availability`，填写影响范围，提交后确认 15 分钟内的下一次通报时间已排程。
3. 在同页追加 Timeline Entry → 发布到 Support Hub；核对 webhook mock 收到 `incident.update` 事件。
4. 应用热修复后，更新状态为 `resolved` 并上传 RCA 文档；Incident Ready Checklist 标记完成。

## 3. 校验 SLA Dashboard & API
1. 运行定时任务模拟：`go run ./backend/cmd/cron --task operations-sla-refresh`（待实现任务）。
2. 打开 `/_p/com.powerx.plugins.base/admin/operations/sla`，确认 Real-time / Transactional / Utility 计划的指标已填充。
3. 调用公共接口：`curl https://localhost:8086/api/v1/marketplace/sla/com.powerx.plugins.base`；返回体需包含 `slaScore` 与 `lastUpdated`。
4. 将 SLA Score 调整到 90，确认 Marketplace Dashboard 展示推荐标识；将 Score 调整到 60，检查处罚提示与通知。

## 4. 审核审计与指标
1. 在数据库查看 `operations_support_ticket_events`、`operations_incident_updates`、`operations_sla_adjustments`，确保带 `tenant_id` 与 `plugin_id` 字段。
2. 使用 `make test` 运行支持与 incident 服务单测、集成测；前端执行 `npm run test -- operations`（待新增）。
3. 通过 `make dist` 生成发布包，核查 `dist/<version>/` 中包含新的 operations UI 构建产物。

> 若需更详细的页面操作指南，请参考 `docs/integration/07_support_and_operations/Operations_Runbook.md`。

> 完成以上步骤后，可在 `/specs/007-support-and-operations/tasks.md` 中拆分开发任务并推动实现。
