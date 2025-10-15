# Implementation Plan: Plugin Runtime & Ops (Env, MCP, Observability, Quotas)

**Branch**: `[003-title-plugin-runtime-ops]` | **Date**: 2025-10-15 | **Spec**: [`specs/003-title-plugin-runtime-ops/spec.md`](specs/003-title-plugin-runtime-ops/spec.md)  
**Input**: Feature specification from `/specs/003-title-plugin-runtime-ops/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

---

## Summary

Define host-side runtime governance for PowerX plugins covering bootstrap sequencing, MCP session lifecycle, observability outputs, and quota/billing enforcement. Implementation will extend RuntimeManager services, MCP controller flows, platform observability/quotas, plus resource isolation, alerting, and audit facilities to ensure deterministic startup, authenticated capability sync, standardized metrics/logs/traces, and multi-tier resource controls with Marketplace reporting. All runtime ops must respect the PowerX installation layout under `plugins/installed/<plugin-id>/<version>/` (binaries, config/host-values.yaml, manifest assets).

---

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24 (backend services), Node 20 + Nuxt 4 with TypeScript 4.x (admin tooling references)  
**Primary Dependencies**: Fiber HTTP stack, gRPC runtime, OpenTelemetry SDK, Prometheus exporter, Loki/Tempo clients, MCP controller libraries, cgroup/iptables helpers  
**Storage**: PostgreSQL schema `powerx_plugin_base` (declared in `plugin.yaml` / `host_values.yaml`); tables for `runtime_assignments`, `port_reservations`, `mcp_sessions`, `runtime_audit_events`, `quota_ledger`, `marketplace_overages` (new migrations required)  
**Testing**: `make test` (Go unit/integration), targeted MCP session simulations, Prometheus scrape smoke tests, chaos tests for restart/backoff  
**Target Platform**: PowerX host Linux nodes running RuntimeManager + MCP controller  
**Project Type**: Backend services with supporting documentation (no new UI screens)  
**Performance Goals**: Heartbeat handling ≤50ms processing, startup READY within 90s, metrics scrape success ≥99.9%  
**Constraints**: P95 request latency <300ms under throttle, memory per plugin instance <= 512MB by default, restart attempts capped by exponential backoff  
**Scale/Scope**: Target up to 5k concurrent plugin instances, per-tenant quota window 5 minutes, hourly Marketplace summaries

### Platform / Hosting Integration (optional)

If this feature runs **under a host platform** (e.g., PowerX or similar), specify:

- **Reverse Proxy & Routes**: Business API under `/v1/**`; Admin endpoints (e.g., `/api/v1/admin/manifest`, `/api/v1/admin/rbac`) if applicable.  
- **Context Signing**: Inbound verification via **JWT (preferred)** or **HMAC**; context payload includes `tenant_id, user_id, permissions, request_id, exp, iat, iss, aud`.  
- **Tenant/RBAC**: Server-side authorization; DB session var for RLS (e.g., `SET LOCAL app.tenant_id = <tenant_id>`).  
- **Outbound Access**: Use **short-lived tokens** (e.g., STS) with explicit scopes to call the host.  
- **Observability**: `/healthz` endpoint; structured logs with `tenant_id`, `request_id`, `plugin_id`.  
- **Packaging**: SemVer; release artifacts checklist (binary, manifest, admin UI bundle, checksums, metadata).

> If there is **no host**, mark N/A and explain the external boundaries (clients, gateways, schedulers, etc.).

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*(If your `constitution.md` defines Gate IDs, reference them in braces for traceability.)*

- [x] **Host Contract First** {PX-HOST-001}  
  `/v1/runtime/**` handlers extended for assignments, health, and quota summaries; manifest/runtime fields (`manifest.runtime`, `manifest.quota`) remain authoritative.
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  REGISTER tokens validated via host-signed JWTs; new tables include `tenant_id` with RLS; tests cover scoped sessions and heartbeat rejection.
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  Transport handlers call `internal/services/runtime_ops` service; MCP controller reuses same service methods for lifecycle events.
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  Admin-only endpoints require `runtime.manage` scope; Marketplace notifications use scoped STS tokens; no broadened permissions.
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  Structured logs, Prometheus endpoints, Tempo spans, restart chaos tests, and CI coverage via `make test` and targeted simulations.
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  Only adds telemetry exporters already in template set; release notes capture new migrations and manifest fields; SemVer minor bump.

> Any unchecked item must be resolved or explicitly justified in **Complexity Tracking** below.

---

## Project Structure

### Documentation (this feature)

```

specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)

```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., internal/services/foo, internal/transport/http/bar).
  The delivered plan must not include Option labels.
-->

```
backend/
├── cmd/
│   ├── plugin/                 # RuntimeManager entrypoint (existing)
│   └── database/               # Migration tooling
├── internal/
│   ├── router/                 # Root router wiring transports
│   ├── transport/
│   │   ├── http/
│   │   │   ├── admin/
│   │   │   │   └── runtime_ops/     # New admin endpoints for runtime management
│   │   │   └── agent/
│   │   └── grpc/
│   ├── services/
│   │   ├── admin/
│   │   │   ├── templates/
│   │   │   └── runtime_ops/        # New service orchestrating bootstrap, MCP, quotas, isolation, alerts
│   │   └── agent/
│   ├── domain/models/
│   │   └── runtime_ops/            # Domain entities for runtime ops
│   ├── domain/repository/
│   │   └── runtime_ops/            # Persistence adapters for runtime ops tables
│   ├── shared/                 # Shared tooling (logging, utils)
│   └── mcp/controller/         # Session lifecycle coordination
├── etc/                        # Configuration incl. host_values.yaml injections & overrides
├── migrations/                 # SQL migrations for runtime ops (audit, assignments, quotas)
└── tests/
    ├── integration/
    ├── services/
    └── fixtures/

docs/
└── integration/03_runtime_and_ops/   # Operational guides & alerting thresholds
```

**Structure Decision**: Extend existing backend layers (`internal/router` orchestrator → `internal/transport/http/admin` for routes → `internal/services/admin`) with a `runtime_ops` subdomain covering bootstrap, isolation, sessions, observability, quotas, and audits, plus matching `domain/models`/`domain/repository` folders and shared migrations under `backend/migrations/`. Documentation continues under `docs/integration/03_runtime_and_ops/`; structure aligns with current admin/agent conventions and remains transport-agnostic.

---

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| [e.g., extra transport] | [current need] | [why single transport insufficient] |
| [e.g., custom repo pattern] | [specific problem] | [why direct DB access insufficient] |
