package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/auth"
)

func testRouter(t *testing.T) *gin.Engine {
	t.Helper()
	a, err := auth.NewJWTAuth(strings.Repeat("x", 32))
	if err != nil {
		t.Fatalf("NewJWTAuth: %v", err)
	}
	// pool nil is safe: no DB-backed route is exercised in these tests
	r, cleanup := NewRouter(a, nil, "/tmp/test-photos")
	t.Cleanup(cleanup)
	return r
}

func TestRouter_Health(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("body = %v", body)
	}
}

func TestRouter_Login_InvalidBody(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouter_Login_BadRequest(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(`{"email":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouter_Refresh_InvalidBody(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouter_Logout_NoHeader(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/logout", strings.NewReader(`{"refresh_token":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRouter_API_RequiresAuth(t *testing.T) {
	r := testRouter(t)
	routes := []string{
		"/api/users",
		"/api/partners",
		"/api/customers",
		"/api/vehicles",
		"/api/visits",
		"/api/tech-roles",
	}
	for _, path := range routes {
		t.Run(path, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, path, nil)
			r.ServeHTTP(w, req)
			if w.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want 401", w.Code)
			}
		})
	}
}

func TestRouter_Me_RequiresAuth(t *testing.T) {
	r := testRouter(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}
