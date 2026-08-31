package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/database"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/logger"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/middleware"

	_ "github.com/Furkangthb/stajKutuphaneUygulamasi/docs"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/core/services"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/handlers"
	"github.com/Furkangthb/stajKutuphaneUygulamasi/internal/repository"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Staj Kütüphane Uygulaması
// @version 1.0
// @description Kütüphane Otomasyonu Rest API Dokümantasyonu
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	log := logger.New()
	slog.SetDefault(log)

	database.VeriTabaniBaglanma()
	database.RedisConnect()

	userRepo := repository.NewUserRepository(database.DB)
	userService := services.NewUserServices(userRepo)
	userHandlers := handlers.NewUserHandlers(userService)

	bookRepo := repository.NewBookRepository(database.DB, database.RedisClient)
	bookService := services.NewBookServices(bookRepo)
	bookHandlers := handlers.NewBookHandlers(bookService)

	if err := bookRepo.WarmupCache(context.Background()); err != nil {
		slog.Warn("redise yuklenemedi")
	}

	reservationRepo := repository.NewResservationRepository(database.DB, database.RedisClient)
	reservationService := services.NewReservationServices(reservationRepo)
	reservationHandlers := handlers.NewReservationHandler(reservationService)

	authRepo := repository.NewAuthRepository(database.RedisClient)
	authService := services.NewAuthServices(authRepo, userRepo, os.Getenv("JWT_SECRET"))
	authHandler := handlers.NewAuthHandlers(authService)

	chatRepo := repository.NewChatRepository(database.DB)
	chatService := services.NewChatServices(os.Getenv("GEMINI_API_KEY"), bookService, chatRepo)
	chatHandlers := handlers.NewChatHandlers(chatService)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://furkan.local:54080", "http://localhost:5173", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/chat", middleware.RequireAuth(authService), chatHandlers.Chat)
	r.GET("/api/chat/history", middleware.RequireAuth(authService), chatHandlers.ChatHistory)

	r.POST("/api/users/register", userHandlers.UserRegister)
	r.POST("/api/auth/login", authHandler.Login)
	r.POST("/api/auth/logout", authHandler.Logout)
	r.GET("/api/users", middleware.RequireAuth(authService), middleware.RequireRole("admin"), userHandlers.UserList)
	r.DELETE("/api/users/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), userHandlers.UserDelete)
	r.PUT("/api/users/:id", middleware.RequireAuth(authService), userHandlers.UserUpdate)
	r.GET("/api/users/:id", middleware.RequireAuth(authService), userHandlers.UserGet)
	r.PUT("/api/users/:id/password", middleware.RequireAuth(authService), userHandlers.UserChangePassword)

	r.POST("/api/books", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookAdd)
	r.GET("/api/books/search", bookHandlers.BookSearch)
	r.GET("/api/books/:id", bookHandlers.BookGet)
	r.DELETE("/api/books/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookDelete)
	r.PUT("/api/books/:id", middleware.RequireAuth(authService), middleware.RequireRole("admin"), bookHandlers.BookUpdate)
	r.GET("/api/books", middleware.RequireAuth(authService), bookHandlers.BookList)

	r.POST("/api/reservation", middleware.RequireAuth(authService), reservationHandlers.ReservationCreate)
	r.PUT("/api/reservation/:id", middleware.RequireAuth(authService), reservationHandlers.ReservationUpdate)
	r.GET("/api/reservation/:id", middleware.RequireAuth(authService), reservationHandlers.ReservationGetUserByID)
	r.GET("/api/reservations", middleware.RequireAuth(authService), middleware.RequireRole("admin"), reservationHandlers.ReservationListAll)

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for {
			count, err := reservationService.ExpireOverdueReservations(context.Background())
			if err != nil {
				slog.Warn("Sureli rezervasyon kontrolu basarisiz", slog.Any("error", err))
			} else if count > 0 {
				slog.Warn(" rezervasyon suresi doldugu icin expired yapildi")
			}
			<-ticker.C
		}
	}()

	r.Run(":8080")
}
