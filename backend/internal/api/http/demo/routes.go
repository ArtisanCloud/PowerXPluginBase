package demo

import (
	powerxclient "scrum-plugin/internal/grpc/client"
	"scrum-plugin/internal/services"
	demohandler "scrum-plugin/internal/transport/http/demo"

	"github.com/gin-gonic/gin"
)

// Routes Demo 路由配置
type Routes struct {
	handler *demohandler.Handler
}

// NewRoutes 创建 Demo 路由
func NewRoutes(powerxClient *powerxclient.PowerX) *Routes {
	// 创建服务层
	demoService := services.NewDemoService(powerxClient)
	
	// 创建 handler 层
	handler := demohandler.NewHandler(demoService)
	
	return &Routes{
		handler: handler,
	}
}

// Register 注册 Demo 路由
func (r *Routes) Register(rg *gin.RouterGroup) {
	demo := rg.Group("/demo/grpc")
	{
		// 健康检查
		demo.GET("/health", r.handler.HealthCheck)
		
		// 成员相关
		demo.GET("/members", r.handler.ListMembers)
		demo.GET("/members/:id", r.handler.GetMember)
		
		// 团队相关
		demo.GET("/teams", r.handler.ListTeams)
		demo.GET("/teams/:id", r.handler.GetTeam)
		
		// 调试信息
		demo.GET("/debug", r.handler.Debug)
	}
}

// GetPrefix 获取路由前缀
func (r *Routes) GetPrefix() string {
	return "/api/v1"
}

// GetName 获取路由名称
func (r *Routes) GetName() string {
	return "Demo API"
}