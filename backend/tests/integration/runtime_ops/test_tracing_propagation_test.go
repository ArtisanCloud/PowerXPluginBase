package runtime_ops_test

import "testing"

func TestTracingPropagation(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}
	t.Skip("tracing propagation integration test pending implementation")
}
