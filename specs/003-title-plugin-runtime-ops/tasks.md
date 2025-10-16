---
description: "Task list template for feature implementation (tool-agnostic)"
---

# Tasks: Plugin Runtime & Ops (Env, MCP, Observability, Quotas)

**Input**: Design docs from `/specs/003-title-plugin-runtime-ops/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Per constitution’s *Observable & Testable Delivery*, each story must include targeted service/integration tests respecting RLS and observability requirements.

**Organization**: Tasks grouped by user story, enabling independent delivery per scenario with shared setup and foundational gates.

> **Format**: `[ID] [P?] [Story] Description`  
>
> - **[P]**: Task can run in parallel (different files/concerns).  
> - **[Story]**: US1/US2/US3 referencing spec priorities.  
> - File paths follow plan structure (e.g., `backend/internal/services/admin/runtime_ops/...`).

## Path Conventions (adjust to your plan.md)

- **Single project**: `src/`, `tests/` at repo root  
- **Web app**: `backend/` (services, transports, db) and `frontend/` (Nuxt app)  
- **Mobile**: `api/`, `ios/` / `android/`  
- Always align paths here with the **“Structure Decision”** in `plan.md`.

<!--
============================================================================
IMPORTANT:
The /speckit.tasks command MUST replace all SAMPLE tasks below with the actual
tasks derived from:
  - User stories & priorities (spec.md)
  - Feature requirements (plan.md)
  - Entities (data-model.md)
  - Endpoints/contracts (contracts/)
Tasks MUST be grouped by story so each story is independently testable/MVP-able.
DO NOT keep these sample tasks in the generated tasks.md.
============================================================================
-->

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: Ensure local dev environment and baseline tooling align with constitution.

- [X] T001 Sync directories per plan (`backend/internal/services/admin/runtime_ops`, `backend/internal/domain/models/runtime_ops`, etc.)  
      `backend/internal/services/admin/runtime_ops/`, `backend/internal/domain/models/runtime_ops/`, `backend/internal/domain/repository/runtime_ops/`, `backend/internal/transport/http/admin/runtime_ops/`
- [X] T002 Install dependencies via `make dev-setup` (Go 1.24 modules, Node 20 setup)  
      `Makefile`, `make-files/dev.mk`
- [X] T003 [P] Verify tooling (`golangci-lint`, `go test`, Prometheus exporters) ready for runtime ops work  
      `make-files/dev.mk`, `make-files/test.mk`

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure gates required before implementing any story.

> ⚠️ Gates: Host Contract {PX-HOST-001}, Tenant Isolation {PX-CTX-001}, Service-Centric {PX-SVC-001}, Observability {PX-OBS-001}, Packaging {PX-PKG-001}.

- [X] T010 Create migration for runtime entities (`runtime_assignments`, `port_reservations`, `mcp_sessions`, `runtime_audit_events`, `quota_ledger`, `marketplace_overages`) with schema `powerx_plugin_base` and RLS policies  
      `backend/migrations/2025Q4_runtime_ops.sql`
- [X] T011 [P] Add runtime ops repositories scaffolding (interfaces + registrations) without business logic  
      `backend/internal/domain/repository/runtime_ops/runtime_ops_repository.go`, `backend/internal/domain/repository/runtime_ops/runtime_ops_repository_impl.go`
- [X] T012 [P] Extend service layer wiring (`services/admin/runtime_ops/service.go`) with constructor, dependencies, and placeholders  
      `backend/internal/services/admin/runtime_ops/service.go`
- [X] T013 Update admin router bootstrap to mount runtime ops handlers behind `runtime.manage` scope  
      `backend/internal/router/router.go`, `backend/internal/transport/http/admin/runtime_ops/routes.go`
- [X] T014 Configure observability exporters baseline (metrics registry, log fields) for runtime ops service  
      `backend/internal/services/admin/runtime_ops/observability.go`
- [X] T015 Ensure configuration surfaces defaults (heartbeat interval, quota window, restart backoff, log retention) while deferring DSN/ports to host-provided `host_values.yaml`  
      `backend/etc/config.yaml`, `backend/etc/config.example.yaml`, `backend/etc/host_values.yaml`, `backend/internal/config/config.go`
- [X] T016 Implement runtime isolation manager (CPU/memory/network limits) reading from manifest + `host_values.yaml`; expose structured violation events  
      `backend/internal/services/admin/runtime_ops/isolation.go`, `backend/internal/config/config.go`
- [X] T017 [P] Create audit repository scaffolding and DTOs for tenant-scoped runtime/MCP events  
      `backend/internal/domain/repository/runtime_ops/audit_repository.go`, `backend/internal/domain/models/runtime_ops/runtime_audit_event.go`
- [X] T018 Deliver debug CLI/scripts to simulate bootstrap & MCP heartbeat using sandbox tokens  
      `scripts/dev/runtime_ops_debug.sh`, `specs/003-title-plugin-runtime-ops/quickstart.md`

**Checkpoint**: Foundation ready → user stories can start (in parallel if staffed).

---

## Phase 3 — User Story 1 — Deterministic Runtime Bootstrap (Priority: P1) 🎯 MVP

**Goal**: Host launches plugin instances consistently across runtime modes with port registry, env injection, and restart policy enforcement.  
**Independent Test**: Use quickstart bootstrap script to start an exec-mode plugin, verify port reservation, env block, `/healthz` readiness, and restart backoff metrics.

### Tests for US1 (write-first, must FAIL initially)

- [X] T101 [P] [US1] Integration test covering bootstrap sequence & port registry lock  
      `backend/tests/integration/runtime_ops/test_bootstrap_port_registry_test.go`
- [X] T102 [US1] Service unit tests validating restart backoff counters & ledger updates  
      `backend/internal/services/admin/runtime_ops/service_bootstrap_test.go`
- [X] T103 [P] [US1] Database migration smoke test ensuring tables and RLS exist  
      `backend/tests/integration/runtime_ops/test_migrations_bootstrap_test.go`
- [X] T104 [P] [US1] Integration test verifying cgroup limits & filesystem sandbox during bootstrap  
      `backend/tests/integration/runtime_ops/test_resource_isolation_test.go`

### Implementation for US1

- [X] T105 [US1] Define domain models (`RuntimeAssignment`, `PortReservation`) and enums  
      `backend/internal/domain/models/runtime_ops/runtime_assignment.go`, `backend/internal/domain/models/runtime_ops/port_reservation.go`
- [X] T106 [US1] Implement repository methods for assignment creation, port locking, restart counters  
      `backend/internal/domain/repository/runtime_ops/runtime_assignment_repository.go`
- [X] T107 [US1] Implement service use-cases: unpack → allocate port → inject env → launch process → register health  
      `backend/internal/services/admin/runtime_ops/bootstrap_usecase.go`
- [X] T108 [US1] Create admin HTTP handlers for runtime launch & health status (thin, RBAC enforced)  
      `backend/internal/transport/http/admin/runtime_ops/handler_bootstrap.go`
- [ ] T109 [US1] Wire metrics/logging for bootstrap (port registry gauge, restart counter)  
      `backend/internal/services/admin/runtime_ops/observability_bootstrap.go`
- [X] T110 [US1] Document runtime mode defaults & health probes, ensuring they layer atop `host_values.yaml` without overriding host assignments  
      `backend/etc/config.example.yaml`, `docs/integration/03_runtime_and_ops/bootstrap.md`
- [X] T111 [US1] Integrate isolation manager application into bootstrap pipeline (CPU/memory/network caps, violation events)  
      `backend/internal/services/admin/runtime_ops/bootstrap_usecase.go`, `backend/internal/services/admin/runtime_ops/isolation.go`
- [X] T112 [US1] Enforce filesystem sandbox (mount whitelist, audit on violation) within runtime ops  
      `backend/internal/services/admin/runtime_ops/filesystem_sandbox.go`, `backend/internal/services/admin/runtime_ops/isolation.go`
- [X] T113 [US1] Documentation update for resource isolation expectations and sandbox directories  
      `docs/integration/03_runtime_and_ops/bootstrap.md`

**Checkpoint**: US1 is independently functional & testable (MVP).

---

## Phase 4 — User Story 2 — Resilient MCP Session Lifecycle (Priority: P1)

**Goal**: Ensure MCP sessions authenticate, exchange capabilities, and transition through lifecycle with heartbeat enforcement.  
**Independent Test**: Simulate REGISTER + CAPABILITY_SYNC, drop heartbeats, confirm session moves to STALE and reconnects respecting retry policy while logging metrics.

### Tests for US2

- [X] T201 [P] [US2] Integration test for REGISTER → READY → STALE lifecycle with JWT failures  
      `backend/tests/integration/runtime_ops/test_mcp_session_lifecycle_test.go`
- [X] T202 [US2] Unit test for heartbeat scheduler enforcing 15s interval & 3-miss timeout  
      `backend/internal/services/admin/runtime_ops/heartbeat_scheduler_test.go`
- [X] T203 [P] [US2] Contract test covering `/api/v1/admin/runtime/sessions` registration endpoint  
      `backend/tests/contract/runtime_ops/test_sessions_contract_test.go`
- [X] T204 [US2] Audit trail test ensuring REGISTER/CAPABILITY_SYNC/STALE events persist to runtime audit store  
      `backend/tests/integration/runtime_ops/test_mcp_audit_log_test.go`

### Implementation for US2

- [X] T205 [US2] Extend domain models with `MCPSession` states and JWT metadata  
      `backend/internal/domain/models/runtime_ops/mcp_session.go`
- [X] T206 [US2] Implement repository methods for session creation, heartbeat tracking, capability hash storage  
      `backend/internal/domain/repository/runtime_ops/mcp_session_repository.go`
- [X] T207 [US2] Implement service logic for REGISTER/ACK/CAPABILITY_SYNC, heartbeat monitor, reconnection  
      `backend/internal/services/admin/runtime_ops/mcp_session_usecase.go`
- [X] T208 [P] [US2] Add JWT validation & TLS enforcement helpers  
      `backend/internal/services/admin/runtime_ops/security.go`
- [X] T209 [US2] Implement `/api/v1/admin/runtime/sessions` handlers with thin transport  
      `backend/internal/transport/http/admin/runtime_ops/handler_sessions.go`
- [X] T210 [US2] Update metrics exposure (`powerx_mcp_sessions_total`, heartbeat latencies)  
      `backend/internal/services/admin/runtime_ops/observability_sessions.go`
- [X] T211 [US2] Documentation for MCP lifecycle and failure handling  
      `docs/integration/03_runtime_and_ops/mcp_sessions.md`
- [X] T212 [US2] Persist runtime audit events for REGISTER/ACK/CAPABILITY_SYNC/STALE transitions  
      `backend/internal/services/admin/runtime_ops/mcp_session_usecase.go`, `backend/internal/domain/repository/runtime_ops/audit_repository.go`

**Checkpoint**: US1 & US2 both independently pass.

---

## Phase 5 — User Story 3 — Unified Observability Surface (Priority: P1)

**Goal**: Standardize logging, metrics, and tracing outputs with Prometheus/Loki/Tempo integration per spec.  
**Independent Test**: Run observability smoke script to confirm `/metrics` exposes required series, logs carry trace IDs, and Tempo receives spans from bootstrap/session flows.

### Tests for US3

- [X] T301 [P] [US3] Integration test verifying `/metrics` scrape contains quota/session/restart metrics  
      `backend/tests/integration/runtime_ops/test_observability_metrics_test.go`
- [X] T302 [US3] Unit test for log formatter ensuring JSON字段（包含 `trace_id`, `tenant_id`）  
      `backend/internal/logger/runtime_test.go`
- [X] T303 [P] [US3] Integration test for OpenTelemetry span propagation from handler to service  
      `backend/tests/integration/runtime_ops/test_tracing_propagation_test.go`

### Implementation for US3

- [X] T304 [US3] Implement log formatter & retention逻辑（7 天归档）  
      `backend/internal/logger/runtime.go`
- [X] T305 [US3] Implement metrics exporter wiring for Prometheus `/metrics` (quota_usage, plugin_cost_total, restart counts)  
      `backend/internal/services/admin/runtime_ops/metrics_exporter.go`
- [X] T306 [US3] Implement OpenTelemetry tracing integration with span naming convention  
      `backend/internal/services/admin/runtime_ops/tracing.go`
- [X] T307 [US3] Update `/metrics` HTTP endpoint to merge runtime ops metrics  
      `backend/internal/transport/http/admin/runtime_ops/handler_metrics.go`
- [X] T308 [P] [US3] Update Loki/Tempo shipping configuration to read from host `POWERX_*` env / `host_values.yaml` fallbacks  
      `backend/internal/config/config.go`, `backend/etc/config.yaml`, `backend/etc/host_values.yaml`
- [X] T309 [US3] Documentation for observability stack usage, dashboards, alert thresholds  
      `docs/integration/03_runtime_and_ops/observability.md`
- [X] T310 [US3] Configure alert thresholds & host alert channel integration (Prometheus rules / notifier)  
      `backend/internal/services/admin/runtime_ops/observability_config.go`, `backend/etc/config.yaml`, `docs/integration/03_runtime_and_ops/Logs_Metrics_and_Tracing.md`

**Checkpoint**: All targeted stories are independently functional.

---

## Phase 6 — User Story 4 — Enforced Quota, Rate, and Cost Controls (Priority: P2)

**Goal**: Apply hierarchical quotas, rate limits, and cost aggregation with Marketplace hourly reporting.  
**Independent Test**: Drive quota breach via load script, confirm 429 or throttle event, `quota_usage` metric progression, and hourly Marketplace summary generation.

### Tests for US4

- [X] T401 [P] [US4] Integration test for quota breach handling + Marketplace summary enqueue  
      `backend/tests/integration/runtime_ops/test_quota_breach_test.go`
- [X] T402 [US4] Unit test for token bucket refill (5 min window) and cost aggregation  
      `backend/internal/services/admin/runtime_ops/quota_manager_test.go`

### Implementation for US4

- [X] T403 [US4] Extend domain models for `QuotaLedger` and `MarketplaceOverage` aggregation  
      `backend/internal/domain/models/runtime_ops/quota_ledger.go`, `backend/internal/domain/models/runtime_ops/marketplace_overage.go`
- [X] T404 [US4] Implement repository methods for quota ledger writes & hourly summaries  
      `backend/internal/domain/repository/runtime_ops/quota_repository.go`
- [X] T405 [US4] Implement rate limiter/token bucket logic in service layer  
      `backend/internal/services/admin/runtime_ops/quota_usecase.go`
- [X] T406 [US4] Integrate breach actions (`throttle`, `circuit_break`, `disable`, `notify_marketplace`)  
      `backend/internal/services/admin/runtime_ops/quota_actions.go`
- [X] T407 [P] [US4] Add admin endpoints or MCP events for quota status & overrides  
      `backend/internal/transport/http/admin/runtime_ops/quota_handler.go`, `backend/internal/mcp/controller/quota_events.go`
- [X] T408 [US4] Update metrics (`quota_usage`, `plugin_cost_total`) and Marketplace report scheduler  
      `backend/internal/services/admin/runtime_ops/observability_quota.go`
- [X] T409 [US4] Documentation for quota policies, breach responses, Marketplace notifications  
      `docs/integration/03_runtime_and_ops/Quotas_RateLimits_and_Costs.md`
- [X] T410 [US4] Persist quota breach audit events & hourly Marketplace summaries to audit store  
      `backend/internal/services/admin/runtime_ops/quota_usecase.go`, `backend/internal/domain/models/runtime_ops/runtime_audit_event.go`

**Checkpoint**: All user stories complete with quota enforcement.

---

## Phase 7 — Polish & Cross-Cutting

- [ ] T900 [P] Finalize docs (`quickstart.md`, `docs/integration/03_runtime_and_ops/*`) to match implementation  
      `specs/003-title-plugin-runtime-ops/quickstart.md`, `docs/integration/03_runtime_and_ops/`
- [ ] T901 Refactor shared utilities (reuse logging helpers, reduce duplication)  
      `backend/internal/shared/*`
- [ ] T902 Performance tuning: DB indexes for runtime ops tables, adjust Prometheus scrape config  
      `backend/internal/domain/repository/runtime_ops/migrations/*.sql`, `backend/etc/config.yaml`
- [ ] T903 [P] Additional regression/unit tests across services  
      `backend/internal/services/admin/runtime_ops/*_test.go`
- [ ] T904 Security review & RBAC audit for new endpoints/roles  
      `backend/internal/transport/http/admin/runtime_ops/routes.go`, `plugin/plugin.yaml`
- [ ] T905 Validate quickstart script & load generator align with latest config  
      `specs/003-title-plugin-runtime-ops/quickstart.md`, `scripts/dev/*`
- [ ] T906 Release prep (`make release`, manifest updates, checksum verification)  
      `make-files/build.mk`, `plugin/plugin.yaml`, `docs/releases/`

---

## Dependencies & Execution Order

### Phase Dependencies

- Setup → Foundational → US1 → US2 → US3 → US4 → Polish (respecting P1 before P2).  
- US1, US2, US3 all P1 stories; can run sequentially or in parallel post-foundation if resources allow, though bootstrap (US1) recommended first to unblock metrics/quota contexts.

### Within a Story

- Write tests first (fail), then implement: models/domain → repositories → services → transports → observability docs.  
- Maintain thin handlers; services remain source of business logic aligning with constitution.

### Parallel Opportunities

- `[P]` tasks across Setup/Foundational (different files).  
- Within stories: tests vs repository vs observability modules operate on distinct files (marked `[P]`).  
- US2 session work may parallel US3 observability once shared observability baseline complete; US4 depends on metrics/quota base from earlier phases.

---

## Notes

- `[P]` only when tasks touch distinct files or can run concurrently without race conditions.  
- Each story produces an independently verifiable increment (per spec’s acceptance tests).  
- Keep tasks concrete with real file paths to aid execution by LLM or engineers.  
- Follow quickstart for validation and update documentation promptly.
