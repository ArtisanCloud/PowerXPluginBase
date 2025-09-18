# PowerX Note Plugin

> 适配 **PowerX 插件生态** 的 Note 任务管理插件

这是一个完整的 Note 任务管理插件，专为 PowerX 平台设计，提供了任务管理、团队协作等核心功能。

## ✨ 功能特性

- 🎯 **任务管理**: 创建、编辑、分配、跟踪任务
- 👥 **多租户支持**: 完整的租户隔离和 RLS 安全策略
- 🔐 **权限控制**: 基于 RBAC 的细粒度权限管理
- 🤖 **Agent 集成**: 支持 PowerX Agent Hub 智能助手
- 📊 **数据分析**: 任务统计、可视化报告
- 🔌 **插件架构**: 完全符合 PowerX 插件规范

## 🏗 项目结构

```
powerx-plugin-note/
├── README.md                    # 项目说明
├── plugin.yaml                  # 插件配置文件
├── Makefile                     # 构建脚本
├── Dockerfile                   # Docker 构建文件
├── .gitignore                   # Git 忽略规则
├── docs/                        # 文档目录
│   └── readme.md                # 详细技术文档
├── backend/                     # 后端服务
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   ├── plugin/
│   │   │   └── main.go
│   │   └── database/
│   │       ├── migrate/
│   │       │   └── migrate.go
│   │       └── seed/
│   │           └── seed.go
│   └── internal/
│       ├── config/
│       │   └── config.go
│       ├── logger/
│       │   └── logger.go
│       ├── db/
│       │   ├── db.go
│       │   └── model.go
│       ├── middleware/
│       │   ├── tenant.go
│       │   ├── rbac.go
│       │   └── common.go
│       ├── contracts/
│       │   ├── manifest.go
│       │   └── api.go
│       ├── domain/
│       │   ├── note.go          # 任务模型
│       │   ├── repository.go
│       │   └── note_repository.go
│       ├── handlers/
│       │   ├── admin_handler.go
│       │   ├── note_handler.go
│       │   └── health_handler.go
│       └── router/
│           └── router.go
└── web-admin/
    ├── README.md
    └── .placeholder
```

## 🚀 快速开始

### 1. 环境准备

```bash
go version  # Go 1.21+
psql --version  # PostgreSQL
```

### 2. 配置环境

```bash
cp backend/etc/config.example.yaml backend/etc/config.yaml
# 或者模拟宿主行为：cp config/values.example.yaml config/host-values.yaml
```

### 3. 安装依赖

```bash
cd backend
go mod download
```

### 4. 数据库迁移

```bash
make migrate
go run ./cmd/database/migrate
go run ./cmd/database/seed   # 可选
```

### 5. 启动服务

```bash
make dev
make run
```

### 6. 验证服务

```bash
curl http://localhost:8091/healthz
curl http://localhost:8091/v1/admin/manifest
curl http://localhost:8091/v1/notes
```

## 🔧 开发指南

### 构建命令

```bash
make build
make test
make lint
make fmt
make clean
make all
```

### Docker 支持

```bash
make docker-build
make docker-run
```

### 开发模式特性

- 🚫 **绕过认证**: 自动注入模拟租户上下文
- 📝 **详细日志**: 输出调试信息
- 🔄 **热重载**: 支持快速重启

## 🧩 PowerX 集成

### 插件注册

1. 将项目目录放入 PowerX 的 `plugins/` 目录
2. PowerX 自动扫描 `plugin.yaml` 并注册插件
3. 插件 API 通过 `/_p/com.powerx.plugins.note/api/v1/*` 访问

### 认证集成

- **HMAC 模式**: 内网环境使用共享密钥
- **JWT 模式**: 推荐生产环境，使用 RS256/ES256

### Agent 能力

- `note.note.create`: 创建任务
- `note.note.query`: 查询任务

## 📋 API 文档

### 管理端 API

- `GET /v1/admin/manifest` - 获取插件清单
- `GET /v1/admin/rbac` - 获取权限信息

### 任务 API

- `POST /v1/notes` - 创建任务
- `GET /v1/notes` - 获取任务列表
- `GET /v1/notes/:id` - 获取任务详情
- `PUT /v1/notes/:id` - 更新任务
- `DELETE /v1/notes/:id` - 删除任务
- `PATCH /v1/notes/:id/status` - 更新任务状态

### 权限要求

| 操作     | 所需权限           |
| -------- | ------------------ |
| 查看任务 | `note:note:read`   |
| 创建任务 | `note:note:create` |
| 更新任务 | `note:note:update` |
| 删除任务 | `note:note:delete` |

## 🗃 数据模型

### Note (任务)

- 基础信息: 标题、描述、状态、优先级
- 分配信息: 分配人、标签
- 时间信息: 创建时间、更新时间、截止时间
- 扩展信息: 估算点数、元数据

## 🔒 安全特性

- **多租户隔离**
- **权限控制**
- **安全头部防护**
- **速率限制**
- **输入验证**

## 📈 监控和日志

- **结构化日志**
- **请求追踪**
- **性能监控**
- **健康检查**

## 🛠 故障排除

1. **数据库连接失败**

   ```bash
   echo $PX_DB_DSN
   psql $PX_DB_DSN -c "SELECT 1"
   ```

2. **权限错误**

   ```bash
   export PX_DEV_MODE=1
   ```

3. **RLS 策略问题**

   ```bash
   make migrate
   ```

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

---

要不要我帮你也把 `backend/domain` 和 `agent` 那些和 `sprint` 有关的文件路径一并裁掉，给你一个干净的目录结构？
