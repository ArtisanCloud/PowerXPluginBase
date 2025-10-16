package runtime_ops

import (
	"net/http"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	runtimeops "github.com/ArtisanCloud/PowerXPlugin/internal/services/admin/runtime_ops"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	"github.com/gin-gonic/gin"
)

// QuotaHandler exposes quota management endpoints.
type QuotaHandler struct {
	svc *runtimeops.QuotaService
}

// NewQuotaHandler constructs handler with service initialized from dependencies.
func NewQuotaHandler(deps *app.Deps, defaults *config.RuntimeOpsDefaults) *QuotaHandler {
	var db = deps.DB
	if defaults == nil && deps != nil && deps.Config != nil {
		defaults = deps.Config.RuntimeOps
	}
	return &QuotaHandler{svc: runtimeops.NewQuotaService(db, defaults)}
}

// GetStatus returns current quota utilization for a tenant scope (placeholder).
func (h *QuotaHandler) GetStatus(c *gin.Context) {
	pluginID := c.Query("plugin_id")
	tenantID := c.Query("tenant_id")
	defaults := h.svc.Defaults()
	window := 5 * time.Minute
	if defaults != nil && defaults.QuotaWindowMinutes > 0 {
		window = time.Duration(defaults.QuotaWindowMinutes) * time.Minute
	}
	if window <= 0 {
		window = 5 * time.Minute
	}
	start := time.Now().Add(-window)
	entries, err := h.svc.ListUsage(c.Request.Context(), "tenant", tenantID, start, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"plugin_id": pluginID,
		"tenant_id": tenantID,
		"ledger":    entries,
	})
}

// SetOverride accepts manual quota overrides (placeholder).
func (h *QuotaHandler) SetOverride(c *gin.Context) {
	var req struct {
		PluginID   string `json:"plugin_id"`
		TenantID   string `json:"tenant_id"`
		Capability string `json:"capability"`
		Action     string `json:"action"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	h.svc.HandleBreach(c.Request.Context(), req.PluginID, req.TenantID, req.Capability, req.Action)
	c.JSON(http.StatusAccepted, gin.H{"status": "override accepted"})
}
