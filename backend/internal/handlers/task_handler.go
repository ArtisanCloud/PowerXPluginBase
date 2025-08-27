package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/powerx-plugins/scrum/internal/contracts"
	"github.com/powerx-plugins/scrum/internal/domain"
	"github.com/powerx-plugins/scrum/internal/logger"
	"github.com/powerx-plugins/scrum/internal/middleware"
)

// TaskHandler 任务处理器
type TaskHandler struct {
	taskService domain.TaskService
}

// NewTaskHandler 创建任务处理器
func NewTaskHandler(taskService domain.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// CreateTask 创建任务
func (h *TaskHandler) CreateTask(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 解析请求
	var req contracts.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeValidationFailed, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 转换为领域请求
	domainReq := &domain.CreateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Status:      domain.TaskStatus(req.Status),
		Priority:    domain.Priority(req.Priority),
		Assignee:    req.Assignee,
		SprintID:    req.SprintID,
		Labels:      req.Labels,
		DueDate:     req.DueDate,
		Estimate:    req.Estimate,
		Meta:        req.Meta,
	}

	// 创建任务
	task, err := h.taskService.CreateTask(c.Request.Context(), tenantID, domainReq)
	if err != nil {
		log.WithError(err).Error("Failed to create task")
		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 转换响应
	response := h.taskToResponse(task)

	log.WithField("task_id", task.ID).Info("Task created successfully")
	c.JSON(http.StatusCreated, contracts.APIResponse{
		Success:   true,
		Data:      response,
		Timestamp: time.Now(),
	})
}

// GetTask 获取任务详情
func (h *TaskHandler) GetTask(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInvalidRequest, Message: "Invalid task ID"},
			Timestamp: time.Now(),
		})
		return
	}

	// 获取任务
	task, err := h.taskService.GetTask(c.Request.Context(), tenantID, uint(id))
	if err != nil {
		log.WithError(err).WithField("task_id", id).Error("Failed to get task")

		// 检查是否为未找到错误
		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, contracts.APIResponse{
				Success:   false,
				Error:     &contracts.APIError{Code: contracts.ErrCodeTaskNotFound, Message: "Task not found"},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 转换响应
	response := h.taskToResponse(task)

	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Data:      response,
		Timestamp: time.Now(),
	})
}

// UpdateTask 更新任务
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInvalidRequest, Message: "Invalid task ID"},
			Timestamp: time.Now(),
		})
		return
	}

	// 解析请求
	var req contracts.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeValidationFailed, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 转换为领域请求
	domainReq := &domain.UpdateTaskRequest{
		Title:       req.Title,
		Description: req.Description,
		Assignee:    req.Assignee,
		SprintID:    req.SprintID,
		Labels:      req.Labels,
		DueDate:     req.DueDate,
		Estimate:    req.Estimate,
		Meta:        req.Meta,
	}

	if req.Status != nil {
		status := domain.TaskStatus(*req.Status)
		domainReq.Status = &status
	}

	if req.Priority != nil {
		priority := domain.Priority(*req.Priority)
		domainReq.Priority = &priority
	}

	// 更新任务
	task, err := h.taskService.UpdateTask(c.Request.Context(), tenantID, uint(id), domainReq)
	if err != nil {
		log.WithError(err).WithField("task_id", id).Error("Failed to update task")

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, contracts.APIResponse{
				Success:   false,
				Error:     &contracts.APIError{Code: contracts.ErrCodeTaskNotFound, Message: "Task not found"},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 转换响应
	response := h.taskToResponse(task)

	log.WithField("task_id", task.ID).Info("Task updated successfully")
	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Data:      response,
		Timestamp: time.Now(),
	})
}

