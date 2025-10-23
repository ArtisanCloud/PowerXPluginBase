# Feature Specification: Dev Console & Admin UI (Plugin Admin, Audit & History, Troubleshooting)

**Feature Branch**: `008-dev-console-admin-ui`  
**Created**: 2025-10-22  
**Status**: Draft  
**Input**: User description: "Title: Dev Console & Admin UI (Plugin Admin, Audit & History, Troubleshooting)\nWHAT/WHY: Provide a consistent Admin Console experience for plugin operators: configuration surfaces, audit/history views, and common troubleshooting tools.\nScope: Console navigation slots & URLs (/_p/<plugin_id>/admin/...); config forms & validation; audit log query & export; job/task history; safe-ops actions (replay, retry, drain); common troubleshooting (health, quota, webhook delivery); RBAC guard integration & permission codes.\nOut-of-Scope: End-customer frontends; theming/branding system beyond required assets.\nDependencies/Assumptions: PowerX Admin (Nuxt) extension points; Reverse-proxy & auth cookies shared per trusted mode; Permission codes aligned with IAM; Observability read APIs exposed."

## Overview

- Deliver a first-party feeling control plane for PowerX plugin operators covering configuration, visibility, and low-risk operations.  
- Ensure every plugin exposes a predictable admin URL pattern, navigation structure, and guard rails that align with platform RBAC.  
- Centralize operational evidence (audit history, task/job runs, webhook delivery) and day-to-day troubleshooting affordances within the console.

## Scope

### In Scope

- Admin console navigation surfaces under `/_p/<plugin_id>/admin/...` with consistent breadcrumbing and slot usage.  
- Configuration forms with inline validation, change history, and safe application of updates to plugin runtime.  
- Audit log browsing, filtering, and export for plugin-scoped admin actions.  
- Job/task execution history, including status, initiator, timestamps, and retry outcomes.  
- Safe-operations actions (replay, retry, drain/disable) with confirmation, impact summary, and eligibility checks.  
- Troubleshooting dashboard exposing health checks, quota consumption, webhook delivery status, and operator guidance.  
- RBAC integration and permission code definitions to ensure only authorized administrators access each capability.

### Out of Scope

- Customer-facing plugin experiences or marketplace storefront changes.  
- Broad theming or branding systems beyond assets required for console framing.  
- New observability signal creation (feature reuses existing metrics and webhook traces).

### Dependencies & Assumptions

- Admin extension points within PowerX are available for plugin modules to register navigation and routes.  
- Platform reverse-proxy shares authenticated sessions with trusted admin modules to avoid duplicate sign-in flows.  
- IAM permission codes will be minted and provisioned alongside plugin installation workflows.  
- Observability APIs already expose health, quota, and webhook delivery data needed for dashboards.  
- Existing integration services handle outbound webhooks; this feature only surfaces diagnostics.  
- Plugin runtime exposes idempotent endpoints for replay, retry, and drain operations.

## Clarifications

### Session 2025-10-22

- Q: 审计日志导出时，默认的文件格式与交付方式应该如何规定？ → A: 控制台下载时允许选择 CSV 或 JSON。
- Q: 安全运维动作（replay/retry/drain）执行时要默认作用于哪个范围？ → A: 操作前需明确选择目标租户或环境，默认不跨租户/环境。
- Q: 故障排查面板中展示的健康与配额等指标，刷新频率应如何设定？ → A: 默认每 5 分钟自动刷新，可手动强制刷新。

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Configure & Govern Plugin Admin Console (Priority: P1)

An Enterprise Plugin Operator needs a predictable admin entrypoint to review plugin status, adjust configuration fields, and apply changes with validation and audit trails.

**Why this priority**: Configuration governance is the primary task administrators perform; without it the console fails to deliver core operational value and onboarding stalls.

**Independent Test**: Navigate to the plugin admin URL, update configuration with validation feedback, apply changes, and verify the action is logged with RBAC enforcement.

**Acceptance Scenarios**:

1. **Given** an operator with `operations.plugin.admin` permission, **When** they open `/_p/<plugin_id>/admin` from the main console, **Then** the plugin dashboard renders with consistent navigation and plugin metadata.  
2. **Given** editable configuration sections, **When** the operator submits a change that violates validation rules, **Then** the form blocks the save, highlights fields with guidance, and records no change in history.  
3. **Given** a valid configuration update, **When** the operator confirms the change, **Then** the console applies it, stamps the update to audit history, and surfaces the resulting state to other viewers.

---

### User Story 2 - Inspect Audit & Activity History (Priority: P2)

Compliance reviewers need to trace who changed what within the plugin, filter events by date, actor, or action type, and export records for quarterly reviews.

**Why this priority**: Transparent auditability is critical for regulated tenants and underpins incident response and governance obligations.

**Independent Test**: Access the audit history tab, filter events for a timeframe, confirm results match expected entries, and export the set for offline review.

**Acceptance Scenarios**:

