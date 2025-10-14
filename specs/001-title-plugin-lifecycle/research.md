# Research Findings — Plugin Lifecycle Governance

## Decision 1: Centralize lifecycle standard in `scaffold/lifecycle/`
- **Decision**: Author the end-to-end lifecycle standard in `scaffold/lifecycle/`, anchored by `scaffold/lifecycle/overview.md` and supporting checklists, then publish curated copies into `docs/integration/01_plugin_lifecycle/` via a sync routine.
- **Rationale**: Keeping the guidance in a scaffold namespace preserves a single source of truth while letting integration docs consume generated outputs for reviewers.
- **Alternatives considered**:
  - Embed guidance directly in `backend/plugin/manifest.yaml` comments — rejected because comments are lossy and not easily referenced by Marketplace reviewers.
  - Publish only in external wiki — rejected due to audit traceability and offline CI review requirements.

## Decision 2: Enforce manifest <-> plugin metadata parity via JSON Schema validation
- **Decision**: Define a JSON Schema representing `manifest.yaml` requirements and add a `make verify-manifest` step that validates both `plugin.yaml` and release `manifest.yaml` against synchronized constraints.
- **Rationale**: JSON Schema offers deterministic validation, integrates with Go tooling, and supports CI enforcement without binding to a single programming language.
- **Alternatives considered**:
  - Hand-written Go validation — rejected for higher maintenance cost and lack of machine-readable spec for Marketplace.
  - Rely on Marketplace API validation alone — rejected because feedback would only arrive post-upload, slowing release cycles.

## Decision 3: Standardize `.pxp` packaging through an orchestrated Make target
- **Decision**: Extend existing `make release` pipeline with a dedicated `make package-pxp` target that stages artefacts in `build/pxp/`, calculates SHA256 hashes, signs payloads, and emits an audit manifest.
- **Rationale**: Make targets are already the canonical automation entry point; isolating packaging steps enables CI reuse and manual invocations while logging every output.
- **Alternatives considered**:
  - Introduce a standalone CLI (e.g., `powerx-cli`) as the only path — rejected to avoid toolchain sprawl and honor the spec’s tool-agnostic promise.
  - Perform signing in ad-hoc scripts per team — rejected due to inconsistent logging and higher security risk.

## Decision 4: Codify lifecycle status propagation through Marketplace APIs
- **Decision**: Document a required sequence for setting `lifecycle.status` transitions (active → deprecated → sunset) and mandate Marketplace API calls that broadcast effective dates alongside tenant notifications.
- **Rationale**: Explicit propagation steps guarantee that hosts and Marketplace catalog stay in sync, reducing surprise installs of deprecated versions.
- **Alternatives considered**:
  - Manual spreadsheet tracking — rejected because it lacks automated enforcement and cannot trigger host visibility updates.
  - Blocking status changes until all tenants upgrade — rejected; instead, the standard provides exception handling guidance within the host tooling.
