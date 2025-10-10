**平台统一 Agent Hub**、**插件只注册能力**、**单实例多租户隔离（RLS）**、**PowerX 反代路径 `/_p/:id/...`**、**安装后动态注入签名材料（HMAC/JWT）**，并覆盖本地联调与生产部署。

---

# powerx-plugin-base

> 适配 **PowerX 插件生态** 的最小可用示例：
>
> - 后端：Gin + GORM + Postgres（**独立 Schema**）
> - 多租户：应用层强制作用域 + Postgres **RLS 兜底**
> - 与 PowerX 的集成：**Plugin Manager 扫描/反代** + 菜单/RBAC 上报
> - Agent：**平台统一 Agent Hub**（插件注册 Profile/Tool/Workflow 能力，UI 由 PowerX 统一提供），插件可选**内嵌挂件**

---

## ✨ 功能概览

- **租户隔离**：所有表包含 `tenant_id`，请求进入后开启事务并 `SET LOCAL app.tenant_id=?`；DB 侧开启 **RLS** 策略
- **Schema 隔离**：插件复用宿主数据库实例，但使用**独立 schema**（例如 `template`），避免污染 `public`
- **对接 PowerX**：

  - 静态管理页：`/_p/:id/admin/*` → 插件前端（可选）
  - 业务 API：`/_p/:id/api/*` → 反代到插件 `/v1/...`
  - 插件上报：`/admin/manifest`、`/admin/rbac`

- **Agent 能力**：插件注册 Agent Profile / Tools / Workflows → 由 **PowerX Agent Hub** 统一呈现与调度
- **打包发布**：`plugin.yaml + backend/bin/plugin (+ web-admin)` → 上传至插件市场

---

## 🗂 目录结构

```
powerx-plugin-base/
├─ README.md
├─ plugin.yaml
├─ Makefile
├─ .gitignore
├─ Dockerfile
├─ backend/
│  ├─ go.mod
│  ├─ go.sum
│  ├─ cmd/
│  │  └─ plugin/
│  │     └─ main.go
│  │  └─ database/
│  │     └─ migrate/
│  │        └─ migrate.go
│  │     └─ seed/
│  │        └─ seed.go
│  ├─ internal/
│  │  ├─ config/        # 读取 POWERX_* 环境变量
│  │  ├─ logger/
│  │  ├─ db/            # Connect/AutoMigrate/RLS/事务封装
│  │  ├─ router/        # 路由与 handler 挂载
│  │  ├─ middleware/    # Tenant 上下文/（可选）RBAC 校验
│  │  ├─ contracts/     # 与 PowerX 的契约模型/事件Topic 等（示例）
│  │  ├─ domain/        # model/repo/service
│  │  └─ handlers/
│  │     ├─ admin_handler.go    # 上报菜单/权限（root 视角）
│  │     └─ template_handler.go     # 租户态业务接口
└─ web-admin/           # 占位；由 Next/Nuxt 产物丢这里（可选）
   ├─ .placeholder
   └─ README.md
```

> ✅ 建议把 `cmd/database/migtate/` 拼写修正为 `migrate/`，防止踩坑。

---

## ⚙️ 配置（ENV）

插件后端读取以下环境变量（由 **PowerX Plugin Manager** 在安装/启用时**注入**）：

| 变量名                   | 说明                                           | 示例                                                        |
| ------------------------ | ---------------------------------------------- | ----------------------------------------------------------- |
| `POWERX_BIND_ADDR`           | 插件监听地址                                   | `:8091`                                                     |
| `POWERX_DB_DSN`              | Postgres 连接串                                | `postgres://user:pwd@127.0.0.1:5432/powerx?sslmode=disable` |
| `POWERX_DB_SCHEMA`           | 插件独立 schema                                | `template`                                                      |
| `POWERX_LOG_LEVEL`           | 日志级别                                       | `info` / `debug`                                            |
| `POWERX_RUN_MIGRATE`         | 仅执行迁移后退出                               | `true/false`                                                |
| `POWERX_DEV_MODE`            | 开启本机调试旁路（仅开发）                     | `1`                                                         |
| `PLUGIN_CTX_HMAC_SECRET` | \*\*（HMAC 模式）\*\*平台 → 插件上下文签名密钥 | base64 编码的 32B                                           |
| `PLUGIN_CTX_KID`         | \*\*（HMAC）\*\*密钥标识                       | `com.xxx:v1`                                                |
| `POWERX_CTX_JWKS_URL`        | \*\*（JWT 模式）\*\*平台公钥 JWKS              | `http://powerx/_p/_internal/jwks`                           |
| `POWERX_CTX_ISSUER`          | 上下文签发者                                   | `powerx-auth`                                               |
| `POWERX_CTX_AUDIENCE`        | 受众                                           | `powerx-plugin`                                             |
| `POWERX_CTX_TTL`             | 上下文有效期                                   | `300s`                                                      |

