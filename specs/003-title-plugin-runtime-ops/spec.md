# Feature Specification: Plugin Runtime & Ops (Env, MCP, Observability, Quotas)

**Feature Branch**: `[003-title-plugin-runtime-ops]`  
**Created**: 2025-10-15  
**Status**: Draft  
**Input**: User description: "Title: Plugin Runtime & Ops (Env, MCP, Observability, Quotas). WHAT/WHY: 规范 PowerX 插件在宿主系统中的运行时行为，包括环境注入、MCP 会话注册、日志与可观测性标准、配额与计费机制，确保插件具备稳定、安全、可监控的运行基础，并支持多语言多租户的统一运维体系。 Scope: 运行环境与端口管理；MCP 会话与注册机制；日志、指标与链路追踪；配额、限速与计费。 Out-of-Scope: 插件能力契约与 Schema；安全与隐私；定价策略与 Marketplace 分润。 Dependencies/Assumptions: 宿主具备 RuntimeManager 与 MCP 控制器；插件 .pxp 包含 manifest/runtime/env 配置；PowerX Observability Stack（Loki / Prometheus / Tempo）已部署；QuotaManager 与 BillingEngine 启用；所有通信需带 POWERX_AUTH_TOKEN 且使用 TLS；插件应在启动后完成 REGISTER + CAPABILITY_SYNC 方可进入 READY 状态；监控指标通过 /metrics 或 MCP 推送汇总。"

## Clarifications

### Session 2025-10-15
- Q: What MCP heartbeat cadence should we standardize before marking a session STALE? → A: Ping every 15s; STALE after 3 missed exchanges (≈45s).
- Q: 租户级配额的令牌桶重置窗口默认设定为多久？ → A: 每 5 分钟刷新窗口。
- Q: 插件健康探针连续失败时，默认重启回退策略设为？ → A: 指数回退：起始 5 秒，乘 2 上限 2 分钟，最多 5 次。
- Q: 日志保留多久后才自动转储或清理？ → A: 本地保留 7 天，随后打包归档。
- Q: 配额或速率超限时，Marketplace 通知节奏默认是什么？ → A: 每小时汇总超限报告给 Marketplace。
- Port governance will be centrally managed by the RuntimeManager via a Port Registry that persists reservations and rejects duplicate allocations; plugins may not hard-code ports.
- MCP session security layers will rely on JWT assertions signed by the host; plugins must verify signature, tenant scope, and expiry before accepting REGISTER completion.
- Observability exports require Prometheus `/metrics` and OpenTelemetry spans; if the runtime cannot host an HTTP endpoint, the plugin must surface metrics via MCP `EVENT` messages in the agreed format.
- Host-provided `host_values.yaml` injects DSN, schema (`px_*`), ports, and security tokens; runtime ops MUST consume these values without overriding host-managed assignments.
- Installed plugins live under `plugins/installed/<plugin-id>/<version>/` with host-provisioned binaries (`backend/bin/plugin`, `backend/bin/migrate`), configuration (`config/host-values.yaml`), and manifest assets; runtime ops MUST coordinate with this layout when reading configs or emitting artefacts.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Deterministic Runtime Bootstrap (Priority: P1)

Runtime operators need a deterministic startup pipeline that unpacks the plugin bundle, allocates resources, injects environment variables, and registers health probes so every plugin instance launches safely across exec, docker, http-proxy, php-fpm, or remote modes.

**Why this priority**: Without a predictable bootstrap sequence, plugins fail inconsistently, consume conflicting ports, or bypass tenancy isolation, creating outages for all tenants.

**Independent Test**: Spin up any supported runtime mode in a staging host, run the documented bootstrap sequence, and verify the plugin reaches `READY` with health probes passing and allocated resources recorded in the Port Registry and runtime ledger.

**Acceptance Scenarios**:

