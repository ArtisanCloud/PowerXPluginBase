package marketplace_test

import "testing"

func TestMarketplaceLicensesContract(t *testing.T) {
	spec := loadOpenAPI(t)

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths section missing or malformed")
	}

	licensesPath, ok := paths["/marketplace/licenses"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/licenses path missing in contract")
	}

	postOp, ok := licensesPath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/licenses operation missing")
	}
	if opID := postOp["operationId"]; opID != "createLicense" {
		t.Fatalf("unexpected operationId for POST /marketplace/licenses: %v", opID)
	}
	resp, ok := postOp["responses"].(map[string]any)
	if !ok {
		t.Fatalf("responses missing for POST /marketplace/licenses")
	}
	status201, ok := resp["201"].(map[string]any)
	if !ok {
		t.Fatalf("201 response missing for POST /marketplace/licenses")
	}
	content, ok := status201["content"].(map[string]any)
	if !ok {
		t.Fatalf("content missing for POST /marketplace/licenses 201 response")
	}
	appJSON, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatalf("application/json schema missing for POST /marketplace/licenses")
	}
	schema, ok := appJSON["schema"].(map[string]any)
	if !ok {
		t.Fatalf("schema missing for POST /marketplace/licenses response")
	}
	if ref := schema["$ref"]; ref != "#/components/schemas/License" {
		t.Fatalf("expected License schema for POST /marketplace/licenses, got %v", ref)
	}
}

func TestMarketplaceLicenseDetailAndRenewContract(t *testing.T) {
	spec := loadOpenAPI(t)
	paths := spec["paths"].(map[string]any)

	licensePath, ok := paths["/marketplace/licenses/{licenseId}"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/licenses/{licenseId} path missing")
	}

	getOp, ok := licensePath["get"].(map[string]any)
	if !ok {
		t.Fatalf("GET /marketplace/licenses/{licenseId} operation missing")
	}
	if opID := getOp["operationId"]; opID != "getLicense" {
		t.Fatalf("unexpected operationId for GET /marketplace/licenses/{licenseId}: %v", opID)
	}

	postOp, ok := licensePath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/licenses/{licenseId} operation missing")
	}
	if opID := postOp["operationId"]; opID != "renewLicense" {
		t.Fatalf("unexpected operationId for POST /marketplace/licenses/{licenseId}: %v", opID)
	}
}

func TestMarketplaceLicenseOfflineExtendContract(t *testing.T) {
	spec := loadOpenAPI(t)
	paths := spec["paths"].(map[string]any)

	extendPath, ok := paths["/marketplace/licenses/{licenseId}/offline-extend"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/licenses/{licenseId}/offline-extend path missing")
	}

	postOp, ok := extendPath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/licenses/{licenseId}/offline-extend operation missing")
	}
	if opID := postOp["operationId"]; opID != "extendOfflineLicense" {
		t.Fatalf("unexpected operationId for POST /marketplace/licenses/{licenseId}/offline-extend: %v", opID)
	}
}
