package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

type ChatHandlers struct {
	chatServices *services.Chatservices
}

func NewChatHandlers(chatServices *services.Chatservices) *ChatHandlers {
	return &ChatHandlers{chatServices: chatServices}
}

type ChatRequist struct {
	Message string `json:"message" binding:"required"`
}

// Chat godoc
// @Summary Sohbet
// @Description Kullanıcı ile sohbet eder. Bearer token gerektirir.
// @Tags Chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChatRequest true "Sohbet İsteği"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/chat [post]

func (h *ChatHandlers) Chat(c *gin.Context) {
	var istek ChatRequist

	if err := c.ShouldBindJSON(&istek); err != nil {
		slog.Warn("chat islemi parse edilemedi", slog.Any("eroro", err))
		c.JSON(http.StatusBadRequest, gin.H{"Hata": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	cevap, err := h.chatServices.Sohbet(c.Request.Context(), userID, istek.Message)
	if err != nil {
		slog.Warn("Sohbet gerceklesemedi", slog.Int("userID", userID), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"Hata": err.Error()})
		return
	}
	slog.Info("Sohbet basariyla gerceklesti", slog.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"Cevap": cevap})
}

func (h *ChatHandlers) ChatHistory(c *gin.Context) {
	userID := c.GetInt("user_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	gecmis, err := h.chatServices.GetHistory(c.Request.Context(), userID, limit)
	if err != nil {
		slog.Warn("gecmis yuklenemedi", slog.Int("userID", userID), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"Hata": err.Error()})
		return
	}
	slog.Info("gecmis basarıyla yuklendi", slog.Int("userID", userID))
	c.JSON(http.StatusOK, gin.H{"data": gecmis})
}
