package integration

import "github.com/gin-gonic/gin"

// RegisterRoutes attaches integration admin endpoints to the supplied router group.
func RegisterRoutes(rg *gin.RouterGroup) {
	if rg == nil {
		return
	}
	// Concrete handlers will be wired in upcoming phases.
}
