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

func (h *UserHandlers) UserUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
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
	user, err := h.userService.UserUpdate(ctx, int(id), rework.FirstName, rework.LastName, rework.Email, rework.Role)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h UserHandlers) UserList(c *gin.Context) {
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
