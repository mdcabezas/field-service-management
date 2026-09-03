package mobilehandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"localis-backend/internal/handler"
	"localis-backend/internal/model/mobile"
)

const maxPhotoSizeBytes = 5 * 1024 * 1024

type PhotoHandler struct {
	uploadDir    string
	thumbnailDir string
	pool         *pgxpool.Pool
}

func NewPhotoHandler(uploadDir string, pool *pgxpool.Pool) *PhotoHandler {
	thumbDir := filepath.Join(uploadDir, "thumbs")
	os.MkdirAll(thumbDir, 0755)
	return &PhotoHandler{
		uploadDir:    uploadDir,
		thumbnailDir: thumbDir,
		pool:         pool,
	}
}

func (h *PhotoHandler) Upload(c *gin.Context) {
	visitIDStr := c.PostForm("visit_id")
	visitID, err := uuid.Parse(visitIDStr)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid visit_id")
		return
	}

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "photo file required")
		return
	}
	defer file.Close()

	if header.Size > maxPhotoSizeBytes {
		handler.RespondError(c, http.StatusBadRequest, "file too large (max 5MB)")
		return
	}

	buf := make([]byte, 512)
	n, _ := file.Read(buf)
	_, _ = file.Seek(0, 0)
	mime := http.DetectContentType(buf[:n])
	if mime != "image/jpeg" && mime != "image/png" {
		handler.RespondError(c, http.StatusBadRequest, "only JPEG/PNG files accepted")
		return
	}

	photoID := uuid.New()
	fullPath := filepath.Join(h.uploadDir, photoID.String()+".jpg")
	thumbPath := filepath.Join(h.thumbnailDir, photoID.String()+"_thumb.jpg")

	out, err := os.Create(fullPath)
	if err != nil {
		handler.RespondError(c, http.StatusInternalServerError, "failed to save file")
		return
	}
	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		os.Remove(fullPath)
		handler.RespondError(c, http.StatusInternalServerError, "failed to write file")
		return
	}
	out.Close()

	if err := generateThumbnail(fullPath, thumbPath); err != nil {
		os.Remove(fullPath)
		handler.RespondError(c, http.StatusInternalServerError, "failed to generate thumbnail")
		return
	}

	timestamp := time.Now()
	if ts := c.PostForm("timestamp"); ts != "" {
		if parsed, err := time.Parse(time.RFC3339, ts); err == nil {
			timestamp = parsed
		}
	}

	var lat, lng *float64
	if latStr := c.PostForm("lat"); latStr != "" {
		if v, err := parseFloat(latStr); err == nil {
			lat = &v
		}
	}
	if lngStr := c.PostForm("lng"); lngStr != "" {
		if v, err := parseFloat(lngStr); err == nil {
			lng = &v
		}
	}

	url := fmt.Sprintf("/api/photos/%s/file", photoID.String())
	thumbURL := fmt.Sprintf("/api/photos/%s/thumb", photoID.String())

	if h.pool != nil {
		ctx := c.Request.Context()
		if lat != nil && lng != nil {
			_, err = h.pool.Exec(ctx,
				`INSERT INTO operations.visit_photos (id, visit_id, geom, timestamp)
				 VALUES ($1, $2, ST_SetSRID(ST_MakePoint($3, $4), 4326), $5)
				 ON CONFLICT (id) DO NOTHING`,
				photoID, visitID, *lng, *lat, timestamp,
			)
		} else {
			_, err = h.pool.Exec(ctx,
				`INSERT INTO operations.visit_photos (id, visit_id, timestamp)
				 VALUES ($1, $2, $3)
				 ON CONFLICT (id) DO NOTHING`,
				photoID, visitID, timestamp,
			)
		}
		if err != nil {
			os.Remove(fullPath)
			os.Remove(thumbPath)
			handler.RespondError(c, http.StatusInternalServerError, "failed to save photo record")
			return
		}
	}

	photo := mobile.MobilePhoto{
		ID:            photoID,
		VisitID:       visitID,
		URL:           url,
		ThumbnailURL:  thumbURL,
		OriginalPath:  fullPath,
		ThumbnailPath: thumbPath,
		Lat:           lat,
		Lng:           lng,
		Timestamp:     timestamp,
		CreatedAt:     time.Now(),
	}

	handler.RespondJSON(c, http.StatusCreated, photo)
}

func (h *PhotoHandler) ServeFile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	pattern := filepath.Join(h.uploadDir, id.String()+".*")
	matches, _ := filepath.Glob(pattern)
	if len(matches) == 0 {
		handler.RespondError(c, http.StatusNotFound, "photo not found")
		return
	}

	c.File(matches[0])
}

func (h *PhotoHandler) ServeThumb(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handler.RespondError(c, http.StatusBadRequest, "invalid id")
		return
	}

	thumbPath := filepath.Join(h.thumbnailDir, id.String()+"_thumb.jpg")
	if _, err := os.Stat(thumbPath); os.IsNotExist(err) {
		handler.RespondError(c, http.StatusNotFound, "thumbnail not found")
		return
	}

	c.File(thumbPath)
}

func generateThumbnail(srcPath, dstPath string) error {
	img, err := imaging.Open(srcPath)
	if err != nil {
		return err
	}
	thumb := imaging.Resize(img, 150, 0, imaging.Lanczos)
	return imaging.Save(thumb, dstPath, imaging.JPEGQuality(75))
}

func parseFloat(s string) (float64, error) {
	var v float64
	_, err := fmt.Sscanf(s, "%f", &v)
	return v, err
}
