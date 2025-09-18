# 插件安装与部署指南

本文档描述如何将 `com.powerx.plugins.base` 插件构建、打包并安装到 PowerX 宿主环境，同时兼顾本地调试与运维交付需求。

## 先决条件

- 已安装 Go 1.22+、Node.js 20+、npm。
- 确保 `make` 可用（macOS 可通过 Xcode Command Line Tools，Linux 可通过对应包管理器安装）。
- 若需访问宿主配置，请准备 PowerX 平台提供的数据库、Redis、消息系统等连接信息。

## 目录约定回顾

- `backend/`: Go 服务代码及依赖资源，产物为 `backend/bin/plugin`。
- `web-admin/`: Nuxt 4 管理端前端，产物位于 `web-admin/.output/`。
- `config/schema.yaml`: 插件声明宿主需要的配置字段。
- `config/values.example.yaml`: 本地调试示例值，可复制为起始模板。
- `backend/etc/config.yaml`: 开发环境默认配置（**不会**随打包一起下发）。

## 构建步骤

1. 安装依赖
   ```bash
   npm --prefix web-admin install
   ```
2. 构建后端与前端
   ```bash
   make dist VERSION=<版本号>
   ```
   该命令会执行：
   - `go build` → 输出 `dist/<版本>/backend/bin/plugin`
   - `nuxt build` → 输出 `dist/<版本>/web-admin/.output`
   - 复制必要的 manifest、schema、README 等文件

## 产物结构

`dist/<版本>/` 目录示例：

```
dist/0.1.0/
├── backend/
│   └── bin/plugin         # Go 二进制
├── web-admin/
│   └── .output/           # Nuxt 构建产物（public + server）
├── config/
│   ├── schema.yaml        # 宿主读取字段描述
│   └── values.example.yaml
├── plugin.yaml            # 插件 Manifest
└── README.md              # 项目说明
```

## 配置注入流程

1. **安装阶段（宿主执行）**
   - PowerX 安装器读取 `config/schema.yaml` → 生成 `host-values.yaml`（实际值保存于宿主控制的目录）。
   - 将生成结果写入 `plugins/installed/<id>/<version>/config/host-values.yaml`。
   - `plugin.yaml` 中的 `runtime.env.CONFIG_PATH` 会指向该目录。

2. **本地调试**
   - 复制示例配置：
     ```bash
     cp config/values.example.yaml backend/etc/config.yaml
     ```
   - 按需修改数据库、PowerX gRPC 等字段。
   - 运行：
     ```bash
     CONFIG_PATH=backend/etc ./backend/bin/plugin
     ```

> ⚠️ 打包目录默认不包含 `backend/etc/config.yaml`，以防敏感信息被发布；发布时请仅提供脱敏模板或按宿主流程注入。

## 前端部署与预览

- **宿主模式**：PowerX 会直接将 `web-admin/.output/public` 挂载到 `/_p/<plugin-id>/admin/`，无需启动 Node 进程。
- **本地预览**：
  ```bash
  PORT=3036 node dist/<版本>/web-admin/.output/server/index.mjs
  ```
  - 当 `NODE_ENV=production`（默认）时，访问路径形如 `http://localhost:3036/_p/com.powerx.plugins.base/admin/`。
  - 若需纯根路径预览，可暂时设 `NODE_ENV=development`。

## 安装到 PowerX 宿主

1. 将 `dist/<版本>/` 整体复制到宿主期望的安装临时目录。
2. 以 `powerx-cli` 或其他官方工具执行插件安装（示例命令，具体以宿主文档为准）：
   ```bash
   powerx plugin install --from ./dist/<版本>
   ```
3. 宿主会：
   - 验证 `plugin.yaml` 与 `config/schema.yaml`。
   - 生成并写入 `host-values.yaml`。
   - 将后端二进制登记到 Supervisor，拷贝前端静态资源。
4. 启用插件时，宿主 Supervisor 会按 `runtime.entry` 启动进程，并注入包括 `CONFIG_PATH` 在内的环境变量。

## 常见问题

- **WARN: failed to parse YAML config**
  - 检查 `host-values.yaml` 或 `backend/etc/config.yaml` 的结构是否与 `config.Config` 字段匹配，尤其是 `grpc_upstream` 需为对象映射。
- **前端打包报 Tailwind Unknown Utility**
  - 确保在带 `scoped` 的样式中使用 `@reference "@/assets/css/main.css";` 引入 Tailwind 基础层。
- **端口冲突**
  - 生产场景由 PowerX 自动在运行时分配端口；本地运行可通过 `PX_BIND_ADDR` 或配置文件修改。

## 进一步操作

- `make package VERSION=<版本号>` 可将 `dist/<版本>` 打包为 zip，方便远程分发。
- 若需要生成离线文档或脚本，请在 `docs/` 目录补充对应说明。

如需更多细节可参考 `docs/plugins/readme.md` 与 `backend/internal/config/config.go` 中的加载逻辑。
