package iam

import (
	"context"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/iam"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberCacheRepository struct {
	*repository.BaseRepository[dbm.MemberCache]
	db *gorm.DB
}

func NewMemberCacheRepository(db *gorm.DB) *MemberCacheRepository {
	return &MemberCacheRepository{
		BaseRepository: repository.NewBaseRepository[dbm.MemberCache](db),
		db:             db,
	}
}

func (r *MemberCacheRepository) FindByID(ctx context.Context, id uint64) (*dbm.MemberCache, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

func (r *MemberCacheRepository) GetByTenantAndMemberID(ctx context.Context, tenantID, memberID uint64) (*dbm.MemberCache, error) {
	cond := map[string]interface{}{"tenant_id": tenantID, "member_id": memberID}
	return r.BaseRepository.GetByCondition(ctx, cond, nil)
}

// Upsert：基于 (tenant_id, member_id)
func (r *MemberCacheRepository) UpsertByTenantAndMemberID(ctx context.Context, obj *dbm.MemberCache) (*dbm.MemberCache, error) {
	unique := []clause.Column{{Name: "tenant_id"}, {Name: "member_id"}}
	return r.BaseRepository.Upsert(ctx, obj, unique)
}
