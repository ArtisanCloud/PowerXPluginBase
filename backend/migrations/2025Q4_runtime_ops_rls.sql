BEGIN;

CREATE OR REPLACE FUNCTION current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE runtime_assignments FORCE ROW LEVEL SECURITY;
CREATE POLICY runtime_assignments_tenant_isolation ON runtime_assignments
    USING (tenant_id IS NULL OR tenant_id = current_tenant());

ALTER TABLE port_reservations FORCE ROW LEVEL SECURITY;
CREATE POLICY port_reservations_tenant_isolation ON port_reservations
    USING (current_tenant() IS NULL OR EXISTS (
        SELECT 1 FROM runtime_assignments ra
        WHERE ra.id = runtime_assignment_id AND (ra.tenant_id IS NULL OR ra.tenant_id = current_tenant())
    ));

ALTER TABLE mcp_sessions FORCE ROW LEVEL SECURITY;
CREATE POLICY mcp_sessions_tenant_isolation ON mcp_sessions
    USING (tenant_id = current_tenant());

ALTER TABLE runtime_audit_events FORCE ROW LEVEL SECURITY;
CREATE POLICY runtime_audit_tenant_isolation ON runtime_audit_events
    USING (tenant_id IS NULL OR tenant_id = current_tenant());

ALTER TABLE quota_ledger FORCE ROW LEVEL SECURITY;
CREATE POLICY quota_ledger_tenant_isolation ON quota_ledger
    USING (scope_type <> 'tenant' OR scope_ref = current_tenant());

ALTER TABLE marketplace_overages FORCE ROW LEVEL SECURITY;
CREATE POLICY marketplace_overages_tenant_isolation ON marketplace_overages
    USING (tenant_id IS NULL OR tenant_id = current_tenant());

COMMIT;
