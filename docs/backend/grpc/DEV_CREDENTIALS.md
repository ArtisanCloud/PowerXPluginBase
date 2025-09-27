# 插件凭证投递（开发联调）

本说明文档用于在宿主未接通通知前，手工向插件投递并持久化（加密）租户凭证，便于联调 STS Exchange。

## 接口

- 方法：POST
- 路径：`/api/v1/agent/tenants/:tenantId/credentials`
- 鉴权：开发模式可免认证；生产需接入 JWT/HMAC/mTLS（待宿主对接）
- 载荷（JSON）：
  - `plugin_id`：字符串（当前插件：`com.powerx.plugins.base`）
  - `client_id`：字符串
  - `client_secret`：字符串（明文，仅此一次传输）

## 入库与加密
- 表：`plugin_credentials`
- 字段：`tenant_id`、`plugin_id` 唯一；`client_id` 明文；`secret_ciphertext`/`iv_nonce` 为 AES-GCM 密文/随机向量；`key_version` 轮换版本
- 主密钥：`server.secret_key`（或环境变量 `POWERX_SERVER_SECRET_KEY`）；通过 SHA-256 导出 32 字节密钥
- AAD：`tenant_id|plugin_id|client_id` 绑定，防止密文移植

## 启动加载
- 启动时若 `POWERX_GRPC_UPSTREAM_TENANT_ID` 已配置，则优先从 DB 解密载入 `client_id/client_secret` 注入 STS；
- 若 DB 无记录且配置/环境变量提供了 `POWERX_STS_CLIENT_ID`/`POWERX_STS_CLIENT_SECRET`，则使用配置值。

## curl 示例

```bash
# 1) 配置 backend/etc/config.yaml：
#    server:
#      secret_key: "dev-change-me"
#    grpc_upstream:
#      tenant_id: 1

# 2) 启动插件（另一个终端）
#   make run

# 3) 投递凭证（仅开发联调）
curl -sS -X POST \
  -H 'Content-Type: application/json' \
  -d '{
        "plugin_id": "com.powerx.plugins.base",
        "client_id": "com.powerx.plugins.base.1",
        "client_secret": "secret-ONLY-FOR-DEV"
      }' \
  http://127.0.0.1:8086/api/v1/agent/tenants/1/credentials

# 4) 查看日志：应提示 Loaded STS credentials for tenant from DB（重启后生效），
#    或在后续通过 STS Exchange 获得短期 token。
```

## 注意事项
- 明文 `client_secret` 不落盘，仅加密密文存储；日志不打印明文
- 生产环境必须配置 `server.secret_key`（config.yaml），并启用认证/签名校验
- 轮换时可重复投递相同路径（后续将扩展 `rotate`/`version` 幂等控制）
