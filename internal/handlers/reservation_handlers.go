package handlers

import (
	"errors"
	"net/http"
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

type ReservationCreatRequest struct {
	UserID int `json:"user_id"`
	BookID int `json:"book_id" binding:"required"`
}

// ReservationCreate godoc
// @Summary Rezervasyon Oluşturma
// @Description Kullanıcı için rezervasyon oluşturur. Bearer token gerektirir.
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ReservationCreatRequest true "Rezervasyon İsteği"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/reservation [post]
func (h *ReservationHandlers) ReservationCreate(c *gin.Context) {
	var reserve ReservationCreatRequest
	if err := c.ShouldBindJSON(&reserve); err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}

	tokenUserID := c.GetInt("user_id")
	tokenRole := c.GetString("role")

	ctx := c.Request.Context()
	reservation, err := h.reserService.ReservationCreate(ctx, tokenUserID, tokenRole, reserve.UserID, reserve.BookID)
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	c.JSON(201, reservation)
}

type ReservationUpdateRequest struct {
	Status string `json:"status" binding:"required"`
}

// ReservationUpdate godoc
// @Summary Rezervasyon Güncelleme
// @Description Kullanıcı için rezervasyon günceller. Bearer token gerektirir.
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Rezervasyon ID"
// @Param request body ReservationUpdateRequest true "Rezervasyon Güncellemesi"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/reservation/{id} [put]
func (h *ReservationHandlers) ReservationUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}

	var rework ReservationUpdateRequest
	if err := c.ShouldBindJSON(&rework); err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}

	tokenUserID := c.GetInt("user_id")
	tokenRole := c.GetString("role")

	reserve, err := h.reserService.ReservationUpdate(c.Request.Context(), tokenUserID, tokenRole, int(id), rework.Status)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrReservationNotFound):
			c.JSON(http.StatusNotFound, gin.H{"Hata": err.Error()})
		case errors.Is(err, services.ErrForbidden):
			c.JSON(http.StatusForbidden, gin.H{"Hata": err.Error()})
		default:
			// ErrInvalidStatus, ErrInvalidTransition ve diğer beklenmeyen hatalar 400 doner
			c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
		}
		return
	}
	c.JSON(200, reserve)
}

// ReservationListAll godoc
// @Summary Tüm Rezervasyonları Listeleme
// @Description Tüm rezervasyonları listeler. Bearer token gerektirir.
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/reservations [get]
func (h *ReservationHandlers) ReservationListAll(c *gin.Context) {
	reservations, err := h.reserService.ReservationListAll(c.Request.Context())
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	c.JSON(200, reservations)
}

// ReservationGetUserByID godoc
// @Summary Kullanıcıya Ait Rezervasyonları Listeleme
// @Description Kullanıcıya ait tüm rezervasyonları listeler. Bearer token gerektirir.
// @Tags Reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Kullanıcı ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/reservation/{id} [get]
func (h *ReservationHandlers) ReservationGetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}

	tokenUserID := c.GetInt("user_id")
	tokenRole := c.GetString("role")

	if tokenRole != "admin" && int64(tokenUserID) != id {
		c.JSON(http.StatusForbidden, gin.H{"Hata": "sadece kendi rezervasyonlarinizi goruntuleyebilirsiniz"})
		return
	}

	reserve, err := h.reserService.ReservationGetByUserID(c.Request.Context(), int(id))
	if err != nil {
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	c.JSON(200, reserve)
}
