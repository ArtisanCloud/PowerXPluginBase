package marketplace_test

import (
	"context"
	"testing"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketplacesvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/marketplace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type fakeVendorGuard struct {
	revoked bool
}

func (f *fakeVendorGuard) VendorRevoked(ctx context.Context, vendorID string) (bool, error) {
	return f.revoked, nil
}

func setupMarketplaceDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file:marketplace_edge_cases?mode=memory&cache=shared"), &gorm.Config{})
	require.NoError(t, err)
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS marketplace_listings (
            id TEXT PRIMARY KEY,
            tenant_id TEXT NOT NULL,
            plugin_id TEXT NOT NULL,
            vendor_id TEXT NOT NULL,
            status TEXT NOT NULL,
            title TEXT NOT NULL,
            slug TEXT NOT NULL,
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
        );`,
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
        );`,
		`CREATE TABLE IF NOT EXISTS marketplace_listing_versions (
            id TEXT PRIMARY KEY,
            listing_id TEXT NOT NULL,
            tenant_id TEXT NOT NULL,
            version TEXT NOT NULL,
            changelog TEXT,
            metadata TEXT,
            submitted_by TEXT NOT NULL,
            review_state TEXT,
            reviewer_id TEXT,
            reviewed_at DATETIME,
            created_at DATETIME
        );`,
	}
	for _, stmt := range stmts {
		require.NoError(t, db.Exec(stmt).Error)
	}
	return db
}

func TestListingSubmissionEdgeCases(t *testing.T) {
	if testing.Short() {
		t.Skip("integration test skipped in short mode")
	}

	ctx := context.Background()
	db := setupMarketplaceDB(t)
	listingRepo := mrepo.NewListingRepository(db)
	checklistRepo := mrepo.NewChecklistRepository(db)
	service := marketplacesvc.NewListingService(listingRepo, checklistRepo, logrus.New().WithField("test", "listing_edge_cases"))

	guard := &fakeVendorGuard{revoked: false}
	service.SetVendorGuard(guard)

	listing := &dbm.Listing{
		ID:       "listing-edge",
		TenantID: "tenant-1",
		PluginID: "com.powerx.plugins.edge",
		VendorID: "vendor-1",
		Title:    "Edge Case Listing",
		Slug:     "edge-case-listing",
		Status:   dbm.ListingStatusDraft,
	}
	require.NoError(t, listingRepo.Create(ctx, listing))

	_, err := service.SubmitForReview(ctx, listing.TenantID, listing.ID, "vendor", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "no assets")

	asset := dbm.ListingAsset{
		ID:         "asset-1",
		ListingID:  listing.ID,
		TenantID:   listing.TenantID,
		AssetType:  "cover",
		StorageURI: "https://cdn.example/cover.png",
		Locale:     "en",
		IsPrimary:  true,
		Weight:     0,
	}
	require.NoError(t, listingRepo.ReplaceAssets(ctx, listing.TenantID, listing.ID, []dbm.ListingAsset{asset}))

	guard.revoked = true
	_, err = service.SubmitForReview(ctx, listing.TenantID, listing.ID, "vendor", nil)
	require.Error(t, err)
	require.Contains(t, err.Error(), "not permitted")

	guard.revoked = false
	updated, err := service.SubmitForReview(ctx, listing.TenantID, listing.ID, "vendor", map[string]any{"changelog": "initial submission"})
	require.NoError(t, err)
	require.Equal(t, dbm.ListingStatusInReview, updated.Status)

	guard.revoked = true
	_, err = service.UpdateListingStatus(ctx, listing.TenantID, listing.ID, "reviewer-1", dbm.ListingStatusPublished, "publishing now")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not permitted")

	guard.revoked = false
	published, err := service.UpdateListingStatus(ctx, listing.TenantID, listing.ID, "reviewer-1", dbm.ListingStatusPublished, "good to go")
	require.NoError(t, err)
	require.Equal(t, dbm.ListingStatusPublished, published.Status)
}
