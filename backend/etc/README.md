# 配置文件说明

## 概述

本项目支持灵活的配置管理，优先级顺序为：

1. **环境变量** (最高优先级)
2. **YAML 配置文件**
3. **默认值** (最低优先级)

## YAML 配置文件

### 配置文件位置

配置文件按以下顺序查找：

1. `./etc/config.yaml` (推荐)
2. `./config.yaml`
3. `../etc/config.yaml`
4. `$CONFIG_PATH/config.yaml`

### 使用方法

1. **复制示例配置**:

   ```bash
   cp backend/etc/config.example.yaml backend/etc/config.yaml
   ```

2. **修改配置**:
   根据你的环境修改 `backend/etc/config.yaml` 中的配置项

3. **启动服务**:
   ```bash
   cd backend
   go run ./cmd/plugin
   ```

### 配置结构

```yaml
# 服务配置
server:
  bind_addr: ":8086" # 服务绑定地址
  log_level: "info" # 日志级别
  dev_mode: false # 开发模式

# 数据库配置
database:
  dsn: "host=localhost..." # 数据库连接字符串
  schema: "px_plugin" # 数据库模式

# 业务配置
business:
  sprint:
    default_capacity: 40 # 默认 Sprint 容量
  note:
    default_priority: "medium" # 默认任务优先级

# 安全配置
security:
  enable_cors: true
  cors_origins:
    - "http://localhost:3036"
```

## 环境变量

环境变量会覆盖 YAML 配置：

| 环境变量         | 说明         | 示例                |
| ---------------- | ------------ | ------------------- |
| `PX_BIND_ADDR`   | 服务绑定地址 | `:8086`             |
| `PX_LOG_LEVEL`   | 日志级别     | `debug`             |
| `PX_DEV_MODE`    | 开发模式     | `true`              |
| `PX_DB_DSN`      | 数据库连接   | `host=localhost...` |
| `PX_DB_SCHEMA`   | 数据库模式   | `px_plugin`         |
| `PX_RUN_MIGRATE` | 运行迁移     | `true`              |

### PowerX 集成环境变量

| 环境变量                 | 说明       |
| ------------------------ | ---------- |
| `PLUGIN_CTX_HMAC_SECRET` | HMAC 密钥  |
| `PLUGIN_CTX_KID`         | 密钥 ID    |
| `PX_CTX_JWKS_URL`        | JWKS URL   |
| `PX_CTX_ISSUER`          | JWT 签发者 |
| `PX_CTX_AUDIENCE`        | JWT 受众   |
| `PX_CTX_TTL`             | 上下文 TTL |

## 开发环境设置

1. **创建开发配置**:

   ```bash
   cp backend/etc/config.example.yaml backend/etc/config.yaml
   ```

2. **修改开发配置**:

   ```yaml
   server:
     dev_mode: true
     log_level: "debug"
   database:
     dsn: "your_dev_database_connection"
   runtime:
     run_migrate: true
   ```

3. **或使用环境变量**:
   ```bash
   export PX_DEV_MODE=true
   export PX_LOG_LEVEL=debug
   export PX_DB_DSN="your_dev_database_connection"
   export PX_RUN_MIGRATE=true
   ```

## 生产环境部署

1. **创建生产配置**:

   ```yaml
   server:
     bind_addr: ":8086"
     log_level: "info"
     dev_mode: false
   database:
     dsn: "your_production_database"
   context:
     hmac_secret: "your_production_secret"
   security:
     cors_origins:
       - "https://your-frontend.com"
   ```

2. **或使用环境变量** (Docker/K8s 推荐):
   ```bash
   export PX_DB_DSN="your_production_database"
   export PLUGIN_CTX_HMAC_SECRET="your_production_secret"
   ```

## 配置验证

启动时会自动验证配置，常见错误：

- ❌ `database DSN is required`: 未设置数据库连接
- ❌ `either HMAC or JWT mode must be configured in production`: 生产环境需要认证配置
- ❌ `sprint default capacity must be positive`: 业务配置无效

## 注意事项

1. **配置文件安全**: `config.yaml` 已加入 `.gitignore`，不会提交到版本控制
2. **敏感信息**: 密码、密钥等敏感信息建议使用环境变量
3. **环境隔离**: 不同环境使用不同的配置文件或环境变量
4. **配置更新**: 修改配置后需要重启服务生效
