package middleware

import (
	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NeedABACFn：可选的 ABAC 触发回调（返回是否需要 ABAC 以及附加属性）
type NeedABACFn func(method, route string) (need bool, attrs map[string]any)

// RBAC：粗粒度 RBAC；命中需要 ABAC 的路由时，调用 PDP 在线校验
func RBAC(cfg *authx.RBACConfig, abac authx.ABACClient, needABAC NeedABACFn) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 预检查
		if cfg == nil || !cfg.Enabled || c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}
		tc, ok := authx.GetTenantContext(c)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}
		if authx.IsSuperAdmin(tc.Roles, cfg.SuperAdminRoles) {
			c.Next()
			return
		}

		full := c.FullPath()
		if full == "" {
			full = c.Request.URL.Path // 非路由表命中的场景
		}
		perm, has := authx.MatchRoute(c.Request.Method, full, cfg.RoutePermissions)

		// 粗粒度判定
		passRBAC := (!has && !cfg.DefaultDeny) || (has && authx.HasPerm(tc.Permissions, perm))
		if !passRBAC {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":             "Insufficient permissions",
				"required_resource": perm.Resource,
				"required_action":   perm.Action,
			})
			return
		}

		// 是否需要触发 ABAC（可选）
		if needABAC != nil && abac != nil {
			if yes, attrs := needABAC(c.Request.Method, full); yes {
				dec, err := abac.Check(c.Request.Context(), authx.ABACInput{
					Subject:  tc,
					Resource: perm.Resource,
					Action:   perm.Action,
					Attrs:    attrs,
				})
				if err != nil {
					// PDP 不可用：按建议返回 503 让上游重试
					c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "PDP unavailable"})
					return
				}
				if !dec.Allowed {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ABAC denied", "reason": dec.Reason})
					return
				}
			}
		}

		c.Next()
	}
}
