package operations_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ArtisanCloud/PowerXPlugin/internal/config"
	basemodels "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models"
	opmodels "github.com/ArtisanCloud/PowerXPlugin/internal/domain/models/operations"
	oprepo "github.com/ArtisanCloud/PowerXPlugin/internal/domain/repository/operations"
	adminmetrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/admin_console"
	opmetrics "github.com/ArtisanCloud/PowerXPlugin/internal/observability/operations"
	operationsvc "github.com/ArtisanCloud/PowerXPlugin/internal/services/operations"
	"github.com/ArtisanCloud/PowerXPlugin/internal/shared/app"
	adminops "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/admin/operations"
	publicmarketplace "github.com/ArtisanCloud/PowerXPlugin/internal/transport/http/public/marketplace"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupSLAEndpoints(t *testing.T) (*gin.Engine, *oprepo.SLARepository) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:sla_flow?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	require.NoError(t, err)
	basemodels.ForceSchemaForTests("")
	require.NoError(t, db.AutoMigrate(
		&opmodels.SLAProfile{},
		&opmodels.SLAAdjustment{},
		&opmodels.ReadinessChecklistItem{},
	))

	repo := oprepo.NewSLARepository(db)
	svc := operationsvc.NewSLAService(repo, &config.Config{}, opmetrics.NewMetrics())

	r := gin.New()
	deps := &app.Deps{DB: db, Config: &config.Config{}, OperationsMetrics: opmetrics.NewMetrics(), AdminConsoleMetrics: adminmetrics.NewMetrics()}
	adminGroup := r.Group("/admin")
	adminops.RegisterRoutes(adminGroup, deps)

	publicGroup := r.Group("/api/v1/marketplace")
	publicmarketplace.Register(publicGroup, publicmarketplace.NewSLAHandler(repo, svc))

	return r, repo
}

func TestSLAAdminAndPublicFlow(t *testing.T) {
	router, repo := setupSLAEndpoints(t)

	targetPayload := map[string]any{
		"planType": "real_time",
		"targets": map[string]any{
			"uptimeTarget":          99.9,
			"responseTargetMs":      600,
			"successTargetPct":      99.5,
			"supportFrtTargetHours": 4,
		},
	}
	body, _ := json.Marshal(targetPayload)
	req := httptest.NewRequest(http.MethodPost, "/admin/operations/sla/profiles", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	actualPayload := map[string]any{
		"planType": "real_time",
		"actuals": map[string]any{
			"uptimeActual":          99.95,
			"responseActualMs":      300,
			"successActualPct":      99.8,
			"supportFrtActualHours": 2,
		},
	}
	body, _ = json.Marshal(actualPayload)
	req = httptest.NewRequest(http.MethodPatch, "/admin/operations/sla/profiles/actuals", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/sla/profiles", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/admin/operations/sla/profiles/recompute", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusAccepted, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/marketplace/sla/"+app.PluginID, nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	adjustments, err := repo.ListAdjustments(context.Background(), app.PluginID, "real_time", 10)
	require.NoError(t, err)
	require.NotEmpty(t, adjustments)
}
