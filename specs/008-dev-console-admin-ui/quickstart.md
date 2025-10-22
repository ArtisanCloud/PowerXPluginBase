# Quickstart — Dev Console & Admin UI

## Prerequisites
1. Run backend setup once: `make dev-setup` (ensures Go modules & tools).  
2. Install admin UI deps: `cd web-admin && npm install`.  
3. Start local dependencies for observability/webhook mocks if needed: `docker compose -f config/docker-compose.integration.yml up -d`.  
4. Copy config sample if not present: `cp backend/etc/config.example.yaml backend/etc/config.yaml` and update `admin_console` section (coming migration adds defaults).

## Apply Migrations
```bash
make migrate
```
This registers:
- `admin_console_audit_events`
- `admin_console_config_changes`
- `admin_console_job_runs`
Indices defined in `backend/migrations/2025Q4_admin_console.sql`.

## Run Dev Services
```bash
make run
```
Backend exposes admin APIs at `http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console`.

Frontend (Nuxt) dev server:
```bash
cd web-admin
npm run dev
```
Access admin console at `http://localhost:3000/_p/com.powerx.plugins.base/admin/dev-console`.

## Smoke Test APIs
### Fetch configuration sections
```bash
curl -H "Authorization: Bearer <token>" \
  "http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console/config/sections?tenant_id=demo-tenant"
```

### Update a configuration section
```bash
curl -X PUT \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"values":{"audit_retention_days":200,"config_change_retention_days":120,"job_history_days":45}}' \
  "http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console/config/sections/admin_console.retention"
```

### Query audit history
```bash
curl -H "Authorization: Bearer <token)" \
  "http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console/audit/events?tenant_id=demo-tenant&occurred_after=2025-10-01T00:00:00Z"
```

### Trigger safe operation replay
```bash
curl -X POST \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"action":"replay","scope_type":"tenant","scope_ref":"demo-tenant","target_id":"hook-evt-123"}' \
  "http://localhost:8086/_p/com.powerx.plugins.base/api/v1/admin/dev-console/safe-ops/actions"
```

## Frontend Testing
```bash
cd web-admin
npm run test tests/dev-console/configure_console.spec.ts
npm run test:e2e -- --grep dev-console  # Playwright smoke for console flows
```

## Backend Testing
```bash
make test ARGS='./backend/internal/services/admin/console/...'
make test ARGS='./backend/tests/integration/admin_console/...'
```

## Deployment Notes
- Update `plugin.yaml` with new permission codes (`operations.plugin.admin`, `operations.plugin.audit`, `operations.plugin.ops`) and admin navigation entries.
- Run `make release` to produce binary + `web-admin/.output` bundle.
- Document new configuration keys in `backend/etc/README.md` (`admin_console.job_history_days`, `admin_console.audit_retention_days`).
