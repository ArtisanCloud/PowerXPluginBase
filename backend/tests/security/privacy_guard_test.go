package security_test

import (
	"testing"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	agentsec "github.com/ArtisanCloud/PowerXPlugin/internal/services/agent/security"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	return db
}

func TestValidateOutboundURL(t *testing.T) {
	cfg := &config.Config{Security: &config.SecurityConfig{
		GatewayAllowlist: []string{"gateway.powerx.dev"},
		RequireTLS13:     true,
	}}
	guard := agentsec.NewPrivacyGuard(newTestDB(t), cfg, nil)

	if err := guard.ValidateOutboundURL("https://gateway.powerx.dev/api"); err != nil {
		t.Fatalf("expected allowlist host to pass: %v", err)
	}

	if err := guard.ValidateOutboundURL("http://gateway.powerx.dev/api"); err == nil {
		t.Fatal("expected non-https to be rejected")
	}

	if err := guard.ValidateOutboundURL("https://example.com/api"); err == nil {
		t.Fatal("expected host outside allowlist to be rejected")
	}
}
