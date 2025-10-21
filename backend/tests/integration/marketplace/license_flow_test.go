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

func setupLicenseFlowDB(t *testing.T) *gorm.DB {
	t.Helper()
	models.ForceSchemaForTests("")
	db, err := gorm.Open(sqlite.Open("file:license_flow?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS marketplace_pricing_plans (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL,
            listing_id TEXT NOT NULL,
            plan_code TEXT NOT NULL,
            plan_type TEXT NOT NULL,
            currency TEXT NOT NULL,
            amount REAL,
            billing_period TEXT,
            trial_period_days INTEGER,
            quota_limit REAL,
            overage_policy TEXT,
            feature_matrix TEXT,
            is_default INTEGER DEFAULT 0,
            status TEXT DEFAULT 'active',
            created_at DATETIME,
            updated_at DATETIME
        );`,
		`CREATE TABLE IF NOT EXISTS marketplace_plan_tiers (
            id TEXT PRIMARY KEY,
            plan_id TEXT NOT NULL,
            tenant_id TEXT NOT NULL,
            metric TEXT NOT NULL,
            range_from REAL NOT NULL,
            range_to REAL,
            unit_amount REAL NOT NULL,
            unit_name TEXT,
            created_at DATETIME,
            updated_at DATETIME
        );`,
		`CREATE TABLE IF NOT EXISTS marketplace_licenses (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL,
            listing_id TEXT NOT NULL,
            plan_id TEXT NOT NULL,
            license_token TEXT NOT NULL,
            status TEXT NOT NULL,
            issued_at DATETIME NOT NULL,
            expires_at DATETIME NOT NULL,
            renewal_token TEXT,
            offline_until DATETIME,
            last_validated_at DATETIME,
            issued_by TEXT,
            metadata TEXT,
            created_at DATETIME,
            updated_at DATETIME
        );`,
		`CREATE TABLE IF NOT EXISTS marketplace_license_events (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL,
            license_id TEXT NOT NULL,
            event_type TEXT NOT NULL,
            event_payload TEXT,
            emitted_at DATETIME,
            actor_id TEXT,
            trace_id TEXT,
            created_at DATETIME
        );`,
		`CREATE TABLE IF NOT EXISTS marketplace_tax_transactions (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL,
            billing_id TEXT NOT NULL,
            external_provider TEXT NOT NULL,
            external_transaction_id TEXT,
            jurisdiction TEXT,
            tax_amount REAL NOT NULL,
            currency TEXT NOT NULL,
            settlement_currency TEXT,
            exchange_rate REAL,
            tax_amount_settlement REAL,
            raw_payload TEXT,
            status TEXT NOT NULL,
            synced_at DATETIME,
            created_at DATETIME,
            updated_at DATETIME
        );`,
	}
	for _, stmt := range stmts {
		require.NoError(t, db.Exec(stmt).Error)
	}
	return db
}

type integrationBillingStub struct {
	calls int
	last  string
}

func (s *integrationBillingStub) ChargeSubscription(ctx context.Context, tenantID string, plan *dbm.PricingPlan, metadata map[string]any) (string, error) {
	s.calls++
	if plan != nil {
		s.last = plan.ID
	}
	return "billing-integration", nil
}

type integrationAuthorityStub struct {
	token string
}

func (s *integrationAuthorityStub) Issue(ctx context.Context, req *marketplacesvc.LicenseIssueRequest) (*marketplacesvc.LicenseIssueResponse, error) {
	return &marketplacesvc.LicenseIssueResponse{
		Token:     s.token,
		ExpiresAt: time.Now().Add(48 * time.Hour),
	}, nil
}

func (s *integrationAuthorityStub) Renew(ctx context.Context, req *marketplacesvc.LicenseRenewRequest) (*marketplacesvc.LicenseIssueResponse, error) {
	return &marketplacesvc.LicenseIssueResponse{
		Token:     s.token + "-renewed",
		ExpiresAt: time.Now().Add(72 * time.Hour),
	}, nil
}

func (s *integrationAuthorityStub) Revoke(ctx context.Context, licenseID string, reason string) error {
	return nil
}

func (s *integrationAuthorityStub) Verify(ctx context.Context, token string) (bool, error) {
	return true, nil
}

type integrationCacheStub struct {
	entries []cacheEntry
}

type cacheEntry struct {
	tenantID  string
	listingID string
	license   *dbm.License
	ttl       time.Duration
}

func (s *integrationCacheStub) Get(ctx context.Context, tenantID, listingID string) (*dbm.License, bool) {
	return nil, false
}

func (s *integrationCacheStub) Set(ctx context.Context, tenantID, listingID string, license *dbm.License, ttl time.Duration) {
	s.entries = append(s.entries, cacheEntry{
		tenantID:  tenantID,
		listingID: listingID,
		license:   license,
		ttl:       ttl,
	})
}

func (s *integrationCacheStub) Delete(ctx context.Context, tenantID, listingID string) {}

func TestLicenseFlow_EndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx := context.Background()
	db := setupLicenseFlowDB(t)
	pricingRepo := mrepo.NewPricingRepository(db)
	licenseRepo := mrepo.NewLicenseRepository(db)

	amount := 19.99
	plan := &dbm.PricingPlan{
		TenantID:  "tenant-1",
		ListingID: "listing-1",
		PlanCode:  "standard",
		PlanType:  dbm.PricingPlanTypeSubscription,
		Currency:  "USD",
		Amount:    &amount,
		Status:    "active",
	}
	require.NoError(t, pricingRepo.CreatePlan(ctx, plan, nil))

	billing := &integrationBillingStub{}
	authority := &integrationAuthorityStub{token: "integration-token"}
	cache := &integrationCacheStub{}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)

	service := marketplacesvc.NewLicenseService(&config.Config{}, pricingRepo, licenseRepo, nil, billing, authority, cache, logger.WithField("test", "license_flow"))

	issued, err := service.IssueLicense(ctx, marketplacesvc.IssueLicenseParams{
		TenantID:  "tenant-1",
		ListingID: "listing-1",
		PlanID:    plan.ID,
		IssuedBy:  "integration",
	})
	require.NoError(t, err)
	require.NotNil(t, issued)
	require.Equal(t, "integration-token", issued.LicenseToken)
	require.Equal(t, 1, billing.calls)
	require.Len(t, cache.entries, 1)

	stored, err := licenseRepo.GetLicense(ctx, "tenant-1", issued.ID)
	require.NoError(t, err)
	require.Equal(t, issued.ID, stored.ID)

	renewed, err := service.RenewLicense(ctx, marketplacesvc.RenewLicenseParams{
		LicenseID: issued.ID,
		TenantID:  "tenant-1",
		IssuedBy:  "integration",
	})
	require.NoError(t, err)
	require.Equal(t, issued.ID, renewed.ID)
	require.Equal(t, "integration-token-renewed", renewed.LicenseToken)

	events, err := licenseRepo.ListEvents(ctx, "tenant-1", issued.ID, 10)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(events), 2)
	require.Equal(t, dbm.LicenseEventRenewed, events[0].EventType)
}
