package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/powerx-plugins/scrum/internal/logger"
)

// Permission 权限结构
type Permission struct {
	Resource string `json:"resource"`
	Action   string `json:"action"`
}

// RBAC 中间件配置
type RBACConfig struct {
	// 是否启用 RBAC
	Enabled bool
	// 默认拒绝访问
	DefaultDeny bool
	// 超级管理员角色
	SuperAdminRoles []string
	// 权限映射：路径 -> 所需权限
	RoutePermissions map[string]Permission
}

// NewRBACConfig 创建默认 RBAC 配置
func NewRBACConfig() *RBACConfig {
	return &RBACConfig{
		Enabled:         true,
		DefaultDeny:     true,
		SuperAdminRoles: []string{"super_admin", "system_admin"},
		RoutePermissions: map[string]Permission{
			// Admin 路由
			"GET:/admin/manifest": {Resource: "system", Action: "read"},
			"GET:/admin/rbac":     {Resource: "system", Action: "read"},

			// Task 管理
			"GET:/tasks":      {Resource: "scrum:task", Action: "read"},
			"POST:/tasks":     {Resource: "scrum:task", Action: "create"},
			"PUT:/tasks/*":    {Resource: "scrum:task", Action: "update"},
			"DELETE:/tasks/*": {Resource: "scrum:task", Action: "delete"},

			// Sprint 管理
			"GET:/sprints":      {Resource: "scrum:sprint", Action: "read"},
			"POST:/sprints":     {Resource: "scrum:sprint", Action: "create"},
			"PUT:/sprints/*":    {Resource: "scrum:sprint", Action: "update"},
			"DELETE:/sprints/*": {Resource: "scrum:sprint", Action: "delete"},
		},
	}
}

// RBACMiddleware 创建 RBAC 权限校验中间件
func RBACMiddleware(config *RBACConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !config.Enabled {
			c.Next()
			return
		}

		tenantCtx, exists := GetTenantContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required for RBAC",
			})
			c.Abort()
			return
		}

		// 检查超级管理员权限
		if hasAnyRole(tenantCtx.Roles, config.SuperAdminRoles) {
			logger.AuthMiddleware().WithFields(logger.Logger.Fields{
				"tenant_id": tenantCtx.TenantID,
				"user_id":   tenantCtx.UserID,
				"roles":     tenantCtx.Roles,
			}).Debug("Super admin access granted")
			c.Next()
			return
		}

		// 获取当前路由所需权限
		requiredPerm := getRequiredPermission(c, config)
		if requiredPerm == nil {
			if config.DefaultDeny {
				logger.AuthMiddleware().WithFields(logger.Logger.Fields{
					"tenant_id": tenantCtx.TenantID,
					"user_id":   tenantCtx.UserID,
					"path":      c.Request.URL.Path,
					"method":    c.Request.Method,
				}).Warn("Access denied: no permission defined and default deny is enabled")

				c.JSON(http.StatusForbidden, gin.H{
					"error": "Access denied: no permission defined for this route",
				})
				c.Abort()
				return
			}
			// 默认允许访问
			c.Next()
			return
		}

		// 检查权限
		if !hasPermission(tenantCtx.Permissions, *requiredPerm) {
			logger.AuthMiddleware().WithFields(logger.Logger.Fields{
				"tenant_id":         tenantCtx.TenantID,
				"user_id":           tenantCtx.UserID,
				"required_resource": requiredPerm.Resource,
				"required_action":   requiredPerm.Action,
				"user_permissions":  tenantCtx.Permissions,
			}).Warn("Access denied: insufficient permissions")

			c.JSON(http.StatusForbidden, gin.H{
				"error":             "Insufficient permissions",
				"required_resource": requiredPerm.Resource,
				"required_action":   requiredPerm.Action,
			})
			c.Abort()
			return
		}

		logger.AuthMiddleware().WithFields(logger.Logger.Fields{
			"tenant_id": tenantCtx.TenantID,
			"user_id":   tenantCtx.UserID,
			"resource":  requiredPerm.Resource,
			"action":    requiredPerm.Action,
		}).Debug("Access granted")

		c.Next()
	}
}

// getRequiredPermission 获取当前路由所需的权限
func getRequiredPermission(c *gin.Context, config *RBACConfig) *Permission {
	path := c.Request.URL.Path
	method := c.Request.Method

	// 构造路由键
	routeKey := method + ":" + path

	// 精确匹配
	if perm, exists := config.RoutePermissions[routeKey]; exists {
		return &perm
	}

	// 通配符匹配
	for pattern, perm := range config.RoutePermissions {
		if matchRoute(pattern, routeKey) {
			return &perm
		}
	}

	return nil
}

// matchRoute 匹配路由模式
func matchRoute(pattern, route string) bool {
	// 简单的通配符匹配实现
	if strings.Contains(pattern, "*") {
		patternParts := strings.Split(pattern, "*")
		if len(patternParts) == 2 {
			prefix := patternParts[0]
			suffix := patternParts[1]
			return strings.HasPrefix(route, prefix) && strings.HasSuffix(route, suffix)
		}
	}
	return false
}

// hasPermission 检查用户是否具有指定权限
func hasPermission(userPerms []string, required Permission) bool {
	requiredPerm := required.Resource + ":" + required.Action

	for _, perm := range userPerms {
		// 完全匹配
		if perm == requiredPerm {
			return true
		}

		// 通配符权限
		if perm == "*" || perm == required.Resource+":*" {
			return true
		}

		// 资源级通配符
		if strings.HasSuffix(perm, ":*") {
			resource := strings.TrimSuffix(perm, ":*")
			if resource == required.Resource {
				return true
			}
		}
	}

	return false
}

// hasAnyRole 检查用户是否具有任一指定角色
func hasAnyRole(userRoles, requiredRoles []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range userRoles {
		roleMap[role] = true
	}

	for _, required := range requiredRoles {
		if roleMap[required] {
			return true
		}
	}

	return false
}

// RequirePermission 创建需要特定权限的中间件
func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCtx, exists := GetTenantContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		required := Permission{Resource: resource, Action: action}
		if !hasPermission(tenantCtx.Permissions, required) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":             "Insufficient permissions",
				"required_resource": resource,
				"required_action":   action,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 创建需要特定角色的中间件
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCtx, exists := GetTenantContext(c)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		if !hasAnyRole(tenantCtx.Roles, roles) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Insufficient roles",
				"required_roles": roles,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
