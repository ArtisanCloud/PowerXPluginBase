package marketplace_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	metrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/marketplace"
)

func TestLicenseVerificationLatencyMetrics(t *testing.T) {
	metrics.ResetMetrics()

	metrics.ObserveLicenseVerification("success", "stripe_tax", "tenant-123", 120*time.Millisecond)
	metrics.ObserveLicenseVerification("cache_hit", "redis", "tenant-123", 15*time.Millisecond)
	metrics.ObserveLicenseVerification("success", "stripe_tax", "tenant-456", 210*time.Millisecond)

	var buf bytes.Buffer
	metrics.RenderMetrics(&buf)
	out := buf.String()

	if !strings.Contains(out, "powerx_marketplace_license_verify_seconds") {
		t.Fatalf("expected license verification histogram in metrics output")
	}
	if !strings.Contains(out, "tenant=\"tenant-123\"") {
		t.Fatalf("expected tenant label in metrics output: %s", out)
	}
	if !strings.Contains(out, "provider=\"stripe_tax\"") {
		t.Fatalf("expected provider label for stripe_tax in metrics output")
	}
}
