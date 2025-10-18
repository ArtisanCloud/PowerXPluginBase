# Data Model – Security & Compliance Enablement

## Entity: PrivacyDataClassification
- **Table**: `privacy_data_classifications`
- **Purpose**: Map tenant datasets to classification tier and lawful processing basis.
- **Fields**:
  - `id UUID` (PK)
  - `tenant_id TEXT` (FK -> tenant context; RLS enforced)
  - `asset_key TEXT` (unique per tenant, e.g., `customer.email`)
  - `category TEXT` (enum: `PII`, `BUSINESS`, `LOG`, `AI_INPUT`, `AI_OUTPUT`)
  - `lawful_basis TEXT` (enum: `CONSENT`, `CONTRACT`, `LEGAL_OBLIGATION`, `LEGITIMATE_INTEREST`)
  - `retention_policy JSONB` (fields: `max_duration`, `event_trigger`, `auto_purge`)
  - `purpose TEXT` (short description matching `manifest.data_usage`)
  - `created_at TIMESTAMP`
  - `updated_at TIMESTAMP`
- **Indexes**: `(tenant_id, asset_key)` unique; `(tenant_id, category)`
- **Relationships**: Referenced by `PrivacyConsentToken.scope` and `PrivacyLifecycleEvent.asset_key`.

## Entity: PrivacyConsentToken
- **Table**: `privacy_consent_tokens`
- **Purpose**: Track tenant-issued consent artifacts authorizing plugin data use.
- **Fields**:
  - `id UUID` (PK)
  - `tenant_id TEXT`
  - `consent_token TEXT` (host-issued identifier; encrypted at rest)
  - `scope JSONB` (array of `asset_key` entries referencing `PrivacyDataClassification`)
  - `expires_at TIMESTAMP`
  - `issued_at TIMESTAMP`
  - `issued_by TEXT` (host user/agent)
  - `status TEXT` (enum: `ACTIVE`, `REVOKED`, `EXPIRED`)
  - `revoked_at TIMESTAMP NULL`
  - `revoked_reason TEXT NULL`
- **Indexes**: `(tenant_id, consent_token)` unique; partial index for `status = 'ACTIVE'`.
- **Relationships**: Joined with requests via middleware; life events recorded in `PrivacyLifecycleEvent`.

## Entity: PrivacyLifecycleEvent
- **Table**: `privacy_lifecycle_events`
- **Purpose**: Evidence log for retention, erasure, export operations.
- **Fields**:
  - `id UUID` (PK)
  - `tenant_id TEXT`
  - `event_type TEXT` (enum: `RETENTION_START`, `RETENTION_PURGE`, `EXPORT`, `ERASURE`, `CONSENT_REVOKE`, `CONSENT_RENEW`)
  - `asset_key TEXT`
  - `payload JSONB` (command metadata, checksum, operator)
  - `occurred_at TIMESTAMP`
  - `recorded_by TEXT` (service identifier)
  - `status TEXT` (enum: `PENDING`, `SUCCEEDED`, `FAILED`)
- **Indexes**: `(tenant_id, event_type, occurred_at)`
- **Relationships**: Derived from consent or data operations; surfaced via admin security APIs.

## Entity: SecurityAuditReport
- **Table**: `security_audit_reports`
- **Purpose**: Persist metadata for each `make security-audit` run.
- **Fields**:
  - `id UUID` (PK)
  - `artifact_path TEXT` (`build/security/<timestamp>/report.json`)
  - `sarif_path TEXT`
  - `report_hash TEXT`
  - `initiated_by TEXT` (CI job, developer)
  - `status TEXT` (enum: `PASSED`, `FAILED`)
  - `findings JSONB` (summary counts per severity)
  - `created_at TIMESTAMP`
- **Indexes**: `(created_at DESC)`
- **Relationships**: Linked to vulnerability advisories for traceability.

## Entity: ToolGrantRevocation
- **Table**: `security_toolgrant_revocations`
- **Purpose**: Record revoked ToolGrant leases for enforcement across replicas.
- **Fields**:
  - `id UUID` (PK)
  - `tenant_id TEXT`
  - `toolgrant_id TEXT` (JWT `jti`)
  - `revoked_at TIMESTAMP`
  - `revoked_by TEXT` (actor or automation)
  - `reason TEXT`
  - `ttl_expiry TIMESTAMP` (original grant expiry)
