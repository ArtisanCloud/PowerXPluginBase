package runtime_ops_test

import "testing"

func TestMigrationsBootstrap(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}
	// TODO: verify migrations via database harness when available
	t.Skip("runtime ops migration smoke test not implemented yet")
}
