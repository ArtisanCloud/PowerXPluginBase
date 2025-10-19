package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	model "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/integration"
	repo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/integration"
	service "github.com/ArtisanCloud/PowerXPlugin/internal/services/integration"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupWebhookService(t *testing.T) (*service.WebhookService, context.Context) {
	t.Helper()

	models.InitSchemaFrom("main")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&model.WebhookSubscription{}, &model.DeliveryAttempt{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	cfg := &config.Config{
		Server: &config.ServerConfig{
			SecretKey: "test-secret",
		},
	}

	subRepo := repo.NewWebhookSubscriptionRepository(db)
	attemptRepo := repo.NewDeliveryAttemptRepository(db)

	return service.NewWebhookService(cfg, subRepo, attemptRepo, nil), context.Background()
}

func TestWebhookService_ReplayAttemptFlow(t *testing.T) {
	svc, ctx := setupWebhookService(t)

	sub, err := svc.CreateSubscription(ctx, service.CreateSubscriptionParams{
		TenantID:  "42",
		EventType: "integration.dispatch",
		TargetURL: "https://example.org/webhook",
	})
	if err != nil {
		t.Fatalf("CreateSubscription failed: %v", err)
	}
	if sub == nil || sub.ID == "" {
		t.Fatalf("CreateSubscription returned empty subscription")
	}

	attempt, err := svc.RecordDeliveryAttempt(ctx, service.DeliveryResultParams{
		SubscriptionID: sub.ID,
		TenantID:       sub.TenantID,
		Status:         model.AttemptStatusDLQ,
		RetryCount:     2,
	})
	if err != nil {
		t.Fatalf("RecordDeliveryAttempt failed: %v", err)
	}
	if attempt.Status != model.AttemptStatusDLQ {
		t.Fatalf("expected attempt status DLQ, got %s", attempt.Status)
	}

	now := time.Now().UTC()
	if err := svc.UpdateAttemptStatus(ctx, attempt.ID, model.AttemptStatusPending, attempt.RetryCount, &now, "", sub.TenantID); err != nil {
		t.Fatalf("UpdateAttemptStatus failed: %v", err)
	}

	updated, err := svc.GetAttempt(ctx, attempt.ID)
	if err != nil {
		t.Fatalf("GetAttempt failed: %v", err)
	}
	if updated.Status != model.AttemptStatusPending {
		t.Fatalf("expected status to be PENDING, got %s", updated.Status)
	}
	if updated.NextDeliveryAt == nil {
		t.Fatalf("expected next_delivery_at to be set")
	}
	if updated.NextDeliveryAt.Before(now) {
		t.Fatalf("expected next_delivery_at >= replay timestamp, got %v", updated.NextDeliveryAt)
	}
}
