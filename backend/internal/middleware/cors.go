package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ",")
	if len(allowedOrigins) == 0 || (len(allowedOrigins) == 1 && allowedOrigins[0] == "") {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	return CORSWithOrigins(allowedOrigins)
}

func CORSWithOrigins(origins []string) gin.HandlerFunc {
	originSet := make(map[string]bool, len(origins))
	for _, o := range origins {
		originSet[strings.TrimSpace(o)] = true
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := origin != "" && originSet[origin]

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Range, X-Total-Count, X-Request-ID")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}
		c.Header("Vary", "Origin")

		if c.Request.Method == http.MethodOptions {
			if !allowed {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
