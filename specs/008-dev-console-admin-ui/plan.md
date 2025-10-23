# Implementation Plan: Dev Console & Admin UI (Plugin Admin, Audit & History, Troubleshooting)

**Branch**: `008-dev-console-admin-ui` | **Date**: 2025-10-22 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/008-dev-console-admin-ui/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

---

## Summary

Deliver a consistent PowerX plugin admin console that lets authorized operators adjust configuration with validation and auditing, inspect activity history, and troubleshoot jobs/webhooks from a single Nuxt-driven UI backed by new Go admin services, RBAC guards, and observability integrations.

---

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.24 backend, Node 20 / TypeScript 4.x / Nuxt 4 frontend  
**Primary Dependencies**: Gin HTTP stack, GORM (PostgreSQL driver), Zap logging, Nuxt UI 3.3 components, Pinia stores, STS client utilities  
**Storage**: PostgreSQL (schema `powerx_plugin_base`); new tables for configuration change history, admin audit events, job run snapshots; relies on existing observability aggregates for health/quota/webhooks  
**Testing**: `make test` (Go unit/integration), targeted service tests under `backend/internal/services/admin`, HTTP contract tests with httptest, frontend unit tests via `npm run test` (Nuxt Vitest) plus Playwright smoke for critical flows  
**Target Platform**: PowerX plugin deployment (Linux container) behind reverse proxy paths `/_p/<plugin_id>/admin/*` and `/_p/<plugin_id>/api/v1/*`  
**Project Type**: Backend + admin web application (Go services + Nuxt admin UI)  
**Performance Goals**: Audit exports up to 1,500 events under 5s, console configuration updates applied within 3 minutes end-to-end, troubleshooting dashboard auto-refresh every 5 minutes without blocking UI, safe-ops actions acknowledge within 30s  
**Constraints**: Enforce tenant RLS and RBAC, operate within STS-scoped credentials, avoid introducing new long-running tasks in request path, keep frontend bundle growth <150KB gz to preserve console load targets  
**Scale/Scope**: Multi-tenant plugin with dozens of operators per plugin instance, ~100 concurrent admin sessions across tenants, job history covering 30 days (~10k entries) per plugin

### Platform / Hosting Integration (optional)

If this feature runs **under a host platform** (e.g., PowerX or similar), specify:

- **Reverse Proxy & Routes**: Admin APIs will mount under `/_p/<plugin_id>/api/v1/admin/dev-console/**`; frontend served at `/_p/<plugin_id>/admin/dev-console/**` with nested tabs for configuration, audit, and troubleshooting. Public runtime APIs remain under `/v1/**`.  
- **Context Signing**: Reuse existing JWT verification middleware; requests include `tenant_id`, `plugin_id`, `user_id`, `permissions`, `request_id`, `exp`, `iat`. Replay-safe operations require validating permission codes such as `operations.plugin.admin`, `operations.plugin.audit`, `operations.plugin.ops`.  
- **Tenant/RBAC**: Each handler enters tenant-scoped transaction via `BeginTenantTx`, sets `app.tenant_id`, and enforces RLS on new tables. Server-side RBAC checks map to manifest-declared permission codes.  
- **Outbound Access**: Continue using STS exchange utilities for observability and webhook diagnostics; no new long-lived credentials introduced.  
- **Observability**: Extend `backend/internal/observability/admin_console` with counters for safe-ops usage, audit export metrics, and dashboard freshness; ensure `/healthz` covers new dependencies.  
- **Packaging**: Update `plugin.yaml` manifest with new admin routes and permission codes; release artifacts include updated Go binaries and Nuxt bundle via `make release`.

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*(If your `constitution.md` defines Gate IDs, reference them in braces for traceability.)*

- [x] **Host Contract First** {PX-HOST-001}  
  Admin APIs planned under `/api/v1/admin/dev-console/**`, manifests and reverse-proxy metadata updated alongside new navigation entries.
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  All new tables include `tenant_id` with RLS, handlers reuse JWT validation, and tests will cover cross-tenant access denials.
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  New `admin/console` handlers delegate to services in `internal/services/admin/console`, backed by repositories in `internal/domain/repository/admin/console`.
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  Permission codes `operations.plugin.admin`, `operations.plugin.audit`, `operations.plugin.ops` anchored in manifest; server-side enforcement precedes operations.
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  Plan includes metrics/structured logs plus Go unit/integration tests and Nuxt Vitest/Playwright coverage; `/healthz` incorporates new dependencies.
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  No new external services; release tasks updated for binaries/UI bundle; SemVer minor bump with changelog.

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
├── internal/
│   ├── domain/models/admin_console/
│   ├── domain/repository/admin_console/
│   ├── services/admin/console/
│   └── transport/http/admin/console/
├── migrations/
└── tests/

web-admin/
├── app/
│   ├── pages/_p/com.powerx.plugins.base/admin/dev-console/
│   ├── components/dev-console/
│   ├── stores/dev-console/
│   └── composables/useDevConsole*
└── tests/operations/dev-console/

```

**Structure Decision**: Extend existing backend admin layering with `console` domain folders and surface Nuxt admin pages/components/stores under `web-admin/app` to match plugin console routing; job coordination code lives within `services/admin/console`.

---

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| (none) |  |  |
