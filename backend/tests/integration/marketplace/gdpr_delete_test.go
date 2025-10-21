package marketplace_test

import (
	"context"
	"testing"
	"time"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketplacesvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/marketplace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestPrivacyService_PurgeUsageData(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx := context.Background()
	db := setupAnalyticsDB(t)
	usageRepo := mrepo.NewUsageRepository(db)
	privacy := marketplacesvc.NewPrivacyService(usageRepo, logrus.New().WithField("test", "gdpr"))

	tenantID := "tenant-gdpr"
	licenseID := "license-gdpr"
	now := time.Now().UTC()

	oldEnvelope := &dbm.UsageEnvelope{
		TenantID:       tenantID,
		LicenseID:      licenseID,
		PluginID:       "plugin.demo",
		TimestampStart: now.Add(-48 * time.Hour),
		TimestampEnd:   now.Add(-47 * time.Hour),
		Signature:      "sig-old",
		Checksum:       "old-checksum",
	}
	require.NoError(t, oldEnvelope.EncodeMetrics([]dbm.UsageMetric{{Name: "calls", Unit: "count", Value: 10}}))

	recentEnvelope := &dbm.UsageEnvelope{
		TenantID:       tenantID,
		LicenseID:      licenseID,
		PluginID:       "plugin.demo",
		TimestampStart: now.Add(-2 * time.Hour),
		TimestampEnd:   now.Add(-time.Hour),
		Signature:      "sig-new",
		Checksum:       "new-checksum",
	}
	require.NoError(t, recentEnvelope.EncodeMetrics([]dbm.UsageMetric{{Name: "calls", Unit: "count", Value: 20}}))

	require.NoError(t, db.Create(oldEnvelope).Error)
	require.NoError(t, db.Create(recentEnvelope).Error)

	oldAggregate := &dbm.UsageAggregate{
		TenantID:   tenantID,
		LicenseID:  licenseID,
		Metric:     "calls",
		Window:     dbm.AggregationWindowDay,
		TimeBucket: now.Add(-48 * time.Hour).Truncate(24 * time.Hour),
		Total:      10,
		Delta:      10,
		Currency:   "USD",
	}
	recentAggregate := &dbm.UsageAggregate{
		TenantID:   tenantID,
		LicenseID:  licenseID,
		Metric:     "calls",
		Window:     dbm.AggregationWindowDay,
		TimeBucket: now.Truncate(24 * time.Hour),
		Total:      20,
		Delta:      20,
		Currency:   "USD",
	}
	require.NoError(t, db.Create(oldAggregate).Error)
	require.NoError(t, db.Create(recentAggregate).Error)

	cutoff := now.Add(-24 * time.Hour)
	result, err := privacy.PurgeUsageData(ctx, tenantID, licenseID, cutoff)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, result.EnvelopesDeleted)
	require.Equal(t, 1, result.AggregatesDeleted)

	var envelopes []dbm.UsageEnvelope
	require.NoError(t, db.Where("tenant_id = ? AND license_id = ?", tenantID, licenseID).Find(&envelopes).Error)
	require.Len(t, envelopes, 1)
	if envelopes[0].Checksum != "new-checksum" {
		t.Fatalf("expected recent envelope to remain, got %s", envelopes[0].Checksum)
	}

	var aggregates []dbm.UsageAggregate
	require.NoError(t, db.Where("tenant_id = ? AND license_id = ?", tenantID, licenseID).Find(&aggregates).Error)
	require.Len(t, aggregates, 1)
	if aggregates[0].TimeBucket.Before(cutoff) {
		t.Fatalf("expected aggregate newer than cutoff to remain")
	}
}
