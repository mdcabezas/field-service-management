package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const testSecret = "this-is-a-test-secret-that-is-long-enough-32-bytes"

func mustJWT(t *testing.T, secret string) *JWTAuth {
	t.Helper()
	a, err := NewJWTAuth(secret)
	if err != nil {
		t.Fatalf("NewJWTAuth: %v", err)
	}
	return a
}

func TestNewJWTAuth_EmptySecret(t *testing.T) {
	if _, err := NewJWTAuth(""); err != ErrJWTSecretEmpty {
		t.Fatalf("err = %v, want ErrJWTSecretEmpty", err)
	}
}

func TestNewJWTAuth_TooShort(t *testing.T) {
	if _, err := NewJWTAuth("short"); err == nil {
		t.Fatal("expected error for short secret")
	}
}

func TestJWTAuth_GenerateAndValidate(t *testing.T) {
	a := mustJWT(t, testSecret)
	tok, err := a.GenerateToken("1001", "admin", TokenTypeAccess, time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}

	claims, err := a.ValidateToken(tok, TokenTypeAccess)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.EmployeeNumber != "1001" || claims.Role != "admin" {
		t.Fatalf("claims = %+v", claims)
	}
}

func TestJWTAuth_ValidateToken_WrongType(t *testing.T) {
	a := mustJWT(t, testSecret)
	tok, _ := a.GenerateToken("1001", "admin", TokenTypeRefresh, time.Hour)
	if _, err := a.ValidateToken(tok, TokenTypeAccess); err != ErrWrongTokenType {
		t.Fatalf("err = %v, want ErrWrongTokenType", err)
	}
}

func TestJWTAuth_ValidateToken_Garbage(t *testing.T) {
	a := mustJWT(t, testSecret)
	if _, err := a.ValidateToken("not.a.token", TokenTypeAccess); err == nil {
		t.Fatal("expected error for garbage token")
	}
}

func TestJWTAuth_ValidateToken_WrongSecret(t *testing.T) {
	a := mustJWT(t, testSecret)
	other := mustJWT(t, strings.Repeat("x", 32))
	tok, _ := a.GenerateToken("1001", "admin", TokenTypeAccess, time.Hour)
	if _, err := other.ValidateToken(tok, TokenTypeAccess); err == nil {
		t.Fatal("expected error for wrong secret")
	}
}

func TestJWTAuth_ValidateToken_Expired(t *testing.T) {
	a := mustJWT(t, testSecret)
	tok, _ := a.GenerateToken("1001", "admin", TokenTypeAccess, -time.Minute)
	if _, err := a.ValidateToken(tok, TokenTypeAccess); err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestRevocationService(t *testing.T) {
	r := NewRevocationService()
	defer r.Stop()

	if r.IsRevoked("jti-1") {
		t.Fatal("should not be revoked initially")
	}
	r.Revoke("jti-1", time.Now().Add(time.Hour))
	if !r.IsRevoked("jti-1") {
		t.Fatal("should be revoked")
	}
}

func TestMiddleware_Validate_NoHeader(t *testing.T) {
	a := mustJWT(t, testSecret)
	m := NewMiddleware(a, nil)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", m.Validate(), func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMiddleware_Validate_BadFormat(t *testing.T) {
	a := mustJWT(t, testSecret)
	m := NewMiddleware(a, nil)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", m.Validate(), func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Basic abc")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMiddleware_Validate_InvalidToken(t *testing.T) {
	a := mustJWT(t, testSecret)
	m := NewMiddleware(a, nil)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", m.Validate(), func(c *gin.Context) { c.Status(http.StatusOK) })
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer bogus")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMiddleware_Validate_ValidToken(t *testing.T) {
	a := mustJWT(t, testSecret)
	m := NewMiddleware(a, nil)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", m.Validate(), func(c *gin.Context) {
		claims, ok := c.Get(ClaimsKey)
		if !ok {
			t.Error("claims not set")
			return
		}
		cl := claims.(*Claims)
		if cl.EmployeeNumber != "1001" {
			t.Errorf("employee = %q", cl.EmployeeNumber)
		}
		c.Status(http.StatusOK)
	})

	tok, _ := a.GenerateToken("1001", "admin", TokenTypeAccess, time.Hour)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestMiddleware_Validate_Revoked(t *testing.T) {
	a := mustJWT(t, testSecret)
	rev := NewRevocationService()
	defer rev.Stop()
	m := NewMiddleware(a, rev)
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/x", m.Validate(), func(c *gin.Context) { c.Status(http.StatusOK) })

	tok, _ := a.GenerateToken("1001", "admin", TokenTypeAccess, time.Hour)
	claims, _ := a.ValidateToken(tok, TokenTypeAccess)
	rev.Revoke(claims.ID, time.Now().Add(time.Hour))

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set("Authorization", "Bearer "+tok)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestRequireRole(t *testing.T) {
	a := mustJWT(t, testSecret)
	m := NewMiddleware(a, nil)
	gin.SetMode(gin.TestMode)

	makeTok := func(role string) string {
		tok, _ := a.GenerateToken("1001", role, TokenTypeAccess, time.Hour)
		return tok
	}

	cases := []struct {
		name   string
		role   string
		status int
	}{
		{name: "admin allowed", role: "admin", status: http.StatusOK},
		{name: "operator forbidden", role: "operator", status: http.StatusForbidden},
		{name: "no claims unauthorized", role: "", status: http.StatusUnauthorized},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			adminOnly := RequireRole("admin")
			if tc.role == "" {
				r.GET("/x", adminOnly, func(c *gin.Context) { c.Status(http.StatusOK) })
			} else {
				r.GET("/x", m.Validate(), adminOnly, func(c *gin.Context) { c.Status(http.StatusOK) })
			}
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			if tc.role != "" {
				req.Header.Set("Authorization", "Bearer "+makeTok(tc.role))
			}
			r.ServeHTTP(w, req)
			if w.Code != tc.status {
				t.Fatalf("status = %d, want %d", w.Code, tc.status)
			}
		})
	}
}

