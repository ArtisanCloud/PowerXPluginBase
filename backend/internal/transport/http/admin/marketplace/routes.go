package marketplace

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes will be populated across user stories; placeholder keeps module wired.
func RegisterRoutes(admin *gin.RouterGroup, deps *app.Deps) {
	if admin == nil || deps == nil {
		return
	}
	// Concrete endpoints (listings, recommendation, pricing, analytics, etc.)
	// will be registered in subsequent implementation phases.
}
