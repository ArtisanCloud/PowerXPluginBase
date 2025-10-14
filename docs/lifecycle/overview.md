# PowerX Plugin Lifecycle — Source Scaffold

This directory captures the canonical lifecycle standard for PowerX plugins from project bootstrap, through manifest governance and release packaging, to deprecation and sunset operations. Downstream integration docs are generated from this directory.

## Lifecycle Phases

1. **Bootstrap** — Repo fork/rename, required directory layout, environment preparation, and migration readiness. See [`bootstrap.md`](./bootstrap.md) and [`checklists/bootstrap-checklist.md`](./checklists/bootstrap-checklist.md).
2. **Manifest & Metadata** — Mapping `plugin.yaml` fields to release `manifest.yaml`, schema validation, and RBAC/config provenance. Refer to [`manifest-mapping.md`](./manifest-mapping.md) and JSON Schema in [`contracts/manifest.schema.json`](./contracts/manifest.schema.json).
3. **Packaging & Publishing** — Deterministic `.pxp` builds, hash + signature capture, Marketplace submission, and audit trails. Guidance stored in [`package.md`](./package.md).
4. **Deprecation & Sunset** — Lifecycle status machine, communication templates, and host visibility controls. Procedures live in [`deprecation.md`](./deprecation.md), notices, and runbooks.

## Operational Workflow

- Author or update documents in `docs/lifecycle/`.
- Run `make sync-lifecycle-docs` to mirror curated copies into `docs/integration/01_plugin_lifecycle/` for reviewers and Marketplace operators.
- Commit both the lifecycle sources and synced outputs to keep history aligned.

## References

- Example metadata: [`examples/plugin.yaml`](./examples/plugin.yaml) and [`examples/manifest.yaml`](./examples/manifest.yaml)
- Communication templates: [`notices/`](./notices/)
- Operational runbooks: [`runbooks/`](./runbooks/)

> 📌 **Note**: The lifecycle directory is authoritative. Do not edit the generated integration docs directly; instead, modify the sources here and re-run the sync task.
