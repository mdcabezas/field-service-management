package geocodinghandler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/geocoding"
	"localis-backend/internal/repository"
)

type AddressHandler struct {
	repo repository.GeocodingAddressRepository
}

func NewAddressHandler(repo repository.GeocodingAddressRepository) *AddressHandler {
	return &AddressHandler{repo: repo}
}

func (h *AddressHandler) GetByID(c *gin.Context) {
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

func (h *AddressHandler) List(c *gin.Context) {
	limit, offset := handler.ParsePagination(c)
	items, err := h.repo.List(c.Request.Context(), limit, offset)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}

func (h *AddressHandler) Create(c *gin.Context) {
	var body geocoding.Address
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	body.ID = uuid.New()
	body.CreatedAt = time.Now()
	if err := h.repo.Create(c.Request.Context(), &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusCreated, body)
}

func (h *AddressHandler) Update(c *gin.Context) {
	id, ok := handler.ParseUUID(c, "id")
	if !ok {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}
	var body geocoding.Address
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.repo.Update(c.Request.Context(), id, &body); err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, body)
}

func (h *AddressHandler) Delete(c *gin.Context) {
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

func (h *AddressHandler) FindNearby(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.Query("radius")

	if latStr == "" || lngStr == "" || radiusStr == "" {
		handler.RespondError(c, http.StatusBadRequest, "lat, lng, and radius are required")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid lat")
		return
	}
	lng, err := strconv.ParseFloat(lngStr, 64)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid lng")
		return
	}
	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid radius")
		return
	}

	if lat < -90 || lat > 90 || lng < -180 || lng > 180 {
		handler.RespondError(c, http.StatusBadRequest, "lat must be in [-90,90] and lng in [-180,180]")
		return
	}
	if radius <= 0 || radius > 50000 {
		handler.RespondError(c, http.StatusBadRequest, "radius must be in (0, 50000] meters")
		return
	}

	items, err := h.repo.FindNearby(c.Request.Context(), lat, lng, radius)
	if err != nil {
		handler.HandleServiceError(c, err)
		return
	}
	handler.RespondJSON(c, http.StatusOK, items)
}
