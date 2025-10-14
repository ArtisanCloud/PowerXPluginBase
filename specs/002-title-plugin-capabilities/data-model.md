# Data Model — Capabilities & Schema Governance

## CapabilityDescriptor
- **Purpose**: Machine-readable definition of a plugin capability.
- **Fields**:
  - `id`: `<domain>.<resource>.<action>` string, unique per plugin.
  - `type`: Enum (`API`, `Event`, `Job`, `Tool`, `Bridge`, `Schema`).
  - `version`: SemVer string (supports parallel `v1`, `v2`).
  - `status`: `active | deprecated | sunset`.
  - `description`: Human-readable summary.
  - `rbac`: Mapping to `resource` + `actions`.
  - `provides`: List of schema IDs produced.
  - `consumes`: List of schema IDs required.
  - `metadata`: Optional map (latency SLAs, throttling hints).
- **Relationships**:
  - References multiple `SchemaDefinition` entities via `provides/consumes`.
  - Linked to `CompatibilityRecord` for version management.

## SchemaDefinition
- **Purpose**: JSON Schema describing input/output structures for capabilities, UI config, or Agents.
- **Fields**:
  - `id`: Unique path (e.g., `schema/input/tool.create.v1`).
  - `version`: SemVer string; used to evaluate compatibility.
  - `kind`: `input | output | config`.
  - `schema`: JSON Schema draft-07 payload.
  - `deprecated`: Optional date or boolean.
- **Relationships**:
  - Associated with one or more `CapabilityDescriptor` entries.
  - Participates in `CompatibilityMatrix` comparisons.

## CompatibilityRecord
- **Purpose**: Track compatibility expectations between schema/capability versions.
- **Fields**:
  - `artifact`: `capability` or `schema`.
  - `id`: Capability ID or Schema ID.
  - `version`: Current version.
  - `baseline`: Previous version used for diff.
  - `change_type`: `major | minor | patch`.
  - `notes`: Summary of detected differences (added fields, removed enums).
- **Relationships**:
  - References `CapabilityDescriptor`/`SchemaDefinition` entries.
  - Stored alongside CI reports for audit.

## RBACMapping
- **Purpose**: Ensure capability IDs map to RBAC resources/actions.
- **Fields**:
  - `capability_id`.
  - `resource`.
  - `actions`: Array of allowed operations.
  - `policy`: Optional least-privilege comments.
- **Relationships**:
  - Derived from capability YAML but verified against `docs/contract/rbac_manifest_spec.md` requirements.
