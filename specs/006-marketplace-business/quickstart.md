# Quickstart – Marketplace & Business Feature

## Prerequisites
- Go 1.24（使用 `make dev-setup` 安装依赖与 `golangci-lint`）
- Node 20 + pnpm（`web-admin` Nuxt 4）
- Postgres ≥ 13 (`powerx_plugin_base` schema) 与 Redis（幂等、License 缓存）
- Stripe Tax / Avalara 测试账户，配置凭证写入 `.env` 或 `config.yaml`

## Backend
```bash
# 1. 同步配置
cp backend/etc/config.example.yaml config.yaml   # 如未存在
```

在 `config.yaml` 中新增：
```yaml
integration:
  billing:
    tax_provider: stripe_tax
    stripe_tax:
      api_key: sk_test_xxx
```

```bash
# 2. 运行数据库迁移
make migrate

# 3. 启动插件后端（含 Marketplace API）
make run
```

- API 入口：`http://localhost:8086/api/v1/marketplace/**`
- `/healthz` 检查服务状态

## Admin 前端
```bash
cd web-admin
pnpm install
pnpm dev   # 本地运行，默认绑定 http://localhost:3000/_p/com.powerx.plugin.base/admin
```

在 `.env` 中设置：
```
NUXT_PUBLIC_API_BASE_URL=http://localhost:8086/api/v1
```

## Sandbox 测试流程
1. 上传 `.pxp` 包 → `POST /marketplace/listings`
2. 触发 Checklist GraphQL `triggerChecklistRun` → 确保全部通过
3. Reviewer 将 Listing 标记为 `published`
4. 使用租户管理员调用 `POST /marketplace/licenses` 购买 Subscription 计划
5. 通过 SDK 上报 Usage → `POST /marketplace/usage`
6. 在 Vendor Dashboard 页面验证趋势图与告警
7. 检查 `GET /marketplace/revenue-share/reports` 的月度报表

## 监控与告警
- 关键指标：`license_verify_latency`, `usage_ingest_lag`, `tax_provider_errors`
- 事件总线：`license.events.*`, `usage.spike.detected`, `billing.tax.failure`

## 下一步
- 参考 `/specs/006-marketplace-business/data-model.md` 完成迁移
- 按 `/specs/006-marketplace-business/contracts/` 实现 API/GraphQL
- 更新 `plugin.yaml` manifest 与 RBAC 映射
