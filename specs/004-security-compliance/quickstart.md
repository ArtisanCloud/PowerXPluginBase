# Quickstart – Security & Compliance Enablement

## 1. Prerequisites
- Ensure Go 1.24, Node 20, Docker (for Trivy), and `cosign` CLI available.
- Copy `backend/etc/config.example.yaml` to `backend/etc/config.yaml`; add `security` block for masking rules and baseline version.
- Export host STS credentials and JWT signing keys via Secrets Manager (DEV mode only for local runs).

## 2. Bootstrap Database Artifacts
```bash
make migrate \
  POWERX_DATABASE_URL=postgres://powerx:powerx@localhost:5432/powerx_plugin_base?sslmode=disable
```
- Confirms new tables (`privacy_*`, `security_*`) created with tenant RLS enabled.
- Seed baseline checklist version using `backend/cmd/database` seed command (run once per environment).

## 3. Run Security Audit Pipeline
```bash
make security-audit
```
- Generates reports under `build/security/<timestamp>/`.
- Inspect `report.json` and `report.sarif`; upload to Marketplace submission if applicable.
- Pipeline fails on High/Critical findings unless explicitly waived.
- Use `scripts/security/audit_export.sh` to sync signed log bundles into `dist/security/` when preparing a release.

## 4. Exercise Consent & ToolGrant Flows
- Start services: `make run`.
- Issue consent via host CLI (stub): `scripts/dev/consent_issue.sh tenant-001 email phone`.
- Issue ToolGrant: `curl -X POST http://localhost:8086/_core/toolgrants -d '{...}'`.
- Invoke protected endpoint with returned JWT; expect `200`.
- Revoke ToolGrant via `scripts/dev/toolgrant_revoke.sh <jti>`; subsequent requests return `403` with `error_code=TOOLGRANT_REVOKED`.
- Inspect active grants and revocations in the admin UI at `/_p/<plugin-id>/admin/security/toolgrants`.

## 5. Trigger Lifecycle Events
- Publish erase event: `scripts/dev/events/emit.sh data.lifecycle.erase tenant-001`.
- Confirm deletion evidence recorded in `/logs/audit.log` and `privacy_lifecycle_events`.
- Use admin UI (`/_p/<plugin-id>/admin/security/events`) to verify status.

## 6. Vulnerability Response Drill
- File mock vulnerability: `scripts/dev/security/report_vuln.sh CRITICAL PX-ADV-2025-TEST`.
- Implement hotfix branch; rebuild `.pxp`: `make build && make package-pxp`.
- Publish advisory via API or UI:
  - CLI: `scripts/dev/security/publish_advisory.sh PX-ADV-2025-TEST`
  - UI: `/_p/<plugin-id>/admin/security/advisories`
- Confirm signed bundles stored in `dist/security/<version>/` and notifications sent (check `security_advisory_distributions` table and audit events).

## 7. Operational Checklist Before Release
- [ ] `make test`, `make security-audit`, and `npm run build` succeed (document failures with waivers if blocked by tooling).
- [ ] Run `make package-pxp` to stage backend binaries, frontend bundle, and `build/security/advisories` into `dist/security/`.
- [ ] Manifest updated with `data_usage`, `security_baseline_version`, consent scopes.
- [ ] `.pxp` and advisory bundles signed (`cosign verify-blob dist/security/<version>/*.json`).
- [ ] Advisory backlog empty or documented with waivers; capture artifact hashes in release checklist.
