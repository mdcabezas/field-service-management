package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"localis-backend/internal/service"
)

const maxPaginationLimit = 100
const maxPaginationOffset = 10000

func RespondJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}

func RespondError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"error": message})
}

func HandleServiceError(c *gin.Context, err error) {
	var notFound *service.NotFoundError
	var validation *service.ValidationError
	var conflict *service.ConflictError
	var transition *service.TransitionError

	switch {
	case errors.As(err, &notFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	case errors.As(err, &validation):
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	case errors.As(err, &conflict):
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
	case errors.As(err, &transition):
		c.JSON(http.StatusConflict, gin.H{"error": "invalid status transition"})
	default:
		if err := c.Error(err).SetType(gin.ErrorTypePrivate); err != nil {
			slog.Warn("failed to set error type", "error", err)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func ParseUUID(c *gin.Context, paramName string) (uuid.UUID, bool) {
	param := c.Param(paramName)
	id, err := uuid.Parse(param)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}

func ParsePagination(c *gin.Context) (limit, offset int) {
	limit = 20
	offset = 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	if limit > maxPaginationLimit {
		limit = maxPaginationLimit
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	if offset > maxPaginationOffset {
		offset = maxPaginationOffset
	}

	return limit, offset
}
