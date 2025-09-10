package iam

import (
	"context"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/iam"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MemberExtRepository struct {
	*repository.BaseRepository[dbm.MemberExt]
	db *gorm.DB
}

func NewMemberExtRepository(db *gorm.DB) *MemberExtRepository {
	return &MemberExtRepository{
		BaseRepository: repository.NewBaseRepository[dbm.MemberExt](db),
		db:             db,
	}
}

func (r *MemberExtRepository) FindByID(ctx context.Context, id uint64) (*dbm.MemberExt, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

func (r *MemberExtRepository) GetByTenantAndMemberID(ctx context.Context, tenantID, memberID uint64) (*dbm.MemberExt, error) {
	cond := map[string]interface{}{"tenant_id": tenantID, "member_id": memberID}
	return r.BaseRepository.GetByCondition(ctx, cond, nil)
}

// Upsert：基于 (tenant_id, member_id)
func (r *MemberExtRepository) UpsertByTenantAndMemberID(ctx context.Context, obj *dbm.MemberExt) (*dbm.MemberExt, error) {
	unique := []clause.Column{{Name: "tenant_id"}, {Name: "member_id"}}
	return r.BaseRepository.Upsert(ctx, obj, unique)
}
