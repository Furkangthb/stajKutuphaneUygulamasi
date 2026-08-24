package handlers

import (
	"log"
	"net/http"

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

func (h *ChatHandlers) Chat(c *gin.Context) {
	var istek ChatRequist

	if err := c.ShouldBindJSON(&istek); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Hata": err.Error(),
		})
		return
	}

	cevap, err := h.chatServices.Sohbet(c.Request.Context(), istek.Message)
	if err != nil {
		log.Printf("Gemini Sohbet Hatası: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Cevap": cevap,
	})
}
