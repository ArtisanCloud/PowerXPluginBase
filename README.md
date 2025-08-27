# PowerX Scrum Plugin

> 适配 **PowerX 插件生态** 的 Scrum 任务管理插件

这是一个完整的 Scrum 任务管理插件，专为 PowerX 平台设计，提供了任务管理、Sprint 规划、团队协作等核心功能。

## ✨ 功能特性

- 🎯 **任务管理**: 创建、编辑、分配、跟踪任务
- 🏃 **Sprint 管理**: 完整的 Sprint 生命周期管理
- 👥 **多租户支持**: 完整的租户隔离和 RLS 安全策略
- 🔐 **权限控制**: 基于 RBAC 的细粒度权限管理
- 🤖 **Agent 集成**: 支持 PowerX Agent Hub 智能助手
- 📊 **数据分析**: 燃尽图、速度图表等可视化报告
- 🔌 **插件架构**: 完全符合 PowerX 插件规范

## 🏗 项目结构

```
powerx-plugin-scrum/
├── README.md                    # 项目说明
├── plugin.yaml                 # 插件配置文件
├── Makefile                     # 构建脚本
├── Dockerfile                   # Docker 构建文件
├── .gitignore                   # Git 忽略规则
├── docs/                        # 文档目录
│   └── readme.md               # 详细技术文档
├── backend/                     # 后端服务
│   ├── go.mod                  # Go 模块文件
│   ├── go.sum                  # Go 依赖锁定
│   ├── cmd/                    # 命令行工具
│   │   ├── plugin/             # 主服务入口
│   │   │   └── main.go
│   │   └── database/           # 数据库工具
│   │       ├── migrate/        # 数据库迁移
│   │       │   └── migrate.go
│   │       └── seed/           # 种子数据
│   │           └── seed.go
│   └── internal/               # 内部包
│       ├── config/             # 配置管理
│       │   └── config.go
│       ├── logger/             # 日志服务
│       │   └── logger.go
│       ├── db/                 # 数据库层
│       │   ├── db.go
│       │   └── model.go
│       ├── middleware/         # 中间件
│       │   ├── tenant.go       # 租户认证
│       │   ├── rbac.go         # 权限控制
│       │   └── common.go       # 通用中间件
│       ├── contracts/          # 契约定义
│       │   ├── manifest.go     # 插件清单
│       │   └── api.go          # API 契约
│       ├── domain/             # 领域层
│       │   ├── task.go         # 任务模型
│       │   ├── sprint.go       # Sprint 模型
│       │   ├── repository.go   # 仓储接口
│       │   └── task_repository.go # 任务仓储实现
│       ├── handlers/           # 处理器层
│       │   ├── admin_handler.go    # 管理接口
│       │   ├── task_handler.go     # 任务接口
│       │   └── health_handler.go   # 健康检查
│       └── router/             # 路由配置
│           └── router.go
└── web-admin/                  # 前端管理界面 (占位)
    ├── README.md               # 前端说明
    └── .placeholder            # 占位文件
```

## 🚀 快速开始

### 1. 环境准备

```bash
# 确保已安装 Go 1.21+
go version

# 确保有 PostgreSQL 数据库
psql --version
```

### 2. 配置环境变量

```bash
export PX_DB_DSN='postgres://user:pass@127.0.0.1:5432/powerx?sslmode=disable'
export PX_DB_SCHEMA='scrum'
export PX_BIND_ADDR=':8091'
export PX_LOG_LEVEL='debug'
export PX_DEV_MODE=1  # 开发模式
```

### 3. 安装依赖

```bash
cd backend
go mod download
```

### 4. 数据库迁移

```bash
# 运行迁移
make migrate

# 或者使用独立命令
go run ./cmd/database/migrate

# 可选：添加种子数据
go run ./cmd/database/seed
```

### 5. 启动服务

```bash
# 开发模式启动
make dev

# 或者直接运行
make run
```

### 6. 验证服务

```bash
# 健康检查
curl http://localhost:8091/healthz

# 获取插件清单
curl http://localhost:8091/v1/admin/manifest

# 测试任务 API (开发模式)
curl http://localhost:8091/v1/tasks
```

## 🔧 开发指南

### 构建命令

```bash
# 构建二进制文件
make build

# 运行测试
make test

# 代码检查
make lint

# 格式化代码
make fmt

# 清理构建文件
make clean

# 完整构建流程
make all
```

