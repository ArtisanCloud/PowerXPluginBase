# Implementation Plan: Security & Compliance (Privacy, ToolGrant, Baseline, Vulnerability Response)

**Branch**: `[004-security-compliance]` | **Date**: 2025-10-17 | **Spec**: specs/004-security-compliance/spec.md  
**Input**: Feature specification from `specs/004-security-compliance/spec.md`

---

## Summary

Deliver tenant-safe security and compliance guardrails for the PowerX base plugin by implementing data privacy enforcement (classification, consent, retention), a security baseline checklist with automated scans and sandbox validation, ToolGrant lifecycle middleware with auditability, and a host-aligned vulnerability response process that can ship signed patches inside mandated SLAs.

---

## Technical Context

**Language/Version**: Go 1.24 (backend runtime) + Node 20 / Nuxt 4 (web-admin)  
**Primary Dependencies**: gin-gonic (HTTP transport), gorm (PostgreSQL ORM), logrus (structured logging), golang-jwt/jwt/v5 (token handling), go-sqlmock/testify (testing), Node/Nuxt security middlewares (helmet-equivalent via Nitro hooks)  
**Storage**: PostgreSQL schema `powerx_plugin_base`; new tables for consent tokens, audit evidence, ToolGrant revocation cache, and vulnerability advisories with migrations managed under `backend/internal/db/migrations`  
**Testing**: `go test ./...`, targeted service/middleware unit tests, migration smoke via `make test`, CLI-driven `make security-audit` (new) plus Nuxt lint/build for admin UX  
**Target Platform**: PowerX managed container runtime on Linux with STS-sealed outbound access, Marketplace distribution  
**Project Type**: Multi-tier plugin (backend APIs + admin UI)  
**Performance Goals**: ToolGrant middleware adds <5% latency to protected routes; security audit pipeline completes <30 minutes; vulnerability notifications propagate within 1 hour  
**Constraints**: ≤24h ToolGrant TTL, per-tenant schema isolation, outbound traffic restricted to host gateway allowlist, logs must exclude raw PII, signed `.pxp` artifacts only  
**Scale/Scope**: Designed for 1000+ active tenants, dozens of agents per tenant, audit logs retained 365 days, marketplace submissions multiple times per week

### Platform / Hosting Integration (optional)

- **Reverse Proxy & Routes**: Business APIs continue under `/_p/<plugin-id>/api/v1/**`; new compliance endpoints exposed via `/_p/<plugin-id>/api/v1/admin/security/**` for audit artifacts; admin UI served at `/_p/<plugin-id>/admin/security`.  
- **Context Signing**: All inbound calls validated via JWT (existing middleware); consent and ToolGrant events require tenant, agent, capability claims; dev bypass remains DEV-mode only.  
- **Tenant/RBAC**: Tenant context injected via middleware sets `SET LOCAL app.tenant_id`; new repositories enforce tenant filters and RLS policies; ToolGrant scopes map to RBAC capabilities before issuance.  
- **Outbound Access**: Integrations (e.g., advisory push, consent event bus) call host services via STS short-lived credentials with scoped permissions.  
- **Observability**: Extend structured logging with consent_token_id, toolgrant_id fields; emit metrics (`security.toolgrant.revocations`, `security.audit.failures`) and events (`plugin.vulnerability.*`).  
- **Packaging**: `.pxp` build gains SARIF/JSON scanner outputs, signed advisory bundles, and updated manifest `data_usage`, `security_baseline` metadata to satisfy marketplace gate.

---

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- [x] **Host Contract First** {PX-HOST-001}  
  Compliance APIs continue under `/v1/**` with admin endpoints namespaced; manifest updates track new security metadata to keep host contract aligned.
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001}  
  Consent, ToolGrant, and audit repositories operate inside tenant-scoped schema with middleware-enforced JWT verification and RLS-backed migrations; tests will cover revocation/consent edge cases.
- [x] **Service-Centric Architecture** {PX-SVC-001}  
  Add dedicated services (`internal/services/admin/security`, `internal/services/agent/toolgrant`) while keeping HTTP handlers thin and reusing services across HTTP/gRPC interceptors.
- [x] **RBAC & Least Privilege** {PX-RBAC-001}  
  ToolGrant issuance derives from RBAC capability registry; new admin surfaces restricted to `security.manage` role; storage enforces least privilege per tenant.
- [x] **Observable & Testable Delivery** {PX-OBS-001}  
  Plan adds structured logging fields, metrics, audit trails, and unit/integration tests for middleware, services, and migrations plus CI hooks for security audit.
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001}  
  Dependencies limited to security/token helpers; release artifacts include new scanner outputs and signed advisories with SemVer minor bump to advertise compliance capabilities.

---

## Project Structure

### Documentation (this feature)

```

specs/004-security-compliance/
├── plan.md              # Implementation plan (this file)
├── research.md          # Phase 0 decisions and references
├── data-model.md        # Phase 1 entity design
├── quickstart.md        # Phase 1 operational runbook
├── contracts/           # Phase 1 API/security contract docs
└── tasks.md             # Phase 2 execution slices (future)

```

### Source Code (repository root)

```

backend/
├── cmd/
│   ├── plugin/
│   └── database/
├── internal/
│   ├── domain/models/
│   │   ├── privacy/
│   │   ├── security/
│   │   └── tool_grant/
│   ├── domain/repository/
│   │   ├── privacy/
│   │   ├── security/
│   │   └── toolgrant/
│   ├── services/
│   │   ├── admin/security/
│   │   └── agent/tool_grant/
│   ├── middleware/
│   │   ├── consent_guard/
│   │   └── tool_grant_verifier/
│   ├── transport/
│   │   ├── http/
│   │   │   ├── admin/security/
│   │   │   └── agent/security/
│   │   └── grpc/interceptors/
│   ├── shared/app/
│   └── observability/security/
├── etc/
│   └── security_baseline.yaml
└── tests/security/

web-admin/
├── app/
│   ├── pages/security/
│   ├── components/security/
│   └── composables/useSecurityBaseline.ts
└── tests/security/

```

**Structure Decision**: Extend existing backend domain/service/middleware layers with dedicated privacy, security, and ToolGrant modules while surfacing admin UI workflows under `web-admin/app/pages/security/` for audit dashboards and vulnerability advisories.

---

## Complexity Tracking

No constitution deviations anticipated; table intentionally left empty.