func TestHandler_Login_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.POST("/login", h.Login)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Login_LDAPUnreachable(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	// ldaps:// to a closed port: Bind dial fails -> 401 invalid credentials
	bad, err := NewLDAPAuth("ldaps://127.0.0.1:1", "dc=workflows,dc=cl", "pw")
	if err != nil {
		t.Fatalf("NewLDAPAuth: %v", err)
	}
	h := NewHandler(bad, a, nil)
	r := gin.New()
	r.POST("/login", h.Login)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(`{"employee_number":"1001","password":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Refresh_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.POST("/refresh", h.Refresh)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Refresh_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	rev := NewRevocationService()
	defer rev.Stop()
	h := NewHandler(nil, a, rev)
	r := gin.New()
	r.POST("/refresh", h.Refresh)

	rt, _ := a.GenerateToken("1001", "admin", TokenTypeRefresh, 24*time.Hour)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{"refresh_token":"`+rt+`"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.AccessToken == "" || resp.RefreshToken == "" {
		t.Fatal("expected tokens in response")
	}
}

func TestHandler_Refresh_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.POST("/refresh", h.Refresh)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/refresh", strings.NewReader(`{"refresh_token":"bogus"}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Logout_InvalidBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.POST("/logout", h.Logout)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Logout_Valid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	rev := NewRevocationService()
	defer rev.Stop()
	h := NewHandler(nil, a, rev)
	r := gin.New()
	r.POST("/logout", h.Logout)

	rt, _ := a.GenerateToken("1001", "admin", TokenTypeRefresh, 24*time.Hour)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/logout", strings.NewReader(`{"refresh_token":"`+rt+`"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer bogus-access")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Me_Unauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.GET("/me", h.Me)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestHandler_Me_Authorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	a := mustJWT(t, testSecret)
	h := NewHandler(nil, a, nil)
	r := gin.New()
	r.GET("/me", func(c *gin.Context) {
		c.Set(ClaimsKey, &Claims{EmployeeNumber: "1001", Role: "admin"})
		h.Me(c)
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/me", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), `"employee_number":"1001"`) {
		t.Fatalf("body = %s", w.Body.String())
	}
}

func TestLDAPAuth_ResolveRole(t *testing.T) {
	cases := []struct {
		groups []string
		role   string
	}{
		{groups: []string{"Admins"}, role: "admin"},
		{groups: []string{"managers"}, role: "manager"},
		{groups: []string{"supervisor"}, role: "supervisor"},
		{groups: []string{"ops"}, role: "technician"},
		{groups: nil, role: "technician"},
	}
	for _, tc := range cases {
		got := (&LDAPAuth{}).ResolveRole(&LDAPUser{Groups: tc.groups})
		if got != tc.role {
			t.Errorf("ResolveRole(%v) = %q, want %q", tc.groups, got, tc.role)
		}
	}
}

func TestEscapeDN(t *testing.T) {
	if got := escapeDN("a,b=c\\d"); got != `a\,b\=c\\d` {
		t.Fatalf("escapeDN = %q", got)
	}
	if got := escapeDN(" plain"); got != `\20plain` {
		t.Fatalf("escapeDN leading space = %q", got)
	}
}

func TestNewLDAPAuth_LDAPSNoDial(t *testing.T) {
	a, err := NewLDAPAuth("ldaps://ldap.example:636", "dc=workflows,dc=cl", "pw")
	if err != nil {
		t.Fatalf("NewLDAPAuth: %v", err)
	}
	if a == nil || a.useStartTLS {
		t.Fatal("ldaps should not use StartTLS")
	}
}

func TestLDAPAuth_Bind_DialFail(t *testing.T) {
	a := &LDAPAuth{url: "ldap://127.0.0.1:1", baseDN: "dc=workflows,dc=cl"}
	if _, err := a.Bind("1001", "pw"); err == nil {
		t.Fatal("expected dial error")
	}
}

func TestLDAPAuth_Bind_NonNumericEmployeeNumber(t *testing.T) {
	// Dial fails before numeric check, so use ldaps URL with no dial
	a := &LDAPAuth{url: "ldaps://127.0.0.1:1", baseDN: "dc=workflows,dc=cl", useStartTLS: false}
	if _, err := a.Bind("abc", "pw"); err == nil {
		t.Fatal("expected error for non-numeric employee number")
	}
}
