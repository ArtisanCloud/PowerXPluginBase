# Tasks: Security & Compliance (Privacy, ToolGrant, Baseline, Vulnerability Response)

**Input**: Design docs from `/specs/004-security-compliance/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: Follow repository constitution (*Observable & Testable Delivery*) when implementing each task; add targeted tests alongside code changes where appropriate.

**Organization**: Tasks are grouped by **user story** so each story can be implemented, inspected, and delivered independently.

> **Format**: `[ID] [P?] [Story] Description`  
>
> - **[P]**: Can run in parallel (different files, no dependencies)  
> - **[Story]**: Which story this task belongs to (e.g., US1, US2)  
> - Include **exact file paths** in descriptions (align with `plan.md` structure)

## Path Conventions

- Backend services live under `backend/internal/**` per plan structure (`domain/models`, `domain/repository`, `services`, `middleware`, `transport`).
- Frontend admin UI updates live under `web-admin/app/**`.
- Build tooling and Make targets reside in `make-files/**` and `Makefile`.
- Security artifacts and docs sit in `build/security/`, `dist/security/`, and `backend/etc/**`.

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: Prepare shared configuration and workspace scaffolding required by all stories.

- [X] T001 [Setup] Create baseline config stub `backend/etc/security_baseline.yaml` with sections for `masking_rules`, `baseline_version`, `tool_grant`, and `consent_defaults`.
- [X] T002 [Setup] Document new security configuration keys in `backend/etc/README.md`, referencing default retention (90 days) and audit pipeline outputs.
- [X] T003 [Setup] Add security artifact directories (`build/security/.gitkeep`, `dist/security/.gitkeep`) and update `.gitignore` to retain generated reports.

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Core plumbing that must exist before user stories start.

- [X] T010 [Foundational] Add empty module scaffolds and package docs for `backend/internal/domain/models/{privacy,security,tool_grant}` and matching `domain/repository`/`services` directories.
- [X] T011 [Foundational] Extend dependency injection in `backend/internal/shared/app/deps.go` to register new security-related repositories/services via placeholders.
- [X] T012 [Foundational] Update router grouping in `backend/internal/transport/http/routes.go` and `backend/internal/transport/http/admin/routes.go` to mount `/admin/security` and `/agent/security` namespaces (handlers wired later).
- [X] T013 [Foundational] Expose configuration handles in `backend/internal/config/config.go` for security baseline, consent defaults, and tool grant TTL (read from `security_baseline.yaml`).
- [X] T014 [Foundational] Seed logging hook in `backend/internal/logger/runtime.go` to accept privacy masking adapters (actual masking logic lands with US1).
- [X] T015 [Foundational] Configure audit log retention & export tooling: add rotation policy (365 天在线) and CLI/export script under `scripts/security/audit_export.sh` documented in `docs/security/audit-logs.md`.

**Checkpoint**: Foundational plumbing ready — user stories can progress independently.

---

## Phase 3 — User Story 1 — Enforce Tenant-Isolated Data Privacy Controls (Priority: P1) 🎯 MVP

**Goal**: Enforce consent-scoped access, retention workflows, and audit evidence so tenant data stays compliant with GDPR/PIPL across its lifecycle.

**Independent Test**: Trigger consent issuance, data erasure, and export events; confirm only approved fields are read, logs are masked, and lifecycle evidence is persisted plus mirrored to `/logs/audit.log`.

### Implementation for US1

- [X] T100 [P] [US1] Implement GORM models in `backend/internal/domain/models/privacy/{classification.go,consent_token.go,lifecycle_event.go}` with enums and JSONB fields per `data-model.md`.
- [X] T101 [US1] Create tenant-scoped migrations in `backend/internal/db/migrations/` for privacy tables and enable RLS policies aligning with `tenant_id`.
- [X] T102 [P] [US1] Build repositories in `backend/internal/domain/repository/privacy/` for classifications, consent tokens, and lifecycle events (CRUD + query helpers).
- [X] T103 [US1] Add services under `backend/internal/services/admin/security/privacy_service.go` and `backend/internal/services/agent/security/privacy_guard.go` handling consent validation, lifecycle logging, and host event hooks.
- [X] T104 [P] [US1] Implement `backend/internal/middleware/consent_guard/guard.go` consuming consent tokens and applying masking adapters.
- [X] T105 [US1] Wire admin HTTP handlers (`backend/internal/transport/http/admin/security/consent_handler.go`, `lifecycle_handler.go`) per OpenAPI contract, returning DTOs and status codes.
- [X] T106 [US1] Introduce agent-facing privacy endpoints in `backend/internal/transport/http/agent/security/privacy_handler.go` for consent status checks and lifecycle acknowledgements.
- [X] T107 [P] [US1] Extend log masking helpers in `backend/internal/logger/privacy_mask.go` and load rules from `backend/etc/security_baseline.yaml`.
- [X] T108 [US1] Append audit emission utility in `backend/internal/observability/security/audit_writer.go` to write consent and lifecycle events to `/logs/audit.log`.
- [X] T109 [US1] Build admin UI views in `web-admin/app/pages/security/consent.vue` and `events.vue` to list consent tokens and lifecycle evidence with masked previews.
- [X] T110 [US1] Implement outbound data gateway guard in `backend/internal/services/agent/security/privacy_guard.go` (host allow-list validation + TLS 1.3 enforcement) with tests under `backend/tests/security/privacy_guard_test.go`.
- [X] T111 [P] [US1] Add AI 数据合规处理：创建 `backend/internal/services/agent/security/ai_filter.go` 负责 prompt/response 脱敏与临时存储清理，并在 `backend/internal/middleware/consent_guard/guard.go` 中挂接。

**Checkpoint**: Consent enforcement, masking, and lifecycle evidence verifiable end-to-end.

---

## Phase 4 — User Story 2 — Apply the Plugin Security Baseline Checklist (Priority: P1)

**Goal**: Provide deterministic security checks (scans, sandbox validation, signature verification) and expose results for marketplace gating.

**Independent Test**: Run `make security-audit`; observe generated SARIF/JSON with no High/Critical findings, stored records, and visible via admin UI/API.

### Implementation for US2

- [X] T200 [P] [US2] Define models `backend/internal/domain/models/security/baseline_checklist.go` and `audit_report.go` reflecting checklist/version metadata.
- [X] T201 [US2] Add migrations for baseline checklists and audit reports with indexes supporting report queries.
- [X] T202 [US2] Implement repositories in `backend/internal/domain/repository/security/{checklist_repository.go,audit_report_repository.go}` to persist run metadata.
- [X] T203 [US2] Create service orchestrators in `backend/internal/services/admin/security/baseline_service.go` to launch scans and evaluate pass/fail thresholds.
- [X] T204 [US2] Author `make-files/security.mk` plus `Makefile` hook for `make security-audit`, invoking golangci-lint, govulncheck, gosec, npm audit, trivy, and cosign verification with outputs in `build/security/`.
- [X] T205 [P] [US2] Implement HTTP handlers `backend/internal/transport/http/admin/security/audit_report_handler.go` for listing audit reports and retrieving artifacts.
- [X] T206 [US2] Surface audit results in `web-admin/app/pages/security/baseline.vue` with download links and status badges.
- [X] T207 [US2] Update CI pipeline templates (`.github/workflows/` or equivalent) to run `make security-audit` and fail on High/Critical findings per FR-011A.
- [X] T208 [US2] Produce SARIF/JSON documentation in `docs/security/audit-pipeline.md`, describing command usage and waiver process.
- [X] T209 [P] [US2] 加固 Nuxt 安全配置：在 `web-admin/nuxt.config.ts` 中启用 CSP、SRI、严格的 `cors`/`headers`，并在 `web-admin/app/plugins/csrf.ts` 增加 CSRF token 注入。
- [X] T210 [US2] 编写前端安全验证脚本 `web-admin/tests/security/headers.spec.ts`，确认 CSRF/CSP/SRI 头与 token 生效。

**Checkpoint**: Security baseline pipeline operational with UI/API visibility.

---

## Phase 5 — User Story 3 — Govern ToolGrant Lifecycle & Consumption (Priority: P1)

**Goal**: Issue, validate, renew, and revoke ToolGrant tokens with least privilege, TTL ≤24h, and comprehensive audit logging.

**Independent Test**: Issue ToolGrant via `/_core/toolgrants`, access protected route, revoke grant, and confirm middleware blocks further access while audit trail records issuance/revocation.

### Implementation for US3

- [X] T300 [P] [US3] Implement models `backend/internal/domain/models/tool_grant/{revocation.go,usage_event.go}` capturing JWT `jti`, tenant, and metadata.
- [X] T301 [US3] Add migrations for tool grant revocations and usage events with TTL cleanup indexes.
- [X] T302 [P] [US3] Build repositories `backend/internal/domain/repository/tool_grant/{revocation_repository.go,usage_repository.go}`.
- [X] T303 [US3] Implement service layer `backend/internal/services/agent/tool_grant/service.go` for issuance, validation, renewal, and revocation orchestration.
- [X] T304 [P] [US3] Develop middleware `backend/internal/middleware/tool_grant_verifier/middleware.go` enforcing token validation and revocation checks.
- [X] T305 [US3] Expose agent verifier endpoint `backend/internal/transport/http/agent/security/toolgrant_handler.go` per contract, and admin revoke endpoint `backend/internal/transport/http/admin/security/toolgrant_handler.go`.
- [X] T306 [US3] Emit ToolGrant audit events via `backend/internal/observability/security/tool_grant_events.go`, forwarding to `/logs/audit.log` and metrics.
- [X] T307 [US3] Update admin UI `web-admin/app/pages/security/toolgrants.vue` to inspect active grants, revocations, and usage events.
- [X] T308 [US3] 对 `backend/internal/middleware/tool_grant_verifier/middleware.go` 进行性能基准测试（基准文件 `backend/internal/middleware/tool_grant_verifier/middleware_benchmark_test.go`），验证延迟增幅 <5%。

**Checkpoint**: ToolGrant lifecycle enforced with middleware and admin control.

---

## Phase 6 — User Story 4 — Execute Vulnerability Response End-to-End (Priority: P2)

**Goal**: Manage vulnerability intake, classification, signed patch publication, and tenant notifications within mandated SLAs.

**Independent Test**: Submit a vulnerability, confirm advisory is created with SLA timers, publish patch, distribute notifications, and verify tenants acknowledge through admin UI/API.

### Implementation for US4

- [X] T400 [P] [US4] Add models `backend/internal/domain/models/security/{advisory.go,distribution.go}` representing advisories and delivery receipts.
- [X] T401 [US4] Create migrations for advisories/distributions with severity/status indexes and notification tracking.
- [X] T402 [P] [US4] Build repositories `backend/internal/domain/repository/security/{advisory_repository.go,distribution_repository.go}`.
- [X] T403 [US4] Implement `backend/internal/services/admin/security/advisory_service.go` handling lifecycle (report → patch → publish → close) and SLA computation.
- [X] T404 [US4] Add admin endpoints `backend/internal/transport/http/admin/security/advisory_handler.go` for create/list/publish operations per OpenAPI contract.
- [X] T405 [US4] Implement event emitters in `backend/internal/observability/security/advisory_events.go` firing `plugin.vulnerability.detected`/`remediated` and queuing marketplace/webhook notifications.
- [X] T406 [US4] Update packaging pipeline (`Makefile`, `build/package.sh` or equivalent) to bundle signed advisories into `dist/security/` alongside `.pxp`.
- [X] T407 [P] [US4] Build admin UI `web-admin/app/pages/security/advisories.vue` for managing advisories and viewing distribution status.
- [X] T408 [US4] Document response playbook in `docs/security/vulnerability-response.md` including SLA tables and communication channels.

**Checkpoint**: Vulnerability response loop validated with signed advisories and notifications.

---

## Phase 7 — Polish & Cross-Cutting

- [X] T900 [P] Polish Update `specs/004-security-compliance/quickstart.md` with final command examples and screenshots as features land.
- [X] T901 Polish Refresh `docs/` and `README.md` sections covering consent workflow, security audit, ToolGrant usage, and vulnerability response.
- [X] T902 [P] Polish Ensure telemetry/metrics dashboards are configured (`backend/internal/observability/security/metrics.go`) and wired to existing exporters.
- [X] T903 Polish Run full validation suite (`make test`, `make security-audit`, Nuxt build) and capture artifact hashes in release checklist.
- [X] T904 Polish Prepare release notes + manifest updates for marketplace submission, including new `manifest.data_usage` and `security_baseline_version`.

---

## Dependencies & Execution Order

- **Phase Order**: Setup → Foundational → US1 → US2 → US3 → US4 → Polish.
- **Story Dependencies**:
  - US1 (privacy) has no upstream story dependencies once Foundational is complete.
  - US2 depends on Foundational only; can run parallel with US1 if bandwidth allows.
  - US3 depends on Foundational and benefits from consent audit utilities delivered in US1 (optional soft dependency).
  - US4 depends on Foundational and reuses audit/report storage patterns from US2.
- Ensure migrations from multiple stories are ordered with timestamps to avoid conflicts.

### Parallel Opportunities

- Setup tasks touch distinct files → can run concurrently.
- Within US1: T100, T102, T107 are parallelizable ([P]) once migrations (T101) land.
- Within US2: T200 and T204 operate independently ([P]) while other tasks chain.
- Within US3: T300, T302, T304 marked [P]; coordinate around migration T301 and service wiring.
- Within US4: T400, T402, T407 flagged [P] after migration T401.

---

## Implementation Strategy

1. **Unlock MVP** by completing Setup → Foundational → **US1** (privacy controls). This satisfies regulatory data handling and masks risk immediately.
2. **Harden security posture** with **US2** to deliver automated audits before enabling advanced capabilities.
3. **Tighten runtime access** via **US3** ToolGrant governance to ensure least privilege enforcement.
4. **Finalize trust loop** with **US4**, enabling rapid vulnerability handling and signed advisories.
5. Finish with polish tasks to align docs, telemetry, and release packaging before shipping.
