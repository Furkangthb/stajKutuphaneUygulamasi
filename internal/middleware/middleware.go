package middleware

import (
	"net/http"
	"strings"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

func RequireAuth(authService *services.AuthServices) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Hata": "Yetkisiz giris",
			})
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := authService.VerifyToken(c.Request.Context(), token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"Hata": "Yetkisiz giris",
			})
			return
		}
		c.Set("user_id", userID)
		c.Set("role", role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		for _, r := range allowedRoles {
			if userRole == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"Hata": "Yetkiniz yok",
		})
	}
}
