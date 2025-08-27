package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/powerx-plugins/scrum/internal/config"
	"github.com/powerx-plugins/scrum/internal/db"
	"github.com/powerx-plugins/scrum/internal/domain"
	"github.com/powerx-plugins/scrum/internal/handlers"
	"github.com/powerx-plugins/scrum/internal/logger"
	"github.com/powerx-plugins/scrum/internal/router"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	logger.Init(cfg.LogLevel)
	logger.Info("Starting PowerX Scrum Plugin...")

	// 连接数据库
	if err := db.Connect(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	// 注意：不再在应用启动时执行任何迁移操作
	// 数据库迁移应该通过独立的 cmd/database/migrate 命令执行

	// 初始化依赖
	taskRepo := domain.NewTaskRepository()
	taskService := NewTaskService(taskRepo)

	// 初始化处理器
	adminHandler := handlers.NewAdminHandler()
	taskHandler := handlers.NewTaskHandler(taskService)
	healthHandler := handlers.NewHealthHandler()

	// 设置路由
	r := router.New(cfg, adminHandler, taskHandler, healthHandler)
	engine := r.Setup()

	// 创建 HTTP 服务器
	server := &http.Server{
		Addr:    cfg.BindAddr,
		Handler: engine,

		// 超时配置
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		IdleTimeout:       120 * time.Second,

		// 最大头部大小
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		logger.WithField("addr", cfg.BindAddr).Info("Starting HTTP server...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("Server forced to shutdown")
	} else {
		logger.Info("Server shutdown completed")
	}
}

// 注意：数据库迁移相关功能已移除
// 请使用独立的迁移命令：go run cmd/database/migrate/migrate.go

// NewTaskService 创建任务服务实例（简化实现）
func NewTaskService(repo domain.TaskRepository) domain.TaskService {
	return &taskService{repo: repo}
}

// taskService 任务服务简化实现
type taskService struct {
	repo domain.TaskRepository
}

// CreateTask 创建任务
func (s *taskService) CreateTask(ctx context.Context, tenantID int64, req *domain.CreateTaskRequest) (*domain.Task, error) {
	// 开始租户事务
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tdb.Rollback()

	// 创建任务模型
	task := &domain.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Assignee:    req.Assignee,
		SprintID:    req.SprintID,
		Labels:      domain.Labels(req.Labels),
		DueDate:     req.DueDate,
		Estimate:    req.Estimate,
		Meta:        domain.Meta(req.Meta),
	}

	// 设置默认值
	if task.Status == "" {
		task.Status = domain.TaskStatusTodo
	}
	if task.Priority == "" {
		task.Priority = domain.PriorityMedium
	}

	// 创建任务
	if err := s.repo.Create(ctx, tdb, task); err != nil {
		return nil, err
	}

	// 提交事务
	if err := tdb.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask 获取任务
func (s *taskService) GetTask(ctx context.Context, tenantID int64, id uint) (*domain.Task, error) {
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tdb.Rollback()

	task, err := s.repo.GetByID(ctx, tdb, id)
	if err != nil {
		return nil, err
	}

	tdb.Commit()
	return task, nil
}

// UpdateTask 更新任务
func (s *taskService) UpdateTask(ctx context.Context, tenantID int64, id uint, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	defer tdb.Rollback()

	// 获取现有任务
	task, err := s.repo.GetByID(ctx, tdb, id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.Assignee != nil {
		task.Assignee = req.Assignee
	}
	if req.SprintID != nil {
		task.SprintID = req.SprintID
	}
	if req.Labels != nil {
		task.Labels = domain.Labels(req.Labels)
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.Estimate != nil {
		task.Estimate = req.Estimate
	}
	if req.Meta != nil {
		task.Meta = domain.Meta(req.Meta)
	}

	// 保存更新
	if err := s.repo.Update(ctx, tdb, task); err != nil {
		return nil, err
	}

	if err := tdb.Commit(); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask 删除任务
func (s *taskService) DeleteTask(ctx context.Context, tenantID int64, id uint) error {
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tdb.Rollback()

	if err := s.repo.Delete(ctx, tdb, id); err != nil {
		return err
	}

	return tdb.Commit()
}

// ListTasks 获取任务列表
func (s *taskService) ListTasks(ctx context.Context, tenantID int64, opts *domain.TaskListOptions) ([]*domain.Task, int64, error) {
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return nil, 0, err
	}
	defer tdb.Rollback()

	tasks, total, err := s.repo.List(ctx, tdb, opts)
	if err != nil {
		return nil, 0, err
	}

	tdb.Commit()
	return tasks, total, nil
}

// UpdateTaskStatus 更新任务状态
func (s *taskService) UpdateTaskStatus(ctx context.Context, tenantID int64, id uint, status domain.TaskStatus) error {
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tdb.Rollback()

	task, err := s.repo.GetByID(ctx, tdb, id)
	if err != nil {
		return err
	}

	task.Status = status
	if err := s.repo.Update(ctx, tdb, task); err != nil {
		return err
	}

	return tdb.Commit()
}

// 其他方法的简化实现（占位）
func (s *taskService) SearchTasks(ctx context.Context, tenantID int64, query string) ([]*domain.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *taskService) AssignTask(ctx context.Context, tenantID int64, id uint, assignee int64) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) UnassignTask(ctx context.Context, tenantID int64, id uint) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) AddTaskToSprint(ctx context.Context, tenantID int64, taskID, sprintID uint) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) RemoveTaskFromSprint(ctx context.Context, tenantID int64, taskID uint) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) AddTaskLabel(ctx context.Context, tenantID int64, id uint, label string) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) RemoveTaskLabel(ctx context.Context, tenantID int64, id uint, label string) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) BatchUpdateTaskStatus(ctx context.Context, tenantID int64, ids []uint, status domain.TaskStatus) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) BatchDeleteTasks(ctx context.Context, tenantID int64, ids []uint) error {
	return fmt.Errorf("not implemented")
}

func (s *taskService) GetTaskStatsByStatus(ctx context.Context, tenantID int64) (map[domain.TaskStatus]int64, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *taskService) GetOverdueTasks(ctx context.Context, tenantID int64) ([]*domain.Task, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *taskService) GetTasksByAssignee(ctx context.Context, tenantID int64, assignee int64) ([]*domain.Task, error) {
	return nil, fmt.Errorf("not implemented")
}
