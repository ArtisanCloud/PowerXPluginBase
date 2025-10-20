package marketplace

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	mrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/marketplace"
	marketobs "github.com/ArtisanCloud/PowerXPlugin/internal/observability/marketplace"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// BillingClient represents a billing engine integration.
type BillingClient interface {
	ChargeSubscription(ctx context.Context, tenantID string, plan *dbm.PricingPlan, metadata map[string]any) (string, error)
}

// LicenseAuthority communicates with the centralized license server.
type LicenseAuthority interface {
	Issue(ctx context.Context, req *LicenseIssueRequest) (*LicenseIssueResponse, error)
	Renew(ctx context.Context, req *LicenseRenewRequest) (*LicenseIssueResponse, error)
	Revoke(ctx context.Context, licenseID string, reason string) error
	Verify(ctx context.Context, token string) (bool, error)
}

// LicenseCache provides temporary license caching for verification.
type LicenseCache interface {
	Get(ctx context.Context, tenantID, listingID string) (*dbm.License, bool)
	Set(ctx context.Context, tenantID, listingID string, license *dbm.License, ttl time.Duration)
	Delete(ctx context.Context, tenantID, listingID string)
}

// LicenseIssueRequest captures issuance payload for authority calls.
type LicenseIssueRequest struct {
	TenantID  string
	ListingID string
	PlanID    string
	Metadata  map[string]any
}

// LicenseRenewRequest captures renew payload for authority calls.
type LicenseRenewRequest struct {
	LicenseID string
	TenantID  string
	PlanID    string
	Metadata  map[string]any
}

// LicenseIssueResponse represents license authority response.
type LicenseIssueResponse struct {
	Token     string
	ExpiresAt time.Time
	Metadata  map[string]any
}

// IssueLicenseParams encapsulates inputs for issuing license.
type IssueLicenseParams struct {
	TenantID  string
	ListingID string
	PlanID    string
	IssuedBy  string
	Trial     bool
	Metadata  map[string]any
	ExpiresAt time.Time
}

// RenewLicenseParams encapsulates inputs for renewal.
type RenewLicenseParams struct {
	LicenseID string
	TenantID  string
	IssuedBy  string
	Metadata  map[string]any
	ExpiresAt time.Time
}

// VerifyResult describes verification result.
type VerifyResult struct {
	License *dbm.License
	Valid   bool
	Reason  string
}

// LicenseService orchestrates plan lookup, billing, license issuance and caching.
type LicenseService struct {
	cfg           *config.Config
	pricingRepo   *mrepo.PricingRepository
	licenseRepo   *mrepo.LicenseRepository
	taxClient     *TaxProviderClient
	billingClient BillingClient
	authority     LicenseAuthority
	cache         LicenseCache
	logger        *logrus.Entry
}

// NewLicenseService constructs the service with dependencies.
func NewLicenseService(cfg *config.Config, pricingRepo *mrepo.PricingRepository, licenseRepo *mrepo.LicenseRepository, taxClient *TaxProviderClient, billing BillingClient, authority LicenseAuthority, cache LicenseCache, logger *logrus.Entry) *LicenseService {
	if logger == nil {
		logger = logrus.New().WithField("component", "marketplace_license_service")
	}
	return &LicenseService{
		cfg:           cfg,
		pricingRepo:   pricingRepo,
		licenseRepo:   licenseRepo,
		taxClient:     taxClient,
		billingClient: billing,
		authority:     authority,
		cache:         cache,
		logger:        logger,
	}
}

