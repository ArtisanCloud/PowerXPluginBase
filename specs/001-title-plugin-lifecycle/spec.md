# Feature Specification: Plugin Lifecycle Governance

**Feature Branch**: `[001-title-plugin-lifecycle]`  
**Created**: 2025-10-14  
**Status**: Draft  
**Input**: User description: "Title: Plugin Lifecycle (Init → Manifest → Versioning/Publishing → Deprecation/Sunset); WHAT/WHY: 定义 PowerX 插件从创建、清单规范、版本与发布，到弃用/退役的全生命周期标准，确保构建产物（.pxp）可验证、可回滚、可审计，并与 Marketplace/宿主保持一致行为。 Scope: 初始化与项目骨架；清单与元数据；版本与发布；弃用与退役。 Out-of-Scope: 运行时/端口/可观测性与配额；协议与外部集成；安全与隐私细则；商业上架/定价/分析。 Dependencies/Assumptions: 交付物为可签名的 .pxp 包；开发以 Makefile/脚本为主；版本遵循 SemVer 并在 Marketplace/宿主同步。"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Bootstrap A Compliant Plugin Workspace (Priority: P1)

Plugin engineers need a step-by-step standard to fork or rename a template, establish the required backend, frontend, and docs directories, run local services, and prepare build folders so their project is audit-ready from day one.

**Why this priority**: Without a shared bootstrap contract, teams create divergent structures that break downstream automation and slow reviews.

**Independent Test**: A reviewer can onboard to a new plugin, follow the documented bootstrap steps, and confirm the workspace meets the mandatory directory, migration, and build-prep checklist without referencing tribal knowledge.

**Acceptance Scenarios**:

1. **Given** a newly forked plugin repository, **When** the engineer applies the lifecycle standard, **Then** the repository contains the prescribed `backend/`, `web-admin/`, `docs/`, and `build/pxp/` layout with environment setup steps captured.
2. **Given** a contributor joining mid-project, **When** they follow the standard’s local run and migration directions, **Then** they can start backend services and seed the database without requesting unpublished instructions.

---

### User Story 2 - Release A Verifiable .pxp Package (Priority: P1)

Release managers need a canonical mapping from development-time `plugin.yaml` metadata to shipping `manifest.yaml`, with guidance for versioning, signing, and marketplace registration so every package is traceable and host-installable.

**Why this priority**: Invalid manifests or unsigned packages block deployments and force emergency rebuilds across teams.

**Independent Test**: A release manager can walk through the standard to produce a .pxp bundle, validate its hash and signature, and register it with the Marketplace without relying on external docs.

**Acceptance Scenarios**:

1. **Given** a plugin ready for release, **When** the manager builds the .pxp following the standard, **Then** the artifact includes the mandated manifest fields, hashed payloads, and signature evidence required by the Marketplace intake checklist.
2. **Given** a marketplace API submission, **When** the standard’s validation pipeline runs, **Then** it rejects duplicated version numbers or incompatible `min_core` ranges before the package is published.

---

### User Story 3 - Manage Deprecation And Sunset Transparently (Priority: P2)

Lifecycle owners must communicate deprecation timelines, recommended replacements, and install visibility rules so tenants, Marketplace operators, and hosts transition safely.

**Why this priority**: Without a coordinated sunset process, tenants may run unsupported plugins, causing compliance gaps and emergency escalations.

**Independent Test**: A lifecycle owner can mark a release as deprecated, trigger notices, and observe host behavior changes using only the standard.

**Acceptance Scenarios**:

1. **Given** an active plugin version slated for retirement, **When** the lifecycle status is updated per the standard, **Then** Marketplace listings and host catalogs show deprecation messaging with links to the migration guide inside one business day.
2. **Given** a plugin marked as sunset with an effective date, **When** a host attempts a fresh install after that date, **Then** installation is blocked or warns according to the prescribed visibility policy.

---

### Edge Cases

- Manifest declares dependencies on a core version range that conflicts with published compatibility data.
- Package signature or hash cannot be verified because of mismatched signing keys or tampered artifacts.
- Rollback requests target a package whose migrations introduce non-reversible schema changes.
- Deprecation date arrives while tenants still have pending upgrades, requiring exception handling.
- Marketplace receives simultaneous submissions for the same version from different branches.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The lifecycle standard MUST define the canonical bootstrap flow, covering repository naming, directory layout (`backend/`, `web-admin/`, `docs/`, `build/pxp/`), environment variables, and migration seeds required before shared tooling runs.
- **FR-002**: The standard MUST document how development-time `plugin.yaml` fields map to release-time `manifest.yaml`, including required metadata (`id`, `name`, `version`, `channel`, `min_core`, `runtime`, `frontends`, `migrations`, `contracts`, `rbac`, `config_schema`, `secrets`, `lifecycle`, `signature`, `build`).
- **FR-003**: The standard MUST prescribe validation checkpoints to ensure `plugin.yaml` and `manifest.yaml` stay synchronized during each release cycle.
- **FR-004**: The standard MUST define the .pxp packaging process end-to-end, including artifact composition, deterministic build ordering, SHA256 hashing, signature generation, and location of logs for audit.
- **FR-005**: The standard MUST codify versioning rules aligned with SemVer and release channels (`stable`, `beta`, `alpha`, `dev`), including eligibility criteria, promotion rules, and rollback constraints.
- **FR-006**: The standard MUST specify Marketplace registration requirements (uniqueness, dependency checks, status propagation, health report submission) before a version is approved for listing.
- **FR-007**: The standard MUST outline host installation, upgrade, and rollback prerequisites, including pre-flight checks, migration execution controls, health verification, and artifact retention.
- **FR-008**: The standard MUST define deprecation and sunset workflows: allowed lifecycle statuses (`active`, `deprecated`, `sunset`), notification cadences, required migration documentation, and visibility or access changes.
- **FR-009**: The standard MUST require audit logging and traceability artifacts (hash manifests, signature receipts, release notes) to support compliance reviews and incident response.
- **FR-010**: The standard MUST cite dependencies on Makefile or scripting automation while remaining tool-agnostic enough to allow alternative orchestrations that meet the same checks.

### Key Entities *(include if feature involves data)*

- **Plugin Package (.pxp)**: Immutable distribution bundle containing runtime binaries, frontends, migrations, contracts, and manifest metadata used by hosts and the Marketplace.
- **Lifecycle Manifest**: Release-time `manifest.yaml` describing identity, compatibility, release channel, cryptographic evidence, and operational hooks for each version.
- **Release Channel Record**: Marketplace entry tying a version to its channel, status, dependencies, and audit trail events.
- **Lifecycle Status Ledger**: Historical record of status changes (active → deprecated → sunset) with effective dates, communication artifacts, and rollback availability.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of newly bootstrapped plugins audited in quarterly reviews conform to the prescribed directory, configuration, and migration readiness checklist.
- **SC-002**: 100% of .pxp packages promoted to `stable` pass automated hash, signature, and manifest-field validation on the first submission.
- **SC-003**: Marketplace release reviews complete within one business day for compliant submissions, measured over a trailing three-month window.
- **SC-004**: Deprecation notices reach all affected tenants and Marketplace listings at least 30 days before sunset, verified across sampled releases each quarter.
- **SC-005**: No critical incidents triggered by incompatible upgrades attributable to missing lifecycle documentation over the first two release cycles post adoption.

## Assumptions

- Marketplace APIs and host tooling can surface lifecycle statuses, notices, and migration guides without additional feature work.
- Signing infrastructure is available to release managers and issues verifiable receipts per package.
- Teams adopting the standard have access to baseline Makefile or equivalent automation to run validation orchestrations.
