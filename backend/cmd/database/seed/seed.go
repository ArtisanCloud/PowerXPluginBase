package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"scrum-plugin/internal/config"
	"scrum-plugin/internal/db"
	"scrum-plugin/internal/domain"
	"scrum-plugin/internal/logger"
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
	logger.Info("Starting database seeding...")

	// 连接数据库
	if err := db.Connect(cfg); err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.WithError(err).Error("Failed to close database connection")
		}
	}()

	// 运行种子数据
	if err := runSeeds(); err != nil {
		logger.WithError(err).Fatal("Failed to run seeds")
	}

	logger.Info("Database seeding completed successfully")
}

// runSeeds 运行种子数据
func runSeeds() error {
	ctx := context.Background()

	// 创建示例租户数据
	tenants := []int64{1, 2} // 示例租户 ID

	for _, tenantID := range tenants {
		logger.WithField("tenant_id", tenantID).Info("Seeding data for tenant")

		if err := seedTenantData(ctx, tenantID); err != nil {
			return fmt.Errorf("failed to seed data for tenant %d: %w", tenantID, err)
		}
	}

	return nil
}

// seedTenantData 为指定租户创建种子数据
func seedTenantData(ctx context.Context, tenantID int64) error {
	// 开始租户事务
	tdb, err := db.BeginTenantTx(ctx, tenantID)
	if err != nil {
		return err
	}
	defer tdb.Rollback()

	// 创建示例 Sprint
	sprints, err := createSampleSprints(ctx, tdb, tenantID)
	if err != nil {
		return fmt.Errorf("failed to create sample sprints: %w", err)
	}

	// 创建示例任务
	if err := createSampleTasks(ctx, tdb, tenantID, sprints); err != nil {
		return fmt.Errorf("failed to create sample tasks: %w", err)
	}

	// 提交事务
	if err := tdb.Commit(); err != nil {
		return err
	}

	logger.WithField("tenant_id", tenantID).Info("Tenant data seeded successfully")
	return nil
}

// createSampleSprints 创建示例 Sprint
func createSampleSprints(ctx context.Context, tdb *db.TenantDB, tenantID int64) ([]*domain.Sprint, error) {
	now := time.Now()

	sprints := []*domain.Sprint{
		{
			BaseModel: db.BaseModel{TenantID: tenantID},
			Name:      "Sprint 1 - MVP 基础功能",
			Goal:      "完成用户管理和基础任务功能",
			StartDate: now.AddDate(0, 0, -14), // 2周前开始
			EndDate:   now,                    // 现在结束
			Capacity:  func() *int { v := 40; return &v }(),
			Status:    domain.SprintStatusCompleted,
		},
		{
			BaseModel: db.BaseModel{TenantID: tenantID},
			Name:      "Sprint 2 - 高级功能",
			Goal:      "实现 Sprint 管理和报告功能",
			StartDate: now,                   // 现在开始
			EndDate:   now.AddDate(0, 0, 14), // 2周后结束
			Capacity:  func() *int { v := 45; return &v }(),
			Status:    domain.SprintStatusActive,
		},
		{
			BaseModel: db.BaseModel{TenantID: tenantID},
			Name:      "Sprint 3 - 优化和集成",
			Goal:      "性能优化和第三方集成",
			StartDate: now.AddDate(0, 0, 14), // 2周后开始
			EndDate:   now.AddDate(0, 0, 28), // 4周后结束
			Capacity:  func() *int { v := 50; return &v }(),
			Status:    domain.SprintStatusPlanning,
		},
	}

	var createdSprints []*domain.Sprint
	for _, sprint := range sprints {
		if err := tdb.Create(sprint).Error; err != nil {
			return nil, fmt.Errorf("failed to create sprint %s: %w", sprint.Name, err)
		}
		createdSprints = append(createdSprints, sprint)
		logger.WithFields(logger.Fields{
			"sprint_id":   sprint.ID,
			"sprint_name": sprint.Name,
			"tenant_id":   tenantID,
		}).Info("Sprint created")
	}

	return createdSprints, nil
}

