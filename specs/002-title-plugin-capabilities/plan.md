# Implementation Plan: Plugin Capabilities & Schema Governance

**Branch**: `[002-title-plugin-capabilities]` | **Date**: 2025-10-14 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/002-title-plugin-capabilities/spec.md`

---

## Summary

Establish a single source of truth for plugin capabilities and IO schemas by standardising capability descriptors (`contracts/capabilities/*.yaml`), enforcing JSON Schema (draft-07) for input/output contracts, and defining compatibility/automation policies. The change introduces manifest references to capability files, RBAC linkage validation, and CI tooling (`make check-capability`, `make check-compat`) to guarantee semantic alignment across host, Marketplace, and Agent layers.

---

## Technical Context

**Language/Version**: Go 1.24 (backend tooling), Node 20 (schema lint scripts), JSON Schema draft-07  
**Primary Dependencies**: `github.com/santhosh-tekuri/jsonschema/v5` (Go validator), `ajv` CLI, existing PowerX contract helpers  
**Storage**: N/A (YAML/JSON files under `contracts/`)  
**Testing**: `go test ./backend/...`, schema validation via `make check-capability`, compatibility diff via `json-diff`, `openapi-diff`  
**Target Platform**: PowerX host (Linux containers), Marketplace registration, Agent tool-generation  
**Project Type**: Backend contracts + documentation  
**Performance Goals**: Validation/diff commands finish <60s in CI; handle ≥100 capability descriptors  
**Constraints**: SemVer compliance; manifest stays lightweight; tooling must operate offline  
**Scale/Scope**: Supports multiple capability versions (v1, v2) per plugin with controlled deprecation windows

### Platform / Hosting Integration

- **Reverse Proxy & Routes**: Existing admin endpoints (`/api/v1/admin/manifest`, `/api/v1/admin/rbac`) remain; manifest now references external capability descriptors consumed by host.
- **Context Signing**: Unchanged; capability governance must not bypass tenant context enforcement.
- **Tenant/RBAC**: Capability IDs mapped to RBAC actions get automated parity checks.
- **Outbound Access**: N/A.
- **Observability**: Tooling emits structured logs/audit artifacts for compatibility checks.
- **Packaging**: `.pxp` bundles must include `contracts/`; manifest references validated in packaging flow.

---

## Constitution Check

- [x] **Host Contract First** {PX-HOST-001} — Manifest references capability descriptors, host contract remains authoritative.
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001} — Capability↔RBAC validations prevent privilege mismatches.
- [x] **Service-Centric Architecture** {PX-SVC-001} — No new handlers; contracts tooling augments existing services.
- [x] **RBAC & Least Privilege** {PX-RBAC-001} — Automated checks ensure capability IDs align with RBAC resources/actions.
- [x] **Observable & Testable Delivery** {PX-OBS-001} — Introduce CI targets (`make check-capability`, `make check-compat`) with audit outputs.
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001} — SemVer policy clarified; capability/schema versions tracked independently with deprecation guidance.

---

## Project Structure

### Documentation (this feature)

```
specs/002-title-plugin-capabilities/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   ├── capability-example.yaml
│   └── schema-example.json
└── tasks.md  (via /speckit.tasks)
```

### Source Code & Artefacts

```
backend/
├── internal/contracts/           # Manifest structs + capability loading helpers
├── cmd/manifestcheck/            # Extend to cover capability/schema checks

contracts/
├── capabilities/                 # New YAML capability descriptors
├── schema/input/                 # JSON Schemas for inbound payloads
└── schema/output/                # JSON Schemas for outbound payloads

docs/
└── lifecycle/                    # Governance docs extended with capability guide

make-files/
├── manifest.mk                   # Manifest/capability validation targets
├── compat.mk                     # New compatibility diff orchestration
└── docs.mk                       # Already handles doc sync
```

**Structure Decision**: Reuse existing backend contract infrastructure, create dedicated `contracts/capabilities` and `contracts/schema` subfolders, and add make targets for capability/schema validation plus compatibility diff. No new runtime packages—focus on contract artefacts, docs, and automation.

---

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| _None_ | — | — |
