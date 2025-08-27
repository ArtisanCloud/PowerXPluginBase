package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/powerx-plugins/scrum/internal/db"
	"github.com/powerx-plugins/scrum/internal/logger"
	"gorm.io/gorm"
)

// taskRepository 任务仓储实现
type taskRepository struct{}

// NewTaskRepository 创建任务仓储实例
func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

// Create 创建任务
func (r *taskRepository) Create(ctx context.Context, tdb *db.TenantDB, task *Task) error {
	log := logger.RepoLogger("task").WithContext(ctx)

	if err := task.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := tdb.WithContext(ctx).Create(task).Error; err != nil {
		log.WithError(err).Error("Failed to create task")
		return fmt.Errorf("failed to create task: %w", err)
	}

	log.WithField("task_id", task.ID).Info("Task created successfully")
	return nil
}

// GetByID 根据 ID 获取任务
func (r *taskRepository) GetByID(ctx context.Context, tdb *db.TenantDB, id uint) (*Task, error) {
	log := logger.RepoLogger("task").WithContext(ctx)

	var task Task
	err := tdb.WithContext(ctx).First(&task, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("task not found")
		}
		log.WithError(err).WithField("task_id", id).Error("Failed to get task")
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// Update 更新任务
func (r *taskRepository) Update(ctx context.Context, tdb *db.TenantDB, task *Task) error {
	log := logger.RepoLogger("task").WithContext(ctx)

	if err := task.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if err := tdb.WithContext(ctx).Save(task).Error; err != nil {
		log.WithError(err).WithField("task_id", task.ID).Error("Failed to update task")
		return fmt.Errorf("failed to update task: %w", err)
	}

	log.WithField("task_id", task.ID).Info("Task updated successfully")
	return nil
}

// Delete 删除任务
func (r *taskRepository) Delete(ctx context.Context, tdb *db.TenantDB, id uint) error {
	log := logger.RepoLogger("task").WithContext(ctx)

	result := tdb.WithContext(ctx).Delete(&Task{}, id)
	if result.Error != nil {
		log.WithError(result.Error).WithField("task_id", id).Error("Failed to delete task")
		return fmt.Errorf("failed to delete task: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("task not found")
	}

	log.WithField("task_id", id).Info("Task deleted successfully")
	return nil
}

// List 获取任务列表
func (r *taskRepository) List(ctx context.Context, tdb *db.TenantDB, opts *TaskListOptions) ([]*Task, int64, error) {
	log := logger.RepoLogger("task").WithContext(ctx)

	if opts == nil {
		opts = &TaskListOptions{Page: 1, Limit: 20}
	}

	// 构建查询
	query := tdb.WithContext(ctx).Model(&Task{})

	// 应用过滤条件
	r.applyTaskFilters(query, opts)

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.WithError(err).Error("Failed to count tasks")
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// 应用排序和分页
	r.applyTaskSorting(query, opts)
	r.applyPagination(query, opts.Page, opts.Limit)

	// 执行查询
	var tasks []*Task
	if err := query.Find(&tasks).Error; err != nil {
		log.WithError(err).Error("Failed to list tasks")
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	log.WithField("count", len(tasks)).Info("Tasks listed successfully")
	return tasks, total, nil
}

// applyTaskFilters 应用任务过滤条件
func (r *taskRepository) applyTaskFilters(query *gorm.DB, opts *TaskListOptions) {
	if opts.Status != nil {
		query.Where("status = ?", *opts.Status)
	}

	if opts.Priority != nil {
		query.Where("priority = ?", *opts.Priority)
	}

	if opts.Assignee != nil {
		query.Where("assignee = ?", *opts.Assignee)
	}

	if opts.SprintID != nil {
		query.Where("sprint_id = ?", *opts.SprintID)
	}

	if len(opts.Labels) > 0 {
		// 使用 PostgreSQL 的 JSON 操作符
		for _, label := range opts.Labels {
			query.Where("labels ? ?", label)
		}
	}

	if opts.Search != "" {
		search := "%" + strings.ToLower(opts.Search) + "%"
		query.Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search)
	}

	if opts.DueBefore != nil {
		query.Where("due_date <= ?", *opts.DueBefore)
	}

	if opts.DueAfter != nil {
		query.Where("due_date >= ?", *opts.DueAfter)
	}
}

// applyTaskSorting 应用任务排序
func (r *taskRepository) applyTaskSorting(query *gorm.DB, opts *TaskListOptions) {
	sortBy := opts.SortBy
	if sortBy == "" {
		sortBy = "created_at"
	}

	sortOrder := strings.ToUpper(opts.SortOrder)
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
		query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		query.Order("created_at DESC")
	}
}

// applyPagination 应用分页
func (r *taskRepository) applyPagination(query *gorm.DB, page, limit int) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit
	query.Offset(offset).Limit(limit)
}

