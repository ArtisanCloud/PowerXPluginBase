package iam

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// Register 注册 IAM 路由
func RegisterAPIRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	iamGroup := rg.Group("/iam")
	{
		// 注册成员路由
		registerMemberRoutes(iamGroup, deps)

		// 注册团队路由（暂时使用占位符）
		registerTeamRoutes(iamGroup, deps)

		// IAM 整体状态检查
		iamGroup.GET("/status", getStatus)
	}
}

// registerMemberRoutes 注册成员管理路由
func registerMemberRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	memberHandler := NewMemberHandler(deps)
	members := rg.Group("/members")
	{
		// 基础 CRUD 操作
		members.GET("", memberHandler.ListMembers)   // GET /admin/iam/members
		members.GET("/:id", memberHandler.GetMember) // GET /admin/iam/members/:id

		// 高级搜索
		members.GET("/search", memberHandler.SearchMembers) // GET /admin/iam/members/search

		// 连接检查
		members.GET("/connection/check", memberHandler.CheckMemberConnection) // GET /admin/iam/members/connection/check
	}
}

// registerTeamRoutes 注册团队管理路由（预留）
func registerTeamRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	teams := rg.Group("/teams")
	{
		// 预留团队相关路由
		teams.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Team routes will be implemented here",
				"status":  "placeholder",
				"module":  "admin.iam.teams",
				"endpoints": []string{
					"GET /admin/iam/teams - 获取团队列表",
					"GET /admin/iam/teams/:id - 获取团队详情",
				},
			})
		})
	}
}

// getStatus IAM 模块状态检查
func getStatus(c *gin.Context) {
	c.JSON(200, gin.H{
		"module":   "admin.iam",
		"version":  "1.0.0",
		"services": []string{"members", "teams"},
		"endpoints": []string{
			"GET /admin/iam/members - 获取成员列表",
			"GET /admin/iam/members/:id - 获取单个成员",
			"GET /admin/iam/members/search - 搜索成员",
			"GET /admin/iam/members/connection/check - 检查 PowerX 连接",
			"GET /admin/iam/teams - 获取团队列表（预留）",
			"GET /admin/iam/status - 模块状态",
		},
		"tip": "集成 PowerX gRPC 服务进行成员和团队管理，支撑模板协作场景",
	})
}
