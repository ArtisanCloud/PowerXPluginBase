---
description: "Task list template for feature implementation (tool-agnostic)"
---

# Tasks: [FEATURE NAME]

**Input**: Design docs from `/specs/[###-feature-name]/`  
**Prerequisites**: `plan.md` (required), `spec.md` (for user stories), `research.md`, `data-model.md`, `contracts/`

**Tests**: Per your project constitution (e.g., *Observable & Testable Delivery*), add the necessary test tasks for **each** change unless the spec explicitly says “no executable code impact”.

**Organization**: Tasks are grouped by **user story** so each story can be implemented, tested, and delivered **independently**.

> **Format**: `[ID] [P?] [Story] Description`  
>
> - **[P]**: Can run in parallel (different files, no dependencies)  
> - **[Story]**: Which story this task belongs to (e.g., US1, US2)  
> - Include **exact file paths** in descriptions (reflecting `plan.md` structure)

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

**Purpose**: Initialize workspace and shared toolchain.

- [ ] T001 Create project structure per `plan.md` (sync real directories, remove unused templates)
- [ ] T002 Initialize language toolchain & dependencies (e.g., Go/Nuxt/Node)
- [ ] T003 [P] Configure lint/format/test runners (e.g., golangci-lint, eslint, prettier, go test)

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete **before any user story**.

> ⚠️ **Gates Reminder** (optional; reference your constitution IDs if any):  
>
> - Context & RLS (e.g., {PX-CTX-001?})  
> - RBAC least privilege (e.g., {PX-RBAC-001?})  
> - Service-centric reuse (e.g., {PX-SVC-001?})  
> - Observability basics (e.g., {PX-OBS-001?})

- [ ] T010 Database migrations & policies (ensure RLS/tenant context variable is honored)  
      `backend/internal/db/migrations/*`
- [ ] T011 [P] Wire **TenantContext / RBAC** guards & middleware stack  
      `backend/internal/transport/http/middleware/*`
- [ ] T012 [P] Service/Repository scaffolding（no business logic yet）  
      `backend/internal/services/...`, `backend/internal/domain/repository/...`
- [ ] T013 Structured logging & `/healthz` probe  
      `backend/internal/transport/http/healthz.go`
- [ ] T014 Configuration & secrets (env, CI). If a host exists, prepare **STS** exchange client

**Checkpoint**: Foundation ready → user stories can start (in parallel if staffed).

---

## (Optional) Platform / Hosting Integration

*(Only if your feature runs under a host platform such as PowerX; otherwise mark N/A and skip.)*

- [ ] T020 Reverse-proxy routes: Business `/v1/**`; Admin contracts `/api/v1/admin/{manifest,rbac}`  
- [ ] T021 Context verification (JWT or HMAC): middleware verifies `tenant_id/user_id/permissions/...`  
- [ ] T022 DB session var for RLS: set on each request (e.g., `SET LOCAL app.tenant_id=?`)  
- [ ] T023 STS client: obtain short-lived token with explicit scopes for outbound calls  
- [ ] T024 Release artifacts checklist (SemVer, `plugin.yaml`, `backend/bin/plugin`, `web-admin/.output`, `checksums.txt`, `manifest.json`)

> These tasks correspond to typical host-related gates; keep them if applicable.

---

## Phase 3 — User Story 1 — [Title] (Priority: P1) 🎯 MVP

**Goal**: [What P1 delivers end-to-end]  
**Independent Test**: [How to verify this story alone]

### Tests for US1 (write-first, must FAIL initially)

- [ ] T101 [P] [US1] Contract test for `[endpoint]`  
      `tests/contract/test_[us1]_[endpoint].go|ts|py`
- [ ] T102 [P] [US1] Integration test for `[journey]`  
      `tests/integration/test_[us1]_[journey].go|ts|py`

### Implementation for US1

- [ ] T103 [P] [US1] DTOs & validators  
      `backend/internal/transport/http/[feature]/dto_[us1].go`
- [ ] T104 [P] [US1] Repository methods (no cross-story deps)  
      `backend/internal/domain/repository/[feature]/repo_[us1].go`
- [ ] T105 [US1] Service use-cases (single source of business logic)  
      `backend/internal/services/[feature]/service_[us1].go`
- [ ] T106 [US1] HTTP handler (thin): validate → RBAC → call service → marshal response  
      `backend/internal/transport/http/[feature]/handler_[us1].go`
- [ ] T107 [US1] (Optional) gRPC transport mapping to the same service  
      `backend/internal/grpc/server/[feature]/[us1]_server.go`
- [ ] T108 [US1] Frontend page/action (if any)  
      `frontend/app/pages/[feature]/[us1].vue`

**Checkpoint**: US1 is independently functional & testable (MVP).

---

## Phase 4 — User Story 2 — [Title] (Priority: P2)

**Goal**: […]  
**Independent Test**: […]

### Tests for US2

- [ ] T201 [P] [US2] Contract test for `[endpoint]`
- [ ] T202 [P] [US2] Integration test for `[journey]`

### Implementation for US2

- [ ] T203 [P] [US2] DTOs & validators
- [ ] T204 [P] [US2] Repository methods
- [ ] T205 [US2] Service use-cases
- [ ] T206 [US2] HTTP handler (thin) / optional gRPC mapping
- [ ] T207 [US2] Frontend page/action (if any)

**Checkpoint**: US1 & US2 both independently pass.

---

## Phase 5 — User Story 3 — [Title] (Priority: P3)

**Goal**: […]  
**Independent Test**: […]

### Tests for US3

- [ ] T301 [P] [US3] Contract test for `[endpoint]`
- [ ] T302 [P] [US3] Integration test for `[journey]`

### Implementation for US3

- [ ] T303 [P] [US3] DTOs & validators
- [ ] T304 [P] [US3] Repository methods
- [ ] T305 [US3] Service use-cases
- [ ] T306 [US3] HTTP handler (thin) / optional gRPC mapping
- [ ] T307 [US3] Frontend page/action (if any)

**Checkpoint**: All targeted stories are independently functional.

---

## Phase N — Polish & Cross-Cutting

- [ ] T900 [P] Docs updates in `docs/`
- [ ] T901 Refactors & dead code cleanup
- [ ] T902 Performance tuning / indexes / caching
- [ ] T903 [P] Additional unit tests in `tests/unit/`
- [ ] T904 Security hardening
- [ ] T905 Validate `quickstart.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup → Foundational → User Stories → Polish**
- No story work before **Foundational** is complete
- After Foundational, stories can proceed **in parallel** (if staffed) or **by priority** (P1 → P2 → P3)

### Within a Story

- **Tests first** (must fail)
- **Models/DTOs → Repo → Service → Transport (HTTP/gRPC) → Frontend (if any)**
- Keep handlers thin; **Service is the single source of business logic**

### Parallel Opportunities

- `[P]` tasks in Setup/Foundational
- Tests marked `[P]` within the same story
- Models/DTOs marked `[P]` within the same story
- Different stories by different contributors in parallel

---

## Notes

- Mark `[P]` only when tasks edit **different files** with **no dependency**.  
- Each story must be independently testable & demo-able.  
- Keep tasks concrete with real file paths.  
- Commit in small increments; stop at checkpoints for validation.
