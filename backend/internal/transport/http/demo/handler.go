package demohandler

import (
	"strconv"

	"scrum-plugin/internal/contracts"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/services"

	"github.com/gin-gonic/gin"
)

// Handler Demo 处理器
type Handler struct {
	demoService *services.DemoService
}

// NewHandler 创建 Demo 处理器
func NewHandler(demoService *services.DemoService) *Handler {
	return &Handler{
		demoService: demoService,
	}
}

// HealthCheck 健康检查
func (h *Handler) HealthCheck(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	req := &services.HealthCheckRequest{}
	resp, err := h.demoService.HealthCheck(c.Request.Context(), req)
	if err != nil {
		log.WithError(err).Error("PowerX gRPC health check failed")
		contracts.ResponseServiceUnavailable(c, "PowerX gRPC service unavailable", map[string]interface{}{
			"error":     err.Error(),
			"connected": false,
		})
		return
	}

	log.Info("PowerX gRPC health check successful")
	contracts.ResponseSuccess(c, resp)
}

// ListMembers 获取成员列表
func (h *Handler) ListMembers(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	req := &services.ListMembersRequest{
		Keyword:   c.Query("keyword"),
		PageIndex: 0,
		PageSize:  20,
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			req.PageIndex = int32(page)
		}
	}
	if sizeStr := c.Query("size"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 32); err == nil && size > 0 && size <= 100 {
			req.PageSize = int32(size)
		}
	}

	resp, err := h.demoService.ListMembers(c.Request.Context(), req)
	if err != nil {
		log.WithError(err).Error("Failed to list members")
		contracts.ResponseInternalError(c, err)
		return
	}

	log.WithField("keyword", req.Keyword).Info("Successfully listed members")
	contracts.ResponseSuccessWithMessage(c, resp, "Successfully retrieved members from PowerX gRPC service")
}

// ListTeams 获取团队列表
func (h *Handler) ListTeams(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	req := &services.ListTeamsRequest{
		Keyword:   c.Query("keyword"),
		PageIndex: 0,
		PageSize:  20,
	}

	// 解析分页参数
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.ParseInt(pageStr, 10, 32); err == nil {
			req.PageIndex = int32(page)
		}
	}
	if sizeStr := c.Query("size"); sizeStr != "" {
		if size, err := strconv.ParseInt(sizeStr, 10, 32); err == nil && size > 0 && size <= 100 {
			req.PageSize = int32(size)
		}
	}

	resp, err := h.demoService.ListTeams(c.Request.Context(), req)
	if err != nil {
		log.WithError(err).Error("Failed to list teams")
		contracts.ResponseInternalError(c, err)
		return
	}

	log.WithField("keyword", req.Keyword).Info("Successfully listed teams")
	contracts.ResponseSuccessWithMessage(c, resp, "Successfully retrieved teams from PowerX gRPC service")
}

// GetMember 获取单个成员
func (h *Handler) GetMember(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	idStr := c.Param("id")

	req := &services.GetMemberRequest{
		ID: 1, // 模拟 ID，实际应该解析 idStr
	}

	resp, err := h.demoService.GetMember(c.Request.Context(), req, idStr)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Failed to get member")
		contracts.ResponseInternalError(c, err)
		return
	}

	log.WithField("id", idStr).Info("Successfully retrieved member")
	contracts.ResponseSuccess(c, resp)
}

// GetTeam 获取单个团队
func (h *Handler) GetTeam(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	idStr := c.Param("id")

	req := &services.GetTeamRequest{
		ID: 1, // 模拟 ID，实际应该解析 idStr
	}

	resp, err := h.demoService.GetTeam(c.Request.Context(), req, idStr)
	if err != nil {
		log.WithError(err).WithField("id", idStr).Error("Failed to get team")
		contracts.ResponseInternalError(c, err)
		return
	}

	log.WithField("id", idStr).Info("Successfully retrieved team")
	contracts.ResponseSuccess(c, resp)
}

// Debug 获取调试信息
func (h *Handler) Debug(c *gin.Context) {
	log := logger.HandlerLogger("demo").WithContext(c.Request.Context())

	resp := h.demoService.GetDebugInfo(c.Request.Context())

	log.Info("Retrieved debug information")
	contracts.ResponseSuccessWithMessage(c, resp, "PowerX gRPC Demo API Debug Information")
}