1. **Given** a plugin package with a declared runtime mode, **When** the host initiates deployment, **Then** the RuntimeManager executes `unpack → port assign → env injection → process launch → health probe registration` and records the assignment.
2. **Given** an existing plugin instance holding port 6100, **When** a second instance of the same plugin requests startup, **Then** the Port Registry rejects the duplicate allocation and the host emits a structured error without launching the process.
3. **Given** a plugin crash triggered by failing `/healthz`, **When** the restart policy attempts recovery, **Then** the RuntimeManager replays the bootstrap flow, revalidates environment variables, and increments restart counters exposed via metrics.

---

### User Story 2 - Resilient MCP Session Lifecycle (Priority: P1)

MCP controllers and plugin runtimes require a handshake that authenticates, registers capabilities, and maintains heartbeats so commands flow reliably and tenants are not exposed to stale sessions.

**Why this priority**: Without strict MCP lifecycle handling, plugins may process commands without valid auth, lose capability sync, or remain in zombie sessions that misroute tenant data.

**Independent Test**: Establish a session with REGISTER/ACK/CAPABILITY_SYNC, force a network interruption, and confirm the plugin transitions through CONNECTING → REGISTERED → READY → STALE → CLOSED while honoring heartbeat backoff and replay rules.

**Acceptance Scenarios**:

1. **Given** a plugin reaching REGISTERED, **When** the host sends CAPABILITY_SYNC with the latest manifest capabilities, **Then** the plugin acknowledges with AGENT_TOOL_SYNC and exposes the updated tool list within one heartbeat interval.
2. **Given** a 2× heartbeat timeout, **When** the plugin misses PING/PONG exchanges, **Then** the session enters STALE, the RuntimeManager triggers reconnection up to the retry budget, and metrics record `powerx_mcp_sessions_total{status="stale"}`.
3. **Given** a JWT presented during REGISTER that fails signature validation or tenant scope, **When** the plugin processes the message, **Then** it MUST reject the session, log a security warning, and refuse to move to READY.

---

### User Story 3 - Unified Observability Surface (Priority: P1)

Site reliability engineers need standardized logging, metrics, and tracing outputs so they can correlate plugin behavior with host infrastructure using Grafana/Loki/Tempo dashboards.

**Why this priority**: Inconsistent observability data blocks alerting, hides tenant-specific incidents, and prevents forensic analysis across distributed components.

**Independent Test**: Deploy a plugin configured with the observability standard, generate synthetic load, and verify logs arrive in Loki with trace-correlated IDs, Prometheus scrapes `/metrics`, and Tempo captures spans with the prescribed naming.

**Acceptance Scenarios**:

1. **Given** a plugin handling tenant requests, **When** it emits logs, **Then** entries follow the JSON schema with `level`, `timestamp`, `trace_id`, `tenant_id`, and `message`, stored under `/var/lib/powerx/plugins/<id>/logs/` and forwarded to Loki.
2. **Given** a tracing-enabled host request with W3C traceparent headers, **When** the plugin generates spans, **Then** it propagates the context, names spans `plugin.<id>.<operation>`, and exports them through OpenTelemetry exporters or MCP events.
3. **Given** the metrics endpoint `/metrics`, **When** Prometheus scrapes it, **Then** gauges for `powerx_plugin_cpu_seconds_total`, `powerx_plugin_memory_bytes`, `powerx_plugin_health_status`, and custom business counters are exposed without scrape errors.

---

### User Story 4 - Enforced Quota, Rate, and Cost Controls (Priority: P2)

Quota managers and billing analysts need multi-tier quotas and usage accounting so plugins cannot exceed tenant allocations and cost reporting remains accurate.

**Why this priority**: Overages lead to noisy-neighbor incidents, billing disputes, and compliance risks across tenants sharing infrastructure.

**Independent Test**: Configure quota thresholds for core, plugin, and tenant levels, execute workload surpassing each threshold, and observe throttling responses, Prometheus counters, and billing ledgers updating with CPU-second and bandwidth charges.

**Acceptance Scenarios**:

