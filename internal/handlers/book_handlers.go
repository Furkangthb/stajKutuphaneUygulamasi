package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/gin-gonic/gin"
)

type BookHandlers struct {
	bookService *services.BookServices
}

func NewBookHandlers(bookService *services.BookServices) *BookHandlers {
	return &BookHandlers{bookService: bookService}
}

type BookAdd struct {
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	PublishDate time.Time `json:"publish_date" binding:"required" `
	Description string    `json:"description" binding:"required"`
	StockCount  int       `json:"stock_count" binding:"required"`
}

func (h *BookHandlers) BookAdd(c *gin.Context) {
	var kitapEkle BookAdd
	if err := c.ShouldBindJSON(&kitapEkle); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	book, err := h.bookService.BookAdd(ctx, kitapEkle.Title, kitapEkle.Author, kitapEkle.Genre, kitapEkle.PublishDate, kitapEkle.Description, kitapEkle.StockCount)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(201, book)
}

func (h *BookHandlers) BookGet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	book, err := h.bookService.BookGet(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandlers) BookDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	err = h.bookService.BookDelete(c.Request.Context(), id)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}

type BookRework struct {
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	PublishDate time.Time `json:"publish_date" binding:"required"`
	Description string    `json:"description" binding:"required"`
	StockCount  int       `json:"stock_count" binding:"required"`
}

func (h *BookHandlers) BookUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	var rework BookRework
	if err := c.ShouldBindJSON(&rework); err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	book, err := h.bookService.BookUpdate(ctx, id, rework.Title, rework.Author, rework.Genre, rework.PublishDate, rework.Description, rework.StockCount)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, book)
}

func (h *BookHandlers) BookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	books, err := h.bookService.BookList(c.Request.Context(), page, page_size)
	if err != nil {
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, books)
}
