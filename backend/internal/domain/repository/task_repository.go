package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"scrum-plugin/internal/domain/models"

	"gorm.io/gorm"
)

// TaskRepository 任务仓储实现
type TaskRepository struct {
	BaseRepository[models.Task]
	db *gorm.DB
}

// NewTaskRepository 创建任务仓储实例
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		BaseRepository: *NewBaseRepository[models.Task](db),
		db:             db,
	}
}

// TaskListOptions 任务列表查询选项
type TaskListOptions struct {
	Page      int
	Limit     int
	Status    *models.TaskStatus
	Priority  *models.Priority
	Assignee  *int64
	SprintID  *uint
	Labels    []string
	Search    string
	DueBefore *time.Time
	DueAfter  *time.Time
	SortBy    string
	SortOrder string
}

// 业务特定方法

// GetBySprintID 根据 Sprint ID 获取任务
func (r *TaskRepository) GetBySprintID(ctx context.Context, sprintID uint) ([]*models.Task, error) {
	conditions := map[string]interface{}{
		"sprint_id": sprintID,
	}

	result, err := r.BaseRepository.FindByCondition(ctx, conditions, 1, 1000, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by sprint: %w", err)
	}

	return result.List, nil
}

// GetByAssignee 根据分配人获取任务
func (r *TaskRepository) GetByAssignee(ctx context.Context, assignee int64) ([]*models.Task, error) {
	conditions := map[string]interface{}{
		"assignee": assignee,
	}

	result, err := r.BaseRepository.FindByCondition(ctx, conditions, 1, 1000, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by assignee: %w", err)
	}

	return result.List, nil
}

// GetByStatus 根据状态获取任务
func (r *TaskRepository) GetByStatus(ctx context.Context, status models.TaskStatus) ([]*models.Task, error) {
	conditions := map[string]interface{}{
		"status": status,
	}

	result, err := r.BaseRepository.FindByCondition(ctx, conditions, 1, 1000, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}

	return result.List, nil
}

// List 获取任务列表（带复杂过滤条件）
func (r *TaskRepository) List(ctx context.Context, opts *TaskListOptions) ([]*models.Task, int64, error) {
	if opts == nil {
		opts = &TaskListOptions{Page: 1, Limit: 20}
	}

	// 使用回调函数来应用复杂的过滤条件
	callback := func(db *gorm.DB, opt interface{}) *gorm.DB {
		options, ok := opt.(*TaskListOptions)
		if !ok {
			return db
		}

		// 应用过滤条件
		if options.Status != nil {
			db = db.Where("status = ?", options.Status)
		}

		if options.Priority != nil {
			db = db.Where("priority = ?", options.Priority)
		}

		if options.Assignee != nil {
			db = db.Where("assignee = ?", *options.Assignee)
		}

		if options.SprintID != nil {
			db = db.Where("sprint_id = ?", *options.SprintID)
		}

		if len(options.Labels) > 0 {
			// 使用 PostgreSQL 的 JSON 操作符
			for _, label := range options.Labels {
				db = db.Where("labels ? ?", label)
			}
		}

		if options.Search != "" {
			search := "%" + strings.ToLower(options.Search) + "%"
			db = db.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
		}

		if options.DueBefore != nil {
			db = db.Where("due_date <= ?", *options.DueBefore)
		}

		if options.DueAfter != nil {
			db = db.Where("due_date >= ?", *options.DueAfter)
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
			"title":      true,
			"priority":   true,
			"due_date":   true,
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
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	return result.List, result.Total, nil
}

// Search 搜索任务
func (r *TaskRepository) Search(ctx context.Context, query string) ([]*models.Task, error) {
	callback := func(db *gorm.DB, opt interface{}) *gorm.DB {
		search := "%" + strings.ToLower(query) + "%"
		return db.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	result, err := r.BaseRepository.FindByCondition(ctx, map[string]interface{}{}, 1, 1000, callback, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}

	return result.List, nil
}

// CountByStatus 按状态统计任务数量
func (r *TaskRepository) CountByStatus(ctx context.Context, status models.TaskStatus) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Task{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by status: %w", err)
	}
	return count, nil
}

// CountBySprintID 按 Sprint 统计任务数量
func (r *TaskRepository) CountBySprintID(ctx context.Context, sprintID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Task{}).Where("sprint_id = ?", sprintID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by sprint: %w", err)
	}
	return count, nil
}

// GetOverdueTasks 获取过期任务
func (r *TaskRepository) GetOverdueTasks(ctx context.Context) ([]*models.Task, error) {
	callback := func(db *gorm.DB, opt interface{}) *gorm.DB {
		now := time.Now()
		return db.Where("due_date < ? AND status != ?", now, models.TaskStatusDone)
	}

	result, err := r.BaseRepository.FindByCondition(ctx, map[string]interface{}{}, 1, 1000, callback, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}

	return result.List, nil
}

// BatchUpdateStatus 批量更新状态
func (r *TaskRepository) BatchUpdateStatus(ctx context.Context, ids []uint, status models.TaskStatus) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Model(&models.Task{}).Where("id IN ?", ids).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("failed to batch update status: %w", err)
	}
	return nil
}

// BatchDelete 批量删除
func (r *TaskRepository) BatchDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.WithContext(ctx).Delete(&models.Task{}, ids).Error
	if err != nil {
		return fmt.Errorf("failed to batch delete tasks: %w", err)
	}
	return nil
}

// UpdateStatus 更新任务状态
func (r *TaskRepository) UpdateStatus(ctx context.Context, id uint, status models.TaskStatus) (*models.Task, error) {
	fields := map[string]interface{}{
		"status": status,
	}

	result, err := r.BaseRepository.UpdateByID(ctx, uint64(id), fields, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update task status: %w", err)
	}

	return result, nil
}

// GetTasksWithActivities 获取带活动记录的任务
func (r *TaskRepository) GetTasksWithActivities(ctx context.Context, ids []uint) ([]*models.Task, error) {
	callback := func(db *gorm.DB, opt interface{}) *gorm.DB {
		return db.Preload("Activities").Where("id IN ?", ids)
	}

	result, err := r.BaseRepository.FindByCondition(ctx, map[string]interface{}{}, 1, 1000, callback, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks with activities: %w", err)
	}

	return result.List, nil
}
