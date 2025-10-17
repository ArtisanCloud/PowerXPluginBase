package security

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires the agent security namespace. Concrete handlers will be
// implemented alongside ToolGrant and privacy middleware in later phases.
func RegisterRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	_ = deps
	sec := rg.Group("/security")
	_ = sec
}