> 说明：HMAC/JWT 二选一。**生产建议 JWT（RS256/ES256）**；HMAC 适合内网/开发期。

---

## 🏃 快速开始

### 1) 数据库准备

```bash
export POWERX_DB_DSN='postgres://user:pass@127.0.0.1:5432/powerx?sslmode=disable'
export POWERX_DB_SCHEMA='px_com_powerx_plugins_base'
```

### 2) 迁移与种子

```bash
# 方式 A：GORM AutoMigrate（main.go 内置）
POWERX_RUN_MIGRATE=true go run ./backend/cmd/plugin/main.go

# 方式 B：自定义命令（可在 cmd/database/migrate/migrate.go 里实现）
go run ./backend/cmd/database/migrate
go run ./backend/cmd/database/seed
```

> 迁移脚本需包含：
>
> - 创建 `template` schema
> - 业务表（含 `tenant_id`）
> - 开启 **RLS** 与策略（见下文）

### 3) 本地启动（插件单跑）

```bash
export POWERX_BIND_ADDR=":8091"
export POWERX_LOG_LEVEL="debug"
export POWERX_DEV_MODE=1   # 开启 dev bypass，仅开发期
go -C backend run ./cmd/plugin
curl :8091/healthz
```

### 4) 与 PowerX 联调（推荐）

1. 启动 **PowerX**（包含 Plugin Manager + DynamicRouter）
2. 将本仓库放入（或软链到）PowerX 的 `plugins/` 目录
3. PowerX 扫描 `plugin.yaml`，挂载反代：

   - `/_p/<plugin-id>/admin/*` → `web-admin` 产物
   - `/_p/<plugin-id>/api/*` → `http://127.0.0.1:8091`

4. 验证：

```bash
curl "http://localhost:8080/_p/com.powerx.plugins.base/api/v1/ping"
curl "http://localhost:8080/_p/com.powerx.plugins.base/api/v1/templates" -H "X-PowerX-CTX: ..." -H "X-PowerX-CTX-JWT: ..."
```

> **注意**：你的 DynamicRouter 不会重复拼 `basePath`，所以前端必须请求 `/_p/:id/api/**v1**/...`，插件内部路由也挂载在 `/v1`。

---

## 🧩 plugin.yaml 关键字段（示例）

```yaml
id: com.powerx.plugins.base
name: Base Template Plugin
version: 0.1.0
backend:
  entry: backend/bin/plugin
  port: 8091
  health: /healthz
routes:
  basePath: /v1 # 插件自身 API 前缀
  adminManifest: /api/v1/admin/manifest
  rbac: /api/v1/admin/rbac
permissions:
  - resource: base:template
    actions: [read, create, update, delete]
menus:
  - id: "plugins.base"
    title: "menu.base.template"
    icon: "i-heroicons-clipboard-document-check"
    path: "/plugins/base"
    order: 20
assets:
  webAdminPath: web-admin/.output # 可选
```

---

## 🔐 多租户与 RLS

**模型要求**：所有业务表必须包含

```sql
tenant_id BIGINT NOT NULL;
CREATE INDEX IF NOT EXISTS idx_<tbl>_tenant ON <schema>.<tbl>(tenant_id);
```

**RLS 迁移（示例）**：

```sql
CREATE SCHEMA IF NOT EXISTS template;

CREATE TABLE IF NOT EXISTS template.template (
  id         BIGSERIAL PRIMARY KEY,
  tenant_id  BIGINT NOT NULL,
  title      VARCHAR(200) NOT NULL,
  status     VARCHAR(24)  NOT NULL DEFAULT 'todo',
  assignee   BIGINT NULL,
  meta       JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ NULL
);

ALTER TABLE template.template ENABLE ROW LEVEL SECURITY;

-- 依赖每请求事务里的：SET LOCAL app.tenant_id = <tenant>
CREATE POLICY p_tenant_isolation ON template.template
  USING (tenant_id::text = current_setting('app.tenant_id', true));
```

