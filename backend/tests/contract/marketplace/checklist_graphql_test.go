package marketplace_test

import (
	"os"
	"strings"
	"testing"
)

func TestReadyChecklistGraphQLContract(t *testing.T) {
	schemaPath := contractPath("specs", "006-marketplace-business", "contracts", "ready-checklist.graphql")
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		t.Fatalf("failed to read GraphQL schema: %v", err)
	}
	text := string(data)
	if !strings.Contains(text, "type Mutation") || !strings.Contains(text, "triggerChecklistRun") {
		t.Fatalf("triggerChecklistRun mutation missing from ready checklist schema")
	}
	if !strings.Contains(text, "type Query") || !strings.Contains(text, "listingChecklistRuns") {
		t.Fatalf("listingChecklistRuns query missing from ready checklist schema")
	}
}
