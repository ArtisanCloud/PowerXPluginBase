package template

import (
	"context"
	"database/sql"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/template"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
)

type TemplateRepository struct {
	*repository.BaseRepository[dbm.Template]
}

func NewTemplateRepository(db *gorm.DB) *TemplateRepository {
	return &TemplateRepository{
		BaseRepository: repository.NewBaseRepository[dbm.Template](db),
	}
}

func (r *TemplateRepository) FindByID(ctx context.Context, id uint64) (*dbm.Template, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

func (r *TemplateRepository) Create(ctx context.Context, t *dbm.Template) (*dbm.Template, error) {
	return r.BaseRepository.Create(ctx, t)
}

func (r *TemplateRepository) UpdateByID(ctx context.Context, id uint64, fields map[string]interface{}) (*dbm.Template, error) {
	return r.BaseRepository.UpdateByID(ctx, id, fields, nil)
}

func (r *TemplateRepository) DeleteByID(ctx context.Context, id uint64) error {
	_, err := r.BaseRepository.Delete(ctx, map[string]interface{}{"id": id}, nil, true)
	return err
}

func (r *TemplateRepository) FindPage(
	ctx context.Context,
	conditions map[string]interface{},
	page, pageSize int,
	cb func(*gorm.DB, interface{}) *gorm.DB,
	opt interface{},
) (*repository.Page[[]*dbm.Template], error) {
	return r.BaseRepository.FindByCondition(ctx, conditions, page, pageSize, cb, opt)
}

func (r *TemplateRepository) CurrentTenantID(ctx context.Context) (uint64, bool, error) {
	var tid sql.NullInt64
	err := r.DB.WithContext(ctx).
		Raw(`SELECT current_setting('app.tenant_id', true)::bigint`).Scan(&tid).Error
	if err != nil {
		return 0, false, err
	}
	if !tid.Valid {
		return 0, false, nil
	}
	return uint64(tid.Int64), true, nil
}