1. **Given** a populated audit log, **When** a reviewer filters by actor and action type, **Then** only matching events display with timestamps and context.  
2. **Given** a need for external review, **When** the reviewer requests an export for a date range and selects CSV or JSON, **Then** the console delivers the chosen format as a direct download with contents matching the on-screen results.  
3. **Given** an unauthorized user lacking audit permissions, **When** they attempt to open the audit tab, **Then** access is denied with guidance to request appropriate roles.

---

### User Story 3 - Troubleshoot Jobs & Webhook Delivery (Priority: P3)

Support engineers require a consolidated view of job/task runs, replay/retry controls, health probes, quota status, and webhook delivery traces to resolve issues without leaving the console.

**Why this priority**: Rapid troubleshooting reduces downtime and support load, but it builds on top of the configuration and audit baselines.

**Independent Test**: From the troubleshooting area, review the latest job runs, execute a safe retry, confirm health and quota indicators render, and inspect webhook delivery details for a specific event.

**Acceptance Scenarios**:

1. **Given** prior job executions, **When** an engineer opens the job history list, **Then** runs display with status, trigger source, duration, and links to detailed context.  
2. **Given** a failed job that is eligible for retry, **When** the engineer initiates a retry, **Then** the console requires confirmation, validates safeguards, and records the outcome.  
3. **Given** webhook delivery diagnostics, **When** the engineer inspects a tenant’s webhook stream, **Then** the console shows recent attempts, response codes, and suggested remediation steps when failures persist.

### Edge Cases

- Console accessed for a plugin that is disabled or lacks required configuration sections should clearly state the limitation and block edits.  
- Audit export requests covering very large ranges must communicate progress, chunk results if needed, and guard against timeouts.  
- Safe-ops actions initiated while an identical action is already in progress should surface a lock notice instead of queueing duplicates.  
- Troubleshooting dashboards should gracefully degrade when upstream observability data is delayed or unavailable, signalling staleness instead of rendering blanks.  
- RBAC checks must prevent mixed-role sessions (e.g., cached cookies) from accidentally inheriting prior privileges.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The admin console MUST expose a stable entrypoint at `/_p/<plugin_id>/admin` with consistent navigation, breadcrumbs, and plugin context headers.  
- **FR-002**: The console MUST render configuration sections with descriptive copy, current values, pending change indicators, and inline validation messaging.  
- **FR-003**: The system MUST enforce role-based permissions for console access, configuration edits, audit review, job controls, and troubleshooting views, returning actionable error messaging when access is denied.  
- **FR-004**: Configuration saves MUST capture actor, timestamp, change summary, and prior values within an immutable audit trail.  
- **FR-005**: The console MUST allow reviewers to filter audit events by time range, actor, action type, and outcome, and to export the filtered set via direct download with an explicit choice between CSV or JSON format.  
- **FR-006**: The system MUST surface job and task history with status, triggering source, execution duration, and links to retries or related artefacts.  
- **FR-007**: Safe-operation actions (replay, retry, drain/disable) MUST require the operator to select the target tenant or environment, and present eligibility checks, impact summaries, confirmation prompts, and result notifications before executing against plugin runtime.  
- **FR-008**: Troubleshooting dashboards MUST display current health check status, quota consumption, and webhook delivery metrics sourced from observability APIs, auto-refresh every 5 minutes, and show both last refresh timestamp and manual refresh control.  
- **FR-009**: The console MUST provide drill-down views for webhook deliveries showing attempts, response codes, payload identifiers, and recommended remediation guidance.  
- **FR-010**: All console actions MUST log to audit history with permission code references to support compliance reviews.  
- **FR-011**: The system MUST surface contextual help and runbooks within the console for each troubleshooting section to guide operators toward resolution.  
- **FR-012**: Navigation slots MUST allow deep-linking to specific tabs or sections via URL parameters while respecting RBAC and providing safe fallbacks.

### Key Entities *(include if feature involves data)*

- **PluginAdminConsole**: Represents the admin experience for a specific plugin, including URL, navigation structure, and registered capabilities.  
- **ConfigurationChange**: Captures configuration section identifier, previous value snapshot, new value summary, actor, timestamp, and validation outcome.  
- **AuditEvent**: Records scoped admin actions with actor, action type, target resource, outcome, and permission code.  
- **JobRun**: Describes scheduled or manual executions with trigger source, timing, status, retry history, and linked artefacts.  
- **TroubleshootingSignal**: Aggregates health checks, quota metrics, and webhook delivery status for a plugin at a specific refresh point.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 90% of pilot plugin operators complete a configuration update within 3 minutes of reaching the console during UAT walkthroughs.  
- **SC-002**: 100% of unauthorized access attempts to audit, job control, or troubleshooting views are blocked and logged with correct permission codes in staging validation.  
- **SC-003**: Audit log queries for a 30-day window return visible results and are exportable within 5 seconds for datasets up to 1,500 events in performance testing.  
- **SC-004**: Support engineers resolve simulated webhook delivery issues using the troubleshooting view in under 10 minutes with no external tools in acceptance testing.  
- **SC-005**: Operators rate the admin console experience at ≥4 out of 5 on post-launch satisfaction surveys covering navigation, clarity, and troubleshooting effectiveness.
