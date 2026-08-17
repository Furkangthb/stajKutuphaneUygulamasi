package handlers

import (
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
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Role      string `json:"role" binding:"required"`
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
	user, err := h.userService.UserRegister(ctx, istek.FirstName, istek.LastName, istek.Phone, istek.Email, istek.Password, istek.Role)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(202, user)
}
