package marketplace_test

import (
	"context"
	"testing"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketplacesvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/marketplace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupAnalyticsDB(t *testing.T) *gorm.DB {
	t.Helper()
	models.ForceSchemaForTests("")
	db, err := gorm.Open(sqlite.Open("file:analytics_flow?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(
		&dbm.Listing{},
		&dbm.PricingPlan{},
		&dbm.PlanTier{},
		&dbm.License{},
		&dbm.UsageEnvelope{},
		&dbm.UsageAggregate{},
		&dbm.RevenueShareReport{},
		&dbm.Notification{},
	))
	return db
}

func TestAnalyticsFlow_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx := context.Background()
	db := setupAnalyticsDB(t)

	usageRepo := mrepo.NewUsageRepository(db)
	revenueRepo := mrepo.NewRevenueRepository(db)
	notificationRepo := mrepo.NewNotificationRepository(db)
	licenseRepo := mrepo.NewLicenseRepository(db)
	pricingRepo := mrepo.NewPricingRepository(db)
	listingRepo := mrepo.NewListingRepository(db)

	cfg := &config.Config{
		Integration: &config.IntegrationConfig{
			Billing: config.IntegrationBillingConfig{
				Reconciliation: config.IntegrationRevenueSplitConfig{
					VendorShare:   0.80,
					PlatformShare: 0.15,
					FeeShare:      0.05,
					Currency:      "USD",
				},
			},
		},
	}

	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	logEntry := logrus.NewEntry(logger)

	analyticsSvc := marketplacesvc.NewAnalyticsService(cfg, usageRepo, revenueRepo, notificationRepo, licenseRepo, pricingRepo, listingRepo, logEntry)
	ingestSvc := marketplacesvc.NewUsageIngestService(cfg, usageRepo, licenseRepo, listingRepo, analyticsSvc, logEntry)

	amount := 199.0
	quota := 200.0
	plan := &dbm.PricingPlan{
		TenantID:  "tenant-1",
		ListingID: "listing-1",
		PlanCode:  "pro",
		PlanType:  dbm.PricingPlanTypeUsage,
		Currency:  "USD",
		Amount:    &amount,
		QuotaLimit: func(v float64) *float64 {
			return &v
		}(quota),
		Status: "active",
	}
	require.NoError(t, pricingRepo.CreatePlan(ctx, plan, nil))

	listing := &dbm.Listing{
		ID:        "listing-1",
		TenantID:  "tenant-1",
		PluginID:  "plugin.demo",
		VendorID:  "vendor-1",
		Status:    "published",
		Title:     "Demo Plugin",
		Locale:    "en-US",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	require.NoError(t, listingRepo.Create(ctx, listing))

	license := &dbm.License{
		ID:           "license-1",
		TenantID:     "tenant-1",
		ListingID:    listing.ID,
		PlanID:       plan.ID,
		LicenseToken: "token-1",
		Status:       dbm.LicenseStatusActive,
		IssuedAt:     time.Now().Add(-24 * time.Hour),
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour),
	}
	require.NoError(t, licenseRepo.CreateLicense(ctx, license, nil))

	now := time.Now().UTC().Truncate(time.Minute)
	batch := []marketplacesvc.UsageEnvelopeInput{
		{
			LicenseID:      "license-1",
			PluginID:       "plugin.demo",
			Metrics:        []marketplacesvc.UsageMetricInput{{Name: "calls", Unit: "count", Value: 120}},
			TimestampStart: now.Add(-30 * time.Minute),
			TimestampEnd:   now.Add(-15 * time.Minute),
			Signature:      "sig-1",
		},
		{
			LicenseID:      "license-1",
			PluginID:       "plugin.demo",
			Metrics:        []marketplacesvc.UsageMetricInput{{Name: "calls", Unit: "count", Value: 150}},
			TimestampStart: now.Add(-10 * time.Minute),
			TimestampEnd:   now,
			Signature:      "sig-2",
		},
	}

	result, err := ingestSvc.IngestBatch(ctx, "tenant-1", batch)
	require.NoError(t, err)
	require.Equal(t, 2, result.Accepted)
	require.Zero(t, result.Duplicates)
	require.Zero(t, result.Failed)

	aggregates, err := usageRepo.ListAggregates(ctx, "tenant-1", "license-1", dbm.AggregationWindowDay)
	require.NoError(t, err)
	require.NotEmpty(t, aggregates)
	require.Equal(t, "calls", aggregates[0].Metric)
	require.InDelta(t, 270.0, aggregates[0].Total, 0.001)
	require.InDelta(t, 150.0, aggregates[0].Delta, 0.001)

	dashboard, err := analyticsSvc.BuildDashboard(ctx, "tenant-1", "license-1", marketplacesvc.UsageDashboardQuery{
		Window: dbm.AggregationWindowDay,
		Metric: "calls",
	})
	require.NoError(t, err)
	require.NotNil(t, dashboard)
	require.NotEmpty(t, dashboard.Series)
	require.InDelta(t, 270.0, dashboard.Series[len(dashboard.Series)-1].Value, 0.001)
	require.NotEmpty(t, dashboard.Alerts)

	alertCodes := make(map[string]struct{})
	for _, alert := range dashboard.Alerts {
		alertCodes[alert.Code] = struct{}{}
	}
	if _, ok := alertCodes["quota_exceeded"]; !ok {
		t.Fatalf("expected quota_exceeded alert, got %#v", alertCodes)
	}

	reports, err := analyticsSvc.ListReports(ctx, "tenant-1", "vendor-1", time.Now().Add(-48*time.Hour), time.Now().Add(48*time.Hour))
	require.NoError(t, err)
	require.Len(t, reports, 1)
	require.InDelta(t, 270.0, reports[0].GrossAmount, 0.001)
	require.InDelta(t, 216.0, reports[0].VendorShare, 0.001)
	require.InDelta(t, 40.5, reports[0].PlatformShare, 0.001)
	require.InDelta(t, 13.5, reports[0].Fees, 0.001)

	notifications, err := notificationRepo.ListByTenant(ctx, "tenant-1")
	require.NoError(t, err)
	require.NotEmpty(t, notifications)
	require.Equal(t, "vendor", notifications[0].RecipientType)
}
