package marketplace

import (
	"context"
	"errors"
	"strings"
	"time"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/marketplace"
	repository "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
)

// LicenseRepository manages license persistence and events.
type LicenseRepository struct {
	*repository.BaseRepository[dbm.License]
}

// NewLicenseRepository constructs repository instance.
func NewLicenseRepository(db *gorm.DB) *LicenseRepository {
	return &LicenseRepository{
		BaseRepository: repository.NewBaseRepository[dbm.License](db),
	}
}

// CreateLicense stores a license and initial event.
func (r *LicenseRepository) CreateLicense(ctx context.Context, license *dbm.License, event *dbm.LicenseEvent) error {
	if license == nil {
		return errors.New("license is required")
	}
	tenantID := strings.TrimSpace(license.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	return r.WithTenantTx(ctx, tenantID, func(tx *gorm.DB) error {
		if err := tx.Create(license).Error; err != nil {
			return err
		}
		if event != nil {
			event.LicenseID = license.ID
			event.TenantID = tenantID
			if err := tx.Create(event).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateLicenseToken updates license token/status/expiry fields.
func (r *LicenseRepository) UpdateLicenseToken(ctx context.Context, tenantID, licenseID string, fields map[string]any) error {
	tenantID = strings.TrimSpace(tenantID)
	licenseID = strings.TrimSpace(licenseID)
	if tenantID == "" || licenseID == "" {
		return errors.New("tenant_id and license_id are required")
	}
	if len(fields) == 0 {
		return errors.New("no fields to update")
	}
	return r.WithTenantTx(ctx, tenantID, func(tx *gorm.DB) error {
		res := tx.Model(&dbm.License{}).
			Where("id = ? AND tenant_id = ?", licenseID, tenantID).
			Updates(fields)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}

// GetLicense returns license by ID.
func (r *LicenseRepository) GetLicense(ctx context.Context, tenantID, licenseID string) (*dbm.License, error) {
	tenantID = strings.TrimSpace(tenantID)
	licenseID = strings.TrimSpace(licenseID)
	if tenantID == "" || licenseID == "" {
		return nil, errors.New("tenant_id and license_id are required")
	}
	var license dbm.License
	err := r.DB.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, licenseID).
		First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

// FindActiveLicense finds active license by listing for a tenant.
func (r *LicenseRepository) FindActiveLicense(ctx context.Context, tenantID, listingID string) (*dbm.License, error) {
	tenantID = strings.TrimSpace(tenantID)
	listingID = strings.TrimSpace(listingID)
	if tenantID == "" || listingID == "" {
		return nil, errors.New("tenant_id and listing_id are required")
	}
	var license dbm.License
	err := r.DB.WithContext(ctx).
		Where("tenant_id = ? AND listing_id = ? AND status IN ?", tenantID, listingID, []string{dbm.LicenseStatusActive, dbm.LicenseStatusTrial}).
		Order("expires_at DESC").
		First(&license).Error
	if err != nil {
		return nil, err
	}
	return &license, nil
}

// CreateEvent appends a license lifecycle event.
func (r *LicenseRepository) CreateEvent(ctx context.Context, event *dbm.LicenseEvent) error {
	if event == nil {
		return errors.New("event is required")
	}
	tenantID := strings.TrimSpace(event.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	return r.WithTenantTx(ctx, tenantID, func(tx *gorm.DB) error {
		return tx.Create(event).Error
	})
}

// ListEvents returns events for a license.
func (r *LicenseRepository) ListEvents(ctx context.Context, tenantID, licenseID string, limit int) ([]*dbm.LicenseEvent, error) {
	tenantID = strings.TrimSpace(tenantID)
	licenseID = strings.TrimSpace(licenseID)
	if tenantID == "" || licenseID == "" {
		return nil, errors.New("tenant_id and license_id are required")
	}
	query := r.DB.WithContext(ctx).
		Where("tenant_id = ? AND license_id = ?", tenantID, licenseID).
		Order("emitted_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	var events []*dbm.LicenseEvent
	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

// UpdateOfflineWindow updates offline_until timestamp for a license.
func (r *LicenseRepository) UpdateOfflineWindow(ctx context.Context, tenantID, licenseID string, offlineUntil *time.Time) error {
	fields := map[string]any{"offline_until": offlineUntil}
	return r.UpdateLicenseToken(ctx, tenantID, licenseID, fields)
}

// RecordTaxTransaction stores a tax transaction entry.
func (r *LicenseRepository) RecordTaxTransaction(ctx context.Context, txn *dbm.TaxTransaction) error {
	if txn == nil {
		return errors.New("transaction is required")
	}
	tenantID := strings.TrimSpace(txn.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	return r.WithTenantTx(ctx, tenantID, func(tx *gorm.DB) error {
		return tx.Create(txn).Error
	})
}
