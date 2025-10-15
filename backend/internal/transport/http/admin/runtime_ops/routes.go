package runtime_ops

import (
	runtimeops "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/runtime_ops"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires runtime ops endpoints behind the admin router.
func RegisterRoutes(router *gin.RouterGroup, deps *app.Deps) {
	bootstrap := NewBootstrapHandler(runtimeops.NewService())
	router.POST("/bootstrap", bootstrap.Bootstrap)

	sessions := NewSessionsHandler(deps)
	router.POST("/sessions/register", sessions.Register)
}
