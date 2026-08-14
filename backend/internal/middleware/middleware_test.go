package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func testEngine(mws ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(mws...)
	r.GET("/x", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
	return r
}

func TestCORS_AllowedOrigin(t *testing.T) {
	r := testEngine(CORSWithOrigins([]string{"http://localhost:3000"}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("allow-origin = %q", got)
	}
}

func TestCORS_DisallowedOrigin(t *testing.T) {
	r := testEngine(CORSWithOrigins([]string{"http://localhost:3000"}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Origin", "http://evil.example")
	r.ServeHTTP(w, req)
	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("allow-origin = %q, want empty", got)
	}
}

func TestCORS_PreflightAllowed(t *testing.T) {
	r := testEngine(CORSWithOrigins([]string{"http://localhost:3000"}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNoContent)
	}
}

func TestCORS_PreflightForbidden(t *testing.T) {
	r := testEngine(CORSWithOrigins([]string{"http://localhost:3000"}))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/x", nil)
	req.Header.Set("Origin", "http://evil.example")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

func TestSecurityHeaders(t *testing.T) {
	r := testEngine(SecurityHeaders())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	for _, h := range []string{"X-Content-Type-Options", "X-Frame-Options", "Content-Security-Policy", "Strict-Transport-Security"} {
		if w.Header().Get(h) == "" {
			t.Fatalf("missing header %s", h)
		}
	}
}

func TestRequestID_Generated(t *testing.T) {
	r := testEngine(RequestID())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if got := w.Header().Get("X-Request-ID"); got == "" {
		t.Fatal("expected generated X-Request-ID")
	}
}

func TestRequestID_Passthrough(t *testing.T) {
	r := testEngine(RequestID())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("X-Request-ID", "abc-123")
	r.ServeHTTP(w, req)
	if got := w.Header().Get("X-Request-ID"); got != "abc-123" {
		t.Fatalf("X-Request-ID = %q, want abc-123", got)
	}
}

func TestBodyLimit(t *testing.T) {
	r := testEngine(BodyLimit(10))
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRateLimiter_AllowAndBlock(t *testing.T) {
	rl := NewRateLimiter(2, time.Hour)
	defer rl.Stop()

	for i := 0; i < 2; i++ {
		if !rl.Allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}
	if rl.Allow("1.2.3.4") {
		t.Fatal("third request should be blocked")
	}
	if !rl.Allow("9.9.9.9") {
		t.Fatal("different IP should be allowed")
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	rl := NewRateLimiter(1, time.Hour)
	defer rl.Stop()
	rl.Allow("1.2.3.4")
	rl.visitors["1.2.3.4"].lastSeen = time.Now().Add(-2 * time.Hour)
	if !rl.Allow("1.2.3.4") {
		t.Fatal("request after window should be allowed")
	}
}

func TestRateLimiter_Middleware(t *testing.T) {
	rl := NewRateLimiter(1, time.Hour)
	defer rl.Stop()
	r := testEngine(rl.Middleware())

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("first status = %d", w.Code)
	}

	w2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("second status = %d, want 429", w2.Code)
	}
	if w2.Header().Get("Retry-After") == "" {
		t.Fatal("expected Retry-After header")
	}
}

func TestLogging(t *testing.T) {
	r := testEngine(RequestID(), Logging())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}
