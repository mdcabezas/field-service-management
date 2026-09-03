package mobilehandler

import (
	"bytes"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/auth"
)

func TestMobileIntegration_VisitFlow(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add auth middleware to set mock claims for tests
	r.Use(func(c *gin.Context) {
		claims := &auth.Claims{
			UserID: "1001",
			Role:           "admin",
		}
		c.Set(auth.ClaimsKey, claims)
		c.Next()
	})

	photoDir := t.TempDir()
	configHandler := NewConfigHandler()
	gpsHandler := NewGPSHandler()
	syncHandler := NewTestSyncHandler()
	photoHandler := NewPhotoHandler(photoDir, nil)

	r.GET("/api/mobile/config", configHandler.GetConfig)
	r.POST("/api/mobile/gps/waivers", gpsHandler.CreateWaiver)
	r.POST("/api/mobile/sync/push", syncHandler.Push)
	r.GET("/api/mobile/sync/pull", syncHandler.Pull)
	r.POST("/api/mobile/photos/upload", photoHandler.Upload)

	t.Run("Get Config", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/mobile/config", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
	})

	t.Run("Create GPS Waiver", func(t *testing.T) {
		body := `{"visit_id":"550e8400-e29b-41d4-a716-446655440000","type":"gps_timeout","timestamp":"2024-01-01T00:00:00Z"}`
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/gps/waivers", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("Upload Photo", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 10, 10))
		for y := 0; y < 10; y++ {
			for x := 0; x < 10; x++ {
				img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
			}
		}
		var imgBuf bytes.Buffer
		if err := jpeg.Encode(&imgBuf, img, nil); err != nil {
			t.Fatalf("failed to create test JPEG: %v", err)
		}

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("visit_id", "550e8400-e29b-41d4-a716-446655440000")

		part, _ := writer.CreateFormFile("photo", "test.jpg")
		part.Write(imgBuf.Bytes())
		writer.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/photos/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("Sync Push", func(t *testing.T) {
		syncBody := `{
			"checkpoints":[],
			"photos":[],
			"usages":[],
			"measurements":[],
			"reports":[],
			"status_changes":[]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(syncBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})

	t.Run("Sync Pull", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/mobile/sync/pull", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}
	})
}

func TestMobileIntegration_SyncWithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	syncHandler := NewTestSyncHandler()
	r.POST("/api/mobile/sync/push", syncHandler.Push)

	t.Run("Push with checkpoints", func(t *testing.T) {
		body := `{
			"checkpoints":[{"local_id":"local-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","type":"arrival","timestamp":"2024-01-01T00:00:00Z","lat":19.4326,"lng":-99.1332,"accuracy":10}],
			"photos":[],
			"usages":[],
			"measurements":[],
			"reports":[],
			"status_changes":[]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}
	})

	t.Run("Push with measurements", func(t *testing.T) {
		body := `{
			"checkpoints":[],
			"photos":[],
			"usages":[],
			"measurements":[{"local_id":"meas-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","key":"voltage","value":"220"}],
			"reports":[],
			"status_changes":[]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}
	})

	t.Run("Push with status changes", func(t *testing.T) {
		body := `{
			"checkpoints":[],
			"photos":[],
			"usages":[],
			"measurements":[],
			"reports":[],
			"status_changes":[{"local_id":"status-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","new_status":"in_progress"}]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}
	})
}

func TestMobileIntegration_ChecklistStageBlock(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	syncHandler := NewTestSyncHandler()
	r.POST("/api/mobile/sync/push", syncHandler.Push)

	t.Run("Status change blocked by incomplete checklist", func(t *testing.T) {
		body := `{
			"checkpoints":[],
			"photos":[],
			"usages":[],
			"measurements":[],
			"reports":[],
			"status_changes":[{"local_id":"status-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","new_status":"completed"}]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		results, ok := result["results"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected results object")
		}

		statusChanges, ok := results["status_changes"].([]interface{})
		if !ok {
			t.Fatalf("expected status_changes array")
		}

		if len(statusChanges) != 1 {
			t.Fatalf("expected 1 status change result")
		}
	})
}

func TestMobileIntegration_SyncConflictResolution(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Add auth middleware to set mock claims for tests
	r.Use(func(c *gin.Context) {
		claims := &auth.Claims{
			UserID: "1001",
			Role:           "admin",
		}
		c.Set(auth.ClaimsKey, claims)
		c.Next()
	})

	syncHandler := NewTestSyncHandler()
	r.POST("/api/mobile/sync/push", syncHandler.Push)
	r.GET("/api/mobile/sync/pull", syncHandler.Pull)

	t.Run("Push multiple item types", func(t *testing.T) {
		body := `{
			"checkpoints":[{"local_id":"cp-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","type":"arrival","timestamp":"2024-01-01T00:00:00Z","lat":19.4326,"lng":-99.1332,"accuracy":10}],
			"photos":[{"local_id":"photo-1","visit_id":"550e8400-e29b-41d4-a716-446655440000"}],
			"usages":[{"local_id":"usage-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","material_id":"550e8400-e29b-41d4-a716-446655440001","quantity":5}],
			"measurements":[{"local_id":"meas-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","key":"voltage","value":"220"}],
			"reports":[{"local_id":"report-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","content":"Test report"}],
			"status_changes":[{"local_id":"status-1","visit_id":"550e8400-e29b-41d4-a716-446655440000","new_status":"in_progress"}]
		}`

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/api/mobile/sync/push", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		results, ok := result["results"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected results object")
		}

		checkpoints, _ := results["checkpoints"].([]interface{})
		photos, _ := results["photos"].([]interface{})
		usages, _ := results["usages"].([]interface{})
		measurements, _ := results["measurements"].([]interface{})
		reports, _ := results["reports"].([]interface{})
		statusChanges, _ := results["status_changes"].([]interface{})

		if len(checkpoints) != 1 {
			t.Errorf("expected 1 checkpoint, got %d", len(checkpoints))
		}
		if len(photos) != 1 {
			t.Errorf("expected 1 photo, got %d", len(photos))
		}
		if len(usages) != 1 {
			t.Errorf("expected 1 usage, got %d", len(usages))
		}
		if len(measurements) != 1 {
			t.Errorf("expected 1 measurement, got %d", len(measurements))
		}
		if len(reports) != 1 {
			t.Errorf("expected 1 report, got %d", len(reports))
		}
		if len(statusChanges) != 1 {
			t.Errorf("expected 1 status change, got %d", len(statusChanges))
		}
	})

	t.Run("Pull returns visit data", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/api/mobile/sync/pull", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d", w.Code)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if _, ok := result["reference_data"]; !ok {
			t.Errorf("expected reference_data in response")
		}
		if _, ok := result["visits"]; !ok {
			t.Errorf("expected visits in response")
		}
		if _, ok := result["synced_at"]; !ok {
			t.Errorf("expected synced_at in response")
		}
	})
}