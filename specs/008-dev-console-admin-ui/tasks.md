# Tasks: Dev Console & Admin UI (Plugin Admin, Audit & History, Troubleshooting)

**Input**: Design docs from `/specs/008-dev-console-admin-ui/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Follow constitution gate {PX-OBS-001}; each story includes Go service/HTTP coverage and Nuxt tests aligned with acceptance criteria.

**Organization**: Tasks grouped by user story so each slice is independently implementable and testable.

> **Format**: `[ID] [P?] [Story] Description`

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: Ensure workspace, dependencies, and directory scaffolding align with the implementation plan.

- [x] T001 [P] Bootstrap backend/frontend deps for feature work (`make dev-setup`, `cd web-admin && npm install`) and verify Go 1.24 / Node 20 toolchains.  
- [x] T002 Establish admin console directory skeletons per plan (`backend/internal/domain/{models,repository}/admin_console`, `backend/internal/services/admin/console`, `backend/internal/transport/http/admin/console`, `web-admin/app/{pages,components,stores}/dev-console`).  
- [x] T003 [P] Configure lint/test scripts for new paths (update `Makefile`, `package.json` scripts, CI manifests) so added packages participate in `make test` / `npm run test`.

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Platform & data groundwork required before any story-specific development.

- [x] T010 Create migration `backend/migrations/2025Q4_admin_console.sql` defining `admin_console_audit_events`, `admin_console_config_changes`, `admin_console_job_runs`, indexes, and RLS policies.  
- [x] T011 Register new models in `backend/cmd/database/migrate/migrate.go` and add table constants in `backend/internal/domain/models/model.go`.  
- [x] T012 [P] Extend configuration defaults: add `AdminConsole` structs in `backend/internal/config/config.go`, update `backend/etc/config.example.yaml`, and document keys in `backend/etc/README.md`.  
- [x] T013 Wire manifest & RBAC: update `plugin.yaml`, `backend/internal/router/router.go`, and `backend/internal/transport/http/registry.go` with `/api/v1/admin/dev-console/**` routes plus permission codes (`operations.plugin.admin`, `.audit`, `.ops`).  
- [x] T014 [P] Introduce admin console observability baseline (`backend/internal/observability/admin_console/metrics.go`, register in `internal/shared/app/deps.go`) for audit export counters and dashboard freshness gauges.  
- [x] T015 Scaffold repository/service entry points (`backend/internal/domain/repository/admin_console/repository.go`, `backend/internal/services/admin/console/service.go`) with tenant-aware constructors, leaving logic TBD.  
- [x] T016 Add HTTP route placeholders (`backend/internal/transport/http/admin/console/routes.go`) and ensure middleware (JWT, tenant context) applied consistently.  
- [x] T017 [P] Update `quickstart.md` prerequisites and smoke commands to include new admin console endpoints and migrations.

**Checkpoint**: Database, config, routing, and observability scaffolds ready for story work.

---

## Phase 3 — User Story 1 — Configure & Govern Plugin Admin Console (Priority: P1) 🎯 MVP

**Goal**: Operators can open the console, edit configuration sections with validation, and audit changes.  
**Independent Test**: Navigate to `/_p/<plugin_id>/admin/dev-console`, update a configuration section, observe validation feedback, save a valid change, and verify the audit entry appears via API.

### Tests for US1

- [x] T101 [P] [US1] Go service unit tests for configuration validation and change history (`backend/internal/services/admin/console/config_service_test.go`).  
- [x] T102 [P] [US1] HTTP handler tests covering `/config/sections` list/update flows with RBAC and validation errors (`backend/internal/transport/http/admin/console/config_handler_test.go`).  
- [x] T103 [P] [US1] Nuxt component/store tests for form validation and audit banner (`web-admin/tests/dev-console/configure_console.spec.ts`).

### Implementation for US1

- [x] T104 [P] [US1] Implement `AdminConsoleConfigChange` model structs and gorm tags (`backend/internal/domain/models/admin_console/config_change.go`) plus repository CRUD (`.../repository/admin_console/config_change_repository.go`).  
- [x] T105 [P] [US1] Build configuration schema loader & validator utilities, including section metadata registry (`backend/internal/services/admin/console/config_schema.go`).  
- [x] T106 [US1] Implement config service use cases for list/update, validation, audit recording (`backend/internal/services/admin/console/config_service.go`).  
- [x] T107 [US1] Add audit event integration: write helpers to persist `AdminConsoleAuditEvent` entries and link to config changes (`backend/internal/services/admin/console/audit_logger.go`).  
- [x] T108 [US1] Implement HTTP handlers and DTOs for `/config/sections` endpoints with tenant scoping (`backend/internal/transport/http/admin/console/config_handler.go`, `dto/config.go`).  
- [x] T109 [US1] Register routes and RBAC guard checks in `routes.go`, ensure middleware enforces `operations.plugin.admin`.  
- [x] T110 [P] [US1] Create Pinia store/composable for configuration state & optimistic concurrency (`web-admin/app/stores/dev-console/config.ts`, `app/composables/useDevConsoleConfig.ts`).  
- [x] T111 [US1] Build Nuxt page & components for configuration dashboard (`web-admin/app/pages/_p/com.powerx.plugins.base/admin/dev-console/index.vue`, `components/dev-console/ConfigSectionCard.vue`).  
- [x] T112 [US1] Surface audit summary UI (last modified, actor) consuming new API (`components/dev-console/AuditSummaryBanner.vue`).  
- [x] T113 [P] [US1] Update `contracts/admin-dev-console-openapi.yaml` examples and ensure generated types (if any) align with implementation.
- [x] T114 [P] [US1] Implement `AdminConsoleAuditEvent` model & repository (`backend/internal/domain/models/admin_console/audit_event.go`, `.../repository/admin_console/audit_event_repository.go`) consistent with migration schema.
- [x] T115 [US1] Extend config/audit services to persist and query audit events via new repository (`backend/internal/services/admin/console/audit_logger.go`, `audit_service.go`).

**Checkpoint**: US1 routes accessible, validation working, audits recorded, UI functional.

---

## Phase 4 — User Story 2 — Inspect Audit & Activity History (Priority: P2)

**Goal**: Compliance reviewers can filter audit events and export CSV/JSON bundles.  
**Independent Test**: Filter events by actor/action, export both CSV and JSON for a date range, confirm files match on-screen entries, RBAC enforced.

### Tests for US2

- [x] T201 [P] [US2] Repository/service tests for audit filtering & pagination (`backend/internal/services/admin/console/audit_service_test.go`).  
- [x] T202 [P] [US2] HTTP export tests validating CSV vs JSON output formats and permissions (`backend/internal/transport/http/admin/console/audit_handler_test.go`).  
- [x] T203 [P] [US2] Nuxt e2e coverage for audit filters & export flow (`web-admin/tests/dev-console/audit_history.spec.ts`).

### Implementation for US2

- [x] T204 [P] [US2] Extend audit repository queries with actor/action filters and pagination (`backend/internal/domain/repository/admin_console/audit_repository.go`).  
- [x] T205 [US2] Implement audit service list/export logic including CSV serializer and download manifest (`backend/internal/services/admin/console/audit_service.go`).  
- [x] T206 [US2] Build HTTP handlers & DTOs for `/audit/events` and `/audit/export` (format negotiation, streaming) (`backend/internal/transport/http/admin/console/audit_handler.go`).  
- [x] T207 [US2] Hook metrics counters (export count, filter usage) in observability module.  
- [x] T208 [P] [US2] Create Nuxt audit history tab with filters, pagination, and permission guard messaging (`web-admin/app/components/dev-console/AuditHistoryTable.vue`).  
- [x] T209 [US2] Implement export drawer/modal enabling CSV/JSON choice and download progress (`web-admin/app/components/dev-console/AuditExportDialog.vue`).  
- [x] T210 [US2] Update Pinia store/composable to fetch audit events and trigger exports (`web-admin/app/stores/dev-console/audit.ts`).  
- [x] T211 [P] [US2] Refresh quickstart troubleshooting steps documenting export endpoints and permission requirements.

**Checkpoint**: US1 + US2 independently demonstrable; audit exports verified.

---

## Phase 5 — User Story 3 — Troubleshoot Jobs & Webhook Delivery (Priority: P3)

**Goal**: Support engineers access job/task history, initiate safe-ops actions, and view health/quota/webhook diagnostics within the console.  
**Independent Test**: View recent job runs, retry a failed job with scope selection, observe refreshed health/quota metrics auto-updating every 5 minutes, inspect webhook delivery attempt details with guidance.

### Tests for US3

- [x] T301 [P] [US3] Service tests for job run persistence, retry eligibility, and advisory lock collisions (`backend/internal/services/admin/console/job_service_test.go`).  
- [x] T302 [P] [US3] Service tests for troubleshooting summary aggregation across health/quota/webhook sources (`backend/internal/services/admin/console/troubleshoot_service_test.go`).  
- [x] T303 [P] [US3] Nuxt tests covering job retry flow and troubleshooting auto-refresh (`web-admin/tests/dev-console/troubleshoot_console.spec.ts`).

### Implementation for US3

- [x] T304 [P] [US3] Implement job run domain model & repository operations (`backend/internal/domain/models/admin_console/job_run.go`, `.../repository/admin_console/job_run_repository.go`).  
- [x] T305 [US3] Create safe-op execution orchestrator tying scope selection to runtime endpoints, enforcing advisory locks, and emitting audit/job run records (`backend/internal/services/admin/console/safe_ops_service.go`).  
- [x] T306 [US3] Add troubleshooting service aggregating health/quota/webhook diagnostics with cache + auto-refresh metadata (`backend/internal/services/admin/console/troubleshoot_service.go`).  
- [x] T307 [US3] Implement HTTP handlers for `/jobs/runs`, `/jobs/runs/{id}/retry`, `/safe-ops/actions`, `/troubleshooting/summary`, `/webhooks/attempts`, `/webhooks/attempts/{id}` (`backend/internal/transport/http/admin/console/troubleshoot_handler.go`, `job_handler.go`).  
- [x] T308 [US3] Integrate webhook attempt repository queries leveraging `integration_webhook_attempts` and prepare data projections for UI drill-downs.  
- [x] T309 [US3] Extend observability metrics for safe-op executions, retry successes/failures, and troubleshooting refresh lag.  
- [x] T310 [P] [US3] Build Nuxt job history table with retry actions and scope selector (`web-admin/app/components/dev-console/JobRunsTable.vue`).  
- [x] T311 [P] [US3] Create troubleshooting dashboard components for health/quota/webhook sections with manual refresh control (`web-admin/app/components/dev-console/TroubleshootingDashboard.vue`).  
- [x] T312 [US3] Implement Pinia store/composables for job runs, safe-ops actions, and troubleshooting data (`web-admin/app/stores/dev-console/troubleshoot.ts`, `app/composables/useSafeOps.ts`).  
- [x] T313 [US3] Wire UI navigation tabs and routing for audit/troubleshoot sections (`web-admin/app/pages/_p/com.powerx.plugins.base/admin/dev-console/index.vue` layout updates).  
- [x] T314 [P] [US3] Update contracts to reflect final payloads and retry eligibility documentation (`contracts/admin-dev-console-openapi.yaml`).  
- [x] T315 [US3] Implement contextual help/runbook service responses for troubleshooting sections (`backend/internal/services/admin/console/help_service.go`, `help_content.go`).  
- [x] T316 [P] [US3] Render troubleshooting help panels and runbook links in Nuxt UI (`web-admin/app/components/dev-console/TroubleshootingHelpPanel.vue`, page layout wiring).  
- [x] T317 [US3] Add migration/index updates for webhook drill-down performance (e.g., `integration_webhook_attempts` composite index) and document in migration file.

**Checkpoint**: Full console delivers configuration, audit, and troubleshooting capabilities; safe-ops actions audited.

---

## Phase 6 — Polish & Cross-Cutting

- [x] T900 [P] Refresh product documentation & runbooks (`docs/observability`, `docs/support`, `web-admin/README.md`) to reflect new console surfaces.  
- [x] T901 Harden security: double-check RBAC enforcement, add threat-model notes in `backend/docs/security/admin_console.md`.  
- [x] T902 [P] Performance review: benchmark audit export for 1,500 events, add index tweaks or streaming adjustments as needed.  
- [x] T903 Validate `quickstart.md` commands, ensure smoke scripts succeed end-to-end.  
- [x] T904 Prep release artifacts (`plugin.yaml` version bump, `make release`, changelog entry).

---

## Dependencies & Execution Order

- **Phase Flow**: Setup → Foundational → US1 → US2 → US3 → Polish.  
- Foundational tasks complete before any user story work.  
- Stories may proceed in priority order; US2/US3 can begin once shared repositories/services introduced in US1 are merged.

### Parallel Opportunities

- Setup tasks T001 & T003 can run alongside directory scaffolding (T002).  
- Foundational tasks T012–T016 are parallel-friendly once migration stub (T010) is authored.  
- Within each story, tests (T101/T102/T103, etc.) can be authored concurrently; backend and frontend implementations marked `[P]` touch distinct files.  
- US2 and US3 frontend work (T208/T209 vs. T310/T311) can proceed in parallel once APIs stabilized.

### Implementation Strategy

- **MVP**: Complete US1 end-to-end (config forms + audits).  
- **Next**: Layer US2 audit history & export.  
- **Final**: Deliver troubleshooting dashboards & safe-ops (US3) followed by polish tasks.

---
