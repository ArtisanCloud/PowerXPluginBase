BEGIN;

CREATE OR REPLACE FUNCTION current_tenant() RETURNS TEXT AS $$
BEGIN
    RETURN current_setting('app.tenant_id', true);
END;
$$ LANGUAGE plpgsql STABLE;

ALTER TABLE privacy_data_classifications FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_classification_tenant_isolation ON privacy_data_classifications
    USING (tenant_id = current_tenant());

ALTER TABLE privacy_consent_tokens FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_consent_tenant_isolation ON privacy_consent_tokens
    USING (tenant_id = current_tenant());

ALTER TABLE privacy_lifecycle_events FORCE ROW LEVEL SECURITY;
CREATE POLICY privacy_lifecycle_tenant_isolation ON privacy_lifecycle_events
    USING (tenant_id = current_tenant());

COMMIT;
