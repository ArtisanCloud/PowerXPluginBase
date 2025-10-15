# Tasks: Plugin Capabilities & Schema Governance

**Input**: Design docs from `/specs/002-title-plugin-capabilities/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: Prioritise Go unit tests for catalog + validator helpers and smoke-run the new Make targets (`make check-capability`, `make check-compat`) against the sample plugin.

**Organization**: Tasks are grouped by user story so each slice can be delivered independently.

> **Format**: `[ID] [P?] [Story] Description`  
> `[P]` marks work that can proceed in parallel (different files, no ordering risk).

## Path Conventions

- Capability assets: `contracts/` (`capabilities/`, `schema/input/`, `schema/output/`)
- Backend helpers & validators: `backend/internal/contracts/`, `backend/cmd/manifestcheck/`
- Automation: `make-files/`, `Makefile`, `scripts/`
- Lifecycle documentation: `docs/lifecycle/`, `docs/integration/02_capabilities_and_schema/`

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: Scaffold shared directories and guidance consumed by all stories.

- [X] T001 [Setup] Create capability/schema directories with `.gitkeep` placeholders to anchor new artefacts  
      `contracts/capabilities/.gitkeep`, `contracts/schema/input/.gitkeep`, `contracts/schema/output/.gitkeep`
- [X] T002 [Setup] Author `contracts/README.md` explaining directory intent, naming rules, and linking to feature quickstart  
      `contracts/README.md`

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Establish reusable capability data structures and loaders before story work.

- [ ] T003 [Foundation] Define `CapabilityDescriptor`, `SchemaRef`, and related structs with YAML/JSON tags  
      `backend/internal/contracts/capability.go`
- [ ] T004 [Foundation] Implement filesystem catalog loader that indexes `contracts/capabilities/*.yaml` and exposes lookup helpers  
      `backend/internal/contracts/capability_loader.go`

**Checkpoint**: Capability catalog primitives exist and can be reused by manifest + validator stories.

---

## Phase 3 — User Story 1 — Canonical Capability Descriptors (Priority: P1) 🎯 MVP

**Goal**: Plugin manifests reference canonical capability descriptors and versioned schemas stored in `contracts/`.

**Independent Test**: Run `go test ./backend/internal/contracts/...` plus `make verify-manifest` to confirm the manifest resolves declared capability IDs and schema paths without manual tweaking.

### Tasks

- [ ] T005 [US1] Extend `PluginManifest` contracts to surface `capabilities.provides/consumes` entries and schema references  
      `backend/internal/contracts/manifest.go`
- [ ] T006 [P] [US1] Author capability descriptors for existing template tools with version + RBAC metadata  
      `contracts/capabilities/base.template.create.yaml`, `contracts/capabilities/base.template.query.yaml`
- [ ] T007 [P] [US1] Create draft-07 JSON Schemas for input/output payloads referenced by the new capability descriptors  
      `contracts/schema/input/base.template.create.v1.json`, `contracts/schema/output/base.template.create.v1.json`, `contracts/schema/output/base.template.query.v1.json`
- [ ] T008 [US1] Add loader-focused unit tests covering happy-path parsing and schema list resolution  
      `backend/internal/contracts/capability_loader_test.go`
- [ ] T009 [US1] Update manifest samples to list capability IDs and schema assets for packaging parity  
      `plugin.yaml`, `docs/lifecycle/examples/manifest.yaml`
- [ ] T010 [US1] Publish governance guidance describing capability authoring + manifest references  
      `docs/lifecycle/capabilities.md`, `docs/integration/02_capabilities_and_schema/Capability_Design_Guide.md`

**Checkpoint**: Capabilities ship as first-class assets and manifest references are traceable to concrete schemas.

---

## Phase 4 — User Story 2 — Automated Capability Validation (Priority: P2)

**Goal**: Provide tooling that validates capability descriptors, schema references, and RBAC alignment via `make check-capability`.

**Independent Test**: Execute `make check-capability` against the sample plugin; the command must pass with canonical assets and fail when capability/RBAC mismatches are introduced.

### Tasks

- [ ] T011 [US2] Implement capability validator that checks descriptor uniqueness, schema existence, and RBAC resource/action parity  
      `backend/internal/contracts/capability_validator.go`
- [ ] T012 [US2] Extend `backend/cmd/manifestcheck` to invoke the validator and emit structured errors  
      `backend/cmd/manifestcheck/main.go`
- [ ] T013 [US2] Add validator unit tests covering missing schema, orphaned RBAC, and duplicate capability scenarios  
      `backend/internal/contracts/capability_validator_test.go`
- [ ] T014 [US2] Wire new validator into automation via `make check-capability`, reusing Go build cache for CLI execution  
      `make-files/manifest.mk`, `Makefile`
- [ ] T015 [US2] Document validator usage, expected outputs, and remediation steps  
      `docs/lifecycle/tooling.md`

**Checkpoint**: CI-ready validator ensures capability descriptors stay in sync with manifest + RBAC contracts.

---

## Phase 5 — User Story 3 — Compatibility Diff Automation (Priority: P3)

**Goal**: Generate compatibility reports for capability/schema version changes via `make check-compat`.

**Independent Test**: Run `make check-compat` to produce diff artefacts under `build/compat/` comparing staged assets with declared baselines.

### Tasks

- [ ] T016 [US3] Implement compatibility diff script (Node 20) that wraps `ajv`, `json-diff`, and `openapi-diff` for schema/capability comparisons  
      `scripts/check-compatibility.mjs`
- [ ] T017 [P] [US3] Declare Node CLI dependencies and npm scripts for compatibility tooling  
      `scripts/package.json`, `scripts/package-lock.json`
- [ ] T018 [US3] Add `make check-compat` target and shared variables for diff output paths  
      `make-files/compat.mk`, `Makefile`
- [ ] T019 [P] [US3] Seed baseline compatibility config and ensure artefact directories are version-controlled  
      `contracts/compatibility.yaml`, `build/compat/.gitkeep`
- [ ] T020 [US3] Document compatibility workflow and update quickstart instructions  
      `docs/lifecycle/package.md`, `specs/002-title-plugin-capabilities/quickstart.md`

**Checkpoint**: Compatibility automation surfaces breaking changes with auditable reports.

---

## Phase 6 — Polish & Cross-Cutting

- [ ] T021 [Polish] Extend lifecycle validation workflow to run `make check-capability` and `make check-compat`  
      `.github/workflows/lifecycle-validation.yml`
- [ ] T022 [P] [Polish] Update release checklist with capability/schema validation gates and compatibility report attachment  
      `docs/lifecycle/checklists/release-checklist.md`
- [ ] T023 [Polish] Refresh release template to call out capability descriptors, schema bundles, and diff artefacts  
      `docs/releases/release-template.md`

---

## Dependencies & Execution Order

### Phase Dependencies

1. Setup → 2. Foundational → 3. US1 → 4. US2 → 5. US3 → 6. Polish  
   - US1 relies on catalog primitives (Phase 2) and directory scaffolding (Phase 1).  
   - US2 requires US1 artefacts to validate manifest ↔ capability parity.  
   - US3 consumes validated capability data before generating compatibility reports.

### Parallel Opportunities

- **US1**: T006 and T007 can proceed concurrently once manifest structures land because they touch independent capability/schema files.  
- **US3**: T017 and T019 can run in parallel after the diff script contract (T016) is defined.  
- **Polish**: T022 can be handled alongside T021 once automation outputs are verified.

### Story Independence

- **US1** delivers the MVP: canonical capability descriptors referenced by manifest + sample schemas.  
- **US2** builds on US1 assets to guarantee automated validation, independently testable via `make check-capability`.  
- **US3** layers compatibility automation atop US2, verified by `make check-compat` outputs.

### Suggested MVP Scope

- Complete Phases 1–3 (through US1) to unlock a canonical capability catalog and manifest references before investing in validator and compatibility automation.
