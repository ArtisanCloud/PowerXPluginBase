# PowerX 插件 Admin 前端

本目录存放 Nuxt 4 Admin UI。**所有页面与静态资源都必须挂载在插件专属前缀 `/_p/<plugin-id>/admin/` 下**，否则浏览器请求会落到宿主兜底返回的 HTML，从而触发白屏、404 或 "module script 的 MIME 是 text/html" 错误。

## 1. 构建期固定 baseURL（最关键）

- `nuxt.config.ts` 的 `app.baseURL` 会在构建期被写死到 HTML 中，运行期无法再覆盖。
- 只能在构建命令里通过环境变量注入：`POWERX_ADMIN_BASE="/_p/<plugin-id>/admin/"`，与 PowerX 的反向代理前缀保持一致。
- 本地直连调试时，直接访问 `http://127.0.0.1:3036/_p/<plugin-id>/admin/`，不用切换到根路径。

```bash
cd web-admin
POWERX_ADMIN_BASE="/_p/com.powerx.plugins.base/admin/" \
NODE_ENV=production \
npx nuxi build
```

构建完成后，务必抽样确认 HTML 中的 baseURL：

```bash
# 方式一：使用现成的 make 目标
make check-base

# 方式二：手动 curl（替换端口为 admin 进程监听端口）
curl http://127.0.0.1:3036/_p/com.powerx.plugins.base/admin/ | \
  grep 'app:{baseURL:"/_p/com.powerx.plugins.base/admin/"}'
```

## 2. PowerX 反向代理规则

- 页面路由：`/_p/<plugin-id>/admin/...` → 原样转发到插件前端进程，可选地在进入插件前先剥离 `/<locale>` 前缀再加回 baseURL 后缀。
- 静态资源：一律走 `/_p/<plugin-id>/admin/assets/...`，不要附带 locale 前缀，也不要兜底到其它插件。
- 遵循以上两条即可消除白屏、404 和 MIME 异常。

## 3. 启动插件 Admin 前端进程

PowerX 在插件生命周期中启动前端时需要：

1. 按 `plugin.yaml` 的 `frontend.admin` 描述执行 `entry + args`（不要写死路径或命令）。
2. 注入 `POWERX_ADMIN_BASE="/_p/<plugin-id>/admin/"`，保持与构建期一致。
3. 可选：同时注入 `POWERX_PROXY=1`，供前端判断自己运行在宿主内。

## 开发调试

```bash
npm install
npm run dev -- --port 3036 --host 0.0.0.0
```

开发模式同样会遵循 `app.baseURL`，因此路由、静态资源均会自动带上 `/_p/<plugin-id>/admin/` 前缀。

## 快速验收

1. **直连插件 Admin 端口**：页面源码需出现 `app:{baseURL:"/_p/com.powerx.plugins.base/admin/"}`。
2. **经 PowerX 访问**：浏览器 Network 面板中所有静态资源路径都应当是 `/_p/com.powerx.plugins.base/admin/assets/...` 且返回 200。

## Dev Console 故障排查面板

新版本引入 `/_p/com.powerx.plugins.base/admin/dev-console` 下的「故障排查」页签，前端约定如下：

- Job Run 表格通过 `/api/v1/admin/dev-console/jobs/runs` 分页加载；过滤条件会同步到 URL 查询参数，便于复制链接给其他值班人员。
- 安全操作表单调用 `/safe-ops/actions`。Dry Run 选项默认关闭，并在成功后自动将最新 Job Run 插入列表顶部。
- 故障排查仪表盘使用 `/troubleshooting/summary`，默认 5 分钟自动刷新，可通过 `useDevConsoleTroubleshootStore` 的 `scheduleAutoRefresh` 修改刷新节奏。
- Webhook 诊断面板依赖 `/webhooks/attempts` 与 `/webhooks/attempts/{id}`；若租户 ID 为空，前端直接提示填写。

调试建议：

```bash
# 启动 Nuxt dev server
npm run dev -- --port 3036

# 触发 mock 操作查看前端是否刷新
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"tenant_id":"demo-tenant","action":"replay","scope_type":"tenant","scope_ref":"demo-tenant"}' \
  "http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console/safe-ops/actions"
```

遇到空白页时优先查看浏览器控制台是否存在跨域或 403，通常意味着权限代码缺失或代理未带上 JWT。

## 常见误区

- 运行期再用 `POWERX_PROXY` 等变量修改 `app.baseURL` → 无效，HTML 中的路径已经写死。
- 把所有 `/assets/*` 都兜到某个插件 → 多插件场景会冲突，且会掩盖构建期 baseURL 错误。
- 手工调整 i18n 路由映射 → 只要 baseURL + strategy 正确，Nuxt 会自己处理语言前缀。

掌握以上三个要点，PowerX 与插件前端的路由、静态资源即可稳定打通。
