BEGIN;

CREATE OR REPLACE FUNCTION powerx_plugin_base.current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE powerx_plugin_base.runtime_assignments FORCE ROW LEVEL SECURITY;
CREATE POLICY runtime_assignments_tenant_isolation ON powerx_plugin_base.runtime_assignments
    USING (tenant_id IS NULL OR tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.port_reservations FORCE ROW LEVEL SECURITY;
CREATE POLICY port_reservations_tenant_isolation ON powerx_plugin_base.port_reservations
    USING (powerx_plugin_base.current_tenant() IS NULL OR EXISTS (
        SELECT 1 FROM powerx_plugin_base.runtime_assignments ra
        WHERE ra.id = runtime_assignment_id AND (ra.tenant_id IS NULL OR ra.tenant_id = powerx_plugin_base.current_tenant())
    ));

ALTER TABLE powerx_plugin_base.mcp_sessions FORCE ROW LEVEL SECURITY;
CREATE POLICY mcp_sessions_tenant_isolation ON powerx_plugin_base.mcp_sessions
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.runtime_audit_events FORCE ROW LEVEL SECURITY;
CREATE POLICY runtime_audit_tenant_isolation ON powerx_plugin_base.runtime_audit_events
    USING (tenant_id IS NULL OR tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.quota_ledger FORCE ROW LEVEL SECURITY;
CREATE POLICY quota_ledger_tenant_isolation ON powerx_plugin_base.quota_ledger
    USING (scope_type <> 'tenant' OR scope_ref = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.marketplace_overages FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_overages_tenant_isolation ON powerx_plugin_base.marketplace_overages
    USING (tenant_id IS NULL OR tenant_id = powerx_plugin_base.current_tenant());

COMMIT;
