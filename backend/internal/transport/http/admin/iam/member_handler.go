package iam

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/contracts"
	iampowerx "github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client/iam"
	"github.com/ArtisanCloud/PowerXPlugin/internal/logger"
	iamsrv "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/iam"

	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"net/http"
	"strconv"
	"time"

	powerxclient "github.com/ArtisanCloud/PowerXPlugin/internal/grpc/client"

	"github.com/gin-gonic/gin"
)

// MemberHandler 成员管理处理器
type MemberHandler struct {
	*app.Deps
	memberService       *iamsrv.MemberService          // 内部业务服务
	memberServiceClient *iampowerx.MemberServiceClient // 内部业务服务
}

// NewMemberHandler 创建成员处理器
func NewMemberHandler(deps *app.Deps) *MemberHandler {
	return &MemberHandler{
		Deps:                deps,
		memberService:       iamsrv.NewMemberService(deps.DB),
		memberServiceClient: iampowerx.NewMemberServiceClient(deps.PowerXClient),
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
	req := &iampowerx.ListMembersRequest{
		Ctx: h.Deps.PowerXClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: int32(pageIndex),
			PageSize:  int32(pageSize),
		},
		Keyword: keyword,
	}

	// 调用 PowerX gRPC 服务
    resp, err := h.memberServiceClient.ListMembers(c.Request.Context(), req)
    if err != nil {
        logger.WithError(err).Error("Failed to call PowerX ListMembers gRPC service")
        contracts.ResponseErrorWithDetails(c, http.StatusInternalServerError, "POWERX_GRPC_ERROR", "Failed to fetch members from PowerX", err.Error())
        return
    }

	logger.WithField("count", len(resp.Members)).Info("Successfully fetched members from PowerX")

	// 构造分页信息
	totalPages := 0
	if resp.PageSize > 0 {
		totalPages = int((resp.TotalCount + int64(resp.PageSize) - 1) / int64(resp.PageSize))
	}

    contracts.ResponseSuccessWithMessage(c, &contracts.ListResponse{
        Data: resp.Members,
        Pagination: &contracts.PaginationResponse{
            Page:       int(resp.PageIndex) + 1,
            Limit:      int(resp.PageSize),
            Total:      resp.TotalCount,
            TotalPages: totalPages,
        },
    }, "Successfully retrieved members list")
}

// GetMember 获取单个成员（通过 PowerX gRPC）
func (h *MemberHandler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        contracts.ResponseErrorWithDetails(c, http.StatusBadRequest, contracts.ErrCodeInvalidRequest, "Invalid member ID", err.Error())
        return
    }

	// 构造 PowerX gRPC 请求
	req := &iampowerx.GetMemberRequest{
		Ctx: h.Deps.PowerXClient.RC(),
		Id:  id,
	}

	// 调用 PowerX gRPC 服务
    resp, err := h.memberServiceClient.GetMember(c.Request.Context(), req)
    if err != nil {
        logger.WithError(err).WithField("member_id", id).Error("Failed to call PowerX GetMember gRPC service")
        contracts.ResponseErrorWithDetails(c, http.StatusInternalServerError, "POWERX_GRPC_ERROR", "Failed to fetch member from PowerX", err.Error())
        return
    }

	logger.WithField("member_id", id).Info("Successfully fetched member from PowerX")

    contracts.ResponseSuccessWithMessage(c, resp.Member, "Successfully retrieved member details")
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
	req := &iampowerx.ListMembersRequest{
		Ctx: h.Deps.PowerXClient.RC(),
		Page: &powerxclient.PageRequest{
			PageIndex: 0,
			PageSize:  50, // 搜索时使用更大的页面大小
		},
		Keyword:  keyword,
		TeamIds:  teamIds,
		Statuses: statuses,
	}

	// 调用 PowerX gRPC 服务
    resp, err := h.memberServiceClient.ListMembers(c.Request.Context(), req)
    if err != nil {
        logger.WithError(err).Error("Failed to search members via PowerX gRPC")
        contracts.ResponseErrorWithDetails(c, http.StatusInternalServerError, "POWERX_GRPC_ERROR", "Failed to search members", err.Error())
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

    contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
        "members": resp.Members,
        "query":   searchQuery,
        "pagination": &contracts.PaginationResponse{
            Page:       int(resp.PageIndex) + 1,
            Limit:      int(resp.PageSize),
            Total:      resp.TotalCount,
            TotalPages: totalPages,
        },
    }, "Successfully searched members")
}

// CheckMemberConnection 检查 PowerX 连接状态
func (h *MemberHandler) CheckMemberConnection(c *gin.Context) {
	// 健康检查
    if err := h.Deps.PowerXClient.HealthCheck(c.Request.Context()); err != nil {
        contracts.ResponseServiceUnavailable(c, "PowerX gRPC service unavailable", err.Error())
        return
    }

    contracts.ResponseSuccessWithMessage(c, &contracts.HealthResponse{
        Status:    "healthy",
        Service:   "PowerX MemberService",
        Version:   "1.0.0",
        Timestamp: time.Now(),
        Checks: map[string]string{
            "grpc_connected": func() string {
                if h.Deps.PowerXClient.IsConnected() {
                    return "ok"
                }
                return "disconnected"
            }(),
            "tenant_id": func() string {
                if h.Deps.PowerXClient.GetTenantID() > 0 {
                    return "configured"
                }
                return "not_configured"
            }(),
            "has_token": func() string {
                if h.Deps.PowerXClient.HasToken() {
                    return "configured"
                }
                return "not_configured"
            }(),
        },
    }, "PowerX MemberService connection is healthy")
}