// GetBySprintID 根据 Sprint ID 获取任务
func (r *taskRepository) GetBySprintID(ctx context.Context, tdb *db.TenantDB, sprintID uint) ([]*Task, error) {
	var tasks []*Task
	err := tdb.WithContext(ctx).Where("sprint_id = ?", sprintID).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by sprint: %w", err)
	}
	return tasks, nil
}

// GetByAssignee 根据分配人获取任务
func (r *taskRepository) GetByAssignee(ctx context.Context, tdb *db.TenantDB, assignee int64) ([]*Task, error) {
	var tasks []*Task
	err := tdb.WithContext(ctx).Where("assignee = ?", assignee).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by assignee: %w", err)
	}
	return tasks, nil
}

// GetByStatus 根据状态获取任务
func (r *taskRepository) GetByStatus(ctx context.Context, tdb *db.TenantDB, status TaskStatus) ([]*Task, error) {
	var tasks []*Task
	err := tdb.WithContext(ctx).Where("status = ?", status).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	return tasks, nil
}

// Search 搜索任务
func (r *taskRepository) Search(ctx context.Context, tdb *db.TenantDB, query string) ([]*Task, error) {
	var tasks []*Task
	search := "%" + strings.ToLower(query) + "%"
	err := tdb.WithContext(ctx).Where("LOWER(title) LIKE ? OR LOWER(description) LIKE ?", search, search).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}
	return tasks, nil
}

// CountByStatus 按状态统计任务数量
func (r *taskRepository) CountByStatus(ctx context.Context, tdb *db.TenantDB, status TaskStatus) (int64, error) {
	var count int64
	err := tdb.WithContext(ctx).Model(&Task{}).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by status: %w", err)
	}
	return count, nil
}

// CountBySprintID 按 Sprint 统计任务数量
func (r *taskRepository) CountBySprintID(ctx context.Context, tdb *db.TenantDB, sprintID uint) (int64, error) {
	var count int64
	err := tdb.WithContext(ctx).Model(&Task{}).Where("sprint_id = ?", sprintID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count tasks by sprint: %w", err)
	}
	return count, nil
}

// GetOverdueTasks 获取过期任务
func (r *taskRepository) GetOverdueTasks(ctx context.Context, tdb *db.TenantDB) ([]*Task, error) {
	var tasks []*Task
	now := time.Now()
	err := tdb.WithContext(ctx).Where("due_date < ? AND status != ?", now, TaskStatusDone).Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	return tasks, nil
}

// BatchUpdateStatus 批量更新状态
func (r *taskRepository) BatchUpdateStatus(ctx context.Context, tdb *db.TenantDB, ids []uint, status TaskStatus) error {
	if len(ids) == 0 {
		return nil
	}

	err := tdb.WithContext(ctx).Model(&Task{}).Where("id IN ?", ids).Update("status", status).Error
	if err != nil {
		return fmt.Errorf("failed to batch update status: %w", err)
	}
	return nil
}

// BatchDelete 批量删除
func (r *taskRepository) BatchDelete(ctx context.Context, tdb *db.TenantDB, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	err := tdb.WithContext(ctx).Delete(&Task{}, ids).Error
	if err != nil {
		return fmt.Errorf("failed to batch delete tasks: %w", err)
	}
	return nil
}
