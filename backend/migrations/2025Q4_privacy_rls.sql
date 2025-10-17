BEGIN;

CREATE OR REPLACE FUNCTION powerx_plugin_base.current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE powerx_plugin_base.privacy_data_classifications FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_classification_tenant_isolation ON powerx_plugin_base.privacy_data_classifications
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.privacy_consent_tokens FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_consent_tenant_isolation ON powerx_plugin_base.privacy_consent_tokens
    USING (tenant_id = powerx_plugin_base.current_tenant());

ALTER TABLE powerx_plugin_base.privacy_lifecycle_events FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_lifecycle_tenant_isolation ON powerx_plugin_base.privacy_lifecycle_events
    USING (tenant_id = powerx_plugin_base.current_tenant());

COMMIT;
