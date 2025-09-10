package notes

import (
	"context"
	"strings"

	dbm "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/note"
	"github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository"
	nrepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/note"
	"gorm.io/gorm"
)

type NoteService struct {
	NoteRepo *nrepo.NoteRepository
}

func NewNoteService(db *gorm.DB) *NoteService {
	return &NoteService{NoteRepo: nrepo.NewNoteRepository(db)}
}

// 列表：分页 + 关键词 + 可选 team/member 过滤
func (s *NoteService) List(
	ctx context.Context,
	q string,
	teamID, memberID uint64,
	page, pageSize int,
) (*repository.Page[[]*dbm.Note], error) {

	conds := map[string]interface{}{}

	// 建议与 RLS 并用；明确添加 tenant_id 能帮助优化器挑选更优索引
	if tid, ok, _ := s.NoteRepo.CurrentTenantID(ctx); ok {
		conds["tenant_id = ?"] = tid
	}
	if teamID > 0 {
		conds["team_id = ?"] = teamID
	}
	if memberID > 0 {
		conds["member_id = ?"] = memberID
	}

	cb := func(db *gorm.DB, opt interface{}) *gorm.DB {
		if kw, _ := opt.(string); strings.TrimSpace(kw) != "" {
			p := "%" + strings.TrimSpace(kw) + "%"
			db = db.Where("(title ILIKE ? OR content ILIKE ?)", p, p)
		}
		return db.Order("id DESC")
	}

	return s.NoteRepo.FindPage(ctx, conds, page, pageSize, cb, q)
}

// 详情
func (s *NoteService) GetByID(ctx context.Context, id uint64) (*dbm.Note, error) {
	n, err := s.NoteRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return n, nil
}

// 新建
func (s *NoteService) Create(
	ctx context.Context,
	title, content, author string,
	teamID, memberID uint64,
) (*dbm.Note, error) {

	tid, ok, err := s.NoteRepo.CurrentTenantID(ctx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, gorm.ErrInvalidData
	}

	n := &dbm.Note{
		Title:    title,
		Content:  content,
		Author:   author,
		TeamID:   teamID,
		MemberID: memberID,
	}
	n.TenantID = tid

	return s.NoteRepo.Create(ctx, n)
}

// 更新（按 ID 局部更新即可，避免先查后存的一次往返）
func (s *NoteService) Update(
	ctx context.Context,
	id uint64,
	title, content, author string,
	teamID, memberID uint64,
) (*dbm.Note, error) {

	fields := map[string]interface{}{
		"title":     title,
		"content":   content,
		"author":    author,
		"team_id":   teamID,
		"member_id": memberID,
	}
	return s.NoteRepo.UpdateByID(ctx, id, fields)
}

// 删除（软删）
func (s *NoteService) Delete(ctx context.Context, id uint64) error {
	return s.NoteRepo.DeleteByID(ctx, id)
}
