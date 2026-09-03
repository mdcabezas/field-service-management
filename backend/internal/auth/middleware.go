package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"localis-backend/internal/service"
)

const ClaimsKey = "claims"

type Middleware struct {
	jwtAuth    *JWTAuth
	revocation *RevocationService
}

func NewMiddleware(jwtAuth *JWTAuth, revocation *RevocationService) *Middleware {
	return &Middleware{
		jwtAuth:    jwtAuth,
		revocation: revocation,
	}
}

func (m *Middleware) Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		claims, err := m.jwtAuth.ValidateToken(parts[1], TokenTypeAccess)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if m.revocation != nil && claims.ID != "" && m.revocation.IsRevoked(claims.ID) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token revoked"})
			return
		}

		c.Set(ClaimsKey, claims)
		c.Request = c.Request.WithContext(
			service.WithUserID(c.Request.Context(), claims.UserID),
		)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	roleSet := make(map[string]bool, len(roles))
	for _, r := range roles {
		roleSet[r] = true
	}
	return func(c *gin.Context) {
		claims, exists := c.Get(ClaimsKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		cl, ok := claims.(*Claims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if !roleSet[cl.Role] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
