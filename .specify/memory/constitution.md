# PowerX Plugin Base Constitution

<!--
Sync Impact Report
Version change: N/A → 1.0.0
Modified principles: Initial publication
Added sections: Core Principles; Operational Constraints; Development Workflow & Quality Gates; Governance
Removed sections: None
Templates requiring updates:
- ✅ .specify/templates/plan-template.md (Constitution Check aligned)
- ✅ .specify/templates/tasks-template.md (Foundational work tightened)
Follow-up TODOs: None
-->

## Core Principles

### I. Host Contract First

- All plugin interfaces MUST conform to the PowerX reverse proxy contract: expose business APIs under `/v1`, surface management endpoints at `/api/v1/admin/{manifest,rbac}`, and keep `plugin.yaml` in sync with the runtime manifest.
- Outbound calls to PowerX MUST use STS-issued credentials; direct coupling to host internals is prohibited.
Citing the host contract first keeps every plugin drop-in compatible with the Platform Router and avoids regressions when PowerX upgrades its routing fabric.

### II. Tenant Isolation & Zero Trust

- Every inbound request MUST validate the PowerX context (JWT or HMAC) before touching application state; `POWERX_DEV_MODE` is limited to local development.
- The data model MUST carry `tenant_id` and Postgres Row Level Security; repositories MUST execute inside `BeginTenantTx` with `SET LOCAL app.tenant_id`.
- Secrets, tokens, and database roles MUST follow least privilege and be rotated via STS or environment management.
Zero trust enforcement is non-negotiable because the plugin frequently runs alongside other tenants within the same cluster and shares host infrastructure.

### III. Service-Centric Architecture

- Transport handlers stay thin: validate input, delegate to services, and translate responses; business orchestration lives exclusively in `internal/services`.
- Repositories encapsulate data access and never leak GORM specifics to services; HTTP and gRPC layers MUST share the same service logic.
- Shared dependencies (config, logger, clients) travel through the application container to keep construction deterministic and testable.
This layering preserves clear failure domains, encourages reuse between protocols, and allows automated testing at each boundary.

### IV. Observable & Testable Delivery

- All features MUST provide structured logging with request IDs, health probes, and metrics hooks needed by PowerX observability.
- Each change MUST ship with automated coverage matching scope: unit tests for services, integration tests for multi-tenant flows, and migration smoke tests when schemas change.
- Migrations MUST be idempotent, reversible, and guarded behind explicit `POWERX_RUN_MIGRATE` toggles.
Evidence-first delivery ensures plugins remain diagnosable after deployment and protects multi-tenant data paths from silent drift.

### V. Minimal Footprint & Versioned Releases

- Keep dependencies minimal, prefer Go/Nuxt stack already curated by the template, and remove dormant code before release.
- Every release MUST document runtime configuration, update manifests, and package artifacts via `make release && make package-release` (or equivalent CI task).
- Breaking changes to APIs, data contracts, or host expectations MUST trigger semantic version bumps and migration guidance.
Lean releases reduce attack surface, simplify audits, and make upgrades safe for operators managing dozens of plugins.

## Operational Constraints

- **Database**: Use Postgres ≥ 13 with a dedicated schema per plugin; enforce RLS and schema migrations through the provided tooling.
- **Runtime**: Production deployments MUST disable `POWERX_DEV_MODE`, configure `POWERX_CTX_*` issuers/audiences, and expose the service on `POWERX_BIND_ADDR`.
- **Networking**: PowerX reverse proxy mounts `/_p/<plugin-id>/admin/*` to `web-admin/.output` and `/_p/<plugin-id>/api/*` to backend `/v1/...`; SDKs and frontends MUST respect this prefix.
- **Secrets & Credentials**: Interactions with PowerX APIs require STS token exchange (`/_p/_internal/sts/exchange`); long-lived credentials are forbidden.
- **Frontend**: When `web-admin` is shipped, the Nuxt runtime MUST derive its base URL and API base from `runtimeConfig.public.apiBaseUrl` to adapt between standalone and proxied modes.

## Development Workflow & Quality Gates

- **Spec → Plan → Tasks**: Each feature starts with a spec capturing independently testable user stories, followed by an implementation plan that passes the Constitution Check, and finally tasks grouped per story to preserve MVP slices.
- **Gate Reviews**: Before implementation, reviewers confirm host contract compliance, tenant isolation coverage, and planned tests; completion reviews verify observability signals and migration discipline.
- **Testing Strategy**: Run `make test`, migration smoke tests, and (when applicable) Nuxt lint/build in CI; no change merges without green automation for its scope.
- **Release Readiness**: Deliveries include updated `plugin.yaml`, manifest/RBAC endpoints, version bumps, and documentation updates in `docs/`.
- **Incident Handling**: On regressions, roll forward with fixes that add missing coverage; rollbacks must preserve schema compatibility.

## Governance

- This constitution supersedes conflicting guidance; deviations require an RFC reviewed by PowerX Core maintainers and recorded in `docs/references/changelog.md`.
- Amendments follow semantic versioning: MAJOR for contract-breaking shifts, MINOR for new principles or mandatory workflow changes, PATCH for clarifications. Each amendment notes rationale and required migrations.
- Compliance is enforced during code review and release sign-off; reviewers MUST document how each principle was satisfied or justified.
- The Docs team ensures downstream templates stay synchronized; any TODOs introduced by amendments carry assigned owners and due dates.

## G8: UI Layer Definition (Optional)

- ID: PX-FE-001
- The term **frontend** in templates refers generically to any UI layer under the plugin:
  `web-admin/`, `web-app/`, `mini-app/`, `mobile-app/`, etc.
- Each project must define its active UI layers in `plan.md` → Project Structure section.

**Version**: 1.0.0 | **Ratified**: 2025-10-11 | **Last Amended**: 2025-10-11
