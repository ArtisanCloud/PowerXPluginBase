package marketplace_test

import (
	"context"
	"testing"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketplacesvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/marketplace"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

type recoveryBillingStub struct {
	calls int
}

func (s *recoveryBillingStub) ChargeSubscription(ctx context.Context, tenantID string, plan *dbm.PricingPlan, metadata map[string]any) (string, error) {
	s.calls++
	return "should-not-be-used", nil
}

func TestLicenseRecovery_ReissuesMissingLicense(t *testing.T) {
	if testing.Short() {
		t.Skip("skip integration test in short mode")
	}

	ctx := context.Background()
	db := setupLicenseFlowDB(t)
	pricingRepo := mrepo.NewPricingRepository(db)
	licenseRepo := mrepo.NewLicenseRepository(db)

	amount := 59.0
	plan := &dbm.PricingPlan{
		TenantID:  "tenant-42",
		ListingID: "listing-99",
		PlanCode:  "premium",
		PlanType:  dbm.PricingPlanTypeSubscription,
		Currency:  "USD",
		Amount:    &amount,
		Status:    "active",
	}
	require.NoError(t, pricingRepo.CreatePlan(ctx, plan, nil))

	billing := &recoveryBillingStub{}
	authority := &integrationAuthorityStub{token: "recovered-token"}
	cache := &integrationCacheStub{}
	logger := logrus.New()
	logger.SetLevel(logrus.WarnLevel)
	service := marketplacesvc.NewLicenseService(
		&config.Config{},
		pricingRepo,
		licenseRepo,
		nil,
		billing,
		authority,
		cache,
		logrus.New().WithField("component", "license_recovery_test_service"),
	)

	recovery := marketplacesvc.NewLicenseRecoveryService(
		service,
		licenseRepo,
		pricingRepo,
		logrus.New().WithField("component", "license_recovery_test"),
	)

	payload := marketplacesvc.RecoveryRequest{
		TenantID:  "tenant-42",
		ListingID: "listing-99",
		PlanID:    plan.ID,
		BillingID: "bill-delayed-1",
		IssuedBy:  "recovery-daemon",
		Metadata: map[string]any{
			"source": "billing-webhook",
		},
	}

	license, recovered, err := recovery.RecoverIssuance(ctx, payload)
	require.NoError(t, err)
	require.True(t, recovered)
	require.NotNil(t, license)
	require.Equal(t, "recovered-token", license.LicenseToken)
	require.Equal(t, dbm.LicenseStatusActive, license.Status)
	require.Equal(t, "bill-delayed-1", license.Metadata["billing_id"])
	require.Equal(t, "billing-webhook", license.Metadata["source"])
	require.Equal(t, "recovery-daemon", *license.IssuedBy)
	require.Contains(t, license.Metadata, "recovery")
	require.Equal(t, 0, billing.calls, "billing should not be re-triggered during recovery")

	events, err := licenseRepo.ListEvents(ctx, payload.TenantID, license.ID, 5)
	require.NoError(t, err)
	require.NotEmpty(t, events)
	require.Equal(t, "bill-delayed-1", events[0].EventPayload["billing_id"])

	require.Len(t, cache.entries, 1)
	require.Equal(t, "tenant-42", cache.entries[0].tenantID)

	licenseAgain, recoveredAgain, err := recovery.RecoverIssuance(ctx, payload)
	require.NoError(t, err)
	require.False(t, recoveredAgain, "second recovery attempt should be a no-op")
	require.Equal(t, license.ID, licenseAgain.ID)
}
