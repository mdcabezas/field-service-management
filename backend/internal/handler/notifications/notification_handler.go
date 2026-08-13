package notificationshandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	notifmodel "localis-backend/internal/model/notifications"
	notifications "localis-backend/internal/service/notifications"
)

type NotificationHandler struct {
	svc *notifications.NotificationService
}

func NewNotificationHandler(svc *notifications.NotificationService) *NotificationHandler {
	return &NotificationHandler{svc: svc}
}

func (h *NotificationHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *NotificationHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	items, err := h.svc.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *NotificationHandler) Create(c *gin.Context) {
	var body notifmodel.Notification
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	created, err := h.svc.Send(c.Request.Context(), notifications.SendParams{
		Type:       body.Type,
		Channel:    body.Channel,
		Recipient:  body.Recipient,
		Subject:    body.Subject,
		Body:       body.Body,
		EntityType: body.EntityType,
		EntityID:   body.EntityID,
	})
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, created)
}

func (h *NotificationHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body notifmodel.Notification
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	updated, err := h.svc.Update(c.Request.Context(), id, &body)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, updated)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