// createSampleTasks 创建示例任务
func createSampleTasks(ctx context.Context, tdb *db.TenantDB, tenantID int64, sprints []*domain.Sprint) error {
	// 为第一个 Sprint（已完成）创建任务
	if len(sprints) > 0 {
		sprint1Tasks := []*domain.Task{
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "设计用户注册登录界面",
				Description: "创建用户友好的注册和登录页面",
				Status:      domain.TaskStatusDone,
				Priority:    domain.PriorityHigh,
				Assignee:    func() *int64 { v := int64(101); return &v }(),
				SprintID:    &sprints[0].ID,
				Labels:      domain.Labels{"ui", "frontend", "user-auth"},
				Estimate:    func() *int { v := 8; return &v }(),
				Meta: domain.Meta{
					"difficulty": "medium",
					"category":   "frontend",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "实现用户认证API",
				Description: "开发用户注册、登录、注销的后端API",
				Status:      domain.TaskStatusDone,
				Priority:    domain.PriorityHigh,
				Assignee:    func() *int64 { v := int64(102); return &v }(),
				SprintID:    &sprints[0].ID,
				Labels:      domain.Labels{"api", "backend", "auth"},
				Estimate:    func() *int { v := 13; return &v }(),
				Meta: domain.Meta{
					"difficulty": "high",
					"category":   "backend",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "创建任务基础CRUD功能",
				Description: "实现任务的创建、读取、更新、删除操作",
				Status:      domain.TaskStatusDone,
				Priority:    domain.PriorityMedium,
				Assignee:    func() *int64 { v := int64(103); return &v }(),
				SprintID:    &sprints[0].ID,
				Labels:      domain.Labels{"crud", "task-management", "backend"},
				Estimate:    func() *int { v := 5; return &v }(),
				Meta: domain.Meta{
					"difficulty": "low",
					"category":   "backend",
				},
			},
		}

		for _, task := range sprint1Tasks {
			if err := tdb.Create(task).Error; err != nil {
				return fmt.Errorf("failed to create task %s: %w", task.Title, err)
			}
			logger.WithFields(logger.Fields{
				"task_id":    task.ID,
				"task_title": task.Title,
				"sprint_id":  task.SprintID,
				"status":     task.Status,
			}).Info("Task created")
		}
	}

	// 为第二个 Sprint（进行中）创建任务
	if len(sprints) > 1 {
		sprint2Tasks := []*domain.Task{
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "实现 Sprint 管理功能",
				Description: "开发 Sprint 的创建、启动、完成等生命周期管理",
				Status:      domain.TaskStatusInProgress,
				Priority:    domain.PriorityHigh,
				Assignee:    func() *int64 { v := int64(102); return &v }(),
				SprintID:    &sprints[1].ID,
				Labels:      domain.Labels{"sprint", "management", "backend"},
				Estimate:    func() *int { v := 21; return &v }(),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
				Meta: domain.Meta{
					"difficulty":      "high",
					"category":        "backend",
					"priority_reason": "核心功能",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "设计任务看板界面",
				Description: "创建拖拽式的任务看板，支持状态变更",
				Status:      domain.TaskStatusTodo,
				Priority:    domain.PriorityMedium,
				Assignee:    func() *int64 { v := int64(101); return &v }(),
				SprintID:    &sprints[1].ID,
				Labels:      domain.Labels{"kanban", "ui", "frontend"},
				Estimate:    func() *int { v := 13; return &v }(),
				DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 10); return &t }(),
				Meta: domain.Meta{
					"difficulty": "medium",
					"category":   "frontend",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "开发燃尽图功能",
				Description: "实现 Sprint 进度的可视化燃尽图",
				Status:      domain.TaskStatusTodo,
				Priority:    domain.PriorityLow,
				Assignee:    func() *int64 { v := int64(104); return &v }(),
				SprintID:    &sprints[1].ID,
				Labels:      domain.Labels{"charts", "analytics", "frontend"},
				Estimate:    func() *int { v := 8; return &v }(),
				Meta: domain.Meta{
					"difficulty": "medium",
					"category":   "frontend",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "编写API文档",
				Description: "为所有API端点编写详细的文档",
				Status:      domain.TaskStatusTodo,
				Priority:    domain.PriorityLow,
				Labels:      domain.Labels{"documentation", "api"},
				Estimate:    func() *int { v := 5; return &v }(),
				Meta: domain.Meta{
					"difficulty": "low",
					"category":   "documentation",
				},
			},
		}

		for _, task := range sprint2Tasks {
			if err := tdb.Create(task).Error; err != nil {
				return fmt.Errorf("failed to create task %s: %w", task.Title, err)
			}
			logger.WithFields(logger.Fields{
				"task_id":    task.ID,
				"task_title": task.Title,
				"sprint_id":  task.SprintID,
				"status":     task.Status,
			}).Info("Task created")
		}
	}

	// 为第三个 Sprint（计划中）创建任务
	if len(sprints) > 2 {
		sprint3Tasks := []*domain.Task{
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "性能优化和缓存实现",
				Description: "优化数据库查询性能，实现Redis缓存",
				Status:      domain.TaskStatusTodo,
				Priority:    domain.PriorityMedium,
				SprintID:    &sprints[2].ID,
				Labels:      domain.Labels{"performance", "cache", "optimization"},
				Estimate:    func() *int { v := 13; return &v }(),
				Meta: domain.Meta{
					"difficulty": "high",
					"category":   "backend",
				},
			},
			{
				BaseModel:   db.BaseModel{TenantID: tenantID},
				Title:       "集成第三方通知服务",
				Description: "集成邮件和Slack通知功能",
				Status:      domain.TaskStatusTodo,
				Priority:    domain.PriorityLow,
				SprintID:    &sprints[2].ID,
				Labels:      domain.Labels{"integration", "notifications", "third-party"},
				Estimate:    func() *int { v := 8; return &v }(),
				Meta: domain.Meta{
					"difficulty": "medium",
					"category":   "integration",
				},
			},
		}

		for _, task := range sprint3Tasks {
			if err := tdb.Create(task).Error; err != nil {
				return fmt.Errorf("failed to create task %s: %w", task.Title, err)
			}
			logger.WithFields(logger.Fields{
				"task_id":    task.ID,
				"task_title": task.Title,
				"sprint_id":  task.SprintID,
				"status":     task.Status,
			}).Info("Task created")
		}
	}

	// 创建一些没有分配到 Sprint 的待办任务
	backlogTasks := []*domain.Task{
		{
			BaseModel:   db.BaseModel{TenantID: tenantID},
			Title:       "用户权限管理系统",
			Description: "实现基于角色的用户权限管理",
			Status:      domain.TaskStatusTodo,
			Priority:    domain.PriorityMedium,
			Labels:      domain.Labels{"auth", "rbac", "security"},
			Estimate:    func() *int { v := 21; return &v }(),
			Meta: domain.Meta{
				"difficulty": "high",
				"category":   "backend",
				"epic":       "用户管理",
			},
		},
		{
			BaseModel:   db.BaseModel{TenantID: tenantID},
			Title:       "移动端适配",
			Description: "优化界面以支持移动设备访问",
			Status:      domain.TaskStatusTodo,
			Priority:    domain.PriorityLow,
			Labels:      domain.Labels{"mobile", "responsive", "ui"},
			Estimate:    func() *int { v := 13; return &v }(),
			Meta: domain.Meta{
				"difficulty": "medium",
				"category":   "frontend",
			},
		},
		{
			BaseModel:   db.BaseModel{TenantID: tenantID},
			Title:       "自动化测试套件",
			Description: "编写单元测试和集成测试",
			Status:      domain.TaskStatusTodo,
			Priority:    domain.PriorityMedium,
			Labels:      domain.Labels{"testing", "automation", "quality"},
			Estimate:    func() *int { v := 34; return &v }(),
			Meta: domain.Meta{
				"difficulty": "high",
				"category":   "testing",
			},
		},
	}

	for _, task := range backlogTasks {
		if err := tdb.Create(task).Error; err != nil {
			return fmt.Errorf("failed to create backlog task %s: %w", task.Title, err)
		}
		logger.WithFields(logger.Fields{
			"task_id":    task.ID,
			"task_title": task.Title,
			"status":     task.Status,
		}).Info("Backlog task created")
	}

	return nil
}
