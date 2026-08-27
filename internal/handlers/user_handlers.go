package handlers

import (
	"net/http"
	"strconv"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	userService *services.UserServices
}

func NewUserHandlers(userServices *services.UserServices) *UserHandlers {
	return &UserHandlers{userService: userServices}
}

type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
}

func (h *UserHandlers) UserRegister(c *gin.Context) {
	var istek RegisterRequest
	if err := c.ShouldBindJSON(&istek); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	user, err := h.userService.UserRegister(ctx, istek.FirstName, istek.LastName, istek.Phone, istek.Email, istek.Password)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(201, user)
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *UserHandlers) UserLogin(c *gin.Context) {
	var istek LoginRequest
	if err := c.ShouldBindJSON(&istek); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	user, err := h.userService.UserLogin(ctx, istek.Email, istek.Password)
	if err != nil {
		c.JSON(401, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(200, user)
}

// UserDelete godoc
// @Summary Kullanıcı Silme
// @Description Kullanıcıyı siler. Bearer token ve Admin yetkisi gerektirir.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kullanıcı ID"
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users/{id} [delete]
func (h *UserHandlers) UserDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	err = h.userService.UserDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

// UserGet godoc
// @Summary Kullanıcı Bilgilerini Getir
// @Description Kullanıcı bilgilerini getirir. Bearer token gerektirir.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kullanıcı ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users/{id} [get]
func (h *UserHandlers) UserGet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	user, err := h.userService.UserGet(c.Request.Context(), id)
	if err != nil {
		c.JSON(404, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

type UserRework struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required" `
	Email     string `json:"email" binding:"required"`
	Role      string `json:"role" binding:"required"`
}

// UserUpdate godoc
// @Summary Kullanıcı Güncelleme
// @Description Kullanıcı bilgilerini günceller. Bearer token gerektirir.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kullanıcı ID"
// @Param request body UserRework true "Kullanıcı Güncellemesi"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users/{id} [put]
func (h *UserHandlers) UserUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	tokenUserID := c.GetInt("user_id")
	tokenRole := c.GetString("role")

	if tokenRole != "admin" && int64(tokenUserID) != id {
		c.JSON(http.StatusForbidden, gin.H{"Hata": "sadece kendi bilgilerinizi guncelleyebilirsiniz"})
		return
	}
	var rework UserRework
	if err := c.ShouldBindJSON(&rework); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	newRole := rework.Role
	if tokenRole != "admin" {
		currentUser, err := h.userService.UserGet(ctx, id)
		if err != nil {
			c.JSON(404, gin.H{
				"Hata": "kullanici bulunamadi",
			})
			return
		}
		newRole = currentUser.Role
	}
	user, err := h.userService.UserUpdate(ctx, int(id), rework.FirstName, rework.LastName, rework.Email, newRole)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UserList godoc
// @Summary Kullanıcıları Listeleme
// @Description Tüm kullanıcıları listeler. Bearer token ve Admin yetkisi gerektirir.
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Sayfa numarası"
// @Param page_size query int false "Sayfa boyutu"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/users [get]
func (h *UserHandlers) UserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	users, err := h.userService.UserList(c.Request.Context(), page, page_size)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"data": users,
		"page": page,
	})
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

func (h *UserHandlers) UserChangePassword(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	tokenUserID := c.GetInt("user_id")
	if int64(tokenUserID) != id {
		c.JSON(http.StatusForbidden, gin.H{"Hata": "sadece kendi sifrenizi degistirebilirsiniz"})
		return
	}
	var istek ChangePasswordRequest
	if err := c.ShouldBindJSON(&istek); err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	if err := h.userService.UserChangePassword(c.Request.Context(), id, istek.CurrentPassword, istek.NewPassword); err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"mesaj": "sifre basariyla guncellendi"})
}
