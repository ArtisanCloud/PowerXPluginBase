package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"scrum-plugin/internal/domain/models"
	"scrum-plugin/internal/domain/repository"
	"scrum-plugin/internal/logger"

	"gorm.io/gorm"
)

// 请求类型定义

// CreateTaskRequest 创建任务请求
type CreateTaskRequest struct {
	Title       string                 `json:"title" binding:"required"`
	Description *string                `json:"description"`
	Status      string                 `json:"status"`
	Priority    string                 `json:"priority"`
	Assignee    *uint64                `json:"assignee"`
	SprintID    *uint64                `json:"sprint_id"`
	Labels      []string               `json:"labels"`
	DueDate     *time.Time             `json:"due_date"`
	Estimate    *int                   `json:"estimate"`
	Meta        map[string]interface{} `json:"meta"`
}

// UpdateTaskRequest 更新任务请求
type UpdateTaskRequest struct {
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	Status      *string                `json:"status"`
	Priority    *string                `json:"priority"`
	Assignee    *uint64                `json:"assignee"`
	SprintID    *uint64                `json:"sprint_id"`
	Labels      []string               `json:"labels"`
	DueDate     *time.Time             `json:"due_date"`
	Estimate    *int                   `json:"estimate"`
	Meta        map[string]interface{} `json:"meta"`
}

// TaskListOptions 任务列表查询选项
type TaskListOptions struct {
	Page      int
	Limit     int
	Status    *string
	Priority  *string
	Assignee  *int64
	SprintID  *uint
	Labels    []string
	Search    string
	SortBy    string
	SortOrder string
}

// TaskService 任务服务结构体
type TaskService struct {
	taskRepo *repository.TaskRepository
}

// NewTaskService 创建任务服务实例
func NewTaskService(db *gorm.DB) *TaskService {
	taskRepo := repository.NewTaskRepository(db)
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// CreateTask 创建任务
func (s *TaskService) CreateTask(ctx context.Context, tenantID int64, req *CreateTaskRequest) (*models.Task, error) {
	log := logger.ServiceLogger("task").WithContext(ctx)

	// 创建领域模型
	task := &models.Task{
		Title:       req.Title,
		TaskType:    models.TaskTypeUserStory, // 默认类型
		Status:      models.TaskStatus(req.Status),
		Priority:    models.Priority(req.Priority),
		AssigneeID:  req.Assignee,
		SprintID:    req.SprintID,
		Labels:      models.Labels(req.Labels),
		DueDate:     req.DueDate,
		StoryPoints: req.Estimate,
	}

	// 处理可选字段
	if req.Description != nil {
		task.Description = *req.Description
	}

	// 处理 Meta 字段
	if req.Meta != nil {
		metaBytes, err := json.Marshal(req.Meta)
		if err != nil {
			log.WithError(err).Error("Failed to marshal meta data")
			return nil, fmt.Errorf("failed to marshal meta data: %w", err)
		}
		task.Meta = metaBytes
	}

	// 设置默认值
	if task.Status == "" {
		task.Status = models.TaskStatusTodo
	}
	if task.Priority == "" {
		task.Priority = models.PriorityMedium
	}
	if task.Labels == nil {
		task.Labels = models.Labels{}
	}

	// 数据验证
	if err := task.Validate(); err != nil {
		log.WithError(err).Error("Task validation failed")
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// 创建任务
	result, err := s.taskRepo.Create(ctx, task)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"title": task.Title,
			"type":  task.TaskType,
		}).Error("Failed to create task")
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	log.WithFields(logger.Fields{
		"task_id": result.ID,
		"title":   result.Title,
	}).Info("Task created successfully")

	return result, nil
}

