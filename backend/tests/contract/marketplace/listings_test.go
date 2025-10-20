package marketplace_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func loadOpenAPI(t *testing.T) map[string]any {
	t.Helper()
	specPath := filepath.Join("specs", "006-marketplace-business", "contracts", "marketplace-openapi.yaml")
	data, err := os.ReadFile(specPath)
	if err != nil {
		t.Fatalf("failed to read OpenAPI contract: %v", err)
	}
	var spec map[string]any
	if err := yaml.Unmarshal(data, &spec); err != nil {
		t.Fatalf("failed to unmarshal OpenAPI contract: %v", err)
	}
	return spec
}

func TestMarketplaceListingsContract(t *testing.T) {
	spec := loadOpenAPI(t)

	paths, ok := spec["paths"].(map[string]any)
	if !ok {
		t.Fatalf("paths section missing or malformed")
	}

	listingsPath, ok := paths["/marketplace/listings"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/listings path not defined in contract")
	}

	postOp, ok := listingsPath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/listings operation missing")
	}
	if opID := postOp["operationId"]; opID != "createListingDraft" {
		t.Fatalf("unexpected operationId for POST /marketplace/listings: %v", opID)
	}

	resp, ok := postOp["responses"].(map[string]any)
	if !ok {
		t.Fatalf("responses for POST /marketplace/listings missing")
	}

	status201, ok := resp["201"].(map[string]any)
	if !ok {
		t.Fatalf("201 response missing for POST /marketplace/listings")
	}

	content, ok := status201["content"].(map[string]any)
	if !ok {
		t.Fatalf("content section missing for POST /marketplace/listings 201 response")
	}

	appJSON, ok := content["application/json"].(map[string]any)
	if !ok {
		t.Fatalf("application/json response missing for POST /marketplace/listings")
	}

	schema, ok := appJSON["schema"].(map[string]any)
	if !ok {
		t.Fatalf("schema not defined for POST /marketplace/listings response")
	}

	if ref := schema["$ref"]; ref != "#/components/schemas/Listing" {
		t.Fatalf("POST /marketplace/listings response expected to reference Listing schema, got %v", ref)
	}

	getOp, ok := listingsPath["get"].(map[string]any)
	if !ok {
		t.Fatalf("GET /marketplace/listings operation missing")
	}
	if opID := getOp["operationId"]; opID != "listListings" {
		t.Fatalf("unexpected operationId for GET /marketplace/listings: %v", opID)
	}
}

func TestMarketplaceListingStatusEndpoint(t *testing.T) {
	spec := loadOpenAPI(t)
	paths := spec["paths"].(map[string]any)
	statusPath, ok := paths["/marketplace/listings/{id}/status"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/listings/{id}/status not documented")
	}
	patchOp, ok := statusPath["patch"].(map[string]any)
	if !ok {
		t.Fatalf("PATCH /marketplace/listings/{id}/status operation missing")
	}
	if opID := patchOp["operationId"]; opID != "updateListingStatus" {
		t.Fatalf("unexpected operationId for PATCH /marketplace/listings/{id}/status: %v", opID)
	}
}

func TestMarketplaceChecklistContract(t *testing.T) {
	schemaPath := filepath.Join("specs", "006-marketplace-business", "contracts", "ready-checklist.graphql")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read GraphQL schema: %v", err)
	}
	text := string(data)

	if !strings.Contains(text, "type Mutation") || !strings.Contains(text, "triggerChecklistRun") {
		t.Fatalf("GraphQL schema missing triggerChecklistRun mutation")
	}
	if !strings.Contains(text, "type Query") || !strings.Contains(text, "listingChecklistRuns") {
		t.Fatalf("GraphQL schema missing listingChecklistRuns query")
	}
}
