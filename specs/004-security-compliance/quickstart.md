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

## 4. Exercise Consent & ToolGrant Flows
- Start services: `make run`.
- Issue consent via host CLI (stub): `scripts/dev/consent_issue.sh tenant-001 email phone`.
- Issue ToolGrant: `curl -X POST http://localhost:8086/_core/toolgrants -d '{...}'`.
- Invoke protected endpoint with returned JWT; expect `200`.
- Revoke ToolGrant via `scripts/dev/toolgrant_revoke.sh <jti>`; subsequent requests return `403` with `error_code=TOOLGRANT_REVOKED`.

## 5. Trigger Lifecycle Events
- Publish erase event: `scripts/dev/events/emit.sh data.lifecycle.erase tenant-001`.
- Confirm deletion evidence recorded in `/logs/audit.log` and `privacy_lifecycle_events`.
- Use admin UI (`/_p/<plugin-id>/admin/security/events`) to verify status.

## 6. Vulnerability Response Drill
- File mock vulnerability: `scripts/dev/security/report_vuln.sh CRITICAL PX-ADV-2025-TEST`.
- Implement hotfix branch; rebuild `.pxp`: `make build && make package`.
- Publish advisory: `scripts/dev/security/publish_advisory.sh PX-ADV-2025-TEST`.
- Confirm signed bundle stored in `dist/security/` and notifications sent (check `security_advisory_distributions`).

## 7. Operational Checklist Before Release
- [ ] `make test` and `make security-audit` green.
- [ ] Nuxt admin build passes `npm run build` and `npm run lint`.
- [ ] Manifest updated with `data_usage`, `security_baseline_version`, consent scopes.
- [ ] `.pxp` package signed; cosign verification succeeds.
- [ ] Advisory backlog empty or documented with waivers.
