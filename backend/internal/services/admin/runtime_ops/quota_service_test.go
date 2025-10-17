package runtime_ops

import (
	"context"
	"testing"
	"time"

	model "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/runtime_ops"
	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRecordBreachPersistsAuditEvent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.RuntimeAuditEvent{}, &model.QuotaLedger{}, &model.MarketplaceOverage{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	svc := NewQuotaService(db, nil)

	ctx := authx.ContextWithTenantID(context.Background(), 42)
	if err := svc.RecordBreach(ctx, "plugin.demo", "tenant-1", "bootstrap", ActionThrottle); err != nil {
		t.Fatalf("RecordBreach failed: %v", err)
	}

	var count int64
	if err := db.Model(&model.RuntimeAuditEvent{}).Count(&count).Error; err != nil {
		t.Fatalf("count audit events: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected 1 audit event, got %d", count)
	}
}

func TestRecordUsageAssignsID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&model.QuotaLedger{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	svc := NewQuotaService(db, nil)
	entry := &model.QuotaLedger{
		ScopeType:      "tenant",
		ScopeRef:       "tenant-1",
		WindowStart:    time.Now().Add(-5 * time.Minute),
		WindowEnd:      time.Now(),
		TokensConsumed: 10,
	}

	result, err := svc.RecordUsage(context.Background(), entry)
	if err != nil {
		t.Fatalf("RecordUsage failed: %v", err)
	}
	if result.ID == "" {
		t.Fatalf("expected RecordUsage to populate ID")
	}
}
