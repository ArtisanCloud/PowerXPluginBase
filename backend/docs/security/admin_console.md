# Admin Console Security Notes

本文件列出 Dev Console（`/_p/<plugin-id>/admin/dev-console` 与 `/api/v1/admin/dev-console/**`）相关的安全考量，供安全评审与日常巡检使用。

## RBAC & 权限映射

| 功能 | HTTP Path | Manifest 权限 |
|------|-----------|----------------|
| 配置管理 | `/config/sections` | `operations.plugin.admin` (`read`/`manage`) |
| 审计查询/导出 | `/audit/events`, `/audit/export` | `operations.plugin.audit` (`read`/`export`) |
| 故障排查面板、任务历史 | `/jobs/runs`, `/troubleshooting/**`, `/webhooks/attempts**` | `operations.plugin.ops` (`read`) |
| 安全操作/重试 | `/jobs/runs/{id}/retry`, `/safe-ops/actions` | `operations.plugin.ops` (`execute`) |

> Manifest 版本 ≥0.7.0 才包含上述政策映射。升级后请运行 `make release` 生成新的 manifest bundle。

## 安全要点

1. **租户隔离**
   - RLS guard：`admin_console_job_runs`、`admin_console_audit_events`、`admin_console_config_changes` 均启用 `FORCE ROW LEVEL SECURITY`，策略 `tenant_id = current_setting('app.tenant_id')`。
   - 服务层 `JobService.ListRuns`、`TroubleshootService.ListWebhookAttempts` 会在查询构建阶段校验/正则化租户 ID，避免空租户绕过。

2. **并发控制**
   - 安全操作采用应用级锁（`LockKey(plugin_id|scope_ref|action)`）防止重复执行，同步写入 Postgres 以便审计。
   - 锁释放在状态变更（`UpdateRunStatus`）时触发，如遇异常需运行：

```sql
SELECT pg_advisory_unlock(hashtext(CONCAT('com.powerx.plugins.base|', scope_ref, '|', action)))
FROM admin_console_job_runs
WHERE status IN ('pending','running');
```

3. **输入验证**
   - Safe Ops 请求要求 `scope_ref`、`action`、`actor.id` 必填，并对 `scope_type` 限定枚举。
   - Webhook Attempt 查询强制提供 `tenant_id`，否则返回 400，防御批量扫描。

4. **审计可追溯性**
   - 所有安全操作与重试都会写入 `admin_console_audit_events`，字段 `permission_code`、`action` 分别标识操作者的权限来源与执行内容。
   - 前端调用完成后建议提示用户保留操作 ID，方便关联 Job Run → Audit Event。

5. **已知风险/缓解**
   - 大批量审计导出可能导致应用内存占用增大。已启用流式导出，每批 300 条并复用 `rows.Next()` 指针，建议限制导出窗口 ≤ 31 天。
   - Webhook drill-down 与运行中 webhook 系统共享数据库，建议保持隔离角色 `powerx_plugin_base_integration_ro`，仅授予 SELECT。

## 日常巡检 Checklist

- [ ] 检查 manifest 是否已包含 `operations.plugin.ops` 权限条目。
- [ ] 验证 `/api/v1/admin/dev-console/jobs/runs` 在权限不足时返回 403（未授权账号使用 curl 测试）。
- [ ] 每周执行一次 RLS 抽样：

```sql
SET app.tenant_id = 'demo-tenant';
SELECT COUNT(*) FROM admin_console_job_runs;
```

结果应只包含当前租户数据。

- [ ] 关注 `powerx_admin_console_safe_op_total{outcome="error"}` 是否突然升高，必要时通知研发检查锁争用。

---

维护者：Platform Security @ PowerX
