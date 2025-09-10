package note

import (
	"context"
	"database/sql"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/note"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	"gorm.io/gorm"
)

type NoteRepository struct {
	*repository.BaseRepository[dbm.Note]
}

func NewNoteRepository(db *gorm.DB) *NoteRepository {
	return &NoteRepository{
		BaseRepository: repository.NewBaseRepository[dbm.Note](db),
	}
}

// 单条
func (r *NoteRepository) FindByID(ctx context.Context, id uint64) (*dbm.Note, error) {
	return r.BaseRepository.GetById(ctx, id, nil)
}

// 创建
func (r *NoteRepository) Create(ctx context.Context, n *dbm.Note) (*dbm.Note, error) {
	return r.BaseRepository.Create(ctx, n)
}

// 按 ID 更新（字段白名单在 BaseRepository 里自动控制）
func (r *NoteRepository) UpdateByID(ctx context.Context, id uint64, fields map[string]interface{}) (*dbm.Note, error) {
	return r.BaseRepository.UpdateByID(ctx, id, fields, nil)
}

// 删除（软删）
func (r *NoteRepository) DeleteByID(ctx context.Context, id uint64) error {
	_, err := r.BaseRepository.Delete(ctx, map[string]interface{}{"id": id}, nil, true)
	return err
}

// 分页查询（支持传入 conditions 与 callback）
func (r *NoteRepository) FindPage(
	ctx context.Context,
	conditions map[string]interface{},
	page, pageSize int,
	cb func(*gorm.DB, interface{}) *gorm.DB,
	opt interface{},
) (*repository.Page[[]*dbm.Note], error) {
	return r.BaseRepository.FindByCondition(ctx, conditions, page, pageSize, cb, opt)
}

// 从 RLS 会话（app.tenant_id）读取当前租户；未设置返回 ok=false
func (r *NoteRepository) CurrentTenantID(ctx context.Context) (uint64, bool, error) {
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
