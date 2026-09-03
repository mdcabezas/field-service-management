package mobilehandler

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func createTestJPEG(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, nil); err != nil {
		t.Fatalf("failed to create test JPEG: %v", err)
	}
	return buf.Bytes()
}

func TestPhotoHandler_Upload(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	jpegData := createTestJPEG(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("visit_id", "550e8400-e29b-41d4-a716-446655440000")

	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write(jpegData)

	writer.Close()

	r := gin.New()
	r.POST("/api/mobile/photos/upload", h.Upload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/photos/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body=%s", w.Code, w.Body.String())
	}
}

func TestPhotoHandler_Upload_InvalidVisitID(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	jpegData := createTestJPEG(t)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("visit_id", "invalid-uuid")
	part, _ := writer.CreateFormFile("photo", "test.jpg")
	part.Write(jpegData)
	writer.Close()

	r := gin.New()
	r.POST("/api/mobile/photos/upload", h.Upload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/photos/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_Upload_NoFile(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("visit_id", "550e8400-e29b-41d4-a716-446655440000")
	writer.Close()

	r := gin.New()
	r.POST("/api/mobile/photos/upload", h.Upload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/photos/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_Upload_InvalidFileType(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("visit_id", "550e8400-e29b-41d4-a716-446655440000")
	part, _ := writer.CreateFormFile("photo", "test.exe")
	part.Write([]byte{0x4D, 0x5A, 0x90, 0x00, 0x03, 0x00, 0x00, 0x00})
	writer.Close()

	r := gin.New()
	r.POST("/api/mobile/photos/upload", h.Upload)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/mobile/photos/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeFile_NotFound(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	id := "550e8400-e29b-41d4-a716-446655440000"
	r := gin.New()
	r.GET("/api/mobile/photos/:id/file", h.ServeFile)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/"+id+"/file", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeFile_InvalidID(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	r := gin.New()
	r.GET("/api/mobile/photos/:id/file", h.ServeFile)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/invalid/file", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeThumb_NotFound(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	id := "550e8400-e29b-41d4-a716-446655440000"
	r := gin.New()
	r.GET("/api/mobile/photos/:id/thumb", h.ServeThumb)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/"+id+"/thumb", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeThumb_InvalidID(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	r := gin.New()
	r.GET("/api/mobile/photos/:id/thumb", h.ServeThumb)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/invalid/thumb", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeFile(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	id := "550e8400-e29b-41d4-a716-446655440000"
	filePath := filepath.Join(dir, id+".jpg")
	os.WriteFile(filePath, []byte("fake image"), 0644)

	r := gin.New()
	r.GET("/api/mobile/photos/:id/file", h.ServeFile)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/"+id+"/file", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestPhotoHandler_ServeThumb(t *testing.T) {
	dir := t.TempDir()
	h := NewPhotoHandler(dir, nil)

	id := "550e8400-e29b-41d4-a716-446655440000"
	thumbPath := filepath.Join(h.thumbnailDir, id+"_thumb.jpg")
	os.WriteFile(thumbPath, []byte("fake thumbnail"), 0644)

	r := gin.New()
	r.GET("/api/mobile/photos/:id/thumb", h.ServeThumb)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/mobile/photos/"+id+"/thumb", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMIMEDetection(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected string
	}{
		{"jpeg", []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10}, "image/jpeg"},
		{"png", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, "image/png"},
		{"text", []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F}, "text/plain; charset=utf-8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := make([]byte, 512)
			n := copy(buf, tt.data)
			got := http.DetectContentType(buf[:n])
			if got != tt.expected {
				t.Errorf("DetectContentType(%q) = %v, want %v", tt.name, got, tt.expected)
			}
		})
	}
}
