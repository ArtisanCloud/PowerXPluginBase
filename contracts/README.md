# Contracts Workspace

This directory houses the canonical capability descriptors and JSON Schemas that PowerX uses to install and validate the plugin.

## Layout

- `capabilities/` — YAML descriptors for each published capability (metadata, RBAC mapping, schema references).
- `schema/input/` — JSON Schema draft-07 files that describe inbound payloads consumed by the plugin.
- `schema/output/` — JSON Schema draft-07 files that describe outbound payloads produced by the plugin.

## Authoring Guidelines

- Capability IDs follow `<domain>.<resource>.<action>` and use SemVer in the descriptor `version` field.
- Keep manifest references lightweight: list capability IDs in `plugin.yaml` / `manifest.yaml`, but store the detailed definition here.
- Schemas must set `$id` to match their file path (e.g. `schema/input/foo.bar.action.v1.json`) so tooling can resolve them unambiguously.

## Next Steps

Implementation tasks for this feature will populate these folders with descriptors, schemas, and validation tooling. See `specs/002-title-plugin-capabilities/quickstart.md` for the end-to-end authoring flow and lifecycle commands (`make check-capability`, `make check-compat`) once they are introduced.
