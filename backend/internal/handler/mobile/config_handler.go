package mobilehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

type ConfigHandler struct {
	config mobile.MobileConfig
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{
		config: mobile.MobileConfig{
			PhotoRetentionDays:      7,
			PhotoMaxStorageMB:       100,
			PhotoWarningThresholdMB: 50,
		},
	}
}

func (h *ConfigHandler) GetConfig(c *gin.Context) {
	handler.RespondJSON(c, http.StatusOK, h.config)
}

func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	var body struct {
		PhotoRetentionDays      *int `json:"photo_retention_days"`
		PhotoMaxStorageMB       *int `json:"photo_max_storage_mb"`
		PhotoWarningThresholdMB *int `json:"photo_warning_threshold_mb"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.PhotoRetentionDays != nil {
		h.config.PhotoRetentionDays = *body.PhotoRetentionDays
	}
	if body.PhotoMaxStorageMB != nil {
		h.config.PhotoMaxStorageMB = *body.PhotoMaxStorageMB
	}
	if body.PhotoWarningThresholdMB != nil {
		h.config.PhotoWarningThresholdMB = *body.PhotoWarningThresholdMB
	}

	handler.RespondJSON(c, http.StatusOK, h.config)
}
