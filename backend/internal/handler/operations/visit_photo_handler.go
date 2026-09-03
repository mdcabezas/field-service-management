package operationshandler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/operations"
	"localis-backend/internal/repository"
)

type VisitPhotoHandler struct {
	repo repository.VisitPhotoRepository
}

func NewVisitPhotoHandler(repo repository.VisitPhotoRepository) *VisitPhotoHandler {
	return &VisitPhotoHandler{repo: repo}
}

func (h *VisitPhotoHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *VisitPhotoHandler) ListByVisit(c *gin.Context) {
	visitID, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	items, err := h.repo.ListByVisit(c.Request.Context(), visitID)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *VisitPhotoHandler) Create(c *gin.Context) {
	var body operations.VisitPhoto
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	body.ID = uuid.New()
	body.CreatedAt = time.Now()
	// URL and ThumbnailURL are constructed by repository, not stored in DB
	if err := h.repo.Create(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	// Return with computed URLs
	urlStr := fmt.Sprintf("/api/photos/%s/file", body.ID.String())
	thumbStr := fmt.Sprintf("/api/photos/%s/thumb", body.ID.String())
	body.URL = &urlStr
	body.ThumbnailURL = &thumbStr
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *VisitPhotoHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