1. **Given** a plugin endpoint exceeding its token-bucket QPS, **When** further INVOKE requests arrive, **Then** the host or plugin responds with HTTP 429 (or MCP throttle event), logs the breach, and emits `quota_usage{scope="plugin"}` metrics approaching 1.0.
2. **Given** tenant-level CPU seconds exceeding the quota window, **When** BillingEngine aggregates usage, **Then** it marks the tenant over-limit, raises an alert, and optionally disables non-essential capabilities per manifest policy.
3. **Given** a marketplace cost plan tied to manifest `costs`, **When** monthly aggregation runs, **Then** `plugin_cost_total{tenant_id}` reflects CPU, bandwidth, and invocation tallies aligned with the pricing table and ready for invoicing.

---

### Edge Cases

- Port allocation exhaustion occurs because of leaked reservations after repeated crash loops.
- Remote runtime delegates are unreachable yet sessions continue to queue requests, risking backlog overflow.
- Plugins deployed in restricted sandboxes cannot open `/metrics`, requiring MCP metric forwarding fallback.
- Tenants sharing the same plugin instance have divergent quota tiers, demanding per-tenant enforcement inside pooled processes.
- Trace sampling is disabled in the host, but critical errors still require forced tracing to capture minimal spans.
- BillingEngine receives delayed usage batches from offline plugins, challenging cutoff windows and proration.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The runtime standard MUST support `exec`, `docker`, `http-proxy`, `php-fpm`, and `remote` execution modes with configuration flags in `manifest.runtime`.
- **FR-002**: The RuntimeManager MUST execute the bootstrap pipeline `unpack → port reservation → environment injection → process launch → health registration` for every deployment and persist the state in a runtime ledger.
- **FR-003**: Environment variables MUST include `POWERX_PLUGIN_ID`, allocated `PORT`, `TENANT_ID` (per session or worker), and `POWERX_AUTH_TOKEN`; missing entries MUST abort startup.
- **FR-004**: A Port Registry MUST prevent port collisions by locking allocations per plugin instance and releasing entries only after confirmed shutdown.
- **FR-005**: Health probes MUST expose `/healthz` with configurable interval, success thresholds, and an exponential restart backoff (start at 5s, double to a 2-minute cap, maximum 5 attempts); failing probes MUST trigger host-managed restarts respecting these limits.
- **FR-006**: Runtime isolation MUST enforce CPU, memory, network, and filesystem limits aligned with manifest declarations; violations MUST surface as structured events.
- **FR-007**: File system access MUST be confined to sandboxed directories including `/var/lib/powerx/plugins/<id>/` and explicitly declared shared mounts.
- **FR-008**: Logs MUST be emitted in JSON with fields `timestamp`, `level`, `plugin_id`, `tenant_id`, `trace_id`, `message`, and optional `error.stack`, stored locally and shipped to Loki.
- **FR-008a**: Local log retention MUST keep `/var/lib/powerx/plugins/<id>/logs/` for 7 days before rotating into compressed archives handed to centralized storage.
- **FR-009**: Metrics MUST be exposed via `/metrics` (Prometheus plaintext) or via MCP `EVENT` payloads when HTTP exposure is impossible, maintaining identical metric names and labels.
- **FR-010**: Trace context MUST propagate using W3C Trace Context headers; spans MUST be named `plugin.<id>.<domain>.<operation>` and include tenant and request identifiers.
- **FR-011**: MCP protocol handling MUST implement REGISTER, ACK, CAPABILITY_SYNC, INVOKE, EVENT, and HEARTBEAT semantics with lifecycle states CONNECTING → REGISTERED → READY → STALE → CLOSED.
- **FR-012**: Session authentication MUST validate host-issued JWTs (signature, expiry, audience, tenant scope) and enforce TLS for every MCP channel.
- **FR-013**: Heartbeats MUST send PING every 15 seconds (±1s jitter) and declare a session STALE after three consecutive missed PONGs (~45s), triggering reconnection or shutdown per policy.
- **FR-014**: Capability synchronization MUST reconcile `manifest.capabilities` with agent tooling via AGENT_TOOL_SYNC events before entering READY.
- **FR-015**: Observability MUST integrate with Grafana/Loki/Tempo dashboards using standardized labels (`plugin_id`, `tenant_id`, `runtime_mode`, `host_id`).
- **FR-016**: Quota enforcement MUST support hierarchical scopes (core, plugin, tenant) with thresholds defined in `manifest.quota` and `manifest.rate_limits`.
- **FR-017**: Rate limiting MUST implement token-bucket or leaky-bucket algorithms with parameters per capability (`qps`, `burst`, `concurrency`) and a default tenant-scope refill window of 5 minutes, returning HTTP 429 or MCP throttle events on exceedance.
- **FR-018**: Cost accounting MUST calculate CPU seconds, bandwidth MB, and invocation counts per tenant, persisting records for BillingEngine aggregation.
- **FR-019**: Over-limit handling MUST support actions `throttle`, `circuit_break`, `disable`, and `notify_marketplace`, configurable in manifest policies.
- **FR-019a**: Marketplace notifications MUST aggregate over-limit events into hourly summaries by plugin and tenant before dispatch.
- **FR-020**: Prometheus metrics MUST include `powerx_mcp_sessions_total`, `quota_usage`, `plugin_cost_total`, `plugin_restart_total`, and latency histograms with P95/P99 buckets.
- **FR-021**: Alerting MUST define thresholds for health check failure rate, P95 latency, error rate, quota exhaustion, and billing anomaly, wiring results to the host alerting stack (Alertmanager/marketplace notifier) using values sourced from `host_values.yaml` and manifest overrides.
- **FR-022**: All runtime and MCP events MUST include tenant-aware audit trails stored for the retention period defined by platform policy, persisting entries to host-audited tables or log streams with immutable IDs.
- **FR-023**: Debug tooling MUST allow developers to simulate MCP sessions and runtime bootstrap locally using sandbox tokens without bypassing security validation, reusing the same CLI flow that reads `config/host-values.yaml`.
- **FR-024**: Runtime ops MUST treat host-generated `host_values.yaml` and environment variables (`POWERX_*`) as source-of-truth for connection info, schemas, and bind ports, augmenting but never conflicting with host startup orchestration.

