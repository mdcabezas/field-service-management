package mobilehandler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

type GPSHandler struct {
	waivers map[string][]mobile.GPSWaiver
}

func NewGPSHandler() *GPSHandler {
	return &GPSHandler{
		waivers: make(map[string][]mobile.GPSWaiver),
	}
}

func (h *GPSHandler) CreateWaiver(c *gin.Context) {
	var body struct {
		VisitID   uuid.UUID `json:"visit_id" binding:"required"`
		Type      string    `json:"type" binding:"required"`
		Timestamp time.Time `json:"timestamp" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	validTypes := map[string]bool{"gps_timeout": true, "no_fix": true, "low_accuracy": true}
	if !validTypes[body.Type] {
		handler.RespondError(c, http.StatusBadRequest, "invalid waiver type")
		return
	}

	userName := c.GetString("user_name")
	if userName == "" {
		userName = "Unknown"
	}

	waiver := mobile.GPSWaiver{
		ID:         uuid.New(),
		VisitID:    body.VisitID,
		Type:       body.Type,
		AcceptedBy: userName,
		Timestamp:  body.Timestamp,
		CreatedAt:  time.Now(),
	}

	h.waivers[body.VisitID.String()] = append(h.waivers[body.VisitID.String()], waiver)

	handler.RespondJSON(c, http.StatusCreated, waiver)
}

func (h *GPSHandler) ListWaivers(c *gin.Context) {
	visitIDStr := c.Query("visit_id")
	if visitIDStr == "" {
		handler.RespondError(c, http.StatusBadRequest, "visit_id required")
		return
	}

	visitID, err := uuid.Parse(visitIDStr)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid visit_id")
		return
	}

	waivers := h.waivers[visitID.String()]
	if waivers == nil {
		waivers = []mobile.GPSWaiver{}
	}

	handler.RespondJSON(c, http.StatusOK, gin.H{"waivers": waivers})
}

func (h *GPSHandler) DeleteWaiver(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	for visitID, waivers := range h.waivers {
		for i, w := range waivers {
			if w.ID == id {
				h.waivers[visitID] = append(waivers[:i], waivers[i+1:]...)
				c.Status(http.StatusNoContent)
				return
			}
		}
	}

	handler.RespondError(c, http.StatusNotFound, "waiver not found")
}
