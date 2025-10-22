# Implementation Plan: Support & Operations (Support Playbook, Incident Handling, SLA/SLO)

**Branch**: `007-support-and-operations` | **Date**: 2025-10-21 | **Spec**: [specs/007-support-and-operations/spec.md](spec.md)  
**Input**: Feature specification from `/specs/007-support-and-operations/spec.md`

---

## Summary

Deliver a full Support & Operations stack for Marketplace plugins: configure multi-channel support playbooks, orchestrate incident lifecycle with SEV gating and stakeholder comms, and publish SLA metrics with incentive/penalty automation. Backend work introduces a new `operations` bounded context (services, repositories, HTTP handlers, scheduled jobs) with normalized ticket/incident/SLA tables and webhook emission. Admin UI gains Operations pages for configuring playbooks, managing incidents, and visualizing SLA dashboards, while existing checklists/metrics components are extended for readiness gating and transparency.

---

## Technical Context

**Language/Version**: Go 1.24 (backend), Node 20 + Nuxt 4.1.3 (web-admin)  
**Primary Dependencies**: Gin HTTP router, GORM ORM, Redis (event retry cache), existing webhook_service abstraction  
**Storage**: PostgreSQL `powerx_plugin_base` schema; new tables `operations_support_*`, `operations_incidents_*`, `operations_sla_*`; migration + RLS policies required  
**Testing**: Go unit/service/integration tests under `internal/services/operations`, contract tests for `operations` HTTP routes, Nuxt component/e2e tests (`pnpm test operations`)  
**Target Platform**: PowerX managed plugin runtime on Linux (backend binary) with admin UI served via host proxy `/_p/com.powerx.plugins.base/admin/**`  
**Project Type**: Hybrid backend + admin web application  
**Performance Goals**: Admin APIs p95 < 300 ms, webhook dispatch latency < 5 s (including retry), SLA recompute job completes < 2 min for 10k tickets/month  
**Constraints**: Enforce tenant isolation (RLS + JWT), audit logging for every support/incident action, SLA incentive applied only once per period, maintain existing release packaging  
**Scale/Scope**: Expect 5k active tenants, ≤20k tickets/月, ≤200 concurrent incident subscribers, public SLA endpoint cached per plugin daily

### Platform / Hosting Integration

- **Reverse Proxy & Routes**: Admin endpoints under `/api/v1/admin/operations/**`; public read-only SLA under `/api/v1/marketplace/sla/{plugin_id}`; webhook topics reuse existing integration dispatcher.
- **Context Signing**: Reuse JWT-based middleware; operations handlers require `integration.operations.*` policies; inbound payload includes tenant & plugin identifiers for RLS.
- **Tenant/RBAC**: Introduce RBAC scopes (`operations.support:*`, `operations.incident:*`, `operations.sla:*`); repositories execute inside `BeginTenantTx` with `plugin_id` guard tables.
- **Outbound Access**: Webhook emissions leverage existing STS-signed requests through `webhook_service`; no new long-lived credentials.
- **Observability**: Extend metrics with `operations_support_ticket_events_total`, `operations_incident_active`, `operations_sla_score`; ensure `/healthz` covers new dependencies; structured logs include `incident_id` / `ticket_id`.
- **Packaging**: Update `plugin.yaml` RBAC, deliver new migrations, backend build, and web-admin `.output`; document release in changelog per SemVer.

---

## Constitution Check

- [x] **Host Contract First** {PX-HOST-001} — Routes stay within `/api/v1/**`, manifest & RBAC updated, public SLA endpoint documented.  
- [x] **Tenant Isolation & Zero Trust** {PX-CTX-001} — JWT auth + tenant scoped repos, new tables include RLS, incidents/tickets require tenant context in services & tests.  
- [x] **Service-Centric Architecture** {PX-SVC-001} — New `operations` services orchestrate logic; HTTP handlers thin, gRPC (future) shares service; repositories encapsulate DB.  
- [x] **RBAC & Least Privilege** {PX-RBAC-001} — Define `operations.support.{read,write}`, `operations.incident.{read,command}`, `operations.sla.{read,manage}`; UI enforces visibility only.  
- [x] **Observable & Testable Delivery** {PX-OBS-001} — Metrics/logs/audit coverage planned; unit + integration + contract tests + checklist gating; `/healthz` covers Redis + DB migrations.  
- [x] **Minimal Footprint & Versioned Releases** {PX-PKG-001} — Reuse existing infra (webhook_service, checklist components), minimal deps, release notes + SemVer bump to ≥0.6.0 with dist assets.

---

## Project Structure

### Documentation (this feature)

```
specs/007-support-and-operations/
├── plan.md
├── research.md
├── data-model.md
├── quickstart.md
├── contracts/
│   └── operations-openapi.yaml
└── tasks.md              # Generated via /speckit.tasks later
```

### Source Code (repository root)

```
backend/
├── internal/domain/models/operations/
│   ├── support_ticket.go
│   ├── incident.go
│   └── sla_profile.go
├── internal/domain/repository/operations/
│   ├── support_repository.go
│   ├── incident_repository.go
│   └── sla_repository.go
├── internal/domain/dto/operations/      # Optional response structs for public SLA
├── internal/services/operations/
│   ├── support_service.go
│   ├── incident_service.go
│   ├── sla_service.go
│   └── jobs/
│       └── sla_recompute_job.go
├── internal/transport/http/admin/operations/
│   ├── support_handler.go
│   ├── incident_handler.go
│   └── sla_handler.go
├── internal/transport/http/public/marketplace/
│   └── sla_handler.go
├── internal/observability/operations/
│   └── metrics.go
├── migrations/
│   └── 2025Q4_operations_support_sla.sql
└── tests/
    ├── integration/operations/
    ├── contract/operations/
    └── jobs/operations/

web-admin/
├── app/pages/_p/com.powerx.plugins.base/admin/operations/
│   ├── support.vue
│   ├── incidents.vue
│   └── sla.vue
├── app/components/operations/
│   ├── IncidentTimeline.vue
│   └── SlaScoreCard.vue
├── app/stores/operations/
│   ├── useOperationsStore.ts
│   └── useSlaStore.ts
└── tests/operations/
    ├── support_playbook.spec.ts
    ├── incident_flow.spec.ts
    └── sla_dashboard.spec.ts
```

**Structure Decision**: Adopt a dedicated `operations` bounded context across backend and admin UI, aligning with Service-Centric architecture and keeping support/incident/SLA assets isolated from existing `integration` 与 `marketplace` 域，实现清晰的依赖与发布边界。

---

## Complexity Tracking

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|--------------------------------------|
| *(None)* | — | — |

---
