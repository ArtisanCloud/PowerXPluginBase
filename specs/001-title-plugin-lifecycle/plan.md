# Implementation Plan: Plugin Lifecycle Governance

**Branch**: `[001-title-plugin-lifecycle]` | **Date**: 2025-10-14 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/001-title-plugin-lifecycle/spec.md`

**Note**: Generated via `/speckit.plan` with repository context applied.

---

## Summary

Establish a unified lifecycle governance standard for PowerX plugins covering workspace bootstrap, manifest metadata, version packaging, and deprecation policy. Delivery will include authoritative documentation, automated manifest/packaging validation hooks, and supporting examples so release managers, engineers, and lifecycle owners can produce verifiable `.pxp` artefacts, register them in the Marketplace, and manage sunsetting without ad-hoc playbooks.

---

## Technical Context

**Language/Version**: Go 1.24 (backend services), TypeScript (Nuxt 4 admin UI)  
**Primary Dependencies**: Echo-based HTTP stack with internal middleware, golang-migrate, Nuxt UI 3.3.x, Makefile automation targets (`make build`, `make release`)  
**Storage**: PostgreSQL per-plugin schema with RLS; migrations defined under `backend/internal/db/migrations`  
**Testing**: `make test` (Go unit/integration), `make frontend-test` (Nuxt), contract validation scripts for manifests  
**Target Platform**: PowerX host runtime (Linux containers) with Marketplace distribution  
**Project Type**: Dual-stack plugin (Go backend + Nuxt admin frontend)  
**Performance Goals**: Packaging and validation pipelines complete within CI time budget (<10 minutes) and manifest checks process within 30 seconds locally  
**Constraints**: Release artefacts must be deterministic, signed, and audit logged; guidance must remain tool-agnostic with Makefile defaults  
**Scale/Scope**: Supports dozens of plugin teams, each managing multiple SemVer channels and lifecycle statuses concurrently

### Platform / Hosting Integration (optional)

This standard governs plugins operating under PowerX host requirements:

- **Reverse Proxy & Routes**: Reinforce `/v1/**` business APIs and `/api/v1/admin/{manifest,rbac}` admin endpoints; manifest updates feed host catalog entries.
- **Context Signing**: Document JWT-based verification with tenant/user claims; endorse HMAC only for legacy sandboxes.
- **Tenant/RBAC**: Require server-side RBAC enforcement and RLS session variables in migrations and smoke tests.
- **Outbound Access**: Mandate STS token exchange before calling core PowerX APIs.
- **Observability**: Standardize `/healthz`, structured logging keys, and lifecycle audit logs during packaging.
- **Packaging**: SemVer-driven channel promotion with signed `.pxp` artefacts, SHA256 manifests, and Marketplace registration receipts.

> If there is **no host**, mark N/A and explain the external boundaries (clients, gateways, schedulers, etc.).

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*  
*(If your `constitution.md` defines Gate IDs, reference them in braces for traceability.)*

- [x] **Host Contract First** {PX-HOST-001}  
  Lifecycle standard will codify routing, manifest obligations, and ensure manifest/admin endpoints stay aligned with host contract examples.
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  Documentation updates include mandatory JWT verification, RLS migration checklists, and validation tasks to confirm tenant context injection.
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  Plan preserves thin transports by referencing service-layer reuse and updating guidelines rather than introducing ad-hoc logic.
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  Standard will enumerate RBAC manifest expectations and release review steps that confirm least-privilege policies before publication.
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  Packaging guidance covers `/healthz`, structured logs, and requires automated validation hooks plus documented test plans.
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  Deliverables include SemVer channel governance, artefact checklists, and dependency minimization guidance to support lean packages.

> Any unchecked item must be resolved or explicitly justified in **Complexity Tracking** below.

---

## Project Structure

### Documentation (this feature)

```

specs/001-title-plugin-lifecycle/
├── plan.md              # Implementation planning (this document)
├── research.md          # Phase 0 findings and decisions
├── data-model.md        # Lifecycle entities, fields, relationships
├── quickstart.md        # How to apply the lifecycle standard
├── contracts/           # Manifest & lifecycle contract artifacts
└── tasks.md             # Delivery breakdown (from /speckit.tasks)

```

### Source Code (repository root)

```

backend/
├── cmd/
│   ├── database/        # Migration entrypoint
│   └── plugin/          # Runtime entrypoint
├── etc/                 # Config samples
├── internal/
│   ├── contracts/       # API contracts (admin + business)
│   ├── db/              # Migrations and schema helpers
│   ├── router/          # HTTP route setup
│   ├── services/        # Domain logic
│   └── transport/       # HTTP & gRPC handlers
└── plugin/              # Plugin metadata (plugin.yaml, manifest)

web-admin/
├── app/                 # Pages, components, layouts, stores
├── i18n/                # Localization resources
└── tests/               # Frontend test suite

docs/
└── contract/            # Existing plugin specs (e.g., plugin_yaml_spec.md)

build/
└── pxp/                 # Release artefact staging (to be standardized)

```

**Structure Decision**: Retain current backend/frontend split aligned with plugin constitution. Lifecycle documentation lives under `docs/lifecycle/`, with a documented sync process that publishes curated outputs into `docs/integration/01_plugin_lifecycle/`. Packaging validation scripts land under `make-files/` targets, and manifest templates within `backend/plugin/`.

---

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| [e.g., extra transport] | [current need] | [why single transport insufficient] |
| [e.g., custom repo pattern] | [specific problem] | [why direct DB access insufficient] |