- **Indexes**: `(tenant_id, toolgrant_id)` unique; `(ttl_expiry)` for TTL sweeper.
- **Relationships**: Consulted by middleware and admin dashboards; correlated with audit events.

## Entity: ToolGrantUsageEvent
- **Table**: `security_toolgrant_usage_events`
- **Purpose**: Audit log for issuance, renewal, consumption, and expiration.
- **Fields**:
  - `id UUID` (PK)
  - `tenant_id TEXT`
  - `toolgrant_id TEXT`
  - `event_type TEXT` (enum: `ISSUED`, `RENEWED`, `CONSUMED`, `REVOKED`, `EXPIRED`)
  - `capability TEXT`
  - `agent_id TEXT`
  - `occurred_at TIMESTAMP`
  - `metadata JSONB` (request_id, ip, user_agent)
- **Indexes**: `(tenant_id, toolgrant_id, event_type)`
- **Relationships**: Links to `ToolGrantRevocation` to provide full lifecycle view.

## Entity: VulnerabilityAdvisory
- **Table**: `security_vulnerability_advisories`
- **Purpose**: Manage vulnerability lifecycle per spec.
- **Fields**:
  - `id UUID` (PK)
  - `reference TEXT` (e.g., PX-ADV-2025-0001 or CVE)
  - `severity TEXT` (enum: `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`)
  - `status TEXT` (enum: `OPEN`, `PATCHED`, `PUBLISHED`, `CLOSED`)
  - `affected_versions TEXT[]`
  - `patched_in_version TEXT`
  - `summary TEXT`
  - `details_markdown TEXT`
  - `published_at TIMESTAMP`
  - `patched_at TIMESTAMP`
  - `closed_at TIMESTAMP`
  - `sla_deadline TIMESTAMP`
  - `created_at TIMESTAMP`
- **Indexes**: `(severity, status)`; gin index on `affected_versions`.
- **Relationships**: Associates with `SecurityAuditReport` and `AdvisoryDistribution`.

## Entity: AdvisoryDistribution
- **Table**: `security_advisory_distributions`
- **Purpose**: Track notification delivery across tenants and channels.
- **Fields**:
  - `id UUID` (PK)
  - `advisory_id UUID` (FK -> `security_vulnerability_advisories`)
  - `tenant_id TEXT`
  - `channel TEXT` (enum: `MARKETPLACE`, `EMAIL`, `WEBHOOK`)
  - `delivered_at TIMESTAMP NULL`
  - `status TEXT` (enum: `PENDING`, `DELIVERED`, `FAILED`, `ACKNOWLEDGED`)
  - `metadata JSONB`
- **Indexes**: `(advisory_id, tenant_id)`
- **Relationships**: Joined with advisories to compute acknowledgement rates (SC-006).

## Entity: SecurityBaselineChecklist
- **Table**: `security_baseline_checklists`
- **Purpose**: Store versioned baseline requirements and results.
- **Fields**:
  - `id UUID` (PK)
  - `version TEXT` (e.g., `2025.10`)
  - `controls JSONB` (list of control IDs and expectations)
  - `created_at TIMESTAMP`
  - `retired_at TIMESTAMP NULL`
- **Indexes**: `(version)` unique.
- **Relationships**: Referenced by `SecurityAuditReport` to know which baseline version executed.

## State Transitions
- **PrivacyConsentToken.status**: `ACTIVE → (REVOKED|EXPIRED)`; renewal creates new row with same `consent_token`.
- **PrivacyLifecycleEvent.status**: `PENDING → SUCCEEDED` or `FAILED`; failures trigger retries and host notifications.
- **ToolGrantRevocation**: remains until `ttl_expiry` passed and sweeper deletes.
- **VulnerabilityAdvisory.status**: `OPEN → PATCHED → PUBLISHED → CLOSED`; re-open triggered if regression detected (optional transition `PUBLISHED → OPEN`).
- **AdvisoryDistribution.status**: `PENDING → DELIVERED → ACKNOWLEDGED`; `FAILED` loops to re-delivery queue.

## Validation Rules
- Consent tokens must reference only asset keys present in `PrivacyDataClassification`.
- Lifecycle events require `payload.checksum` to match hashed dataset snapshot for erasure/export verification.
- ToolGrant revocations must include reason and actor for audit acceptance.
- Vulnerability advisories cannot close until all associated distributions marked `ACKNOWLEDGED` or manually waived by Security.
