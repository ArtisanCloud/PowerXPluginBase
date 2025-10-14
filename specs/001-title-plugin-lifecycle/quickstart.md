# Quickstart — Applying the Plugin Lifecycle Standard

Follow this guide to adopt the lifecycle governance standard for a PowerX plugin.

## 1. Bootstrap the Workspace
1. Fork or scaffold the plugin repository and rename to `com.powerx.plugin.<slug>`.
2. Ensure the canonical structure exists:
   - `backend/`, `web-admin/`, `docs/`, `build/pxp/`
   - Copy `backend/etc/config.example.yaml` to `config.yaml` (gitignored).
3. Run `make dev-setup` then `make migrate` to initialize schema.
4. Document environment variables and migration seeds in `docs/lifecycle/bootstrap.md`.

## 2. Maintain Manifest Parity
1. Update `backend/plugin/plugin.yaml` for development metadata.
2. Mirror required fields in `backend/plugin/manifest.yaml`; run `make verify-manifest` to validate JSON Schema compliance.
3. Record RBAC, config schema, and secrets definitions alongside API contracts.

## 3. Produce a Signed `.pxp` Package
1. Build backend and admin assets: `make build && make frontend-build`.
2. Invoke `make package-pxp` to populate `build/pxp/` with artefacts, hashes, and signing receipts.
3. Archive generated audit log in `docs/lifecycle/releases/<version>.md`.

## 4. Register with the Marketplace
1. Submit `.pxp` via Marketplace API along with manifest metadata and hash proof.
2. Confirm uniqueness of `(plugin_id, version, channel)` and resolve validation feedback.
3. Publish release notes in `docs/releases/<version>.md` referencing lifecycle status.

## 5. Manage Deprecation & Sunset
1. When planning deprecation, update `manifest.yaml:lifecycle.status` and `effective_date`.
2. Call Marketplace lifecycle endpoint to broadcast status and optional replacement version.
3. Notify tenants using templated communications stored in `docs/lifecycle/notices/`.
4. After the sunset date, verify host catalog hides the version and archive artefacts per compliance retention rules.

## 6. Validate Continuously
1. Add `make verify-manifest` and `make package-pxp` to CI pipelines.
2. Run `make test`, `make frontend-test`, and migration smoke tests before promotions.
3. Review lifecycle checklists during quarterly audits to ensure ongoing conformance.

## 7. Publish Updated Lifecycle Docs
1. Run `make sync-lifecycle-docs` (or equivalent) to mirror `docs/lifecycle/` into `docs/integration/01_plugin_lifecycle/`.
2. Commit both lifecycle source files and synced integration outputs for review.
