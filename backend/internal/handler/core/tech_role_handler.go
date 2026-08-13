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

type TechRoleHandler struct {
	repo repository.CoreTechRoleRepository
}

func NewTechRoleHandler(repo repository.CoreTechRoleRepository) *TechRoleHandler {
	return &TechRoleHandler{repo: repo}
}

func (h *TechRoleHandler) GetByID(c *gin.Context) {
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

func (h *TechRoleHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	result, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, result)
}

func (h *TechRoleHandler) Create(c *gin.Context) {
	var role core.TechRole
	if err := c.ShouldBindJSON(&role); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	role.ID = uuid.New()
	role.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &role); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, role)
}

func (h *TechRoleHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var role core.TechRole
	if err := c.ShouldBindJSON(&role); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &role); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, role)
}

func (h *TechRoleHandler) Delete(c *gin.Context) {
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
