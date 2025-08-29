package iam

import (
	"net/http"
	"strconv"
	"time"

	"scrum-plugin/internal/contracts"
	powerxclient "scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/logger"
	"scrum-plugin/internal/services"

	"github.com/gin-gonic/gin"
)

// MemberHandler 成员管理处理器
type MemberHandler struct {
	memberService *services.MemberService // 内部业务服务
	powerxClient  *powerxclient.PowerX    // PowerX gRPC 客户端
}

// NewMemberHandler 创建成员处理器
func NewMemberHandler(memberService *services.MemberService, powerxClient *powerxclient.PowerX) *MemberHandler {
	return &MemberHandler{
		memberService: memberService,
		powerxClient:  powerxClient,
	}
}

// ListMembers 获取成员列表（通过 PowerX gRPC）
func (h *MemberHandler) ListMembers(c *gin.Context) {
	// 解析查询参数
	keyword := c.Query("keyword")
	pageSizeStr := c.DefaultQuery("page_size", "20")
	pageIndexStr := c.DefaultQuery("page_index", "0")

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		pageSize = 20
	}

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil || pageIndex < 0 {
		pageIndex = 0
	}

	// 构造 PowerX gRPC 请求
	req := &powerxclient.ListMembersRequest{
		Ctx: h.powerxClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: int32(pageIndex),
			PageSize:  int32(pageSize),
		},
		Keyword: keyword,
	}

	// 调用 PowerX gRPC 服务
	resp, err := h.powerxClient.ListMembers(c.Request.Context(), req)
	if err != nil {
		logger.WithError(err).Error("Failed to call PowerX ListMembers gRPC service")
		c.JSON(http.StatusInternalServerError, &contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:    "POWERX_GRPC_ERROR",
				Message: "Failed to fetch members from PowerX",
				Details: err.Error(),
			},
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	logger.WithField("count", len(resp.Members)).Info("Successfully fetched members from PowerX")

	// 构造分页信息
	totalPages := 0
	if resp.PageSize > 0 {
		totalPages = int((resp.TotalCount + int64(resp.PageSize) - 1) / int64(resp.PageSize))
	}

	c.JSON(http.StatusOK, &contracts.APIResponse{
		Success: true,
		Data: &contracts.ListResponse{
			Data: resp.Members,
			Pagination: &contracts.PaginationResponse{
				Page:       int(resp.PageIndex) + 1, // 转换为 1 基础的页码
				Limit:      int(resp.PageSize),
				Total:      resp.TotalCount,
				TotalPages: totalPages,
			},
		},
		Message:   "Successfully retrieved members list",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	})
}

// GetMember 获取单个成员（通过 PowerX gRPC）
func (h *MemberHandler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, &contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:    contracts.ErrCodeInvalidRequest,
				Message: "Invalid member ID",
				Details: err.Error(),
			},
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// 构造 PowerX gRPC 请求
	req := &powerxclient.GetMemberRequest{
		Ctx: h.powerxClient.RC(),
		Id:  id,
	}

	// 调用 PowerX gRPC 服务
	resp, err := h.powerxClient.GetMember(c.Request.Context(), req)
	if err != nil {
		logger.WithError(err).WithField("member_id", id).Error("Failed to call PowerX GetMember gRPC service")
		c.JSON(http.StatusInternalServerError, &contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:    "POWERX_GRPC_ERROR",
				Message: "Failed to fetch member from PowerX",
				Details: err.Error(),
			},
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	logger.WithField("member_id", id).Info("Successfully fetched member from PowerX")

	c.JSON(http.StatusOK, &contracts.APIResponse{
		Success:   true,
		Data:      resp.Member,
		Message:   "Successfully retrieved member details",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	})
}

// SearchMembers 搜索成员（带高级参数）
func (h *MemberHandler) SearchMembers(c *gin.Context) {
	// 解析复杂查询参数
	keyword := c.Query("keyword")
	department := c.Query("department")
	status := c.Query("status")

	// 解析团队ID列表
	var teamIds []int64
	if teamIdsStr := c.Query("team_ids"); teamIdsStr != "" {
		// 简化版本，实际可以解析逗号分隔的ID列表
		if teamId, err := strconv.ParseInt(teamIdsStr, 10, 64); err == nil {
			teamIds = append(teamIds, teamId)
		}
	}

	// 解析状态列表
	var statuses []string
	if status != "" {
		statuses = append(statuses, status)
	}

	// 构造 PowerX gRPC 请求
	req := &powerxclient.ListMembersRequest{
		Ctx: h.powerxClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: 0,
			PageSize:  50, // 搜索时使用更大的页面大小
		},
		Keyword:  keyword,
		TeamIds:  teamIds,
		Statuses: statuses,
	}

	// 调用 PowerX gRPC 服务
	resp, err := h.powerxClient.ListMembers(c.Request.Context(), req)
	if err != nil {
		logger.WithError(err).Error("Failed to search members via PowerX gRPC")
		c.JSON(http.StatusInternalServerError, &contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:    "POWERX_GRPC_ERROR",
				Message: "Failed to search members",
				Details: err.Error(),
			},
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	// 可以在这里添加额外的业务逻辑处理
	// 比如与内部服务结合，获取成员的任务统计等

	logger.WithField("keyword", keyword).
		WithField("department", department).
		WithField("count", len(resp.Members)).
		Info("Successfully searched members")

	// 构造搜索结果和分页信息
	totalPages := 0
	if resp.PageSize > 0 {
		totalPages = int((resp.TotalCount + int64(resp.PageSize) - 1) / int64(resp.PageSize))
	}

	// 构造搜索查询信息
	searchQuery := map[string]interface{}{
		"keyword":    keyword,
		"department": department,
		"status":     status,
		"team_ids":   teamIds,
	}

	c.JSON(http.StatusOK, &contracts.APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"members": resp.Members,
			"query":   searchQuery,
			"pagination": &contracts.PaginationResponse{
				Page:       int(resp.PageIndex) + 1,
				Limit:      int(resp.PageSize),
				Total:      resp.TotalCount,
				TotalPages: totalPages,
			},
		},
		Message:   "Successfully searched members",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	})
}

// CheckMemberConnection 检查 PowerX 连接状态
func (h *MemberHandler) CheckMemberConnection(c *gin.Context) {
	// 健康检查
	if err := h.powerxClient.HealthCheck(c.Request.Context()); err != nil {
		c.JSON(http.StatusServiceUnavailable, &contracts.APIResponse{
			Success: false,
			Error: &contracts.APIError{
				Code:    "POWERX_CONNECTION_ERROR",
				Message: "PowerX gRPC service unavailable",
				Details: err.Error(),
			},
			Timestamp: time.Now(),
			RequestID: c.GetString("request_id"),
		})
		return
	}

	c.JSON(http.StatusOK, &contracts.APIResponse{
		Success: true,
		Data: &contracts.HealthResponse{
			Status:    "healthy",
			Service:   "PowerX MemberService",
			Version:   "1.0.0",
			Timestamp: time.Now(),
			Checks: map[string]string{
				"grpc_connected": func() string {
					if h.powerxClient.IsConnected() {
						return "ok"
					}
					return "disconnected"
				}(),
				"tenant_id": func() string {
					if h.powerxClient.GetTenantID() > 0 {
						return "configured"
					}
					return "not_configured"
				}(),
				"has_token": func() string {
					if h.powerxClient.GetToken() != "" {
						return "configured"
					}
					return "not_configured"
				}(),
			},
		},
		Message:   "PowerX MemberService connection is healthy",
		Timestamp: time.Now(),
		RequestID: c.GetString("request_id"),
	})
}
