# Plugin Lifecycle Documentation

> Source-of-truth documentation for PowerX plugin lifecycle governance. Files under this directory are synchronized into `docs/integration/01_plugin_lifecycle/` via `make sync-lifecycle-docs`.

## Contents

- [overview.md](./overview.md) — end-to-end lifecycle narrative and entry points
- [checklists/](./checklists/) — bootstrap, release, and deprecation checklists
- [examples/](./examples/) — canonical metadata snapshots (`plugin.yaml`, `manifest.yaml`)
- [contracts/](./contracts/) — machine-readable schemas and Marketplace API definitions
- [notices/](./notices/) — communication templates for tenants and marketplace operators
- [runbooks/](./runbooks/) — operational procedures for lifecycle events

Keep edits in this directory; run `make sync-lifecycle-docs` to publish curated copies into integration docs.
