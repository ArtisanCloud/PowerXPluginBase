package marketplace_test

import (
	"context"
	"testing"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketplacesvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/marketplace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestUsageSpikeDetection(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx := context.Background()
	db := setupAnalyticsDB(t)
	if db.Dialector.Name() == "sqlite" {
		t.Skip("analytics flow integration test requires PostgreSQL features")
	}

	usageRepo := mrepo.NewUsageRepository(db)
	revenueRepo := mrepo.NewRevenueRepository(db)
	notificationRepo := mrepo.NewNotificationRepository(db)
	licenseRepo := mrepo.NewLicenseRepository(db)
	pricingRepo := mrepo.NewPricingRepository(db)
	listingRepo := mrepo.NewListingRepository(db)

	cfg := &config.Config{}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logEntry := logrus.NewEntry(logger)

	analyticsSvc := marketplacesvc.NewAnalyticsService(cfg, usageRepo, revenueRepo, notificationRepo, licenseRepo, pricingRepo, listingRepo, logEntry)
	ingestSvc := marketplacesvc.NewUsageIngestService(cfg, usageRepo, licenseRepo, listingRepo, analyticsSvc, logEntry)

	quota := 500.0
	plan := &dbm.PricingPlan{
		TenantID:   "tenant-spike",
		ListingID:  "listing-spike",
		PlanCode:   "pro",
		PlanType:   dbm.PricingPlanTypeUsage,
		Currency:   "USD",
		QuotaLimit: &quota,
		Status:     "active",
	}
	require.NoError(t, pricingRepo.CreatePlan(ctx, plan, nil))

	listing := &dbm.Listing{
		ID:        "listing-spike",
		TenantID:  "tenant-spike",
		PluginID:  "plugin.demo",
		VendorID:  "vendor-spike",
		Status:    "published",
		Title:     "Spike Plugin",
		Locale:    "en-US",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, listingRepo.Create(ctx, listing))

	license := &dbm.License{
		ID:           "license-spike",
		TenantID:     "tenant-spike",
		ListingID:    listing.ID,
		PlanID:       plan.ID,
		LicenseToken: "token-spike",
		Status:       dbm.LicenseStatusActive,
		IssuedAt:     time.Now().Add(-24 * time.Hour),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	require.NoError(t, licenseRepo.CreateLicense(ctx, license, nil))

	now := time.Now().UTC()
	batch1 := []marketplacesvc.UsageEnvelopeInput{
		{
			LicenseID:      license.ID,
			PluginID:       listing.PluginID,
			Metrics:        []marketplacesvc.UsageMetricInput{{Name: "calls", Unit: "count", Value: 50}},
			TimestampStart: now.Add(-2 * time.Hour),
			TimestampEnd:   now.Add(-90 * time.Minute),
			Signature:      "sig-1",
		},
	}
	res, err := ingestSvc.IngestBatch(ctx, listing.TenantID, batch1)
	require.NoError(t, err)
	require.Equal(t, 1, res.Accepted)

	batch2 := []marketplacesvc.UsageEnvelopeInput{
		{
			LicenseID:      license.ID,
			PluginID:       listing.PluginID,
			Metrics:        []marketplacesvc.UsageMetricInput{{Name: "calls", Unit: "count", Value: 200}},
			TimestampStart: now.Add(-30 * time.Minute),
			TimestampEnd:   now,
			Signature:      "sig-2",
		},
	}
	res2, err := ingestSvc.IngestBatch(ctx, listing.TenantID, batch2)
	require.NoError(t, err)
	require.Equal(t, 1, res2.Accepted)

	dashboard, err := analyticsSvc.BuildDashboard(ctx, listing.TenantID, license.ID, marketplacesvc.UsageDashboardQuery{
		Window: dbm.AggregationWindowDay,
		Metric: "calls",
	})
	require.NoError(t, err)
	require.NotNil(t, dashboard)

	foundSpike := false
	for _, alert := range dashboard.Alerts {
		if alert.Code == "usage_spike" {
			foundSpike = true
		}
	}
	if !foundSpike {
		t.Fatalf("expected usage_spike alert after ingest spike")
	}

	notifications, err := notificationRepo.ListByTenant(ctx, listing.TenantID)
	require.NoError(t, err)
	require.NotEmpty(t, notifications)
	if notifications[0].TemplateCode != "marketplace.usage.alert" {
		t.Fatalf("expected usage alert notification, got %s", notifications[0].TemplateCode)
	}
}
