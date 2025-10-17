package security

import (
	authx "github.com/ArtisanCloud/PowerXPlugin/internal/middleware"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the admin security namespace. Actual handlers will be
// added in subsequent implementation tasks.
func RegisterRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	_ = deps
	sec := rg.Group("/security")
	_ = sec
}

// RBACEntries returns RBAC metadata for admin security endpoints. Placeholder
// returns nil until concrete routes are defined.
func RBACEntries(prefix string) map[string]authx.Permission {
	_ = prefix
	return nil
}
