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
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS marketplace_listings (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      plugin_id TEXT NOT NULL,
      vendor_id TEXT,
      status TEXT,
      title TEXT,
      slug TEXT,
      summary TEXT,
      description TEXT,
      cover_asset_id TEXT,
      hero_video_asset_id TEXT,
      categories TEXT,
      tags TEXT,
      locale TEXT,
      version TEXT,
      ready_checklist_score INTEGER DEFAULT 0,
      recommended_weight REAL DEFAULT 0,
      published_at DATETIME,
      reviewed_at DATETIME,
      reviewer_id TEXT,
      audit_notes TEXT,
      branding_theme TEXT,
      created_at DATETIME,
      updated_at DATETIME,
      deleted_at DATETIME
	    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_listing_assets (
	      id TEXT PRIMARY KEY,
	      listing_id TEXT NOT NULL,
	      tenant_id TEXT NOT NULL,
	      asset_type TEXT NOT NULL,
	      storage_uri TEXT NOT NULL,
	      checksum TEXT,
	      is_primary INTEGER DEFAULT 0,
	      locale TEXT,
	      weight INTEGER DEFAULT 0,
	      metadata TEXT,
	      created_at DATETIME,
	      updated_at DATETIME
	    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_pricing_plans (
      id TEXT PRIMARY KEY,
      listing_id TEXT NOT NULL,
      tenant_id TEXT NOT NULL,
      plan_code TEXT,
      plan_type TEXT,
      currency TEXT,
      amount REAL,
      billing_period TEXT,
      trial_period_days INTEGER,
      quota_limit REAL,
      overage_policy TEXT,
      feature_matrix TEXT,
      is_default INTEGER,
      status TEXT,
      created_at DATETIME,
      updated_at DATETIME
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_plan_tiers (
      id TEXT PRIMARY KEY,
      plan_id TEXT NOT NULL,
      tenant_id TEXT NOT NULL,
      metric TEXT,
      range_from REAL,
      range_to REAL,
      unit_amount REAL,
      unit_name TEXT,
      created_at DATETIME,
      updated_at DATETIME
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_licenses (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      listing_id TEXT NOT NULL,
      plan_id TEXT NOT NULL,
      license_token TEXT,
      status TEXT,
      issued_at DATETIME,
      expires_at DATETIME,
      renewal_token TEXT,
      offline_until DATETIME,
      last_validated_at DATETIME,
      issued_by TEXT,
      metadata TEXT,
      created_at DATETIME,
      updated_at DATETIME
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_usage_envelopes (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      license_id TEXT NOT NULL,
      plugin_id TEXT NOT NULL,
      metrics TEXT,
      timestamp_start DATETIME,
      timestamp_end DATETIME,
      signature TEXT,
      checksum TEXT UNIQUE,
      ingest_status TEXT,
      ingested_at DATETIME,
      created_at DATETIME,
      updated_at DATETIME
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_usage_aggregates (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      license_id TEXT NOT NULL,
      metric TEXT NOT NULL,
      window TEXT NOT NULL,
      time_bucket DATETIME NOT NULL,
      total REAL,
      delta REAL,
      currency TEXT,
      revenue REAL,
      created_at DATETIME,
      updated_at DATETIME,
      UNIQUE (tenant_id, license_id, metric, window, time_bucket)
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_revenue_share_reports (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      vendor_id TEXT NOT NULL,
      period_start DATETIME,
      period_end DATETIME,
      gross_amount REAL,
      vendor_share REAL,
      platform_share REAL,
      fees REAL,
      currency TEXT,
      status TEXT,
      generated_at DATETIME,
      export_uri TEXT,
      created_at DATETIME,
      updated_at DATETIME,
      UNIQUE (tenant_id, vendor_id, period_start, period_end)
    )`,
		`CREATE TABLE IF NOT EXISTS marketplace_notifications (
      id TEXT PRIMARY KEY,
      tenant_id TEXT NOT NULL,
      recipient_type TEXT,
      recipient_id TEXT,
      channel TEXT,
      template_code TEXT,
      payload TEXT,
      scheduled_at DATETIME,
      sent_at DATETIME,
      status TEXT,
      created_at DATETIME,
      updated_at DATETIME
    )`,
	}
	for _, stmt := range stmts {
		require.NoError(t, db.Exec(stmt).Error)
	}
	return db
}

func TestAnalyticsFlow_EndToEnd(t *testing.T) {
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
