# 配置文件说明

## 总览

Base 插件的配置来源按照以下优先级生效：

1. **PowerX 宿主注入的环境变量**（如 `POWERX_BIND_ADDR`、`POWERX_LOG_LEVEL`）
2. **YAML 配置文件**（`config/host-values.yaml` 或 `config.yaml`）
3. **内置默认值**（`backend/internal/config/getDefaultConfig`）

宿主环境会在安装阶段生成配置文件并通过 `CONFIG_PATH` 环境变量传递目录；
本地开发可直接使用仓库内的示例文件。

## 配置查找优先级

插件启动时会依次尝试加载以下文件，找到首个存在且可读的 YAML：

1. `$CONFIG_PATH/host-values.yaml`
2. `$CONFIG_PATH/config.yaml`
3. `./config/host-values.yaml`
4. `./config/config.yaml`
5. `./etc/config.yaml`
6. `./config.yaml`
7. `../config/host-values.yaml`
8. `../config/config.yaml`
9. `../etc/config.yaml`

如果以上文件均不存在，则仅使用内置默认值并在日志中给出警告。

## PowerX 宿主配合流程

1. 读取仓库中的 `config/schema.yaml`，根据当前 PowerX 主系统的数据库、Redis、鉴权等配置
   生成 `config/host-values.yaml`。
2. 将生成结果写入插件安装目录，例如 `plugins/installed/<id>/<version>/config/host-values.yaml`。
3. 启动插件进程前设置：
   - `CONFIG_PATH=/path/to/plugins/<id>/<version>/config`
   - `POWERX_BIND_ADDR`、`POWERX_PLUGIN_PORT` 等运行时变量（端口由宿主分配）。
4. 重新启用插件即可让新配置生效。

## 本地开发流程

1. 复制示例配置：

   ```bash
   cp backend/etc/config.example.yaml backend/etc/config.yaml
   ```

2. 如需模拟宿主行为，可先生成一份自定义 `config/host-values.yaml`，再通过 `CONFIG_PATH` 指向该目录：

   ```bash
   CONFIG_PATH=./config go run ./cmd/plugin
   ```

3. 若直接使用 `backend/etc/config.yaml`，可以省略 `CONFIG_PATH`：

   ```bash
   cd backend
   go run ./cmd/plugin
   ```

## 安全基线配置

- `backend/etc/security_baseline.yaml` 提供默认的安全与合规基线，可在部署时由宿主覆盖。
  - `baseline_version`：基线版本号，便于审计追踪。
  - `masking_rules`：日志和数据脱敏规则（PII 列表、掩码占位符等）。
  - `tool_grant`：ToolGrant 生命周期配置（TTL、续期阈值、登出即撤销）。
  - `consent_defaults`：宿主未下发策略时的隐私默认值，例如 90 天数据保留、审计日志通道、导出目的地。
- 宿主可在其配置包中注入替代文件，并结合 `CONFIG_PATH` 定位，实现环境化定制。

## 常用环境变量

| 变量名 | 说明 |
| ------ | ---- |
| `CONFIG_PATH` | 指定配置目录，优先读取 `host-values.yaml` |
| `POWERX_BIND_ADDR` | HTTP 监听地址（宿主自动分配端口） |
| `POWERX_LOG_LEVEL` | 覆盖日志级别，默认 `info` |
| `POWERX_DEV_MODE` | 是否启用开发模式，`true/false` |
| `POWERX_DB_DSN` | 若需要从环境层覆盖数据库连接，可设置该变量 |
| `POWERX_DB_SCHEMA` | 覆盖数据库 Schema 名称 |
| `POWERX_RUN_MIGRATE` | 设置为 `true` 时强制执行数据库迁移 |
| `PLUGIN_CTX_HMAC_SECRET` | PowerX 下发的 HMAC 签名密钥 |
| `POWERX_CTX_JWKS_URL` | JWT 模式下的 JWKS 地址 |

> ✅ 建议在生产环境通过 `config/host-values.yaml` 写入敏感配置，仅在必要时才用环境变量覆盖。

## Admin Console 配置

`admin_console` 节点控制插件运营控制台的保留策略与刷新节奏：

- `audit_retention_days` / `config_change_retention_days`：审计日志与配置历史保留天数，默认 365。
- `job_history_days`：安全操作与任务历史保留天数，默认 45。
- `export.max_range_days`：单次审计导出的最大时间跨度，默认 31 天；`default_format` 支持 `csv` 或 `json`。
- `troubleshooting.refresh_interval_seconds`：故障排查仪表盘自动刷新间隔；`cache_ttl_seconds` 控制后端缓存 TTL。
- `safe_ops.lock_ttl_seconds` / `max_concurrent_ops`：安全操作的锁超时时长及并发上限，避免重复执行。
- `navigation_hidden_routes`：可选，填入路径以在特定环境隐藏部分控制台菜单。

## 常见错误与排查

- `database DSN is required`：未在配置文件或环境变量中提供数据库连接。
- `either HMAC or JWT mode must be configured in production`：生产模式至少需要配置 HMAC 或 JWT 方案。
- `logging level must be one of: debug, info, warn, error`：日志级别写错。

## 版本控制与安全

- `backend/etc/config.yaml` 与 `config/host-values.yaml` 默认被列入 `.gitignore`，不会提交到仓库。
- 请勿在示例配置中填写真实生产密钥；通过 CI/CD 或 PowerX 安装器在部署时注入。