// DeleteTask 删除任务
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInvalidRequest, Message: "Invalid task ID"},
			Timestamp: time.Now(),
		})
		return
	}

	// 删除任务
	err = h.taskService.DeleteTask(c.Request.Context(), tenantID, uint(id))
	if err != nil {
		log.WithError(err).WithField("task_id", id).Error("Failed to delete task")

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, contracts.APIResponse{
				Success:   false,
				Error:     &contracts.APIError{Code: contracts.ErrCodeTaskNotFound, Message: "Task not found"},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	log.WithField("task_id", id).Info("Task deleted successfully")
	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Message:   "Task deleted successfully",
		Timestamp: time.Now(),
	})
}

// ListTasks 获取任务列表
func (h *TaskHandler) ListTasks(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 解析查询参数
	var req contracts.TaskListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.WithError(err).Error("Failed to bind query parameters")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeValidationFailed, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	// 转换为领域查询选项
	opts := &domain.TaskListOptions{
		Page:      req.Page,
		Limit:     req.Limit,
		Labels:    req.Labels,
		Search:    req.Search,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	if req.Status != "" {
		status := domain.TaskStatus(req.Status)
		opts.Status = &status
	}

	if req.Priority != "" {
		priority := domain.Priority(req.Priority)
		opts.Priority = &priority
	}

	if req.Assignee != nil {
		opts.Assignee = req.Assignee
	}

	if req.SprintID != nil {
		opts.SprintID = req.SprintID
	}

	// 获取任务列表
	tasks, total, err := h.taskService.ListTasks(c.Request.Context(), tenantID, opts)
	if err != nil {
		log.WithError(err).Error("Failed to list tasks")
		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
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

	log.WithField("count", len(responses)).Info("Tasks listed successfully")
	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Data:      listResponse,
		Timestamp: time.Now(),
	})
}

// UpdateTaskStatus 更新任务状态
func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	log := logger.HandlerLogger("task").WithContext(c.Request.Context())

	// 获取租户 ID
	tenantID, err := middleware.GetTenantID(c)
	if err != nil {
		log.WithError(err).Error("Failed to get tenant ID")
		c.JSON(http.StatusUnauthorized, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeUnauthorized, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 获取任务 ID
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Invalid task ID")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInvalidRequest, Message: "Invalid task ID"},
			Timestamp: time.Now(),
		})
		return
	}

	// 解析请求
	var req struct {
		Status string `json:"status" binding:"required,oneof=todo in_progress done"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		log.WithError(err).Error("Failed to bind request")
		c.JSON(http.StatusBadRequest, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeValidationFailed, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	// 更新状态
	err = h.taskService.UpdateTaskStatus(c.Request.Context(), tenantID, uint(id), domain.TaskStatus(req.Status))
	if err != nil {
		log.WithError(err).WithField("task_id", id).Error("Failed to update task status")

		if err.Error() == "task not found" {
			c.JSON(http.StatusNotFound, contracts.APIResponse{
				Success:   false,
				Error:     &contracts.APIError{Code: contracts.ErrCodeTaskNotFound, Message: "Task not found"},
				Timestamp: time.Now(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, contracts.APIResponse{
			Success:   false,
			Error:     &contracts.APIError{Code: contracts.ErrCodeInternalError, Message: err.Error()},
			Timestamp: time.Now(),
		})
		return
	}

	log.WithField("task_id", id).WithField("status", req.Status).Info("Task status updated successfully")
	c.JSON(http.StatusOK, contracts.APIResponse{
		Success:   true,
		Message:   "Task status updated successfully",
		Timestamp: time.Now(),
	})
}

// taskToResponse 将领域模型转换为响应
func (h *TaskHandler) taskToResponse(task *domain.Task) *contracts.TaskResponse {
	return &contracts.TaskResponse{
		ID:          task.ID,
		TenantID:    task.TenantID,
		Title:       task.Title,
		Description: task.Description,
		Status:      string(task.Status),
		Priority:    string(task.Priority),
		Assignee:    task.Assignee,
		SprintID:    task.SprintID,
		Labels:      []string(task.Labels),
		DueDate:     task.DueDate,
		Estimate:    task.Estimate,
		Meta:        map[string]interface{}(task.Meta),
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
