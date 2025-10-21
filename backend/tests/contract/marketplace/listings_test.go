package marketplace_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func contractPath(parts ...string) string {
	candidates := [][]string{
		parts,
		append([]string{".."}, parts...),
		append([]string{"..", ".."}, parts...),
		append([]string{"..", "..", ".."}, parts...),
		append([]string{"..", "..", "..", ".."}, parts...),
	}
	for _, pathParts := range candidates {
		p := filepath.Join(pathParts...)
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return filepath.Join(parts...)
}

func loadOpenAPI(t *testing.T) map[string]any {
	t.Helper()
	specPath := contractPath("specs", "006-marketplace-business", "contracts", "marketplace-openapi.yaml")
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
	statusPath, ok := paths["/marketplace/listings/{listingId}/status"].(map[string]any)
	if !ok {
		t.Fatalf("/marketplace/listings/{listingId}/status not documented")
	}
	patchOp, ok := statusPath["post"].(map[string]any)
	if !ok {
		t.Fatalf("POST /marketplace/listings/{listingId}/status operation missing")
	}
	if opID := patchOp["operationId"]; opID != "updateListingStatus" {
		t.Fatalf("unexpected operationId for POST /marketplace/listings/{listingId}/status: %v", opID)
	}
}

func TestMarketplaceChecklistContract(t *testing.T) {
	schemaPath := contractPath("specs", "006-marketplace-business", "contracts", "ready-checklist.graphql")
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
