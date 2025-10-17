# Research & Decisions – Security & Compliance

## Decision 1 – Persist tenant consent & classification in plugin schema
- **Decision**: Introduce GORM models `privacy_data_classifications`, `privacy_consent_tokens`, and `privacy_lifecycle_events` under schema `powerx_plugin_base` to track consent scope, retention windows, and erasure/export evidence per tenant.
- **Rationale**: Database-backed records guarantee auditability, allow RLS enforcement, and align with GDPR/PIPL logging obligations while keeping evidence co-located with plugin data.
- **Alternatives considered**:
  - Store consent metadata only in host services → rejected because plugins still need local enforcement context and offline evidence.
  - Persist consent in flat config/YAML → rejected due to lack of per-tenant granularity and rotation headaches.

## Decision 2 – Enforce log redaction via middleware + structured helpers
- **Decision**: Extend `internal/middleware/request_trace` with configurable field-level masks sourced from `backend/etc/security_baseline.yaml`, and wrap logrus with helper `privacy.MaskFields` to sanitize PII before emission.
- **Rationale**: Central middleware ensures all transports inherit masking without duplicating logic; config-driven rules let compliance teams update masking without code redeploy.
- **Alternatives considered**:
  - Rely on downstream log pipeline scrubbing → rejected because it cannot prevent transient exposure or satisfy “privacy by design”.
  - Scatter manual masking at call sites → rejected for maintainability and higher regression risk.

## Decision 3 – ToolGrant lifecycle backed by JWT + revocation table
- **Decision**: Use `golang-jwt/jwt/v5` to sign ToolGrant JWTs with host-provided keys, cache active grants in-memory per tenant, and persist revocations in `security_toolgrant_revocations` to enforce ≤24h TTL and immediate invalidation.
- **Rationale**: JWT keeps runtime lean, while persisted revocations support horizontal scaling and audit trace; aligns with spec’s non-sharing and traceability requirements.
- **Alternatives considered**:
  - Host introspection endpoint on every call → rejected due to latency and host coupling.
  - Long-lived opaque tokens → rejected for violating least-privilege/TTL requirements.

## Decision 4 – Security baseline automation via single Make target
- **Decision**: Implement `make security-audit` orchestrating `golangci-lint`, `govulncheck`, `gosec`, `npm audit --production`, `trivy fs` (filesystem scan), and signature verification scripts, emitting SARIF/JSON into `build/security/`.
- **Rationale**: Consolidated command enables Marketplace gating, produces machine-readable artifacts, and matches spec requirement for standardized reporting.
- **Alternatives considered**:
  - Separate ad-hoc scripts per team → rejected for inconsistent outputs and reviewer burden.
  - Outsource to external SaaS scanner → rejected due to network restrictions and reproducibility needs.

## Decision 5 – Vulnerability response workflow anchored in advisory registry
- **Decision**: Create service `internal/services/admin/security/advisory_service` managing `security_vulnerability_advisories` table, auto-generating signed advisory bundles and orchestrating event emission (`plugin.vulnerability.*`) plus marketplace webhook pushes.
- **Rationale**: Local registry ensures traceable lifecycle (report → fix → publish), enables SLA tracking, and integrates with packaging pipeline for signed `.pxp` updates.
- **Alternatives considered**:
  - Store advisories only in docs repo → rejected because it lacks automation, signatures, and runtime enforcement hooks.
  - Depend entirely on host marketplace APIs for storage → rejected since plugins must prepare metadata pre-submission and supply evidence for review.