### Docker 支持

```bash
# 构建镜像
make docker-build

# 运行容器
make docker-run
```

### 开发模式特性

在开发模式下 (`PX_DEV_MODE=1`)，插件提供以下便利功能：

- 🚫 **绕过认证**: 自动注入模拟租户上下文
- 📝 **详细日志**: 输出详细的调试信息
- 🔄 **热重载**: 支持代码变更后快速重启

### 数据库管理

```bash
# 重置数据库 (危险操作)
make reset-db

# 仅运行迁移
PX_RUN_MIGRATE=true go run ./cmd/plugin/main.go
```

## 🧩 PowerX 集成

### 插件注册

1. 将项目目录放入 PowerX 的 `plugins/` 目录
2. PowerX 会自动扫描 `plugin.yaml` 并注册插件
3. 插件 API 通过 `/_p/com.powerx.plugins.scrum/api/v1/*` 访问

### 认证集成

插件支持两种认证模式：

- **HMAC 模式**: 适合内网环境，使用共享密钥
- **JWT 模式**: 推荐生产环境，使用 RS256/ES256

### Agent 能力

插件注册了以下 Agent 工具：

- `scrum.task.create`: 创建任务
- `scrum.task.query`: 查询任务
- `scrum.sprint.create`: 创建 Sprint
- `scrum.plan.generate`: 生成 Sprint 计划

## 📋 API 文档

### 管理端 API

- `GET /v1/admin/manifest` - 获取插件清单
- `GET /v1/admin/rbac` - 获取权限信息

### 任务 API

- `POST /v1/tasks` - 创建任务
- `GET /v1/tasks` - 获取任务列表
- `GET /v1/tasks/:id` - 获取任务详情
- `PUT /v1/tasks/:id` - 更新任务
- `DELETE /v1/tasks/:id` - 删除任务
- `PATCH /v1/tasks/:id/status` - 更新任务状态

### 权限要求

| 操作 | 所需权限 |
|------|----------|
| 查看任务 | `scrum:task:read` |
| 创建任务 | `scrum:task:create` |
| 更新任务 | `scrum:task:update` |
| 删除任务 | `scrum:task:delete` |
| 管理 Sprint | `scrum:sprint:*` |

## 🗃 数据模型

### Task (任务)
- 基础信息: 标题、描述、状态、优先级
- 分配信息: 分配人、Sprint、标签
- 时间信息: 创建时间、更新时间、截止时间
- 扩展信息: 估算点数、元数据

### Sprint
- 基础信息: 名称、目标、状态
- 时间信息: 开始日期、结束日期
- 容量管理: 团队容量（故事点）
- 关联信息: 包含的任务列表

## 🔒 安全特性

- **多租户隔离**: 应用层 + 数据库 RLS 双重保障
- **权限控制**: 基于资源和动作的细粒度权限
- **安全头部**: 防止 XSS、点击劫持等攻击
- **速率限制**: 防止 API 滥用
- **输入验证**: 严格的请求参数验证

## 📈 监控和日志

- **结构化日志**: 使用 logrus 提供结构化日志
- **请求追踪**: 自动记录所有 HTTP 请求
- **性能监控**: 请求延迟和状态码统计
- **健康检查**: 提供服务和依赖健康状态

## 🛠 故障排除

### 常见问题

1. **数据库连接失败**
   ```bash
   # 检查环境变量
   echo $PX_DB_DSN
   
   # 测试连接
   psql $PX_DB_DSN -c "SELECT 1"
   ```

2. **权限错误**
   ```bash
   # 开启开发模式绕过认证
   export PX_DEV_MODE=1
   ```

3. **RLS 策略问题**
   ```bash
   # 重新运行迁移
   make migrate
   ```

### 调试技巧

- 设置 `PX_LOG_LEVEL=debug` 获取详细日志
- 使用 `PX_DEV_MODE=1` 进行本地调试
- 检查 `/healthz` 端点确认服务状态

## 📚 相关文档

- [PowerX 插件开发指南](docs/readme.md)
- [API 参考文档](docs/api.md)
- [部署指南](docs/deployment.md)
- [贡献指南](docs/contributing.md)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目！

## 📄 许可证

MIT License

---

**注意**: 本项目是一个初始化的代码骨架，提供了完整的项目结构和核心功能实现。web-admin 前端部分将在后续单独实现。