package security

import (
	"context"
	"io"
	"testing"
	"time"

	secmodel "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/security"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return logger.WithField("component", "test")
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&secmodel.Advisory{}, &secmodel.AdvisoryDistribution{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestCreateAdvisoryComputesSLADefault(t *testing.T) {
	db := setupTestDB(t)
	logger := newTestLogger()
	service := NewAdvisoryService(db, logger)

	fixed := time.Date(2025, 10, 17, 12, 0, 0, 0, time.UTC)
	service.WithClock(func() time.Time { return fixed })

	advisory, err := service.CreateAdvisory(context.Background(), CreateAdvisoryParams{
		Reference:        "PX-ADV-2025-0001",
		Severity:         secmodel.AdvisorySeverityCritical,
		Summary:          "Critical vulnerability in parsing pipeline",
		AffectedVersions: []string{"1.0.0", "1.1.0"},
	})
	if err != nil {
		t.Fatalf("CreateAdvisory failed: %v", err)
	}
	if advisory.Status != secmodel.AdvisoryStatusOpen {
		t.Fatalf("expected OPEN status, got %s", advisory.Status)
	}
	if advisory.SlaDeadline == nil {
		t.Fatalf("expected SLA deadline to be populated")
	}
	expectedDeadline := fixed.Add(24 * time.Hour)
	if !advisory.SlaDeadline.Equal(expectedDeadline) {
		t.Fatalf("expected SLA deadline %v, got %v", expectedDeadline, advisory.SlaDeadline)
	}
}

func TestPublishAdvisoryUpdatesLifecycle(t *testing.T) {
	db := setupTestDB(t)
	logger := newTestLogger()
	service := NewAdvisoryService(db, logger)

	fixed := time.Date(2025, 11, 1, 9, 30, 0, 0, time.UTC)
	service.WithClock(func() time.Time { return fixed })

	advisory, err := service.CreateAdvisory(context.Background(), CreateAdvisoryParams{
		Reference:        "PX-ADV-2025-0002",
		Severity:         "high",
		Summary:          "High severity advisory",
		AffectedVersions: []string{"2.0.0"},
	})
	if err != nil {
		t.Fatalf("CreateAdvisory failed: %v", err)
	}

	published, distributions, err := service.PublishAdvisory(context.Background(), PublishAdvisoryParams{
		AdvisoryID:       advisory.ID,
		PatchedInVersion: "2.0.1",
		NotifyChannels:   []string{"marketplace", "webhook"},
	})
	if err != nil {
		t.Fatalf("PublishAdvisory failed: %v", err)
	}
	if published.Status != secmodel.AdvisoryStatusPublished {
		t.Fatalf("expected status PUBLISHED, got %s", published.Status)
	}
	if published.PatchedInVersion != "2.0.1" {
		t.Fatalf("expected patched_in_version to be updated")
	}
	if published.PublishedAt == nil || !published.PublishedAt.Equal(fixed) {
		t.Fatalf("expected PublishedAt to match fixed time")
	}
	if published.PatchedAt == nil || !published.PatchedAt.Equal(fixed) {
		t.Fatalf("expected PatchedAt to match fixed time")
	}
	if len(distributions) != 2 {
		t.Fatalf("expected 2 distribution records, got %d", len(distributions))
	}

	var count int64
	if err := db.Model(&secmodel.AdvisoryDistribution{}).Count(&count).Error; err != nil {
		t.Fatalf("count distributions: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected persisted distribution records, got %d", count)
	}
}
