package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"localis-backend/internal/testutil"
)

func TestNewRouter(t *testing.T) {
	pool := testutil.TestPool(t)
	defer pool.Close()

	jwtAuth := testutil.NewMockJWTAuth(t)

	router, cleanup := NewRouter(jwtAuth, pool, "/tmp/test-photos")
	defer cleanup()

	require.NotNil(t, router)

	// Test health endpoint
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "ok")
}

func TestRouter_WithNilPool(t *testing.T) {
	jwtAuth := testutil.NewMockJWTAuth(t)

	router, cleanup := NewRouter(jwtAuth, nil, "/tmp/test-photos")
	defer cleanup()

	require.NotNil(t, router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}

func TestRouter_AuthRoutes(t *testing.T) {
	jwtAuth := testutil.NewMockJWTAuth(t)

	router, cleanup := NewRouter(jwtAuth, nil, "/tmp/test-photos")
	defer cleanup()

	// Test login endpoint exists (400 for missing body is expected)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	router.ServeHTTP(w, req)

	require.NotEqual(t, http.StatusNotFound, w.Code)
}
