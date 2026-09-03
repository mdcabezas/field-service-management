package auth

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"localis-backend/internal/repository"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Handler struct {
	userRepo   repository.CoreUserRepository
	jwtAuth    *JWTAuth
	revocation *RevocationService
}

func NewHandler(userRepo repository.CoreUserRepository, jwtAuth *JWTAuth, revocation *RevocationService) *Handler {
	return &Handler{
		userRepo:   userRepo,
		jwtAuth:    jwtAuth,
		revocation: revocation,
	}
}

const accessTokenTTLSeconds = 3600

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password required"})
		return
	}

	userWithPassword, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		slog.Warn("login failed: user not found", "email", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if userWithPassword.PasswordHash == nil {
		slog.Warn("login failed: no password set", "email", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*userWithPassword.PasswordHash), []byte(req.Password)); err != nil {
		slog.Warn("login failed: wrong password", "email", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	userID := userWithPassword.ID.String()
	role := userWithPassword.Role

	accessToken, err := h.jwtAuth.GenerateToken(userID, role, TokenTypeAccess, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	refreshToken, err := h.jwtAuth.GenerateToken(userID, role, TokenTypeRefresh, 7*24*time.Hour)
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

	accessToken, err := h.jwtAuth.GenerateToken(claims.UserID, claims.Role, TokenTypeAccess, 1*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	newRefreshToken, err := h.jwtAuth.GenerateToken(claims.UserID, claims.Role, TokenTypeRefresh, 7*24*time.Hour)
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

	user, err := h.userRepo.GetByID(context.Background(), cl.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID.String(),
		"email": user.Email,
		"name":  user.Name,
		"role":  cl.Role,
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