// IssueLicense issues a new license for the given plan.
func (s *LicenseService) IssueLicense(ctx context.Context, params IssueLicenseParams) (*dbm.License, error) {
	if hasEmpty(params.TenantID, params.ListingID, params.PlanID) {
		return nil, errors.New("tenant_id, listing_id, plan_id are required")
	}

	plan, err := s.pricingRepo.GetPlan(ctx, params.TenantID, params.PlanID)
	if err != nil {
		return nil, err
	}

	billingID := ""
	if s.billingClient != nil {
		if id, err := s.billingClient.ChargeSubscription(ctx, params.TenantID, plan, params.Metadata); err == nil {
			billingID = id
		} else {
			return nil, fmt.Errorf("billing charge failed: %w", err)
		}
	}

	var taxResult *TaxChargeResult
	if s.taxClient != nil && plan.Amount != nil {
		cents := int64(*plan.Amount * 100)
		req := &TaxChargeRequest{
			TenantID:          params.TenantID,
			ListingID:         params.ListingID,
			Currency:          plan.Currency,
			AmountCents:       cents,
			Jurisdiction:      "",
			ExternalReference: billingID,
			Items: []TaxLineItem{
				{SKU: plan.PlanCode, Quantity: 1, AmountCents: cents, TaxCode: plan.PlanType},
			},
		}
		if res, err := s.taxClient.CreateTransaction(ctx, req); err == nil {
			taxResult = res
		} else {
			s.logger.WithError(err).Warn("tax transaction failed")
		}
	}

	expiry := params.ExpiresAt
	if expiry.IsZero() {
		expiry = time.Now().Add(30 * 24 * time.Hour)
	}

	issueRes := &LicenseIssueResponse{}
	if s.authority != nil {
		payload := &LicenseIssueRequest{
			TenantID:  params.TenantID,
			ListingID: params.ListingID,
			PlanID:    params.PlanID,
			Metadata:  params.Metadata,
		}
		resp, err := s.authority.Issue(ctx, payload)
		if err != nil {
			return nil, fmt.Errorf("license authority issue failed: %w", err)
		}
		issueRes = resp
		if !resp.ExpiresAt.IsZero() {
			expiry = resp.ExpiresAt
		}
	} else {
		token := generateToken(params.TenantID, params.ListingID, params.PlanID)
		issueRes.Token = token
		issueRes.Metadata = params.Metadata
	}

	offlineUntil := timePtr(minTime(expiry, time.Now().Add(72*time.Hour)))
	license := &dbm.License{
		TenantID:     params.TenantID,
		ListingID:    params.ListingID,
		PlanID:       params.PlanID,
		LicenseToken: issueRes.Token,
		Status:       statusFromTrial(params.Trial),
		IssuedAt:     time.Now(),
		ExpiresAt:    expiry,
		OfflineUntil: offlineUntil,
		IssuedBy:     stringPtr(params.IssuedBy),
		Metadata:     toJSONMap(issueRes.Metadata),
	}

	event := &dbm.LicenseEvent{
		TenantID:  params.TenantID,
		EventType: dbm.LicenseEventIssued,
		EventPayload: toJSONMap(map[string]any{
			"plan_id":    params.PlanID,
			"billing_id": billingID,
		}),
	}

	if err := s.licenseRepo.CreateLicense(ctx, license, event); err != nil {
		return nil, err
	}

	if taxResult != nil {
		txn := &dbm.TaxTransaction{
			TenantID:              params.TenantID,
			BillingID:             billingID,
			ExternalProvider:      s.taxClient.Provider(),
			ExternalTransactionID: stringPtr(taxResult.ExternalTransactionID),
			Jurisdiction:          taxResult.Jurisdiction,
			TaxAmount:             float64(taxResult.TaxAmountCents) / 100,
			Currency:              taxResult.Currency,
			RawPayload:            jsonFromBytes(taxResult.RawPayload),
			Status:                "completed",
		}
		_ = s.licenseRepo.RecordTaxTransaction(ctx, txn)
	}

	provider := ""
	if s.taxClient != nil {
		provider = s.taxClient.Provider()
	}
	marketobs.ObserveLicenseVerification("issue", provider, params.TenantID, 0)
	if s.cache != nil {
		s.cache.Set(ctx, params.TenantID, params.ListingID, license, time.Until(expiry))
	}
	return license, nil
}

// RenewLicense renews an existing license.
func (s *LicenseService) RenewLicense(ctx context.Context, params RenewLicenseParams) (*dbm.License, error) {
	if hasEmpty(params.TenantID, params.LicenseID) {
		return nil, errors.New("tenant_id and license_id are required")
	}

	license, err := s.licenseRepo.GetLicense(ctx, params.TenantID, params.LicenseID)
	if err != nil {
		return nil, err
	}

	expiry := params.ExpiresAt
	if expiry.IsZero() {
		expiry = time.Now().Add(30 * 24 * time.Hour)
	}

	if s.authority != nil {
		resp, err := s.authority.Renew(ctx, &LicenseRenewRequest{
			LicenseID: license.ID,
			TenantID:  params.TenantID,
			PlanID:    license.PlanID,
			Metadata:  params.Metadata,
		})
		if err != nil {
			return nil, err
		}
		if resp.Token != "" {
			license.LicenseToken = resp.Token
		}
		if !resp.ExpiresAt.IsZero() {
			expiry = resp.ExpiresAt
		}
	}

	fields := map[string]any{
		"expires_at":        expiry,
		"status":            dbm.LicenseStatusActive,
		"last_validated_at": time.Now(),
	}
	if err := s.licenseRepo.UpdateLicenseToken(ctx, params.TenantID, license.ID, fields); err != nil {
		return nil, err
	}

	event := &dbm.LicenseEvent{
		TenantID:  params.TenantID,
		LicenseID: license.ID,
		EventType: dbm.LicenseEventRenewed,
		EventPayload: toJSONMap(map[string]any{
			"expires_at": expiry,
		}),
	}
	_ = s.licenseRepo.CreateEvent(ctx, event)

	if s.cache != nil {
		license.ExpiresAt = expiry
		s.cache.Set(ctx, params.TenantID, license.ListingID, license, time.Until(expiry))
	}

	return license, nil
}

