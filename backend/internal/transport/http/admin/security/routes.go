package security

import (
	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the admin security namespace.
func RegisterRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	if rg == nil || deps == nil {
		return
	}
	auditWriter := CreateAuditWriter(deps.Config)
	h := NewConsentHandler(deps, auditWriter)
	sec := rg.Group("/security")
	{
		sec.GET("/consent-tokens", h.ListConsentTokens)
		sec.POST("/consent-tokens/:tokenId/revoke", h.RevokeConsentToken)
		sec.GET("/lifecycle-events", h.ListLifecycleEvents)
	}
}

// RBACEntries returns RBAC metadata for admin security endpoints.
func RBACEntries(prefix string) map[string]authx.Permission {
	if prefix == "" {
		prefix = "/api/v1"
	}
	return map[string]authx.Permission{
		prefix + "/admin/security/consent-tokens":          {Resource: "admin.security.consent", Action: "read"},
		prefix + "/admin/security/consent-tokens/:tokenId": {Resource: "admin.security.consent", Action: "write"},
		prefix + "/admin/security/lifecycle-events":        {Resource: "admin.security.lifecycle", Action: "read"},
	}
}
