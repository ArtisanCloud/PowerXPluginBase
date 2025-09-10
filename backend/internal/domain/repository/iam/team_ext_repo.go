package iam

import (
	"context"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/iam"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamExtRepository struct {
	*repository.BaseRepository[dbm.TeamExt]
	db *gorm.DB
}

func NewTeamExtRepository(db *gorm.DB) *TeamExtRepository {
	return &TeamExtRepository{
		BaseRepository: repository.NewBaseRepository[dbm.TeamExt](db),
		db:             db,
	}
}

func (r *TeamExtRepository) FindByID(ctx context.Context, id uint64) (*dbm.TeamExt, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

func (r *TeamExtRepository) GetByTenantAndTeamID(ctx context.Context, tenantID, teamID uint64) (*dbm.TeamExt, error) {
	cond := map[string]interface{}{"tenant_id": tenantID, "team_id": teamID}
	return r.BaseRepository.GetByCondition(ctx, cond, nil)
}

// Upsert：基于 (tenant_id, team_id) 复合唯一键
func (r *TeamExtRepository) UpsertByTenantAndTeamID(ctx context.Context, obj *dbm.TeamExt) (*dbm.TeamExt, error) {
	unique := []clause.Column{{Name: "tenant_id"}, {Name: "team_id"}}
	return r.BaseRepository.Upsert(ctx, obj, unique)
}
