package iam

import (
	"context"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/iam"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TeamCacheRepository struct {
	*repository.BaseRepository[dbm.TeamCache]
	db *gorm.DB
}

func NewTeamCacheRepository(db *gorm.DB) *TeamCacheRepository {
	return &TeamCacheRepository{
		BaseRepository: repository.NewBaseRepository[dbm.TeamCache](db),
		db:             db,
	}
}

func (r *TeamCacheRepository) FindByID(ctx context.Context, id uint64) (*dbm.TeamCache, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

func (r *TeamCacheRepository) GetByTenantAndTeamID(ctx context.Context, tenantID, teamID uint64) (*dbm.TeamCache, error) {
	cond := map[string]interface{}{"tenant_id": tenantID, "team_id": teamID}
	return r.BaseRepository.GetByCondition(ctx, cond, nil)
}

// Upsert：基于 (tenant_id, team_id)
func (r *TeamCacheRepository) UpsertByTenantAndTeamID(ctx context.Context, obj *dbm.TeamCache) (*dbm.TeamCache, error) {
	unique := []clause.Column{{Name: "tenant_id"}, {Name: "team_id"}}
	return r.BaseRepository.Upsert(ctx, obj, unique)
}
