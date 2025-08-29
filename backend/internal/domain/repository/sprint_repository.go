package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scrum-plugin/internal/domain/models"

	"gorm.io/gorm"
)

// SprintRepository Sprint 仓储实现
type SprintRepository struct {
	BaseRepository[models.Sprint]
	db *gorm.DB
}

// NewSprintRepository 创建 Sprint 仓储实例
func NewSprintRepository(db *gorm.DB) *SprintRepository {
	return &SprintRepository{
		BaseRepository: *NewBaseRepository[models.Sprint](db),
		db:             db,
	}
}

// SprintListOptions Sprint 列表查询选项
type SprintListOptions struct {
	Page      int
	Limit     int
	Status    *models.SprintStatus
	Name      string
	Search    string
	SortBy    string
	SortOrder string
}

// List 获取 Sprint 列表（带复杂过滤条件）
func (r *SprintRepository) List(ctx context.Context, opts *SprintListOptions) ([]*models.Sprint, int64, error) {
	if opts == nil {
		opts = &SprintListOptions{Page: 1, Limit: 20}
	}

	// 使用回调函数来应用复杂的过滤条件
	callback := func(db *gorm.DB, opt interface{}) *gorm.DB {
		options, ok := opt.(*SprintListOptions)
		if !ok {
			return db
		}

		// 应用过滤条件
		if options.Status != nil {
			db = db.Where("status = ?", options.Status)
		}

		if options.Name != "" {
			db = db.Where("name = ?", options.Name)
		}

		if options.Search != "" {
			search := "%" + strings.ToLower(options.Search) + "%"
			db = db.Where("LOWER(name) LIKE ? OR LOWER(goal) LIKE ?", search, search)
		}

		// 应用排序
		sortBy := options.SortBy
		if sortBy == "" {
			sortBy = "created_at"
		}

		sortOrder := strings.ToUpper(options.SortOrder)
		if sortOrder != "ASC" && sortOrder != "DESC" {
			sortOrder = "DESC"
		}

		// 验证排序字段
		validSortFields := map[string]bool{
			"created_at": true,
			"updated_at": true,
			"name":       true,
			"start_date": true,
			"end_date":   true,
			"status":     true,
		}

		if validSortFields[sortBy] {
			db = db.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
		} else {
			db = db.Order("created_at DESC")
		}

		return db
	}

	result, err := r.BaseRepository.FindByCondition(ctx, map[string]interface{}{}, opts.Page, opts.Limit, callback, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sprints: %w", err)
	}

	return result.List, result.Total, nil
}

// GetActiveSprint 获取当前活跃的 Sprint
func (r *SprintRepository) GetActiveSprint(ctx context.Context) (*models.Sprint, error) {
	var sprint models.Sprint
	err := r.db.WithContext(ctx).Where("status = ?", models.SprintStatusActive).First(&sprint).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active sprint: %w", err)
	}
	return &sprint, nil
}

// GetByDateRange 根据日期范围获取 Sprint
func (r *SprintRepository) GetByDateRange(ctx context.Context, startDate, endDate time.Time) ([]*models.Sprint, error) {
	var sprints []*models.Sprint
	err := r.db.WithContext(ctx).
		Where("start_date <= ? AND end_date >= ?", endDate, startDate).
		Find(&sprints).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get sprints by date range: %w", err)
	}
	return sprints, nil
}
