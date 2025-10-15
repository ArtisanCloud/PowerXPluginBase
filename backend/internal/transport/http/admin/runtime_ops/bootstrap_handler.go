package runtime_ops

import (
	"net/http"

	runtimeops "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/runtime_ops"
	"github.com/gin-gonic/gin"
)

// Handler exposes runtime ops HTTP endpoints.
type Handler struct {
	svc *runtimeops.Service
}

// NewHandler builds a runtime ops handler.
func NewHandler(svc *runtimeops.Service) *Handler {
	if svc == nil {
		svc = runtimeops.NewService()
	}
	return &Handler{svc: svc}
}

// Bootstrap is a placeholder handler for launching plugin runtime instances.
func (h *Handler) Bootstrap(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "runtime bootstrap endpoint not implemented"})
}
