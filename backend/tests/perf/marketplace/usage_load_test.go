package marketplace_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	metrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/marketplace"
)

func TestUsageIngestLagMetric(t *testing.T) {
	metrics.ResetMetrics()

	batches := 10000
	lag := 15 * time.Millisecond
	for i := 0; i < batches; i++ {
		metrics.ObserveUsageLag("tenant-perf", lag)
	}

	var buf bytes.Buffer
	metrics.RenderMetrics(&buf)
	out := buf.String()

	if !strings.Contains(out, "powerx_marketplace_usage_ingest_lag_seconds") {
		t.Fatalf("expected usage ingest lag histogram in metrics output")
	}
	if !strings.Contains(out, "tenant=\"tenant-perf\"") {
		t.Fatalf("expected tenant label for usage ingest lag metric")
	}
}
