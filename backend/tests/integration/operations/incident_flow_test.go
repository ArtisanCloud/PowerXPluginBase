package operations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	basemodels "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	opmodels "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/operations"
	oprepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/operations"
	adminmetrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/admin_console"
	opmetrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/operations"
	operationsvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/operations"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	adminops "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/operations"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupIncidentServer(t *testing.T) (*gin.Engine, *operationsvc.IncidentService) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:incident_flow?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	require.NoError(t, err)
	basemodels.ForceSchemaForTests("")
	require.NoError(t, db.AutoMigrate(
		&opmodels.Incident{},
		&opmodels.IncidentTimelineEntry{},
		&opmodels.IncidentChecklistItem{},
		&opmodels.ReadinessChecklistItem{},
	))

	repo := oprepo.NewIncidentRepository(db)
	svc := operationsvc.NewIncidentService(repo, &config.Config{}, opmetrics.NewMetrics(), nil)

	r := gin.New()
	deps := &app.Deps{DB: db, Config: &config.Config{}, OperationsMetrics: opmetrics.NewMetrics(), AdminConsoleMetrics: adminmetrics.NewMetrics()}
	adminGroup := r.Group("/admin")
	adminops.RegisterRoutes(adminGroup, deps)

	return r, svc
}

func TestIncidentEndpoints(t *testing.T) {
	router, svc := setupIncidentServer(t)

	createPayload := map[string]any{
		"severity":         "sev1",
		"detection_source": "monitoring",
		"summary":          "API latency spike",
	}
	body, _ := json.Marshal(createPayload)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/incidents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var resp struct {
		Data struct {
			Incident struct {
				ID     string `json:"id"`
				Status string `json:"status"`
			} `json:"incident"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	incidentID := resp.Data.Incident.ID
	require.NotEmpty(t, incidentID)
	require.Equal(t, "/admin/operations/incidents/"+incidentID, rec.Header().Get("Location"))

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/incidents", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	updatePayload := map[string]any{
		"status":         "acknowledged",
		"mitigation":     "rerouted traffic",
		"next_update_at": time.Now().Add(30 * time.Minute).Format(time.RFC3339),
	}
	body, _ = json.Marshal(updatePayload)
	req = httptest.NewRequest(http.MethodPatch, "/admin/operations/incidents/"+incidentID, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var updateResp struct {
		Data struct {
			Incident struct {
				Status string `json:"status"`
			} `json:"incident"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &updateResp))
	require.Equal(t, "acknowledged", updateResp.Data.Incident.Status)

	timelinePayload := map[string]any{
		"entry_type":          "announcement",
		"message":             "Issue acknowledged",
		"stakeholder_channel": "status_page",
	}
	body, _ = json.Marshal(timelinePayload)
	req = httptest.NewRequest(http.MethodPost, "/admin/operations/incidents/"+incidentID+"/timeline", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	view, err := svc.GetIncident(context.Background(), incidentID)
	require.NoError(t, err)
	require.Equal(t, "announcement", view.Timeline[0].EntryType)
}
