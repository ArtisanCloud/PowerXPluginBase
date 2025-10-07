package templates

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(rg *gin.RouterGroup, deps *app.Deps) {
	h := NewTemplateHandler(deps)

	g := rg.Group("/templates")
	{
		g.GET("", h.GetTemplates)
		g.GET("/:id", h.GetTemplate)
		g.POST("", h.CreateTemplate)
		g.PUT("/:id", h.UpdateTemplate)
		g.DELETE("/:id", h.DeleteTemplate)
	}
}