**应用层约束**：

- 路由启用 **TenantContext** 中间件（JWT/HMAC 验签）→ 取出 `tenant_id`
- **每个请求开启事务**并执行

  ```sql
  SET LOCAL app.tenant_id = ?
  ```

- 所有 Repo 使用同一事务句柄；若遗漏 where 条件，RLS 仍会兜底

> 开发期可用 `POWERX_DEV_MODE=1` 开启受控旁路（仅本机，**上线务必关闭**）。

---

## 🗺 管理端接口（供 PowerX 拉取）

- `GET /api/v1/admin/manifest`
  返回前端菜单、入口信息
- `GET /api/v1/admin/rbac`
  返回资源与动作清单（如 `base:template` 的 CRUD）

### RBAC 如何与宿主合并？

1. **插件上报资源/动作**：`/api/v1/admin/rbac` 只需要返回插件自有的资源树，例如 `scrm:lead`、`ecommerce:order` 等。
2. **Plugin Manager 聚合**：PowerX 会定期（或在安装时）拉取所有启用插件的 RBAC 描述，并与平台内建资源合并到统一的权限仓库。
3. **Settings → Permission UI 映射**：宿主的「系统设置/权限」页面会读取聚合后的资源树，把插件声明的 `resource + action` 暴露给租户管理员配置角色。插件无需重复实现角色/团队管理页，只要信任宿主会在授权时回传带有权限数组的上下文。
4. **请求态校验**：插件在业务接口中读取上下文的 `permissions`（或调用宿主的权限校验 SDK）来决定是否放行。模板工程在中间件中预留了钩子，可按需补充。

> ✅ 结论：插件必须提供 `/admin/rbac` 接口，但**不需要**自建一套“角色/团队” UI；PowerX Settings 会统一展示并合并权限配置。

---

## 🤖 与 PowerX Agent Hub 的集成

本插件**不自带 Chat UI**。把 Agent 能力注册到平台，统一在 **Agent Hub** 呈现：

### 1) 注册清单（示例）

```yaml
agents:
  - id: "base.assistant"
    plugin_id: "com.powerx.plugins.base"
    name: "Note 助理"
    description: "创建/分配任务，生成 Sprint 计划"
    default_tools: ["template.template.create", "template.template.query"]

tools:
  - id: "template.template.create"
    plugin_id: "com.powerx.plugins.base"
    name: "创建任务"
    transport: "grpc" # 或 http
    endpoint: "127.0.0.1:51031"
    rbac_resource: "base:template"
    input_schema: { ... }
    output_schema: { ... }

workflows:
  - id: "template.plan.generate"
    plugin_id: "com.powerx.plugins.base"
    name: "生成 Sprint 计划"
    endpoint: "grpc://127.0.0.1:51031/workflows/plan_generate"
```

> 注册方式：
>
> - 启动时调用 PowerX 的注册接口（例如 `POST /api/v1/agents/register`），或
> - 在 `plugin.yaml`/`contracts/` 提供清单，由 **Plugin Manager** 代为注册。

### 2) 插件内嵌（可选）

PowerX 提供 `<AgentWidget />`（或 iframe 微前端）以便在插件页面内**嵌入同一套会话**。

---

## 🧪 本地调试建议

- **模式 A（推荐）**：PowerX + 插件一起跑，经 `/_p/:id/api/...` 链路调试（与上线一致）
- **模式 B**：插件单跑直连 `:8091/v1/...`，开发期用 `POWERX_DEV_MODE=1` 旁路注入 `tenant_id`
- **工具**：在 PowerX 增加一个仅开发可用的 `/_p/_dev/mount`，可热挂载/更新插件反代（可选）

### 本地管理端菜单如何布置？

模板工程默认提供以下页面占位：

| 菜单             | 推荐用途                         | 说明 |
| ---------------- | -------------------------------- | ---- |
| `Intro`          | 插件概览、接入指引               | 可保留，帮助运营快速了解插件价值 |
| `Templates`      | 业务主列表（CRUD 示例）          | 建议替换成插件的核心业务模块 |
| *(其它菜单)*   | 视业务所需自行扩展             | 不再内置团队协作、报表或系统设置示例 |

