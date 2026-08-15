package operationshandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	"localis-backend/internal/testutil/mocks"
)

func mockCtx() context.Context { return context.Background() }

func TestChecklistTemplateEPPHandler(t *testing.T) {
	repo := mocks.NewChecklistTemplateEPPRepository(t)
	id := uuid.New()
	tid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.ChecklistTemplateEPP{ID: id}, nil)
	repo.EXPECT().ListByTemplate(mockCtx(), tid).Return([]operations.ChecklistTemplateEPP{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ChecklistTemplateEPP")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewChecklistTemplateEPPHandler(repo)
	r := gin.New()
	r.GET("/cte/:id", h.GetByID)
	r.GET("/templates/:id/epps", h.ListByTemplate)
	r.POST("/cte", h.Create)
	r.DELETE("/cte/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/cte/"+id.String(), nil), http.StatusOK},
		{"ListByTemplate", httptest.NewRequest(http.MethodGet, "/templates/"+tid.String()+"/epps", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/cte", `{"template_id":"`+tid.String()+`","epp_id":"`+uuid.New().String()+`","default_quantity":1,"required":true}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/cte/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestChecklistTemplateMaterialHandler(t *testing.T) {
	repo := mocks.NewChecklistTemplateMaterialRepository(t)
	id := uuid.New()
	tid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.ChecklistTemplateMaterial{ID: id}, nil)
	repo.EXPECT().ListByTemplate(mockCtx(), tid).Return([]operations.ChecklistTemplateMaterial{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ChecklistTemplateMaterial")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewChecklistTemplateMaterialHandler(repo)
	r := gin.New()
	r.GET("/ctm/:id", h.GetByID)
	r.GET("/templates/:id/materials", h.ListByTemplate)
	r.POST("/ctm", h.Create)
	r.DELETE("/ctm/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/ctm/"+id.String(), nil), http.StatusOK},
		{"ListByTemplate", httptest.NewRequest(http.MethodGet, "/templates/"+tid.String()+"/materials", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/ctm", `{"template_id":"`+tid.String()+`","material_id":"`+uuid.New().String()+`","default_quantity":1,"required":true}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/ctm/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestChecklistTemplateToolHandler(t *testing.T) {
	repo := mocks.NewChecklistTemplateToolRepository(t)
	id := uuid.New()
	tid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.ChecklistTemplateTool{ID: id}, nil)
	repo.EXPECT().ListByTemplate(mockCtx(), tid).Return([]operations.ChecklistTemplateTool{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ChecklistTemplateTool")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewChecklistTemplateToolHandler(repo)
	r := gin.New()
	r.GET("/ctt/:id", h.GetByID)
	r.GET("/templates/:id/tools", h.ListByTemplate)
	r.POST("/ctt", h.Create)
	r.DELETE("/ctt/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/ctt/"+id.String(), nil), http.StatusOK},
		{"ListByTemplate", httptest.NewRequest(http.MethodGet, "/templates/"+tid.String()+"/tools", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/ctt", `{"template_id":"`+tid.String()+`","tool_id":"`+uuid.New().String()+`","default_quantity":1,"required":true}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/ctt/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestReportEntryHandler(t *testing.T) {
	repo := mocks.NewReportEntryRepository(t)
	id := uuid.New()
	rid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.ReportEntry{ID: id}, nil)
	repo.EXPECT().ListByReport(mockCtx(), rid).Return([]operations.ReportEntry{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ReportEntry")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewReportEntryHandler(repo)
	r := gin.New()
	r.GET("/re/:id", h.GetByID)
	r.GET("/reports/:id/entries", h.ListByReport)
	r.POST("/re", h.Create)
	r.DELETE("/re/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/re/"+id.String(), nil), http.StatusOK},
		{"ListByReport", httptest.NewRequest(http.MethodGet, "/reports/"+rid.String()+"/entries", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/re", `{"report_id":"`+rid.String()+`","data_json":{"a":1}}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/re/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestReportImageHandler(t *testing.T) {
	repo := mocks.NewReportImageRepository(t)
	id := uuid.New()
	rid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.ReportImage{ID: id}, nil)
	repo.EXPECT().ListByReport(mockCtx(), rid).Return([]operations.ReportImage{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.ReportImage")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewReportImageHandler(repo)
	r := gin.New()
	r.GET("/ri/:id", h.GetByID)
	r.GET("/reports/:id/images", h.ListByReport)
	r.POST("/ri", h.Create)
	r.DELETE("/ri/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/ri/"+id.String(), nil), http.StatusOK},
		{"ListByReport", httptest.NewRequest(http.MethodGet, "/reports/"+rid.String()+"/images", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/ri", `{"report_id":"`+rid.String()+`","url":"http://img","pages":2}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/ri/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitChecklistEPPHandler(t *testing.T) {
	repo := mocks.NewVisitChecklistEPPRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitChecklistEPP{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitChecklistEPP{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitChecklistEPP")).Return(nil)
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VisitChecklistEPP")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitChecklistEPPHandler(repo)
	r := gin.New()
	r.GET("/vce/:id", h.GetByID)
	r.GET("/visits/:id/epps", h.ListByVisit)
	r.POST("/vce", h.Create)
	r.PUT("/vce/:id", h.Update)
	r.DELETE("/vce/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","epp_id":"` + uuid.New().String() + `","planned_quantity":2,"confirmed":true}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vce/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/epps", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vce", body), http.StatusCreated},
		{"Update", reqJSON(http.MethodPut, "/vce/"+id.String(), body), http.StatusOK},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vce/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitChecklistMaterialHandler(t *testing.T) {
	repo := mocks.NewVisitChecklistMaterialRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitChecklistMaterial{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitChecklistMaterial{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitChecklistMaterial")).Return(nil)
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VisitChecklistMaterial")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitChecklistMaterialHandler(repo)
	r := gin.New()
	r.GET("/vcm/:id", h.GetByID)
	r.GET("/visits/:id/materials", h.ListByVisit)
	r.POST("/vcm", h.Create)
	r.PUT("/vcm/:id", h.Update)
	r.DELETE("/vcm/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","material_id":"` + uuid.New().String() + `","planned_quantity":2,"confirmed":true}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vcm/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/materials", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vcm", body), http.StatusCreated},
		{"Update", reqJSON(http.MethodPut, "/vcm/"+id.String(), body), http.StatusOK},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vcm/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitChecklistToolHandler(t *testing.T) {
	repo := mocks.NewVisitChecklistToolRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitChecklistTool{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitChecklistTool{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitChecklistTool")).Return(nil)
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VisitChecklistTool")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitChecklistToolHandler(repo)
	r := gin.New()
	r.GET("/vct/:id", h.GetByID)
	r.GET("/visits/:id/tools", h.ListByVisit)
	r.POST("/vct", h.Create)
	r.PUT("/vct/:id", h.Update)
	r.DELETE("/vct/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","tool_id":"` + uuid.New().String() + `","planned_quantity":2,"confirmed":true}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vct/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/tools", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vct", body), http.StatusCreated},
		{"Update", reqJSON(http.MethodPut, "/vct/"+id.String(), body), http.StatusOK},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vct/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitCheckpointHandler(t *testing.T) {
	repo := mocks.NewVisitCheckpointRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitCheckpoint{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitCheckpoint{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitCheckpoint")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitCheckpointHandler(repo)
	r := gin.New()
	r.GET("/vcp/:id", h.GetByID)
	r.GET("/visits/:id/checkpoints", h.ListByVisit)
	r.POST("/vcp", h.Create)
	r.DELETE("/vcp/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","type":"arrival","geom":[1,1],"timestamp":"2026-01-01T00:00:00Z"}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vcp/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/checkpoints", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vcp", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vcp/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitEPPUsageHandler(t *testing.T) {
	repo := mocks.NewVisitEPPUsageRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitEPPUsage{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitEPPUsage{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitEPPUsage")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitEPPUsageHandler(repo)
	r := gin.New()
	r.GET("/veu/:id", h.GetByID)
	r.GET("/visits/:id/epp-usages", h.ListByVisit)
	r.POST("/veu", h.Create)
	r.DELETE("/veu/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","epp_id":"` + uuid.New().String() + `","quantity":1,"status":"used"}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/veu/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/epp-usages", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/veu", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/veu/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitMaterialUsageHandler(t *testing.T) {
	repo := mocks.NewVisitMaterialUsageRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitMaterialUsage{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitMaterialUsage{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitMaterialUsage")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitMaterialUsageHandler(repo)
	r := gin.New()
	r.GET("/vmu/:id", h.GetByID)
	r.GET("/visits/:id/material-usages", h.ListByVisit)
	r.POST("/vmu", h.Create)
	r.DELETE("/vmu/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","material_id":"` + uuid.New().String() + `","quantity":1.5}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vmu/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/material-usages", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vmu", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vmu/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitMeasurementHandler(t *testing.T) {
	repo := mocks.NewVisitMeasurementRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitMeasurement{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitMeasurement{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitMeasurement")).Return(nil)
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*operations.VisitMeasurement")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitMeasurementHandler(repo)
	r := gin.New()
	r.GET("/vm/:id", h.GetByID)
	r.GET("/visits/:id/measurements", h.ListByVisit)
	r.POST("/vm", h.Create)
	r.PUT("/vm/:id", h.Update)
	r.DELETE("/vm/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","type":"` + uuid.New().String() + `","result":"approved"}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vm/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/measurements", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vm", body), http.StatusCreated},
		{"Update", reqJSON(http.MethodPut, "/vm/"+id.String(), body), http.StatusOK},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vm/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitPhotoHandler(t *testing.T) {
	repo := mocks.NewVisitPhotoRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitPhoto{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitPhoto{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitPhoto")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitPhotoHandler(repo)
	r := gin.New()
	r.GET("/vp/:id", h.GetByID)
	r.GET("/visits/:id/photos", h.ListByVisit)
	r.POST("/vp", h.Create)
	r.DELETE("/vp/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","url":"http://photo","geom":[1,1],"timestamp":"2026-01-01T00:00:00Z","stage":"diagnosis"}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vp/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/photos", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vp", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vp/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitRentalHandler(t *testing.T) {
	repo := mocks.NewVisitRentalRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitRental{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitRental{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitRental")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitRentalHandler(repo)
	r := gin.New()
	r.GET("/vr/:id", h.GetByID)
	r.GET("/visits/:id/rentals", h.ListByVisit)
	r.POST("/vr", h.Create)
	r.DELETE("/vr/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","rental_id":"` + uuid.New().String() + `","start_time":"2026-01-01T08:00:00Z","total_cost":100}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vr/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/rentals", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vr", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vr/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func TestVisitToolUsageHandler(t *testing.T) {
	repo := mocks.NewVisitToolUsageRepository(t)
	id := uuid.New()
	vid := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitToolUsage{ID: id}, nil)
	repo.EXPECT().ListByVisit(mockCtx(), vid).Return([]operations.VisitToolUsage{}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitToolUsage")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitToolUsageHandler(repo)
	r := gin.New()
	r.GET("/vtu/:id", h.GetByID)
	r.GET("/visits/:id/tool-usages", h.ListByVisit)
	r.POST("/vtu", h.Create)
	r.DELETE("/vtu/:id", h.Delete)

	body := `{"visit_id":"` + vid.String() + `","tool_id":"` + uuid.New().String() + `"}`
	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/vtu/"+id.String(), nil), http.StatusOK},
		{"ListByVisit", httptest.NewRequest(http.MethodGet, "/visits/"+vid.String()+"/tool-usages", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/vtu", body), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/vtu/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

func reqJSON(method, path, body string) *http.Request {
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func TestVisitTypeHandler(t *testing.T) {
	repo := mocks.NewVisitTypeRepository(t)
	id := uuid.New()
	repo.EXPECT().GetByID(mockCtx(), id).Return(&operations.VisitType{ID: id}, nil)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[operations.VisitType]{Total: 0}, nil)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*operations.VisitType")).Return(nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	h := NewVisitTypeHandler(repo)
	r := gin.New()
	r.GET("/visit-types/:id", h.GetByID)
	r.GET("/visit-types", h.List)
	r.POST("/visit-types", h.Create)
	r.DELETE("/visit-types/:id", h.Delete)

	cases := []struct {
		name string
		req  *http.Request
		code int
	}{
		{"GetByID", httptest.NewRequest(http.MethodGet, "/visit-types/"+id.String(), nil), http.StatusOK},
		{"List", httptest.NewRequest(http.MethodGet, "/visit-types", nil), http.StatusOK},
		{"Create", reqJSON(http.MethodPost, "/visit-types", `{"code":"INSTALL","name":"Installation","active":true}`), http.StatusCreated},
		{"Delete", httptest.NewRequest(http.MethodDelete, "/visit-types/"+id.String(), nil), http.StatusNoContent},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r.ServeHTTP(w, tc.req)
			if w.Code != tc.code {
				t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
			}
		})
	}
}

var _ = shared.CheckpointType("arrival")
var _ = shared.EPPUsageStatus("used")
var _ = shared.MeasurementResult("approved")
var _ = shared.PhotoStage("diagnosis")
