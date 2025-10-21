package marketplace_test

import "testing"

func TestMarketplaceUsageContract(t *testing.T) {
	spec := loadOpenAPI(t)

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths section missing or malformed")
	}

	usagePath, ok := paths["/marketplace/usage"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/usage path not defined in contract")
	}

	postOp, ok := usagePath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/usage operation missing")
	}
	if opID := postOp["operationId"]; opID != "ingestUsage" {
		t.Fatalf("unexpected operationId for POST /marketplace/usage: %v", opID)
	}
	if responses, ok := postOp["responses"].(map[string]any); !ok {
		t.Fatalf("responses missing for POST /marketplace/usage")
	} else if _, ok := responses["202"]; !ok {
		t.Fatalf("202 response missing for POST /marketplace/usage")
	}

	metricsPath, ok := paths["/marketplace/usage/tenants/{tenantId}/licenses/{licenseId}/metrics"].(map[string]any)
	if !ok {
		t.Fatalf("usage metrics path not documented")
	}
	getOp, ok := metricsPath["get"].(map[string]any)
	if !ok {
		t.Fatalf("GET /marketplace/usage/tenants/{tenantId}/licenses/{licenseId}/metrics operation missing")
	}
	if opID := getOp["operationId"]; opID != "getUsageMetrics" {
		t.Fatalf("unexpected operationId for usage metrics GET: %v", opID)
	}
	if params, ok := getOp["parameters"].([]any); ok {
		foundWindow := false
		for _, param := range params {
			obj, _ := param.(map[string]any)
			if obj["name"] == "window" {
				foundWindow = true
				break
			}
		}
		if !foundWindow {
			t.Fatalf("window query parameter expected for usage metrics endpoint")
		}
	} else {
		t.Fatalf("parameters missing for usage metrics endpoint")
	}

	reportsPath, ok := paths["/marketplace/revenue-share/reports"].(map[string]any)
	if !ok {
		t.Fatalf("revenue share reports path not documented")
	}
	reportsGet, ok := reportsPath["get"].(map[string]any)
	if !ok {
		t.Fatalf("GET /marketplace/revenue-share/reports operation missing")
	}
	if opID := reportsGet["operationId"]; opID != "listRevenueReports" {
		t.Fatalf("unexpected operationId for revenue reports GET: %v", opID)
	}
}

func TestMarketplaceUsageSchemas(t *testing.T) {
	spec := loadOpenAPI(t)
	components, ok := spec["components"].(map[string]any)
	if !ok {
		t.Fatalf("components section missing")
	}
	schemas, ok := components["schemas"].(map[string]any)
	if !ok {
		t.Fatalf("schemas section missing")
	}

	requiredSchemas := []string{
		"UsageEnvelopeBatch",
		"UsageEnvelope",
		"UsageMetricsResponse",
		"RevenueReportSummary",
	}
	for _, name := range requiredSchemas {
		if _, ok := schemas[name]; !ok {
			t.Fatalf("expected schema %s to be documented", name)
		}
	}
}
