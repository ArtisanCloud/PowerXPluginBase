package runtime_ops

import (
	runtimeops "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/runtime_ops"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires runtime ops endpoints behind the admin router.
func RegisterRoutes(router gin.IRoutes) {
	handler := NewHandler(runtimeops.NewService())
	router.POST("/bootstrap", handler.Bootstrap)
}
