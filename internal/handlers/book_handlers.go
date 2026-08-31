package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/domain"
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
	ISBN        string    `json:"isbn" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	PublishDate time.Time `json:"publish_date" binding:"required" `
	Description string    `json:"description" binding:"required"`
}

// BookAdd godoc
// @Summary Kitap Ekleme
// @Description Admin tarafından kitap eklenir
// @Tags Books
// @Accept json
// @Produce json
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param book body BookAdd true "Kitap Bilgileri"
// @Security BearerAuth
// @Router /api/books  [post]
func (h *BookHandlers) BookAdd(c *gin.Context) {
	var kitapEkle BookAdd
	if err := c.ShouldBindJSON(&kitapEkle); err != nil {
		slog.Warn("kitap ekleme parse edilemedi", slog.String("kitapISBN", kitapEkle.ISBN))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	book, err := h.bookService.BookAdd(ctx, kitapEkle.ISBN, kitapEkle.Title, kitapEkle.Author, kitapEkle.Genre, kitapEkle.PublishDate, kitapEkle.Description)
	if err != nil {
		slog.Warn("Kitap ekleme basarısız", slog.String("kitapISBN", kitapEkle.ISBN))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	slog.Info("Kitap basarili bir sekilde eklendi", slog.String("kitapISBN", kitapEkle.ISBN))
	c.JSON(201, book)
}

// BookGet godoc
// @Summary Kitap Getirme
// @Description Seçilen Kitabı getirir
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param id path int true "Kitap ID"
// @Router /api/books/{id}  [get]
func (h *BookHandlers) BookGet(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Warn(" id donusumu yapilamadi", slog.String("id_param", c.Param("id")))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	book, err := h.bookService.BookGet(c.Request.Context(), id)
	if err != nil {
		slog.Warn("Kitap getirelemedi", slog.Int("kitapID", int(id)))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	slog.Info("kitap getirme basarili", slog.Int("kitapID", book.ID))
	c.JSON(http.StatusOK, book)
}

// BookDelete godoc
// @Summary Kitap Silme
// @Description Admin tarafından kitap silinir
// @Tags Books
// @Accept json
// @Produce json
// @Success 204 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param id path int true "Kitap ID"
// @Security BearerAuth
// @Router /api/books/{id}  [delete]
func (h *BookHandlers) BookDelete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Warn("id donusumu yapilamadi", slog.String("id_param", c.Param("id")))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	err = h.bookService.BookDelete(c.Request.Context(), id)
	if err != nil {
		slog.Warn("Kitap silinemedi", slog.Int("kitapID", int(id)))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	slog.Info("Kitap basarili bir sekilde silindi", slog.String("kitapID", c.Param("id")))
	c.Status(http.StatusNoContent)
}

type BookRework struct {
	ISBN        string    `json:"isbn" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Author      string    `json:"author" binding:"required"`
	Genre       string    `json:"genre" binding:"required"`
	PublishDate time.Time `json:"publish_date" binding:"required"`
	Description string    `json:"description" binding:"required"`
}

// BookUpdate godoc
// @Summary Kitap Güncelleme
// @Description Admin tarafından kitap güncellenir
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param id path int true "Kitap ID"
// @Param book body BookRework true "Kitap Bilgileri"
// @Security BearerAuth
// @Router /api/books/{id}  [put]
func (h *BookHandlers) BookUpdate(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		slog.Warn("id donusumu yapilamadi", slog.String("id_param", c.Param("id")))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	var rework BookRework
	if err := c.ShouldBindJSON(&rework); err != nil {
		slog.Warn("update isteği parse edilemedi")
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	ctx := c.Request.Context()
	book, err := h.bookService.BookUpdate(ctx, id, rework.ISBN, rework.Title, rework.Author, rework.Genre, rework.PublishDate, rework.Description)
	if err != nil {
		slog.Warn("Kitap guncellemedi", slog.String("kitapISBN", rework.ISBN))
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	slog.Info("Kitap basarılı bir sekilde guncellendi", slog.String("BookISBN", book.ISBN))
	c.JSON(http.StatusOK, book)
}

// BookList godoc
// @Summary Kitap Listeleme
// @Description Kitapları listeler
// @Tags Books
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Param page query int false "Sayfa"
// @Param page_size query int false "Sayfa Boyutu"
// @Router /api/books  [get]
// @Security BearerAuth
func (h *BookHandlers) BookList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	page_size, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	books, err := h.bookService.BookList(c.Request.Context(), page, page_size)
	if err != nil {
		slog.Warn("Kitaplar listelenemedi")
		c.JSON(400, gin.H{
			"Hata": err.Error(),
		})
		return
	}
	slog.Info("Kitaplar basarılı bir sekilde listelendi")
	c.JSON(http.StatusOK, books)
}

func (h *BookHandlers) BookSearch(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		slog.Info("bos arama sorgusu")
		c.JSON(200, []*domain.Book{})
		return
	}

	keywords := strings.Fields(q)

	books, err := h.bookService.BookSearch(c.Request.Context(), keywords, 50)
	if err != nil {
		slog.Warn("kitap aranamadi")
		c.JSON(400, gin.H{"Hata": err.Error()})
		return
	}
	slog.Info("Kitap basarili bir şekilde arandi")
	c.JSON(200, books)
}
