package iam

import (
	"scrum-plugin/internal/transport/http/admin/iam"

	"github.com/gin-gonic/gin"
)

// Routes IAM 路由配置
type Routes struct {
	memberHandler *iam.MemberHandler
}

// NewRoutes 创建 IAM 路由
func NewRoutes(memberHandler *iam.MemberHandler) *Routes {
	return &Routes{
		memberHandler: memberHandler,
	}
}

// Register 注册 IAM 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	iamGroup := rg.Group("/iam")
	{
		// 注册成员路由
		r.registerMemberRoutes(iamGroup)

		// 注册团队路由（暂时使用占位符）
		r.registerTeamRoutes(iamGroup)

		// IAM 整体状态检查
		iamGroup.GET("/status", r.getStatus)
	}
}

// registerMemberRoutes 注册成员管理路由
func (r *Routes) registerMemberRoutes(rg *gin.RouterGroup) {
	members := rg.Group("/members")
	{
		// 基础 CRUD 操作
		members.GET("", r.memberHandler.ListMembers)   // GET /admin/iam/members
		members.GET("/:id", r.memberHandler.GetMember) // GET /admin/iam/members/:id

		// 高级搜索
		members.GET("/search", r.memberHandler.SearchMembers) // GET /admin/iam/members/search

		// 连接检查
		members.GET("/connection/check", r.memberHandler.CheckMemberConnection) // GET /admin/iam/members/connection/check
	}
}

// registerTeamRoutes 注册团队管理路由（预留）
func (r *Routes) registerTeamRoutes(rg *gin.RouterGroup) {
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
func (r *Routes) getStatus(c *gin.Context) {
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
		"note": "集成 PowerX gRPC 服务进行成员和团队管理",
	})
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/admin"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Admin IAM API"
}
