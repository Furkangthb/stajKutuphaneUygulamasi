package handlers

import (
	"net/http"
	"strings"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

type AuthHandlers struct {
	authService *services.AuthServices
}

func NewAuthHandlers(authService *services.AuthServices) *AuthHandlers {
	return &AuthHandlers{authService: authService}
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Login godoc
// @Summary Kullanıcı Giriş Yap
// @Description Sisteme giriş yapmak için kullanılır
// @Tags Auth
// @Accept  json
// @Produce json
// @Param request body LoginRequest true "Kullanıcı Giriş Bilgileri"
// @Success 200 {object}  map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// Failure 401 {object} map[string]interface {}
// @Router /api/auth/login  [post]
func (h *AuthHandlers) Login(c *gin.Context) {
	var istek loginRequest

	if err := c.ShouldBindJSON(&istek); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	token, err := h.authService.Login(c.Request.Context(), istek.Email, istek.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Token": token,
	})
}

// Logout godoc
// @Summary Kullanıcı Çıkış Yap
// @Description Sistemden çıkış yapmak için kullanılır
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/auth/logout [post]
func (h *AuthHandlers) Logout(c *gin.Context) {
	autHeader := c.GetHeader("Authorization")
	if autHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Hata": "token bulunamadı",
		})
		return
	}
	parts := strings.SplitN(autHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"Hata": "gecersiz authorization header",
		})
		return
	}
	token := parts[1]
	if err := h.authService.Logout(c.Request.Context(), token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Mesaj": "Cıkıs basarali",
	})
}