> **模板整理建议**：基础模板现在仅保留 `Intro` 与 `Templates` 页面，已移除 Team Collaboration、Reports、System Settings 等占位。
> 若你的插件需要额外的报表或配置页，请自行实现后再在 `manifest` 中声明菜单，避免出现空白入口。

实践建议：

1. **只暴露业务相关入口**：保留 `Intro` + 1~2 个核心业务菜单即可；其它页面在本地调试时可通过 Feature Flag 隐藏，避免上线出现空白导航。
2. **系统级设置交回宿主**：权限、团队等能力由 PowerX Settings 提供统一 UI，插件无需上架同类菜单。
3. **自定义配置集中在插件 Settings**：若插件确有特定配置（如外部 API Key、回调 URL），可在 `System Settings` 菜单中实现表单，其权限应绑定插件自有的 `resource`（例如 `base:settings:update`）。

开发时可通过 `manifest` 返回如下精简菜单示例：

```json
{
  "menus": [
    {
      "id": "plugins.base",
      "title": "menu.base.intro",
      "icon": "i-heroicons-sparkles",
      "path": "/plugins/base/intro"
    },
    {
      "id": "plugins.base.templates",
      "title": "menu.base.templates",
      "path": "/plugins/base/templates",
      "required_permissions": ["base:template:read"]
    }
  ]
}
```

这样在本地也能验证 RBAC 绑定逻辑，同时避免出现无内容的 Team/Permission 页面。

---

## 🛠 Makefile 常用目标

```makefile
build:      # 构建 backend/bin/plugin
run:        # 启动（POWERX_BIND_ADDR=:8091 POWERX_DB_SCHEMA=template）
migrate:    # 运行迁移
package:    # 打包 plugin.zip（plugin.yaml + backend/bin/... + web-admin）
docker:     # 构建镜像 powerx-plugin-base:<version>
```

---

## 🐳 Docker 运行

```bash
docker build -t powerx-plugin-base:0.1.0 .
docker run --rm -p 8091:8091 \
  -e POWERX_BIND_ADDR=":8091" \
  -e POWERX_DB_DSN="postgres://user:pwd@host:5432/powerx?sslmode=disable" \
  -e POWERX_DB_SCHEMA="px_com_powerx_plugins_base" \
  -e POWERX_CTX_JWKS_URL="http://powerx/_p/_internal/jwks" \
  -e POWERX_CTX_ISSUER="powerx-auth" \
  -e POWERX_CTX_AUDIENCE="powerx-plugin" \
  powerx-plugin-base:0.1.0
```

---

## ❓FAQ

**Q: 多租户密钥要不要“每租户一把”？**
A: 不需要。上下文签名是 **平台 ↔ 插件** 的信任，不是租户级。建议**平台级 JWT（RS256/ES256）**，或为**每插件派生 HMAC**。

**Q: `/_p/:id/api` 与插件的 `/v1` 会不会重复？**
A: 不会。反代仅拼接客户端传入的 `/*filepath`，因此前端必须带 `/v1/...`，插件内部路由也挂 `/v1`。

**Q: 必须上 RLS 吗？**
A: 强烈建议。即便应用层写错 where 条件，RLS 仍能兜底，避免跨租户读写。

---

## ✅ 待办清单

- [ ] 修正 `cmd/database/migtate/` → `migrate/` 命名
- [ ] 完成 RLS 迁移脚本并在 `migrate.go` 中执行
- [ ] 实现 `TenantContext`（HMAC/JWT）与 `BeginTenantTx`（`SET LOCAL`）
- [ ] 补齐 `/api/v1/admin/manifest` 与 `/api/v1/admin/rbac` 返回内容
- [ ] 注册 Agent Profile/Tools/Workflows（可由 `contracts/` 自动上报）
- [ ] 打通 PowerX 联调链路（反代 + 上下文注入）
- [ ] 关闭 `POWERX_DEV_MODE` 上线

---

如需，我可以把 README 中提到的 **`TenantContext` 中间件、`BeginTenantTx` 封装、RLS 迁移 SQL** 直接补到你的代码骨架里（保持你当前包结构）。
