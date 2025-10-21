package marketplace_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	metrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/marketplace"
)

func TestListingSubmissionSLAInstrumentation(t *testing.T) {
	metrics.ResetMetrics()
	metrics.ObserveListingSubmission("success", "tenant-1", 45*time.Second)
	metrics.ObserveListingSubmission("success", "tenant-1", 4*time.Minute)

	var buf bytes.Buffer
	metrics.RenderMetrics(&buf)
	output := buf.String()

	if !strings.Contains(output, "powerx_marketplace_listing_submission_seconds") {
		t.Fatalf("expected submission latency histogram in metrics output")
	}
	if !strings.Contains(output, "status=\"timeout\"") {
		t.Fatalf("expected timeout status counter in metrics output")
	}
}
