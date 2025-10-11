# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]  
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

---

## Summary

[Extract from feature spec: primary requirement + technical approach from research]

---

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: [e.g., Go 1.22, Python 3.11, Swift 5.9, Rust 1.75 or NEEDS CLARIFICATION]  
**Primary Dependencies**: [e.g., Gin/Fiber/FastAPI/Actix or NEEDS CLARIFICATION]  
**Storage**: [e.g., PostgreSQL schema: <name>, migrations needed or N/A]  
**Testing**: [e.g., go test/pytest/XCTest; unit/integration/contract]  
**Target Platform**: [e.g., Linux server, iOS 15+, WASM or NEEDS CLARIFICATION]  
**Project Type**: [single/web/mobile - determines source structure]  
**Performance Goals**: [domain-specific, e.g., 1000 req/s, 60 fps]  
**Constraints**: [e.g., <200ms p95, <100MB memory, offline-capable]  
**Scale/Scope**: [e.g., 10k users, 50 screens, 1M rows/day]

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

- [ ] **Host Contract First** {PX-HOST-001?}  
  Document how the feature respects `/v1/**` routing; list admin endpoints if any and required manifest/metadata updates.
- [ ] **Tenant Isolation & Zero Trust** {PX-CTX-001?}  
  Explain context verification (JWT/HMAC), RLS-covered migrations, and how tests/tooling verify tenant safety.
- [ ] **Service-Centric Architecture** {PX-SVC-001?}  
  Identify impacted services/repositories; confirm HTTP/gRPC/MCP handlers stay thin and reuse the same Service layer.
- [ ] **RBAC & Least Privilege** {PX-RBAC-001?}  
  Map resources/actions and enforce server-side authorization (UI only controls visibility).
- [ ] **Observable & Testable Delivery** {PX-OBS-001?}  
  Outline logging/metrics, `/healthz`, and automated tests (unit/integration/contract/migration smoke).
- [ ] **Minimal Footprint & Versioned Releases** {PX-PKG-001?}  
  Call out dependency changes, packaging tasks, SemVer bump rationale, and release artifacts checklist.

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

# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)

src/
├── models/
├── services/
├── api/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (frontend + backend)

backend/
├── internal/
│   ├── domain/models/
│   ├── domain/repository/
│   ├── services/
│   └── transport/{http,grpc}/
└── tests/

frontend/
├── app/ (pages/layouts/components/stores)
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API

api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]

```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

---

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| [e.g., extra transport] | [current need] | [why single transport insufficient] |
| [e.g., custom repo pattern] | [specific problem] | [why direct DB access insufficient] |
