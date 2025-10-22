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
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupOperationsServer(t *testing.T) (*gin.Engine, *operationsvc.SupportService) {
	gin.SetMode(gin.TestMode)
	db, err := gorm.Open(sqlite.Open("file:operations_support?mode=memory&cache=shared"), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true})
	require.NoError(t, err)
	basemodels.ForceSchemaForTests("")
	require.NoError(t, db.AutoMigrate(
		&opmodels.SupportChannel{},
		&opmodels.SupportTicket{},
		&opmodels.SupportTicketEvent{},
		&opmodels.ReadinessChecklistItem{},
	))

	repo := oprepo.NewSupportRepository(db)
	svc := operationsvc.NewSupportService(repo, &config.Config{}, opmetrics.NewMetrics(), nil)

	r := gin.New()
	deps := &app.Deps{DB: db, Config: &config.Config{}, OperationsMetrics: opmetrics.NewMetrics(), AdminConsoleMetrics: adminmetrics.NewMetrics()}
	adminGroup := r.Group("/admin")
	adminops.RegisterRoutes(adminGroup, deps)

	return r, svc
}

func TestSupportPlaybookEndpoints(t *testing.T) {
	router, svc := setupOperationsServer(t)

	payload := map[string]any{
		"channels": []map[string]any{
			{"channel": "marketplace_ticket", "address": "https://support.local", "escalates": []string{"agent"}},
		},
		"knowledge_base": []map[string]any{{"label": "FAQ", "url": "https://docs.local/faq"}},
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPut, "/admin/operations/support/playbook", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/support/playbook", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/admin/operations/support/channels/test", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/admin/operations/support/metrics", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	_, err := svc.CreateTicket(context.Background(), operationsvc.CreateTicketRequest{
		TenantID: "tenant-1",
		Subject:  "Outage",
		Priority: "P1",
	})
	require.NoError(t, err)
}
