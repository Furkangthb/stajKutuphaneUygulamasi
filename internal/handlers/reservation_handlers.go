package handlers

import (
	"strconv"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

type ReservationHandlers struct {
	reserService *services.ReservationServices
}

func NewReservationHandler(reserService *services.ReservationServices) *ReservationHandlers {
	return &ReservationHandlers{reserService: reserService}
}

type ReservationCreat struct {
	UserID int    `json:"user_id" binding:"required"`
	BookID int    `json:"book_id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

func (h *ReservationHandlers) ReservationCreate(c *gin.Context) {
	var reserve ReservationCreat
	if err := c.ShouldBindJSON(&reserve); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	reservation, err := h.reserService.ReservationCreate(ctx, reserve.UserID, reserve.BookID, reserve.Status)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(201, reservation)
}

type ReservationUpdate struct {
	Status string `json:"status" binding:"required"`
}

func (h *ReservationHandlers) ReservationUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	var rework ReservationUpdate
	if err := c.ShouldBindJSON(&rework); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	reserve, err := h.reserService.ReservationUpdate(ctx, int(id), rework.Status)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(200, reserve)
}

func (h *ReservationHandlers) ReservationGetUserByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	reserve, err := h.reserService.ReservationGetByUserID(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(200, reserve)
}
