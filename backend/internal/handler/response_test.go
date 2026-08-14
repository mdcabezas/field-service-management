package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/service"
)

func TestParseUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x/:id", func(c *gin.Context) {
		id, ok := ParseUUID(c, "id")
		if !ok {
			c.String(http.StatusOK, "invalid")
			return
		}
		c.String(http.StatusOK, id.String())
	})

	cases := []struct {
		name  string
		param string
		valid bool
	}{
		{name: "valid uuid", param: "53000000-0000-0000-0000-000000000001", valid: true},
		{name: "invalid", param: "nope", valid: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/x/"+tc.param, nil)
			r.ServeHTTP(w, req)
			if tc.valid && w.Body.String() == "invalid" {
				t.Fatal("expected valid uuid, got invalid")
			}
			if !tc.valid && w.Body.String() != "invalid" {
				t.Fatalf("expected invalid, got %q", w.Body.String())
			}
		})
	}
}

func TestParsePagination(t *testing.T) {
	r := gin.New()
	r.GET("/x", func(c *gin.Context) {
		limit, offset := ParsePagination(c)
		c.JSON(http.StatusOK, map[string]int{"limit": limit, "offset": offset})
	})

	cases := []struct {
		name       string
		query      string
		wantLimit  int
		wantOffset int
	}{
		{name: "defaults", query: "", wantLimit: 20, wantOffset: 0},
		{name: "custom", query: "?limit=50&offset=10", wantLimit: 50, wantOffset: 10},
		{name: "limit capped", query: "?limit=9999", wantLimit: 100, wantOffset: 0},
		{name: "offset capped", query: "?offset=999999", wantLimit: 20, wantOffset: 10000},
		{name: "invalid ignored", query: "?limit=abc&offset=-5", wantLimit: 20, wantOffset: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/x"+tc.query, nil)
			r.ServeHTTP(w, req)
			body := w.Body.String()
			if !strings.Contains(body, `"limit":`+itoa(tc.wantLimit)) {
				t.Fatalf("limit not %d in %s", tc.wantLimit, body)
			}
			if !strings.Contains(body, `"offset":`+itoa(tc.wantOffset)) {
				t.Fatalf("offset not %d in %s", tc.wantOffset, body)
			}
		})
	}
}

func itoa(i int) string {
	return strconv.Itoa(i)
}

func TestHandleServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := []struct {
		name   string
		err    error
		status int
	}{
		{name: "not found", err: &service.NotFoundError{Resource: "v", ID: "1"}, status: http.StatusNotFound},
		{name: "validation", err: &service.ValidationError{Field: "f", Message: "m"}, status: http.StatusBadRequest},
		{name: "conflict", err: &service.ConflictError{Message: "c"}, status: http.StatusConflict},
		{name: "transition", err: &service.TransitionError{Entity: "v", From: "a", To: "b"}, status: http.StatusConflict},
		{name: "generic", err: errors.New("boom"), status: http.StatusInternalServerError},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.GET("/x", func(c *gin.Context) {
				HandleServiceError(c, tc.err)
			})
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			r.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d", w.Code, tc.status)
			}
		})
	}
}

func TestRespondError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", func(c *gin.Context) {
		RespondError(c, http.StatusBadRequest, "bad input")
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "bad input") {
		t.Fatalf("body = %q", w.Body.String())
	}
}
