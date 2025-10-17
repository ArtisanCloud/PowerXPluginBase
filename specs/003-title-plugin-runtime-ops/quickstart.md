# Quickstart: Plugin Runtime & Ops Validation

1. **Bootstrap environment**
   - Run `make dev-setup` to install Go modules and lint tools.
   - Copy `backend/etc/config.example.yaml` to `config.yaml`; merge with host-generated `host_values.yaml` so DSN、端口等仍由安装器注入；仅补充本地缺省值（如可观测性端点）。
   - Familiarize with the installed layout `plugins/installed/<plugin-id>/<version>/` (`backend/bin/{plugin,migrate}`, `config/host-values.yaml`, `web-admin/`), because runtime ops tooling reads and emits artefacts relative to this structure.

2. **Apply database migrations**
   - Execute `make migrate` to create `runtime_assignments`, `port_reservations`, `mcp_sessions`, `quota_ledger`, and `marketplace_overages` tables.
   - Verify RLS policies attach `tenant_id` filters on new tables.

3. **Start runtime manager in dev mode**
   - Run `make run` or `go run ./backend/cmd/plugin --config ./config.yaml`.
   - Confirm `/v1/runtime/healthz` responds `200 OK`.
   - Use `scripts/dev/runtime_ops_debug.sh inspect` to确认 `runtime_ops` 与 `monitoring.metrics.path=/api/v1/admin/runtime/metrics` 已生效。

4. **Launch sample plugin instance**
   - Use `scripts/dev/register_plugin.sh --plugin-id example.plugin --runtime exec`.
   - Observe logs for `unpack → port assign → env injection → process launch → health probe registration`.
   - Check `backend/logs/runtime-manager.log` includes assigned port and environment block.
   - Optional：执行 `scripts/dev/runtime_ops_debug.sh bootstrap` 以读取 `config.yaml` + `host-values.yaml` 并打印即将执行的启动参数。

5. **Validate MCP session lifecycle**
   - Trigger REGISTER via `scripts/dev/mcp_register.sh --plugin-id example.plugin --tenant demo`.
   - Simulate heartbeat loss with `scripts/dev/mcp_drop.sh` and confirm state transitions to `stale` within ~45s.

6. **Inspect observability outputs**
   - Scrape metrics: `scripts/dev/runtime_ops_debug.sh metrics` 或手动请求 `http://127.0.0.1:8086/api/v1/admin/runtime/metrics`，确认 `powerx_plugin_quota_usage`、`powerx_plugin_cost_total`、`powerx_mcp_sessions_total` 等系列存在。
   - Tail logs under `/var/lib/powerx/plugins/example.plugin/logs/`; confirm JSON schema with `trace_id`, `tenant_id`.
   - Use `scripts/dev/emit_trace.sh` to push a test span and verify it appears in Tempo with name `plugin.example.plugin.bootstrap`.

7. **Exercise quota breach flow**
   - Run load generator `scripts/dev/quota_burst.sh --tenant demo --capability bootstrap --qps 20 --duration 15`.
   - Confirm HTTP 429 or MCP throttle events after token bucket depletion（日志中可见 `quota_breach` 审计事件）。
   - Check Prometheus `powerx_plugin_quota_usage` 接近 1.0，`powerx_plugin_cost_total` 按租户累积；必要时再运行 `scripts/dev/runtime_ops_debug.sh metrics`。

8. **Run test suites**
   - `make test` for unit/integration coverage.
   - `go test ./backend/internal/services/admin/runtime_ops -run Test -count=1` for focused service tests.
   - Optional: `make test-coverage` to review `backend/coverage.html`.
