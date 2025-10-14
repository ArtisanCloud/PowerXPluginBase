# Tasks: Plugin Lifecycle Governance

**Input**: Design docs from `/specs/001-title-plugin-lifecycle/`  
**Prerequisites**: `plan.md`, `spec.md`, `research.md`, `data-model.md`, `contracts/`

**Tests**: Focus on documentation and automation validation. Add executable tests only where packaging/validation scripts produce code paths.

**Organization**: Tasks are grouped by user story so each slice can be delivered independently.

> **Format**: `[ID] [P?] [Story] Description`  
> `[P]` marks work that can proceed in parallel (different files, no ordering risk).

## Path Conventions

- Backend runtime + metadata: `backend/`, `backend/plugin/`
- Build automation: `make-files/`, `Makefile`
- Lifecycle docs: `docs/lifecycle/`
- Published integration docs: `docs/integration/01_plugin_lifecycle/`
- Contracts & references: `docs/contract/`
- Examples & templates: `build/pxp/`, `docs/lifecycle/examples/`

---

## Phase 1 — Setup (Shared Infrastructure)

**Purpose**: Establish lifecycle documentation assets used by downstream integration guides.

- [X] T001 [US-All] Create lifecycle documentation hub (`docs/lifecycle/overview.md`, `docs/lifecycle/_index.md`) and add reference pointer in `docs/readme.md`  
      `docs/lifecycle/`, `docs/readme.md`
- [X] T002 [P] [US-All] Set up examples workspace under `docs/lifecycle/examples/` and stage current `backend/plugin/plugin.yaml` + `manifest.yaml` snapshots for reference  
      `docs/lifecycle/examples/`
- [X] T003 [P] [US-All] Update developer quick references to point at lifecycle sources (`docs/contract/plugin_yaml_spec.md`, `docs/contract/rbac_manifest_spec.md`)  
      `docs/contract/`

---

## Phase 2 — Foundational (Blocking Prerequisites)

**Purpose**: Shared groundwork before individual stories branch off.

- [X] T010 [US-All] Capture current build + release flow by diagramming `make run`, `make build`, `make release` dependencies for inclusion in lifecycle docs  
      `docs/lifecycle/overview.md`
- [X] T011 [P] [US-All] Draft lifecycle glossary aligning terminology across spec, plan, and docs (Plugin Package, Lifecycle Manifest, Release Channel Record, Status Ledger)  
      `docs/lifecycle/glossary.md`
- [X] T012 [US-All] Inventory existing automation scripts and note extension points for lifecycle tasks (Makefiles, backend/plugin metadata)  
      `make-files/`, `Makefile`, `backend/plugin/`
- [X] T013 [US-All] Implement sync routine (e.g., `make sync-lifecycle-docs`) and document how lifecycle content is published into `docs/integration/01_plugin_lifecycle/`  
      `make-files/docs.mk`, `docs/integration/01_plugin_lifecycle/README.md`

**Checkpoint**: Common understanding of assets, terminology, and publication flow established.

---

## Phase 3 — User Story 1 — Bootstrap A Compliant Plugin Workspace (Priority: P1) 🎯 MVP

**Goal**: Deliver an authoritative onboarding playbook so teams create audit-ready plugin repositories.

**Independent Test**: Follow the bootstrap checklist against a fresh clone and confirm all required directories, configs, and migrations are present without tribal knowledge.

### Tasks

- [X] T101 [US1] Author `docs/lifecycle/bootstrap.md` covering repo fork/rename, directory layout, environment setup, and migration execution  
      `docs/lifecycle/bootstrap.md`
- [X] T102 [P] [US1] Produce bootstrap checklist (`docs/lifecycle/checklists/bootstrap-checklist.md`) referencing RLS, RBAC, and build prep items  
      `docs/lifecycle/checklists/bootstrap-checklist.md`
- [X] T103 [US1] Update `docs/lifecycle/overview.md` with init → bootstrap timeline and link to checklist/quickstart  
      `docs/lifecycle/overview.md`
- [X] T104 [P] [US1] Embed bootstrap quickstart snippet into `docs/lifecycle/quickstart.md` and cross-link from repository `README.md`  
      `docs/lifecycle/quickstart.md`, `README.md`
- [X] T105 [US1] Align `docs/contract/plugin_yaml_spec.md` to reference lifecycle bootstrap requirements and examples  
      `docs/contract/plugin_yaml_spec.md`

**Checkpoint**: US1 documentation reviewed and validated via dry run.

---

## Phase 4 — User Story 2 — Release A Verifiable .pxp Package (Priority: P1)

**Goal**: Provide deterministic packaging instructions, validation tooling, and manifest parity guidance.

**Independent Test**: Execute `make verify-manifest && make package-pxp` on a plugin and achieve Marketplace-ready artefacts with signed hashes and synced metadata.

### Tasks

- [ ] T201 [US2] Author manifest mapping guide bridging `backend/plugin/plugin.yaml` to release-time `manifest.yaml`  
      `docs/lifecycle/manifest-mapping.md`
