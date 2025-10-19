# Quickstart – Protocols & Integrations Feature

## 1. 环境准备
- Go 1.24、Node 20、Redis（可选，用于幂等缓存；无需时以 PostgreSQL 回退）。
- 运行 `make dev-setup` 安装后端依赖，`npm install --prefix web-admin` 安装前端依赖。
- 确保本地或测试环境可访问 PowerX STS / Secrets Manager / EventBus。

## 2. 数据库迁移
```bash
make migrate \
  POWERX_DATABASE_URL=postgres://powerx:powerx@localhost:5432/powerx_plugin_base?sslmode=disable
```
- 新增表：`integration_grant_matrix_overrides`、`integration_webhook_attempts`、`integration_webhook_dlq`、`integration_secrets`。
- 如果 Redis 不可用，确认 PostgreSQL 中存在幂等缓存回退表。

## 3. 配置文件
1. 复制 `backend/etc/config.example.yaml` 为 `backend/etc/config.yaml`。  
2. 在 `integration` 节添加：
   ```yaml
   integration:
     idempotency:
       provider: redis # 或 postgres
       redis_url: redis://localhost:6379
       ttl_hours: 24
     envelope:
       payload_threshold_bytes: 1048576
     webhook:
       retry_policy: [60, 300, 900] # 秒
       dlq_topic: plugin.webhook.dlq
     secrets:
       rotation_days_default: 30
   ```
3. 更新 `plugin.yaml` 的 `data_usage`、`security_baseline_version`。

## 4. 开发与运行
```bash
make run            # 启动后端
cd web-admin && npm run dev  # 启动管理界面
```
- 管理 UI 路径：`http://localhost:3000/_p/com.powerx.plugins.base/admin/integration`.
- API 基础路径：`http://localhost:8086/_p/com.powerx.plugins.base/api/v1/integration`.

## 5. 关键场景验证
1. **Envelope 调用**  
   - 使用示例请求（见 contracts/openapi）调用 `/dispatch` 接口，确认响应携带 `trace_id` 与幂等信息。
   - 样例（默认代理路径）：
     ```bash
     curl -X POST \
       http://localhost:8086/_p/com.powerx.plugins.base/api/v1/integration/dispatch \
       -H 'Content-Type: application/json' \
       -H 'Authorization: Bearer <token>' \
       -d '{
         "message_id": "3f4f4f44-9d0a-4a79-9a8b-2c8b95f6b2de",
         "trace_id": "a2f47d69-4e3e-4b9a-8b78-41df6a5c76c1",
         "correlation_id": "51f26de9-4c5c-4a1f-8fd4-d8cbbdb1dca3",
         "tenant_id": "demo-tenant",
         "tool_scope": "integration.dispatch",
         "issued_at": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
         "idempotency_key": "sample-key-001",
         "payload_ref": "{\"resource\":\"/powerx/example\",\"input\":{\"foo\":\"bar\"}}",
         "metadata": {"channel": "HTTP"},
         "signature": "BASE64_SIGNATURE"
       }'
     ```
   - 重复发送相同 `idempotency_key` 可观察响应中的 `"replay": true`，同时在日志中确认幂等重放记录。
2. **GrantMatrix**  
   - 在管理界面或 API 上传 YAML 基础配置，添加数据库覆盖项，验证审批记录与缓存刷新。
3. **Webhook 重试**  
   - 运行 `./scripts/mock-webhook-target.sh 8089` 启动本地 webhook mock（默认端口 8089，可自定义）；脚本会打印请求头和 Body，便于排查。
   - 创建订阅后模拟目标 500；检查重试日志与 DLQ 记录，执行联合处理流程。
   - 通过管理端 API 创建订阅示例：
     ```bash
     curl -X POST \
       "$API_BASE/admin/integration/webhooks" \
       -H "Authorization: Bearer $TOKEN" \
       -H "Content-Type: application/json" \
       -d '{
         "event_type": "integration.envelope.dispatch",
         "target_url": "https://localhost:8089/webhooks",
         "retry_policy": [60,300,900],
         "secret": "demo-secret"
       }'
     ```
   - 通过 `GET $API_BASE/admin/integration/webhooks/{id}/attempts` 查看最近投递记录，若状态为 `DLQ` 可调用 `POST $API_BASE/admin/integration/webhooks/attempts/{attemptId}/replay` 重新排队。
4. **Secrets 轮换**  
   - 创建外部凭证，触发轮换，确认双密钥过渡和审计日志。

## 6. 测试
```bash
make test                # 后端单元、服务、仓储测试
make integration-test    # 含适配器/幂等/GrantMatrix 流程（需新增脚本）
npm run test --prefix web-admin  # 前端测试
```
- 建议使用 `./scripts/mock-webhook-target.sh` 模拟订阅方。
- 运行 `make lint` + `npm run lint --prefix web-admin` 保持风格一致。

## 7. 打包与发布
```bash
make package-pxp
```
- 确认 `dist/security/` 中生成的 webhook/Secrets 报表、OpenAPI、Nuxt 构建产物。
- 更新 release notes（包含成功率指标、轮换流程描述、迁移步骤）。
