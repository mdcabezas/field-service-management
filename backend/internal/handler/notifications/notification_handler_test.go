package notificationshandler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	notifmodel "localis-backend/internal/model/notifications"
	"localis-backend/internal/model/shared"
	"localis-backend/internal/repository"
	services "localis-backend/internal/service/notifications"
	"localis-backend/internal/testutil/mocks"
)

func newNotifRouter(t *testing.T, h *NotificationHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/notifications", h.List)
	r.GET("/notifications/:id", h.GetByID)
	r.POST("/notifications", h.Create)
	r.PUT("/notifications/:id", h.Update)
	r.DELETE("/notifications/:id", h.Delete)
	return r
}

func newTemplateRouter(t *testing.T, h *NotificationTemplateHandler) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/notification-templates", h.List)
	r.GET("/notification-templates/:id", h.GetByID)
	r.POST("/notification-templates", h.Create)
	r.PUT("/notification-templates/:id", h.Update)
	r.DELETE("/notification-templates/:id", h.Delete)
	return r
}

func mockCtx() context.Context { return context.Background() }

func newNotifHandler(t *testing.T) (*NotificationHandler, *mocks.NotificationRepository) {
	repo := mocks.NewNotificationRepository(t)
	return NewNotificationHandler(services.NewNotificationService(repo)), repo
}

func TestNotificationHandler_GetByID(t *testing.T) {
	h, repo := newNotifHandler(t)
	id := uuid.New()
	n := &notifmodel.Notification{ID: id, Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(n, nil)

	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_GetByID_InvalidUUID(t *testing.T) {
	h, _ := newNotifHandler(t)
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_List(t *testing.T) {
	h, repo := newNotifHandler(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[notifmodel.Notification]{Total: 0}, nil)

	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_Create_DispatchFails(t *testing.T) {
	h, repo := newNotifHandler(t)
	repo.EXPECT().Create(mockCtx(), mock.Anything).Return(nil)
	repo.EXPECT().Update(mockCtx(), mock.Anything, mock.Anything).Return(nil)

	body := `{"type":"other","channel":"email","recipient":"a@b.cl","status":"pending"}`
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500 (dispatch not implemented)", w.Code)
	}
}

func TestNotificationHandler_Create_InvalidBody(t *testing.T) {
	h, _ := newNotifHandler(t)
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notifications", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_Update(t *testing.T) {
	h, repo := newNotifHandler(t)
	id := uuid.New()
	existing := &notifmodel.Notification{ID: id, Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(existing, nil)
	repo.EXPECT().Update(mockCtx(), id, mock.Anything).Return(nil)

	body := `{"type":"other","channel":"email","recipient":"a@b.cl","status":"pending"}`
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notifications/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestNotificationHandler_Update_InvalidUUID(t *testing.T) {
	h, _ := newNotifHandler(t)
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notifications/bad", bytes.NewBufferString(`{"type":"other","channel":"email","recipient":"a@b.cl","status":"pending"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_Update_Error(t *testing.T) {
	h, repo := newNotifHandler(t)
	id := uuid.New()
	existing := &notifmodel.Notification{ID: id, Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(existing, nil)
	repo.EXPECT().Update(mockCtx(), id, mock.Anything).Return(&repoErr{})

	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notifications/"+id.String(), bytes.NewBufferString(`{"type":"other","channel":"email","recipient":"a@b.cl","status":"pending"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_Delete(t *testing.T) {
	h, repo := newNotifHandler(t)
	id := uuid.New()
	existing := &notifmodel.Notification{ID: id, Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Recipient: "a@b.cl"}
	repo.EXPECT().GetByID(mockCtx(), id).Return(existing, nil)
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/notifications/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationHandler_Delete_InvalidUUID(t *testing.T) {
	h, _ := newNotifHandler(t)
	r := newNotifRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/notifications/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func newTemplateHandler(t *testing.T) (*NotificationTemplateHandler, *mocks.NotificationTemplateRepository) {
	repo := mocks.NewNotificationTemplateRepository(t)
	return NewNotificationTemplateHandler(repo), repo
}

func TestNotificationTemplateHandler_GetByID(t *testing.T) {
	h, repo := newTemplateHandler(t)
	id := uuid.New()
	tmpl := &notifmodel.NotificationTemplate{ID: id, Type: shared.NotificationTypeOther, Channel: shared.NotificationChannelEmail, Active: true}
	repo.EXPECT().GetByID(mockCtx(), id).Return(tmpl, nil)

	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notification-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_GetByID_InvalidUUID(t *testing.T) {
	h, _ := newTemplateHandler(t)
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notification-templates/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_List(t *testing.T) {
	h, repo := newTemplateHandler(t)
	repo.EXPECT().List(mockCtx(), 20, 0).Return(&repository.ListResult[notifmodel.NotificationTemplate]{Total: 0}, nil)

	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/notification-templates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Create(t *testing.T) {
	h, repo := newTemplateHandler(t)
	repo.EXPECT().Create(mockCtx(), mock.AnythingOfType("*notifications.NotificationTemplate")).Return(nil)

	body := `{"type":"other","channel":"email","active":true}`
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notification-templates", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Create_InvalidBody(t *testing.T) {
	h, _ := newTemplateHandler(t)
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/notification-templates", bytes.NewBufferString("{"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Update(t *testing.T) {
	h, repo := newTemplateHandler(t)
	id := uuid.New()
	repo.EXPECT().Update(mockCtx(), id, mock.AnythingOfType("*notifications.NotificationTemplate")).Return(nil)

	body := `{"type":"other","channel":"email","active":true}`
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notification-templates/"+id.String(), bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Update_InvalidUUID(t *testing.T) {
	h, _ := newTemplateHandler(t)
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/notification-templates/bad", bytes.NewBufferString(`{"type":"other","channel":"email","active":true}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Delete(t *testing.T) {
	h, repo := newTemplateHandler(t)
	id := uuid.New()
	repo.EXPECT().Delete(mockCtx(), id).Return(nil)

	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/notification-templates/"+id.String(), nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestNotificationTemplateHandler_Delete_InvalidUUID(t *testing.T) {
	h, _ := newTemplateHandler(t)
	r := newTemplateRouter(t, h)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/notification-templates/bad", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

type repoErr struct{}

func (*repoErr) Error() string { return "boom" }