// GetTask 根据ID获取任务
func (s *TaskService) GetTask(ctx context.Context, tenantID int64, id uint64) (*models.Task, error) {
	log := logger.ServiceLogger("task").WithContext(ctx)

	task, err := s.taskRepo.GetById(ctx, uint64(id), nil)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"task_id": id,
		}).Error("Failed to get task by ID")
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	if task == nil {
		log.WithFields(logger.Fields{
			"task_id": id,
		}).Warn("Task not found")
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

// UpdateTask 更新任务
func (s *TaskService) UpdateTask(ctx context.Context, tenantID int64, id uint64, req *UpdateTaskRequest) (*models.Task, error) {
	log := logger.ServiceLogger("task").WithContext(ctx)

	// 获取原任务
	originalTask, err := s.GetTask(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// 应用更新
	if req.Title != nil {
		originalTask.Title = *req.Title
	}
	if req.Description != nil {
		originalTask.Description = *req.Description
	}
	if req.Status != nil {
		originalTask.Status = models.TaskStatus(*req.Status)
	}
	if req.Priority != nil {
		originalTask.Priority = models.Priority(*req.Priority)
	}
	if req.Assignee != nil {
		originalTask.AssigneeID = req.Assignee
	}
	if req.SprintID != nil {
		originalTask.SprintID = req.SprintID
	}
	if req.Labels != nil {
		originalTask.Labels = models.Labels(req.Labels)
	}
	if req.DueDate != nil {
		originalTask.DueDate = req.DueDate
	}
	if req.Estimate != nil {
		originalTask.StoryPoints = req.Estimate
	}
	if req.Meta != nil {
		metaBytes, err := json.Marshal(req.Meta)
		if err != nil {
			log.WithError(err).Error("Failed to marshal meta data")
			return nil, fmt.Errorf("failed to marshal meta data: %w", err)
		}
		originalTask.Meta = metaBytes
	}

	// 数据验证
	if err := originalTask.Validate(); err != nil {
		log.WithError(err).Error("Task validation failed")
		return nil, fmt.Errorf("task validation failed: %w", err)
	}

	// 更新任务
	result, err := s.taskRepo.Update(ctx, originalTask)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"task_id": originalTask.ID,
			"title":   originalTask.Title,
		}).Error("Failed to update task")
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	log.WithFields(logger.Fields{
		"task_id": result.ID,
		"title":   result.Title,
	}).Info("Task updated successfully")

	return result, nil
}

// DeleteTask 删除任务
func (s *TaskService) DeleteTask(ctx context.Context, tenantID int64, id uint64) error {
	log := logger.ServiceLogger("task").WithContext(ctx)

	// 检查任务是否存在
	task, err := s.GetTask(ctx, tenantID, id)
	if err != nil {
		return err
	}

	// 删除任务
	_, err = s.taskRepo.Delete(ctx, nil, &models.Task{BaseModel: models.BaseModel{ID: id}}, true)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"task_id": id,
			"title":   task.Title,
		}).Error("Failed to delete task")
		return fmt.Errorf("failed to delete task: %w", err)
	}

	log.WithFields(logger.Fields{
		"task_id": id,
		"title":   task.Title,
	}).Info("Task deleted successfully")

	return nil
}

// ListTasks 获取任务列表
func (s *TaskService) ListTasks(ctx context.Context, tenantID int64, opts *TaskListOptions) ([]*models.Task, int64, error) {
	log := logger.ServiceLogger("task").WithContext(ctx)

	// 转换为仓储层选项
	repoOpts := &repository.TaskListOptions{
		Page:      opts.Page,
		Limit:     opts.Limit,
		Labels:    opts.Labels,
		Search:    opts.Search,
		SortBy:    opts.SortBy,
		SortOrder: opts.SortOrder,
	}

	if opts.Status != nil {
		status := models.TaskStatus(*opts.Status)
		repoOpts.Status = &status
	}

	if opts.Priority != nil {
		priority := models.Priority(*opts.Priority)
		repoOpts.Priority = &priority
	}

	if opts.Assignee != nil {
		repoOpts.Assignee = opts.Assignee
	}

	if opts.SprintID != nil {
		repoOpts.SprintID = opts.SprintID
	}

	tasks, total, err := s.taskRepo.List(ctx, repoOpts)
	if err != nil {
		log.WithError(err).Error("Failed to list tasks")
		return nil, 0, fmt.Errorf("failed to list tasks: %w", err)
	}

	log.WithFields(logger.Fields{
		"count": len(tasks),
		"total": total,
	}).Debug("Tasks listed successfully")

	return tasks, total, nil
}

// UpdateTaskStatus 更新任务状态
func (s *TaskService) UpdateTaskStatus(ctx context.Context, tenantID int64, id uint, status string) error {
	log := logger.ServiceLogger("task").WithContext(ctx)

	// 转换为领域模型状态
	domainStatus := models.TaskStatus(status)

	_, err := s.taskRepo.UpdateStatus(ctx, id, domainStatus)
	if err != nil {
		log.WithError(err).WithFields(logger.Fields{
			"task_id": id,
			"status":  status,
		}).Error("Failed to update task status")
		return fmt.Errorf("failed to update task status: %w", err)
	}

	log.WithFields(logger.Fields{
		"task_id": id,
		"status":  status,
	}).Info("Task status updated successfully")

	return nil
}