- [ ] T202 [P] [US2] Introduce JSON Schema for manifest validation and store under lifecycle documentation set (align with contract asset)  
      `docs/lifecycle/contracts/manifest.schema.json`
- [ ] T203 [US2] Add `make verify-manifest` target wiring schema validation for `plugin.yaml` + `manifest.yaml`, including CI hook  
      `Makefile`, `make-files/manifest.mk`
- [ ] T204 [P] [US2] Implement `make package-pxp` pipeline to stage artefacts in `build/pxp/`, compute SHA256, and collect signature placeholders  
      `make-files/release.mk`, `build/pxp/`
- [ ] T205 [US2] Document packaging workflow, audit logging, and rollback expectations (`docs/lifecycle/package.md`)  
      `docs/lifecycle/package.md`
- [ ] T206 [P] [US2] Provide Marketplace submission checklist and log retention template  
      `docs/lifecycle/checklists/release-checklist.md`
- [ ] T207 [US2] Update `docs/contract/rbac_manifest_spec.md` with cross references to lifecycle manifest requirements  
      `docs/contract/rbac_manifest_spec.md`

**Checkpoint**: Packaging validation runs clean on sample plugin and docs guide the flow end-to-end.

---

## Phase 5 — User Story 3 — Manage Deprecation And Sunset Transparently (Priority: P2)

**Goal**: Define lifecycle status governance, notification paths, and Marketplace integration for deprecation/sunset.

**Independent Test**: Change a manifest’s lifecycle status and follow the documented procedure to publish notices, update Marketplace, and verify host catalog behavior.

### Tasks

- [ ] T301 [US3] Document lifecycle status state machine, required effective dates, and replacement guidance  
      `docs/lifecycle/deprecation.md`
- [ ] T302 [P] [US3] Add tenant & Marketplace communication templates (email, in-app notice)  
      `docs/lifecycle/notices/deprecation-email.md`, `docs/lifecycle/notices/in-app.md`
- [ ] T303 [US3] Provide operational runbook covering host install visibility toggles and rollback policy  
      `docs/lifecycle/runbooks/deprecation-runbook.md`
- [ ] T304 [P] [US3] Publish Marketplace lifecycle OpenAPI contract and reference usage in docs  
      `docs/lifecycle/contracts/marketplace-lifecycle.openapi.yaml`, `docs/lifecycle/deprecation.md`
- [ ] T305 [US3] Update manifest sample with lifecycle block examples for active → deprecated → sunset transitions  
      `docs/lifecycle/examples/manifest-lifecycle.yaml`
- [ ] T306 [US3] Extend release checklist with deprecation review gates (permissions, compliance archiving)  
      `docs/lifecycle/checklists/release-checklist.md`

**Checkpoint**: Lifecycle transitions demonstrably propagate to Marketplace/host with documented playbook.

---

## Phase 6 — Polish & Cross-Cutting

- [ ] T901 [P] [US-All] Run editorial review (style, terminology) and align glossary across lifecycle + published docs  
      `docs/lifecycle/`, `docs/integration/01_plugin_lifecycle/`
- [ ] T902 [US-All] Update `docs/releases/` templates to reference lifecycle governance steps  
      `docs/releases/`
- [ ] T903 [US-All] Wire CI job to call new Make targets (`make verify-manifest`, `make package-pxp`) and capture artefacts  
      `.github/workflows/*.yml`
- [ ] T904 [P] [US-All] Provide sample audit trail bundle in `docs/lifecycle/examples/pxp-audit/`  
      `docs/lifecycle/examples/pxp-audit/`

---

## Dependencies & Execution Order

### Phase Dependencies

1. Setup → 2. Foundational → 3. US1 → 4. US2 → 5. US3 → 6. Polish  
   - US1 (Bootstrap) must complete before US2 (packaging) to ensure workspace standards exist.  
   - US2 must complete before US3 (deprecation) so lifecycle states reference validated manifests and packaging artefacts.

### Parallel Opportunities

- Setup: T002 and T003 can proceed after T001 establishes the lifecycle docs.  
- Foundational: T011 can run alongside T010; T012 can proceed independently once directories exist.  
- US1: T102 and T104 parallel after T101 outlines bootstrap narrative.  
- US2: T202, T204, T206 can run in parallel after T201 defines manifest mapping.  
- US3: T302 and T304 can execute concurrently post-T301.  
- Polish: T901 and T904 parallelized when story phases close.

### Story Independence

- **US1** delivers the MVP bootstrap playbook; verify by running bootstrap checklist on a clean repo.  
- **US2** depends on US1 baseline but is independently testable via package + manifest validation runs.  
- **US3** builds on US2 artefacts and is validated by executing lifecycle status changes and communications.

### Suggested MVP Scope

- Complete Phases 1–3 (through US1) to unlock onboarding consistency and validate workspace bootstrap without waiting on packaging or deprecation flows.
