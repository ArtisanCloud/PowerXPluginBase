package marketplace

import (
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	httpmw "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes wires marketplace admin endpoints.
func RegisterRoutes(admin *gin.RouterGroup, deps *app.Deps) {
	if admin == nil || deps == nil {
		return
	}
	handler := NewListingHandler(deps)
	checklist := NewChecklistGraphQLHandler(handler.Service())

	group := admin.Group("/marketplace", httpmw.EnsureTenant())
	{
		listings := group.Group("/listings")
		{
			listings.GET("", handler.List)
			listings.POST("", handler.Create)
			listings.GET("/:id", handler.Get)
			listings.PATCH("/:id", handler.Update)
			listings.POST("/:id/review", handler.SubmitForReview)
			listings.POST("/:id/publish", handler.Publish)
			listings.POST("/:id/suspend", handler.Suspend)
		}

		group.POST("/checklist/graphql", checklist.Resolve)
	}
}
