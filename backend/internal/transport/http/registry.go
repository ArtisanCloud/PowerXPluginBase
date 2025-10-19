package http

import (
	"fmt"
	"strings"

	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin"
	adminintegration "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/integration"
	adminruntime "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/runtime_ops"
	adminsecurity "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/security"
	"github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/templates"
	agentapi "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/agent"
	integrationapi "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/integration"
	"github.com/gin-gonic/gin"
)

// Registry API 注册器
type Registry struct {
	engine *gin.Engine
	deps   *app.Deps
	rbac   map[string]authx.Permission
}

// NewRegistry 创建 API 注册器
func NewRegistry(engine *gin.Engine, deps *app.Deps) *Registry {
	return &Registry{
		engine: engine,
		deps:   deps,
		rbac:   map[string]authx.Permission{},
	}
}

// RegisterRoutes 注册所有路由
func (r *Registry) RegisterAPIRoutes(gApi *gin.RouterGroup) {
	admin.RegisterAPIRoutes(gApi, r.deps)
	agentapi.RegisterAPIRoutes(gApi, r.deps)
	templates.RegisterAPIRoutes(gApi, r.deps)
	integrationapi.RegisterAPIRoutes(gApi, r.deps)

	r.mergeRBAC(adminruntime.RBACEntries(r.apiPrefix()))
	r.mergeRBAC(adminsecurity.RBACEntries(r.apiPrefix()))
	r.mergeRBAC(adminintegration.RBACEntries(r.apiPrefix()))
	r.mergeRBAC(integrationRBACEntries(r.apiPrefix()))
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

// RBACMap 汇总所有模块的 RBAC 声明。
func (r *Registry) RBACMap() map[string]authx.Permission {
	out := make(map[string]authx.Permission, len(r.rbac))
	for route, perm := range r.rbac {
		out[route] = perm
	}
	return out
}

func (r *Registry) mergeRBAC(entries map[string]authx.Permission) {
	if entries == nil {
		return
	}
	for route, perm := range entries {
		r.rbac[route] = perm
	}
}

func (r *Registry) apiPrefix() string {
	prefix := "/api/v1"
	if r.deps != nil && r.deps.Config != nil && r.deps.Config.Server != nil {
		if p := strings.TrimSpace(r.deps.Config.Server.APIPrefix); p != "" {
			if !strings.HasPrefix(p, "/") {
				p = "/" + p
			}
			prefix = p
		}
	}
	return prefix
}
