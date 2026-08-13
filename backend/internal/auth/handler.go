package auth

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	EmployeeNumber string `json:"employee_number" binding:"required"`
	Password       string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Handler struct {
	ldapAuth   *LDAPAuth
	jwtAuth    *JWTAuth
	revocation *RevocationService
}

func NewHandler(ldapAuth *LDAPAuth, jwtAuth *JWTAuth, revocation *RevocationService) *Handler {
	return &Handler{
		ldapAuth:   ldapAuth,
		jwtAuth:    jwtAuth,
		revocation: revocation,
	}
}

const accessTokenTTLSeconds = 3600

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee_number and password required"})
		return
	}

	user, err := h.ldapAuth.Bind(req.EmployeeNumber, req.Password)
	if err != nil {
		slog.Warn("login failed", "error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	role := h.ldapAuth.ResolveRole(user)

	accessToken, err := h.jwtAuth.GenerateToken(user.EmployeeNumber, role, TokenTypeAccess, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.jwtAuth.GenerateToken(user.EmployeeNumber, role, TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    accessTokenTTLSeconds,
		TokenType:    "Bearer",
	})
}

func (h *Handler) Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	claims, err := h.jwtAuth.ValidateToken(req.RefreshToken, TokenTypeRefresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if h.revocation != nil && claims.ID != "" && h.revocation.IsRevoked(claims.ID) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked"})
		return
	}

	if h.revocation != nil && claims.ID != "" && claims.ExpiresAt != nil {
		h.revocation.Revoke(claims.ID, claims.ExpiresAt.Time)
	}

	accessToken, err := h.jwtAuth.GenerateToken(claims.EmployeeNumber, claims.Role, TokenTypeAccess, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	newRefreshToken, err := h.jwtAuth.GenerateToken(claims.EmployeeNumber, claims.Role, TokenTypeRefresh, 7*24*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    accessTokenTTLSeconds,
		TokenType:    "Bearer",
	})
}

func (h *Handler) Me(c *gin.Context) {
	claims, exists := c.Get(ClaimsKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	cl, ok := claims.(*Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"employee_number": cl.EmployeeNumber,
		"role":            cl.Role,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "refresh_token required"})
		return
	}

	claims, err := h.jwtAuth.ValidateToken(req.RefreshToken, TokenTypeRefresh)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if h.revocation != nil && claims.ID != "" {
		if claims.ExpiresAt != nil {
			h.revocation.Revoke(claims.ID, claims.ExpiresAt.Time)
		}
	}

	if authHeader := c.GetHeader("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")
		if accessClaims, err := h.jwtAuth.ValidateToken(accessToken, TokenTypeAccess); err == nil {
			if h.revocation != nil && accessClaims.ID != "" && accessClaims.ExpiresAt != nil {
				h.revocation.Revoke(accessClaims.ID, accessClaims.ExpiresAt.Time)
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