// VerifyLicense checks the cache and repository to ensure license validity.
func (s *LicenseService) VerifyLicense(ctx context.Context, tenantID, listingID string) (*VerifyResult, error) {
	if hasEmpty(tenantID, listingID) {
		return nil, errors.New("tenant_id and listing_id are required")
	}
	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, tenantID, listingID); ok && cached != nil {
			return &VerifyResult{License: cached, Valid: !cached.ExpiresAt.Before(time.Now())}, nil
		}
	}
	license, err := s.licenseRepo.FindActiveLicense(ctx, tenantID, listingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &VerifyResult{Valid: false, Reason: "license not found"}, nil
		}
		return nil, err
	}
	valid := !license.ExpiresAt.Before(time.Now())
	if s.cache != nil {
		s.cache.Set(ctx, tenantID, listingID, license, freshTTL(license.ExpiresAt))
	}
	return &VerifyResult{License: license, Valid: valid}, nil
}

// RevokeLicense transitions a license to revoked state.
func (s *LicenseService) RevokeLicense(ctx context.Context, tenantID, licenseID, reason string) error {
	if hasEmpty(tenantID, licenseID) {
		return errors.New("tenant_id and license_id are required")
	}
	if s.authority != nil {
		if err := s.authority.Revoke(ctx, licenseID, reason); err != nil {
			return err
		}
	}
	if err := s.licenseRepo.UpdateLicenseToken(ctx, tenantID, licenseID, map[string]any{
		"status":            dbm.LicenseStatusRevoked,
		"last_validated_at": time.Now(),
	}); err != nil {
		return err
	}
	event := &dbm.LicenseEvent{
		TenantID:  tenantID,
		LicenseID: licenseID,
		EventType: dbm.LicenseEventRevoked,
		EventPayload: toJSONMap(map[string]any{
			"reason": reason,
		}),
	}
	_ = s.licenseRepo.CreateEvent(ctx, event)
	if s.cache != nil {
		s.cache.Delete(ctx, tenantID, licenseID)
	}
	return nil
}

// ExtendOffline updates offline window within 72 hour constraint.
func (s *LicenseService) ExtendOffline(ctx context.Context, tenantID, licenseID string, until time.Time) error {
	if hasEmpty(tenantID, licenseID) {
		return errors.New("tenant_id and license_id are required")
	}
	target := until
	max := time.Now().Add(72 * time.Hour)
	if until.After(max) {
		target = max
	}
	if err := s.licenseRepo.UpdateOfflineWindow(ctx, tenantID, licenseID, &target); err != nil {
		return err
	}
	event := &dbm.LicenseEvent{
		TenantID:  tenantID,
		LicenseID: licenseID,
		EventType: dbm.LicenseEventOfflineExtend,
		EventPayload: toJSONMap(map[string]any{
			"offline_until": target,
		}),
	}
	return s.licenseRepo.CreateEvent(ctx, event)
}

func hasEmpty(values ...string) bool {
	for _, v := range values {
		if strings.TrimSpace(v) == "" {
			return true
		}
	}
	return false
}

func statusFromTrial(trial bool) string {
	if trial {
		return dbm.LicenseStatusTrial
	}
	return dbm.LicenseStatusActive
}

func generateToken(tenantID, listingID, planID string) string {
	payload := fmt.Sprintf("%s|%s|%s|%s", tenantID, listingID, planID, uuid.NewString())
	return base64.RawURLEncoding.EncodeToString([]byte(payload))
}

func stringPtr(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	value := s
	return &value
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func minTime(a, b time.Time) time.Time {
	if b.Before(a) {
		return b
	}
	return a
}

func freshTTL(expiry time.Time) time.Duration {
	if expiry.IsZero() {
		return time.Hour
	}
	ttl := time.Until(expiry)
	if ttl <= 0 {
		return time.Minute
	}
	return ttl
}