### Key Entities *(include if feature involves data)*

- **RuntimeAssignment**: Records runtime mode, process identifiers, allocated ports, resource limits, and health status for each plugin instance.
- **PortReservation**: Persistent mapping of plugin instance IDs to port ranges with timestamps, owner host, and lifecycle events.
- **MCPSession**: Represents the handshake, JWT claims, capability list, heartbeat stats, and state transitions for each active session.
- **ObservabilityEnvelope**: Normalized structure for logs, metrics, and traces including plugin, tenant, trace IDs, and correlation metadata.
- **QuotaLedger**: Aggregated usage record per scope (core, plugin, tenant) capturing tokens consumed, CPU seconds, bandwidth, and billing status.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 99% of plugin startups across supported runtime modes reach `READY` within 90 seconds without manual intervention over a 30-day window.
- **SC-002**: 100% of MCP sessions authenticate JWTs successfully, with zero unauthorized REGISTER events recorded in `powerx_mcp_sessions_total{status="rejected"}` per month.
- **SC-003**: 95% of logs ingested into Loki contain valid `trace_id` and `tenant_id` fields, enabling end-to-end correlation in Grafana dashboards.
- **SC-004**: Prometheus scrape error rate for plugin `/metrics` endpoints remains below 0.1% per day across all tenants.
- **SC-005**: 100% of quota breaches emit alerts within 1 minute and enforce the configured mitigation (throttle/circuit break) during targeted chaos tests.
- **SC-006**: Monthly billing reconciliation variance between QuotaLedger and BillingEngine invoices stays under 1% of total plugin revenue.

## Assumptions

- RuntimeManager can enforce container, cgroup, and filesystem constraints in the hosting infrastructure.
- MCP transport supports TLS termination and exposes certificate rotation hooks without additional feature work.
- Observability stack (Prometheus, Loki, Tempo) is reachable from plugin runtimes with required credentials pre-provisioned.
- QuotaManager and BillingEngine expose APIs to ingest usage snapshots and return enforcement decisions synchronously.
- Plugins can be patched to embed the PowerX logging, metrics, and tracing SDKs or use sidecar shims where direct integration is not possible.
