package corehandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/core"
	"localis-backend/internal/repository"
)

type UserHandler struct {
	repo repository.CoreUserRepository
}

func NewUserHandler(repo repository.CoreUserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	item, err := h.repo.GetByID(c.Request.Context(), id.String())
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, item)
}

func (h *UserHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *UserHandler) Create(c *gin.Context) {
	var user core.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &user); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, user)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var user core.User
	if err := c.ShouldBindJSON(&user); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id.String(), &user); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.repo.Delete(c.Request.Context(), id.String()); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusNoContent, nil)
}
