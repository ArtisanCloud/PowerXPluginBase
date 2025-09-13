package http

import (
    "fmt"
    "github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
    "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin"
    agentapi "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/agent"
    "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/iam"
    "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/notes"
    "github.com/gin-gonic/gin"
)

// Registry API 注册器
type Registry struct {
	engine *gin.Engine
	deps   *app.Deps
}

// NewRegistry 创建 API 注册器
func NewRegistry(engine *gin.Engine, deps *app.Deps) *Registry {
	return &Registry{
		engine: engine,
		deps:   deps,
	}
}

// RegisterRoutes 注册所有路由
func (r *Registry) RegisterAPIRoutes(gApi *gin.RouterGroup) {
    admin.RegisterAPIRoutes(gApi, r.deps)
    agentapi.RegisterAPIRoutes(gApi, r.deps)
    iam.RegisterAPIRoutes(gApi, r.deps)
    notes.RegisterAPIRoutes(gApi, r.deps)
}

func (r *Registry) PrintRegisteredRoutes() {
	routes := r.engine.Routes()
	fmt.Println("==== Registered Routes ====")
	for _, route := range routes {
		// 格式化输出：方法、路径、处理函数
		fmt.Printf("%-6s %-30s %s\n", route.Method, route.Path, route.Handler)
	}
	fmt.Println("===========================")
}
