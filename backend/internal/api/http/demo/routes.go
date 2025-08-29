package demo

import (
	"scrum-plugin/internal/contracts"
	powerxclient "scrum-plugin/internal/grpc/client"

	"github.com/gin-gonic/gin"
)

// GRPCDemoRoutes 设置 gRPC 演示路由
func GRPCDemoRoutes(r *gin.Engine, pxc *powerxclient.PowerX) {
	demo := r.Group("/api/v1/demo/grpc")
	{
		// 检查 PowerX gRPC 连接状态
		demo.GET("/health", func(c *gin.Context) {
			if err := pxc.HealthCheck(c.Request.Context()); err != nil {
				contracts.ResponseServiceUnavailable(c, "PowerX gRPC service unavailable", map[string]interface{}{
					"error":     err.Error(),
					"connected": false,
				})
				return
			}

			contracts.ResponseSuccess(c, map[string]interface{}{
				"status":    "ok",
				"connected": pxc.IsConnected(),
				"tenant_id": pxc.GetTenantID(),
				"has_token": pxc.GetToken() != "",
			})
		})

		// 调用 PowerX gRPC 服务获取成员列表
		demo.GET("/members", func(c *gin.Context) {
			// 构造请求
			req := &powerxclient.ListMembersRequest{
				Ctx: pxc.RC(),
				Page: &powerxclient.PageRequest{
					PageIndex: 0,
					PageSize:  20,
				},
			}

			// 获取查询参数
			if keyword := c.Query("keyword"); keyword != "" {
				req.Keyword = keyword
			}

			// 调用 gRPC 服务
			resp, err := pxc.ListMembers(c.Request.Context(), req)
			if err != nil {
				contracts.ResponseInternalError(c, err)
				return
			}

			contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
				"members": resp,
				"grpc": map[string]interface{}{
					"connected": pxc.IsConnected(),
					"tenant_id": pxc.GetTenantID(),
				},
			}, "Successfully retrieved members from PowerX gRPC service")
		})

		// 调用 PowerX gRPC 服务获取团队信息
		demo.GET("/teams", func(c *gin.Context) {
			// 构造请求
			req := &powerxclient.ListTeamsRequest{
				Ctx: pxc.RC(),
				Page: &powerxclient.PageRequest{
					PageIndex: 0,
					PageSize:  20,
				},
			}

			// 获取查询参数
			if keyword := c.Query("keyword"); keyword != "" {
				req.Keyword = keyword
			}

			// 调用 gRPC 服务
			resp, err := pxc.ListTeams(c.Request.Context(), req)
			if err != nil {
				contracts.ResponseInternalError(c, err)
				return
			}

			contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
				"teams": resp,
				"grpc": map[string]interface{}{
					"connected": pxc.IsConnected(),
					"tenant_id": pxc.GetTenantID(),
				},
			}, "Successfully retrieved teams from PowerX gRPC service")
		})

		// 获取单个成员信息
		demo.GET("/members/:id", func(c *gin.Context) {
			id := c.Param("id")
			req := &powerxclient.GetMemberRequest{
				Ctx: pxc.RC(),
				Id:  1, // 模拟 ID
			}
			resp, err := pxc.GetMember(c.Request.Context(), req)
			if err != nil {
				contracts.ResponseInternalError(c, err)
				return
			}
			contracts.ResponseSuccess(c, map[string]interface{}{
				"member":       resp,
				"requested_id": id,
			})
		})

		// 获取单个团队信息
		demo.GET("/teams/:id", func(c *gin.Context) {
			id := c.Param("id")
			req := &powerxclient.GetTeamRequest{
				Ctx: pxc.RC(),
				Id:  1, // 模拟 ID
			}
			resp, err := pxc.GetTeam(c.Request.Context(), req)
			if err != nil {
				contracts.ResponseInternalError(c, err)
				return
			}
			contracts.ResponseSuccess(c, map[string]interface{}{
				"team":         resp,
				"requested_id": id,
			})
		})

		// 调试信息接口
		demo.GET("/debug", func(c *gin.Context) {
			contracts.ResponseSuccessWithMessage(c, map[string]interface{}{
				"grpc_connection": map[string]interface{}{
					"connected": pxc.IsConnected(),
					"tenant_id": pxc.GetTenantID(),
					"has_token": pxc.GetToken() != "",
				},
				"endpoints": []string{
					"GET /api/v1/demo/grpc/health - 检查 gRPC 连接状态",
					"GET /api/v1/demo/grpc/members - 获取成员列表",
					"GET /api/v1/demo/grpc/members/{id} - 获取单个成员",
					"GET /api/v1/demo/grpc/teams - 获取团队列表",
					"GET /api/v1/demo/grpc/teams/{id} - 获取单个团队",
					"GET /api/v1/demo/grpc/debug - 查看调试信息",
				},
				"note": "当前使用模拟数据，可以通过这些接口测试 gRPC 连接和功能",
			}, "PowerX gRPC Demo API Debug Information")
		})
	}
}
