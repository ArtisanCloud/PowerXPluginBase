package iam

import (
	"scrum-plugin/internal/transport/http/admin/iam"

	"github.com/gin-gonic/gin"
)

// RegisterMemberRoutes 注册成员管理路由
func RegisterMemberRoutes(r *gin.RouterGroup, handler *iam.MemberHandler) {
	members := r.Group("/members")
	{
		// 基础 CRUD 操作
		members.GET("", handler.ListMembers)   // GET /api/v1/admin/iam/members
		members.GET("/:id", handler.GetMember) // GET /api/v1/admin/iam/members/:id

		// 高级搜索
		members.GET("/search", handler.SearchMembers) // GET /api/v1/admin/iam/members/search

		// 连接检查
		members.GET("/connection/check", handler.CheckMemberConnection) // GET /api/v1/admin/iam/members/connection/check
	}
}

// RegisterTeamRoutes 注册团队管理路由（预留）
func RegisterTeamRoutes(r *gin.RouterGroup, handler interface{}) {
	teams := r.Group("/teams")
	{
		// 预留团队相关路由
		teams.GET("", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Team routes will be implemented here",
				"status":  "placeholder",
			})
		})
	}
}

// RegisterIAMRoutes 注册所有 IAM 相关路由
func RegisterIAMRoutes(r *gin.RouterGroup, memberHandler *iam.MemberHandler) {
	iamGroup := r.Group("/iam")
	{
		// 注册成员路由
		RegisterMemberRoutes(iamGroup, memberHandler)

		// 注册团队路由（暂时使用占位符）
		RegisterTeamRoutes(iamGroup, nil)

		// IAM 整体状态检查
		iamGroup.GET("/status", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"module":   "iam",
				"version":  "1.0.0",
				"services": []string{"members", "teams"},
				"endpoints": []string{
					"GET /iam/members - 获取成员列表",
					"GET /iam/members/:id - 获取单个成员",
					"GET /iam/members/search - 搜索成员",
					"GET /iam/members/connection/check - 检查连接",
					"GET /iam/status - 模块状态",
				},
				"note": "集成 PowerX gRPC 服务进行成员和团队管理",
			})
		})
	}
}
