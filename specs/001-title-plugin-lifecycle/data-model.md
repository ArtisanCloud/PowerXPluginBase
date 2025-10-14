# Data Model — Plugin Lifecycle Governance

## Plugin Package (`.pxp`)
- **Purpose**: Immutable bundle installed by PowerX host containing runtime assets and metadata.
- **Key Fields**:
  - `manifest.yaml` — canonical metadata (see Lifecycle Manifest)
  - `backend/` artefacts — compiled binaries or container images
  - `web-admin/.output/` — admin UI build
  - `migrations/` — database migration scripts with RLS enforcement
  - `contracts/` — API and RBAC contract definitions
  - `signature/` — detached signature plus certificate chain
- **Relationships**:
  - Contains exactly one **Lifecycle Manifest**.
  - References zero or more **Release Channel Records** during publication.
- **Validation Rules**:
  - SHA256 hash recorded in manifest checksum table.
  - All contained migrations must be idempotent and reversible (where feasible).
  - Bundle must be reproducible given tagged source state.

## Lifecycle Manifest (`manifest.yaml`)
- **Purpose**: Source of truth for plugin identity, compatibility, runtime assets, and compliance evidence.
- **Key Fields**:
  - `id`, `name`, `description`
  - `version` (SemVer), `channel` (`stable|beta|alpha|dev`)
  - `min_core`, `max_core` (optional) compatibility bounds
  - `runtime` descriptors (backend binary, container digest)
  - `frontends` list (admin URL, hash, asset manifest)
  - `migrations` metadata (sequence, checksum, reversible flag)
  - `contracts` references (API, RBAC, config schemas)
  - `rbac` roles/actions matrix
  - `config_schema`, `secrets` definitions
  - `lifecycle` block (`status`, `effective_date`, `replacement`, `notes`)
  - `signature` block (hashes, signer, certificate, timestamp)
  - `build` provenance (git commit, builder identity, CI run)
- **Relationships**:
  - Linked one-to-many with **Release Channel Records** across the Marketplace.
  - Lifecycle status changes replicated into **Lifecycle Status Ledger**.
- **Validation Rules**:
  - `version` must increment per SemVer and be unique per `channel`.
  - `min_core`/`max_core` range must overlap supported host versions.
  - `lifecycle.status` transitions must follow allowed state machine (Active → Deprecated → Sunset).
  - Signatures must match `.pxp` bundle hash table.

## Release Channel Record
- **Purpose**: Marketplace representation of a published plugin version within a specific channel.
- **Key Fields**:
  - `plugin_id`, `version`, `channel`
  - `status` (pending_review, published, revoked)
  - `submitted_at`, `published_at`
  - `dependencies` (core version range, other plugins)
  - `health_checks` summary (latest status, timestamp)
- **Relationships**:
  - References exactly one **Lifecycle Manifest** by `version`.
  - Feeds updates into **Lifecycle Status Ledger** when deprecated/sunset.
- **Validation Rules**:
  - Version uniqueness per `(plugin_id, channel)` enforced via Marketplace API.
  - Only one `published` record per channel at a time unless intentionally marked `beta`/`alpha`.
  - `health_checks` must report success before promotion to stable.

## Lifecycle Status Ledger
- **Purpose**: Historical log of lifecycle status changes for auditing and tenant communication.
- **Key Fields**:
  - `plugin_id`, `version`
  - `previous_status`, `new_status`
  - `effective_date`
  - `replacement_version` (optional)
  - `announcement_links` (docs, migration guides)
  - `recorded_by` (user/service)
- **Relationships**:
  - Receives events from **Release Channel Records** and host tooling updates.
  - Referenced by tenant notification systems and compliance auditors.
- **Validation Rules**:
  - Status transitions must follow allowed state machine; no direct `active` → `sunset`.
  - `effective_date` cannot precede `recorded_at`.
  - Entries must retain immutable history (append-only).
