package task

import (
	"strconv"

	"scrum-plugin/internal/contracts"
	"scrum-plugin/internal/domain/models"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/middleware"
	"scrum-plugin/internal/services"

	"github.com/gin-gonic/gin"
)

// Handler Task 后台管理处理器
type Handler struct {
	taskService *services.TaskService
}

// NewHandler 创建 Task 后台管理处理器
func NewHandler(taskService *services.TaskService) *Handler {
	return &Handler{
		taskService: taskService,
	}
}

// ListTasks 获取任务列表（后台管理版本）
func (h *Handler) ListTasks(c *gin.Context) {
	log := logger.HandlerLogger("admin.task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		contracts.ResponseUnauthorized(c, err.Error())
		return
	}

	// 解析分页参数
	var req contracts.PaginationRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.WithError(err).Error("Failed to bind pagination request")
		contracts.ResponseBadRequest(c, "Invalid pagination parameters")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// 调用 task service
	tasks, total, err := h.taskService.ListTasks(c.Request.Context(), tenantID, &services.TaskListOptions{
		Page:     req.Page,
		Limit:    req.Limit,
		Status:   stringPtr(c.Query("status")),
		Priority: stringPtr(c.Query("priority")),
		// Assignee: c.Query("assignee"),  // 需要转换为 *int64
		// SprintID: c.Query("sprint_id"), // 需要转换为 *uint
	})
	if err != nil {
		log.WithError(err).Error("Failed to list tasks")
		contracts.ResponseInternalError(c, err)
		return
	}

	// 转换响应
	var responses []*contracts.TaskResponse
	for _, task := range tasks {
		responses = append(responses, h.taskToResponse(task))
	}

	// 计算分页信息
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	listResponse := contracts.ListResponse{
		Data: responses,
		Pagination: &contracts.PaginationResponse{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	log.WithField("count", len(responses)).Info("Admin tasks listed successfully")
	contracts.ResponseSuccess(c, listResponse)
}

// GetTask 获取单个任务（后台管理版本）
func (h *Handler) GetTask(c *gin.Context) {
	log := logger.HandlerLogger("admin.task").WithContext(c.Request.Context())

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // 改为 64 位
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		contracts.ResponseBadRequest(c, "Invalid task ID")
		return
	}

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		contracts.ResponseUnauthorized(c, err.Error())
		return
	}

	// 获取任务
	task, err := h.taskService.GetTask(c.Request.Context(), tenantID, id) // 使用 uint64
	if err != nil {
		if err.Error() == "task not found" {
			log.WithField("task_id", id).Warn("Task not found")
			contracts.ResponseNotFound(c, "Task not found")
			return
		}
		log.WithError(err).WithField("task_id", id).Error("Failed to get task")
		contracts.ResponseInternalError(c, err)
		return
	}

	// 转换响应
	response := h.taskToResponse(task)

	log.WithField("task_id", id).Info("Admin task retrieved successfully")
	contracts.ResponseSuccess(c, response)
}

// UpdateTaskStatus 更新任务状态（后台管理版本）
func (h *Handler) UpdateTaskStatus(c *gin.Context) {
	log := logger.HandlerLogger("admin.task").WithContext(c.Request.Context())

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64) // 改为 64 位
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		contracts.ResponseBadRequest(c, "Invalid task ID")
		return
	}

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		contracts.ResponseUnauthorized(c, err.Error())
		return
	}

	// 解析请求
	var req struct {
		Status string `json:"status" binding:"required,oneof=todo in_progress done"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Failed to bind request")
		contracts.ResponseBadRequest(c, "Invalid request parameters")
		return
	}

	// 更新状态
	err = h.taskService.UpdateTaskStatus(c.Request.Context(), tenantID, uint(id), req.Status) // 转换为 uint
	if err != nil {
		if err.Error() == "task not found" {
			log.WithField("task_id", id).Warn("Task not found")
			contracts.ResponseNotFound(c, "Task not found")
			return
		}
		log.WithError(err).WithField("task_id", id).Error("Failed to update task status")
		contracts.ResponseInternalError(c, err)
		return
	}

	log.WithField("task_id", id).WithField("status", req.Status).Info("Admin task status updated successfully")
	contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
		"task_id": id,
		"status":  req.Status,
	}, "Task status updated successfully")
}

// GetTaskStats 获取任务统计信息
func (h *Handler) GetTaskStats(c *gin.Context) {
	log := logger.HandlerLogger("admin.task").WithContext(c.Request.Context())

	// TODO: 当实现统计逻辑时需要获取租户 ID
	// tenantID, err := middleware.GetTenantID(c)
	// if err != nil {
	// 	log.WithError(err).Error("Failed to get tenant ID")
	// 	contracts.ResponseUnauthorized(c, err.Error())
	// 	return
	// }

	// TODO: 实现统计逻辑
	// stats, err := h.taskService.GetTaskStats(c.Request.Context(), tenantID)
	// if err != nil {
	//     log.WithError(err).Error("Failed to get task statistics")
	//     contracts.ResponseInternalError(c, err)
	//     return
	// }

	// 临时返回占位符统计
	log.Info("Admin task statistics requested (placeholder implementation)")
	contracts.ResponseSuccess(c, map[string]interface{}{
		"total_tasks":       0,
		"todo_tasks":        0,
		"in_progress_tasks": 0,
		"done_tasks":        0,
		"by_priority": map[string]int{
			"low":    0,
			"medium": 0,
			"high":   0,
			"urgent": 0,
		},
		"message": "Task statistics will be implemented here",
	})
}

// BatchUpdateStatus 批量更新任务状态
func (h *Handler) BatchUpdateStatus(c *gin.Context) {
	log := logger.HandlerLogger("admin.task").WithContext(c.Request.Context())

	// 解析请求
	var req struct {
		TaskIDs []uint64 `json:"task_ids" binding:"required,min=1"` // 改为 uint64
		Status  string   `json:"status" binding:"required,oneof=todo in_progress done"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Failed to bind batch update request")
		contracts.ResponseBadRequest(c, "Invalid request parameters")
		return
	}

	// 获取租户 ID
	// tenantID, err := middleware.GetTenantID(c)
	// if err != nil {
	// 	log.WithError(err).Error("Failed to get tenant ID")
	// 	contracts.ResponseUnauthorized(c, err.Error())
	// 	return
	// }

	// 转换 TaskIDs 为 uint 类型
	taskIDs := make([]uint, len(req.TaskIDs))
	for i, id := range req.TaskIDs {
		taskIDs[i] = uint(id)
	}

	// TODO: 实现批量更新逻辑
	// err = h.taskService.BatchUpdateStatus(c.Request.Context(), tenantID, taskIDs, req.Status)
	// if err != nil {
	//     log.WithError(err).Error("Failed to batch update task status")
	//     contracts.ResponseInternalError(c, err)
	//     return
	// }

	log.WithField("task_count", len(req.TaskIDs)).WithField("status", req.Status).Info("Admin batch task status update requested (placeholder implementation)")
	contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
		"updated_count": len(req.TaskIDs),
		"status":        req.Status,
		"message":       "Batch update will be implemented here",
	}, "Batch update requested successfully")
}

// taskToResponse 将领域模型转换为响应（复用现有的转换逻辑）
func (h *Handler) taskToResponse(task *models.Task) *contracts.TaskResponse {
	// TODO: 实现转换逻辑，可以参考现有的 task_handler.go 中的实现
	return &contracts.TaskResponse{
		// 临时占位符
		ID: task.ID,
	}
}

// stringPtr 将字符串转换为字符串指针
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
